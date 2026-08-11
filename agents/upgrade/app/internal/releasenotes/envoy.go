package releasenotes

import (
	"context"
	"fmt"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/semverx"
)

// EnvoyClient fetches and parses Envoy's per-release version history
// pages (RST source, served by envoyproxy.io's Sphinx docs build).
type EnvoyClient struct {
	Fetcher *Fetcher

	// RSTURLTemplate accepts three %s args, all the same "MAJOR.MINOR.PATCH"
	// version string repeated (docs tree version, minor dir, file name) -
	// e.g. https://www.envoyproxy.io/docs/envoy/v1.32.0/_sources/version_history/v1.32/v1.32.0.rst.txt
	RSTURLTemplate string

	Repo string
}

// NewEnvoyClient returns a client wired to the real envoyproxy.io / GitHub
// endpoints.
func NewEnvoyClient(f *Fetcher) *EnvoyClient {
	return &EnvoyClient{
		Fetcher:        f,
		RSTURLTemplate: "https://www.envoyproxy.io/docs/envoy/v%s/_sources/version_history/v%s/v%s.rst.txt",
		Repo:           "envoyproxy/envoy",
	}
}

// FetchMinor builds the VersionReport for one Major.Minor Envoy series.
// For a plain minor-range request it reports the ".0" release notes
// (what's new in that minor); for a patch-scoped request (explicit
// --envoy-from/--envoy-to with a patch, or a single-minor hotfix range)
// it additionally reports every patch strictly after FromPatch up to
// ToPatch (or the latest known patch).
func (c *EnvoyClient) FetchMinor(ctx context.Context, plan semverx.MinorPlan) VersionReport {
	minor := plan.MinorString()
	report := VersionReport{Component: ComponentEnvoy, Minor: minor}

	baseFrom := 0
	if plan.IncludePatchDetail {
		baseFrom = plan.FromPatch
	}
	c.fetchPatch(ctx, &report, plan.Major, plan.Minor, baseFrom)

	if !plan.IncludePatchDetail {
		return report
	}

	toPatch := plan.ToPatch
	if toPatch == -1 {
		tags, err := listTags(ctx, c.Fetcher, c.Repo)
		if err != nil {
			report.Warnings = append(report.Warnings, fmt.Sprintf("listing tags for %s: %v", c.Repo, err))
			return report
		}
		latest, ok := latestPatchForMinor(tags, plan.Major, plan.Minor)
		if !ok {
			report.Warnings = append(report.Warnings, fmt.Sprintf("could not determine latest patch for envoy %s.x", minor))
			return report
		}
		toPatch = latest
	}
	for p := baseFrom + 1; p <= toPatch; p++ {
		c.fetchPatch(ctx, &report, plan.Major, plan.Minor, p)
	}
	return report
}

func (c *EnvoyClient) fetchPatch(ctx context.Context, report *VersionReport, major, minor, patch int) {
	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	minorDir := fmt.Sprintf("%d.%d", major, minor)
	url := fmt.Sprintf(c.RSTURLTemplate, version, minorDir, version)

	body, err := c.Fetcher.Get(ctx, url)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			report.Warnings = append(report.Warnings, fmt.Sprintf("no version history page for envoy %s (skipped)", version))
			return
		}
		report.Warnings = append(report.Warnings, fmt.Sprintf("fetching %s: %v", url, err))
		return
	}

	for _, sec := range parseRSTSections(string(body)) {
		cat, breaking := envoyCategoryFor(sec.Title)
		for _, bullet := range sec.Bullets {
			report.Changes = append(report.Changes, Change{
				Component: ComponentEnvoy,
				Version:   version,
				Category:  cat,
				Summary:   bullet,
				Breaking:  breaking,
				Source:    url,
			})
		}
	}
}

func envoyCategoryFor(heading string) (Category, bool) {
	switch normalizeHeading(heading) {
	case "incompatible behavior changes":
		return CategoryBreaking, true
	case "removed config or runtime":
		return CategoryDeprecated, true
	case "deprecated":
		return CategoryDeprecated, false
	case "bug fixes":
		return CategoryBugFix, false
	case "new features":
		return CategoryFeature, false
	case "minor behavior changes":
		return CategoryConfig, false
	default:
		return CategoryOther, false
	}
}
