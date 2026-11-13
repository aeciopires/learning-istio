#!/usr/bin/env bash
# Shared helpers + config loading for agents/upgrade/scripts/*.sh.
# This file is meant to be SOURCED by every other script in this
# directory, not executed directly - it owns cluster/Istio automation,
# deliberately kept separate from the Go CLI (which only analyzes
# release notes and updates version pins - see ../README.md).

set -euo pipefail

SCRIPT_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
AGENT_DIR="$(cd "$SCRIPT_LIB_DIR/../.." && pwd)"
# shellcheck disable=SC2034 # available to every script that sources this file
REPO_ROOT="$(cd "$AGENT_DIR/.." && pwd)"

VERSIONS_ENV="$AGENT_DIR/configs/versions.env"
if [[ -f "$VERSIONS_ENV" ]]; then
  # shellcheck disable=SC1090
  source "$VERSIONS_ENV"
else
  echo "ERROR: $VERSIONS_ENV not found." >&2
  echo "Run 'make docs-update-apply' (or 'go run . docs update --apply') from $AGENT_DIR first." >&2
  exit 1
fi

: "${ISTIO_VERSION:?ISTIO_VERSION not set in $VERSIONS_ENV}"
: "${ISTIO_MINOR:?ISTIO_MINOR not set in $VERSIONS_ENV}"
: "${GATEWAY_API_VERSION:?GATEWAY_API_VERSION not set in $VERSIONS_ENV}"
: "${METALLB_VERSION:?METALLB_VERSION not set in $VERSIONS_ENV}"

# Overridable via environment; sensible defaults matching day1.md/day2.md
# and INSTALL_REQUIREMENTS.md.
CLUSTER_NAME="${CLUSTER_NAME:-istio-upgrade}"
MY_NAMESPACE="${MY_NAMESPACE:-myapp}"
ISTIO_MODE="${ISTIO_MODE:-ambient}" # ambient | sidecar

# Must match the kindest/node image pinned in ../../INSTALL_REQUIREMENTS.md
# ("#kind" section). docsupdater deliberately never rewrites this (the
# tag and sha256 digest must be bumped together, verified against
# https://github.com/kubernetes-sigs/kind/releases) - so keep these two
# in sync by hand when you do bump it.
KIND_NODE_IMAGE="${KIND_NODE_IMAGE:-kindest/node:v${KIND_NODE_K8S_VERSION:-1.37.0}@sha256:a1ed56cfb0e7b93589bdf97c8cd566405a265939e3620fc4f5de89adff580ae5}"
POD_SUBNET="${POD_SUBNET:-10.244.0.0/16}"
SERVICE_SUBNET="${SERVICE_SUBNET:-10.96.0.0/12}"
METALLB_ADDRESS_END="${METALLB_ADDRESS_END:-100}"
METALLB_NETMASK="${METALLB_NETMASK:-27}"

ISTIO_BASE_URL="https://raw.githubusercontent.com/istio/istio/release-${ISTIO_MINOR}/samples"
ISTIO_BOOKINFO_URL="${ISTIO_BASE_URL}/bookinfo"
# Used by install-addons.sh (Grafana/Prometheus/Kiali/Jaeger) - the same
# URLs day1.md/day1-sidecar.md use manually.
ISTIO_ADDONS_URL="${ISTIO_BASE_URL}/addons"

log_info()  { echo -e "\033[1;34m[INFO]\033[0m  $*"; }
log_warn()  { echo -e "\033[1;33m[WARN]\033[0m  $*" >&2; }
log_error() { echo -e "\033[1;31m[ERROR]\033[0m $*" >&2; }
log_step()  { echo -e "\n\033[1;36m==> $*\033[0m"; }

require_bin() {
  local missing=0
  for bin in "$@"; do
    if ! command -v "$bin" >/dev/null 2>&1; then
      log_error "required binary not found in PATH: $bin"
      missing=1
    fi
  done
  [[ $missing -eq 0 ]] || exit 1
}

require_mode_arg() {
  local mode="${1:-}"
  if [[ "$mode" != "ambient" && "$mode" != "sidecar" ]]; then
    log_error "usage: $(basename "$0") [ambient|sidecar]"
    exit 1
  fi
}

# --- Existence checks shared by the install-*.sh / deploy-*.sh scripts -----
# Every install step in this repo follows the same rule: check first, never
# fail just because the thing already exists, log an INFO line either way,
# then keep going (reconciling in place via `helm upgrade --install` /
# `kubectl apply`, which are themselves idempotent). These helpers only
# control the log message - the caller is responsible for still running the
# idempotent install command afterwards.

# helm_release_note <release> <namespace>
# Logs whether a Helm release already exists, without side effects.
helm_release_note() {
  local release="$1" ns="$2"
  if helm -n "$ns" status "$release" >/dev/null 2>&1; then
    log_info "Helm release '$release' already exists in namespace '$ns' - reconciling/upgrading it in place."
  else
    log_info "Helm release '$release' does not exist yet in namespace '$ns' - installing it."
  fi
}

# kube_resource_note <description> <kind> <name> [namespace]
# Logs whether a Kubernetes resource already exists, without side effects.
kube_resource_note() {
  local desc="$1" kind="$2" name="$3" ns="${4:-}"
  local ns_args=()
  [[ -n "$ns" ]] && ns_args=(-n "$ns")
  if kubectl get "$kind" "$name" "${ns_args[@]}" >/dev/null 2>&1; then
    log_info "$desc already exists - reconciling/updating it."
  else
    log_info "$desc does not exist yet - creating it."
  fi
}

# --- MetalLB pieces shared by install-metallb.sh (Helm) and argocd.sh -------
# Neither is part of the MetalLB chart: strict ARP is a kube-proxy setting,
# and the address pool depends on the subnet Docker gave THIS host's "kind"
# network, so it can't live in a static values file or Git manifest.

# metallb_enable_strict_arp
metallb_enable_strict_arp() {
  log_step "Enabling strict ARP mode on kube-proxy"
  kubectl get configmap kube-proxy -n kube-system -o yaml |
    sed -e "s/strictARP: false/strictARP: true/" |
    kubectl apply -f - -n kube-system
}

# metallb_configure_address_pool
# Requires the MetalLB CRDs and webhook to be ready.
metallb_configure_address_pool() {
  log_step "Configuring the MetalLB address pool"
  local kind_address kind_lb_range
  kind_address="$(docker network inspect kind | jq -r '.[].IPAM.Config[] | select(.Subnet | test("^[0-9]+\\.")) | .Subnet' | cut -d"." -f1-3 | head -1)"
  if [[ -z "$kind_address" ]]; then
    log_error "Could not determine the kind docker network's IPv4 subnet."
    exit 1
  fi
  kind_lb_range="${kind_address}.${METALLB_ADDRESS_END}/${METALLB_NETMASK}"
  log_info "Using load balancer address range: $kind_lb_range"

  kube_resource_note "IPAddressPool 'first-pool'" ipaddresspool first-pool metallb-system
  cat <<EOF | kubectl apply --namespace metallb-system -f -
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: first-pool
spec:
  addresses:
  - $kind_lb_range
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: example
spec:
  ipAddressPools:
  - first-pool
EOF
}
