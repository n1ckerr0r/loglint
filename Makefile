GO ?= go
GOLANGCI_LINT ?= golangci-lint
GO_PACKAGES := ./...

.PHONY: all
all: verify

.PHONY: fmt
fmt:
	$(GO) fmt $(GO_PACKAGES)

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: test
test:
	$(GO) test $(GO_PACKAGES)

.PHONY: vet
vet:
	$(GO) vet $(GO_PACKAGES)

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run

.PHONY: vulncheck
vulncheck:
	govulncheck $(GO_PACKAGES)

.PHONY: verify
verify: fmt tidy vet test

.PHONY: build-plugin
build-plugin:
	golangci-lint custom
