# SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
#
# SPDX-License-Identifier: AGPL-3.0-only

VERSION 0.8
FROM golang:1.25.4
WORKDIR /workspace

# Runs all recipes. Do this before you commit your changes to ensure that nothing broke.
all:
    WAIT
        BUILD +lint
    END
    BUILD +build

clean:
    LOCALLY
    RUN rm -rf build/

deps-go:
    COPY go.mod go.sum ./
    RUN go mod download
    # Output these back in case go mod download changes them.
    SAVE ARTIFACT go.mod
    SAVE ARTIFACT go.sum
    RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2

lint:
    # TODO: Fix Dockerfile lint issues and then reenable.
    #BUILD +lint-dockerfile
    BUILD +lint-go
    BUILD +lint-license-headers

lint-go:
    FROM +deps-go
    COPY . .
    RUN golangci-lint run --verbose --timeout 5m

lint-dockerfile:
    FROM ghcr.io/hadolint/hadolint:latest-alpine
    COPY . .
    RUN find . -type f -name 'Dockerfile*' -exec hadolint --verbose {} +

lint-license-headers:
    FROM fsfe/reuse
    COPY . .
    RUN reuse lint

fix-license-headers:
    LOCALLY
    RUN reuse annotate --copyright="maxhash.io <dev@maxhash.io>" --license="AGPL-3.0-only" --fallback-dot-license --skip-existing --recursive

gen-static:
    FROM rust:latest
    ENV PATH=/root/.cargo/bin:$PATH
    RUN cargo install minhtml
    COPY http/ http/
    RUN minhtml --minify-css --minify-js http/*.html http/*.css
    SAVE ARTIFACT --force http/*.html AS LOCAL http/static/
    SAVE ARTIFACT --force http/*.css AS LOCAL http/static/

build:
    COPY . .
    ENV GOOS=linux
    ENV GOARCH=amd64
    ENV CGO_ENABLED=0
    RUN go build -a -ldflags '-s -w -extldflags "-static"' -o build/dashboard ./cmd/dashboard/main.go
    SAVE ARTIFACT --force build/dashboard AS LOCAL build/dashboard

deploy-dashboard-fly:
    WAIT
        BUILD +build
    END
    LOCALLY
    RUN fly deploy --config infra/dashboard.fly.toml

deploy-ckpassthrough-fly:
    LOCALLY
    RUN fly deploy --config infra/ckpassthrough.fly.toml
