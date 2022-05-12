.DEFAULT_GOAL := help
LINT_ARGS ?= --timeout 5m
FRMD_CMD ?=
SHELL = /bin/bash
TAILWINDCSS_ARGS ?= 
GORELEASER_VERSION ?= v1.8.3
GORELEASER_ARGS ?= --auto-snapshot --rm-dist
GITCHLOG_ARGS ?=
SHELL := /bin/bash

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

watch: deps ## Watching updated files - live reload
	( set -o allexport && source .env && set +o allexport && go run -mod=readonly github.com/cortesi/modd/cmd/modd@latest )

test: deps ## Executing tests
	( set -o allexport && source .env && set +o allexport && go test -v -race -count=1 $(GOTEST_ARGS) ./... )

lint: ## Lint sources code
	golangci-lint run --enable-all $(LINT_ARGS)

build: build-frmd ## Build artefacts

build-frmd: deps tailwind ## Build executable
	CGO_ENABLED=0 go build -v -o ./bin/frmd ./cmd/frmd

.PHONY: tailwind
tailwind: deps
	npx tailwindcss -i ./internal/server/assets/src/main.css -o ./internal/server/assets/dist/main.css $(TAILWINDCSS_ARGS)

internal/server/assets/dist/main.css: tailwind

.env:
	cp .env.dist .env

.PHONY: deps
deps: .env node_modules

node_modules:
	npm ci

.PHONY: release
release: deps
	( set -o allexport && source .env && set +o allexport && VERSION=$(GORELEASER_VERSION) curl -sfL https://goreleaser.com/static/run | bash /dev/stdin $(GORELEASER_ARGS) )

.PHONY: changelog
changelog:
	go run -mod=readonly github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.1 $(GITCHLOG_ARGS)

install-git-hooks:
	git config core.hooksPath .githooks