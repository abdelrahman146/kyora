---
description: Backend Core Architecture — Go Monolith Patterns
applyTo: "backend/**"
---

# Backend Core Architecture

**Stack**: Go 1.25.3 monolith, Cobra CLI, Gin, GORM, Viper, Stripe, Memcached  
**Entry**: `backend/main.go` → `cmd/root.go` → HTTP via `cmd/server.go` → `internal/server/server.go`

## Project Structure

```
backend/
├── internal/
│   ├── platform/     # Infrastructure (config, DB, cache, auth, logging, types, utils)
│   ├── domain/       # Business logic modules (per-domain: model, storage, service, errors, handler_http)
│   └── tests/        # e2e/ (integration), testutils/ (helpers) — See backend-testing.instructions.md
└── .kyora.yaml       # Config (copy from .kyora.yaml.example)
```

## Run & Configure

**Local Dev**: `make dev.server` (requires `air`). Watches, builds to `tmp/main`, runs `kyora server`.

**Config SSOT**: `internal/platform/config/config.go` — All keys as constants.

**Key Settings**:

- App: `app.port`, `app.domain`, `app.notifications_email`
- HTTP: `http.port`, `http.base_url`, `http.trace_id_header` (default: `X-Trace-ID`)
- Database: `database.dsn`, `database.max_open_conns`, `database.log_level`
- Cache: `cache.hosts` (memcached)
- JWT: `auth.jwt.secret`, `auth.jwt.expiry_seconds`, `auth.jwt.issuer`
- Tokens: `auth.password_reset_ttl_seconds`, `auth.verify_email_ttl_seconds`, `auth.invitation_token_ttl_seconds`
- OAuth: `auth.google_oauth.{client_id, client_secret, redirect_url}`
- Stripe: `billing.stripe.{api_key, webhook_secret}`
- Email: `email.provider` (resend/mock), `email.resend.api_key`, `email.from_email`
- Logging: `log.format`, `log.level`

## Database Access

**Repository Pattern**: Generic `database.Repository[T]` per domain.

```go
// storage.go
type Storage struct {
    order     *database.Repository[Order]
    orderItem *database.Repository[OrderItem]
    cache     *cache.Cache
}
func NewStorage(db *database.Database, cache *cache.Cache) *Storage {
    return &Storage{
        order:     database.NewRepository[Order](db), // auto-migrates
        orderItem: database.NewRepository[OrderItem](db),
        cache:     cache,
    }
}
```

**Connection**: `db.Conn(ctx)` — transaction-aware via `TxKey` context.

**Scopes** (chainable filters):

- Workspace: `ScopeWorkspaceID(id)`, `ScopeBusinessID(id)` (legacy)
- ID: `ScopeID(id)`, `ScopeIDs(ids...)`
- Comparison: `ScopeEquals(field, val)`, `ScopeIn(field, vals)`, `ScopeNotIn`, `ScopeGreaterThan`, `ScopeLessThan`, `ScopeBetween`
- Time: `ScopeTime(field, from, to)`, `ScopeCreatedAt(from, to)`
- Search: `ScopeSearchTerm(term, fields...)`, `ScopeIsNull(field)`

**Query Modifiers**:

- Load: `WithPreload(assocs...)`, `WithJoins(assocs...)`
- Page: `WithPagination(page, pageSize)`, `WithLimit(n)`
- Order: `WithOrderBy(fields...)` — SQL format: `["total DESC", "createdAt ASC"]`
- Lock: `WithLockingStrength(database.LockUpdate)` → FOR UPDATE / SHARE / SKIP LOCKED / NOWAIT
- Return: `WithReturning(fields...)` — Postgres RETURNING clause

**CRUD**:

- Create: `CreateOne(ctx, model)`, `CreateMany(ctx, models)`
- Update: `UpdateOne(ctx, model, scopes...)`, `UpdateMany(ctx, updates, scopes...)`
- Delete: `DeleteOne(ctx, scopes...)`, `DeleteMany(ctx, scopes...)`
- Find: `FindByID(ctx, id)`, `FindOne(ctx, scopes...)`, `FindMany(ctx, scopes...)`
- Count: `Count(ctx, scopes...)`

**Aggregates**:

