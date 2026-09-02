package docsupdater

import (
	"fmt"
	"regexp"
)

// FileUpdate is the result of applying every applicable rule to one file.
type FileUpdate struct {
	Path    string
	Before  string
	After   string
	Changed bool
	Notes   []string // human-readable summary of what changed, for dry-run/--apply output
}

func newFileUpdate(path, content string) *FileUpdate {
	return &FileUpdate{Path: path, Before: content, After: content}
}

func (u *FileUpdate) apply(newContent string, c *change) {
	if c == nil {
		return
	}
	u.After = newContent
	u.Changed = true
	u.Notes = append(u.Notes, fmt.Sprintf("%s: %q -> %q", c.Description, c.Before, c.After))
}

// UpdateREADME rewrites the "Configurações do Ambiente" version bullets
// and any istio.io/vX.Y/ documentation links in README.md.
func UpdateREADME(content string, rv ResolvedVersions) *FileUpdate {
	u := newFileUpdate("README.md", content)
	cur := content

	if rv.KindNodeK8s != "" {
		next, c := applyRule(cur, "Kubernetes (kind) version bullet",
			regexp.MustCompile(`(?m)^- Kubernetes \(kind\): [\d.]+$`),
			"- Kubernetes (kind): "+rv.KindNodeK8s)
		u.apply(next, c)
		cur = next
	}
	if rv.Kubectl != "" {
		next, c := applyRule(cur, "kubectl version bullet",
			regexp.MustCompile(`(?m)^- kubectl: [\d.]+$`),
			"- kubectl: "+rv.Kubectl)
		u.apply(next, c)
		cur = next
	}
	if rv.Helm != "" {
		next, c := applyRule(cur, "Helm version bullet",
			regexp.MustCompile(`(?m)^- Helm: [\d.]+$`),
			"- Helm: "+rv.Helm)
		u.apply(next, c)
		cur = next
	}
	if rv.IstioMinor != "" {
		next, c := applyRule(cur, "Istio version bullet",
			regexp.MustCompile(`(?m)^- Istio: [\d.]+$`),
			"- Istio: "+rv.IstioMinor)
		u.apply(next, c)
		cur = next

		next, c = rewriteIstioLinks(cur, rv.IstioMinor)
		u.apply(next, c)
		cur = next
	}
	if rv.IstioPatch != "" {
		// istioctl's version always mirrors the resolved Istio version - see
		// the note on UpdateInstallRequirements below.
		next, c := applyRule(cur, "istioctl version bullet",
			regexp.MustCompile(`(?m)^- istioctl: [\d.]+$`),
			"- istioctl: "+rv.IstioPatch)
		u.apply(next, c)
	}
	return u
}

// UpdateDayFile rewrites the ISTIO_RELEASE/VERSION_ISTIO/GATEWAY_API_VERSION
// env var exports and istio.io/vX.Y/ links shared by day1.md,
// day1-sidecar.md and day2.md.
func UpdateDayFile(path, content string, rv ResolvedVersions) *FileUpdate {
	u := newFileUpdate(path, content)
	cur := content

	if rv.IstioMinor != "" {
		next, c := applyRule(cur, "ISTIO_RELEASE env var",
			regexp.MustCompile(`export ISTIO_RELEASE=\d+\.\d+`),
			"export ISTIO_RELEASE="+rv.IstioMinor)
		u.apply(next, c)
		cur = next

		next, c = rewriteIstioLinks(cur, rv.IstioMinor)
		u.apply(next, c)
		cur = next
	}
	if rv.IstioPatchNum != "" {
		// NB: "$$" is regexp.ReplaceAllString's escape for a literal "$" -
		// without it, "${ISTIO_RELEASE}" in the replacement would be
		// parsed as a (nonexistent) named capture group and vanish.
		next, c := applyRule(cur, "VERSION_ISTIO patch suffix",
			regexp.MustCompile(`export VERSION_ISTIO="\$\{ISTIO_RELEASE\}\.\d+"`),
			`export VERSION_ISTIO="$${ISTIO_RELEASE}.`+rv.IstioPatchNum+`"`)
		u.apply(next, c)
		cur = next
	}
	if rv.GatewayAPI != "" {
		next, c := applyRule(cur, "GATEWAY_API_VERSION env var",
			regexp.MustCompile(`export GATEWAY_API_VERSION="v[\d.]+"`),
			`export GATEWAY_API_VERSION="`+rv.GatewayAPI+`"`)
		u.apply(next, c)
	}
	return u
}

var (
	helmVersionRe     = regexp.MustCompile(`(?s)VERSION="[^"]+"(\n\nmise latest helm)`)
	kindVersionRe     = regexp.MustCompile(`(?s)VERSION="[^"]+"(\nmise latest kind)`)
	kubectlVersionRe  = regexp.MustCompile(`VERSION_OPTION_1="[^"]+"`)
	istioctlVersionRe = regexp.MustCompile(`(?s)VERSION="[^"]+"(\n\nmise latest istioctl)`)
	metallbVersionRe  = regexp.MustCompile(`export HELM_CHART_VERSION="[^"]+"`)
)

