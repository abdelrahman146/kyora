---
description: Go Backend Project Patterns (Kyora-style, reusable)
applyTo: "backend/**"
---

# Go Backend Project Patterns (Reusable SSOT)

This file documents a **reusable Go backend pattern** based on Kyora’s backend.
Use it when creating or extending any Go API in this repo now or in the future (e.g., an admin backend).

Scope:

- **How to structure the project** (folders, dependencies, layering)
- **How to wire HTTP, config, logging, and persistence**
- **What rules are strict** (security, request validation, errors, transactions)

If you are modifying the Kyora backend specifically, also read:

- `.github/instructions/backend-core.instructions.md` (Kyora backend: ground truth wiring)

If you are writing tests, also read:

- `.github/instructions/backend-testing.instructions.md`

## Core architecture (layered, dependency direction)

Recommended baseline layout:

```
<project>/
  cmd/                 # Cobra CLI entrypoints (server, jobs, scripts)
  internal/
    server/            # HTTP server wiring and route registration (DI only)
    platform/          # Cross-cutting infrastructure (config, db, cache, auth, request/response, logger, types, utils)
    domain/            # Business modules (model/storage/service/http)
    tests/             # E2E/integration test harness (optional)
  main.go              # CLI entry (delegates to cmd)
```

Strict dependency rules:

- `internal/domain/**` may depend on `internal/platform/**`.
- `internal/platform/**` must not depend on `internal/domain/**`.
- Prefer cross-domain calls via **domain services** (not storage).
- Keep infrastructure decisions in `platform` and business invariants in `domain`.

## Domain module structure (strict file organization)

Each domain lives under `internal/domain/<name>/` and should contain:

- `model.go`: GORM models, enums, constants, and schema definitions
  - Define table/struct name constants: `OrderTable`, `OrderStruct`, `OrderPrefix`
  - Define schema mappings: `OrderSchema` with field definitions
  - Simple domains may keep request DTOs here
- `model_request.go` (recommended for complex domains): request DTOs only
  - All `CreateXRequest`, `UpdateXRequest`, etc. structures
  - Isolates API contract from domain models
- `model_response.go` (required): response DTOs only
  - All `XResponse` structures that map GORM models to API responses
  - Must use camelCase JSON tags
  - Must never expose GORM internals (`gorm.Model`, `DeletedAt`, etc.)
- `storage.go`: repositories and cache wiring
  - Keep caching logic here
  - Define storage indexes and search configuration
- `service.go`: business logic (source of truth for invariants)
  - Services orchestrate; they never construct queries
  - Use repository scopes and schema constants only
- `errors.go`: domain-specific `problem.*` constructors
- `handler_http.go`: HTTP handlers (parse/validate → service → respond)
  - Handlers must map GORM models to response DTOs before returning
- `middleware_http.go`: optional domain-specific middleware
- `state_machine.go`: optional explicit transition rules

### No magic strings rule (strictly enforced)

**ALL string literals for database identifiers must be constants or schema references:**

❌ **BAD** (magic strings):

```go
// In service.go
s.storage.order.WithPreload("ShippingAddress")
s.storage.order.WithPreload("Items.Product")
s.storage.order.ScopeWhere("business_id = ?", bizID)
```

✅ **GOOD** (constants and schema):

```go
// In model.go
const (
    OrderTable  = "orders"
    OrderStruct = "Order"
    OrderPrefix = "ord"
    OrderItemStruct = "OrderItem"
    OrderNoteStruct = "OrderNote"
)

var OrderSchema = struct {
    ID         schema.Field
    BusinessID schema.Field
    // ... all fields
}{
    ID:         schema.NewField("id", "id"),
    BusinessID: schema.NewField("business_id", "businessId"),
    // ...
}

// In service.go
s.storage.order.WithPreload(OrderItemStruct)
s.storage.order.WithPreload("ShippingAddress") // Only acceptable if relation name is not a simple field
s.storage.order.ScopeBusinessID(bizID) // Use scope helper instead of raw WHERE
s.storage.order.ScopeEquals(OrderSchema.BusinessID, bizID) // When scope helper not available
```

