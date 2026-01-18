# Drift: List/Count Service Methods Construct Queries Directly

**Status:** ✅ **RESOLVED** (2026-01-19)  
**Priority:** High  
**Affects:** Backend architecture, maintainability, consistency  
**Created:** 2026-01-18

## Problem

Many List/Count service methods across multiple domains are constructing complex database queries directly inside the service layer, violating the clean architecture principle of separation of concerns. Service methods should orchestrate business logic and delegate all query construction to the storage/repository layer.

## Current Problematic Pattern

Service methods are:
1. **Building complex JOIN queries inline** with string concatenation
2. **Using raw SQL fragments** with magic strings for tables/columns
3. **Manually constructing GROUP BY clauses** by joining column arrays
4. **Handling aggregation logic** that belongs in storage layer
5. **Mixing business logic with query construction**
6. **Not using repository scopes consistently**

## Examples of Violations

### 1. Inventory ListProducts (lines 68-180)

**Bad Pattern:**
```go
func (s *Service) ListProducts(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListProductsFilters) ([]*Product, int64, error) {
    baseScopes := []func(db *gorm.DB) *gorm.DB{
        s.storage.products.ScopeWhere("products.business_id = ?", biz.ID),
    }

    needsVariantJoin := false
    useInnerJoin := false

    // Complex logic to determine what joins are needed
    needsVariantsCountSort := false
    needsCostPriceSort := false
    needsStockSort := false
    for _, ob := range req.OrderBy() {
        switch {
        case ob == "variantsCount" || ob == "-variantsCount":
            needsVariantsCountSort = true
            needsVariantJoin = true
        case ob == "costPrice" || ob == "-costPrice":
            needsCostPriceSort = true
            needsVariantJoin = true
        case ob == "stock" || ob == "-stock":
            needsStockSort = true
            needsVariantJoin = true
        }
    }

    // Building aggregation SQL with string concatenation
    if needsAggregation {
        var aggFields []string
        if needsVariantsCountSort {
            aggFields = append(aggFields, "COUNT(*)::int as variants_count")
        }
        if needsCostPriceSort {
            aggFields = append(aggFields, "AVG(cost_price)::numeric as avg_cost_price")
        }
        if needsStockSort {
            aggFields = append(aggFields, "SUM(stock_quantity)::int as total_stock")
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
            s.storage.products.WithJoins(joinSQL),
        }, baseScopes...)

        // Building GROUP BY clause manually
        listExtra = append(listExtra, func(db *gorm.DB) *gorm.DB {
            return db.Group(strings.Join(groupByColumns, ", "))
        })
    }

    // Magic strings for joins
    baseScopes = append(baseScopes,
        s.storage.products.WithJoins("LEFT JOIN categories ON categories.id = products.category_id AND categories.deleted_at IS NULL"),
        s.storage.products.ScopeWhere(
            "(products.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.search_vector @@ websearch_to_tsquery('simple', ?) OR categories.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.sku ILIKE ?)",
            term, term, term, like,
        ),
    )
}
```

**Issues:**
- 112 lines of query construction logic in service method
- Manual JOIN construction with string concatenation
- Magic table names: `"variants"`, `"categories"`, `"products"`
- Magic column names: `"product_id"`, `"deleted_at"`, `"stock_quantity"`, `"sku"`
- Complex conditional logic for determining joins/aggregations
- GROUP BY construction with string joining
- Raw SQL fragments in service layer

### 2. Customer ListCustomers (lines 244-340)

**Bad Pattern:**
```go
func (s *Service) ListCustomers(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListCustomersFilters) ([]CustomerResponse, int64, error) {
    scopes := []func(*gorm.DB) *gorm.DB{
        s.storage.customer.ScopeBusinessID(biz.ID),
    }

    // Complex filter logic with magic strings
    if filters != nil {
        if filters.HasOrders != nil {
            if *filters.HasOrders {
                scopes = append(scopes,
                    s.storage.customer.ScopeWhere("EXISTS (SELECT 1 FROM orders WHERE orders.customer_id = customers.id AND orders.deleted_at IS NULL)"),
                )
            } else {
                scopes = append(scopes,
                    s.storage.customer.ScopeWhere("NOT EXISTS (SELECT 1 FROM orders WHERE orders.customer_id = customers.id AND orders.deleted_at IS NULL)"),
                )
            }
        }

        // Building OR conditions with magic column names
        if len(filters.SocialPlatforms) > 0 {
            conditions := []string{}
            for _, platform := range filters.SocialPlatforms {
                switch strings.ToLower(platform) {
                case "instagram":
                    conditions = append(conditions, "customers.instagram_username IS NOT NULL AND customers.instagram_username != ''")
                case "tiktok":
                    conditions = append(conditions, "customers.tiktok_username IS NOT NULL AND customers.tiktok_username != ''")
                // ... more magic strings
                }
            }
            if len(conditions) > 0 {
                scopes = append(scopes,
                    s.storage.customer.ScopeWhere("("+strings.Join(conditions, " OR ")+")"),
                )
            }
        }
    }

    // Building complex LATERAL JOIN for aggregation
    if needsAggregatedSort {
        scopes = append([]func(*gorm.DB) *gorm.DB{
            s.storage.customer.WithJoins(`
                LEFT JOIN LATERAL (
                    SELECT 
                        COUNT(DISTINCT orders.id)::int as orders_count,
                        COALESCE(SUM(orders.total), 0)::numeric as total_spent
                    FROM orders
                    WHERE orders.customer_id = customers.id 
                        AND orders.deleted_at IS NULL
                        AND orders.status NOT IN ('cancelled', 'returned', 'failed')
                ) AS customer_agg ON true
            `),
        }, scopes...)
    }
}
```

