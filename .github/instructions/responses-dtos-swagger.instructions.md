---
description: Responses, DTOs & Swagger Generation SSOT (backend + portal-web)
---

# Responses, DTOs & Swagger (SSOT)

**SSOT Hierarchy:**

- Parent: `.github/copilot-instructions.md`
- Backend peers: `.github/instructions/backend-core.instructions.md`, `.github/instructions/go-backend-patterns.instructions.md`
- Portal peers: `.github/instructions/portal-web-architecture.instructions.md`, `.github/instructions/ky.instructions.md`
- Portal HTTP + TanStack Query: `.github/instructions/http-tanstack-query.instructions.md`

**When to Read:**

- Adding a new endpoint or changing response shape
- Fixing “frontend expected nested object but got null/missing”
- Fixing timestamp casing (`CreatedAt` vs `createdAt`)
- Updating OpenAPI/Swagger (`backend/docs/swagger.{json,yaml}`)
- Implementing/adjusting portal query/mutation hooks that consume these responses

---

## 0) The core rule (SSOT)

**Backend is the source of truth for response contracts.**

- Portal-web must model responses based on the backend contract (ideally regenerated OpenAPI).
- Backend must _not_ leak storage/ORM implementation details (GORM internals) to clients.

---

## 1) Backend response contract

### 1.1 Always respond via the shared response layer

- Success: `response.SuccessJSON(c, status, payload)`
- Errors: `response.Error(c, err)` (see `.github/instructions/errors-handling.instructions.md`)

Do not hand-roll success envelope formats; most endpoints return the payload directly.

### 1.2 Pagination/list responses

Kyora’s list responses use `backend/internal/platform/types/list/ListResponse[T]`:

```json
{
  "items": [],
  "totalCount": 123,
  "page": 1,
  "pageSize": 20,
  "totalPages": 7,
  "hasMore": true
}
```

Portal-web should reuse this shape rather than inventing a competing pagination contract.

### 1.3 Response naming + casing standards

**All JSON fields returned to clients must be `camelCase`.**

- ✅ `createdAt`, `updatedAt`, `deletedAt`
- ❌ `CreatedAt`, `UpdatedAt`, `DeletedAt`

Rule: do not expose GORM’s embedded `gorm.Model` fields directly to JSON.

---

## 2) DTO layer (what to do, and what NOT to do)

### 2.1 Don't return GORM models directly (CRITICAL RULE)

**MANDATORY:** Never return a GORM model directly in an HTTP response.

Many Kyora domain models embed `gorm.Model`. When a handler returns that struct directly, Go's `encoding/json` will serialize embedded fields as **PascalCase** (`CreatedAt`, `UpdatedAt`, etc.).

This causes:

- **Broken response standards** (casing violates API contract)
- **Swagger drift** (OpenAPI includes the wrong fields)
- **Portal types becoming inconsistent** or duplicative
- **GORM internals leak** to clients (violates abstraction)

**The Rule:**

- **Always create an explicit response DTO** for each model that is returned to clients
- **Always use `To<Model>Response()` converter functions** in handlers
- **Never wrap raw models** in `gin.H{}` or other envelopes
- **Every response DTO must use camelCase** JSON tags exclusively

**Anti-Pattern Examples (DO NOT DO):**

```go
// ❌ WRONG: Returns raw model (leaks GORM internals with PascalCase)
func (h *Handler) GetUser(c *gin.Context) {
    user := h.service.GetUser(id)
    response.SuccessJSON(c, http.StatusOK, user)  // BAD: user has CreatedAt, UpdatedAt
}

// ❌ WRONG: Wraps raw model in gin.H (still leaks GORM)
func (h *Handler) CompleteOnboarding(c *gin.Context) {
    user := h.service.CreateUser(...)
    response.SuccessJSON(c, http.StatusOK, gin.H{"user": user})  // BAD: user has CreatedAt, UpdatedAt
}

// ❌ WRONG: Uses a response struct that directly embeds the model
type UserResponse struct {
    *User  // BAD: embeds all GORM fields including CreatedAt
}
```

**Correct Pattern:**

```go
// ✅ CORRECT: Explicit response DTO with camelCase
type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`      // camelCase
    UpdatedAt time.Time `json:"updatedAt"`      // camelCase
}

// ✅ CORRECT: Converter function maps model to DTO
func ToUserResponse(user *User) *UserResponse {
    return &UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}

// ✅ CORRECT: Handler uses converter
func (h *Handler) GetUser(c *gin.Context) {
    user := h.service.GetUser(id)
    response.SuccessJSON(c, http.StatusOK, ToUserResponse(user))  // GOOD
}