- Scalar: `Sum(ctx, column, scopes...)`, `Avg(ctx, column, scopes...)`
- Group: `SumBy(ctx, valueCol, groupCol, scopes...)`, `CountBy(ctx, groupCol, scopes...)`, `AvgBy` → `[]keyvalue.KeyValue`
- Time Series: `TimeSeriesSum(ctx, valueCol, timeCol, granularity, opts...)`, `TimeSeriesCount` → `*timeseries.TimeSeries`

## Transactions

**Atomic Operations**: Multi-step writes MUST use atomic processor.

```go
err := s.db.AtomicProcess.Exec(ctx, func(tctx context.Context) error {
    // All repo calls use tctx (transaction injected via TxKey)
    if err := s.storage.order.CreateOne(tctx, order); err != nil {
        return err
    }
    return s.storage.inventory.UpdateMany(tctx, updates, scopes...)
}, atomic.WithIsolationLevel(atomic.LevelSerializable), atomic.WithRetries(3))
```

**Reference**: `internal/platform/database/atomic.go`, `internal/platform/types/atomic/atomic.go`

## HTTP Layer

**Auth Middleware Chain** (workspace-based):

1. `auth.EnforceAuthentication` — JWT from `Authorization: Bearer <token>`, sets claims
2. `account.EnforceValidActor(svc)` — loads user, sets `ActorKey`, enriches logger
3. `account.EnforceWorkspaceMembership(svc)` — validates `workspaceId` param, sets `WorkspaceKey`
4. `account.EnforceActorPermissions(action, resource)` — RBAC check
5. Optional: `billing.EnforceActiveSubscription(svc)`, `billing.EnforcePlanWorkspaceLimits(limit, counterFunc)`

**Extract Context**:

- User: `account.ActorFromContext(c)`
- Workspace: `account.WorkspaceFromContext(c)`
- Claims: `auth.ClaimsFromContext(c)`
- JWT Token: `auth.JwtFromContext(c)` — extracts from Authorization header ONLY (no cookies)

**Token Flow**:

- Access: Short-lived JWT in `Authorization` header
- Refresh: Opaque server-side session (hashed token, expiry, metadata)
- Rotation: `POST /v1/auth/refresh` consumes `refreshToken`, returns new `{token, refreshToken}`
- Logout: `POST /v1/auth/logout` (revoke session), `/logout-others` (keep current), `/logout-all` (revoke all + bump `User.AuthVersion`)
- Invalidation: JWT claims include `authVersion`; middleware rejects if `claims.authVersion != user.AuthVersion`

**Request Validation**: `request.ValidBody(c, &req)` — bind + validate JSON body

**Responses**:

- Success: `response.SuccessJSON(c, status, data)`, `response.SuccessEmpty(c, status)`, `response.SuccessText`
- Error: `response.Error(c, err)` — auto-converts to RFC 7807 Problem JSON

**Error Handling**:

- Problem types: `problem.BadRequest()`, `Unauthorized()`, `Forbidden()`, `NotFound()`, `Conflict()`, `InternalError()`, `ValidationError(fields...)`
- DB errors: `database.IsRecordNotFound(err)`, `IsUniqueViolation(err)`
- Domain errors: Define in `errors.go` using problem constructors with `.With(key, val)`, `.WithError(err)`

**Logging**: Structured `slog` logger. Context-aware: `logger.FromContext(ctx)`. Middleware adds `traceId` (from `http.trace_id_header`), actor info.

## Domain Conventions

**model.go**:

- Embed `gorm.Model` (adds ID, CreatedAt, UpdatedAt, DeletedAt)
- Override ID: `ID string \`gorm:"column:id;primaryKey;type:text"\``
- Generate IDs in `BeforeCreate` hook: `id.KsuidWithPrefix("ord")` (external), `id.Base62(6)` (short codes)
- Constants: `TableName`, `StructName`, `Prefix`
- Schema: `var OrderSchema = struct { Total schema.Field }{Total: schema.NewField("total", "total")}` — maps DB column → JSON field

**storage.go**:

- Typed repositories + cache
- Auto-migrate on construction: `database.NewRepository[Model](db)`

**service.go**:

- Signature: `(ctx context.Context, actor *account.User, workspace *account.Workspace, ...)`
- **Multi-tenancy**: ALWAYS scope by `WorkspaceID`: `s.storage.order.FindMany(ctx, s.storage.order.ScopeWorkspaceID(workspace.ID), ...)`
- Use `AtomicProcess.Exec` for multi-entity ops
- Pattern: Validate → Load related → Calculate → Persist → Return
- Updates: Load with locking → Validate transitions → Update → Persist

**errors.go**:

- Domain-specific problems: `problem.NotFound("order not found").WithError(err).With("orderId", id)`

**handler_http.go**:

- Extract actor/workspace from context
- Validate: `request.ValidBody`
- Call service
- Respond: `response.Success*` or `response.Error`

**middleware_http.go** (optional):

- Domain-specific checks (see `account/middleware_http.go`)

**state_machine.go** (optional):

- Valid state transitions + timestamp updates (see `order/state_machine.go`)

## Pagination & Analytics

**Lists**:

- Accept `*list.ListRequest` (page, pageSize, orderBy, searchTerm)
- Parse ordering: `req.ParsedOrderBy(DomainSchema)` → `["total DESC", "createdAt ASC"]` (JSON fields → SQL)
- Apply: `s.storage.order.WithOrderBy(req.ParsedOrderBy(OrderSchema))`
- Return: `list.NewListResponse(items, page, pageSize, totalCount, hasMore)`

**Time Series**:

- `TimeSeriesSum(ctx, valueCol, timeCol, granularity, opts...)`, `TimeSeriesCount(ctx, timeCol, granularity, opts...)`
- Granularity: `timeseries.GranularityHour/Day/Week/Month` or auto: `timeseries.GetTimeGranularityByDateRange(from, to)`
- Returns: `*timeseries.TimeSeries` (bucketed data points)

**Aggregations**:

- Group: `CountBy(ctx, groupCol, scopes...)`, `SumBy(ctx, valueCol, groupCol, scopes...)`, `AvgBy` → `[]keyvalue.KeyValue`
- Scalar: `Sum`, `Avg`, `Count`
- Time filter: `ScopeTime(field, from, to)`

## Integrations

**Event Bus** (in-memory pub/sub):

- Emit: `bus.Emit(topic, payload)` — non-blocking, async
- Listen: `unsub := bus.Listen(topic, func(payload any) {...})` — returns unsubscribe func
- Topics: `bus.VerifyEmailTopic`, `ResetPasswordTopic`, `WorkspaceInvitationTopic`
- Shutdown: `bus.Close()` — waits for handlers

**Stripe**:

- Domain: `internal/domain/billing` (plans, subscriptions, invoices, webhooks, customer management)
- Webhooks: `webhooks.go`, `stripe_retry.go`
- Middleware: Plan limits enforcement
- Config: `billing.stripe.{api_key, webhook_secret}`

**Email**:

- Platform: `internal/platform/email`
- Providers: `resend` (prod), `mock` (dev/test)
- Send: `client.SendTemplate(ctx, templateID, to, from, subject, data)`
- Templates: HTML in `internal/platform/email/templates/` (embedded via `go:embed`)
- Available: `TemplateForgotPassword`, `TemplateEmailVerification`, `TemplateWelcome`, `TemplateSubscriptionConfirmed`, etc.
- Reference: `internal/platform/email/README.md`

**OAuth**:

- Google flow: `internal/platform/auth/google_oauth.go`
- Config: `auth.google_oauth.{client_id, client_secret, redirect_url}`
- JWT helpers: `internal/platform/auth/jwt.go`

## Caching (Memcached)

- **Where**: Storage layer ONLY (keep service clean)
- **When**: Frequently accessed data, infrequent changes
- **Pattern**: Read-through (cache on first access, invalidate on update)
- **Keys**: Unique, descriptive (avoid collisions)
- **Service Use**: Expirable memory (auth tokens, password reset, email verification, invitations)
- **Reference**: `internal/platform/cache`

## Multi-Tenancy

- **Workspace**: Groups users under single subscription
- **User**: Belongs to one workspace (`WorkspaceID`), has role
- **Business**: Legacy (transitioning to workspace; some code uses `BusinessID`)
- **Scoping**: ALWAYS filter by `WorkspaceID` for data isolation
- **Invitations**: Users invited with role (admin/member)
- **Permissions**: `internal/platform/types/role` (actions: view/manage, resources: account/billing/orders)
- **Middleware**: `EnforceWorkspaceMembership` validates user belongs to workspace in URL param