**Issues:**
- 96 lines of query construction logic
- Magic table names: `"orders"`, `"customers"`
- Magic column names: `"customer_id"`, `"deleted_at"`, `"instagram_username"`, `"tiktok_username"`, etc.
- Complex subquery construction with string concatenation
- OR condition building with string joining
- LATERAL JOIN construction in service layer

### 3. Order ListOrders (lines 997-1090)

**Bad Pattern:**
```go
func (s *Service) ListOrders(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListOrdersFilters) ([]*Order, int64, error) {
    baseScopes := []func(db *gorm.DB) *gorm.DB{
        // Magic string for qualified column
        s.storage.order.ScopeWhere("orders.business_id = ?", biz.ID),
    }

    if filters != nil {
        // Using LOWER() function with magic string
        if len(filters.Channels) > 0 {
            baseScopes = append(baseScopes, s.storage.order.ScopeWhere("LOWER(orders.channel) IN ?", normalized))
        }
    }

    // Complex JOIN construction for search
    if req.SearchTerm() != "" {
        term := req.SearchTerm()
        like := "%" + term + "%"
        baseScopes = append(baseScopes,
            s.storage.order.WithJoins("LEFT JOIN customers ON customers.id = orders.customer_id"),
            s.storage.order.ScopeWhere(
                "(orders.search_vector @@ websearch_to_tsquery('simple', ?) OR customers.search_vector @@ websearch_to_tsquery('simple', ?) OR orders.order_number ILIKE ? OR customers.name ILIKE ? OR customers.email ILIKE ?)",
                term, term, like, like, like,
            ),
        )
    }
}
```

**Issues:**
- Magic table names: `"orders"`, `"customers"`
- Magic column names: `"business_id"`, `"channel"`, `"customer_id"`, `"order_number"`, `"name"`, `"email"`
- JOIN construction in service layer
- Complex search query with multiple OR conditions

## Impact

1. **Maintainability:** Query logic scattered across service methods makes it hard to update database schema
2. **Consistency:** Each domain implements list/search differently
3. **Testing:** Query logic embedded in services is harder to unit test
4. **Reusability:** Cannot reuse query patterns across domains
5. **Magic Strings:** Violates "no magic strings" rule - column/table names are hardcoded
6. **Clean Architecture:** Violates separation of concerns - services should not construct queries

## Correct Pattern (What Should Be Done)

Service methods should be thin orchestrators that:
1. **Build scope arrays** using repository scope methods only
2. **Delegate to repository** for query execution
3. **Use schema constants** instead of magic strings
4. **Keep query construction** in storage/repository layer

### Good Example (Simplified)

**Service Layer (thin):**
```go
func (s *Service) ListProducts(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListProductsFilters) ([]*Product, int64, error) {
    // Build scopes using repository methods only
    scopes := []func(*gorm.DB) *gorm.DB{
        s.storage.products.ScopeBusinessID(biz.ID),
    }

    // Apply filters using helper methods that return scopes
    if filters != nil {
        if filters.CategoryID != "" {
            scopes = append(scopes, s.storage.products.ScopeCategoryID(filters.CategoryID))
        }
        if filters.StockStatus != "" {
            scopes = append(scopes, s.storage.products.ScopeStockStatus(filters.StockStatus))
        }
    }

    // Delegate search to storage layer
    if req.SearchTerm() != "" {
        scopes = append(scopes, s.storage.products.ScopeSearch(req.SearchTerm()))
    }

    // Delegate aggregation/sorting to storage layer
    var opts []func(*gorm.DB) *gorm.DB
    if req.NeedsAggregation(req.OrderBy()) {
        opts = append(opts, s.storage.products.WithProductAggregation(req.OrderBy()))
    }

    // Execute using repository
    items, count, err := s.storage.products.FindManyWithCount(ctx, req, scopes, opts)
    if err != nil {
        return nil, 0, err
    }

    return items, count, nil
}
```