// ✅ CORRECT: Even in complex responses like login
func (h *Handler) Login(c *gin.Context) {
    user := h.service.Login(email, password)
    response.SuccessJSON(c, http.StatusOK, ToLoginResponse(user, token, refreshToken))
}
```

**Policy:** handlers should return explicit response DTOs.

### 2.2 DTOs must be explicit about nested data

GORM preloading only affects what the server has in memory.

**Clients only get what you serialize.**

- If an endpoint intends to return nested objects (e.g., `order.items`, `order.notes`, `customer.addresses`), the handler must:
  - preload the relation in storage/service, **and**
  - map it into the DTO explicitly.

If nested objects are not intended, return only IDs and omit the relation fields entirely.

### 2.3 Recommended backend structure

Per domain, keep:

- request DTOs in the domain model file (already common): `CreateXRequest`, `UpdateXRequest`
- response types in a dedicated file:
  - `backend/internal/domain/<domain>/model_response.go`

Mapping patterns:

- `func ToXResponse(m *X) XResponse`
- `func ToXResponses(items []*X) []XResponse`

Naming rule:

- Use `XResponse` types (no `DTO` suffix).

DTO fields must:

- use explicit JSON tags
- exclude internal-only fields
- avoid exposing embedded pointers that can recurse (e.g. `OrderNote.Order`)

### 2.4 Timestamps and soft-delete fields

- Prefer returning `createdAt`/`updatedAt`.
- Only include `deletedAt` when the endpoint explicitly needs it (e.g. admin views). Otherwise omit it.
- Use RFC3339 timestamps (Go `time.Time` JSON encoding) consistently.

---

## 3) Swagger/OpenAPI generation (backend)

### 3.1 How it’s generated

Swagger is generated via Swaggo from Go annotations in handlers.

Command (repo SSOT): `make openapi`

This runs:

- `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init`
- entrypoint: `backend/main.go`
- output: `backend/docs/swagger.json` and `backend/docs/swagger.yaml`

### 3.2 Swagger must describe what the endpoint actually returns

Because portal-web often types based on Swagger:

- `@Success` should reference the **response type** (`<domain>.<XResponse>`), not the GORM model.
- If the handler returns `list.ListResponse[customer.CustomerResponse]`, Swagger should reflect that.

Rule: if you change a response shape, regenerate OpenAPI (`make openapi`) and update portal types/schemas.

---

## 4) Portal-web alignment guidance

### 4.1 Don’t assume nested objects exist

If portal-web needs nested objects (e.g. `order.notes`), ensure the backend endpoint explicitly returns them.

- If backend returns IDs only, portal should not type nested objects as present.
- Prefer marking optional nested relations as `?` and only using them in UI when present.

### 4.2 Prefer Zod schemas (SSOT inside portal)

Portal should prefer Zod schemas in `portal-web/src/api/types/**` for response validation and inference.

- If a backend endpoint is known to be inconsistent today, model it explicitly (with optional fields) and log the drift.

### 4.3 Timestamp casing

Portal must target the backend standard (`createdAt`, `updatedAt`).

If an endpoint currently returns PascalCase timestamps due to leaked `gorm.Model`, do not normalize silently in random components. Log drift and fix backend to return DTOs.

### 4.4 Portal-web types MUST use camelCase to match backend responses (CRITICAL)

**Rule:** All TypeScript interfaces in `portal-web/src/api/` must use camelCase field names to match backend JSON responses.

**MANDATORY:**

- ✅ `pageSize`, `totalCount`, `totalPages`, `hasMore` (ListResponse fields)
- ✅ `productId`, `businessId`, `categoryId` (foreign keys)
- ✅ `createdAt`, `updatedAt`, `deletedAt` (timestamps)
- ❌ `page_size`, `total_count`, `total_pages`, `has_more` (snake_case)
- ❌ `product_id`, `business_id` (snake_case)

**Why this matters:**

When portal-web types use snake_case, TypeScript will fail at runtime when accessing the actual camelCase data from backend responses. Example:

```typescript
// ❌ WRONG: Type defines snake_case
interface ListResponse {
  page_size: number; // Type says page_size
}

// But backend sends camelCase
const response = await api.list(); // { pageSize: 20, ... }
console.log(response.page_size); // undefined! ❌

// ✅ CORRECT: Type matches backend
interface ListResponse {
  pageSize: number; // Type matches backend
}
const response = await api.list(); // { pageSize: 20, ... }
console.log(response.pageSize); // 20 ✅
```

**When adding new types:**

1. Check backend `model_response.go` or OpenAPI schema
2. Copy field names and casing exactly
3. Use camelCase consistently
4. Verify with TypeScript type check (`npm run type-check`)

---

## 5) Known systemic drift areas

- Endpoints returning GORM models embed `gorm.Model` → PascalCase timestamps leak.
- Swagger generated from those models encodes the wrong field casing, so portal types drift.
- Relations (`Business`, `Order.Notes`, `OrderItem.Product`) are declared on models but are only present when explicitly preloaded and returned.

When you discover a new example, log it in `DRIFT_TODO.md` with the endpoint + expected vs actual response shape.
