---
description: Go Backend Patterns — Error Handling, Middleware, Context, GORM (Reusable)
applyTo: "backend/**"
---

# Go Backend Patterns

**Reusable Go implementation patterns for backend services.**

Use when: Writing Go backend code (any project).

See also:

- `architecture.instructions.md` — Domain-driven architecture patterns
- `testing.instructions.md` — Testing strategies
- `errors.instructions.md` — Error handling patterns
- `api-contracts.instructions.md` — Response/DTO patterns

---

## Error Handling

### Error Wrapping

Always wrap errors with context:

```go
if err != nil {
    return problem.Wrap(err, "failed to create order")
}
```

**Rules:**

- Wrap at boundaries (service → storage, external API calls)
- Preserve original error for stack traces
- Add user-facing context message

### Domain Errors

Create domain-specific errors in `errors.go`:

```go
func OrderNotFound(orderID string) *problem.Problem {
    return problem.NotFound("order not found").
        WithCode("order.not_found").
        With("orderId", orderID)
}

func InsufficientStock(productID string, available, requested int) *problem.Problem {
    return problem.BadRequest("insufficient stock").
        WithCode("inventory.insufficient_stock").
        With("productId", productID).
        With("available", available).
        With("requested", requested)
}
```

**Rules:**

- All errors must include `extensions.code`
- Use problem constructors (`problem.BadRequest`, `problem.Forbidden`, etc.)
- Enrich with structured context via `.With(key, value)`
- Never include secrets/PII in error details

### Error Return Pattern

```go
// ✅ GOOD: Check and return immediately
order, err := s.storage.GetOrder(ctx, orderID)
if err != nil {
    return nil, problem.Wrap(err, "failed to get order")
}

// ❌ BAD: Nested error handling
order, err := s.storage.GetOrder(ctx, orderID)
if err == nil {
    // ... lots of code
} else {
    return nil, err
}
```

---

## Context Usage

### Context Propagation

Always pass context as first parameter:

```go
func (s *Service) CreateOrder(ctx context.Context, actor *account.User, biz *business.Business, req *CreateOrderRequest) (*Order, error) {
    // Use ctx for all downstream calls
    order, err := s.storage.order.Create(ctx, order)
    if err != nil {
        return nil, err
    }
    return order, nil
}
```

**Rules:**

- Context is always first parameter
- Never store context in struct fields
- Pass context to all DB/cache/external API calls
- Use `context.Background()` only in tests or top-level initialization

### Context Values

Extract values from context using typed helpers:

```go
// In middleware: store in context
c.Set(account.ActorKey, user)

// In handlers/services: extract from context
actor := account.ActorFromContext(c)
workspace := account.WorkspaceFromContext(c)
biz := business.BusinessFromContext(c)
```

**Rules:**

- Use typed extraction helpers (not `c.MustGet()` directly)
- Always check for nil when extracting
- Use constants for context keys (not magic strings)

---

## GORM Patterns

### Model Definition

```go
type Order struct {
    ID         string          `gorm:"type:varchar(50);primaryKey"`
    BusinessID string          `gorm:"type:varchar(50);not null;index"`
    Total      decimal.Decimal `gorm:"type:numeric(12,2);not null"`
    Status     OrderStatus     `gorm:"type:varchar(20);not null;default:'pending'"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  gorm.DeletedAt  `gorm:"index"`

    // Relations (not serialized)
    Items []OrderItem `gorm:"foreignKey:OrderID" json:"-"`
}
```

**Rules:**

- Use `varchar(50)` for IDs (UUIDs or prefixed IDs)
- Use `numeric(12,2)` for money (never float)
- Use `varchar(20)` for enums
- Add `index` for foreign keys and frequently queried columns
- Use `json:"-"` for relations (prevent accidental serialization)
- Include `CreatedAt`, `UpdatedAt`, `DeletedAt` (soft deletes)

### Query Construction

Always use scopes (never raw SQL in services):

```go
// ✅ GOOD: Scopes in storage layer
func (s *Storage) ScopeBusinessID(bizID string) func(*gorm.DB) *gorm.DB {
    return s.order.ScopeWhere(fmt.Sprintf("%s.%s = ?", OrderTable, OrderSchema.BusinessID.DB), bizID)
}

// Service uses scopes
orders, err := s.storage.order.FindMany(ctx,
    s.storage.ScopeBusinessID(biz.ID),
    s.storage.order.WithPreload(OrderItemStruct),
)

// ❌ BAD: Raw SQL in service
db.Where("business_id = ?", biz.ID).Preload("Items").Find(&orders)
```

### Preloading Relations

Use constants for relation names:

```go
// In model.go
const (
    OrderItemStruct = "Items"
    OrderNoteStruct = "Notes"
)

