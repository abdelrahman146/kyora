---
description: Backend Core Architecture — Go Monolith Patterns
applyTo: "backend/**"
---

# Backend Core Architecture (SSOT)

This file is the **source of truth** for how the Go backend is structured and how to extend it safely.
It intentionally documents **only patterns that exist in this repo**.

**Stack**: Go 1.25.3 monolith, Cobra CLI, Gin, GORM/Postgres, Viper, Memcached, Stripe, Resend, slog

If you are applying these patterns to a new Go backend project (not Kyora’s main backend), start with:

- `.github/instructions/go-backend-patterns.instructions.md`

## Non‑negotiables (must follow)

- **Multi-tenancy isolation**: never allow cross-workspace/business access.
- **Middleware-first security**: all protected routes must apply the correct middleware chain.
- **No raw SQL for domain work**: use the generic repository + schema/scopes.
- **No manual transactions**: use the atomic processor.
- **Strict JSON bodies**: use `request.ValidBody` (disallows unknown fields).
- **Consistent error shape**: always return RFC7807 Problem JSON via `response.Error`.

## Runtime entrypoints

- CLI entry: `backend/main.go` → `backend/cmd/root.go`.
- HTTP server: `kyora server` (`backend/cmd/server.go`) → `backend/internal/server/server.go`.

Key lifecycle facts:

- Config defaults are prepared via `internal/platform/config.Configure()`.
- In normal CLI runs, config is loaded once in Cobra `PersistentPreRunE` (`config.Load()`), then logging is initialized (`logger.Init()`).
- `server.New()` calls `config.Configure()` again to ensure defaults exist (important for tests / non-Cobra usage).

## Folder structure and dependency direction

```
backend/
  cmd/                    # CLI commands (server, seed, sync_plans, ...)
  internal/
    server/               # Gin engine wiring, DI, and route registration
    platform/             # Infrastructure primitives (config, db, cache, auth, request/response, logger, types, utils)
    domain/               # Business modules (account, billing, inventory, ...)
    tests/                # E2E/integration tests (see backend-testing.instructions.md)
```

**Dependency rules**:

- `internal/domain/**` may depend on `internal/platform/**`.
- `internal/platform/**` must not depend on `internal/domain/**`.
- Domains may depend on other **domain services** when needed (example: orders depend on inventory/customer/business services).
- Avoid “cross-domain storage access”: do not import another domain’s `Storage` and query its tables directly.

## HTTP server wiring (what actually happens)

Server initialization is in `internal/server/server.go`:

- Creates DB connection (`database.NewConnection`) and cache connection (`cache.NewConnection`).
- Creates a shared `AtomicProcess` (`database.NewAtomicProcess`).
- Creates shared integrations: `bus.New()`, `email.New()`, `blob.FromConfig()`.
- Constructs each domain’s `Storage` then `Service`, then `HttpHandler`.
- Builds a Gin engine with:
  - `logger.Middleware()` (trace id + request lifecycle logs)
  - `request.LimitBodySize(http.max_body_bytes)` (early body limit enforcement)
  - `gin.Recovery()`
- Adds health endpoints: `GET /healthz`, `GET /livez`.
- Adds a catch-all `OPTIONS /*path` handler that applies:
  - public CORS for `/v1/storefront...`
  - regular CORS for everything else

Routes are registered in `internal/server/routes.go`.

## CORS (important invariants)

Implementation: `internal/platform/middleware/cors.go`

- Uses Bearer auth (Authorization header), not cookies.
- `AllowCredentials` is **false**.
- Allowed origins are configured via `cors.allowed_origins` (defaults to `*`).
- Public endpoints (storefront + some public APIs) use a middleware that always allows `*`.

## Auth + tenancy middleware (required chains)

### Context keys you will rely on

- `auth.ClaimsKey` → JWT claims
- `account.ActorKey` → authenticated user
- `account.WorkspaceKey` → workspace loaded from the authenticated user
- `business.BusinessKey` → business loaded for `:businessDescriptor` within the user’s workspace
- `database.TxKey` → injected gorm transaction used by `Database.Conn(ctx)`

### Middleware chain (workspace-scoped protected routes)

From `internal/server/routes.go` and domain middleware implementations:

1. `middleware.NewCORSMiddleware()`
2. `auth.EnforceAuthentication`
   - JWT is taken from `Authorization: Bearer <token>`.
   - Parsed claims are stored in Gin context.
3. `account.EnforceValidActor(accountService)`
   - Loads the user from DB.
   - Enforces `claims.AuthVersion == user.AuthVersion` (access token invalidation).
   - Enriches the slog logger with actor fields.
4. `account.EnforceWorkspaceMembership(accountService)`
   - Loads workspace by `user.WorkspaceID`.
   - **Never** trusts a workspace id from URL params.
5. Optional per-endpoint:
   - `account.EnforceActorPermissions(role.Action*, role.Resource*)`
   - `billing.EnforceActiveSubscription(billingService)`
   - `billing.EnforcePlanWorkspaceLimits(...)`

### Middleware chain (business-scoped protected routes)

Same as workspace chain, plus:

6. `business.EnforceBusinessValidity(businessService)`
   - Loads business by `:businessDescriptor` but only for the authenticated user’s workspace.

### Public routes

- Storefront is mounted under `/v1/storefront` and intentionally uses public CORS.
- Stripe webhooks are public (`POST /webhooks/stripe`) and must verify signature inside the handler.

## Request validation (strict)

Use `internal/platform/request.ValidBody(c, &req)` in handlers.

It:

- Requires a request body.
- Uses a JSON decoder with `DisallowUnknownFields()`.
- Rejects trailing JSON tokens.
- Runs Gin’s validator on `binding:` tags.

Do not use `c.BindJSON`, `ShouldBindJSON`, or ad-hoc decoding.

## Responses and errors (RFC7807)

Success:

- `response.SuccessJSON(c, status, data)`
- `response.SuccessEmpty(c, status)`
- `response.SuccessText(c, status, text)`

Errors:

- Always use `response.Error(c, err)`.
- If `err` is not a `*problem.Problem`, `response.Error` maps:
  - record not found → 404
  - unique violation → 409
  - everything else → 500

Problem shape is `internal/platform/types/problem.Problem`.

## Configuration (SSOT is code constants)

All config keys must be constants in `internal/platform/config/config.go`.
Do not introduce “magic strings” in code.

Notable defaults set by `config.Configure()`:

- `http.max_body_bytes` defaults to 1 MiB.
- `database.auto_migrate` defaults to true.
- `billing.auto_sync_plans` defaults to true.
- `cors.allowed_origins` defaults to `[*]`.
- `auth.refresh_token_ttl_seconds` defaults to 30 days.
- Storage / uploads defaults (local provider under `./tmp/assets`).

## Database access (GORM + repository)

### Database connection

- `database.NewConnection(dsn, logLevel)` opens Postgres and configures pooling.
- Ensures `pg_trgm` extension exists (best-effort).

### Transaction-aware connections

Use `db.Conn(ctx)` everywhere.

- If `database.TxKey` exists in the context, it uses that transaction.
- Otherwise it uses the normal DB connection.

### Repository pattern

Domains create typed repositories in `storage.go`:

- `database.NewRepository[T](db)` (auto-migrates the model when enabled).

Important: auto-migration is controlled by `database.auto_migrate`.

### Common repository options (real APIs)

Query modifiers:

- `WithPreload(associations...)`
- `WithJoins(joins...)`
- `WithOrderBy([]string{"created_at DESC"})` (already-parsed SQL order strings)
- `WithOrderByExpr(clause.Expr)` for complex, injection-safe ordering
- `WithPagination(offset, limit)`
- `WithLimit(limit)`
- `WithLockingStrength(database.LockingStrengthUpdate|Share|SkipLocked|NoWait)`
- `WithReturning(value any)` (applies `RETURNING` and scans into `value`)

Scopes (filters):

- `ScopeWorkspaceID(workspaceID)` / `ScopeBusinessID(businessID)`
- `ScopeID(id)` / `ScopeIDs([]any{...})`
- `ScopeEquals(field, value)`, `ScopeNotEquals`, `ScopeIn`, `ScopeNotIn`
- `ScopeGreaterThan`, `ScopeLessThan`, `ScopeBetween`, and `...OrEqual` variants
- `ScopeCreatedAt(from, to)` / `ScopeTime(field, from, to)`
- `ScopeSearchTerm(term, fields...)`
- `ScopeWhere(sql, vars...)` (bound vars; use sparingly for specialized queries)

