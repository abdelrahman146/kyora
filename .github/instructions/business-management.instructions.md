---
description: "Kyora business management SSOT (backend + portal-web): businesses, descriptors, archive, shipping zones, payment methods, tenant scoping"
---

# Kyora Business Management — Single Source of Truth (SSOT)

This file documents **business management** behavior implemented today across:

- Backend (source of truth): `backend/internal/domain/business/**` + wiring in `backend/internal/server/routes.go`
- Portal Web (current consumer): `portal-web/src/api/business.ts`, `portal-web/src/stores/businessStore.ts`, `portal-web/src/routes/business/$businessDescriptor.tsx`

If you change the business contract (JSON shapes, validation, plan gates, or route paths), update backend + portal-web together.

## Non-negotiables

- **Workspace is the tenant:** all business CRUD is scoped to the authenticated actor’s workspace.
- **Never accept `workspaceId` from the client** for business routes; it must come from the authenticated actor.
- **Business scoping is via `:businessDescriptor`:** business-scoped routes must apply `business.EnforceBusinessValidity` so the business is loaded for the actor’s workspace.
- **RBAC is enforced on every business route:** view vs manage are distinct actions.

## Backend: route surface (authoritative)

### Workspace-scoped business CRUD (`/v1/businesses`)

Middleware chain:

- `auth.EnforceAuthentication`
- `account.EnforceValidActor(accountService)`
- `account.EnforceWorkspaceMembership(accountService)`

Routes:

- `GET /v1/businesses`
  - Permission: `role.ActionView` on `role.ResourceBusiness`
  - Returns: `{ businesses: Business[] }`

- `GET /v1/businesses/descriptor/availability?descriptor=...`
  - Permission: `role.ActionView` on `role.ResourceBusiness`
  - Returns: `{ available: boolean }`

- `GET /v1/businesses/:businessDescriptor`
  - Permission: `role.ActionView` on `role.ResourceBusiness`
  - Returns: `{ business: Business }`

- `POST /v1/businesses`
  - Permission: `role.ActionManage` on `role.ResourceBusiness`
  - Plan gates:
    - `billing.EnforceActiveSubscription`
    - `billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxBusinesses, businessService.MaxBusinessesEnforceFunc)`
  - Body: `CreateBusinessInput`
  - Returns: `201 { business: Business }`

- `PATCH /v1/businesses/:businessDescriptor`
  - Permission: `role.ActionManage` on `role.ResourceBusiness`
  - Body: `UpdateBusinessInput`
  - Returns: `{ business: Business }`

- `POST /v1/businesses/:businessDescriptor/archive`
  - Permission: `role.ActionManage` on `role.ResourceBusiness`
  - Returns: `204`

- `POST /v1/businesses/:businessDescriptor/unarchive`
  - Permission: `role.ActionManage` on `role.ResourceBusiness`
  - Returns: `204`

- `DELETE /v1/businesses/:businessDescriptor`
  - Permission: `role.ActionManage` on `role.ResourceBusiness`
  - Returns: `204`

### Business-scoped settings (`/v1/businesses/:businessDescriptor/...`)

Middleware chain:

- `auth.EnforceAuthentication`
- `account.EnforceValidActor(accountService)`
- `account.EnforceWorkspaceMembership(accountService)`
- `business.EnforceBusinessValidity(businessService)`

#### Shipping zones

- `GET /v1/businesses/:businessDescriptor/shipping-zones`
  - Permission: `role.ActionView` on `role.ResourceBusiness`
  - Returns: `ShippingZone[]`

- `GET /v1/businesses/:businessDescriptor/shipping-zones/:zoneId`
  - Permission: `role.ActionView` on `role.ResourceBusiness`
  - Returns: `{ zone: ShippingZone }`

Manage operations are plan gated:

- Permission: `role.ActionManage` on `role.ResourceBusiness`
- `billing.EnforceActiveSubscription`

- `POST /v1/businesses/:businessDescriptor/shipping-zones` (create)
- `PATCH /v1/businesses/:businessDescriptor/shipping-zones/:zoneId` (update)
- `DELETE /v1/businesses/:businessDescriptor/shipping-zones/:zoneId` (delete)

#### Payment methods

