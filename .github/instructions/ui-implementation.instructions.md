---
description: Kyora UI Implementation - RTL, daisyUI, Icons, Accessibility
applyTo: "portal-web/**,storefront-web/**"
---

# UI Implementation Patterns

**SSOT Hierarchy:**

- Parent: copilot-instructions.md
- Peers: forms.instructions.md, charts.instructions.md
- Required Reading: design-tokens.instructions.md

**When to Read:**

- Implementing UI components
- RTL/LTR layout work
- Icon usage decisions
- Accessibility requirements
- daisyUI component selection

**Portal Web UX/UI SSOT:** `.github/instructions/portal-web-ui-guidelines.instructions.md`

---

## 1. RTL-First Rules (NON-NEGOTIABLE)

### Forbidden Classes

```
NEVER USE:
  left, right, ml-*, mr-*, pl-*, pr-*,
  float-left, float-right,
  text-left, text-right,
  border-left, border-right
```

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

### Directional Icons (CRITICAL)

**Must rotate 180° in RTL:**

```tsx
import { useLanguage } from '@/hooks/useLanguage'

const { isRTL } = useLanguage()

// ✅ CORRECT
<ArrowLeft className={isRTL ? 'rotate-180' : ''} />
<ArrowRight className={isRTL ? 'rotate-180' : ''} />
```

**Auto-flip (no rotation needed):**
`ChevronLeft`, `ChevronRight`, `ChevronsLeft`, `ChevronsRight`

**Never rotate:**
`ChevronDown`, `ChevronUp`, `Plus`, `X`, `Check`, `Trash2`, `Settings`, `Menu`, `Search`, `Filter`

### Special Cases

```tsx
// Phone numbers always LTR
<span dir="ltr">{phoneNumber}</span>

// Dates in locale format (auto-handled by i18n)
<span>{formatDate(date, locale)}</span>
```

---

## 2. daisyUI Component Usage

**Official Docs:** https://daisyui.com/components/

### Core Rules

1. **Always use semantic classes** (`.btn`, `.input`, `.card`)
2. **Never override with Tailwind utilities** (exception: spacing/layout only)
3. **Use modifiers for variants** (`.btn-primary`, `.btn-lg`)
4. **Custom colors via design tokens** (see design-tokens.instructions.md)

### Component Selection Guide

```
Need interactive element?
├─ Form input? → See forms.instructions.md
├─ Button? → Use daisyUI .btn
├─ Data viz? → See charts.instructions.md
└─ Custom? → Build with daisyUI base + Tailwind spacing
```

### Common Components

**Buttons:**

```tsx
<button className="btn btn-primary btn-lg">Submit</button>
```

Variants: `btn-primary`, `btn-secondary`, `btn-ghost`, `btn-outline`
Sizes: `btn-xs`, `btn-sm`, default, `btn-lg`

**Inputs (use form system instead):**

```tsx
// ❌ Don't use raw daisyUI inputs
<input className="input input-bordered" />

// ✅ Use form system
<form.TextField name="email" />
// See: forms.instructions.md
```

**Cards:**

```tsx
<div className="card bg-base-100 border border-base-300">
  <div className="card-body">{/* Content */}</div>
</div>
```

**Badges:**

```tsx
<span className="badge badge-primary">New</span>
```

Variants: `badge-primary`, `badge-secondary`, `badge-success`, `badge-error`, `badge-warning`, `badge-info`

**Modals (use BottomSheet/Dialog instead):**

```tsx
// ❌ Don't use raw daisyUI modals
<dialog className="modal">...</dialog>

// ✅ Use custom components
<BottomSheet isOpen={isOpen} onClose={onClose}>
  {/* Content */}
</BottomSheet>
// Note: file placement is governed by portal-web-code-structure.instructions.md
// Current implementation location may change during refactors.
```

**Loading:**

```tsx
<span className="loading loading-spinner loading-lg"></span>
```

**Progress:**

```tsx
<progress className="progress progress-primary" value="70" max="100"></progress>
```

### Spacing with Tailwind (ALLOWED)

