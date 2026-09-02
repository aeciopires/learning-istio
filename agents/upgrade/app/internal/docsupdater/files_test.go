package docsupdater

import (
	"strings"
	"testing"
)

func testResolvedVersions() ResolvedVersions {
	return ResolvedVersions{
		IstioMinor:    "1.31",
		IstioPatch:    "1.31.1",
		IstioPatchNum: "1",
		GatewayAPI:    "v1.7.0",
		Kind:          "0.33.0",
		Helm:          "3.22.0",
		Kubectl:       "1.37.0",
		MetalLB:       "0.17.0",
		KindNodeK8s:   "1.36.1",
		Grafana:       "12.0.0",
		Prometheus:    "3.0.0",
		Kiali:         "2.0.0",
		Jaeger:        "2.0.0",
	}
}

const readmeFixture = `# learning-istio

# Configurações do Ambiente

As versões dos softwares utilizados foram:

- Docker: 29.2.1
- Kubernetes (kind): 1.35.5
- kubectl: 1.36.2
- Helm: 3.21.2
- Istio: 1.30
- istioctl: 1.30.3

# Requisitos

Veja: https://istio.io/v1.30/docs/setup/platform-setup/docker/
Veja também: https://istio.io/v1.30/docs/releases/supported-releases/
`

func TestUpdateREADME(t *testing.T) {
	rv := testResolvedVersions()
	u := UpdateREADME(readmeFixture, rv)

	if !u.Changed {
		t.Fatal("expected README.md to be marked as changed")
	}
	wantLines := []string{
		"- Kubernetes (kind): 1.36.1",
		"- kubectl: 1.37.0",
		"- Helm: 3.22.0",
		"- Istio: 1.31",
		"- istioctl: 1.31.1",
		"https://istio.io/v1.31/docs/setup/platform-setup/docker/",
		"https://istio.io/v1.31/docs/releases/supported-releases/",
	}
	for _, want := range wantLines {
		if !strings.Contains(u.After, want) {
			t.Errorf("expected updated README to contain %q, got:\n%s", want, u.After)
		}
	}
	// Docker isn't tracked by toolversions and must be left untouched.
	if !strings.Contains(u.After, "- Docker: 29.2.1") {
		t.Errorf("expected Docker line to be left untouched, got:\n%s", u.After)
	}
}

func TestUpdateREADME_NoOpWhenAlreadyCurrent(t *testing.T) {
	rv := ResolvedVersions{IstioMinor: "1.30", IstioPatch: "1.30.3", Kubectl: "1.36.2", Helm: "3.21.2", KindNodeK8s: "1.35.5"}
	u := UpdateREADME(readmeFixture, rv)
	if u.Changed {
		t.Errorf("expected no change when resolved versions match the doc already, got notes: %v", u.Notes)
	}
}

const dayFileFixture = `export ISTIO_RELEASE=1.30
export VERSION_ISTIO="${ISTIO_RELEASE}.2"
export GATEWAY_API_VERSION="v1.6.0"

# https://istio.io/v1.30/docs/setup/getting-started/
`

func TestUpdateDayFile(t *testing.T) {
	rv := testResolvedVersions()
	u := UpdateDayFile("day1.md", dayFileFixture, rv)

	if !u.Changed {
		t.Fatal("expected day1.md to be marked as changed")
	}
	wantLines := []string{
		"export ISTIO_RELEASE=1.31",
		`export VERSION_ISTIO="${ISTIO_RELEASE}.1"`,
		`export GATEWAY_API_VERSION="v1.7.0"`,
		"https://istio.io/v1.31/docs/setup/getting-started/",
	}
	for _, want := range wantLines {
		if !strings.Contains(u.After, want) {
			t.Errorf("expected updated day1.md to contain %q, got:\n%s", want, u.After)
		}
	}
}

