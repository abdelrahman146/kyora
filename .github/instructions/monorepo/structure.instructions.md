---
description: "Kyora monorepo structure — directory layout, app/service naming, code placement, AGENTS.md cascade"
applyTo: "**/AGENTS.md,Makefile,docker-compose*.yml"
---

# Kyora Monorepo Structure — Single Source of Truth (SSOT)

**SSOT Hierarchy:**

- Parent: `.github/copilot-instructions.md`, `AGENTS.md`
- Related: `monorepo/workflows.instructions.md`, `monorepo/adding-projects.instructions.md`

**When to Read:**

- Starting work on Kyora for the first time
- Adding a new app or service
- Deciding where to put new code
- Understanding the AGENTS.md cascade pattern

---

## 1) Directory Structure

```
kyora/                              # Monorepo root
├── .github/
│   ├── instructions/               # Shared instruction files (SSOT)
│   │   ├── _meta/                  # Meta-guidelines (AI artifacts)
│   │   ├── kyora/                  # Kyora business/brand guidelines
│   │   ├── backend/                # Backend general patterns
│   │   │   ├── _general/           # Cross-backend patterns
│   │   │   └── projects/           # Backend project-specific
│   │   │       └── kyora-backend/  # Kyora backend specifics
│   │   ├── frontend/               # Frontend general patterns
│   │   │   ├── _general/           # Cross-frontend patterns
│   │   │   └── projects/           # Frontend project-specific
│   │   │       └── portal-web/     # Portal web specifics
│   │   ├── domain/                 # Domain module SSOTs (account, orders, etc.)
│   │   ├── platform/               # Platform/infra SSOTs (errors, assets, tenant isolation)
│   │   └── monorepo/               # Monorepo structure & workflows
│   ├── agents/                     # Custom GitHub Copilot agents
│   ├── prompts/                    # Reusable prompt files
│   ├── skills/                     # Multi-step workflow skills
│   └── workflows/                  # GitHub Actions CI/CD
│
├── AGENTS.md                       # Root agent manifest
├── KYORA_AGENT_OS.md              # Agent operating model
├── Makefile                        # Unified dev/build/test commands
├── docker-compose.dev.yml         # Local dev infrastructure
│
├── backend/                        # Go API monolith
│   ├── AGENTS.md                   # Backend-specific agent manifest
│   ├── cmd/                        # CLI commands (server, seed, sync_plans, etc.)
│   ├── docs/                       # Swagger/OpenAPI
│   ├── internal/
│   │   ├── domain/                 # Business logic (account, order, inventory, etc.)
│   │   ├── platform/               # Infrastructure (auth, database, email, etc.)
│   │   ├── server/                 # HTTP server + routes
│   │   └── tests/                  # E2E and test utilities
│   ├── go.mod
│   └── main.go
│
├── portal-web/                     # React dashboard SPA
│   ├── AGENTS.md                   # Portal-web-specific agent manifest
│   ├── src/
│   │   ├── api/                    # API client modules
│   │   ├── components/             # Shared UI components
│   │   ├── features/               # Feature modules
│   │   ├── routes/                 # Page routes (TanStack Router)
│   │   ├── i18n/                   # Translations (ar/, en/)
│   │   ├── lib/                    # Utilities
│   │   ├── stores/                 # Global state (TanStack Store)
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
│
└── (future apps/services)
    ├── storefront-web/             # Public storefront (React SPA)
    ├── mobile-web/                 # Mobile web app (React SPA)
    ├── admin-api/                  # Admin API service (Go)
    └── shared/                     # Shared libraries (future)
```

---

## 2) App/Service Naming Conventions

### Backend Services

- **Pattern**: `<domain>-<language>` or `<domain>-<purpose>`
- **Examples**:
  - `backend/` — Main API monolith (Go)
  - `admin-api/` — Admin API service (Go)
  - `webhook-worker/` — Background job processor (Go)

### Frontend Apps

- **Pattern**: `<product>-<platform>`
- **Examples**:
  - `portal-web/` — Seller dashboard (React + Vite)
  - `storefront-web/` — Public storefront (React + Vite)
  - `mobile-web/` — Mobile web app (React + Vite)

### Shared Libraries

- **Pattern**: `shared/` or `libs/`
- **Structure**:
  ```
  shared/
  ├── ts-utils/      # Shared TypeScript utilities
  ├── go-utils/      # Shared Go utilities
  └── design-tokens/ # Shared design tokens (future)
  ```

---

## 3) Where to Put New Code

### Decision Tree

```
Is it business logic?
├─ Yes → Is it backend or frontend?
│   ├─ Backend → `backend/internal/domain/<module>/`
│   └─ Frontend → `portal-web/src/features/<module>/`
│
├─ No → Is it infrastructure/platform?
│   ├─ Backend → `backend/internal/platform/<category>/`
│   └─ Frontend → `portal-web/src/lib/<category>/`
│
└─ Is it reusable across apps?
    ├─ Backend → `backend/internal/platform/<category>/`
    ├─ Frontend → `portal-web/src/lib/<category>/`
    └─ Cross-app → `shared/<package>/` (future)
```

