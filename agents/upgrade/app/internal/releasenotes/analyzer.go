package releasenotes

import (
	"context"
	"fmt"
	"time"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/semverx"
)

// Options configures one Analyzer.Run invocation.
type Options struct {
	// Component selects what to analyze: "istio", "envoy", or "all".
	Component string
	// From/To bound the Istio range (required unless Component == "envoy").
	From, To string
	// EnvoyFrom/EnvoyTo optionally bound an independent Envoy range.
	// Required when Component == "envoy"; optional when Component == "all"
	// (Istio's exact bundled Envoy version can't be reliably derived from
	// public metadata, so without this the report just explains why).
	EnvoyFrom, EnvoyTo string
}

// Analyzer ties semverx range-expansion to the Istio/Envoy clients to
// build a full upgrade-planning Report.
type Analyzer struct {
	Istio *IstioClient
	Envoy *EnvoyClient
}

// NewAnalyzer wires an Analyzer to the real istio.io/envoyproxy.io/GitHub
// endpoints, sharing one on-disk cache.
func NewAnalyzer(f *Fetcher) *Analyzer {
	return &Analyzer{
		Istio: NewIstioClient(f),
		Envoy: NewEnvoyClient(f),
	}
}

// Run executes opts and returns the assembled Report.
func (a *Analyzer) Run(ctx context.Context, opts Options) (*Report, error) {
	report := &Report{
		Generated: time.Now().UTC().Format(time.RFC3339),
		From:      opts.From,
		To:        opts.To,
	}

	wantIstio := opts.Component == "" || opts.Component == "all" || opts.Component == "istio"
	wantEnvoy := opts.Component == "all" || opts.Component == "envoy"

	if wantIstio {
		if opts.From == "" || opts.To == "" {
			return nil, fmt.Errorf("--from and --to are required to analyze Istio")
		}
		fromV, err := semverx.Parse(opts.From)
		if err != nil {
			return nil, fmt.Errorf("--from: %w", err)
		}
		toV, err := semverx.Parse(opts.To)
		if err != nil {
			return nil, fmt.Errorf("--to: %w", err)
		}
		plans, err := semverx.ExpandRange(fromV, toV)
		if err != nil {
			return nil, err
		}
		for _, p := range plans {
			report.Istio = append(report.Istio, a.Istio.FetchMinor(ctx, p))
		}
	}

	if wantEnvoy {
		envFrom, envTo := opts.EnvoyFrom, opts.EnvoyTo
		if envFrom == "" || envTo == "" {
			if opts.Component == "envoy" {
				return nil, fmt.Errorf("--envoy-from and --envoy-to are required when --component=envoy")
			}
			report.EnvoyNote = "Envoy was not analyzed: Istio's bundled Envoy version is pinned via an istio/proxy " +
				"fork commit (see istio.deps), not a clean upstream Envoy release number, so it can't be derived " +
				"reliably. Pass --envoy-from/--envoy-to explicitly to analyze Envoy's own release notes."
		} else {
			fromV, err := semverx.Parse(envFrom)
			if err != nil {
				return nil, fmt.Errorf("--envoy-from: %w", err)
			}
			toV, err := semverx.Parse(envTo)
			if err != nil {
				return nil, fmt.Errorf("--envoy-to: %w", err)
			}
			plans, err := semverx.ExpandRange(fromV, toV)
			if err != nil {
				return nil, err
			}
			for _, p := range plans {
				report.Envoy = append(report.Envoy, a.Envoy.FetchMinor(ctx, p))
			}
		}
	}

	return report, nil
}
