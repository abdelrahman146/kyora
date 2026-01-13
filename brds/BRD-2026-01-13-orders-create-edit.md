---
status: draft
owner: product
created_at: 2026-01-13
updated_at: 2026-01-13
stakeholders:
  - name: ""
    role: "PM"
  - name: ""
    role: "Engineering Manager"
areas:
  - portal-web
  - backend
kpis:
  - name: "Time to create an order"
    definition: "Median time from tapping ‘Add order’ → order successfully created."
    baseline: "N/A (not implemented)"
    target: "≤ 60s on mobile"
  - name: "Order creation completion rate"
    definition: "% of started order drafts that successfully create an order."
    baseline: "N/A"
    target: "≥ 90%"
  - name: "Status update speed"
    definition: "Median time from opening quick actions → status updated."
    baseline: "TBD"
    target: "≤ 10s"
---

# BRD: Orders — Create, Edit, and Quick Updates (Mini-POS)

## 1) Problem (in plain language)

- Sellers live in DMs and need to record orders fast, while they’re chatting with customers.
- Today portal-web has an orders list, but no real “create order” or “edit order” experience.
- Without this, Kyora can’t become a daily tool (orders are the main input that powers stock, money in/out, and insights).
- When users try to do the “right thing” (status/payment updates), they shouldn’t hit backend errors — the UI must guide them to allowed choices.

## 2) Customer & Context

- Primary user persona(s):
  - Solo seller (admin) managing orders from WhatsApp/Instagram.
  - Team member (member role) helping during busy hours.
- Where does this happen:
  - WhatsApp / Instagram / TikTok / Facebook / Snapchat / X.
- Device context (mobile-first):
  - Used one-handed, during live chats, under time pressure.
- Language direction: RTL/Arabic-first requirements:
  - UI must be RTL-safe and fully translated (en/ar parity).
  - Phone numbers and payment references must render `dir="ltr"`.

## 3) Goals (what success looks like)

- Users can create an order in a single, calm flow with minimal steps (customer → items → total → create).
- Users can update order status and payment quickly from the list without “invalid transition” errors.
- Users can create customer + address on the fly during order creation.
- The UI prevents common backend errors by only showing valid actions and by validating basic requirements before submission.
- The overall look and flow feels consistent across Kyora (same patterns and sheets everywhere), so users don’t have to relearn similar actions in different places.
- This feature does not introduce “second versions” of existing flows (customer add, address add, selects, sheets). If something is missing, we improve the existing flow so the whole product benefits.

## 4) Non-goals (explicitly out of scope)

- Barcode scanning, receipt printing, or hardware POS integrations.
- Returns/refunds “full workflow” beyond what backend already supports.
- Complex discount rules (percent discounts, per-line discounts, promo codes).
- Shipping labels, carrier integrations.

Note: We do intend to support both fixed-amount and percent discounts for the whole order. Line-level discounts and promo codes remain out of scope.

## 5) User journey (happy path)

1. Seller opens Orders list.
2. Taps “Add order”.
3. Selects an existing customer or taps “Add new customer” (inline sheet).
4. Selects an existing address or taps “Add address” (inline sheet).
5. Adds 1–N products by searching inventory variants, adjusts quantities.
6. Optionally sets discount, ordered date/time, payment info, and note.
7. Taps “Create order”.
8. Sees success confirmation and the order appears in the list.
9. Later, seller quickly updates status/payment from the list (only valid next choices).

## 6) Edge cases & failure handling

- Out of stock / stock would go negative (backend returns `409`):
  - Expected behavior: keep the draft, highlight which items are the issue, suggest reducing quantity.
  - What the user sees: “Some items are out of stock. Adjust quantities and try again.”

- Order plan gate / subscription required (backend returns plan restriction / forbidden):
  - Expected behavior: show a calm blocker with a single next step (go to billing / plan).
  - What the user sees: “Creating orders is not available on your plan. Upgrade to continue.” or "You have reached the maximum number of orders for your current plan"