### Examples

| Code Type               | Location                                                     |
| ----------------------- | ------------------------------------------------------------ |
| Order creation logic    | `backend/internal/domain/order/service.go`                   |
| Order list page         | `portal-web/src/features/orders/`                            |
| JWT validation          | `backend/internal/platform/auth/jwt.go`                      |
| HTTP client             | `portal-web/src/api/client.ts`                               |
| Form utilities          | `portal-web/src/lib/form/`                                   |
| Email templates         | `backend/internal/platform/email/templates/`                 |
| Shared TypeScript types | `portal-web/src/types/` (or `shared/ts-utils/` if cross-app) |

---

## 4) AGENTS.md Cascade Pattern

Kyora uses a **cascading AGENTS.md pattern** to scope agent behavior:

```
Root (kyora/AGENTS.md)
  ├─ Monorepo-level: Tech stack, structure, SSOT references
  │
  ├─ Backend (backend/AGENTS.md)
  │   └─ Backend-specific: Go patterns, domain modules, testing
  │
  └─ Portal Web (portal-web/AGENTS.md)
      └─ Frontend-specific: React patterns, TanStack, i18n, RTL
```

### Cascade Rules

1. **Root AGENTS.md** provides:
   - High-level project overview
   - Tech stack summary
   - Build/test/validate commands
   - SSOT entry points (organized by category)

2. **App-specific AGENTS.md** provides:
   - App-specific patterns
   - App-specific commands
   - App-specific SSOT references
   - **Must reference root AGENTS.md** (avoid duplication)

3. **Always read root AGENTS.md first**, then app-specific

### Example: Backend AGENTS.md

```markdown
# Backend (Go API)

For monorepo structure and SSOT references, see [root AGENTS.md](../AGENTS.md).

## Backend-Specific Patterns

- Domain modules: `backend/internal/domain/**`
- Platform/infra: `backend/internal/platform/**`
- Tests: `backend/internal/tests/**`

## Commands

See root Makefile for all commands. Backend-specific:

- `make dev.server` — Run API only
- `make test` — Backend tests
- `make openapi` — Regenerate Swagger docs

## SSOT References

- [backend/_general/architecture.instructions.md] — Backend domain structure
- [backend/_general/go-patterns.instructions.md] — Go patterns
- [backend/projects/kyora-backend/domain-modules.instructions.md] — Domain modules
```

---

## 5) Adding New Apps/Services

See `monorepo/adding-projects.instructions.md` for detailed scaffolding instructions.

**Quick checklist**:

- [ ] Create app directory (`<name>/`)
- [ ] Add app-specific `AGENTS.md`
- [ ] Add Makefile targets (`dev.<name>`, `test.<name>`, `build.<name>`)
- [ ] Add docker-compose service (if needed)
- [ ] Add GitHub Actions workflow (if separate deployment)
- [ ] Update root `AGENTS.md` with app references

---

## 6) Monorepo Boundaries

### What Belongs in Root

- Shared instruction files (`.github/instructions/`)
- Shared workflows (`.github/workflows/`)
- Root documentation (`AGENTS.md`, `KYORA_AGENT_OS.md`, `README.md`)
- Infrastructure config (`docker-compose.dev.yml`, `Makefile`)

### What Belongs in Apps

- App-specific code (`backend/`, `portal-web/`, etc.)
- App-specific config (`go.mod`, `package.json`, `vite.config.ts`)
- App-specific documentation (`backend/AGENTS.md`, `portal-web/AGENTS.md`)

### What Belongs in Shared (Future)

- Reusable libraries used by multiple apps
- Shared design tokens (CSS variables, color palettes)
- Shared TypeScript types (if not API-generated)

---

## 7) Navigation Guide

**Finding the right file:**

1. **Starting point**: Root `AGENTS.md` or `.github/copilot-instructions.md`
2. **Backend patterns**: `backend/_general/` or `backend/projects/kyora-backend/`
3. **Frontend patterns**: `frontend/_general/` or `frontend/projects/portal-web/`
4. **Domain logic**: `domain/<module>.instructions.md`
5. **Platform/infra**: `platform/<topic>.instructions.md`
6. **Kyora business**: `kyora/<topic>.instructions.md`
7. **Meta (AI artifacts)**: `_meta/<artifact-type>.instructions.md`

**Quick search:**

```bash
# Find instruction files by keyword
grep -r "keyword" .github/instructions/

# List all instruction files
find .github/instructions -name "*.instructions.md"

# Find AGENTS.md files
find . -name "AGENTS.md"
```

---

## 8) References

**Monorepo patterns:**

- [monorepo/workflows.instructions.md] — Makefile, docker-compose, CI/CD
- [monorepo/adding-projects.instructions.md] — Scaffolding new apps

**App-specific:**

- [backend/_general/architecture.instructions.md] — Backend structure
- [frontend/projects/portal-web/architecture.instructions.md] — Portal web structure

**Root:**

- [AGENTS.md] — Root agent manifest
- [KYORA_AGENT_OS.md] — Agent operating model
- [.github/copilot-instructions.md] — Copilot baseline
