---
description: Portal Web Architecture and Guidelines
applyTo: "portal-web/**"
---

# Portal Web Architecture

## Overview

**Portal Web** is the business management dashboard for Kyora. It's a modern React SPA built with TanStack Router, TanStack Query, and TanStack Form, designed specifically for Arabic-first social commerce entrepreneurs. The application is mobile-first, RTL-native, and prioritizes simplicity over feature bloat.

**Key Philosophy**: "Professional tools that feel effortless." Every UI decision removes friction for non-technical business owners.

## Tech Stack

### Core Framework

- **React 19.2.0**: Latest React with concurrent features
- **TypeScript 5.7.2**: Strict mode enabled, full type safety
- **Vite 7.1.7**: Lightning-fast dev server and build tool

### Routing & Data Fetching

- **TanStack Router 1.132.0**: File-based routing with type-safe navigation and route-level code splitting
- **TanStack Query 5.66.5**: Server state management with automatic caching, refetching, and optimistic updates
- **TanStack Router SSR Query**: Integration for prefetching data during route transitions

### Form Management

- **TanStack Form 1.27.7**: Powerful form state management
- **TanStack Store**: Reactive state management for forms
- **Zod 4.2.1**: Schema validation and type inference
- **Custom `useKyoraForm` Hook**: Zero-boilerplate composition layer over TanStack Form

### Styling & UI

- **Tailwind CSS 4.0.6**: CSS-first (no JS config), utility-first styling
- **daisyUI 5.5.14**: Semantic component classes built on Tailwind
- **clsx + tailwind-merge**: Dynamic className composition
- **lucide-react + react-icons**: Icon libraries

### HTTP Client

- **Ky 1.14.2**: Modern fetch wrapper with retry logic, hooks, and TypeScript support

### Internationalization

- **i18next 25.7.3**: Internationalization framework
- **react-i18next 16.5.0**: React bindings for i18next

### State Management

- **TanStack Store 0.8.0**: Reactive state management (used for auth store)
- **No Redux or Zustand**: We use TanStack Query for server state and TanStack Store for client state

### Developer Tools

- **TanStack Router Devtools**: Route debugging
- **TanStack Query Devtools**: Query cache visualization
- **React Devtools**: Component tree inspection
- **ESLint + Prettier**: Code linting and formatting

## Project Structure

