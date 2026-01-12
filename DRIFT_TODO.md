# Drift / Bugs Backlog (Found During SSOT Analysis)

This file is a **tracking list only**. These items were discovered while creating SSOT instruction files (accounting/customers/inventory/orders).

Guiding rule: **backend is the source of truth for API contracts**. Portal-web should align to backend JSON shapes, query param encoding, and validation semantics.

## High priority (likely user-visible)

- [x] **Portal-web code structure violates target architecture (routes contain pages; shared components contain feature UI)**
  - **SSOT target:** `.github/instructions/portal-web-code-structure.instructions.md`
  - **Routes contain page logic:** business routes like `portal-web/src/routes/business/$businessDescriptor/orders/index.tsx` and `portal-web/src/routes/business/$businessDescriptor/inventory/index.tsx` are large “page implementations” inside the route file.
    - **Impact:** routing layer is not thin; hard to reuse/test; encourages copy-paste.
    - **Fix:** move page components into `portal-web/src/features/<feature>/components/**` and keep route files as wrappers (schema/loader + render).
  - **Resource-specific UI lives in shared components (misclassified atomic design):**
    - Molecules: `OrderCard`, `CustomerCard`, `InventoryCard`, `ShippingZoneSelect`, `BusinessSwitcher`.
    - Organisms: product/customer sheets and forms under `components/organisms/**`.
    - Atoms: resource-specific skeletons under `components/atoms/skeletons/**`.
    - **Fix:** move these into their owning features (e.g. `features/orders/components`, `features/inventory/components`, `features/customers/components`).
  - **Cross-cutting/global functionality is not modeled as features yet (but should be):**
    - Examples currently living in shared `components/**`: `BusinessSwitcher`, `LanguageSwitcher`, auth widgets like `LoginForm`.
    - **Impact:** shared component layers become “feature soup”; unclear ownership; hard to refactor without breaking unrelated UI.
    - **Fix:** create first-class feature modules such as:
      - `portal-web/src/features/business-switcher/**`
      - `portal-web/src/features/language/**`
      - `portal-web/src/features/auth/**`
      and move their state/components/forms there.
  - **Layout/template complexity is mixed into shared templates/components:**
    - Examples: complex app chrome/layout composition (e.g. dashboard/app shell) is spread across `components/templates/**` + `components/organisms/**`.
    - **Fix:** treat complex layouts as features (e.g. `portal-web/src/features/dashboard-layout/**` or `features/app-shell/**`) so they can own supporting components + local context/state.
  - **Charts + form controls are stored as atoms today (should be in dedicated folders):**
    - Chart components exist under `portal-web/src/components/atoms/*Chart.tsx`.
    - Form inputs exist under `portal-web/src/components/atoms/Form*.tsx`.
    - **Fix:** create and migrate to `portal-web/src/components/charts/**` and `portal-web/src/components/form/**` (shared, generic).
  - **Legacy file exists (violates no-legacy):** `portal-web/src/components/molecules/InventoryCard.old.tsx`.
    - **Fix:** delete it (or migrate callers and delete old in the same PR).
  - **Feature-specific utilities in `lib/`:** `portal-web/src/lib/inventoryUtils.ts`, `portal-web/src/lib/onboarding.ts`.
    - **Fix:** move to `portal-web/src/features/inventory/utils/**` and `portal-web/src/features/onboarding/utils/**` and update imports.

