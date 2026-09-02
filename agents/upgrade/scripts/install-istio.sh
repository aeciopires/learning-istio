#!/usr/bin/env bash
# Installs Istio in either "sidecar" or "ambient" mode, mirroring
# day1.md's manual Helm install steps. The CNI node agent (istio-cni) is
# installed in BOTH modes - for sidecar it replaces the privileged
# istio-init container for traffic redirection; for ambient it's what
# captures pods into the ztunnel data plane.
# Usage: install-istio.sh [ambient|sidecar]
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

MODE="${1:-$ISTIO_MODE}"
require_mode_arg "$MODE"
require_bin helm kubectl

log_step "Installing Istio $ISTIO_VERSION (profile=$MODE)"

# helm repo add is itself idempotent (a repo already registered with the
# same URL is a no-op, exit 0) - --force-update makes that explicit
# instead of relying on the default behavior.
helm repo add istio https://blob.istio.io/istio-release/charts --force-update >/dev/null
helm repo update istio >/dev/null

log_step "Installing base CRDs (istio-base)"
helm_release_note istio-base istio-system
helm -n istio-system upgrade --install istio-base istio/base \
  --version "$ISTIO_VERSION" --create-namespace --wait --timeout 900s

if ! kubectl get crd gateways.gateway.networking.k8s.io >/dev/null 2>&1; then
  log_step "Installing Kubernetes Gateway API CRDs $GATEWAY_API_VERSION"
  kubectl apply -f "https://github.com/kubernetes-sigs/gateway-api/releases/download/${GATEWAY_API_VERSION}/standard-install.yaml"
else
  log_info "Gateway API CRDs already installed, skipping."
fi

log_step "Installing istiod (control plane)"
# Istio's Helm charts only have a "profile=ambient" value for ambient mode;
# there is no "profile=sidecar" - sidecar injection is the chart's default
# profile. For sidecar mode we instead set pilot.cni.enabled=true, which
# tells the sidecar injector to skip the privileged istio-init init
# container and rely on the CNI node agent (installed below, for BOTH
# modes) for traffic redirection instead - see
# https://istio.io/latest/docs/setup/additional-setup/cni/
ISTIOD_ARGS=()
if [[ "$MODE" == "ambient" ]]; then
  ISTIOD_ARGS=(--set profile=ambient)
else
  ISTIOD_ARGS=(--set pilot.cni.enabled=true)
fi
helm_release_note istiod istio-system
helm -n istio-system upgrade --install istiod istio/istiod \
  --version "$ISTIO_VERSION" "${ISTIOD_ARGS[@]}" --wait --timeout 900s

log_step "Installing the CNI node agent"
# Used in both modes: traffic redirection for sidecar pods (since istiod
# was installed with pilot.cni.enabled=true above) and pod-to-ztunnel
# capture for ambient pods.
CNI_ARGS=()
if [[ "$MODE" == "ambient" ]]; then
  CNI_ARGS=(--set profile=ambient)
fi
helm_release_note istio-cni istio-system
helm -n istio-system upgrade --install istio-cni istio/cni \
  --version "$ISTIO_VERSION" "${CNI_ARGS[@]}" --wait --timeout 900s

if [[ "$MODE" == "ambient" ]]; then
  log_step "Installing ztunnel (ambient node proxy)"
  helm_release_note ztunnel istio-system
  helm -n istio-system upgrade --install ztunnel istio/ztunnel \
    --version "$ISTIO_VERSION" --wait --timeout 900s
fi

log_step "Installing the ingress gateway"
helm_release_note istio-ingress istio-system
helm -n istio-system upgrade --install istio-ingress istio/gateway \
  --version "$ISTIO_VERSION" --wait --timeout 900s

log_step "Istio ($MODE) install status"
helm -n istio-system ls
kubectl -n istio-system get pods -o wide

log_info "Istio ($MODE) installed."
