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

### Spacing and density

- Use the existing 4px spacing grid (see `design-tokens.instructions.md`).
- Avoid “tight” UIs: breathing room increases comprehension for non-technical users.

### Lists

- Lists should be scannable: title + 1–2 supporting lines, optional right-side meta.
- Entire card row can be clickable, but don’t hide the primary action.

---

## 3) Mobile Interactions

### Touch targets

- Touch targets must be at least **44×44px** (prefer 48px+).
- Do not hardcode pixel heights; use Tailwind scale (`min-h-12`, `min-w-12`, etc.) and daisyUI sizing (`btn-lg`).

### Primary actions

- Primary CTA on mobile is typically `btn btn-primary btn-lg w-full`.
- When actions live at the bottom (BottomSheet), keep **Submit** last and visually primary.

### Overlays

- Mobile: BottomSheet.
- Desktop: Dialog where appropriate.
- All overlays must trap focus and support Escape to close.

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

- Currency/amounts should be formatted via locale-aware helpers (prefer `Intl.NumberFormat`).
- Don’t rely on string concatenation that breaks in RTL. Prefer translation templates.

### Icon direction

- Back/forward arrows must match reading direction (see `ui-implementation.instructions.md`).
- Don’t rotate icons that are not directional.

---

## 5) Copy, Labels, and i18n

- Never hardcode user-facing strings.
- Arabic copy should be **short**, **direct**, and **non-technical**.
- Prefer **verb-first CTAs**: “Add order”, “Save”, “Mark as paid”.
- Error messages: what happened + how to fix (no stack-trace vibes).

---

## 6) Feedback & States

### Loading

- Prefer skeletons to spinners for page-level loading.
- Keep layout stable (avoid jumping content).

### Saving

- Buttons show loading state; prevent double-submit.
- On success: toast + UI updates.

### Empty states

- Always include: what it means + a clear next action.

---

## 7) Visual Style Constraints

- Use daisyUI semantic classes (`btn`, `card`, etc.).
- Use Tailwind utilities for spacing/layout only.
- **No shadows, no gradients.** Separation uses `border border-base-300` and background surfaces (`bg-base-100`, `bg-base-200`).
- Focus styles must remain visible (accessibility).

---

## 8) Portal Web UI Review Checklist

Before finishing any portal-web UI work:

- Mobile-first layout and interactions (one-column, full-width CTAs, BottomSheet)
- RTL: no left/right utilities; directional icons correct
- Copy: translation keys used; plain language; no accounting jargon
- Touch targets meet minimum (no hardcoded px heights)
- Loading/empty/error states present
- Accessible: focus visible, aria-labels for icon-only buttons, dialogs trap focus
