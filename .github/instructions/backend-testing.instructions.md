---
description: Backend Testing Guidelines — E2E, Unit, Integration, Coverage
applyTo: "backend/**"
---

# Backend Testing Guidelines

**Coverage Target**: 80%+ for all new code  
**Framework**: `github.com/stretchr/testify` (suite-based)  
**Current State**: E2E tests exist, unit tests to be added following patterns below

## Test Organization

**Unit Tests**:

- Location: Within domain/platform folders (e.g., `internal/domain/order/service_test.go`)
- Purpose: Isolated component testing
- Pattern: Single responsibility, mocked dependencies
- Status: Not yet implemented — follow this pattern when adding

**Integration/E2E Tests**:

- Location: `backend/internal/tests/e2e/` (separate from production code)
- Structure: One suite per API endpoint, grouped by domain
- Files: `account_login_test.go`, `account_workspace_test.go`, `order_create_test.go`
- Entry: `main_test.go` contains `TestMain` (global setup/teardown)

**Utilities**:

- Helpers: `backend/internal/tests/testutils/` (fixtures, setup, common functions)
- Mocks: `backend/internal/tests/mocks/` (external dependency mocks)

## E2E Infrastructure

**Testcontainers**: Ephemeral Docker containers (Postgres, Memcached, Stripe-mock) for isolation.

**Global Setup** (`main_test.go`):

- `TestMain` starts containers, initializes server, tears down after all tests
- Shared resources: `testEnv` (containers), `testServer` (HTTP server on port 18080)
- Config: Mock email provider, test configuration

**Container Init** (via testutils):

- `CreateDatabaseCtx(ctx)` → Postgres container
- `CreateCacheCtx(ctx)` → Memcached container
- `CreateStripeMockCtx(ctx)` → Stripe mock server
- `InitEnvironment(ctx)` → aggregates all containers

**Requirements**: Docker Desktop running (testcontainers dependency)

## Suite Structure

**Standard Suite Pattern**:

```go
// File: backend/internal/tests/e2e/account_login_test.go
package e2e_test

import (
    "context"
    "net/http"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/abdelrahman146/kyora/internal/tests/testutils"
    "github.com/abdelrahman146/kyora/internal/domain/account"
)

// LoginSuite tests POST /v1/auth/login
type LoginSuite struct {
    suite.Suite
    client         *testutils.HTTPClient
    accountStorage *account.Storage
}

func (s *LoginSuite) SetupSuite() {
    // Initialize HTTP client and storage layers
    s.client = testutils.NewHTTPClient("http://localhost:18080")
    cache := cache.NewConnection([]string{"localhost:11211"})
    s.accountStorage = account.NewStorage(testEnv.Database, cache)
}

func (s *LoginSuite) SetupTest() {
    // CRITICAL: Clear database before each test
    testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
}

func (s *LoginSuite) TearDownTest() {
    // CRITICAL: Clean up after each test
    testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
}

func TestLoginSuite(t *testing.T) {
    if testServer == nil {
        t.Skip("Test server not initialized")
    }
    suite.Run(t, new(LoginSuite))
}
```

**Naming Convention**:

- Suite: `<Endpoint><Operation>Suite` (e.g., `CreateOrderSuite`, `LoginSuite`)
- File: `<domain>_<endpoint>_test.go` (e.g., `order_create_test.go`, `account_login_test.go`)
- Test function: Descriptive, indicates scenario (e.g., `TestLogin_Success`, `TestLogin_InvalidCredentials`)

## Critical Best Practices

### 1. Use Domain Storage Layers (NEVER Raw SQL)

**ALWAYS** use domain storage repositories for database operations.

```go
// ❌ BAD — Raw SQL
db.Exec("UPDATE onboarding_sessions SET stage = ? WHERE token = ?", stage, token)

// ✅ GOOD — Domain storage
onboardingStorage := onboarding.NewStorage(testEnv.Database, cache)
sess, _ := onboardingStorage.GetByToken(ctx, token)
sess.Stage = onboarding.SessionStage(stage)
onboardingStorage.UpdateSession(ctx, sess)
```

**Initialize Storage in SetupSuite**:

```go
func (s *LoginSuite) SetupSuite() {
    cache := cache.NewConnection([]string{"localhost:11211"})
    s.accountStorage = account.NewStorage(testEnv.Database, cache)
    s.onboardingStorage = onboarding.NewStorage(testEnv.Database, cache)
}
```

