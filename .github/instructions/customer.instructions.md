---
description: "Kyora customer management SSOT (backend + portal-web): customers, addresses, notes, search/filters, RBAC, tenant isolation"
---

# Kyora Customer Management SSOT (Backend + Portal Web)

This file is the **single source of truth** for customer management behavior that is **implemented today** across:

- Backend: `backend/internal/domain/customer/**` + wiring in `backend/internal/server/routes.go`
- Portal Web: `portal-web/src/api/customer.ts`, `portal-web/src/api/address.ts`, and `portal-web/src/routes/business/$businessDescriptor/customers/**`

If you change customer management behavior, keep backend + portal-web consistent.

## Non-negotiables

- **Business-scoped always:** all customer resources are scoped under `/v1/businesses/:businessDescriptor/...` and must not leak across businesses/workspaces.
- **RBAC on every route:** customer endpoints are guarded by `role.ResourceCustomer` with `ActionView` vs `ActionManage`.
- **Backend is the API contract:** portal-web must follow backend JSON shapes and semantics.
- **ProblemDetails errors:** errors are returned as RFC7807 ProblemDetails via the shared response layer (see `.github/instructions/ky.instructions.md` for portal handling patterns).

## Backend: route surface (authoritative)

All routes are under:

- `/v1/businesses/:businessDescriptor/customers`

### Customers

- `GET /customers`
  - Pagination: `page` (default 1), `pageSize` (default 20, max 100)
  - Sorting: `orderBy` (repeatable). Use `-field` for DESC.
  - Search: `search` (normalized via `list.NormalizeSearchTerm`; overly-long values return `400`)
  - Filters:
    - `countryCode` (2-letter)
    - `hasOrders` (boolean)
    - `socialPlatforms` (repeatable) values: `instagram|tiktok|facebook|x|snapchat|whatsapp`
  - Response: `list.ListResponse<CustomerResponse>` (camelCase list metadata)

- `GET /customers/:customerId`
  - Returns customer including `addresses[]` and `notes[]` (preloaded)

- `POST /customers`
  - Creates a customer. Email uniqueness is enforced per business: `(business_id, email)`.

- `PATCH /customers/:customerId`
  - Updates a customer.

- `DELETE /customers/:customerId`
  - Soft deletes a customer.

### Customer addresses

- `GET /customers/:customerId/addresses` → returns array of `CustomerAddress`
- `POST /customers/:customerId/addresses` → creates address
- `PATCH /customers/:customerId/addresses/:addressId` → updates address
- `DELETE /customers/:customerId/addresses/:addressId` → soft deletes address

### Customer notes

- `GET /customers/:customerId/notes` → returns array of notes
- `POST /customers/:customerId/notes` → creates note
- `DELETE /customers/:customerId/notes/:noteId` → soft deletes note

**Important:** There is no PATCH/update endpoint for customer notes. Notes are create-only (immutable after creation) or can be deleted.

## Backend: RBAC and isolation rules (enforced)

Routes are guarded like:

- View: `account.EnforceActorPermissions(role.ActionView, role.ResourceCustomer)`
- Manage: `account.EnforceActorPermissions(role.ActionManage, role.ResourceCustomer)`

E2E tests confirm:

- Cross-workspace and cross-business access attempts return `404` (not found), not forbidden.
- Role `user` can view (e.g., list notes) but cannot manage (create/delete notes, create customers, etc.).

## Backend: data model semantics (what’s actually stored)

### Customer

- `Customer.Email` field:
  - **Model:** stored as `nullable.String` in GORM model (to support optional email during creation)
  - **Response:** `CustomerResponse.Email` is typed as `string` (non-null) and populated from `customer.Email.String`
  - **Validation:** `CreateCustomerRequest.Email` is optional (`binding:"omitempty,email"`), but backend requires `PhoneNumber` and `PhoneCode` as mandatory identification
  - **Uniqueness:** enforced per business via unique constraint `(business_id, email)`
- `CountryCode` is normalized to uppercase on create/update.
- Social handles are nullable strings:
  - `instagramUsername`, `tiktokUsername`, `facebookUsername`, `xUsername`, `snapchatUsername`, `whatsappNumber`

### CustomerResponse (list-only computed fields)

Customer list responses include computed aggregation fields:

- `ordersCount`
- `totalSpent`

These are computed from orders excluding statuses: `cancelled`, `returned`, `failed`.

## Backend: search + ordering (important)

### Search

`GET /customers?search=...` matches via:

- `customers.search_vector @@ websearch_to_tsquery('simple', term)`
- OR `customers.name ILIKE %term%`
- OR `customers.email ILIKE %term%`

Indexes are ensured at storage init:

- Generated `search_vector` with weighted fields (name/email highest)
- GIN on `search_vector`
- Trigram GIN indexes on `name` and `email`

When a search term is provided and there is **no explicit orderBy**, backend orders by websearch rank.

### Ordering

Supported `orderBy` fields include normal schema fields and special aggregated sorts:

- `ordersCount` / `-ordersCount`
- `totalSpent` / `-totalSpent`

