---
description: Errors & Failure Handling SSOT (backend + portal-web)
applyTo: "backend/internal/platform/types/problem/**,backend/internal/platform/response/**,backend/internal/platform/request/**,backend/internal/platform/auth/**,portal-web/src/api/client.ts,portal-web/src/lib/errorParser.ts,portal-web/src/lib/translateError.ts"
---

# Errors & Failure Handling (SSOT)

**SSOT Hierarchy:**

- Parent: `.github/copilot-instructions.md`
- Backend peers: `.github/instructions/backend-core.instructions.md`, `.github/instructions/go-backend-patterns.instructions.md`
- Portal peers: `.github/instructions/portal-web-architecture.instructions.md`, `.github/instructions/ky.instructions.md`, `.github/instructions/forms.instructions.md`

**When to Read:**

- Adding/modifying API error responses
- Handling request validation failures
- Implementing retries/timeouts/offline UX
- Token refresh / 401 flows
- Showing user-friendly errors in portal-web

---

## 1) Backend error contract (authoritative)

### 1.1 Problem JSON shape

Kyora APIs return errors as **Problem JSON** (`Content-Type: application/problem+json`) using `backend/internal/platform/types/problem/problem.go`.

Canonical shape:

```json
{
  "status": 400,
  "title": "Bad Request",
  "detail": "invalid request body",
  "type": "about:blank",
  "instance": "/v1/...",
  "extensions": {
    "code": "request.invalid_body",
    "field": "email"
  }
}
```

Notes:

- `instance` is set automatically in `response.Error` from the request path when missing.
- `extensions.code` is **required** and must be stable + machine-readable.
- Other `extensions` are optional; domains may add structured context (e.g. `field`, `action`, `resource`, `assetId`).
- Internal errors are attached via `WithError(err)` but **never serialized**.

### 1.2 How errors are emitted

**Always** emit errors via `backend/internal/platform/response/response.go`:

- `response.Error(c, err)` is the single entry point.
- If `err` is not a `*problem.Problem`, it is normalized:
  - DB not found → `404 Not Found` (`resource not found`)
  - DB unique violation → `409 Conflict` (`resource already exists`)
  - everything else → `500 Internal Server Error` (generic detail)

**Code rule (strict):** all errors must include `extensions.code`.

- Domain errors: call `problem.*(...).WithCode("<area>.<reason>")`.
- Normalized errors in `response.Error` already include stable codes:
  - `resource.not_found`
  - `resource.conflict`
  - `generic.internal`

Do not write ad-hoc error responses in handlers (no `c.JSON(status, ...)` for errors).

### 1.3 Request body validation (JSON)

`backend/internal/platform/request/valid_body.go` is the SSOT for JSON body handling:

- Uses `json.Decoder` with `DisallowUnknownFields()`.
- Invalid JSON / missing body / unknown fields → `400 Bad Request` with `detail: "invalid request body"` and `extensions.code: "request.invalid_body"`.
- Oversized bodies (via `http.MaxBytesReader`) → `413 Payload Too Large`.
- Gin struct validation errors (`binding.Validator.ValidateStruct`) currently also map to `400 Bad Request` with the same generic detail.

If you need field-level UX in portal-web, return **structured** details via `Problem.Extensions` (see “Drift” section).

### 1.4 Common status code semantics (backend)

Backend currently uses:

- `400` for invalid JSON, missing required fields, and most validation failures
- `401` for missing/invalid JWT (auth middleware)
- `403` for permission failures (RBAC)
- `404` for scoped resource not found
- `409` for uniqueness/conflict
- `413` for body too large
- `429` for rate-limits (some endpoints add `extensions.retryAfterSeconds`)
- `5xx` for unexpected failures

`422 Unprocessable Entity` exists in `problem` but is not used today.

### 1.5 Middleware-driven failures

Backend failures can happen before handlers run:

- `auth.EnforceAuthentication` → `401` when JWT missing/invalid.
- RBAC helpers typically return `403` with extensions like `action` and `resource`.

Rule: middleware must call `response.Error` and return/abort early.

### 1.6 Logging and privacy

- Never put secrets/tokens/PII into `Problem.Detail` or `Problem.Extensions`.
- Use `.WithError(err)` for internal diagnostics; it is not serialized.
- Prefer stable, user-friendly `detail` strings (portal uses them as fallback text).

---

## 2) Portal-web error handling (alignment guidance)

### 2.1 Single HTTP client (SSOT)

Portal-web must use `portal-web/src/api/client.ts` (`apiClient`) and its typed helpers.

Key behavior:

- Adds `Authorization: Bearer <token>` automatically (if present).
- Adds `X-Request-ID` automatically.
- Retries are enabled only for **idempotent** methods (`GET/PUT/DELETE/...`) on `[408, 413, 429, 5xx]`.
- `401` handling:
  - For non-auth endpoints, the client attempts refresh (`POST /v1/auth/refresh`) and retries the original request.
  - If refresh fails, tokens are cleared and the browser is redirected to `/auth/login`.
  - Auth endpoints under `/v1/auth/*` are excluded from refresh retry logic.

Do not re-implement refresh logic in feature code.

### 2.2 Translating and showing errors

SSOT utilities:

- Parser: `portal-web/src/lib/errorParser.ts` (`parseProblemDetails`)
- Translator: `portal-web/src/lib/translateError.ts` (`translateErrorAsync`)

UI pattern:

- In `catch`, call `translateErrorAsync(error, t)` and show a toast or inline message.
- Do not show raw backend errors directly; always go through i18n.

**Mapping rule (strict):** portal-web maps backend codes to i18n keys as:

- `extensions.code = "customer.not_found"` → `errors.backend.customer.not_found`

This mapping is implemented in `portal-web/src/lib/errorParser.ts`.

### 2.3 Failure scenarios to handle explicitly

Handle these consistently in UI:

- Network offline / DNS / CORS: treat as connection error.
- Timeouts: show timeout message; allow user retry.
- `401` after refresh attempt: user must re-login (client redirects).
- `403`: show “no permission” UX; do not offer actions that will always fail.
- `409`: show conflict message (e.g. descriptor taken).
- `429`: if `retryAfterSeconds` exists in Problem `extensions`, surface it to the user.

### 2.4 Form validation errors (current behavior)

Today:

- Backend often returns `400` with generic `detail`.
- Some domain services include `extensions.field` for a single field.

Frontend guidance:

- Prefer showing a toast for generic failures.
- If `extensions.field` is present, you may focus the field and/or show a field-level error message.
- Do not assume a full `{ fieldName: message }` map exists unless the backend explicitly provides it.

---

## 3) Known drift / gaps (log to DRIFT_TODO.md when touching)

These are current backend↔portal misalignments that affect error UX:

- Portal includes `parseValidationErrors()` expecting `extensions.errors | validationErrors | fieldErrors`, but backend does not emit a full field→message map today.
- Backend maps Gin struct validation errors to a generic `400 invalid request body` (code: `request.invalid_body`) without exposing which fields failed.

If you change either side, keep this document updated and log new drift in `DRIFT_TODO.md`.
