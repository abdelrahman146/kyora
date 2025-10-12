# Kyora — Copilot Instructions

Concise, codebase-specific guidance for AI agents.

## Big picture

- Go monolith with a Cobra CLI and two primary commands:
  - `web`: HTTP server (Gin) + Templ views. See `cmd/web.go`.
  - `recurring`: daily recurring expenses job. See `cmd/recurring.go`.
- Config via Viper from `config.yaml` (AutomaticEnv enabled but no key replacer; prefer YAML).
- Data: Postgres via GORM wrapped by `internal/db.Postgres`; caching via `internal/db.Memcache`.
- Domains under `internal/domain/**`: Repository + Service layers. Repos expose GORM scopes; Services orchestrate logic/transactions.

## App startup (what happens)

- `main.go` → `cmd.Execute()`.
- `cmd/web.go`: read config → setup logging (`utils.Log`) → create Postgres + Memcache → create `db.AtomicProcess` → init domains → build Gin router (`webrouter.NewRouter()`) → instantiate handlers and `RegisterRoutes` → `router.Run(server.port)`.
- `cmd/recurring.go`: similar boot, runs `postgres.AutoMigrate(expense.Expense, expense.RecurringExpense)` then `expense.Service.ProcessRecurringExpensesDaily(ctx)`.

## Data access patterns (do this)

- Always pass a `context.Context`; use `Postgres.Conn(ctx, opts...)`. It auto-binds to ambient tx in ctx.
- Transactions: `atomic.Exec(ctx, func(txCtx) error { ... })`; inside use `txCtx` so repos share the tx.
- Compose queries with repo scopes + Postgres options:
  - Scopes: repo methods like `ScopeID`, `ScopeStoreID`, `ScopeFilter` etc.
  - Options: `db.WithPagination`, `db.WithSorting`, `db.WithPreload`, `db.WithLock`, `db.WithLimit`, `db.WithJoins`.
- Examples: see `inventory/product_repository.go` (scopes, upsert) and `order/order_service.go` (calc totals, preloads, transactions).

## Web conventions (Gin + Templ)

- Router: `webrouter.NewRouter()` adds recovery, `middleware.LoggerMiddleware()`, and `registerRoutes(r)`; handlers register routes in `cmd/web.go`.
- Handlers: implement `RegisterRoutes(r gin.IRoutes)` and actions like `Index/New/Show/Edit`.
  - Each action builds `webcontext.PageInfo` → `webcontext.SetupPageInfo` on the request context → `webutils.Render(c, status, templComponent)`.
- Fragments/redirects: `webutils.RenderFragments(c, status, component, keys...)`; `webutils.Redirect(c, location)` sets `HX-Redirect`.
- Auth: `AuthRequiredMiddleware` requires JWT cookie `jwt` (see `utils.JWT`); `GuestRequiredMiddleware` redirects logged-in users from guest pages.

## Errors, logging, IDs

- Structured logging: use `utils.Log.FromContext(ctx)`; `LoggerMiddleware` injects `X-Trace-ID` (configurable via `server.trace_id_header`).
- DB errors → RFC7807: `db.HandleDBError(err)` returns `utils.ProblemDetails` (NotFound/Conflict/Internal).
- IDs: `utils.ID.NewKsuid()` (trace); `utils.ID.NewBase62WithPrefix(prefix, n)` for human-friendly codes (orders).

## Dev workflows

- Make:
  - `make templates` (runs `templ generate`).
  - `make dev.web` (requires `air`; clears `tmp/` then `air web`).
  - `make dev.css` (`yarn css:watch` → builds Tailwind v4 + daisyUI v5 from `public/css/tw-input.css` → `public/css/base.css`).
- Without Make: `go run . web` or `go run . recurring`; CSS via `yarn css:watch`.
- Config: edit `config.yaml` (Postgres DSN, memcache hosts, JWT secret/expiry, ports, site info).

## Frontend styling

- Tailwind CSS v4 + daisyUI v5 (see `package.json`). Prefer daisyUI component classes in Templ views; avoid custom CSS when possible.
- For daisyUI specifics (components, themes, upgrade notes), see `.github/instructions/daisyui.instructions.md` in this repo.

## Adding features (happy path)

- Repo: add under the relevant domain with `Scope*` and CRUD methods using `db.Postgres`.
- Service: orchestrate, use `AtomicProcess` for multi-write operations, use `decimal` for money.
- Web: new handler + `RegisterRoutes`; set `PageInfo`; render components from `internal/web/views/**` via `webutils.Render`.
- Wire the handler in `cmd/web.go`.

## External deps

Gin, Templ, GORM, Cobra, Viper, gomemcache, govalues/decimal, golang-jwt/jwt.

Notes: env overrides are limited (no key replacer); migrations aren’t centralized (only `recurring` migrates expense tables).
