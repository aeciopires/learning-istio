package toolversions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/releasenotes"
)

func TestClient_LatestRelease_FromReleasesAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/istio/istio/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name": "1.30.3"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f := releasenotes.NewFetcher(t.TempDir(), false)
	f.NoCache = true
	c := &Client{Fetcher: f, APIBase: srv.URL}

	got, err := c.LatestRelease(context.Background(), "istio/istio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "1.30.3" {
		t.Errorf("got %q, want %q", got, "1.30.3")
	}
}

func TestClient_LatestRelease_FallsBackToTags(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/metallb/metallb/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	mux.HandleFunc("/repos/metallb/metallb/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"v0.16.1"},{"name":"v0.15.9"},{"name":"v0.17.0-rc1"},{"name":"v0.16.0"}]`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f := releasenotes.NewFetcher(t.TempDir(), false)
	f.NoCache = true
	c := &Client{Fetcher: f, APIBase: srv.URL}

	got, err := c.LatestRelease(context.Background(), "metallb/metallb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "v0.16.1" {
		t.Errorf("got %q, want %q (highest stable tag, ignoring the -rc1 prerelease)", got, "v0.16.1")
	}
}

func TestClient_LatestAll_BestEffort(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/istio/istio/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name": "1.30.3"}`))
	})
	// Every other repo 404s on both releases and tags, to prove one
	// failure doesn't abort the whole batch.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f := releasenotes.NewFetcher(t.TempDir(), false)
	f.NoCache = true
	c := &Client{Fetcher: f, APIBase: srv.URL}

	results := c.LatestAll(context.Background())
	if len(results) != len(Catalog) {
		t.Fatalf("expected %d results, got %d", len(Catalog), len(results))
	}

	var istioResult *Result
	failures := 0
	for i := range results {
		if results[i].Tool.Key == "istio" {
			istioResult = &results[i]
		}
		if results[i].Err != nil {
			failures++
		}
	}
	if istioResult == nil || istioResult.Err != nil || istioResult.Version != "1.30.3" {
		t.Fatalf("expected istio to resolve successfully, got %+v", istioResult)
	}
	if failures != len(Catalog)-1 {
		t.Fatalf("expected every other tool to fail gracefully, got %d failures", failures)
	}
}
