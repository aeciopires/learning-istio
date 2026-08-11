#!/usr/bin/env bash
# Runs the full pipeline: preflight -> create cluster -> MetalLB -> Istio
# -> Bookinfo -> validate. Usage: bootstrap.sh [ambient|sidecar]
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"

MODE="${1:-$ISTIO_MODE}"
require_mode_arg "$MODE"

log_step "Bootstrapping a full Istio ($MODE) environment"

"$SCRIPT_DIR/preflight.sh"
"$SCRIPT_DIR/create-cluster.sh"
"$SCRIPT_DIR/install-metallb.sh"
"$SCRIPT_DIR/install-istio.sh" "$MODE"
"$SCRIPT_DIR/deploy-bookinfo.sh" "$MODE"
"$SCRIPT_DIR/validate-cluster.sh" "$MODE"

log_step "Bootstrap complete"
log_info "Cluster:   $CLUSTER_NAME (context: kind-$CLUSTER_NAME)"
log_info "Mode:      $MODE"
log_info "Namespace: $MY_NAMESPACE"
log_info ""
log_info "Get the gateway address:"
log_info "  kubectl -n $MY_NAMESPACE get gateway -o jsonpath='{.items[*].status.addresses[*].value}'"
log_info "Or port-forward:"
log_info "  kubectl -n $MY_NAMESPACE port-forward svc/bookinfo-gateway-istio 8080:80"
log_info ""
log_info "Tear it down with: make destroy-cluster"
