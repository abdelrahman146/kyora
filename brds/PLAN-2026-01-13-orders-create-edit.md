---
status: draft
created_at: 2026-01-13
updated_at: 2026-01-13
brd_ref: "brds/BRD-2026-01-13-orders-create-edit.md"
owners:
  - area: backend
    agent: Feature Builder
  - area: portal-web
    agent: Feature Builder
  - area: tests
    agent: Feature Builder
risk_level: high
---

# Engineering Plan: Orders — Create/Edit + Quick Updates (Mini-POS)

## 0) Inputs

- BRD: brds/BRD-2026-01-13-orders-create-edit.md
- SSOT references:
  - Orders SSOT: .github/instructions/orders.instructions.md
  - Backend architecture SSOT: .github/instructions/backend-core.instructions.md
  - Portal architecture SSOT: .github/instructions/portal-web-architecture.instructions.md
  - Portal HTTP + Query SSOT: .github/instructions/http-tanstack-query.instructions.md
  - Portal code structure SSOT: .github/instructions/portal-web-code-structure.instructions.md

Assumptions:
- Backend is source of truth for state machines and validation.
- Portal-web has no automated UI test suite today; verification will be manual + backend E2E.
- “Consistency/DRY” is a product requirement: reuse existing sheets/patterns so the experience feels unified.
- “Minimize experience changes” means we prefer improving an existing shared flow over introducing a second similar flow.

## 1) Confirmation Gate (must be approved before implementation)

This plan proposes contract and behavior changes that must be explicitly approved:

- Backend API contract change: `POST /v1/businesses/:businessDescriptor/orders` accepts and applies optional target `status` and `paymentStatus` in the same request.
- Backend data/behavior change: order-level discount supports `amount|percent` (new fields + new totals logic).
- Backend API behavior change: `POST /orders` accepts an optional single note content and creates it as part of the create operation.
- Backend API behavior change: `PATCH /orders/:orderId` supports updating `shippingAddressId` (restricted to pre-shipped statuses).

No new libraries or new apps are proposed.

## 2) Architecture summary (high level)

### Options (pick one)

- Option A — Minimal (fastest): keep backend as-is; portal-web sequences multiple calls on create (create → status → payment status → payment details → note).
  - Pros: quickest; lowest backend risk.
  - Cons: violates the product requirement (single action, fewer API calls); more failure points; harder UX.
  - Complexity: M

- Option B — Robust (preferred default): extend backend create endpoint to accept “final desired state” (status/payment/discount/note) and apply internally in one atomic operation.
  - Pros: matches product requirement; fewer UI failure modes; easier to keep consistent.
  - Cons: requires backend model/handler/service changes + new E2E coverage.
  - Complexity: L

- Option C — Strategic expansion: Option B + add “draft autosave” in portal-web (local-only) and “duplicate order” for speed.
  - Pros: best daily usability for power users.
  - Cons: scope creep; more QA surface.
  - Complexity: L+

### Chosen approach

- Proposed: Option B.

### Backend

- Extend order create DTO + service to apply:
  - optional target order status
  - optional target payment status + required payment details when needed
  - optional discount type/value
  - optional note creation
- Maintain existing plan gates and inventory atomicity.
- Preserve existing state machine invariants (same as SSOT).

### Portal-web

- Add a single Create/Edit Order sheet under `portal-web/src/features/orders/**`.
- Reuse existing customer/address sheets (`features/customers/...`) for inline creation, and improve them if needed (no parallel implementations).
- Enforce valid transitions in UI (guided choices) to prevent backend errors.
- Add business payment methods API module to only show enabled methods.

### Data model

- Orders: add fields for discount type/value (or a new embedded struct) and potentially computed/normalized discount amount for storage.
- Potentially update order request/response DTOs in swagger.

### Security/tenancy

- No cross-business access: business-scoped routes only.
- RBAC: ActionView for list/detail; ActionManage + plan gates for create/update/delete.

## 3) Work breakdown (handoff-ready)

### Milestone 0 — UX consistency audit (shippable)

Goal: Ensure we reuse existing Kyora patterns and avoid duplicating flows.

Portal-web tasks:
- Audit existing customer/address create/edit UX and selects used across the app.
- Define the “single source” components that order creation will reuse (customer sheet, address sheet, select fields, BottomSheet layout, validation/copy).
- If order creation needs a missing capability (e.g., returning the created entity, preselecting values, better error copy), implement it by enhancing the existing shared flow.

