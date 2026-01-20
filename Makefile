# =============================================================================
# SubFlow: Enterprise Construction Financial Ledger
# Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
# Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
# =============================================================================

# Build variables
BINARY_NAME=subflow
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Docker variables
DOCKER_IMAGE=subflow
DOCKER_TAG?=latest

# Database
DB_URL?=postgres://subflow:subflow@localhost:5432/subflow?sslmode=disable

.PHONY: help build run test clean docker-build docker-run migrate-up migrate-down lint fmt

# Default target
help: ## Show this help message
	@echo "SubFlow - Enterprise Construction Financial Ledger"
	@echo "=================================================="
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the application binary
	@echo "🏗️  Building $(BINARY_NAME)..."
	@go build $(GOFLAGS) -o bin/$(BINARY_NAME) ./cmd/api

build-linux: ## Cross-compile for Linux
	@echo "🐧 Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/api

run: ## Run the application locally
	@echo "🚀 Starting $(BINARY_NAME)..."
	@go run ./cmd/api

dev: ## Run with hot reload (requires air)
	@air -c .air.toml

# Testing
test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "📊 Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Code quality
lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@gofumpt -l -w .

vet: ## Run go vet
	@go vet ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

# Database migrations
migrate-up: ## Apply database migrations
	@echo "📈 Applying migrations..."
	@migrate -path migrations -database "$(DB_URL)" up

migrate-down: ## Rollback last migration
	@echo "📉 Rolling back migration..."
	@migrate -path migrations -database "$(DB_URL)" down 1

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@echo "📝 Creating migration: $(NAME)"
	@migrate create -ext sql -dir migrations -seq $(NAME)

# Docker targets
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	@echo "🐳 Running container..."
	@docker run -p 3000:3000 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Start all services with Docker Compose
	@echo "🐳 Starting services..."
	@docker-compose up -d

docker-compose-down: ## Stop all Docker Compose services
	@docker-compose down

docker-compose-logs: ## View Docker Compose logs
	@docker-compose logs -f

# Cleanup
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Development setup
setup: ## Install development dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install mvdan.cc/gofumpt@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅ Setup complete!"

# Generate
generate: ## Run go generate
	@echo "⚙️  Generating code..."
	@go generate ./...

# SQL (for sqlc)
sqlc: ## Generate SQL type-safe code
	@echo "📊 Generating sqlc code..."
	@sqlc generate

# Documentation
docs: ## Generate API documentation
	@echo "📚 Generating docs..."
	@swag init -g cmd/api/main.go -o docs

# Version info
version: ## Show version info  
	@echo "SubFlow v$(VERSION)"
	@echo "Architect: Muhammet Ali Büyük"
	@echo "Contact: iletisim@alibuyuk.net"
