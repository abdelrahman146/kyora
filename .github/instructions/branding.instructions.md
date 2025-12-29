---
description: KYORA DESIGN SYSTEM (KDS) - Portal Web Implementation
applyTo: "portal-web/**,storefront-web/**"
---

# KYORA DESIGN SYSTEM (KDS) - Portal Web

**Target Platform:** Web (Desktop & Mobile), PWA
**Primary Audience:** Middle East (Arabic First / RTL)  
**Core Philosophy:** "Professional tools that feel effortless."

---

## 1. Brand Identity & Philosophy

### 1.1. Core Values

1. **Tamkeen (Empowerment):** Transform social media entrepreneurs into professional business owners
2. **Basata (Simplicity):** Every interaction must be intuitive and purposeful
3. **Thiqah (Trust):** Professional, stable, and secure design that handles money and inventory

### 1.2. Voice & Tone

- **Language:** Modern Standard Arabic (MSA) - approachable and business-standard
- **Tone:** Direct, warm, respectful, and professional
- **English:** Full support with proper LTR layout
- **RTL-First:** All layouts designed for RTL (Arabic) first, LTR (English) is the mirror

---

## 2. Design Tokens

All design tokens are defined in `/portal-web/src/index.css` using daisyUI theme system.

### 2.1. Color Palette

**Primary (Teal - The Brand)**
- Primary: `#0D9488` (oklch(52% 0.11 192))
- Usage: Main actions, links, brand elements
- DaisyUI: `bg-primary`, `text-primary`, `btn-primary`

**Secondary (Gold - Action/Accent)**
- Secondary: `#EAB308` (oklch(75% 0.15 90))
- Usage: Secondary CTAs, highlighting important items
- DaisyUI: `bg-secondary`, `text-secondary`, `btn-secondary`

**Accent (Yellow - CTAs)**
- Accent: `#FACC15` (oklch(87% 0.14 100))
- Usage: Attention-grabbing CTAs, badges
- DaisyUI: `bg-accent`, `text-accent`

**Base Colors (Backgrounds & Surfaces)**
- `base-100`: `#FFFFFF` (Main backgrounds, cards)
- `base-200`: `#F8FAFC` (App background, secondary surfaces)
- `base-300`: `#E2E8F0` (Borders, dividers)
- `base-content`: `#0F172A` (Primary text)

**Semantic Colors**
- Success: `#10B981` (Completed, paid, active)
- Error: `#EF4444` (Failed, invalid, critical)
- Warning: `#F59E0B` (Low stock, ending soon)
- Info: `oklch(60% 0.12 192)` (Informational states)

**Text Colors**
- Primary text: `text-base-content` (`#0F172A`)
- Secondary text: `text-base-content/70` (70% opacity)
- Tertiary text: `text-base-content/60` (60% opacity)
- Disabled text: `text-base-content/40` (40% opacity)
- Icon color: `text-base-content/50` (50% opacity)
- Placeholder: `placeholder:text-base-content/40`

### 2.2. Typography

**Font Family:** `IBM Plex Sans Arabic` (via Google Fonts)
**Fallback:** `Almarai`, `-apple-system`, `sans-serif`

**Scale (Fluid Responsive)**
- Display: 32px / Bold (Page marketing headers)
- H1: 24px / Bold (Page titles)
- H2: 20px / SemiBold (Section headers)
- H3: 18px / Medium (Card titles)
- Body-L: 16px / Regular (Default text - `text-base`)
- Body-M: 14px / Regular (Secondary text - `text-sm`)
- Caption: 12px / Medium (Labels, timestamps - `text-xs`)
- Micro: 10px / Bold (Small badges)

**Line Height:**
- Headings: 1.2-1.3
- Body: 1.5-1.6

### 2.3. Spacing (4px Baseline Grid)