**Rationale**: Magic strings are brittle, hard to refactor, and bypass compile-time checking. Constants provide:

- Single source of truth for identifiers
- Compile-time validation
- IDE autocomplete and refactoring support
- Easy global search and replace

## HTTP server wiring (DI only)

`internal/server` should:

- Create shared platform dependencies once (config, logger, DB, cache, integrations).
- Construct domain storages → services → handlers.
- Register routes and middleware chains.

Strict rules:

- Do not embed business logic in routing.
- Keep handlers thin: parse/validate → call service → respond.

## Configuration (SSOT is code constants)

Pattern:

- Define every config key as a constant in `internal/platform/config`.
- Provide sane defaults in a single `Configure()` function.
- Load config once at CLI startup (Cobra `PersistentPreRunE`) and avoid process exits inside config helpers.

Strict rules:

- No “magic strings” for config keys.
- Defaults must be safe for local dev.

## Logging (structured and request-scoped)

Pattern:

- Use structured logging (`log/slog`).
- Add a request middleware that:
  - generates/propagates a trace id
  - attaches a logger into the request context
  - logs request start and completion

Strict rules:

- Never log secrets (tokens, API keys, refresh tokens).
- Enrich logs with actor/tenant context _after_ authentication.

## Request validation (strict JSON)

Pattern:

- Use a single helper to decode JSON bodies that:
  - rejects unknown fields (`DisallowUnknownFields`)
  - rejects trailing tokens
  - validates using struct tags
- **Organize request DTOs**: all request structs go in `model_request.go` per domain (or remain in `model.go` for simple domains)

Strict rules:

- Do not use ad-hoc decoding per handler.
- Prefer explicit request DTO structs with `binding:` tags.
- **No magic strings in validation**: enum values should reference constants.

## Responses and errors (RFC 7807)

Pattern:

- Represent API errors using a Problem JSON type (`application/problem+json`).
- Route all errors through a single helper, which:
  - writes Problem JSON
  - aborts the request
  - maps common DB errors to stable HTTP codes
- **Enforce response DTO pattern**: never return GORM models directly to clients.
  - Response DTOs go in `model_response.go` per domain
  - Map GORM models to response DTOs in handlers or service layer

Strict rules:

- Never return inconsistent error shapes.
- Domain errors should be created once (domain `errors.go`) and enriched with `.With(key, value)`.
- **Never expose GORM internals** (e.g., `gorm.Model`, PascalCase timestamps) to API responses.
- Response DTOs must use **camelCase** for JSON fields to match frontend expectations.

## Persistence (repository + scopes)

Pattern:

- Use a typed repository wrapper per model and a set of reusable query scopes.
- Use schema objects to map API-facing JSON fields to DB columns.
- Define constants for table names, struct names, and prefixes in `model.go`.

Strict rules (separation of concerns):

- **Service methods must never construct queries**: services call repository methods with scopes.
- **No magic strings for tables/columns**: always use constants from schema or model constants:
  - For preloads: use struct name constants (e.g., `OrderItemStruct`, `customer.CustomerStruct`)
  - For column references: use schema fields (e.g., `OrderSchema.ID`, `CustomerSchema.Email`)
  - For table references: use table constants (e.g., `OrderTable`, `CustomerTable`)
- Do not accept raw DB column names from clients.
- Avoid raw SQL for domain logic; if you must use custom WHERE clauses, bind variables (no string concatenation).

## List/Count methods pattern (critical separation of concerns)

List and Count methods are particularly prone to violating clean architecture by constructing complex queries directly in service methods. This section defines the correct pattern.

### Service layer responsibilities (thin orchestration only)

Service methods for listing/counting must:

1. Build scope arrays using repository scope methods
2. Delegate query execution to storage/repository
3. Never construct SQL fragments, JOINs, or aggregations
4. Keep method length under 50 lines (ideally 20-30)

