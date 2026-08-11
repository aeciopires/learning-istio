package releasenotes

import "testing"

const sampleEnvoyRST = `1.32.0 (October 15, 2024)
==========================

Incompatible behavior changes
------------------------------
*Changes that are expected to cause an incompatibility if applicable*

* eds: enabling caching of EDS assignments when used with ADS by default.
* golang: filter methods now provide access to request/response headers
  and trailers during logging.

Bug fixes
---------
*Changes expected to improve the state of the world*

* http: fixed a bug where incomplete requests could be proxied upstream.

New features
------------
* access log: added a new JSON formatter.
`

func TestParseRSTSections(t *testing.T) {
	sections := parseRSTSections(sampleEnvoyRST)
	if len(sections) != 4 {
		t.Fatalf("expected 4 sections (version heading + 3 categories), got %d: %+v", len(sections), sections)
	}

	// First section is the version title itself; it has no top-level bullets.
	if len(sections[0].Bullets) != 0 {
		t.Errorf("expected version heading to have no bullets, got %+v", sections[0].Bullets)
	}

	incompatible := sections[1]
	if incompatible.Title != "Incompatible behavior changes" {
		t.Fatalf("unexpected title: %q", incompatible.Title)
	}
	if len(incompatible.Bullets) != 2 {
		t.Fatalf("expected 2 bullets, got %d: %+v", len(incompatible.Bullets), incompatible.Bullets)
	}
	if incompatible.Bullets[1] != "golang: filter methods now provide access to request/response headers and trailers during logging." {
		t.Errorf("wrapped continuation line not joined correctly: %q", incompatible.Bullets[1])
	}

	bugFixes := sections[2]
	if bugFixes.Title != "Bug fixes" || len(bugFixes.Bullets) != 1 {
		t.Fatalf("unexpected bug fixes section: %+v", bugFixes)
	}

	features := sections[3]
	if features.Title != "New features" || len(features.Bullets) != 1 {
		t.Fatalf("unexpected new features section: %+v", features)
	}
}

func TestIsUnderline(t *testing.T) {
	cases := map[string]bool{
		"===":     true,
		"---":     true,
		"~~~~":    true,
		"--":      false,
		"":        false,
		"abc":     false,
		"-=-":     false,
		"  ---  ": true,
	}
	for in, want := range cases {
		if got := isUnderline(in); got != want {
			t.Errorf("isUnderline(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestEnvoyCategoryFor(t *testing.T) {
	cat, breaking := envoyCategoryFor("Incompatible behavior changes")
	if cat != CategoryBreaking || !breaking {
		t.Errorf("got %v/%v, want CategoryBreaking/true", cat, breaking)
	}
	cat, breaking = envoyCategoryFor("Bug fixes")
	if cat != CategoryBugFix || breaking {
		t.Errorf("got %v/%v, want CategoryBugFix/false", cat, breaking)
	}
	cat, _ = envoyCategoryFor("Something unexpected")
	if cat != CategoryOther {
		t.Errorf("got %v, want CategoryOther", cat)
	}
}
