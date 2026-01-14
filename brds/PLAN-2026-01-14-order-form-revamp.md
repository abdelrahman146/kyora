---
status: draft
created_at: 2026-01-14
updated_at: 2026-01-14
brd_ref: BRD-2026-01-14-order-form-revamp.md
owners:
  - area: backend
    agent: Feature Builder
  - area: portal-web
    agent: Feature Builder
  - area: tests
    agent: Feature Builder
risk_level: medium
---

# Engineering Plan: Order Create/Update Revamp

## 0) Inputs

- BRD: [brds/BRD-2026-01-14-order-form-revamp.md](brds/BRD-2026-01-14-order-form-revamp.md)
- UX Spec: [brds/UX-2026-01-14-order-form-revamp.md](brds/UX-2026-01-14-order-form-revamp.md)
- Assumptions:
  - Dry-run must reuse the same validation/totals logic as create/update and be plan-gated the same.
  - Line price overrides remain allowed; stock allocation occurs only on create/update (not on dry-run).
  - Shipping zones/payment methods are preconfigured via Business Settings; UI only consumes.
  - Draft state is in-memory only while the sheet is open (no server draft persistence).

## 1) Confirmation Gate (must be approved before implementation)

- New dependency/library? No.
- New project/app? No.
- Breaking change? None intended; new endpoint is additive.
- Migration? None.
- Data model change with customer impact? None (additive endpoint only). No confirmation-gated changes proposed.

## 2) Architecture summary (high level)

- Backend: Add additive dry-run endpoint under business-scoped orders that reuses order service validation/totals without persisting or adjusting stock; plan-gated and RBAC-managed identically to create/update.
- Portal-web: Refactor orders create/edit into a guided 3-card layout with inline sheets (customer/address), item picker, live summary, preview pill, sticky action bar; integrate dry-run preview API with stale/error banners.
- Data model: Unchanged tables; reuse existing order, item, shipping zone, payment method, customer/address models.
- Security/tenancy: Keep business-scoped routes under `/v1/businesses/:businessDescriptor`; enforce RBAC (`ResourceOrder` view/manage) and plan gates on dry-run; no cross-tenant leakage.

## 3) Step-based execution plan (handoff-ready)

Execution protocol:
- Feature Builder will implement **one step per request**.
- Each step below is sized to be completed “perfectly” in a single AI request.
- Do not start Step N+1 before Step N is merged/verified.

### Step Index

- Step 0 — Repo alignment + DRY map
- Step 1 — Backend dry-run endpoint
- Step 2 — Portal API + preview plumbing
- Step 3 — Create sheet 3-card layout + summary
- Step 4 — Inline customer/address + item picker revamp
- Step 5 — Edit sheet status-aware + sticky actions
- Step 6 — i18n, analytics, polish
- Step 7 — Tests (backend + portal) and verification

### Step 0 — Repo alignment + DRY map (required)

