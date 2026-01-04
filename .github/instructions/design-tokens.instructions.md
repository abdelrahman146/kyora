---
description: Kyora Design Tokens - Single Source of Truth
applyTo: "portal-web/**,storefront-web/**"
---

# Design Tokens Reference

**Philosophy**: "Professional tools that feel effortless" for Middle East entrepreneurs
**Primary Language**: Arabic (RTL-first), English (LTR)
**Token Location**: `/portal-web/src/index.css` (daisyUI theme), `/storefront-web/src/index.css`

---

## Colors

### Semantic Palette

| Token       | Value                  | daisyUI Class      | Usage                           |
| ----------- | ---------------------- | ------------------ | ------------------------------- |
| Primary     | `#0D9488` (Teal)       | `bg-primary`       | Brand, main actions, links      |
| Secondary   | `#EAB308` (Gold)       | `bg-secondary`     | Secondary CTAs, accents         |
| Accent      | `#FACC15` (Yellow)     | `bg-accent`        | Attention, badges               |
| Success     | `#10B981` (Green)      | `bg-success`       | Completed, paid, active         |
| Error       | `#EF4444` (Red)        | `bg-error`         | Failed, invalid, critical       |
| Warning     | `#F59E0B` (Orange)     | `bg-warning`       | Low stock, caution              |
| Info        | `oklch(60% 0.12 192)`  | `bg-info`          | Informational                   |
| Base-100    | `#FFFFFF`              | `bg-base-100`      | Main backgrounds, cards         |
| Base-200    | `#F8FAFC`              | `bg-base-200`      | App background                  |
| Base-300    | `#E2E8F0`              | `bg-base-300`      | Borders, dividers               |
| Base-Content| `#0F172A`              | `text-base-content`| Primary text                    |

### Text Opacity Scale

```
text-base-content       (100% - primary text)
text-base-content/70    (secondary text)
text-base-content/60    (tertiary text)
text-base-content/50    (icons)
text-base-content/40    (disabled, placeholders)
text-base-content/30    (subtle decorative)
```

---

## Typography

**Font**: `IBM Plex Sans Arabic`
**Fallback**: `Almarai, -apple-system, sans-serif`

| Scale    | Size | Weight   | Usage                        | Tailwind Class |
| -------- | ---- | -------- | ---------------------------- | -------------- |
| Display  | 32px | Bold     | Marketing headers            | `text-[32px]`  |
| H1       | 24px | Bold     | Page titles                  | `text-2xl`     |
| H2       | 20px | SemiBold | Section headers              | `text-xl`      |
| H3       | 18px | Medium   | Card titles                  | `text-lg`      |
| Body-L   | 16px | Regular  | Default text                 | `text-base`    |
| Body-M   | 14px | Regular  | Secondary text               | `text-sm`      |
| Caption  | 12px | Medium   | Labels, timestamps           | `text-xs`      |
| Micro    | 10px | Bold     | Small badges                 | `text-[10px]`  |

**Line Heights**: Headings `1.2-1.3`, Body `1.5-1.6`

---

## Spacing (4px Grid)

```
gap-1  → 4px     gap-2  → 8px     gap-3  → 12px
gap-4  → 16px (standard)          gap-6  → 24px
gap-8  → 32px
```

**Mobile Safe Area**: `px-4` (16px horizontal), `py-4` (16px vertical)

---

## Border Radius

```
rounded-sm   → 4px  (tags, small elements)
rounded-md   → 8px  (default inner elements)
rounded-lg   → 12px (cards, modals, containers)
rounded-xl   → 16px (buttons, bottom sheets)
rounded-box  → 12px (daisyUI card utility)
rounded-full → pills, avatars, circular buttons
```

---

## Shadows

```
shadow-sm  → Subtle card elevation
shadow     → Standard cards
shadow-md  → Elevated elements
shadow-lg  → Floating elements, dropdowns
shadow-xl  → Modals, important overlays
```

---

## Borders

```
Default: border border-base-300
Focus:   focus:border-primary
Error:   border-error
Width:   1px (default)
```

---

## Transitions

```
Standard: transition-all duration-200
Colors:   transition-colors duration-200
Opacity:  transition-opacity duration-200
Shadow:   transition-shadow duration-200
```

---

## Responsive Breakpoints

```
sm: 640px  (tablets portrait)
md: 768px  (tablets landscape, small desktops)
lg: 1024px (desktops)
xl: 1280px (large desktops)
```

**Hook**: `useMediaQuery("(max-width: 768px)")` for JS logic

---

## Touch Targets

**Minimum**: 44-50px height
- Buttons: `h-[52px]` (default)
- Inputs: `h-[50px]`
- Icon buttons: `min-w-[44px] min-h-[44px]`

---

## Animations

```
Pulse (skeletons):    animate-pulse
Spin (loaders):       animate-spin
Press (buttons):      active:scale-95
Hover (cards):        hover:shadow-md transition-shadow
```

---

## Icon Sizing

```
Inline with text:     size={18} or size={20}
Large interactive:    size={24}
Headers:              size={28} to size={32}
Empty states:         size={48} to size={64}
```

**Icon Colors**:
```
Default:     text-base-content/50
Hover:       group-hover:text-primary
Active:      text-primary
Disabled:    text-base-content/30
Semantic:    text-success, text-error, text-warning
```