```
portal-web/
├── src/
│   ├── api/                    # HTTP client and API modules
│   │   ├── client.ts           # Ky HTTP client with auth interceptors
│   │   ├── auth.ts             # Authentication API calls
│   │   ├── user.ts             # User management API
│   │   ├── business.ts         # Business/workspace API
│   │   ├── customer.ts         # Customer API
│   │   ├── onboarding.ts       # Onboarding API
│   │   ├── metadata.ts         # Metadata API (countries, currencies)
│   │   ├── address.ts          # Address API
│   │   └── types/              # API response/request types
│   ├── components/             # UI components (Atomic Design)
│   │   ├── atoms/              # Basic building blocks
│   │   ├── molecules/          # Composed components
│   │   ├── organisms/          # Complex UI sections
│   │   ├── templates/          # Page layouts
│   │   ├── icons/              # Custom icon components
│   │   └── index.ts            # Public exports
│   ├── hooks/                  # Custom React hooks
│   │   ├── useAuth.ts          # Authentication hook
│   │   ├── useLanguage.ts      # i18n language switching
│   │   ├── useMediaQuery.ts    # Responsive breakpoints
│   │   └── index.ts            # Public exports
│   ├── i18n/                   # Internationalization
│   │   ├── init.ts             # i18next configuration
│   │   ├── ar/                 # Arabic translations
│   │   │   ├── common.json
│   │   │   ├── auth.json
│   │   │   ├── errors.json
│   │   │   ├── validation.json
│   │   │   └── ...
│   │   └── en/                 # English translations (fallback)
│   ├── integrations/           # Third-party integrations
│   │   └── tanstack-query/     # Query client config
│   ├── lib/                    # Utility functions and helpers
│   │   ├── auth.ts             # Authentication logic
│   │   ├── cookies.ts          # Cookie management
│   │   ├── errorParser.ts      # RFC 7807 problem details parser
│   │   ├── formatCurrency.ts   # Currency formatting
│   │   ├── phone.ts            # Phone number utilities
│   │   ├── queryInvalidation.ts # Query cache invalidation helpers
│   │   ├── queryKeys.ts        # Centralized query key factory
│   │   ├── routeGuards.ts      # Route-level auth guards
│   │   ├── sessionStorage.ts   # Session storage helpers
│   │   ├── storePersistence.ts # Store persistence helpers
│   │   ├── toast.ts            # Toast notification helpers
│   │   ├── translateError.ts   # Error translation
│   │   ├── translateValidationError.ts # Zod error translation
│   │   ├── utils.ts            # General utilities
│   │   └── form/               # Form system utilities
│   │       ├── createFormHook.tsx # Form hook factory
│   │       ├── fieldComponents.tsx # Field component registry
│   │       ├── formComponents.tsx # Form component registry
│   │       └── index.ts
│   ├── routes/                 # File-based routes (TanStack Router)
│   │   ├── __root.tsx          # Root layout
│   │   ├── index.tsx           # Homepage
│   │   ├── _auth/              # Auth layout group
│   │   ├── _app/               # Authenticated app layout group
│   │   └── ...
│   ├── schemas/                # Zod validation schemas
│   │   ├── auth.ts             # Auth form schemas
│   │   ├── business.ts         # Business form schemas
│   │   └── ...
│   ├── stores/                 # TanStack Store instances
│   │   └── authStore.ts        # Authentication state store
│   ├── types/                  # TypeScript type definitions
│   ├── data/                   # Static data (mock data, constants)
│   ├── main.tsx                # Application entry point
│   ├── router.tsx              # Router configuration
│   ├── routeTree.gen.ts        # Generated route tree (auto-generated)
│   └── styles.css              # Global styles and Tailwind imports
├── docs/                       # Documentation
│   ├── FORM_SYSTEM.md          # Form system documentation
│   └── form-features.plan.md   # Form feature planning
├── public/                     # Static assets
├── index.html                  # HTML entry point
├── package.json                # Dependencies and scripts
├── tsconfig.json               # TypeScript configuration
├── vite.config.ts              # Vite configuration
├── eslint.config.js            # ESLint configuration
├── prettier.config.js          # Prettier configuration
└── README.md                   # Project README
```

## Core Patterns

### 1. Authentication Flow

**Token Management**:

- **Access Token**: Short-lived JWT stored in memory (`accessToken` variable in `client.ts`)
- **Refresh Token**: Long-lived opaque token stored in HTTP-only cookie (`kyora_refresh_token`)
- **Token Rotation**: On token expiration, automatically refresh using `POST /v1/auth/refresh`
- **Concurrent Requests**: Single refresh promise prevents thundering herd

**Auth State**:

- Managed by `authStore` (TanStack Store) in `src/stores/authStore.ts`
- State shape: `{ user: User | null, isAuthenticated: boolean, isLoading: boolean }`
- Initialized on app mount via `initializeAuth()` which calls `restoreSession()`
- No local storage persistence - session restored from refresh token cookie

**HTTP Client with Auto-Refresh** (`src/api/client.ts`):

```typescript
import ky from "ky";

// 1. Before request: inject access token
hooks: {
  beforeRequest: [
    (request) => {
      const token = getAccessToken();
      if (token) {
        request.headers.set("Authorization", `Bearer ${token}`);
      }
    },
  ];
}

// 2. After response: auto-refresh on 401
hooks: {
  afterResponse: [
    async (request, options, response) => {
      if (response.status === 401) {
        const newToken = await refreshAccessToken();
        if (newToken) {
          request.headers.set("Authorization", `Bearer ${newToken}`);
          return ky(request); // Retry with new token
        }
        clearTokens();
        window.location.href = "/auth/login";
      }
    },
  ];
}
```

**Route Guards** (`src/lib/routeGuards.ts`):

- `requireAuth`: Redirects to login if not authenticated
- `requireGuest`: Redirects to dashboard if authenticated
- Applied in route `beforeLoad` hooks

### 2. Data Fetching with TanStack Query

**Query Keys Factory** (`src/lib/queryKeys.ts`):

```typescript
export const queryKeys = {
  auth: {
    me: () => ["auth", "me"] as const,
  },
  users: {
    all: () => ["users"] as const,
    detail: (id: string) => ["users", id] as const,
  },
  customers: {
    all: (params?: Record<string, any>) => ["customers", params] as const,
    detail: (id: string) => ["customers", id] as const,
  },
  // ... more query keys
};
```