- Invalid status/payment transitions (backend would return `400`):
  - Expected behavior: prevent these by only showing allowed transitions.
  - What the user sees: Ideally never sees this; if it happens, show “This change isn’t allowed for this order right now.”

- Shipping zones not configured (customer address creation depends on zone):
  - Expected behavior: explain why, link to shipping zones setup.
  - What the user sees: “To add an address, set up shipping zones first.”

- Creating with a chosen status/payment fails due to business rules:
  - Expected behavior: keep the draft, show a clear reason and the closest allowed choice.
  - What the user sees: “This order can’t be marked as paid until it’s placed. Change status to ‘Placed’ or keep payment as ‘Pending’.”

- Offline / network failures:
  - Expected behavior: keep draft locally in memory (not persisted), allow retry.
  - What the user sees: “Connection issue. Your order draft is still here. Try again.”

- Large cart (near 100 items limit):
  - Expected behavior: show soft warning as cart grows.
  - What the user sees: “This order is getting large. Max is 100 items.”

## 7) UX / IA (mobile-first)

### Pages / Surfaces

#### Consistency principle (product requirement)

- Actions that “feel the same” must “work the same” across Kyora.
- Order creation must reuse the same customer/address add flows users already see elsewhere (same layout, validation, wording, and success behavior), so the platform feels coherent and changes/improvements apply everywhere.
- If an existing surface is close-but-not-perfect for this flow, we improve it so it serves both the existing customer feature and order creation.
- Do not introduce parallel copies of the same behavior (e.g., a second customer form or second address form). This is a user-experience consistency requirement.

#### A) Orders List (existing) — Enhanced with create + safer quick actions

- Purpose:
  - Review orders, search/filter, and do the most common updates fast.
- Primary action:
  - “Add order”
- Secondary actions:
  - Quick review
  - Quick update status
  - Quick update payment
  - Quick update shipping address (only before shipped)
- Content (what must be shown):
  - Order number, customer, total, ordered date, status + payment status.
  - Quick actions menu.
- Empty state:
  - Friendly message + “Add order”.
- Loading state:
  - Skeleton cards/table.
- Error state:
  - Plain language: “Couldn’t load orders. Try again.”
- i18n keys needed (en/ar parity):
  - Existing orders namespace + new keys for create/edit/validation copy.

**Important UX rule: state-machine-guided choices**
- Status sheet must only show allowed next statuses for the current order status.
  - Backend allowed transitions:
    - `pending → placed | cancelled`
    - `placed → ready_for_shipment | shipped | cancelled`
    - `ready_for_shipment → shipped | cancelled`
    - `shipped → fulfilled`
    - `fulfilled → returned`
    - `cancelled` and `returned` are terminal
- Payment sheet must only show allowed payment transitions AND must respect the invariant:
  - Payment changes only allowed when order status is in `{placed, shipped, fulfilled}`.
  - Allowed transitions:
    - `pending → paid | failed`
    - `failed → pending`
    - `paid → refunded`
    - `refunded` is terminal

#### B) Create Order (new) — “Mini POS” bottom sheet (full-height on mobile)

- Purpose:
  - Create an order in one calm flow during a DM conversation.
- Surface type:
  - BottomSheet / full-screen sheet on mobile.
- Primary action:
  - “Create order” (sticky at bottom, shows total)
- Secondary actions:
  - Cancel
  - Add new customer
  - Add new address
  - Add note

**Layout (single-page, sectioned, minimal steps)**
1) Customer
- Customer select (search + recent)
- Inline CTA: “Add new customer” opens AddCustomerSheet
- On customer selected: show Address section

2) Address
- Address select (list addresses for selected customer)
- Inline CTA: “Add address” opens AddressSheet
- If no addresses exist: show empty hint + CTA

3) Items
- Search inventory variants
- Add item to cart
- Cart list:
  - Variant name + product name + available stock
  - Quantity stepper (min 1)
  - Price (default from variant sale price, editable)
- Optional (advanced): show cost as read-only (from variant cost) to avoid asking user for it