When sorting by these aggregated fields, backend uses a LATERAL join for ordering and then fetches aggregation values via a single query (`GetCustomerAggregations`).

## Backend: filters (exact behavior)

- `countryCode` filter is case-insensitive (uppercased before applying).
- `hasOrders=true` means: there exists at least one non-deleted order for that customer.
- `hasOrders=false` means: no non-deleted orders exist for that customer.
- `socialPlatforms` filters by “field is non-empty”, e.g. `instagram` means `customers.instagram_username IS NOT NULL AND != ''`.

## Backend: addresses + notes (validation/normalization)

### Address create

Create request requires:

- `countryCode` (len=2), `state`, `city`, `phoneCode`, and `phone`.

Backend stores:

- `phoneCode` (string)
- `phoneNumber` (from the request field `phone`)

The backend does **not** currently validate that `phone` is E.164; portal-web currently chooses to send an E.164-ish value.

### Address update (repo reality)

Backend `UpdateCustomerAddress` now updates:

- `street`, `city`, `state`, `zipCode`, `countryCode`
- `phoneCode`, `phoneNumber` (trimmed; still optional and not validated for E.164)

Portal-web request types already include these fields; backend now persists them when provided.

### Note JSON shape

- Backend note responses use **camelCase** timestamp keys: `createdAt`, `updatedAt`.

## Backend: used by other domains (important security invariant)

`Service.GetCustomerAddressByID(...)` enforces that:

- the customer exists in this business
- the address belongs to that customer

Other domains (e.g., orders) should use this helper to prevent cross-customer / cross-business reference attacks.

## Portal Web: current implementation (how it works today)

### Where it lives

**File placement SSOT:** `.github/instructions/portal-web-code-structure.instructions.md`

Current implementation follows the target structure:

- Customer API + hooks: `portal-web/src/api/customer.ts`
- Address API + hooks: `portal-web/src/api/address.ts`
- Route wrappers: `portal-web/src/routes/business/$businessDescriptor/customers/**`
- Feature components: `portal-web/src/features/customers/components/**`
- Search schemas: `portal-web/src/features/customers/schema/**`

**Note:** No legacy `portal-web/src/components/organisms/customers/*` exists; implementation is already in features.

### URL-driven list state

Customers list uses TanStack Router search params:

- `search`, `page`, `pageSize`
- `sortBy`, `sortOrder` → translated to backend `orderBy` (e.g., `-joinedAt` by default)
- `countryCode`, `hasOrders`, `socialPlatforms[]`

### Data fetching and invalidation

- Uses TanStack Query keys from `portal-web/src/lib/queryKeys.ts` under `queryKeys.customers.*`.
- Detail page invalidates customer detail after note mutations.
- Delete customer invalidates the customers list.

### Known drifts (do not propagate)

Portal-web includes some type/contract drift vs backend. Don’t copy these into new code; align to backend reality when touching customer management:

- **`Customer.email` typing:** Portal-web types `Customer.email` as `string` (non-null), which matches `CustomerResponse` from backend but not the underlying GORM model (`nullable.String`). This is acceptable as the API contract uses response DTOs.
- **`CustomerNote` timestamp casing:** `CustomerNote` type in `portal-web/src/api/customer.ts` currently uses `createdAt/updatedAt` (camelCase), which correctly matches backend `CustomerNoteResponse` (camelCase).
- **`CustomerAddress.shippingZoneId` field:** Portal-web address types and forms include `shippingZoneId`; backend `CustomerAddress` model and DTOs do not currently support shipping zones, and responses will not include `shippingZoneId`. This is used in portal-web for UI pre-selection but not sent to backend.
- **`orderBy` encoding:** backend binds `orderBy` as a repeatable query param (e.g., `orderBy=-joinedAt&orderBy=name`). Portal currently serializes it as CSV (via `join(',')`).
  - This works only when a single field is provided; multiple fields will not be parsed correctly by backend schema mapping.
  - If you touch customers list sorting, switch portal to `searchParams.append('orderBy', value)` per item.
- **`socialPlatforms` encoding:** backend binds `socialPlatforms` as repeatable values. Portal currently sends CSV (`socialPlatforms=instagram,whatsapp`), which does not match backend binding and can break filtering.
  - If you touch customers list filters, switch portal to repeatable `searchParams.append('socialPlatforms', platform)`.

## Change checklist (when extending customer management)

Backend:

- Update handlers in `backend/internal/domain/customer/handler_http.go`.
- Keep business scoping via `ScopeBusinessID(biz.ID)` and prevent reference attacks (use `GetCustomerByID` / `GetCustomerAddressByID`).
- If you add new list filters/sorts, ensure indexes exist and add/extend e2e tests under `backend/internal/tests/e2e/customer_*_test.go`.

Portal Web:

- Update API client + hooks in `portal-web/src/api/customer.ts` / `portal-web/src/api/address.ts`.
- Keep list state URL-driven (TanStack Router search schema).
- Use existing form system + UI primitives (see `.github/instructions/forms.instructions.md` and `.github/instructions/ui-implementation.instructions.md`).
