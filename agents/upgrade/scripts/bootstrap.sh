#!/usr/bin/env bash
# Runs the full pipeline: preflight -> create cluster -> MetalLB -> Istio
# -> Bookinfo -> addons -> validate. Usage: bootstrap.sh [ambient|sidecar]
#
# INSTALL_BOOKINFO and INSTALL_ADDONS (default "true") let you install
# everything in one shot (the default) or skip a step while still running
# the rest, e.g.:
#   INSTALL_ADDONS=false make bootstrap-sidecar
# Each step also has its own standalone make target (create-cluster,
# install-metallb, install-istio-{sidecar,ambient},
# deploy-bookinfo-{sidecar,ambient}, install-addons) for running it alone.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"

MODE="${1:-$ISTIO_MODE}"
require_mode_arg "$MODE"

INSTALL_BOOKINFO="${INSTALL_BOOKINFO:-true}"
INSTALL_ADDONS="${INSTALL_ADDONS:-true}"

log_step "Bootstrapping a full Istio ($MODE) environment"
log_info "INSTALL_BOOKINFO=$INSTALL_BOOKINFO INSTALL_ADDONS=$INSTALL_ADDONS"

"$SCRIPT_DIR/preflight.sh"
"$SCRIPT_DIR/create-cluster.sh"
"$SCRIPT_DIR/install-metallb.sh"
"$SCRIPT_DIR/install-istio.sh" "$MODE"

if [[ "$INSTALL_BOOKINFO" == "true" ]]; then
  "$SCRIPT_DIR/deploy-bookinfo.sh" "$MODE"
else
  log_info "Skipping Bookinfo (INSTALL_BOOKINFO=$INSTALL_BOOKINFO)."
fi

if [[ "$INSTALL_ADDONS" == "true" ]]; then
  "$SCRIPT_DIR/install-addons.sh"
else
  log_info "Skipping addons (INSTALL_ADDONS=$INSTALL_ADDONS)."
fi

"$SCRIPT_DIR/validate-cluster.sh" "$MODE"

log_step "Bootstrap complete"
log_info "Cluster:   $CLUSTER_NAME (context: kind-$CLUSTER_NAME)"
log_info "Mode:      $MODE"
log_info "Namespace: $MY_NAMESPACE"
log_info ""
if [[ "$INSTALL_BOOKINFO" == "true" ]]; then
  log_info "Get the gateway address:"
  log_info "  kubectl -n $MY_NAMESPACE get gateway -o jsonpath='{.items[*].status.addresses[*].value}'"
  log_info "Or port-forward:"
  log_info "  kubectl -n $MY_NAMESPACE port-forward svc/bookinfo-gateway-istio 8080:80"
  log_info ""
fi
if [[ "$INSTALL_ADDONS" == "true" ]]; then
  log_info "Addons: run 'make install-addons' output above for the port-forward commands."
  log_info ""
fi
log_info "Tear it down with: make destroy-cluster"
