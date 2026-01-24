---
description: Kyora Design System — Colors, typography, spacing, components, RTL patterns, accessibility
applyTo: "portal-web/**,storefront-web/**,mobile-web/**"
---

# Kyora Design System (SSOT)

**Purpose**: Single source of truth for all UI implementation across Kyora frontend applications.

**Philosophy**: "Professional tools that feel effortless" — clean, minimal, calm. No shadows, no gradients.

**Primary audience**: Middle East entrepreneurs (Arabic-first, mobile-first, low-moderate tech literacy).

---

## 1. Design Principles

### Core Principles

| Principle               | What it means in practice                                        |
| ----------------------- | ---------------------------------------------------------------- |
| Mobile-first            | Design for 375px viewport first, enhance for larger screens      |
| Arabic-first            | RTL is default (not afterthought); Arabic phrasing drives design |
| Plain language          | No jargon ("Profit" not "EBITDA", "Cash in hand" not "Accrual")  |
| Clarity over complexity | Simple flows, obvious next steps, no hidden features             |
| Minimal aesthetic       | Borders + spacing for separation (no shadows, no gradients)      |
| Confidence-building     | Show progress, confirm actions, explain impact                   |

### Visual Language

- **Calm and discreet**: Subdued colors, generous spacing
- **Purposeful hierarchy**: Typography and spacing create structure
- **Touch-optimized**: 44×44px minimum targets, swipe-friendly
- **Scannable content**: Icons + labels, short sentences, visual grouping

**Wrong aesthetics to avoid**:

- ❌ Dense dashboards (too much info at once)
- ❌ Flashy gradients or heavy shadows (feels consumer-app, not professional)
- ❌ Overly technical UI (charts without context, jargon labels)

---

## 2. Color System

### Semantic Palette (daisyUI Theme)

| Token        | Value (OKLCH)         | Hex Equivalent | daisyUI Class       | Usage                       |
| ------------ | --------------------- | -------------- | ------------------- | --------------------------- |
| Primary      | `oklch(55% 0.12 188)` | `#0D9488`      | `bg-primary`        | Brand actions, links, focus |
| Secondary    | `oklch(75% 0.15 90)`  | `#EAB308`      | `bg-secondary`      | Secondary CTAs, accents     |
| Accent       | `oklch(85% 0.15 95)`  | `#FACC15`      | `bg-accent`         | Attention, badges           |
| Success      | `oklch(65% 0.15 150)` | `#10B981`      | `bg-success`        | Completed, paid, active     |
| Error        | `oklch(60% 0.22 25)`  | `#EF4444`      | `bg-error`          | Failed, invalid, critical   |
| Warning      | `oklch(70% 0.15 70)`  | `#F59E0B`      | `bg-warning`        | Low stock, caution          |
| Info         | `oklch(60% 0.12 192)` | `#3B82F6`      | `bg-info`           | Informational               |
| Base-100     | `oklch(100% 0 0)`     | `#FFFFFF`      | `bg-base-100`       | Main backgrounds, cards     |
| Base-200     | `oklch(98% 0 0)`      | `#F8FAFC`      | `bg-base-200`       | App background              |
| Base-300     | `oklch(92% 0 0)`      | `#E2E8F0`      | `bg-base-300`       | Borders, dividers           |
| Base-Content | `oklch(15% 0 0)`      | `#0F172A`      | `text-base-content` | Primary text                |

### Text Opacity Scale

Use semantic text colors with opacity modifiers:

```
text-base-content       (100% — primary text, headings)
text-base-content/70    (70% — secondary text, descriptions)
text-base-content/60    (60% — tertiary text, captions)
text-base-content/50    (50% — icons, decorative elements)
text-base-content/40    (40% — disabled states, placeholders)
text-base-content/30    (30% — subtle decorative, backgrounds)
```

### Color Usage Rules

**DO:**

- Use daisyUI semantic classes: `bg-primary`, `text-success`, `border-base-300`
- Use opacity modifiers for text hierarchy: `text-base-content/70`
- Use semantic colors for status: `badge-success` (paid), `badge-warning` (low stock)

**DON'T:**

