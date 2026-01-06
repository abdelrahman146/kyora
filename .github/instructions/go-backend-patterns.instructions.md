---
description: Go Backend Project Patterns (Kyora-style, reusable)
applyTo: "backend/**/*.go,backend/**/*_test.go,admin-backend/**/*.go,admin-backend/**/*_test.go,admin-backend-api/**/*.go,admin-backend-api/**/*_test.go"
---

# Go Backend Project Patterns (Reusable SSOT)

This file documents a **reusable Go backend pattern** based on Kyora’s backend.
Use it when creating or extending any Go API in this repo now or in the future (e.g., an admin backend).

Scope:

- **How to structure the project** (folders, dependencies, layering)
- **How to wire HTTP, config, logging, and persistence**
- **What rules are strict** (security, request validation, errors, transactions)

If you are modifying the Kyora backend specifically, also read:

- `.github/instructions/backend-core.instructions.md` (Kyora backend: ground truth wiring)

If you are writing tests, also read:

- `.github/instructions/backend-testing.instructions.md`

## Core architecture (layered, dependency direction)

Recommended baseline layout:

```
<project>/
  cmd/                 # Cobra CLI entrypoints (server, jobs, scripts)
  internal/
    server/            # HTTP server wiring and route registration (DI only)
    platform/          # Cross-cutting infrastructure (config, db, cache, auth, request/response, logger, types, utils)
    domain/            # Business modules (model/storage/service/http)
    tests/             # E2E/integration test harness (optional)
  main.go              # CLI entry (delegates to cmd)
```

Strict dependency rules:

- `internal/domain/**` may depend on `internal/platform/**`.
- `internal/platform/**` must not depend on `internal/domain/**`.
- Prefer cross-domain calls via **domain services** (not storage).
- Keep infrastructure decisions in `platform` and business invariants in `domain`.

## HTTP server wiring (DI only)

`internal/server` should:

- Create shared platform dependencies once (config, logger, DB, cache, integrations).
- Construct domain storages → services → handlers.
- Register routes and middleware chains.

Strict rules:

- Do not embed business logic in routing.
- Keep handlers thin: parse/validate → call service → respond.

## Configuration (SSOT is code constants)

Pattern:

- Define every config key as a constant in `internal/platform/config`.
- Provide sane defaults in a single `Configure()` function.
- Load config once at CLI startup (Cobra `PersistentPreRunE`) and avoid process exits inside config helpers.

Strict rules:

- No “magic strings” for config keys.
- Defaults must be safe for local dev.

## Logging (structured and request-scoped)

Pattern:

- Use structured logging (`log/slog`).
- Add a request middleware that:
  - generates/propagates a trace id
  - attaches a logger into the request context
  - logs request start and completion

Strict rules:

- Never log secrets (tokens, API keys, refresh tokens).
- Enrich logs with actor/tenant context *after* authentication.

## Request validation (strict JSON)

Pattern:

- Use a single helper to decode JSON bodies that:
  - rejects unknown fields (`DisallowUnknownFields`)
  - rejects trailing tokens
  - validates using struct tags

Strict rules:

- Do not use ad-hoc decoding per handler.
- Prefer explicit request DTO structs with `binding:` tags.

## Responses and errors (RFC 7807)

Pattern:

- Represent API errors using a Problem JSON type (`application/problem+json`).
- Route all errors through a single helper, which:
  - writes Problem JSON
  - aborts the request
  - maps common DB errors to stable HTTP codes

Strict rules:

- Never return inconsistent error shapes.
- Domain errors should be created once (domain `errors.go`) and enriched with `.With(key, value)`.

## Persistence (repository + scopes)

Pattern:

- Use a typed repository wrapper per model and a set of reusable query scopes.
- Use schema objects to map API-facing JSON fields to DB columns.

Strict rules:

- Do not accept raw DB column names from clients.
- Avoid raw SQL for domain logic; if you must use custom WHERE clauses, bind variables (no string concatenation).

## Transactions (atomic processor)

Pattern:

- Provide an “atomic processor” abstraction that runs a callback inside a transaction.
- Support:
  - isolation level selection
  - retries for retryable errors
  - reuse of an outer transaction if one is already in the context

Strict rules:

- Multi-entity writes must run inside an atomic transaction.
- Do not manually call begin/commit/rollback in domain services.

## Auth (Bearer access JWT + rotating refresh token sessions)

Pattern:

- Access token: short-lived JWT passed via `Authorization: Bearer <token>`.
- Refresh token: long-lived opaque token stored **hashed** server-side and rotated on refresh.
- Invalidate access tokens by including a server-checked version (`authVersion`) in JWT claims.

Strict rules:

- Never accept auth tokens via cookies unless the project explicitly chooses a cookie strategy.
- Treat refresh tokens like passwords: never log, store only hashes, revoke on use when rotating.

## Multi-tenancy (when applicable)

If the backend is multi-tenant, enforce boundaries at multiple layers:

- Middleware must load the authenticated actor and the tenant scope.
- Services must scope every query by tenant id.

Strict rules:

- Never trust tenant/workspace ids from URL params for authorization decisions.
- Provide “scoped getters” that prevent ID probing (BOLA), e.g., `GetWorkspaceUserByID(tenantID, userID)`.

## Testing strategy

Pattern:

- Prefer E2E/integration tests that boot real dependencies (DB/cache) for critical paths.
- Keep unit tests focused and isolated (mock external integrations).

Strict rules:

- Tests must be isolated and idempotent (truncate/cleanup between tests).
- Do not use raw SQL in tests when a domain storage/service exists.

## Public reference implementation (Kyora)

If you need a concrete example of this pattern as implemented today:

- Entry + lifecycle: `backend/main.go`, `backend/cmd/root.go`, `backend/cmd/server.go`
- DI + engine setup: `backend/internal/server/server.go`
- Routes + middleware: `backend/internal/server/routes.go`
- Auth middleware + JWT: `backend/internal/platform/auth/*`
- Request/response/problem: `backend/internal/platform/request/*`, `backend/internal/platform/response/*`, `backend/internal/platform/types/problem/*`
- DB repo + atomic tx: `backend/internal/platform/database/*`

## Anti-patterns (avoid)

- Putting business rules in handlers/routes.
- Returning ad-hoc error JSON.
- Accepting tenant/workspace ids from clients for authorization.
- Building SQL strings from user input.
- Using `float64` for money.
- Logging secrets.
- Calling another domain’s storage layer directly (bypass service invariants).
