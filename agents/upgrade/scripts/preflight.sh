#!/usr/bin/env bash
# Checks that every binary the cluster automation scripts shell out to is
# installed, and prints the resolved versions.env, before doing anything
# that creates real infrastructure.
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/common.sh"

log_step "Checking required binaries"
require_bin docker kind helm kubectl istioctl curl jq

log_info "docker:   $(docker version --format '{{.Client.Version}}' 2>/dev/null || echo 'unknown')"
log_info "kind:     $(kind version 2>/dev/null || echo 'unknown')"
log_info "helm:     $(helm version --short 2>/dev/null || echo 'unknown')"
log_info "kubectl:  $(kubectl version --client 2>/dev/null | tr '\n' ' ')"
# istioctl isn't used by the install scripts (Helm does the install), but
# it's required for troubleshooting the mesh after it's up - see
# INSTALL_REQUIREMENTS.md#istioctl.
log_info "istioctl: $(istioctl version --remote=false --short 2>/dev/null || echo 'unknown')"

log_step "Checking Docker daemon is reachable"
if ! docker info >/dev/null 2>&1; then
  log_error "Docker daemon not reachable. Is Docker running and is your user in the 'docker' group?"
  exit 1
fi
log_info "Docker daemon OK."

# istio-cni (ambient mode, and sidecar mode with pilot.cni.enabled=true)
# watches the node's CNI config directory via inotify; on a default-tuned
# Linux host the fs.inotify limits are too low and every istio-cni-node
# pod CrashLoopBackOffs with "couldn't initialize inotify: too many open
# files" (documented in ../../day1.md). Fix it proactively instead of
# letting `make bootstrap-*` hit it. fs.inotify is a Linux kernel feature
# with no macOS equivalent on the host - on macOS the relevant kernel is
# inside Docker Desktop's Linux VM, not reachable via the host's sysctl,
# so this check only applies on Linux.
OS_NAME="$(uname -s)"
if [[ "$OS_NAME" == "Linux" ]]; then
  log_step "Checking fs.inotify limits (needed for istio-cni)"
  MIN_INOTIFY_INSTANCES=8192
  CUR_INSTANCES="$(sysctl -n fs.inotify.max_user_instances 2>/dev/null || echo 0)"
  if [[ "$CUR_INSTANCES" -lt "$MIN_INOTIFY_INSTANCES" ]]; then
    log_warn "fs.inotify.max_user_instances=$CUR_INSTANCES is too low for istio-cni; raising it."
    if sudo -n sysctl -w fs.inotify.max_user_watches=2099999999 >/dev/null 2>&1 &&
      sudo -n sysctl -w fs.inotify.max_user_instances=2099999999 >/dev/null 2>&1 &&
      sudo -n sysctl -w fs.inotify.max_queued_events=2099999999 >/dev/null 2>&1; then
      log_info "Raised fs.inotify limits for this boot (non-persistent)."
    else
      log_warn "Could not raise fs.inotify limits automatically (no passwordless sudo)."
      log_warn "If installing Istio, run this first:"
      log_warn "  sudo sysctl -w fs.inotify.max_user_watches=2099999999"
      log_warn "  sudo sysctl -w fs.inotify.max_user_instances=2099999999"
      log_warn "  sudo sysctl -w fs.inotify.max_queued_events=2099999999"
    fi
  else
    log_info "fs.inotify.max_user_instances=$CUR_INSTANCES OK."
  fi
elif [[ "$OS_NAME" == "Darwin" ]]; then
  log_step "Checking fs.inotify limits (needed for istio-cni)"
  log_info "macOS detected - this host-level check doesn't apply (the relevant kernel is"
  log_info "Docker Desktop's Linux VM). If istio-cni-node CrashLoopBackOffs with"
  log_info "'too many open files', increase Docker Desktop's resources (Settings > Resources)"
  log_info "and, if needed, raise the limits inside the VM: docker run --rm --privileged --pid=host"
  log_info "justincormack/nsenter1 sysctl -w fs.inotify.max_user_instances=2099999999"
fi

log_step "Configured versions (from configs/versions.env)"
log_info "ISTIO_VERSION=$ISTIO_VERSION ISTIO_MINOR=$ISTIO_MINOR GATEWAY_API_VERSION=$GATEWAY_API_VERSION METALLB_VERSION=$METALLB_VERSION"
log_info "KIND_NODE_IMAGE=$KIND_NODE_IMAGE"
log_info "CLUSTER_NAME=$CLUSTER_NAME MY_NAMESPACE=$MY_NAMESPACE ISTIO_MODE=$ISTIO_MODE"

log_info "Preflight checks passed."
