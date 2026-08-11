package releasenotes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLatestPatchForMinor(t *testing.T) {
	tags := []string{"1.30.0", "1.30.1", "1.30.2", "v1.29.5", "1.30.10", "1.30.3-rc.1", "2.0.0"}
	patch, ok := latestPatchForMinor(tags, 1, 30)
	if !ok {
		t.Fatal("expected to find a patch")
	}
	if patch != 10 {
		t.Fatalf("expected latest patch 10 (must not be confused by the 1.30.3-rc.1 prerelease), got %d", patch)
	}

	if _, ok := latestPatchForMinor(tags, 9, 9); ok {
		t.Fatal("expected no match for a minor series with no tags")
	}
}

func TestListTags_Pagination(t *testing.T) {
	page1 := make([]ghTag, 100)
	for i := range page1 {
		page1[i] = ghTag{Name: fmt.Sprintf("1.30.%d", i)}
	}
	page2 := []ghTag{{Name: "1.30.100"}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")
		switch page {
		case "1":
			_ = json.NewEncoder(w).Encode(page1)
		case "2":
			_ = json.NewEncoder(w).Encode(page2)
		default:
			_ = json.NewEncoder(w).Encode([]ghTag{})
		}
	}))
	defer srv.Close()

	f := NewFetcher(t.TempDir(), false)
	f.NoCache = true
	orig := githubAPIBase
	githubAPIBase = srv.URL
	defer func() { githubAPIBase = orig }()

	tags, err := listTags(context.Background(), f, "istio/istio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 101 {
		t.Fatalf("expected 101 tags across 2 pages, got %d", len(tags))
	}
}

func TestCompareCommitSubjects(t *testing.T) {
	resp := ghCompare{
		TotalCommits: 2,
		Commits: []ghCommit{
			{Commit: struct {
				Message string `json:"message"`
			}{Message: "fix: CVE-2026-1234 auth bypass\n\nmore details here"}},
			{Commit: struct {
				Message string `json:"message"`
			}{Message: "chore: bump deps"}},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	f := NewFetcher(t.TempDir(), false)
	f.NoCache = true
	orig := githubAPIBase
	githubAPIBase = srv.URL
	defer func() { githubAPIBase = orig }()

	subjects, err := compareCommitSubjects(context.Background(), f, "istio/istio", "1.30.0", "1.30.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subjects) != 2 {
		t.Fatalf("expected 2 subjects, got %d: %v", len(subjects), subjects)
	}
	if subjects[0] != "fix: CVE-2026-1234 auth bypass" {
		t.Errorf("expected only the first line of the commit message, got %q", subjects[0])
	}
}