// In service
order, err := s.storage.order.FindOne(ctx,
    s.storage.order.ScopeID(orderID),
    s.storage.order.WithPreload(OrderItemStruct),
    s.storage.order.WithPreload(OrderNoteStruct),
)
```

### Transaction-Aware Queries

Always use `db.Conn(ctx)` (not `db` directly):

```go
func (r *Repository[T]) Create(ctx context.Context, model *T) error {
    // Conn() uses transaction if present in context
    return r.db.Conn(ctx).Create(model).Error
}
```

---

## Middleware Patterns

### Standard Middleware Signature

```go
func EnforceAuthentication(c *gin.Context) {
    token := extractBearerToken(c)
    if token == "" {
        response.Error(c, problem.Unauthorized("missing authorization header"))
        c.Abort()
        return
    }

    claims, err := validateJWT(token)
    if err != nil {
        response.Error(c, problem.Unauthorized("invalid token"))
        c.Abort()
        return
    }

    c.Set(auth.ClaimsKey, claims)
    c.Next()
}
```

**Rules:**

- Return early on errors (call `c.Abort()`)
- Store validated data in context via `c.Set()`
- Call `c.Next()` to continue chain
- Use `response.Error()` for consistent error responses

### Middleware with Dependencies

```go
func EnforceValidActor(accountService AccountService) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := auth.ClaimsFromContext(c)
        user, err := accountService.GetUser(c.Request.Context(), claims.UserID)
        if err != nil {
            response.Error(c, problem.Unauthorized("invalid actor"))
            c.Abort()
            return
        }

        c.Set(account.ActorKey, user)
        c.Next()
    }
}
```

---

## Configuration Patterns

### Config Keys as Constants

```go
// In config/config.go
const (
    HTTPPort         = "http.port"
    DatabaseURL      = "database.url"
    StripeAPIKey     = "billing.stripe.api_key"
    EmailProvider    = "email.provider"
)
```

**Rules:**

- Never use magic strings for config keys
- Group related keys by prefix
- Document defaults in `Configure()` function

### Config Loading

```go
func Load() (*viper.Viper, error) {
    v := viper.New()

    // Set defaults
    Configure(v)

    // Load from file
    v.SetConfigName("config")
    v.AddConfigPath(".")
    _ = v.ReadInConfig()

    // Override with env vars
    v.AutomaticEnv()

    return v, nil
}
```

---

## Logging Patterns

### Structured Logging

```go
logger := logger.FromContext(ctx)
logger.Info("order created",
    slog.String("orderId", order.ID),
    slog.String("businessId", order.BusinessID),
    slog.Float64("total", order.Total),
)
```

**Rules:**

- Use structured logging (`slog.String`, `slog.Int`, etc.)
- Extract logger from context (includes request ID, actor info)
- Never log secrets, tokens, or PII
- Use appropriate levels: Debug, Info, Warn, Error

### Request Logging Middleware

```go
func Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // Generate trace ID
        traceID := generateTraceID()
        c.Set(logger.TraceIDKey, traceID)

        // Create request logger
        reqLogger := slog.Default().With(
            slog.String("traceId", traceID),
            slog.String("method", c.Request.Method),
            slog.String("path", c.Request.URL.Path),
        )
        c.Set(logger.LoggerKey, reqLogger)

        reqLogger.Info("request started")

        c.Next()

        reqLogger.Info("request completed",
            slog.Int("status", c.Writer.Status()),
            slog.Duration("duration", time.Since(start)),
        )
    }
}
```

---

## Decimal Handling

Always use `shopspring/decimal` for money:

```go
import "github.com/shopspring/decimal"

// ✅ GOOD: Decimal for money
type Order struct {
    Total decimal.Decimal `gorm:"type:numeric(12,2);not null"`
}

func calculateTotal(items []OrderItem) decimal.Decimal {
    total := decimal.Zero
    for _, item := range items {
        total = total.Add(item.Price.Mul(decimal.NewFromInt(int64(item.Quantity))))
    }
    return total
}

// ❌ BAD: Float for money
type Order struct {
    Total float64 `gorm:"type:float"`
}
```

---

## Time Handling

Always use UTC:

```go
// ✅ GOOD: UTC timestamps
now := time.Now().UTC()
expiresAt := time.Now().UTC().Add(24 * time.Hour)

// ❌ BAD: Local time
now := time.Now()
```

**Rules:**

- Store timestamps as UTC in database
- Convert to user's timezone only in UI layer
- Use `time.Time` (not string timestamps)

---

## ID Generation

Use prefixed IDs for readability:

```go
func generateID(prefix string) string {
    return prefix + "_" + uuid.New().String()[:16]
}

// Usage
orderID := generateID("ord")     // ord_abc123def456
customerID := generateID("cus")  // cus_xyz789ghi012
```

---

## HTTP Request/Response Patterns

### Request Validation

```go
type CreateOrderRequest struct {
    Total      float64 `json:"total" binding:"required,gt=0"`
    CustomerID string  `json:"customerId" binding:"required"`
    Items      []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
    var req CreateOrderRequest
    if err := request.ValidBody(c, &req); err != nil {
        response.Error(c, err)
        return
    }

    // Use validated req...
}
```

### Response Mapping

```go
// ✅ GOOD: Map to response DTO
response.SuccessJSON(c, http.StatusOK, ToOrderResponse(order))

// ❌ BAD: Return raw model
response.SuccessJSON(c, http.StatusOK, order)
```

---

## Pointer vs Value

Use pointers for: modifying fields, optional data, large structs, receiver methods that modify.  
Use values for: small structs (<100 bytes), immutable data, read-only ops.

---

## Anti-Patterns

❌ `float64` for money, secrets in code, logging PII, manual transactions, raw SQL in services, magic strings, ignoring errors, local time, context in structs, raw GORM models in responses

---

## Quick Reference

**Error Handling**: Wrap errors, use problem constructors, include `.WithCode()`  
**Context**: Always first parameter, pass to all downstream calls  
**GORM**: Use scopes, constants for relations, `db.Conn(ctx)` for transactions  
**Middleware**: Store in context via `c.Set()`, abort on errors with `c.Abort()`  
**Config**: Constants for keys, defaults in `Configure()`, env var overrides  
**Logging**: Structured logging with `slog`, extract logger from context  
**Money**: Use `decimal.Decimal`, never `float64`  
**Time**: Always UTC, use `time.Time`  
**IDs**: Prefixed IDs for readability  
**Requests**: Use `request.ValidBody()`, validate with struct tags  
**Responses**: Map to DTOs with camelCase, never return raw models
