---
description: Portal Web Architecture - Tech Stack, Patterns, State Management
applyTo: "portal-web/**"
---

# Portal Web Architecture

**SSOT Hierarchy:**

- Parent: copilot-instructions.md
- Peers: portal-web-development.instructions.md
- Related: forms.instructions.md, ui-implementation.instructions.md, ky.instructions.md

**When to Read:**

- Starting portal-web work
- Understanding tech stack decisions
- Auth/routing/state patterns
- Architecture questions

---

## 1. Overview

**Portal Web** is the business management dashboard for Kyora - a React SPA for Arabic-first social commerce entrepreneurs.

**Philosophy:** "Professional tools that feel effortless" - zero accounting knowledge required.

**Key Characteristics:**

- Mobile-first, RTL-native
- Arabic primary, English fallback
- Middle East social commerce focus (Instagram/WhatsApp/TikTok sellers)
- Workspace-based multi-tenancy with business sub-scope (workspace can have multiple businesses)
- RBAC: admin/member roles

**Tenancy & Scoping (SSOT):**

- **Workspace-scoped domains** (global to the workspace; not tied to a single business):
  - Account/auth/session, workspace membership & RBAC
  - Billing/subscription (Stripe)
  - Workspace-level settings/metadata
- **Business-owned domains** (must always be scoped to the current business):
  - Orders
  - Inventory
  - Customers
  - Analytics
  - Accounting
  - Assets
  - Storefront
  - Onboarding (when creating/configuring the business)

**Frontend scoping rule:**

- Business-owned API calls must include `businessDescriptor` in the URL: `v1/businesses/${businessDescriptor}/...`.
- Business-owned UI routes must be nested under `/business/$businessDescriptor/...`.
- Never fetch or mutate business-owned data without an explicit `businessDescriptor`.

---

## 2. Tech Stack (Definitive)

### Core

- **React 19.2.0** - Concurrent features enabled
- **TypeScript 5.7.2** - Strict mode, full type safety
- **Vite 7.1.7** - Dev server, build tool

### Routing & Data

- **TanStack Router 1.132.0** - File-based routing, type-safe navigation, code splitting
- **TanStack Query 5.66.5** - Server state, caching, auto-refetch, optimistic updates
- **Ky 1.14.2** - HTTP client with retry logic, auth interceptors

### Forms & Validation

- **TanStack Form 1.27.7** - Form state management
- **TanStack Store 0.8.0** - Reactive state (auth, business, metadata, onboarding stores)
- **Zod 4.2.1** - Schema validation, type inference
- **Custom `useKyoraForm`** - Zero-boilerplate composition layer

### Styling & UI

- **Tailwind CSS 4.0.6** - CSS-first (no JS config)
- **daisyUI 5.5.14** - Semantic component classes
- **clsx 2.1.1 + tailwind-merge 3.4.0** - Dynamic classNames
- **lucide-react 0.561.0** - Icons

### Data Visualization

- **Chart.js 4.5.1** - Canvas-based charts
- **react-chartjs-2 5.3.1** - React wrapper
- **chartjs-adapter-date-fns 3.0.0** - Time-series adapter
- **date-fns 4.1.0** - Date formatting

### UI Enhancements

- **react-hot-toast 2.6.0** - Toast notifications
- **react-day-picker 9.13.0** - Date/range picker
- **@dnd-kit 6.3.1** - Drag-drop (FieldArray reordering)
- **@ffmpeg/ffmpeg 0.12.15** - Video thumbnails (client-side)

### i18n

- **i18next 25.7.3** - i18n framework
- **react-i18next 16.5.0** - React bindings

### Developer Tools

- **TanStack Router Devtools** - Route debugging
- **TanStack Query Devtools** - Query cache visualization
- **ESLint + Prettier** - Linting, formatting
- **Vitest 3.0.5** - Unit testing

---

## 3. Project Structure

