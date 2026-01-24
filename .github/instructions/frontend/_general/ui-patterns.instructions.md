---
description: Frontend UI patterns - Component structure, RTL layout, accessibility, icon usage (general, reusable across portal-web, storefront-web)
applyTo: "portal-web/**,storefront-web/**"
---

# Frontend UI Patterns

General UI component patterns, RTL support, and accessibility.

**Cross-refs:**

- i18n: `./i18n.instructions.md` (RTL detection)
- Forms: `./forms.instructions.md` (form components)

---

## 1. RTL-First Rules (Non-Negotiable)

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

---

## 2. Directional Icons

### Must Rotate 180° in RTL

```tsx
import { useLanguage } from "@/hooks/useLanguage";
import { ArrowLeft, ArrowRight } from "lucide-react";

function BackButton() {
  const { isRTL } = useLanguage();

  return (
    <button>
      <ArrowLeft className={isRTL ? "rotate-180" : ""} />
      Back
    </button>
  );
}
```

**Icons that need rotation:**

- `ArrowLeft`, `ArrowRight`
- `ArrowUpLeft`, `ArrowUpRight`, `ArrowDownLeft`, `ArrowDownRight`

### Auto-Flip (No Rotation Needed)

- `ChevronLeft`, `ChevronRight`
- `ChevronsLeft`, `ChevronsRight`

### Never Rotate

- `ChevronDown`, `ChevronUp`
- `Plus`, `X`, `Check`
- `Trash2`, `Settings`, `Menu`, `Search`, `Filter`
- Icons without directional meaning

---

## 3. Component Structure (Atomic Design)

### Atoms

Basic UI elements:

- Buttons
- Inputs
- Badges
- Tooltips
- Skeletons

**Rules:**

- No business logic
- Pure UI
- No TanStack Query
- No route params

### Molecules

Composed components:

- SearchInput (Input + Icon)
- Pagination (Buttons + Text)
- ConfirmDialog (Dialog + Buttons)
- BottomSheet (Modal + Header + Footer)

**Rules:**

- Composition of atoms
- Still generic and reusable
- No feature-specific logic

### Organisms

Complex sections:

- Table (with sorting/filtering)
- Header (app chrome)
- Sidebar (navigation)
- FilterButton (complex filter UI)

**Rules:**

- Higher-level reusable composites
- Still resource-agnostic
- If feature-specific, move to `features/**`

### Templates

Layout shells:

- AuthLayout
- AppLayout
- TwoColumnLayout

**Rules:**

- Layout-level wrappers with slots
- No resource-specific behavior

---

## 4. Common UI Patterns

### Loading State (Skeleton)

```tsx
{
  isLoading ? (
    <div className="space-y-3">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="h-24 bg-base-200 animate-pulse rounded-lg" />
      ))}
    </div>
  ) : (
    <ItemsList items={data} />
  );
}
```

### Empty State

```tsx
<div className="flex flex-col items-center justify-center py-12 text-center">
  <Package size={48} className="text-base-content/30 mb-4" />
  <h3 className="text-lg font-semibold mb-2">{t("no_products_yet")}</h3>
  <p className="text-base-content/70 mb-6">
    {t("get_started_by_adding_first_product")}
  </p>
  <button className="btn btn-primary gap-2">
    <Plus size={20} />
    {t("add_product")}
  </button>
</div>
```

### Error State

```tsx
<div className="flex flex-col items-center justify-center py-12 text-center">
  <AlertCircle size={48} className="text-error mb-4" />
  <h3 className="text-lg font-semibold mb-2">{t("error_loading_data")}</h3>
  <p className="text-base-content/70 mb-6">{error.message}</p>
  <button onClick={refetch} className="btn btn-primary">
    {t("try_again")}
  </button>
</div>
```

### Card List

```tsx
<div className="space-y-3">
  {items.map((item) => (
    <div
      key={item.id}
      className="card bg-base-100 border border-base-300 cursor-pointer hover:shadow-md transition-shadow"
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

---

## 5. Accessibility

### ARIA Attributes

| Element      | Required Attributes                                     |
| ------------ | ------------------------------------------------------- |
| Icon-only    | `aria-label="Action description"`                       |
| Form input   | `aria-invalid`, `aria-describedby`, `aria-required`     |
| Modal/Dialog | `role="dialog"`, `aria-modal="true"`, `aria-labelledby` |
| Decorative   | `aria-hidden="true"`                                    |
| Error        | `role="alert"`                                          |
| Live region  | `aria-live="polite"` or `aria-live="assertive"`         |

### Examples

```tsx
// Icon-only button (MUST have aria-label)
<button aria-label="Delete customer">
  <Trash2 size={20} />
</button>

// Icon with text (icon is decorative)
<button>
  <Plus size={20} aria-hidden="true" />
  <span>Add</span>
