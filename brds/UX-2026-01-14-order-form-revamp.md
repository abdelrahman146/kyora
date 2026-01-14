---
status: draft
created_at: 2026-01-14
updated_at: 2026-01-14
brd_ref: BRD-2026-01-14-order-form-revamp.md
owners:
  - area: portal-web
    agent: UI/UX Designer
stakeholders:
  - name: Product
    role: PM
  - name: Design
    role: UX/UI
areas:
  - portal-web
---

# UX Spec: Order Create/Update Revamp — Guided, Calm, Mobile-first

## 0) Inputs & scope

- BRD: [brds/BRD-2026-01-14-order-form-revamp.md](brds/BRD-2026-01-14-order-form-revamp.md)
- Goals (what should feel better for the user):
  - Guided, decluttered flow with 3 stacked cards (Customer & Channel → Items & Discounts → Summary & Payment) and a sticky action bar.
  - Inline customer/address creation without leaving the sheet; never lose draft state.
  - Live, trustworthy summary (items, discount, shipping, VAT, total) powered by dry-run previews so there are zero surprises on submit.
  - Mobile-first, RTL-native experience with large tap targets and no hidden controls.
- Non-goals:
  - No new payment/checkout flow; no state machine changes.
  - No shipping zone/payment method management UI (use Business Settings).
- Assumptions:
  - Dry-run shares validation/totals code with create/update and is plan-gated the same.
  - Line price overrides remain allowed; stock is allocated only on create/update, not on dry-run.
  - Shipping zones and payment methods are already configured in Business Settings when reachable.
  - Draft-saving is implicit: the form state remains in memory while sheets are open; no separate draft persistence.

## 1) Reuse Map (Evidence)

List concrete existing files/components to reuse.

- Existing routes/pages:
  - Orders list page: [portal-web/src/routes/business/$businessDescriptor/orders/index.tsx](portal-web/src/routes/business/$businessDescriptor/orders/index.tsx) (entry to sheets, list refresh).
- Existing feature components:
  - Create sheet baseline: [portal-web/src/features/orders/components/CreateOrderSheet.tsx](portal-web/src/features/orders/components/CreateOrderSheet.tsx)
  - Edit sheet baseline: [portal-web/src/features/orders/components/EditOrderSheet.tsx](portal-web/src/features/orders/components/EditOrderSheet.tsx)
  - Review sheet: [portal-web/src/features/orders/components/OrderReviewSheet.tsx](portal-web/src/features/orders/components/OrderReviewSheet.tsx)
  - Order cards/list scaffolding: [portal-web/src/features/orders/components/OrderCard.tsx](portal-web/src/features/orders/components/OrderCard.tsx), [portal-web/src/features/orders/components/OrdersListPage.tsx](portal-web/src/features/orders/components/OrdersListPage.tsx)
- Existing shared components (components/lib):
  - BottomSheet: [portal-web/src/components](portal-web/src/components)
  - Shipping zone info helper: [portal-web/src/components/organisms/ShippingZoneInfo.tsx](portal-web/src/components/organisms/ShippingZoneInfo.tsx)
  - Form system: [portal-web/src/lib/form](portal-web/src/lib/form)
  - Currency/format helpers: [portal-web/src/lib/formatCurrency.ts](portal-web/src/lib/formatCurrency.ts), [portal-web/src/lib/formatDate.ts](portal-web/src/lib/formatDate.ts)

**Do-not-duplicate list** (must reuse these if applicable):
- BottomSheet / sheet patterns: BottomSheet (feature uses it already).
- Form fields/selects: use form system (`form.AppField`, `SelectField`, `PriceField`, `TextareaField`).
- Empty/loading/error state patterns: list/card skeletons ([portal-web/src/features/orders/components/OrderListSkeleton.tsx](portal-web/src/features/orders/components/OrderListSkeleton.tsx)), `form.FormError`, inline helper text.
- Resource list layout patterns: card list in OrdersListPage and OrderCard; keep the same list density and action placement.

## 2) IA + Surfaces