- [x] **Portal-web translations have competing/duplicated structures (agents pick wrong key path)**
  - **Where:**
    - `portal-web/src/i18n/{en,ar}/common.json` (namespace `common`)
    - `portal-web/src/i18n/{en,ar}/translation.json` (namespace `translation`, contains nested `common` object and also feature strings)
  - **Impact:**
    - encourages inconsistent usage (`t('common.save')` vs `useTranslation('common'); t('save')`)
    - duplicates shared keys across files/namespaces (invalid by SSOT requirements)
    - is a primary cause of “page created but translations broken” on first pass
  - **Fix:** pick one structure (preferred: true namespaces like `common`, `orders`, etc.), remove duplicated nested objects from `translation.json`, and update call sites accordingly.
  - **Migration decision (SSOT target):**
    - **Portal-web uses true i18next namespaces only**: `common`, `errors`, `onboarding`, `inventory`, `orders`, `analytics`, `upload`.
    - **New UI code must always use explicit namespaces**: `useTranslation('<namespace>'); t('<key>')`.
    - **No legacy buckets:** fully migrate all keys out of `portal-web/src/i18n/{en,ar}/translation.json` and then delete it.
    - Add explicit namespaces for anything currently living in `translation.json` (e.g. `auth`, `dashboard`, `customers`) and load them in `portal-web/src/i18n/init.ts`.
    - **No duplicated meaning** across namespaces: shared strings live in `common` only.
    - **Locale parity is mandatory**: any migrated key must exist in both `en` and `ar`.
  - **Migration steps (suggested):**
    - Inventory which keys in `portal-web/src/i18n/{en,ar}/translation.json` overlap with existing namespace files.
    - Move overlapping/shared keys into the correct namespace JSON (usually `common.json`) and delete the duplicates.
    - For remaining feature keys in `translation.json`, create the correct namespace JSON(s) (e.g. `auth.json`, `dashboard.json`, `customers.json`) and move keys there.
    - Update call sites to the single allowed pattern (`useTranslation('<ns>'); t('<key>')`).
    - Update `portal-web/src/i18n/init.ts` to remove `translation` from `resources`, `ns`, and `defaultNS` (prefer `common` as `defaultNS` once fully migrated).
    - Delete `portal-web/src/i18n/{en,ar}/translation.json` after it is empty.

- [x] **Portal-web likely contains duplicate keys inside translation JSON (invalid per SSOT)**
  - **Where:** `portal-web/src/i18n/{en,ar}/translation.json`
  - **Observation:** file content shows repeated keys in the same object (e.g. dashboard entries) and mixed formatting; even when runtime survives, duplicates are invalid by policy.
  - **Fix:** remove duplicates and enforce a single authoritative location for each key.

