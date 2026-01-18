---
description: Portal Web Development - Workflow, Testing, Deployment
applyTo: "portal-web/**"
---

# Portal Web Development Guide

**SSOT Hierarchy:**

- Parent: copilot-instructions.md
- Peers: portal-web-architecture.instructions.md
- Implementation Guides: forms.instructions.md, ui-implementation.instructions.md

**When to Read:**

- Starting development work
- Adding new features
- Testing strategies
- Deployment processes

---

## 1. Development Workflow

### Setup

```bash
cd portal-web
npm install
npm run dev  # Start dev server (default port: 3000, configurable via VITE_DEV_PORT)
```

### Dev Server Features

- **HMR:** Hot Module Replacement (instant updates)
- **TypeScript:** Type checking in IDE (not during dev server)
- **TanStack Devtools:** Route/Query debugging overlays
- **Port:** Default 3000 (configurable via `VITE_DEV_PORT` env var)
- **Network Access:** Server allows access from other devices on same network (`host: true`)
- **Mobile Testing:** Use `VITE_DEV_HOST` env var to set explicit HMR websocket host when needed

### Type Checking

```bash
npm run type-check  # Run TypeScript compiler
```

**CI Integration:** Type check runs in GitHub Actions before merge.

---

## 2. Adding New Features

**File placement SSOT:** `.github/instructions/portal-web-code-structure.instructions.md`

- `src/routes/**` must stay a thin routing layer (schema/loader + render wrapper).
- Page implementations, feature-specific components, feature-specific schemas, and feature utilities must live under `src/features/<feature>/**`.

### Step-by-Step Process

#### 1. Create API Module

**Location:** `src/api/{feature}.ts`

```typescript
// src/api/product.ts
import { apiClient } from "./client";
import type { Product, CreateProductRequest } from "./types/product";

export async function getProducts(params: { page: number }) {
  return apiClient
    .get("v1/products", { searchParams: params })
    .json<Product[]>();
}

export async function getProduct(id: string) {
  return apiClient.get(`v1/products/${id}`).json<Product>();
}

export async function createProduct(data: CreateProductRequest) {
  return apiClient.post("v1/products", { json: data }).json<Product>();
}
```

#### 2. Add Query Keys

**Location:** `src/lib/queryKeys.ts`

```typescript
export const queryKeys = {
  // ... existing
  products: {
    all: (params?: Record<string, any>) => ["products", params] as const,
    detail: (id: string) => ["products", id] as const,
  },
};
```

#### 3. Create Types

**Location:** `src/api/types/{feature}.ts`

```typescript
// src/api/types/product.ts
export interface Product {
  id: string;
  name: string;
  price: number;
  // ...
}

export interface CreateProductRequest {
  name: string;
  price: number;
  // ...
}
```

#### 4. Create Feature Schema (if feature-specific)

**Location:** `src/features/<feature>/schema/*`

If the schema is only used by one feature, keep it inside that feature (recommended for new code).

**Legacy Note:** Some older schemas exist in `src/schemas/` (e.g., `auth.ts`, `onboarding.ts`, `inventory.ts`). These are gradually being migrated to feature-specific locations per `.github/instructions/portal-web-code-structure.instructions.md`.

For new feature-specific schemas, always create them under the feature:

```typescript
// src/features/products/schema/createProduct.ts
import { z } from "zod";

export const createProductSchema = z.object({
  name: z
    .string()
    .min(2, "validation.min_length")
    .max(100, "validation.max_length"),
  price: z.number().positive("validation.positive_number"),
});

export type CreateProductFormData = z.infer<typeof createProductSchema>;
```

**Cross-cutting schemas** (used by 2+ distinct features, not resource-specific) may live in `src/schemas/`, but prefer feature-specific placement.

#### 5. Create Feature Components / Page

**Location:** `src/features/<feature>/components/*`