List every UI surface you’re specifying.

- Page: Orders list (existing) with entry buttons and success refresh.
- Sheet: Create Order (BottomSheet mobile, side-sheet desktop) with 3 stacked cards + sticky action bar.
- Sheet: Edit Order (inherits create layout; locks fields per status rules).
- Sheet: Inline Add Customer.
- Sheet: Inline Add Address.
- Sheet: Item Picker (add/edit item row with variant search and quantity/price).
- Section component: Live Summary card (inline, sticky on desktop).
- Banner: Dry-run preview banner (non-blocking, retry inline).
- Chip row: Preview status + last-updated time + currency shown above summary.

## 3) User flows (step-by-step)

### Flow A — Create order (guided cards)

1. From orders list, tap "Add order" → Create Order sheet slides up (mobile) with focus on Customer card.
2. Card 1: Customer & Channel — pick existing customer or add new inline; pick address inline; channel default Instagram; address selection auto-infers zone.
3. Card 2: Items & Discounts — tap "Add item" → Item Picker sheet; confirm → row appears with stock badge; repeat; toggle discount type (amount/percent) and enter value; validation runs immediately.
4. Card 3: Shipping & Payment — choose shipping mode (zone/manual), see computed fee/free badge; pick enabled payment method; optional reference.
5. Dry-run auto-fires on changes to items/discount/shipping/payment; preview pill shows timestamp; summary updates totals.
6. Sticky action bar shows primary CTA "Save order"; enabled only after successful preview + required fields; tap to submit → success toast → sheet closes → list refetch.

### Flow B — Edit order (status-aware)

1. From Order card quick actions or Review sheet, tap "Edit" → Edit Order sheet opens with same 3-card layout.
2. If status ∈ {shipped, fulfilled, cancelled, returned}, lock items/pricing; show notice chip "Editing limited"; allow address/discount/shipping only when permitted (per state machine); dry-run still allowed for visible fields.
3. Preview updates; sticky CTA "Save changes"; on submit success → toast → sheet closes → list refresh.

### Flow C — Inline add customer/address

1. Tap "Add customer" → Inline sheet with 3 fields; save → parent auto-selects customer and prompts address.
2. Tap "Add address" → Inline sheet; save → address selected; zone auto-inferred; summary preview reruns.

### Flow D — Dry-run retry

1. If preview fails (network/validation), show inline banner + "Try again"; submit disabled until preview succeeds.
2. User retries or adjusts fields → preview reruns → banner clears on success.

## 4) Per-surface specification (implementation-ready)

### Surface: Create Order Sheet (new decluttered layout)

- Entry points: Orders list primary CTA; empty state CTA.
- Layout structure:
  - Header strip: order number placeholder (for edit), business currency, preview status pill.
  - Card 1 (Customer & Channel): compact card with customer selector + "Add customer" button; address selector with inline "Add address"; channel buttons row (icon + label pills). Collapse/expand arrow (auto-expands when empty or error).
  - Card 2 (Items & Discounts): list of item rows; top-right "Add item" button; each row shows product thumbnail, name, variant chip, stock badge, quantity stepper, price override field; discount toggle (percent/amount) sits under the list in same card.
  - Card 3 (Shipping & Payment): radio for shipping mode (zone/manual); when zone, show zone card (cost, free threshold, currency lock); manual fee field inline; payment method select with enabled methods first; optional reference field; note textarea at bottom of card.
  - Summary card: visually separated card with totals rows, free-shipping badge, VAT line; sticky on desktop right; on mobile, placed above sticky action bar.
  - Sticky action bar: full-width primary CTA on mobile; secondary ghost button; shows last preview time.
- Layout annotations (mobile):
  - Bottom sheet height: 88% viewport; card stack scrolls; sticky action bar fixed to bottom with safe-area padding.
  - Spacing: `px-4 py-4`, `gap-4` between cards; card bodies `space-y-3`.
  - Channel pills: 2 columns on small screens; wrap with `gap-2`.