4) Order info (optional, collapsed by default)
- Channel (defaults to WhatsApp; fixed list only: WhatsApp, Instagram, TikTok, Facebook, Snapchat, X)
- Ordered date/time (defaults to now; can set old orders)
- Discount (either fixed amount or percent)
- Note (plain text, optional; **only one note at create time**)

5) Payment (optional, collapsed by default)
- “Payment status” options for create:
  - Default: Pending
  - Allow: any payment status supported by backend (advanced)
- If Paid:
  - Require payment method (only methods enabled for this business)
  - Optional payment reference

6) Status (optional, collapsed by default)
- For create:
  - Default: Pending
  - Allow: any order status supported by backend (advanced)
  - UI should guide users toward safe defaults, but allow experienced users to set a final state in one flow.

**Create behavior aligned with backend**

- Backend create request today requires:
  - `customerId`, `shippingAddressId`, `channel`, `items[]`.

- Product requirement: users may optionally choose a target order status and payment status at creation time, and the backend should adapt to this in **one create request** (no “multiple requests to update status/payment” from the UI).
  - The backend should internally apply the necessary state machine transitions and invariants (including payment invariants) while keeping the UI interaction as a single “Create order” action.

Backend contract requirement (so portal-web stays simple and consistent):
- `POST /orders` should accept optional fields:
  - `status` (target order status)
  - `paymentStatus` (target payment status)
  - `paymentMethod` + optional `paymentReference` (required when paymentStatus implies paid)
  - `discountType` (`amount|percent`) + `discountValue`
  - `note` (single plain-text note content)

If backend cannot reach the requested final state (due to state machine rules), it must return a clear, user-friendly error category the UI can translate into guidance.

- Product requirement: discount supports either fixed amount or percent at order level.
  - Backend should accept a discount type + value (amount/percent) and compute totals accordingly.

- Product requirement: create order accepts an optional single note content.
  - The note is created as part of the create request.

- Empty state:
  - Item list empty: “Add at least one item.”
- Loading state:
  - During create: disable actions, show spinner in button.
- Error state:
  - Out of stock: highlight item rows.
  - Validation: show inline field errors.
  - Plan gate: show blocker + upgrade action.
- i18n keys needed (en/ar parity):
  - Orders create form labels, help text, error copy.

#### C) Edit Order (new) — Same sheet, but “editable only when allowed”

- Purpose:
  - Fix mistakes without leaving the flow (especially before shipping).
- Entry points:
  - From Order review sheet: “Edit order”
  - From quick actions: “Edit order”
- Primary action:
  - “Save changes”
- Secondary actions:
  - Cancel

Notes:
- This edit surface does **not** show or edit notes. Notes will be a separate feature.

**Rules (backend-aligned constraints)**
- Editing items is NOT allowed when status is `shipped | fulfilled | cancelled | returned`.
  - In these states: show items as read-only and explain why.
- Updating payment details is not allowed when status is `cancelled | returned`.
- Deleting an order is only allowed when status is `pending | cancelled`.

- Customer is not editable in edit order (to avoid confusion and prevent accidental reassignment).
- Shipping address can only be edited while the order is **before shipped** (i.e., when status is not `shipped` or beyond).

**What can be edited (when allowed)**
- Items (when allowed), discount, orderedAt, channel.
- Status via the status state machine.
- Payment status via the payment state machine (only in allowed statuses).
- Payment method + reference.
- Shipping address (only before shipped).

Backend contract requirement:
- `PATCH /orders/:orderId` must support updating `shippingAddressId` (only while pre-shipped), so the portal can offer a consistent address picker experience.

### Copy principles

- Use plain language and show next step.
  - Good: “Mark as paid”, “Move to shipped”, “Add customer”, “Add address”.
  - Avoid: “Invalid transition”, “Mutation failed”, “Bad request”.
- Always reassure user on errors:
  - “Your draft is still here.”

## 8) Functional requirements

