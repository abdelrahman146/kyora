---
description: Portal Web architecture - Business context routing, onboarding flow, auth flow, portal-specific navigation (portal-web only)
applyTo: "portal-web/**"
---

# Portal Web Architecture

Portal-specific architecture patterns (business management dashboard).

**Cross-refs:**

- General frontend: `../../_general/architecture.instructions.md`
- Code structure: `./code-structure.instructions.md`
- Charts: `./charts.instructions.md`

---

## 1. Overview

Portal Web is the business management dashboard for Kyora - a React SPA for Arabic-first social commerce entrepreneurs.

**Philosophy:** "Professional tools that feel effortless" - zero accounting knowledge required.

**Key Characteristics:**

- Mobile-first, RTL-native
- Arabic primary, English fallback
- Middle East social commerce focus (Instagram/WhatsApp/TikTok sellers)
- Workspace-based multi-tenancy with business sub-scope

---

## 2. Tenancy & Scoping

### Workspace vs Business

- **Workspace-scoped domains** (global to workspace):
  - Account/auth/session
  - Workspace membership & RBAC
  - Billing/subscription (Stripe)
  - Workspace-level settings

- **Business-owned domains** (must be scoped to current business):
  - Orders
  - Inventory
  - Customers
  - Analytics
  - Accounting
  - Assets
  - Storefront
  - Onboarding

### Frontend Scoping Rules

- Business-owned API calls MUST include `businessDescriptor` in URL:
  - `v1/businesses/${businessDescriptor}/orders`
  - `v1/businesses/${businessDescriptor}/inventory`
- Business-owned UI routes MUST be nested under `/business/$businessDescriptor/...`
- NEVER fetch/mutate business-owned data without explicit `businessDescriptor`

---

## 3. Route Structure

```
routes/
├── __root.tsx              → Root layout
├── index.tsx               → / (homepage/redirect)
├── auth/                   → Auth routes (redirectIfAuthenticated guard)
│   ├── login.tsx           → /auth/login
│   ├── forgot-password.tsx → /auth/forgot-password
│   ├── reset-password.tsx  → /auth/reset-password
│   └── oauth/              → OAuth callbacks
├── onboarding/             → Onboarding flow
│   ├── index.tsx           → /onboarding (session required)
│   ├── plan.tsx            → /onboarding/plan
│   ├── business-info.tsx   → /onboarding/business-info
│   └── ...                 → Other steps
└── business/               → Business-scoped routes
    └── $businessDescriptor/
        ├── index.tsx       → /business/:businessDescriptor (dashboard)
        ├── orders/         → Orders management
        ├── customers/      → Customers management
        ├── inventory/      → Inventory management
        ├── accounting/     → Accounting management
        └── reports/        → Financial reports
```

---

## 4. Authentication Flow

### Session Restoration

```tsx
// src/stores/authStore.ts
export async function initializeAuth() {
  if (isInitialized) return;

  authStore.setState({ isLoading: true });

  try {
    const user = await restoreSession(); // GET /v1/auth/me
    authStore.setState({ user, isAuthenticated: true });
  } catch {
    authStore.setState({ user: null, isAuthenticated: false });
  } finally {
    authStore.setState({ isLoading: false });
    isInitialized = true;
  }
}
```

### Route Guards

```tsx
// src/lib/routeGuards.ts
export async function requireAuth() {
  await initializeAuth();
  const { isAuthenticated } = authStore.state;
  if (!isAuthenticated) {
    throw redirect({
      to: "/auth/login",
      search: { redirect: window.location.pathname },
    });
  }
}

export async function redirectIfAuthenticated() {
  await initializeAuth();
  const { isAuthenticated } = authStore.state;
  if (isAuthenticated) {
    throw redirect({ to: "/" });
  }
}
```

---

## 5. Business Context Flow

### Business Selection

1. User lands on `/` (after login)
2. `businessStore` checks for businesses
3. If single business → redirect to `/business/$descriptor/`
4. If multiple → show business switcher
5. If none → redirect to onboarding

### Business Switching

```tsx
// src/features/business-switcher/components/BusinessSwitcher.tsx
function BusinessSwitcher() {
  const businesses = useQuery(businessQueries.list());
  const navigate = useNavigate();

  return (
    <Select
      value={currentDescriptor}
      onChange={(descriptor) => {
        navigate({ to: `/business/${descriptor}/` });
      }}
    >
      {businesses.map((b) => (
        <option key={b.descriptor} value={b.descriptor}>
          {b.name}
        </option>
      ))}
    </Select>
  );
}
```

---

## 6. Onboarding Flow

### Session-Based Flow

- **Session token** in URL: `/onboarding?session=<token>`
- **Stage progression**: `GET /v1/onboarding/session/{token}`
- **Stage completion**: `POST /v1/onboarding/session/{token}/complete-{stage}`
- **Finalization**: `POST /v1/onboarding/session/{token}/finalize`

### Route Pattern

