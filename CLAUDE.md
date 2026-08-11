# CLAUDE.md

Guidance for Claude Code (and other AI agents) working in this repository.

## What this repo is

`learning-istio` is Aécio Pires' **pt-BR study/teaching material** for
Istio (sidecar and ambient modes) and Envoy, built around a local
[kind](https://kind.sigs.k8s.io) cluster. It originated from the LinuxTips
"Descomplicando Istio" course. Prose, headings and comments in the
top-level `.md` files are **Portuguese**; keep new prose in the same
files in Portuguese unless told otherwise. Code and comments inside
`agents/upgrade/` (the Go CLI, see below) are in **English**, matching Go
community convention - don't "fix" that to Portuguese.

## Repository layout

- `README.md` - entry point; pins the versions used throughout (Docker,
  kind's Kubernetes version, kubectl, Helm, Istio, istioctl).
- `day1.md` (ambient mode) / `day1-sidecar.md` (sidecar mode) / `day2.md`
  - copy-pasteable walkthroughs: create the kind cluster → Helm → MetalLB
  → Istio → Bookinfo → addons (Grafana/Prometheus/Kiali/Jaeger) → traffic
  management tasks. These are the canonical source for "how do the manual
  commands go" - `agents/upgrade/scripts/*.sh` mirror them for automation
  (`docs update` keeps `ISTIO_RELEASE`/`VERSION_ISTIO`/`GATEWAY_API_VERSION`
  in sync across all three files identically).
- `INSTALL_REQUIREMENTS.md` - OS packages, [mise](https://mise.jdx.dev)
  (the asdf-compatible version manager used here - not asdf itself),
  Docker, Helm, kubectl, istioctl (required, mainly for troubleshooting -
  it mirrors the pinned Istio version, not a separately-tracked upstream),
  kind (incl. the kind cluster config and the `kindest/node` image
  tag+sha256 digest), MetalLB installation. This is the source of truth
  for exact tool versions/pins.
- `UNINSTALL_ISTIO.md` - teardown instructions.
- `ISTIO_ENVOY_GUIA_DE_ESTUDOS.md` - a long-form study guide (glossary,
  mermaid diagrams, Envoy xDS/request-lifecycle deep dive, troubleshooting
  exercises). Treat this as curated prose - don't run automated rewrites
  against it.
- `.tool-versions` - version pins for `kind`/`helm`/`kubectl`/`istioctl`,
  in the asdf `.tool-versions` format (mise reads/writes this format
  natively - see `INSTALL_REQUIREMENTS.md#mise`).
- `files/`, `images/` - sample manifests (Bookinfo, mirroring) and
  screenshots referenced by the walkthroughs.
- `OLD/` - **archived/deprecated** content (old auth-policy, egress,
  retry-policy examples, an older/stale `day1-sidecar.md` - not the same
  file as the root `day1-sidecar.md`, which is current and CLI-managed -
  etc.). Don't update it as part of routine version-sync work; it's kept
  for historical reference only.
- `agents/upgrade/` - a self-contained Go CLI + shell automation. See its
  own [agents/upgrade/README.md](agents/upgrade/README.md) for full
  details; summary below.

## `agents/upgrade/` - the istio-upgrade CLI

Two deliberately separate pieces - **don't blur this split**, it was a
correction made explicitly during this tool's development:

1. **The Go CLI** (`istio-upgrade` binary; source lives entirely under
   `agents/upgrade/app/` - `app/main.go` + `app/cmd/` + `app/internal/` -
   so it's obvious what's the CLI vs. what's orchestration):
   analyzes Istio/Envoy release notes across a version range (`analyze`,
   supports both minor ranges like `1.10`→`1.30` and patch/hotfix ranges
   like `1.29.0`→`1.29.3`) and keeps this repo's version pins in sync
   with upstream (`docs update`, `versions latest`). **It does not touch
   Kubernetes, Docker, or Istio itself.**
2. **The cluster automation** (`Makefile` targets + `scripts/*.sh`):
   creates the kind cluster, installs MetalLB, installs Istio (sidecar
   *or* ambient), deploys Bookinfo, installs the addons (Grafana,
   Prometheus, Kiali, Jaeger), and validates the result. This is plain
   bash, not the Go CLI, because it drives the host's Docker daemon -
   something that doesn't containerize/CLI-ify cleanly. It reads
   `agents/upgrade/configs/versions.env`, which the CLI's `docs update
   --apply` generates - run that before running `make bootstrap-*`.
   `bootstrap-sidecar`/`bootstrap-ambient` install everything by default;
   `INSTALL_BOOKINFO=false`/`INSTALL_ADDONS=false` env vars skip a step,
   and each step also has its own standalone `make` target
   (`install-addons`, `deploy-bookinfo-{sidecar,ambient}`, etc.) - see
   `scripts/bootstrap.sh`. Works on Linux and macOS (Intel/arm64).

Key commands (run from `agents/upgrade/` - the `Makefile` there is the
single entry point and internally does `go build -C app`/`go test -C
app`/etc., so you never need to `cd app` yourself):

```bash
make build && ./bin/istio-upgrade --help   # binary lands in agents/upgrade/bin/
make build-all                # cross-compile linux/darwin x amd64/arm64 into bin/
make test                    # unit tests, offline
make lint                    # gofmt -l + go vet + golangci-lint if present
make docs-update-apply       # sync version pins repo-wide from GitHub
make bootstrap-sidecar       # real kind cluster + Istio sidecar + Bookinfo + addons + validation
make bootstrap-ambient       # same, ambient mode
make install-addons          # just Grafana/Prometheus/Kiali/Jaeger, standalone
make validate MODE=sidecar   # re-run just the health checks
make destroy-cluster
```

The CLI's own version lives in the single file `agents/upgrade/app/VERSION`
(embedded via `go:embed`) - bump that, not a hardcoded string somewhere,
and update `agents/upgrade/app/CHANGELOG.md` alongside it.

### Non-obvious gotchas learned the hard way

- **Istio's Helm charts have no `profile=sidecar`.** Sidecar mode is the
  chart's *default* profile; only pass `--set profile=ambient` for
  ambient mode, and pass no `profile` flag at all for sidecar mode
  (`scripts/install-istio.sh` handles this - if you're scripting Istio
  Helm installs elsewhere in this repo, replicate that, not
  `profile=sidecar`).
- **MetalLB's GitHub "latest release" tag can be `metallb-chart-X.Y.Z`**,
  not `vX.Y.Z` - MetalLB publishes separate app and chart releases, and
  GitHub's "latest" sometimes points at the chart one.
  `app/internal/docsupdater`'s `normalizeVersionTag` extracts the trailing
  `X.Y.Z` regardless of prefix; don't assume a plain `strings.TrimPrefix(v,
  "v")` is enough when adding new tools to the catalog
  (`app/internal/toolversions.Catalog`).
- **A MAJOR version bump is never auto-applied** by `docs update` (e.g.
  Helm 3 → 4, which is genuinely available upstream) - it's reported as a
  warning and the existing pin is kept, since it routinely breaks the
  exact commands this repo documents. Pass `--allow-major-bump` to
  override deliberately.
- **The `kindest/node` image tag+sha256 digest** in
  `INSTALL_REQUIREMENTS.md` is never rewritten automatically - a digest
  can't be safely guessed/validated without pulling it, so bumping the
  kind-provisioned Kubernetes version is a manual, deliberate edit (get
  the digest from https://github.com/kubernetes-sigs/kind/releases).
- Istio's exact **bundled Envoy version isn't published as a clean
  version number** (it's an `istio/proxy` fork commit SHA in
  `istio.deps`) - `analyze --component all` won't guess it; pass
  `--envoy-from`/`--envoy-to` explicitly to analyze Envoy.
- **`istio-cni` (installed in BOTH modes - it replaces the privileged
  `istio-init` container in sidecar mode, and captures pods into
  `ztunnel` in ambient mode) CrashLoopBackOffs on a default-tuned Linux
  host** with `couldn't initialize inotify: too many open files` - the
  host's `fs.inotify.max_user_instances`/`max_user_watches` are too low
  (already documented as a manual fix in `day1.md`). Verified live during
  this CLI's development, in both modes. `scripts/preflight.sh` now checks this and
  raises the limits automatically via `sudo -n sysctl -w ...` when
  passwordless sudo is available, else prints the exact commands to run.
  **`fs.inotify` is Linux-only** - on macOS `preflight.sh` detects
  `uname -s == Darwin` and skips the host-level check (the relevant
  kernel is Docker Desktop's Linux VM, not the host).
- **macOS/BSD `mktemp` isn't a drop-in for GNU `mktemp`**: `mktemp -t
  NAME-XXXXXX.yaml` (GNU) silently produces the wrong filename on BSD
  mktemp (macOS), which treats the whole `-t` argument as a literal
  prefix instead of a template. Use `mktemp "${TMPDIR:-/tmp}/NAME.XXXXXX"`
  (no `-t`) instead - portable across both - see
  `scripts/create-cluster.sh`.
- **Addon versions in `configs/versions.env`** (`GRAFANA_VERSION`,
  `PROMETHEUS_VERSION`, `KIALI_VERSION`, `JAEGER_VERSION`) are
  **informational only** - `scripts/install-addons.sh` installs them from
  Istio's own bundled sample manifests (tracks `ISTIO_RELEASE`, not an
  independently selectable version), so these fields don't control what
  gets installed; they just report the latest upstream release of each,
  for reference (`internal/toolversions.Catalog`,
  `internal/docsupdater.BuildResolvedVersions`).

## Working conventions

- Version pins are duplicated across `README.md`, `day1.md`,
  `day1-sidecar.md`, `day2.md`, `INSTALL_REQUIREMENTS.md` and
  `.tool-versions` by design (each file is meant to be readable/runnable
  standalone). Keep them consistent - prefer running `agents/upgrade`'s
  `docs update` over hand-editing when bumping tool versions.
- This repo has no CI configured; `agents/upgrade/Makefile`'s
  `test`/`lint` targets are the closest thing - run them after touching
  Go code there.
- Commits: per `CONTRIBUTING.md`, prefer one commit per file/topic
  changed, in a feature branch off `main`, opened as a PR - not direct
  pushes to `main`.
- Don't commit or push automatically; the user in this repo reviews
  diffs (including generated ones, e.g. from `docs update --apply`)
  before asking for a commit.
