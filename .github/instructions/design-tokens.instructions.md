---
description: Kyora Design Tokens - Single Source of Truth
applyTo: "portal-web/**,storefront-web/**"
---

# Design Tokens Reference

**Philosophy**: "Professional tools that feel effortless" for Middle East entrepreneurs
**Primary Language**: Arabic (RTL-first), English (LTR)
**Token Location**: `/portal-web/src/styles.css` (daisyUI theme), `/storefront-web/src/index.css`

---

## Portal Web Visual Language (Minimal)

Kyora Portal UI should feel **calm, minimal, and effortless**.

- **No gradients. No shadows.**
- Separation comes from **spacing + typography + borders**.
- Use **daisyUI semantic tokens** (Primary/Secondary/Base) instead of ad-hoc colors.

## Colors

### Semantic Palette

| Token        | Value                 | daisyUI Class       | Usage                      |
| ------------ | --------------------- | ------------------- | -------------------------- |
| Primary      | `#0D9488` (Teal)      | `bg-primary`        | Brand, main actions, links |
| Secondary    | `#EAB308` (Gold)      | `bg-secondary`      | Secondary CTAs, accents    |
| Accent       | `#FACC15` (Yellow)    | `bg-accent`         | Attention, badges          |
| Success      | `#10B981` (Green)     | `bg-success`        | Completed, paid, active    |
| Error        | `#EF4444` (Red)       | `bg-error`          | Failed, invalid, critical  |
| Warning      | `#F59E0B` (Orange)    | `bg-warning`        | Low stock, caution         |
| Info         | `oklch(60% 0.12 192)` | `bg-info`           | Informational              |
| Base-100     | `#FFFFFF`             | `bg-base-100`       | Main backgrounds, cards    |
| Base-200     | `#F8FAFC`             | `bg-base-200`       | App background             |
| Base-300     | `#E2E8F0`             | `bg-base-300`       | Borders, dividers          |
| Base-Content | `#0F172A`             | `text-base-content` | Primary text               |

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
**Fallback**: `-apple-system, BlinkMacSystemFont, sans-serif`

**Implementation**: Font family is defined in `portal-web/src/styles.css` and applied globally via the `body` tag.

| Scale   | Size | Weight   | Usage              | Tailwind Class |
| ------- | ---- | -------- | ------------------ | -------------- |
| Display | 32px | Bold     | Marketing headers  | `text-[32px]`  |
| H1      | 24px | Bold     | Page titles        | `text-2xl`     |
| H2      | 20px | SemiBold | Section headers    | `text-xl`      |
| H3      | 18px | Medium   | Card titles        | `text-lg`      |
| Body-L  | 16px | Regular  | Default text       | `text-base`    |
| Body-M  | 14px | Regular  | Secondary text     | `text-sm`      |
| Caption | 12px | Medium   | Labels, timestamps | `text-xs`      |
| Micro   | 10px | Bold     | Small badges       | `text-[10px]`  |

**Usage Notes**:

- `text-2xl` (24px) is used for page titles in ResourceListLayout
- `text-xs` (12px) is heavily used for captions, labels, and secondary text
- `text-[32px]` and `text-[10px]` are defined but not currently used in the codebase
- Most components use the semantic classes (`text-2xl`, `text-xl`, `text-lg`, `text-base`, `text-sm`, `text-xs`) rather than pixel values

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

## Shadows (No Shadows) - Updated 2026-01-19

**Rule (Explicit):**

- MUST NOT: Use any shadow utilities (`shadow-*`, `drop-shadow-*`)
- MUST: Use borders for elevation and separation
- WHY: Maintains minimal, calm aesthetic; prevents visual clutter

**Common Mistakes:**

❌ **Wrong - Using shadows for elevation:**

```tsx
// FAB buttons
className = "btn-circle shadow-lg";

// Hover states
className = "hover:shadow-sm";

// Modals
className = "shadow-2xl";
```

✅ **Correct - Using borders:**

```tsx
// FAB buttons
className = "btn-circle border-2 border-base-300";

// Hover states - use border color change
className = "border border-base-300 hover:border-primary/50";

// Modals
className = "border-2 border-base-300";
```

**Affected Areas:**

