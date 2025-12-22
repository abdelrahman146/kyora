# Kyora

A business management assistant for social media entrepreneurs and small teams.

## Monorepo Structure

This repository is organized as a monorepo containing multiple projects:

```
kyora/
├── backend/          # Go backend API server
├── Makefile          # Root-level build commands
├── LICENSE           # Project license
└── README.md         # This file
```

### Backend

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

## Development

### Prerequisites

- Go 1.25.3 or higher
- Docker Desktop (for running tests)
- Air (for hot reload): `go install github.com/air-verse/air@latest`

### Available Make Commands

Run `make help` to see all available commands:

```bash
make help
```

## Future Projects

This monorepo will include:
- **Frontend**: Web application (coming soon)
- **Mobile**: Mobile applications (coming soon)
- **Services**: Additional microservices (coming soon)

## License

See [LICENSE](LICENSE) file for details.