**Usage Pattern**:

```typescript
// In component
const { data, isLoading, error } = useQuery({
  queryKey: queryKeys.customers.all({ page: 1 }),
  queryFn: () => getCustomers({ page: 1 }),
});

// Invalidation after mutation
const queryClient = useQueryClient();
await queryClient.invalidateQueries({ queryKey: queryKeys.customers.all() });
```

**Mutation Pattern**:

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

### 3. Form Management with useKyoraForm

**The Problem**: TanStack Form is powerful but requires boilerplate for every field.

**The Solution**: `useKyoraForm` composition layer that pre-binds components and eliminates manual wiring.

**Architecture**:

- `createFormHook` factory creates the `useKyoraForm` hook
- Returns form instance with:
  - `form.AppForm`: Form context provider (MUST wrap all form components)
  - `form.AppField`: Field component with pre-bound field context
  - `form.FormRoot`, `form.SubmitButton`, `form.FormError`: Form-level components
  - Field components: `field.TextField`, `field.PasswordField`, `field.SelectField`, etc.

**Critical Rule**: ALL components that use form context MUST be inside `<form.AppForm>`. If you see `Error: formContext only works within formComponent`, you have a component using form context outside `<form.AppForm>`.

**Basic Usage**:

```typescript
import { useKyoraForm } from "@/lib/form";
import { z } from "zod";

function LoginForm() {
  const { t } = useTranslation();

  const form = useKyoraForm({
    defaultValues: {
      email: "",
      password: "",
    },
    onSubmit: async ({ value }) => {
      await loginUser(value.email, value.password);
    },
  });

  return (
    <form.AppForm>
      <form.FormRoot className="space-y-4">
        <form.FormError />

        <form.AppField
          name="email"
          validators={{
            onBlur: z.string().email("invalid_email"),
          }}
        >
          {(field) => (
            <field.TextField
              type="email"
              label={t("auth.email")}
              autoComplete="email"
            />
          )}
        </form.AppField>

        <form.AppField
          name="password"
          validators={{
            onBlur: z.string().min(8, "password_too_short"),
          }}
        >
          {(field) => (
            <field.PasswordField
              label={t("auth.password")}
              autoComplete="current-password"
            />
          )}
        </form.AppField>

        <form.SubmitButton variant="primary">
          {t("auth.login")}
        </form.SubmitButton>
      </form.FormRoot>
    </form.AppForm>
  );
}
```

**Key Features**:

- ✅ **Zero Boilerplate**: Components pre-bound to field context
- ✅ **Progressive Validation**: Smart revalidation (submit → blur → change modes)
- ✅ **Auto-Translation**: Zod error keys automatically translated via i18n
- ✅ **Focus Management**: Automatic focus on first invalid field
- ✅ **Server Errors**: RFC7807 problem details integration
- ✅ **Type-Safe**: Full TypeScript inference from `defaultValues`

**Available Field Components**:

- `field.TextField` - Text input
- `field.PasswordField` - Password input with show/hide toggle
- `field.EmailField` - Email input with validation
- `field.TextareaField` - Multi-line text input
- `field.SelectField` - Dropdown select
- `field.CheckboxField` - Single checkbox
- `field.RadioGroupField` - Radio button group
- `field.SwitchField` - Toggle switch

**Server Error Handling**:

```typescript
const form = useKyoraForm({
  defaultValues: { email: "", password: "" },
  onSubmit: async ({ value }) => {
    try {
      await loginUser(value.email, value.password);
    } catch (error) {
      // RFC 7807 errors are automatically parsed and displayed
      throw error;
    }
  },
});
```

**See `portal-web/docs/FORM_SYSTEM.md` for comprehensive documentation.**

### 4. Routing with TanStack Router

**File-Based Routing**:

- Routes are defined in `src/routes/` directory
- File structure determines URL structure
- Auto-generated route tree in `src/routeTree.gen.ts`

**Route Patterns**:

```
src/routes/
├── __root.tsx               # Root layout (shared across all routes)
├── index.tsx                # / (homepage)
├── _auth/                   # Layout group (not in URL)
│   ├── login.tsx            # /login
│   ├── register.tsx         # /register
│   └── forgot-password.tsx  # /forgot-password
├── _app/                    # Authenticated layout group
│   ├── dashboard.tsx        # /dashboard
│   ├── customers/
│   │   ├── index.tsx        # /customers
│   │   └── $id.tsx          # /customers/:id (dynamic route)
│   └── settings/
│       ├── index.tsx        # /settings
│       ├── profile.tsx      # /settings/profile
│       └── business.tsx     # /settings/business
```