```typescript
import { useKyoraForm } from "@/lib/form";
import { createProductSchema } from "@/features/products/schema/createProduct";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createProduct } from "@/api/product";
import { queryKeys } from "@/lib/queryKeys";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";

export function AddProductForm({ onSuccess }: { onSuccess: () => void }) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: createProduct,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.products.all() });
      toast.success(t("product.created"));
      onSuccess();
    },
  });

  const form = useKyoraForm({
    defaultValues: {
      name: "",
      price: 0,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync(value);
    },
  });

  return (
    <form.AppForm>
      <form.FormRoot className="space-y-4">
        <form.FormError />

        <form.AppField
          name="name"
          validators={{
            onBlur: createProductSchema.shape.name,
          }}
        >
          {(field) => <field.TextField label={t("product.name")} required />}
        </form.AppField>

        <form.AppField
          name="price"
          validators={{
            onBlur: createProductSchema.shape.price,
          }}
        >
          {(field) => (
            <field.TextField
              type="number"
              label={t("product.price")}
              required
            />
          )}
        </form.AppField>

        <form.SubmitButton variant="primary" className="w-full">
          {t("product.add")}
        </form.SubmitButton>
      </form.FormRoot>
    </form.AppForm>
  );
}
```

#### 6. Create Route

**Location:** `src/routes/**` (thin wrapper)

Routes must only wire URL/search, loaders/prefetch, and render the feature page.

```typescript
import { createFileRoute } from "@tanstack/react-router";
import { ProductsListPage } from "@/features/products/components/ProductsListPage";

export const Route = createFileRoute("/business/$businessDescriptor/products/")(
  {
    component: ProductsRoute,
  }
);

function ProductsRoute() {
  return <ProductsListPage />;
}
```

#### 7. Add Translations (Match Existing i18next Setup)

**Rule:** Never hardcode user-facing strings (Arabic-first). Use `t("...")` keys with explicit namespace binding.

For example, update UI labels in the route/component examples:

```tsx
// Always bind to explicit namespace (SSOT pattern)
const { t: tInventory } = useTranslation("inventory");
const { t: tCommon } = useTranslation("common");

<h1 className="text-2xl font-bold">{tInventory("title")}</h1>

<button onClick={() => setIsAddOpen(true)} className="btn btn-primary">
  {tInventory("add_product")}
</button>

<button onClick={onClose} className="btn btn-ghost">
  {tCommon("actions.cancel")}
</button>
```

Translations live in `src/i18n/<locale>/<namespace>.json`.

Choose the right place based on existing patterns (portal-web is namespace-only):

- Shared UI primitives: `src/i18n/*/common.json`
- Feature-specific UI strings: `src/i18n/*/<feature>.json` (e.g. `inventory.json`, `orders.json`, `customers.json`)
- Errors and API error messaging: `src/i18n/*/errors.json`

Example (feature namespace):

```json
// src/i18n/en/inventory.json
{
  "title": "Inventory",
  "add_product": "Add product"
}
```

Usage:

```tsx
import { useTranslation } from "react-i18next";

// Always use explicit namespace binding (not default namespace)
const { t: tInventory } = useTranslation("inventory");
return <h1>{tInventory("title")}</h1>;
```

**Important:** Do not use `t("namespace:key")` pattern (multi-colon). Always bind the namespace with `useTranslation("namespace")` and use bare keys.

See `.github/instructions/i18n-translations.instructions.md` for the SSOT rules (explicit namespaces, locale parity, no duplication).

#### 8. Add a New Namespace (Only If Needed)

If you create a new namespace file, you must register it in `src/i18n/init.ts`.

Checklist:

1. Create both files:
   - `src/i18n/ar/<namespace>.json`
   - `src/i18n/en/<namespace>.json`

2. Import both JSON files in `src/i18n/init.ts`:

   ```ts
   import arNamespace from "./ar/namespace.json";
   import enNamespace from "./en/namespace.json";
   ```

3. Add them to the `resources` object:

   ```ts
   resources: {
     ar: {
       // ... existing
       namespace: arNamespace,
     },
     en: {
       // ... existing
       namespace: enNamespace,
     },
   }
   ```