- ❌ Hardcode hex colors: `bg-[#0d9488]` (use `bg-primary`)
- ❌ Use arbitrary colors: `text-blue-500` (not in design system)
- ❌ Override daisyUI colors with Tailwind utilities (breaks theming)

---

## 3. Typography

### Font Family

**Primary**: IBM Plex Sans Arabic  
**Fallback**: `-apple-system, BlinkMacSystemFont, sans-serif`

**Implementation**: Loaded via Google Fonts, applied globally in `styles.css`.

### Type Scale

| Scale   | Size | Weight   | Line Height | Tailwind Class | Usage                    |
| ------- | ---- | -------- | ----------- | -------------- | ------------------------ |
| Display | 32px | Bold     | 1.2         | `text-[32px]`  | Marketing headers (rare) |
| H1      | 24px | Bold     | 1.3         | `text-2xl`     | Page titles              |
| H2      | 20px | SemiBold | 1.3         | `text-xl`      | Section headers          |
| H3      | 18px | Medium   | 1.4         | `text-lg`      | Card titles, subsections |
| Body-L  | 16px | Regular  | 1.5         | `text-base`    | Default body text        |
| Body-M  | 14px | Regular  | 1.5         | `text-sm`      | Secondary text, labels   |
| Caption | 12px | Medium   | 1.5         | `text-xs`      | Timestamps, small labels |
| Micro   | 10px | Bold     | 1.6         | `text-[10px]`  | Tiny badges (rare)       |

### Usage Patterns

```tsx
// Page title
<h1 className="text-2xl font-bold">Orders</h1>

// Section header
<h2 className="text-xl font-semibold">Recent Activity</h2>

// Card title
<h3 className="text-lg font-medium">Product Details</h3>

// Body text
<p className="text-base">This is the main content.</p>

// Secondary text
<p className="text-sm text-base-content/70">Last updated 2 hours ago</p>

// Caption
<span className="text-xs text-base-content/60">Order #12345</span>
```

### Typography Rules

**DO:**

- Use semantic classes (`text-2xl`, `text-base`, `text-xs`)
- Pair font size with appropriate weight (H1 is bold, body is regular)
- Use opacity for hierarchy (`text-base-content/70` for secondary text)

**DON'T:**

- ❌ Use arbitrary sizes: `text-[15px]` (not in scale)
- ❌ Mix too many weights (stick to regular, medium, semibold, bold)
- ❌ Use small text (<12px) for critical actions

---

## 4. Spacing System

### 4px Grid

All spacing uses multiples of 4px:

```
gap-0  → 0px      gap-1  → 4px      gap-2  → 8px
gap-3  → 12px     gap-4  → 16px     gap-5  → 20px
gap-6  → 24px     gap-8  → 32px     gap-12 → 48px
gap-16 → 64px
```

### Common Spacing Patterns

| Pattern           | Spacing     | Example                        |
| ----------------- | ----------- | ------------------------------ |
| Container padding | `px-4 py-4` | 16px horizontal/vertical       |
| Card body         | `card-body` | daisyUI default (1rem padding) |
| List item gap     | `gap-3`     | 12px between cards             |
| Form field gap    | `gap-4`     | 16px between inputs            |
| Button icon gap   | `gap-2`     | 8px between icon + text        |
| Section spacing   | `space-y-6` | 24px vertical rhythm           |

### Mobile Safe Area

- **Horizontal**: `px-4` (16px) — prevents edge cutoff
- **Vertical**: `py-4` (16px) — breathing room at top/bottom
- **Bottom nav**: Account for 60px + safe area on mobile

---

## 5. Layout Patterns

### Responsive Breakpoints

```
sm:  640px  (Small tablets portrait)
md:  768px  (Tablets landscape, small desktops)
lg:  1024px (Desktops)
xl:  1280px (Large desktops)
2xl: 1536px (Extra large)
```

**Mobile detection**: Use `useMediaQuery("(max-width: 768px)")` hook.

### Container Patterns