const installRequirementsFixture = `# Helm

` + "```bash" + `
VERSION="3.21.2"

mise latest helm
mise install helm@$VERSION
` + "```" + `

# kubectl

` + "```bash" + `
VERSION_OPTION_1="1.36.2"

mise latest kubectl
mise install kubectl@$VERSION_OPTION_1
` + "```" + `

# istioctl

` + "```bash" + `
VERSION="1.30.3"

mise latest istioctl
mise install istioctl@$VERSION
` + "```" + `

# kind

` + "```bash" + `
VERSION="0.32.0"
mise latest kind
mise install kind@$VERSION
` + "```" + `

    image: kindest/node:v1.35.5@sha256:ce977ae6d65918d0b58a5f8b5e940429c2ce42fa3a5619ec2bbc60b949c0ac95

# MetalLB

` + "```bash" + `
export HELM_CHART_VERSION="0.16.1"

helm upgrade --install metallb metallb/metallb \
  --create-namespace \
  --namespace metallb-system \
  --version $HELM_CHART_VERSION \
  --set "controller.tolerations[0].key=node-role.kubernetes.io/master" \
` + "```" + `
`

func TestUpdateInstallRequirements(t *testing.T) {
	rv := testResolvedVersions()
	u := UpdateInstallRequirements(installRequirementsFixture, rv)

	if !u.Changed {
		t.Fatal("expected INSTALL_REQUIREMENTS.md to be marked as changed")
	}
	wantLines := []string{
		`VERSION="3.22.0"` + "\n\nmise latest helm",
		`VERSION_OPTION_1="1.37.0"`,
		`VERSION="1.31.1"` + "\n\nmise latest istioctl",
		`VERSION="0.33.0"` + "\nmise latest kind",
		`export HELM_CHART_VERSION="0.17.0"`,
	}
	for _, want := range wantLines {
		if !strings.Contains(u.After, want) {
			t.Errorf("expected updated INSTALL_REQUIREMENTS.md to contain %q, got:\n%s", want, u.After)
		}
	}
	// The kind node image tag+digest must be left untouched.
	if !strings.Contains(u.After, "kindest/node:v1.35.5@sha256:ce977ae6d65918d0b58a5f8b5e940429c2ce42fa3a5619ec2bbc60b949c0ac95") {
		t.Errorf("expected kind node image to be left untouched, got:\n%s", u.After)
	}
}

func TestExtractKindNodeK8sVersion(t *testing.T) {
	got, ok := ExtractKindNodeK8sVersion(installRequirementsFixture)
	if !ok {
		t.Fatal("expected to find the kind node image version")
	}
	if got != "1.35.5" {
		t.Errorf("got %q, want %q", got, "1.35.5")
	}

	if _, ok := ExtractKindNodeK8sVersion("no image here"); ok {
		t.Error("expected no match when there's no kindest/node image line")
	}
}

const toolVersionsFixture = `kind 0.32.0
helm 3.21.2
kubectl 1.36.2
istioctl 1.30.3
`

func TestUpdateToolVersions(t *testing.T) {
	rv := testResolvedVersions()
	u := UpdateToolVersions(toolVersionsFixture, rv)

	if !u.Changed {
		t.Fatal("expected .tool-versions to be marked as changed")
	}
	want := "kind 0.33.0\nhelm 3.22.0\nkubectl 1.37.0\nistioctl 1.31.1\n"
	if u.After != want {
		t.Errorf("got:\n%s\nwant:\n%s", u.After, want)
	}
}

func TestRenderVersionsEnv(t *testing.T) {
	rv := testResolvedVersions()
	env := RenderVersionsEnv(rv)
	wantLines := []string{
		"ISTIO_MINOR=1.31",
		"ISTIO_VERSION=1.31.1",
		"ISTIOCTL_VERSION=1.31.1",
		"GATEWAY_API_VERSION=v1.7.0",
		"KIND_VERSION=0.33.0",
		"HELM_VERSION=3.22.0",
		"KUBECTL_VERSION=1.37.0",
		"METALLB_VERSION=0.17.0",
		"KIND_NODE_K8S_VERSION=1.36.1",
		"GRAFANA_VERSION=12.0.0",
		"PROMETHEUS_VERSION=3.0.0",
		"KIALI_VERSION=2.0.0",
		"JAEGER_VERSION=2.0.0",
	}
	for _, want := range wantLines {
		if !strings.Contains(env, want) {
			t.Errorf("expected versions.env to contain %q, got:\n%s", want, env)
		}
	}
}
