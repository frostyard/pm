SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c

# Install dev tools into a repo-local bin dir so contributors don't need sudo.
BIN_DIR := $(CURDIR)/bin
export PATH := $(BIN_DIR):$(PATH)

GO ?= go
GOLANGCI_LINT_VERSION ?= latest
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

GOFILES := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_\/-]+:.*##/ {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: tools
tools: $(GOLANGCI_LINT) ## Install required external tools
	@echo "Tools installed into $(BIN_DIR)"

$(GOLANGCI_LINT): ## Install golangci-lint
	@mkdir -p "$(BIN_DIR)"
	@if [ "$(GOLANGCI_LINT_VERSION)" = "latest" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(BIN_DIR)"; \
	else \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(BIN_DIR)" "$(GOLANGCI_LINT_VERSION)"; \
	fi
	@test -x "$(GOLANGCI_LINT)" || { echo "ERROR: golangci-lint install failed"; exit 1; }

.PHONY: ensure-go
ensure-go: ## Verify Go toolchain is available
	@command -v $(GO) >/dev/null 2>&1 || { echo "ERROR: 'go' not found in PATH"; exit 1; }

.PHONY: ensure-mod
ensure-mod: ## Verify go.mod exists
	@test -f go.mod || { \
		echo "ERROR: go.mod not found."; \
		echo "Run: $(GO) mod init <module-path>"; \
		echo "Then rerun: make tidy"; \
		exit 1; \
	}

.PHONY: fmt
fmt: ensure-go ## Format Go code (gofmt)
	@if [ -n "$(GOFILES)" ]; then \
		echo "gofmt"; \
		gofmt -w $(GOFILES); \
	else \
		echo "No Go files to format."; \
	fi

.PHONY: tidy
tidy: ensure-go ensure-mod ## Run go mod tidy
	$(GO) mod tidy

.PHONY: vet
vet: ensure-go ensure-mod ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ensure-go ensure-mod tools ## Run golangci-lint
	$(GOLANGCI_LINT) run ./...

.PHONY: test
test: ensure-go ensure-mod ## Run unit tests
	$(GO) test ./...

.PHONY: build
build: ensure-go ensure-mod ## Compile all packages (no tests)
	$(GO) test -run=^$ ./...

.PHONY: check
check: fmt lint test ## Format, lint, and test (main local gate)

.PHONY: ci
ci: tools check ## CI entrypoint
