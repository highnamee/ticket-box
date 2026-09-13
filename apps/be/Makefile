.PHONY: help run build test test-coverage lint fmt vet tidy clean migrate

# Default binary name
BIN_DIR = bin
BIN_NAME = api

help: ## Display available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

run: ## Run application locally
	go run ./cmd/api/main.go

migrate: ## Run database migrations
	go run ./cmd/migrate/main.go

build: ## Build binary executable
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) ./cmd/api/main.go

test: ## Run unit tests with race detection
	go test -v -race ./...

test-coverage: ## Run tests with code coverage report
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "golangci-lint is not installed. Install via: brew install golangci-lint (macOS) or go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet analysis
	go vet ./...

tidy: ## Clean up and download Go module dependencies
	go mod tidy

clean: ## Clean build artifacts and coverage files
	rm -rf $(BIN_DIR) coverage.out coverage.html
