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

## Monorepo Structure

Kyora is organized as a **monorepo** to support multiple projects (backend, frontend, mobile, etc.):

```
kyora/
├── backend/              # Go backend API server (current focus)
│   ├── cmd/             # CLI commands
│   ├── internal/        # Internal packages
│   ├── main.go          # Entry point
│   ├── go.mod           # Go dependencies
│   ├── .air.toml        # Hot reload config
│   └── .kyora.yaml      # Backend configuration
├── Makefile             # Root-level build commands
├── README.md            # Monorepo overview
└── STRUCTURE.md         # Monorepo guidelines
```

**Important**: All Go backend code is in the `backend/` directory. When referencing files or paths, always include the `backend/` prefix.

## Backend Architecture

- Go 1.25.3 monolith using Cobra CLI framework. Entry point: `backend/main.go` → `backend/cmd/root.go`. HTTP server runs via `server` subcommand in `backend/cmd/server.go` which delegates to `backend/internal/server/server.go`.
- Go dependencies: Gin web framework, GORM ORM, Viper config, Stripe SDK, gomemcache, JWT, OAuth2, and more (see `backend/go.mod`).
- Layered `backend/internal/` architecture:
  - **platform**: Infrastructure layer (config, database, cache, logging, auth, event bus, request/response helpers, shared types).
  - **domain**: Business logic modules (account, accounting, analytics, billing, business, customer, inventory, order, onboarding). Standard files per domain: `model.go`, `storage.go`, `service.go`, `errors.go`, `handler_http.go`, optional `middleware_http.go` and `state_machine.go`.

Run and configure

- Local dev: `make dev.server` from root (or `cd backend && air server`). Requires `air` tool installed. Air watches for changes, builds to `backend/tmp/main`, and runs `kyora server` with live reload.
- Config: Viper loads `backend/.kyora.yaml` + environment variables. **Source of truth**: `backend/internal/platform/config/config.go` defines all config keys as constants.
- Example config: Copy `backend/.kyora.yaml.example` to `backend/.kyora.yaml` and customize for local development.
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
- See `backend/internal/platform/database/atomic.go` and `backend/internal/platform/types/atomic/atomic.go`.

**HTTP middleware & auth:**

- Middleware chain: `logger.Middleware()` logs requests/responses and adds `traceId` (header: config `http.trace_id_header`, default `X-Trace-ID`).
- Auth flow (workspace-based):
  1. `auth.EnforceAuthentication` - validates JWT token from Authorization header (format: `Bearer <token>`), sets claims in context
  2. `account.EnforceValidActor(accountService)` - loads user from claims, sets `ActorKey` in context, enriches logger
  3. `account.EnforceWorkspaceMembership(accountService)` - validates `workspaceId` param matches user's workspace, sets `WorkspaceKey` in context
  4. `account.EnforceActorPermissions(action, resource)` - validates user role has permission
  5. Optional: `billing.EnforceActiveSubscription(billingService)` - ensures workspace has active subscription
  6. Optional: `billing.EnforcePlanWorkspaceLimits(planLimit, counterFunc)` - checks plan limits before operations
- Extract from context: `account.ActorFromContext(c)`, `account.WorkspaceFromContext(c)`, `auth.ClaimsFromContext(c)`
- JWT helpers in `backend/internal/platform/auth/jwt.go`. `JwtFromContext` extracts token from Authorization header only.

**Responses & errors:**

- Success: `response.SuccessJSON(c, status, data)`, `response.SuccessText(c, status, text)`, `response.SuccessEmpty(c, status)`
- Errors: `response.Error(c, err)` - automatically converts to RFC 7807 Problem JSON
- Problem types in `backend/internal/platform/types/problem`: `BadRequest()`, `Unauthorized()`, `Forbidden()`, `NotFound()`, `Conflict()`, `InternalError()`, `ValidationError()`
- DB error normalization: `database.IsRecordNotFound(err)`, `database.IsUniqueViolation(err)`
- Domain errors: Each domain defines errors in `errors.go` using problem constructors (see `backend/internal/domain/order/errors.go`)

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
- See `backend/internal/domain/order/service.go` for comprehensive examples

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
- See `backend/internal/domain/account/middleware_http.go` for patterns