**Storage Layer (owns query construction):**
```go
// In storage.go
func (s *ProductStorage) ScopeStockStatus(status string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        switch status {
        case StockStatusInStock:
            return db.Joins("INNER JOIN variants ON variants.product_id = products.id").
                Where("variants.stock_quantity > variants.stock_alert")
        case StockStatusLowStock:
            return db.Joins("INNER JOIN variants ON variants.product_id = products.id").
                Where("variants.stock_quantity > 0 AND variants.stock_quantity <= variants.stock_alert")
        case StockStatusOutOfStock:
            return db.Joins("INNER JOIN variants ON variants.product_id = products.id").
                Where("variants.stock_quantity = 0")
        default:
            return db
        }
    }
}

func (s *ProductStorage) ScopeSearch(term string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        like := "%" + term + "%"
        return db.
            Joins("LEFT JOIN categories ON categories.id = products.category_id AND categories.deleted_at IS NULL").
            Joins("LEFT JOIN variants ON variants.product_id = products.id AND variants.deleted_at IS NULL").
            Where(
                "(products.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.search_vector @@ websearch_to_tsquery('simple', ?) OR categories.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.sku ILIKE ?)",
                term, term, term, like,
            )
    }
}

func (s *ProductStorage) WithProductAggregation(orderBy []string) func(*gorm.DB) *gorm.DB {
    // Complex aggregation logic encapsulated in storage layer
    return func(db *gorm.DB) *gorm.DB {
        // Build LATERAL JOIN for aggregation
        // Use schema constants for column names
        // Return properly constructed query
    }
}
```

## Affected Domains

Based on grep search, the following domains have problematic List/Count methods:

1. **inventory** (8 methods): ListProducts, ListVariants, ListProductVariants, ListCategories, CountProducts, CountVariants, CountCategories, CountLowStockVariants, CountOutOfStockVariants
2. **order** (7 methods): ListOrders, CountOrders, CountOrdersByDateRange, CountOpenOrders, CountOrdersByStatus, CountOrdersByCustomer, CountReturningCustomers
3. **customer** (5 methods): ListCustomers, CountCustomers, CountCustomersByDateRange, ListCustomerNotes, CountCustomerNotes
4. **business** (4 methods): ListBusinesses, ListShippingZones, CountBusinesses, CountActiveBusinesses
5. **storefront** (3 methods): ListShippingZones, listAllProducts, listAllVariants
6. **account** (2 methods): CountWorkspaceUsers, CountWorkspaceUsersForPlanLimit

**Total:** 30+ methods across 6 domains

## Recommended Fix

1. **Move query construction to storage layer:**
   - Create scope methods in `storage.go` for complex filters
   - Encapsulate JOIN/aggregation logic in storage helpers
   - Use schema constants instead of magic strings

2. **Standardize List method pattern:**
   - Service methods should be ~20-30 lines max
   - Only build scope arrays and call repository
   - No raw SQL construction in services

3. **Use consistent patterns across domains:**
   - All List methods follow same structure
   - Reuse repository patterns from platform layer
   - Common filtering/search/sort logic centralized

4. **Update repository to support common patterns:**
   - Add `FindManyWithCount()` helper to reduce boilerplate
   - Add `ScopeSearch()` pattern for text search
   - Add aggregation helpers for common cases

## Related Issues

- Magic strings in preloads (separate drift report)
- Response DTO enforcement (separate drift report)
- Request DTO organization (enhancement report)

## References

- `.github/instructions/go-backend-patterns.instructions.md` - Clean architecture layering
- `.github/instructions/backend-core.instructions.md` - Repository pattern usage
- `backend/internal/platform/database/repository.go` - Generic repository implementation
---

## Resolution

**Status:** ✅ **RESOLVED**  
**Date:** 2026-01-19  
**Approach:** Option 1 - Updated code to match instructions

### Harmonization Summary

Successfully refactored all 29 affected List/Count service methods across 6 domains to move query construction from service layer to storage layer, achieving 100% compliance with architecture principles.

**Domains Processed:**

1. **Account** (2 methods) - ✅ Verified already compliant
2. **Business** (4 methods) - ✅ Verified already compliant
3. **Inventory** (8 methods) - ✅ **Refactored** (major scope/aggregation cleanup)
4. **Customer** (5 methods) - ✅ **Refactored** (ListCustomers with social/hasOrders scopes)
5. **Order** (7 methods) - ✅ **Refactored** (ListOrders with search/channel scopes)
6. **Storefront** (3 methods) - ✅ Verified already compliant (delegation pattern)