**❌ BAD - Service constructs queries directly:**

```go
func (s *Service) ListProducts(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListProductsFilters) ([]*Product, int64, error) {
    baseScopes := []func(db *gorm.DB) *gorm.DB{
        s.storage.products.ScopeWhere("products.business_id = ?", biz.ID), // Magic string
    }

    // 🚫 BAD: Complex join logic in service
    needsVariantJoin := false
    useInnerJoin := false
    needsVariantsCountSort := false
    for _, ob := range req.OrderBy() {
        if ob == "variantsCount" || ob == "-variantsCount" {
            needsVariantsCountSort = true
            needsVariantJoin = true
        }
    }

    // 🚫 BAD: Building SQL with string concatenation in service
    if needsAggregation {
        var aggFields []string
        if needsVariantsCountSort {
            aggFields = append(aggFields, "COUNT(*)::int as variants_count")
        }
        joinSQL := fmt.Sprintf(`
            LEFT JOIN LATERAL (
                SELECT %s
                FROM variants
                WHERE variants.product_id = products.id
                    AND variants.deleted_at IS NULL
            ) AS product_agg ON true
        `, strings.Join(aggFields, ", "))

        baseScopes = append([]func(db *gorm.DB) *gorm.DB{
            s.storage.products.WithJoins(joinSQL), // 🚫 Magic strings everywhere
        }, baseScopes...)

        // 🚫 BAD: GROUP BY construction in service
        listExtra = append(listExtra, func(db *gorm.DB) *gorm.DB {
            return db.Group(strings.Join(groupByColumns, ", "))
        })
    }

    // 🚫 BAD: Search JOIN construction in service
    if req.SearchTerm() != "" {
        baseScopes = append(baseScopes,
            s.storage.products.WithJoins("LEFT JOIN categories ON categories.id = products.category_id"), // Magic strings
            s.storage.products.ScopeWhere(
                "(products.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.search_vector @@ websearch_to_tsquery('simple', ?))",
                term, term,
            ),
        )
    }

    // 🚫 BAD: 112 lines of query construction logic in service
    findOpts := append([]func(*gorm.DB) *gorm.DB{}, baseScopes...)
    // ... more query building ...
}
```

**✅ GOOD - Service delegates to storage:**

```go
func (s *Service) ListProducts(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListProductsFilters) ([]*Product, int64, error) {
    // Build scopes only - no query construction
    scopes := []func(*gorm.DB) *gorm.DB{
        s.storage.products.ScopeBusinessID(biz.ID), // Uses schema constant internally
    }

    // Apply filters using storage scope methods
    if filters != nil {
        if filters.CategoryID != "" {
            scopes = append(scopes, s.storage.products.ScopeCategoryID(filters.CategoryID))
        }
        if filters.StockStatus != "" {
            scopes = append(scopes, s.storage.products.ScopeStockStatus(filters.StockStatus))
        }
    }

    // Delegate search to storage
    if req.SearchTerm() != "" {
        scopes = append(scopes, s.storage.products.ScopeSearch(req.SearchTerm()))
    }

    // Delegate to repository for execution
    items, err := s.storage.products.FindMany(ctx,
        append(scopes,
            s.storage.products.WithPagination(req.Offset(), req.Limit()),
            s.storage.products.WithOrderBy(req.ParsedOrderBy(ProductSchema)),
            s.storage.products.WithPreload(VariantStruct), // Constant, not magic string
        )...,
    )
    if err != nil {
        return nil, 0, err
    }

    count, err := s.storage.products.Count(ctx, scopes...)
    if err != nil {
        return nil, 0, err
    }

    return items, count, nil
}
```

### Storage layer responsibilities (owns query construction)

Storage methods must encapsulate all query construction logic:

**✅ GOOD - Complex query logic in storage:**

