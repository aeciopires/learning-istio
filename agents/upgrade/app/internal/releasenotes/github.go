package releasenotes

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ghTagsMaxPages bounds how many pages of the GitHub tags API we will
// walk when hunting for the latest patch tag of a minor series. GitHub
// paginates 100 per page; 10 pages covers even projects with thousands
// of tags without risking a runaway loop against a rate-limited API.
const ghTagsMaxPages = 10

// githubAPIBase is a package-level var (not a const) so unit tests can
// point it at an httptest server instead of the real GitHub API.
var githubAPIBase = "https://api.github.com"

type ghTag struct {
	Name string `json:"name"`
}

// listTags returns every tag name for repo (e.g. "istio/istio"),
// walking pages until GitHub returns an empty page or ghTagsMaxPages is
// reached.
func listTags(ctx context.Context, f *Fetcher, repo string) ([]string, error) {
	var all []string
	for page := 1; page <= ghTagsMaxPages; page++ {
		url := fmt.Sprintf("%s/repos/%s/tags?per_page=100&page=%d", githubAPIBase, repo, page)
		body, err := f.Get(ctx, url)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				break
			}
			return all, err
		}
		var tags []ghTag
		if err := json.Unmarshal(body, &tags); err != nil {
			return all, fmt.Errorf("parsing tags from %s: %w", url, err)
		}
		if len(tags) == 0 {
			break
		}
		for _, t := range tags {
			all = append(all, t.Name)
		}
		if len(tags) < 100 {
			break
		}
	}
	return all, nil
}

// latestPatchForMinor scans tag names (which may or may not have a "v"
// prefix) for the highest PATCH of MAJOR.MINOR.PATCH, ignoring
// pre-release suffixes such as "-rc.1" or "-beta.0".
func latestPatchForMinor(tags []string, major, minor int) (int, bool) {
	prefix := fmt.Sprintf("%d.%d.", major, minor)
	best := -1
	for _, t := range tags {
		t = strings.TrimPrefix(t, "v")
		if !strings.HasPrefix(t, prefix) {
			continue
		}
		rest := strings.TrimPrefix(t, prefix)
		// Stop at the first non-digit rune (pre-release/build metadata).
		end := 0
		for end < len(rest) && rest[end] >= '0' && rest[end] <= '9' {
			end++
		}
		if end == 0 {
			continue
		}
		patch, err := strconv.Atoi(rest[:end])
		if err != nil {
			continue
		}
		if patch > best {
			best = patch
		}
	}
	return best, best >= 0
}

type ghCommit struct {
	Commit struct {
		Message string `json:"message"`
	} `json:"commit"`
}

type ghCompare struct {
	TotalCommits int        `json:"total_commits"`
	Commits      []ghCommit `json:"commits"`
}

// compareCommitSubjects returns the first line of every commit message
// between fromTag (exclusive) and toTag (inclusive) via the GitHub
// compare API, most-recent-last (GitHub's natural order).
func compareCommitSubjects(ctx context.Context, f *Fetcher, repo, fromTag, toTag string) ([]string, error) {
	url := fmt.Sprintf("%s/repos/%s/compare/%s...%s", githubAPIBase, repo, fromTag, toTag)
	body, err := f.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	var cmp ghCompare
	if err := json.Unmarshal(body, &cmp); err != nil {
		return nil, fmt.Errorf("parsing compare response from %s: %w", url, err)
	}
	subjects := make([]string, 0, len(cmp.Commits))
	for _, c := range cmp.Commits {
		msg := strings.TrimSpace(c.Commit.Message)
		if nl := strings.IndexByte(msg, '\n'); nl != -1 {
			msg = msg[:nl]
		}
		if msg != "" {
			subjects = append(subjects, msg)
		}
	}
	return subjects, nil
}
