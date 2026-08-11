# Changelog

All notable changes to the `istio-upgrade` CLI are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/). The
version below must always match [VERSION](VERSION).

## [Unreleased]

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

### Added

- `make install-tools`: installs the `kind`/`helm`/`kubectl`/`istioctl`
  versions pinned in the repo's `.tool-versions` via `mise install
  <tool>...` (scoped to those tools on purpose - a bare `mise install`
  also installs whatever is pinned in the user's *global* mise config).
- `docs update` now also manages an `istioctl` pin in `.tool-versions`,
  `INSTALL_REQUIREMENTS.md`, and `configs/versions.env`
  (`ISTIOCTL_VERSION`) - always mirroring the resolved Istio version
  (istioctl ships from the same `istio/istio` release), not a separate
  GitHub lookup.

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
  API) and syncs README.md, day1.md, day2.md, INSTALL_REQUIREMENTS.md,
  `.tool-versions` and `configs/versions.env` - dry-run by default,
  `--apply` to write. Never applies a MAJOR version bump automatically
  (guards against e.g. Helm 3 -> 4) unless `--allow-major-bump` is passed.
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

[Unreleased]: https://github.com/aeciopires/learning-istio/compare/main...HEAD
