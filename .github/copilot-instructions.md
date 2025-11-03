Kyora codebase guide for AI coding agents

Purpose and shape

- Go monolith (Go 1.25) using Cobra CLI (entry: `main.go` -> `cmd/root.go`) and a layered internal package structure under `internal/`.
- Two major layers:
  - platform: cross-cutting infra (config, DB, cache, logging, request/response, auth, event bus, types).
  - domain: business modules (account, accounting, analytics, billing, business, customer, inventory, order). Each domain typically has `model.go`, `service.go`, `storage.go`, `errors.go` and sometimes a `state_machine.go`.

Run and config

- Local dev uses Air (live-reload). Command: `make dev.server` (requires `air` installed). Air builds `./tmp/main` and runs it with args; the Makefile passes `server` to the binary. If you add an HTTP server, implement a Cobra subcommand named `server` under `cmd/`.
- Configuration is loaded via Viper from `.kyora.yaml` at repo root and env vars (see keys in `internal/platform/config/config.go`). Example file exists in root; prefer the code constants over guessing key names.
- Data store: Postgres via GORM; memcached for cache. Stripe (billing) and Google OAuth are wired but opt‑in via config.

Core platform patterns

- Database access
  - `internal/platform/database.Database` wraps a GORM connection and exposes `Conn(ctx)` which picks up a transaction from context when present (key: `database.TxKey`).
  - Generic repository `database.Repository[T]` provides CRUD and common scopes: `ScopeBusinessID/WorkspaceID/Equals/In/Time/CreatedAt`, preloading (`WithPreload`), joins, ordering (`WithOrderBy`), pagination, `WithReturning`, aggregates (`Sum/Avg/Count`, time series helpers).
  - Transactions: use an `atomic.AtomicProcessor` (implemented by `database.AtomicProcess`). Wrap multi-step operations with `Exec(ctx, func(ctx) error, atomic.WithIsolationLevel(...), atomic.WithRetries(n))`. The transaction is injected into `ctx` so all repository calls in the closure are consistent.
- HTTP and errors
  - Logging middleware (`internal/platform/logger/middleware.go`) attaches `traceId` (header from `http.trace_id_header`, default `X-Trace-ID`) and logs JSON via slog; use `logger.FromContext(ctx)` inside services.
  - Auth uses JWT cookies (`jwt` cookie). `request.EnforceAuthentication` parses JWT and stores claims in Gin context; follow with `request.EnforceValidActor` and `request.EnforceBusinessValidity` to resolve `*account.User` and `*business.Business` from context.
  - Responses: use `response.SuccessJSON/SuccessText/SuccessEmpty` and `response.Error`. Errors are normalized to RFC 7807 Problem JSON (`internal/platform/types/problem`). DB errors are mapped (NotFound/Conflict) via `database.IsRecordNotFound/IsUniqueViolation`.
- Event bus
  - Lightweight async pub/sub in `internal/platform/bus` with `Bus.Emit(topic, payload)` and `Bus.Listen(topic, handler)`. Built-in topics: `VerifyEmailTopic`, `ResetPasswordTopic`.
- Lists and analytics
  - For pagination/sorting, pass `*list.ListRequest` through services; build DB ordering via `req.ParsedOrderBy(schemaDef)` where `schemaDef` is the domain `...Schema` (see `internal/platform/types/schema/field.go`). Wrap results with `list.NewListResponse`.
  - Time series helpers in `database.Repository` (`TimeSeriesSum/TimeSeriesCount`) plus label formatting in `internal/platform/types/timeseries`.

Domain conventions (copy when adding a module)

- model.go: GORM models with `gorm.Model` and string `ID` generated via `utils/id` (e.g., `id.KsuidWithPrefix("ord")`). Define a `...Schema` struct of `schema.Field` to map DB columns to JSON field names for ordering/search.
- storage.go: construct `database.Repository[T]` per aggregate; `NewStorage` wires repositories (and cache if needed). Repos auto‑migrate their model on construction.
- service.go: business methods have signature `(ctx, actor *account.User, biz *business.Business, ...)` to enforce tenancy. Use repository scopes like `ScopeBusinessID` in all reads/writes. Use atomic processor for multi-entity updates. Prefer returning domain-specific problems from `errors.go`.
- state_machine.go: encode allowed transitions (see `order/state_machine.go`) and update timestamp fields on transitions.

Integration points and examples

- Billing/Stripe: `internal/domain/billing` injects `*stripe.Client` and reads secrets from `billing.stripe.*` config keys.
- OAuth: `internal/platform/auth/google_oauth.go` for Google OAuth; JWT helpers in `internal/platform/auth/jwt.go` (cookie name `jwt`).
- Cache: `internal/platform/cache.Cache` (memcached). Use `Marshal/Unmarshal` helpers and `Increment` for counters.

Gotchas you’ll likely hit

- Config keys: treat `internal/platform/config/config.go` as the source of truth for key names. The sample `.kyora.yaml` is illustrative and may diverge; e.g., DB idle/connection timing keys.
- Multi‑tenancy: always scope by `BusinessID` (and/or `WorkspaceID`) using repository scopes; many services assume this contract.
- JWT source: `JwtFromContext` reads only the cookie, not the Authorization header.
- Ordering: `ListRequest.OrderBy()` expects JSON field names (e.g., `-createdAt`), which are translated via `...Schema`. Don’t pass DB column names directly.

File map to learn by example

- Repos/scopes: `internal/platform/database/repository.go`
- Transactions: `internal/platform/database/atomic.go`, `internal/platform/types/atomic/atomic.go`
- Problem errors: `internal/platform/types/problem/problem.go`; domain examples: `internal/domain/order/errors.go`
- Time series: `internal/platform/types/timeseries/timeseries.go` and usage in `internal/domain/order/service.go`
- Middleware chain: `internal/platform/request/*.go`, `internal/platform/logger/middleware.go`, responses in `internal/platform/response/response.go`
