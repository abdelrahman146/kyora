---
description: Portal Web code structure - Routes, features, components organization, placement rules (portal-web only)
applyTo: "portal-web/**"
---

# Portal Web Code Structure

File placement rules for portal-web.

**Cross-refs:**

- General architecture: `../../_general/architecture.instructions.md`
- Portal architecture: `./architecture.instructions.md`

---

## 1. Directory Layout

```
portal-web/src/
├── api/                    # HTTP client + domain APIs
├── components/             # Shared UI (resource-agnostic)
│   ├── atoms/
│   ├── molecules/
│   ├── organisms/
│   ├── templates/
│   ├── form/
│   ├── charts/
│   └── icons/
├── features/               # Feature modules
│   ├── orders/
│   ├── inventory/
│   ├── customers/
│   ├── accounting/
│   ├── analytics/
│   ├── auth/
│   ├── business-switcher/
│   ├── dashboard-layout/
│   └── onboarding/
├── hooks/                  # Custom hooks
├── i18n/                   # Translations (ar/, en/)
├── lib/                    # Cross-cutting utils
├── routes/                 # File-based routes
├── stores/                 # TanStack Store instances
└── types/                  # TypeScript types
```

---

## 2. Routes (Thin Routing Layer)

**Location:** `src/routes/**`

**Allowed:**

- `Route = createFileRoute()` config
- `validateSearch` Zod schema
- `loader` / `beforeLoad` / `loaderDeps`
- `staticData` (e.g., title keys)
- Thin `component` wrapper that renders Feature Page

**Forbidden:**

- Page UI implementation (tables/cards/sheets/layouts)
- Domain/resource business logic
- Feature-specific components

**Example:**

```tsx
// src/routes/business/$businessDescriptor/orders/index.tsx
export const Route = createFileRoute("/business/$businessDescriptor/orders/")({
  validateSearch: z.object({
    page: z.number().default(1),
    status: z.string().optional(),
  }),
  loader: async ({ context, params }) => {
    return context.queryClient.ensureQueryData(
      orderQueries.list(params.businessDescriptor, { page: 1 }),
    );
  },
  staticData: {
    titleKey: "pages.orders", // From common.pages.orders
  },
  component: OrdersRoute,
});

function OrdersRoute() {
  return <OrdersListPage />; // From features/orders/
}
```

---

## 3. Shared Components (Resource-Agnostic)

**Location:** `src/components/**`

