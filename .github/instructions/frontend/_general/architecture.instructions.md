---
description: Frontend architecture patterns - TanStack stack, routing, state, feature modules (reusable across all Kyora frontends)
applyTo: "portal-web/**,storefront-web/**,mobile-web/**"
---

# Frontend Architecture

Kyora frontend apps share a common tech stack and architectural patterns.

**Cross-refs:**

- Portal-specific: `../projects/portal-web/architecture.instructions.md`
- State management: Covered in this file (see State Management section)

---

## 1. Tech Stack

### Core

- **React 19+** - Concurrent features, RSC-ready
- **TypeScript 5.7+** - Strict mode
- **Vite** - Dev server + build

### Routing & Data

- **TanStack Router** - File-based routing, type-safe navigation
- **TanStack Query** - Server state, caching, auto-refetch
- **Ky** - HTTP client (retry, timeout, hooks)

### Forms & State

- **TanStack Form** - Form state
- **TanStack Store** - Client state (auth, preferences)
- **Zod** - Validation

### Styling

- **Tailwind CSS** - Utility-first
- **CSS-first config** - No JS config
- **RTL-native** - Arabic-first

### i18n

- **i18next** - Translation framework
- **react-i18next** - React bindings

---

## 2. File Structure

```
src/
├── api/                    # HTTP client + domain APIs
│   ├── client.ts           # Ky instance
│   ├── {domain}.ts         # Domain API methods
│   └── types/              # API types
├── components/             # Shared UI (Atomic Design)
│   ├── atoms/              # Buttons, inputs, badges
│   ├── molecules/          # Composed components
│   ├── organisms/          # Complex sections
│   └── templates/          # Layout shells
├── features/               # Feature modules
│   └── {feature}/
│       ├── components/     # Feature UI
│       ├── schema/         # Zod schemas
│       ├── state/          # Feature state
│       └── utils/          # Feature helpers
├── hooks/                  # Custom hooks
├── i18n/                   # Translations
│   ├── ar/                 # Arabic (primary)
│   └── en/                 # English (fallback)
├── lib/                    # Cross-cutting utils
├── routes/                 # File-based routes
├── stores/                 # TanStack Store instances
└── types/                  # TypeScript types
```

---

## 3. Feature Modules

Features are cohesive slices of functionality:

- **Resource features**: orders, inventory, customers, accounting
- **Cross-cutting features**: auth, business-switcher, language, onboarding
- **Layout features**: dashboard-layout, app-shell

**Feature structure:**

```
features/{feature}/
├── components/     # Feature-specific UI
├── schema/         # Zod schemas (if feature-specific)
├── state/          # Feature state (optional)
├── utils/          # Feature helpers
└── types/          # Feature types
```

**Rules:**

- Features may compose shared components
- Features may use Query hooks from `api/**`
- Feature code is reusable across the app
- If used by 1 feature only → keep it in that feature
- If truly cross-cutting → move to `lib/`

---

## 4. Routing

### File-Based Routes

TanStack Router auto-generates route tree from `src/routes/**`.

**Pattern:**

```tsx
// src/routes/business/$businessDescriptor/orders/index.tsx
export const Route = createFileRoute("/business/$businessDescriptor/orders/")({
  validateSearch: z.object({
    page: z.number().default(1),
    status: z.string().optional(),
  }),
  loader: async ({ context, params }) => {
    return context.queryClient.ensureQueryData(
      ordersQueries.list(params.businessDescriptor, { page: 1 }),
    );
  },
  component: OrdersRoute,
});

function OrdersRoute() {
  return <OrdersListPage />;
}
```

**Route responsibilities:**

- URL/search validation
- Data prefetch (loader)
- Render wrapper (thin)

**Page implementation lives in:** `features/{feature}/components/`

---

## 5. Data Fetching

### Query Pattern

```typescript
// src/api/customer.ts
export const customerApi = {
  async list(params?: ListParams): Promise<Customer[]> {
    return get<Customer[]>(`v1/customers`, { searchParams: params });
  },
};

export const customerQueries = {
  all: (params?: ListParams) => ({
    queryKey: ["customers", params] as const,
    queryFn: () => customerApi.list(params),
  }),
};

// In component
const { data, isLoading } = useQuery(customerQueries.all({ page: 1 }));
```

### Mutation Pattern

```typescript
const mutation = useMutation({
  mutationFn: customerApi.create,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["customers"] });
    toast.success(t("customer.created"));
  },
});
```

**Query keys:** Centralized factory in `lib/queryKeys.ts`

---

## 6. State Management

### Decision Tree

```
What kind of state?
├─ From backend? → TanStack Query
├─ In URL (filters/page)? → TanStack Router
├─ Form fields? → TanStack Form
└─ Client-only? → TanStack Store
```

### TanStack Store

**Use for:**

- Auth session
- User preferences
- UI state (sidebar collapsed)
- "Last selected" convenience state

**Don't:**

- Mirror Query data (creates dual caches)
- Store form values (use TanStack Form)
- Store route state (use URL search params)

**Pattern:**

```typescript
import { Store } from "@tanstack/store";

export const authStore = new Store({
  user: null as User | null,
  isAuthenticated: false,
});

// Update
authStore.setState((state) => ({
  ...state,
  user: newUser,
  isAuthenticated: true,
}));
```

---

## 7. HTTP Client

### Configuration

```typescript
import ky from "ky";

export const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL,
  timeout: 30000,
  retry: {
    limit: 2,
    methods: ["get", "put", "head", "delete", "options", "trace"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
  },
  hooks: {
    beforeRequest: [
      (request) => {
        const token = getAccessToken();
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (request, options, response) => {
        if (response.status === 401) {
          const newToken = await refreshAccessToken();
          if (newToken) {
            request.headers.set("Authorization", `Bearer ${newToken}`);
            return ky(request);
          }
        }
      },
    ],
  },
});
```