**Route Context** (`src/router.tsx`):

```typescript
export interface RouterContext {
  auth: {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
  };
  queryClient: QueryClient;
}
```

**Route Guards**:

```typescript
// Protected route
export const Route = createFileRoute("/_app/dashboard")({
  beforeLoad: ({ context }) => {
    requireAuth(context.auth);
  },
  component: Dashboard,
});

// Guest-only route
export const Route = createFileRoute("/_auth/login")({
  beforeLoad: ({ context }) => {
    requireGuest(context.auth);
  },
  component: Login,
});
```

**Data Prefetching**:

```typescript
export const Route = createFileRoute("/_app/customers/$id")({
  loader: ({ params, context }) => {
    return context.queryClient.ensureQueryData({
      queryKey: queryKeys.customers.detail(params.id),
      queryFn: () => getCustomer(params.id),
    });
  },
  component: CustomerDetail,
});
```

**Navigation**:

```typescript
import { Link, useNavigate } from "@tanstack/react-router";

// Declarative
<Link to="/customers/$id" params={{ id: "123" }}>
  View Customer
</Link>;

// Imperative
const navigate = useNavigate();
navigate({ to: "/customers", search: { page: 1 } });
```

### 5. Internationalization

**Setup** (`src/i18n/init.ts`):

```typescript
import i18n from "i18next";
import { initReactI18next } from "react-i18next";

i18n.use(initReactI18next).init({
  resources: {
    ar: {
      /* Arabic translations */
    },
    en: {
      /* English translations */
    },
  },
  lng: "ar", // Default language
  fallbackLng: "en",
  interpolation: {
    escapeValue: false, // React already escapes
  },
});
```

**Translation Structure**:

```
i18n/
├── ar/
│   ├── common.json          # Common UI text
│   ├── auth.json            # Authentication
│   ├── errors.json          # Error messages
│   ├── validation.json      # Validation errors
│   ├── customer.json        # Customer module
│   ├── order.json           # Order module
│   └── ...
└── en/                      # Same structure
```

**Usage**:

```typescript
import { useTranslation } from "react-i18next";

function MyComponent() {
  const { t } = useTranslation();

  return (
    <div>
      <h1>{t("common.welcome")}</h1>
      <p>{t("customer.greeting", { name: "Ahmad" })}</p>
    </div>
  );
}
```

**Language Switching** (`src/hooks/useLanguage.ts`):

```typescript
const { language, setLanguage } = useLanguage();

// Switch to Arabic
setLanguage("ar");

// Switch to English
setLanguage("en");
```

**RTL Support**:

- Language direction automatically set on `<html dir="rtl">` or `<html dir="ltr">`
- Use logical properties in CSS: `start` instead of `left`, `end` instead of `right`
- Tailwind utilities: `ms-4` (margin-start), `pe-2` (padding-end)

### 6. Error Handling

**RFC 7807 Problem Details** (`src/lib/errorParser.ts`):

```typescript
export interface ProblemDetails {
  type: string;
  title: string;
  status: number;
  detail?: string;
  instance?: string;
  traceId?: string;
  [key: string]: any;
}

export function parseProblemDetails(error: any): ProblemDetails | null;
```

**Error Translation** (`src/lib/translateError.ts`):

```typescript
import { translateError } from "@/lib/translateError";

try {
  await api.createCustomer(data);
} catch (error) {
  const message = translateError(error);
  toast.error(message); // Translated user-friendly message
}
```

**Validation Error Translation** (`src/lib/translateValidationError.ts`):

```typescript
// Zod error codes are automatically translated
const schema = z.object({
  email: z.string().email("invalid_email"), // Translates to ar.validation.invalid_email
});
```

**Toast Notifications** (`src/lib/toast.ts`):

```typescript
import { toast } from "@/lib/toast";

toast.success(t("customer.created"));
toast.error(t("errors.network_error"));
toast.info(t("common.processing"));
toast.loading(t("common.loading"));
```

### 7. Styling Guidelines

**Tailwind CSS 4 (CSS-First)**:

- No `tailwind.config.js` - configuration is in CSS via `@theme` directive
- Theme tokens defined in `src/styles.css`
- Import in components: `import '@/styles.css'`

