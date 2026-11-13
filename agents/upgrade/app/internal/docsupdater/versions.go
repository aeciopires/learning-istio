// Package docsupdater keeps this repository's installation/configuration
// files (README.md, day1.md, day1-sidecar.md, day2.md,
// INSTALL_REQUIREMENTS.md, mise.toml, and
// agents/upgrade/configs/versions.env) in sync with the latest stable
// tool versions resolved by internal/toolversions.
//
// It is deliberately conservative: only version *tokens* it can match
// unambiguously are rewritten (env var assignments, README bullet
// lines, istio.io/vX.Y/ link prefixes, mise-managed VERSION= blocks).
// Narrative prose and the kind node image digest (which can't be safely rewritten
// without an independently verified sha256) are left untouched, and a
// MAJOR version bump (e.g. Helm 3 -> 4) is never applied automatically -
// see Plan's Notes/Warnings for what was skipped and why.
package docsupdater

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/toolversions"
)

// ResolvedVersions is the set of concrete version strings a Plan run
// will try to apply across the repo's docs.
type ResolvedVersions struct {
	IstioMinor    string // "1.30"
	IstioPatch    string // "1.30.3"
	IstioPatchNum string // "3" (the part after "${ISTIO_RELEASE}.")
	GatewayAPI    string // "v1.6.0"
	Kind          string // "0.32.0" (no "v" prefix)
	Helm          string // "3.21.2"
	Kubectl       string // "1.36.2"
	MetalLB       string // "0.16.1"
	ArgoCD        string // "3.6.1" (no "v" prefix; scripts/argocd.sh adds it)

	// KindNodeK8s is the Kubernetes version baked into the kind node
	// image currently pinned in INSTALL_REQUIREMENTS.md (e.g. "1.35.5").
	// It is extracted, not resolved from GitHub, and is empty if it
	// couldn't be found - see ExtractKindNodeK8sVersion.
	KindNodeK8s string

	// Grafana/Prometheus/Kiali/Jaeger are the latest upstream versions of
	// the Istio addons (see scripts/install-addons.sh) - informational
	// only. Nothing in this repo's docs pins these as a literal token to
	// rewrite (the addons are installed from Istio's own bundled sample
	// manifests, tracking ISTIO_RELEASE), so they only flow into
	// configs/versions.env and `versions latest`, never into a file-rewrite
	// rule. No major-bump guard applies to them for the same reason.
	Grafana    string
	Prometheus string
	Kiali      string
	Jaeger     string
}

// CurrentVersions is what's presently pinned in the repo's docs, used
// only to guard against silently applying a MAJOR version bump.
type CurrentVersions struct {
	IstioMinor string
	Kind       string
	Helm       string
	Kubectl    string
}

// trailingSemverRe extracts a "MAJOR.MINOR.PATCH" suffix from an
// arbitrary tag name, e.g. "v1.30.3" -> "1.30.3", and, importantly,
// "metallb-chart-0.16.1" -> "0.16.1" (some repos, MetalLB included,
// tag chart releases with a non-version prefix; GitHub's "latest
// release" can point at one of those instead of a bare "vX.Y.Z" tag).
var trailingSemverRe = regexp.MustCompile(`(\d+\.\d+\.\d+)$`)

func normalizeVersionTag(tag string) string {
	if m := trailingSemverRe.FindStringSubmatch(tag); m != nil {
		return m[1]
	}
	return strings.TrimPrefix(tag, "v")
}

