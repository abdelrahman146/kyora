---
description: Backend Error Handling — Problem JSON (RFC 7807), Domain Errors, Error Wrapping (Reusable)
applyTo: "backend/**"
---

# Backend Error Handling

**Reusable error handling patterns for Go backends.**

Use when: Returning errors from services, handling errors in handlers, creating domain errors.

See also:

- `architecture.instructions.md` — Domain architecture patterns
- `go-patterns.instructions.md` — Go implementation patterns
- `api-contracts.instructions.md` — Response/DTO patterns

---

## Problem JSON (RFC 7807)

All API errors use **Problem JSON** (`Content-Type: application/problem+json`):

```json
{
  "status": 400,
  "title": "Bad Request",
  "detail": "invalid request body",
  "type": "about:blank",
  "instance": "/v1/orders",
  "extensions": {
    "code": "request.invalid_body",
    "field": "email"
  }
}
```

**Required fields:**

- `status`: HTTP status code
- `title`: Human-readable status
- `detail`: Specific error message
- `extensions.code`: Stable machine-readable error code

**Optional fields:**

- `type`: URI reference (defaults to `about:blank`)
- `instance`: Request path (auto-filled)
- `extensions.*`: Additional context

---

## Problem Type

```go
type Problem struct {
    Status     int                    `json:"status"`
    Title      string                 `json:"title"`
    Detail     string                 `json:"detail"`
    Type       string                 `json:"type,omitempty"`
    Instance   string                 `json:"instance,omitempty"`
    Extensions map[string]interface{} `json:"extensions,omitempty"`
    Err        error                  `json:"-"` // Internal only
}
```

**Key methods:**

- `WithCode(code string)` — Set `extensions.code`
- `With(key, value)` — Add to `extensions`
- `WithError(err)` — Attach internal error (not serialized)

---

## Error Constructors

### Standard HTTP Errors

```go
// 400 Bad Request
problem.BadRequest("invalid email format")

// 401 Unauthorized
problem.Unauthorized("invalid credentials")

// 403 Forbidden
problem.Forbidden("insufficient permissions")

// 404 Not Found
problem.NotFound("order not found")

// 409 Conflict
problem.Conflict("email already exists")

// 422 Unprocessable Entity
problem.UnprocessableEntity("validation failed")

// 429 Too Many Requests
problem.TooManyRequests("rate limit exceeded")

// 500 Internal Server Error
problem.InternalServerError("unexpected error")
```

### Domain-Specific Errors

Create in `errors.go`:

```go
package order

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

func InvalidOrderStatus(orderID string, currentStatus, requestedStatus string) *problem.Problem {
    return problem.BadRequest("invalid status transition").
        WithCode("order.invalid_status_transition").
        With("orderId", orderID).
        With("currentStatus", currentStatus).
        With("requestedStatus", requestedStatus)
}
```

**Rules:**

- All errors MUST include `extensions.code`
- Use domain-specific prefixes (e.g., `order.`, `inventory.`, `customer.`)
- Enrich with structured context via `.With(key, value)`
- Never include secrets, tokens, or PII

---

## Error Wrapping

Wrap errors at boundaries to add context:

```go
// Service → Storage boundary
order, err := s.storage.order.FindOne(ctx, s.storage.order.ScopeID(orderID))
if err != nil {
    return nil, problem.Wrap(err, "failed to get order")
}

// External API call
resp, err := stripe.CreateSubscription(params)
if err != nil {
    return nil, problem.Wrap(err, "failed to create Stripe subscription")
}
```

**`problem.Wrap()` behavior:**

- If `err` is already a `*problem.Problem`, returns it unchanged
- Otherwise, wraps in `InternalServerError` with given detail
- Preserves original error via `.WithError(err)` for logging

---

## HTTP Response Layer

Always use `response.Error()` to emit errors:

```go
func (h *Handler) GetOrder(c *gin.Context) {
    orderID := c.Param("orderId")

    order, err := h.service.GetOrder(c.Request.Context(), orderID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.SuccessJSON(c, http.StatusOK, ToOrderResponse(order))
}
```

**`response.Error()` behavior:**

