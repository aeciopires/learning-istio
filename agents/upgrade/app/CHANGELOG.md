# Changelog

All notable changes to the `istio-upgrade` CLI are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/). The
version below must always match [VERSION](VERSION).

## [0.4.0] - 2026-11-13

### Added

- `docs update`: the later walkthroughs `day3.md`, `day4.md`, `day5.md` and
  `day6.md` (multicluster, observability, EnvoyFilter/compression, Envoy
  Gateway) now get the same `ISTIO_RELEASE` / `VERSION_ISTIO` /
  `GATEWAY_API_VERSION` / `istio.io/vX.Y/` rewrites as `day1.md`,
  `day1-sidecar.md` and `day2.md`. Unlike those three, they are optional:
  a missing file is skipped instead of failing the run
  (`docsupdater.OptionalDayFiles`).

## [0.3.0] - 2026-09-01

### Fixed

- `docs update`: `UpdateREADME` had no rule for the `- istioctl: X.Y.Z`
  bullet in README.md's "Configurações do Ambiente" section, so it was
  silently left stale (still pinned to the old Istio patch) while every
  other istioctl pin (`.tool-versions`, `INSTALL_REQUIREMENTS.md`) was
  correctly kept in sync with the resolved Istio version. Found while
  bumping this repo to Istio 1.31.0.
- `scripts/create-cluster.sh` (and its mirrored example in
  `INSTALL_REQUIREMENTS.md`): the control-plane node's
  `kubeadmConfigPatches` set `kubeletExtraArgs: [{name: nodeapp, value:
  loadbalancer}]`, which kubelet rejects outright (`unknown flag:
  --nodeapp`) - `node-labels` is the actual kubelet flag that accepts a
  `key=value` label pair, so the node never finished starting and `kind
  create cluster` hung indefinitely at "Starting control-plane". Now
  `{name: node-labels, value: "nodeapp=loadbalancer"}`. Found by actually
  running `make bootstrap-ambient` end to end against the Istio 1.31.0 /
  kind 0.33.0 pins - this bug predates that bump and would have hit
  every prior pin too, it just hadn't been exercised.

### Changed

- `Dockerfile`: bumped the builder base image from `golang:1.26-alpine`
  to `golang:1.27-alpine` and the runtime base from `alpine:3.22` to
  `alpine:3.24.1`. `go.mod`'s `go` directive now pins `1.27.1` exactly
  (was `1.26.5`); since Docker Hub's `golang` image hasn't published a
  matching `1.27.1-alpine` tag yet (only `1.27.0` at the time of this
  change), the Dockerfile now sets `ENV GOTOOLCHAIN=auto` so the build
  fetches the pinned 1.27.1 toolchain itself instead of failing with
  "go.mod requires go >= 1.27.1 (running go 1.27.0; GOTOOLCHAIN=local)" -
  the default `GOTOOLCHAIN=local` in Docker's official Go images exists
  precisely to avoid a surprise network fetch mid-build, so this is an
  explicit, deliberate opt-in. No CLI code changes were needed for the
  Go 1.27.1 bump - `go build`/`go vet`/`go test` all pass unchanged.
