---
description: Portal Web UI components - daisyUI theme, component library usage, portal-specific icons (portal-web only)
applyTo: "portal-web/**"
---

# Portal Web UI Components

daisyUI-based component system for portal-web.

**Cross-refs:**

- General UI patterns: `../../_general/ui-patterns.instructions.md`
- Design tokens: `../../../design-tokens.instructions.md`

---

## 1. daisyUI Theme

### Kyora Portal Theme

**Primary:** `#0d9488` (Teal)
**Accent:** `#eab308` (Gold)
**Base:** `#ffffff` (White background)

**Full spec:** See `design-tokens.instructions.md`

### Theme Usage

```tsx
// ✅ CORRECT - Use semantic classes
<button className="btn btn-primary">Submit</button>

// ❌ WRONG - Don't override with Tailwind colors
<button className="btn btn-primary bg-red-500!">Submit</button>

// ✅ CORRECT - Spacing/layout utilities are OK
<button className="btn btn-primary mt-4 gap-2 w-full">Submit</button>
```

---

## 2. Component Reference

### Buttons

```tsx
// Primary action
<button className="btn btn-primary">Submit</button>

// Secondary action
<button className="btn btn-secondary">Cancel</button>

// Ghost (minimal)
<button className="btn btn-ghost">Back</button>

// Outline
<button className="btn btn-outline">More</button>

// With icon
<button className="btn btn-primary gap-2">
  <Plus size={20} />
  Add Order
</button>

// Sizes
<button className="btn btn-xs">Extra Small</button>
<button className="btn btn-sm">Small</button>
<button className="btn">Default</button>
<button className="btn btn-lg">Large</button>

// Full width (mobile)
<button className="btn btn-primary w-full">Submit</button>
```

### Cards

```tsx
<div className="card bg-base-100 border border-base-300">
  <div className="card-body">
    <h2 className="card-title">Card Title</h2>
    <p>Card content goes here</p>
    <div className="card-actions justify-end">
      <button className="btn btn-primary">Action</button>
    </div>
  </div>
</div>
```

### Badges

```tsx
<span className="badge badge-primary">New</span>
<span className="badge badge-secondary">Featured</span>
<span className="badge badge-success">Active</span>
<span className="badge badge-error">Failed</span>
<span className="badge badge-warning">Pending</span>
<span className="badge badge-info">Draft</span>

// Sizes
<span className="badge badge-xs">Extra Small</span>
<span className="badge badge-sm">Small</span>
<span className="badge">Default</span>
<span className="badge badge-lg">Large</span>
```

### Loading

```tsx
<span className="loading loading-spinner loading-lg"></span>
<span className="loading loading-dots loading-md"></span>
<span className="loading loading-ring loading-sm"></span>
```

### Progress

```tsx
<progress className="progress progress-primary" value="70" max="100"></progress>
<progress className="progress progress-success" value="40" max="100"></progress>
```

---

## 3. Custom Components

### BottomSheet

```tsx
import { BottomSheet } from "@/components/molecules/BottomSheet";

<BottomSheet
  isOpen={isOpen}
  onClose={onClose}
  title="Add Customer"
  footer={
    <div className="flex gap-2">
      <Button variant="ghost" onClick={onClose}>
        Cancel
      </Button>
      <Button variant="primary">Save</Button>
    </div>
  }
>
  <div className="p-4">{/* Content */}</div>
</BottomSheet>;
```

### StatCard

```tsx
import { StatCard, StatCardGroup } from "@/components";

<StatCardGroup cols={3}>
  <StatCard
    label="Total Revenue"
    value="$12,450"
    icon={<DollarSign className="h-5 w-5" />}
    trend="up"
    trendValue="+12.5%"
    variant="success"
  />
  <StatCard
    label="Orders"
    value="150"
    icon={<ShoppingCart className="h-5 w-5" />}
    trend="down"
    trendValue="-5%"
    variant="warning"
  />
</StatCardGroup>;
```

### ComplexStatCard

```tsx
<ComplexStatCard
  label="Cash on Hand"
  value="$25,000"
  comparisonText="vs last month: +$2,500"
  statusBadge={{ label: "Healthy", variant: "success" }}
  secondaryMetrics={[
    { label: "Cash In", value: "$30,000" },
    { label: "Cash Out", value: "$5,000" },
  ]}
/>
```

---

## 4. Icon System (Lucide React)

### Common Portal Icons

