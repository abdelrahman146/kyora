---
description: "Kyora orders SSOT (backend + portal-web): endpoints, list/search, status/payment state machines, inventory adjustments, plan gates"
---

# Kyora Orders — Single Source of Truth (SSOT)

This file documents **orders behavior implemented today** across:

- Backend (source of truth): `backend/internal/domain/order/**` + wiring in `backend/internal/server/routes.go`
- Portal Web (current consumer): `portal-web/src/api/order.ts` + `portal-web/src/routes/business/$businessDescriptor/orders/**`
- Storefront (order entry path): `backend/internal/domain/storefront/**` uses order service to create **pending** orders

If you change order behavior or contracts, keep backend + portal-web consistent.

## Non-negotiables

- **Business-scoped always:** all order data is scoped to a business under `/v1/businesses/:businessDescriptor/...`.
- **No cross-tenant leaks:** cross-workspace/business access must return **404** (not found), not forbidden.
- **RBAC on every route:** orders use `role.ResourceOrder` with `ActionView` vs `ActionManage`.
- **Plan gates on “manage” operations:** most write operations require an active subscription and the `OrderManagement` feature.
- **Inventory adjustments are part of order create/update/delete:** stock is decremented on create, and restocked on delete (and on item replacements).
- **Notes are plain text:** treat `OrderNote.content` as text, never as HTML.

## Backend: route surface (authoritative)

All routes below are under:

- `/v1/businesses/:businessDescriptor/orders`

### View routes (permission: `ActionView`)

- `GET /orders` → `list.ListResponse<OrderResponse>`
- `GET /orders/by-number/:orderNumber` → `OrderResponse`
- `GET /orders/:orderId` → `OrderResponse` (includes `items[]` and `notes[]`)

### Manage routes (permission: `ActionManage` + plan gates)

These routes are grouped under middleware:

- `billing.EnforceActiveSubscription`
- `billing.EnforcePlanFeatureRestriction(billing.PlanSchema.OrderManagement)`

Endpoints:

- `POST /orders`
  - Additional plan limit gate: `billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxOrdersPerMonth, billingService.CountMonthlyOrdersForPlanLimit)`
- `PATCH /orders/:orderId`
- `DELETE /orders/:orderId`
- `PATCH /orders/:orderId/status`
- `PATCH /orders/:orderId/payment-status`
- `PATCH /orders/:orderId/payment-details`
- Notes (also manage + plan-gated):
  - `POST /orders/:orderId/notes`
  - `PATCH /orders/:orderId/notes/:noteId`
  - `DELETE /orders/:orderId/notes/:noteId`

## Backend: list/search/filter/sort contract

### `GET /orders`

Query params (Gin binding):

- `page` (default 1), `pageSize` (default 20, max 100)
- `orderBy` (repeatable, not CSV)
  - Example: `orderBy=-orderedAt&orderBy=-createdAt`
  - Fields are validated through `OrderSchema` mapping; unknown fields are ignored.
- `search` (normalized via `list.NormalizeSearchTerm`; too long → `400`)
  - Search matches:
    - `orders.search_vector` (generated)
    - `customers.search_vector` (join)
    - `orders.order_number ILIKE %term%`
    - `customers.name ILIKE %term%`
    - `customers.email ILIKE %term%`
  - If `search` is present and no explicit `orderBy`, backend uses rank ordering.

Filters:

- `status` (repeatable) values: `pending|placed|ready_for_shipment|shipped|fulfilled|cancelled|returned`
- `paymentStatus` (repeatable) values: `pending|paid|failed|refunded`
- `socialPlatforms` (repeatable) — actually filters `orders.channel` (case-insensitive)
- `customerId` (exact)
- `orderNumber` (exact)
- `from`, `to` filter `orderedAt` by RFC3339 timestamps
  - Gin expects the format `2006-01-02T15:04:05Z07:00`.

Response:

- `list.ListResponse<OrderResponse>` (camelCase list metadata): `items`, `page`, `pageSize`, `totalCount`, `totalPages`, `hasMore`.

## Backend: core order semantics

### Order number

- Generated as `id.Base62(6)`.
- Unique per business (`order_number_business_id_idx`).
- Create uses savepoints and retries (max 5) on unique violations.

### Totals

- `subtotal` = sum of `items[].total`
- `cogs` = sum of `items[].totalCost`
- `vat` = `subtotal * biz.VatRate`
- `total` = `subtotal + vat + shippingFee - discount`

### Shipping fee

Two mutually exclusive modes:

- **Manual:** `shippingFee` is accepted (must be ≥ 0) and clears `shippingZoneId`.
- **Shipping zone:** if `shippingZoneId` is provided and non-empty:
  - Zone must belong to this business.
  - Zone currency must match business currency.
  - Customer shipping address country must be included in zone countries.
  - `shippingFee` is computed as:
    - base = `subtotal - discount` (floored at 0)
    - if `freeShippingThreshold > 0` and `base >= threshold` → fee = `0`
    - else fee = `shippingCost`

