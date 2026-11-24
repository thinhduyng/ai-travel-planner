# One command per check. CI runs exactly these targets (.github/workflows/ci.yml).
SHELL := /bin/bash
API   := api
BIN   := $(CURDIR)/bin
GOLANGCI_LINT_VERSION := v2.15.0

.PHONY: tools lint fmt fmt-check typecheck generate generate-check test contract check run

tools: $(BIN)/golangci-lint
$(BIN)/golangci-lint:
	GOBIN=$(BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

## Lint
lint: tools
	cd $(API) && $(BIN)/golangci-lint run ./...

## Format (writes files)
fmt:
	cd $(API) && gofmt -w .

## Format check
fmt-check:
	@cd $(API) && files="$$(gofmt -l .)" && if [ -n "$$files" ]; then echo "Not gofmt'd:"; echo "$$files"; exit 1; fi

## Type check (compile + vet)
typecheck:
	cd $(API) && go vet ./...

## Regenerate code from the OpenAPI spec
generate:
	cd $(API) && go generate ./...

## Fail if generated code is stale (regenerates, then checks git sees no change)
generate-check: generate
	@if [ -n "$$(git status --porcelain -- $(API))" ]; then echo "Generated code is out of date: run 'make generate' and commit"; git status --porcelain -- $(API); exit 1; fi

## Unit tests
test:
	cd $(API) && go test -race $$(go list ./... | grep -v /contract)

## Contract tests (over real HTTP)
contract:
	cd $(API) && go test -race -count=1 -v ./contract/...

## Everything CI runs (generate-check needs a clean git tree, so it runs in CI only)
check: lint fmt-check typecheck test contract

run:
	cd $(API) && go run ./cmd/server
