---
description: Errors & Failure Handling SSOT (backend + portal-web)
applyTo: "**/*"
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

### 2.2 Global React Query error handling (SSOT)

Portal-web uses **global error handlers** configured via `QueryClient` with `QueryCache.onError` and `MutationCache.onError` callbacks. This is the **default and recommended** way to handle backend errors.

**Implementation:** `portal-web/src/main.tsx`

```typescript
const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      if (shouldIgnoreGlobalError(error)) return;
      const meta = query.meta as undefined | { errorToast?: "global" | "off" };
      if (meta?.errorToast === "off") return;

      // Dedupe by queryHash for 30s to prevent toast spam on background refetches
      const now = Date.now();
      const last = queryErrorToastDeduper.get(query.queryHash) ?? 0;
      if (now - last < QUERY_ERROR_TOAST_DEDUPE_MS) return;
      queryErrorToastDeduper.set(query.queryHash, now);

      void showErrorFromException(error, i18n.t);
    },
  }),
  mutationCache: new MutationCache({
    onError: (error, _variables, _context, mutation) => {
      if (shouldIgnoreGlobalError(error)) return;
      const meta = mutation.meta as
        | undefined
        | { errorToast?: "global" | "off" };
      if (meta?.errorToast === "off") return;

      void showErrorFromException(error, i18n.t);
    },
  }),
});
```

**Key behaviors:**

1. **Automatic error translation:** All backend errors are translated via `translateErrorAsync` and shown as toast notifications
2. **Query error deduplication:** Query errors are deduplicated by `queryHash` for 30 seconds to prevent toast spam on background refetches
3. **Mutation errors always show:** Mutation errors show toast immediately (no deduplication)
4. **AbortError filtering:** `shouldIgnoreGlobalError()` filters out `AbortError` and cancel errors
5. **Opt-out mechanism:** Use `meta: { errorToast: 'off' }` to suppress global toast for special cases

**When to use global handler (default):**

✅ All mutations and queries by default  
✅ Background refetches  
✅ User-initiated actions (create/update/delete)

**When to opt out (`meta: { errorToast: 'off' }`):**

- Inline form validation where errors should appear next to fields
- Silent background operations where errors shouldn't interrupt user
- Custom error handling that requires different UX (e.g., redirect on auth failure)

**Example: Relying on global handler (recommended):**

```tsx
// In API client (portal-web/src/api/customer.ts)
export function useUpdateCustomerMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Customer, Error, UpdateCustomerRequest>,
) {
  return useMutation({
    mutationFn: (data) => customerApi.update(businessDescriptor, data),
    // No onError needed - global handler shows toast automatically
    ...options,
  });
}

// In component
const updateMutation = useUpdateCustomerMutation(businessDescriptor, {
  onSuccess: (updated) => {
    showSuccessToast(t("customers.update_success"));
    onClose();
  },
  // Error toast shown automatically by global handler
});
```

**Example: Opting out for inline errors:**

```tsx
// In API client (portal-web/src/api/auth.ts)
export function useLoginMutation(
  options?: UseMutationOptions<LoginResponse, Error, LoginRequest>,
) {
  return useMutation({
    mutationFn: (data) => authApi.login(data),
    meta: { errorToast: "off" }, // Suppress global toast
    ...options,
  });
}

// In component
const loginMutation = useLoginMutation({
  onError: async (error) => {
    // Custom inline error display
    const translated = await translateErrorAsync(error, t);
    setFormError(translated);
  },
});
```

### 2.3 Error translation and display utilities

**SSOT utilities:**

- Parser: `portal-web/src/lib/errorParser.ts` (`parseProblemDetails`)
- Translator: `portal-web/src/lib/translateError.ts` (`translateErrorAsync`)
- Toast: `portal-web/src/lib/toast.ts` (`showErrorFromException`, `showErrorToast`)

**UI pattern:**

- **Default:** Rely on global handler (no manual error handling needed)
- **Custom handling:** Call `translateErrorAsync(error, t)` and show inline message
- Do not show raw backend errors directly; always go through i18n

**Mapping rule (strict):** portal-web maps backend codes to i18n keys as:

- `extensions.code = "customer.not_found"` → `errors.backend.customer.not_found`

This mapping is implemented in `portal-web/src/lib/errorParser.ts`.

### 2.4 Explicit error handling scenarios

These scenarios require **explicit handling beyond the global handler**:

Handle these consistently in UI:

- **Network offline / DNS / CORS:** treat as connection error (global handler shows translated toast).
- **Timeouts:** show timeout message; allow user retry (global handler shows translated toast).
- **`401` after refresh attempt:** user must re-login (client redirects automatically in `portal-web/src/api/client.ts`).
- **`403`:** show "no permission" UX; do not offer actions that will always fail (global handler shows toast, but UI should prevent the action).
- **`409`:** show conflict message (e.g. descriptor taken) (global handler shows translated toast).
- **`429`:** if `retryAfterSeconds` exists in Problem `extensions`, surface it to the user (global handler shows toast, custom handling can extract retry time).

**Inline form validation errors (opt-out scenario):**

When backend returns field-level errors (rare), you may need to opt out of global handler and show errors inline:

```tsx
// If backend includes extensions.field for a single field
const mutation = useMutation({
  meta: { errorToast: "off" },
  onError: async (error) => {
    const problem = parseProblemDetails(error);
    if (problem?.extensions?.field) {
      // Show field-level error
      form.setFieldError(
        problem.extensions.field,
        await translateErrorAsync(error, t),
      );
    } else {
      // Fall back to generic error display
      showErrorToast(await translateErrorAsync(error, t));
    }
  },
});
```

**Loading states during error recovery:**

When a mutation fails and user retries:

- Show loading state: `mutation.isPending`
- Disable retry button: `disabled={mutation.isPending}`
- Clear previous error on retry: `mutation.reset()`

**Query error states (less common with global handler):**

Queries with `enabled: false` or special UX needs may require inline error display:

```tsx
const query = useQuery({
  // ... query config
  meta: { errorToast: "off" }, // Suppress global toast
});

// In JSX
{
  query.isError && (
    <div className="alert alert-error">
      <span>{await translateErrorAsync(query.error, t)}</span>
    </div>
  );
}
```

### 2.5 Current portal behavior and known gaps

**Implemented (production-ready):**

- ✅ Global error handling via QueryCache/MutationCache onError
- ✅ Automatic error translation (`translateErrorAsync`)
- ✅ Toast notifications for all errors by default
- ✅ Opt-out mechanism (`meta: { errorToast: 'off' }`)
- ✅ Query error deduplication (30s window by queryHash)
- ✅ AbortError filtering
- ✅ 401 auto-refresh in ky client (`portal-web/src/api/client.ts`)

**Current gaps/inconsistencies:**

- Some components bypass global handler and show `mutation.error.message` directly (see drift report: `backlog/drifts/2026-01-18-onboarding-components-bypass-global-error-handler.md`)
- Backend validation errors don't include full field→message map; most return generic `400` with a single `extensions.field`
- Not all backend error codes have i18n translations (missing keys fall back to generic messages)

**When touching error handling:**

- Default: rely on global handler (no manual error display)
- Special case: use `meta: { errorToast: 'off' }` + custom inline display only when UX specifically requires it
- Always translate errors before displaying (`translateErrorAsync`)
- Never show raw `error.message` or backend `detail` without translation

---

## 3) Known drift / gaps (log to DRIFT_TODO.md when touching)

These are current backend↔portal misalignments that affect error UX:

- Portal includes `parseValidationErrors()` expecting `extensions.errors | validationErrors | fieldErrors`, but backend does not emit a full field→message map today.
- Backend maps Gin struct validation errors to a generic `400 invalid request body` (code: `request.invalid_body`) without exposing which fields failed.
- **Onboarding components bypass global error handler:** Several onboarding pages manually display `mutation.error.message` instead of relying on the global QueryCache/MutationCache error handlers (see drift report: `backlog/drifts/2026-01-18-onboarding-components-bypass-global-error-handler.md`).

If you change either side, keep this document updated and log new drift in `DRIFT_TODO.md`.
