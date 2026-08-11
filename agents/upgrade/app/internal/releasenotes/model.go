// Package releasenotes fetches and parses Istio and Envoy release notes
// and renders an upgrade-planning report across a version range.
package releasenotes

// Category buckets a Change for reporting purposes. Keeping this as a
// small closed set makes the rendered report predictable and lets tests
// assert on exact section membership.
type Category string

const (
	CategoryBreaking     Category = "Breaking Changes"
	CategoryDeprecated   Category = "Deprecated / Removed"
	CategoryInstallation Category = "Installation / Helm"
	CategoryConfig       Category = "Configuration / API / CRDs"
	CategoryMetrics      Category = "Metrics / Telemetry"
	CategoryComponent    Category = "Components"
	CategorySecurity     Category = "Security / CVE Fixes"
	CategoryBugFix       Category = "Bug Fixes"
	CategoryFeature      Category = "New Features"
	CategoryOther        Category = "Other Changes"
)

// Component identifies which project a Change belongs to.
type Component string

const (
	ComponentIstio Component = "Istio"
	ComponentEnvoy Component = "Envoy"
)

// Change is a single, reportable release-note entry.
type Change struct {
	Component Component
	Version   string // e.g. "1.30" for a minor-level note, "1.30.2" for a patch note
	Category  Category
	Summary   string
	Breaking  bool
	Source    string // URL to the original note
}

// VersionReport groups every Change found for one Major.Minor series
// (plus any patch-scoped changes within it).
type VersionReport struct {
	Component Component
	Minor     string // "1.30"
	EnvoyRef  string // best-effort note about the bundled Envoy version/commit, Istio only
	Changes   []Change
	Warnings  []string // non-fatal issues while fetching/parsing this series (e.g. 404s)
}

// Report is the full output of an analyze run, potentially spanning both
// components.
type Report struct {
	Generated string
	From      string
	To        string
	Istio     []VersionReport
	Envoy     []VersionReport

	// EnvoyNote explains, when component=all and no explicit
	// --envoy-from/--envoy-to was given, why Envoy wasn't analyzed
	// automatically (Istio's exact bundled Envoy version isn't reliably
	// derivable from public metadata - see internal/releasenotes/istio.go).
	EnvoyNote string
}