```tsx
// ✅ Spacing/layout utilities are OK
<button className="btn btn-primary mt-4 gap-2 w-full">
  Submit
</button>

// ❌ Never override daisyUI colors/sizes
<button className="btn btn-primary bg-red-500!">  {/* WRONG! */}
  Submit
</button>
```

### Custom Theme (Kyora)

**Primary:** `#0d9488` (Teal)
**Accent:** `#eab308` (Gold)
**Full spec:** See design-tokens.instructions.md

---

## 3. Icon System (Lucide React)

**Library:** `lucide-react` ONLY (no other icon libraries)

### Icon Reference Table

| Context            | Icon              | RTL Rotate? | Usage                      |
| ------------------ | ----------------- | ----------- | -------------------------- |
| **Actions**        |                   |             |                            |
| Add/Create         | `Plus`            | No          | Add Order, Create Customer |
| Edit               | `Edit`            | No          | Edit Product               |
| Delete             | `Trash2`          | No          | Remove Item                |
| View               | `Eye`             | No          | View Details               |
| Save/Confirm       | `Check`           | No          | Saved Successfully         |
| Close/Cancel       | `X`               | No          | Close Modal                |
| Search             | `Search`          | No          | Search Bar                 |
| Filter             | `Filter`          | No          | Filter Button              |
| **Navigation**     |                   |             |                            |
| Back               | `ArrowLeft`       | **YES**     | Go Back Button             |
| Forward            | `ArrowRight`      | **YES**     | Next Page                  |
| Next               | `ChevronRight`    | Auto        | Next Step                  |
| Previous           | `ChevronLeft`     | Auto        | Previous Step              |
| Expand             | `ChevronDown`     | No          | Expand Menu                |
| Collapse           | `ChevronUp`       | No          | Collapse Section           |
| Menu               | `Menu`            | No          | Hamburger Menu             |
| **Entities**       |                   |             |                            |
| Business           | `Building2`       | No          | Business Profile           |
| Customers (plural) | `Users`           | No          | Customers Page             |
| Customer (single)  | `User`            | No          | Customer Detail            |
| Product            | `Package`         | No          | Inventory Page             |
| Order              | `ShoppingCart`    | No          | Orders List                |
| Address            | `MapPin`          | No          | Location Pin               |
| **Communication**  |                   |             |                            |
| Email              | `Mail`            | No          | Email Field                |
| Phone              | `Phone`           | No          | Phone Number               |
| Website            | `Globe`           | No          | Website Link               |
| Language           | `Languages`       | No          | Language Switcher          |
| **Financial**      |                   |             |                            |
| Money/Price        | `DollarSign`      | No          | Revenue, Price             |
| Accounting         | `Calculator`      | No          | Accounting Page            |
| Billing            | `CreditCard`      | No          | Payment Methods            |
| Analytics          | `BarChart3`       | No          | Analytics Dashboard        |
| Dashboard          | `LayoutDashboard` | No          | Main Dashboard             |
| **Files**          |                   |             |                            |
| Image              | `ImageIcon`       | No          | Upload Photo               |
| Camera             | `Camera`          | No          | Take Picture               |
| Upload             | `CloudUpload`     | No          | File Upload                |
| Document           | `FileText`        | No          | PDF, Docs                  |
| Video              | `Video`           | No          | Video Files                |
| **States**         |                   |             |                            |
| Loading            | `Loader2`         | No          | `animate-spin` required    |
| Success            | `CheckCircle`     | No          | Success Message            |
| Error              | `AlertCircle`     | No          | Error Message              |
| Warning            | `AlertTriangle`   | No          | Warning State              |
| Info               | `HelpCircle`      | No          | Help Tooltip               |
| **Form Elements**  |                   |             |                            |
| Password Show      | `Eye`             | No          | Toggle Password            |
| Password Hide      | `EyeOff`          | No          | Hide Password              |
| Calendar           | `Calendar`        | No          | Date Picker                |
| Clock              | `Clock`           | No          | Time Picker                |
| Clear              | `X`               | No          | Clear Input                |
| Drag Handle        | `GripVertical`    | No          | Reorder Items              |
| **Misc**           |                   |             |                            |
| Settings           | `Settings`        | No          | Settings Page              |
| Logout             | `LogOut`          | No          | Sign Out                   |

