.PHONY: help lint lint-fix lint-integration test test-verbose test-coverage test-integration build fmt deps clean check ci dev

# Colors
BLUE=\033[0;34m
GREEN=\033[0;32m
NC=\033[0m

.DEFAULT_GOAL := help

help:
	@echo "$(BLUE)🚀 Available commands:$(NC)"
	@echo "  make lint      - Run golangci-lint"
	@echo "  make lint-fix  - Run golangci-lint with fixes"
	@echo "  make lint-integration - Type-check the integration suite (part of check/ci)"
	@echo "  make test      - Run tests (unit only, integration suite is behind a build tag)"
	@echo "  make test-cov  - Run tests with coverage"
	@echo "  make test-integration - Run the integration suite against the real Pingen staging API (needs .env)"
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

lint-integration:
	@echo "$(BLUE)🔍 Running staticcheck (integration build tag)...$(NC)"
	GOFLAGS="-buildvcs=false" go run honnef.co/go/tools/cmd/staticcheck@latest -tags integration ./...

# Tests
test:
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	go test ./...

test-verbose:
	@echo "$(BLUE)🧪 Running tests (verbose)...$(NC)"
	go test -v ./...

# Hits the real Pingen staging API. Reads credentials from .env
# (see .env.example); skips itself when they are absent.
test-integration:
	@echo "$(BLUE)🧪 Running integration tests (real staging API)...$(NC)"
	go test -tags integration -v -timeout 30m ./integration/...

test-cov:
	@echo "$(BLUE)🧪 Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./...
	@echo "$(GREEN)📊 Coverage Report:"
	go tool cover -func=coverage.out

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
check: lint lint-integration test
	@echo "$(GREEN)✅ All checks passed!$(NC)"

# CI pipeline
ci: deps lint lint-integration test-cov
	@echo "$(GREEN)🚀 CI pipeline completed!$(NC)"

# Development workflow
dev: fmt lint-fix test
	@echo "$(GREEN)🎉 Development workflow completed!$(NC)"