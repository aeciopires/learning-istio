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
- `day3.md` (multicluster / single mesh, sidecar and ambient) / `day4.md`
  (metrics export to Prometheus/Grafana and Datadog) / `day5.md`
  (EnvoyFilter, Envoy filters, compression) / `day6.md` (Envoy Gateway) -
  intermediate/advanced walkthroughs added on top of the course material.
  Every command in them is taken from the official docs linked in each
  section (istio.io, envoyproxy.io, gateway.envoyproxy.io,
  docs.datadoghq.com) - keep it that way: don't add behaviour you can't
  point to an official source for. `docs update` syncs their
  `ISTIO_RELEASE`/`VERSION_ISTIO`/`GATEWAY_API_VERSION` exports and
  `istio.io/vX.Y/` links too (they're optional for the CLI: a missing
  one is skipped). `day6.md`'s `ENVOY_GATEWAY_VERSION` is **not**
  CLI-managed - bump it by hand against Envoy Gateway's compatibility
  matrix.
- `day7.md` (cross-cluster security: SPIFFE identity, AuthorizationPolicy
  between apps in different clusters, JWT authorization by claims) -
  builds on `day3.md`'s multicluster environment (sidecar option) and
  mints its tokens with the CLI's `jwt` subcommand (`make jwt-*`). Same
  "official sources only" rule and same optional `docs update` sync as
  `day3.md`-`day6.md`.
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
- `.jwt/` (git-ignored, created by `make jwt-keygen`) - the lab's RSA
  private key, public key and JWKS. Never commit it.
- `mise.toml` - version pins for `kind`/`helm`/`kubectl`/`istioctl`/`go`, in
  [mise](https://mise.jdx.dev)'s native TOML format (see
  `INSTALL_REQUIREMENTS.md#mise`). This repo does not use the asdf
  `.tool-versions` format. `go` is there only to build the CLI and is
  **not** rewritten by `docs update`: keep it equal to the `go` directive
  in `agents/upgrade/app/go.mod` when bumping it by hand.
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

Deliberately separate pieces - **don't blur this split**, it was a
correction made explicitly during this tool's development:

