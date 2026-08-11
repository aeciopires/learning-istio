package releasenotes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// NotFoundError is returned by Fetcher.Get when the remote resource
// responds with 404. Callers can type-assert on this to skip a source
// gracefully (e.g. a patch release that never got its own announcement
// page) instead of failing the whole analyze run.
type NotFoundError struct {
	URL string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.URL)
}

// Fetcher performs HTTP GETs with a simple on-disk cache, so repeated
// `analyze` runs over the same version range don't re-download the same
// pages, and so the tool degrades gracefully when offline (stale cache
// beats a hard failure for a demo/classroom environment).
type Fetcher struct {
	Client    *http.Client
	CacheDir  string
	Refresh   bool // ignore cache reads, but still write fresh responses
	NoCache   bool // never read or write the cache (used by tests)
	UserAgent string
}

// NewFetcher builds a Fetcher. If cacheDir is empty, a default under the
// user's cache directory is used.
func NewFetcher(cacheDir string, refresh bool) *Fetcher {
	if cacheDir == "" {
		if base, err := os.UserCacheDir(); err == nil {
			cacheDir = filepath.Join(base, "istio-upgrade-cli")
		} else {
			cacheDir = filepath.Join(os.TempDir(), "istio-upgrade-cli-cache")
		}
	}
	return &Fetcher{
		Client:    &http.Client{Timeout: 30 * time.Second},
		CacheDir:  cacheDir,
		Refresh:   refresh,
		UserAgent: "istio-upgrade-cli",
	}
}

func (f *Fetcher) cachePath(url string) string {
	sum := sha256.Sum256([]byte(url))
	return filepath.Join(f.CacheDir, hex.EncodeToString(sum[:])+".cache")
}

// Get fetches url, transparently using/populating the on-disk cache. It
// returns *NotFoundError for HTTP 404 responses and a plain error for
// anything else abnormal.
func (f *Fetcher) Get(ctx context.Context, url string) ([]byte, error) {
	cachePath := f.cachePath(url)

	if !f.NoCache && !f.Refresh {
		if data, err := os.ReadFile(cachePath); err == nil {
			return data, nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for %s: %w", url, err)
	}
	if f.UserAgent != "" {
		req.Header.Set("User-Agent", f.UserAgent)
	}
	req.Header.Set("Accept", "application/vnd.github+json, text/plain, */*")

	resp, err := f.Client.Do(req)
	if err != nil {
		// Network failure: fall back to a stale cache entry if we have one.
		if !f.NoCache {
			if data, cacheErr := os.ReadFile(cachePath); cacheErr == nil {
				return data, nil
			}
		}
		return nil, fmt.Errorf("fetching %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &NotFoundError{URL: url}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body from %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s: %s", resp.StatusCode, url, truncate(string(body), 300))
	}

	if !f.NoCache {
		if err := os.MkdirAll(f.CacheDir, 0o755); err == nil {
			_ = os.WriteFile(cachePath, body, 0o644)
		}
	}
	return body, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