Rollout notes:
- This milestone is complete when the plan’s component reuse map is clear enough that implementation cannot accidentally clone existing behaviors.

### Milestone 1 — Backend contract + E2E foundations (shippable)

Goal: Make the backend capable of the desired create/edit behavior safely (so portal-web can remain simple).

Backend tasks:
- Update `CreateOrderRequest` (backend order domain) to accept optional:
  - `status` (target)
  - `paymentStatus` (target)
  - `paymentMethod` + `paymentReference` (required if paymentStatus implies paid)
  - `discountType` + `discountValue` (percent or amount)
  - `note` (single note content)
- Maintain backwards compatibility with existing portal fields:
  - Keep existing `discount` amount behavior; if new discount fields are present, they take precedence.
- Apply the target states atomically inside the create service:
  - Validate requested target states against current state machines.
  - Enforce payment invariant (payment status changes only allowed when order status is in `{placed, shipped, fulfilled}`), even on create.
  - Ensure inventory adjustments remain all-or-nothing.
- Add `PATCH /orders/:orderId` support for `shippingAddressId` updates:
  - Allowed only when status is before shipped.
  - Ensure address belongs to the same customer and business.
- Expose payment methods enablement already exists on business domain; ensure it is used in order validation.
- Update swagger docs and run `make openapi`.

Tests:
- Add/extend backend E2E tests (`backend/internal/tests/e2e/order_*_test.go`) for:
  - Create order with percent discount.
  - Create order with target status + payment status (valid path).
  - Create order with invalid requested target state returns 400 with actionable problem details.
  - Create order out-of-stock returns 409 and no partial stock changes.
  - Update shipping address allowed pre-shipped; blocked post-shipped.

Rollout notes:
- No UI is depending on new fields until Milestone 2; safe to deploy backend first.

### Milestone 2 — Create Order sheet (shippable)

Goal: Deliver the “daily driver” flow: create order in one calm sheet, inline customer/address add, minimal steps.

Portal-web tasks:
- Implement Create Order sheet under `portal-web/src/features/orders/**` using existing patterns (BottomSheet, useKyoraForm, resource list layout conventions).
- Reuse existing sheets:
  - `features/customers/components/AddCustomerSheet.tsx` for inline add.
  - `features/customers/components/AddressSheet.tsx` for inline add.
  - Existing select fields (CustomerSelectField, AddressSelectField) rather than duplicating.
- Reuse existing “resource list” visual language (headers, actions, buttons, loading/empty states) so the new flow feels native.
- Inventory picker:
  - Variant search + add to cart.
  - Quantity controls, price default, editable price.
  - Inline out-of-stock messaging consistent with global error handling.
- Payment section:
  - Fetch enabled payment methods via new portal API module calling `GET /v1/businesses/:businessDescriptor/payment-methods`.
  - Only show enabled methods.
- Advanced section:
  - Fixed channel list.
  - Optional orderedAt.
  - Optional target status/payment status (advanced UI).
  - Discount type selector (amount/percent) + value field.
  - Optional single note field.
- Wire `OrdersListPage` “Add order” CTA to open the sheet.

Backend tasks:
- None expected beyond Milestone 1 changes.

Tests:
- Manual QA checklist (mobile/RTL):
  - create via existing customer + address
  - create customer inline
  - create address inline
  - percent discount totals
  - create with advanced status/payment
  - plan-gated behavior (upgrade CTA)

Rollout notes:
- Keep advanced status/payment collapsed by default to protect novice users.

### Milestone 3 — Edit Order sheet + safer quick actions (shippable)

Goal: Enable correcting orders without breaking state machine rules or shipped orders.

Portal-web tasks:
- Implement Edit Order sheet (same component) with rules:
  - Customer is read-only.
  - Address editable only before shipped.
  - Items editable only when allowed; otherwise read-only with explanation.
  - Notes are not shown/edited.
- Update `OrderQuickActions` to:
  - Only show allowed next status values.
  - Only show payment status options allowed for current status.
  - Only show “update address” when allowed (pre-shipped).
  - Use backend payment-details endpoint for method/reference when required (avoid incorrect “update items to update payment” workarounds).

Backend tasks:
- If portal requires a dedicated endpoint for payment details (already exists per orders SSOT), ensure portal consumes it.

