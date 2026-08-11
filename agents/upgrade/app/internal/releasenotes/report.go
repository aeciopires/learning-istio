package releasenotes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// categoryOrder controls section ordering in the rendered markdown
// report: the things an upgrade planner cares about most come first.
var categoryOrder = []Category{
	CategoryBreaking,
	CategoryDeprecated,
	CategoryInstallation,
	CategoryConfig,
	CategoryMetrics,
	CategoryComponent,
	CategorySecurity,
	CategoryBugFix,
	CategoryFeature,
	CategoryOther,
}

// RenderMarkdown renders r as a human-readable markdown report.
func RenderMarkdown(r *Report) string {
	var b strings.Builder
	b.WriteString("# Istio / Envoy Upgrade Report\n\n")
	fmt.Fprintf(&b, "- Generated: %s\n", r.Generated)
	if r.From != "" || r.To != "" {
		fmt.Fprintf(&b, "- Istio range: %s -> %s\n", r.From, r.To)
	}
	b.WriteString("\n")

	if len(r.Istio) > 0 {
		b.WriteString("## Istio\n\n")
		for _, vr := range r.Istio {
			renderVersionReport(&b, vr)
		}
	}

	if len(r.Envoy) > 0 {
		b.WriteString("## Envoy\n\n")
		for _, vr := range r.Envoy {
			renderVersionReport(&b, vr)
		}
	}

	if r.EnvoyNote != "" {
		b.WriteString("## Envoy\n\n")
		fmt.Fprintf(&b, "> %s\n\n", r.EnvoyNote)
	}

	return b.String()
}

func renderVersionReport(b *strings.Builder, vr VersionReport) {
	fmt.Fprintf(b, "### %s %s\n\n", vr.Component, vr.Minor)

	byCategory := make(map[Category][]Change)
	for _, c := range vr.Changes {
		byCategory[c.Category] = append(byCategory[c.Category], c)
	}

	if len(vr.Changes) == 0 {
		b.WriteString("_No changes parsed for this version (see warnings below, if any)._\n\n")
	}

	for _, cat := range categoryOrder {
		changes := byCategory[cat]
		if len(changes) == 0 {
			continue
		}
		fmt.Fprintf(b, "**%s**\n\n", cat)
		for _, c := range changes {
			marker := ""
			if c.Breaking && cat != CategoryBreaking {
				marker = " ⚠️"
			}
			fmt.Fprintf(b, "- [%s]%s %s ([source](%s))\n", c.Version, marker, c.Summary, c.Source)
		}
		b.WriteString("\n")
	}

	if len(vr.Warnings) > 0 {
		b.WriteString("<details><summary>Warnings</summary>\n\n")
		for _, w := range vr.Warnings {
			fmt.Fprintf(b, "- %s\n", w)
		}
		b.WriteString("\n</details>\n\n")
	}
}

// RenderJSON renders r as indented JSON.
func RenderJSON(r *Report) (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