```
portal-web/
├── src/
│   ├── api/                    # HTTP client + API modules
│   │   ├── client.ts           # Ky client with auth interceptors
│   │   ├── auth.ts             # Auth API
│   │   ├── user.ts             # User API
│   │   ├── business.ts         # Business/workspace API
│   │   ├── customer.ts         # Customer API
│   │   ├── assets.ts           # File upload API
│   │   └── types/              # API types
│   ├── components/             # UI components (Atomic Design)
│   │   ├── atoms/              # Buttons, inputs, badges
│   │   ├── molecules/          # SearchInput, BottomSheet
│   │   ├── organisms/          # LoginForm, FilterButton
│   │   └── templates/          # Page layouts
│   ├── hooks/                  # Custom hooks
│   │   ├── useAuth.ts
│   │   ├── useLanguage.ts
│   │   └── useMediaQuery.ts
│   ├── i18n/                   # Translations
│   │   ├── ar/                 # Arabic (primary)
│   │   └── en/                 # English (fallback)
│   ├── lib/                    # Utilities
│   │   ├── auth.ts             # Auth logic
│   │   ├── errorParser.ts      # RFC7807 parser
│   │   ├── queryKeys.ts        # Query key factory
│   │   ├── routeGuards.ts      # Route auth guards
│   │   ├── form/               # Form system
│   │   ├── charts/             # Chart.js utils
│   │   └── upload/             # File upload utils
│   ├── routes/                 # File-based routes
│   │   ├── __root.tsx          # Root layout
│   │   ├── _auth/              # Auth layout group
│   │   └── _app/               # App layout group
│   ├── schemas/                # Zod schemas
│   ├── stores/                 # TanStack Store instances
│   │   ├── authStore.ts
│   │   ├── businessStore.ts
│   │   ├── metadataStore.ts
│   │   └── onboardingStore.ts
│   ├── types/                  # TypeScript types
│   └── main.tsx                # Entry point
├── public/                     # Static assets
└── vite.config.ts              # Vite config
```

---

## 4. Authentication Flow

### Token Management

| Token Type    | Storage           | Lifespan      | Purpose        |
| ------------- | ----------------- | ------------- | -------------- |
| Access Token  | Memory (variable) | Short (15min) | API requests   |
| Refresh Token | HTTP-only cookie  | Long (30d)    | Token rotation |

### HTTP Client Auto-Refresh

**Implementation:** `src/api/client.ts`

```typescript
import ky from "ky";

const apiClient = ky.create({
  hooks: {
    // 1. Inject access token
    beforeRequest: [
      (request) => {
        const token = getAccessToken();
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`);
        }
      },
    ],

    // 2. Auto-refresh on 401
    afterResponse: [
      async (request, options, response) => {
        if (response.status === 401) {
          const newToken = await refreshAccessToken(); // POST /v1/auth/refresh
          if (newToken) {
            request.headers.set("Authorization", `Bearer ${newToken}`);
            return ky(request); // Retry
          }
          clearTokens();
          window.location.href = "/auth/login";
        }
      },
    ],
  },
});
```

**Concurrency Protection:** Single refresh promise prevents thundering herd.

### Auth State

**Store:** `authStore` (TanStack Store) in `src/stores/authStore.ts`

```typescript
{
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}
```

**Initialization:**

1. App mounts → `initializeAuth()` called
2. Calls `restoreSession()` → `GET /v1/auth/me`
3. If refresh cookie valid → user restored
4. If not → `isAuthenticated = false`

**No Local Storage:** Session restored from HTTP-only refresh cookie only.

### Route Guards

**Location:** `src/lib/routeGuards.ts`

```typescript
// Require authentication
export const requireAuth = () => {
  const { isAuthenticated } = authStore.state;
  if (!isAuthenticated) {
    throw redirect({ to: "/auth/login" });
  }
};

// Require guest (redirect if logged in)
export const requireGuest = () => {
  const { isAuthenticated } = authStore.state;
  if (isAuthenticated) {
    throw redirect({ to: "/" });
  }
};
```

**Usage in routes:**

```typescript
// src/routes/_app/dashboard.tsx
export const Route = createFileRoute("/_app/dashboard")({
  beforeLoad: requireAuth,
});
```

---

## 5. Data Fetching (TanStack Query)

### Query Keys Factory

**Location:** `src/lib/queryKeys.ts`

```typescript
export const queryKeys = {
  auth: {
    me: () => ["auth", "me"] as const,
  },
  customers: {
    all: (params?: Record<string, any>) => ["customers", params] as const,
    detail: (id: string) => ["customers", id] as const,
  },
  // ... more
};
```

### Query Pattern

```typescript
const { data, isLoading, error } = useQuery({
  queryKey: queryKeys.customers.all({ page: 1 }),
  queryFn: () => getCustomers({ page: 1 }),
});
```

### Mutation Pattern

```typescript
const mutation = useMutation({
  mutationFn: createCustomer,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: queryKeys.customers.all() });
    toast.success(t("customer.created"));
  },
  onError: (error) => {
    toast.error(translateError(error));
  },
});
```

### Cache Invalidation

**Helper:** `src/lib/queryInvalidation.ts`

```typescript
// Invalidate related queries after action
await invalidateCustomerQueries(customerId);
```

---

## 6. State Management

### Server State

**TanStack Query** - API data, caching, auto-refetch

### Client State

**TanStack Store** - Local UI state with persistence

**Stores:**

1. **authStore** (`src/stores/authStore.ts`)

   - User session state
   - No persistence (restored from cookie)

2. **businessStore** (`src/stores/businessStore.ts`)

- Business list + selected business (`selectedBusinessDescriptor`)
- Sidebar UI state
- Persists preferences to `localStorage` via `createPersistencePlugin`

3. **metadataStore** (`src/stores/metadataStore.ts`)

   - Countries/currencies cache
   - Persisted to `localStorage`

4. **onboardingStore** (`src/stores/onboardingStore.ts`)
   - Onboarding wizard state
   - Persisted to `sessionStorage`

**Store Pattern:**

```typescript
import { Store } from "@tanstack/store";