```tsx
// Page container (mobile-first)
<div className="px-4 py-4 max-w-7xl mx-auto">
  {/* Content */}
</div>

// Card container
<div className="card bg-base-100 border border-base-300">
  <div className="card-body">{/* Content */}</div>
</div>

// Grid layout (responsive)
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {/* Cards */}
</div>

// Flex layout (mobile stack, desktop row)
<div className="flex flex-col md:flex-row gap-4">
  {/* Items */}
</div>
```

---

## 6. Component Patterns

### Buttons

**Variants**:

```tsx
<button className="btn btn-primary">Primary Action</button>
<button className="btn btn-secondary">Secondary Action</button>
<button className="btn btn-ghost">Tertiary Action</button>
<button className="btn btn-outline">Outline Action</button>
```

**Sizes**:

```tsx
<button className="btn btn-xs">Extra Small</button>
<button className="btn btn-sm">Small</button>
<button className="btn">Default</button>
<button className="btn btn-lg">Large</button>
```

**Mobile primary CTA**:

```tsx
<button className="btn btn-primary btn-lg w-full">Submit</button>
```

### Cards

**Standard card**:

```tsx
<div className="card bg-base-100 border border-base-300">
  <div className="card-body">
    <h3 className="card-title">Card Title</h3>
    <p>Card content goes here.</p>
  </div>
</div>
```

**Clickable card**:

```tsx
<div className="card bg-base-100 border border-base-300 cursor-pointer hover:border-primary/50 transition-colors">
  <div className="card-body">{/* Content */}</div>
</div>
```

### Badges

```tsx
<span className="badge badge-primary">New</span>
<span className="badge badge-success">Paid</span>
<span className="badge badge-warning">Low Stock</span>
<span className="badge badge-error">Failed</span>
```

### Inputs (via Form System)

```tsx
// ✅ Use form system (see forms.instructions.md)
<form.TextField name="email" label="Email" required />

// ❌ Don't use raw daisyUI inputs
<input className="input input-bordered" /> {/* WRONG */}
```

### Loading States

```tsx
// Spinner
<span className="loading loading-spinner loading-lg"></span>

// Skeleton (prefer for page-level loading)
<div className="skeleton h-32 w-full rounded-box"></div>

// Progress bar
<progress className="progress progress-primary" value="70" max="100"></progress>
```

---

## 7. RTL-First Rules (NON-NEGOTIABLE)

### Forbidden Classes

**NEVER USE**: `left`, `right`, `ml-*`, `mr-*`, `pl-*`, `pr-*`, `float-left`, `float-right`, `text-left`, `text-right`, `border-l`, `border-r`

### Required Replacements

| Forbidden    | Use Instead  | Purpose              |
| ------------ | ------------ | -------------------- |
| `ml-4`       | `ms-4`       | margin-inline-start  |
| `mr-4`       | `me-4`       | margin-inline-end    |
| `pl-4`       | `ps-4`       | padding-inline-start |
| `pr-4`       | `pe-4`       | padding-inline-end   |
| `left-0`     | `start-0`    | positioning          |
| `right-0`    | `end-0`      | positioning          |
| `text-left`  | `text-start` | text alignment       |
| `text-right` | `text-end`   | text alignment       |
| `border-l`   | `border-s`   | border-inline-start  |
| `border-r`   | `border-e`   | border-inline-end    |

### Directional Icons

**Must rotate 180° in RTL**:

```tsx
import { useLanguage } from '@/hooks/useLanguage'

const { isRTL } = useLanguage()

// ✅ CORRECT
<ArrowLeft className={isRTL ? 'rotate-180' : ''} />
<ArrowRight className={isRTL ? 'rotate-180' : ''} />
```

**Auto-flip (no rotation needed)**: `ChevronLeft`, `ChevronRight`, `ChevronsLeft`, `ChevronsRight`

**Never rotate**: `ChevronDown`, `ChevronUp`, `Plus`, `X`, `Check`, `Trash2`, `Settings`, `Menu`, `Search`, `Filter`

### Mixed-Direction Content

```tsx
// Phone numbers, order IDs, IBANs always LTR
<span dir="ltr">{phoneNumber}</span>
<span dir="ltr">{orderId}</span>

// Currency inputs (price, quantity)
<input dir="ltr" type="number" />
```

