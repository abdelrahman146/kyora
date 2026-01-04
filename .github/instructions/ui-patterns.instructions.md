---
description: Kyora UI Implementation Patterns
applyTo: "portal-web/**,storefront-web/**"
---

# UI Implementation Patterns

**Prerequisite**: Read `design-tokens.instructions.md` first
**Component Location**: `/portal-web/src/components/`, `/storefront-web/src/components/`

---

## RTL-First Rules (MANDATORY)

### Forbidden Classes

```
NEVER: left, right, ml-*, mr-*, pl-*, pr-*, float-left, float-right,
       text-left, text-right, border-left, border-right
```

### Required Replacements

```
ms-*    (margin-inline-start)
me-*    (margin-inline-end)
ps-*    (padding-inline-start)
pe-*    (padding-inline-end)
start-* (positioning)
end-*   (positioning)
text-start, text-end
border-s-*, border-e-*
```

### Directional Icons

**Must rotate 180° in RTL**:

```tsx
import { useLanguage } from "@/hooks/useLanguage";
const { isRTL } = useLanguage();

<ArrowLeft className={isRTL ? "rotate-180" : ""} />;
```

**Auto-flip** (no rotation needed): `ChevronLeft`, `ChevronRight`, `ChevronsLeft`, `ChevronsRight`

**Never rotate**: `ChevronDown`, `ChevronUp`, `Plus`, `X`, `Check`, `Trash2`

### Special Cases

```tsx
// Phone numbers always LTR
<span dir="ltr">{phoneNumber}</span>
```

---

## Component Patterns

### Buttons

**Implementation**: `/portal-web/src/components/atoms/Button.tsx`

```tsx
"btn rounded-xl font-semibold transition-all active:scale-95
 disabled:opacity-50 disabled:cursor-not-allowed"
```

**Variants**: `btn-primary`, `btn-secondary`, `btn-ghost`, `btn-outline`
**Sizes**: `btn-sm` (40px), default (52px), `btn-lg` (56px)

---

### Form Inputs

**Implementation**: `/portal-web/src/components/atoms/Input.tsx`, `FormInput.tsx`

```tsx
"input input-bordered w-full h-[50px] bg-base-100 text-base-content
 focus:border-primary focus:ring-2 focus:ring-primary/20
 transition-all duration-200"
```

**Error state**: `border-error focus:border-error focus:ring-error/20`
**Icon padding**: Start icon `ps-10`, end icon `pe-10`
**Accessibility**: Label + `aria-describedby` for errors + `aria-required` for required fields

---

### Cards

```tsx
<div className="card bg-base-100 border border-base-300 shadow-sm">
  <div className="card-body">{/* Content */}</div>
</div>
```

**Hover**: `hover:shadow-md transition-shadow`
**Clickable**: Add `cursor-pointer`

---

### Modals & Bottom Sheets

**Modal**: `/portal-web/src/components/atoms/Modal.tsx`, `Dialog.tsx`
**Bottom Sheet**: `/portal-web/src/components/molecules/BottomSheet.tsx`

**Rules**:

- Portal-based rendering
- Backdrop: `backdrop-blur-sm`
- Focus trap enabled
- Escape to close
- Mobile: Bottom sheet (85% max height), Desktop: Modal

---

### Badges

**Implementation**: `/portal-web/src/components/atoms/Badge.tsx`

```tsx
"badge badge-{variant}";
```

**Variants**: `default`, `primary`, `secondary`, `success`, `error`, `warning`, `info`

---

### Skeletons

```tsx
<Skeleton variant="circular|rectangular|text" width={X} height={Y} />
// Classes: bg-base-300 animate-pulse
```

---

### Search Input

**Implementation**: `/portal-web/src/components/molecules/SearchInput.tsx`

- Debounced (300ms default)
- Search icon at `start` position (RTL-aware)
- Clear button with loading state

---

### Filter Button

**Implementation**: `/portal-web/src/components/organisms/FilterButton.tsx`

- Trigger button + drawer
- Active filter count badge
- `h-[50px]` (matches form inputs)
- Apply/Reset callbacks

---

## Accessibility Requirements

### ARIA Attributes

```
Icon-only buttons:     aria-label="Action description"
Form inputs:           aria-invalid, aria-describedby, aria-required
Modals:                role="dialog", aria-modal="true", aria-labelledby
Decorative icons:      aria-hidden="true"
Live regions:          role="alert" (for errors)
```

### Focus Management

```
Standard focus ring:   focus:ring-2 focus:ring-primary/20 focus:border-primary focus:outline-none
Focus trap:            All modals/bottom sheets trap focus
Keyboard nav:          Tab, Escape, Enter/Space, Arrow keys
```

### Color Contrast

```
Primary text:          text-base-content (WCAG AA compliant)
Secondary:             Minimum 70% opacity
Never remove focus:    Always provide visible focus indicator
```

---

## Icon System (Lucide React)

**Library**: `lucide-react` only

### Icon Usage by Context

