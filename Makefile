GO ?= go
MIGRATE ?= migrate
SQLC ?= sqlc
DATABASE_URL ?= postgres://genieacs:genieacs@localhost:5432/genieacs_platform?sslmode=disable

.PHONY: all build test tidy run-manager run-agent sqlc migrate-up migrate-down migrate-create

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

sqlc:
	$(SQLC) generate

migrate-up:
	$(MIGRATE) -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	$(MIGRATE) -path db/migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	$(MIGRATE) create -ext sql -dir db/migrations -seq $(name)
