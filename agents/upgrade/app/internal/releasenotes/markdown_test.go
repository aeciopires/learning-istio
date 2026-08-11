package releasenotes

import "testing"

const sampleFrontmatterDoc = `---
title: Announcing Istio 1.30.0
release: 1.30.0
---

The Istio project has released version 1.30.

## Gateway API CRDs must be upgraded

Istio 1.30 upgrades to Gateway API v1.5.1.

## CNI config file permissions changed

Default permissions shifted from 0644 to 0600.
`

func TestSplitMarkdownSections(t *testing.T) {
	sections := splitMarkdownSections(sampleFrontmatterDoc)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d: %+v", len(sections), sections)
	}
	if sections[0].Title != "Gateway API CRDs must be upgraded" {
		t.Errorf("unexpected title: %q", sections[0].Title)
	}
	if sections[1].Title != "CNI config file permissions changed" {
		t.Errorf("unexpected title: %q", sections[1].Title)
	}
	if sections[0].Body == "" || sections[1].Body == "" {
		t.Errorf("expected non-empty bodies, got %+v", sections)
	}
}

func TestSplitMarkdownSections_NoFrontmatter(t *testing.T) {
	doc := "# Title\n\nlede\n\n## First\n\nbody1\n\n## Second\n\nbody2\n"
	sections := splitMarkdownSections(doc)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d: %+v", len(sections), sections)
	}
}

func TestSplitMarkdownSections_Empty(t *testing.T) {
	if got := splitMarkdownSections(""); len(got) != 0 {
		t.Errorf("expected no sections for empty input, got %+v", got)
	}
}
