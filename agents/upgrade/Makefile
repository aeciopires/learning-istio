SHELL := /usr/bin/env bash

APP_DIR        := app
BINARY_NAME    := istio-upgrade
VERSION        := $(shell cat $(APP_DIR)/VERSION)
GIT_COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DIR      := bin
DOCKER_IMAGE   := istio-upgrade-cli:$(VERSION)
REPO_ROOT      := $(abspath ../..)
MODE           ?= ambient
ARGS           ?=

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@echo "istio-upgrade CLI - version $(VERSION)"
	@echo "(Go source lives in $(APP_DIR)/; see README.md for details)"
	@echo ""
	@echo "Targets:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2}'

## --- Build / test the Go CLI (source in app/) ------------------------------

.PHONY: build
build: ## Build the istio-upgrade binary for the host OS/arch into bin/
	go build -C $(APP_DIR) -ldflags "-X main.gitCommit=$(GIT_COMMIT)" -o ../$(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: install
install: ## go install the CLI into $GOBIN/$GOPATH/bin
	go install -C $(APP_DIR) .

## Cross-compilation: pure Go, no cgo, so linux/darwin x amd64/arm64 all
## build from any host with no extra toolchain. Docker stays Linux-only on
## purpose (Docker Desktop on macOS runs Linux containers via a VM - there
## is no such thing as a "macOS Docker image").
.PHONY: build-linux-amd64
build-linux-amd64: ## Cross-build for Linux amd64 into bin/
	GOOS=linux GOARCH=amd64 go build -C $(APP_DIR) -ldflags "-X main.gitCommit=$(GIT_COMMIT)" -o ../$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

.PHONY: build-linux-arm64
build-linux-arm64: ## Cross-build for Linux arm64 into bin/
	GOOS=linux GOARCH=arm64 go build -C $(APP_DIR) -ldflags "-X main.gitCommit=$(GIT_COMMIT)" -o ../$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

.PHONY: build-darwin-amd64
build-darwin-amd64: ## Cross-build for macOS amd64 (Intel) into bin/
	GOOS=darwin GOARCH=amd64 go build -C $(APP_DIR) -ldflags "-X main.gitCommit=$(GIT_COMMIT)" -o ../$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Cross-build for macOS arm64 (Apple Silicon) into bin/
	GOOS=darwin GOARCH=arm64 go build -C $(APP_DIR) -ldflags "-X main.gitCommit=$(GIT_COMMIT)" -o ../$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

.PHONY: build-all
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 ## Cross-build for all supported OS/arch combos into bin/

.PHONY: run
run: ## Run the CLI from source, e.g. `make run ARGS="analyze --from 1.10 --to 1.30"`
	go run -C $(APP_DIR) . $(ARGS)

.PHONY: fmt
fmt: ## gofmt every .go file in place
	gofmt -w $(APP_DIR)

.PHONY: fmt-check
fmt-check: ## Fail if any .go file is not gofmt'd
	@test -z "$$(gofmt -l $(APP_DIR))" || { echo "not gofmt'd:"; gofmt -l $(APP_DIR); exit 1; }

.PHONY: vet
vet: ## go vet ./...
	go vet -C $(APP_DIR) ./...

.PHONY: test
test: ## Run unit tests (fast, offline, no docker/kind/helm/kubectl required)
	go test -C $(APP_DIR) ./... -count=1

.PHONY: test-verbose
test-verbose: ## Run unit tests with -v
	go test -C $(APP_DIR) ./... -count=1 -v

.PHONY: coverage
coverage: ## Run unit tests with coverage report
	go test -C $(APP_DIR) ./... -count=1 -coverprofile=coverage.out
	go -C $(APP_DIR) tool cover -func=coverage.out

.PHONY: lint
lint: fmt-check vet ## fmt-check + vet, plus golangci-lint if installed
	@if command -v golangci-lint >/dev/null 2>&1; then \
		(cd $(APP_DIR) && golangci-lint run ./...); \
	else \
		echo "golangci-lint not installed, skipping (fmt-check + vet already ran)"; \
	fi

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) $(APP_DIR)/coverage.out

.PHONY: version
version: ## Print the CLI version (from app/VERSION)
	@echo $(VERSION)

## --- Docker (build context is app/) ----------------------------------------

.PHONY: docker-build
docker-build: ## Build the istio-upgrade-cli Docker image
	docker build -f $(APP_DIR)/Dockerfile -t $(DOCKER_IMAGE) $(APP_DIR)

.PHONY: docker-run
docker-run: ## Run the Docker image, e.g. `make docker-run ARGS="--help"`
	docker run --rm $(DOCKER_IMAGE) $(ARGS)

