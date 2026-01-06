---
description: Portal Web state management SSOT (TanStack Store vs Router vs Query vs Form)
applyTo: "portal-web/src/stores/**,portal-web/src/routes/**,portal-web/src/components/**,portal-web/src/hooks/**,portal-web/src/main.tsx,portal-web/src/router.tsx"
---

# Portal Web State Management (SSOT)

This document defines the single correct way state is owned and propagated in `portal-web/**`.

**SSOT Hierarchy**

- Parent: `.github/copilot-instructions.md`
- Portal HTTP + TanStack Query rules: `.github/instructions/http-tanstack-query.instructions.md`
- Portal architecture overview: `.github/instructions/portal-web-architecture.instructions.md`
- Onboarding flow contracts: `.github/instructions/onboarding.instructions.md`

---

## 0) Non‑negotiables (Kyora rules)

1. **Server state lives in TanStack Query.**

- Anything fetched from the backend (including “reference” lists like countries) is **server state**.
- UI code must not mirror server state into a TanStack Store as a second cache.

2. **Route state lives in the URL (path/search).**

- If a piece of state is required to render a route and should survive reload/share links, it must be expressed in the URL.
- Do not duplicate route state in a Store.

3. **Form state lives in TanStack Form.**

- Don’t keep form field values in a Store.

4. **TanStack Store is only for true client state.**

Allowed examples:

- Auth session (current user, loading flags)
- UI preferences (sidebar collapsed/open)
- “Last selected” preferences used for convenience redirects (optional)

5. **If state has two homes, it’s drift.**

- Pick the correct owner using this doc, refactor, and/or log it in `DRIFT_TODO.md`.

---

## 1) State taxonomy (what goes where)

### 1.1 URL (TanStack Router)

Use URL state for:

- Business scope: `/business/$businessDescriptor/...`
- Onboarding session token: `?session=<token>` (see onboarding SSOT)
- Onboarding plan descriptor: `?plan=<descriptor>`
- List filters/sort/pagination that should be shareable/bookmarkable

Rules:

- Prefer `validateSearch` (Zod) + typed `Route.useSearch()`.
- Route loaders may prefetch using `queryClient.ensureQueryData()`.

### 1.2 TanStack Query

Use Query for:

- All backend reads (queries)
- All backend writes (mutations)
- Cache invalidation and refetching

Rules:

- Do not copy query results into a Store for “convenience”.
- If something needs to be shared across many components, that’s a **query key** / hook concern, not a Store.

### 1.3 TanStack Store

Use Store for:

- Cross-route client state that is not derived from server data or URL.

Rules:

- Store state must be **serializable** if you persist it (no functions, no class instances).
- Persistence is for preferences only. Never persist secrets or tokens.
- Avoid “god stores”. Keep each store small and purpose-built.

### 1.4 Cookies / sessionStorage / localStorage

- **Cookies**: used for language (`kyora_language`) and refresh token (HTTP-only) by design.
- **sessionStorage**: only for short-lived cross-navigation bridging (e.g. OAuth redirect) when the URL cannot carry enough state.
- **localStorage**: preferences only. Use `portal-web/src/lib/storePersistence.ts` for persistence.

---

## 2) Kyora canonical sources of truth (today)

### 2.1 Authenticated user

- Source of truth: `authStore` (`portal-web/src/stores/authStore.ts`).
- Session restoration: `initializeAuth()` uses refresh token cookie via `restoreSession()`.
- Do not persist auth state in localStorage.

Required invariants:

- `authStore.state` shape must be stable and match its declared TypeScript state type.
- `initializeAuth()` must be idempotent and must always settle `isLoading=false`.

### 2.2 Selected business

- Source of truth when in business-scoped pages: the URL path param `$businessDescriptor`.
- `businessStore` may store:
  - the fetched businesses list (non-persisted),
  - UI preferences (sidebar),
  - optionally “last selected business” for redirecting from `/`.

Rule:

- Do not treat both the URL and `businessStore.selectedBusinessDescriptor` as independent sources of truth.

### 2.3 Language

- Source of truth: i18next language + cookie (`kyora_language`).
- Do not create a TanStack Store for language.

### 2.4 Onboarding

- Source of truth:
  - session token: URL search param `session`.
  - stage + session details: `onboardingQueries.session(sessionToken)` (TanStack Query).

Rules:

- Onboarding pages must be fully reconstructible from the URL + Query.
- Any extra browser storage must be treated as an exception and must be documented in `.github/instructions/onboarding.instructions.md`.

### 2.5 Metadata (countries/currencies)

- Countries/currencies are server state.
- Source of truth must be TanStack Query.

Rule:

- Do not maintain an additional Store cache of the same data.

---

## 3) Common anti-patterns (log as drift)

- **Mirroring Query data into Store** (creates dual caches and stale bugs).
- **Persisting functions in Store state** (JSON serialization drops them silently).
- **Keeping onboarding session token in localStorage** (token should be URL-driven).
- **Duplicating business selection** (URL param + store selection both editable).
- **Creating a Store for language** (cookie/i18n already covers it).
