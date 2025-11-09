Kyora — AI coding agent quickstart

Big picture

- Go 1.25.3 monolith using Cobra CLI framework. Entry point: `main.go` → `cmd/root.go`. HTTP server runs via `server` subcommand in `cmd/server.go` which delegates to `internal/server/server.go`.
- Go dependencies: Gin web framework, GORM ORM, Viper config, Stripe SDK, gomemcache, JWT, OAuth2, and more (see `go.mod`).
- Layered `internal/` architecture:
  - **platform**: Infrastructure layer (config, database, cache, logging, auth, event bus, request/response helpers, shared types).
  - **domain**: Business logic modules (account, accounting, analytics, billing, business, customer, inventory, order, onboarding). Standard files per domain: `model.go`, `storage.go`, `service.go`, `errors.go`, `handler_http.go`, optional `middleware_http.go` and `state_machine.go`.

Run and configure

- Local dev: `make dev.server` (requires `air` tool installed). Air watches for changes, builds to `./tmp/main`, and runs `kyora server` with live reload.
- Config: Viper loads `.kyora.yaml` + environment variables. **Source of truth**: `internal/platform/config/config.go` defines all config keys as constants.
- Key config categories:
  - App: `app.name`, `app.port`, `app.domain`, `app.notifications_email`
  - HTTP: `http.port`, `http.base_url`, `http.trace_id_header` (default `X-Trace-ID`)
  - Database: `database.dsn`, `database.max_open_conns`, `database.max_idle_conns`, `database.max_idle_time`, `database.log_level`
  - Cache: `cache.hosts` (memcached servers)
  - JWT: `auth.jwt.secret`, `auth.jwt.expiry_seconds`, `auth.jwt.issuer`, `auth.jwt.audience`
  - Auth tokens: `auth.password_reset_ttl_seconds`, `auth.verify_email_ttl_seconds`, `auth.invitation_token_ttl_seconds`
  - Google OAuth: `auth.google_oauth.client_id`, `auth.google_oauth.client_secret`, `auth.google_oauth.redirect_url`
  - Stripe: `billing.stripe.api_key`, `billing.stripe.webhook_secret`
  - Email: `email.provider` (resend/mock), `email.resend.api_key`, `email.resend.base_url`, `email.from_email`, `email.from_name`, `email.support_email`, `email.help_url`
  - Logging: `log.format`, `log.level`
- Data stores: Postgres (GORM), Memcached (gomemcache). Stripe Billing and Google OAuth are opt-in via config.

Core patterns you must follow

**Database access:**

- Use `database.Database.Conn(ctx)` for transaction-aware connections. Context may contain transaction from `TxKey`.
- Use generic `database.Repository[T]` pattern. Each domain creates typed repos in `storage.go` via `database.NewRepository[Model](db)`.
- Repositories auto-migrate models on construction.
- **Scopes** (chainable filters): `ScopeBusinessID`, `ScopeWorkspaceID`, `ScopeID`, `ScopeIDs`, `ScopeEquals`, `ScopeIn`, `ScopeNotIn`, `ScopeGreaterThan`, `ScopeLessThan`, `ScopeBetween`, `ScopeTime`, `ScopeCreatedAt`, `ScopeSearchTerm`, `ScopeIsNull`.
- **Query modifiers**: `WithPreload`, `WithJoins`, `WithPagination`, `WithOrderBy`, `WithLimit`, `WithReturning`, `WithLockingStrength` (UPDATE, SHARE, SKIP LOCKED, NOWAIT).
- **CRUD**: `CreateOne`, `CreateMany`, `UpdateOne`, `UpdateMany`, `DeleteOne`, `DeleteMany`, `FindByID`, `FindOne`, `FindMany`, `Count`.
- **Aggregates**: `Sum`, `Avg`, `CountBy`, `SumBy`, `AvgBy`, `TimeSeriesSum`, `TimeSeriesCount`.