### Icon Sizing

```tsx
// Standard sizes
<Plus size={16} />  // Small (inline with text)
<Plus size={20} />  // Default (buttons, inputs)
<Plus size={24} />  // Large (page headers)
<Plus size={32} />  // Extra large (empty states)
```

### Icon Accessibility

```tsx
// Icon with text (icon is decorative)
<button className="btn gap-2">
  <Plus size={20} aria-hidden="true" />
  <span>Add</span>
</button>

// Icon-only (MUST have aria-label)
<button className="btn btn-ghost btn-sm" aria-label="Delete customer">
  <Trash2 size={20} />
</button>

// Icon in loading state
<Loader2 size={20} className="animate-spin" aria-label="Loading" />
```

### RTL Icon Example

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

## 4. Accessibility Requirements

### ARIA Attributes

| Element Type     | Required Attributes                                     |
| ---------------- | ------------------------------------------------------- |
| Icon-only button | `aria-label="Action description"`                       |
| Form input       | `aria-invalid`, `aria-describedby`, `aria-required`     |
| Modal/Dialog     | `role="dialog"`, `aria-modal="true"`, `aria-labelledby` |
| Decorative icon  | `aria-hidden="true"`                                    |
| Error message    | `role="alert"`                                          |
| Live region      | `aria-live="polite"` or `aria-live="assertive"`         |
| List             | `role="list"` on container, `role="listitem"` on items  |

### Focus Management

**Standard focus ring:**

```tsx
className =
  "focus:ring-2 focus:ring-primary/20 focus:border-primary focus:outline-none";
```

**Focus trap:** All modals/bottom sheets must trap focus
**Keyboard nav:** Support Tab, Escape, Enter/Space, Arrow keys

### Color Contrast

- **Primary text:** `text-base-content` (WCAG AA compliant)
- **Secondary text:** Minimum 70% opacity
- **Never remove focus:** Always provide visible focus indicator
- **Test with:** Browser devtools accessibility panel

### Touch Targets

**Minimum**: 44×44px (prefer 48px+) for all interactive elements.

- Avoid hardcoded pixel heights.
- Prefer Tailwind scale (`min-h-12`, `min-w-12`) and daisyUI sizing (`btn-lg`).
- Mobile primary actions should usually be full-width: `w-full`.

---

## 5. Common UI Patterns

### Standard Form Layout

```tsx
<form className="space-y-4">
  <form.TextField name="email" label="Email" required />
  <form.TextField name="name" label="Full Name" required />
  <form.SubmitButton variant="primary" className="w-full">
    Submit
  </form.SubmitButton>
</form>
```

### Card List

```tsx
<div className="space-y-3">
  {items.map((item) => (
    <div
      key={item.id}
      className="card bg-base-100 border border-base-300 cursor-pointer"
      onClick={() => navigate(`/items/${item.id}`)}
    >
      <div className="card-body p-4">
        <h3 className="font-semibold">{item.name}</h3>
        <p className="text-sm text-base-content/70">{item.description}</p>
      </div>
    </div>
  ))}
</div>
```

### Loading State (Skeleton)

```tsx
{
  isLoading ? (
    <div className="space-y-3">
      {Array.from({ length: 5 }).map((_, i) => (
        <div
          key={i}
          className="card bg-base-100 border border-base-300 h-24 animate-pulse"
        />
      ))}
    </div>
  ) : (
    <div className="space-y-3">
      {items.map((item) => (
        <ItemCard key={item.id} item={item} />
      ))}
    </div>
  );
}
```

### Empty State

```tsx
<div className="flex flex-col items-center justify-center py-12 text-center">
  <Package size={48} className="text-base-content/30 mb-4" />
  <h3 className="text-lg font-semibold mb-2">No products yet</h3>
  <p className="text-base-content/70 mb-6">
    Get started by adding your first product
  </p>
  <button className="btn btn-primary gap-2">
    <Plus size={20} />
    Add Product
  </button>
</div>
```

