#!/usr/bin/env bash
# Deploys the Bookinfo sample app + gateway, mirroring day1.md. Labels the
# target namespace for the requested data plane mode BEFORE creating any
# pods, so they're captured by the mesh immediately (no manual restart
# needed, unlike the interactive walkthrough in day1.md).
# Usage: deploy-bookinfo.sh [ambient|sidecar]
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

MODE="${1:-$ISTIO_MODE}"
require_mode_arg "$MODE"
require_bin kubectl

log_step "Preparing namespace '$MY_NAMESPACE' (mode=$MODE)"
kubectl get namespace "$MY_NAMESPACE" >/dev/null 2>&1 || kubectl create namespace "$MY_NAMESPACE"

if [[ "$MODE" == "ambient" ]]; then
  kubectl label namespace "$MY_NAMESPACE" istio.io/dataplane-mode=ambient --overwrite
  # Removing a label that isn't set is harmless but noisy ("label ... not
  # found", printed on stdout, not stderr) - silence it on a clean namespace.
  kubectl label namespace "$MY_NAMESPACE" istio-injection- --overwrite >/dev/null 2>&1 || true
else
  kubectl label namespace "$MY_NAMESPACE" istio-injection=enabled --overwrite
  kubectl label namespace "$MY_NAMESPACE" istio.io/dataplane-mode- --overwrite >/dev/null 2>&1 || true
fi

log_step "Deploying Bookinfo sample app"
kubectl -n "$MY_NAMESPACE" apply -f "$ISTIO_BOOKINFO_URL/platform/kube/bookinfo.yaml"
kubectl -n "$MY_NAMESPACE" apply -f "$ISTIO_BOOKINFO_URL/platform/kube/bookinfo-versions.yaml"

log_step "Creating the ingress Gateway + HTTPRoute"
kubectl -n "$MY_NAMESPACE" apply -f "$ISTIO_BOOKINFO_URL/gateway-api/bookinfo-gateway.yaml"

log_step "Waiting for Bookinfo deployments to become available"
kubectl -n "$MY_NAMESPACE" wait --for=condition=available --timeout=300s deployment --all

log_info "Bookinfo deployed to namespace '$MY_NAMESPACE' (mode=$MODE)."
log_info "Get the gateway address with:"
log_info "  kubectl -n $MY_NAMESPACE get gateway -o jsonpath='{.items[*].status.addresses[*].value}'"
