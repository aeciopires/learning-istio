package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/releasenotes"
)

var analyzeOpts struct {
	component string
	from      string
	to        string
	envoyFrom string
	envoyTo   string
	format    string
	output    string
	cacheDir  string
	refresh   bool
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze Istio/Envoy release notes across a version range",
	Long: `analyze fetches and parses real Istio and Envoy release notes across a
version range and reports breaking changes, deprecations, installation/config
changes, metrics/telemetry changes, component changes, and (for patch/hotfix
ranges) bug-fix and security notes - one section per intermediate version.

--from/--to accept either a minor version ("1.10") or a full patch version
("1.29.3"); a patch-scoped range (e.g. "1.29.0" -> "1.29.3") additionally
reports per-patch bug-fix/CVE notes for that series.

Istio's exact bundled Envoy version can't be reliably derived from public
metadata (see internal/releasenotes/istio.go), so with --component=all and no
--envoy-from/--envoy-to, the report explains this instead of guessing; pass
--envoy-from/--envoy-to explicitly (or --component=envoy) to analyze Envoy too.

Examples:
  istio-upgrade analyze --from 1.10 --to 1.30
  istio-upgrade analyze --from 1.29.0 --to 1.29.3
  istio-upgrade analyze --component envoy --envoy-from 1.28 --envoy-to 1.32
  istio-upgrade analyze --from 1.10 --to 1.30 --format json --output report.json`,
	RunE: runAnalyze,
}

func init() {
	analyzeCmd.Flags().StringVar(&analyzeOpts.component, "component", "all", `what to analyze: "istio", "envoy", or "all"`)
	analyzeCmd.Flags().StringVar(&analyzeOpts.from, "from", "", `Istio range start, e.g. "1.10" or "1.29.0"`)
	analyzeCmd.Flags().StringVar(&analyzeOpts.to, "to", "", `Istio range end, e.g. "1.30" or "1.29.3"`)
	analyzeCmd.Flags().StringVar(&analyzeOpts.envoyFrom, "envoy-from", "", "Envoy range start (required if --component=envoy)")
	analyzeCmd.Flags().StringVar(&analyzeOpts.envoyTo, "envoy-to", "", "Envoy range end (required if --component=envoy)")
	analyzeCmd.Flags().StringVar(&analyzeOpts.format, "format", "markdown", `output format: "markdown" or "json"`)
	analyzeCmd.Flags().StringVar(&analyzeOpts.output, "output", "", "write the report to this file instead of stdout")
	analyzeCmd.Flags().StringVar(&analyzeOpts.cacheDir, "cache-dir", "", "on-disk cache directory (default: OS user cache dir)")
	analyzeCmd.Flags().BoolVar(&analyzeOpts.refresh, "refresh-cache", false, "bypass the on-disk cache and re-fetch everything")

	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	if analyzeOpts.component != "istio" && analyzeOpts.component != "envoy" && analyzeOpts.component != "all" {
		return fmt.Errorf(`--component must be "istio", "envoy", or "all", got %q`, analyzeOpts.component)
	}
	if analyzeOpts.format != "markdown" && analyzeOpts.format != "json" {
		return fmt.Errorf(`--format must be "markdown" or "json", got %q`, analyzeOpts.format)
	}

	fetcher := releasenotes.NewFetcher(analyzeOpts.cacheDir, analyzeOpts.refresh)
	analyzer := releasenotes.NewAnalyzer(fetcher)

	report, err := analyzer.Run(cmd.Context(), releasenotes.Options{
		Component: analyzeOpts.component,
		From:      analyzeOpts.from,
		To:        analyzeOpts.to,
		EnvoyFrom: analyzeOpts.envoyFrom,
		EnvoyTo:   analyzeOpts.envoyTo,
	})
	if err != nil {
		return err
	}

	var rendered string
	switch analyzeOpts.format {
	case "json":
		rendered, err = releasenotes.RenderJSON(report)
		if err != nil {
			return err
		}
	default:
		rendered = releasenotes.RenderMarkdown(report)
	}

	if analyzeOpts.output == "" {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), rendered)
		return err
	}
	if err := os.WriteFile(analyzeOpts.output, []byte(rendered), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", analyzeOpts.output, err)
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Report written to %s\n", analyzeOpts.output)
	return err
}