### Modal/Bottom Sheet Pattern

```tsx
import { BottomSheet } from "@/components/molecules/BottomSheet";
import { useState } from "react";

function MyComponent() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <button onClick={() => setIsOpen(true)} className="btn btn-primary">
        Open
      </button>

      <BottomSheet
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Add Customer"
      >
        <div className="p-4">{/* Form content */}</div>
      </BottomSheet>
    </>
  );
}
```

### Toast Notifications

```tsx
import toast from 'react-hot-toast'
import { useLanguage } from '@/hooks/useLanguage'

// Setup (in App.tsx)
const { isRTL } = useLanguage()
<Toaster position={isRTL ? 'top-right' : 'top-left'} />

// Usage
toast.success('Order created successfully')
toast.error('Failed to save changes')
toast.loading('Saving...')
```

---

## 6. Responsive Design

### Mobile-First Breakpoints

```
sm:  640px  (Small tablets)
md:  768px  (Tablets)
lg:  1024px (Desktops)
xl:  1280px (Large desktops)
2xl: 1536px (Extra large)
```

### Mobile vs Desktop Patterns

| Component    | Mobile                   | Desktop           |
| ------------ | ------------------------ | ----------------- |
| Modal        | BottomSheet (85% height) | Dialog (centered) |
| Select       | Action sheet             | Dropdown          |
| Date picker  | Full-screen              | Popover           |
| Navigation   | Bottom bar               | Sidebar           |
| Form buttons | Full width (`w-full`)    | Auto width        |

### Portal Web Style Constraints

- **No shadows. No gradients.** Prefer borders + spacing for separation.
- Use daisyUI semantic classes for components; Tailwind utilities for layout/spacing only.

### Example

```tsx
<button className="btn btn-primary w-full md:w-auto">Submit</button>
```

---

## 7. Performance Patterns

- **Code splitting:** Use `lazy()` for route-level components
- **Skeleton loaders:** Prevent layout shift during loading
- **Optimize images:** WebP format, proper sizing
- **Defer non-critical:** Scripts load after initial render
- **Transitions:** Use `transition-all` only on interactive elements

---

## Agent Validation Checklist

Before completing UI task:

- ☑ RTL: No `left`/`right` classes (use `start`/`end`, `ms-*`, `me-*`, `ps-*`, `pe-*`)
- ☑ RTL: `ArrowLeft`/`ArrowRight` icons rotate 180° with `isRTL` check
- ☑ daisyUI: Using semantic classes (`.btn`, `.card`, `.input`)
- ☑ daisyUI: Not overriding colors (only spacing/layout utilities)
- ☑ Icons: From `lucide-react` only, correct icon per context
- ☑ Icons: Decorative icons have `aria-hidden="true"`
- ☑ Icons: Icon-only buttons have `aria-label`
- ☑ Accessibility: Focus ring on all interactive elements
- ☑ Accessibility: Form errors have `aria-describedby`
- ☑ Touch targets: Minimum 50px height
- ☑ Responsive: Mobile-first approach (base styles, then `md:` overrides)
- ☑ Loading: Skeleton states implemented
- ☑ Empty states: Helpful messaging with action button

---

## See Also

- **Design Tokens:** `.github/instructions/design-tokens.instructions.md` → Colors, typography, spacing
- **Forms:** `.github/instructions/forms.instructions.md` → Form components, validation
- **Charts:** `.github/instructions/charts.instructions.md` → Data visualization
- **Backend:** `.github/instructions/backend-core.instructions.md` → API contracts

---

## Resources

- daisyUI Docs: https://daisyui.com/components/
- Lucide Icons: https://lucide.dev/icons/
- Tailwind CSS: https://tailwindcss.com/docs
- WCAG Guidelines: https://www.w3.org/WAI/WCAG21/quickref/
- Implementation: `portal-web/src/components/`
