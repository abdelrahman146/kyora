# portal-web/AGENTS.md

## Scope

React SPA dashboard — business management portal for Kyora merchants.

**Parent AGENTS.md**: [../AGENTS.md](../AGENTS.md) (read first for project context and global boundaries)

## Tech Stack

- **Language**: TypeScript 5.7+
- **Framework**: React 19
- **Build**: Vite 7
- **Routing**: TanStack Router (file-based, type-safe)
- **Server State**: TanStack Query
- **Client State**: TanStack Store
- **Forms**: TanStack Form + Zod + `useKyoraForm`
- **HTTP Client**: Ky
- **Styling**: Tailwind CSS 4 + daisyUI 5
- **i18n**: i18next (Arabic-first, RTL-native)
- **Charts**: Chart.js + react-chartjs-2
- **Icons**: lucide-react

## Setup Commands

```bash
# Install dependencies
make portal.install
# or: cd portal-web && npm ci

# Run dev server (default port 3000)
make dev.portal
# or: cd portal-web && npm run dev

# Run dev server on custom port
PORTAL_PORT=3001 make dev.portal

# Lint + type check
make portal.check
# or: cd portal-web && npm run lint && npm run typecheck

# Build for production
make portal.build
# or: cd portal-web && npm run build

# Preview production build
make portal.preview
```

## Structure

```
portal-web/src/
├── api/                    # HTTP client + API modules
│   ├── client.ts           # Ky instance with auth interceptors
│   ├── types/              # API response types
│   └── *.ts                # Domain-specific API modules
├── components/             # Shared UI (Atomic Design)
│   ├── atoms/              # Buttons, badges, inputs
│   ├── molecules/          # SearchInput, BottomSheet
│   ├── organisms/          # App chrome, complex composites
│   ├── templates/          # Page layouts
│   ├── charts/             # Chart.js wrappers
│   └── form/               # Generic form controls
├── features/               # Feature modules
│   ├── auth/               # Login, registration
│   ├── onboarding/         # Business setup
│   ├── dashboard/          # Main dashboard
│   ├── orders/             # Order management
│   ├── inventory/          # Products, variants
│   ├── customers/          # Customer management
│   ├── accounting/         # Expenses, investments
│   └── reports/            # Analytics, reports
├── hooks/                  # Custom hooks
├── i18n/                   # Translations
│   ├── ar/                 # Arabic (primary)
│   └── en/                 # English (fallback)
├── lib/                    # Utilities
│   ├── form/               # Form system (useKyoraForm)
│   ├── upload/             # File upload utilities
│   ├── charts/             # Chart.js utilities
│   └── *.ts                # Other utilities
├── routes/                 # File-based routes
├── schemas/                # Zod validation schemas
├── stores/                 # TanStack Store instances
└── types/                  # TypeScript types
```

## Code Style

### Component Pattern

```tsx
// ✅ Good: RTL-safe, uses design tokens, handles states
function OrderCard({ order }: { order: Order }) {
  const { t } = useTranslation();
  
  return (
    <div className="card card-compact bg-base-100 shadow">
      <div className="card-body">
        <div className="flex justify-between items-center">
          <h3 className="card-title text-base">
            {t('orders:order_number', { number: order.number })}
          </h3>
          <Badge variant={statusVariant[order.status]}>
            {t(`orders:status.${order.status}`)}
          </Badge>
        </div>
        <p className="text-base-content/70">
          {formatCurrency(order.total, order.currency)}
        </p>
      </div>
    </div>
  );
}
```

### Form Pattern

```tsx
// ✅ Good: useKyoraForm with proper field pattern
function CreateCustomerForm() {
  const { t } = useTranslation();
  const form = useKyoraForm({
    defaultValues: { name: '', phone: '' },
    onSubmit: async ({ value }) => {
      await createCustomer(value);
    },
  });

  return (
    <form.AppForm>
      <form.FormRoot className="space-y-4">
        <form.AppField
          name="name"
          validators={{ onBlur: z.string().min(1, 'validation.required') }}
        >
          {(field) => (
            <field.TextField label={t('customers:name')} required />
          )}
        </form.AppField>
        
        <form.AppField
          name="phone"
          validators={{ onBlur: z.string().min(1, 'validation.required') }}
        >
          {(field) => (
            <field.TextField
              type="tel"
              label={t('customers:phone')}
              dir="ltr"
              required
            />
          )}
        </form.AppField>
        
        <form.SubmitButton variant="primary">
          {t('common:save')}
        </form.SubmitButton>
      </form.FormRoot>
    </form.AppForm>
  );
}
```

