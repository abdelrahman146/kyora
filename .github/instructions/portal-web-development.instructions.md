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
npm run dev  # Start dev server (http://localhost:3000)
```

### Dev Server Features

- **HMR:** Hot Module Replacement (instant updates)
- **TypeScript:** Type checking in IDE (not during dev server)
- **TanStack Devtools:** Route/Query debugging overlays
- **Port:** 3000 (see `portal-web/package.json` dev script)

### Type Checking

```bash
npm run type-check  # Run TypeScript compiler
```

**CI Integration:** Type check runs in GitHub Actions before merge.

---

## 2. Adding New Features

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

#### 4. Create Zod Schema

**Location:** `src/schemas/{feature}.ts`

```typescript
// src/schemas/product.ts
import { z } from "zod";

export const createProductSchema = z.object({
  name: z.string().min(2, "min_length").max(100, "max_length"),
  price: z.number().positive("positive_number"),
});

export type CreateProductFormData = z.infer<typeof createProductSchema>;
```

#### 5. Create Component

**Location:** `src/components/organisms/AddProductForm.tsx`

```typescript
import { useKyoraForm } from "@/lib/form";
import { createProductSchema } from "@/schemas/product";
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

**Location:** `src/routes/_app/products/index.tsx`

```typescript
import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { getProducts } from "@/api/product";
import { queryKeys } from "@/lib/queryKeys";
import { AddProductForm } from "@/components/organisms/AddProductForm";
import { BottomSheet } from "@/components/molecules/BottomSheet";
import { useState } from "react";

export const Route = createFileRoute("/_app/products/")({
  component: ProductsPage,
});

function ProductsPage() {
  const [isAddOpen, setIsAddOpen] = useState(false);

  const { data, isLoading } = useQuery({
    queryKey: queryKeys.products.all(),
    queryFn: () => getProducts({ page: 1 }),
  });

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Products</h1>
        <button onClick={() => setIsAddOpen(true)} className="btn btn-primary">
          Add Product
        </button>
      </div>

      {isLoading ? (
        <div>Loading...</div>
      ) : (
        <div className="space-y-3">
          {data?.map((product) => (
            <div key={product.id} className="card">
              {product.name}
            </div>
          ))}
        </div>
      )}

      <BottomSheet
        isOpen={isAddOpen}
        onClose={() => setIsAddOpen(false)}
        title="Add Product"
      >
        <AddProductForm onSuccess={() => setIsAddOpen(false)} />
      </BottomSheet>
    </div>
  );
}
```

#### 7. Add Translations (Match Existing i18next Setup)

**Rule:** Never hardcode user-facing strings (Arabic-first). Use `t("...")` keys.

For example, update UI labels in the route/component examples:

```tsx
const { t } = useTranslation()

<h1 className="text-2xl font-bold">{t("products.title")}</h1>

<button onClick={() => setIsAddOpen(true)} className="btn btn-primary">
  {t("products.add")}
</button>
```

Translations live in `src/i18n/<locale>/<namespace>.json`.

Choose the right place based on existing patterns:

- General UI text: add under `src/i18n/*/translation.json` (default namespace, dotted keys)
- Feature-heavy screens: add to a feature namespace file (e.g. `src/i18n/*/inventory.json`) using flat snake_case keys
- Validation and API error messages: `src/i18n/*/errors.json` (see forms.instructions.md for key prefixes)

Example (default namespace):

```json
// src/i18n/en/translation.json
{
  "inventory": {
    "title": "Inventory"
  }
}
```

Usage:

```tsx
import { useTranslation } from "react-i18next";

const { t } = useTranslation();
return <h1>{t("inventory.title")}</h1>;
```

#### 8. Add a New Namespace (Only If Needed)

If you create a new namespace file, you must register it in `src/i18n/init.ts`.

Checklist:

1. Create both files:

- `src/i18n/ar/<namespace>.json`
- `src/i18n/en/<namespace>.json`

2. Import both JSON files in `src/i18n/init.ts`

3. Add them to `resources.ar` and `resources.en`

4. Add `<namespace>` to the `ns` array

Usage patterns:

```tsx
const { t } = useTranslation("<namespace>");
t("some_key");

// default namespace (translation)
const { t: tTranslation } = useTranslation();
tTranslation("some.dotted_key");
```

Rules:

- Do not use `t("ns:key")` (including multi-colon strings).
- In UI components, do not use `t(key, { ns: "..." })`; bind `useTranslation("<namespace>")` instead.

---

## 3. Testing Strategy

### Unit Tests (Vitest)

**Run Tests:**

```bash
npm run test          # Run once
npm run test:watch    # Watch mode
npm run test:coverage # With coverage
```

**Test File Location:** Co-located with source (`*.test.ts`, `*.test.tsx`)

**Example:**

```typescript
// src/stores/authStore.test.ts
import { describe, it, expect, beforeEach } from "vitest";
import { authStore, loginUser, logoutUser } from "./authStore";

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

**Not yet implemented.** Planned for:

- Critical user flows (onboarding, order creation)
- Payment flows
- Multi-workspace scenarios

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

```typescript
const mutation = useMutation({
  mutationFn: deleteCustomer,
  onSuccess: (_, customerId) => {
    queryClient.invalidateQueries({ queryKey: queryKeys.customers.all() });
    queryClient.invalidateQueries({
      queryKey: queryKeys.customers.detail(customerId),
    });
    toast.success(t("customer.deleted"));
  },
});
```

### Optimistic Update

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
      queryKeys.customers.detail(updatedCustomer.id)
    );

    // Optimistically update
    queryClient.setQueryData(
      queryKeys.customers.detail(updatedCustomer.id),
      updatedCustomer
    );

    return { previous };
  },
  onError: (err, variables, context) => {
    // Rollback on error
    queryClient.setQueryData(
      queryKeys.customers.detail(variables.id),
      context.previous
    );
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
VITE_CDN_BASE_URL=https://cdn.kyora.app
```

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