---

## 8. Border Radius

```
rounded-sm   → 4px  (tags, small elements)
rounded-md   → 8px  (default inner elements)
rounded-lg   → 12px (cards, modals, containers)
rounded-xl   → 16px (buttons, bottom sheets)
rounded-box  → 12px (daisyUI card utility)
rounded-full → pills, avatars, circular buttons
```

**Usage**:

```tsx
<div className="card rounded-lg">Card</div>
<button className="btn rounded-xl">Button</button>
<span className="badge rounded-full">Badge</span>
```

---

## 9. Borders & Shadows

### Borders

**Default**: `border border-base-300`  
**Focus**: `focus:border-primary`  
**Error**: `border-error`  
**Width**: 1px (default), 2px for emphasis (`border-2`)

### Shadows (FORBIDDEN)

**Rule**: MUST NOT use shadows (`shadow-*`, `drop-shadow-*`)

**Why**: Maintains minimal, calm aesthetic; prevents visual clutter.

**Wrong**:

```tsx
// ❌ Don't use shadows
<div className="card shadow-lg">Card</div>
<button className="btn hover:shadow-sm">Button</button>
```

**Correct**:

```tsx
// ✅ Use borders for elevation
<div className="card border-2 border-base-300">Card</div>
<button className="btn border border-base-300 hover:border-primary/50">Button</button>
```

---

## 10. Transitions & Animations

### Standard Transitions

```
Standard: transition-all duration-200
Colors:   transition-colors duration-200
Opacity:  transition-opacity duration-200
```

### Animations

```tsx
// Pulse (skeleton loading)
<div className="animate-pulse skeleton h-32" />

// Spin (loading icons)
<Loader2 className="animate-spin" />

// Press feedback (buttons)
<button className="active:scale-95 transition-all">Press Me</button>

// Hover (interactive elements)
<div className="hover:border-primary/50 transition-colors">Card</div>
```

---

## 11. Icons (Lucide React)

### Icon Sizing

```tsx
<Plus size={16} />  // Small (inline with text, table icons)
<Plus size={18} />  // Default small (form actions, drag handles)
<Plus size={20} />  // Default (buttons, inputs, sidebar nav)
<Plus size={24} />  // Large (page headers, interactive elements)
<Plus size={32} />  // Extra large (empty states)
<Plus size={48} />  // Huge (marketing, large empty states)
```

### Icon Colors

```tsx
// Default (decorative)
<Plus className="text-base-content/50" />

// Active/interactive
<Plus className="text-primary" />

// Disabled
<Plus className="text-base-content/30" />

// Semantic
<CheckCircle className="text-success" />
<AlertCircle className="text-error" />
<AlertTriangle className="text-warning" />
```

### Icon Accessibility

```tsx
// Icon with text (decorative)
<button className="btn gap-2">
  <Plus size={20} aria-hidden="true" />
  <span>Add</span>
</button>

// Icon-only button (MUST have aria-label)
<button className="btn btn-ghost" aria-label="Delete">
  <Trash2 size={20} />
</button>

// Loading icon
<Loader2 size={20} className="animate-spin" aria-label="Loading" />
```

---

## 12. Touch Targets & Accessibility

### Touch Target Minimums

**Minimum**: 44×44px (prefer 48px+)

**Implementation**:

```tsx
// ✅ Correct — Use Tailwind scale
<button className="btn min-h-12 min-w-12">Icon</button>

// ❌ Wrong — Hardcoded pixels
<button className="h-[40px] w-[40px]">Icon</button> {/* Too small */}
```

### ARIA Attributes

| Element Type     | Required Attributes                                     |
| ---------------- | ------------------------------------------------------- |
| Icon-only button | `aria-label="Action description"`                       |
| Form input       | `aria-invalid`, `aria-describedby`, `aria-required`     |
| Modal/Dialog     | `role="dialog"`, `aria-modal="true"`, `aria-labelledby` |
| Decorative icon  | `aria-hidden="true"`                                    |
| Error message    | `role="alert"`                                          |
| Live region      | `aria-live="polite"` or `aria-live="assertive"`         |

