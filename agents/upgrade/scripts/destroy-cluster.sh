#!/usr/bin/env bash
# Deletes the local kind cluster. See ../../UNINSTALL_ISTIO.md for the
# manual equivalent.
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

require_bin kind

if ! kind get clusters 2>/dev/null | grep -qx "$CLUSTER_NAME"; then
  log_info "kind cluster '$CLUSTER_NAME' does not exist, nothing to do."
  exit 0
fi

log_step "Deleting kind cluster '$CLUSTER_NAME'"
kind delete cluster --name "$CLUSTER_NAME"
log_info "Cluster '$CLUSTER_NAME' deleted."
