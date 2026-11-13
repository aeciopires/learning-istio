#!/usr/bin/env bash
# GitOps alternative to bootstrap.sh: installs Argo CD on the kind cluster
# and lets IT deploy every component (MetalLB, Gateway API CRDs, Istio,
# Bookinfo, addons) from argocd/application.yaml - official charts combined
# with the values files of the LOCAL directory argocd/values/, which this
# script embeds in each Application when rendering it.
#
# Usage: argocd.sh <command> [ambient|sidecar]
#   install             install Argo CD itself (official install manifest)
#   render  <mode>      print the rendered Applications, apply nothing
#   deploy  <mode>      apply the Applications step by step and wait
#   status              list the Applications with sync/health status
#   password            print the initial admin password + how to open the UI
#   undeploy            delete the Applications (and what they deployed)
#   uninstall           remove Argo CD itself
#
# The existing Helm/kubectl flow (bootstrap.sh and friends) is untouched;
# use one flow or the other on a given cluster, not both.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"

: "${ARGOCD_VERSION:?ARGOCD_VERSION not set in $VERSIONS_ENV - run 'make docs-update-apply'}"

ARGOCD_NAMESPACE="argocd"
ARGOCD_INSTALL_URL="https://raw.githubusercontent.com/argoproj/argo-cd/v${ARGOCD_VERSION}/manifests/install.yaml"
APPLICATION_TEMPLATE="$AGENT_DIR/argocd/application.yaml"

# Local directory holding one values file per chart. Argo CD runs inside
# the cluster and cannot read this disk, so render_applications embeds each
# file's content in its Application (helm.valuesObject) instead of pointing
# Argo CD at a Git branch - what you have on disk is what gets deployed.
ARGOCD_VALUES_DIR="${ARGOCD_VALUES_DIR:-$AGENT_DIR/argocd/values}"
ARGOCD_WAIT_TIMEOUT="${ARGOCD_WAIT_TIMEOUT:-900s}"

# Highest "learning-istio/step" label used in application.yaml.
LAST_STEP=6
APP_SELECTOR="app.kubernetes.io/part-of=learning-istio"

usage() {
  sed -n '2,/^set -euo/p' "${BASH_SOURCE[0]}" | sed -e '$d' -e 's/^# \{0,1\}//'
}

# render_applications <mode>
# Prints application.yaml with every ${PLACEHOLDER} replaced and every
# "# values-from: <file>" line replaced by that file's content. Plain bash
# on a fixed list of names - no envsubst dependency (not installed by
# default on macOS), and nothing else in the file is touched.
render_applications() {
  local mode="$1" label_key label_value content name
  if [[ "$mode" == "ambient" ]]; then
    label_key="istio.io/dataplane-mode"
    label_value="ambient"
  else
    label_key="istio-injection"
    label_value="enabled"
  fi

  content="$(<"$APPLICATION_TEMPLATE")"
  local -a names=(ISTIO_VERSION METALLB_VERSION GATEWAY_API_VERSION MY_NAMESPACE
    ISTIO_MODE BOOKINFO_NAMESPACE_LABEL_KEY BOOKINFO_NAMESPACE_LABEL_VALUE)
  local -a values=("$ISTIO_VERSION" "$METALLB_VERSION" "$GATEWAY_API_VERSION" "$MY_NAMESPACE"
    "$mode" "$label_key" "$label_value")
  local i
  for i in "${!names[@]}"; do
    name="\${${names[$i]}}"
    content="${content//"$name"/${values[$i]}}"
  done

  # Only the header comment may still mention placeholders by name.
  if grep -v '^#' <<<"$content" | grep -q '\${'; then
    log_error "unresolved placeholder left in $APPLICATION_TEMPLATE:"
    grep -v '^#' <<<"$content" | grep -n '\${' >&2
    exit 1
  fi

  # Embed the local values files. The directive must be indented (the
  # header comment, which starts at column 0, only documents it).
  local line indent file path values_line
  local include_re='^([[:space:]]+)# values-from: ([^[:space:]]+)$'
  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ "$line" =~ $include_re ]]; then
      indent="${BASH_REMATCH[1]}"
      file="${BASH_REMATCH[2]}"
      path="$ARGOCD_VALUES_DIR/$file"
      if [[ ! -f "$path" ]]; then
        log_error "values file not found: $path (referenced by $APPLICATION_TEMPLATE)"
        exit 1
      fi
      printf '%s# from argocd/values/%s (local file)\n' "$indent" "$file"
      while IFS= read -r values_line || [[ -n "$values_line" ]]; do
        if [[ -z "$values_line" ]]; then
          printf '\n'
        else
          printf '%s%s\n' "$indent" "$values_line"
        fi
      done <"$path"
    else
      printf '%s\n' "$line"
    fi
  done <<<"$content"
}