**Transactions:**

- Multi-step writes must use `database.AtomicProcess.Exec(ctx, func(tctx context.Context) error {...}, opts...)`.
- Transaction is injected into context as `TxKey`, so all repo calls inside the closure use the same transaction.
- Options: `atomic.WithIsolationLevel(atomic.LevelSerializable)`, `atomic.WithRetries(3)`, `atomic.WithReadOnly(true)`.
- See `internal/platform/database/atomic.go` and `internal/platform/types/atomic/atomic.go`.

**HTTP middleware & auth:**

- Middleware chain: `logger.Middleware()` logs requests/responses and adds `traceId` (header: config `http.trace_id_header`, default `X-Trace-ID`).
- Auth flow (workspace-based):
  1. `auth.EnforceAuthentication` - validates JWT cookie (name: `jwt`), sets claims in context
  2. `account.EnforceValidActor(accountService)` - loads user from claims, sets `ActorKey` in context, enriches logger
  3. `account.EnforceWorkspaceMembership(accountService)` - validates `workspaceId` param matches user's workspace, sets `WorkspaceKey` in context
  4. `account.EnforceActorPermissions(action, resource)` - validates user role has permission
  5. Optional: `billing.EnforceActiveSubscription(billingService)` - ensures workspace has active subscription
  6. Optional: `billing.EnforcePlanWorkspaceLimits(planLimit, counterFunc)` - checks plan limits before operations
- Extract from context: `account.ActorFromContext(c)`, `account.WorkspaceFromContext(c)`, `auth.ClaimsFromContext(c)`
- JWT helpers in `internal/platform/auth/jwt.go`. Note: `JwtFromContext` reads cookie only, not `Authorization` header.

**Responses & errors:**

- Success: `response.SuccessJSON(c, status, data)`, `response.SuccessText(c, status, text)`, `response.SuccessEmpty(c, status)`
- Errors: `response.Error(c, err)` - automatically converts to RFC 7807 Problem JSON
- Problem types in `internal/platform/types/problem`: `BadRequest()`, `Unauthorized()`, `Forbidden()`, `NotFound()`, `Conflict()`, `InternalError()`, `ValidationError()`
- DB error normalization: `database.IsRecordNotFound(err)`, `database.IsUniqueViolation(err)`
- Domain errors: Each domain defines errors in `errors.go` using problem constructors (see `internal/domain/order/errors.go`)

**Request validation:**

- Use `request.ValidBody(c, &req)` to bind and validate JSON request bodies

Domain conventions (copy when adding modules)

**model.go:**

- GORM models with `gorm.Model` embedded (adds ID, CreatedAt, UpdatedAt, DeletedAt)
- Override ID with string type: `ID string \`gorm:"column:id;primaryKey;type:text"\``
- Generate IDs in `BeforeCreate` hook using `id.KsuidWithPrefix("prefix")` or `id.Base62(length)`
- Define constants: `TableName`, `StructName`, `Prefix`
- Define `Schema` struct with `schema.Field` for each column: `schema.NewField("db_column", "jsonField")`
- Schema maps DB columns → JSON field names for ordering, searching, and aggregations
- Example: `OrderSchema.Total` = `schema.NewField("total", "total")`

**storage.go:**

- Create typed repositories: `order *database.Repository[Order]`, `orderItem *database.Repository[OrderItem]`
- Initialize in constructor: `database.NewRepository[Order](db)` (auto-migrates on construction)
- Store cache reference if needed: `cache *cache.Cache`
- Return struct with all repositories

**service.go:**

- Methods signature: `(ctx context.Context, actor *account.User, workspace *account.Workspace, ...)`
- **Always scope by WorkspaceID** for multi-tenancy: `s.storage.order.FindMany(ctx, s.storage.order.ScopeWorkspaceID(workspace.ID), ...)`
- Use `AtomicProcess.Exec` for multi-entity operations
- Return domain errors from `errors.go`
- Business logic patterns:
  - Validation → Load related entities → Calculate derived values → Persist → Return
  - For updates: Load with locking → Validate state transitions → Update → Persist