## --- Release note analysis / repo doc sync (the Go CLI) -------------------

.PHONY: analyze
analyze: ## Analyze release notes, e.g. `make analyze ARGS="--from 1.10 --to 1.30"`
	go run -C $(APP_DIR) . analyze $(ARGS)

.PHONY: versions-latest
versions-latest: ## Print the latest stable version of every tracked tool
	go run -C $(APP_DIR) . versions latest

.PHONY: docs-update
docs-update: ## Dry-run: show what `docs update` would change in the repo
	go run -C $(APP_DIR) . docs update --repo-root $(REPO_ROOT)

.PHONY: docs-update-apply
docs-update-apply: ## Apply `docs update` changes to the repo (writes files)
	go run -C $(APP_DIR) . docs update --apply --repo-root $(REPO_ROOT)

## --- JWT helper (the Go CLI's "jwt" subcommand; offline) -------------------
## Mints one token per application for the JWT authorization walkthrough
## (day7.md). Built on github.com/golang-jwt/jwt. Keys and tokens are for
## study clusters only; JWT_DIR is git-ignored.

JWT_DIR      ?= $(REPO_ROOT)/.jwt
JWT_KID      ?= learning-istio
JWT_ISSUER   ?= learning-istio@example.com
JWT_SUBJECT  ?= curl
JWT_EXPIRES  ?= 1h
JWT_CLAIMS   ?=
JWT_AUDIENCE ?=
TOKEN        ?=

.PHONY: jwt-keygen
jwt-keygen: ## Generate the RSA key pair + JWKS into JWT_DIR (refuses to overwrite; add ARGS=--force to replace)
	@go run -C $(APP_DIR) . jwt keygen --out-dir $(JWT_DIR) --kid $(JWT_KID) $(ARGS)

.PHONY: jwt-jwks
jwt-jwks: ## Print the JWKS (for RequestAuthentication's "jwks" field)
	@go run -C $(APP_DIR) . jwt jwks --public-key $(JWT_DIR)/public.pem --kid $(JWT_KID)

.PHONY: jwt-sign
jwt-sign: ## Print a signed token, e.g. `make jwt-sign JWT_SUBJECT=curl JWT_CLAIMS="--claim app=curl --claim-list groups=readers"`
	@go run -C $(APP_DIR) . jwt sign --private-key $(JWT_DIR)/private.pem --kid $(JWT_KID) \
		--issuer $(JWT_ISSUER) --subject $(JWT_SUBJECT) --expires-in $(JWT_EXPIRES) \
		$(if $(JWT_AUDIENCE),--audience $(JWT_AUDIENCE)) $(JWT_CLAIMS)

.PHONY: jwt-decode
jwt-decode: ## Decode AND verify a token against JWT_DIR's public key, e.g. `make jwt-decode TOKEN=eyJ...`
	@test -n "$(TOKEN)" || { echo "usage: make jwt-decode TOKEN=<jwt>"; exit 1; }
	@go run -C $(APP_DIR) . jwt decode --token $(TOKEN) --public-key $(JWT_DIR)/public.pem

## --- Host tooling (kind/helm/kubectl, via mise) -----------------------------
## See INSTALL_REQUIREMENTS.md#mise for how to install mise itself.

.PHONY: install-tools
install-tools: ## Install kind/helm/kubectl/istioctl pinned in $(REPO_ROOT)/mise.toml via mise
	@command -v mise >/dev/null 2>&1 || { echo "mise not found on PATH - see INSTALL_REQUIREMENTS.md#mise"; exit 1; }
	# Scoped to these 4 tools on purpose: a bare "mise install" also installs
	# whatever is pinned in the user's *global* mise config, not just this repo.
	# mise itself is idempotent here: a tool version already installed is
	# reported (not re-downloaded) and this never fails because of it.
	cd $(REPO_ROOT) && mise install helm kubectl kind istioctl

.PHONY: install-go
install-go: ## Install the Go toolchain pinned in $(REPO_ROOT)/mise.toml via mise (only needed to build the CLI)
	@command -v mise >/dev/null 2>&1 || { echo "mise not found on PATH - see INSTALL_REQUIREMENTS.md#mise"; exit 1; }
	cd $(REPO_ROOT) && mise install go

## --- Cluster lifecycle (kind + MetalLB + Istio + Bookinfo) -----------------
## NOTE: these targets shell out to scripts/*.sh, NOT the Go CLI - they
## create real Docker containers and a real Kubernetes cluster.
## All of them are safe to re-run: each checks whether what it's about to
## create/install already exists, logs an INFO line either way, and
## reconciles in place instead of failing - see scripts/lib/common.sh.

