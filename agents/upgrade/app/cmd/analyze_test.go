package cmd

import "testing"

// resetAnalyzeOpts restores analyzeOpts to its flag defaults, since it's
// package-level state shared with the real CLI flags.
func resetAnalyzeOpts() {
	analyzeOpts.component = "all"
	analyzeOpts.from = ""
	analyzeOpts.to = ""
	analyzeOpts.envoyFrom = ""
	analyzeOpts.envoyTo = ""
	analyzeOpts.format = "markdown"
	analyzeOpts.output = ""
	analyzeOpts.cacheDir = ""
	analyzeOpts.refresh = false
}

func TestRunAnalyze_InvalidComponent(t *testing.T) {
	resetAnalyzeOpts()
	t.Cleanup(resetAnalyzeOpts)
	analyzeOpts.component = "bogus"

	if err := runAnalyze(analyzeCmd, nil); err == nil {
		t.Fatal("expected an error for an invalid --component")
	}
}

func TestRunAnalyze_InvalidFormat(t *testing.T) {
	resetAnalyzeOpts()
	t.Cleanup(resetAnalyzeOpts)
	analyzeOpts.format = "yaml"

	if err := runAnalyze(analyzeCmd, nil); err == nil {
		t.Fatal("expected an error for an invalid --format")
	}
}

// Both validations must run before any network call, so these tests stay
// fast and offline (see runAnalyze in analyze.go).