export const myStore = new Store({
  user: null,
  settings: {},
});

// Usage
myStore.setState((state) => ({
  ...state,
  user: newUser,
}));

// Subscribe
myStore.subscribe(() => {
  console.log(myStore.state.user);
});
```

**Persistence:** See `src/lib/storePersistence.ts`

---

## 7. Routing (TanStack Router)

### File-Based Routes

```
src/routes/
├── __root.tsx              → Root layout
├── index.tsx               → / (homepage)
├── _auth/                  → Auth layout group (requireGuest)
│   ├── login.tsx           → /auth/login
│   └── register.tsx        → /auth/register
└── business/                → Business-scoped routes
  └── $businessDescriptor/ → /business/:businessDescriptor (requireAuth)
    ├── index.tsx        → Business home
    ├── orders/          → /business/:businessDescriptor/orders
    ├── customers/       → /business/:businessDescriptor/customers
    └── inventory/       → /business/:businessDescriptor/inventory
```

### Route Patterns

**Index Route:**

```typescript
// src/routes/index.tsx
export const Route = createFileRoute("/")({
  component: HomePage,
});
```

**Dynamic Route:**

```typescript
// src/routes/_app/customers/$id.tsx
export const Route = createFileRoute("/_app/customers/$id")({
  loader: async ({ params }) => {
    return await getCustomer(params.id);
  },
  component: CustomerDetail,
});
```

**Layout Group:**

```typescript
// src/routes/_app.tsx (app layout with auth requirement)
export const Route = createFileRoute("/_app")({
  beforeLoad: requireAuth,
  component: AppLayout,
});
```

### Navigation

**Type-safe navigation:**

```typescript
import { Link, useNavigate } from "@tanstack/react-router";

// Link
<Link to="/customers/$id" params={{ id: "123" }}>
  View Customer
</Link>;

// Programmatic
const navigate = useNavigate();
await navigate({ to: "/customers" });
```

---

## 8. Internationalization (i18next)

### Configuration (Source of Truth)

**Location:** `src/i18n/init.ts`

Kyora initializes i18next at module load time and:

- Detects language from cookie `kyora_language`, then browser language, then falls back to `en`
- Sets `document.documentElement.lang` + `document.documentElement.dir` before React renders
- Updates `lang/dir` on `languageChanged`

Key settings used by the app today:

- `defaultNS: "translation"`
- `fallbackLng: "en"`
- Explicit namespaces list (`ns`) and explicit `resources` imports

### Translation Files / Namespaces

Translations live in `src/i18n/<locale>/<namespace>.json` and are registered in `src/i18n/init.ts`.

Current namespaces:

- `translation` (default)
- `common`
- `errors`
- `onboarding`
- `upload`
- `analytics`
- `inventory`
- `orders`

### Key Conventions (Match Existing Code)

There are two patterns in use; follow the pattern of the namespace you are editing:

1. `translation` namespace: dotted keys backed by nested objects

- Example key usage: `t("dashboard.title")`, `t("common.save")`
- Example file shape: `src/i18n/en/translation.json` contains `{ "dashboard": { "title": "..." } }`

2. Feature namespaces (e.g. `inventory`, `orders`): flat, snake_case keys

- Example key usage:
  - `const { t } = useTranslation("inventory"); t("product_deleted")`

### Canonical Usage (Enforced)

- UI components must bind a translator to a namespace and call keys without any `ns:` prefix:
  - `const { t } = useTranslation("orders"); t("order_number")`
  - `const { t: tErrors } = useTranslation("errors"); tErrors("generic.unexpected")`
- Dotted keys are allowed/expected when the JSON uses nested objects (e.g. `errors.route.retry`, `common.date.selectDate`).

### Forbidden Patterns

- Do not use `t("ns:key")` anywhere (e.g. `t("orders:order_number")`).
- Do not use multi-colon strings like `t("errors:route:retry")`.
- Avoid `t(key, { ns: "..." })` in UI components; reserve it for shared utilities where the namespace is truly dynamic.

### Adding / Updating Translations (Failure-Proof Checklist)

- Add the key in **both** `ar` and `en` JSON files for the same namespace
- If you add a **new namespace file**:
  - Create `src/i18n/ar/<namespace>.json` and `src/i18n/en/<namespace>.json`
  - Import both in `src/i18n/init.ts`
  - Register them under `resources.{ar,en}`
  - Add the namespace to the `ns` array

### RTL/LTR

- Do not hardcode layout direction.
- Use `i18n.language.toLowerCase().startsWith("ar")` only when you need conditional logic.

---

## 9. Error Handling

### RFC 7807 Problem Details

**Backend Contract:** All API errors return RFC 7807 format:

```json
{
  "type": "https://api.kyora.app/errors/validation-error",
  "title": "Validation Error",
  "status": 422,
  "detail": "The request contains invalid fields.",
  "instance": "/v1/customers",
  "errors": {
    "email": "Email is already taken",
    "phone": "Invalid phone number"
  }
}
```

**Parser:** `src/lib/errorParser.ts`

```typescript
export function parseError(error: unknown): ProblemDetails {
  // Extract RFC 7807 structure
  // Handle network errors, timeouts, etc.
}
```

**Translation:** `src/lib/translateError.ts`

```typescript
export function translateError(error: unknown, t: TFunction): string {
  const problem = parseError(error);
  return t(`errors.${problem.type}`, { defaultValue: problem.detail });
}
```

**Form Integration:** See forms.instructions.md → Server Errors

---

## 10. File Upload System

**Reference:** `.github/instructions/asset_upload.instructions.md`

**Architecture:**

1. Client validates files (size, type)
2. Request upload URL: `POST /v1/businesses/{id}/assets/upload-url`
3. Upload to CDN (presigned URL)
4. Generate thumbnail (client-side: Canvas/FFmpeg)
5. Upload thumbnail
6. Return AssetReference to form

**Integration:** FileUploadField, ImageUploadField (see forms.instructions.md)

**Context Requirement:**

```typescript
import { BusinessContext } from "@/contexts/BusinessContext";