**atoms/**: Small UI primitives

- Button, Badge, Input, Dialog, Tooltip, Skeleton
- No TanStack Query, no route params, no feature logic

**molecules/**: Compositions of atoms

- SearchInput, Pagination, ConfirmDialog, BottomSheet
- Still generic and reusable

**organisms/**: Complex sections

- Table, Header, Sidebar (app chrome)
- Still resource-agnostic
- If resource-specific → move to `features/**`

**templates/**: Layout shells

- AuthLayout, TwoColumnLayout
- No resource-specific behavior

**form/**: Generic form controls

- Usable with any form library
- Not tied to TanStack Form

**charts/**: Reusable Chart.js components

- Generic, receive `data`/`options` via props
- Not analytics-specific

**icons/**: Custom icons

- Only when `lucide-react` doesn't provide what we need

**Strict rule:** Shared component names MUST NOT include resource names (Order, Customer, Product, etc.)

---

## 4. Feature Modules

**Location:** `src/features/<feature>/`

**Feature types:**

- **Resource features**: orders, inventory, customers, accounting, analytics, reports
- **Cross-cutting features**: auth, business-switcher, language, onboarding, home
- **Layout features**: dashboard-layout, app-shell

**Feature structure:**

```
features/<feature>/
├── components/     # Feature-specific UI + pages
├── schema/         # Zod schemas (if feature-specific)
├── state/          # Feature-local state (optional)
├── utils/          # Feature-specific helpers
└── types/          # Feature-specific types
```

**Rules:**

- Feature components may compose shared components
- Feature code may use Query hooks from `api/**`
- Feature code is reusable across the app
- If used by 1 feature only → keep it in that feature
- If truly cross-cutting → move to `lib/`

---

## 5. lib/ (Cross-Cutting Utils)

**Location:** `src/lib/**`

A utility belongs in `lib/` only if ALL are true:

- Used by 2+ distinct features OR clearly cross-cutting
- Not tied to single resource vocabulary
- Has stable API and is expected to be shared

**Examples of what belongs in lib/:**

- `errorParser.ts` (RFC7807 parsing)
- `queryKeys.ts` (centralized query keys)
- `routeGuards.ts` (auth guards)
- `form/` (form system)
- `charts/` (chart utilities)
- `upload/` (file upload utilities)

**Examples of what does NOT belong in lib/:**

- `shippingZone.ts` → `features/orders/utils/`
- `inventoryUtils.ts` → `features/inventory/utils/`
- `onboarding.ts` → `features/onboarding/utils/`

---

## 6. Decision Tree

```
Need to add code?
├─ Backend call? → api/
├─ Route? → routes/ (thin wrapper only)
├─ Page implementation? → features/<feature>/components/
├─ Shared UI (resource-agnostic)? → components/
├─ Feature-specific component? → features/<feature>/components/
├─ Zod schema (1 feature)? → features/<feature>/schema/
├─ Zod schema (2+ features)? → schemas/
├─ Feature utility? → features/<feature>/utils/
└─ Cross-cutting utility? → lib/
```

---

## 7. Known Drifts

These patterns violate current structure and should be fixed:

1. **Resource-specific organisms**:
   - `components/organisms/ShippingZoneInfo.tsx` → `features/orders/components/`
   - `components/organisms/Notes.tsx` → `features/customers/components/` or `features/notes/`

2. **Resource-specific form fields in lib/**:
   - `lib/form/components/CategorySelectField.tsx` → `features/inventory/components/fields/`
   - `lib/form/components/ProductVariantSelectField.tsx` → `features/inventory/components/fields/`
   - `lib/form/components/CustomerSelectField.tsx` → `features/customers/components/fields/`

3. **Feature-specific utilities in lib/**:
   - `lib/shippingZone.ts` → `features/orders/utils/` or `features/business/utils/`

When adding similar code, follow correct patterns, not drift patterns.

---

## 8. File Naming

| Type       | Convention                  | Example             |
| ---------- | --------------------------- | ------------------- |
| Components | PascalCase                  | `CustomerCard.tsx`  |
| Hooks      | camelCase with `use` prefix | `useAuth.ts`        |
| Utils      | camelCase                   | `formatCurrency.ts` |
| Types      | PascalCase                  | `Customer.ts`       |
| Constants  | UPPER_SNAKE_CASE            | `API_BASE_URL`      |

---

## 9. Imports Order

```typescript
// 1. External libraries
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

// 2. Internal modules (@/ alias)
import { getCustomers } from "@/api/customer";
import { queryKeys } from "@/lib/queryKeys";
import { Button } from "@/components/atoms/Button";

// 3. Types
import type { Customer } from "@/api/types/customer";

// 4. Styles (if any)
import "./styles.css";
```

---

## Agent Validation

Before completing code structure task:

- ☑ Routes are thin wrappers (no page UI implementation)
- ☑ Page implementations in `features/<feature>/components/`
- ☑ Shared components are resource-agnostic
- ☑ Feature-specific code in `features/<feature>/`
- ☑ Feature utilities NOT in `lib/`
- ☑ Resource names NOT in shared component names
- ☑ File naming conventions followed
- ☑ Imports ordered correctly

---

## Resources

- Portal routes: `src/routes/`
- Shared components: `src/components/`
- Feature modules: `src/features/`
- Cross-cutting utils: `src/lib/`
