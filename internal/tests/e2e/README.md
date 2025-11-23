# E2E Testing Setup

This directory contains end-to-end tests for the Kyora API using testcontainers and testify suites.

## Overview

The e2e test setup automatically:
- Spins up ephemeral Docker containers (Postgres, Memcached, Stripe-mock)
- Initializes the HTTP server with test configuration
- Runs test suites against the live server
- Gracefully tears down all resources

## Prerequisites

- Docker Desktop running
- Go 1.21+
- All project dependencies installed (`go mod download`)

## Running Tests

Run all e2e tests:
```bash
go test ./internal/tests/e2e -v
```

Run with coverage:
```bash
go test ./internal/tests/e2e -v -cover -coverprofile=coverage.out
```

Run specific suite:
```bash
go test ./internal/tests/e2e -v -run TestOnboardingSuite
```

## Architecture

### TestMain (`main_test.go`)
- **Purpose**: Sets up global test environment before any tests run
- **Containers**: Postgres, Memcached, Stripe-mock
- **Server**: Starts on port 18080 with mock email provider
- **Cleanup**: Gracefully stops server and terminates containers

### Test Suites
Each domain has its own testify suite (e.g., `OnboardingSuite`):
- `SetupSuite()`: Runs once before all tests in suite
- `SetupTest()`: Runs before each test
- `TearDownTest()`: Runs after each test
- `TearDownSuite()`: Runs once after all tests

### Shared Resources
Global variables available to all test suites:
- `testEnv`: Contains database, cache, and Stripe mock references
- `testServer`: Running HTTP server instance

## Writing New Tests

### 1. Create a new test suite

```go
type MySuite struct {
    suite.Suite
    baseURL    string
    httpClient *http.Client
}

func (s *MySuite) SetupSuite() {
    s.baseURL = "http://localhost:18080"
    s.httpClient = &http.Client{}
}

func (s *MySuite) TestSomething() {
    resp, err := s.httpClient.Get(s.baseURL + "/api/endpoint")
    s.NoError(err)
    s.Equal(http.StatusOK, resp.StatusCode)
    defer resp.Body.Close()
}

func TestMySuite(t *testing.T) {
    if testServer == nil {
        t.Skip("Test server not initialized")
    }
    suite.Run(t, new(MySuite))
}
```

### 2. Use helper methods

The `OnboardingSuite` provides reusable helpers:
- `postJSON(path, payload)`: Make JSON POST requests
- `authenticatedRequest(method, path, body, token)`: Make authenticated requests

### 3. Access test environment

Access shared resources via global variables:
```go
// Direct database access
db := testEnv.Database.GetDB()

// Cache access
cache := testEnv.Cache

// Stripe mock base URL
stripeURL := testEnv.StripeMockBase
```

## Configuration

Test-specific overrides in `TestMain`:
```go
viper.Set(config.EmailProvider, "mock")  // Use mock email
viper.Set(config.HTTPPort, "18080")      // Test port
```

Main config file (`.kyora.yaml`) is auto-discovered from project root.

## Troubleshooting

### Tests hang or timeout
- Ensure Docker Desktop is running
- Check container logs: `docker ps` and `docker logs <container-id>`
- Increase timeout: `go test ./internal/tests/e2e -timeout=60s`

### Port conflicts
- Change test port in `TestMain`: `server.WithServerAddress(":19080")`

### Database errors
- Testcontainers creates fresh Postgres for each run
- Migrations run automatically on server start
- Some migration warnings are expected (e.g., pg_trgm extension)

### Email provider errors
- Verify `viper.Set(config.EmailProvider, "mock")` is set before `server.New()`
- Check `.kyora.yaml` has proper structure (nested `email.resend.*` keys)

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Use `defer resp.Body.Close()` for HTTP responses
3. **Assertions**: Use testify assertions for clear failure messages
4. **Database**: Consider database transactions/rollbacks for test isolation
5. **Parallel**: Avoid `t.Parallel()` unless you implement request isolation

## Example Test Flow

```go
func (s *MySuite) TestUserRegistration() {
    // 1. Register user
    payload := map[string]interface{}{
        "email": "test@example.com",
        "password": "SecurePass123!",
    }
    resp, err := s.postJSON("/api/onboarding/start", payload)
    s.NoError(err)
    s.Equal(http.StatusCreated, resp.StatusCode)
    defer resp.Body.Close()
    
    // 2. Parse response
    var result map[string]interface{}
    s.NoError(json.NewDecoder(resp.Body).Decode(&result))
    
    // 3. Verify email sent (mock captures it)
    // 4. Complete onboarding flow
}
```

## CI/CD Integration

For GitHub Actions:
```yaml
- name: Run E2E Tests
  run: go test ./internal/tests/e2e -v -timeout=120s
  env:
    DOCKER_HOST: unix:///var/run/docker.sock
```

## Further Reading

- [Testcontainers Go](https://golang.testcontainers.org/)
- [Testify Suite](https://pkg.go.dev/github.com/stretchr/testify/suite)
- [Project architecture](../../../README.md)