**daisyUI Component Classes**:

- Use semantic classes: `.btn`, `.btn-primary`, `.input`, `.select`, `.card`, etc.
- Never style daisyUI components with Tailwind utilities - use daisyUI modifiers
- Example: `.btn.btn-primary.btn-lg` instead of `.btn.bg-teal-600.px-8.py-4`

**RTL-First Design**:

- **Never use**: `left`, `right`, `ml-*`, `mr-*`, `pl-*`, `pr-*`, `border-l`, `border-r`
- **Always use**: `start`, `end`, `ms-*`, `me-*`, `ps-*`, `pe-*`, `border-s`, `border-e`
- **Directional icons**: Use `ChevronStart`, `ChevronEnd` instead of `ChevronLeft`, `ChevronRight`

**Design Tokens** (from `src/styles.css`):

```css
@theme {
  /* Colors */
  --color-primary: #0d9488; /* Teal */
  --color-secondary: #eab308; /* Gold */
  --color-success: #22c55e;
  --color-error: #ef4444;
  --color-warning: #f59e0b;
  --color-info: #3b82f6;

  /* Spacing */
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 2rem;

  /* Typography */
  --font-family-base: "IBM Plex Sans Arabic", sans-serif;
  --font-size-xs: 0.75rem;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
}
```

**Component Styling Pattern**:

```typescript
import { cn } from "@/lib/utils";

function Button({ variant = "primary", className, ...props }) {
  return (
    <button
      className={cn(
        "btn", // daisyUI base class
        variant === "primary" && "btn-primary",
        variant === "secondary" && "btn-secondary",
        className // Allow overrides
      )}
      {...props}
    />
  );
}
```

### 8. Component Architecture (Atomic Design)

**Atoms** (`src/components/atoms/`):

- Basic building blocks: Button, Input, Label, Badge, Avatar, etc.
- Reusable across the entire application
- Minimal to no business logic
- Example: `Button.tsx`, `Input.tsx`, `Badge.tsx`

**Molecules** (`src/components/molecules/`):

- Compositions of atoms: FormField, SearchBar, Pagination, etc.
- Moderate complexity
- Reusable with slight variations
- Example: `FormField.tsx`, `SearchBar.tsx`, `DataTable.tsx`

**Organisms** (`src/components/organisms/`):

- Complex UI sections: Navbar, Sidebar, CustomerTable, OrderForm, etc.
- May contain business logic
- Domain-specific
- Example: `Navbar.tsx`, `CustomerTable.tsx`, `OrderForm.tsx`

**Templates** (`src/components/templates/`):

- Page layouts: AuthLayout, AppLayout, DashboardLayout, etc.
- Structural composition
- No direct data fetching
- Example: `AuthLayout.tsx`, `AppLayout.tsx`

**Component Export Pattern** (`src/components/index.ts`):

```typescript
// Re-export only public components
export { Button } from "./atoms/Button";
export { FormField } from "./molecules/FormField";
export { Navbar } from "./organisms/Navbar";
export { AppLayout } from "./templates/AppLayout";
```

### 9. State Management Strategy

**Server State**: TanStack Query

- API data, remote state
- Automatic caching, refetching, invalidation
- Query keys factory in `src/lib/queryKeys.ts`

**Client State**: TanStack Store (minimal use)

- Authentication state (`authStore`)
- Avoid local state when possible - prefer React state or URL state

**URL State**: TanStack Router

- Pagination, filters, search params
- Shareable links, browser back/forward
- Example: `/customers?page=2&search=ahmad`

**Form State**: TanStack Form

- Form values, validation, errors
- Ephemeral, isolated to form component

**React State**: useState, useReducer

- Transient UI state (modals, dropdowns, toggles)
- Component-local state

### 10. API Client Patterns

**Base Client** (`src/api/client.ts`):

```typescript
import ky from "ky";

const client = ky.create({
  prefixUrl: "http://localhost:8080",
  timeout: 30000,
  hooks: {
    beforeRequest: [
      /* Auth token injection */
    ],
    afterResponse: [
      /* Auto-refresh on 401 */
    ],
  },
});

export default client;
```

**API Module Pattern** (`src/api/customers.ts`):

