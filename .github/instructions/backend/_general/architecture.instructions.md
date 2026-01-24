---
description: Backend Architecture — Domain-Driven Design, Layered Architecture, Service Patterns (Reusable)
applyTo: "backend/**"
---

# Backend Architecture

**Reusable domain-driven architecture patterns for Go backends.**

Use when: Creating new Go backend services or extending existing ones.

See also:

- `go-patterns.instructions.md` — Go-specific implementation patterns
- `testing.instructions.md` — Testing strategies
- `errors.instructions.md` — Error handling patterns
- `api-contracts.instructions.md` — Response/DTO patterns

## Core Principles

**Layered architecture with strict dependency direction:**

- `platform/` (infrastructure) ← `domain/` (business logic) ← `server/` (HTTP wiring)
- Platform provides primitives (DB, cache, auth, logging)
- Domain contains business rules and invariants
- Server wires dependencies and registers routes

**Multi-tenancy isolation (when applicable):**

- All queries scoped by tenant/workspace/business ID
- Never trust tenant IDs from client input for authorization
- Enforce boundaries in middleware + service layer

**Transaction safety:**

- Multi-entity writes run in transactions
- Use atomic processor abstraction (no manual begin/commit)
- Support retries for retryable errors

**Clean separation of concerns:**

- Storage layer: data access only (repositories, scopes)
- Service layer: business logic only (orchestration, validation)
- Handler layer: HTTP I/O only (parse, call service, respond)

---

## Folder Structure

```
backend/
  cmd/                    # CLI commands (server, seed, migrations)
  internal/
    server/               # HTTP server wiring (DI, routes)
    platform/             # Infrastructure (config, db, cache, auth, request/response, logger, types, utils)
    domain/               # Business modules (account, inventory, order, customer, etc.)
      <domain>/
        model.go          # GORM models, enums, schema definitions
        model_request.go  # Request DTOs (optional for complex domains)
        model_response.go # Response DTOs (required)
        storage.go        # Repositories, cache, query scopes
        service.go        # Business logic, orchestration
        errors.go         # Domain-specific problem constructors
        handler_http.go   # HTTP handlers
        middleware_http.go # Domain-specific middleware (optional)
    tests/                # E2E/integration tests
```

**Dependency rules (strict):**

- ✅ `domain/` may depend on `platform/`
- ❌ `platform/` must NOT depend on `domain/`
- ✅ Domains may call other domain services (avoid storage cross-access)
- ❌ Do not bypass service layer to access another domain's storage directly

---

## Domain Module Conventions

Each domain under `internal/domain/<name>/` typically contains:

### Model Files

**`model.go`**: GORM models, enums, constants, schema definitions

```go
// Constants for table/struct references (no magic strings)
const (
    OrderTable  = "orders"
    OrderStruct = "Order"
    OrderPrefix = "ord"
)

// Schema mappings (JSON field ↔ DB column translation)
var OrderSchema = struct {
    ID         schema.Field
    BusinessID schema.Field
    Total      schema.Field
    // ... all fields
}{
    ID:         schema.NewField("id", "id"),
    BusinessID: schema.NewField("business_id", "businessId"),
    Total:      schema.NewField("total", "total"),
}

// GORM model
type Order struct {
    ID         string          `gorm:"type:varchar(50);primaryKey"`
    BusinessID string          `gorm:"type:varchar(50);not null;index"`
    Total      decimal.Decimal `gorm:"type:numeric(12,2);not null"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  gorm.DeletedAt `gorm:"index"`
}
```

**`model_request.go`** (optional for complex domains): Request DTOs

```go
type CreateOrderRequest struct {
    Total      float64 `json:"total" binding:"required,gt=0"`
    CustomerID string  `json:"customerId" binding:"required"`
}
```

**`model_response.go`** (required): Response DTOs with camelCase JSON tags

```go
type OrderResponse struct {
    ID         string    `json:"id"`
    Total      string    `json:"total"`
    CustomerID string    `json:"customerId"`
    CreatedAt  time.Time `json:"createdAt"`
}

func ToOrderResponse(o *Order) *OrderResponse {
    return &OrderResponse{
        ID:         o.ID,
        Total:      o.Total.String(),
        CustomerID: o.CustomerID,
        CreatedAt:  o.CreatedAt,
    }
}
```

### Storage Layer

**`storage.go`**: Repositories, cache, query scopes

```go
type Storage struct {
    db    *gorm.DB
    cache *cache.Connection
    order *database.Repository[Order]
}

func NewStorage(db *gorm.DB, cache *cache.Connection) *Storage {
    return &Storage{
        db:    db,
        cache: cache,
        order: database.NewRepository[Order](db),
    }
}

// Scope methods encapsulate query construction
func (s *Storage) ScopeBusinessID(bizID string) func(*gorm.DB) *gorm.DB {
    return s.order.ScopeWhere(fmt.Sprintf("%s.%s = ?", OrderTable, OrderSchema.BusinessID.DB), bizID)
}

