.PHONY: help build run test test-coverage test-integration lint fmt clean docker-up docker-down migrate mock-gen

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build commands
build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api-server cmd/api-server/main.go

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o bin/api-server-linux cmd/api-server/main.go

# Run commands
run: ## Run the application
	@echo "Running application..."
	@go run cmd/api-server/main.go

dev: ## Run in development mode with auto-reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

# Test commands
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -tags=integration -v ./api/test/integration/...

test-contract: ## Run contract tests
	@echo "Running contract tests..."
	@go test -v ./api/test/contract/...

test-performance: ## Run performance tests
	@echo "Running performance tests..."
	@go test -v ./api/test/performance/...

test-all: test test-integration test-contract test-performance ## Run all tests

# Quality commands
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Mock generation
mock-gen: ## Generate mocks using mockery
	@echo "Generating mocks..."
	@if command -v mockery > /dev/null; then \
		mockery --all --output=./internal/mocks --case=underscore; \
	else \
		echo "mockery not installed. Install with: go install github.com/vektra/mockery/v2@latest"; \
		exit 1; \
	fi

# Database commands
migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run cmd/api-server/main.go migrate

migrate-reset: ## Reset database (development only)
	@echo "Resetting database..."
	@go run cmd/api-server/main.go migrate-reset

seed: ## Seed database with test data
	@echo "Seeding database..."
	@go run cmd/api-server/main.go seed

# Docker commands
docker-build: ## Build docker image
	@echo "Building docker image..."
	@docker build -t go-api-server-sample .

docker-up: ## Start docker compose services
	@echo "Starting docker services..."
	@docker-compose up -d

docker-down: ## Stop docker compose services
	@echo "Stopping docker services..."
	@docker-compose down

docker-logs: ## Show docker logs
	@docker-compose logs -f

# Development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/cosmtrek/air@latest

# Cleanup commands
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html
	@go clean -cache
	@go clean -modcache

clean-docker: ## Clean docker resources
	@echo "Cleaning docker resources..."
	@docker-compose down -v
	@docker system prune -f

# Check commands
check: lint vet test ## Run all checks (lint, vet, test)

ci: check test-coverage ## Run CI pipeline

# Coverage analysis
coverage: test-coverage ## Generate and open coverage report
	@go tool cover -func=coverage.out
	@echo "Opening coverage report..."
	@if command -v open > /dev/null; then \
		open coverage.html; \
	elif command -v xdg-open > /dev/null; then \
		xdg-open coverage.html; \
	else \
		echo "Coverage report available at: coverage.html"; \
	fi

# Security scan
security: ## Run security scan with gosec
	@echo "Running security scan..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
		exit 1; \
	fi

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

tidy: ## Tidy up dependencies
	@echo "Tidying up dependencies..."
	@go mod tidy

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v godoc > /dev/null; then \
		echo "Starting godoc server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Git hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@cp scripts/hooks/* .git/hooks/ 2>/dev/null || echo "No git hooks found"
	@chmod +x .git/hooks/* 2>/dev/null || true

# Environment setup
env-example: ## Create example environment file
	@echo "Creating .env.example..."
	@echo "# Database Configuration" > .env.example
	@echo "DB_HOST=localhost" >> .env.example
	@echo "DB_PORT=5432" >> .env.example
	@echo "DB_USER=postgres" >> .env.example
	@echo "DB_PASSWORD=password" >> .env.example
	@echo "DB_NAME=go_api_server" >> .env.example
	@echo "DB_SSLMODE=disable" >> .env.example
	@echo "DB_MAX_IDLE_CONNS=10" >> .env.example
	@echo "DB_MAX_OPEN_CONNS=100" >> .env.example
	@echo "DB_CONN_MAX_LIFETIME=1h" >> .env.example
	@echo "" >> .env.example
	@echo "# Server Configuration" >> .env.example
	@echo "SERVER_HOST=localhost" >> .env.example
	@echo "SERVER_PORT=8080" >> .env.example
	@echo "" >> .env.example
	@echo "# Application Configuration" >> .env.example
	@echo "APP_ENV=development" >> .env.example
	@echo "LOG_LEVEL=debug" >> .env.example
	@echo ".env.example created"

# Quick start
quickstart: env-example docker-up migrate ## Setup development environment
	@echo "Development environment ready!"
	@echo "Run 'make run' to start the server"

# Version info
version: ## Show version information
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build date: $(shell date)"