- `GET /v1/businesses/:businessDescriptor/payment-methods`
  - Permission: `role.ActionView` on `role.ResourceBusiness`
  - Returns: `{ paymentMethods: PaymentMethodView[] }`

Manage operations are plan gated:

- Permission: `role.ActionManage` on `role.ResourceBusiness`
- `billing.EnforceActiveSubscription`

- `PATCH /v1/businesses/:businessDescriptor/payment-methods/:descriptor` (update override)

## Backend: business descriptor rules

- Descriptor normalization happens server-side:
  - `strings.TrimSpace(strings.ToLower(descriptor))`
- Descriptor must match regex: `^[a-z0-9][a-z0-9-]{1,62}$`
  - Implicitly: length 2–63, lowercase letters/digits/hyphen, cannot start with hyphen.
- Descriptor uniqueness is scoped to the workspace.

## Backend: create/update semantics

### CreateBusinessInput (required fields)

- `name` (required)
- `descriptor` (required; normalized + validated)
- `countryCode` (required; normalized to uppercase; must be length 2)
- `currency` (required; normalized to uppercase; must be length 3)

Additional supported fields exist (brand, logo, storefront config, contact/social fields, vatRate/safetyBuffer/establishedAt).

Important behavior:

- Create is transactional and **always creates a default shipping zone**:
  - `name = <business countryCode>`
  - `countries = [<business countryCode>]`
  - `shippingCost = 0`, `freeShippingThreshold = 0`
  - `currency = <business currency>`

### UpdateBusinessInput

- Partial update (all fields optional).
- `descriptor` can be changed; it is re-normalized and must remain unique in the workspace.

## Backend: shipping zone rules

- Zone `countries` are normalized to uppercase and de-duplicated.
- Zone `currency` is always the business currency (service overwrites it on update).
- Zone name must be unique per business.
- Create/update/delete are rate limited (cache-backed), returning a business-specific rate limit error.

## Backend: payment method rules

- Payment methods are a **global catalog** with per-business overrides.
- Known global descriptors (stable identifiers used in URLs and DB):
  - `cash_on_delivery`, `bank_transfer`, `credit_card`, `tamara`, `tabby`, `paypal`
- `GET /payment-methods` returns an effective view:
  - Defaults from the global catalog
  - Overridden `enabled`, `feePercent`, `feeFixed` from the business override row (if present)
- Update request validation:
  - `feePercent` must be between 0 and 1 (inclusive)
  - `feeFixed` must be non-negative
- Update is rate limited (cache-backed).

## Portal Web: expected client behavior

### Business selection + persistence

- Business preferences are stored in a TanStack Store in `portal-web/src/stores/businessStore.ts`.
- Only these are persisted to localStorage:
  - `selectedBusinessDescriptor`
  - `sidebarCollapsed`

### Business route loader pattern

- The business layout route (`portal-web/src/routes/business/$businessDescriptor.tsx`) should:
  - Ensure auth (route guard)
  - Prefetch `businessQueries.detail(descriptor)`
  - Call `selectBusiness(descriptor)` so the UI remembers the chosen business

When adding UI to switch businesses, ensure business-scoped query keys include the descriptor and invalidate/refetch when switching.

### API client pattern

- Co-locate query options factories in `portal-web/src/api/business.ts` (`businessQueries.*`) and use them in route loaders and components.
- Prefer aligning Zod schemas with backend JSON shapes and keeping a **single SSOT schema** per domain.

## Frontend features not implemented yet (but backend supports)

If/when building business settings UI, align to backend routes above:

- Descriptor availability check (`GET /v1/businesses/descriptor/availability`)
- Archive/unarchive (`POST /v1/businesses/:descriptor/archive|unarchive`)
- Shipping zone CRUD (create/update/delete are plan-gated)
- Payment methods list + per-business override update (plan-gated)

## Known portal drift (must fix before building on it)

- **Multiple business schema sources**
  - Portal has overlapping schemas in `portal-web/src/api/business.ts` and `portal-web/src/api/types/business.ts` with different optionality/theme shapes.
  - Consolidate to a single schema SSOT to avoid silent contract drift.

- **Portal does not expose several backend business endpoints yet**
  - No client methods for descriptor availability, archive/unarchive
  - Shipping zone mutations (create/update/delete) are not exposed yet
  - Payment method update (per-business override) is not exposed yet
