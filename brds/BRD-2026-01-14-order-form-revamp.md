---
status: draft
owner: product
created_at: 2026-01-14
updated_at: 2026-01-14
stakeholders:
  - name: Product
    role: PM
  - name: Design
    role: UX/UI
  - name: Backend
    role: Lead Engineer
areas:
  - portal-web
  - backend
kpis:
  - name: Time to submit order
    definition: Median time from opening create form to successful submit
    baseline: TBD (instrument after revamp)
    target: -30%
  - name: On-the-fly customer creation rate
    definition: % of orders created with new customer/address from the form
    baseline: TBD
    target: +25%
  - name: Order accuracy confidence
    definition: % of submits without correction within 10 minutes (proxy via edit/delete)
    baseline: TBD
    target: -20% corrections
---

# BRD: Order Create/Update Revamp (mobile-first, Arabic-first)

## 1) Problem (in plain language)

- Sellers cannot create customers/addresses while creating an order; they must leave the flow.
- No live summary of items, discounts, shipping, VAT, and total, so users guess the final amount.
- UI is cluttered and not mobile-friendly; error states are unclear.
- Backend lacks a dry-run endpoint to preview totals/validation without persisting, so users risk stock/validation errors late.

## 2) Customer & Context

- Primary user persona(s): solo seller or small team member taking orders via WhatsApp/Instagram DMs.
- Where does this happen: During live DM chat while confirming items/price.
- Device context: Mobile-first (one hand use, quick taps).
- Language direction: RTL/Arabic-first; all labels/messages require en/ar parity and RTL-safe layout.

## 3) Goals (what success looks like)

- Goal 1: User can create/edit an order in ≤3 tap steps without leaving the form (customer, address, payment, items).
- Goal 2: User sees a live, accurate order summary (items, discounts, shipping, VAT, total) before submitting.
- Goal 3: Reduce failed submissions due to validation/stock/shipping zone/payment method issues by surfacing them early via dry-run.

## 4) Non-goals (explicitly out of scope)

- Adding a checkout/online payment flow (orders remain DM-driven).
- Changing order state/payment state machines.
- Implementing shipping zone CRUD or payment method settings UI (already handled in business settings).

## 5) User journey (happy path)

1. User opens create/edit order form from orders list while chatting with a customer.
2. Selects or creates customer inline; adds/edits shipping address inline.
3. Adds items (product/variant, quantity, price) with auto stock validation.
4. Applies optional discount and selects shipping mode (manual fee or shipping zone) + payment method.
5. Live summary shows subtotal, discount, shipping, VAT, total; user reviews.
6. User taps “Save order”; backend validates and persists; success toast and returns to orders list with refreshed data.

## 6) Edge cases & failure handling

- Stock insufficient for a variant:
  - Expected behavior: show inline error on the item and prevent submit; suggest lowering quantity.
  - What the user sees: “Quantity exceeds available stock for this item.”
- Shipping address country not in selected shipping zone:
  - Expected behavior: show blocking error, prompt to pick a valid zone or change address country.
  - What the user sees: “This address is outside the selected shipping zone.”
- Payment method disabled for business:
  - Expected behavior: disallow selection and show info on how to enable in settings.
  - What the user sees: “Payment method not enabled for this business. Enable it in Business Settings.”
- Discount drives subtotal negative:
  - Expected behavior: clamp at zero; show warning if discount exceeds subtotal.
  - What the user sees: “Discount can’t be more than subtotal.”
- Offline/timeout during dry-run:
  - Expected behavior: show non-blocking banner and allow retry; keep local form state.
  - What the user sees: “Preview unavailable. Check connection and retry.”
- Edit of shipped/fulfilled/cancelled/returned order:
  - Expected behavior: disable item edits and price changes; allow note/status/payment actions per state machine.
  - What the user sees: “This order can’t be edited in its current status.”

## 7) UX / IA (mobile-first)

### Pages / Surfaces

- Purpose: Create/edit an order quickly during a DM conversation.
- Primary action: Save order (create or update).
- Secondary actions: Save as draft (pending), delete (when allowed), add note, cancel/back.
- Content (what must be shown):
  - Customer selector with inline “Add customer” and “Add address” sheet; default to last used address.
  - Channel selector (Instagram/WhatsApp/TikTok/Facebook/Other) with icons.
  - Items list with add/edit/remove rows, showing variant stock.
  - Discount input (amount or percent) with toggle.
  - Shipping section: choose shipping zone or manual fee; show computed shipping fee + free-shipping note.
  - Payment method selector filtered to enabled methods.
  - Live summary card: subtotal, discount, shipping fee, VAT, total; currency shown; badge for free shipping.
  - Status (for edit) with allowed transitions only.