| Context            | Icon              | RTL Rotate? | Example                 |
| ------------------ | ----------------- | ----------- | ----------------------- |
| **Actions**        |                   |             |                         |
| Add/Create         | `Plus`            | No          | Add Order               |
| Edit               | `Edit`            | No          | Edit Customer           |
| Delete             | `Trash2`          | No          | Delete                  |
| View               | `Eye`             | No          | View Details            |
| Save/Confirm       | `Check`           | No          | Saved                   |
| Close/Cancel       | `X`               | No          | Close                   |
| Search             | `Search`          | No          | Search                  |
| Filter             | `Filter`          | No          | Filter                  |
| **Navigation**     |                   |             |                         |
| Back               | `ArrowLeft`       | **Yes**     | Go Back                 |
| Next/Forward       | `ChevronRight`    | Auto        | Next Step               |
| Previous           | `ChevronLeft`     | Auto        | Previous                |
| Expand             | `ChevronDown`     | No          | Expand Menu             |
| Collapse           | `ChevronUp`       | No          | Collapse                |
| Menu               | `Menu`            | No          | Open Menu               |
| **Entities**       |                   |             |                         |
| Business           | `Building2`       | No          | Business Profile        |
| Customers (plural) | `Users`           | No          | Customers Page          |
| Customer (single)  | `User`            | No          | Customer Detail         |
| Product            | `Package`         | No          | Inventory               |
| Order              | `ShoppingCart`    | No          | Orders                  |
| Address            | `MapPin`          | No          | Location                |
| **Communication**  |                   |             |                         |
| Email              | `Mail`            | No          | Email Field             |
| Phone              | `Phone`           | No          | Phone Number            |
| Website            | `Globe`           | No          | Website Link            |
| Language Switcher  | `Languages`       | No          | Language Menu           |
| **Financial**      |                   |             |                         |
| Money/Price        | `DollarSign`      | No          | Total Revenue           |
| Accounting         | `Calculator`      | No          | Accounting Page         |
| Billing            | `CreditCard`      | No          | Billing Settings        |
| Analytics          | `BarChart3`       | No          | Analytics               |
| Dashboard          | `LayoutDashboard` | No          | Dashboard               |
| **Files**          |                   |             |                         |
| Image              | `ImageIcon`       | No          | Upload Photo            |
| Camera             | `Camera`          | No          | Take Photo              |
| Upload             | `CloudUpload`     | No          | Upload Files            |
| Document           | `FileText`        | No          | PDF                     |
| Video              | `Video`           | No          | Video File              |
| **States**         |                   |             |                         |
| Loading            | `Loader2`         | No          | `animate-spin` required |
| Success            | `CheckCircle`     | No          | Completed               |
| Error              | `AlertCircle`     | No          | Error Message           |
| Warning            | `AlertTriangle`   | No          | Warning State           |
| Info               | `HelpCircle`      | No          | Help Tooltip            |
| **Form Elements**  |                   |             |                         |
| Password Show      | `Eye`             | No          | Show Password           |
| Password Hide      | `EyeOff`          | No          | Hide Password           |
| Calendar           | `Calendar`        | No          | Date Picker             |
| Clock              | `Clock`           | No          | Time Picker             |
| Clear              | `X`               | No          | Clear Input             |
| Drag Handle        | `GripVertical`    | No          | Reorder                 |
| **Misc**           |                   |             |                         |
| Settings           | `Settings`        | No          | Settings                |
| Logout             | `LogOut`          | No          | Sign Out                |

### Icon Accessibility

```tsx
// Icon with text (decorative)
<button>
  <Plus size={20} aria-hidden="true" />
  <span>Add</span>
</button>

// Icon-only (must have aria-label)
<button aria-label="Delete customer">
  <Trash2 size={20} />
</button>
```

---

## Validation Checklist

Before shipping any UI component, verify:

```
☑ Uses daisyUI tokens (bg-primary, text-base-content)
☑ RTL: No left/right, only start/end and ms-*/me-*/ps-*/pe-*
☑ RTL: ArrowLeft rotates 180° with isRTL check
☑ Touch targets: min 50px height
☑ Focus ring: focus:ring-2 focus:ring-primary/20
☑ Icon-only buttons: aria-label present
☑ Form errors: aria-describedby links to error message
☑ Transitions: transition-all duration-200
☑ Hover/active states defined
☑ Loading state implemented (spinner or skeleton)
☑ Icons from lucide-react only
☑ Icons: aria-hidden="true" when decorative
```

---

## Common Patterns

### Standard Form

```tsx
<form className="space-y-4">
  <FormInput
    label="Email"
    type="email"
    required
    error={errors.email?.message}
    startIcon={<Mail size={20} />}
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

### Card List

```tsx
<div className="space-y-3">
  {items.map((item) => (
    <div
      key={item.id}
      className="card bg-base-100 border border-base-300 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
    >
      <div className="card-body p-4">{/* Content */}</div>
    </div>
  ))}
</div>
```

### Loading State

```tsx
{
  isLoading ? (
    <div className="space-y-3">
      {Array.from({ length: 5 }).map((_, i) => (
        <Skeleton key={i} variant="rectangular" height={100} />
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

### RTL-Aware Back Button

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

## Toast Notifications

**Library**: `react-hot-toast`

```tsx
import toast from "react-hot-toast";
import { useLanguage } from "@/hooks/useLanguage";

const { isRTL } = useLanguage();

// Setup
<Toaster position={isRTL ? "top-right" : "top-left"} />;

// Usage
toast.success("Order created");
toast.error("Failed to save");
toast.loading("Saving...");
```

**Position**: Desktop `top-right/left` (language-aware), Mobile `bottom-center`
**Auto-dismiss**: 4 seconds

---

## Performance

- Use `lazy()` for route-level code splitting
- Skeleton loaders prevent layout shift
- `transition-all` limited to interactive elements only
- Optimize images with proper formats (WebP)
- Defer non-critical scripts