### Payment method enablement

On create (and on `payment-details` update), backend validates the payment method is **enabled for the business** via business settings (if the business service is configured).

### Anti-abuse throttles (best-effort)

- Create order: cache-backed throttle per `bizId + actorId`.
- Create order note: cache-backed throttle per `bizId + actorId + orderId`.

## Backend: inventory adjustments (stock semantics)

- Creating an order **allocates stock** by decrementing each variant’s `stockQuantity`.
- If any adjustment would drive stock below 0, the whole operation fails with `409` (conflict) and stock remains unchanged.
- Updating items (when allowed) does:
  1. Delete existing order items and **restock** inventory.
  2. Create new items and **allocate** inventory.
  3. Recompute totals.
- Deleting an order (when allowed) deletes items and **restocks** inventory.

## Backend: status state machine (order lifecycle)

Allowed `status` transitions:

- `pending → placed|cancelled`
- `placed → ready_for_shipment|shipped|cancelled` (placed→shipped kept for backward compatibility)
- `ready_for_shipment → shipped|cancelled`
- `shipped → fulfilled`
- `fulfilled → returned`
- `cancelled` and `returned` are terminal

Timestamps are set at transition time:

- `placedAt`, `readyForShipmentAt`, `shippedAt`, `fulfilledAt`, `cancelledAt`.

## Backend: payment state machine

Allowed `paymentStatus` transitions:

- `pending → paid|failed`
- `failed → pending`
- `paid → refunded`
- `refunded` is terminal

Additional invariant:

- Payment status changes are only allowed when `order.status ∈ {placed, ready_for_shipment, shipped, fulfilled}`.

Event-driven automation:

- When an order becomes `paid` (and was not previously `paid`), backend emits `bus.OrderPaymentSucceededTopic`.
  - This is used by accounting automation (transaction fee upsert).

## Backend: mutation constraints

- Max items per create/update request: **100**.
- Item validation:
  - `quantity >= 1`
  - `unitPrice > 0`
  - `unitCost >= 0` (omitted defaults to 0)
  - Variant must exist in this business.
- Updating items is not allowed when status is `shipped|fulfilled|cancelled|returned`.
- Deleting an order is only allowed when status is `pending|cancelled`.
- Updating payment details is not allowed when status is `cancelled|returned`.
- Notes:
  - `content` is required and max length is **2000** (handler-level check).

## Storefront order creation (public)

Public endpoint (no auth):

- `POST /v1/storefront/:storefrontPublicId/orders`

Behavior:

- Creates an order as **pending** and **unpaid**.
- Uses server-side pricing from inventory variants.
- Uses `channel = "storefront"`.
- Sets `shippingFee = 0`, `discount = 0`.
- May create a single consolidated note if provided.

## Portal Web: implemented behavior and known gaps

### Implemented today

**File placement SSOT:** `.github/instructions/portal-web-code-structure.instructions.md`

- The items below are _current code locations_, not a requirement.
- Any new/refactored orders UI must live under `portal-web/src/features/orders/**`.

- List page: `portal-web/src/routes/business/$businessDescriptor/orders/index.tsx`
  - URL-driven search params using Zod.
  - Uses `orderQueries.list(...)` and `useOrdersQuery(...)`.
- API client: `portal-web/src/api/order.ts`
  - Implements: list, get-by-id, create, update, update status, update payment status, delete, notes CRUD.

### Known drift / footguns (document, don’t hide)

- **Date filters:** portal currently writes `from/to` as `yyyy-MM-dd` strings, but backend expects RFC3339 timestamps.
  - If you touch date filtering, align portal to backend’s RFC3339 expectation.
- **`UpdateOrderRequest` type includes `shippingAddressId`:** backend supports updating shipping address via `UpdateOrderRequest.ShippingAddressID` (validated and tested), but portal type documents it as unsupported.
  - Portal can safely use this field; backend validates that address belongs to customer and business, and only allows updates before order is shipped.
- **Missing API surface in portal:** portal API does not expose `GET /orders/by-number/:orderNumber` nor `PATCH /payment-details` yet.
- **`orderBy` encoding:** backend binds `orderBy` as a repeatable query param. Portal currently serializes it as CSV (via `join(',')`).
  - This works when the list only uses a single sort field; multiple sorts will not be parsed correctly.
  - If you touch orders list sorting, switch portal to `searchParams.append('orderBy', value)` per item.

## Change checklist (when touching orders)

Backend:

- Keep all queries business-scoped and tenant-safe.
- Keep list responses as `list.ListResponse<T>`.
- Preserve and test state machine invariants (status + payment status).
- Preserve inventory adjustment atomicity (stock should never partially update).
- If you change plan gating semantics, update route wiring in `backend/internal/server/routes.go` and extend E2E coverage.

Portal Web:

- Treat `OrderNote.content` as plain text and escape on render.
- Keep URL search params in sync with backend query params.
- Use `orderQueries` keys consistently and invalidate/refetch on mutations.
