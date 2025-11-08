Kyora — AI coding agent quickstart

Big picture

- Go 1.25 monolith using Cobra CLI (entry: `main.go` → `cmd/root.go`). HTTP runs via the `server` subcommand in `cmd/server.go` which delegates to `internal/server/server.go`.
- Layered `internal/` packages:
  - platform: infra (config, DB, cache, logging, request/response, auth, event bus, types).
  - domain: business modules (account, accounting, analytics, billing, business, customer, inventory, order, onboarding). Typical files: `model.go`, `storage.go`, `service.go`, `errors.go`, optional `state_machine.go`.

Run and configure

- Local dev: `make dev.server` (needs `air`). Air builds to `./tmp/main` and runs `kyora server` for live reload.
- Config: Viper loads `.kyora.yaml` + env. Treat `internal/platform/config/config.go` as the source of truth for keys (app/http/db/cache/jwt/google oauth/stripe/email). Example keys: `http.port`, `database.dsn`, `auth.jwt.secret`, `billing.stripe.api_key`, `email.provider`.
- Data stores: Postgres via GORM; memcached cache. Stripe Billing and Google OAuth are opt‑in via config.

Core patterns you must follow

- DB access: use `internal/platform/database.Database.Conn(ctx)` (transaction-aware). Prefer generic `database.Repository[T]` with scopes: `ScopeBusinessID/WorkspaceID/Equals/In/Time/CreatedAt`, plus `WithPreload`, joins, `WithOrderBy`, pagination, `WithReturning`, aggregates and time series.
- Transactions: wrap multi-step writes with `database.AtomicProcess.Exec(ctx, func(ctx) error, ...)`. The transaction is injected into `ctx` so repos inside the closure are consistent.
- HTTP: use middleware chain in `internal/platform/logger/middleware.go` (adds `traceId` from `http.trace_id_header`, default `X-Trace-ID`). Auth uses JWT cookie named `jwt`; call `request.EnforceAuthentication`, then `request.EnforceValidActor` and `request.EnforceBusinessValidity` to resolve user/business.
- Responses & errors: return via `response.SuccessJSON/SuccessText/SuccessEmpty` and `response.Error`. Errors are RFC 7807 Problems (`internal/platform/types/problem`) and DB errors normalize via `database.IsRecordNotFound/IsUniqueViolation`.

Domain conventions (copy when adding modules)

- `model.go`: GORM model with string `ID` from `utils/id` (e.g., `id.KsuidWithPrefix("ord")`). Define `...Schema` using `types/schema.Field` to map DB columns → JSON names for ordering/search.
- `storage.go`: construct `database.Repository[T]` per aggregate; auto-migrate on construction.
- `service.go`: methods `(ctx, actor *account.User, biz *business.Business, ...)` to enforce tenancy. Always scope queries by `BusinessID`. Use atomic processor for multi-entity updates. Prefer domain errors from `errors.go`.
- `state_machine.go`: encode transitions (see `internal/domain/order/state_machine.go`) and maintain transition timestamps.

Lists & analytics

- Pagination/sorting: pass `*list.ListRequest`; build ordering with `req.ParsedOrderBy(domainSchema)`. The order-by value expects JSON field names (e.g., `-createdAt`), not DB columns. Wrap results with `list.NewListResponse`.
- Time series: use repo helpers `TimeSeriesSum/TimeSeriesCount`; labels via `internal/platform/types/timeseries`.

Events & integrations

- Event bus: `internal/platform/bus` with `Bus.Emit(topic, payload)` and `Bus.Listen(topic, handler)`. Built-in topics: `VerifyEmailTopic`, `ResetPasswordTopic`.
- Billing/Stripe: domain code in `internal/domain/billing` (handlers, services, webhooks). Config keys: `billing.stripe.api_key`, `billing.stripe.webhook_secret`.
- Email: `internal/platform/email` with templates under `internal/platform/email/templates`. Provider via config: `email.provider` = `resend` or `mock` (see also `EmailFrom*` keys). Use `SendTemplate` for built-in templates (see package README).
- OAuth: `internal/platform/auth/google_oauth.go`; JWT helpers in `internal/platform/auth/jwt.go` (cookie name `jwt`).

Gotchas

- Multi‑tenancy: always scope by `BusinessID` (and/or `WorkspaceID`) using repository scopes; many services assume this.
- Config truth: rely on constants in `internal/platform/config/config.go` instead of guessing key names.
- JWT source: `JwtFromContext` reads only the cookie, not `Authorization` header.
- Ordering: use JSON field names mapped by your domain `...Schema`.

Pointers (good examples)

- Repos/scopes: `internal/platform/database/repository.go`
- Transactions: `internal/platform/database/atomic.go`, `internal/platform/types/atomic/atomic.go`
- Errors/Problems: `internal/platform/types/problem/problem.go`, e.g., `internal/domain/order/errors.go`
- Time series: `internal/platform/types/timeseries/timeseries.go` and `internal/domain/order/service.go`
- HTTP helpers: `internal/platform/request/*.go`, `internal/platform/logger/middleware.go`, `internal/platform/response/response.go`
