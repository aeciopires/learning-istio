package releasenotes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/semverx"
)

func newTestIstioClient(t *testing.T, mux *http.ServeMux) *IstioClient {
	t.Helper()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	f := NewFetcher(t.TempDir(), false)
	f.NoCache = true

	orig := githubAPIBase
	githubAPIBase = srv.URL
	t.Cleanup(func() { githubAPIBase = orig })

	c := NewIstioClient(f)
	c.HighlightsURLTemplate = srv.URL + "/highlights/%s/%s"
	c.UpgradeNotesURLTemplate = srv.URL + "/upgrade-notes/%s/%s"
	return c
}

func TestIstioClient_FetchMinor_HighlightsAndUpgradeNotes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/highlights/1.30/1.30", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("---\ntitle: x\n---\n\n## New agentgateway support\n\nExperimental Gateway API implementation.\n"))
	})
	mux.HandleFunc("/upgrade-notes/1.30/1.30", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("## Gateway API CRDs must be upgraded\n\nUpgrade CRDs to v1.5.x before upgrading Istio.\n"))
	})

	c := newTestIstioClient(t, mux)
	plan := semverx.MinorPlan{Major: 1, Minor: 30, ToPatch: -1}
	report := c.FetchMinor(context.Background(), plan)

	if len(report.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", report.Warnings)
	}
	if len(report.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d: %+v", len(report.Changes), report.Changes)
	}

	var sawBreaking, sawFeature bool
	for _, c := range report.Changes {
		if c.Breaking {
			sawBreaking = true
			if c.Category != CategoryConfig && c.Category != CategoryBreaking {
				t.Errorf("expected breaking change to classify as Config or Breaking, got %v", c.Category)
			}
		} else {
			sawFeature = true
		}
	}
	if !sawBreaking || !sawFeature {
		t.Errorf("expected both a breaking and a non-breaking change, got %+v", report.Changes)
	}
}

func TestIstioClient_FetchMinor_MissingHighlightsPage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/highlights/1.10/1.10", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	mux.HandleFunc("/upgrade-notes/1.10/1.10", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("## Something changed\n\ndetails\n"))
	})

	c := newTestIstioClient(t, mux)
	plan := semverx.MinorPlan{Major: 1, Minor: 10, ToPatch: -1}
	report := c.FetchMinor(context.Background(), plan)

	if len(report.Warnings) != 1 {
		t.Fatalf("expected exactly 1 warning for the missing page, got %v", report.Warnings)
	}
	if len(report.Changes) != 1 {
		t.Fatalf("expected 1 change from the upgrade-notes page, got %d", len(report.Changes))
	}
}

func TestIstioClient_FetchMinor_PatchDetailBounded(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/highlights/1.29/1.29", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	mux.HandleFunc("/upgrade-notes/1.29/1.29", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	mux.HandleFunc("/repos/istio/istio/compare/1.29.0...1.29.3", func(w http.ResponseWriter, r *http.Request) {
		resp := ghCompare{Commits: []ghCommit{
			{Commit: struct {
				Message string `json:"message"`
			}{Message: "fix: security CVE-2026-9999 patch"}},
			{Commit: struct {
				Message string `json:"message"`
			}{Message: "fix: minor bug in retry logic"}},
		}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	c := newTestIstioClient(t, mux)
	plan := semverx.MinorPlan{Major: 1, Minor: 29, IncludePatchDetail: true, FromPatch: 0, ToPatch: 3}
	report := c.FetchMinor(context.Background(), plan)

	if len(report.Changes) != 2 {
		t.Fatalf("expected 2 patch-level changes, got %d: %+v", len(report.Changes), report.Changes)
	}
	for _, ch := range report.Changes {
		if ch.Version != "1.29.3" {
			t.Errorf("expected patch changes tagged with the target patch version, got %q", ch.Version)
		}
	}
	if report.Changes[0].Category != CategorySecurity {
		t.Errorf("expected CVE-mentioning commit to classify as Security, got %v", report.Changes[0].Category)
	}
	if report.Changes[1].Category != CategoryBugFix {
		t.Errorf("expected plain fix commit to classify as BugFix, got %v", report.Changes[1].Category)
	}
}

func TestIstioClient_FetchMinor_PatchDetailUnboundedDiscoversLatest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/highlights/1.29/1.29", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	mux.HandleFunc("/upgrade-notes/1.29/1.29", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	mux.HandleFunc("/repos/istio/istio/tags", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")
		if page == "1" {
			_ = json.NewEncoder(w).Encode([]ghTag{{Name: "1.29.0"}, {Name: "1.29.1"}, {Name: "1.29.4"}})
			return
		}
		_ = json.NewEncoder(w).Encode([]ghTag{})
	})
	var gotComparePath string
	mux.HandleFunc("/repos/istio/istio/compare/", func(w http.ResponseWriter, r *http.Request) {
		gotComparePath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ghCompare{})
	})

	c := newTestIstioClient(t, mux)
	plan := semverx.MinorPlan{Major: 1, Minor: 29, IncludePatchDetail: true, FromPatch: 2, ToPatch: -1}
	c.FetchMinor(context.Background(), plan)

	want := "/repos/istio/istio/compare/1.29.2...1.29.4"
	if gotComparePath != want {
		t.Fatalf("expected compare call %q (latest patch discovered via tags), got %q", want, gotComparePath)
	}
}

func TestIsBoilerplateHeading(t *testing.T) {
	boilerplate := []string{
		"Upgrading to 1.29",
		"Get Involved",
		"Tell Us How We're Doing",
		"Istio 1.11 Upgrade Survey",
		"Join the Istio community",
	}
	for _, title := range boilerplate {
		if !isBoilerplateHeading(title) {
			t.Errorf("expected %q to be classified as boilerplate", title)
		}
	}

	real := []string{
		"Gateway API CRDs must be upgraded",
		"WebAssembly API",
		"And much, much more",
	}
	for _, title := range real {
		if isBoilerplateHeading(title) {
			t.Errorf("expected %q to NOT be classified as boilerplate", title)
		}
	}
}

func TestIstioClient_ClassifyBreakingSection(t *testing.T) {
	if got := classifyByKeywords("Gateway API CRDs must be upgraded", CategoryBreaking); got != CategoryConfig {
		t.Errorf("expected Gateway API heading to classify as Config, got %v", got)
	}
	if got := classifyByKeywords("XDS debug endpoints now require authentication", CategoryBreaking); got != CategorySecurity {
		t.Errorf("expected auth heading to classify as Security, got %v", got)
	}
	if got := classifyByKeywords("Something with no keyword match", CategoryBreaking); got != CategoryBreaking {
		t.Errorf("expected fallback to CategoryBreaking, got %v", got)
	}
}
