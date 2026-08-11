package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRepoRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# x"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "agents", "upgrade", "cmd")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := findRepoRoot(nested)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Resolve symlinks (e.g. /tmp -> /private/tmp on macOS) before comparing.
	wantResolved, _ := filepath.EvalSymlinks(root)
	gotResolved, _ := filepath.EvalSymlinks(got)
	if gotResolved != wantResolved {
		t.Errorf("got %q, want %q", got, root)
	}
}

func TestFindRepoRoot_NotFound(t *testing.T) {
	// A fresh temp dir has neither .git nor README.md anywhere above it
	// within the temp hierarchy (walking up eventually hits "/").
	dir := t.TempDir()
	if _, err := findRepoRoot(dir); err == nil {
		t.Error("expected an error when no repo root can be found")
	}
}

func TestFindRepoRoot_GitWithoutReadmeIsNotEnough(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	// No README.md - should NOT be recognized as the repo root.
	if _, err := findRepoRoot(root); err == nil {
		t.Error("expected an error when .git exists but README.md does not")
	}
}