Use Tailwind spacing scale:
- **Gap-1:** 4px (`gap-1`, `space-y-1`)
- **Gap-2:** 8px (`gap-2`, `space-y-2`)
- **Gap-3:** 12px (`gap-3`, `space-y-3`)
- **Gap-4:** 16px (`gap-4`, `space-y-4`) - **Standard spacing**
- **Gap-6:** 24px (`gap-6`, `space-y-6`)
- **Gap-8:** 32px (`gap-8`, `space-y-8`)

**Safe Area Padding:** Always use `p-4` (16px) for mobile content padding

### 2.4. Border Radius

- `rounded-sm`: 4px (Small elements, tags)
- `rounded-md`: 8px (Default, inner elements)
- `rounded-lg`: 12px (Cards, modals, containers)
- `rounded-xl`: 16px (Buttons, bottom sheets, prominent cards)
- `rounded-box`: 12px (daisyUI card utility)
- `rounded-full`: Pills, avatars, circular buttons

### 2.5. Shadows

- `shadow-sm`: Subtle card elevation (`0 1px 2px 0 rgb(0 0 0 / 0.05)`)
- `shadow`: Standard card shadow
- `shadow-md`: Elevated elements
- `shadow-lg`: Floating elements, dropdowns
- `shadow-xl`: Modals, important overlays

### 2.6. Borders

- Default border: `border border-base-300`
- Focus border: `focus:border-primary`
- Error border: `border-error`
- Border width: 1px (default)

### 2.7. Transitions

All interactive elements use smooth transitions:
- Standard: `transition-all duration-200`
- Colors: `transition-colors duration-200`
- Opacity: `transition-opacity duration-200`

---

## 3. Component Standards

### 3.1. Buttons

**Implementation:** `/portal-web/src/components/atoms/Button.tsx`

**Variants:**
- `primary`: Primary brand color, main actions (`btn-primary`)
- `secondary`: Secondary style with primary-50 background
- `ghost`: Transparent background (`btn-ghost`)
- `outline`: Outline style (`btn-outline`)

**Sizes:**
- `sm`: 40px height (`btn-sm`)
- `md`: 52px height (default, min-height for touch targets)
- `lg`: 56px height (`btn-lg`)

**States:**
- Active: `active:scale-95` (press animation)
- Loading: Spinner with disabled state
- Disabled: 50% opacity, cursor-not-allowed

**Key Classes:**
```tsx
"btn rounded-xl font-semibold transition-all active:scale-95 
 disabled:opacity-50 disabled:cursor-not-allowed"
```

### 3.2. Form Inputs

**Base Implementation:** `/portal-web/src/components/atoms/Input.tsx`, `FormInput.tsx`, `PasswordInput.tsx`

**Design Standards:**
- Min-height: `50px` (`h-[50px]`) - Critical for touch targets
- Border: `border-base-300`
- Background: `bg-base-100`
- Text: `text-base text-base-content`
- Placeholder: `placeholder:text-base-content/40`
- Focus: `focus:border-primary focus:ring-2 focus:ring-primary/20`
- Error: `border-error focus:border-error focus:ring-error/20`
- Icons: `text-base-content/50` with `aria-hidden="true"`
- Transitions: `transition-all duration-200`

**Key Classes:**
```tsx
"input input-bordered w-full h-[50px] bg-base-100 text-base-content
 focus:border-primary focus:ring-2 focus:ring-primary/20
 transition-all duration-200"
```

**Icon Positioning:**
- Start icon: `ps-10` (padding-inline-start)
- End icon: `pe-10` (padding-inline-end)
- Icon wrapper: `z-10`, input: `z-0` (proper layering)

**Accessibility:**
- All inputs have labels or `aria-label`
- Error messages linked via `aria-describedby`
- Required fields marked with `aria-required`
- Icons wrapped in spans with `aria-hidden="true"`

### 3.3. Textarea

**Implementation:** `/portal-web/src/components/atoms/FormTextarea.tsx`

**Standards:**
- Min-height: `min-h-[120px]` (default)
- Resize: `resize-y` (vertical only)
- Character counter: Optional with live updates
- Same color scheme as inputs

### 3.4. Search Input

