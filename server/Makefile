# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=habit-tracker
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the project
all: test build

build:
	$(GOBUILD) -o out/ -v ./...

# Test commands
test: test-unit test-integration test-router test-auth test-reminder

test-all:
	@echo "Running all tests..."
	$(GOTEST) -v ./tests/...

test-unit:
	@echo "Running database unit tests..."
	$(GOTEST) -v ./tests/db/...

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -run "TestIntegrationTestSuite" ./tests/

test-router:
	@echo "Running router tests..."
	$(GOTEST) -v -run "TestCreateRouter|TestRouterHandle|TestMatchFunction|TestRouterServeHTTP|TestCORSHeaders|TestOPTIONSRequest" ./tests/

test-auth:
	@echo "Running authentication tests..."
	$(GOTEST) -v ./tests/auth/...

test-reminder:
	@echo "Running reminder service tests..."
	$(GOTEST) -v ./tests/reminder/...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-verbose:
	@echo "Running all tests with verbose output..."
	$(GOTEST) -v ./tests/...

test-short:
	@echo "Running short tests..."
	$(GOTEST) -short ./tests/...

# Specific test suites
test-inmem:
	@echo "Running in-memory database tests..."
	$(GOTEST) -v -run "TestInMemoryDBTestSuite|TestHabitCopyIntegrity|TestTrackingEntryCopyIntegrity|TestReminderCopyIntegrity|TestCalculateNextReminderTime" ./tests/db/

test-auth-core:
	@echo "Running core authentication tests..."
	$(GOTEST) -v -run "TestAuthTestSuite" ./tests/auth/

test-auth-middleware:
	@echo "Running authentication middleware tests..."
	$(GOTEST) -v -run "TestMiddlewareTestSuite" ./tests/auth/

test-auth-handlers:
	@echo "Running authentication handler tests..."
	$(GOTEST) -v -run "TestHandlersTestSuite" ./tests/auth/

# Clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out coverage.html

# Run with different database configurations
run:
	$(GOCMD) run . -db-driver=memory

run-sqlite:
	$(GOCMD) run . -db-driver=sqlite -sqlite-path=./test.db

# Development tasks
deps:
	$(GOMOD) download
	$(GOMOD) tidy

deps-upgrade:
	$(GOMOD) get -u all
	$(GOMOD) tidy

fmt:
	$(GOCMD) fmt ./...

vet:
	$(GOCMD) vet ./...

lint:
	golangci-lint run

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Docker related (if needed)
docker-build:
	docker build -t $(BINARY_NAME) .

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Build commands:"
	@echo "  make build          - Build the binary"
	@echo "  make build-linux    - Cross compile for Linux"
	@echo ""
	@echo "Test commands:"
	@echo "  make test           - Run all major test suites"
	@echo "  make test-all       - Run all tests in tests/ directory"
	@echo "  make test-unit      - Run database unit tests only"
	@echo "  make test-integration - Run integration tests only" 
	@echo "  make test-router    - Run router tests only"
	@echo "  make test-auth      - Run all authentication tests"
	@echo "  make test-reminder  - Run reminder service tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-short     - Run short tests"
	@echo ""
	@echo "Specific test suites:"
	@echo "  make test-inmem     - Run in-memory database tests"
	@echo "  make test-auth-core - Run core authentication tests"
	@echo "  make test-auth-middleware - Run auth middleware tests"
	@echo "  make test-auth-handlers - Run auth handler tests"
	@echo ""
	@echo "Development commands:"
	@echo "  make run            - Run with in-memory database"
	@echo "  make run-sqlite     - Run with SQLite database"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make deps-upgrade   - Upgrade all dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run linter (requires golangci-lint)"
	@echo "  make help           - Show this help"

.PHONY: all build test test-all test-unit test-integration test-router test-auth test-reminder test-coverage test-verbose test-short test-inmem test-auth-core test-auth-middleware test-auth-handlers clean run run-sqlite deps deps-upgrade fmt vet lint build-linux docker-build help 