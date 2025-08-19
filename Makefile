.PHONY: all fmt lint test build ci

GO  ?= go
PKG ?= ./...

all: test

fmt:
	$(GO) fmt $(PKG)

# Requires golangci-lint installed locally (brew install golangci-lint or see docs)
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed"; exit 1; }
	golangci-lint run

test:
	$(GO) test -race -cover $(PKG)

build:
	$(GO) build .

# What CI runs locally
ci: fmt lint test build