**Implementation:** `/portal-web/src/components/molecules/SearchInput.tsx`

**Features:**
- Debounced input (300ms default)
- Search icon at start position (RTL-aware)
- Clear button with loading state indicator
- Same styling as form inputs for consistency

### 3.5. Cards

**DaisyUI Classes:** `card`, `card-body`

**Standard Card Pattern:**
```tsx
<div className="card bg-base-100 border border-base-300 shadow-sm">
  <div className="card-body">
    {/* Content */}
  </div>
</div>
```

**Variants:**
- Default: White background, subtle border and shadow
- Hover: `hover:shadow-md transition-shadow`
- Clickable: `cursor-pointer` with hover effect

### 3.6. Badges

**Implementation:** `/portal-web/src/components/atoms/Badge.tsx`

**Variants:**
- `default`: Neutral gray
- `primary`: Brand teal
- `secondary`: Gold/yellow
- `success`: Green
- `error`: Red
- `warning`: Orange
- `info`: Blue

**DaisyUI Classes:** `badge`, `badge-primary`, `badge-success`, etc.

### 3.7. Modals & Dialogs

**Implementation:** 
- `/portal-web/src/components/atoms/Modal.tsx` - Basic modal
- `/portal-web/src/components/atoms/Dialog.tsx` - Dialog with footer actions

**Behavior:**
- Portal-based rendering
- Backdrop blur: `backdrop-blur-sm`
- Click outside to close (optional)
- Escape key to close
- Focus trap for accessibility
- Sizes: `sm`, `md`, `lg`, `xl`, `full`

### 3.8. Bottom Sheets

**Implementation:** `/portal-web/src/components/molecules/BottomSheet.tsx`

**Behavior:**
- Mobile: Slides up from bottom (85% max height)
- Desktop: Side drawer from start/end
- Drag handle on mobile (gray pill)
- Scrollable content area
- Sticky footer for actions
- RTL-aware positioning

### 3.9. Skeletons

**Implementation:** `/portal-web/src/components/atoms/Skeleton.tsx`

**Standards:**
- Color: `bg-base-300`
- Animation: `animate-pulse`
- Variants: `text`, `circular`, `rectangular`
- Match content shape (circle for avatars, rectangles for text)

**Usage:**
```tsx
<Skeleton variant="circular" width={40} height={40} />
<Skeleton variant="rectangular" height={100} />
<Skeleton variant="text" width="60%" />
```

### 3.10. Filter Button

**Implementation:** `/portal-web/src/components/organisms/FilterButton.tsx`

**Features:**
- Unified component with trigger button + drawer
- Active filter count badge
- Consistent styling with form fields
- Internal state management
- Apply and Reset callbacks

**Design Consistency:**
- Same min-height as form inputs (`50px`)
- Same focus states and transitions
- Badge indicator for active filters

---

## 4. RTL-First Design Rules

### 4.1. Logical Properties (MANDATORY)

**NEVER use:**
- `left`, `right`
- `margin-left`, `margin-right`, `padding-left`, `padding-right`
- `float-left`, `float-right`
- `text-align: left`, `text-align: right`
- `border-left`, `border-right`

**ALWAYS use:**
- `start`, `end` (Tailwind utilities)
- `ms-*`, `me-*` (margin-inline-start/end)
- `ps-*`, `pe-*` (padding-inline-start/end)
- `start-*`, `end-*` (positioning)
- `text-start`, `text-end`
- `border-s-*`, `border-e-*`

**Examples:**
```tsx
// ❌ WRONG
<div className="ml-4 text-left float-left" />

// ✅ CORRECT
<div className="ms-4 text-start float-start" />
```

### 4.2. Directional Icons

Icons like arrows, chevrons, and back buttons must be mirrored in RTL:

```tsx
import { useLanguage } from "@/hooks/useLanguage";

const { isRTL } = useLanguage();

<ArrowLeft className={isRTL ? "rotate-180" : ""} />
```

### 4.3. Phone Numbers