4. Add `<namespace>` to the `ns` array:
   ```ts
   ns: [
     "common",
     "errors",
     // ... existing
     "namespace",
   ];
   ```

Usage patterns:

```tsx
// Always use explicit namespace binding
const { t: tNamespace } = useTranslation("namespace");
tNamespace("some_key");

// For common namespace
const { t: tCommon } = useTranslation("common");
tCommon("actions.save");
```

**Forbidden patterns:**

- ❌ `t("ns:key")` - Multi-colon pattern is forbidden
- ❌ `t(key, { ns: "..." })` - Do not use namespace option in UI code
- ❌ `useTranslation()` without namespace - Always specify explicit namespace

---

## 3. Testing Strategy

### Unit Tests (Vitest)

**Run Tests:**

```bash
npm run test          # Run once (vitest run)
```

**Note:** `test:watch` and `test:coverage` scripts are not currently configured. To run in watch mode or with coverage, use:

```bash
npx vitest watch        # Watch mode
npx vitest --coverage    # With coverage
```

**Test Configuration:** Vitest is configured via `vite.config.ts` (no separate `vitest.config.ts`).

**Test File Location:** Co-located with source (`*.test.ts`, `*.test.tsx`)

**Current Test Coverage:**

- Only one test file exists: `src/stores/authStore.test.ts`
- Tests must include `// @vitest-environment jsdom` comment at the top for DOM-dependent tests
- Use `beforeEach` to ensure consistent state between tests

**Example:**

```typescript
// src/stores/authStore.test.ts
// @vitest-environment jsdom

import { describe, it, expect, beforeEach } from "vitest";
import { authStore, initializeAuth } from "./authStore";

describe("authStore", () => {
  beforeEach(() => {
    authStore.setState({ user: null, isAuthenticated: false });
  });

  it("should set user on login", async () => {
    const mockUser = { id: "1", email: "test@example.com" };
    await loginUser(mockUser);

    expect(authStore.state.user).toEqual(mockUser);
    expect(authStore.state.isAuthenticated).toBe(true);
  });

  it("should clear user on logout", () => {
    authStore.setState({ user: mockUser, isAuthenticated: true });
    logoutUser();

    expect(authStore.state.user).toBeNull();
    expect(authStore.state.isAuthenticated).toBe(false);
  });
});
```

### E2E Tests (Playwright - Future)

**Status:** Not yet implemented. No Playwright configuration exists in the repository.

**Planned for:**

- Critical user flows (onboarding, order creation)
- Payment flows
- Multi-workspace scenarios
- Cross-browser testing

**When implementing:**

1. Add `@playwright/test` to devDependencies
2. Create `playwright.config.ts` in `portal-web/`
3. Add test scripts to `package.json`: `test:e2e`, `test:e2e:ui`
4. Create test files in `portal-web/tests/` or `portal-web/e2e/`

---

## 4. Code Organization

### Atomic Design

```
components/
├── atoms/          # Basic UI elements (Button, Input, Badge)
├── molecules/      # Composed components (SearchInput, BottomSheet)
├── organisms/      # Complex sections (LoginForm, FilterButton)
└── templates/      # Page layouts (AuthLayout, AppLayout)
```

**Guidelines:**

- **Atoms:** No business logic, pure UI
- **Molecules:** Composition of atoms with minimal logic
- **Organisms:** Feature-specific, contains business logic
- **Templates:** Layout shells, no feature-specific content

### File Naming

| Type       | Convention                  | Example             |
| ---------- | --------------------------- | ------------------- |
| Components | PascalCase                  | `CustomerCard.tsx`  |
| Hooks      | camelCase with `use` prefix | `useAuth.ts`        |
| Utils      | camelCase                   | `formatCurrency.ts` |
| Types      | PascalCase                  | `Customer.ts`       |
| Constants  | UPPER_SNAKE_CASE            | `API_BASE_URL`      |

