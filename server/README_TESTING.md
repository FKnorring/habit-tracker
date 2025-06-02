# Testing Framework Documentation

This document explains the comprehensive testing framework for the Habit Tracker Go server.

## Overview

Our testing framework provides multiple layers of testing with a hybrid directory structure:

1. **Router Tests** (`router_test.go`) - Test routing patterns, parameter extraction, and HTTP handling
2. **Integration Tests** (`integration_test.go`) - Test the complete application stack with in-memory database
3. **Unit Tests** (`tests/db/inmemdb_test.go`) - Test individual database operations and business logic

## Directory Structure

```
server/
├── router_test.go              # Router functionality tests (main package)
├── integration_test.go         # End-to-end integration tests (main package)
├── tests/
│   └── db/
│       └── inmemdb_test.go     # Database unit tests (separate package)
├── main.go                     # Main application
├── router.go                   # Custom HTTP router
├── db/                         # Database implementations
├── Makefile                    # Test automation
├── .gitignore                  # Git ignore patterns
└── README_TESTING.md          # This documentation
```

## Testing Framework Stack

- **Go's built-in `testing` package** - Foundation
- **`testify/assert`** - Rich assertions
- **`testify/suite`** - Organized test suites with setup/teardown
- **`httptest`** - HTTP testing utilities

## Test Types

### 1. Router Tests (`router_test.go`)

Tests the custom HTTP router implementation in the same package as the main application:

```bash
# Run router tests only
make test-router
```

**What it tests:**
- Route registration with `Handle()`
- Pattern matching (exact paths, parameters)
- Parameter extraction from URLs (e.g., `/habits/:id`)
- HTTP method routing
- CORS header injection
- OPTIONS preflight requests
- 404 handling for unmatched routes

**Example test cases:**
- `TestMatchFunction` - Tests URL pattern matching with various scenarios
- `TestRouterServeHTTP` - Tests complete HTTP request handling
- `TestCORSHeaders` - Verifies CORS headers are properly set

### 2. Integration Tests (`integration_test.go`)

Tests the complete application using `testify/suite` in the same package as the main application:

```bash
# Run integration tests only
make test-integration
```

**What it tests:**
- End-to-end HTTP API functionality
- All endpoints with real HTTP requests
- Database integration with in-memory database
- Complete workflows (create → read → update → delete)
- Error scenarios and edge cases
- JSON request/response handling

**Key features:**
- Uses `testify/suite` for organized test structure
- Automatic setup/teardown of test environment
- Fresh database state for each test
- Real HTTP server with `httptest.Server`

**Example test cases:**
- `TestCreateHabit` - Tests POST /habits endpoint
- `TestFullWorkflow` - Tests complete CRUD workflow
- `TestGetTrackingEmpty` - Tests edge case scenarios

### 3. Unit Tests (`tests/db/inmemdb_test.go`)

Tests the in-memory database implementation in a separate package:

```bash
# Run unit tests only
make test-unit
```

**What it tests:**
- Database interface compliance
- CRUD operations for habits and tracking entries
- Error handling (not found, duplicates)
- Data isolation and copying
- Edge cases and boundary conditions

**Example test cases:**
- `TestCreateHabit` - Tests habit creation and storage
- `TestGetTrackingEntriesByHabitID` - Tests filtering tracking entries
- `TestHabitCopyIntegrity` - Ensures data integrity through copying

## Running Tests

### Quick Commands

```bash
# Run all tests
make test

# Run specific test types
make test-router       # Router tests only
make test-unit         # Unit tests only
make test-integration  # Integration tests only
make test-inmem        # In-memory DB tests only

# Run with coverage
make test-coverage     # Generates coverage.html report

# Run with verbose output
make test-verbose

# Run short tests (excludes long-running tests)
make test-short
```

### Manual Test Running

```bash
# Run all tests with verbose output
go test -v ./... ./tests/...

# Run specific test suites
go test -v -run "TestIntegrationTestSuite" .
go test -v -run "TestInMemoryDBTestSuite" ./tests/db/

# Run router tests
go test -v . -run "TestCreateRouter|TestRouterHandle"

# Run with coverage
go test -v -coverprofile=coverage.out ./... ./tests/...
go tool cover -html=coverage.out -o coverage.html
```

## Test Structure

### Integration Test Suite Structure

