# Kyora Backend

Go-based API server for Kyora business management platform.

## Architecture

Kyora backend follows a layered architecture with clear separation of concerns:

```
backend/
├── cmd/                    # CLI commands (server, sync_plans, etc.)
├── internal/
│   ├── domain/            # Business logic modules
│   │   ├── account/       # User accounts and authentication
│   │   ├── billing/       # Stripe billing and subscriptions
│   │   ├── customer/      # Customer management
│   │   ├── inventory/     # Product inventory tracking
│   │   ├── order/         # Order management
│   │   └── ...
│   ├── platform/          # Infrastructure layer
│   │   ├── auth/          # JWT, OAuth
│   │   ├── bus/           # Event bus
│   │   ├── cache/         # Memcached
│   │   ├── config/        # Configuration
│   │   ├── database/      # Database connections, repositories
│   │   ├── email/         # Email providers
│   │   ├── logger/        # Structured logging
│   │   └── types/         # Shared types
│   ├── server/            # HTTP server setup and routing
│   └── tests/             # E2E and integration tests
├── main.go                # Application entry point
├── go.mod                 # Go module dependencies
└── .air.toml              # Hot reload configuration
```

## Tech Stack

- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: PostgreSQL
- **Cache**: Memcached
- **Auth**: Access JWT (Authorization: Bearer), refresh-token sessions (rotation/revocation), Google OAuth
- **Payments**: Stripe SDK
- **Email**: Resend
- **Logging**: slog
- **CLI**: Cobra

## Development

### Prerequisites

- Go 1.25.3 or higher
- Docker Desktop
- Air (for hot reload): `go install github.com/air-verse/air@latest`

### Local Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Configure environment:**
   Copy `.kyora.yaml.example` to `.kyora.yaml` and configure:
   - Database DSN
   - Cache hosts
   - JWT secret
    - Refresh token TTL (optional)
   - Stripe API keys
   - Email provider settings

3. **Run development server:**
   ```bash
   # From repository root
   make dev.server
   
   # Or from backend directory
   cd backend && air server
   ```

### Configuration

Configuration is managed through `.kyora.yaml` and environment variables using Viper.

Key configuration sections:
- `app.*` - Application settings
- `http.*` - HTTP server settings
- `database.*` - PostgreSQL connection
- `cache.*` - Memcached settings
- `auth.*` - JWT (access), refresh token TTL, and OAuth settings
- `billing.*` - Stripe configuration
- `email.*` - Email provider settings

See `internal/platform/config/config.go` for all available configuration options.

## Testing

### Running Tests

```bash
# All tests
make test

# Unit tests only
make test.unit

# E2E tests only
make test.e2e

# Quick run (no verbose)
make test.quick
```

### Coverage Reports

```bash
# Generate coverage report
make test.coverage

# Generate HTML report
make test.coverage.html

# Open HTML report in browser
make test.coverage.view

# E2E coverage
make test.e2e.coverage
```

### Test Organization

- **Unit tests**: Located within domain/platform folders (e.g., `internal/domain/order/service_test.go`)
- **E2E tests**: Located in `internal/tests/e2e/`
- **Test utilities**: `internal/tests/testutils/`
- **Mocks**: `internal/tests/mocks/`

E2E tests use testcontainers to spin up isolated Docker containers for complete integration testing.

## Domain Conventions

Each domain module follows a standard structure:

- `model.go` - GORM models, schemas, DTOs
- `storage.go` - Data access layer using generic repositories
- `service.go` - Business logic
- `errors.go` - Domain-specific errors (RFC 7807 Problem types)
- `handler_http.go` - HTTP handlers
- `middleware_http.go` - Domain middleware (optional)
- `state_machine.go` - State transitions (optional)

### Core Patterns

**Repository Pattern:**
```go
repo := database.NewRepository[Order](db)
orders, err := repo.FindMany(ctx, 
    repo.ScopeWorkspaceID(workspaceID),
    repo.WithOrderBy([]string{"createdAt DESC"}),
    repo.WithPagination(page, pageSize),
)
```

**Transactions:**
```go
err := atomic.Exec(ctx, func(tctx context.Context) error {
    // All operations use tctx for transaction
    order, err := s.storage.order.CreateOne(tctx, order)
    if err != nil {
        return err
    }
    return s.storage.orderItem.CreateMany(tctx, items)
})
```

**Error Handling:**
```go
// Define domain errors
var ErrOrderNotFound = problem.NotFound("order not found")

// Return with context
return ErrOrderNotFound.With("orderId", id).WithError(err)
```

## Multi-Tenancy

Kyora uses workspace-based multi-tenancy:

- Each **workspace** represents one business with its own subscription
- **Users** belong to a single workspace with role-based permissions
- All data is strictly isolated by `WorkspaceID`
- Always scope queries: `repo.ScopeWorkspaceID(workspace.ID)`

## API Documentation

API documentation is generated using Swagger/OpenAPI.

HTTP endpoints follow REST conventions with middleware chain:
1. Logger - Adds trace ID
2. Authentication - Validates access JWT from `Authorization: Bearer <token>`
3. Actor validation - Loads user
4. Workspace membership - Validates workspace access
5. Permission check - Validates role permissions
6. Subscription check - Ensures active subscription (optional)
7. Plan limits - Enforces usage limits (optional)

Auth token flow:
- Login/onboarding/invitation acceptance return `{ user, token, refreshToken }`
- Use `token` as the access JWT for authenticated requests
- Use `refreshToken` with `POST /v1/auth/refresh` to rotate and get new tokens
- Logout endpoints revoke refresh sessions: `POST /v1/auth/logout`, `/v1/auth/logout-others`, `/v1/auth/logout-all`

## Commands

The backend uses Cobra CLI framework:

```bash
# Run HTTP server
go run main.go server

# Sync Stripe plans
go run main.go sync-plans

# Cleanup onboarding sessions
go run main.go onboarding-cleanup
```

## Best Practices

1. **Always scope by WorkspaceID** for multi-tenancy
2. **Use domain storage layers** - never raw SQL
3. **Use atomic transactions** for multi-step operations
4. **Return domain errors** (RFC 7807 Problems)
5. **Use decimal for money** - never float64
6. **Store timestamps in UTC**
7. **Use structured logging** via slog
8. **Document with godoc style** comments
9. **Aim for 80%+ test coverage**

## Contributing

1. Follow existing code patterns and conventions
2. Add tests for new features
3. Update documentation
4. Ensure all tests pass: `make test`
5. Check coverage: `make test.coverage`

## License

See [LICENSE](../LICENSE) file for details.