**state_machine.go (optional):**

- Encode valid state transitions (e.g., order status flow)
- Update timestamp fields on transitions
- See `backend/internal/domain/order/state_machine.go`

Lists & analytics

**Pagination & sorting:**

- Accept `*list.ListRequest` with page, pageSize, orderBy (JSON field names with optional `-` prefix for DESC)
- Create: `list.NewListRequest(page, pageSize, orderBy, searchTerm)`
- Parse ordering: `req.ParsedOrderBy(DomainSchema)` converts JSON fields to SQL: `["total DESC", "createdAt ASC"]`
- Apply to repo: `s.storage.order.WithOrderBy(req.ParsedOrderBy(OrderSchema))`
- Return: `list.NewListResponse(items, page, pageSize, totalCount, hasMore)`

**Time series:**

- Repo helpers: `TimeSeriesSum(ctx, valueColumn, timeColumn, granularity, opts...)`, `TimeSeriesCount(ctx, timeColumn, granularity, opts...)`
- Granularity from `backend/internal/platform/types/timeseries`: `GranularityHour`, `GranularityDay`, `GranularityWeek`, `GranularityMonth`
- Auto-determines granularity: `timeseries.GetTimeGranularityByDateRange(from, to)`
- Returns `*timeseries.TimeSeries` with bucketed data points

**Aggregations:**

- Group by: `CountBy`, `SumBy`, `AvgBy` return `[]keyvalue.KeyValue`
- Direct: `Sum`, `Avg`, `Count` return scalar values
- Scoping with time ranges: `ScopeTime(field, from, to)`

Events & integrations

**Event bus:**

- In-memory pub/sub: `backend/internal/platform/bus`
- Emit: `bus.Emit(topic, payload)` (non-blocking, async)
- Listen: `unsub := bus.Listen(topic, func(payload any) {...})` (returns unsubscribe function)
- Built-in topics in `backend/internal/platform/bus/events.go`: `VerifyEmailTopic`, `ResetPasswordTopic`, `WorkspaceInvitationTopic`
- Graceful shutdown: `bus.Close()` waits for handlers to finish

**Billing/Stripe:**

- Stripe SDK initialized in `backend/internal/server/server.go` with API key from config
- Domain: `backend/internal/domain/billing` (plans, subscriptions, invoices, webhooks, customer management)
- Services: `service.go`, `service_customer.go`, `service_invoices.go`, `service_subscription.go`, `service_usage_tax.go`
- Webhooks: `webhooks.go` handles Stripe events, `stripe_retry.go` for retry logic
- Plan limits enforcement: middleware in `middleware_http.go`
- Config: `billing.stripe.api_key`, `billing.stripe.webhook_secret`

**Email:**

- Platform: `backend/internal/platform/email`
- Providers: `resend` (production) or `mock` (development/testing)
- Config: `email.provider`, `email.resend.api_key`, `email.from_email`, `email.from_name`
- Send template: `client.SendTemplate(ctx, templateID, to, from, subject, data)`
- Templates: HTML files in `backend/internal/platform/email/templates/` embedded with `go:embed`
- Available templates: `TemplateForgotPassword`, `TemplateEmailVerification`, `TemplateWelcome`, `TemplateSubscriptionConfirmed`, etc.
- See `backend/internal/platform/email/README.md` for full template list and usage

**OAuth:**

- Google OAuth flow in `backend/internal/platform/auth/google_oauth.go`
- Config: `auth.google_oauth.client_id`, `auth.google_oauth.client_secret`, `auth.google_oauth.redirect_url`
- JWT creation/parsing in `backend/internal/platform/auth/jwt.go`

**Workspace-based multi-tenancy**