<BusinessContext.Provider value={business.descriptor}>
  <form.ImageUploadField name="logo" />
</BusinessContext.Provider>;
```

---

## Architecture Decision Records

### Why TanStack over Alternatives?

| Decision      | Chosen             | Rejected                | Rationale                          |
| ------------- | ------------------ | ----------------------- | ---------------------------------- |
| Router        | TanStack Router    | React Router            | Type-safe, file-based, better DX   |
| Data Fetching | TanStack Query     | RTK Query, SWR          | Best caching, React 19 support     |
| Forms         | TanStack Form      | React Hook Form, Formik | Most flexible, composable          |
| State         | TanStack Store     | Redux, Zustand          | Reactive, minimal API              |
| Styling       | Tailwind + daisyUI | Styled Components, MUI  | Utility-first, semantic components |

### No Redux

**Reason:** TanStack Query handles server state, TanStack Store handles client state. Redux is overkill.

### No Axios

**Reason:** Ky is modern, TypeScript-first, smaller bundle, better retry logic.

---

## Performance Patterns

1. **Route-Level Code Splitting:** TanStack Router auto-splits by route
2. **Lazy Components:** Use `React.lazy()` for large components
3. **Query Stale Time:** Configure per query (default: 0ms)
4. **Optimistic Updates:** Mutate cache before server response
5. **Skeleton Loaders:** Prevent layout shift during loading

---

## Agent Validation Checklist

Before completing architecture task:

- ☑ Using TanStack Router file-based routing
- ☑ Query keys from `queryKeys` factory
- ☑ Auth logic through `authStore`, not custom hooks
- ☑ Route guards (`requireAuth`, `requireGuest`) applied
- ☑ HTTP requests via `apiClient` from `src/api/client.ts`
- ☑ Error handling via `parseError` + `translateError`
- ☑ State management: TanStack Query (server), TanStack Store (client)
- ☑ No Redux, no Zustand, no Axios
- ☑ Translation keys, not hardcoded strings
- ☑ RTL-first CSS (see ui-implementation.instructions.md)

---

## See Also

- **Forms:** `.github/instructions/forms.instructions.md`
- **UI Patterns:** `.github/instructions/ui-implementation.instructions.md`
- **HTTP Client:** `.github/instructions/ky.instructions.md`
- **File Upload:** `.github/instructions/asset_upload.instructions.md`
- **Charts:** `.github/instructions/charts.instructions.md`
- **Development:** `.github/instructions/portal-web-development.instructions.md`

---

## Resources

- TanStack Router: https://tanstack.com/router/latest
- TanStack Query: https://tanstack.com/query/latest
- TanStack Form: https://tanstack.com/form/latest
- TanStack Store: https://tanstack.com/store/latest
- Implementation: `portal-web/src/`