```go
// In storage.go
type ProductStorage struct {
    repository *database.Repository[Product]
}

// Scope methods encapsulate query logic
func (s *ProductStorage) ScopeStockStatus(status string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        switch status {
        case StockStatusInStock:
            // Uses schema constants, not magic strings
            return db.Joins(fmt.Sprintf("INNER JOIN %s ON %s.%s = %s.%s",
                VariantTable,
                VariantTable, VariantSchema.ProductID,
                ProductTable, ProductSchema.ID,
            )).Where(fmt.Sprintf("%s.%s > %s.%s",
                VariantTable, VariantSchema.StockQuantity,
                VariantTable, VariantSchema.StockAlert,
            ))
        case StockStatusLowStock:
            // Encapsulate complex logic here
            return db.Joins(/* ... */).Where(/* ... */)
        default:
            return db
        }
    }
}

func (s *ProductStorage) ScopeSearch(term string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        like := "%" + term + "%"
        // All JOIN and search logic encapsulated
        return db.
            Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
                CategoryTable,
                CategoryTable, CategorySchema.ID,
                ProductTable, ProductSchema.CategoryID,
            )).
            Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
                VariantTable,
                VariantTable, VariantSchema.ProductID,
                ProductTable, ProductSchema.ID,
            )).
            Where(
                fmt.Sprintf("(%s.search_vector @@ websearch_to_tsquery('simple', ?) OR %s.search_vector @@ websearch_to_tsquery('simple', ?))",
                    ProductTable, VariantTable),
                term, term,
            )
    }
}

func (s *ProductStorage) WithProductAggregation(orderBy []string) func(*gorm.DB) *gorm.DB {
    // Complex aggregation logic fully encapsulated
    needsVariantsCount := false
    needsCostPrice := false
    for _, ob := range orderBy {
        if ob == "variantsCount" || ob == "-variantsCount" {
            needsVariantsCount = true
        }
        if ob == "costPrice" || ob == "-costPrice" {
            needsCostPrice = true
        }
    }

    return func(db *gorm.DB) *gorm.DB {
        if !needsVariantsCount && !needsCostPrice {
            return db
        }

        var aggFields []string
        if needsVariantsCount {
            aggFields = append(aggFields, "COUNT(*)::int as variants_count")
        }
        if needsCostPrice {
            aggFields = append(aggFields, "AVG(cost_price)::numeric as avg_cost_price")
        }

        joinSQL := fmt.Sprintf(`
            LEFT JOIN LATERAL (
                SELECT %s
                FROM %s
                WHERE %s.%s = %s.%s
            ) AS product_agg ON true
        `, strings.Join(aggFields, ", "),
            VariantTable,
            VariantTable, VariantSchema.ProductID,
            ProductTable, ProductSchema.ID,
        )

        return db.Joins(joinSQL).Group(fmt.Sprintf("%s.%s", ProductTable, ProductSchema.ID))
    }
}
```

### Checklist for List/Count methods

When implementing or reviewing List/Count methods:

✅ Service method is under 50 lines (ideally 20-30)  
✅ Service only builds scope arrays and calls repository  
✅ No SQL fragments constructed in service  
✅ No JOINs constructed in service  
✅ No GROUP BY/HAVING constructed in service  
✅ No magic strings for tables/columns  
✅ Complex query logic lives in storage scope methods  
✅ Storage scope methods use schema constants  
✅ Each scope method has a single, clear purpose  
❌ Service does not call `db.Where()`, `db.Joins()`, `db.Group()`, `db.Having()` directly  
❌ Service does not build SQL with string concatenation/formatting  
❌ Service does not contain complex conditional logic for query construction

### Common anti-patterns to avoid

1. **Complex conditional query building in service:**

   ```go
   // 🚫 BAD: Don't do this in service
   needsVariantJoin := false
   for _, ob := range req.OrderBy() {
       if ob == "variantsCount" { needsVariantJoin = true }
   }
   if needsVariantJoin {
       // Build complex JOIN...
   }
   ```

