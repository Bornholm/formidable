.DEFAULT_GOAL := help
LINT_ARGS ?= --timeout 5m
FRMD_CMD ?=
SHELL = /bin/bash
TAILWINDCSS_ARGS ?= 
GORELEASER_VERSION ?= v1.8.3
GORELEASER_ARGS ?= --auto-snapshot --rm-dist
GITCHLOG_ARGS ?=
SHELL := /bin/bash
RUN_INSTALL_TESTS ?= yes

FORMIDABLE_VERSION ?= 

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

watch: deps ## Watching updated files - live reload
	( set -o allexport && source .env && set +o allexport && go run -mod=readonly github.com/cortesi/modd/cmd/modd@latest )

.PHONY: test
test: test-go ## Executing tests

ifeq ($(RUN_INSTALL_TESTS), yes)
test: test-install-script
endif

test-go: deps
	( set -o allexport && source .env && set +o allexport && go test -v -race -count=1 $(GOTEST_ARGS) ./... )

test-install-script: tools/bin/bash_unit
	tools/bin/bash_unit ./misc/script/test_install.sh

tools/bin/bash_unit:
	mkdir -p tools/bin
	cd tools/bin && bash <(curl -s https://raw.githubusercontent.com/pgrange/bash_unit/master/install.sh)	

lint: ## Lint sources code
	golangci-lint run --enable-all $(LINT_ARGS)

build: build-frmd ## Build artefacts

build-frmd: deps tailwind ## Build executable
	CGO_ENABLED=0 go build \
		-v \
		-ldflags "\
			-X 'main.GitRef=$(shell git rev-parse --short HEAD)' \
			-X 'main.ProjectVersion=$(shell git describe --always)' \
			-X 'main.BuildDate=$(shell date --utc --rfc-3339=seconds)' \
		" \
		-o ./bin/frmd \
		./cmd/frmd

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

.PHONY: start-release
start-release:
	if [ -z "$(FORMIDABLE_VERSION)" ]; then echo "You must define environment variable FORMIDABLE_VERSION"; exit 1; fi
	
	git flow release start $(FORMIDABLE_VERSION)

	# Update package.json version
	jq '.version = "$(FORMIDABLE_VERSION)"' package.json | sponge package.json
	git add  package.json
	git commit -m "chore: bump to version $(FORMIDABLE_VERSION)"

	# Generate updated changelog
	$(MAKE) GITCHLOG_ARGS='--next-tag $(FORMIDABLE_VERSION)' changelog
	git add CHANGELOG.md
	git commit -m "chore: update changelog for version $(FORMIDABLE_VERSION)"

	echo "Commit you additional modifications then execute 'make finish-release'"

.PHONY: finish-release
finish-release:
	git flow release finish -m "v$(FORMIDABLE_VERSION)"
	git push --all
	git push --tags

.PHONY: changelog
changelog:
	go run -mod=readonly github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.1 $(GITCHLOG_ARGS) > CHANGELOG.md

install-git-hooks:
	git config core.hooksPath .githooks