# step_selector <step> <mode>
step_selector() {
  echo "learning-istio/step=$1,learning-istio/mode in (all,$2)"
}

# wait_for_step <step> <mode>
# Waits until every Application of the step finished syncing and is Healthy.
wait_for_step() {
  local selector apps
  selector="$(step_selector "$1" "$2")"
  apps="$(kubectl -n "$ARGOCD_NAMESPACE" get applications.argoproj.io -l "$selector" -o name)"
  if [[ -z "$apps" ]]; then
    return 0
  fi
  log_info "Waiting (up to $ARGOCD_WAIT_TIMEOUT) for: $(tr '\n' ' ' <<<"$apps")"
  # shellcheck disable=SC2086 # one resource name per word, on purpose
  if ! kubectl -n "$ARGOCD_NAMESPACE" wait $apps \
    --for=jsonpath='{.status.operationState.phase}'=Succeeded --timeout="$ARGOCD_WAIT_TIMEOUT" ||
    ! kubectl -n "$ARGOCD_NAMESPACE" wait $apps \
      --for=jsonpath='{.status.health.status}'=Healthy --timeout="$ARGOCD_WAIT_TIMEOUT"; then
    log_error "step $1 did not become Healthy. Current status:"
    cmd_status || true
    log_error "Inspect one with: kubectl -n $ARGOCD_NAMESPACE describe application <name>"
    exit 1
  fi
}

require_argocd() {
  if ! kubectl get crd applications.argoproj.io >/dev/null 2>&1; then
    log_error "Argo CD is not installed on this cluster. Run 'make argocd-install' first."
    exit 1
  fi
}

cmd_install() {
  require_bin kubectl
  log_step "Installing Argo CD v$ARGOCD_VERSION"
  if kubectl get namespace "$ARGOCD_NAMESPACE" >/dev/null 2>&1; then
    log_info "Namespace '$ARGOCD_NAMESPACE' already exists, skipping creation."
  else
    kubectl create namespace "$ARGOCD_NAMESPACE"
  fi
  kube_resource_note "Argo CD (deployment argocd-server)" deployment argocd-server "$ARGOCD_NAMESPACE"
  # --server-side: some Argo CD CRDs exceed the size limit of client-side
  # apply's annotation. --force-conflicts: needed on re-runs/upgrades.
  # Both come from Argo CD's own "Getting Started".
  kubectl apply -n "$ARGOCD_NAMESPACE" --server-side --force-conflicts -f "$ARGOCD_INSTALL_URL"

  log_step "Waiting for Argo CD to become available"
  kubectl -n "$ARGOCD_NAMESPACE" wait --for=condition=available --timeout=600s deployment --all
  kubectl -n "$ARGOCD_NAMESPACE" rollout status statefulset/argocd-application-controller --timeout=600s

  log_info "Argo CD v$ARGOCD_VERSION installed in namespace '$ARGOCD_NAMESPACE'."
  log_info "Next: make argocd-deploy-ambient   (or: make argocd-deploy-sidecar)"
  log_info "UI:   make argocd-password"
}

cmd_render() {
  render_applications "$1"
}

cmd_deploy() {
  local mode="$1" rendered step
  require_bin docker kubectl jq sed
  require_argocd

  log_step "Deploying all components through Argo CD (mode=$mode)"
  log_info "Charts: official upstream repositories. Values: local directory $ARGOCD_VALUES_DIR"

  rendered="$(mktemp "${TMPDIR:-/tmp}/learning-istio-argocd.XXXXXX")"
  # shellcheck disable=SC2064 # expand $rendered now, not at exit
  trap "rm -f '$rendered'" EXIT
  render_applications "$mode" >"$rendered"

  # Switching modes on a live cluster: drop the Application that only the
  # other mode has, instead of leaving it behind.
  if [[ "$mode" == "sidecar" ]]; then
    kubectl -n "$ARGOCD_NAMESPACE" delete applications.argoproj.io \
      -l "$APP_SELECTOR,learning-istio/mode=ambient" --ignore-not-found --wait=true
  fi

  metallb_enable_strict_arp

  for step in $(seq 1 "$LAST_STEP"); do
    log_step "Argo CD step $step/$LAST_STEP"
    kubectl apply -f "$rendered" -l "$(step_selector "$step" "$mode")"
    wait_for_step "$step" "$mode"
    if [[ "$step" == "1" ]]; then
      # LoadBalancer Services (the ingress gateway, Bookinfo's Gateway)
      # only turn Healthy once MetalLB has addresses to hand out.
      metallb_configure_address_pool
    fi
  done

  cmd_status
  log_step "Argo CD deploy complete"
  log_info "Mode:      $mode"
  log_info "Namespace: $MY_NAMESPACE"
  log_info "Validate with:   make validate MODE=$mode"
  log_info "Gateway address: kubectl -n $MY_NAMESPACE get gateway -o jsonpath='{.items[*].status.addresses[*].value}'"
  log_info "Remove with:     make argocd-undeploy"
}