- Layout annotations (desktop md+):
  - Two-column grid: left 70% (cards stack), right 30% (sticky summary); `gap-4`.
  - Side sheet width ~640px if used; action bar anchors to sheet footer.
  - Keep card headers H3 (text-lg), rows text-sm; totals row text-xl.
- Primary CTA:
  - label: tOrders('create_submit')
  - enabled/disabled rules: requires customer, at least one item, successful preview in last 30s (stale after 30s or any change), no blocking validation errors; disable when mutation pending.
  - loading behavior: spinner inside button; bar stays visible; prevent double submit.
- Secondary actions: cancel (tCommon('cancel')); optional "Save pending" hidden by default (requires PM sign-off; default off).
- Fields (detailed):
  - Customer selector: searchable, shows name + phone; supports clear; button "Add customer" opens inline sheet; required.
  - Address selector: filtered by customer; required; shows country/city; "Add address" button; when none, show inline hint "Add a shipping address".
  - Channel: pill buttons with icons (Instagram, WhatsApp, TikTok, Facebook, In-person, Other); default Instagram; single select.
  - Items list: rows show thumbnail (if available), product name, variant chip, stock badge (warning color if <5), quantity stepper (min 1), unit price override field (PriceField, `dir="ltr"`), remove icon button; item total shown on right.
  - Discount: segmented control (Amount | Percent); input with validation (<= subtotal); helper text when clamped; percent input capped at 100.
  - Shipping: radio buttons (Shipping zone | Manual fee). Zone mode shows dropdown of zones (name + cost + free threshold). Manual fee uses PriceField min 0. Show computed shipping fee read-only line under selection.
  - Payment method: select showing enabled first; disabled ones greyed with reason "Enable in Business Settings"; cannot select disabled.
  - Payment reference: optional TextField, `dir="ltr"`, placeholder like "#1234".
  - Notes: TextareaField with character count, max 500; helper "Visible to your team".
- Empty states: Items card shows icon + "No items yet" + primary ghost "Add item" button. Address empty shows subtle callout with "Add address".
- Loading states: Skeleton cards for customer/address selects, items rows, summary rows; disable actions until base data loaded; show inline spinner on preview pill during dry-run.
- Error states (mapped to user-friendly copy categories):
  - Stock: inline under quantity "Quantity exceeds available stock"; highlight row.
  - Shipping zone mismatch: inline callout under address/zone.
  - Payment disabled: helper under select with link to settings (non-navigating tooltip note).
  - Discount too high: inline helper under discount input, clamps value.
  - Preview unavailable: top-of-summary banner + retry button.
- Success behavior: toast success; sheet closes; list + relevant queries refetch; last preview pill resets.
- Accessibility notes: BottomSheet focus trap; all icon-only buttons have aria-label; quantity steppers reachable by keyboard; phone/reference fields `dir="ltr"`; preview banner uses `role="alert"`; sticky bar reachable with tab.

### Surface: Edit Order Sheet (status-aware)

- Entry points: Order card quick actions, review sheet.
- Layout structure: mirrors Create layout; items/pricing card switches to read-only when status locked; status/payment chips displayed at top of header strip.
- Primary CTA:
  - label: tOrders('update_submit')
  - enabled/disabled rules: requires successful preview and allowed status; disable when status ∈ {shipped, fulfilled, cancelled, returned} for item/price changes (CTA still available if only shipping/discount editable and preview valid).
  - loading: spinner.
- Secondary actions: cancel.
- Fields: same as create; locked rows show lock icon + tooltip "Editing limited by status"; payment method select disabled if paymentStatus = paid (unless backend allows changes—follow state machine).
- States: same loading/error/preview behavior; notice chip in header when editing limited.
- Success behavior: toast, close, refetch list/detail.
- Layout annotations:
  - Locked sections use `bg-base-200` and `opacity-60`; keep text legible; add lock icon inline with section title.
  - Header strip shows status/payment chips aligned end; preview pill aligned start.

### Surface: Inline Add Customer Sheet

