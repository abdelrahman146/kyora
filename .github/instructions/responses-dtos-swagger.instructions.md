---
description: Responses, DTOs & Swagger Generation SSOT (backend + portal-web)
applyTo: "backend/internal/platform/response/**,backend/internal/platform/types/list/**,backend/internal/domain/**/handler_http.go,backend/internal/domain/**/model*.go,backend/docs/swagger.json,backend/docs/swagger.yaml,backend/main.go,Makefile,portal-web/src/api/**,portal-web/src/api/types/**"
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

### 2.1 Don’t return GORM models directly

Many Kyora domain models embed `gorm.Model`. When a handler returns that struct directly, Go’s `encoding/json` will serialize embedded fields as **PascalCase** (`CreatedAt`, `UpdatedAt`, etc.).

This causes:

- Broken response standards (casing)
- Swagger drift (OpenAPI includes the wrong fields)
- Portal types becoming inconsistent or duplicative

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

---

## 5) Known systemic drift areas

- Endpoints returning GORM models embed `gorm.Model` → PascalCase timestamps leak.
- Swagger generated from those models encodes the wrong field casing, so portal types drift.
- Relations (`Business`, `Order.Notes`, `OrderItem.Product`) are declared on models but are only present when explicitly preloaded and returned.

When you discover a new example, log it in `DRIFT_TODO.md` with the endpoint + expected vs actual response shape.