### Focus Management

**Standard focus ring**:

```tsx
className =
  "focus:ring-2 focus:ring-primary/20 focus:border-primary focus:outline-none";
```

**Keyboard navigation**: All modals/overlays must trap focus and support Tab, Escape, Enter/Space.

---

## 13. Responsive Patterns

### Mobile vs Desktop

| Component    | Mobile (<768px)          | Desktop (≥768px)  |
| ------------ | ------------------------ | ----------------- |
| Modal        | BottomSheet (85% height) | Dialog (centered) |
| Select       | Action sheet             | Dropdown          |
| Date picker  | Full-screen              | Popover           |
| Navigation   | Bottom bar               | Sidebar           |
| Form buttons | Full width (`w-full`)    | Auto width        |
| Page layout  | Single column            | Multi-column      |

### Example

```tsx
// Mobile: full-width button, Desktop: auto-width
<button className="btn btn-primary w-full md:w-auto">Submit</button>;

// Mobile: BottomSheet, Desktop: Dialog
{
  isMobile ? (
    <BottomSheet isOpen={isOpen} onClose={onClose}>
      <Form />
    </BottomSheet>
  ) : (
    <Dialog isOpen={isOpen} onClose={onClose}>
      <Form />
    </Dialog>
  );
}
```

---

## 14. Implementation Files

### CSS Files

- **Portal Web**: `/portal-web/src/styles.css` (Tailwind v4 + daisyUI plugin)
- **Storefront Web**: `/storefront-web/src/index.css` (Tailwind v4 + daisyUI plugin)

### daisyUI Theme Configuration

Both apps define custom `kyora` theme using daisyUI v5 plugin syntax:

```css
@plugin "daisyui/theme" {
  name: "kyora";
  default: true;
  /* ... OKLCH color definitions ... */
}
```

### Font Loading

- **Portal Web**: Implicit via browser (Google Fonts fallback)
- **Storefront Web**: Explicit `@import` from Google Fonts

---

## 15. Chart.js Integration

Design tokens integrate with Chart.js visualizations via dedicated theme configuration.

**Theme File**: `/portal-web/src/lib/charts/chartTheme.ts`  
**Plugin File**: `/portal-web/src/lib/charts/chartPlugins.ts`

All Chart.js instances use IBM Plex Sans Arabic font and pull colors from daisyUI theme.

**See**: `.github/instructions/charts.instructions.md` for full Chart.js patterns.

---

## Agent Validation Checklist

Before completing any UI implementation:

- ☑ **RTL**: No `left`/`right` classes (use `start`/`end`, `ms-*`, `me-*`, `ps-*`, `pe-*`)
- ☑ **RTL**: Directional arrows rotate 180° with `isRTL` check
- ☑ **Colors**: Using daisyUI semantic classes (`bg-primary`, `text-base-content/70`)
- ☑ **Typography**: Using type scale classes (`text-2xl`, `text-base`, `text-xs`)
- ☑ **Spacing**: Using 4px grid (`gap-4`, `px-4`, `space-y-6`)
- ☑ **Borders**: Using borders for separation (no shadows)
- ☑ **Icons**: From Lucide React only, correct size, RTL handled
- ☑ **Accessibility**: Focus visible, aria-labels on icon-only buttons, touch targets ≥44px
- ☑ **Mobile-first**: Base styles mobile, `md:` for desktop enhancements
- ☑ **Components**: Using daisyUI base classes (`.btn`, `.card`, `.badge`)

---

## Related Documentation

- **UX Strategy**: `.github/instructions/kyora/ux-strategy.instructions.md` — Interaction patterns, content strategy
- **Target Customer**: `.github/instructions/kyora/target-customer.instructions.md` — Who we're designing for
- **Brand Key**: `.github/instructions/kyora/brand-key.instructions.md` — Voice and tone guidelines
- **Forms**: `.github/instructions/forms.instructions.md` — Form-specific patterns
- **Charts**: `.github/instructions/charts.instructions.md` — Data visualization
- **i18n**: `.github/instructions/frontend/_general/i18n.instructions.md` — Translation rules
