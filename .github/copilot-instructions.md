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

- Backend architecture/patterns: `.github/instructions/backend-core.instructions.md`, `.github/instructions/go-backend-patterns.instructions.md`
- Backend testing: `.github/instructions/backend-testing.instructions.md`
- Portal architecture/dev workflow: `.github/instructions/portal-web-architecture.instructions.md`, `.github/instructions/portal-web-development.instructions.md`
- Portal UI/RTL: `.github/instructions/ui-implementation.instructions.md`, `.github/instructions/design-tokens.instructions.md`
- Portal forms/HTTP/i18n: `.github/instructions/forms.instructions.md`, `.github/instructions/http-tanstack-query.instructions.md`, `.github/instructions/i18n-translations.instructions.md`
