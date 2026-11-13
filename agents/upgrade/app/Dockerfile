# syntax=docker/dockerfile:1

# istio-upgrade CLI image.
#
# This image packages ONLY the Go CLI (analyze / docs update / versions
# latest) - it does NOT include kind/helm/kubectl/docker. Cluster
# lifecycle (create/install/deploy/validate/destroy) is intentionally a
# host-side Makefile+scripts concern (see ../README.md and
# scripts/*.sh), since it needs to drive a kind cluster via the host's
# Docker daemon - something that doesn't containerize cleanly. Run
# `make bootstrap-sidecar` / `make bootstrap-ambient` on the host instead.

FROM golang:1.27-alpine AS builder

# go.mod pins the exact toolchain patch (go 1.27.4); the golang:1.27-alpine
# base image ships whatever 1.27.x point release Docker Hub currently
# publishes (its official images can lag a Go patch release by a bit) and
# defaults GOTOOLCHAIN=local, which would hard-fail on an older patch
# instead of fetching 1.27.4 - see https://go.dev/doc/toolchain.
ENV GOTOOLCHAIN=auto

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION_OVERRIDE=""
RUN GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) && \
    CGO_ENABLED=0 go build \
      -ldflags "-s -w -X main.gitCommit=${GIT_COMMIT}" \
      -o /out/istio-upgrade .

FROM alpine:3.24.1

RUN apk add --no-cache ca-certificates && \
    addgroup -S istio-upgrade && adduser -S istio-upgrade -G istio-upgrade

COPY --from=builder /out/istio-upgrade /usr/local/bin/istio-upgrade

USER istio-upgrade
WORKDIR /home/istio-upgrade

ENTRYPOINT ["istio-upgrade"]
CMD ["--help"]
