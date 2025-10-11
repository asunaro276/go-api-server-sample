.PHONY: help run dev build test test-coverage test-integration test-performance test-all lint fmt vet check ci migrate migrate-reset docker-up docker-down docker-logs install-tools mock-gen quickstart

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

pre-commit:
	vet fmt test

run:
	go run ./cmd/api-server

dev:
	air

build: ## Build application
	mkdir -p bin
	go build -o bin/api-server ./cmd/api-server

# Testing
test: ## Run unit tests
	go test -parallel 4 -v ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests
	go test -v -parallel 2 -timeout 10m ./cmd/api-server/test/integration/...

test-performance: ## Run performance tests
	go test -bench=. -benchmem -benchtime=100x ./cmd/api-server/test/performance/...

test-all: test test-integration test-performance ## Run all tests

# Code quality
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

check: lint vet test ## Run all code quality checks

ci: check test-all ## Run CI pipeline

# Database
migrate: ## Run database migrations
	go run ./cmd/api-server -migrate

migrate-reset: ## Reset database (development only)
	go run ./cmd/api-server -migrate-reset

# Docker
docker-up: ## Start PostgreSQL container
	docker run --name postgres_api \
		-e POSTGRES_DB=api_db \
		-e POSTGRES_USER=api_user \
		-e POSTGRES_PASSWORD=api_password \
		-p 5432:5432 \
		-d postgres:15

docker-down: ## Stop PostgreSQL container
	docker stop postgres_api || true
	docker rm postgres_api || true

docker-logs: ## Show PostgreSQL container logs
	docker logs postgres_api -f

# Development tools
install-tools: ## Install development tools
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

mock-gen: ## Generate mocks
	go generate ./...

# Quick setup
quickstart: docker-up migrate ## Quick environment setup
	@echo "Environment setup completed!"
	@echo "Run 'make dev' to start the development server"
