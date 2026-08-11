package releasenotes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyzer_Run_RequiresFromTo(t *testing.T) {
	a := &Analyzer{Istio: NewIstioClient(NewFetcher(t.TempDir(), false)), Envoy: NewEnvoyClient(NewFetcher(t.TempDir(), false))}
	_, err := a.Run(context.Background(), Options{Component: "istio", From: "1.30"})
	if err == nil {
		t.Fatal("expected error when --to is missing")
	}
}

func TestAnalyzer_Run_EnvoyRequiresExplicitRange(t *testing.T) {
	a := &Analyzer{Istio: NewIstioClient(NewFetcher(t.TempDir(), false)), Envoy: NewEnvoyClient(NewFetcher(t.TempDir(), false))}
	_, err := a.Run(context.Background(), Options{Component: "envoy"})
	if err == nil {
		t.Fatal("expected error when component=envoy without --envoy-from/--envoy-to")
	}
}

func TestAnalyzer_Run_AllWithoutEnvoyRangeSetsNote(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f := NewFetcher(t.TempDir(), false)
	f.NoCache = true
	istioClient := NewIstioClient(f)
	istioClient.HighlightsURLTemplate = srv.URL + "/h/%s/%s"
	istioClient.UpgradeNotesURLTemplate = srv.URL + "/u/%s/%s"

	a := &Analyzer{Istio: istioClient, Envoy: NewEnvoyClient(f)}
	report, err := a.Run(context.Background(), Options{Component: "all", From: "1.30", To: "1.30"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.EnvoyNote == "" {
		t.Error("expected EnvoyNote to explain why Envoy wasn't analyzed")
	}
	if len(report.Envoy) != 0 {
		t.Errorf("expected no Envoy version reports, got %+v", report.Envoy)
	}
	if len(report.Istio) != 1 {
		t.Errorf("expected 1 Istio version report, got %d", len(report.Istio))
	}
}

func TestAnalyzer_Run_InvalidVersion(t *testing.T) {
	a := &Analyzer{Istio: NewIstioClient(NewFetcher(t.TempDir(), false)), Envoy: NewEnvoyClient(NewFetcher(t.TempDir(), false))}
	_, err := a.Run(context.Background(), Options{Component: "istio", From: "not-a-version", To: "1.30"})
	if err == nil {
		t.Fatal("expected error for invalid --from")
	}
}
