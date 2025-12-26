---
description: KYORA DESIGN SYSTEM (KDS)
applyTo: "portal-web/**,storefront-web/**"
---

# KYORA DESIGN SYSTEM (KDS) v1.0

**Target Platform:** Mobile Web (PWA) & Mobile Native (React Native)
**Primary Audience:** Middle East (Arabic First / RTL)
**Core Philosophy:** "Complexity made invisible."

---

## 1. Brand Identity & Voice

### 1.1. Core Values

1. **Tamkeen (Empowerment):** We turn small merchants into business tycoons. The UI must feel professional, not "cute."
2. **Basata (Simplicity):** Our users are busy. They stand in shops all day. Every tap must be purposeful. No hidden menus.
3. **Thiqah (Trust):** We handle money and inventory. The design must feel stable, secure, and precise.

### 1.2. Voice & Tone (Arabic Centric)

- **Language:** Modern Standard Arabic (MSA) but approachable (White Arabic). Avoid overly classical terms; use business-standard terms used in the GCC/Levant.
- **Tone:** Direct, warm, and respectful.
- **English Fallback:** English is the secondary citizen. UI layouts are designed for RTL (Right-to-Left) first. LTR (Left-to-Right) is the mirror.

---

## 2. Design Tokens (The DNA)

All AI agents must strictly utilize these tokens. Do not hardcode magic numbers.

**Source of truth in this repo (storefront-web):**

- daisyUI theme tokens are defined in `storefront-web/src/index.css` via `@plugin "daisyui/theme"`.
- KDS raw CSS variables (e.g. `--primary-600`) also exist in `storefront-web/src/index.css` under `:root`.
- Prefer daisyUI semantic tokens in components: `bg-base-100`, `text-base-content`, `btn btn-primary`, etc.

### 2.1. Color Palette

The palette is inspired by the Middle East landscape (Teal/Sea) and Modern Fintech (Trust).

**Primary (The Brand)**

- `primary-50`: `#F0FDFA` (Backgrounds)
- `primary-100`: `#CCFBF1` (Interactive bg)
- `primary-500`: `#14B8A6` (Icons, Highlights)
- `primary-600`: `#0D9488` (Main Actions - **Key Color**)
- `primary-700`: `#0F766E` (Active States)
- `primary-900`: `#134E4A` (Text on Light)

**Secondary (Action/Accent - Gold/Sand)**

- `secondary-400`: `#FACC15`
- `secondary-500`: `#EAB308` (CTAs requiring attention)
- `secondary-600`: `#CA8A04`

**Neutral (Text & Surface)**

- `neutral-0`: `#FFFFFF` (Card Backgrounds)
- `neutral-50`: `#F8FAFC` (App Background - Off White to reduce glare)
- `neutral-200`: `#E2E8F0` (Borders/Dividers)
- `neutral-500`: `#64748B` (Secondary Text)
- `neutral-900`: `#0F172A` (Primary Text - Almost Black)

**Semantic**

- `success`: `#10B981` (Completed orders, Paid invoices)
- `error`: `#EF4444` (Failed payment, Out of stock)
- `warning`: `#F59E0B` (Low stock, Subscription ending)

**Third-party brand colors (exception):**

- WhatsApp green `#25D366` is allowed only for WhatsApp-specific UI (e.g., the floating WhatsApp button).

### 2.2. Typography (Arabic First)

**Font Family:** `IBM Plex Sans Arabic` (Google Fonts). It is geometric, modern, and legible at small mobile sizes.
**Fallback:** `Almarai`.

**Scale (Mobile)**

- `Display`: 32px / 1.2 / Bold (Marketing Headers)
- `H1`: 24px / 1.3 / Bold (Page Titles)
- `H2`: 20px / 1.3 / SemiBold (Section Headers)
- `H3`: 18px / 1.4 / Medium (Card Titles)
- `Body-L`: 16px / 1.5 / Regular (Default Text)
- `Body-M`: 14px / 1.5 / Regular (Secondary descriptions)
- `Caption`: 12px / 1.4 / Medium (Labels, Timestamps)
- `Micro`: 10px / 1.2 / Bold (Badges)

### 2.3. Spacing & Grid (4px Baseline)

- **Safe Area:** Left/Right padding is always `16px` (mobile standard).
- **Gap-XS:** 4px
- **Gap-S:** 8px
- **Gap-M:** 16px
- **Gap-L:** 24px
- **Gap-XL:** 32px
- **Touch Target:** Minimum tappable area is **44x44px**.