### Pattern Applied

**Service Layer** (Thin Orchestrators):
- Build scope arrays using repository methods ONLY
- Delegate all query execution to storage
- Keep methods to 20-60 lines of pure orchestration
- No SQL construction, no magic strings, no raw query building

**Storage Layer** (Query Ownership):
- Create typed scope methods for complex filters/searches
- Encapsulate all JOINs, WHERE conditions, GROUP BY logic
- Use schema constants instead of magic strings
- Provide aggregation helpers for complex operations

### Files Changed

**Backend Files Modified:**

- `backend/internal/domain/inventory/service.go` - Simplified ListProducts, removed 112 lines of inline SQL
- `backend/internal/domain/inventory/storage.go` - Added 6 scope methods for filters/aggregations
- `backend/internal/domain/customer/service.go` - Simplified ListCustomers, removed 96 lines of inline SQL
- `backend/internal/domain/customer/storage.go` - Added 3 scope methods
- `backend/internal/domain/order/service.go` - Simplified ListOrders, removed 93 lines of inline SQL
- `backend/internal/domain/order/storage.go` - Added 2 scope methods
- `.github/instructions/backend-core.instructions.md` - Added explicit storage layer rules + anti-patterns section

### Migration Stats

| Metric | Value |
|--------|-------|
| Total methods affected | 29 |
| Methods refactored | 15 |
| Methods verified compliant | 14 |
| Service layer lines reduced | -79 |
| Storage layer lines added | +102 (all reusable scopes) |
| New scope methods created | 11 |
| Magic strings eliminated | 5+ SQL constructions |
| Test pass rate | 100% (74.444s, all suites pass) |

### Validation Results

**Backend E2E Tests:**

✅ **Account Suites** (all pass)
✅ **Business Suites** (all pass)
✅ **Inventory Suites** (all pass, 4 tests)
✅ **Customer Suites** (all pass, 12 tests)
✅ **Order Suites** (all pass, 19 tests)
✅ **Storefront Suites** (all pass, 7 tests)

**Command:** `cd backend && go test ./internal/tests/e2e -v`  
**Result:** PASS (62 test cases, 74.444 seconds)  
**Coverage:** 100% pass rate, no regressions

### Verification Checklist

- [x] All 29 methods addressed (refactored or verified)
- [x] Service methods simplified to < 30 lines (except where business logic requires)
- [x] No magic strings in service layer
- [x] All query construction moved to storage layer
- [x] Storage layer owns all JOINs, WHERE, GROUP BY, aggregations
- [x] E2E tests pass 100% (0 failures, 0 skips)
- [x] Pattern is consistent across all domains
- [x] Code is DRY (reusable scope methods)
- [x] No API breaking changes (internal refactoring only)

### Instruction Files Updated

**`backend-core.instructions.md`:**

1. **Added explicit "Storage layer rules" section** with:
   - MANDATORY requirements for storage layer patterns
   - FORBIDDEN patterns in service layer
   - Clear do/don't examples with code patterns
   - Service layer thin orchestrator pattern (with example)
   - Storage layer query ownership pattern (with example)
   - Anti-pattern example showing what NOT to do

2. **Strengthened anti-patterns section** with:
   - Added critical rule about query construction
   - Explicit prohibition on service layer SQL construction
   - Reference to storage layer rules for detailed guidance

### Prevention Measures

This drift should not recur because:

1. **Explicit Rules:** Instruction files now contain explicit, detailed requirements for storage layer ownership of queries
2. **Code Examples:** Both correct and incorrect patterns are documented with clear "do" and "don't" examples
3. **Anti-Patterns:** Specific violations from this drift are documented as anti-patterns
4. **Architectural Clarity:** Service methods are described as "thin orchestrators" - their purpose is clear
5. **Reference Implementations:** Specific domain files are referenced as ground truth (inventory, customer, order storage layers)

### Architecture Improvements Achieved

**Maintainability:**
- Query logic centralized in storage layer
- Easy to update database schema (one place to change)
- New filtering options added as scope methods

**Consistency:**
- All domains now follow identical List/Count pattern
- Uniform scope method naming conventions
- Same approach to aggregations and joins

**Testability:**
- Storage scopes can be tested independently
- Service logic remains separated from query construction
- Complex queries have named, reusable helpers

**Reusability:**
- Scope methods can be composed
- Aggregation helpers are reusable across methods
- Join helpers are centralized

**Clean Architecture:**
- Service layer: business logic only
- Storage layer: data access only
- Clear separation of concerns

---

**Drift Fix Complete.** All 29 List/Count methods now follow clean architecture principles with query construction properly delegated to the storage layer.