- [x] **Portal orders date filters don't match backend**
  - **Where:** `portal-web/src/routes/business/$businessDescriptor/orders/index.tsx`
  - **Backend expects:** `from/to` as RFC3339 timestamps (`2006-01-02T15:04:05Z07:00`) via Gin `time_format`.
  - **Portal sends today:** `yyyy-MM-dd` strings.
  - **Fix:** serialize `from/to` as RFC3339 (or update backend to accept date-only, but that's a contract change; prefer aligning portal).
  - **Resolution:** Updated `OrdersListPage.tsx` to use `formatISO()` from date-fns, which produces RFC3339/ISO 8601 timestamps (e.g., `2024-01-15T00:00:00Z`).

- [x] **Portal list sorting uses CSV `orderBy`, backend binds repeatable `orderBy`**
  - **Where:**
    - `portal-web/src/api/order.ts`
    - `portal-web/src/api/customer.ts`
    - `portal-web/src/api/inventory.ts`
  - **Backend binds:** `OrderBy []string \`form:"orderBy"\`` (repeatable query param).
  - **Portal sends today:** `orderBy=-foo,bar` (CSV) via `join(',')`.  
  - **Impact:** multi-sort will not parse correctly; backend schema mapping will see a single unknown field.
  - **Fix:** send repeatable params using `searchParams.append('orderBy', value)` for each entry.
  - **Resolution:** Updated all three API files to use `params.orderBy.forEach((o) => searchParams.append('orderBy', o))` instead of CSV join. This now produces `orderBy=value1&orderBy=value2` format that matches backend expectations.

- [x] **Portal customers `socialPlatforms` filter uses CSV, backend binds repeatable**
  - **Where:** `portal-web/src/api/customer.ts`
  - **Backend binds:** `SocialPlatforms []string \`form:"socialPlatforms"\`` (repeatable query param).
  - **Portal sends today:** `socialPlatforms=instagram,whatsapp`.
  - **Fix:** use `searchParams.append('socialPlatforms', platform)` per value.
  - **Resolution:** Updated to use `params.socialPlatforms.forEach((p) => searchParams.append('socialPlatforms', p))` instead of CSV join, matching the backend's repeatable query param binding.

- [ ] **Portal analytics sidebar link points to a missing route**
  - **Where:** `portal-web/src/components/organisms/Sidebar.tsx`
  - **Portal links today:** `/analytics`
  - **Backend reality:** analytics is business-scoped under `/v1/businesses/:businessDescriptor/analytics/*`.
  - **Impact:** clicking Analytics will 404/blank until a route exists.
  - **Fix:** implement `portal-web/src/routes/business/$businessDescriptor/analytics/**` (preferred), or remove the nav item until implemented.

- [x] **Portal asset uploads `partUrls` shape mismatches backend (upload breaks)**
  - **Where:**
    - `portal-web/src/types/asset.ts` (`UploadDescriptor.partUrls`)
    - `portal-web/src/api/assets.ts` (`uploadMultipart`)
  - **Backend returns:** `partUrls: [{ partNumber, url }, ...]`
  - **Portal expects:** `partUrls: string[]`
  - **Impact:** multipart uploads will fail at runtime (URL becomes `[object Object]`).
  - **Fix:** update portal types + upload logic to use `partUrls[i].url` (and optionally validate `partNumber`).
  - **Resolution:** Updated `UploadDescriptor.partUrls` type to `Array<{ partNumber: number; url: string }>` and modified `uploadMultipart` to use `partUrls[i].url` instead of treating array as strings.

- [x] **Portal cannot upload thumbnails on S3 (backend uses single pre-signed PUT)**
  - **Where:** `portal-web/src/api/assets.ts` (`uploadToStorage`)
  - **Backend returns (thumbnail):** `method: "PUT"` + `url` + `headers` (no `partUrls`/`partSize`)
  - **Portal supports today:** multipart PUT (`partUrls`+`partSize`) or local POST (`url`)
  - **Impact:** thumbnail upload will throw `Invalid upload descriptor` in production storage.
  - **Fix:** add support for simple pre-signed `PUT` (single request) in `uploadToStorage`.
  - **Resolution:** Added `uploadSimplePut` function and updated `uploadToStorage` to handle `method === 'PUT' && url` case (simple pre-signed PUT without multipart). Now supports three upload modes: multipart PUT, simple PUT, and local POST.

- [ ] **Backend leaks GORM `gorm.Model` fields into API JSON (PascalCase timestamps)**
  - **Where:** many domain models embed `gorm.Model` and are returned directly from handlers (e.g. orders/customers/billing).
  - **Symptom:** responses include `CreatedAt` / `UpdatedAt` / `DeletedAt` fields (PascalCase) in addition to (or instead of) `createdAt` / `updatedAt`.
  - **Impact:** breaks response casing standards, pollutes Swagger, and forces portal types to be inconsistent.
  - **Fix (preferred):** introduce per-domain response DTOs and return DTOs from handlers; update Swagger `@Success` types to DTOs.

## Medium priority (contract/type drift, correctness + maintenance)

- [x] **Portal state ownership is inconsistent (Query vs Store vs URL)**
  - **SSOT target:** `.github/instructions/state-management.instructions.md`
  - **Auth drift:** `portal-web/src/stores/authStore.ts` writes `isInitialized` into store state, but `AuthState` does not declare it.
    - **Resolution (✅):** Removed undeclared `isInitialized` write from `login()` action. State shape now matches TypeScript interface.
    - **Impact:** persistence correctness issues + hard-to-debug runtime state; also encourages a second cache.
    - **Fix:** persist only serializable data (countries/currencies/lastFetched) or remove the store and use TanStack Query as the only cache.
  - **Query↔Store mirroring drift:** metadata is maintained both in TanStack Query and in `metadataStore` (dual caches).
    - **Fix:** pick Query as the source of truth and remove store mirroring.
  - **Onboarding drift:** `portal-web/src/stores/onboardingStore.ts` exists and persists onboarding state in localStorage, but it has no known import/usage in `portal-web/src/**`.
    - **Impact:** dead code + misleading architecture; risks new code adopting the wrong pattern (localStorage onboarding token).
    - **Fix:** delete the store entirely, or intentionally integrate it (but SSOT prefers URL `session` + Query session as source of truth).
  - **Selected business drift:** `businessStore.selectedBusinessDescriptor` is persisted in localStorage, while business pages are already scoped by the URL path param `$businessDescriptor`.
    - **Impact:** two sources of truth for “current business”.
    - **Fix:** treat URL as the single source of truth; only keep a “last selected business” preference for convenience redirects (if needed).

- [ ] **No automated enforcement for i18n locale parity + duplicate-key detection**
  - **Where:** repo tooling (no check currently blocks missing keys or duplicate keys across `en`/`ar`)
  - **Impact:** missing Arabic/English keys ship silently; i18next fallback masks bugs; duplicate keys silently override.
  - **Fix:** add a CI/dev script that:
    - parses all locale resources for portal-web + storefront-web,
    - fails on duplicate keys,
    - fails on any key missing in either locale,
    - optionally enforces allowed namespaces and prevents strings from living in multiple namespaces.

- [ ] **Storefront-web uses a different translation system than portal-web (TS object vs JSON namespaces)**
  - **Where:** `storefront-web/src/i18n/translations.ts` vs portal-web JSON namespaces
  - **Impact:** inconsistent contributor behavior; higher chance of duplicated keys and missing locale parity checks.
  - **Fix:** either (A) migrate storefront to the same JSON namespace structure, or (B) document storefront as legacy and ensure validation covers both representations.

- [ ] **Backend lacks stable machine-readable error codes for consistent translation mapping**
  - **Where:** backend Problem JSON generally lacks an `extensions.code` (portal often maps by status and sometimes URL)
  - **Impact:** translation mappings become endpoint-specific and brittle; hard to provide precise, consistent localized messages.
  - **Fix:** standardize on a small set of stable error codes (e.g. `extensions.code`) and map them to `errors.*` keys in one place in portal.

- [ ] **Portal TanStack Query has no global error handler (errors handled ad-hoc)**
  - **Where:** `portal-web/src/main.tsx` (QueryClient config)
  - **Current behavior:** QueryClient defaultOptions exist, but there is no `QueryCache`/`MutationCache` `onError` that translates and shows a consistent toast for HTTP 4xx/5xx.
  - **Impact:** inconsistent UX; many mutations show generic errors; background queries may fail silently or surface inconsistent messages.
  - **Fix:** add a global mutation error handler (toast translated message via `translateErrorAsync`) and define a conservative policy for query errors (avoid toast spam; use route boundaries/inline UI).

- [ ] **Onboarding verify route bypasses TanStack Query/mutations for Google OAuth URL**
  - **Where:** `portal-web/src/routes/onboarding/verify.tsx` (`handleGoogleOAuth`)
  - **Current behavior:** calls `authApi.getGoogleAuthUrl()` directly (not via `useMutation`) and manually assigns `sendOTPMutation.error`.
  - **Impact:** violates SSOT “no direct API calls from routes/components”; error handling becomes non-standard and harder to reason about.
  - **Fix:** model Google OAuth URL fetch as a dedicated mutation hook (or query if appropriate) and handle translated errors via the standard path.

- [ ] **Some portal mutations discard backend error details (generic toast instead of translated Problem JSON)**
  - **Where:** `portal-web/src/api/address.ts` mutation hooks (`onError` uses `showErrorToast(t('errors.generic.unexpected'))`)
  - **Expected:** show translated error based on backend Problem JSON (`showErrorFromException(error, t)` / `translateErrorAsync`).
  - **Impact:** users get unhelpful generic messages even when backend provides specific details.
  - **Fix:** replace generic onError toast with translated error, and/or rely on a global mutation error handler.

- [x] **Portal auth email verification request endpoint does not exist in backend**
  - **Where:** `portal-web/src/api/auth.ts`
  - **Resolution (✅):** Updated endpoint from `POST /v1/auth/request-email-verification` to `POST /v1/auth/verify-email/request` to match backend implementation.

- [x] **Portal reset-password request body mismatches backend**
  - **Where:**
    - `portal-web/src/api/types/auth.ts`
    - `portal-web/src/features/auth/components/ResetPasswordPage.tsx`
  - **Resolution (✅):** Updated `ResetPasswordRequestSchema` to use `newPassword` field instead of `password`. Updated `ResetPasswordPage` to send `newPassword` in API call.

- [x] **Portal create/update business request body mismatches backend (`country` vs `countryCode`)**
  - **Where:** `portal-web/src/api/business.ts`
  - **Resolution (✅):** Updated `CreateBusinessRequestSchema` and `UpdateBusinessRequestSchema` to use `countryCode` field instead of `country`. Schemas now match backend expectations.

- [x] **Portal has two competing Business schemas (risk of silent contract drift)**
  - **Where:**
    - `portal-web/src/api/business.ts`
    - `portal-web/src/api/types/business.ts`
  - **Resolution (✅):** Schemas have been aligned. `business.ts` contains request/response schemas with full validation. `types/business.ts` serves as backward compatibility layer. Both now use consistent field names (`countryCode`, matching storefront theme shapes).

- [ ] **Analytics backend JSON uses `businessID` instead of `businessId` (inconsistent casing)**
  - **Where:** `backend/internal/domain/analytics/model.go`
  - **Backend returns today:** `businessID` in multiple analytics response payloads.
  - **Risk:** portal conventions (and other backend domains) typically use `businessId`.
  - **Fix:** consider migrating analytics JSON tags to `businessId` across models (breaking change; requires portal alignment).

- [ ] **Customer lifetime value (CLV) uses all-time AOV but period purchase frequency**
  - **Where:** `backend/internal/domain/analytics/service.go`
  - **Current behavior:** `averageOrderValue` is computed with an all-time date range, while purchase frequency is computed for `[from,to]`.
  - **Risk:** CLV becomes hard to interpret when users change date filters.
  - **Fix:** either compute AOV for the same `[from,to]` range, or rename the field/label to reflect mixed horizons.

- [ ] **Portal `UpdateOrderRequest` includes `shippingAddressId` but backend does not support updating it**
  - **Where:** `portal-web/src/api/order.ts` (type definition)
  - **Backend:** `backend/internal/domain/order/model.go` `UpdateOrderRequest` does not include `shippingAddressId` and service ignores it.
  - **Fix:** add backend support intentionally (must include validation + tests).

- [ ] **Portal inventory API types use snake_case keys while backend returns camelCase**
  - **Where:** `portal-web/src/api/inventory.ts` (interfaces)
  - **Backend reality:** `list.ListResponse` and domain models are camelCase (`pageSize`, `totalCount`, `productId`, etc.).
  - **Fix:** migrate portal inventory types to camelCase and update call sites; use `portal-web/src/api/order.ts` as the “correct” pattern.

- [ ] **Portal customer note timestamps use PascalCase, backend returns camelCase**
  - **Where:** `portal-web/src/api/customer.ts` (note type)
  - **Backend reality:** `createdAt`, `updatedAt`.
  - **Fix:** rename portal type fields and update any UI usage.

- [ ] **Portal customer model typing: `email` nullable vs backend create requirement**
  - **Where:** `portal-web/src/api/customer.ts`
  - **Backend reality:** create requires email; stored as non-null.
  - **Fix:** either (A) tighten portal type for created/fetched customer email, or (B) explicitly model nullable only where backend truly returns null.

- [ ] **Portal customer address drift: includes `shippingZoneId` but backend does not**
  - **Where:** portal customer address types/forms (referenced in customer SSOT)
  - **Backend reality:** customer address DTO/model does not include shipping zone today.
  - **Fix:** remove from portal forms/types, or add backend support intentionally (needs contract + tests).

- [ ] **Backend does not return structured field-level validation errors (portal can’t map form errors)**
  - **Where:** `backend/internal/platform/request/valid_body.go`
  - **Backend behavior:** struct validation errors map to `400 Bad Request` with generic `detail: "invalid request body"` and no per-field details.
  - **Portal has (unused):** `portal-web/src/lib/errorParser.ts` `parseValidationErrors()` expecting `extensions.errors|validationErrors|fieldErrors`.
  - **Impact:** portal forms can’t highlight which fields failed; users get generic toasts.
  - **Fix (preferred):** standardize backend to return a map (e.g. `extensions.fieldErrors: { field: message }`) and optionally use `422` for validation.
  - **Fix (alternate):** remove/avoid `parseValidationErrors()` assumptions in portal and keep toast-only UX.

- [ ] **Portal customer note timestamps use PascalCase fields (DTO drift)**
  - **Where:** `portal-web/src/api/customer.ts` (`CustomerNote` uses `CreatedAt/UpdatedAt/DeletedAt`).
  - **Backend models:** customer/order notes embed `gorm.Model`, so PascalCase fields can leak when returning models directly.
  - **Impact:** portal types encode backend leakage instead of the intended response standard.
  - **Fix:** after backend moves to DTOs with `createdAt/updatedAt`, migrate portal types accordingly.

- [ ] **Nested relations are declared on backend models but not guaranteed in responses**
  - **Where:** e.g. `backend/internal/domain/order/model.go` includes `Order.Notes`, `Order.Items[*].Product`, `Order.Items[*].Variant`, but handlers may not preload/return them consistently.
  - **Impact:** portal may expect nested data (notes, product info) that is missing/null depending on endpoint.
  - **Fix:** define explicit response DTOs per endpoint that state exactly what is included; preload + map accordingly; update portal types to match.

## Low priority / completeness

- [ ] **Backend does not persist multipart completion state**
  - **Where:** `backend/internal/domain/asset/storage.go` (`MarkUploadComplete`)
  - **Current behavior:** `POST /assets/uploads/:assetId/complete` completes multipart on the blob provider, but DB completion tracking is a no-op.
  - **Impact:** the system can’t reliably answer “is this upload complete?” from DB.
  - **Fix:** implement a real persisted completion marker (or completed parts JSONB) and update service/storage accordingly.

- [ ] **Portal orders API does not expose all backend endpoints**
  - **Where:** `portal-web/src/api/order.ts`
  - **Backend has:**
    - `GET /v1/businesses/:businessDescriptor/orders/by-number/:orderNumber`
    - `PATCH /v1/businesses/:businessDescriptor/orders/:orderId/payment-details`
  - **Fix:** add these methods + hooks if/when UI needs them.

- [ ] **Portal business API does not expose several backend business endpoints**
  - **Where:** `portal-web/src/api/business.ts`
  - **Backend has:**
    - `GET /v1/businesses/descriptor/availability?descriptor=...`
    - `POST /v1/businesses/:businessDescriptor/archive`
    - `POST /v1/businesses/:businessDescriptor/unarchive`
    - Shipping zone mutations under `/v1/businesses/:businessDescriptor/shipping-zones` (plan gated)
    - Payment methods list + update under `/v1/businesses/:businessDescriptor/payment-methods` (plan gated)
  - **Fix:** add client methods + hooks as UI needs them.

- [ ] **Portal has no analytics API client yet**
  - **Where:** `portal-web/src/api/**`
  - **Backend has:**
    - `GET /v1/businesses/:businessDescriptor/analytics/dashboard`
    - `GET /v1/businesses/:businessDescriptor/analytics/sales?from&to`
    - `GET /v1/businesses/:businessDescriptor/analytics/inventory?from&to`
    - `GET /v1/businesses/:businessDescriptor/analytics/customers?from&to`
    - Reports under `/v1/businesses/:businessDescriptor/analytics/reports/*?asOf`
  - **Fix:** add `portal-web/src/api/analytics.ts` + Zod schemas and query hooks when implementing UI.

## Notes

- These are intentionally tracked as a backlog. The SSOT instruction files now document the backend contract and call out the known portal drift so future work doesn’t reintroduce inconsistencies.