func majorOf(version string) (int, bool) {
	parts := strings.SplitN(version, ".", 2)
	if len(parts) == 0 {
		return 0, false
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	return n, true
}

// BuildResolvedVersions turns raw toolversions.Result tags into the
// normalized strings docsupdater's rewrite rules expect. A tool is left
// unresolved (existing pin untouched, with a warning) when: it failed
// to resolve, or resolving it would silently apply a MAJOR version bump
// and allowMajorBump is false (major bumps routinely break the exact
// commands this repo documents, e.g. Helm 3 -> 4, so they need a human
// to review first).
func BuildResolvedVersions(results []toolversions.Result, installRequirementsContent string, current CurrentVersions, allowMajorBump bool) (ResolvedVersions, []string) {
	var rv ResolvedVersions
	var warnings []string

	byKey := make(map[string]toolversions.Result, len(results))
	for _, r := range results {
		byKey[r.Tool.Key] = r
	}

	// get resolves key's target version. On any reason to *not* apply the
	// resolved latest (fetch failure, or a blocked major bump), it falls
	// back to currentVersion (if known) rather than "" - callers such as
	// versions.env want the version actually in effect, not a blank
	// field, even when nothing changed.
	get := func(key, currentVersion string) string {
		r, ok := byKey[key]
		if !ok || r.Err != nil {
			warnings = append(warnings, "could not resolve latest "+key+"; leaving existing pin untouched")
			return currentVersion
		}
		latest := normalizeVersionTag(r.Version)
		if !allowMajorBump && currentVersion != "" {
			curMajor, curOK := majorOf(currentVersion)
			latestMajor, latestOK := majorOf(latest)
			if curOK && latestOK && latestMajor > curMajor {
				warnings = append(warnings, "latest "+key+" is "+latest+", a MAJOR version bump from "+currentVersion+
					"; leaving existing pin untouched (pass --allow-major-bump to apply it anyway)")
				return currentVersion
			}
		}
		return latest
	}

	if istio := get("istio", current.IstioMinor); istio != "" {
		rv.IstioPatch = istio
		parts := strings.Split(istio, ".")
		if len(parts) == 3 {
			rv.IstioMinor = parts[0] + "." + parts[1]
			rv.IstioPatchNum = parts[2]
		} else {
			rv.IstioMinor = istio
		}
	}
	if gw, ok := byKey["gateway-api"]; ok && gw.Err == nil {
		rv.GatewayAPI = gw.Version // keep the "v" prefix, matches existing doc style
	} else {
		warnings = append(warnings, "could not resolve latest gateway-api; leaving existing pin untouched")
	}
	rv.Kind = get("kind", current.Kind)
	rv.Helm = get("helm", current.Helm)
	rv.Kubectl = get("kubectl", current.Kubectl)
	rv.MetalLB = get("metallb", "")
	rv.ArgoCD = get("argocd", "")

	// Informational only (see the ResolvedVersions doc comment) - no
	// current-version guard, since there's no pin to compare against.
	rv.Grafana = get("grafana", "")
	rv.Prometheus = get("prometheus", "")
	rv.Kiali = get("kiali", "")
	rv.Jaeger = get("jaeger", "")

	if k8s, ok := ExtractKindNodeK8sVersion(installRequirementsContent); ok {
		rv.KindNodeK8s = k8s
	} else {
		warnings = append(warnings, "could not extract the kind node image's Kubernetes version from INSTALL_REQUIREMENTS.md")
	}

	return rv, warnings
}

var kindNodeImageRe = regexp.MustCompile(`kindest/node:v(\d+\.\d+\.\d+)@sha256:`)

// ExtractKindNodeK8sVersion reads the Kubernetes version out of the
// kindest/node image tag pinned in INSTALL_REQUIREMENTS.md. The image
// digest is intentionally never rewritten by docsupdater - bumping the
// node version safely requires the matching sha256 from
// https://github.com/kubernetes-sigs/kind/releases, which is a manual
// step (documented in agents/upgrade/README.md).
func ExtractKindNodeK8sVersion(installRequirementsContent string) (string, bool) {
	m := kindNodeImageRe.FindStringSubmatch(installRequirementsContent)
	if m == nil {
		return "", false
	}
	return m[1], true
}

var (
	readmeIstioLineRe = regexp.MustCompile(`(?m)^- Istio: ([\d.]+)$`)
	miseTomlKindLine  = regexp.MustCompile(`(?m)^kind\s*=\s*"([\d.]+)"$`)
	miseTomlHelmLine  = regexp.MustCompile(`(?m)^helm\s*=\s*"([\d.]+)"$`)
	miseTomlKubeLine  = regexp.MustCompile(`(?m)^kubectl\s*=\s*"([\d.]+)"$`)
)

// ExtractCurrentVersions reads the versions presently pinned in
// README.md and mise.toml, used only for the major-version-bump guard in
// BuildResolvedVersions.
func ExtractCurrentVersions(readmeContent, miseTomlContent string) CurrentVersions {
	var cv CurrentVersions
	if m := readmeIstioLineRe.FindStringSubmatch(readmeContent); m != nil {
		cv.IstioMinor = m[1]
	}
	if m := miseTomlKindLine.FindStringSubmatch(miseTomlContent); m != nil {
		cv.Kind = m[1]
	}
	if m := miseTomlHelmLine.FindStringSubmatch(miseTomlContent); m != nil {
		cv.Helm = m[1]
	}
	if m := miseTomlKubeLine.FindStringSubmatch(miseTomlContent); m != nil {
		cv.Kubectl = m[1]
	}
	return cv
}
