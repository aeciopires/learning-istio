// Command istio-upgrade analyzes Istio and Envoy release notes across a
// version range and keeps this repository's installation/configuration
// files pinned to the latest stable tool versions. See agents/upgrade/README.md.
package main

import (
	_ "embed"
	"strings"

	"github.com/aeciopires/learning-istio/agents/upgrade/cmd"
)

// version is embedded from the VERSION file at build time (and at `go
// run`/`go test` time too, since go:embed doesn't need ldflags) - the
// single source of truth for the CLI's own version, per the project's
// requirement to control it centrally.
//
//go:embed VERSION
var version string

// gitCommit is optionally set at build time via:
//
//	go build -ldflags "-X main.gitCommit=$(git rev-parse --short HEAD)"
//
// (the Makefile's `build` target does this). It's informational only -
// the CLI's version identity is always the VERSION file above.
var gitCommit = "unknown"

func main() {
	cmd.Execute(strings.TrimSpace(version), gitCommit)
}
