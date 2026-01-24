# Kyora

Kyora is a DM-first commerce platform for micro-businesses in the Middle East. It helps sellers feel “handled” by keeping orders, stock, and money organized in plain language.

Current repo has **two active apps**: `backend/` (Go API) and `portal-web/` (React dashboard). **More apps/services will be added later** (e.g., admin portal/service, public storefront, mobile app, website), so never assume this list is exhaustive.

## Tech stack

- Backend: Go (Gin, GORM/Postgres), PostgreSQL, Memcached, Stripe, Resend
- Frontend: React + Vite + TanStack (Router/Query/Form/Store), i18n (Arabic/RTL-first), Tailwind + daisyUI
- Tests: Go tests under `backend/` (E2E in `backend/internal/tests/e2e/`)

## Repo map (current)

- `backend/`: Go API monolith (domain + platform layers), Swagger/OpenAPI under `backend/docs/`
- `portal-web/`: dashboard SPA
- `.github/instructions/`: scoped SSOT rules (prefer these over duplicating guidance here)
- `docker-compose.dev.yml`: local Postgres + Memcached + Stripe mock
- `Makefile`: verified local dev/test commands

## Build / test / validate (verified)

- See all targets: `make help`
- Tooling sanity: `make doctor`

- Infra (local Postgres/Memcached/Stripe mock): `make infra.up` (or `make infra.down`, `make infra.reset`, `make infra.logs`)
- DB shell: `make db.psql`

- Dev (API + portal): `make dev` (or `make dev.infra` to start infra first)
- Dev (API only): `make dev.server`
- Dev (portal only): `make dev.portal` (optionally `PORTAL_PORT=3001 make dev.portal`)

- Backend tests: `make test` (or `make test.unit`, `make test.e2e`, `make test.quick`)
- Backend OpenAPI: `make openapi` (and `make openapi.check` / `make openapi.verify` to ensure docs are up-to-date)

- Portal deps: `make portal.install`
- Portal checks: `make portal.check`
- Portal build: `make portal.build` (preview: `make portal.preview`)

## Coding rules (high-signal)

- MUST keep strict tenant isolation: workspace is top-level, business is second-level; no cross-scope reads/writes.
- MUST validate inputs and return domain errors (use Kyora Problem/RFC7807 patterns where applicable).
- MUST avoid accounting jargon in UI text; prefer “Profit”, “Cash in hand”, “Money in/out”, “Best seller”, “What to do next”.
- MUST be Arabic/RTL-first in portal UI; don’t assume left/right or English-only labels.
- MUST not add TODO/FIXME placeholders; ship complete implementations.
- SHOULD prefer existing shared utils (`backend/internal/platform/utils/`, `portal-web/src/lib/`) over duplication.
- SHOULD follow existing module boundaries: backend domain logic in `backend/internal/domain/**`, infra in `backend/internal/platform/**`.
- SHOULD use the portal HTTP/TanStack Query layer (don’t ad-hoc `fetch`); see the scoped HTTP instructions.
- MAY introduce breaking changes (repo is under heavy development), but keep changes focused.

## Resources (SSOT)

**Kyora Business/Brand:**

- `.github/instructions/kyora/brand-key.instructions.md` — Brand key model
- `.github/instructions/kyora/business-model.instructions.md` — Business model
- `.github/instructions/kyora/target-customer.instructions.md` — Target customer
- `.github/instructions/kyora/ux-strategy.instructions.md` — UX strategy
- `.github/instructions/kyora/design-system.instructions.md` — Design system

**Backend (General Patterns):**

- `.github/instructions/backend/_general/architecture.instructions.md` — Architecture
- `.github/instructions/backend/_general/go-patterns.instructions.md` — Go patterns
- `.github/instructions/backend/_general/api-contracts.instructions.md` — API contracts
- `.github/instructions/backend/_general/errors.instructions.md` — Error handling
- `.github/instructions/backend/_general/testing.instructions.md` — Testing

**Backend (Kyora Project-Specific):**

- `.github/instructions/backend/projects/kyora-backend/domain-modules.instructions.md` — Domain modules
- `.github/instructions/backend/projects/kyora-backend/integrations.instructions.md` — Integrations
- `.github/instructions/backend/projects/kyora-backend/testing-specifics.instructions.md` — Testing specifics

**Frontend (General Patterns):**

- `.github/instructions/frontend/_general/architecture.instructions.md` — Architecture
- `.github/instructions/frontend/_general/ui-patterns.instructions.md` — UI/RTL patterns
- `.github/instructions/frontend/_general/forms.instructions.md` — Forms (core)
- `.github/instructions/frontend/_general/forms-validation.instructions.md` — Forms validation
- `.github/instructions/frontend/_general/http-client.instructions.md` — HTTP/Ky/TanStack Query
- `.github/instructions/frontend/_general/i18n.instructions.md` — i18n
- `.github/instructions/frontend/_general/testing.instructions.md` — Testing

**Frontend (Portal Web Project-Specific):**

- `.github/instructions/frontend/projects/portal-web/architecture.instructions.md` — Architecture
- `.github/instructions/frontend/projects/portal-web/code-structure.instructions.md` — Code structure
- `.github/instructions/frontend/projects/portal-web/ui-components.instructions.md` — UI components (daisyUI)
- `.github/instructions/frontend/projects/portal-web/charts.instructions.md` — Charts (Chart.js)
- `.github/instructions/frontend/projects/portal-web/development.instructions.md` — Development workflow

**Domain Modules:**

- `.github/instructions/domain/account.instructions.md` — Account/auth/workspace
- `.github/instructions/domain/accounting.instructions.md` — Accounting
- `.github/instructions/domain/analytics.instructions.md` — Analytics
- `.github/instructions/domain/billing.instructions.md` — Billing
- `.github/instructions/domain/business.instructions.md` — Business management
- `.github/instructions/domain/customer.instructions.md` — Customer management
- `.github/instructions/domain/inventory.instructions.md` — Inventory
- `.github/instructions/domain/onboarding.instructions.md` — Onboarding
- `.github/instructions/domain/orders.instructions.md` — Orders

**Platform/Infrastructure:**

- `.github/instructions/platform/asset-upload.instructions.md` — Asset upload
- `.github/instructions/platform/errors-handling.instructions.md` — Error handling (cross-cutting)
- `.github/instructions/platform/tenant-isolation.instructions.md` — Tenant isolation

**Monorepo:**

- `.github/instructions/monorepo/structure.instructions.md` — Structure
- `.github/instructions/monorepo/workflows.instructions.md` — Workflows
- `.github/instructions/monorepo/adding-projects.instructions.md` — Adding projects

**Meta (AI Artifacts):**

- `.github/instructions/_meta/ai-artifacts.instructions.md` — Artifact selection matrix
- `.github/instructions/_meta/agents.instructions.md` — Agent authoring
- `.github/instructions/_meta/prompts.instructions.md` — Prompt authoring
- `.github/instructions/_meta/skills.instructions.md` — Skill authoring
- `.github/instructions/_meta/writing-instructions.instructions.md` — Writing instructions
- `.github/instructions/_meta/copilot-instructions.instructions.md` — Copilot instructions
