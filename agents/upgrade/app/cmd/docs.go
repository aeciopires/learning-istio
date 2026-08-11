package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/docsupdater"
	"github.com/aeciopires/learning-istio/agents/upgrade/internal/releasenotes"
	"github.com/aeciopires/learning-istio/agents/upgrade/internal/toolversions"
)

var docsOpts struct {
	repoRoot       string
	apply          bool
	cacheDir       string
	refresh        bool
	allowMajorBump bool
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Update this repo's installation/configuration files",
}

var docsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Sync README.md, day1.md, day1-sidecar.md, day2.md, INSTALL_REQUIREMENTS.md, .tool-versions and configs/versions.env to the latest stable tool versions",
	Long: `update resolves the latest stable version of every tool this repo pins
(Istio, Gateway API, kind, Helm, kubectl, MetalLB - see "istio-upgrade versions
latest") and rewrites the matching version tokens in:

  README.md, day1.md, day1-sidecar.md, day2.md, INSTALL_REQUIREMENTS.md,
  .tool-versions, and agents/upgrade/configs/versions.env (consumed by scripts/*.sh)

It never touches the kind node image's tag+digest in INSTALL_REQUIREMENTS.md
(that needs a manually verified sha256 - see agents/upgrade/README.md) or any
narrative prose.

A MAJOR version bump for kind/Helm/kubectl/Istio (e.g. Helm 3 -> 4) is never
applied automatically, since it routinely breaks the exact commands this repo
documents - it's reported as a warning instead. Pass --allow-major-bump to
apply it anyway.

By default this is a DRY RUN: it prints the proposed changes and writes
nothing. Pass --apply to actually write the files.`,
	RunE: runDocsUpdate,
}

func init() {
	docsUpdateCmd.Flags().StringVar(&docsOpts.repoRoot, "repo-root", "", "learning-istio repo root (default: auto-detected from the current directory)")
	docsUpdateCmd.Flags().BoolVar(&docsOpts.apply, "apply", false, "write the changes (default: dry-run diff only)")
	docsUpdateCmd.Flags().StringVar(&docsOpts.cacheDir, "cache-dir", "", "on-disk cache directory (default: OS user cache dir)")
	docsUpdateCmd.Flags().BoolVar(&docsOpts.refresh, "refresh-cache", false, "bypass the on-disk cache and re-fetch everything")
	docsUpdateCmd.Flags().BoolVar(&docsOpts.allowMajorBump, "allow-major-bump", false, "also apply MAJOR version bumps (e.g. Helm 3 -> 4) instead of warning and skipping them")

	docsCmd.AddCommand(docsUpdateCmd)
	rootCmd.AddCommand(docsCmd)
}

func runDocsUpdate(cmd *cobra.Command, args []string) error {
	repoRoot := docsOpts.repoRoot
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		found, err := findRepoRoot(cwd)
		if err != nil {
			return fmt.Errorf("%w (pass --repo-root explicitly)", err)
		}
		repoRoot = found
	}

	fetcher := releasenotes.NewFetcher(docsOpts.cacheDir, docsOpts.refresh)
	tvClient := toolversions.NewClient(fetcher)
	results := tvClient.LatestAll(cmd.Context())

	updater := &docsupdater.Updater{RepoRoot: repoRoot, AllowMajorBump: docsOpts.allowMajorBump}
	plan, err := updater.Plan(results)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(out, "Resolved versions: istio=%s gateway-api=%s kind=%s helm=%s kubectl=%s metallb=%s (kind node k8s=%s)\n\n",
		plan.Versions.IstioPatch, plan.Versions.GatewayAPI, plan.Versions.Kind, plan.Versions.Helm,
		plan.Versions.Kubectl, plan.Versions.MetalLB, plan.Versions.KindNodeK8s)

	anyChanged := false
	for _, f := range plan.Files {
		if !f.Changed {
			_, _ = fmt.Fprintf(out, "= %s (already up to date)\n", f.Path)
			continue
		}
		anyChanged = true
		_, _ = fmt.Fprintf(out, "* %s\n", f.Path)
		for _, n := range f.Notes {
			_, _ = fmt.Fprintf(out, "    - %s\n", n)
		}
	}
	for _, w := range plan.Warnings {
		_, _ = fmt.Fprintf(out, "! warning: %s\n", w)
	}

	if !anyChanged {
		_, _ = fmt.Fprintln(out, "\nNothing to update.")
		return nil
	}

	if !docsOpts.apply {
		_, _ = fmt.Fprintln(out, "\nDry run only - re-run with --apply to write these changes.")
		return nil
	}

	if err := updater.Apply(plan); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, "\nChanges written.")
	return nil
}

// findRepoRoot walks up from start looking for the learning-istio repo
// root: a directory containing both a ".git" entry and a "README.md".
func findRepoRoot(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "README.md")); err == nil {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not auto-detect the learning-istio repo root from %s", start)
		}
		dir = parent
	}
}