- Empty state: “No items yet. Add products to start.” with CTA “Add item”.
- Loading state: skeletons for selectors and summary; disable submit.
- Error state: inline field errors + top banner for blocking errors; plain language.
- i18n keys needed (en/ar parity):
  - order.form.customer.add, order.form.address.add, order.form.items.empty, order.form.discount.label,
    order.form.shipping.label, order.form.shipping.free, order.form.paymentMethod.disabled,
    order.form.summary.subtotal, order.form.summary.discount, order.form.summary.shipping,
    order.form.summary.vat, order.form.summary.total, order.form.submit, order.form.saveDraft,
    order.form.error.stock, order.form.error.zone, order.form.error.paymentDisabled,
    order.form.error.discountTooHigh, order.form.error.previewUnavailable, order.form.status.locked.

### Copy principles

- Use plain language; avoid accounting jargon.
- CTAs are actionable: “Save order”, “Add customer”, “Preview again”.

## 8) Functional requirements

- FR-1: User can create a new customer inline from the order form (name, phone, email optional) without leaving the page.
- FR-2: User can create a shipping address inline (country, city, line1, phone) and set it as default for the order.
- FR-3: Item picker validates stock against current inventory and blocks quantities that exceed available stock.
- FR-4: Discount supports fixed amount or percentage; discount cannot exceed subtotal (post-clamp at zero base).
- FR-5: Shipping supports two modes: (a) manual fee input, (b) shipping zone selection that auto-computes shipping fee based on subtotal-discount and zone rules.
- FR-6: Payment method selector only lists methods enabled for the business; disabled methods show guidance and cannot be chosen.
- FR-7: Live summary updates on every change (items, discount, shipping, VAT) showing subtotal, discount, shipping fee, VAT, total in business currency.
- FR-8: New backend dry-run endpoint returns computed totals, shipping fee, VAT, and validation errors without mutating stock/orders; UI uses it for live preview.
- FR-9: Create/Update submit uses the same validation logic as dry-run and shows inline errors mapped to fields.
- FR-10: Edit form respects status/payment state machine: prevent item/price edits when status ∈ {shipped, fulfilled, cancelled, returned}.
- FR-11: All operations are business-scoped and enforce RBAC + plan gates for order management.
- FR-12: Form must be fully usable on mobile (≤3 taps to add item and preview total; primary actions reachable without scroll).

## 9) Data & permissions

- Tenant scoping (workspace + business): All calls go through `/v1/businesses/:businessDescriptor/...`; no cross-business data visible.
- Roles (admin/member): Manage actions require `ActionManage` on orders; view uses `ActionView`.
- Plan gates: Create/update/dry-run follow `OrderManagement` feature and active subscription checks; dry-run must use same gates to avoid surprise at submit.
- No data leaks: customer/address lists are scoped to the business; shipping zones/payment methods only for the selected business.

## 10) Analytics & KPIs

- Event(s) to track: `order_form_open`, `order_form_add_customer_inline`, `order_form_add_address_inline`, `order_form_dry_run`, `order_form_submit_success`, `order_form_submit_error` (with error category: stock, zone, payment, validation, network).
- KPI impact expectation: Faster submit time, higher inline creation rate, fewer submission errors.

## 11) Rollout & risks

- Rollout plan: Behind a feature flag per workspace/business; allow opt-in; migrate gradually.
- Risks:
  - Dry-run divergence from actual create/update logic.
  - Performance on low-end mobile devices with live preview.
  - Confusion if shipping zones/payment methods misconfigured.
- Mitigations:
  - Share validation/totals code between dry-run and mutate endpoints.
  - Debounce preview calls and cache last successful preview per form state hash.
  - Provide clear guidance links to business settings when config is missing.

## 12) Open questions

- Should dry-run be rate-limited or throttled per actor to avoid abuse? No
- Do we allow item price overrides per line when variant has a base price? (today portal does; confirm policy) No.
- Should draft-saving be supported explicitly (pending without stock allocation) or always allocate stock on create? allocate stock only on create

## 13) Acceptance criteria (definition of done)

- [ ] Works end-to-end on mobile
- [ ] RTL/Arabic parity verified
- [ ] Clear empty/loading/error states
- [ ] No confusing jargon
- [ ] KPIs/events defined
- [ ] Multi-tenant safety respected

## Handoff Notes for Engineering

- Suggested owners: Backend (order service + new dry-run endpoint + shared validation), Portal-web (orders feature refactor + form UX), QA (state machine + stock validation cases), Product/Design for copy/i18n.
- Risky areas: Stock adjustment must not happen on dry-run; ensure shipping fee logic matches existing order service; state machine restrictions on edit; plan gate alignment for dry-run.
- Dependencies: Business shipping zones + payment methods APIs; inventory variant data; i18n key additions (en/ar).