### Pagination, ordering, search (canonical pattern)

For list endpoints, use `internal/platform/types/list.ListRequest`:

- Convert page/pageSize → `Offset()` and `Limit()`.
- Convert API-facing orderBy fields to SQL columns via `ParsedOrderBy(DomainSchema)`.
- Normalize search terms with `list.NormalizeSearchTerm`.

Do not accept DB column names directly from clients.

## Transactions (atomic processor)

Use `database.NewAtomicProcess(db)` and call `AtomicProcess.Exec(ctx, fn, opts...)`.

Key behavior:

- If `database.TxKey` already exists in the context, `Exec` reuses it (no nested transaction).
- Retries retryable transaction failures with exponential backoff + jitter.
- Supports isolation level and read-only options via `internal/platform/types/atomic`.

Never call `db.Begin()/Commit()/Rollback()` directly in domain services.

## Auth model (access JWT + rotating refresh sessions)

Access token:

- JWT created by `auth.NewJwtToken(userID, workspaceID, authVersion)`.
- Required on protected endpoints via `Authorization: Bearer <token>`.

Refresh token:

- Generated via `auth.NewRefreshToken()`.
- Stored **hashed** in DB (see `account.Service.issueTokensForUser`).
- Rotated on refresh (`/v1/auth/refresh`): old session revoked, new session created.

Invalidation:

- `account.EnforceValidActor` rejects access tokens when `claims.AuthVersion != user.AuthVersion`.
- `/v1/auth/logout-all` increments `user.AuthVersion` and revokes all sessions.

Strict rule: treat `refreshToken` as a password. Never log it.

## Domain module conventions (how to add/change behavior)

Each domain lives under `internal/domain/<name>/` and commonly contains:

- `model.go`: GORM models + schema definitions (JSON field ↔ DB column mapping).
- `storage.go`: repositories and cache wiring. Keep caching here.
- `service.go`: business logic. This is where invariants live.
- `errors.go`: domain-specific `problem.*` constructors.
- `handler_http.go`: HTTP handlers (parse/validate → service → respond).
- `middleware_http.go`: optional domain-specific middleware.
- `state_machine.go`: optional explicit transition rules.

Service rules:

- Always scope reads/writes by workspace/business.
- Prefer helper methods that prevent BOLA patterns (example: `account.Service.GetWorkspaceUserByID`).
- Use `shopspring/decimal` for money-like values (never `float64`).
- Use UTC for timestamps.

Handler rules:

- Always use `request.ValidBody`.
- Always use `response.Success*` / `response.Error`.
- Extract actor/workspace/business via the `FromContext` helpers.

## Integrations (SSOT pointers)

- Email (Resend/mock): see resend.instructions.md.
- Stripe billing: see stripe.instructions.md.
- Asset uploads/blob storage: see asset_upload.instructions.md.
- Testing (E2E, unit patterns): see backend-testing.instructions.md.

## Anti-patterns (avoid)

- Accepting `workspaceId` from URL/body for authorization decisions.
- Skipping `DisallowUnknownFields` decoding by using Gin bind helpers.
- Building SQL strings from user input (order/search filters) without schema mapping.
- Using floats for money.
- Performing multi-entity writes without `AtomicProcess.Exec`.
- Calling another domain’s storage layer directly.

## Reference implementations (ground truth)

- Server wiring + DI: `backend/internal/server/server.go`
- Routes and middleware chains: `backend/internal/server/routes.go`
- Auth middleware + JWT helpers: `backend/internal/platform/auth/middleware.go`, `backend/internal/platform/auth/jwt.go`
- Actor/workspace/business middleware: `backend/internal/domain/account/middleware_http.go`, `backend/internal/domain/business/middleware_http.go`
- Request validation: `backend/internal/platform/request/valid_body.go`
- Response + error mapping: `backend/internal/platform/response/response.go`
- Repository + scopes: `backend/internal/platform/database/repository.go`
- Atomic transactions: `backend/internal/platform/database/atomic.go`