- **Workspace model**: Groups users under a single subscription/billing entity
- **User model**: Belongs to one workspace (`WorkspaceID`), has a role
- **Business model**: Deprecated/transitioning to workspace (some code still uses `BusinessID`)
- **Scoping**: Always filter by `WorkspaceID` in queries for data isolation
- **Invitations**: Users can be invited to workspaces with a role (admin/member)
- **Role permissions**: Defined in `backend/internal/platform/types/role` (actions: view/manage, resources: account/billing/orders/etc.)
- **Middleware**: `EnforceWorkspaceMembership` ensures user belongs to workspace in URL param

**Gotchas & best practices**

- **Multi-tenancy**: ALWAYS scope queries by `WorkspaceID` (or legacy `BusinessID`). Never return data across workspaces.
- **Config constants**: Use constants from `backend/internal/platform/config/config.go`, never hardcode config keys.
- **JWT source**: `JwtFromContext` extracts JWT token from the Authorization header (format: `Bearer <token>`).
- **Ordering/search**: Use JSON field names (from Schema) in requests, not DB column names. Schema handles mapping.
- **ID generation**: Prefer KSUID with prefix for external IDs (`id.KsuidWithPrefix("ord")`), Base62 for short codes (`id.Base62(6)`).
- **Transactions**: Use atomic processor for multi-step writes. Don't manually begin/commit transactions.
- **Error handling**: Return domain errors (RFC 7807 Problems). Use `database.IsRecordNotFound/IsUniqueViolation` for DB errors.
- **Decimal math**: Use `shopspring/decimal` for financial calculations. Never use float64 for money.
- **Time zones**: Store timestamps in UTC. GORM uses `timestamptz` by default.
- **Logging**: Use structured logging via `slog`. Logger is context-aware and enriched by middleware.
- **State machines**: For entities with state transitions (orders, subscriptions), define allowed transitions and timestamp updates.

**Pointers (reference implementations)**

- **Repository pattern**: `backend/internal/platform/database/repository.go` (all scopes and methods)
- **Transactions**: `backend/internal/platform/database/atomic.go`, `backend/internal/platform/types/atomic/atomic.go`
- **Service layer**: `backend/internal/domain/order/service.go` (comprehensive CRUD, aggregations, time series)
- **HTTP handlers**: `backend/internal/domain/order/handler_http.go`, `backend/internal/domain/billing/handler_http.go`
- **Middleware**: `backend/internal/domain/account/middleware_http.go`, `backend/internal/domain/billing/middleware_http.go`
- **Errors**: `backend/internal/domain/order/errors.go` (domain-specific problem definitions)
- **Models & schemas**: `backend/internal/domain/order/model.go` (GORM models, schemas, DTOs)
- **State machines**: `backend/internal/domain/order/state_machine.go` (status transitions)
- **Routes**: `backend/internal/server/routes.go` (middleware chains, grouping, permissions)
- **Server setup**: `backend/internal/server/server.go` (DI, service initialization, graceful shutdown)
- **Response helpers**: `backend/internal/platform/response/response.go`
- **Request helpers**: `backend/internal/platform/request/valid_body.go`
- **Logger middleware**: `backend/internal/platform/logger/middleware.go`
- **List types**: `backend/internal/platform/types/list/list_request.go`, `backend/internal/platform/types/list/list_response.go`
- **Time series**: `backend/internal/platform/types/timeseries/timeseries.go`
- **Problems**: `backend/internal/platform/types/problem/problem.go`

**HTTP handlers and middleware and routes and permissions and limits**

