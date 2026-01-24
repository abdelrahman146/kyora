---
description: Kyora Backend Testing Specifics — Test Database, Seed Data, E2E Helpers
applyTo: "backend/**"
---

# Kyora Backend Testing Specifics

**Kyora-specific testing patterns and helpers.**

Use when: Writing E2E/integration tests for Kyora backend.

See also:

- `../general/testing.instructions.md` — General testing patterns
- `domain-modules.instructions.md` — Domain overview
- `.github/instructions/backend-testing.instructions.md` — Full testing guide

---

## Test Infrastructure

### Testcontainers Setup

Kyora E2E tests use Docker containers for isolation:

```go
// main_test.go
var (
    testEnv    *testutils.Environment
    testServer *server.Server
)

func TestMain(m *testing.M) {
    ctx := context.Background()

    // Start containers
    env, err := testutils.InitEnvironment(ctx)
    if err != nil {
        panic(err)
    }
    testEnv = env

    // Start HTTP server on port 18080
    srv, _ := server.New(testConfig())
    testServer = srv
    go srv.Start(":18080")

    // Run tests
    code := m.Run()

    // Cleanup
    env.Cleanup()
    os.Exit(code)
}
```

**Containers initialized:**

- PostgreSQL 16
- Memcached 1.6
- Stripe mock server

**Requirements**: Docker Desktop running

---

## Test Environment

```go
type Environment struct {
    Database   *gorm.DB
    CacheHosts []string
    StripeMock *StripeMockContainer
    Cleanup    func()
}

func InitEnvironment(ctx context.Context) (*Environment, error) {
    // Postgres
    pgContainer, _ := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:16"),
        postgres.WithDatabase("kyora_test"),
    )

    // Memcached
    cacheContainer, _ := testcontainers.GenericContainer(ctx, ...)

    // Stripe mock
    stripeMockContainer, _ := testcontainers.GenericContainer(ctx,
        testcontainers.GenericContainerRequest{
            ContainerRequest: testcontainers.ContainerRequest{
                Image: "stripe/stripe-mock:latest",
            },
        },
    )

    // Connect to database
    dsn, _ := pgContainer.ConnectionString(ctx)
    db := database.NewConnection(dsn)

    return &Environment{
        Database: db,
        CacheHosts: []string{cacheHost},
        StripeMock: stripeMock,
        Cleanup: func() {
            pgContainer.Terminate(ctx)
            cacheContainer.Terminate(ctx)
            stripeMockContainer.Terminate(ctx)
        },
    }, nil
}
```

---

## Domain-Specific Helpers

### Account Helpers

**File**: `backend/internal/tests/e2e/account_helpers_test.go`

```go
type AccountTestHelper struct {
    storage *account.Storage
    client  *testutils.HTTPClient
}

func NewAccountTestHelper(storage *account.Storage, client *testutils.HTTPClient) *AccountTestHelper {
    return &AccountTestHelper{storage: storage, client: client}
}

// CreateTestUser creates a user with workspace
func (h *AccountTestHelper) CreateTestUser(email, password string) *account.User {
    workspace := &account.Workspace{
        ID:   generateID("wks"),
        Name: "Test Workspace",
    }
    h.storage.CreateWorkspace(context.Background(), workspace)

    user := &account.User{
        ID:          generateID("usr"),
        Email:       email,
        Password:    hashPassword(password),
        FirstName:   "Test",
        LastName:    "User",
        WorkspaceID: workspace.ID,
    }
    h.storage.CreateUser(context.Background(), user)

    return user
}

// LoginUser performs login and returns JWT token
func (h *AccountTestHelper) LoginUser(email, password string) string {
    payload := map[string]interface{}{
        "email":    email,
        "password": password,
    }

    resp, _ := h.client.Post("/v1/auth/login", payload)
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    return result["token"].(string)
}

// CreateInvitation creates a workspace invitation
func (h *AccountTestHelper) CreateInvitation(workspaceID, email string) *account.Invitation {
    invitation := &account.Invitation{
        ID:          generateID("inv"),
        WorkspaceID: workspaceID,
        Email:       email,
        Token:       generateToken(),
        ExpiresAt:   time.Now().UTC().Add(24 * time.Hour),
    }
    h.storage.CreateInvitation(context.Background(), invitation)

    return invitation
}
```

### Business Helpers

**File**: `backend/internal/tests/e2e/business_helpers_test.go`

```go
type BusinessTestHelper struct {
    storage *business.Storage
}

// CreateBusiness creates a business for testing
func (h *BusinessTestHelper) CreateBusiness(workspaceID, descriptor string) *business.Business {
    biz := &business.Business{
        ID:          generateID("biz"),
        WorkspaceID: workspaceID,
        Descriptor:  descriptor,
        Name:        "Test Business",
        Currency:    "USD",
    }
    h.storage.CreateBusiness(context.Background(), biz)

    return biz
}

// CreateShippingZone creates a shipping zone for testing
func (h *BusinessTestHelper) CreateShippingZone(businessID string) *business.ShippingZone {
    zone := &business.ShippingZone{
        ID:         generateID("shz"),
        BusinessID: businessID,
        Name:       "Default Zone",
        Price:      decimal.NewFromFloat(5.00),
    }
    h.storage.CreateShippingZone(context.Background(), zone)

    return zone
}
```