```go
type IntegrationTestSuite struct {
    suite.Suite
    router   *Router
    server   *httptest.Server
    origDB   db.Database
}

func (suite *IntegrationTestSuite) SetupSuite() {
    // One-time setup for the entire suite
}

func (suite *IntegrationTestSuite) SetupTest() {
    // Setup before each test (fresh database)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
    // Cleanup after all tests
}
```

### Test Data Patterns

**Habits:**
```go
habit := &db.Habit{
    ID:          "test-habit-1",
    Name:        "Exercise",
    Description: "Daily workout routine",
    Frequency:   "daily",
    StartDate:   "2024-01-01",
}
```

**Tracking Entries:**
```go
entry := &db.TrackingEntry{
    ID:        "entry-1",
    HabitID:   "habit-1",
    Timestamp: time.Now().Format(time.RFC3339),
    Note:      "Great workout!",
}
```

## Best Practices

### 1. Test Isolation
- Each test has a fresh database state
- Tests don't depend on each other
- Clean setup and teardown

### 2. Comprehensive Coverage
- Test success paths
- Test error conditions
- Test edge cases
- Test invalid input

### 3. Realistic Test Data
- Use meaningful test data
- Test with various data combinations
- Include edge cases (empty strings, large values)

### 4. Clear Test Names
- Descriptive test function names
- Clear test case descriptions
- Organized test groups

## Example Test Workflows

### Testing a New Endpoint

1. **Add router test** in `router_test.go` for URL pattern matching
2. **Add integration test** in `integration_test.go` for complete HTTP flow
3. **Add unit tests** in `tests/db/` for any new database operations
4. **Run coverage** to ensure good test coverage

### Testing Error Handling

```go
func (suite *IntegrationTestSuite) TestCreateHabitInvalidJSON() {
    resp, err := http.Post(suite.server.URL+"/habits", 
        "application/json", 
        bytes.NewBuffer([]byte("invalid json")))
    suite.NoError(err)
    defer resp.Body.Close()
    
    suite.Equal(http.StatusBadRequest, resp.StatusCode)
}
```

### Testing Database Errors

```go
func (suite *InMemoryDBTestSuite) TestCreateHabitDuplicate() {
    // Create habit first time
    err := suite.db.CreateHabit(habit)
    suite.NoError(err)
    
    // Try to create same habit again
    err = suite.db.CreateHabit(habit)
    suite.Equal(db.ErrDuplicate, err)
}
```

## Coverage Reports

Generate and view coverage reports:

```bash
make test-coverage
# Opens coverage.html in your browser
```

Coverage reports help identify:
- Untested code paths
- Areas needing more tests
- Overall test quality

## Continuous Integration

The testing framework is designed for CI/CD with GitHub Actions:

```bash
# CI-friendly commands
make test          # Run all tests
make test-coverage # Generate coverage for reporting
make build         # Ensure code compiles
```

## Development Workflow

1. **Write failing test** for new feature
2. **Implement feature** to make test pass
3. **Run specific tests** for fast feedback
4. **Run full test suite** before committing
5. **Check coverage** to ensure adequate testing

## Common Test Patterns

### HTTP Request Testing
```go
req := httptest.NewRequest("GET", "/habits/123", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
```

### JSON Request/Response Testing
```go
// Marshal request
jsonData, _ := json.Marshal(habitData)
resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

// Unmarshal response
var result db.Habit
json.NewDecoder(resp.Body).Decode(&result)
```

### Database State Verification
```go
// Verify creation
stored, err := database.GetHabit(habitID)
suite.NoError(err)
suite.Equal(expectedName, stored.Name)

// Verify deletion
_, err = database.GetHabit(habitID)
suite.Equal(db.ErrNotFound, err)
```

## File Organization

The tests use a hybrid approach for optimal organization:

- **Main directory tests** (`router_test.go`, `integration_test.go`) - Can access unexported functions and variables from the main package
- **Separate directory tests** (`tests/db/inmemdb_test.go`) - Provide clean separation for database layer testing

This structure provides:
- **Direct access** to main package internals for router and integration tests
- **Clean separation** for database layer unit tests
- **Easy navigation** to find specific test types
- **Maintainable organization** as the test suite grows
- **CI/CD friendly** structure for automated testing

## .gitignore Coverage

The comprehensive `.gitignore` file includes:
- Build artifacts and binaries
- Test coverage reports
- Database files (SQLite, test databases)
- IDE and editor files
- OS-generated files
- Environment variables and configuration files
- Logs and temporary files

This testing framework provides comprehensive coverage for both the router functionality and the complete application stack, ensuring reliability and maintainability of your Go HTTP server. 