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

- **Primary Navigation:** Bottom Tab Bar (Fixed). Labels must be present (Icons + Text).
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

- **Behavior:** Slides up from bottom. Backdop blur.
- **Handle:** Gray pill at the top center (`w-12 h-1 bg-neutral-300 rounded-full`).
- **Content:** Scrollable content area.
- **Footer:** Sticky container for "Save" or "Apply" buttons.

### 4.6. Skeleton Loaders

- **Mandatory:** Never show white screens. Use pulsing gray blocks (`animate-pulse bg-neutral-200`) matching the shape of the content loading (Circle for avatars, Rect for text).

---

## 5. Storefront Specific Guidelines (The Revamp)

The User-Facing Storefront (`storefront-web`) needs to feel like a high-end native app (PWA).

### 5.1. Storefront Header

- **Left (RTL):** Hamburger Menu (Categories).
- **Center:** Brand Logo.
- **Right (RTL):** Search Icon + Cart Bag (with Badge).
- **Behavior:** Sticky on scroll. transforms to valid title on scroll down.

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

- **Floating Button:** Fixed bottom-left (RTL) or bottom-right depending on layout.
- **Product Page:** "Order via WhatsApp" button variant (secondary action).
- **Message Pre-fill:** "Hi [Store Name], I am interested in [Product Name]..."

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

1. **Tailwind Config:** Use the color tokens defined above in `tailwind.config.js`.
2. **Icons:** Use `lucide-react`.

- _Rule:_ For RTL, wrap icons in a component that applies `transform: scaleX(-1)` only for directional icons (arrow-left, chevron-right, etc.).

3. **Dates/Currency:**

- Use `Intl.NumberFormat` with `ar-AE` (or relevant country) for currency (e.g., "د.إ.‏ 150.00").
- Use Hijri/Gregorian toggle if requested, but Gregorian (English numbers) is standard for business in most GCC apps.

4. **Error Handling:**

- Network Error -> Show "Retry" button component.
- 404 -> Friendly illustration + "Go Home" button.

5. **Files:** When creating a component, strictly separate logic (hook) from view (JSX).

### Example Component Structure (React + Tailwind)

```tsx
// components/Button.tsx
import { cva } from "class-variance-authority";
import { Loader2 } from "lucide-react";

// Define variants strictly based on Design System
const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-xl text-sm font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 disabled:pointer-events-none disabled:opacity-50 h-12 px-6 w-full active:scale-95",
  {
    variants: {
      variant: {
        primary: "bg-primary-600 text-white hover:bg-primary-700 shadow-sm",
        secondary: "bg-primary-50 text-primary-900 hover:bg-primary-100",
        outline:
          "border border-neutral-200 bg-transparent hover:bg-neutral-50 text-neutral-900",
        ghost: "hover:bg-neutral-100 text-neutral-700",
      },
      size: {
        default: "h-12 px-6",
        sm: "h-9 rounded-lg px-3",
        lg: "h-14 rounded-2xl px-8 text-base", // Fat finger friendly
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "default",
    },
  }
);

// ... implementation details
```

---

## 7. Accessibility (A11y)

1. **Contrast:** Ensure `primary-600` is used on white. Do not use `primary-400` for text.
2. **Screen Readers:** All images must have `alt` text. Use `aria-label` for icon-only buttons (like the Hamburger menu or Cart icon).
3. **Focus:** Never remove outline unless replacing with a custom ring.

---

## 8. Dashboard / Merchant App Specifics

The "Admin" side for the merchant.

### 8.1. Dashboard Home

- **Greeting:** "Sabah el Kheir, [Name]" (Time aware).
- **Quick Actions Grid:** 4 buttons at top (New Sale, Add Product, Share Store, Expenses).
- **Stats Cards:** Scrolled horizontally. (Sales Today, Orders Pending).

### 8.2. Inventory Management

- **Visuals:** Small thumbnails are mandatory.
- **Stock Levels:** Color coded.
- 0: Red background/text badge "Out of Stock".
- < 5: Orange badge "Low Stock".
- 5+: Green text.

### 8.3. Order Processing

- **Status Timeline:** Vertical stepper (RTL).
- Pending (Yellow) -> Processing (Blue) -> Shipped (Purple) -> Delivered (Green).

- **WhatsApp Action:** Button to "Send Receipt" or "Ask Location" via WhatsApp directly from the Order Detail view.

---

## 9. Animation Guidelines

- **Page Transitions:** Fade in + Slide Up slightly (`y-4` to `y-0`).
- **Modals/Sheets:** Spring physics (damping: 20, stiffness: 300).
- **Buttons:** `active:scale-95`.
- **Lists:** Staggered fade-in for items.
