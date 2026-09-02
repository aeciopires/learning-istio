package docsupdater

import (
	"errors"
	"strings"
	"testing"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/toolversions"
)

func result(key, repo, version string, err error) toolversions.Result {
	return toolversions.Result{Tool: toolversions.Tool{Key: key, Repo: repo, Label: key}, Version: version, Err: err}
}

func TestNormalizeVersionTag(t *testing.T) {
	cases := map[string]string{
		"v1.30.3":              "1.30.3",
		"1.30.3":               "1.30.3",
		"metallb-chart-0.16.1": "0.16.1",
		"v4.2.3":               "4.2.3",
	}
	for in, want := range cases {
		if got := normalizeVersionTag(in); got != want {
			t.Errorf("normalizeVersionTag(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildResolvedVersions_SkipsMajorBumpByDefault(t *testing.T) {
	results := []toolversions.Result{
		result("istio", "istio/istio", "1.30.3", nil),
		result("helm", "helm/helm", "v4.2.3", nil), // major bump from 3.x
		result("kind", "kubernetes-sigs/kind", "v0.32.0", nil),
		result("kubectl", "kubernetes/kubernetes", "v1.36.3", nil),
		result("metallb", "metallb/metallb", "metallb-chart-0.16.1", nil),
		result("gateway-api", "kubernetes-sigs/gateway-api", "v1.6.1", nil),
	}
	current := CurrentVersions{IstioMinor: "1.30", Kind: "0.32.0", Helm: "3.21.2", Kubectl: "1.36.2"}

	rv, warnings := BuildResolvedVersions(results, "", current, false)

	if rv.Helm != "3.21.2" {
		t.Errorf("expected Helm major bump to be skipped and fall back to the current pin, got %q", rv.Helm)
	}
	if rv.Kubectl != "1.36.3" {
		t.Errorf("expected kubectl patch bump to apply, got %q", rv.Kubectl)
	}
	if rv.MetalLB != "0.16.1" {
		t.Errorf("expected MetalLB tag to be normalized to 0.16.1, got %q", rv.MetalLB)
	}

	foundWarning := false
	for _, w := range warnings {
		if strings.Contains(w, "MAJOR version bump") {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Errorf("expected a MAJOR version bump warning, got %v", warnings)
	}
}

func TestBuildResolvedVersions_AllowMajorBump(t *testing.T) {
	results := []toolversions.Result{
		result("helm", "helm/helm", "v4.2.3", nil),
	}
	current := CurrentVersions{Helm: "3.21.2"}

	rv, _ := BuildResolvedVersions(results, "", current, true)
	if rv.Helm != "4.2.3" {
		t.Errorf("expected major bump to apply when allowed, got %q", rv.Helm)
	}
}

func TestBuildResolvedVersions_ToolErrorLeavesFieldEmpty(t *testing.T) {
	results := []toolversions.Result{
		result("kind", "kubernetes-sigs/kind", "", errors.New("network error")),
	}
	rv, warnings := BuildResolvedVersions(results, "", CurrentVersions{}, false)
	if rv.Kind != "" {
		t.Errorf("expected empty Kind on error, got %q", rv.Kind)
	}
	if len(warnings) == 0 {
		t.Error("expected a warning for the failed tool")
	}
}

func TestBuildResolvedVersions_AddonsAreInformationalOnly(t *testing.T) {
	results := []toolversions.Result{
		result("grafana", "grafana/grafana", "v12.0.1", nil),
		result("prometheus", "prometheus/prometheus", "v3.1.0", nil),
		result("kiali", "kiali/kiali", "v2.1.0", nil),
		result("jaeger", "jaegertracing/jaeger", "v2.1.0", nil),
	}

	// No "current" pin exists for addons (nothing rewrites a version
	// token for them), so even a big jump must resolve normally - the
	// MAJOR-bump guard only applies when there's an existing pin to
	// compare against.
	rv, warnings := BuildResolvedVersions(results, "", CurrentVersions{}, false)

	if rv.Grafana != "12.0.1" || rv.Prometheus != "3.1.0" || rv.Kiali != "2.1.0" || rv.Jaeger != "2.1.0" {
		t.Errorf("unexpected addon versions: %+v", rv)
	}
	for _, w := range warnings {
		if strings.Contains(w, "MAJOR version bump") {
			t.Errorf("did not expect a MAJOR version bump warning for addons (no pin to compare against), got %v", warnings)
		}
	}
}

func TestExtractCurrentVersions(t *testing.T) {
	readme := "- Istio: 1.30\n- Helm: 3.21.2\n"
	miseToml := "[tools]\nkind = \"0.32.0\"\nhelm = \"3.21.2\"\nkubectl = \"1.36.2\"\n"

	cv := ExtractCurrentVersions(readme, miseToml)
	if cv.IstioMinor != "1.30" || cv.Kind != "0.32.0" || cv.Helm != "3.21.2" || cv.Kubectl != "1.36.2" {
		t.Errorf("unexpected CurrentVersions: %+v", cv)
	}
}