- See `internal/domain/order/service.go` for comprehensive examples

**errors.go:**

- Define domain-specific errors using problem constructors
- Add context with `.With(key, value)` and `.WithError(err)`
- Example: `problem.NotFound("order not found").WithError(err).With("orderId", id)`

**handler_http.go:**

- Gin handlers that extract actor/workspace from context
- Validate request body with `request.ValidBody`
- Call service methods
- Return responses via `response.Success*` or `response.Error`

**middleware_http.go (optional):**

- Domain-specific middleware (e.g., enforce workspace membership, check permissions)
- See `internal/domain/account/middleware_http.go` for patterns

**state_machine.go (optional):**

- Encode valid state transitions (e.g., order status flow)
- Update timestamp fields on transitions
- See `internal/domain/order/state_machine.go`

Lists & analytics

**Pagination & sorting:**

- Accept `*list.ListRequest` with page, pageSize, orderBy (JSON field names with optional `-` prefix for DESC)
- Create: `list.NewListRequest(page, pageSize, orderBy, searchTerm)`
- Parse ordering: `req.ParsedOrderBy(DomainSchema)` converts JSON fields to SQL: `["total DESC", "createdAt ASC"]`
- Apply to repo: `s.storage.order.WithOrderBy(req.ParsedOrderBy(OrderSchema))`
- Return: `list.NewListResponse(items, page, pageSize, totalCount, hasMore)`

**Time series:**

- Repo helpers: `TimeSeriesSum(ctx, valueColumn, timeColumn, granularity, opts...)`, `TimeSeriesCount(ctx, timeColumn, granularity, opts...)`
- Granularity from `internal/platform/types/timeseries`: `GranularityHour`, `GranularityDay`, `GranularityWeek`, `GranularityMonth`
- Auto-determines granularity: `timeseries.GetTimeGranularityByDateRange(from, to)`
- Returns `*timeseries.TimeSeries` with bucketed data points

**Aggregations:**

- Group by: `CountBy`, `SumBy`, `AvgBy` return `[]keyvalue.KeyValue`
- Direct: `Sum`, `Avg`, `Count` return scalar values
- Scoping with time ranges: `ScopeTime(field, from, to)`

Events & integrations

**Event bus:**

- In-memory pub/sub: `internal/platform/bus`
- Emit: `bus.Emit(topic, payload)` (non-blocking, async)
- Listen: `unsub := bus.Listen(topic, func(payload any) {...})` (returns unsubscribe function)
- Built-in topics in `internal/platform/bus/events.go`: `VerifyEmailTopic`, `ResetPasswordTopic`, `WorkspaceInvitationTopic`
- Graceful shutdown: `bus.Close()` waits for handlers to finish

**Billing/Stripe:**

- Stripe SDK initialized in `internal/server/server.go` with API key from config
- Domain: `internal/domain/billing` (plans, subscriptions, invoices, webhooks, customer management)
- Services: `service.go`, `service_customer.go`, `service_invoices.go`, `service_subscription.go`, `service_usage_tax.go`
- Webhooks: `webhooks.go` handles Stripe events, `stripe_retry.go` for retry logic
- Plan limits enforcement: middleware in `middleware_http.go`
- Config: `billing.stripe.api_key`, `billing.stripe.webhook_secret`

**Email:**

- Platform: `internal/platform/email`
- Providers: `resend` (production) or `mock` (development/testing)
- Config: `email.provider`, `email.resend.api_key`, `email.from_email`, `email.from_name`
- Send template: `client.SendTemplate(ctx, templateID, to, from, subject, data)`
- Templates: HTML files in `internal/platform/email/templates/` embedded with `go:embed`
- Available templates: `TemplateForgotPassword`, `TemplateEmailVerification`, `TemplateWelcome`, `TemplateSubscriptionConfirmed`, etc.
- See `internal/platform/email/README.md` for full template list and usage