| Context        | Icon              | RTL Rotate? | Usage           |
| -------------- | ----------------- | ----------- | --------------- |
| **Navigation** |                   |             |                 |
| Dashboard      | `LayoutDashboard` | No          | Dashboard page  |
| Orders         | `ShoppingCart`    | No          | Orders page     |
| Customers      | `Users`           | No          | Customers page  |
| Products       | `Package`         | No          | Inventory page  |
| Analytics      | `BarChart3`       | No          | Analytics page  |
| Accounting     | `Calculator`      | No          | Accounting page |
| Reports        | `FileText`        | No          | Reports page    |
| **Actions**    |                   |             |                 |
| Add            | `Plus`            | No          | Add button      |
| Edit           | `Edit`            | No          | Edit button     |
| Delete         | `Trash2`          | No          | Delete button   |
| Back           | `ArrowLeft`       | **YES**     | Back button     |
| Next           | `ArrowRight`      | **YES**     | Next button     |
| **Financial**  |                   |             |                 |
| Money          | `DollarSign`      | No          | Revenue, price  |
| Billing        | `CreditCard`      | No          | Payment methods |
| **States**     |                   |             |                 |
| Success        | `CheckCircle`     | No          | Success message |
| Error          | `AlertCircle`     | No          | Error message   |
| Warning        | `AlertTriangle`   | No          | Warning state   |
| Info           | `HelpCircle`      | No          | Help tooltip    |

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

## 5. Portal-Specific Patterns

### Dashboard Cards

```tsx
<div className="card bg-base-100 border border-base-300">
  <div className="card-body p-4">
    <div className="flex items-center justify-between mb-4">
      <h3 className="font-semibold">{t("recent_orders")}</h3>
      <button className="btn btn-ghost btn-sm gap-1">
        {t("view_all")}
        <ChevronRight size={16} />
      </button>
    </div>

    <div className="space-y-3">
      {orders.map((order) => (
        <OrderCard key={order.id} order={order} />
      ))}
    </div>
  </div>
</div>
```

### Empty State (Portal Style)

```tsx
<div className="flex flex-col items-center justify-center py-12 text-center">
  <div className="w-16 h-16 rounded-full bg-base-200 flex items-center justify-center mb-4">
    <Package size={32} className="text-base-content/30" />
  </div>
  <h3 className="text-lg font-semibold mb-2">
    {t("inventory:no_products_yet")}
  </h3>
  <p className="text-base-content/70 mb-6 max-w-sm">
    {t("inventory:get_started_by_adding_first_product")}
  </p>
  <button className="btn btn-primary gap-2">
    <Plus size={20} />
    {t("inventory:add_product")}
  </button>
</div>
```

### List Item Card

```tsx
<div className="card bg-base-100 border border-base-300 hover:shadow-md transition-shadow cursor-pointer">
  <div className="card-body p-4">
    <div className="flex items-start gap-3">
      <div className="avatar placeholder">
        <div className="bg-base-200 text-base-content rounded-full w-12">
          <span className="text-lg">{initials}</span>
        </div>
      </div>

      <div className="flex-1 min-w-0">
        <h3 className="font-semibold truncate">{name}</h3>
        <p className="text-sm text-base-content/70 truncate">{email}</p>
      </div>

      <ChevronRight size={20} className="text-base-content/30 flex-shrink-0" />
    </div>
  </div>
</div>
```

---

## 6. Portal Color Usage

### Status Colors

```tsx
// Success (green)
<span className="text-success">Active</span>
<div className="bg-success/10 text-success">Completed</div>

// Warning (yellow)
<span className="text-warning">Pending</span>
<div className="bg-warning/10 text-warning">Draft</div>

// Error (red)
<span className="text-error">Failed</span>
<div className="bg-error/10 text-error">Cancelled</div>

// Info (blue)
<span className="text-info">Processing</span>
<div className="bg-info/10 text-info">Shipped</div>
```

### Text Colors

```tsx
// Primary text
<p className="text-base-content">Main content</p>

// Secondary text (70% opacity)
<p className="text-base-content/70">Helper text</p>

// Muted text (50% opacity)
<p className="text-base-content/50">Disabled text</p>
```

---

## Agent Validation

Before completing UI task:

- ☑ Using daisyUI semantic classes (`.btn`, `.card`, `.badge`)
- ☑ NOT overriding daisyUI colors (only spacing/layout)
- ☑ Icons from `lucide-react` only
- ☑ Directional icons rotate with `isRTL` check
- ☑ Icon-only buttons have `aria-label`
- ☑ Status colors use semantic variants
- ☑ Portal-specific patterns followed (Dashboard cards, empty states)
- ☑ RTL-safe layout (use `start`/`end`, not `left`/`right`)

---

## Resources

- daisyUI Docs: https://daisyui.com/components/
- Lucide Icons: https://lucide.dev/icons/
- Design Tokens: `../../../design-tokens.instructions.md`
