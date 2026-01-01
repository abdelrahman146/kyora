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
- **clsx 2.1.1 + tailwind-merge 3.4.0**: Dynamic className composition
- **lucide-react 0.561.0 + react-icons 5.5.0**: Icon libraries

### Data Visualization

- **Chart.js 4.5.1**: Canvas-based charting library
- **react-chartjs-2 5.3.1**: React wrapper for Chart.js
- **chartjs-adapter-date-fns 3.0.0**: Date adapter for time-series charts
- **date-fns 4.1.0**: Date formatting and manipulation

### HTTP Client

- **Ky 1.14.2**: Modern fetch wrapper with retry logic, hooks, and TypeScript support

### Internationalization

- **i18next 25.7.3**: Internationalization framework
- **react-i18next 16.5.0**: React bindings for i18next

### UI Enhancements

- **react-hot-toast 2.6.0**: Toast notifications
- **react-day-picker 9.13.0**: Date and date range picker component
- **@dnd-kit 6.3.1**: Drag and drop functionality (sortable, utilities)
- **@ffmpeg/ffmpeg 0.12.15**: Video thumbnail generation (client-side)

### State Management

- **TanStack Store 0.8.0**: Reactive state management with persistence layer
- **Multiple stores**: authStore, businessStore, metadataStore, onboardingStore
- **No Redux or Zustand**: We use TanStack Query for server state and TanStack Store for client state

### Developer Tools

- **TanStack Router Devtools 1.132.0**: Route debugging and visualization
- **TanStack Query Devtools 5.84.2**: Query cache visualization
- **TanStack React Devtools 0.7.0**: General TanStack devtools
- **React Devtools**: Component tree inspection
- **ESLint + Prettier**: Code linting and formatting
- **Vitest 3.0.5**: Unit testing framework
- **TypeScript 5.7.2**: Static type checking

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
│   │   ├── assets.ts           # Asset/file upload API
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
│   │   ├── formErrors.ts       # Form error utilities
│   │   ├── onboarding.ts       # Onboarding utilities
│   │   ├── utils.ts            # General utilities
│   │   ├── form/               # Form system utilities
│   │   │   ├── createFormHook.tsx # Form hook factory
│   │   │   ├── fieldComponents.tsx # Field component registry
│   │   │   ├── formComponents.tsx # Form component registry
│   │   │   ├── contexts.tsx    # Form/field context providers
│   │   │   ├── types.ts        # Form type definitions
│   │   │   ├── useServerErrors.ts # Server error injection
│   │   │   ├── useFocusManagement.ts # Focus on error
│   │   │   ├── useKyoraForm.tsx # Main form hook
│   │   │   ├── components.tsx  # Form UI components
│   │   │   └── validation/     # Array field validation
│   │   ├── charts/             # Chart.js integration utilities
│   │   │   ├── chartTheme.ts   # Theme token resolver & hook
│   │   │   ├── chartPlugins.ts # Custom Chart.js plugins
│   │   │   ├── chartUtils.ts   # Data transformers & formatters
│   │   │   ├── rtlSupport.ts   # RTL configuration helpers
│   │   │   └── index.ts        # Public exports
│   │   └── upload/             # File upload utilities
│   │       ├── constants.ts    # File size limits, MIME types
│   │       ├── fileValidator.ts # File validation logic
│   │       ├── metadataExtractor.ts # Extract file metadata
│   │       ├── filePreviewManager.ts # Preview URL management
│   │       ├── thumbnailWorker.ts # Canvas thumbnail generation
│   │       ├── videoThumbnail.ts # Video frame extraction
│   │       ├── useFileUpload.ts # File upload hook
│   │       └── index.ts        # Public expormponents.tsx # Form component registry
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
│   │   ├── onboarding.ts       # Onboarding form schemas
│   │   ├── upload.ts           # File upload schemas
│   │   └── ...
│   ├── stores/                 # TanStack Store instances
│   │   ├── authStore.ts        # Authentication state store
│   │   ├── businessStore.ts    # Business selection & sidebar state
│   │   ├── metadataStore.ts    # Countries & currencies cache
│   │   ├── onboardingStore.ts  # Onboarding session state
│   │   └── authStore.test.ts   # Auth store unit tests
│   ├── types/                  # TypeScript type definitions
│   ├── data/                   # Static data (mock data, constants)
│   ├── main.tsx                # Application entry point
│   ├── router.tsx              # Router configuration
│   ├── routeTree.gen.ts        # Generated route tree (auto-generated)
│   └── styles.css              # Global styles and Tailwind imports
├── docs/                       # Documentation
│   ├── FORM_SYSTEM.md          # Form system documentation
│   ├── CHARTS.md               # Chart.js integration guide
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

