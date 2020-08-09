BIN_NAME = "kill-switch"

ROOT = $(shell pwd)
GO ?= go
OS = $(shell uname -s | tr A-Z a-z)
export GOBIN = ${ROOT}/bin

LINT = ${GOBIN}/golangci-lint
LINT_DOWNLOAD = curl --progress-bar -SfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.28.1
VERSION_TAG = $(shell git describe --tags --abbrev=0 --always)
VERSION_COMMIT = $(shell git rev-parse --short HEAD)
VERSION_DATE = $(shell git show -s --format=%cI HEAD)
VERSION = -X main.versionTag=$(VERSION_TAG) -X main.versionCommit=$(VERSION_COMMIT) -X main.versionDate=$(VERSION_DATE)
PATH := $(PATH):$(GOBIN)

.PHONY: help
help: ## Display this help message
	@ cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build development binary file
	@ $(GO) build -ldflags '$(VERSION)' -o ./bin/${BIN_NAME} ./cmd/...

.PHONY: mod
mod: ## Get dependency packages
	@ $(GO) mod tidy

.PHONY: test
test:create-env ## Run unit tests
	echo $(TPARSE)
	@ test -e $(TPARSE) || $(TPARSE_DOWNLOAD)
	@ $(GO) test -failfast -count=1 ./... -json -cover | $(TPARSE) -all -smallscreen

.PHONY: race
race:create-env ## Run data race detector
	@ test -e $(TPARSE) || $(TPARSE_DOWNLOAD)
	@ $(GO) test -short -race ./... -json -cover | $(TPARSE) -all -smallscreen

.PHONY: coverage
coverage:create-env ## check coverage test code of sample https://penkovski.com/post/gitlab-golang-test-coverage/
	@ $(GO) test ./... -coverprofile=coverage.out
	@ $(GO) tool cover -func=coverage.out
	@ $(GO) tool cover -html=coverage.out -o coverage.html;

.PHONY: lint
lint: ## Lint the files
	@ test -e $(LINT) || $(LINT_DOWNLOAD)
	@ $(LINT) version
	@ $(LINT) --timeout 10m run


