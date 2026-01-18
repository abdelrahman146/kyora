---
description: Portal Web UX/UI System — Mobile-first, Arabic/RTL-first, Minimal
applyTo: "portal-web/**"
---

# Portal Web UX/UI Guidelines (SSOT)

Kyora is built for **Middle East social commerce entrepreneurs** (Instagram/WhatsApp/TikTok sellers). Users are **non-technical**, often working from a phone, and many prefer **Arabic (RTL)**.

**Design philosophy:** “Professional tools that feel effortless” — no accounting knowledge required.

**SSOT Hierarchy:**

- Parent: `copilot-instructions.md`
- Required Reading: `design-tokens.instructions.md`, `ui-implementation.instructions.md`, `forms.instructions.md`

---

## 1) Non-Negotiables

### Mobile-first (90% of usage)

- Default to **one-column layouts**.
- Default primary actions to **thumb-friendly, full-width** buttons.
- Default overlays to **BottomSheet** on mobile.
- Ensure all interactive elements meet **touch target minimums** (see `design-tokens.instructions.md`).

### Arabic/RTL-first

- UI must look “native” in Arabic: spacing, alignment, icon direction, and mixed-direction content all handled.
- Never use left/right utility classes; use logical properties (see `ui-implementation.instructions.md`).

### Minimal, elegant, calm

- **No shadows. No gradients.**
- Visual hierarchy comes from **spacing, typography, and borders**.
- Prefer fewer elements on screen; avoid dense layouts.

### Plain language (no accounting jargon)

- Use words like **“Profit”**, **“Cash in hand”**, **“Best seller”**.
- Avoid jargon like “Ledger”, “Accrual”, “EBITDA”, “COGS”.

---

## 2) Screen Structure & Hierarchy

### Page hierarchy

- **1 primary action** per screen (mobile). Secondary actions should be `btn-ghost` / `btn-outline`.
- Keep critical information above the fold: “What should I do next?”
- Prefer **cards** and **section blocks** over tables on mobile.

**Implementation examples:**

- ResourceListLayout: Single primary action in header (`btn btn-primary gap-2`)
- Empty states: Single primary action with icon (`btn btn-primary btn-sm gap-2`)
- BottomSheet footers: Primary action visually distinguished from secondary

### Spacing and density

- Use the existing 4px spacing grid (see `design-tokens.instructions.md`).
- Avoid "tight" UIs: breathing room increases comprehension for non-technical users.

**Common spacing patterns:**

- Container padding: `px-4` (16px), `py-4` (16px) on mobile
- Card body: uses daisyUI `card-body` class
- Flex gaps: `gap-2` (8px), `gap-4` (16px)

### Lists

- Lists should be scannable: title + 1–2 supporting lines, optional right-side meta.
- Entire card row can be clickable, but don't hide the primary action.

**Implementation:**

- ResourceListLayout provides configurable list view with cards
- Table component used for desktop-optimized data display
- Mobile-first card layouts in features (inventory, customers, orders)

---

## 3) Mobile Interactions

### Touch targets

- Touch targets must be at least **44×44px** (prefer 48px+).
- Do not hardcode pixel heights; use Tailwind scale (`min-h-12`, `min-w-12`, etc.) and daisyUI sizing (`btn-lg`).

**Implementation examples:**

- `BottomNav` uses `min-h-11` (44px) for navigation items (`portal-web/src/components/organisms/BottomNav.tsx`)
- Buttons use daisyUI sizing: `btn-sm`, `btn` (default), `btn-lg`
- Interactive elements in cards/lists ensure adequate spacing for touch

### Primary actions

**Current implementation patterns:**

- Primary CTA: `btn btn-primary` with size variants (`btn-sm`, `btn-lg`)
- Full-width mobile CTAs: `btn btn-primary btn-lg w-full` (used in auth flows)
- Bottom sheet footers: actions in flex container, Submit last and visually primary

**Example from production code:**

```tsx
// BottomSheet footer pattern
<div className="flex gap-2">
  <Button variant="ghost" className="flex-1" onClick={onClose}>
    {t("common:cancel")}
  </Button>
  <form.SubmitButton className="flex-1" variant="primary">
    {t("common:save")}
  </form.SubmitButton>
</div>
```

### Overlays

**Mobile: BottomSheet** (`portal-web/src/components/molecules/BottomSheet.tsx`)

- Slides up from bottom on mobile (< 768px)
- Side drawer on desktop (≥ 768px)
- Focus trap and Escape to close built-in
- Prevents body scroll when open
- Used for forms, filters, notes, and CRUD operations

**Desktop: Dialog** (`portal-web/src/components/molecules/Dialog.tsx`)

- Centered modal on all screen sizes
- Used for confirmations, alerts, and simple dialogs
- Also supports focus trap and Escape to close

**Modal** (`portal-web/src/components/molecules/Modal.tsx`)

- Alternative modal component with different API
- Portal-based rendering
- Used in some legacy components

**Rule:** All overlays must trap focus, support Escape to close, and prevent body scroll.

### Keyboard UX

- Use proper `inputMode` and `autoComplete`.
- Avoid layouts that break when the keyboard opens; keep submit reachable.

---

## 4) RTL + Mixed Direction Content

### Mixed content rules

- Phone numbers, IBANs, order IDs, coupon codes, tracking numbers: render as LTR:

```tsx
<span dir="ltr">{orderId}</span>
```

**Implementation examples from codebase:**

