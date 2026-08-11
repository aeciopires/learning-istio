// Package cmd implements the istio-upgrade CLI's commands (cobra).
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// cliVersion/cliGitCommit are set by Execute from the embedded VERSION
// file and (optionally) an ldflags-injected git commit.
var cliVersion = "0.0.0-dev"
var cliGitCommit = "unknown"

var rootCmd = &cobra.Command{
	Use:   "istio-upgrade",
	Short: "Upgrade-planning assistant for Istio and Envoy, and repo doc/version sync",
	Long: `istio-upgrade analyzes Istio and Envoy release notes across a version range
(minor-to-minor, e.g. 1.10 -> 1.30, or patch/hotfix, e.g. 1.29.0 -> 1.29.3) and
reports breaking changes, deprecations, installation/config changes, metrics
and component changes, and security fixes - to help plan an upgrade.

It also keeps this repository's installation/configuration files (README.md,
day1.md, day2.md, INSTALL_REQUIREMENTS.md, .tool-versions, and
agents/upgrade/configs/versions.env) pinned to the latest stable tool
versions from their official sources.

It does NOT create Kubernetes clusters or install Istio itself - that is
the job of the Makefile and scripts/ in this same directory (see
agents/upgrade/README.md).`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the CLI. version is embedded from the VERSION file by
// main.go; gitCommit is optionally injected at build time via ldflags.
func Execute(version, gitCommit string) {
	cliVersion = version
	cliGitCommit = gitCommit
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{.Name}} version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