2. **String concatenation for SQL fragments:**

   ```go
   // 🚫 BAD: Don't build SQL with strings in service
   conditions := []string{}
   for _, platform := range filters.SocialPlatforms {
       conditions = append(conditions, "customers."+platform+"_username IS NOT NULL")
   }
   query := "(" + strings.Join(conditions, " OR ") + ")"
   ```

3. **Manual GROUP BY construction:**

   ```go
   // 🚫 BAD: Don't construct GROUP BY in service
   listExtra = append(listExtra, func(db *gorm.DB) *gorm.DB {
       return db.Group(strings.Join(groupByColumns, ", "))
   })
   ```

4. **Inline LATERAL JOIN construction:**
   ```go
   // 🚫 BAD: Don't build complex JOINs in service
   joinSQL := fmt.Sprintf(`
       LEFT JOIN LATERAL (
           SELECT COUNT(*) FROM orders WHERE orders.customer_id = customers.id
       ) AS customer_agg ON true
   `)
   ```

### Benefits of correct pattern

1. **Maintainability:** Query logic centralized in storage layer
2. **Testability:** Service methods are simple and easy to unit test
3. **Reusability:** Storage scope methods can be composed and reused
4. **Consistency:** All List methods follow the same pattern
5. **No magic strings:** Schema constants enforced at storage layer
6. **Clean architecture:** Clear separation between service and persistence

## Transactions (atomic processor)

Pattern:

- Provide an “atomic processor” abstraction that runs a callback inside a transaction.
- Support:
  - isolation level selection
  - retries for retryable errors
  - reuse of an outer transaction if one is already in the context

Strict rules:

- Multi-entity writes must run inside an atomic transaction.
- Do not manually call begin/commit/rollback in domain services.

## Auth (Bearer access JWT + rotating refresh token sessions)

Pattern:

- Access token: short-lived JWT passed via `Authorization: Bearer <token>`.
- Refresh token: long-lived opaque token stored **hashed** server-side and rotated on refresh.
- Invalidate access tokens by including a server-checked version (`authVersion`) in JWT claims.

Strict rules:

- Never accept auth tokens via cookies unless the project explicitly chooses a cookie strategy.
- Treat refresh tokens like passwords: never log, store only hashes, revoke on use when rotating.

## Multi-tenancy (when applicable)

If the backend is multi-tenant, enforce boundaries at multiple layers:

- Middleware must load the authenticated actor and the tenant scope.
- Services must scope every query by tenant id.

Strict rules:

- Never trust tenant/workspace ids from URL params for authorization decisions.
- Provide “scoped getters” that prevent ID probing (BOLA), e.g., `GetWorkspaceUserByID(tenantID, userID)`.

## Testing strategy

Pattern:

- Prefer E2E/integration tests that boot real dependencies (DB/cache) for critical paths.
- Keep unit tests focused and isolated (mock external integrations).

Strict rules:

- Tests must be isolated and idempotent (truncate/cleanup between tests).
- Do not use raw SQL in tests when a domain storage/service exists.

## Clean architecture layering (separation of concerns)

The backend follows a strict layering pattern to maintain clean separation of concerns:

### Layer 1: Platform (Infrastructure)

Located in `internal/platform/**`:

- **Database repository** (`database/repository.go`): generic CRUD + scopes
- **Schema mappings** (`types/schema`): JSON ↔ DB column translation
- **Atomic transactions** (`database/atomic.go`): transaction management
- **Request/response** (`request/`, `response/`): HTTP I/O primitives
- **Problem JSON** (`types/problem`): RFC7807 error types

**Rules**:

- Platform knows nothing about domains
- Platform provides reusable primitives
- No business logic in platform code

### Layer 2: Domain Storage

Located in `internal/domain/<name>/storage.go`:

- Creates typed repositories: `database.NewRepository[Order](db)`
- Configures search indexes and caching
- Exposes repository instances via `Storage` struct

**Rules**:

- Storage layer is data access only
- No business logic (no validation, no authorization)
- No direct queries from services; services call repository methods

### Layer 3: Domain Service

