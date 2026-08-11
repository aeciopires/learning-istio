package releasenotes

import (
	"context"
	"fmt"
	"strings"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/semverx"
)

// IstioClient fetches and parses Istio release notes. Every URL template
// is a struct field (not a constant) so tests can point it at a local
// httptest server instead of the real istio.io / GitHub hosts.
type IstioClient struct {
	Fetcher *Fetcher

	// HighlightsURLTemplate/UpgradeNotesURLTemplate accept one %s arg: the
	// "MAJOR.MINOR" version, e.g. "1.30".
	HighlightsURLTemplate   string
	UpgradeNotesURLTemplate string

	// Repo is the GitHub "owner/name" used for tag discovery and the
	// compare API (patch-level bug-fix/CVE listing).
	Repo string
}

// NewIstioClient returns a client wired to the real istio.io / GitHub
// endpoints.
func NewIstioClient(f *Fetcher) *IstioClient {
	return &IstioClient{
		Fetcher:                 f,
		HighlightsURLTemplate:   "https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/%s.x/announcing-%s/_index.md",
		UpgradeNotesURLTemplate: "https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/%s.x/announcing-%s/upgrade-notes/index.md",
		Repo:                    "istio/istio",
	}
}

// FetchMinor builds the VersionReport for one Major.Minor Istio series,
// including patch-scoped bug-fix/security entries when plan requests them.
func (c *IstioClient) FetchMinor(ctx context.Context, plan semverx.MinorPlan) VersionReport {
	minor := plan.MinorString()
	report := VersionReport{Component: ComponentIstio, Minor: minor}

	c.fetchAndAppend(ctx, &report, fmt.Sprintf(c.HighlightsURLTemplate, minor, minor), CategoryFeature, false)
	c.fetchAndAppend(ctx, &report, fmt.Sprintf(c.UpgradeNotesURLTemplate, minor, minor), CategoryBreaking, true)

	if plan.IncludePatchDetail {
		c.fetchPatchDetail(ctx, &report, plan)
	}

	return report
}

func (c *IstioClient) fetchAndAppend(ctx context.Context, report *VersionReport, url string, defaultCategory Category, breaking bool) {
	body, err := c.Fetcher.Get(ctx, url)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			report.Warnings = append(report.Warnings, fmt.Sprintf("no page at %s (skipped)", url))
			return
		}
		report.Warnings = append(report.Warnings, fmt.Sprintf("fetching %s: %v", url, err))
		return
	}
	for _, sec := range splitMarkdownSections(string(body)) {
		if sec.Title == "" || sec.Body == "" || isBoilerplateHeading(sec.Title) {
			continue
		}
		report.Changes = append(report.Changes, Change{
			Component: ComponentIstio,
			Version:   report.Minor,
			Category:  classifyByKeywords(sec.Title, defaultCategory),
			Summary:   fmt.Sprintf("%s: %s", sec.Title, firstSentence(sec.Body)),
			Breaking:  breaking,
			Source:    url,
		})
	}
}

func (c *IstioClient) fetchPatchDetail(ctx context.Context, report *VersionReport, plan semverx.MinorPlan) {
	toPatch := plan.ToPatch
	if toPatch == -1 {
		tags, err := listTags(ctx, c.Fetcher, c.Repo)
		if err != nil {
			report.Warnings = append(report.Warnings, fmt.Sprintf("listing tags for %s: %v", c.Repo, err))
			return
		}
		latest, ok := latestPatchForMinor(tags, plan.Major, plan.Minor)
		if !ok {
			report.Warnings = append(report.Warnings, fmt.Sprintf("could not determine latest patch for %s.x", report.Minor))
			return
		}
		toPatch = latest
	}
	if toPatch <= plan.FromPatch {
		return
	}

	fromTag := fmt.Sprintf("%d.%d.%d", plan.Major, plan.Minor, plan.FromPatch)
	toTag := fmt.Sprintf("%d.%d.%d", plan.Major, plan.Minor, toPatch)
	subjects, err := compareCommitSubjects(ctx, c.Fetcher, c.Repo, fromTag, toTag)
	if err != nil {
		report.Warnings = append(report.Warnings, fmt.Sprintf("comparing %s...%s: %v", fromTag, toTag, err))
		return
	}
	compareURL := fmt.Sprintf("https://github.com/%s/compare/%s...%s", c.Repo, fromTag, toTag)
	for _, subject := range subjects {
		report.Changes = append(report.Changes, Change{
			Component: ComponentIstio,
			Version:   toTag,
			Category:  classifyByKeywords(subject, CategoryBugFix),
			Summary:   subject,
			Breaking:  false,
			Source:    compareURL,
		})
	}
}

// boilerplateHeadingPrefixes lists H2 sections that appear on every
// istio.io announcement page but aren't actual release notes (a
// closing "please give us feedback" call-to-action), so they'd
// otherwise show up as a spurious "New Features" entry.
var boilerplateHeadingPrefixes = []string{
	"upgrading to ",
	"get involved",
	"thank you",
	"acknowledgements",
	"tell us how",
	"join the istio community",
	"join the community",
}

// boilerplateHeadingSubstrings catches boilerplate whose exact wording
// varies by release but always includes one of these phrases, e.g.
// "Istio 1.11 Upgrade Survey".
var boilerplateHeadingSubstrings = []string{
	"upgrade survey",
}

func isBoilerplateHeading(title string) bool {
	t := strings.ToLower(strings.TrimSpace(title))
	for _, prefix := range boilerplateHeadingPrefixes {
		if strings.HasPrefix(t, prefix) {
			return true
		}
	}
	for _, sub := range boilerplateHeadingSubstrings {
		if strings.Contains(t, sub) {
			return true
		}
	}
	return false
}

func firstSentence(s string) string {
	// Keep summaries compact: take up to the first blank line / ~280 chars.
	const maxLen = 280
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