### 2. Assert ALL Response Fields

Verify EVERY field in API responses to catch unexpected changes.

```go
// Decode response
var result map[string]interface{}
s.NoError(testutils.DecodeJSON(resp, &result))

// Assert exact field count
s.Len(result, 2, "response should have exactly 2 fields")

// Assert field presence
s.Contains(result, "user")
s.Contains(result, "token")

// Assert exact values
s.NotEmpty(result["token"])
user := result["user"].(map[string]interface{})
s.Equal("test@example.com", user["email"])
s.Equal("John", user["firstName"])
s.Equal("Doe", user["lastName"])
s.Equal(true, user["isEmailVerified"])
s.NotEmpty(user["id"])
s.NotEmpty(user["workspaceId"])
```

**Field Count**: Use `s.Len(response, N)` to catch unexpected additions/removals.

**Nested Objects**: Validate all nested fields recursively.

**IDs/UUIDs**: Use `s.NotEmpty()` to verify generation.

**Timestamps**: Verify non-nil and reasonable ranges.

### 3. Test Isolation is Critical

Each test MUST be completely independent.

**Database Cleanup Pattern**:

```go
func (s *LoginSuite) SetupTest() {
    // Truncate BEFORE test runs
    testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
}

func (s *LoginSuite) TearDownTest() {
    // Truncate AFTER test runs
    testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
}
```

**Table-Driven Tests**: Create fresh data per iteration to avoid state conflicts.

```go
for i, tt := range tests {
    s.Run(tt.name, func() {
        // Create unique test data per iteration
        email := fmt.Sprintf("user%d@example.com", i)
        user := createTestUser(email, tt.password)
        // ... test logic
    })
}
```

**Never Reuse**: Don't reuse session tokens, users, or data across tests.

### 4. Proper Time Handling

- ALWAYS use UTC: `time.Now().UTC()`
- Expired times: `time.Now().UTC().Add(-20 * time.Minute)`
- Context: Use `context.Background()` for test operations

### 5. Domain-Specific Helpers

Create `<domain>_helpers_test.go` for domain-specific utilities.

```go
// File: backend/internal/tests/e2e/account_helpers_test.go
package e2e_test

type AccountTestHelper struct {
    storage *account.Storage
}

func (h *AccountTestHelper) CreateTestUser(email, password string) *account.User {
    user := &account.User{
        Email:    email,
        Password: hashPassword(password),
        // ... other fields
    }
    h.storage.CreateUser(context.Background(), user)
    return user
}

func (h *AccountTestHelper) CreateInvitation(workspaceID, email string) *account.Invitation {
    // ...
}
```

**Generic Helpers**: Place in `backend/internal/tests/testutils/` (setup, fixtures, HTTP client).

**Initialize in SetupSuite**:

```go
func (s *LoginSuite) SetupSuite() {
    cache := cache.NewConnection([]string{"localhost:11211"})
    s.accountHelper = &AccountTestHelper{storage: account.NewStorage(testEnv.Database, cache)}
}
```

### 6. Table-Driven Test Organization

Use sub-tests with `s.Run()` for clear output.

```go
func (s *LoginSuite) TestLogin_InputValidation() {
    tests := []struct {
        name           string
        email          string
        password       string
        expectedStatus int
        expectedError  string
    }{
        {"missing email", "", "Pass123!", 400, "email is required"},
        {"invalid email format", "not-an-email", "Pass123!", 400, "invalid email"},
        {"missing password", "user@example.com", "", 400, "password is required"},
    }

    for i, tt := range tests {
        s.Run(tt.name, func() {
            // Create unique test user per iteration
            email := fmt.Sprintf("user%d@example.com", i)
            user := s.accountHelper.CreateTestUser(email, "ValidPass123!")

            // Make request
            payload := map[string]interface{}{
                "email":    tt.email,
                "password": tt.password,
            }
            resp, err := s.client.Post("/v1/auth/login", payload)
            s.NoError(err)
            defer resp.Body.Close()

            // Assert status
            s.Equal(tt.expectedStatus, resp.StatusCode)

            // Assert error message
            if tt.expectedError != "" {
                var result map[string]interface{}
                s.NoError(testutils.DecodeJSON(resp, &result))
                s.Contains(result["detail"], tt.expectedError)
            }
        })
    }
}
```