```tsx
// src/routes/onboarding/index.tsx
export const Route = createFileRoute("/onboarding/")({
  validateSearch: z.object({
    session: z.string(),
    plan: z.string().optional(),
  }),
  loader: async ({ context, search }) => {
    return context.queryClient.ensureQueryData(
      onboardingQueries.session(search.session),
    );
  },
  beforeLoad: requireAuth,
  component: OnboardingRouter,
});

function OnboardingRouter() {
  const { session } = Route.useSearch();
  const { data } = useQuery(onboardingQueries.session(session));

  // Redirect based on stage
  if (data.stage === "plan_selection") {
    return <Navigate to="/onboarding/plan" search={{ session }} />;
  }
  if (data.stage === "business_info") {
    return <Navigate to="/onboarding/business-info" search={{ session }} />;
  }
  // ... more stages
}
```

---

## 7. Dashboard Layout

### App Shell

```tsx
// src/features/dashboard-layout/components/DashboardLayout.tsx
export function DashboardLayout() {
  const { businessDescriptor } = Route.useParams();
  const { data: business } = useQuery(
    businessQueries.detail(businessDescriptor),
  );

  return (
    <div className="min-h-screen flex flex-col">
      <Header business={business} />

      <div className="flex-1 flex">
        <Sidebar />

        <main className="flex-1 p-4">
          <Outlet />
        </main>
      </div>

      <MobileNav />
    </div>
  );
}
```

### Sidebar Navigation

```tsx
const navItems = [
  {
    to: "/business/$businessDescriptor/",
    icon: LayoutDashboard,
    label: "dashboard",
  },
  {
    to: "/business/$businessDescriptor/orders",
    icon: ShoppingCart,
    label: "orders",
  },
  {
    to: "/business/$businessDescriptor/customers",
    icon: Users,
    label: "customers",
  },
  {
    to: "/business/$businessDescriptor/inventory",
    icon: Package,
    label: "inventory",
  },
  {
    to: "/business/$businessDescriptor/accounting",
    icon: Calculator,
    label: "accounting",
  },
  {
    to: "/business/$businessDescriptor/reports",
    icon: BarChart3,
    label: "reports",
  },
];
```

---

## 8. Data Prefetching

### Route-Level Prefetch

```tsx
export const Route = createFileRoute("/business/$businessDescriptor/orders/")({
  loader: async ({ context, params }) => {
    return context.queryClient.ensureQueryData(
      orderQueries.list(params.businessDescriptor, { page: 1 }),
    );
  },
  component: OrdersPage,
});
```

### Manual Prefetch

```tsx
function OrdersListRow({ order }) {
  const queryClient = useQueryClient();

  return (
    <tr
      onMouseEnter={() => {
        queryClient.prefetchQuery(
          orderQueries.detail(order.businessDescriptor, order.id),
        );
      }}
    >
      {/* ... */}
    </tr>
  );
}
```

---

## 9. Feature Modules

### Feature Structure

```
features/orders/
├── components/
│   ├── OrdersListPage.tsx    → Main page
│   ├── OrderDetailsPage.tsx  → Detail page
│   ├── CreateOrderSheet.tsx  → Bottom sheet
│   └── OrderCard.tsx         → Card component
├── schema/
│   ├── createOrder.ts        → Zod schema
│   └── updateOrder.ts        → Zod schema
└── utils/
    └── orderStatus.ts        → Status helpers
```

---

## 10. RBAC & Permissions

### Permission Check

```tsx
function DeleteButton({ order }) {
  const { user } = authStore.state;
  const canDelete = user?.role === "admin";

  if (!canDelete) return null;

  return <button onClick={handleDelete}>Delete</button>;
}
```

### Plan Gates

```tsx
function AdvancedAnalytics() {
  const { data: subscription } = useQuery(billingQueries.subscription());
  const hasAdvanced =
    subscription?.plan.features.includes("advanced_analytics");

  if (!hasAdvanced) {
    return <UpgradeBanner feature="Advanced Analytics" />;
  }

  return <AdvancedDashboard />;
}
```

---

## 11. Analytics Dashboard

### Stats Cards

```tsx
import { StatCard, StatCardGroup } from "@/components";

<StatCardGroup cols={4}>
  <StatCard
    label="Total Revenue"
    value="$12,450"
    icon={<DollarSign />}
    trend="up"
    trendValue="+12.5%"
    variant="success"
  />
  <StatCard
    label="Orders"
    value="150"
    icon={<ShoppingCart />}
    trend="down"
    trendValue="-5%"
    variant="warning"
  />
</StatCardGroup>;
```

### Charts

```tsx
import { ChartCard, LineChart } from "@/components/charts";

<ChartCard title="Revenue Over Time" height={320}>
  <LineChart data={chartData} enableArea />
</ChartCard>;
```

---

## Agent Validation

Before completing portal architecture task:

- ☑ Business-owned routes under `/business/$businessDescriptor/...`
- ☑ Business-owned API calls include `businessDescriptor`
- ☑ Route guards applied (`requireAuth`, `redirectIfAuthenticated`)
- ☑ Onboarding flow uses session token in URL
- ☑ Auth state managed via `authStore`
- ☑ Business selection logic implemented
- ☑ Dashboard layout includes Header + Sidebar + Main
- ☑ RBAC checks for admin-only features
- ☑ Plan gates for premium features

---

## Resources

- Portal routing: `src/routes/`
- Dashboard layout: `src/features/dashboard-layout/`
- Auth: `src/stores/authStore.ts`, `src/lib/routeGuards.ts`
- Business: `src/stores/businessStore.ts`, `src/features/business-switcher/`