- For every domain we have httpHandler and we might have other API interfaces like grpc or kafka events or internal event bus events handling, we should keep the same structure and patterns as the httpHandler for consistency.
- we should always define the routes for every domain in the backend/internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
- we should always define the middleware chain for every domain routes in the backend/internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
- we should always enforce the actor permissions for every domain routes in the backend/internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
- we should always enforce the plan workspace limits for every domain routes in the backend/internal/server/routes.go file and keep the same patterns and structure as other domain routes for consistency.
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
- we should use the cache package from backend/internal/platform/cache for caching data. which is usually a wrapper around memcached. and accessible in every service domain by the service storage (ex: s.storage.cache).
- we should always set an appropriate expiration time for cached data to ensure that the data is fresh and up to date.
- we should always follow a cache invalidation strategy to ensure that cached data is invalidated when the underlying data changes.
- we should always use cache keys that are unique and descriptive to avoid collisions and make it easy to identify cached data.
- for data retrieval caching we should always follow the read-through caching pattern to ensure that data is cached on first access and subsequent accesses are served from the cache with proper invalidation strategy
- caching should be used judiciously and only for data that is frequently accessed and doesn't change frequently to avoid cache staleness and unnecessary complexity.
- caching should be done in storage layer to keep the service layer clean and focused on business logic. never ever do caching in service layer.
- in service layer we can only use the cache for business logic that requires expirable memory storage like authentication tokens and password reset tokens and email verification tokens and invitation tokens and similar use cases.

**Seperation of concerns**

- handler layer should should only handle requests and responses and request validation and maybe some transformations to match the service interface and delegate the business logic to the service layer.
- service layer should only handle business logic and data transformations and delegate the data access to the storage layer and should never ever handle request or response or request validation.
- storage layer should only handle data access and data transformations and should never ever handle business logic or request or response or request validation.
- for models they might have some methods that are related to the model itself like state transitions or derived values calculations but they should never ever have business logic or data access or request or response handling.
- for errors they should only define the domain specific errors.
- for every domain it should never have circular dependencies with other domains. if you find yourself needing to use another domain service or storage you should always use event bus to communicate between domains or you can move the shared functionality to a common utils package or helpers package.
- domains should be independent and self contained and should not depend on other domains to function.
- domains should never access other domain storages directly. if you need to access another domain storage you should always use the other domain service to access the data.

**Test cases**

- In this Project, we should always reach 80%+ test coverage for any new code we add. the code currently doesn't have comprehensive test coverage so we should add tests as we develop new features.
- we should always use `github.com/stretchr/testify` for structuring our test cases and assertions with testify suite.

**Test Organization:**

- **Unit tests**: Domain-specific unit tests can be placed within the domain folder (e.g., `backend/internal/domain/order/service_test.go`) or platform folder for testing individual components in isolation. Currently, we don't have unit tests but they should follow this pattern when added.
- **Integration/E2E tests**: All integration and end-to-end tests are located in `backend/internal/tests/e2e/` directory, separate from domain code. This ensures clean separation between production code and integration tests.
- **Test file structure**: Tests are organized by domain with each API route/endpoint having its own testify suite within domain-specific test files (e.g., `backend/internal/tests/e2e/account_login_test.go`, `backend/internal/tests/e2e/account_workspace_test.go`, `backend/internal/tests/e2e/order_create_test.go`). This provides granular test isolation and clear organization.
- **Test utilities**: Common test helpers, fixtures, and setup functions are in `backend/internal/tests/testutils/` package for reuse across all test suites.
- **Mocks**: All mocks for external dependencies are placed in `backend/internal/tests/mocks/` package to keep them centralized and avoid duplication.

**E2E Test Infrastructure:**

- E2E tests use testcontainers to spin up ephemeral Docker containers (Postgres, Memcached, Stripe-mock) for complete isolation.
- `backend/internal/tests/e2e/main_test.go` contains `TestMain` which sets up the global test environment, starts containers, initializes the server, and tears down after all tests complete.
- Each API route/endpoint has its own testify suite grouped by domain (e.g., `LoginSuite` in `account_login_test.go`, `CreateOrderSuite` in `order_create_test.go`). Suite names should clearly indicate the functionality being tested.
- Test suites access shared resources via global variables: `testEnv` (containers) and `testServer` (HTTP server instance).
- Server runs on isolated port (18080) with mock email provider and test configuration.
- Testutils provides context-based container initialization functions (`CreateDatabaseCtx`, `CreateCacheCtx`, `CreateStripeMockCtx`) and aggregated `InitEnvironment` for TestMain.
- Each test suite MUST clean database tables in `SetupTest()` and `TearDownTest()` hooks to ensure complete isolation between tests. Use testutils helper functions like `ClearTable(db, "users", "workspaces", "orders")` for consistent cleanup.

