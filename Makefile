# EVA-Mind-FZPN Makefile
# Commands for testing, building, and running the project

.PHONY: test test-unit test-integration test-coverage test-all lint build clean help

# Default target
all: test build

# ============================================================================
# TESTING
# ============================================================================

# Run all tests
test:
	go test ./... -v -race

# Run unit tests only (fast, no database required)
test-unit:
	go test ./internal/cortex/scales/... \
		./internal/metrics/... \
		./internal/audit/... \
		./internal/cortex/alert/... \
		./internal/mocks/... \
		-v -race

# Run integration tests (requires database)
test-integration:
	go test ./test/integration/... -v -race -timeout 5m

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	go test ./... \
		-coverprofile=coverage/coverage.out \
		-covermode=atomic \
		-race \
		-timeout 10m
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	go tool cover -func=coverage/coverage.out | tail -1
	@echo ""
	@echo "Coverage report generated: coverage/coverage.html"

# Run tests with detailed coverage per package
test-coverage-detailed:
	@echo "Running detailed coverage analysis..."
	@mkdir -p coverage
	@for pkg in $$(go list ./internal/... | grep -v /mocks); do \
		name=$$(echo $$pkg | sed 's/.*\///'); \
		echo "Testing $$name..."; \
		go test -coverprofile=coverage/$$name.out -covermode=atomic $$pkg 2>/dev/null || true; \
	done
	@echo "Combining coverage..."
	@echo "mode: atomic" > coverage/combined.out
	@for f in coverage/*.out; do \
		if [ "$$f" != "coverage/combined.out" ]; then \
			tail -n +2 "$$f" >> coverage/combined.out 2>/dev/null || true; \
		fi \
	done
	go tool cover -func=coverage/combined.out | grep -E "^total:|coverage:" || true
	@echo ""
	@echo "Per-package coverage:"
	@go tool cover -func=coverage/combined.out | grep -v "^total:" | sort -t: -k3 -n | tail -20

# Run all tests (unit + integration)
test-all: test-unit test-integration

# Run critical tests only (C-SSRS, alerts, scales)
test-critical:
	go test ./internal/cortex/scales/... \
		./internal/cortex/alert/... \
		-v -race -run "CSSRS|Alert|Critical"

# Run benchmark tests
test-bench:
	go test ./... -bench=. -benchmem -run=^$$

# ============================================================================
# LINTING & FORMATTING
# ============================================================================

# Run linter
lint:
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Verify code
vet:
	go vet ./...

# ============================================================================
# BUILD
# ============================================================================

# Build the main application
build:
	go build -o bin/eva-mind ./main.go

# Build with optimizations
build-prod:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/eva-mind ./main.go

# Build integration service
build-integration:
	go build -o bin/integration-service ./cmd/integration_service/

# ============================================================================
# DATABASE
# ============================================================================

# Run migrations
migrate:
	@echo "Running migrations..."
	@for f in migrations/*.sql; do \
		echo "Applying $$f..."; \
		psql "$$DATABASE_URL" -f "$$f" || true; \
	done

# Verify database tables
verify-db:
	go run scripts/check_migrations.go

# ============================================================================
# DOCKER
# ============================================================================

# Start monitoring stack
monitoring-up:
	docker-compose -f deployments/docker-compose.monitoring.yml up -d

# Stop monitoring stack
monitoring-down:
	docker-compose -f deployments/docker-compose.monitoring.yml down

# ============================================================================
# UTILITIES
# ============================================================================

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf coverage/
	go clean -cache

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate mocks (if mockgen is installed)
mocks:
	@if command -v mockgen &> /dev/null; then \
		mockgen -source=internal/mocks/interfaces.go -destination=internal/mocks/generated_mocks.go -package=mocks; \
	else \
		echo "mockgen not installed. Run: go install github.com/golang/mock/mockgen@latest"; \
	fi

# Show help
help:
	@echo "EVA-Mind-FZPN Makefile"
	@echo ""
	@echo "Testing:"
	@echo "  make test              - Run all tests"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-critical     - Run critical tests (C-SSRS, alerts)"
	@echo "  make test-bench        - Run benchmark tests"
	@echo ""
	@echo "Building:"
	@echo "  make build             - Build the application"
	@echo "  make build-prod        - Build for production"
	@echo ""
	@echo "Database:"
	@echo "  make migrate           - Run database migrations"
	@echo "  make verify-db         - Verify database tables"
	@echo ""
	@echo "Other:"
	@echo "  make lint              - Run linter"
	@echo "  make fmt               - Format code"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make deps              - Download dependencies"
