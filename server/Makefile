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
	$(GOBUILD) -o $(BINARY_NAME).out -v ./...

# Test commands
test: test-unit test-integration test-router

test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v ./tests/db/...

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -run "TestIntegrationTestSuite" ./tests/

test-router:
	@echo "Running router tests..."
	$(GOTEST) -v -run "TestCreateRouter|TestRouterHandle|TestMatchFunction|TestRouterServeHTTP|TestCORSHeaders|TestOptionsRequest" ./tests/

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./... ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-verbose:
	@echo "Running all tests with verbose output..."
	$(GOTEST) -v ./... ./tests/...

test-short:
	@echo "Running short tests..."
	$(GOTEST) -short ./... ./tests/...

# Database specific tests
test-inmem:
	@echo "Running in-memory database tests..."
	$(GOTEST) -v -run "TestInMemoryDBTestSuite|TestHabitCopyIntegrity|TestTrackingEntryCopyIntegrity" ./tests/db/

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
	@echo "  make build          - Build the binary"
	@echo "  make test           - Run all tests (unit + integration + router)"
	@echo "  make test-unit      - Run unit tests only"
	@echo "  make test-integration - Run integration tests only" 
	@echo "  make test-router    - Run router tests only"
	@echo "  make test-inmem     - Run in-memory database tests only"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-short     - Run short tests"
	@echo "  make run            - Run with in-memory database"
	@echo "  make run-sqlite     - Run with SQLite database"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make deps-upgrade   - Upgrade all dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run linter (requires golangci-lint)"
	@echo "  make build-linux    - Cross compile for Linux"
	@echo "  make help           - Show this help"

.PHONY: all build test test-unit test-integration test-router test-coverage test-verbose test-short test-inmem clean run run-sqlite deps deps-upgrade fmt vet lint build-linux docker-build help 