- `scripts/install-istio.sh`, `day1.md`, `day1-sidecar.md`,
  `OLD/day1-sidecar.md`: Istio's official Helm chart repository moved
  from `https://istio-release.storage.googleapis.com/charts` to
  `https://blob.istio.io/istio-release/charts` - the old URL's index
  still exists but stops at chart version 1.30.4, so `helm upgrade
  --install ... --version 1.31.0` failed with `no chart version found
  for base-1.31.0` even though the 1.31.0 GitHub release itself is out
  (a similar publish lag to the `golang`/Docker Hub one above, just on
  Istio's chart host instead). Confirmed against the live repo before
  switching: the new URL serves 1.31.0 (and every version back through
  1.29.x), the old one doesn't.
- **Replaced `.tool-versions` (asdf format) with `mise.toml` (mise's
  native TOML format) as the single source of truth for the
  `kind`/`helm`/`kubectl`/`istioctl` pins.** `docs update` now reads and
  writes `mise.toml` exclusively - `.tool-versions` has been removed from
  the repo and is no longer read, written, or referenced anywhere (docs,
  scripts, `Makefile`, CLI help text). `ExtractCurrentVersions` and
  `UpdateMiseToml` (renamed from `UpdateToolVersions`) now parse/emit
  TOML (`kind = "X.Y.Z"`) instead of asdf's `kind X.Y.Z`. Confirmed
  locally that `mise` resolves every pin correctly from `mise.toml` alone
  (`mise current`, `mise which <tool>`, `mise install helm kubectl kind
  istioctl`, `make install-tools`) with no `.tool-versions` present. See
  `UpdateMiseToml` in `internal/docsupdater/files.go`.
  - This briefly went through an intermediate state (keeping both files
    in sync) before this same session's follow-up request to drop
    `.tool-versions` entirely - no intermediate release ever shipped, so
    there's a single 0.3.0 entry rather than two.

## [0.2.0] - 2026-08-19

### Added

- `scripts/lib/common.sh`: two shared logging helpers, `helm_release_note`
  and `kube_resource_note`, used by every install script to print an
  `[INFO]` line stating whether the Helm release / Kubernetes resource
  about to be installed already exists, before reconciling it in place.
  Purely informational - the actual install command (`helm upgrade
  --install`, `kubectl apply`) was already idempotent and unchanged.
- `scripts/create-cluster.sh`: detects and removes leftover Docker
  containers from a previous, incomplete kind cluster creation attempt
  (interrupted by Ctrl-C, a host reboot, a port conflict, ...) before
  retrying, instead of colliding with them. Also warns (does not fail)
  if host ports 80/443 are already held by an unrelated container,
  naming that container instead of letting `kind create cluster` fail
  with Docker's cryptic "port is already allocated" error.
- `scripts/install-metallb.sh` / `scripts/install-istio.sh`: `helm repo
  add` now passes `--force-update`, and every `helm upgrade --install`
  (MetalLB, istio-base, istiod, istio-cni, ztunnel, istio-ingress) is now
  preceded by a `helm_release_note` call.
- `scripts/deploy-bookinfo.sh`: logs `[INFO] Namespace '...' already
  exists, skipping creation.` instead of silently no-op'ing.
- `scripts/install-addons.sh`: logs whether each addon's Deployment
  (Grafana/Prometheus/Kiali/Jaeger) already exists before applying its
  manifest.

### Changed

- `INSTALL_REQUIREMENTS.md`, `day1.md`, `day1-sidecar.md`: aligned the
  manual copy-pasteable commands with the same "check first, never fail
  on an already-existing resource" rule now enforced in
  `scripts/*.sh` - `kubectl create namespace` guarded the same way the
  Gateway API CRD step already was, `kind create cluster
  --name kind-multinodes` guarded with a `kind get clusters` check,
  `helm repo add` gained `--force-update`, and the `sudo rm
  /usr/local/bin/helm|kind` / apt sources cleanup commands switched to
  `rm -f` so they don't error when there's nothing to remove.
- `agents/upgrade/README.md` and the `Makefile`: added short notes
  clarifying that every cluster-lifecycle target (`create-cluster`,
  `install-metallb`, `install-istio-*`, `install-addons`, ...) and
  `make install-tools` (mise) are safe to re-run.

## [0.1.0] - 2026-08-11

### Added

- `analyze` command: fetches and parses real Istio (istio.io + GitHub) and
  Envoy (envoyproxy.io) release notes across a version range, supporting
  both minor-to-minor ranges (e.g. `1.10` -> `1.30`) and patch/hotfix
  ranges (e.g. `1.29.0` -> `1.29.3`). Reports breaking changes,
  deprecations, installation/Helm changes, config/API/CRD changes,
  metrics/telemetry changes, component changes, and security/CVE fixes,
  as markdown or JSON.
- `docs update` command: resolves the latest stable version of Istio,
  Envoy, kind, Helm, kubectl, MetalLB and Gateway API (GitHub Releases
  API) and syncs README.md, day1.md, day1-sidecar.md, day2.md,
  INSTALL_REQUIREMENTS.md, `.tool-versions` and `configs/versions.env` -
  dry-run by default, `--apply` to write. Never applies a MAJOR version
  bump automatically (guards against e.g. Helm 3 -> 4) unless
  `--allow-major-bump` is passed.
- `versions latest` command: prints the latest stable version of every
  tracked tool.
- `version` command and `--version`/`--help` flags.
- On-disk HTTP response cache (`--cache-dir`, `--refresh-cache`) shared by
  every network-backed command.
- `scripts/*.sh` + `Makefile` cluster lifecycle automation: create a kind
  cluster, install MetalLB, install Istio in sidecar or ambient mode,
  deploy Bookinfo, and validate the result end to end
  (`make bootstrap-sidecar` / `make bootstrap-ambient` / `make validate` /
  `make destroy-cluster`) - kept deliberately separate from the Go CLI.
- `Dockerfile` + `make docker-build`/`docker-run` for the CLI.
- Unit test suite (`make test`) covering version-range expansion,
  release-note parsing, GitHub API interaction, and doc-rewrite rules,
  all offline via `httptest` and fixture data.
- `make install-tools`: installs the `kind`/`helm`/`kubectl`/`istioctl`
  versions pinned in the repo's `.tool-versions` via `mise install
  <tool>...` (scoped to those tools on purpose - a bare `mise install`
  also installs whatever is pinned in the user's *global* mise config).
- `docs update` now also manages an `istioctl` pin in `.tool-versions`,
  `INSTALL_REQUIREMENTS.md`, and `configs/versions.env`
  (`ISTIOCTL_VERSION`) - always mirroring the resolved Istio version
  (istioctl ships from the same `istio/istio` release), not a separate
  GitHub lookup.
- `day1-sidecar.md` (root of the repo): a sidecar-mode counterpart to
  `day1.md`, restored/updated from the stale `OLD/day1-sidecar.md`.
  `docs update` now manages version pins in it exactly like `day1.md`/
  `day2.md`.
- Cross-compilation: `make build-linux-amd64`, `build-linux-arm64`,
  `build-darwin-amd64`, `build-darwin-arm64`, and `build-all` - the CLI
  is pure Go (no cgo), so it cross-builds for **macOS (Intel and Apple
  Silicon) and Linux (amd64/arm64)** from any host with no extra
  toolchain. The Docker image stays Linux-only on purpose (there's no
  such thing as a "macOS Docker image" - Docker Desktop on macOS runs
  Linux containers via a VM).
- macOS support for the cluster automation scripts: `scripts/preflight.sh`
  detects `uname -s == Darwin` and skips the Linux-only `fs.inotify`
  sysctl check (informational message instead); `scripts/create-cluster.sh`
  now uses a `mktemp` invocation portable between GNU (Linux) and BSD
  (macOS) mktemp. `INSTALL_REQUIREMENTS.md` gained macOS
  (Homebrew/Docker Desktop) instructions alongside the existing
  Ubuntu/Debian ones.
- `scripts/install-addons.sh` + `make install-addons`: installs the Istio
  observability addons (Grafana, Prometheus, Kiali, Jaeger) from
  `day1.md`/`day1-sidecar.md`'s manifests, waiting for each rollout.
  Grafana/Prometheus/Kiali/Jaeger were added to
  `internal/toolversions.Catalog` and `configs/versions.env`
  (`GRAFANA_VERSION`, etc.) as **informational-only** references (the
  addons install from Istio's own bundled manifests, tracking
  `ISTIO_RELEASE` - there's no independent version to select there).
- `scripts/bootstrap.sh` now also runs Bookinfo and the addons by
  default, controllable via `INSTALL_BOOKINFO`/`INSTALL_ADDONS` env vars
  (default `true`) - e.g. `INSTALL_ADDONS=false make bootstrap-sidecar`.
  Each step keeps its own standalone `make` target too.


### Changed

- Moved every CLI-specific file (`main.go`, `cmd/`, `internal/`, `go.mod`,
  `go.sum`, `VERSION`, `Dockerfile`, `.dockerignore`, `testdata/`,
  `examples/`, `CHANGELOG.md`) into `agents/upgrade/app/`, to make the
  split between "CLI code" and "cluster automation" (`Makefile`,
  `scripts/*.sh`, `configs/versions.env`, which stay at
  `agents/upgrade/`) explicit in the directory layout, not just in
  documentation. The `Makefile` remains the single entry point (now uses
  `go build/test/vet/run/install -C app`); the compiled binary still
  lands in `agents/upgrade/bin/`. See `plans/plan2.md`.
- `docs update` now rewrites `INSTALL_REQUIREMENTS.md`'s Helm/kubectl/kind
  sections to match [mise](https://mise.jdx.dev) commands instead of
  `asdf` (the repo switched version managers; `.tool-versions` itself is
  unchanged, mise reads/writes that same format natively).
- Fixed `REPO_ROOT` in the `Makefile` (`docs-update`/`docs-update-apply`
  targets): it resolved one directory short of the actual repo root and
  had never been exercised until now.
- Fixed the MetalLB `--version` rewrite rule in `docs update`: the doc
  moved the pin into an `export HELM_CHART_VERSION="..."` line instead of
  a literal `--version X.Y.Z`, which silently stopped matching the old
  regex (found and fixed by testing against the real file, not just the
  unit test fixture).
- `make install-tools`: now also installs `istioctl` (added to
  `.tool-versions` as a required tool, mainly for troubleshooting the
  mesh - `istioctl proxy-status`, `istioctl analyze`, etc.).