**Auto-refresh:** 401 triggers token refresh + retry automatically

---

## 8. Authentication

### Token Storage

| Token         | Storage           | Lifespan | Purpose      |
| ------------- | ----------------- | -------- | ------------ |
| Access Token  | Memory (variable) | 15min    | API requests |
| Refresh Token | HTTP-only cookie  | 30d      | Rotation     |

**No localStorage** - session restored from cookie only

### Auth State

```typescript
// src/stores/authStore.ts
{
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}
```

**Initialization:**

1. App mounts → `initializeAuth()`
2. Calls `restoreSession()` → `GET /v1/auth/me`
3. If refresh cookie valid → user restored

### Route Guards

```typescript
// src/lib/routeGuards.ts
export async function requireAuth() {
  await initializeAuth();
  const { isAuthenticated } = authStore.state;
  if (!isAuthenticated) {
    throw redirect({ to: "/auth/login" });
  }
}

// In route
export const Route = createFileRoute("/protected")({
  beforeLoad: requireAuth,
  component: ProtectedPage,
});
```

---

## 9. Internationalization

### Configuration

```typescript
// src/i18n/init.ts
i18next.init({
  resources: {
    en: { common, errors /* ... */ },
    ar: { common, errors /* ... */ },
  },
  lng: detectLanguage(), // Cookie → browser → 'en'
  fallbackLng: "en",
  defaultNS: "common",
  ns: ["common", "errors" /* ... */],
});

// Set document attributes
document.documentElement.lang = i18next.language;
document.documentElement.dir = i18next.language === "ar" ? "rtl" : "ltr";
```

### Usage

```tsx
// Always use explicit namespace
const { t } = useTranslation("orders");
t("order_number"); // orders:order_number

// For errors
const { t: tErrors } = useTranslation("errors");
tErrors("http.404"); // errors:http.404
```

**Forbidden:** `t('ns:key')` pattern

---

## 10. Error Handling

### RFC 7807 Problem JSON

Backend returns:

```json
{
  "status": 401,
  "title": "Unauthorized",
  "detail": "Invalid credentials",
  "extensions": {
    "code": "account.invalid_credentials"
  }
}
```

**Portal parsing:**

1. Ky `beforeError` extracts code
2. `translateError()` maps to i18n key: `errors.backend.account.invalid_credentials`
3. Global handler shows toast

---

## 11. Performance

### Code Splitting

- **Auto:** TanStack Router splits by route
- **Manual:** `React.lazy()` for heavy components

### Query Optimization

```typescript
const { data } = useQuery({
  queryKey: ["metadata", "countries"],
  queryFn: getCountries,
  staleTime: 1000 * 60 * 60, // 1 hour (rarely changes)
});
```

### Bundle Analysis

```bash
npm run build -- --mode analyze
```

---

## 12. Testing

### Unit Tests (Vitest)

```typescript
// Co-located with source: *.test.ts
describe("authStore", () => {
  it("should set user on login", () => {
    authStore.setState({ user: mockUser, isAuthenticated: true });
    expect(authStore.state.user).toEqual(mockUser);
  });
});
```

### E2E Tests (Playwright)

```typescript
// tests/e2e/orders.spec.ts
test("should create order", async ({ page }) => {
  await page.goto("/orders");
  await page.click("text=Add Order");
  await page.fill('[name="customer"]', "John Doe");
  await page.click("text=Submit");
  await expect(page.locator("text=Order created")).toBeVisible();
});
```

---

## 13. Common Patterns

### Mutation with Invalidation

```typescript
const mutation = useMutation({
  mutationFn: deleteCustomer,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["customers"] });
    toast.success(t("customer.deleted"));
  },
});
```

### Loading + Empty + Error

```tsx
const { data, isLoading, error } = useQuery(query);

if (error) return <ErrorMessage error={error} />;
if (isLoading) return <SkeletonList />;
if (!data || data.length === 0) return <EmptyState />;

return <List items={data} />;
```

### Conditional Rendering (Mobile vs Desktop)

```tsx
import { useMediaQuery } from "@/hooks/useMediaQuery";

function MyComponent() {
  const isMobile = useMediaQuery("(max-width: 768px)");
  return isMobile ? <BottomSheet /> : <Dialog />;
}
```

---

## 14. Architecture Decisions

### Why TanStack?

| Tool          | Chosen          | Rationale                |
| ------------- | --------------- | ------------------------ |
| Router        | TanStack Router | Type-safe, file-based    |
| Data Fetching | TanStack Query  | Best caching, React 19   |
| Forms         | TanStack Form   | Flexible, composable     |
| State         | TanStack Store  | Reactive, minimal API    |
| HTTP          | Ky              | Modern, TypeScript-first |

### No Redux

**Reason:** TanStack Query handles server state, TanStack Store handles client state. Redux is overkill.

### No Axios

**Reason:** Ky is modern, smaller bundle, better retry logic.

---

## Agent Validation

Before completing architecture task:

- ☑ Using TanStack Router file-based routing
- ☑ Query keys from centralized factory
- ☑ Auth logic through store, not custom hooks
- ☑ Route guards applied (`requireAuth`)
- ☑ HTTP via configured Ky client
- ☑ State: Query (server), Router (URL), Store (client)
- ☑ Translation keys, not hardcoded strings
- ☑ RTL-first CSS

---

## Resources

- TanStack Router: https://tanstack.com/router/latest
- TanStack Query: https://tanstack.com/query/latest
- TanStack Form: https://tanstack.com/form/latest
- TanStack Store: https://tanstack.com/store/latest