.PHONY: preflight
preflight: ## Check docker/kind/helm/kubectl/istioctl are installed and reachable
	scripts/preflight.sh

.PHONY: create-cluster
create-cluster: ## Create the local kind cluster
	scripts/create-cluster.sh

.PHONY: install-metallb
install-metallb: ## Install and configure MetalLB on the kind cluster
	scripts/install-metallb.sh

.PHONY: install-istio-sidecar
install-istio-sidecar: ## Install Istio in sidecar mode
	scripts/install-istio.sh sidecar

.PHONY: install-istio-ambient
install-istio-ambient: ## Install Istio in ambient mode
	scripts/install-istio.sh ambient

.PHONY: deploy-bookinfo-sidecar
deploy-bookinfo-sidecar: ## Deploy the Bookinfo sample app for sidecar mode
	scripts/deploy-bookinfo.sh sidecar

.PHONY: deploy-bookinfo-ambient
deploy-bookinfo-ambient: ## Deploy the Bookinfo sample app for ambient mode
	scripts/deploy-bookinfo.sh ambient

.PHONY: install-addons
install-addons: ## Install the Istio addons (Grafana, Prometheus, Kiali, Jaeger)
	scripts/install-addons.sh

.PHONY: validate
validate: ## Validate the cluster/Istio/Bookinfo, e.g. `make validate MODE=sidecar`
	scripts/validate-cluster.sh $(MODE)

.PHONY: destroy-cluster
destroy-cluster: ## Delete the local kind cluster
	scripts/destroy-cluster.sh

.PHONY: bootstrap-sidecar
bootstrap-sidecar: ## Full pipeline: cluster + MetalLB + Istio (sidecar) + Bookinfo + validate
	scripts/bootstrap.sh sidecar

.PHONY: bootstrap-ambient
bootstrap-ambient: ## Full pipeline: cluster + MetalLB + Istio (ambient) + Bookinfo + validate
	scripts/bootstrap.sh ambient

## --- GitOps alternative: deploy everything through Argo CD -----------------
## Same components as bootstrap-sidecar/bootstrap-ambient (MetalLB, Gateway
## API CRDs, Istio, Bookinfo, addons), but installed BY Argo CD from
## argocd/application.yaml: official charts + the values in argocd/values/.
## The Helm/kubectl targets above are unchanged - use one flow or the other
## on a given cluster. Needs an existing cluster (make create-cluster).
## Argo CD reads the values files from Git, so they must be pushed:
##   ARGOCD_VALUES_REPO_URL (default: this repo on GitHub)
##   ARGOCD_VALUES_REVISION (default: main)

.PHONY: argocd-install
argocd-install: ## Install Argo CD on the cluster (official install manifest, version from configs/versions.env)
	scripts/argocd.sh install

.PHONY: argocd-render
argocd-render: ## Print the rendered Argo CD Applications without applying, e.g. `make argocd-render MODE=sidecar`
	@scripts/argocd.sh render $(MODE)

.PHONY: argocd-deploy-sidecar
argocd-deploy-sidecar: ## Deploy ALL components through Argo CD, Istio in sidecar mode
	scripts/argocd.sh deploy sidecar

.PHONY: argocd-deploy-ambient
argocd-deploy-ambient: ## Deploy ALL components through Argo CD, Istio in ambient mode
	scripts/argocd.sh deploy ambient

.PHONY: argocd-status
argocd-status: ## Show sync/health status of the Argo CD Applications
	scripts/argocd.sh status

.PHONY: argocd-password
argocd-password: ## Print Argo CD's initial admin password and how to open the UI
	scripts/argocd.sh password

.PHONY: argocd-undeploy
argocd-undeploy: ## Remove ALL components deployed through Argo CD (cascade delete of the Applications)
	scripts/argocd.sh undeploy

.PHONY: argocd-uninstall
argocd-uninstall: ## Remove Argo CD itself (run argocd-undeploy first)
	scripts/argocd.sh uninstall

.PHONY: argocd-bootstrap-sidecar
argocd-bootstrap-sidecar: ## Full GitOps pipeline: cluster + Argo CD + all components (sidecar) + validate
	scripts/preflight.sh
	scripts/create-cluster.sh
	scripts/argocd.sh install
	scripts/argocd.sh deploy sidecar
	scripts/validate-cluster.sh sidecar

.PHONY: argocd-bootstrap-ambient
argocd-bootstrap-ambient: ## Full GitOps pipeline: cluster + Argo CD + all components (ambient) + validate
	scripts/preflight.sh
	scripts/create-cluster.sh
	scripts/argocd.sh install
	scripts/argocd.sh deploy ambient
	scripts/validate-cluster.sh ambient