```typescript
import client from "./client";
import type { Customer, CreateCustomerRequest } from "./types/customer";

export async function getCustomers(params?: {
  page?: number;
  search?: string;
}) {
  return client
    .get("v1/customers", { searchParams: params })
    .json<Customer[]>();
}

export async function getCustomer(id: string) {
  return client.get(`v1/customers/${id}`).json<Customer>();
}

export async function createCustomer(data: CreateCustomerRequest) {
  return client.post("v1/customers", { json: data }).json<Customer>();
}

export async function updateCustomer(
  id: string,
  data: Partial<CreateCustomerRequest>
) {
  return client.patch(`v1/customers/${id}`, { json: data }).json<Customer>();
}

export async function deleteCustomer(id: string) {
  return client.delete(`v1/customers/${id}`).json();
}
```

**Type Safety**:

- Response types in `src/api/types/`
- Match backend Go structs exactly
- Use `snake_case` for JSON fields (backend convention)
- Convert to `camelCase` only if necessary via mappers

## Development Workflow

### Running Locally

```bash
# Install dependencies
npm install

# Start dev server (localhost:3000)
npm run dev

# Type check
npm run type-check

# Lint
npm run lint

# Format
npm run format

# Format and fix
npm run check

# Build for production
npm run build

# Preview production build
npm run preview
```

### Code Quality Checks

**Before Committing**:

1. Run `npm run type-check` - fix TypeScript errors
2. Run `npm run lint` - fix linting errors
3. Run `npm run format` - format code
4. Ensure no console errors in browser
5. Test critical user flows

### Adding New Features

**1. Add Route** (`src/routes/`):

```bash
# Create route file
touch src/routes/_app/new-feature.tsx
```

**2. Add API Module** (`src/api/`):

```typescript
// src/api/newFeature.ts
import client from "./client";

export async function getNewFeature() {
  return client.get("v1/new-feature").json();
}
```

**3. Add Query Keys** (`src/lib/queryKeys.ts`):

```typescript
export const queryKeys = {
  // ...
  newFeature: {
    all: () => ["newFeature"] as const,
    detail: (id: string) => ["newFeature", id] as const,
  },
};
```

**4. Add Translations** (`src/i18n/ar/`, `src/i18n/en/`):

```json
// src/i18n/ar/newFeature.json
{
  "title": "الميزة الجديدة",
  "description": "وصف الميزة"
}
```

**5. Add Components** (`src/components/`):

```typescript
// src/components/organisms/NewFeatureTable.tsx
export function NewFeatureTable() {
  const { t } = useTranslation();
  const { data } = useQuery({
    queryKey: queryKeys.newFeature.all(),
    queryFn: getNewFeature,
  });

  return <table>{/* ... */}</table>;
}
```

## Best Practices

### General Code Quality

1. **TypeScript**: Use strict mode, avoid `any`, prefer interfaces over types
2. **Immutability**: Never mutate objects/arrays - use spread operators
3. **Pure Functions**: Avoid side effects in utility functions
4. **Single Responsibility**: Each function/component does one thing well
5. **DRY**: Extract repeated logic into utilities or hooks
6. **YAGNI**: Don't add features until needed
7. **Comments**: Explain "why", not "what" (code should be self-documenting)

### React Patterns

1. **Hooks**: Prefer hooks over class components
2. **Custom Hooks**: Extract reusable logic into hooks
3. **Memoization**: Use `useMemo` and `useCallback` for expensive computations
4. **Lazy Loading**: Use `React.lazy()` for route-level code splitting
5. **Error Boundaries**: Wrap risky components in error boundaries
6. **Accessibility**: Always add ARIA labels, keyboard navigation, focus management

### Form Best Practices

1. **Always wrap in `<form.AppForm>`**: Form context is required for all form components
2. **Use `form.AppField`**: Never use `<form.Field>` directly - always use `<form.AppField>`
3. **Field Component Pattern**: Always use `{(field) => <field.ComponentName />}` pattern
4. **Validation**: Use Zod schemas with i18n-friendly error keys
5. **Auto-Complete**: Always add `autoComplete` attribute to inputs
6. **Progressive Validation**: Let TanStack Form handle validation timing
7. **Server Errors**: Let `useKyoraForm` handle RFC 7807 errors automatically

### Performance

1. **Lazy Load Routes**: Use `React.lazy()` for route components
2. **Query Stale Time**: Set reasonable `staleTime` to avoid excessive refetches
3. **Query Caching**: Use query keys factory for consistent cache keys
4. **Debounce Search**: Debounce search inputs to reduce API calls
5. **Virtual Scrolling**: Use virtual scrolling for large lists (react-window)
6. **Image Optimization**: Use WebP format, lazy load images