cmd_status() {
  require_bin kubectl
  require_argocd
  kubectl -n "$ARGOCD_NAMESPACE" get applications.argoproj.io -l "$APP_SELECTOR" \
    -o custom-columns='NAME:.metadata.name,STEP:.metadata.labels.learning-istio/step,SYNC:.status.sync.status,HEALTH:.status.health.status,REVISION:.status.sync.revision'
}

cmd_password() {
  require_bin kubectl
  require_argocd
  local password
  if ! password="$(kubectl -n "$ARGOCD_NAMESPACE" get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' 2>/dev/null)" || [[ -z "$password" ]]; then
    log_error "Secret argocd-initial-admin-secret not found (it is deleted once the password is changed)."
    exit 1
  fi
  log_info "User:     admin"
  log_info "Password: $(base64 -d <<<"$password")"
  log_info "UI:       kubectl -n $ARGOCD_NAMESPACE port-forward svc/argocd-server 8443:443"
  log_info "          then open https://localhost:8443 (self-signed certificate)"
}

cmd_undeploy() {
  require_bin kubectl
  require_argocd
  log_step "Removing every component deployed through Argo CD"
  local step selector
  # Reverse order. Each Application carries Argo CD's resources finalizer,
  # so deleting it also deletes what it deployed (cascade delete).
  for step in $(seq "$LAST_STEP" -1 1); do
    selector="$APP_SELECTOR,learning-istio/step=$step"
    if [[ "$step" == "1" ]]; then
      # Applied by kubectl (not by an Application) in cmd_deploy; it has to
      # go before MetalLB's CRDs do.
      kubectl -n metallb-system delete ipaddresspools.metallb.io first-pool --ignore-not-found 2>/dev/null || true
      kubectl -n metallb-system delete l2advertisements.metallb.io example --ignore-not-found 2>/dev/null || true
    fi
    log_info "Deleting Applications of step $step"
    kubectl -n "$ARGOCD_NAMESPACE" delete applications.argoproj.io -l "$selector" \
      --ignore-not-found --wait=true --timeout="$ARGOCD_WAIT_TIMEOUT"
  done
  log_info "Applications removed. Namespaces created by Argo CD (CreateNamespace=true) are kept:"
  log_info "  kubectl delete namespace $MY_NAMESPACE istio-system metallb-system --ignore-not-found"
  log_info "Argo CD itself is still installed - remove it with: make argocd-uninstall"
}

cmd_uninstall() {
  require_bin kubectl
  if kubectl get crd applications.argoproj.io >/dev/null 2>&1 &&
    [[ -n "$(kubectl -n "$ARGOCD_NAMESPACE" get applications.argoproj.io -o name 2>/dev/null)" ]]; then
    log_error "There are still Applications in '$ARGOCD_NAMESPACE'. Run 'make argocd-undeploy' first:"
    log_error "without Argo CD running, their finalizers can never complete."
    exit 1
  fi
  log_step "Uninstalling Argo CD v$ARGOCD_VERSION"
  kubectl delete -n "$ARGOCD_NAMESPACE" --ignore-not-found -f "$ARGOCD_INSTALL_URL"
  kubectl delete namespace "$ARGOCD_NAMESPACE" --ignore-not-found
  log_info "Argo CD removed."
}

COMMAND="${1:-}"
case "$COMMAND" in
  install | status | password | undeploy | uninstall)
    "cmd_$COMMAND"
    ;;
  render | deploy)
    MODE="${2:-$ISTIO_MODE}"
    if [[ "$MODE" != "ambient" && "$MODE" != "sidecar" ]]; then
      log_error "usage: $(basename "$0") $COMMAND [ambient|sidecar]"
      exit 1
    fi
    "cmd_$COMMAND" "$MODE"
    ;;
  -h | --help | help)
    usage
    ;;
  *)
    usage >&2
    exit 1
    ;;
esac