- Entry points: Customer card "Add customer" button.
- Layout structure: compact bottom sheet with tight spacing; max 3 inputs.
- Primary CTA: tCustomers('inline.save_customer'); disabled until name present; shows loading spinner.
- Fields:
  - Name (required) TextField, autoFocus, enterKeyHint="next".
  - Phone (optional) `type="tel"` `inputMode="tel"` `dir="ltr"` with country code helper if available.
  - Email (optional) `type="email"` autoComplete="email".
- States: inline validation messages; error banner on network failure; on success closes and selects customer; toast optional.
- Accessibility: focus trap; Escape/overlay closes; first field focused.
- Layout annotations: max height 60% viewport; `px-4 py-4`; fields `space-y-3`; CTA full width; secondary close via top-left/back icon.

### Surface: Inline Add Address Sheet

- Entry points: Address button when customer selected.
- Layout structure: bottom sheet; fields grouped: Location (country select, city, state), Details (street/line1 textarea with count), Contact (phone `dir="ltr"`), Default toggle.
- Primary CTA: tAddress('inline.save_address'); disabled until country + city + street + phone present and valid.
- Fields: country select uses business country as default; phone uses tel keyboard; default toggle text "Use for this order".
- States: inline validation; error callout for invalid country; success selects address, triggers zone inference + preview rerun.
- Accessibility: focus first field; respect RTL; phone LTR.
- Layout annotations: `space-y-3` groups; use section labels (Location/Details/Contact) as caption text-sm; CTA full width; maintain `px-4 py-4`.

### Surface: Item Picker Sheet

- Entry points: "Add item" button, edit row action.
- Layout structure: bottom sheet with 3 stacked controls: Variant select (searchable, shows stock + price), quantity stepper, unit price override (PriceField) + optional unit cost.
- Primary CTA: label switches between "Add item" / "Update item"; disabled until variant selected and quantity ≥1.
- Fields: variant select shows product thumbnail, variant name, stock; search by name/sku; quantity stepper min 1 max stock (if stock known); price override accepts decimal, `dir="ltr"`; unit cost optional.
- Error states: inline stock error; price must be >0; show helper if stock data missing (allow entry but warn).
- Success: closes, updates parent list, reruns preview.
- Layout annotations: three controls stacked with `space-y-3`; variant select at top; quantity and price in a 2-column grid on md+, stacked on mobile; CTA full width; sheet height ~70% viewport to show search results comfortably.

### Surface: Live Summary Section

- Entry points: inside create/edit sheets; sticky on desktop (right column) and above action bar on mobile.
- Layout structure: bordered card with labeled rows (subtotal, discount, shipping, VAT, total). Total row larger weight; free-shipping badge when fee=0; currency label at top right.
- Actions: preview pill (icon + "Updated HH:MM"), Refresh button when stale/error.
- States: skeleton rows on initial load; shimmer on updates; error banner with retry; stale state after 30s or form change.
- Success behavior: enables submit; shows last updated time and who triggered (auto/manual refresh not needed to display actor).
- Layout annotations: rows use `grid grid-cols-2` with labels start-aligned and values end-aligned (use logical properties); total row `text-lg font-bold`; badge inline next to shipping row; card padding `p-4` mobile, `p-5` desktop.

### Surface: Dry-run Banner

- Placement: above summary card and below header strip when preview failed/stale.
- Content: icon + message (orders.form.preview.unavailable) + "Try again" button; note that submit disabled until preview succeeds; banner disappears on next success.
- Layout annotations: full-width within sheet; `border border-warning bg-warning/10 px-3 py-2 rounded-lg`; icon size 18; button small `btn-ghost btn-sm` aligned inline end on desktop, stacked below message on mobile.

## 5) Responsiveness + RTL rules