// UpdateInstallRequirements rewrites the mise-managed VERSION= pins for
// Helm, kubectl, istioctl and kind, and the MetalLB Helm chart --version,
// in INSTALL_REQUIREMENTS.md. The kind node image (tag + sha256 digest)
// is intentionally left untouched - see ExtractKindNodeK8sVersion.
func UpdateInstallRequirements(content string, rv ResolvedVersions) *FileUpdate {
	u := newFileUpdate("INSTALL_REQUIREMENTS.md", content)
	cur := content

	if rv.Helm != "" {
		next, c := applyRule(cur, "Helm VERSION pin", helmVersionRe, `VERSION="`+rv.Helm+`"${1}`)
		u.apply(next, c)
		cur = next
	}
	if rv.Kind != "" {
		next, c := applyRule(cur, "kind VERSION pin", kindVersionRe, `VERSION="`+rv.Kind+`"${1}`)
		u.apply(next, c)
		cur = next
	}
	if rv.Kubectl != "" {
		next, c := applyRule(cur, "kubectl VERSION_OPTION_1 pin", kubectlVersionRe, `VERSION_OPTION_1="`+rv.Kubectl+`"`)
		u.apply(next, c)
		cur = next
	}
	// istioctl's version always mirrors the resolved Istio version - it's
	// released from the same istio/istio repo/tag, so there's no separate
	// GitHub lookup for it (see ResolvedVersions.IstioPatch).
	if rv.IstioPatch != "" {
		next, c := applyRule(cur, "istioctl VERSION pin", istioctlVersionRe, `VERSION="`+rv.IstioPatch+`"${1}`)
		u.apply(next, c)
		cur = next
	}
	if rv.MetalLB != "" {
		next, c := applyRule(cur, "MetalLB Helm chart --version", metallbVersionRe, `export HELM_CHART_VERSION="`+rv.MetalLB+`"`)
		u.apply(next, c)
	}
	return u
}

var (
	miseTomlKindRe     = regexp.MustCompile(`(?m)^kind\s*=\s*"[\d.]+"$`)
	miseTomlHelmRe     = regexp.MustCompile(`(?m)^helm\s*=\s*"[\d.]+"$`)
	miseTomlKubectlRe  = regexp.MustCompile(`(?m)^kubectl\s*=\s*"[\d.]+"$`)
	miseTomlIstioctlRe = regexp.MustCompile(`(?m)^istioctl\s*=\s*"[\d.]+"$`)
)

// UpdateMiseToml rewrites the "[tools]" table in mise.toml (mise's native
// config format): the single source of truth for the kind/helm/kubectl/
// istioctl pins this repo tracks.
func UpdateMiseToml(content string, rv ResolvedVersions) *FileUpdate {
	u := newFileUpdate("mise.toml", content)
	cur := content

	if rv.Kind != "" {
		next, c := applyRule(cur, "kind pin", miseTomlKindRe, `kind = "`+rv.Kind+`"`)
		u.apply(next, c)
		cur = next
	}
	if rv.Helm != "" {
		next, c := applyRule(cur, "helm pin", miseTomlHelmRe, `helm = "`+rv.Helm+`"`)
		u.apply(next, c)
		cur = next
	}
	if rv.Kubectl != "" {
		next, c := applyRule(cur, "kubectl pin", miseTomlKubectlRe, `kubectl = "`+rv.Kubectl+`"`)
		u.apply(next, c)
		cur = next
	}
	// istioctl's version always mirrors the resolved Istio version - see
	// the comment in UpdateInstallRequirements.
	if rv.IstioPatch != "" {
		next, c := applyRule(cur, "istioctl pin", miseTomlIstioctlRe, `istioctl = "`+rv.IstioPatch+`"`)
		u.apply(next, c)
	}
	return u
}

// RenderVersionsEnv renders agents/upgrade/configs/versions.env: the
// single source of truth the Makefile-driven cluster scripts source.
func RenderVersionsEnv(rv ResolvedVersions) string {
	return fmt.Sprintf(`# Generated by "istio-upgrade docs update --apply". Do not edit by hand -
# re-run the command instead. Sourced by agents/upgrade/scripts/*.sh.

ISTIO_MINOR=%s
ISTIO_VERSION=%s
ISTIOCTL_VERSION=%s
GATEWAY_API_VERSION=%s
KIND_VERSION=%s
HELM_VERSION=%s
KUBECTL_VERSION=%s
METALLB_VERSION=%s
KIND_NODE_K8S_VERSION=%s

# Latest upstream addon versions (informational only - see
# scripts/install-addons.sh; not what controls what actually gets
# installed, which tracks ISTIO_VERSION's bundled sample manifests instead).
GRAFANA_VERSION=%s
PROMETHEUS_VERSION=%s
KIALI_VERSION=%s
JAEGER_VERSION=%s
`, rv.IstioMinor, rv.IstioPatch, rv.IstioPatch, rv.GatewayAPI, rv.Kind, rv.Helm, rv.Kubectl, rv.MetalLB, rv.KindNodeK8s,
		rv.Grafana, rv.Prometheus, rv.Kiali, rv.Jaeger)
}