### Imports Order

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

## 5. Common Patterns

### Mutation with Invalidation

Current pattern used in the codebase:

```typescript
const mutation = useMutation({
  mutationFn: deleteCustomer,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: queryKeys.customers.all() });
    queryClient.invalidateQueries({
      queryKey: queryKeys.customers.detail(customerId),
    });
    toast.success(t("customer.deleted"));
  },
});
```

**Note:** The codebase currently uses `void queryClient.invalidateQueries()` pattern in some places (e.g., accounting API), but both with/without `void` work.

### Optimistic Update

**Status:** Pattern documented but not currently implemented in the codebase.

The codebase mentions "optimistic updates" in comments (FileUploadField, Notes component) but does not use the full `onMutate`/`onError`/`onSettled` pattern shown below. Most mutations use simple `onSuccess` invalidation.

If implementing optimistic updates in the future, follow this pattern:

```typescript
const mutation = useMutation({
  mutationFn: updateCustomer,
  onMutate: async (updatedCustomer) => {
    // Cancel outgoing refetches
    await queryClient.cancelQueries({
      queryKey: queryKeys.customers.detail(updatedCustomer.id),
    });

    // Snapshot previous value
    const previous = queryClient.getQueryData(
      queryKeys.customers.detail(updatedCustomer.id),
    );

    // Optimistically update
    queryClient.setQueryData(
      queryKeys.customers.detail(updatedCustomer.id),
      updatedCustomer,
    );

    return { previous };
  },
  onError: (err, variables, context) => {
    // Rollback on error
    if (context?.previous) {
      queryClient.setQueryData(
        queryKeys.customers.detail(variables.id),
        context.previous,
      );
    }
  },
  onSettled: (_, __, variables) => {
    queryClient.invalidateQueries({
      queryKey: queryKeys.customers.detail(variables.id),
    });
  },
});
```

### Loading + Empty + Error States

```typescript
const { data, isLoading, error } = useQuery({
  queryKey: queryKeys.customers.all(),
  queryFn: getCustomers,
});

if (error) {
  return <ErrorMessage error={error} />;
}

if (isLoading) {
  return <SkeletonList count={5} />;
}

if (!data || data.length === 0) {
  return <EmptyState message={t("customer.no_customers")} />;
}

return (
  <div className="space-y-3">
    {data.map((customer) => (
      <CustomerCard key={customer.id} customer={customer} />
    ))}
  </div>
);
```

### Conditional Rendering (Mobile vs Desktop)

```typescript
import { useMediaQuery } from "@/hooks/useMediaQuery";

function MyComponent() {
  const isMobile = useMediaQuery("(max-width: 768px)");

  return isMobile ? <BottomSheet>...</BottomSheet> : <Dialog>...</Dialog>;
}
```

---

## 6. Build & Deployment

### Build Commands

```bash
npm run build         # Production build → dist/
npm run preview       # Preview prod build locally
```

### Build Output

```
dist/
├── index.html
├── assets/
│   ├── index-abc123.js
│   ├── index-def456.css
│   └── chunks/          # Code-split chunks
└── ...
```

### Environment Variables

**Location:** `.env.local` (not committed)

```env
VITE_API_BASE_URL=https://api.kyora.app
```

**Available Environment Variables:**

- `VITE_API_BASE_URL` - Backend API server URL (default: `http://localhost:8080`)
- `VITE_DEV_PORT` - Dev server port (default: `3000`)
- `VITE_DEV_HOST` - HMR websocket host override (optional, for mobile testing)

**Usage:**

```typescript
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
```

**Production:** Set via Vercel/Netlify environment variables.

---

## 7. Performance Optimization

### Code Splitting

**Automatic:** TanStack Router splits by route.

**Manual:**

```typescript
import { lazy } from "react";

const HeavyComponent = lazy(() => import("./HeavyComponent"));

function MyPage() {
  return (
    <Suspense fallback={<Skeleton />}>
      <HeavyComponent />
    </Suspense>
  );
}
```