### 5. State Management with TanStack Store

**Store Pattern**: Kyora uses TanStack Store for reactive client state management with optional localStorage persistence.

#### Auth Store

**Purpose**: Manages authentication state (user, isAuthenticated, isLoading)

**Location**: `src/stores/authStore.ts`

**State Shape**:

```typescript
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}
```

**Key Features**:

- No persistence - session restored from refresh token cookie
- Initializes on app mount via `initializeAuth()`
- Prevents concurrent initialization with shared promise

**Usage**:

```typescript
import { authStore, setUser, clearAuth } from "@/stores/authStore";
import { useStore } from "@tanstack/react-store";

// In component
const user = useStore(authStore, (state) => state.user);
const isAuthenticated = useStore(authStore, (state) => state.isAuthenticated);

// Actions
setUser(userData);
clearAuth();
```

#### Business Store

**Purpose**: Manages business list, selected business, and sidebar preferences

**Location**: `src/stores/businessStore.ts`

**State Shape**:

```typescript
interface BusinessState {
  businesses: Array<Business>;
  selectedBusinessDescriptor: string | null;
  sidebarCollapsed: boolean; // Desktop: icon-only mode
  sidebarOpen: boolean; // Mobile: drawer open state
}
```

**Persistence**: Only `selectedBusinessDescriptor` and `sidebarCollapsed` are persisted to localStorage (key: `kyora_business_prefs`)

**Usage**:

```typescript
import {
  businessStore,
  setBusinesses,
  selectBusiness,
  toggleSidebar,
} from "@/stores/businessStore";
import { useStore } from "@tanstack/react-store";

// In component
const businesses = useStore(businessStore, (state) => state.businesses);
const selected = useStore(
  businessStore,
  (state) => state.selectedBusinessDescriptor
);

// Actions
setBusinesses(businessesArray);
selectBusiness("my-business");
toggleSidebar(); // Desktop: collapse/expand, Mobile: open/close
```

#### Metadata Store

**Purpose**: Caches countries and currencies reference data with 24-hour TTL

**Location**: `src/stores/metadataStore.ts`

**State Shape**:

```typescript
interface MetadataState {
  countries: Array<CountryMetadata>;
  currencies: Array<CurrencyInfo>;
  lastFetched: number | null;
  status: "idle" | "loading" | "loaded" | "error";
  loadCountries: () => Promise<void>;
}
```

**Persistence**: Full state persisted to localStorage (key: `kyora_metadata`) with 24-hour TTL. Auto-clears stale data.

**Usage**:

```typescript
import { metadataStore, loadCountries } from "@/stores/metadataStore";
import { useStore } from "@tanstack/react-store";

// In component
const countries = useStore(metadataStore, (state) => state.countries);
const status = useStore(metadataStore, (state) => state.status);

// Load metadata (auto-cached for 24 hours)
await loadCountries();
```

#### Onboarding Store

**Purpose**: Tracks multi-step onboarding flow state (email verification, plan selection, business creation, payment)

**Location**: `src/stores/onboardingStore.ts`

**State Shape**:

```typescript
interface OnboardingState {
  sessionToken: string | null;
  stage: SessionStage | null; // 'plan_selected', 'identity_verified', etc.
  email: string | null;
  planId: string | null;
  planDescriptor: string | null;
  isPaidPlan: boolean;
  businessData: BusinessData | null;
  paymentCompleted: boolean;
  checkoutUrl: string | null;
}
```

**Persistence**: Full state persisted to localStorage (key: `kyora_onboarding_session`) with no TTL. Must be manually cleared after completion.

**Usage**:

```typescript
import {
  onboardingStore,
  setOnboardingEmail,
  setSelectedPlan,
  clearOnboardingSession,
} from "@/stores/onboardingStore";
import { useStore } from "@tanstack/react-store";

// In component
const stage = useStore(onboardingStore, (state) => state.stage);
const businessData = useStore(onboardingStore, (state) => state.businessData);

// Actions
setOnboardingEmail("user@example.com");
setSelectedPlan("plan_abc", "pro", true);
clearOnboardingSession(); // After successful completion
```

**Store Persistence Pattern**:

All stores use `createPersistencePlugin` from `src/lib/storePersistence.ts` which provides:

- Selective state persistence (choose what to save)
- TTL support (auto-expire cached data)
- Type-safe restore logic
- Automatic localStorage management

### 6. Internationalization

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

### 7. Chart System

**Overview**: Portal-web includes a comprehensive Chart.js integration with theme integration, RTL support, custom plugins, and reusable chart components.

**Full documentation**: See `portal-web/docs/CHARTS.md` for detailed API reference and advanced patterns.