- Goal: Confirm current entry points, avoid duplication, and finalize DRY targets before coding.
- Scope (in/out): In: inventory of orders UI/backend surfaces; Out: code changes.
- Repo Recon (Evidence):
  - Portal-web entry points: orders list [portal-web/src/routes/business/$businessDescriptor/orders/index.tsx](portal-web/src/routes/business/$businessDescriptor/orders/index.tsx); create/edit sheets [portal-web/src/features/orders/components/CreateOrderSheet.tsx](portal-web/src/features/orders/components/CreateOrderSheet.tsx), [portal-web/src/features/orders/components/EditOrderSheet.tsx](portal-web/src/features/orders/components/EditOrderSheet.tsx); review sheet [portal-web/src/features/orders/components/OrderReviewSheet.tsx](portal-web/src/features/orders/components/OrderReviewSheet.tsx); cards [portal-web/src/features/orders/components/OrderCard.tsx](portal-web/src/features/orders/components/OrderCard.tsx); shipping zone helper [portal-web/src/components/organisms/ShippingZoneInfo.tsx](portal-web/src/components/organisms/ShippingZoneInfo.tsx).
  - Backend entry points: orders domain [backend/internal/domain/order/**]; routes wiring [backend/internal/server/routes.go]; current create/update handlers (no dry-run endpoint yet).
- Code Structure & Reuse:
  - New files (if any) + reuse scope: LiveSummaryCard, PreviewStatusPill, ItemPickerSheet under `portal-web/src/features/orders/components/`.
  - Modified files grouped by responsibility: backend order handler/service; portal API order client; create/edit sheets; review sheet; orders list triggers; i18n resources.
  - Do-Not-Duplicate list: BottomSheet; form system (`form.AppField` et al.); ShippingZoneInfo; formatCurrency/formatDate; existing order queries/mutations; loading skeletons.
- Tasks:
  - Backend: confirm no dry-run route exists; identify reusable validation/totals functions in order service.
  - Portal-web: confirm current form fields and gaps vs UX spec; list translation namespaces to touch (`orders`, `common`, `customers`, `address`).
  - Tests: identify existing order E2E coverage (if any) under `backend/internal/tests/e2e/` and portal smoke coverage (none today).
- Verification: Document findings in this plan; no code changes.
- Definition of done: DRY map agreed; no pending unknowns.

### Step 1 — Backend dry-run endpoint

- Goal (user-visible outcome): Provide a non-persisting preview endpoint that returns totals/fees/validation for the form.
- Scope (explicitly in/out): In: new endpoint under orders; reuse validation/totals; plan gates/RBAC; response DTO. Out: state machine changes; data model changes.
- Targets (files/symbols/endpoints): `backend/internal/domain/order/handler_http.go`, `service.go`, DTOs/models; `backend/internal/server/routes.go`; add tests under `backend/internal/tests/e2e/` or integration.
- Tasks (detailed checklist):
  - Backend:
    - Add `POST /v1/businesses/:businessDescriptor/orders/preview` (dry-run) with `ActionManage` + `OrderManagement` plan gates.
    - Parse payload same as create/update; call shared validation/totals logic without mutating DB or stock; ensure shipping zone/payment method checks executed.
    - Return `OrderPreviewResponse` with items, subtotal, discount, shipping fee, VAT, total, currency, free-shipping flag, and echo of normalized inputs (discountType/value, shipping mode).
    - Prevent stock adjustment; wrap in transactionless flow.
    - Add throttling (reuse existing order create throttle helper) to avoid abuse (best-effort).
    - Wire route in `routes.go` alongside manage routes.
  - Tests:
    - Add E2E/integration test: valid payload returns totals; mismatched zone/currency returns 400; disabled payment method returns 400; stock-insufficient returns 409-like problem consistent with create.
- Edge cases + error handling: zone-country mismatch; negative discount; shipping mode mutual exclusivity; max items (100) enforced; plan gate rejects without subscription.
- Verification checklist: Swagger (if required) updated; tests green; response matches frontend needs; no DB writes.
- Definition of done: Dry-run endpoint live, documented, and covered by tests.

### Step 2 — Portal API + preview plumbing

- Goal: Expose preview endpoint in portal API client with TanStack Query mutation/query helpers and error mapping.
- Scope: In: `portal-web/src/api/order.ts` query/mutation factory; types; error handling. Out: UI wiring (later steps).
- Targets: `portal-web/src/api/order.ts`, `portal-web/src/api/types/order.ts` (or existing type file), `portal-web/src/lib/queryKeys.ts` if needed.
- Tasks:
  - Add `orders.preview` ky client + Zod schema; include businessDescriptor scoping; map response to typed preview result.
  - Add query/mutation helpers (e.g., `orderMutations.preview`) with debounce guidance for UI.
  - Ensure error parsing uses existing helpers; align status codes to user-facing categories (stock, zone, payment disabled).
- Edge cases: RFC3339 date compliance not needed here; ensure shippingZoneId/manual fee mutual exclusivity; handle 409 stock as validation error category.
- Verification: Type checks; manual call against backend; unit test for schema parse if lightweight.
- Definition of done: Preview API callable with typed result and correct error surfaces.

### Step 3 — Create sheet 3-card layout + summary

- Goal: Implement guided, decluttered create flow with 3 cards, sticky action bar, preview pill, and LiveSummaryCard.
- Scope: In: layout refactor; summary card; preview pill; sticky actions; wiring to preview API. Out: inline customer/address creation (next step), item picker revamp (next step).
- Targets: `CreateOrderSheet.tsx`; new components `LiveSummaryCard.tsx`, `PreviewStatusPill.tsx`; shared styles if needed.
- Tasks:
  - Recompose layout into header strip + Card1 (customer/channel), Card2 (items/discount), Card3 (shipping/payment/notes) + summary card + sticky action bar.
  - Insert PreviewStatusPill showing success/stale/error + timestamp; hook into preview call results.
  - Implement LiveSummaryCard rendering subtotal/discount/shipping/VAT/total, free-shipping badge, currency label; shows skeleton and error banner.
  - Wire preview calls: auto-trigger on relevant field changes with debounce; mark stale on change; disable submit until latest success.
  - Sticky action bar: primary `create_submit`, secondary cancel; display last preview time.
- Edge cases: prevent submit when no customer/items; block when preview stale/error; handle loading states on initial data fetch.
- Verification: UI renders on mobile/desktop; preview triggers; submit stays disabled until preview success; no regressions to existing mutations.
- Definition of done: Create sheet matches UX layout and preview gating; new components in place and reused within create.
- Refer to `brds/UX-2026-01-14-order-form-revamp.md` whenever you need more clarity about UI/UX
- Refer to `brds/BRD-2026-01-14-order-form-revamp` whenever you need more about the business requirements.

### Step 4 — Inline customer/address + item picker revamp

- Goal: Deliver inline creation sheets and clean item picker per UX spec.
- Scope: In: Add/rewire inline customer and address sheets; new ItemPickerSheet; integrate into Card2 items list. Out: edit-specific locking (handled in Step 5).
- Targets: `CreateOrderSheet.tsx`, new `ItemPickerSheet.tsx`; customer/address sheets (existing or new lightweight components under orders feature or shared); forms libs.
- Tasks:
  - Implement Inline Add Customer sheet (3 fields, phone LTR) launched from Card1; on success, select customer and prompt address add.
  - Implement Inline Add Address sheet (country/city/state/street/phone/default toggle) launched from Card1; on success, select address and re-run preview with inferred zone.
  - Replace inline item row entry with ItemPickerSheet (searchable variant, quantity stepper, price override, stock badge); connect add/edit per row; keep unit cost optional.
  - Ensure dry-run reruns after item/discount/ship changes; keep form state when sheets close.
- Edge cases: missing addresses; stock warnings; percent discount cap; address-country mismatch with zone triggers error callout.
- Verification: Create flow allows adding customer/address without leaving; item picker works and updates list; preview reruns; RTL/mobile layout intact.
- Definition of done: Inline sheets and item picker functional and integrated in create flow.
- Refer to `brds/UX-2026-01-14-order-form-revamp.md` whenever you need more clarity about UI/UX
- Refer to `brds/BRD-2026-01-14-order-form-revamp` whenever you need more about the business requirements.

### Step 5 — Edit sheet status-aware + sticky actions

- Goal: Bring edit flow to new layout with status-aware locking and shared components.
- Scope: In: EditOrderSheet refactor to reuse new layout/components; enforce status locks; sticky action bar; preview gating. Out: order status/payment status mutation flows (unchanged).
- Targets: `EditOrderSheet.tsx`, reuse `LiveSummaryCard`, `PreviewStatusPill`, `ItemPickerSheet` (when allowed), `Inline` sheets.
- Tasks:
  - Mirror 3-card layout; show status/payment chips in header strip; lock items/price when status ∈ {shipped, fulfilled, cancelled, returned}; show lock notice.
  - Allow shipping/discount edits only when permitted by backend rules; ensure preview honors same validation.
  - Wire submit gating to preview success + allowed status; secondary cancel.
- Edge cases: paid orders with payment method changes (respect backend rules); ensure preview errors still show even when items locked.
- Verification: Edit flow renders; locked states visually clear; preview required; submit disabled when disallowed.
- Definition of done: Edit sheet matches UX; uses shared components; status locks enforced.
- Refer to `brds/UX-2026-01-14-order-form-revamp.md` whenever you need more clarity about UI/UX
- Refer to `brds/BRD-2026-01-14-order-form-revamp` whenever you need more about the business requirements.

### Step 6 — i18n, analytics, polish

- Goal: Add/align translation keys, analytics events, and UX polish states.
- Scope: In: i18n keys per UX inventory; analytics events from BRD; minor UI polish (badges, skeletons, banners). Out: net-new features.
- Targets: `portal-web/src/i18n/{en,ar}/orders.json` (+ common/customers/address); analytics hook points in create/edit flows.
- Tasks:
  - Add new keys listed in UX spec; ensure en/ar parity; avoid duplication; route titles remain under common if touched.
  - Emit events: `order_form_open`, `order_form_add_customer_inline`, `order_form_add_address_inline`, `order_form_dry_run`, `order_form_submit_success`, `order_form_submit_error` (with category) via existing analytics helper if available.
  - Polish states: empty/loading/error visuals per spec; badge for free shipping; preview stale timer (30s).
- Edge cases: avoid breaking existing translations; ensure namespace usage explicit.
- Verification: i18n build passes; manual check en/ar; analytics calls fire in devtools (if instrumented); UI shows polish elements.
- Definition of done: Keys added with parity; analytics events wired; polish items visible.
- Refer to `brds/UX-2026-01-14-order-form-revamp.md` whenever you need more clarity about UI/UX
- Refer to `brds/BRD-2026-01-14-order-form-revamp` whenever you need more about the business requirements.

### Step 7 — Tests (backend + portal) and verification

- Goal: Validate flows end-to-end and unit where possible.
- Scope: In: backend E2E for preview; portal integration/unit for preview hook and layout gating; manual smoke for mobile/RTL.
- Targets: backend `internal/tests/e2e/order_*_test.go` (add new preview-focused file); portal tests if framework available; manual checklist.
- Tasks:
  - Backend: cover preview happy path, stock conflict, zone mismatch, payment disabled, plan gate denial.
  - Portal: add lightweight tests for preview hook (staleness gating) if test infra exists; otherwise document manual QA checklist for mobile/RTL.
  - Manual QA: mobile viewport, RTL, sticky bar, inline sheets, preview gating, error banners.
- Edge cases: 100-item limit; discount edge; free shipping threshold.
- Verification: Tests green; QA checklist executed; no regressions reported.
- Definition of done: All planned tests implemented/passing; QA checklist signed off.
- Refer to `brds/UX-2026-01-14-order-form-revamp.md` whenever you need more clarity about UI/UX
- Refer to `brds/BRD-2026-01-14-order-form-revamp` whenever you need more about the business requirements.

## 4) API contracts (high level)

- Endpoints:
  - `POST /v1/businesses/:businessDescriptor/orders/preview` (new, manage + plan gated) — request body same as create; response `OrderPreviewResponse { items?, subtotal, discount, shippingFee, vat, total, currency, freeShipping: bool, discountType?, discountValue?, shippingZoneId?, shippingFeeMode }` (exact shape to align to service output and frontend needs).
  - Existing create/update unchanged.
- DTOs: Reuse `CreateOrderInput`-like structure for request; new response DTO for preview.
- Error cases: 400 validation; 401/403 auth/plan; 404 tenant; 409 stock conflict; problem+json format.
- Refer to `brds/UX-2026-01-14-order-form-revamp.md` whenever you need more clarity about UI/UX
- Refer to `brds/BRD-2026-01-14-order-form-revamp` whenever you need more about the business requirements.

## 5) Data model & migrations

- Tables/models: None changed.
- Indexing: None needed (read-only preview reuses existing lookups).
- Migration plan: None.

## 6) Security & privacy

- Tenant scoping: Business descriptor scoping with `business.EnforceBusinessValidity`.
- RBAC: `role.ResourceOrder` with `ActionManage` for preview; `ActionView` unaffected.
- Abuse prevention: Reuse throttle used for create (per biz+actor) on preview; ensure no data leakage in errors.

## 7) Observability & KPIs

- Events/metrics: Track preview usage and error categories in portal; backend logs preview calls and validation failures.
- Dashboards/alerts: None required now; rely on existing logging/metrics.

## 8) Test strategy

- E2E: Backend preview paths (happy, stock conflict, zone mismatch, payment disabled, plan gate).
- Integration/unit: Preview schema parsing in portal; preview hook staleness gating if test infra exists.
- Edge cases: max items=100; discount > subtotal; free shipping threshold; currency mismatch zone.

## 9) Risks & mitigations

- Risk: Divergence between preview and create/update logic.
  - Mitigation: Share service functions; add tests comparing outputs for same input (preview vs create dry execution).
- Risk: Mobile performance with frequent preview calls.
  - Mitigation: Debounce and cancel in-flight previews; throttle server.
- Risk: User confusion when status locks edits.
  - Mitigation: Clear lock chips/tooltips and disabled states.

## 10) Definition of done

- [ ] Meets BRD acceptance criteria
- [ ] Mobile-first UX verified
- [ ] RTL/i18n parity verified
- [ ] Multi-tenancy verified
- [ ] Error handling + empty/loading states complete
- [ ] No TODO/FIXME
