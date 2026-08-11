package releasenotes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/semverx"
)

func newTestEnvoyClient(t *testing.T, mux *http.ServeMux) *EnvoyClient {
	t.Helper()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	f := NewFetcher(t.TempDir(), false)
	f.NoCache = true

	orig := githubAPIBase
	githubAPIBase = srv.URL
	t.Cleanup(func() { githubAPIBase = orig })

	c := NewEnvoyClient(f)
	c.RSTURLTemplate = srv.URL + "/rst/%s/%s/%s"
	return c
}

const envoyRSTFixture = `1.32.0 (October 15, 2024)
==========================

Incompatible behavior changes
------------------------------

* eds: enabling caching of EDS assignments when used with ADS by default.

Bug fixes
---------

* http: fixed a bug where incomplete requests could be proxied upstream.
`

func TestEnvoyClient_FetchMinor_NoPatchDetail(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rst/1.32.0/1.32/1.32.0", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(envoyRSTFixture))
	})

	c := newTestEnvoyClient(t, mux)
	plan := semverx.MinorPlan{Major: 1, Minor: 32, ToPatch: -1}
	report := c.FetchMinor(context.Background(), plan)

	if len(report.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", report.Warnings)
	}
	if len(report.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d: %+v", len(report.Changes), report.Changes)
	}
	if report.Changes[0].Category != CategoryBreaking || !report.Changes[0].Breaking {
		t.Errorf("expected first change to be breaking, got %+v", report.Changes[0])
	}
	if report.Changes[1].Category != CategoryBugFix {
		t.Errorf("expected second change to be a bug fix, got %+v", report.Changes[1])
	}
}

func TestEnvoyClient_FetchMinor_PatchDetailBounded(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rst/1.32.0/1.32/1.32.0", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(envoyRSTFixture))
	})
	mux.HandleFunc("/rst/1.32.1/1.32/1.32.1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("1.32.1\n=======\n\nBug fixes\n---------\n\n* quic: fixed a crash.\n"))
	})
	mux.HandleFunc("/rst/1.32.2/1.32/1.32.2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	c := newTestEnvoyClient(t, mux)
	plan := semverx.MinorPlan{Major: 1, Minor: 32, IncludePatchDetail: true, FromPatch: 0, ToPatch: 2}
	report := c.FetchMinor(context.Background(), plan)

	if len(report.Warnings) != 1 {
		t.Fatalf("expected 1 warning for the missing 1.32.2 page, got %v", report.Warnings)
	}
	// 2 changes from 1.32.0 baseline + 1 change from 1.32.1 = 3
	if len(report.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d: %+v", len(report.Changes), report.Changes)
	}
	if report.Changes[2].Version != "1.32.1" {
		t.Errorf("expected the patch-level change to be tagged 1.32.1, got %q", report.Changes[2].Version)
	}
}
