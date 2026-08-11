package releasenotes

import (
	"encoding/json"
	"strings"
	"testing"
)

func sampleReport() *Report {
	return &Report{
		Generated: "2026-08-11T00:00:00Z",
		From:      "1.29",
		To:        "1.30",
		Istio: []VersionReport{
			{
				Component: ComponentIstio,
				Minor:     "1.30",
				Changes: []Change{
					{Component: ComponentIstio, Version: "1.30", Category: CategoryBreaking, Summary: "CRDs must be upgraded", Breaking: true, Source: "https://example.com/a"},
					{Component: ComponentIstio, Version: "1.30", Category: CategoryFeature, Summary: "New agentgateway support", Source: "https://example.com/b"},
				},
				Warnings: []string{"no page at https://example.com/missing (skipped)"},
			},
		},
	}
}

func TestRenderMarkdown(t *testing.T) {
	md := RenderMarkdown(sampleReport())

	mustContain := []string{
		"# Istio / Envoy Upgrade Report",
		"Istio range: 1.29 -> 1.30",
		"## Istio",
		"### Istio 1.30",
		"**Breaking Changes**",
		"CRDs must be upgraded",
		"**New Features**",
		"New agentgateway support",
		"Warnings",
		"no page at https://example.com/missing",
	}
	for _, s := range mustContain {
		if !strings.Contains(md, s) {
			t.Errorf("expected rendered markdown to contain %q\n---\n%s", s, md)
		}
	}

	// Breaking Changes section must be rendered before New Features.
	if strings.Index(md, "**Breaking Changes**") > strings.Index(md, "**New Features**") {
		t.Errorf("expected Breaking Changes section before New Features section")
	}
}

func TestRenderMarkdown_EnvoyNote(t *testing.T) {
	r := &Report{Generated: "now", EnvoyNote: "explicit range required"}
	md := RenderMarkdown(r)
	if !strings.Contains(md, "explicit range required") {
		t.Errorf("expected EnvoyNote to be rendered, got:\n%s", md)
	}
}

func TestRenderJSON_RoundTrip(t *testing.T) {
	original := sampleReport()
	data, err := RenderJSON(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Report
	if err := json.Unmarshal([]byte(data), &decoded); err != nil {
		t.Fatalf("failed to unmarshal rendered JSON: %v", err)
	}
	if decoded.From != original.From || decoded.To != original.To {
		t.Errorf("round-tripped report mismatch: %+v vs %+v", decoded, original)
	}
	if len(decoded.Istio) != 1 || len(decoded.Istio[0].Changes) != 2 {
		t.Errorf("unexpected round-tripped structure: %+v", decoded)
	}
}