- All UI components requiring visual separation
- Interactive elements (buttons, cards, panels)
- Modal overlays and floating elements

**Related Drift Reports:**

- `backlog/drifts/2026-01-18-portal-web-shadow-usage-violates-design-tokens.md` - Resolved 2026-01-19

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
```

---

## Responsive Breakpoints

```
sm: 640px  (tablets portrait)
md: 768px  (tablets landscape, small desktops)
lg: 1024px (desktops)
xl: 1280px (large desktops)
```

**React Hook**: `useMediaQuery("(max-width: 768px)")` for mobile detection
**Implementation**: Hook is located at `/portal-web/src/hooks/useMediaQuery.ts`

Common usage patterns:

```tsx
const isMobile = useMediaQuery("(max-width: 768px)");
const isDesktop = useMediaQuery("(min-width: 768px)");
const isMobileSmall = useMediaQuery("(max-width: 640px)");
```

---

## Touch Targets

**Minimum**: 44×44px (prefer 48px+)

- Avoid hardcoded pixel heights.
- Prefer Tailwind scale (`min-h-12`, `min-w-12`) and daisyUI sizing (`btn-lg`) to keep UI consistent.
- Mobile primary CTA should typically be `btn btn-primary btn-lg w-full`.

---

## Animations

```
Pulse (skeletons):    animate-pulse
Spin (loaders):       animate-spin
Press (buttons):      active:scale-95 transition-all
Hover (interactive):  transition-colors duration-200
Standard:             transition-all duration-200
```

**Implementation Notes**:

- `animate-pulse` is used in Skeleton components, StatCardSkeleton, ChartSkeleton, FormSkeleton
- `animate-spin` is used with Loader2 icons in loading states
- `active:scale-95` is used in Button component and BottomNav for touch feedback
- Standard transition is `transition-all duration-200` (defined in main Transitions section)

---

## Icon Sizing

```
Inline with text:     size={18} or size={20}
Large interactive:    size={24}
Headers:              size={28} to size={32}
Empty states:         size={48} to size={64}
```

**Common Icon Sizes in Use**:

- 16: ChevronUp/ChevronDown in Table sorting
- 18: Plus, Trash2, GripVertical in forms and actions
- 20: Sidebar navigation, Package, Plus in headers
- 24: Default for large interactive elements
- 32: Notes empty state
- 48-64: Empty states in ResourceListLayout

**Icon Colors**:

```
Default:     text-base-content/50
Active:      text-primary
Disabled:    text-base-content/30
Semantic:    text-success, text-error, text-warning
```

---

## Chart.js Integration

Design tokens are integrated with Chart.js visualizations through dedicated theme configuration.

**Theme File**: `/portal-web/src/lib/charts/chartTheme.ts`
**Plugin File**: `/portal-web/src/lib/charts/chartPlugins.ts`

All Chart.js instances use IBM Plex Sans Arabic font family and pull colors from the daisyUI theme. See `.github/instructions/charts.instructions.md` for full Chart.js patterns.

---

## Implementation Notes

### Actual CSS File Locations

- **Portal Web**: `/portal-web/src/styles.css` (uses Tailwind v4 + daisyUI plugin)
- **Storefront Web**: `/storefront-web/src/index.css` (uses Tailwind v4 + daisyUI plugin)

### Font Loading

Portal web loads IBM Plex Sans Arabic from Google Fonts (implicit via browser).
Storefront web explicitly imports from Google Fonts via `@import` in CSS.

### daisyUI Theme Configuration

Both portal and storefront define custom `kyora` theme using daisyUI v5 plugin syntax:

```css
@plugin "daisyui/theme" {
  name: "kyora";
  default: true;
  /* ... color definitions using oklch ... */
}
```

All colors are defined using OKLCH color space for better perceptual uniformity.

### Pattern Usage

The codebase primarily uses **daisyUI semantic classes** (`bg-primary`, `text-base-content`, etc.) rather than direct hex values or arbitrary values. This ensures consistency and makes theme changes easier.

---

## Related Documentation

- **UI Implementation**: `.github/instructions/ui-implementation.instructions.md`
- **Charts**: `.github/instructions/charts.instructions.md`
- **Portal Architecture**: `.github/instructions/portal-web-architecture.instructions.md`