**Descriptive Names**: `"missing email field"`, `"invalid email format"`, `"wrong password"`

**Fresh Data**: Create new user/session per iteration to avoid conflicts.

### 7. Security Testing Checklist

- **Authentication**: Missing/invalid/expired tokens
- **Authorization**: Permission boundaries (admin vs member)
- **Workspace & Business Isolation**: Users can't access other workspace data; business-owned resources are isolated per business
- **Input Validation**: SQL injection, XSS, malformed input
- **Rate Limiting**: Token reuse, expired tokens, multiple requests
- **CSRF Protection**: State tokens in OAuth flows
- **Enumeration Prevention**: Consistent error messages for existing/non-existing resources

Example:

```go
func (s *OrderSuite) TestCreateOrder_BusinessIsolation() {
    // Create two workspaces
    workspace1 := s.createTestWorkspace()
    workspace2 := s.createTestWorkspace()

    // Create a business in each workspace
    biz1 := s.createTestBusiness(workspace1.ID, "biz1")
    _ = s.createTestBusiness(workspace2.ID, "biz2")

    // User1 creates order in biz1
    user1Token := s.loginUser(workspace1.AdminEmail)
    orderPayload := map[string]interface{}{"total": 100}
    resp1, _ := s.client.Post("/v1/businesses/"+biz1.Descriptor+"/orders", orderPayload,
        testutils.WithAuth(user1Token))
    s.Equal(http.StatusCreated, resp1.StatusCode)

    // User2 tries to access biz1 (should fail)
    user2Token := s.loginUser(workspace2.AdminEmail)
    resp2, _ := s.client.Get("/v1/businesses/"+biz1.Descriptor+"/orders",
        testutils.WithAuth(user2Token))
    s.Equal(http.StatusForbidden, resp2.StatusCode)
}
```

### 8. External Dependencies (OAuth, Stripe)

External services may not be configured in test environment.

**Graceful Handling**:

```go
func (s *OAuthSuite) TestGoogleOAuth_Redirect() {
    resp, _ := s.client.Get("/v1/auth/google")

    // Handle missing config gracefully
    if resp.StatusCode >= 500 {
        s.T().Skip("Google OAuth not configured in test environment")
    }

    // Test with flexible assertions
    s.GreaterOrEqual(resp.StatusCode, 400)
}
```

**Consider Mocking**: For reliable test execution without external dependencies.

### 9. Response Validation Best Practices

- **Exact Field Count**: `s.Len(response, N)` catches unexpected changes
- **Field Presence**: `s.Contains(response, "fieldName")`
- **Exact Values**: `s.Equal(expected, actual)` where possible
- **IDs/UUIDs**: `s.NotEmpty(response["id"])` verifies generation
- **Timestamps**: Verify non-nil, reasonable ranges
- **Nested Objects**: Validate all nested fields recursively

## Running Tests

**Make Commands** (from repository root):

```bash
make test                 # All backend tests (verbose)
make test.unit            # Unit tests only (domain + platform)
make test.e2e             # E2E tests (120s timeout)
make test.quick           # All tests (no verbose, faster)
make test.coverage        # Tests + coverage report
make test.coverage.html   # Generate HTML coverage report
make test.coverage.view   # Generate + open coverage in browser
make test.e2e.coverage    # E2E tests + coverage
make clean.coverage       # Remove all coverage files
make help                 # Display all commands
```

**Go Commands** (from `backend/` directory):

```bash
cd backend && go test ./...                              # All tests
cd backend && go test ./internal/tests/e2e -v            # E2E tests
cd backend && go test ./... -cover -coverprofile=coverage.out  # With coverage
cd backend && go test ./internal/tests/e2e -v -run TestLoginSuite  # Specific suite
cd backend && go test ./internal/tests/e2e/account_login_test.go -v  # Specific file
cd backend && go test ./internal/tests/e2e -race         # Race detection
cd backend && go test -fuzz=FuzzFunctionName -fuzztime=30s  # Fuzzing
```

**Coverage Analysis**:

```bash
# Generate coverage report
cd backend && go test ./... -coverprofile=coverage.out

# View summary
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html  # macOS
```

## Test Quality Standards

**Coverage**:

- Target: 80%+ for all new code
- Monitor: `go test -cover`
- Detailed: `go test -coverprofile=coverage.out`
- Tools: `go tool cover -func=coverage.out`, `-html=coverage.out`

**Test Scenarios**:

- Positive cases (happy path)
- Negative cases (validation failures, not found, unauthorized)
- Edge cases (boundary conditions, empty inputs, large datasets)
- Error handling (DB errors, external service failures)

**Testing Techniques**:

- **Table-driven**: Multiple scenarios, input variations, edge cases
- **Fuzzing**: `testing.F` for user input, parsing, untrusted data
- **Race detection**: `go test -race` for concurrency issues

**Documentation**:

- Test purpose explained
- Setup requirements documented
- Descriptive test names
- Comments for complex logic

**Maintenance**:

- Review tests in code reviews
- Update with feature changes
- Keep tests <500 lines (split if needed)
- Include in Definition of Done

**Code Quality**:

- Fast-running tests (quick feedback)
- Descriptive function names
- Assertions (not print statements)
- Package suffix `_test` to avoid circular dependencies (e.g., `package e2e_test`)
- Test isolation (independent, stateless)
- No test logic in production code

## Example: Complete E2E Test Suite

```go
// File: backend/internal/tests/e2e/order_create_test.go
package e2e_test

import (
    "context"
    "fmt"
    "net/http"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/abdelrahman146/kyora/internal/tests/testutils"
    "github.com/abdelrahman146/kyora/internal/domain/order"
    "github.com/abdelrahman146/kyora/internal/domain/account"
)

// CreateOrderSuite tests POST /v1/businesses/{businessDescriptor}/orders
type CreateOrderSuite struct {
    suite.Suite
    client         *testutils.HTTPClient
    orderStorage   *order.Storage
    accountStorage *account.Storage
    accountHelper  *AccountTestHelper
}

func (s *CreateOrderSuite) SetupSuite() {
    s.client = testutils.NewHTTPClient("http://localhost:18080")
    cache := cache.NewConnection([]string{"localhost:11211"})
    s.orderStorage = order.NewStorage(testEnv.Database, cache)
    s.accountStorage = account.NewStorage(testEnv.Database, cache)
    s.accountHelper = &AccountTestHelper{storage: s.accountStorage}
}

func (s *CreateOrderSuite) SetupTest() {
    testutils.TruncateTables(testEnv.Database, "orders", "order_items", "users", "workspaces")
}

func (s *CreateOrderSuite) TearDownTest() {
    testutils.TruncateTables(testEnv.Database, "orders", "order_items", "users", "workspaces")
}

func (s *CreateOrderSuite) TestCreateOrder_Success() {
    ctx := context.Background()

    // Create user using helper
    user := s.accountHelper.CreateTestUser("admin@example.com", "Pass123!")
    token := s.accountHelper.LoginUser("admin@example.com", "Pass123!")

    // Create a business in the user's workspace (business-scoped APIs require businessDescriptor)
    // Implement this via API (`POST /v1/businesses`) or a helper that uses `business.NewStorage`.
    biz := s.createTestBusiness(user.WorkspaceID, "demo")

    // Prepare order payload
    payload := map[string]interface{}{
        "total":      100.50,
        "currency":   "USD",
        "customerID": "cust_123",
        "items": []map[string]interface{}{
            {"productID": "prod_1", "quantity": 2, "price": 50.25},
        },
    }

    // Make request
    url := fmt.Sprintf("/v1/businesses/%s/orders", biz.Descriptor)
    resp, err := s.client.Post(url, payload, testutils.WithAuth(token))
    s.NoError(err)
    defer resp.Body.Close()

    // Assert response
    s.Equal(http.StatusCreated, resp.StatusCode)

    var result map[string]interface{}
    s.NoError(testutils.DecodeJSON(resp, &result))

    // Assert exact structure
    s.Len(result, 6, "response should have 6 fields")
    s.Contains(result, "id")
    s.Contains(result, "total")
    s.Contains(result, "currency")
    s.Contains(result, "customerID")
    s.Contains(result, "items")
    s.Contains(result, "createdAt")

    // Assert exact values
    s.NotEmpty(result["id"])
    s.Equal(100.50, result["total"])
    s.Equal("USD", result["currency"])
    s.Equal("cust_123", result["customerID"])

    items := result["items"].([]interface{})
    s.Len(items, 1)

    // Verify database state
    // Verify database state (business-owned resources are scoped by business_id)
    orders, _ := s.orderStorage.FindMany(ctx, s.orderStorage.ScopeBusinessID(biz.ID))
    s.Len(orders, 1)
    s.Equal("100.50", orders[0].Total.String())
}

func (s *CreateOrderSuite) TestCreateOrder_ValidationErrors() {
    tests := []struct {
        name           string
        payload        map[string]interface{}
        expectedStatus int
        expectedField  string
    }{
        {
            "missing total",
            map[string]interface{}{"currency": "USD"},
            400,
            "total",
        },
        {
            "negative total",
            map[string]interface{}{"total": -10, "currency": "USD"},
            400,
            "total",
        },
        {
            "invalid currency",
            map[string]interface{}{"total": 100, "currency": "INVALID"},
            400,
            "currency",
        },
    }

    for i, tt := range tests {
        s.Run(tt.name, func() {
            // Create unique user per iteration
            email := fmt.Sprintf("user%d@example.com", i)
            user := s.accountHelper.CreateTestUser(email, "Pass123!")
            token := s.accountHelper.LoginUser(email, "Pass123!")

            // Make request
            url := fmt.Sprintf("/v1/businesses/%s/orders", biz.Descriptor)
            resp, err := s.client.Post(url, tt.payload, testutils.WithAuth(token))
            s.NoError(err)
            defer resp.Body.Close()

            // Assert error
            s.Equal(tt.expectedStatus, resp.StatusCode)

            var result map[string]interface{}
            s.NoError(testutils.DecodeJSON(resp, &result))
            s.Contains(result, "invalidFields")

            fields := result["invalidFields"].([]interface{})
            found := false
            for _, f := range fields {
                field := f.(map[string]interface{})
                if field["name"] == tt.expectedField {
                    found = true
                    break
                }
            }
            s.True(found, "expected validation error for field: %s", tt.expectedField)
        })
    }
}

func (s *CreateOrderSuite) TestCreateOrder_BusinessIsolation() {
    // Business-owned resources must not be accessible across workspaces.
    user1 := s.accountHelper.CreateTestUser("user1@example.com", "Pass123!")
    token1 := s.accountHelper.LoginUser("user1@example.com", "Pass123!")
    biz1 := s.createTestBusiness(user1.WorkspaceID, "biz1")

    user2 := s.accountHelper.CreateTestUser("user2@example.com", "Pass123!")
    token2 := s.accountHelper.LoginUser("user2@example.com", "Pass123!")

    // User1 creates order in biz1
    payload := map[string]interface{}{"total": 100, "currency": "USD"}
    url1 := fmt.Sprintf("/v1/businesses/%s/orders", biz1.Descriptor)
    resp1, _ := s.client.Post(url1, payload, testutils.WithAuth(token1))
    s.Equal(http.StatusCreated, resp1.StatusCode)

    // User2 tries to access biz1 (should fail)
    url2 := fmt.Sprintf("/v1/businesses/%s/orders", biz1.Descriptor)
    resp2, _ := s.client.Get(url2, testutils.WithAuth(token2))
    s.Equal(http.StatusForbidden, resp2.StatusCode)
}

func TestCreateOrderSuite(t *testing.T) {
    if testServer == nil {
        t.Skip("Test server not initialized")
    }
    suite.Run(t, new(CreateOrderSuite))
}
```

## Quick Reference

**Test File**: `<domain>_<endpoint>_test.go` → Suite per API endpoint  
**Storage Init**: `domain.NewStorage(testEnv.Database, cache)` in `SetupSuite()`  
**Cleanup**: `testutils.TruncateTables(db, "table1", "table2")` in `SetupTest()` + `TearDownTest()`  
**HTTP Client**: `testutils.NewHTTPClient("http://localhost:18080")`  
**Requests**: `s.client.Get/Post/Put/Delete(url, payload, testutils.WithAuth(token))`  
**Assertions**: `s.Equal()`, `s.Contains()`, `s.Len()`, `s.NotEmpty()`, `s.NoError()`  
**Helpers**: Domain-specific in `<domain>_helpers_test.go`, generic in `testutils/`  
**Context**: `context.Background()` for test operations  
**Time**: `time.Now().UTC()` for all timestamps  
**Isolation**: Fresh data per test, truncate tables before/after