- Accepts `*problem.Problem` or any `error`
- Normalizes non-Problem errors:
  - `gorm.ErrRecordNotFound` → `404 Not Found` (`resource.not_found`)
  - Unique constraint → `409 Conflict` (`resource.conflict`)
  - Everything else → `500 Internal Server Error` (`generic.internal`)
- Sets `instance` from request path if missing
- Logs error with trace ID

**Never do this:**

```go
// ❌ BAD: Ad-hoc error response
c.JSON(400, gin.H{"error": "invalid input"})

// ❌ BAD: Inconsistent error shape
c.JSON(500, map[string]string{"message": err.Error()})
```

---

## Error Code Naming

Use hierarchical naming with dots:

| Pattern                 | Example                | Description            |
| ----------------------- | ---------------------- | ---------------------- |
| `request.*`             | `request.invalid_body` | Request parsing errors |
| `auth.*`                | `auth.invalid_token`   | Authentication errors  |
| `permission.*`          | `permission.denied`    | Authorization errors   |
| `<domain>.*`            | `order.not_found`      | Domain-specific errors |
| `<domain>.<resource>.*` | `order.item.invalid`   | Nested resource errors |
| `generic.*`             | `generic.internal`     | Generic errors         |

**Stability rule:** Once an error code is in production, never change it. Frontend depends on these for i18n.

---

## Request Validation Errors

Use `request.ValidBody()` for JSON validation:

```go
var req CreateOrderRequest
if err := request.ValidBody(c, &req); err != nil {
    response.Error(c, err)
    return
}
```

**Validation errors map to:**

- Invalid JSON: `400` with `request.invalid_body`
- Unknown fields: `400` with `request.invalid_body`
- Missing required fields: `400` with `request.invalid_body`

**For field-level errors (if backend supports):**

```json
{
  "status": 400,
  "title": "Validation Failed",
  "detail": "invalid request body",
  "extensions": {
    "code": "request.validation_failed",
    "invalidFields": [
      { "name": "email", "reason": "invalid format" },
      { "name": "password", "reason": "too short" }
    ]
  }
}
```

---

## Middleware Errors

Middleware should emit errors via `response.Error()`:

```go
func EnforceAuthentication(c *gin.Context) {
    token := extractBearerToken(c)
    if token == "" {
        response.Error(c, problem.Unauthorized("missing authorization header").
            WithCode("auth.missing_token"))
        c.Abort()
        return
    }

    claims, err := validateJWT(token)
    if err != nil {
        response.Error(c, problem.Unauthorized("invalid token").
            WithCode("auth.invalid_token").
            WithError(err))
        c.Abort()
        return
    }

    c.Set(auth.ClaimsKey, claims)
    c.Next()
}
```

---

## Error Logging

Use structured logging with error context:

```go
logger := logger.FromContext(ctx)

order, err := s.storage.order.FindOne(ctx, s.storage.order.ScopeID(orderID))
if err != nil {
    logger.Error("failed to get order",
        slog.String("orderId", orderID),
        slog.Any("error", err),
    )
    return nil, problem.Wrap(err, "failed to get order")
}
```

**Rules:**

- Log errors at point of failure
- Include relevant context (IDs, parameters)
- Never log secrets, tokens, or PII
- Use error levels appropriately:
  - `Error`: Unexpected failures
  - `Warn`: Expected but notable failures (rate limits)
  - `Info`: Normal operation events

---

## Common Error Patterns

### Not Found

```go
order, err := s.storage.order.FindOne(ctx, s.storage.order.ScopeID(orderID))
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, OrderNotFound(orderID)
    }
    return nil, problem.Wrap(err, "failed to get order")
}
```

### Validation Failure

```go
if req.Total <= 0 {
    return nil, problem.BadRequest("total must be positive").
        WithCode("order.invalid_total").
        With("total", req.Total)
}
```

### Permission Denied

```go
if !actor.HasPermission(role.ActionRead, role.ResourceOrder) {
    return nil, problem.Forbidden("insufficient permissions").
        WithCode("permission.denied").
        With("action", role.ActionRead).
        With("resource", role.ResourceOrder)
}
```

