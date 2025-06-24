.PHONY: help lint lint-fix test test-coverage build fmt deps clean check ci dev

# Colors
BLUE=\033[0;34m
GREEN=\033[0;32m
NC=\033[0m

.DEFAULT_GOAL := help

help:
	@echo "$(BLUE)🚀 Available commands:$(NC)"
	@echo "  make lint      - Run golangci-lint"
	@echo "  make lint-fix  - Run golangci-lint with fixes"
	@echo "  make test      - Run tests"
	@echo "  make test-cov  - Run tests with coverage"
	@echo "  make check     - Run lint + tests"
	@echo "  make fmt       - Format code"
	@echo "  make build     - Build project"
	@echo "  make deps      - Download dependencies"
	@echo "  make clean     - Clean temporary files"
	@echo "  make dev       - Development workflow (fmt + lint-fix + test)"
	@echo "  make ci        - CI pipeline (deps + fmt + lint + test)"

# Linting
lint:
	@echo "$(BLUE)🔍 Running staticcheck...$(NC)"
	GOFLAGS="-buildvcs=false" go run honnef.co/go/tools/cmd/staticcheck@latest ./...

lint-fix:
	@echo "$(BLUE)🔧 Running go fmt...$(NC)"
	go fmt ./...
	@echo "$(BLUE)🔍 Running staticcheck...$(NC)"
	GOFLAGS="-buildvcs=false" go run honnef.co/go/tools/cmd/staticcheck@latest ./...

# Tests
test:
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	go test ./...

test-verbose:
	@echo "$(BLUE)🧪 Running tests (verbose)...$(NC)"
	go test -v ./...

test-coverage:
	@echo "$(BLUE)🧪 Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)📊 Coverage saved to coverage.html$(NC)"

# Build
build:
	@echo "$(BLUE)🔨 Building project...$(NC)"
	go build ./...

# Formatting
fmt:
	@echo "$(BLUE)✨ Formatting code...$(NC)"
	go fmt ./...

# Dependencies
deps:
	@echo "$(BLUE)📦 Downloading dependencies...$(NC)"
	go mod download
	go mod tidy

# Cleanup
clean:
	@echo "$(BLUE)🧹 Cleaning up...$(NC)"
	go clean ./...
	rm -f coverage.out coverage.html

# Combined checks
check: lint test
	@echo "$(GREEN)✅ All checks passed!$(NC)"

# CI pipeline
ci: deps lint test
	@echo "$(GREEN)🚀 CI pipeline completed!$(NC)"

# Development workflow
dev: fmt lint-fix test
	@echo "$(GREEN)🎉 Development workflow completed!$(NC)"