### Bundle Analysis

```bash
npm run build -- --mode analyze
```

Opens bundle visualizer in browser.

### Query Stale Time

```typescript
const { data } = useQuery({
  queryKey: queryKeys.metadata.countries(),
  queryFn: getCountries,
  staleTime: 1000 * 60 * 60, // 1 hour (countries rarely change)
});
```

### Image Optimization

```tsx
// Use WebP with fallback
<picture>
  <source srcSet="/image.webp" type="image/webp" />
  <img src="/image.jpg" alt="Description" />
</picture>
```

---

## 8. Debugging

### React Devtools

Install extension: https://react.dev/learn/react-developer-tools

**Features:**

- Component tree inspection
- Props/state visualization
- Profiler for performance

### TanStack Devtools

**Location:** Floating button in bottom-right corner (dev mode only)

**Router Devtools:**

- View route tree
- Active route params
- Loader data

**Query Devtools:**

- Query cache state
- Mutation status
- Refetch/invalidate manually

### Network Debugging

**Browser DevTools → Network tab:**

- Filter by XHR
- Inspect request/response payloads
- Check auth headers

**Common Issues:**

- 401: Access token expired → check auto-refresh logic
- 422: Validation error → check RFC 7807 `errors` field
- 500: Server error → check backend logs

---

## 9. Common Issues & Solutions

### Issue: "Cannot read property of undefined"

**Cause:** Data not loaded yet, accessing `data.field` before `data` exists.

**Solution:**

```typescript
// ❌ Bad
<div>{data.customer.name}</div>

// ✅ Good
<div>{data?.customer.name ?? 'Loading...'}</div>

// ✅ Better
if (!data) return <Skeleton />
return <div>{data.customer.name}</div>
```

### Issue: "Query not updating after mutation"

**Cause:** Missing invalidation.

**Solution:**

```typescript
onSuccess: () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.customers.all() });
};
```

### Issue: "Form not submitting"

**Cause:** Missing `<form.AppForm>` wrapper.

**Solution:** See forms.instructions.md → Critical Rules

### Issue: "Icons not rotating in RTL"

**Cause:** Missing `isRTL` check.

**Solution:** See ui-implementation.instructions.md → RTL Rules

---

## 10. CI/CD Pipeline

### GitHub Actions

**Location:** `.github/workflows/portal-web.yml`

**Steps:**

1. Install dependencies
2. Type check (`npm run type-check`)
3. Lint (`npm run lint`)
4. Test (`npm run test`)
5. Build (`npm run build`)
6. Deploy to Vercel (main branch only)

**Branch Protection:**

- All checks must pass
- No direct push to main

---

## Agent Validation Checklist

Before completing development task:

- ☑ API module created in `src/api/`
- ☑ Types defined in `src/api/types/`
- ☑ Query keys added to `queryKeys` factory
- ☑ Zod schema created in `src/schemas/`
- ☑ Component follows Atomic Design hierarchy
- ☑ Translations added (Arabic + English)
- ☑ Route created with TanStack Router file convention
- ☑ Mutation invalidates related queries
- ☑ Loading/empty/error states implemented
- ☑ TypeScript: No `any` types
- ☑ RTL: No `left`/`right` classes
- ☑ Accessibility: ARIA labels on icon-only buttons
- ☑ Tests: Unit tests for business logic (stores, utils)

---

## See Also

- **Architecture:** `.github/instructions/portal-web-architecture.instructions.md`
- **Forms:** `.github/instructions/forms.instructions.md`
- **UI Patterns:** `.github/instructions/ui-implementation.instructions.md`
- **HTTP Client:** `.github/instructions/ky.instructions.md`

---

## Resources

- Vite Docs: https://vite.dev/
- Vitest Docs: https://vitest.dev/
- TanStack Router: https://tanstack.com/router/latest
- TanStack Query: https://tanstack.com/query/latest
- Project README: `portal-web/README.md`
