Kyora — AI coding agent quickstart

## What is Kyora?

Kyora is a business management assistant specifically designed for solo social media entrepreneurs and small teams who sell products on social media platforms. It's built for business owners who are excellent at creating and selling their products but may not have expertise in accounting or business management. Kyora handles all the complex financial tracking and business operations behind the scenes, presenting everything in a simple, easy-to-understand way so owners can focus on what they do best: creating and selling.

### The Problem Kyora Solves

Social media business owners face a unique challenge: they're passionate about their products and skilled at selling on platforms like Instagram, Facebook, TikTok, and WhatsApp, but they struggle with:

- **Financial Confusion**: Not knowing if they're actually making profit or just generating revenue
- **Order Chaos**: Tracking orders scattered across DMs, comments, and multiple platforms
- **Inventory Blindness**: Not knowing what's in stock, what's selling, or when to reorder
- **Cash Flow Mystery**: Unclear about their actual financial position and available cash
- **Tax Anxiety**: Worried about missing important financial records for tax time
- **Growth Paralysis**: Unable to make data-driven decisions because they lack proper business insights

Kyora becomes their silent business partner that handles all this complexity automatically, without requiring them to become accounting experts or spend hours on administrative tasks.

### Core Value Proposition

**Simplicity First**: Every feature is designed to be intuitive and require zero accounting knowledge. Complex financial concepts are translated into simple, actionable insights.

**Automatic Heavy Lifting**: Kyora automatically tracks revenue recognition, calculates profitability, monitors inventory levels, and maintains complete financial records without manual bookkeeping.

**Social Media Native**: Understanding that orders come through DMs, comments, and chat apps, Kyora makes it easy to quickly log orders from any source and keep everything organized in one place.

**Peace of Mind**: Owners can sleep well knowing their business finances are properly tracked, they have records for tax purposes, and they truly understand their financial position.

### Key Features

- **Simple Order Management**: Quickly add orders from any social media platform, track their status from payment to delivery, and automatically recognize revenue.
- **Effortless Customer Tracking**: Automatically build a customer database from orders, see purchase history, and identify your best customers without manual data entry.
- **Intelligent Inventory**: Know what's in stock, get alerts when items are running low, and understand which products are your best sellers.
- **Clear Financial Picture**: See your actual profit (not just revenue), understand your cash flow, and get simple reports that show how your business is really doing.
- **Automated Accounting**: Revenue recognition, expense tracking, and financial reporting happen automatically in the background.
- **Business Insights in Plain English**: No confusing charts or accounting jargon—just clear answers like "You made $X profit this month" and "Your best-selling product is Y."

### Multi-Tenancy Model

Kyora uses a workspace-based multi-tenancy architecture where:

- Each **workspace** represents one business owner or small team with their own subscription, users, and data.
- **Users** belong to a single workspace and have role-based permissions (admin/member for team collaboration).
- All data is strictly isolated by workspace to ensure complete privacy and security.
- Billing and subscriptions are managed at the workspace level, with affordable pricing tiers designed for solo entrepreneurs and small teams.

### Target Users

- **Solo social media sellers**: Individuals selling handmade goods, fashion, beauty products, food, or services through Instagram, Facebook, TikTok, WhatsApp, etc.
- **Small social commerce teams**: 2-5 person teams managing a social media-based business together
- **Side hustlers**: People running a business alongside their day job who need dead-simple management tools
- **Product creators**: Artisans, designers, bakers, makers who want to focus on creation, not administration
- **Non-technical entrepreneurs**: Business owners who are intimidated by complex software and just want something that works

### Key Business Flows

1. **Quick Onboarding**: Sign up with email or Google, verify email, name your workspace, and start adding orders immediately—no complex setup required.
2. **Simple Order Entry**: Add an order in seconds (customer name, product, price, payment received) → Kyora automatically tracks it through delivery → automatically recognizes revenue and updates inventory.
3. **Automatic Financial Tracking**: As orders are added → revenue is recognized → inventory is adjusted → profit is calculated → financial position is updated in real-time.
4. **Instant Business Insights**: Open the dashboard → see profit this month, total revenue, best customers, top products—all in plain language with simple visuals.
5. **Team Collaboration (Optional)**: Invite a helper or partner → assign them as admin or member → they can help manage orders while you focus on production.
6. **Subscription Management**: Start with a free or basic plan → as business grows, upgrade to handle more orders → billing happens automatically through Stripe.

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

**Workspace-based multi-tenancy**

- **Workspace model**: Groups users under a single subscription/billing entity
- **User model**: Belongs to one workspace (`WorkspaceID`), has a role
- **Business model**: Deprecated/transitioning to workspace (some code still uses `BusinessID`)
- **Scoping**: Always filter by `WorkspaceID` in queries for data isolation
- **Invitations**: Users can be invited to workspaces with a role (admin/member)
- **Role permissions**: Defined in `internal/platform/types/role` (actions: view/manage, resources: account/billing/orders/etc.)
- **Middleware**: `EnforceWorkspaceMembership` ensures user belongs to workspace in URL param

**Gotchas & best practices**

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

**Pointers (reference implementations)**

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

**HTTP handlers and middleware and routes and permissions and limits**

- For every domain we have httpHandler and we might have other API interfaces like grpc or kafka events or internal event bus events handling, we should keep the same structure and patterns as the httpHandler for consistency.
- we should always define the routes for every domain in the internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
- we should always define the middleware chain for every domain routes in the internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
- we should always enforce the actor permissions for every domain routes in the internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
- we should always enforce the plan workspace limits for every domain routes in the internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
- we should always keep the http handlers concise and delegate the business logic to the domain services.
- we should always handle errors in the http handlers using the response.Error function to ensure consistent error handling across the codebase.
- we should always validate the request body in the http handlers using the request.ValidBody function to ensure consistent request validation across the codebase.
- we should always return the response in the http handlers using the response.SuccessJSON function to ensure consistent response formatting across the codebase.
- we should always keep the middleware concise and delegate the business logic to the domain services.