</button>

// Form input with error
<input
  aria-invalid={!!error}
  aria-describedby={error ? 'error-email' : undefined}
  aria-required="true"
/>
{error && <span id="error-email" role="alert">{error}</span>}

// Modal
<dialog role="dialog" aria-modal="true" aria-labelledby="modal-title">
  <h2 id="modal-title">Add Customer</h2>
  {/* content */}
</dialog>
```

### Focus Management

```tsx
// Standard focus ring
className="focus:ring-2 focus:ring-primary/20 focus:border-primary focus:outline-none"

// Focus trap (all modals)
<FocusTrap>
  <dialog>{/* content */}</dialog>
</FocusTrap>
```

### Keyboard Navigation

- **Tab**: Move between interactive elements
- **Escape**: Close modal/dropdown
- **Enter/Space**: Activate button/link
- **Arrow keys**: Navigate lists/menus

---

## 6. Touch Targets

**Minimum**: 44×44px for all interactive elements

```tsx
// ❌ WRONG - Too small
<button className="w-8 h-8">...</button>

// ✅ CORRECT - Adequate size
<button className="min-w-12 min-h-12">...</button>

// ✅ CORRECT - Full width on mobile
<button className="w-full md:w-auto">...</button>
```

---

## 7. Responsive Design

### Mobile-First Breakpoints

```
sm:  640px  (Small tablets)
md:  768px  (Tablets)
lg:  1024px (Desktops)
xl:  1280px (Large desktops)
2xl: 1536px (Extra large)
```

### Common Patterns

```tsx
// Stack on mobile, grid on desktop
<div className="space-y-4 md:grid md:grid-cols-2 md:gap-4 md:space-y-0">
  <div>Column 1</div>
  <div>Column 2</div>
</div>

// Full width on mobile, auto on desktop
<button className="w-full md:w-auto">Submit</button>

// Bottom sheet on mobile, dialog on desktop
const isMobile = useMediaQuery('(max-width: 768px)');
return isMobile ? <BottomSheet /> : <Dialog />;
```

---

## 8. Special Cases

### Phone Numbers (Always LTR)

```tsx
<span dir="ltr">{phoneNumber}</span>
```

### Dates (Locale Format)

```tsx
// Auto-handled by i18n
<span>{formatDate(date, locale)}</span>
```

### Numbers (Western Arabic Numerals)

```tsx
// Use Western Arabic numerals (0-9) even in Arabic UI
<span>{count}</span> // Not Eastern Arabic (٠-٩)
```

---

## 9. Common Component Patterns

### Modal/Bottom Sheet

```tsx
import { BottomSheet } from "@/components/molecules/BottomSheet";

function MyComponent() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <button onClick={() => setIsOpen(true)}>Open</button>

      <BottomSheet
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Add Customer"
      >
        <div className="p-4">{/* Content */}</div>
      </BottomSheet>
    </>
  );
}
```

### Toast Notifications

```tsx
import toast from "react-hot-toast";
import { useLanguage } from "@/hooks/useLanguage";

// Setup (in App.tsx)
const { isRTL } = useLanguage();
<Toaster position={isRTL ? "top-right" : "top-left"} />;

// Usage
toast.success("Order created successfully");
toast.error("Failed to save changes");
toast.loading("Saving...");
```

### Conditional Rendering

```tsx
import { useMediaQuery } from "@/hooks/useMediaQuery";

function MyComponent() {
  const isMobile = useMediaQuery("(max-width: 768px)");

  return isMobile ? <MobileView /> : <DesktopView />;
}
```

---

## 10. Performance

- **Code splitting**: `React.lazy()` for route-level components
- **Skeleton loaders**: Prevent layout shift
- **Optimize images**: WebP format, proper sizing
- **Defer non-critical**: Load after initial render
- **Transitions**: Use `transition-all` only on interactive elements

---

## Agent Validation

Before completing UI task:

- ☑ RTL: No `left`/`right` classes (use `start`/`end`, `ms-*`, `me-*`)
- ☑ RTL: Directional icons rotate 180° with `isRTL` check
- ☑ Icons: Decorative icons have `aria-hidden="true"`
- ☑ Icons: Icon-only buttons have `aria-label`
- ☑ Accessibility: Focus ring on all interactive elements
- ☑ Accessibility: Form errors have `aria-describedby`
- ☑ Touch targets: Minimum 44×44px (prefer 48px+)
- ☑ Responsive: Mobile-first approach
- ☑ Loading/empty/error states implemented
- ☑ Translation keys, not hardcoded strings

---

## Resources

- WCAG Guidelines: https://www.w3.org/WAI/WCAG21/quickref/
- Tailwind CSS: https://tailwindcss.com/docs
- Lucide Icons: https://lucide.dev/icons/