- Phone numbers: `<span dir="ltr">{phone}</span>` (AddressCard, CustomerDetailPage)
- Price inputs: `dir="ltr"` on input element (PriceInput, QuantityInput)
- Social media handles: `dir="ltr"` (SocialMediaHandles, SocialMediaInputs)
- Order totals: `<span dir="ltr">{amount}</span>` (OrderReviewSheet)

- Currency/amounts should be formatted via locale-aware helpers (prefer `Intl.NumberFormat`).
- Don't rely on string concatenation that breaks in RTL. Prefer translation templates.

### Icon direction

- Back/forward arrows must match reading direction (see `ui-implementation.instructions.md`).
- Don't rotate icons that are not directional.

**Implementation:**

```tsx
import { useLanguage } from '@/hooks/useLanguage'

const { isRTL } = useLanguage()

// ✅ CORRECT - Directional arrows rotate
<ArrowLeft className={isRTL ? 'rotate-180' : ''} />
<ArrowRight className={isRTL ? 'rotate-180' : ''} />

// ✅ CORRECT - Chevrons auto-flip (no rotation)
<ChevronLeft />
<ChevronRight />

// ✅ CORRECT - Non-directional icons never rotate
<Plus />
<X />
<Check />
```

---

## 5) Copy, Labels, and i18n

- Never hardcode user-facing strings.
- Arabic copy should be **short**, **direct**, and **non-technical**.
- Prefer **verb-first CTAs**: “Add order”, “Save”, “Mark as paid”.
- Error messages: what happened + how to fix (no stack-trace vibes).
  **Translation enforcement:**

- All validation messages must use `validation.*` keys from `src/i18n/*/errors.json`
- All UI labels must use translation keys (never hardcode English/Arabic)
- See `.github/instructions/i18n-translations.instructions.md` for complete rules
- See `.github/instructions/forms.instructions.md` for form-specific translation requirements

---

## 6) Feedback & States

### Loading

- Prefer skeletons to spinners for page-level loading.
- Keep layout stable (avoid jumping content).

**Implementation patterns:**

- ResourceListLayout uses skeleton placeholders: `<div className="skeleton h-32 rounded-box" />`
- Table component uses `Skeleton` atoms for row-level loading
- Custom skeletons per feature: `InventoryListSkeleton`, `CustomerListSkeleton`, `CustomerDetailSkeleton`
- Suspense boundaries use feature-specific skeleton components

### Saving

- Buttons show loading state; prevent double-submit.
- On success: toast + UI updates.

**Implementation:**

- Form buttons use `disabled={mutation.isPending}` or `disabled={isSubmitting}`
- Loading indicators: `<Loader2 className="w-5 h-5 animate-spin" />`
- Success toasts via global error handler (see `.github/instructions/http-tanstack-query.instructions.md`)

### Empty states

- Always include: what it means + a clear next action.

**Implementation pattern (from ResourceListLayout):**

```tsx
<div className="card bg-base-100 shadow">
  <div className="card-body items-center text-center">
    {emptyIcon && <div className="mb-4">{emptyIcon}</div>}
    <h3 className="card-title">{emptyTitle}</h3>
    <p className="text-base-content/70 mb-4">{emptyMessage}</p>
    {emptyActionText && onEmptyAction && (
      <button onClick={onEmptyAction} className="btn btn-primary btn-sm gap-2">
        <Plus size={16} />
        {emptyActionText}
      </button>
    )}
  </div>
</div>
```

---

## 7) Visual Style Constraints

- Use daisyUI semantic classes (`btn`, `card`, etc.).
- Use Tailwind utilities for spacing/layout only.
- **No shadows, no gradients.** Separation uses `border border-base-300` and background surfaces (`bg-base-100`, `bg-base-200`).
- Focus styles must remain visible (accessibility).

**Card pattern (standard):**

```tsx
<div className="card bg-base-100 border border-base-300">
  <div className="card-body">{/* content */}</div>
</div>
```

**Known drift:** Some components use shadows (documented in `backlog/drifts/2026-01-18-portal-web-shadow-usage-violates-design-tokens.md`). Do not replicate this pattern in new code.

---

## 8) Portal Web UI Review Checklist

Before finishing any portal-web UI work:

- **Mobile-first** layout and interactions (one-column, full-width CTAs where appropriate, BottomSheet for forms/actions)
- **Touch targets** meet minimum (no hardcoded px heights, use Tailwind scale: `min-h-11`+, `min-w-11`+)
- **RTL:** no left/right utilities; directional icons correct; mixed-direction content uses `dir="ltr"` where needed
- **Copy:** translation keys used; plain language; no accounting jargon
- **Loading states:** skeletons for page-level loading, stable layouts
- **Empty states:** include what it means + clear CTA
- **Error states:** present and user-friendly (via global error handler or inline)
- **Visual style:** `card bg-base-100 border border-base-300` pattern; no shadows; no gradients
- **Accessible:** focus visible, aria-labels for icon-only buttons, dialogs/sheets trap focus
- **Components:** Use BottomSheet for mobile forms/actions, Dialog for confirmations, ResourceListLayout for lists

**Implementation references:**

- BottomSheet: `portal-web/src/components/molecules/BottomSheet.tsx`
- Dialog: `portal-web/src/components/molecules/Dialog.tsx`
- ResourceListLayout: `portal-web/src/components/templates/ResourceListLayout.tsx`
- Form system: `.github/instructions/forms.instructions.md`
- RTL patterns: `.github/instructions/ui-implementation.instructions.md`