1. **The Go CLI** (`istio-upgrade` binary; source lives entirely under
   `agents/upgrade/app/` - `app/main.go` + `app/cmd/` + `app/internal/` -
   so it's obvious what's the CLI vs. what's orchestration):
   analyzes Istio/Envoy release notes across a version range (`analyze`,
   supports both minor ranges like `1.10`→`1.30` and patch/hotfix ranges
   like `1.29.0`→`1.29.3`) and keeps this repo's version pins in sync
   with upstream (`docs update`, `versions latest`). It also has an
   offline `jwt` subcommand (`keygen`/`jwks`/`sign`/`decode`, in
   `app/internal/jwttool`, built on `github.com/golang-jwt/jwt/v5`) that
   plays the token issuer for `day7.md`. **It does not touch Kubernetes,
   Docker, or Istio itself** - `jwt` included: it only writes key files
   and prints tokens.
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
3. **The GitOps alternative** (`make argocd-*` targets +
   `scripts/argocd.sh` + `argocd/`): installs Argo CD on the kind cluster
   and has *it* deploy the same components from
   `argocd/application.yaml` - one multi-document file of Argo CD
   `Application`s, each combining an **official chart** with a values
   file from the local `argocd/values/` directory; Gateway API CRDs,
   Bookinfo and the addons have no official chart, so they are read from
   their upstream Git repositories instead. It was **added next to** the
   Helm/kubectl flow, which must keep working unchanged - never replace
   one with the other. Things to know before touching it:
   - `application.yaml` is a **template**: `argocd.sh` substitutes a
     fixed list of `${PLACEHOLDERS}` (plain bash, no `envsubst` - it
     isn't on macOS by default) from `configs/versions.env`. Adding a
     placeholder means adding it to `render_applications` too; the script
     fails if one is left unresolved. `make argocd-render` shows the
     result without a cluster.
   - Values are read from the **local directory, not from a Git branch**
     (an explicit request): each `# values-from: <file>` line in the
     template is replaced by the content of `argocd/values/<file>` and
     ends up embedded under `helm.valuesObject`. Don't reintroduce
     multi-source `$values`/`ref` pointing at this repo - that made the
     deploy depend on what was pushed to a branch. To apply a values
     change: edit the file and re-run `make argocd-deploy-*`.
   - Ordering is done by the script, not Argo CD: Applications carry a
     `learning-istio/step` label, are applied step by step
     (`kubectl apply -l`), and each step must finish syncing and be
     `Healthy` before the next. `learning-istio/mode: ambient` marks
     ambient-only Applications (ztunnel). Bump `LAST_STEP` in the script
     if you add a step.
   - Strict ARP and the MetalLB address pool stay in bash
     (`metallb_enable_strict_arp`/`metallb_configure_address_pool` in
     `scripts/lib/common.sh`, shared with `install-metallb.sh`): the pool
     range depends on the host's Docker `kind` network, so it can't be a
     static manifest in Git.

Key commands (run from `agents/upgrade/` - the `Makefile` there is the
single entry point and internally does `go build -C app`/`go test -C
app`/etc., so you never need to `cd app` yourself):

```bash
make build && ./bin/istio-upgrade --help   # binary lands in agents/upgrade/bin/
make build-all                # cross-compile linux/darwin x amd64/arm64 into bin/
make test                    # unit tests, offline
make lint                    # gofmt -l + go vet + golangci-lint (auto-built into bin/, see below)
make docs-update-apply       # sync version pins repo-wide from GitHub
make bootstrap-sidecar       # real kind cluster + Istio sidecar + Bookinfo + addons + validation
make bootstrap-ambient       # same, ambient mode
make install-addons          # just Grafana/Prometheus/Kiali/Jaeger, standalone
make validate MODE=sidecar   # re-run just the health checks
make destroy-cluster
make argocd-install          # Argo CD on an existing cluster (make create-cluster first)
make argocd-deploy-ambient   # all components through Argo CD (or argocd-deploy-sidecar)
make argocd-render MODE=sidecar  # print the rendered Applications, needs no cluster
make argocd-undeploy         # cascade-delete everything Argo CD deployed
make jwt-keygen              # RSA key pair + JWKS into <repo>/.jwt (git-ignored)
make jwt-sign JWT_SUBJECT=curl JWT_CLAIMS="--claim app=curl --claim-list groups=readers"
make jwt-decode TOKEN=...    # decode AND verify against .jwt/public.pem
make install-go              # the Go toolchain pinned in mise.toml, via mise
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
- **A MAJOR version bump is never auto-applied** by `docs update` - it's
  reported as a warning and the existing pin is kept, since it routinely
  breaks the exact commands this repo documents. Pass
  `--allow-major-bump` to override deliberately. (That is how Helm went
  3 → 4 here: a deliberate, reviewed bump.)
- **golangci-lint must be built with a Go at least as new as
  `app/go.mod`'s `go` directive**, or it refuses to run ("the Go language
  version used to build golangci-lint is lower than the targeted Go
  version"). A prebuilt binary, and even a plain `go install ...@version`
  (which follows golangci-lint's *own* toolchain line), can be too old.
  `make install-lint` therefore builds the version pinned in the
  Makefile (`GOLANGCI_LINT_VERSION`) with `GOTOOLCHAIN=go<go.mod
  version>` into `bin/`, and `make lint` rebuilds it automatically when
  either pin changes. Bump the lint pin by hand - `docs update` doesn't.
- **This repo is on Helm 4.** Every flag the docs/scripts use (`--wait`,
  `--timeout`, `--create-namespace`, `--debug`, `--force-update`) still
  exists; what changed is documented in `INSTALL_REQUIREMENTS.md#helm`
  (server-side apply by default for *new* releases, `--atomic` →
  `--rollback-on-failure`, `--force` → `--force-replace`). Don't
  reintroduce Helm 3-only advice, and check Helm's version-skew table
  before bumping the kind node's Kubernetes version (Helm 4.3.x supports
  Kubernetes 1.34-1.37).
- **Istio via Argo CD needs `base.validationFailurePolicy: Fail`** (in
  `argocd/values/istio-base.yaml` and `istiod-*.yaml`). Argo CD renders
  with `helm template` and applies server-side; without it the chart and
  istiod fight over the `ValidatingWebhookConfiguration`'s
  `failurePolicy` (Istio's own Helm install docs call this out). The
  `ignoreDifferences` on webhook `caBundle` fields are there for the same
  class of reason - don't drop them as "noise".
- **`RequestAuthentication` alone does not require a token** - it only
  rejects invalid ones. `day7.md` relies on this being explained
  correctly; a "require JWT" example needs an `AuthorizationPolicy` with
  `requestPrincipals` (or `jwtRules[].presence: REQUIRED`, Istio 1.32+).
- **In ambient mode, L7 policy attributes (JWT, methods, paths) are only
  enforced by a waypoint**, attached with `targetRefs`. A selector-based
  policy with L7 attributes is skipped by ztunnel: an `ALLOW` then denies
  everything and an L7-only `DENY` silently protects nothing. Don't write
  ambient JWT examples with a `selector`.
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
- **`EnvoyFilter` is not supported on ambient waypoints** (stated in
  Istio's ambient L7-features and migration docs) and doesn't apply to
  `ztunnel` at all (it isn't Envoy). Don't write ambient examples that
  patch a waypoint with `EnvoyFilter` - `day5.md` points to
  `TrafficExtension` (alpha) instead.
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
  `day1-sidecar.md`, `day2.md`, `day3.md`-`day7.md`, `INSTALL_REQUIREMENTS.md` and
  `mise.toml` by design (each file is meant to be readable/runnable
  standalone). Keep them consistent - prefer running `agents/upgrade`'s
  `docs update` over hand-editing when bumping tool versions.
- This repo has no CI configured; `agents/upgrade/Makefile`'s
  `test`/`lint` targets are the closest thing - run them after touching
  Go code there. After touching `argocd/` or `scripts/argocd.sh`, at
  least run `make argocd-render MODE=sidecar` and `MODE=ambient` (no
  cluster needed).
- Commits: per `CONTRIBUTING.md`, prefer one commit per file/topic
  changed, in a feature branch off `main`, opened as a PR - not direct
  pushes to `main`.
- Don't commit or push automatically; the user in this repo reviews
  diffs (including generated ones, e.g. from `docs update --apply`)
  before asking for a commit.
