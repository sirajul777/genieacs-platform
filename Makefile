GO ?= go

.PHONY: all build test tidy run-manager run-agent

all: tidy test build

build:
	$(GO) build ./...

test:
	$(GO) test ./...

tidy:
	$(GO) mod tidy

run-manager:
	$(GO) run ./cmd/manager

run-agent:
	$(GO) run ./cmd/agent
