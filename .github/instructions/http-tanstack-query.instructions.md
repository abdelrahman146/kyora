---
description: Portal Web HTTP layer + TanStack Query SSOT (queries, mutations, global error handling)
applyTo: "portal-web/src/api/**,portal-web/src/lib/queryKeys.ts,portal-web/src/lib/queryInvalidation.ts,portal-web/src/main.tsx,portal-web/src/router.tsx,portal-web/src/routes/**,portal-web/src/components/**"
---

# Portal Web HTTP + TanStack Query (SSOT)

**SSOT Hierarchy**

- Parent: `.github/copilot-instructions.md`
- Backend error contract: `.github/instructions/errors-handling.instructions.md`
- Backend response/DTO contract + Swagger workflow: `.github/instructions/responses-dtos-swagger.instructions.md`
- Portal HTTP client details (Ky wrapper): `.github/instructions/ky.instructions.md`

This doc is the SSOT for how `portal-web/**` must call the backend.

---

## 0) Non-negotiables

1. **No direct backend calls from UI code.**

- Routes and components (`portal-web/src/routes/**`, `portal-web/src/components/**`) must not call `fetch`, `ky`, or `apiClient` directly.
- All reads/writes must go through **TanStack Query** (`useQuery`/`useMutation`/`useSuspenseQuery`) using the `portal-web/src/api/**` modules.

2. **Backend contracts are authoritative.**

- Response shape, casing, and error format are backend SSOT.
- Portal must not silently “normalize” inconsistent backend responses in random UI code.
- If backend is inconsistent, **log drift in `DRIFT_TODO.md`** and fix at the contract layer.

3. **All user-facing errors must be meaningful + translated.**

- Backend errors are surfaced as Problem JSON and parsed by portal.
- Default: show a translated message (see `portal-web/src/lib/errorParser.ts` and `portal-web/src/lib/translateError.ts`).

---

## 1) Current portal layering (what exists today)

### 1.1 HTTP client

- Ky instance: `portal-web/src/api/client.ts` exports:
  - `apiClient` (ky instance with auth + retry + refresh flow)
  - typed helpers: `get/post/postVoid/put/patch/del/delVoid`

Important behavior:

- **401 refresh flow** is implemented in Ky `afterResponse`, excluding `/v1/auth/*`.
- **Problem JSON parsing** is done in Ky `beforeError` by calling `parseProblemDetails()`.
  - `error.message` may be set to a safe fallback string, but user-facing messages must still go through `translateErrorAsync` (single source of truth).

### 1.2 TanStack Query integration

- QueryClient is created in `portal-web/src/main.tsx` and provided via `<QueryClientProvider>`.
- TanStack Router context includes `queryClient` (`portal-web/src/router.tsx`) so route loaders can prefetch if/when used.

### 1.3 Query organization pattern

Most API modules follow this structure:

- `*Api`: plain async functions that call `get/post/patch/...`.
- `*Queries`: `queryOptions({ queryKey, queryFn, staleTime, enabled })` factories.
- `use*Query` hooks that call `useQuery(*Queries.*)`.
- `use*Mutation` hooks that call `useMutation({ mutationFn, onSuccess, onError })`.

Query keys are centralized under `portal-web/src/lib/queryKeys.ts`.

---

## 2) Required SSOT pattern going forward

### 2.1 Where network calls are allowed

Allowed:

- `portal-web/src/api/client.ts` (Ky + wrappers)
- `portal-web/src/api/**/*.ts` (domain API modules; still must be invoked via Query/Muation hooks from UI)

Forbidden (unless explicitly documented as an exception):

- Any `fetch(...)` in routes/components
- Any direct `apiClient.*` / `get/post/...` in routes/components

If you need a new backend call:

1. Add it to the domain API module (e.g. `portal-web/src/api/customer.ts`).
2. Add a queryOptions factory (`customerQueries.*`) or a mutation hook.
3. Use the hook from the route/component.

### 2.2 Reads must be queries

- Reads are `useQuery`/`useSuspenseQuery` with a stable `queryKey` from `queryKeys`.
- Use `queryOptions(...)` factories so the same query can be reused for:
  - routes/components
  - prefetching (router loaders)
  - invalidation

### 2.3 Writes must be mutations

- Writes are `useMutation`.
- Mutation side-effects must be expressed via query invalidation:
  - prefer targeted invalidation (specific list/detail keys)
  - avoid “invalidate everything” unless the user explicitly requests refresh
  - use helpers in `portal-web/src/lib/queryInvalidation.ts` when needed

### 2.4 Global error handling (SSOT)

**Goal:** Any $\ge 400$ backend error should surface a translated, meaningful message automatically.

SSOT rule:

- Use a **global TanStack Query handler** for error-to-toast mapping.

Policy:

- **Mutations:** global onError should show a translated toast by default.
  - Most mutations represent user-triggered actions; failing silently is poor UX.
- **Queries:** global onError should be conservative.
  - Queries can be background (prefetch/refetch) and should not spam toasts.
  - Prefer letting route error boundaries / inline UI handle query errors.

### 2.5 Global toast behavior (current implementation)

The global handlers are implemented in `portal-web/src/main.tsx` via `QueryCache` and `MutationCache`.

Rules:

- Abort/cancel errors are ignored globally.
- Query error toasts are **deduped** by `queryHash` for 30s to prevent spam from refetch loops.
- Both queries and mutations support an opt-out meta flag:
  - set `meta: { errorToast: 'off' }` to suppress the global toast.

Use the shared helper `portal-web/src/lib/toast.ts` (`showErrorFromException`) which translates via `translateErrorAsync`.

**Local (per-hook) onError is allowed only when you need special behavior**, such as:

- branching UI state (e.g. show login CTA on 409)
- handling a specific status in a specific screen (e.g. rate-limit countdown from `extensions.retryAfterSeconds`)
- suppressing toast (rare; must be intentional)

Implementation guidance:

- Use `portal-web/src/lib/toast.ts` + `translateErrorAsync(error, t)`.
- Never hard-code `t('errors.generic.unexpected')` for HTTP errors; that discards backend specificity.

---

## 3) Backend interaction notes (frontend-facing)

### 3.1 Error contract

Portal assumes the backend returns Problem JSON (`application/problem+json`).

- Parsing SSOT: `portal-web/src/lib/errorParser.ts` (`parseProblemDetails`, `parseValidationErrors`)
- Backend SSOT: `.github/instructions/errors-handling.instructions.md`

### 3.2 DTO/response contract

Portal must not assume:

- timestamp casing differences (`CreatedAt` vs `createdAt`)
- nested relations exist unless the endpoint DTO explicitly includes them

Response/DTO SSOT: `.github/instructions/responses-dtos-swagger.instructions.md`

---

## 4) What to do when adding a new endpoint (end-to-end)

- Backend:
  - return a DTO (not a GORM model)
  - update handler swagger annotations
  - run `make openapi`
- Portal:
  - add schema/type under `portal-web/src/api/types/**` if applicable
  - add API method + queryOptions factory + hook under `portal-web/src/api/**`
  - consume only through `useQuery`/`useMutation`

If anything doesn’t match, log drift in `DRIFT_TODO.md`.
