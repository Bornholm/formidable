.DEFAULT_GOAL := help
LINT_ARGS ?= --timeout 5m
FRMD_CMD ?=
SHELL = /bin/bash

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

watch: ## Watching updated files - live reload
	go run -mod=readonly github.com/cortesi/modd/cmd/modd@latest

test: .env ## Executing tests
	( set -o allexport && source .env && set +o allexport && go test -v -race -count=1 $(GOTEST_ARGS) ./... )

lint: ## Lint sources code
	golangci-lint run --enable-all $(LINT_ARGS)

build: build-frmd ## Build artefacts

build-frmd: ## Build executable
	CGO_ENABLED=0 go build -v -o ./bin/frmd ./cmd/frmd

.env:
	cp .env.dist .env

deps:

.PHONY: release
release:
	./misc/script/release