func (s *Storage) ScopeSearch(term string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        like := "%" + term + "%"
        return db.Where(fmt.Sprintf("%s.search_vector @@ websearch_to_tsquery('simple', ?)", OrderTable), term)
    }
}
```

**Rules:**

- All query construction belongs in storage layer
- Use schema constants (no magic strings for tables/columns)
- Encapsulate JOINs, aggregations, complex filters in scope methods
- Storage layer is data access only (no business logic)

### Service Layer

**`service.go`**: Business logic, orchestration, validation

```go
type Service struct {
    storage *Storage
    atomic  *database.AtomicProcess
    bus     *bus.Bus
}

func NewService(storage *Storage, atomic *database.AtomicProcess, bus *bus.Bus) *Service {
    return &Service{storage: storage, atomic: atomic, bus: bus}
}

// Service methods are thin orchestrators (20-60 lines max)
func (s *Service) CreateOrder(ctx context.Context, actor *account.User, biz *business.Business, req *CreateOrderRequest) (*Order, error) {
    // Validation
    if req.Total <= 0 {
        return nil, problem.BadRequest("total must be positive")
    }

    // Build model
    order := &Order{
        ID:         generateID(OrderPrefix),
        BusinessID: biz.ID,
        Total:      decimal.NewFromFloat(req.Total),
    }

    // Persist
    if err := s.storage.order.Create(ctx, order); err != nil {
        return nil, problem.Wrap(err, "failed to create order")
    }

    // Emit event
    s.bus.Emit(bus.OrderCreatedTopic, &bus.OrderCreatedEvent{
        OrderID:    order.ID,
        BusinessID: order.BusinessID,
    })

    return order, nil
}

// List pattern: build scopes, delegate to storage
func (s *Service) ListOrders(ctx context.Context, biz *business.Business, req *list.ListRequest) ([]*Order, int64, error) {
    scopes := []func(*gorm.DB) *gorm.DB{
        s.storage.ScopeBusinessID(biz.ID),
    }

    if req.SearchTerm() != "" {
        scopes = append(scopes, s.storage.ScopeSearch(req.SearchTerm()))
    }

    items, err := s.storage.order.FindMany(ctx,
        append(scopes,
            s.storage.order.WithPagination(req.Offset(), req.Limit()),
            s.storage.order.WithOrderBy(req.ParsedOrderBy(OrderSchema)),
        )...,
    )
    if err != nil {
        return nil, 0, problem.Wrap(err, "failed to list orders")
    }

    count, err := s.storage.order.Count(ctx, scopes...)
    if err != nil {
        return nil, 0, problem.Wrap(err, "failed to count orders")
    }

    return items, count, nil
}
```

**Rules:**

- Services orchestrate; they never construct queries
- Use repository scopes and schema constants only
- Keep methods thin (20-60 lines)
- Multi-entity writes must use atomic processor
- Always scope by tenant/business

### Handler Layer

**`handler_http.go`**: HTTP handlers (parse → call service → respond)

```go
type HttpHandler struct {
    service *Service
}

func NewHttpHandler(service *Service) *HttpHandler {
    return &HttpHandler{service: service}
}

