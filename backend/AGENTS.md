# backend/AGENTS.md

## Scope

Go API monolith — domain services, platform infrastructure, HTTP handlers, and tests.

**Parent AGENTS.md**: [../AGENTS.md](../AGENTS.md) (read first for project context and global boundaries)

## Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: Gin
- **ORM**: GORM (PostgreSQL)
- **Cache**: Memcached
- **Config**: Viper
- **Logging**: slog (structured JSON)
- **CLI**: Cobra
- **Payments**: Stripe SDK
- **Email**: Resend SDK
- **Storage**: Blob abstraction (local filesystem or S3-compatible)

## Setup Commands

```bash
# Start infra (Postgres, Memcached, Stripe mock)
make infra.up

# Run API server (hot reload via air)
make dev.server

# Run all tests
make test

# Run unit tests only (fast)
make test.quick

# Run E2E tests only
make test.e2e

# Regenerate Swagger/OpenAPI docs
make openapi

# Verify OpenAPI docs are up-to-date
make openapi.check
```

## Structure

```
backend/
├── cmd/                    # Cobra CLI commands
│   ├── root.go             # Root command
│   ├── server.go           # `kyora server` - starts HTTP server
│   ├── seed.go             # `kyora seed` - seeds test data
│   └── sync_plans.go       # `kyora sync_plans` - syncs Stripe plans
├── docs/                   # Swagger/OpenAPI generated docs
├── internal/
│   ├── server/             # HTTP server bootstrap
│   │   ├── server.go       # Gin engine, DI, middleware
│   │   └── routes.go       # Route registration
│   ├── domain/             # Business modules (DDD-ish)
│   │   ├── account/        # Users, workspaces, RBAC
│   │   ├── business/       # Business entities
│   │   ├── inventory/      # Products, variants, categories
│   │   ├── order/          # Orders
│   │   ├── customer/       # Customers, addresses
│   │   ├── accounting/     # Expenses, investments, withdrawals
│   │   ├── analytics/      # Dashboards, reports
│   │   ├── billing/        # Stripe subscriptions
│   │   ├── asset/          # File uploads
│   │   ├── storefront/     # Public storefront API
│   │   ├── onboarding/     # Onboarding flow
│   │   └── metadata/       # System metadata
│   ├── platform/           # Infrastructure (shared)
│   │   ├── auth/           # JWT, tokens
│   │   ├── blob/           # File storage abstraction
│   │   ├── bus/            # Event bus (in-process)
│   │   ├── cache/          # Memcached client
│   │   ├── config/         # Viper config
│   │   ├── database/       # GORM, migrations, atomic processing
│   │   ├── email/          # Resend client
│   │   ├── logger/         # slog setup + middleware
│   │   ├── middleware/     # HTTP middleware (CORS, rate limit)
│   │   ├── request/        # Request parsing helpers
│   │   ├── response/       # Response helpers (RFC7807)
│   │   ├── types/          # Shared types (Problem, etc.)
│   │   └── utils/          # Utilities (slugs, pagination, etc.)
│   └── tests/
│       ├── e2e/            # End-to-end HTTP tests
│       └── testutils/      # Shared test helpers
└── tmp/                    # Hot reload artifacts (gitignored)
```

## Code Style

### Domain Module Pattern

Each domain follows this layout:

```
domain/<name>/
├── handler.go         # HTTP handlers (Gin)
├── service.go         # Business logic
├── storage.go         # Database access (GORM)
├── model.go           # GORM models
├── model_response.go  # Response DTOs
├── dto.go             # Request DTOs
├── middleware.go      # Domain-specific middleware (if any)
└── events.go          # Domain events (if any)
```

### HTTP Handler Pattern

```go
// ✅ Good: Structured handler with validation
func (h *Handler) CreateOrder(c *gin.Context) {
    business := business.FromContext(c)  // middleware-injected
    
    var req CreateOrderRequest
    if err := request.ValidBody(c, &req); err != nil {
        response.Error(c, err)
        return
    }
    
    order, err := h.service.CreateOrder(c.Request.Context(), business.ID, req)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.SuccessJSON(c, http.StatusCreated, ToOrderResponse(order))
}
```

### Service Pattern

```go
// ✅ Good: Tenant-scoped with validation
func (s *Service) CreateOrder(ctx context.Context, businessID string, req CreateOrderRequest) (*Order, error) {
    if businessID == "" {
        return nil, types.NewValidationError("business_id required")
    }
    
    // Business logic here
    order := &Order{
        BusinessID: businessID,
        // ...
    }
    
    if err := s.storage.Create(ctx, order); err != nil {
        return nil, err
    }
    
    return order, nil
}
```

### Response DTO Pattern

```go
// model_response.go

type OrderResponse struct {
    ID          string    `json:"id"`           // Always camelCase
    BusinessID  string    `json:"businessId"`
    Status      string    `json:"status"`
    TotalAmount float64   `json:"totalAmount"`
    CreatedAt   time.Time `json:"createdAt"`
}

func ToOrderResponse(o *Order) OrderResponse {
    return OrderResponse{
        ID:          o.ID,
        BusinessID:  o.BusinessID,
        Status:      string(o.Status),
        TotalAmount: o.TotalAmount,
        CreatedAt:   o.CreatedAt,
    }
}
```

## Boundaries (Backend-Specific)

### ✅ Always do

- Use `request.ValidBody(c, &req)` for all request parsing (strict JSON)
- Return `response.Error(c, err)` for all errors (RFC7807 Problem JSON)
- Convert models to response DTOs via `To*Response()` functions
- Scope all queries by `businessID` or `workspaceID` (tenant isolation)
- Use `database.AtomicProcess` for multi-table transactions
- Add Swagger annotations for new endpoints
- Write E2E tests for new endpoints

### ⚠️ Ask first

- New database migrations
- New domain modules
- Cross-domain service dependencies
- New middleware in the auth chain
- Changes to billing/subscription logic

### 🚫 Never do

- Return raw GORM models in responses
- Use raw SQL for domain logic (use repository pattern)
- Access another domain's storage directly (use service)
- Skip tenant scoping (businessID/workspaceID)
- Use `c.BindJSON()` or `ShouldBindJSON()` (use `request.ValidBody`)
- Hardcode config values (use Viper constants)

## SSOT Entry Points

- [.github/instructions/backend-core.instructions.md](../.github/instructions/backend-core.instructions.md) — Architecture
- [.github/instructions/go-backend-patterns.instructions.md](../.github/instructions/go-backend-patterns.instructions.md) — Go patterns
- [.github/instructions/backend-testing.instructions.md](../.github/instructions/backend-testing.instructions.md) — Testing
- [.github/instructions/errors-handling.instructions.md](../.github/instructions/errors-handling.instructions.md) — Error patterns
- [.github/instructions/responses-dtos-swagger.instructions.md](../.github/instructions/responses-dtos-swagger.instructions.md) — DTOs/OpenAPI

## Agent Routing Hints

**Backend Lead** (`@Backend Lead`): Architecture, API contracts, domain modeling
**Backend Implementer** (`@Backend Implementer`): Code changes, tests, OpenAPI updates
**QA/Test Specialist** (`@QA/Test Specialist`): Test coverage, E2E tests
