# AGENTS.md

## Project Overview

Kyora is a DM-first commerce platform for micro-businesses in the Middle East. It helps sellers feel "handled" by keeping orders, stock, and money organized in plain language. This monorepo currently includes `backend/` (Go API) and `portal-web/` (React dashboard), with more apps/services to be added.

For the operating model governing agent collaboration, lanes, routing, and handoffs, see [KYORA_AGENT_OS.md](./KYORA_AGENT_OS.md).

## Tech Stack

- **Backend**: Go 1.22+ (Gin, GORM/Postgres), PostgreSQL, Memcached, Stripe, Resend
- **Frontend**: React + Vite + TanStack (Router/Query/Form/Store), i18n (Arabic/RTL-first), Tailwind + daisyUI
- **Tests**: Go tests under `backend/` (E2E in `backend/internal/tests/e2e/`)

## Setup Commands

```bash
# Tooling sanity check
make doctor

# Start infra (Postgres + Memcached + Stripe mock)
make infra.up

# Run both API + portal in dev mode
make dev

# Run API only
make dev.server

# Run portal only
make dev.portal

# Backend tests (all)
make test

# Backend tests (quick - unit only)
make test.quick

# Backend E2E tests
make test.e2e

# Backend OpenAPI regeneration + verification
make openapi
make openapi.check

# Portal dependency install
make portal.install

# Portal lint + type check
make portal.check

# Portal build
make portal.build

# See all available targets
make help
```

## Project Structure

```
backend/
  cmd/                    # CLI commands (server, seed, sync_plans)
  docs/                   # Swagger/OpenAPI
  internal/
    domain/               # Business logic (account, accounting, analytics, asset, billing, business, customer, inventory, metadata, onboarding, order, storefront)
    platform/             # Infrastructure (auth, blob, bus, cache, config, database, email, logger, middleware, request, response, types, utils)
    server/               # HTTP server + routes
    tests/e2e/            # End-to-end tests

portal-web/
  src/
    api/                  # API client modules
    components/           # Shared UI components (atoms, molecules, organisms, templates)
    features/             # Feature modules (accounting, app-shell, auth, customers, dashboard, inventory, onboarding, orders, reports)
    hooks/                # Custom hooks
    i18n/                 # Translations (ar/, en/)
    lib/                  # Utilities (form, upload, charts)
    routes/               # Page routes
    schemas/              # Zod schemas
    stores/               # TanStack Store instances
    types/                # TypeScript types

.github/
  agents/                 # Custom Copilot agents
  prompts/                # Reusable prompts
  skills/                 # Multi-step workflow skills
  instructions/           # Always-on SSOT coding standards
```

## Code Style

### Backend (Go)

```go
// ✅ Good: Tenant-scoped query with validation
func (s *OrderService) GetOrderByID(ctx context.Context, businessID, orderID string) (*Order, error) {
    if businessID == "" || orderID == "" {
        return nil, domain.NewValidationError("business_id and order_id required")
    }
    
    var order Order
    err := s.db.WithContext(ctx).
        Where("business_id = ? AND id = ?", businessID, orderID).
        First(&order).Error
    if err != nil {
        return nil, domain.WrapError(err, "failed to get order")
    }
    return &order, nil
}
```

### Frontend (React/TypeScript)

```tsx
// ✅ Good: RTL-safe, uses design tokens, handles all states
function OrderStatusBadge({ status }: { status: OrderStatus }) {
  const { t } = useTranslation();
  
  return (
    <Badge variant={statusVariantMap[status]} className="gap-1">
      <StatusIcon status={status} className="size-3" />
      {t(`orders:status.${status}`)}
    </Badge>
  );
}
```

## Boundaries

### ✅ Always do

- Run `make test.quick` or `make portal.check` before committing changes
- Follow existing module boundaries (domain logic in `domain/**`, infra in `platform/**`)
- Use existing shared utils (`backend/internal/platform/utils/`, `portal-web/src/lib/`)
- Validate inputs and return domain errors (Kyora Problem/RFC7807 patterns)
- Use translation keys for all user-facing strings
- Handle loading/empty/error states in UI components
- Write tests for new backend functionality
- Verify RTL layout when touching portal UI

### ⚠️ Ask first (PO gate required)

- DB schema changes or migrations
- New dependencies (any app)
- Auth/RBAC/tenant boundary changes
- Breaking API contract changes
- Major UX redesign or new UI primitives
- Payments/billing/Stripe integration changes

### 🚫 Never do

- Commit secrets, tokens, or credentials
- Cross-tenant reads/writes (violate workspace > business isolation)
- Hardcode UI strings (use i18n keys)
- Add TODO/FIXME placeholders (ship complete implementations)
- Use accounting jargon in UI text (prefer "Profit", "Cash in hand", "Money in/out")
- Assume left/right positioning (use RTL-safe start/end)
- Create/expand docs unless explicitly requested ("No surprise docs")
- Modify `node_modules/`, `vendor/`, or generated files

## SSOT Entry Points

Core artifact guidance:
- [.github/copilot-instructions.md](.github/copilot-instructions.md) — Repo baseline
- [.github/instructions/ai-artifacts.instructions.md](.github/instructions/ai-artifacts.instructions.md) — Artifact selection matrix

Backend:
- [.github/instructions/backend-core.instructions.md](.github/instructions/backend-core.instructions.md) — Architecture
- [.github/instructions/go-backend-patterns.instructions.md](.github/instructions/go-backend-patterns.instructions.md) — Go patterns
- [.github/instructions/backend-testing.instructions.md](.github/instructions/backend-testing.instructions.md) — Testing
- [.github/instructions/errors-handling.instructions.md](.github/instructions/errors-handling.instructions.md) — Error patterns
- [.github/instructions/responses-dtos-swagger.instructions.md](.github/instructions/responses-dtos-swagger.instructions.md) — DTOs/OpenAPI

Portal:
- [.github/instructions/portal-web-architecture.instructions.md](.github/instructions/portal-web-architecture.instructions.md) — Architecture
- [.github/instructions/portal-web-code-structure.instructions.md](.github/instructions/portal-web-code-structure.instructions.md) — Code structure
- [.github/instructions/ui-implementation.instructions.md](.github/instructions/ui-implementation.instructions.md) — UI/RTL
- [.github/instructions/forms.instructions.md](.github/instructions/forms.instructions.md) — Forms
- [.github/instructions/i18n-translations.instructions.md](.github/instructions/i18n-translations.instructions.md) — i18n
- [.github/instructions/http-tanstack-query.instructions.md](.github/instructions/http-tanstack-query.instructions.md) — HTTP/TanStack Query

Agent OS:
- [KYORA_AGENT_OS.md](./KYORA_AGENT_OS.md) — Operating model, routing, lanes, handoffs