### Inventory Helpers

```go
type InventoryTestHelper struct {
    storage *inventory.Storage
}

// CreateProduct creates a product with variant
func (h *InventoryTestHelper) CreateProduct(businessID, name string) (*inventory.Product, *inventory.Variant) {
    product := &inventory.Product{
        ID:         generateID("prd"),
        BusinessID: businessID,
        Name:       name,
    }
    h.storage.CreateProduct(context.Background(), product)

    variant := &inventory.Variant{
        ID:            generateID("var"),
        ProductID:     product.ID,
        BusinessID:    businessID,
        SKU:           "TEST-SKU",
        Price:         decimal.NewFromFloat(100.00),
        StockQuantity: 10,
    }
    h.storage.CreateVariant(context.Background(), variant)

    return product, variant
}
```

### Customer Helpers

```go
type CustomerTestHelper struct {
    storage *customer.Storage
}

// CreateCustomer creates a customer for testing
func (h *CustomerTestHelper) CreateCustomer(businessID, email string) *customer.Customer {
    cust := &customer.Customer{
        ID:         generateID("cus"),
        BusinessID: businessID,
        Email:      email,
        FirstName:  "Test",
        LastName:   "Customer",
    }
    h.storage.CreateCustomer(context.Background(), cust)

    return cust
}

// CreateAddress creates a customer address
func (h *CustomerTestHelper) CreateAddress(customerID, businessID string) *customer.Address {
    addr := &customer.Address{
        ID:         generateID("adr"),
        CustomerID: customerID,
        BusinessID: businessID,
        Street:     "123 Test St",
        City:       "Test City",
        Country:    "US",
    }
    h.storage.CreateAddress(context.Background(), addr)

    return addr
}
```

---

## Seed Data Patterns

### Minimal Test Data

Create only what's needed per test:

```go
func (s *Suite) TestCreateOrder() {
    // Create minimal required data
    user := s.accountHelper.CreateTestUser("test@example.com", "Pass123!")
    token := s.accountHelper.LoginUser("test@example.com", "Pass123!")
    biz := s.businessHelper.CreateBusiness(user.WorkspaceID, "demo")
    customer := s.customerHelper.CreateCustomer(biz.ID, "customer@example.com")
    product, variant := s.inventoryHelper.CreateProduct(biz.ID, "Test Product")

    // Test order creation
    payload := map[string]interface{}{
        "customerId": customer.ID,
        "items": []map[string]interface{}{
            {"variantId": variant.ID, "quantity": 2},
        },
    }

    resp, _ := s.client.Post("/v1/businesses/"+biz.Descriptor+"/orders", payload,
        testutils.WithAuth(token))

    s.Equal(http.StatusCreated, resp.StatusCode)
}
```

### Reusable Fixtures (Avoid)

Don't create shared fixtures across tests (breaks isolation):

```go
// ❌ BAD: Shared fixtures
var sharedUser *account.User

func (s *Suite) SetupSuite() {
    sharedUser = s.accountHelper.CreateTestUser("shared@example.com", "Pass123!")
}

func (s *Suite) TestA() {
    // Uses sharedUser - breaks isolation
}

// ✅ GOOD: Fresh data per test
func (s *Suite) TestA() {
    user := s.accountHelper.CreateTestUser(fmt.Sprintf("test%d@example.com", time.Now().UnixNano()), "Pass123!")
}
```

---

## Database Cleanup

Truncate all tables in `SetupTest/TearDownTest`:

```go
testutils.TruncateTables(testEnv.Database,
    "orders", "order_items", "products", "variants",
    "customers", "addresses", "businesses",
    "users", "workspaces", "sessions",
)
```

**Key tables**: users, workspaces, sessions, businesses, products, variants, orders, order_items, customers, addresses, assets, expenses, subscriptions

---

## Anti-Patterns

❌ Shared fixtures, skipping truncation, hardcoded IDs, raw SQL helpers, missing cleanup, assuming order, not verifying isolation

---

## Quick Reference

**Testcontainers**: Postgres + Memcached + Stripe mock  
**Helpers**: Account, business, inventory, customer helpers  
**Seed Data**: Create minimal data per test (avoid shared fixtures)  
**Cleanup**: Truncate all tables in SetupTest/TearDownTest  
**Isolation**: Test workspace/business boundaries  
**Plan Limits**: Test free vs paid plan enforcement  
**Commands**: `make test.e2e`, `make test.quick`, `make test.coverage`