// @Summary Create order
// @Tags orders
// @Accept json
// @Produce json
// @Param businessDescriptor path string true "Business descriptor"
// @Param request body CreateOrderRequest true "Order details"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} problem.Problem
// @Router /v1/businesses/{businessDescriptor}/orders [post]
func (h *HttpHandler) CreateOrder(c *gin.Context) {
    var req CreateOrderRequest
    if err := request.ValidBody(c, &req); err != nil {
        response.Error(c, err)
        return
    }

    actor := account.ActorFromContext(c)
    biz := business.BusinessFromContext(c)

    order, err := h.service.CreateOrder(c.Request.Context(), actor, biz, &req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.SuccessJSON(c, http.StatusCreated, ToOrderResponse(order))
}
```

**Rules:**

- Handlers are thin (no business logic)
- Always use `request.ValidBody` for JSON
- Always map GORM models to response DTOs
- Always use `response.SuccessJSON` or `response.Error`

---

## HTTP Server Wiring

**`internal/server/server.go`**: Dependency injection and initialization

```go
func New(cfg *viper.Viper) (*Server, error) {
    // Initialize platform dependencies
    db := database.NewConnection(dsn)
    cache := cache.NewConnection(cacheHosts)
    atomic := database.NewAtomicProcess(db)
    bus := bus.New()

    // Initialize domains (storage → service → handler)
    orderStorage := order.NewStorage(db, cache)
    orderService := order.NewService(orderStorage, atomic, bus)
    orderHandler := order.NewHttpHandler(orderService)

    // Create HTTP engine
    engine := gin.Default()
    engine.Use(logger.Middleware())

    // Register routes
    routes.Register(engine, orderHandler, ...)

    return &Server{engine: engine}, nil
}
```

**`internal/server/routes.go`**: Route registration and middleware chains

```go
func Register(engine *gin.Engine, orderHandler *order.HttpHandler, ...) {
    v1 := engine.Group("/v1")

    // Business-scoped protected routes
    businesses := v1.Group("/businesses/:businessDescriptor")
    businesses.Use(
        middleware.CORS(),
        auth.EnforceAuthentication,
        account.EnforceValidActor(accountService),
        account.EnforceWorkspaceMembership(accountService),
        business.EnforceBusinessValidity(businessService),
    )

    orders := businesses.Group("/orders")
    orders.POST("", orderHandler.CreateOrder)
    orders.GET("", orderHandler.ListOrders)
}
```

---

## Repository Pattern

Use typed repository wrapper per model:

```go
type Repository[T any] struct {
    db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
    // Auto-migrate if enabled
    return &Repository[T]{db: db}
}
```

**Common methods:**

- `Create(ctx, model)` — Insert single record
- `Update(ctx, model)` — Update single record
- `Delete(ctx, model)` — Soft delete
- `FindOne(ctx, scopes...)` — Fetch single record
- `FindMany(ctx, scopes...)` — Fetch multiple records
- `Count(ctx, scopes...)` — Count records
- `Exists(ctx, scopes...)` — Check existence

**Common scopes:**

- `ScopeID(id)` / `ScopeIDs([]id)`
- `ScopeWorkspaceID(wsID)` / `ScopeBusinessID(bizID)`
- `ScopeEquals(field, value)` / `ScopeIn(field, values)`
- `ScopeSearchTerm(term, fields...)`
- `ScopeCreatedAt(from, to)` / `ScopeTime(field, from, to)`
- `WithPreload(associations...)` / `WithJoins(joins...)`
- `WithPagination(offset, limit)` / `WithOrderBy([]string)`

**When JOINs exist, use qualified table names:**

```go
// ❌ BAD: Ambiguous when JOINs present
s.storage.order.ScopeBusinessID(bizID)

// ✅ GOOD: Qualified table name
s.storage.order.ScopeWhere("orders.business_id = ?", bizID)
```

---

## Transaction Pattern

Use atomic processor abstraction:

```go
err := s.atomic.Exec(ctx, func(txCtx context.Context) error {
    order, err := s.storage.order.Create(txCtx, order)
    if err != nil {
        return err
    }

    err = s.storage.orderItem.CreateMany(txCtx, items)
    if err != nil {
        return err
    }

    return nil
})
```

**Rules:**

- Multi-entity writes must use atomic processor
- Use `txCtx` (not `ctx`) inside transaction callback
- Return errors directly (processor handles rollback)
- Processor reuses existing transaction if present in context

---

## Request Validation

Always use strict JSON body validation:

```go
var req CreateOrderRequest
if err := request.ValidBody(c, &req); err != nil {
    response.Error(c, err)
    return
}
```

**Features:**

- Requires request body
- `DisallowUnknownFields()`
- Rejects trailing JSON tokens
- Runs struct tag validation (`binding:`)

---

## Multi-Tenancy Middleware

Standard middleware chain for protected routes:

1. **CORS** — Allow origins
2. **Authentication** — Verify JWT
3. **Actor** — Load user from JWT
4. **Workspace** — Load workspace (never trust URL param)
5. **Business** — Load business (scoped to workspace)
6. **RBAC** (optional) — Check permissions
7. **Plan Gates** (optional) — Check subscription/limits

**Context keys:**

- `auth.ClaimsKey` → JWT claims
- `account.ActorKey` → authenticated user
- `account.WorkspaceKey` → user's workspace
- `business.BusinessKey` → business (scoped to workspace)

---

## Event-Driven Automation

Use internal event bus for cross-domain workflows:

```go
// Publish event
bus.Emit(bus.OrderCreatedTopic, &bus.OrderCreatedEvent{
    OrderID:    order.ID,
    BusinessID: biz.ID,
})

// Subscribe to event
bus.Listen(bus.OrderCreatedTopic, func(payload any) {
    event := payload.(*bus.OrderCreatedEvent)
    // Handle event
})
```

**Rules:**

- Events dispatched asynchronously (non-blocking)
- Handlers must be idempotent
- Handler panics caught and logged

---

## Checklist

✅ Models in `model.go` with constants + schema  
✅ Response DTOs in `model_response.go` (camelCase)  
✅ Storage scopes encapsulate queries  
✅ Service methods thin (<60 lines)  
✅ Handlers parse/validate/call/respond only  
✅ No magic strings  
✅ Multi-entity writes use atomic processor  
✅ Queries scoped by tenant/business

## Anti-Patterns

❌ Queries in services, magic strings, business logic in handlers, raw GORM models in responses, cross-domain storage, manual transactions, trusting client tenant IDs, `float64` for money
