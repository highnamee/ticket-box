.PHONY: help setup run build test test-coverage lint fmt vet tidy clean migrate swagger

# Default binary name
BIN_DIR = bin
BIN_NAME = api

help: ## Display available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

setup: ## Install required development tools (swag, golangci-lint) and dependencies
	@echo "📦 Downloading Go dependencies..."
	go mod download
	@echo "🛠️ Installing developer CLI tools..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin; \
	go install github.com/swaggo/swag/cmd/swag@v1.16.4; \
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "📖 Generating initial Swagger documentation..."
	@$(MAKE) swagger
	@echo "✅ Setup completed successfully!"

run: ## Run application locally
	go run ./cmd/api/main.go

migrate: ## Run database migrations
	go run ./cmd/migrate/main.go

swagger: ## Generate OpenAPI / Swagger documentation
	@export PATH=$$PATH:$$(go env GOPATH)/bin; \
	which swag > /dev/null || (echo "Installing swag CLI..." && go install github.com/swaggo/swag/cmd/swag@v1.16.4); \
	swag init -g cmd/api/main.go -o docs

build: ## Build binary executable
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) ./cmd/api/main.go

test: ## Run unit and integration tests with race detection
	go test -v -race ./internal/...

test-coverage: ## Run tests with code coverage report
	go test -race -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run golangci-lint
	@export PATH=$$PATH:$$(go env GOPATH)/bin; \
	which golangci-lint > /dev/null || (echo "golangci-lint is not installed. Run 'make setup' or install via: brew install golangci-lint (macOS)" && exit 1); \
	golangci-lint run ./...

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet analysis
	go vet ./...

tidy: ## Clean up and download Go module dependencies
	go mod tidy

clean: ## Clean build artifacts and coverage files
	rm -rf $(BIN_DIR) coverage.out coverage.html
