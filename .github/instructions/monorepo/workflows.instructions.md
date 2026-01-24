---
description: "Kyora monorepo workflows — Makefile targets, docker-compose, validation scripts, development workflow"
applyTo: "Makefile,docker-compose*.yml,scripts/**,.github/workflows/**"
---

# Kyora Monorepo Workflows — Single Source of Truth (SSOT)

**SSOT Hierarchy:**

- Parent: `.github/copilot-instructions.md`, `AGENTS.md`
- Related: `monorepo/structure.instructions.md`, `monorepo/adding-projects.instructions.md`

**When to Read:**

- Setting up local development
- Running tests or builds
- Understanding CI/CD triggers
- Debugging infrastructure issues

---

## 1) Makefile Targets

The root `Makefile` provides unified commands for all apps/services.

### Core Commands

```bash
make help           # Show all available targets
make doctor         # Check tooling (go, node, docker, etc.)
```

### Infrastructure

```bash
make infra.up       # Start Postgres + Memcached + Stripe mock
make infra.down     # Stop all infrastructure
make infra.reset    # Stop, remove volumes, restart
make infra.logs     # Tail infrastructure logs
make db.psql        # Open PostgreSQL shell
```

**Infrastructure stack** (via `docker-compose.dev.yml`):

- **Postgres** (5432) — Main database
- **Memcached** (11211) — Cache layer
- **Stripe CLI** (mock webhook endpoint) — Payment testing

### Development

```bash
make dev            # Run API + portal in parallel
make dev.infra      # Start infra, then run dev
make dev.server     # Run API only (backend)
make dev.portal     # Run portal only (frontend)
```

**Environment variables:**

- `PORTAL_PORT=3001 make dev.portal` — Run portal on custom port

### Backend (Go API)

```bash
make test           # All backend tests (unit + E2E)
make test.quick     # Unit tests only (fast)
make test.e2e       # E2E tests only (requires infra)
make test.unit      # Alias for test.quick

make openapi        # Regenerate Swagger docs (swag init)
make openapi.check  # Verify docs are up-to-date
make openapi.verify # Alias for openapi.check

make build          # Build backend binary
make run            # Run backend server
```

### Frontend (Portal Web)

```bash
make portal.install # Install npm dependencies
make portal.check   # Lint + type check
make portal.build   # Production build
make portal.preview # Preview production build
```

### Validation

```bash
make agent.os.check # Validate KYORA_AGENT_OS.md manifest
```

---

## 2) Development Workflow

### Initial Setup

```bash
# 1. Check tooling
make doctor

# 2. Install dependencies
make portal.install

# 3. Start infrastructure
make infra.up

# 4. Run dev servers
make dev
```

Access:

- **API**: http://localhost:8080
- **Portal**: http://localhost:3000
- **Swagger**: http://localhost:8080/swagger/index.html

### Daily Development

```bash
# Start infra (if not running)
make infra.up

# Run dev servers
make dev
# or
make dev.server  # Backend only
make dev.portal  # Frontend only
```

### Before Committing

```bash
# Backend: run quick tests + verify OpenAPI
make test.quick
make openapi.check

# Frontend: lint + type check
make portal.check

# Full validation (slower)
make test
```

### Troubleshooting

```bash
# Reset infrastructure (clean slate)
make infra.reset

# Check infra logs
make infra.logs

# Open database shell
make db.psql

# Rebuild backend
make build
```

---

## 3) Docker Compose

### File Structure

```
docker-compose.dev.yml    # Local dev infrastructure
```

**Services:**

| Service      | Port  | Purpose             |
| ------------ | ----- | ------------------- |
| `postgres`   | 5432  | PostgreSQL database |
| `memcached`  | 11211 | Cache layer         |
| `stripe-cli` | -     | Stripe webhook mock |

### Configuration

**Environment variables** (in `docker-compose.dev.yml`):

- `POSTGRES_DB=kyora_dev`
- `POSTGRES_USER=kyora`
- `POSTGRES_PASSWORD=kyora_dev_password`

**Volumes:**

- `kyora_postgres_data` — Persistent database storage
- Cleared by `make infra.reset`

### Manual Docker Compose Commands

```bash
# Start services
docker compose -f docker-compose.dev.yml up -d

# Stop services
docker compose -f docker-compose.dev.yml down

# Remove volumes (data loss!)
docker compose -f docker-compose.dev.yml down -v

# View logs
docker compose -f docker-compose.dev.yml logs -f

# Check service status
docker compose -f docker-compose.dev.yml ps
```

---

## 4) CI/CD (GitHub Actions)

### Workflow Triggers

**Location**: `.github/workflows/`

**Trigger patterns**:

- **Push to `main`**: Full CI (backend tests, frontend checks, build)
- **Pull requests**: Same as push to `main`
- **Specific paths**: Workflows can be scoped to changed files
  - `backend/**` → Backend CI only
  - `portal-web/**` → Frontend CI only

### CI Steps

**Backend CI**:

1. Setup Go
2. Install dependencies (`go mod download`)
3. Run tests (`make test`)
4. Verify OpenAPI docs (`make openapi.check`)
5. Build binary (`make build`)

**Frontend CI**:

1. Setup Node.js
2. Install dependencies (`make portal.install`)
3. Lint + type check (`make portal.check`)
4. Build production (`make portal.build`)

### Deployment

**Not yet implemented** (future):

- Backend: Deploy to Cloud Run / ECS
- Portal: Deploy to Vercel / CloudFront
- Staging environment for PR previews

---

## 5) Validation Scripts

### Agent OS Validation

```bash
make agent.os.check
# or
./scripts/agent-os/validate.sh
```

**Checks**:

- KYORA_AGENT_OS.md manifest structure
- Agent role definitions
- Cross-references

**Exit codes**:

- `0` — All checks passed
- `1` — Validation failed

### OpenAPI Validation

```bash
make openapi.check
# or
make openapi.verify
```

**Checks**:

- Swagger docs are up-to-date with code
- Re-runs `swag init` and diffs output

**Exit codes**:

- `0` — Docs are up-to-date
- `1` — Docs need regeneration (run `make openapi`)

---

## 6) Common Commands Reference

| Task                 | Command               |
| -------------------- | --------------------- |
| **Setup**            |                       |
| Check tooling        | `make doctor`         |
| Install portal deps  | `make portal.install` |
| Start infrastructure | `make infra.up`       |
| **Development**      |                       |
| Run all dev servers  | `make dev`            |
| Run API only         | `make dev.server`     |
| Run portal only      | `make dev.portal`     |
| **Testing**          |                       |
| Quick tests          | `make test.quick`     |
| All tests            | `make test`           |
| E2E tests            | `make test.e2e`       |
| **Validation**       |                       |
| Backend OpenAPI      | `make openapi.check`  |
| Portal checks        | `make portal.check`   |
| Agent OS             | `make agent.os.check` |
| **Build**            |                       |
| Build backend        | `make build`          |
| Build portal         | `make portal.build`   |
| **Infrastructure**   |                       |
| Reset infra          | `make infra.reset`    |
| View infra logs      | `make infra.logs`     |
| Database shell       | `make db.psql`        |

---

## 7) Environment Variables

### Backend (Go)

**Location**: `backend/.env` (gitignored) or environment

| Variable                | Default                                                                        | Purpose               |
| ----------------------- | ------------------------------------------------------------------------------ | --------------------- |
| `PORT`                  | `8080`                                                                         | API server port       |
| `DATABASE_URL`          | `postgres://kyora:kyora_dev_password@localhost:5432/kyora_dev?sslmode=disable` | PostgreSQL connection |
| `MEMCACHED_SERVERS`     | `localhost:11211`                                                              | Memcached servers     |
| `JWT_SECRET`            | (required)                                                                     | JWT signing key       |
| `STRIPE_API_KEY`        | (required)                                                                     | Stripe API key        |
| `STRIPE_WEBHOOK_SECRET` | (required)                                                                     | Stripe webhook secret |
| `RESEND_API_KEY`        | (required)                                                                     | Resend email API key  |

### Frontend (Portal Web)

**Location**: `portal-web/.env` (gitignored) or environment

| Variable                      | Default                 | Purpose                |
| ----------------------------- | ----------------------- | ---------------------- |
| `VITE_API_URL`                | `http://localhost:8080` | Backend API base URL   |
| `VITE_STRIPE_PUBLISHABLE_KEY` | (required)              | Stripe publishable key |

---

## 8) Debugging Tips

### Backend Issues

**API won't start:**

```bash
# Check if infra is running
docker compose -f docker-compose.dev.yml ps

# Check database connection
make db.psql

# Check logs
make dev.server  # Run in foreground to see logs
```

**Tests failing:**

```bash
# Ensure infra is running for E2E tests
make infra.up

# Run quick tests only (no infra required)
make test.quick

# Check test output
go test -v ./backend/internal/tests/e2e/...
```

### Frontend Issues

**Portal won't start:**

```bash
# Check dependencies
make portal.install

# Check API connection
curl http://localhost:8080/health

# Run in dev mode to see errors
cd portal-web && npm run dev
```

**Build failures:**

```bash
# Type check
cd portal-web && npm run typecheck

# Lint
cd portal-web && npm run lint

# Clean install
rm -rf portal-web/node_modules portal-web/dist
make portal.install
```

### Infrastructure Issues

**Postgres connection errors:**

```bash
# Reset infra
make infra.reset

# Check logs
make infra.logs

# Manual psql
psql -h localhost -U kyora -d kyora_dev
```

**Memcached issues:**

```bash
# Check if running
docker compose -f docker-compose.dev.yml ps memcached

# Restart
make infra.reset
```

---

## 9) References

**Monorepo patterns:**

- [monorepo/structure.instructions.md] — Directory structure
- [monorepo/adding-projects.instructions.md] — Scaffolding new apps

**Backend:**

- [backend/_general/architecture.instructions.md] — Backend structure
- [backend/_general/testing.instructions.md] — Testing guidelines

**Frontend:**

- [frontend/projects/portal-web/development.instructions.md] — Portal dev workflow

**Root:**

- [Makefile] — Command definitions
- [docker-compose.dev.yml] — Infrastructure config
- [AGENTS.md] — Root agent manifest
