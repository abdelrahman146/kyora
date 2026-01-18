# Drift: List/Count Service Methods Construct Queries Directly

**Status:** Open  
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