### Business Rule Violation

```go
if order.Status == StatusCancelled {
    return nil, problem.BadRequest("cannot modify cancelled order").
        WithCode("order.already_cancelled").
        With("orderId", order.ID)
}
```

### External Service Failure

```go
resp, err := stripe.CreateSubscription(params)
if err != nil {
    return nil, problem.InternalServerError("payment provider error").
        WithCode("billing.stripe_error").
        WithError(err)
}
```

---

## Rate Limiting Errors

```go
if rateLimitExceeded {
    return problem.TooManyRequests("rate limit exceeded").
        WithCode("rate_limit.exceeded").
        With("retryAfterSeconds", retryAfter)
}
```

Frontend can extract `retryAfterSeconds` and show countdown timer.

---

## Database Errors

**GORM error mapping (auto-handled by `response.Error`):**

| GORM Error                  | Status | Code                            |
| --------------------------- | ------ | ------------------------------- |
| `gorm.ErrRecordNotFound`    | 404    | `resource.not_found`            |
| Unique constraint violation | 409    | `resource.conflict`             |
| Foreign key violation       | 400    | `resource.invalid_reference`    |
| Check constraint            | 400    | `resource.constraint_violation` |
| Other DB errors             | 500    | `generic.internal`              |

**Custom handling:**

```go
err := s.storage.order.Create(ctx, order)
if err != nil {
    // Check for specific constraint
    if strings.Contains(err.Error(), "duplicate key") {
        return nil, problem.Conflict("order already exists").
            WithCode("order.duplicate").
            With("orderId", order.ID)
    }
    return nil, problem.Wrap(err, "failed to create order")
}
```

---

## Transaction Errors

Atomic processor handles rollback automatically:

```go
err := s.atomic.Exec(ctx, func(txCtx context.Context) error {
    order, err := s.storage.order.Create(txCtx, order)
    if err != nil {
        return err // Transaction rolls back
    }

    err = s.storage.orderItem.CreateMany(txCtx, items)
    if err != nil {
        return err // Transaction rolls back
    }

    return nil // Transaction commits
})

if err != nil {
    return nil, problem.Wrap(err, "failed to create order")
}
```

---

## Testing Error Handling

### E2E Tests

```go
func (s *Suite) TestCreateOrder_InvalidTotal() {
    payload := map[string]interface{}{
        "total": -10,
    }

    resp, _ := s.client.Post("/v1/orders", payload, testutils.WithAuth(token))

    s.Equal(http.StatusBadRequest, resp.StatusCode)

    var result map[string]interface{}
    testutils.DecodeJSON(resp, &result)

    // Verify error structure
    s.Equal(400, result["status"])
    s.Contains(result, "title")
    s.Contains(result, "detail")

    // Verify error code
    extensions := result["extensions"].(map[string]interface{})
    s.Equal("order.invalid_total", extensions["code"])
}
```

### Unit Tests

```go
func TestCreateOrder_ValidationError(t *testing.T) {
    service := NewService(mockStorage, mockAtomic, mockBus)

    req := &CreateOrderRequest{Total: -10}

    _, err := service.CreateOrder(ctx, actor, biz, req)

    assert.Error(t, err)

    // Check if it's a Problem
    problem, ok := err.(*problem.Problem)
    assert.True(t, ok)
    assert.Equal(t, 400, problem.Status)
    assert.Equal(t, "order.invalid_total", problem.Extensions["code"])
}
```

---

## Anti-Patterns

❌ Inconsistent shapes, ignoring wrapping, missing codes, logging secrets, ad-hoc JSON, generic messages, changing codes, exposing internals

---

## Quick Reference

**Problem JSON**: RFC 7807 (`status`, `detail`, `extensions.code`)  
**Constructors**: `problem.BadRequest()`, `problem.NotFound()`, etc.  
**Domain**: Create in `errors.go`, use `.WithCode()` + `.With()`  
**Wrapping**: `problem.Wrap(err, "context")`  
**Response**: `response.Error(c, err)` only  
**Codes**: `<domain>.<error>`, never change  
**Logging**: Context, no secrets  
**Testing**: Assert status, detail, code