#### Chart Components

Located in `src/components/atoms/`:

- **LineChart**: Line charts for trends over time
- **BarChart**: Bar charts for comparisons
- **PieChart**: Pie charts for proportions
- **DoughnutChart**: Doughnut charts with center labels
- **MixedChart**: Combined chart types (line + bar)
- **ChartCard**: Wrapper component with title, loading, empty states
- **ChartSkeleton**: Loading skeleton
- **ChartEmptyState**: Empty state placeholder

#### Theme Integration

**Theme Hook** (`src/lib/charts/chartTheme.ts`):

```typescript
import { useChartTheme } from "@/lib/charts";

function MyChart() {
  const theme = useChartTheme();

  const options = {
    plugins: {
      legend: {
        labels: {
          color: theme.text.primary,
          font: { family: theme.fontFamily },
        },
      },
    },
    scales: {
      x: {
        ticks: { color: theme.text.secondary },
        grid: { color: theme.grid },
      },
    },
  };

  return <Line data={data} options={options} />;
}
```

**Theme Structure**:

```typescript
interface ChartTheme {
  colors: {
    primary: string;
    secondary: string;
    success: string;
    error: string;
    warning: string;
    info: string;
  };
  text: {
    primary: string;
    secondary: string;
    tertiary: string;
  };
  background: {
    base100: string;
    base200: string;
    base300: string;
  };
  grid: string;
  fontFamily: string;
}
```

Theme automatically resolves from CSS variables (daisyUI tokens) and updates reactively when theme changes.

#### RTL Support

**Auto-Configuration** (`src/lib/charts/rtlSupport.ts`):

```typescript
import { getChartRTLConfig } from "@/lib/charts";
import { useLanguage } from "@/hooks/useLanguage";

function MyChart() {
  const { isRTL } = useLanguage();

  const options = {
    ...getChartRTLConfig(isRTL),
    // Your other options
  };

  return <Line data={data} options={options} />;
}
```

RTL configuration handles:

- Text alignment (right-aligned in RTL)
- Legend positioning (right side in RTL)
- Tooltip alignment
- Axis label positioning

#### Custom Plugins

**Available Plugins** (`src/lib/charts/chartPlugins.ts`):

```typescript
import { centerTextPlugin, gradientPlugin } from "@/lib/charts";

// Center text in doughnut/pie charts
const options = {
  plugins: {
    centerText: {
      text: "Total",
      value: "1,234",
      color: theme.text.primary,
    },
  },
};

<Doughnut data={data} options={options} plugins={[centerTextPlugin]} />

// Gradient backgrounds
<Bar data={data} options={options} plugins={[gradientPlugin]} />
```

#### Data Utilities

**Transformers** (`src/lib/charts/chartUtils.ts`):

```typescript
import {
  formatChartCurrency,
  formatChartDate,
  generateChartColors,
  aggregateChartData,
} from "@/lib/charts";

// Currency formatting
const formatted = formatChartCurrency(12500, "SAR"); // "SAR 12,500"

// Date formatting (respects locale)
const formatted = formatChartDate(new Date(), "MMM dd"); // "Jan 15" or "15 يناير"

// Generate color palette
const colors = generateChartColors(5, theme.colors.primary); // Returns 5 shades

// Aggregate data by time period
const aggregated = aggregateChartData(rawData, "day"); // Group by day
```

#### ChartCard Component

**Wrapper with Loading/Empty States**:

```typescript
import { ChartCard } from "@/components/atoms/ChartCard";

<ChartCard
  title="Sales Trend"
  subtitle="Last 30 days"
  isLoading={isLoading}
  isEmpty={data.length === 0}
  emptyMessage="No sales data available"
>
  <Line data={chartData} options={chartOptions} />
</ChartCard>;
```

### 8. File Upload System

**Overview**: Portal-web includes a comprehensive file upload system with validation, metadata extraction, preview management, and thumbnail generation for images and videos.

**Architecture**: Located in `src/lib/upload/`

#### File Validation

**Constants** (`src/lib/upload/constants.ts`):

```typescript
// MIME type validation
ALLOWED_IMAGE_TYPES = [
  "image/jpeg",
  "image/png",
  "image/webp",
  "image/gif",
  "image/heic",
  "image/heif",
];
ALLOWED_VIDEO_TYPES = ["video/mp4", "video/quicktime", "video/webm"];

// Size limits
MAX_IMAGE_SIZE = 10 * 1024 * 1024; // 10 MB
MAX_VIDEO_SIZE = 100 * 1024 * 1024; // 100 MB
```

**Validator** (`src/lib/upload/fileValidator.ts`):

