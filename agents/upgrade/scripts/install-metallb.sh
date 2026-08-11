#!/usr/bin/env bash
# Installs and configures MetalLB, mirroring INSTALL_REQUIREMENTS.md#metallb.
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

require_bin docker helm kubectl jq sed

log_step "Enabling strict ARP mode on kube-proxy"
kubectl get configmap kube-proxy -n kube-system -o yaml |
  sed -e "s/strictARP: false/strictARP: true/" |
  kubectl apply -f - -n kube-system

log_step "Installing MetalLB $METALLB_VERSION"
helm repo add metallb https://metallb.github.io/metallb >/dev/null
helm repo update metallb >/dev/null

helm upgrade --install metallb metallb/metallb \
  --create-namespace \
  --namespace metallb-system \
  --version "$METALLB_VERSION" \
  --set "controller.tolerations[0].key=node-role.kubernetes.io/master" \
  --set "controller.tolerations[0].effect=NoSchedule" \
  --set speaker.tolerateMaster=true \
  --wait --timeout=900s

log_step "Configuring the MetalLB address pool"
KIND_ADDRESS="$(docker network inspect kind | jq -r '.[].IPAM.Config[] | select(.Subnet | test("^[0-9]+\\.")) | .Subnet' | cut -d"." -f1-3 | head -1)"
if [[ -z "$KIND_ADDRESS" ]]; then
  log_error "Could not determine the kind docker network's IPv4 subnet."
  exit 1
fi
KIND_LB_RANGE="${KIND_ADDRESS}.${METALLB_ADDRESS_END}/${METALLB_NETMASK}"
log_info "Using load balancer address range: $KIND_LB_RANGE"

cat <<EOF | kubectl apply --namespace metallb-system -f -
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: first-pool
spec:
  addresses:
  - $KIND_LB_RANGE
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: example
spec:
  ipAddressPools:
  - first-pool
EOF

log_info "MetalLB installed and configured."