**OAuth:**

- Google OAuth flow in `internal/platform/auth/google_oauth.go`
- Config: `auth.google_oauth.client_id`, `auth.google_oauth.client_secret`, `auth.google_oauth.redirect_url`
- JWT creation/parsing in `internal/platform/auth/jwt.go`
- Cookie name: `jwt`

Workspace-based multi-tenancy

- **Workspace model**: Groups users under a single subscription/billing entity
- **User model**: Belongs to one workspace (`WorkspaceID`), has a role
- **Business model**: Deprecated/transitioning to workspace (some code still uses `BusinessID`)
- **Scoping**: Always filter by `WorkspaceID` in queries for data isolation
- **Invitations**: Users can be invited to workspaces with a role (admin/member)
- **Role permissions**: Defined in `internal/platform/types/role` (actions: view/manage, resources: account/billing/orders/etc.)
- **Middleware**: `EnforceWorkspaceMembership` ensures user belongs to workspace in URL param

Gotchas & best practices

- **Multi-tenancy**: ALWAYS scope queries by `WorkspaceID` (or legacy `BusinessID`). Never return data across workspaces.
- **Config constants**: Use constants from `internal/platform/config/config.go`, never hardcode config keys.
- **JWT source**: `JwtFromContext` reads only the cookie named `jwt`, not `Authorization` header.
- **Ordering/search**: Use JSON field names (from Schema) in requests, not DB column names. Schema handles mapping.
- **ID generation**: Prefer KSUID with prefix for external IDs (`id.KsuidWithPrefix("ord")`), Base62 for short codes (`id.Base62(6)`).
- **Transactions**: Use atomic processor for multi-step writes. Don't manually begin/commit transactions.
- **Error handling**: Return domain errors (RFC 7807 Problems). Use `database.IsRecordNotFound/IsUniqueViolation` for DB errors.
- **Decimal math**: Use `shopspring/decimal` for financial calculations. Never use float64 for money.
- **Time zones**: Store timestamps in UTC. GORM uses `timestamptz` by default.
- **Logging**: Use structured logging via `slog`. Logger is context-aware and enriched by middleware.
- **State machines**: For entities with state transitions (orders, subscriptions), define allowed transitions and timestamp updates.

Pointers (reference implementations)

- **Repository pattern**: `internal/platform/database/repository.go` (all scopes and methods)
- **Transactions**: `internal/platform/database/atomic.go`, `internal/platform/types/atomic/atomic.go`
- **Service layer**: `internal/domain/order/service.go` (comprehensive CRUD, aggregations, time series)
- **HTTP handlers**: `internal/domain/order/handler_http.go`, `internal/domain/billing/handler_http.go`
- **Middleware**: `internal/domain/account/middleware_http.go`, `internal/domain/billing/middleware_http.go`
- **Errors**: `internal/domain/order/errors.go` (domain-specific problem definitions)
- **Models & schemas**: `internal/domain/order/model.go` (GORM models, schemas, DTOs)
- **State machines**: `internal/domain/order/state_machine.go` (status transitions)
- **Routes**: `internal/server/routes.go` (middleware chains, grouping, permissions)
- **Server setup**: `internal/server/server.go` (DI, service initialization, graceful shutdown)
- **Response helpers**: `internal/platform/response/response.go`
- **Request helpers**: `internal/platform/request/valid_body.go`
- **Logger middleware**: `internal/platform/logger/middleware.go`
- **List types**: `internal/platform/types/list/list_request.go`, `internal/platform/types/list/list_response.go`
- **Time series**: `internal/platform/types/timeseries/timeseries.go`
- **Problems**: `internal/platform/types/problem/problem.go`