### Query Pattern

```tsx
// ✅ Good: Query with proper keys and error handling
function useOrders(businessDescriptor: string) {
  return useQuery({
    queryKey: queryKeys.orders.list(businessDescriptor),
    queryFn: () => ordersApi.list(businessDescriptor),
    enabled: !!businessDescriptor,
  });
}
```

### Mutation Pattern

```tsx
// ✅ Good: Mutation with invalidation and toast
function useCreateOrder(businessDescriptor: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateOrderInput) => 
      ordersApi.create(businessDescriptor, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: queryKeys.orders.list(businessDescriptor) 
      });
      toast.success(t('orders:created'));
    },
  });
}
```

## Boundaries (Portal-Specific)

### ✅ Always do

- Use `useKyoraForm` for all forms (never raw TanStack Form)
- Use `<form.AppField>` pattern for all form fields
- Use translation keys for ALL user-facing text
- Use `PriceField` for money inputs (never TextField with type="number")
- Handle loading/empty/error states in every component
- Use RTL-safe classes (`start/end`, `ms/me`, not `left/right`, `ml/mr`)
- Use `dir="ltr"` for LTR-only content (phone numbers, codes)
- Use `queryKeys` factory for all query keys
- Invalidate queries after mutations

### ⚠️ Ask first

- New shared components in `components/`
- New dependencies
- Changes to auth flow
- New feature modules
- Changes to form system (`lib/form/`)

### 🚫 Never do

- Use form components directly (must use `<form.AppField>` pattern)
- Use `TextField` for money (use `PriceField`)
- Hardcode strings (use `t()` function)
- Use `left/right` or `ml/mr` classes (use `start/end`, `ms/me`)
- Skip loading/empty/error states
- Make API calls without going through `api/*.ts` modules
- Use raw `ky` or `fetch` (use `apiClient` from `api/client.ts`)

## i18n Conventions

- **Primary**: Arabic (`ar/`)
- **Fallback**: English (`en/`)
- **Namespace files**: `common.json`, `errors.json`, `orders.json`, etc.
- **Validation keys**: Must use `validation.*` prefix
- **Key format**: `namespace:key` or `namespace:nested.key`

```tsx
// ✅ Correct usage
t('orders:status.pending')           // Namespaced key
t('common:save')                     // Common namespace
t('validation.required')             // Validation (errors.json)

// ❌ Wrong
t('Save')                            // Hardcoded
t('Order Status')                    // Hardcoded
```

## SSOT Entry Points

- [.github/instructions/frontend/projects/portal-web/architecture.instructions.md](../.github/instructions/frontend/projects/portal-web/architecture.instructions.md) — Architecture
- [.github/instructions/frontend/projects/portal-web/code-structure.instructions.md](../.github/instructions/frontend/projects/portal-web/code-structure.instructions.md) — Code structure
- [.github/instructions/frontend/_general/ui-patterns.instructions.md](../.github/instructions/frontend/_general/ui-patterns.instructions.md) — UI/RTL
- [.github/instructions/frontend/_general/forms.instructions.md](../.github/instructions/frontend/_general/forms.instructions.md) — Forms
- [.github/instructions/frontend/_general/i18n.instructions.md](../.github/instructions/frontend/_general/i18n.instructions.md) — i18n
- [.github/instructions/frontend/_general/http-client.instructions.md](../.github/instructions/frontend/_general/http-client.instructions.md) — HTTP/TanStack Query
- [.github/instructions/kyora/design-system.instructions.md](../.github/instructions/kyora/design-system.instructions.md) — Design system

## Agent Routing Hints

**Web Lead** (`@Web Lead`): Architecture, component design, state patterns
**Web Implementer** (`@Web Implementer`): UI implementation, i18n, API integration
**Design/UX Lead** (`@Design/UX Lead`): UX specs, states/variants, RTL notes
**i18n/Localization Lead** (`@i18n/Localization Lead`): Translation keys, Arabic copy
