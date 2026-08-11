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

FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION_OVERRIDE=""
RUN GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) && \
    CGO_ENABLED=0 go build \
      -ldflags "-s -w -X main.gitCommit=${GIT_COMMIT}" \
      -o /out/istio-upgrade .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates && \
    addgroup -S istio-upgrade && adduser -S istio-upgrade -G istio-upgrade

COPY --from=builder /out/istio-upgrade /usr/local/bin/istio-upgrade

USER istio-upgrade
WORKDIR /home/istio-upgrade

ENTRYPOINT ["istio-upgrade"]
CMD ["--help"]