Tests:
- Backend E2E: confirm payment status invariant is enforced.
- Manual QA:
  - edit order in each status bucket
  - attempt forbidden edits (ensure UI blocks and backend rejects)

Rollout notes:
- Treat quick actions as the primary daily updates path; prioritize smooth error copy.

## 4) API contracts (high level)

Endpoints impacted/added:
- Orders:
  - `POST /v1/businesses/:businessDescriptor/orders` (extend request)
  - `PATCH /v1/businesses/:businessDescriptor/orders/:orderId` (support shippingAddressId update + discount type/value)
  - `PATCH /v1/businesses/:businessDescriptor/orders/:orderId/status` (no change; still supported)
  - `PATCH /v1/businesses/:businessDescriptor/orders/:orderId/payment-status` (no change)
  - `PATCH /v1/businesses/:businessDescriptor/orders/:orderId/payment-details` (ensure portal uses)
- Business:
  - `GET /v1/businesses/:businessDescriptor/payment-methods` (portal consumes; implement portal API module)

DTO highlights:
- Create order request additions (proposed):
  - `status?: OrderStatus`
  - `paymentStatus?: OrderPaymentStatus`
  - `paymentMethod?: OrderPaymentMethod`
  - `paymentReference?: string`
  - `discountType?: 'amount'|'percent'`
  - `discountValue?: string`
  - `note?: string`

Error cases:
- 400: invalid requested state transition / invalid discount / invalid payment method.
- 409: stock would go below 0.
- Plan gate responses as per billing middleware.

## 5) Data model & migrations

- Orders table:
  - Add discount fields to support percent (exact schema to be decided in backend implementation).
  - Ensure totals logic remains consistent and stored values remain correct.

Migration plan:
- Add nullable fields for the new discount representation.
- Backfill existing rows to “amount” using current `discount`.

Indexing:
- No new indexes expected.

## 6) Security & privacy

- Tenant scoping:
  - All operations must remain business-scoped and enforce business validity middleware.
- RBAC:
  - View vs Manage permissions on all endpoints.
- Abuse prevention:
  - Preserve existing create throttle and note throttle behavior.

## 7) Observability & KPIs

- Portal events (tracked client-side):
  - `orders_create_opened`, `orders_create_submitted`, `orders_create_succeeded`, `orders_create_failed` (reason buckets)
  - `orders_quick_status_updated`, `orders_quick_payment_updated`

## 8) Test strategy

- E2E (backend testcontainers):
  - Create with percent discount
  - Create with target status/payment
  - Out-of-stock 409 + no partial stock
  - Update shipping address pre/post shipped
  - Plan gating for manage endpoints

- Integration tests:
  - Order totals calculation with percent discount and VAT.

- Manual QA (portal-web):
  - Mobile-first, RTL-first flows
  - Inline customer/address create
  - Error copy and recovery (draft preserved)

## 9) Risks & mitigations

- Risk: “Single create request” becomes complex and error-prone.
  - Mitigation: keep backend state application logic centralized in order service; cover with E2E + integration tests.

- Risk: Discount percent introduces rounding disputes.
  - Mitigation: define rounding rules explicitly (backend source of truth) and display totals consistently.

- Risk: Payment invariants cause confusing failures.
  - Mitigation: UI guides choices; backend returns structured, translatable problem details.

- Risk: Portal accidentally duplicates customer/address flows.
  - Mitigation: enforce Milestone 0 reuse audit; explicitly reuse existing sheets and select fields; adjust them if needed rather than cloning.

- Risk: New order creation introduces one-off UI patterns that don’t match Kyora.
  - Mitigation: reuse existing BottomSheet/layout/form components; keep copy style consistent; avoid new one-off controls unless promoted into a shared component used elsewhere.

## 10) Definition of done

- [ ] Meets BRD acceptance criteria
- [ ] Backend supports single-request create with optional status/payment/discount/note
- [ ] Portal create/edit sheets shipped and consistent with existing UX
- [ ] Quick actions only allow valid transitions
- [ ] No duplicated customer/address/order form implementations; any needed changes were made by enhancing the existing shared flows
- [ ] Mobile-first UX verified
- [ ] RTL/i18n parity verified
- [ ] Multi-tenancy verified
- [ ] Error handling + empty/loading states complete
- [ ] No TODO/FIXME