### 2.4. Radius (Soft & Friendly)

- `rounded-sm`: 4px (Checkboxes, Tags)
- `rounded-md`: 8px (Inner cards, Inputs)
- `rounded-lg`: 12px (Standard Cards, Modals)
- `rounded-xl`: 16px (Bottom Sheets)
- `rounded-full`: 9999px (Buttons, Avatars)

### 2.5. Shadows (Elevation)

- `shadow-sm`: `0 1px 2px 0 rgb(0 0 0 / 0.05)` (Cards)
- `shadow-float`: `0 10px 15px -3px rgb(0 0 0 / 0.1)` (Sticky Buttons, Bottom Sheets)

---

## 3. Global UX Behaviors (Mobile First)

### 3.1. Right-to-Left (RTL) Mechanics

Since Arabic is 1st class, the AI **must** implement logical properties.

- **NEVER use:** `margin-left`, `padding-right`, `float-left`.
- **ALWAYS use:** `margin-inline-start`, `padding-inline-end`, `text-align: start`.
- **Icons:** Directional icons (arrows, chevrons, back buttons) must be mirrored.
- _Example:_ A "Next" arrow points Left in LTR, but must point **Left** in RTL? **NO.**
- _Correction:_ In RTL, the flow is Right -> Left. A "Next" arrow points **Left** (towards the future). A "Back" arrow points **Right** (towards the start).

- **Phone Numbers:** Always display LTR standard (e.g., +971 50...), even in Arabic text blocks.

### 3.2. Navigation Topology (The "Thumb Zone")

- **Primary Navigation:** Prefer simple, thumb-friendly navigation. If a tab bar exists, labels must be present (Icons + Text).
- **Secondary Actions:** Floating Action Button (FAB) or Sticky Bottom Bar.
- **Back Navigation:** Top-Right (in RTL) header arrow.
- **Modals:** DO NOT use centered modals. Use **Bottom Sheets** (Slide up panels) for filters, confirmations, and quick forms. This is easier for thumb reach.

### 3.3. Input & Forms (The "Fat Finger" Rule)

- **Inputs:** Minimum height `50px`.
- **Labels:** Use "Floating Labels" or Top-aligned labels. Never use placeholder text as the only label.
- **Keyboard:**
- Phone fields -> Numeric keypad.
- Email fields -> Email keyboard.
- Search -> 'Search' action key.

- **Validation:** Inline, immediate validation. Red border + text message below input.

---

## 4. Component Dictionary (Detailed Specs)

AI Agents: Implement components exactly as described below.

### 4.1. Buttons (`Button`)

- **Primary Button:**
- Bg: `primary-600`
- Text: `neutral-0` / SemiBold
- Height: `52px` (Full width on mobile usually)
- Radius: `rounded-xl`
- State: `active:scale-95` (Press animation)

- **Secondary Button:**
- Bg: `primary-50`
- Text: `primary-700`

- **Ghost Button:**
- Bg: Transparent
- Text: `neutral-500`

**Implementation note (storefront-web):** Prefer daisyUI button variants (`btn`, `btn-primary`, `btn-secondary`, `btn-ghost`, `btn-outline`) and reuse the shared interaction utilities (`focus-ring`, `active-scale`).

### 4.2. Cards (`Card`)

- **Container:** `bg-white`, `rounded-lg`, `border border-neutral-100`, `shadow-sm`.
- **Padding:** `p-4` (16px).
- **Clickable Cards:** Should have a `chevron-left` (in RTL) icon to indicate navigation.

### 4.3. List Items (`ListItem`)

Used for Inventory lists, Customer lists, etc.

- **Layout:** Flex row.
- **Left (RTL Start):** Image/Icon (40x40px, rounded-md, object-cover).
- **Middle:** Title (Body-L, SemiBold) + Subtitle (Body-M, neutral-500).
- **Right (RTL End):** Value/Price (SemiBold) or Status Badge.
- **Separator:** Full width or indented divider `border-b-neutral-100`.

### 4.4. Input Fields (`Input`)

- **Base:** `bg-neutral-50`, `border-transparent`, `focus:bg-white`, `focus:border-primary-500`, `focus:ring-2`.
- **RTL Specific:** Text aligns Right. Caret aligns Right.
- **Icons:** Input may have `startIcon` (Payment card icon) or `endIcon` (Eye toggle).