Always display LTR (left-to-right) even in Arabic text:

```tsx
<span dir="ltr">{phoneNumber}</span>
```

### 4.4. Document Direction

Automatically set via `useLanguage` hook:
- HTML `dir` attribute: `rtl` or `ltr`
- HTML `lang` attribute: `ar` or `en`

---

## 5. Mobile-First & Touch Standards

### 5.1. Touch Targets

**Minimum Size:** 44x44px (Apple HIG) or 50x50px (Material/KDS)
- All interactive elements: `min-h-[50px]`
- Buttons: Default `h-[52px]`
- Form inputs: `h-[50px]`
- Icon buttons: `min-w-[44px] min-h-[44px]`

### 5.2. Safe Areas

**Mobile Padding:**
- Horizontal: `px-4` (16px)
- Vertical: `py-4` (16px)
- iOS safe areas: Use `safe-bottom`, `safe-top` classes where needed

### 5.3. Responsive Breakpoints

Use Tailwind default breakpoints:
- `sm`: 640px (tablets portrait)
- `md`: 768px (tablets landscape, small desktops)
- `lg`: 1024px (desktops)
- `xl`: 1280px (large desktops)

**Hook:** Use `useMediaQuery` for JavaScript-based responsive logic:
```tsx
import { useMediaQuery } from "@/hooks/useMediaQuery";

const isMobile = useMediaQuery("(max-width: 768px)");
```

---

## 6. Interaction Patterns

### 6.1. Navigation

**Mobile:**
- Bottom navigation bar for primary actions
- Hamburger menu for secondary navigation
- Top header: Business switcher, search, notifications

**Desktop:**
- Left sidebar with collapsible navigation
- Top header: Business switcher, search, user menu

### 6.2. Forms

**Bottom Sheets (Mobile) / Modals (Desktop):**
- Create/Edit actions: Use `BottomSheet` component
- Quick filters: Use `FilterButton` component
- Confirmations: Use `Dialog` component

**Validation:**
- Inline, real-time validation
- Error messages below fields
- Field-level error highlighting
- Translated error messages

### 6.3. Loading States

**Skeletons (Preferred):**
- Match content shape
- Use for initial page loads
- Prevent layout shift

**Spinners:**
- Button loading states
- Inline actions
- Small content areas

**Progress Indicators:**
- Multi-step flows (onboarding)
- File uploads
- Long-running operations

### 6.4. Toast Notifications

**Implementation:** `react-hot-toast`

**Position:**
- Desktop: `top-right` (RTL), `top-left` (LTR)
- Mobile: `bottom-center`

**Types:**
- `toast.success()` - Green checkmark
- `toast.error()` - Red X
- `toast.loading()` - Spinner
- Custom with icons

**Auto-dismiss:** 4 seconds (default)

**Usage:**
```tsx
import toast from "react-hot-toast";
import { useLanguage } from "@/hooks/useLanguage";

const { isRTL } = useLanguage();

// Position based on language
<Toaster position={isRTL ? "top-right" : "top-left"} />

// Show toast
toast.success("Order created successfully");
toast.error("Failed to save changes");
```

---

## 7. Accessibility Standards

### 7.1. ARIA Attributes

**Required:**
- All icon-only buttons: `aria-label`
- Form inputs: `aria-invalid`, `aria-describedby`, `aria-required`
- Modals: `role="dialog"`, `aria-modal="true"`, `aria-labelledby`
- Icons: `aria-hidden="true"` (decorative)
- Live regions: `role="alert"` for errors

### 7.2. Keyboard Navigation

**Standards:**
- All interactive elements: Focusable via Tab
- Focus visible: Clear focus ring (`focus:ring-2 focus:ring-primary/20`)
- Escape: Close modals, dropdowns
- Enter/Space: Activate buttons, toggles
- Arrow keys: Navigate dropdowns, radio groups

### 7.3. Color Contrast

**WCAG AA Compliance:**
- Text on base-100: `text-base-content` (high contrast)
- Secondary text: Minimum 70% opacity
- Primary brand on white: AA compliant
- Never use low-contrast text for critical content

