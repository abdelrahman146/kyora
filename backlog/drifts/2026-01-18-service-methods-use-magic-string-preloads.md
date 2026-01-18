---
title: "Service Methods Use Magic String Preloads Instead of Constants"
date: 2026-01-18
priority: medium
category: consistency
status: open
domain: backend
---

# Service Methods Use Magic String Preloads Instead of Constants

## Summary

Multiple service methods use magic string literals for GORM preload associations instead of defined struct name constants. This violates the separation of concerns principle and the "no magic strings" rule that should be enforced in backend patterns.

## Current State

Service methods across multiple domains use string literals for preloads:

### Locations Found

1. **backend/internal/domain/order/service.go** (lines 50-54, 56):
   ```go
   func (s *Service) orderListPreloads() []func(*gorm.DB) *gorm.DB {
       return []func(*gorm.DB) *gorm.DB{
           s.storage.order.WithPreload(customer.CustomerStruct), // ✅ GOOD - uses constant
           s.storage.order.WithPreload("ShippingAddress"),      // ❌ BAD - magic string
           s.storage.order.WithPreload("ShippingZone"),         // ❌ BAD - magic string
           s.storage.order.WithPreload(OrderItemStruct),        // ✅ GOOD - uses constant
           s.storage.order.WithPreload("Items.Product"),        // ❌ BAD - magic string
           s.storage.order.WithPreload("Items.Variant"),        // ❌ BAD - magic string
           s.storage.order.WithPreload(OrderNoteStruct),        // ✅ GOOD - uses constant
       }
   }
   ```

2. **backend/internal/domain/customer/service.go** (line 50):
   ```go
   func (s *Service) GetCustomerByID(...) (*Customer, error) {
       return s.storage.customer.FindOne(ctx,
           s.storage.customer.ScopeBusinessID(biz.ID),
           s.storage.customer.ScopeID(id),
           s.storage.customer.WithPreload("Notes"), // ❌ BAD - magic string
       )
   }
   ```

3. **backend/internal/domain/account/service.go** (line 66):
   ```go
   func (s *Service) GetWorkspaceByID(ctx context.Context, id string) (*Workspace, error) {
       return s.storage.workspace.FindByID(ctx, id, s.storage.workspace.WithPreload("Users")) // ❌ BAD
   }
   ```

4. **backend/internal/domain/account/service.go** (line 760):
   ```go
   scopes = append(scopes, s.storage.invitation.WithPreload("Inviter")) // ❌ BAD
   ```

5. **backend/internal/domain/inventory/service.go** (lines 43, 257):
   ```go
   s.storage.products.WithPreload("Variants") // ❌ BAD - should use VariantStruct constant
   ```

### Additional Magic String Pattern

6. **backend/internal/domain/inventory/service.go** (lines 366, 377):
   ```go
   s.storage.variants.ScopeWhere("variants.business_id = ?", biz.ID) // ❌ BAD - magic string table/column
   s.storage.variants.WithJoins("LEFT JOIN products ON products.id = variants.product_id AND products.deleted_at IS NULL") // ❌ BAD
   ```

## Expected State

Per the updated `go-backend-patterns.instructions.md`:

1. **All struct references in preloads must use constants** defined in `model.go`:
   ```go
   const (
       OrderItemStruct = "OrderItem"
       OrderNoteStruct = "OrderNote"
       ShippingAddressStruct = "ShippingAddress"
       ShippingZoneStruct = "ShippingZone"
   )
   ```

2. **All table/column references must use schema fields**:
   ```go
   s.storage.variants.ScopeBusinessID(biz.ID) // Use built-in scope helper
   // OR
   s.storage.variants.ScopeEquals(VariantSchema.BusinessID, biz.ID) // Use schema field
   ```

3. **Services must never construct raw SQL strings** for table/column names.

## Impact

- **Maintainability**: Harder to refactor struct/field names
- **Consistency**: Mixed usage of constants vs magic strings across codebase
- **Type Safety**: No compile-time checking for preload association names
- **Search/Replace**: Can't reliably find all usages of a struct name
- **Documentation**: Pattern is unclear to new contributors

## Suggested Fix

### Step 1: Define missing constants in model.go files

For each domain, add struct name constants for all associations:

```go
// backend/internal/domain/order/model.go
const (
    OrderTable  = "orders"
    OrderStruct = "Order"
    OrderPrefix = "ord"
    OrderItemStruct = "OrderItem"
    OrderNoteStruct = "OrderNote"
    
    // Add missing:
    ShippingAddressStruct = "ShippingAddress"
    ShippingZoneStruct = "ShippingZone"
    ItemsProductStruct = "Items.Product"  // For nested preloads
    ItemsVariantStruct = "Items.Variant"
)
```

### Step 2: Replace all magic string preloads

Search and replace pattern:
- `WithPreload("ShippingAddress")` → `WithPreload(ShippingAddressStruct)`
- `WithPreload("Notes")` → `WithPreload(CustomerNoteStruct)`
- `WithPreload("Variants")` → `WithPreload(VariantStruct)`
- etc.

### Step 3: Replace ScopeWhere with schema-based scopes

Replace:
```go
s.storage.variants.ScopeWhere("variants.business_id = ?", biz.ID)
```

With:
```go
s.storage.variants.ScopeBusinessID(biz.ID)
```

Or when no helper exists:
```go
s.storage.variants.ScopeEquals(VariantSchema.BusinessID, biz.ID)
```

### Step 4: Add E2E test coverage

Add a test that verifies preload associations work correctly (this will catch typos in constants during compile time).

## Related Issues

- [2026-01-18-gorm-model-exposed-in-responses.md](./2026-01-18-gorm-model-exposed-in-responses.md) - Response DTO pattern drift

## Priority Justification

**Medium** because:
- Does not break functionality (magic strings work at runtime)
- Affects maintainability and consistency
- Should be fixed before patterns become more entrenched
- Relatively easy to fix with search/replace + constant definitions
