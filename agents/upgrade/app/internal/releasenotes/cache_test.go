package releasenotes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestFetcher_GetAndCache(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	f := NewFetcher(t.TempDir(), false)
	body, err := f.Get(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("unexpected body: %q", body)
	}

	// Second call should be served from cache, not hit the server again.
	body2, err := f.Get(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body2) != "hello" {
		t.Fatalf("unexpected body: %q", body2)
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("expected exactly 1 server hit due to caching, got %d", got)
	}
}

func TestFetcher_Refresh(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("v"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := NewFetcher(dir, false)
	if _, err := f.Get(context.Background(), srv.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f2 := NewFetcher(dir, true) // Refresh=true
	if _, err := f2.Get(context.Background(), srv.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("expected 2 server hits with refresh enabled, got %d", got)
	}
}

func TestFetcher_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	f := NewFetcher(t.TempDir(), false)
	_, err := f.Get(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	if _, ok := err.(*NotFoundError); !ok {
		t.Fatalf("expected *NotFoundError, got %T: %v", err, err)
	}
}

func TestFetcher_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	f := NewFetcher(t.TempDir(), false)
	_, err := f.Get(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
	if _, ok := err.(*NotFoundError); ok {
		t.Fatalf("500 should not be treated as NotFoundError")
	}
}