**Test Best Practices:**

- test cases should cover both positive and negative scenarios to ensure robustness and should validate edge cases as well and ensure error handling works as expected.
- use **table-driven tests** extensively for testing multiple scenarios, input variations, and edge cases. Each table entry should represent a complete test case with inputs, expected outputs, and validation logic.
- always aim for fast-running tests to ensure quick feedback during development and CI.
- use **fuzzing tests** (Go's native fuzzing with `testing.F`) for functions that process user input, parse data, or handle untrusted inputs to discover edge cases and potential bugs.
- monitor **test coverage** using `go test -cover` and aim for 80%+ coverage. Use `go test -coverprofile=coverage.out` to generate detailed coverage reports.
- document test cases to explain the purpose of each test and any setup required to run them.
- review and update test cases regularly to ensure they remain relevant as the codebase evolves.
- strive for high-quality tests that are reliable and maintainable.
- test files should have `_test.go` suffix and use package name suffixed with `_test` to avoid circular dependencies (e.g., `package e2e_test`).
- use descriptive names for test functions to clearly indicate what is being tested.
- **test isolation is critical**: Each test MUST be completely independent and not rely on the state of other tests. Use `SetupTest()` to clear database tables and set up fresh fixtures, and `TearDownTest()` to clean up after each test.
- **database cleanup pattern**: In `SetupTest()`, truncate all relevant tables before setting up test data. In `TearDownTest()`, truncate tables again to ensure clean state. Use testutils helpers like `testutils.TruncateTables(testEnv.Database, "users", "workspaces", "orders")`.
- use assertions to validate expected outcomes and avoid using print statements for validation.
- test cases should be reviewed as part of code reviews to ensure quality and coverage.
- test cases should be included in the definition of done for any new feature or bug fix.
- test cases should be written with the mindset of future maintainers to ensure they are easy to understand and maintain over time.
- test files should not be very long (maximum 500 lines) to keep them maintainable. Split into smaller files based on functionality if needed (e.g., separate files per API endpoint).
- don't let main code files have any logic related to test cases. All test logic belongs in test files, testutils, and mocks.

**Running Tests:**

**Make commands (recommended):**

- `make test` - run all backend tests with verbose output (from repository root)
- `make test.unit` - run only backend unit tests (domain and platform packages)
- `make test.e2e` - run backend E2E tests with 120s timeout
- `make test.quick` - run all backend tests without verbose output (faster feedback)
- `make test.coverage` - run backend tests with coverage report and summary
- `make test.coverage.html` - generate HTML coverage report (`backend/coverage.html`)
- `make test.coverage.view` - generate and open HTML coverage report in browser
- `make test.e2e.coverage` - run backend E2E tests with coverage reporting
- `make clean.coverage` - remove all backend coverage report files
- `make help` - display all available Makefile commands

**Go test commands (alternative - run from backend/ directory):**

- `cd backend && go test ./...` - run all tests
- `cd backend && go test ./internal/tests/e2e -v` - run e2e tests
- `cd backend && go test ./... -cover -coverprofile=coverage.out` - run with coverage
- `cd backend && go test ./internal/tests/e2e -v -run TestLoginSuite` - run specific suite
- `cd backend && go test ./internal/tests/e2e/account_login_test.go -v` - run specific test file
- `cd backend && go test ./internal/tests/e2e -race` - run with race detection
- `cd backend && go test -fuzz=FuzzFunctionName -fuzztime=30s` - run fuzzing tests

**Requirements:**

- tests require Docker Desktop running for testcontainers

**Test Best Practices:**

1. **Use Domain Storage Layers - NEVER Raw SQL**

   - ALWAYS use domain storage repositories for database operations in tests
   - Initialize storage layers properly: `onboarding.NewStorage(db, cache)`, `account.NewStorage(db, cache)`, etc.
   - Follow the same CRUD patterns as production code: `storage.GetByToken(ctx, token)`, `storage.UpdateSession(ctx, sess)`
   - Use proper domain models (structs) instead of maps or raw SQL
   - Example:

     ```go
     // ❌ BAD - Raw SQL
     db.Exec("UPDATE onboarding_sessions SET stage = ? WHERE token = ?", stage, token)

     // ✅ GOOD - Domain storage
     sess, _ := onboardingStorage.GetByToken(ctx, token)
     sess.Stage = onboarding.SessionStage(stage)
     onboardingStorage.UpdateSession(ctx, sess)
     ```

2. **Assert ALL Expected Response Fields**

   - Verify EVERY field in API responses matches expected values
   - Assert exact field count to catch unexpected additions/removals: `s.Len(response, 3, "response should have exactly 3 fields")`
   - Use `s.Contains()` to verify field presence: `s.Contains(response, "sessionToken")`
   - Assert exact values, not just presence: `s.Equal("plan_selected", response["stage"])`
   - For nested objects (like user in complete response), validate all nested fields
   - Example:

     ```go
     // Assert response structure
     s.Len(result, 2, "response should have exactly 2 fields")
     s.Contains(result, "user")
     s.Contains(result, "token")

     // Assert all user fields
     user := result["user"].(map[string]interface{})
     s.Equal("test@example.com", user["email"])
     s.Equal("John", user["firstName"])
     s.Equal("Doe", user["lastName"])
     s.Equal(true, user["isEmailVerified"])
     s.NotEmpty(user["id"])
     s.NotEmpty(user["workspaceId"])
     s.NotEmpty(user["role"])
     ```

3. **Test Isolation is Critical**

   - Each test MUST be completely independent
   - Use `SetupTest()` to clear database tables and set up fresh fixtures
   - Use `TearDownTest()` to clean up after each test
   - Never reuse session tokens, users, or data across tests
   - For table-driven tests with loops, create new session/user per iteration to avoid state conflicts
   - Use indexed emails for unique test data: `fmt.Sprintf("user%d@example.com", i)`

4. **Proper Time Handling**

   - Always use UTC for database timestamps: `time.Now().UTC()`
   - Be aware of timezone differences between local time and database storage
   - When setting expired times, use UTC: `time.Now().UTC().Add(-20 * time.Minute)`

5. **Context Usage**

   - Always use `context.Background()` for test operations
   - Pass context to all storage layer methods for consistency with production code

6. **Testing External Dependencies (OAuth, Stripe, etc.)**

   - External dependencies may not be configured in test environment
   - Tests should gracefully handle missing configuration by skipping or using flexible assertions
   - Example: Google OAuth tests skip when OAuth returns 500 (missing config)
   - Use flexible status code assertions (>= 400) when testing with external dependencies
   - Consider mocking external services for more reliable test execution

7. **Domain-Specific Helper Pattern**

   - Create `<domain>_helpers_test.go` for domain-specific test utilities
   - Keep generic helpers in `backend/internal/tests/testutils/` package
   - Domain helpers should provide: CreateTestEntity, GetEntity, UpdateEntity, specialized setup functions
   - Example: `AccountTestHelper` provides CreateTestUser, CreateInvitation, CreatePasswordResetToken, etc.
   - Initialize helpers in `SetupSuite()` and reuse across all tests in that domain

8. **Security Testing Checklist**

   - **Authentication**: Test missing/invalid/expired tokens
   - **Authorization**: Test permission boundaries (admin vs member)
   - **Workspace Isolation**: Verify users can't access other workspace data
   - **Input Validation**: Test SQL injection, XSS, malformed input
   - **Rate Limiting**: Test token reuse, expired tokens, multiple requests
   - **CSRF Protection**: Test state tokens in OAuth flows
   - **Enumeration Prevention**: Ensure consistent error messages for existing/non-existing resources

9. **Table-Driven Test Organization**

   - Use sub-tests with `s.Run()` for clear test output
   - Group related scenarios together (e.g., all invalid email formats)
   - Create fresh data per iteration to avoid state conflicts
   - Use descriptive test case names: `"missing email field"`, `"invalid email format"`, etc.

10. **Response Validation Best Practices**
    - Always assert exact field count to catch unexpected changes
    - Verify presence of all expected fields with `s.Contains()`
    - Assert exact values where possible, not just presence
    - For IDs/UUIDs: use `s.NotEmpty()` to verify generation
    - For timestamps: verify non-nil and reasonable ranges
    - For nested objects: validate all nested fields recursively

**Test Suite Example:**

Each API endpoint should have its own suite following this pattern:

```go
// File: backend/internal/tests/e2e/account_login_test.go
package e2e_test

import (
    "context"
    "net/http"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/abdelrahman146/kyora/internal/tests/testutils"
    "github.com/abdelrahman146/kyora/internal/domain/account"
)

// LoginSuite tests the POST /v1/auth/login endpoint
type LoginSuite struct {
    suite.Suite
    client         *testutils.HTTPClient
    accountStorage *account.Storage
}

func (s *LoginSuite) SetupSuite() {
    s.client = testutils.NewHTTPClient("http://localhost:18080")
    // Initialize storage for test data setup
    cache := cache.NewConnection([]string{"localhost:11211"})
    s.accountStorage = account.NewStorage(testEnv.Database, cache)
}

func (s *LoginSuite) SetupTest() {
    // Clear database tables before each test for isolation
    testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
}

func (s *LoginSuite) TearDownTest() {
    // Clean up after each test
    testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
}

func (s *LoginSuite) TestLogin_Success() {
    ctx := context.Background()

    // Create test user using storage layer (not raw SQL)
    user := &account.User{
        Email:    "user@example.com",
        Password: "hashedPassword",
        // ... other fields
    }
    s.NoError(s.accountStorage.CreateUser(ctx, user))

    // Make API request
    payload := map[string]interface{}{
        "email":    "user@example.com",
        "password": "ValidPass123!",
    }
    resp, err := s.client.Post("/v1/auth/login", payload)
    s.NoError(err)
    defer resp.Body.Close()
    s.Equal(http.StatusOK, resp.StatusCode)

    // Assert ALL response fields
    var result map[string]interface{}
    s.NoError(testutils.DecodeJSON(resp, &result))

    // Verify exact structure
    s.Len(result, 2, "response should have exactly 2 fields")
    s.Contains(result, "user")
    s.Contains(result, "token")

    // Verify all values
    s.NotEmpty(result["token"])
    user := result["user"].(map[string]interface{})
    s.Equal("user@example.com", user["email"])
    s.NotEmpty(user["id"])
}

func (s *LoginSuite) TestLogin_TableDriven() {
    tests := []struct {
        name           string
        email          string
        password       string
        expectedStatus int
    }{
        {"valid credentials", "user@example.com", "ValidPass123!", 200},
        {"wrong password", "user@example.com", "WrongPass!", 401},
        {"non-existent user", "nobody@example.com", "Pass123!", 404},
    }

    for _, tt := range tests {
        s.Run(tt.name, func() {
            // Test implementation
        })
    }
}

func TestLoginSuite(t *testing.T) {
    if testServer == nil {
        t.Skip("Test server not initialized")
    }
    suite.Run(t, new(LoginSuite))
}
```

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
- don't leave deprecated code in the codebase, we should always remove it to keep the code clean and maintainable.
- always do cleanups for no longer needed code or functions.
- never ever brief any implementation. always provide complete and thorough implementations.
- never ever settle on examples or partial implementations. always provide complete and thorough implementations.
- always aim for high-quality code that is secure, robust, maintainable, and production-ready