**OpenAPI docs**

- we should always document the http handlers using OpenAPI specs to ensure that the API is well documented and easy to understand for other developers.
- we will use the go package github.com/swaggo/swag to generate the OpenAPI docs from the code comments in the http handlers.
- the code comments in the http handlers should follow the swaggo format to ensure that the OpenAPI docs are generated correctly.
- the code comments should include information about the request parameters, request body, response body, and possible error responses. and it should be accurate and complete.
- we should always keep the OpenAPI docs up to date with the code changes.
- we should always check the generated OpenAPI docs to ensure that they are correct and complete and match the code.

**Caching (Memcached)**

- we should always use caching for frequently accessed data to improve performance and reduce database load.
- we should use the cache package from internal/platform/cache for caching data. which is usually a wrapper around memcached. and accessible in every service domain by the service storage (ex: s.storage.cache).
- we should always set an appropriate expiration time for cached data to ensure that the data is fresh and up to date.
- we should always follow a cache invalidation strategy to ensure that cached data is invalidated when the underlying data changes.
- we should always use cache keys that are unique and descriptive to avoid collisions and make it easy to identify cached data.
- for data retrieval caching we should always follow the read-through caching pattern to ensure that data is cached on first access and subsequent accesses are served from the cache with proper invalidation strategy
- caching should be used judiciously and only for data that is frequently accessed and doesn't change frequently to avoid cache staleness and unnecessary complexity.
- caching should be done in storage layer to keep the service layer clean and focused on business logic. never ever do caching in service layer.
- in service layer we can only use the cache for business logic that requires expirable memory storage like authentication tokens and password reset tokens and email verification tokens and invitation tokens and similar use cases.

  **Test cases**

- In this Project, we should always reach 80%+ test coverage for any new code we add. the code currently doesn't have any test cases so we should start adding them as we go.
- we should always use `github.com/stretchr/testify` for structuring our test cases and assertions with testify suite.
- every domain should have its own test files for service, storage, handler, and any other significant components.
- we should use test containers for setting up temporary databases for running our test cases to ensure isolation and consistency.
- we should mock external dependencies like email providers, payment gateways, and any other third-party services to avoid making real calls during tests.
- mocks should be shared across domains and placed in a common mocks package to avoid duplication.
- we should always run our test cases in CI to ensure that new changes don't break existing functionality.
- test cases should cover both positive and negative scenarios to ensure robustness and should validate edge cases as well and should ensure that error handling works as expected.
- we should use table-driven tests for functions with multiple input scenarios to keep the test code clean and maintainable.
- we should always aim for fast-running tests to ensure quick feedback during development and CI.
- we should use code coverage tools to monitor our test coverage and identify areas that need more tests.
- we should document our test cases to explain the purpose of each test and any setup required to run them.
- we should have integration tests that cover end-to-end scenarios to ensure that different components work together as expected.
- we should review and update our test cases regularly to ensure they remain relevant as the codebase evolves.
- we should strive for high-quality tests that are reliable and maintainable.

**General Instructions**

- the code is maintained by a single developer so we should always aim for simplicity and clarity in the code we write.
- the code is still under heavy development so we should always write code that is flexible and easy to change as requirements evolve and we can do breaking changes as needed.
- never leave any TODOs or FIXMEs in the code, we should always address them before finalizing the code.
- we should always follow the SOLID principles and best practices when writing code.
- we should always follow the existing code style and conventions used in the project to maintain consistency across the codebase.
- we should always write clear and concise comments and documentation for the code we write to ensure that it's easy to understand for other developers.
- we should always consider performance implications when writing code and optimize for efficiency where necessary.
- we should always consider security implications when writing code and ensure that the code is secure and follows best security practices.
- we should always consider scalability implications when writing code and ensure that the code can handle increased load and scale as needed.
- we should always consider maintainability implications when writing code and ensure that the code is easy to maintain and extend in the future
- the code should be clear and human maintainable and concise and follow best practices.
- when we have a generic sharable functionality we should add it in utils package either in helpers package or create its own package if its big enough.
- when creating domain services we should follow the same pattern as other domain services to keep consistency across the codebase.
- when handling errors we should use the problem types defined in the platform/types/problem package to have consistent error handling across the codebase.
- for logging we should always use the slog logger from context using the logger.FromContext(ctx) function, it is already enriched with useful information like trace id and actor info.
- for every domain we have httpHandler and we might have other API interfaces like grpc or kafka events or internal event bus events handling, we should keep the same structure and patterns as the httpHandler for consistency.
- whenever you find inefficient or duplicate code across domains we should refactor it into a shared utility function in the utils package or helpers package.
- whenever you find a code that doesn't follow best practices or has potential bugs we should fix it to follow best practices and avoid potential issues in the future.
- the code output should be always secure, robust, and production-ready and follows best practices and standards and should be 100% complete.
- we should alwasys look for smart simple and elegant solutions to complex problems leaving the code clean and maintainble and very easy to fix and extend in the future.
- whenever we update or create a function or struct or a field or a package we should always update or add a standard comment or documentation to explain its purpose and usage using godoc style.
- everytime we write code we should always think about the future and how this code will be used and maintained in the future and we should always write code that is easy to understand and maintain in the future.
- whenever we introduce new design pattern or archiectural changes we need to make sure it aligns with the overall architecture and design principles of the project and we should document the changes properly to explain the reasoning behind them.
- whenever you find a uncommented code or a function or a struct or a package without documentation we should always add a proper comment or documentation to explain its purpose and usage using godoc style.