### 7.4. Focus Management

**Focus Ring:**
- Default: `focus:ring-2 focus:ring-primary/20 focus:border-primary focus:outline-none`
- Never remove without replacement
- Consistent across all interactive elements

**Focus Trap:**
- Modals and bottom sheets trap focus
- Tab cycles within modal
- Escape returns focus to trigger

---

## 8. Animation Guidelines

### 8.1. Interactive Animations

**Button Press:**
```tsx
active:scale-95
```

**Hover Effects:**
```tsx
hover:shadow-md transition-shadow
hover:bg-base-200 transition-colors
```

### 8.2. Transitions

**Standard:**
```tsx
transition-all duration-200
transition-colors duration-200
transition-opacity duration-200
transition-shadow duration-200
```

### 8.3. Loading Animations

**Pulse (Skeletons):**
```tsx
animate-pulse
```

**Spin (Loaders):**
```tsx
animate-spin
```

### 8.4. Page Transitions

- Subtle fade-in for new content
- Avoid jarring animations
- Keep under 300ms for perceived performance

---

## 9. Implementation Checklist

When creating new components, ensure:

### ✅ Design Consistency
- [ ] Uses daisyUI semantic tokens (`bg-primary`, `text-base-content`)
- [ ] Follows color scheme (primary, secondary, base colors)
- [ ] Uses consistent border radius (`rounded-lg`, `rounded-xl`)
- [ ] Applies standard shadows (`shadow-sm`, `shadow-md`)
- [ ] Uses 4px spacing grid (`gap-4`, `space-y-4`)

### ✅ Typography
- [ ] Uses IBM Plex Sans Arabic font family
- [ ] Correct font sizes (`text-base`, `text-lg`)
- [ ] Proper font weights (Regular, Medium, SemiBold, Bold)
- [ ] Appropriate line heights (1.5 for body text)

### ✅ RTL-First
- [ ] Uses logical properties (`start`, `end`, `ms-*`, `ps-*`)
- [ ] No `left`, `right`, `margin-left`, etc.
- [ ] Directional icons rotate in RTL (180deg)
- [ ] Phone numbers use `dir="ltr"`
- [ ] Layout mirrors correctly in Arabic

### ✅ Mobile-First
- [ ] Touch targets minimum 50px height
- [ ] Responsive breakpoints (`sm:`, `md:`, `lg:`)
- [ ] Safe area padding (`px-4`)
- [ ] Works on small screens (320px+)

### ✅ Accessibility
- [ ] All interactive elements keyboard accessible
- [ ] Focus rings visible (`focus:ring-2 focus:ring-primary/20`)
- [ ] ARIA labels on icon-only buttons
- [ ] Error messages linked to inputs
- [ ] Color contrast WCAG AA compliant

### ✅ States & Feedback
- [ ] Hover states defined
- [ ] Active/pressed states (`active:scale-95`)
- [ ] Disabled states (50% opacity)
- [ ] Loading states (spinner or skeleton)
- [ ] Error states (red border, error message)
- [ ] Focus states (ring, border color)

### ✅ Transitions
- [ ] Smooth transitions (`transition-all duration-200`)
- [ ] Consistent animation timing
- [ ] No jarring animations

---

## 10. Common Patterns & Examples

### 10.1. Standard Form

```tsx
<form className="space-y-4">
  <FormInput
    label="Email"
    type="email"
    required
    error={errors.email?.message}
    startIcon={<Mail size={20} />}
  />
  
  <PasswordInput
    label="Password"
    required
    error={errors.password?.message}
  />
  
  <Button
    type="submit"
    variant="primary"
    size="lg"
    fullWidth
    loading={isSubmitting}
  >
    Submit
  </Button>
</form>
```

### 10.2. Card List

