.PHONY: all build run test clean migrate migrate-down dev install lint fmt vet check

# Binary name
BINARY_NAME=radix
MAIN_PACKAGE=./cmd/api
MIGRATE_PACKAGE=./cmd/migrate

# Build directory
BUILD_DIR=./build

# Default target
all: build

# Build the API binary
build:
	@echo "Building API server..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_NAME)"

# Build the migrate binary
build-migrate:
	@echo "Building migrate tool..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/migrate $(MIGRATE_PACKAGE)
	@echo "✓ Built $(BUILD_DIR)/migrate"

# Build both binaries
build-all: build build-migrate

# Run the API server
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode with hot reload (requires air)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air is not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Database migrations
migrate: build-migrate
	@echo "Running migrations up..."
	@$(BUILD_DIR)/migrate up

migrate-down: build-migrate
	@echo "Running migrations down..."
	@$(BUILD_DIR)/migrate down

migrate-version: build-migrate
	@echo "Checking migration version..."
	@$(BUILD_DIR)/migrate version

# Create new migration file using golang-migrate
migrate-new:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-new NAME=create_users_table"; \
		exit 1; \
	fi
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "Installing golang-migrate..."; \
		go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@migrate create -ext sql -dir migrations -seq $(NAME)
	@echo "✓ Created migration files for: $(NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

# Update dependencies
update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✓ Dependencies updated"

# Code quality checks
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Go vet passed"

check: fmt vet
	@echo "✓ All checks passed"

# Docker commands
docker-build:
	@docker build -t $(BINARY_NAME):latest .

docker-run:
	@docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

# Generate code (for future use)
generate:
	@go generate ./...

# Help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the API server"
	@echo "  make build-migrate - Build the migration tool"
	@echo "  make build-all     - Build both binaries"
	@echo "  make run           - Build and run the API server"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make migrate       - Run database migrations up"
	@echo "  make migrate-down  - Run database migrations down"
	@echo "  make migrate-new   - Create new migration files (NAME=your_migration_name)"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install       - Install dependencies"
	@echo "  make update        - Update dependencies"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make check         - Run fmt and vet"
	@echo "  make help          - Show this help message"
