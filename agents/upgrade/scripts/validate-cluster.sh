#!/usr/bin/env bash
# Validates that the cluster + Istio + Bookinfo are actually working: the
# "tests for cluster/Istio installation" piece, runnable standalone
# (make validate) or as the last step of make bootstrap-{sidecar,ambient}.
# Usage: validate-cluster.sh [ambient|sidecar]
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

MODE="${1:-$ISTIO_MODE}"
require_mode_arg "$MODE"
require_bin kind kubectl

FAILED=0

check() {
  local desc="$1"
  shift
  if "$@" >/dev/null 2>&1; then
    log_info "OK   - $desc"
  else
    log_error "FAIL - $desc"
    FAILED=1
  fi
}

log_step "Validating the kind cluster"
check "kind cluster '$CLUSTER_NAME' exists" bash -c "kind get clusters | grep -qx '$CLUSTER_NAME'"
check "kubectl can reach the API server" kubectl --context "kind-$CLUSTER_NAME" get nodes

log_step "Validating the Istio control plane (mode=$MODE)"
check "istiod deployment available" kubectl -n istio-system wait --for=condition=available --timeout=120s deployment/istiod
check "istio-ingress deployment available" kubectl -n istio-system wait --for=condition=available --timeout=120s deployment/istio-ingress

check "istio-cni-node daemonset ready" kubectl -n istio-system rollout status daemonset/istio-cni-node --timeout=120s
if [[ "$MODE" == "ambient" ]]; then
  check "ztunnel daemonset ready" kubectl -n istio-system rollout status daemonset/ztunnel --timeout=120s
fi

log_step "Validating the Bookinfo sample app in namespace '$MY_NAMESPACE'"
check "productpage deployment available" kubectl -n "$MY_NAMESPACE" wait --for=condition=available --timeout=120s deployment/productpage-v1
check "reviews deployment available" kubectl -n "$MY_NAMESPACE" wait --for=condition=available --timeout=120s deployment/reviews-v1
check "ratings deployment available" kubectl -n "$MY_NAMESPACE" wait --for=condition=available --timeout=120s deployment/ratings-v1
check "details deployment available" kubectl -n "$MY_NAMESPACE" wait --for=condition=available --timeout=120s deployment/details-v1

log_step "Validating traffic through the mesh"
RATINGS_POD="$(kubectl -n "$MY_NAMESPACE" get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)"
if [[ -n "$RATINGS_POD" ]]; then
  RATINGS_CONTAINER="ratings"
  if [[ "$MODE" == "sidecar" ]]; then
    check "productpage reachable from inside the mesh (via sidecar pod)" \
      kubectl -n "$MY_NAMESPACE" exec "$RATINGS_POD" -c "$RATINGS_CONTAINER" -- curl -sSf -o /dev/null http://productpage:9080/productpage
  else
    check "productpage reachable from inside the mesh (ambient)" \
      kubectl -n "$MY_NAMESPACE" exec "$RATINGS_POD" -c "$RATINGS_CONTAINER" -- curl -sSf -o /dev/null http://productpage:9080/productpage
  fi
else
  log_error "FAIL - could not find a ratings pod to test from"
  FAILED=1
fi

check "gateway programmed" bash -c \
  "kubectl -n '$MY_NAMESPACE' get gateway -o jsonpath='{.items[0].status.conditions[?(@.type==\"Programmed\")].status}' | grep -qx True"

echo
if [[ $FAILED -ne 0 ]]; then
  log_error "Validation FAILED - see the FAIL lines above."
  exit 1
fi
log_info "Validation PASSED - cluster, Istio ($MODE) and Bookinfo are healthy."