```tsx
<div className="space-y-3">
  {items.map((item) => (
    <div
      key={item.id}
      className="card bg-base-100 border border-base-300 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
      onClick={() => handleClick(item)}
    >
      <div className="card-body p-4">
        {/* Content */}
      </div>
    </div>
  ))}
</div>
```

### 10.3. Bottom Sheet Form

```tsx
<BottomSheet
  isOpen={isOpen}
  onClose={onClose}
  title="Add Customer"
  size="md"
  footer={
    <div className="flex gap-2">
      <button className="btn btn-ghost flex-1" onClick={onClose}>
        Cancel
      </button>
      <button className="btn btn-primary flex-1" onClick={handleSave}>
        Save
      </button>
    </div>
  }
>
  <form className="space-y-4">
    {/* Form fields */}
  </form>
</BottomSheet>
```

### 10.4. Loading State

```tsx
{isLoading ? (
  <div className="space-y-3">
    {Array.from({ length: 5 }).map((_, i) => (
      <Skeleton key={i} variant="rectangular" height={100} />
    ))}
  </div>
) : (
  <div className="space-y-3">
    {items.map((item) => <ItemCard key={item.id} item={item} />)}
  </div>
)}
```

### 10.5. RTL-Aware Icon

```tsx
import { useLanguage } from "@/hooks/useLanguage";
import { ArrowLeft } from "lucide-react";

function BackButton() {
  const { isRTL } = useLanguage();
  
  return (
    <button className="btn btn-ghost btn-sm gap-2">
      <ArrowLeft size={18} className={isRTL ? "rotate-180" : ""} />
      <span>Back</span>
    </button>
  );
}
```

---

## 11. Don'ts - Common Mistakes to Avoid

### ❌ Don't Do This:

```tsx
// Wrong: Hardcoded colors
<div className="bg-[#0D9488] text-[#FFFFFF]" />

// Wrong: Using left/right
<div className="ml-4 text-left" />

// Wrong: No touch target height
<button className="h-8 px-2">Click</button>

// Wrong: No transitions
<button className="bg-primary">Click</button>

// Wrong: Missing focus ring
<button className="focus:outline-none">Click</button>

// Wrong: Icon without aria-hidden
<button><Mail size={20} /></button>

// Wrong: No error message
<input className="border-error" />

// Wrong: Inline styles
<div style={{ marginLeft: "16px" }}>Content</div>
```

### ✅ Do This Instead:

```tsx
// Correct: Semantic tokens
<div className="bg-primary text-primary-content" />

// Correct: Logical properties
<div className="ms-4 text-start" />

// Correct: Proper touch target
<button className="btn min-h-[50px]">Click</button>

// Correct: Smooth transitions
<button className="btn btn-primary transition-all duration-200">Click</button>

// Correct: Clear focus ring
<button className="btn focus:ring-2 focus:ring-primary/20 focus:outline-none">
  Click
</button>

// Correct: Icon with ARIA
<button aria-label="Send email">
  <Mail size={20} aria-hidden="true" />
</button>

// Correct: Error with message
<div>
  <input className="input input-error" aria-describedby="email-error" />
  <span id="email-error" className="text-error text-sm">
    Invalid email
  </span>
</div>

// Correct: Tailwind utilities
<div className="ms-4">Content</div>
```

---

## 12. Resources & References

- **DaisyUI Documentation:** https://daisyui.com/
- **Tailwind CSS:** https://tailwindcss.com/
- **IBM Plex Sans Arabic:** https://fonts.google.com/specimen/IBM+Plex+Sans+Arabic
- **Lucide Icons:** https://lucide.dev/
- **WCAG Guidelines:** https://www.w3.org/WAI/WCAG21/quickref/

---

## Revision History

- **v2.0** (2025-12-29): Complete rewrite based on actual portal-web implementation
  - Updated all component patterns to match real codebase
  - Added detailed form field standards
  - Documented RTL-first patterns as implemented
  - Added FilterButton component documentation
  - Removed deprecated patterns
  - Added comprehensive accessibility guidelines
  - Included real code examples from the project

- **v1.0** (Initial): Original KDS specification
