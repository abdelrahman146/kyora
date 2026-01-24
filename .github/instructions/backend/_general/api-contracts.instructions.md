---
description: API Contracts — Response DTOs, Swagger/OpenAPI, camelCase Standards (Reusable)
applyTo: "backend/**"
---

# API Contracts

**Response DTOs, Swagger/OpenAPI generation, and API contract standards.**

Use when: Creating/modifying API endpoints, defining response shapes, generating OpenAPI docs.

See also:

- `architecture.instructions.md` — Domain architecture patterns
- `errors.instructions.md` — Error response patterns
- `go-patterns.instructions.md` — Go implementation patterns

---

## Core Rule

**Backend is the source of truth for API contracts.**

- Never return raw GORM models
- Always use explicit response DTOs
- All JSON fields MUST be camelCase
- Generate OpenAPI from code annotations

---

## Response DTO Pattern

### File Organization

Per domain, keep:

- Request DTOs: `model_request.go` (or `model.go` for simple domains)
- Response DTOs: `model_response.go` (required)

### Response DTO Structure

```go
// model_response.go
package order

type OrderResponse struct {
    ID         string    `json:"id"`
    BusinessID string    `json:"businessId"`     // camelCase!
    Total      string    `json:"total"`
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"createdAt"`      // camelCase!
    UpdatedAt  time.Time `json:"updatedAt"`      // camelCase!
}

func ToOrderResponse(o *Order) *OrderResponse {
    return &OrderResponse{
        ID:         o.ID,
        BusinessID: o.BusinessID,
        Total:      o.Total.String(),
        Status:     string(o.Status),
        CreatedAt:  o.CreatedAt,
        UpdatedAt:  o.UpdatedAt,
    }
}

