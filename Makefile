include .env
export $(shell sed 's/=.*//' .env)

MIGRATE_CMD=migrate -path ./db/migrations \
	-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

run:
	go run cmd/main.go

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

create-migration:
	$(MIGRATE_CMD) create -ext sql -dir ./db/migrations -seq $(NAME)

.PHONY: all lint test lint-fix install-tools

# Default target
all: install-tools lint test

# Install required tools
install-tools:
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/mgechev/revive@latest

# Run all linters
lint:
	@echo "Running linters..."
	golangci-lint run ./...
	revive -config revive.toml -formatter friendly ./...

# Run tests
test:
	go test -v -race ./...

# Fix common linting issues
lint-fix:
	@echo "Fixing common issues..."
	gofmt -w .
	golangci-lint run --fix ./...

# Show lint errors in real-time (for development)
lint-watch:
	@echo "Watching for changes..."
	find . -name "*.go" | entr -c make lint

# Run specific linter
revive:
	@echo "Running revive..."
	revive -config revive.toml -formatter friendly ./...

# Help command
help:
	@echo "Available commands:"
	@echo "  make install-tools  - Install required linting tools"
	@echo "  make lint          - Run all linters"
	@echo "  make lint-fix      - Fix common linting issues"
	@echo "  make test          - Run tests"
	@echo "  make lint-watch    - Watch for changes and run lint"
	@echo "  make revive        - Run only revive linter"