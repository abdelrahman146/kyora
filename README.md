# Kyora

A business management assistant for social media entrepreneurs and small teams.

## Monorepo Structure

This repository is organized as a monorepo containing multiple projects.

```
kyora/
├── backend/                 # Go backend API server
├── storefront-web/          # Customer-facing storefront (React)
├── docker-compose.dev.yml   # Local infra (db/cache/etc.) for development
├── Makefile                 # Root-level build commands
├── README.md                # This file
└── STRUCTURE.md             # Monorepo guidelines
```

## Backend

The backend is a Go-based API server that handles all business logic, data persistence, and integrations.

**Tech Stack:**
- Go 1.25.3
- Gin Web Framework
- GORM ORM
- PostgreSQL
- Memcached
- Stripe SDK

**Quick Start:**
```bash
# Start local infra (db/cache/etc.)
docker compose -f docker-compose.dev.yml up -d

# Run development server with hot reload
make dev.server

# Run tests
make test

# Run E2E tests
make test.e2e

# View coverage report
make test.coverage.view
```

See [backend/README.md](backend/README.md) for detailed backend documentation.

## Storefront Web

`storefront-web/` is the customer-facing storefront (mobile-first, RTL-aware) built with React + TypeScript.

**Tech Stack:**
- React 19
- React Router v7
- Tailwind CSS v4 (CSS-first)
- daisyUI v5
- TanStack Query + Zustand
- i18next (Arabic/English, RTL-first)

**Quick Start:**
```bash
cd storefront-web
npm install
npm run dev
```

See [storefront-web/DESIGN_SYSTEM.md](storefront-web/DESIGN_SYSTEM.md) for the storefront design system implementation details.

## Development

### Prerequisites

- Go 1.25.3 or higher
- Node.js (for `storefront-web/`)
- Docker Desktop (for running tests)
- Air (for hot reload): `go install github.com/air-verse/air@latest`

### Available Make Commands

Run `make help` to see all available commands:

```bash
make help
```

## License

See [LICENSE](LICENSE) file for details.
