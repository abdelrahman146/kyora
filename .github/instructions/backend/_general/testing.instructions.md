---
description: Backend Testing — E2E, Unit, Integration, Test Suite Patterns (Reusable)
applyTo: "backend/**"
---

# Backend Testing

**Reusable testing patterns for Go backends.**

Use when: Writing tests for backend services.

See also:

- `architecture.instructions.md` — Domain architecture patterns
- `go-patterns.instructions.md` — Go implementation patterns
- `errors.instructions.md` — Error handling patterns

---

## Test Organization

```
backend/
  internal/
    domain/<name>/
      service_test.go        # Unit tests
      storage_test.go        # Repository tests
    platform/<name>/
      *_test.go              # Platform tests
    tests/
      e2e/                   # End-to-end tests
        main_test.go         # Global setup/teardown
        account_login_test.go
        order_create_test.go
      testutils/             # Test helpers
        fixtures.go
        http_client.go
        database.go
      mocks/                 # External dependency mocks
```

**Test Types:**

- **Unit**: Single function/method, mocked dependencies
- **Integration**: Storage + database interaction
- **E2E**: Full HTTP stack, real dependencies

**Coverage Target**: 80%+ for all new code

---

## E2E Test Infrastructure

### Global Setup

**`main_test.go`**: Single TestMain for all E2E tests

```go
package e2e_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
)

var (
    testEnv    *testutils.Environment
    testServer *server.Server
)

func TestMain(m *testing.M) {
    ctx := context.Background()

    // Start containers (Postgres, Memcached, Stripe mock)
    env, err := testutils.InitEnvironment(ctx)
    if err != nil {
        panic(err)
    }
    testEnv = env

    // Start HTTP server
    srv, err := server.New(testConfig())
    if err != nil {
        panic(err)
    }
    testServer = srv
    go srv.Start(":18080")

    // Run tests
    code := m.Run()

    // Cleanup
    env.Cleanup()
    os.Exit(code)
}
```

### Container Initialization

Use testcontainers for ephemeral dependencies:

```go
func InitEnvironment(ctx context.Context) (*Environment, error) {
    // Postgres container
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:16"),
        postgres.WithDatabase("test_db"),
    )
    if err != nil {
        return nil, err
    }

    // Memcached container
    cacheContainer, err := testcontainers.GenericContainer(ctx,
        testcontainers.GenericContainerRequest{
            ContainerRequest: testcontainers.ContainerRequest{
                Image: "memcached:1.6",
            },
        },
    )
    if err != nil {
        return nil, err
    }

    // Connect to database
    dsn, _ := pgContainer.ConnectionString(ctx)
    db := database.NewConnection(dsn)

    return &Environment{
        Database: db,
        CacheHosts: []string{cacheHost},
        Cleanup: func() {
            pgContainer.Terminate(ctx)
            cacheContainer.Terminate(ctx)
        },
    }, nil
}
```

**Requirements**: Docker Desktop running

---

## Suite Pattern

### Standard Suite Structure

```go
package e2e_test

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

// CreateOrderSuite tests POST /v1/businesses/{businessDescriptor}/orders
type CreateOrderSuite struct {
    suite.Suite
    client         *testutils.HTTPClient
    orderStorage   *order.Storage
    accountStorage *account.Storage
}

func (s *CreateOrderSuite) SetupSuite() {
    // Initialize once per suite
    s.client = testutils.NewHTTPClient("http://localhost:18080")
    cache := cache.NewConnection([]string{"localhost:11211"})
    s.orderStorage = order.NewStorage(testEnv.Database, cache)
    s.accountStorage = account.NewStorage(testEnv.Database, cache)
}

func (s *CreateOrderSuite) SetupTest() {
    // CRITICAL: Clear database before each test
    testutils.TruncateTables(testEnv.Database, "orders", "order_items", "users", "workspaces")
}

func (s *CreateOrderSuite) TearDownTest() {
    // CRITICAL: Clean up after each test
    testutils.TruncateTables(testEnv.Database, "orders", "order_items", "users", "workspaces")
}

func TestCreateOrderSuite(t *testing.T) {
    if testServer == nil {
        t.Skip("Test server not initialized")
    }
    suite.Run(t, new(CreateOrderSuite))
}
```

**Naming:**

- Suite: `<Endpoint><Operation>Suite`
- File: `<domain>_<endpoint>_test.go`
- Test: `Test<Operation>_<Scenario>`

---

## Critical Best Practices

### 1. Use Domain Storage Layers

**NEVER use raw SQL in tests**. Always use domain storage repositories.

```go
// ❌ BAD: Raw SQL
db.Exec("UPDATE orders SET status = ? WHERE id = ?", "completed", orderID)

// ✅ GOOD: Domain storage
order, _ := s.orderStorage.FindOne(ctx, s.orderStorage.ScopeID(orderID))
order.Status = order.StatusCompleted
s.orderStorage.Update(ctx, order)
```

### 2. Assert ALL Response Fields

Verify every field to catch unexpected changes:

```go
var result map[string]interface{}
s.NoError(testutils.DecodeJSON(resp, &result))

// Assert exact field count
s.Len(result, 5, "response should have exactly 5 fields")

// Assert field presence
s.Contains(result, "id")
s.Contains(result, "total")
s.Contains(result, "status")

// Assert exact values
s.NotEmpty(result["id"])
s.Equal(100.50, result["total"])
s.Equal("pending", result["status"])
```

### 3. Test Isolation

Each test MUST be independent:

```go
func (s *Suite) SetupTest() {
    // Truncate BEFORE
    testutils.TruncateTables(testEnv.Database, "orders", "customers")
}

func (s *Suite) TearDownTest() {
    // Truncate AFTER
    testutils.TruncateTables(testEnv.Database, "orders", "customers")
}
```