- Mobile-first layout rules: one column; cards separated by `gap-4`; sticky action bar with `pb-safe` padding; BottomSheets for pickers.
- Tablet/desktop layout rules: md+ uses two-column split (70/30) with summary sticky; side-sheet instead of bottom sheet for main create/edit if viewport wide.
- RTL rules: only logical spacing classes (ms/me/ps/pe); rotate ArrowLeft/ArrowRight via `useLanguage`; icons inside pills follow text direction; channel handles and phones use `dir="ltr"` spans.
- `dir="ltr"` fields: phone, payment reference, order numbers, discount numeric, price/quantity, item totals.

## 6) Copy & i18n keys inventory

- Namespaces to use (prefer existing): `orders`, `common`, `errors`, `customers`, `address` (if not present add under customers/orders appropriately).
- New keys needed (en/ar parity required):
  - orders.form.header.preview_pill, orders.form.header.currency_label
  - orders.form.customer.add, orders.form.customer.select_placeholder, orders.form.address.add, orders.form.address.empty_hint
  - orders.form.channel.instagram, whatsapp, tiktok, facebook, in_person, other
  - orders.form.items.empty, orders.form.items.add, orders.form.items.edit, orders.form.items.stock_badge_low
  - orders.form.discount.label, orders.form.discount.amount, orders.form.discount.percent, orders.form.discount.too_high
  - orders.form.shipping.mode_zone, orders.form.shipping.mode_manual, orders.form.shipping.zone_note, orders.form.shipping.free_badge, orders.form.shipping.fee_label
  - orders.form.payment.method_disabled, orders.form.payment.reference_placeholder
  - orders.form.summary.subtotal, orders.form.summary.discount, orders.form.summary.shipping, orders.form.summary.vat, orders.form.summary.total, orders.form.summary.updated_at
  - orders.form.preview.unavailable, orders.form.preview.retry, orders.form.preview.stale
  - orders.form.status.locked, orders.form.status.editing_limited
  - customers.inline.create.title, customers.inline.create.success, customers.inline.create.phone_label
  - address.inline.create.title, address.inline.create.success, address.inline.create.default_toggle
  - common.cta.save_order, common.cta.save_changes, common.cta.save_customer, common.cta.save_address, common.cta.add_item, common.cta.refresh_preview

## 7) Component gaps / enhancements

### Missing components (if any)

- LiveSummaryCard (orders feature):
  - Why: centralize preview display (subtotal, discount, shipping, VAT, total, free badge, updated pill).
  - Where: `portal-web/src/features/orders/components/LiveSummaryCard.tsx`.
  - Reuse: Create + Edit sheets; optional in Review sheet.

- ItemPickerSheet (orders feature):
  - Why: cleaner item add/edit with search, stock, price override; reduces clutter in main form.
  - Where: `portal-web/src/features/orders/components/ItemPickerSheet.tsx`.
  - Reuse: Create/Edit; future quick-add.

- PreviewStatusPill (shared small component):
  - Why: show preview success/stale/error with timestamp and spinner.
  - Where: `portal-web/src/features/orders/components/PreviewStatusPill.tsx` (or shared if reused elsewhere).
  - Reuse: header strip, summary card.

### Enhancements to existing components

- CreateOrderSheet/EditOrderSheet:
  - Enhancement: adopt 3-card layout, integrate PreviewStatusPill + LiveSummaryCard, enforce preview-before-submit, inline customer/address sheets, item picker trigger per row, sticky action bar.
  - Call sites to update: orders list CTA and quick actions to pass new props if needed (e.g., onPreviewError callback).

- OrderReviewSheet:
  - Enhancement: optionally embed LiveSummaryCard for consistency; add action buttons matching new CTA labels.
  - Call sites: review triggers from OrderCard.

## 8) Acceptance checklist (for Engineering Manager)

- [ ] Reuses existing Kyora UI patterns; no parallel “second versions” of forms/sheets
- [ ] Calm UI: minimal sections, clear spacing, advanced options collapsed
- [ ] Mobile-first verified (one-handed, sticky primary CTA where appropriate)
- [ ] RTL-first verified; `dir="ltr"` applied where needed
- [ ] All empty/loading/error states specified and implemented
- [ ] i18n keys listed; en/ar parity planned and added
- [ ] Any missing components/enhancements explicitly documented
