// Package toolversions discovers the latest stable version of every
// tool this repository's installation docs pin (Istio, Envoy, kind,
// Helm, kubectl, MetalLB, Gateway API), via the GitHub Releases API,
// falling back to the tags API for repos that don't tag GitHub
// Releases. It reuses releasenotes.Fetcher for HTTP + on-disk caching
// rather than duplicating that logic.
package toolversions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/releasenotes"
	"github.com/aeciopires/learning-istio/agents/upgrade/internal/semverx"
)

// Tool describes one piece of software this repo's docs pin a version for.
type Tool struct {
	// Key is the short, stable identifier used in configs/versions.env
	// and CLI output (e.g. "istio", "kind").
	Key string
	// Repo is the "owner/name" GitHub repository whose releases/tags
	// determine the latest stable version.
	Repo string
	// Label is a human-friendly name for report/table output.
	Label string
}

// Catalog is every tool `versions latest` / `docs update` knows about.
//
// Grafana/Prometheus/Kiali/Jaeger are installed as Istio addons (see
// agents/upgrade/scripts/install-addons.sh) via static manifests fetched
// straight from the istio/istio repo for the pinned Istio release - there
// is no independent version number to select there, so these 4 entries
// are informational only: they let `versions latest`/`configs/versions.env`
// report the current upstream release of each addon (useful to notice
// when the version bundled by Istio's sample manifests has fallen behind
// upstream), without controlling what actually gets installed.
var Catalog = []Tool{
	{Key: "istio", Repo: "istio/istio", Label: "Istio"},
	{Key: "envoy", Repo: "envoyproxy/envoy", Label: "Envoy"},
	{Key: "kind", Repo: "kubernetes-sigs/kind", Label: "kind"},
	{Key: "helm", Repo: "helm/helm", Label: "Helm"},
	{Key: "kubectl", Repo: "kubernetes/kubernetes", Label: "kubectl / Kubernetes"},
	{Key: "metallb", Repo: "metallb/metallb", Label: "MetalLB"},
	{Key: "gateway-api", Repo: "kubernetes-sigs/gateway-api", Label: "Gateway API"},
	{Key: "argocd", Repo: "argoproj/argo-cd", Label: "Argo CD"},
	{Key: "grafana", Repo: "grafana/grafana", Label: "Grafana (addon, informational)"},
	{Key: "prometheus", Repo: "prometheus/prometheus", Label: "Prometheus (addon, informational)"},
	{Key: "kiali", Repo: "kiali/kiali", Label: "Kiali (addon, informational)"},
	{Key: "jaeger", Repo: "jaegertracing/jaeger", Label: "Jaeger (addon, informational)"},
}

// Result is the outcome of resolving one Tool's latest stable version.
type Result struct {
	Tool    Tool
	Version string // raw tag name, e.g. "v1.36.2" or "1.30.3"
	Err     error
}

// Client resolves latest-stable-version lookups against the GitHub API.
type Client struct {
	Fetcher *releasenotes.Fetcher
	APIBase string
}

// NewClient wires a Client to the real GitHub API, sharing f's cache.
func NewClient(f *releasenotes.Fetcher) *Client {
	return &Client{Fetcher: f, APIBase: "https://api.github.com"}
}

type ghRelease struct {
	TagName string `json:"tag_name"`
}

type ghTagEntry struct {
	Name string `json:"name"`
}

// LatestRelease returns the newest stable tag for repo: the GitHub
// "latest release" (which GitHub itself excludes drafts/prereleases
// from) if the repo publishes Releases, else the highest semver-looking
// tag from the tags API.
func (c *Client) LatestRelease(ctx context.Context, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", c.APIBase, repo)
	body, err := c.Fetcher.Get(ctx, url)
	if err == nil {
		var rel ghRelease
		if jsonErr := json.Unmarshal(body, &rel); jsonErr == nil && rel.TagName != "" {
			return rel.TagName, nil
		}
	} else if _, ok := err.(*releasenotes.NotFoundError); !ok {
		return "", err
	}
	return c.latestStableTag(ctx, repo)
}

func (c *Client) latestStableTag(ctx context.Context, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/tags?per_page=100", c.APIBase, repo)
	body, err := c.Fetcher.Get(ctx, url)
	if err != nil {
		return "", err
	}
	var tags []ghTagEntry
	if err := json.Unmarshal(body, &tags); err != nil {
		return "", fmt.Errorf("parsing tags from %s: %w", url, err)
	}

	var best string
	var bestVersion semverx.Version
	for _, t := range tags {
		clean := strings.TrimPrefix(t.Name, "v")
		if strings.ContainsAny(clean, "-+") {
			continue // skip pre-releases/build metadata (e.g. "1.30.0-rc.1")
		}
		v, err := semverx.Parse(clean)
		if err != nil {
			continue
		}
		if best == "" || semverx.Compare(v, bestVersion) > 0 {
			best = t.Name
			bestVersion = v
		}
	}
	if best == "" {
		return "", fmt.Errorf("no stable tags found for %s", repo)
	}
	return best, nil
}

// LatestAll resolves every tool in Catalog, best-effort: a failure for
// one tool doesn't stop the others (it's recorded in Result.Err).
func (c *Client) LatestAll(ctx context.Context) []Result {
	results := make([]Result, 0, len(Catalog))
	for _, tool := range Catalog {
		version, err := c.LatestRelease(ctx, tool.Repo)
		results = append(results, Result{Tool: tool, Version: version, Err: err})
	}
	return results
}
