package docsupdater

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/toolversions"
)

// writeFixtureRepo materializes a minimal repo checkout (the files Plan
// reads) in a temp directory, reusing the same fixtures files_test.go
// exercises directly, so Plan/Apply are tested against realistic content.
func writeFixtureRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	files := map[string]string{
		"README.md":               readmeFixture,
		"day1.md":                 dayFileFixture,
		"day1-sidecar.md":         dayFileFixture,
		"day2.md":                 dayFileFixture,
		"INSTALL_REQUIREMENTS.md": installRequirementsFixture,
		"mise.toml":               miseTomlFixture,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			t.Fatalf("writing fixture %s: %v", name, err)
		}
	}
	return root
}

func fixtureToolversionsResults() []toolversions.Result {
	return []toolversions.Result{
		{Tool: toolversions.Tool{Key: "istio", Repo: "istio/istio"}, Version: "1.31.1"},
		{Tool: toolversions.Tool{Key: "gateway-api", Repo: "kubernetes-sigs/gateway-api"}, Version: "v1.7.0"},
		{Tool: toolversions.Tool{Key: "kind", Repo: "kubernetes-sigs/kind"}, Version: "v0.33.0"},
		{Tool: toolversions.Tool{Key: "helm", Repo: "helm/helm"}, Version: "v3.22.0"},
		{Tool: toolversions.Tool{Key: "kubectl", Repo: "kubernetes/kubernetes"}, Version: "v1.37.0"},
		{Tool: toolversions.Tool{Key: "metallb", Repo: "metallb/metallb"}, Version: "v0.17.0"},
	}
}

func TestUpdater_Plan_IncludesDay1Sidecar(t *testing.T) {
	root := writeFixtureRepo(t)
	u := &Updater{RepoRoot: root}

	plan, err := u.Plan(fixtureToolversionsResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var day1Sidecar *FileUpdate
	for _, f := range plan.Files {
		if f.Path == "day1-sidecar.md" {
			day1Sidecar = f
		}
	}
	if day1Sidecar == nil {
		t.Fatal("expected day1-sidecar.md to be part of the plan")
	}
	if !day1Sidecar.Changed {
		t.Errorf("expected day1-sidecar.md to be marked as changed, notes: %v", day1Sidecar.Notes)
	}
}

func TestUpdater_PlanAndApply_WritesDay1Sidecar(t *testing.T) {
	root := writeFixtureRepo(t)
	u := &Updater{RepoRoot: root}

	plan, err := u.Plan(fixtureToolversionsResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := u.Apply(plan); err != nil {
		t.Fatalf("unexpected error applying plan: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "day1-sidecar.md"))
	if err != nil {
		t.Fatalf("unexpected error reading day1-sidecar.md: %v", err)
	}
	if !strings.Contains(string(got), "export ISTIO_RELEASE=1.31") {
		t.Errorf("expected day1-sidecar.md to be rewritten with the new Istio release, got:\n%s", got)
	}

	gotMiseToml, err := os.ReadFile(filepath.Join(root, "mise.toml"))
	if err != nil {
		t.Fatalf("unexpected error reading mise.toml: %v", err)
	}
	if !strings.Contains(string(gotMiseToml), `kind = "0.33.0"`) {
		t.Errorf("expected mise.toml to be rewritten with the new kind pin, got:\n%s", gotMiseToml)
	}
}
