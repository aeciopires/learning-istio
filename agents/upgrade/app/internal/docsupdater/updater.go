package docsupdater

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/toolversions"
)

// Updater applies docsupdater's rewrite rules across a checkout of the
// learning-istio repository.
type Updater struct {
	// RepoRoot is the root of the learning-istio repo (parent of
	// agents/upgrade).
	RepoRoot string

	// AllowMajorBump, when true, lets Plan apply a MAJOR version bump
	// (e.g. Helm 3 -> 4) instead of leaving the existing pin untouched
	// with a warning. Defaults to false (safe).
	AllowMajorBump bool
}

// PlanResult is everything Plan computed: proposed file changes plus any
// non-fatal warnings (e.g. a tool whose latest version failed to
// resolve, so its existing pin was left untouched).
type PlanResult struct {
	Versions ResolvedVersions
	Files    []*FileUpdate
	Warnings []string
}

// Plan resolves the latest tool versions and computes every file's
// proposed new content, without writing anything. Callers should
// present PlanResult before calling Apply (dry-run by default).
func (u *Updater) Plan(results []toolversions.Result) (*PlanResult, error) {
	installReqPath := filepath.Join(u.RepoRoot, "INSTALL_REQUIREMENTS.md")
	installReqContent, err := os.ReadFile(installReqPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", installReqPath, err)
	}

	readmePath := filepath.Join(u.RepoRoot, "README.md")
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", readmePath, err)
	}

	toolVersionsPath := filepath.Join(u.RepoRoot, ".tool-versions")
	toolVersionsContent, err := os.ReadFile(toolVersionsPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", toolVersionsPath, err)
	}

	current := ExtractCurrentVersions(string(readme), string(toolVersionsContent))
	rv, warnings := BuildResolvedVersions(results, string(installReqContent), current, u.AllowMajorBump)

	plan := &PlanResult{Versions: rv, Warnings: warnings}

	plan.Files = append(plan.Files, UpdateREADME(string(readme), rv))

	for _, name := range []string{"day1.md", "day2.md"} {
		path := filepath.Join(u.RepoRoot, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		plan.Files = append(plan.Files, UpdateDayFile(name, string(content), rv))
	}

	plan.Files = append(plan.Files, UpdateInstallRequirements(string(installReqContent), rv))
	plan.Files = append(plan.Files, UpdateToolVersions(string(toolVersionsContent), rv))

	versionsEnvPath := filepath.Join(u.RepoRoot, "agents", "upgrade", "configs", "versions.env")
	existingEnv := ""
	if b, err := os.ReadFile(versionsEnvPath); err == nil {
		existingEnv = string(b)
	}
	rendered := RenderVersionsEnv(rv)
	envUpdate := &FileUpdate{Path: filepath.Join("agents", "upgrade", "configs", "versions.env"), Before: existingEnv, After: rendered}
	if existingEnv != rendered {
		envUpdate.Changed = true
		envUpdate.Notes = []string{"regenerated from resolved tool versions"}
	}
	plan.Files = append(plan.Files, envUpdate)

	return plan, nil
}

// Apply writes every changed file in plan to disk.
func (u *Updater) Apply(plan *PlanResult) error {
	for _, f := range plan.Files {
		if !f.Changed {
			continue
		}
		fullPath := f.Path
		if !filepath.IsAbs(fullPath) {
			fullPath = filepath.Join(u.RepoRoot, f.Path)
		}
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(f.After), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", fullPath, err)
		}
	}
	return nil
}
