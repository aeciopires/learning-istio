package cmd

import "testing"

func TestRunVersionsLatest_InvalidFormat(t *testing.T) {
	orig := versionsOpts.format
	t.Cleanup(func() { versionsOpts.format = orig })
	versionsOpts.format = "xml"

	if err := runVersionsLatest(versionsLatestCmd, nil); err == nil {
		t.Fatal("expected an error for an invalid --format")
	}
}