**Never reuse** data across tests:

```go
// ❌ BAD: Shared user
user := s.createUser("test@example.com")

// ✅ GOOD: Unique user per test
user := s.createUser(fmt.Sprintf("user%d@example.com", time.Now().UnixNano()))
```

### 4. Time Handling

Always use UTC:

```go
// ✅ GOOD: UTC timestamps
now := time.Now().UTC()
expiresAt := time.Now().UTC().Add(-20 * time.Minute)

// ❌ BAD: Local time
now := time.Now()
```

### 5. Domain-Specific Helpers

Create helper methods for common operations:

```go
// account_helpers_test.go
type AccountTestHelper struct {
    storage *account.Storage
}

func (h *AccountTestHelper) CreateTestUser(email, password string) *account.User {
    user := &account.User{
        ID:       generateID("usr"),
        Email:    email,
        Password: hashPassword(password),
        // ...
    }
    h.storage.CreateUser(context.Background(), user)
    return user
}

func (h *AccountTestHelper) LoginUser(email, password string) string {
    // Make login request, return token
    resp, _ := h.client.Post("/v1/auth/login", map[string]interface{}{
        "email":    email,
        "password": password,
    })
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    return result["token"].(string)
}
```

### 6. Table-Driven Tests

Use sub-tests for multiple scenarios:

```go
func (s *Suite) TestCreateOrder_Validation() {
    tests := []struct {
        name           string
        payload        map[string]interface{}
        expectedStatus int
        expectedError  string
    }{
        {
            "missing total",
            map[string]interface{}{"customerId": "cus_123"},
            400,
            "total is required",
        },
        {
            "negative total",
            map[string]interface{}{"total": -10, "customerId": "cus_123"},
            400,
            "total must be positive",
        },
    }

    for i, tt := range tests {
        s.Run(tt.name, func() {
            // Create unique user per iteration
            user := s.createUser(fmt.Sprintf("user%d@example.com", i))
            token := s.loginUser(user.Email, "Pass123!")

            resp, _ := s.client.Post("/v1/orders", tt.payload,
                testutils.WithAuth(token))

            s.Equal(tt.expectedStatus, resp.StatusCode)

            if tt.expectedError != "" {
                var result map[string]interface{}
                testutils.DecodeJSON(resp, &result)
                s.Contains(result["detail"], tt.expectedError)
            }
        })
    }
}
```

### 7. Security Testing

Test authentication, authorization, and isolation:

```go
func (s *Suite) TestListOrders_BusinessIsolation() {
    // Create two workspaces
    user1 := s.createUser("user1@example.com")
    user2 := s.createUser("user2@example.com")

    biz1 := s.createBusiness(user1.WorkspaceID, "biz1")
    biz2 := s.createBusiness(user2.WorkspaceID, "biz2")

    // User1 creates order in biz1
    token1 := s.loginUser(user1.Email, "Pass123!")
    s.createOrder(biz1.Descriptor, token1)

    // User2 tries to access biz1 (should fail)
    token2 := s.loginUser(user2.Email, "Pass123!")
    resp, _ := s.client.Get("/v1/businesses/"+biz1.Descriptor+"/orders",
        testutils.WithAuth(token2))

    s.Equal(http.StatusForbidden, resp.StatusCode)
}
```

**Security Checklist:**

- ✅ Missing/invalid/expired tokens
- ✅ Permission boundaries (admin vs member)
- ✅ Workspace isolation
- ✅ Business isolation
- ✅ Input validation (SQL injection, XSS)

---

## Unit Testing

Mock dependencies, test business logic in isolation. Use testcontainers for repository tests with real database.

---

## Testing Utilities

Common helpers: HTTP client with auth, database truncation, JSON decoding.

---

## Running Tests

```bash
# All tests
go test ./...

# E2E only
go test ./internal/tests/e2e -v

# Specific suite
go test ./internal/tests/e2e -v -run TestLoginSuite

# With coverage
go test ./... -cover -coverprofile=coverage.out

# Race detection
go test ./... -race

# Quick tests (skip E2E)
go test ./internal/domain/... ./internal/platform/...
```

---

## Coverage Analysis

```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View summary
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html
```

**Coverage Targets:**

- Overall: 80%+
- Critical paths (auth, payments): 90%+
- New code: must not decrease coverage

---

## Test Quality Standards

✅ **Fast**: Unit tests <1s, E2E <5s per test  
✅ **Isolated**: No shared state between tests  
✅ **Reliable**: No flaky tests (race conditions, timing issues)  
✅ **Readable**: Descriptive names, clear assertions  
✅ **Maintainable**: DRY (use helpers), <500 lines per file

---

## Anti-Patterns

❌ Raw SQL in tests (use domain storage)  
❌ Shared test data across tests  
❌ Ignoring cleanup (database state leaks)  
❌ Missing field assertions (incomplete validation)  
❌ Skipping security tests (auth, isolation)  
❌ Local time instead of UTC  
❌ Hardcoded IDs/emails (causes conflicts)  
❌ Tests that depend on execution order

---

## Quick Reference

**E2E Setup**: `main_test.go` with testcontainers, global server  
**Suite Pattern**: `SetupSuite/SetupTest/TearDownTest`, truncate tables  
**Storage**: Always use domain storage, never raw SQL  
**Assertions**: Verify all fields, exact counts, nested objects  
**Isolation**: Unique data per test, truncate before/after  
**Helpers**: Domain-specific helpers in `*_helpers_test.go`  
**Table-Driven**: Use `s.Run()` for multiple scenarios  
**Security**: Test auth, authz, workspace/business isolation  
**Coverage**: 80%+ target, use `go tool cover`
