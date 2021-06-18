GO := go
GO_BUILD := $(GO) build
export GOBIN ?= $(shell pwd)/bin
export GO111MODULE := on
# export GOPRIVATE := github.com/ebi-yade/go-json-sandbox

rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))
SRC := $(call rwildcard,,*.go) go.mod go.sum

MAKE2HELP := $(GOBIN)/make2help
GOLANGCI_LINT := $(GOBIN)/golangci-lint

$(GOBIN)/%:
	@scripts/tools.sh $(notdir $@)

.DEFAULT_GOAL := help
.PHONY: fmt lint test tools clean help

## Format files
fmt:
	$(GO) fmt ./...

## Run linters
lint: $(GOLANGCI_LINT)
	$(GO) vet ./...
	$(GOLANGCI_LINT) run

## Run tests concurrently
test:
	$(GO) test -race ./...

## Install tools
tools:
	@scripts/tools.sh

## Clean up artifacts
clean:
	$(GO) clean

## Show help via make2help
help: $(MAKE2HELP)
	@$(MAKE2HELP)