Located in `internal/domain/<name>/service.go`:

- **Business invariants live here** (validation, authorization, state transitions)
- Orchestrates storage operations
- Calls atomic processor for transactions
- Emits domain events via bus

**Rules** (strict separation):

- ✅ Services call `storage.repo.FindOne(ctx, scopes...)`
- ✅ Services call `storage.repo.ScopeBusinessID(bizID)` to build scopes
- ❌ Services **never** construct queries with `db.Where("column = ?", val)`
- ❌ Services **never** reference table names as strings
- ❌ Services **never** call `db.Preload("Association")` directly

**Example of correct service pattern**:

```go
// ✅ GOOD: Service uses repository + scopes
func (s *Service) GetOrder(ctx context.Context, actor *account.User, biz *business.Business, orderID string) (*Order, error) {
    return s.storage.order.FindOne(ctx,
        s.storage.order.ScopeBusinessID(biz.ID),
        s.storage.order.ScopeID(orderID),
        s.storage.order.WithPreload(OrderItemStruct),
        s.storage.order.WithPreload(OrderNoteStruct),
    )
}

// ❌ BAD: Service constructs query
func (s *Service) GetOrder(ctx context.Context, biz *business.Business, orderID string) (*Order, error) {
    var order Order
    err := s.storage.db.Where("business_id = ? AND id = ?", biz.ID, orderID).
        Preload("Items").
        Preload("Notes").
        First(&order).Error
    return &order, err
}
```

### Layer 4: HTTP Handler

Located in `internal/domain/<name>/handler_http.go`:

- Parse/validate request (via `request.ValidBody`)
- Call service method
- Map result to response DTO (from `model_response.go`)
- Return via `response.SuccessJSON` or `response.Error`

**Rules**:

- Handlers are thin (no business logic)
- Always map GORM models to response DTOs
- Never return `gorm.Model` embedded structs
- Use response DTOs from `model_response.go`

**Example of correct handler pattern**:

```go
// ✅ GOOD: Handler delegates to service + maps to DTO
func (h *HttpHandler) GetOrder(c *gin.Context) {
    orderID := c.Param("orderId")
    actor := account.ActorFromContext(c)
    biz := business.BusinessFromContext(c)

    order, err := h.service.GetOrder(c.Request.Context(), actor, biz, orderID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.SuccessJSON(c, http.StatusOK, ToOrderResponse(order))
}

// ❌ BAD: Handler returns GORM model directly
func (h *HttpHandler) GetOrder(c *gin.Context) {
    order, err := h.service.GetOrder(...)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, order) // Exposes gorm.Model, wrong casing, etc.
}
```

### Why This Matters

**Maintainability**:

- Clear responsibilities per layer
- Easy to find where logic lives
- Services can be unit tested without HTTP

**Refactoring safety**:

- Can change storage implementation without touching services
- Can change DB schema without touching handlers
- Constants provide compile-time safety

**Security**:

- Repository scopes enforce tenant isolation at the data layer
- Services enforce authorization before calling storage
- Handlers can't bypass service-layer checks

## Public reference implementation (Kyora)

If you need a concrete example of this pattern as implemented today:

- Entry + lifecycle: `backend/main.go`, `backend/cmd/root.go`, `backend/cmd/server.go`
- DI + engine setup: `backend/internal/server/server.go`
- Routes + middleware: `backend/internal/server/routes.go`
- Auth middleware + JWT: `backend/internal/platform/auth/*`
- Request/response/problem: `backend/internal/platform/request/*`, `backend/internal/platform/response/*`, `backend/internal/platform/types/problem/*`
- DB repo + atomic tx: `backend/internal/platform/database/*`

## Anti-patterns (avoid)

- Putting business rules in handlers/routes.
- Returning ad-hoc error JSON.
- Accepting tenant/workspace ids from clients for authorization.
- Building SQL strings from user input.
- Using `float64` for money.
- Logging secrets.
- Calling another domain’s storage layer directly (bypass service invariants).