- FR-1: From Orders list, user can open a Create Order sheet.
- FR-2: Create Order requires selecting/creating a customer without leaving the flow.
- FR-3: Create Order requires selecting/creating a shipping address for that customer without leaving the flow.
- FR-4: User can add 1–100 items from inventory variants with quantity control.
- FR-5: Unit price defaults from inventory sale price and is editable.
- FR-6: User can optionally enter an order discount as either fixed amount or percent.
- FR-7: User can optionally set ordered date/time (for backfilling old orders).
- FR-8: User can optionally add one note at create time only.
- FR-9: User can optionally choose target order status and payment status at create time (advanced), and the backend processes this in a single create request.
- FR-10: If user chooses a “paid” payment state on create, payment method must be required and limited to enabled methods.
- FR-11: Quick actions (status/payment) must only show valid transitions based on current state.
- FR-12: Payment methods presented in UI must be sourced from business payment methods configuration (enabled-only).
- FR-13: All failures must keep the draft and explain next step.
- FR-14: Member/admin permissions must be respected; users without manage permission cannot create/edit/update.
- FR-15: Edit order does not allow changing customer.
- FR-16: Edit order allows changing shipping address only before shipped.
- FR-17: Edit order does not show or edit notes.

## 9) Data & permissions

- Tenant scoping (workspace + business):
  - All operations are business-scoped under `/v1/businesses/:businessDescriptor/...`.
  - Never show data from another business/workspace.
- Roles (admin/member):
  - View-only users can list and view orders.
  - Manage users can create/edit/update/delete.
- What must never leak across tenants:
  - Customers, addresses, inventory variants, orders.

## 10) Analytics & KPIs

- Events to track (minimal, actionable):
  - `orders_create_opened`
  - `orders_create_customer_added_inline`
  - `orders_create_address_added_inline`
  - `orders_create_submitted`
  - `orders_create_succeeded`
  - `orders_create_failed` (with reason buckets: `out_of_stock`, `plan_gate`, `validation`, `network`, `unknown`)
  - `orders_quick_status_updated`
  - `orders_quick_payment_updated`
- KPI impact expectation:
  - Higher weekly active use; fewer missed orders; fewer stock surprises.

## 11) Rollout & risks

- Rollout plan:
  - Phase 1: Create order + quick status/payment transitions (guided).
  - Phase 2: Edit order (items/discount/orderedAt/channel) + notes from the same sheet.
- Risks:
  - Plan gates can confuse users if not explained.
  - Inventory conflicts can feel like “Kyora is blocking me”.
  - Payment methods enablement mismatch can cause errors.
- Mitigations:
  - Friendly blockers with one clear CTA.
  - Out-of-stock UX that points to exact items.
  - Only show enabled payment methods.

## 12) Open questions

Resolved:
- Create flow supports choosing any order/payment status (advanced); backend adapts in one create request.
- Channel is a fixed list (no free-text “Other”).
- Discount supports fixed amount or percent at whole-order level.

## 13) Acceptance criteria (definition of done)

- [ ] Orders list has an Add Order flow end-to-end.
- [ ] Create Order supports inline customer creation and inline address creation.
- [ ] Create Order supports adding items from inventory variants and prevents empty cart submission.
- [ ] Status and payment quick actions only allow valid transitions.
- [ ] Payment status changes are blocked when order status disallows them (UI-level).
- [ ] Payment methods shown are enabled for the business.
- [ ] Clear empty/loading/error states for all new surfaces.
- [ ] Works end-to-end on mobile.
- [ ] RTL/Arabic parity verified (including `dir="ltr"` for phone/reference).
- [ ] Multi-tenant safety respected.
- [ ] Customer creation and address creation inside order creation reuse the same existing sheets/fields used in the Customers feature (no duplicate implementations).
- [ ] Any needed behavior changes are made by improving the existing shared flows so the experience stays consistent everywhere.

## Handoff Notes for Engineering

- Portal-web will need a new orders feature surface (Create/Edit Order sheet) under the existing orders feature area.
- Portal-web likely needs to add API support for business payment methods and order payment-details update.
- The UI must implement status/payment transitions as state-machine-guided options to avoid backend 400s.
- Key backend constraints to encode in UI:
  - Stock conflict returns 409.
  - Create/edit/write routes are plan-gated.
  - Payment status changes only allowed in statuses `{placed, shipped, fulfilled}`.
