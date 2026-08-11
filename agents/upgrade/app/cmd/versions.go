package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/releasenotes"
	"github.com/aeciopires/learning-istio/agents/upgrade/internal/toolversions"
)

var versionsOpts struct {
	format   string
	cacheDir string
	refresh  bool
}

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Look up the latest stable version of each tool this repo pins",
}

var versionsLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Print the latest stable version of Istio, Envoy, kind, Helm, kubectl, MetalLB and Gateway API",
	Long: `latest queries the GitHub Releases API (falling back to tags for repos that
don't publish GitHub Releases) for the latest STABLE version of every tool this
repo's docs pin. Pre-releases (e.g. "1.31.0-alpha.0") are always excluded.`,
	RunE: runVersionsLatest,
}

func init() {
	versionsLatestCmd.Flags().StringVar(&versionsOpts.format, "format", "table", `output format: "table" or "json"`)
	versionsLatestCmd.Flags().StringVar(&versionsOpts.cacheDir, "cache-dir", "", "on-disk cache directory (default: OS user cache dir)")
	versionsLatestCmd.Flags().BoolVar(&versionsOpts.refresh, "refresh-cache", false, "bypass the on-disk cache and re-fetch everything")

	versionsCmd.AddCommand(versionsLatestCmd)
	rootCmd.AddCommand(versionsCmd)
}

func runVersionsLatest(cmd *cobra.Command, args []string) error {
	if versionsOpts.format != "table" && versionsOpts.format != "json" {
		return fmt.Errorf(`--format must be "table" or "json", got %q`, versionsOpts.format)
	}

	fetcher := releasenotes.NewFetcher(versionsOpts.cacheDir, versionsOpts.refresh)
	client := toolversions.NewClient(fetcher)
	results := client.LatestAll(cmd.Context())

	if versionsOpts.format == "json" {
		type entry struct {
			Key     string `json:"key"`
			Label   string `json:"label"`
			Repo    string `json:"repo"`
			Version string `json:"version,omitempty"`
			Error   string `json:"error,omitempty"`
		}
		out := make([]entry, 0, len(results))
		for _, r := range results {
			e := entry{Key: r.Tool.Key, Label: r.Tool.Label, Repo: r.Tool.Repo, Version: r.Version}
			if r.Err != nil {
				e.Error = r.Err.Error()
			}
			out = append(out, e)
		}
		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return err
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "TOOL\tLATEST\tSOURCE"); err != nil {
		return err
	}
	for _, r := range results {
		version := r.Version
		if r.Err != nil {
			version = "ERROR: " + r.Err.Error()
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\thttps://github.com/%s\n", r.Tool.Label, version, r.Tool.Repo); err != nil {
			return err
		}
	}
	return w.Flush()
}