## Routes & Permissions

- **Routes**: Define in `internal/server/routes.go` per domain
- **Middleware chain**: Authentication → Actor → Workspace → Permissions → Plan limits
- **Permissions**: `EnforceActorPermissions(action, resource)` in route chain
- **Plan limits**: `EnforcePlanWorkspaceLimits(limit, counterFunc)` before operations

## OpenAPI Documentation

- **Tool**: `github.com/swaggo/swag` — generates from code comments
- **Format**: Swaggo comment format in handlers
- **Include**: Request params, body, response, errors
- **Maintenance**: Keep docs in sync with code

## Separation of Concerns

**Handler**: Request/response, validation, transformations → service call  
**Service**: Business logic, data transformations → storage call  
**Storage**: Data access, caching — NO business logic  
**Model**: State transitions, derived values — NO business logic or data access

**Domain Independence**:

- No circular dependencies between domains
- Cross-domain communication via event bus
- Shared functionality → `internal/platform/utils/` or helpers
- Never access other domain storages directly — use their service

## Best Practices

**Security**:

- ALWAYS scope by `WorkspaceID` (or legacy `BusinessID`) — no cross-workspace leaks
- Use config constants from `internal/platform/config/config.go` — never hardcode
- JWT from Authorization header ONLY (format: `Bearer <token>`) — no cookies
- Treat `refreshToken` like password (log/redact carefully) — store only hashes
- Session invalidation: logout endpoints + `User.AuthVersion` bump on sensitive changes

**Data Handling**:

- Ordering/search: JSON field names (Schema) in requests, not DB columns
- ID generation: KSUID with prefix for external IDs, Base62 for short codes
- Decimal math: `shopspring/decimal` for financial calculations — NEVER `float64` for money
- Time zones: Store UTC — GORM uses `timestamptz`

**Code Quality**:

- Transactions: Use atomic processor — never manual begin/commit
- Errors: Return domain problems (RFC 7807) — use `database.IsRecordNotFound/IsUniqueViolation` for DB errors
- Logging: Structured `slog` via `logger.FromContext(ctx)` — context-aware, middleware-enriched
- State machines: Define allowed transitions + timestamp updates for entities with state
- Comments: Godoc style for all exported functions, structs, packages
- Deprecated code: DELETE immediately — no "marked as deprecated"
- Refactoring: Extract duplicates to `internal/platform/utils/` or helpers
- No TODOs/FIXMEs: Complete implementations only

**Code Standards**:

- KISS: Simple, unambiguous requirements satisfaction
- DRY: Extract, generalize, reuse — never repeat
- Readability: Junior developer understands immediately
- Pillars: 100% Robust, Reliable, Secure, Scalable, Optimized, Traceable, Testable
- Self-documenting: Clear names, minimal comments

## Reference Implementations

**Core Patterns**:

- Repository: `internal/platform/database/repository.go` (all scopes/methods)
- Transactions: `internal/platform/database/atomic.go`, `internal/platform/types/atomic/atomic.go`
- Service: `internal/domain/order/service.go` (CRUD, aggregations, time series)
- HTTP: `internal/domain/order/handler_http.go`, `internal/domain/billing/handler_http.go`
- Middleware: `internal/domain/account/middleware_http.go`, `internal/domain/billing/middleware_http.go`
- Errors: `internal/domain/order/errors.go` (domain problems)
- Models: `internal/domain/order/model.go` (GORM models, schemas, DTOs)
- State: `internal/domain/order/state_machine.go` (status transitions)

**Platform Helpers**:

- Routes: `internal/server/routes.go` (middleware chains, grouping, permissions)
- Server: `internal/server/server.go` (DI, service init, graceful shutdown)
- Response: `internal/platform/response/response.go`
- Request: `internal/platform/request/valid_body.go`
- Logger: `internal/platform/logger/middleware.go`
- List: `internal/platform/types/list/list_request.go`, `list_response.go`
- Time Series: `internal/platform/types/timeseries/timeseries.go`
- Problems: `internal/platform/types/problem/problem.go`

**Testing**: See [backend-testing.instructions.md](./backend-testing.instructions.md) for comprehensive test guidelines.