```typescript
import { validateFile, FileValidationError } from "@/lib/upload";

try {
  const result = validateFile(file, {
    allowedTypes: ["image/jpeg", "image/png"],
    maxSize: 5 * 1024 * 1024, // 5 MB
  });

  if (result.isValid) {
    // File is valid
  }
} catch (error) {
  if (error instanceof FileValidationError) {
    toast.error(error.message);
  }
}
```

#### Metadata Extraction

**Extract File Metadata** (`src/lib/upload/metadataExtractor.ts`):

```typescript
import { extractImageMetadata, extractVideoMetadata } from "@/lib/upload";

// Image metadata
const metadata = await extractImageMetadata(file);
// Returns: { width, height, aspectRatio, format, size }

// Video metadata
const metadata = await extractVideoMetadata(file);
// Returns: { width, height, duration, format, size }
```

#### Preview Management

**Object URL Management** (`src/lib/upload/filePreviewManager.ts`):

```typescript
import { createFilePreview, revokeFilePreview } from "@/lib/upload";

// Create preview URL
const previewUrl = createFilePreview(file);

// Display in UI
<img src={previewUrl} alt="Preview" />;

// Cleanup when done (prevent memory leaks)
useEffect(() => {
  return () => revokeFilePreview(previewUrl);
}, [previewUrl]);
```

#### Thumbnail Generation

**Image Thumbnails** (`src/lib/upload/thumbnailWorker.ts`):

```typescript
import { generateImageThumbnail } from "@/lib/upload";

const thumbnail = await generateImageThumbnail(file, {
  maxWidth: 200,
  maxHeight: 200,
  quality: 0.8,
});

// Returns: Blob
```

Uses canvas-based resizing for fast, high-quality thumbnails.

**Video Thumbnails** (`src/lib/upload/videoThumbnail.ts`):

```typescript
import { generateVideoThumbnail } from "@/lib/upload";

const thumbnail = await generateVideoThumbnail(file, {
  timeInSeconds: 1.0, // Extract frame at 1 second
  width: 200,
  height: 200,
});

// Returns: Blob
```

Uses FFmpeg.wasm to extract video frames client-side without server upload.

#### Upload Hook

**useFileUpload Hook** (`src/lib/upload/useFileUpload.ts`):

```typescript
import { useFileUpload } from "@/lib/upload";

function UploadForm() {
  const {
    files,
    isUploading,
    progress,
    errors,
    addFiles,
    removeFile,
    uploadFiles,
    reset,
  } = useFileUpload({
    maxFiles: 10,
    maxSize: 10 * 1024 * 1024,
    allowedTypes: ["image/jpeg", "image/png"],
    generateThumbnails: true,
    onUploadComplete: (uploadedFiles) => {
      console.log("Uploaded:", uploadedFiles);
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  return (
    <div>
      <input
        type="file"
        multiple
        accept="image/*"
        onChange={(e) => addFiles(Array.from(e.target.files || []))}
      />

      {files.map((file) => (
        <div key={file.id}>
          <img src={file.preview} alt={file.name} />
          <button onClick={() => removeFile(file.id)}>Remove</button>
        </div>
      ))}

      <button onClick={uploadFiles} disabled={isUploading}>
        {isUploading ? `Uploading... ${progress}%` : "Upload"}
      </button>
    </div>
  );
}
```

**Hook Features**:

- Automatic file validation
- Preview URL generation and cleanup
- Thumbnail generation (images + videos)
- Upload progress tracking
- Error handling with translated messages
- File management (add, remove, reset)

#### Integration with Backend

See `.github/instructions/asset_upload.instructions.md` for backend contract and upload flow.

**Upload Flow**:

1. **Generate Upload URLs**: `POST /v1/businesses/{businessDescriptor}/assets/uploads`
2. **Upload File Bytes**: Direct upload to storage using pre-signed URL
3. **Complete Multipart** (S3 only): `POST /v1/businesses/{businessDescriptor}/assets/uploads/{assetId}/complete`
4. **Save AssetReference**: Store returned `{ url, assetId, thumbnailUrl }` in business/product/variant

**AssetReference Type**:

```typescript
interface AssetReference {
  url: string; // CDN/public URL
  originalUrl?: string; // Storage URL (fallback)
  thumbnailUrl?: string; // Thumbnail CDN URL
  thumbnailOriginalUrl?: string; // Thumbnail storage URL
  assetId: string; // Backend asset ID (required for GC)
  metadata?: {
    altText?: string;
    caption?: string;
    width?: number;
    height?: number;
  };
}
```

### 9. Error Handling

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

### 10. Styling Guidelines

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

### 11. Component Architecture (Atomic Design)

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

### 12. State Management Strategy

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

### 13. API Client Patterns

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