### 4.5. Bottom Sheet (`BottomSheet`)

- **Behavior:** Slides up from bottom. Backdrop blur.
- **Handle:** Gray pill at the top center (`w-12 h-1 bg-neutral-300 rounded-full`).
- **Content:** Scrollable content area.
- **Footer:** Sticky container for "Save" or "Apply" buttons.

**Implementation note (storefront-web):** Keep it keyboard accessible (Escape to close) and restore scroll on close via effect cleanup.

### 4.6. Skeleton Loaders

- **Mandatory:** Never show white screens. Use pulsing gray blocks (`animate-pulse bg-neutral-200`) matching the shape of the content loading (Circle for avatars, Rect for text).

---

## 5. Storefront Specific Guidelines (The Revamp)

The User-Facing Storefront (`storefront-web`) needs to feel like a high-end native app (PWA).

### 5.1. Storefront Header

Current implementation uses a minimal, safe-area-aware header with:

- Language switcher on one side
- Cart button (with badge) on the other side

Brand identity (logo/name) is displayed in a separate centered brand header below.

### 5.2. Product Card (The Grid)

- **Layout:** 2 Columns on Mobile (Grid).
- **Image:** Aspect Ratio `1:1` or `3:4`. `rounded-lg`.
- **Title:** Max 2 lines. `text-sm`.
- **Price:** Bold. `text-primary-700`.
- **Add Button:** A distinct circular button with a "+" icon floating on the bottom-left (RTL) of the image, OR a full-width "Add" button below. **Decision:** Floating button on image corner for cleaner look.

### 5.3. Sticky Cart Bar

- If items > 0 in cart, show a sticky bar at the bottom of the screen (above tab bar if exists).
- **Content:** "3 Items • 150 AED" (Right) .... "View Cart >" (Left).
- **Color:** `bg-primary-900` text white.

### 5.4. WhatsApp Integration (Critical)

- **Floating Button:** Use logical positioning (Tailwind `end-*`) so it naturally flips in RTL.
- **Product Page:** "Order via WhatsApp" button variant (secondary action).
- **Message Pre-fill:** "Hi [Store Name], I am interested in [Product Name]..."

**Implementation note (storefront-web):** The floating WhatsApp FAB uses WhatsApp green `#25D366` (allowed exception).

### 5.5. Checkout Flow

1. **Phone Number First:** No email requirement. Just phone.
2. **Location:** Use "Current Location" GPS button to fill address.
3. **Payment:**

- Option 1: Cash on Delivery (COD) - Default for many.
- Option 2: Online Payment (Stripe).

4. **Success:** Digital Receipt animation.

---

## 6. Implementation Rules for AI Agents

When writing code (React SPA):

1. **Tailwind + daisyUI (storefront-web):** Tailwind v4 is configured CSS-first. Use daisyUI semantic tokens and the theme defined in `storefront-web/src/index.css`. Do not introduce new colors in component classNames.
2. **Icons:** Prefer `lucide-react` for consistency.
3. **RTL & logical properties:** Prefer logical utilities (`start-*`, `end-*`, `text-start`). For directional icons, apply the existing `rtl-mirror` class.
4. **Dates/Currency:**

- Use `Intl.NumberFormat` with `ar-AE` (or relevant country) for currency (e.g., "د.إ.‏ 150.00").
- Use Hijri/Gregorian toggle if requested, but Gregorian (English numbers) is standard for business in most GCC apps.

5. **Error Handling:**

- Network Error -> Show "Retry" button component.
- 404 -> Friendly illustration + "Go Home" button.

6. **Files:** Keep components small and composable; prefer shared atoms/molecules/organisms where they already exist.

---

## 7. Accessibility (A11y)

1. **Contrast:** Ensure `primary-600` is used on white. Do not use `primary-400` for text.
2. **Screen Readers:** All images must have `alt` text. Use `aria-label` for icon-only buttons (like the Hamburger menu or Cart icon).
3. **Focus:** Never remove outline unless replacing with a custom ring.

---

## 8. Animation Guidelines

- **Page Transitions:** Fade in + Slide Up slightly (`y-4` to `y-0`).
- **Modals/Sheets:** Spring physics (damping: 20, stiffness: 300).
- **Buttons:** `active:scale-95`.
- **Lists:** Staggered fade-in for items.
