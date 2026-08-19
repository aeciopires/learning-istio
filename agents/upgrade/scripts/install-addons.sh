#!/usr/bin/env bash
# Installs the Istio observability addons (Grafana, Prometheus, Kiali,
# Jaeger), mirroring the "Instale os seguintes addons" section of
# day1.md/day1-sidecar.md. Works the same regardless of sidecar/ambient
# mode - the addon manifests aren't mode-specific.
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

require_bin kubectl

log_step "Installing Istio addons (Grafana, Prometheus, Kiali, Jaeger)"
# kubectl apply is idempotent by itself (create-or-update); the notes below
# just make the already-exists/first-install distinction visible in the
# logs before each apply reconciles it.
kube_resource_note "Grafana deployment" deployment grafana istio-system
kubectl apply -f "$ISTIO_ADDONS_URL/grafana.yaml"
kube_resource_note "Prometheus deployment" deployment prometheus istio-system
kubectl apply -f "$ISTIO_ADDONS_URL/prometheus.yaml"
kube_resource_note "Kiali deployment" deployment kiali istio-system
kubectl apply -f "$ISTIO_ADDONS_URL/kiali.yaml"
kube_resource_note "Jaeger deployment" deployment jaeger istio-system
kubectl apply -f "$ISTIO_ADDONS_URL/jaeger.yaml"

log_step "Waiting for addon rollouts"
kubectl rollout status deployment/grafana -n istio-system --timeout=300s
kubectl rollout status deployment/prometheus -n istio-system --timeout=300s
kubectl rollout status deployment/kiali -n istio-system --timeout=300s
kubectl rollout status deployment/jaeger -n istio-system --timeout=300s

log_info "Addons installed."
log_info "Port-forward to access them locally:"
log_info "  kubectl -n istio-system port-forward service/grafana 3000:3000     # http://localhost:3000 (admin/admin)"
log_info "  kubectl -n istio-system port-forward service/prometheus 9090:9090  # http://localhost:9090"
log_info "  kubectl -n istio-system port-forward service/kiali 20001:20001    # http://localhost:20001"
log_info "  kubectl -n istio-system port-forward service/tracing 8081:80      # http://localhost:8081"
