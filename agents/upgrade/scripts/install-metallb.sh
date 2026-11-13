#!/usr/bin/env bash
# Installs and configures MetalLB, mirroring INSTALL_REQUIREMENTS.md#metallb.
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

require_bin docker helm kubectl jq sed

metallb_enable_strict_arp

log_step "Installing MetalLB $METALLB_VERSION"
# helm repo add is itself idempotent (a repo already registered with the
# same URL is a no-op, exit 0) - --force-update makes that explicit
# instead of relying on the default behavior.
helm repo add metallb https://metallb.github.io/metallb --force-update >/dev/null
helm repo update metallb >/dev/null

helm_release_note metallb metallb-system
helm upgrade --install metallb metallb/metallb \
  --create-namespace \
  --namespace metallb-system \
  --version "$METALLB_VERSION" \
  --set "controller.tolerations[0].key=node-role.kubernetes.io/master" \
  --set "controller.tolerations[0].effect=NoSchedule" \
  --set speaker.tolerateMaster=true \
  --wait --timeout=900s

metallb_configure_address_pool

log_info "MetalLB installed and configured."