### Security

1. **XSS Prevention**: React escapes by default - never use `dangerouslySetInnerHTML`
2. **CSRF Protection**: Backend handles CSRF tokens
3. **Content Security Policy**: Configured in `index.html`
4. **HTTP-Only Cookies**: Refresh tokens stored in HTTP-only cookies
5. **Token Expiry**: Access tokens are short-lived, auto-refresh on expiry
6. **Sensitive Data**: Never log sensitive data (passwords, tokens)

### Accessibility (a11y)

1. **Semantic HTML**: Use `<button>`, `<nav>`, `<main>`, `<article>`, etc.
2. **ARIA Labels**: Add `aria-label`, `aria-describedby` where needed
3. **Keyboard Navigation**: Ensure all interactive elements are keyboard-accessible
4. **Focus Management**: Focus first invalid field on form submission
5. **Color Contrast**: Ensure WCAG AA compliance (4.5:1 ratio)
6. **Screen Reader Testing**: Test with NVDA/JAWS/VoiceOver

## Testing Strategy

### Unit Tests (Vitest)

- Test utility functions in `src/lib/`
- Test custom hooks in `src/hooks/`
- Test form validation schemas in `src/schemas/`

### Integration Tests

- Test API client error handling
- Test form submission flows
- Test authentication flows

### E2E Tests (Future)

- Test critical user journeys
- Test authentication flows
- Test order creation flows

## Deployment

### Production Build

```bash
npm run build
```

Output: `dist/` directory with optimized assets

### Environment Variables

Create `.env.production`:

```
VITE_API_BASE_URL=https://api.kyora.app
```

### Hosting

- **Recommended**: Vercel, Netlify, Cloudflare Pages
- **SPA Fallback**: Ensure all routes serve `index.html`
- **CDN**: Enable CDN caching for static assets
- **Compression**: Enable gzip/brotli compression

## Troubleshooting

### Common Issues

**1. Form context error**:

```
Error: formContext only works when within a formComponent passed to createFormHook
```

**Fix**: Wrap form components in `<form.AppForm>`, use `<form.AppField>` instead of `<form.Field>`

**2. Query not updating after mutation**:
**Fix**: Invalidate query cache after mutation using `queryClient.invalidateQueries()`

**3. RTL layout broken**:
**Fix**: Use logical properties (`ms-*`, `me-*`) instead of directional (`ml-*`, `mr-*`)

**4. Token refresh loop**:
**Fix**: Ensure refresh token endpoint doesn't trigger another refresh on 401

**5. Translation missing**:
**Fix**: Add translation key to both `i18n/ar/` and `i18n/en/` files

## Migration Notes

### From Old Portal Web

The new portal-web replaces the legacy React Router v6 implementation with:

1. **TanStack Router** instead of React Router v6
2. **TanStack Form** with `useKyoraForm` instead of React Hook Form
3. **Tailwind CSS 4** (CSS-first) instead of Tailwind CSS 3
4. **File-based routing** instead of route configuration
5. **TanStack Store** for auth state instead of Context API

**Breaking Changes**:

- All routes must be migrated to file-based routing
- All forms must be migrated to TanStack Form with `useKyoraForm`
- All route guards must use `beforeLoad` instead of route wrappers
- All translations must use new i18n structure

## Resources

- [TanStack Router Docs](https://tanstack.com/router)
- [TanStack Query Docs](https://tanstack.com/query)
- [TanStack Form Docs](https://tanstack.com/form)
- [Tailwind CSS 4 Docs](https://tailwindcss.com/docs)
- [daisyUI Docs](https://daisyui.com/)
- [Ky HTTP Client Docs](https://github.com/sindresorhus/ky)
- [Zod Docs](https://zod.dev/)
- [i18next Docs](https://www.i18next.com/)

## Contributing

When contributing to portal-web:

1. **Read this document** thoroughly
2. **Follow existing patterns** - consistency is critical
3. **Write clear commit messages** - use conventional commits
4. **Test thoroughly** - verify all affected features work
5. **Update documentation** - keep this file up-to-date with changes
6. **Ask questions** - better to ask than to assume

---

**Remember**: The goal of portal-web is to make complex business management feel effortless for non-technical users. Every feature should remove friction, not add complexity.