func ToOrderResponses(orders []*Order) []*OrderResponse {
    responses := make([]*OrderResponse, len(orders))
    for i, o := range orders {
        responses[i] = ToOrderResponse(o)
    }
    return responses
}
```

**Rules:**

- ✅ Use explicit JSON tags with camelCase
- ✅ Convert `decimal.Decimal` to `string` for precision
- ✅ Convert enums to strings
- ✅ Omit GORM internals (`gorm.Model`, `DeletedAt` unless needed)
- ✅ Use `time.Time` for timestamps (auto-serialized as RFC3339)
- ❌ Never embed GORM models
- ❌ Never use PascalCase JSON tags

---

## Handler Pattern

Always map models to DTOs before responding:

```go
// ✅ CORRECT: Map to DTO
func (h *Handler) GetOrder(c *gin.Context) {
    orderID := c.Param("orderId")

    order, err := h.service.GetOrder(c.Request.Context(), orderID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.SuccessJSON(c, http.StatusOK, ToOrderResponse(order))
}

// ❌ WRONG: Return raw model
func (h *Handler) GetOrder(c *gin.Context) {
    order, err := h.service.GetOrder(c.Request.Context(), orderID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.SuccessJSON(c, http.StatusOK, order) // Leaks GORM internals!
}

// ❌ WRONG: Wrap raw model in gin.H
response.SuccessJSON(c, http.StatusOK, gin.H{"order": order}) // Still leaks GORM!
```

---

## Nested Relations

### Explicit Nested DTOs

If endpoint returns nested relations, define them explicitly:

```go
type OrderResponse struct {
    ID    string              `json:"id"`
    Total string              `json:"total"`
    Items []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
    ID        string `json:"id"`
    ProductID string `json:"productId"`
    Quantity  int    `json:"quantity"`
    Price     string `json:"price"`
}

func ToOrderResponse(o *Order) *OrderResponse {
    return &OrderResponse{
        ID:    o.ID,
        Total: o.Total.String(),
        Items: ToOrderItemResponses(o.Items),
    }
}
```

**Rules:**

- ✅ Preload relations in service layer
- ✅ Define explicit nested DTOs
- ✅ Convert nested models to DTOs
- ❌ Don't assume relations exist without preloading
- ❌ Don't return relation IDs without preloading full objects

### Optional Nested Relations

Use pointers for optional nested data:

```go
type OrderResponse struct {
    ID       string             `json:"id"`
    Customer *CustomerResponse  `json:"customer,omitempty"`
}

func ToOrderResponse(o *Order, includeCustomer bool) *OrderResponse {
    resp := &OrderResponse{ID: o.ID}

    if includeCustomer && o.Customer != nil {
        resp.Customer = ToCustomerResponse(o.Customer)
    }

    return resp
}
```

---

## List Response Pattern

Use `list.ListResponse[T]` for paginated endpoints:

```go
type ListResponse[T any] struct {
    Items      []T   `json:"items"`
    TotalCount int64 `json:"totalCount"`
    Page       int   `json:"page"`
    PageSize   int   `json:"pageSize"`
    TotalPages int   `json:"totalPages"`
    HasMore    bool  `json:"hasMore"`
}

// Handler
func (h *Handler) ListOrders(c *gin.Context) {
    req := list.ParseListRequest(c)

    items, count, err := h.service.ListOrders(c.Request.Context(), req)
    if err != nil {
        response.Error(c, err)
        return
    }

    resp := list.NewListResponse(
        ToOrderResponses(items),
        count,
        req.Page(),
        req.PageSize(),
    )

    response.SuccessJSON(c, http.StatusOK, resp)
}
```

**Rules:**

- ✅ Use consistent pagination structure
- ✅ Include `totalCount`, `page`, `pageSize`, `totalPages`, `hasMore`
- ✅ Map items to DTOs before wrapping
- ❌ Don't invent custom pagination shapes per endpoint

---

## Casing Standards

### JSON Field Names (CRITICAL)

**All JSON fields MUST be camelCase:**

| Model Field  | JSON Tag     | ❌ Wrong                    |
| ------------ | ------------ | --------------------------- |
| `CreatedAt`  | `createdAt`  | `CreatedAt`, `created_at`   |
| `BusinessID` | `businessId` | `BusinessID`, `business_id` |
| `TotalCount` | `totalCount` | `TotalCount`, `total_count` |

**Why this matters:**

- Frontend expects camelCase (TypeScript/JavaScript convention)
- Inconsistent casing causes type errors and runtime bugs
- GORM models use PascalCase; DTOs must convert

### Response DTO Checklist

When creating a response DTO:

✅ All JSON tags are camelCase  
✅ No embedded GORM models (`gorm.Model`)  
✅ Timestamps converted: `CreatedAt` → `createdAt`  
✅ Foreign keys converted: `BusinessID` → `businessId`  
✅ Relations omitted with `json:"-"` or mapped to nested DTOs  
✅ Enums converted to strings  
✅ Decimals converted to strings

---

## Swagger/OpenAPI Generation

### Annotations

Use Swaggo annotations in handlers:

```go
// @Summary Create order
// @Description Create a new order for a business
// @Tags orders
// @Accept json
// @Produce json
// @Param businessDescriptor path string true "Business descriptor"
// @Param request body CreateOrderRequest true "Order details"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} problem.Problem
// @Failure 401 {object} problem.Problem
// @Failure 403 {object} problem.Problem
// @Router /v1/businesses/{businessDescriptor}/orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
    // ...
}
```

**Rules:**

- ✅ Reference response DTOs (not GORM models)
- ✅ Document all path/query/body parameters
- ✅ Document all success and error responses
- ✅ Use `@Router` with correct path and method
- ❌ Don't reference GORM models in `@Success`

### Generate OpenAPI

```bash
# Generate swagger docs
make openapi

# Verify no uncommitted changes
make openapi.check
```

**Output:**

- `backend/docs/swagger.json`
- `backend/docs/swagger.yaml`

**CI Rule**: OpenAPI must be up-to-date (checked in CI).

---

## Timestamp Handling

### Server-Side (Backend)

Always use UTC:

```go
order := &Order{
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
}
```

### Client-Side (Frontend)

Backend returns timestamps as RFC3339 strings:

```json
{
  "createdAt": "2026-01-24T10:30:00Z",
  "updatedAt": "2026-01-24T10:30:00Z"
}
```

Frontend converts to local timezone for display.

---

## Money Fields

Always use `decimal.Decimal` in models, convert to `string` in DTOs:

```go
type OrderResponse struct {
    Total string `json:"total"`
}

func ToOrderResponse(o *Order) *OrderResponse {
    return &OrderResponse{
        Total: o.Total.String(), // "100.50"
    }
}
```

**Rules:**

- ✅ Use `string` for money in DTOs (preserves precision)
- ✅ Frontend parses as string or BigDecimal
- ❌ Never use `float64` for money

---

## Enum Fields

Convert enums to strings:

```go
type OrderStatus string

const (
    StatusPending   OrderStatus = "pending"
    StatusCompleted OrderStatus = "completed"
)

type OrderResponse struct {
    Status string `json:"status"`
}

func ToOrderResponse(o *Order) *OrderResponse {
    return &OrderResponse{
        Status: string(o.Status),
    }
}
```

---

## Optional Fields

Use pointers for optional fields:

```go
type OrderResponse struct {
    ID          string  `json:"id"`
    Notes       *string `json:"notes,omitempty"`
    CompletedAt *time.Time `json:"completedAt,omitempty"`
}

func ToOrderResponse(o *Order) *OrderResponse {
    resp := &OrderResponse{ID: o.ID}

    if o.Notes != "" {
        resp.Notes = &o.Notes
    }

    if !o.CompletedAt.IsZero() {
        resp.CompletedAt = &o.CompletedAt
    }

    return resp
}
```

---

## Error Responses

See `errors.instructions.md` for full details.

All errors return Problem JSON:

```json
{
  "status": 400,
  "title": "Bad Request",
  "detail": "total must be positive",
  "extensions": {
    "code": "order.invalid_total"
  }
}
```

Document in Swagger:

```go
// @Failure 400 {object} problem.Problem
// @Failure 404 {object} problem.Problem
```

---

## Versioning

Use URL path versioning:

- ✅ `/v1/orders`
- ❌ `/orders?version=1`
- ❌ `Accept: application/vnd.api+json; version=1`

**Breaking changes require new version**: `/v2/orders`

---

## Testing Response Contracts

### E2E Tests

```go
func (s *Suite) TestCreateOrder_ResponseShape() {
    resp, _ := s.client.Post("/v1/orders", payload, testutils.WithAuth(token))

    s.Equal(http.StatusCreated, resp.StatusCode)

    var result OrderResponse
    s.NoError(testutils.DecodeJSON(resp, &result))

    // Verify all fields present
    s.NotEmpty(result.ID)
    s.NotEmpty(result.Total)
    s.Equal("pending", result.Status)
    s.NotZero(result.CreatedAt)

    // Verify casing (this catches PascalCase leaks)
    var raw map[string]interface{}
    s.NoError(testutils.DecodeJSON(resp, &raw))
    s.Contains(raw, "createdAt") // Not "CreatedAt"
    s.Contains(raw, "businessId") // Not "BusinessID"
}
```

---

## Frontend Integration

### TypeScript Types

Frontend should generate types from OpenAPI or create manual interfaces:

```typescript
interface OrderResponse {
  id: string;
  businessId: string;
  total: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}
```

**Critical**: Frontend types MUST use camelCase to match backend DTOs.

---

## Anti-Patterns

❌ Raw GORM models, wrapping in `gin.H{}`, embedded GORM, PascalCase JSON, `float64` for money, GORM refs in Swagger, inconsistent pagination, missing OpenAPI regen

---

## Checklist

✅ Request DTO in `model_request.go`, response DTO in `model_response.go` (camelCase), converter functions, handler maps before responding, Swagger annotations, `make openapi`, E2E tests, frontend types updated

---

## Quick Reference

**DTOs**: Explicit response structs in `model_response.go`, camelCase JSON tags  
**Converters**: `To<Model>Response()` functions  
**Handlers**: Always map models to DTOs before `response.SuccessJSON()`  
**Lists**: Use `list.ListResponse[T]` for pagination  
**Nested**: Define explicit nested DTOs, preload in service  
**Money**: Convert `decimal.Decimal` to `string`  
**Enums**: Convert to strings  
**Timestamps**: RFC3339 strings (auto-serialized)  
**Swagger**: Annotate handlers, reference DTOs, regenerate with `make openapi`  
**Testing**: Verify all fields, check camelCase in raw JSON
