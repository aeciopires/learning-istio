package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecute_VersionFlag(t *testing.T) {
	origVersion, origGitCommit := cliVersion, cliGitCommit
	t.Cleanup(func() { cliVersion, cliGitCommit = origVersion, origGitCommit })

	cliVersion = "9.9.9"
	rootCmd.Version = cliVersion

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "9.9.9") {
		t.Errorf("expected version output to contain 9.9.9, got %q", buf.String())
	}
}

func TestExecute_HelpFlag(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "istio-upgrade") {
		t.Errorf("expected help output to mention istio-upgrade, got %q", buf.String())
	}
}

func TestVersionCommand(t *testing.T) {
	origVersion, origGitCommit := cliVersion, cliGitCommit
	t.Cleanup(func() { cliVersion, cliGitCommit = origVersion, origGitCommit })

	cliVersion = "1.2.3"
	cliGitCommit = "abc123"

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "1.2.3") || !strings.Contains(got, "abc123") {
		t.Errorf("unexpected version command output: %q", got)
	}
}
