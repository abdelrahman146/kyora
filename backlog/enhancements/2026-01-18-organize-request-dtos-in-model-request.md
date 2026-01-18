---
title: "Request DTOs Not Organized in Dedicated model_request.go Files"
date: 2026-01-18
priority: low
category: organization
status: open
domain: backend
---

# Request DTOs Not Organized in Dedicated model_request.go Files

## Summary

Request DTOs (Create/Update request structs) are currently defined in `model.go` files alongside GORM models, constants, and schemas. While this works for simple domains, the updated `go-backend-patterns.instructions.md` recommends organizing request DTOs in dedicated `model_request.go` files for complex domains to improve separation of concerns and maintainability.

## Current State

All domains keep request DTOs in `model.go`:

### Locations Found

1. **backend/internal/domain/order/model.go** (lines 135-175):
   ```go
   type CreateOrderRequest struct { ... }
   type UpdateOrderRequest struct { ... }
   type AddOrderPaymentDetailsRequest struct { ... }
   ```

2. **backend/internal/domain/inventory/model.go** (lines ~100-250):
   ```go
   type CreateProductRequest struct { ... }
   type UpdateProductRequest struct { ... }
   type CreateVariantRequest struct { ... }
   type UpdateVariantRequest struct { ... }
   type CreateCategoryRequest struct { ... }
   type UpdateCategoryRequest struct { ... }
   ```

3. **backend/internal/domain/customer/model.go**:
   ```go
   type CreateCustomerRequest struct { ... }
   type UpdateCustomerRequest struct { ... }
   type CreateCustomerAddressRequest struct { ... }
   type UpdateCustomerAddressRequest struct { ... }
   type CreateCustomerNoteRequest struct { ... }
   type UpdateCustomerNoteRequest struct { ... }
   ```

4. **backend/internal/domain/accounting/model.go**:
   ```go
   type CreateAssetRequest struct { ... }
   type UpdateAssetRequest struct { ... }
   type CreateInvestmentRequest struct { ... }
   // ... many more request types
   ```

5. **Other domains**: Similar pattern across all domains

## Expected State (Recommendation)

Per the updated `go-backend-patterns.instructions.md`:

For **simple domains** (< 5 request types):
- Keeping requests in `model.go` is acceptable

For **complex domains** (≥ 5 request types):
- Create `model_request.go` and move all request DTOs there
- Keeps API contract separate from domain models
- Improves file organization and navigation

Example structure for order domain:
```go
// backend/internal/domain/order/model.go
- GORM models (Order, OrderItem, OrderNote)
- Enums and constants
- Schema definitions

// backend/internal/domain/order/model_request.go
- CreateOrderRequest
- UpdateOrderRequest
- AddOrderPaymentDetailsRequest
- CreateOrderItemRequest
- CreateOrderNoteRequest
- UpdateOrderNoteRequest

// backend/internal/domain/order/model_response.go (already exists)
- OrderResponse
- OrderItemResponse
- OrderNoteResponse
```

## Impact

**Current impact is minimal:**
- Functionally correct (no bugs)
- Slightly harder to navigate large `model.go` files
- Mixed concerns (data model + API contract + schema definitions)

**Benefits of organizing:**
- Clearer separation between internal models and API contract
- Easier to find request validation rules
- Better IDE navigation (jump to request definition)
- Consistent with response DTO pattern (already in `model_response.go`)

## Suggested Fix

### Phase 1: Identify complex domains (≥ 5 request types)

Domains to prioritize:
- `order` (6+ request types)
- `inventory` (6+ request types)
- `customer` (6+ request types)
- `accounting` (10+ request types)
- `business` (5+ request types)

### Phase 2: Create model_request.go files

For each complex domain:

1. Create `model_request.go`
2. Move all `*Request` structs from `model.go`
3. Verify imports (may need to import other domain types)
4. Run tests to ensure no breakage

### Phase 3: Update instructions (already done)

The `go-backend-patterns.instructions.md` now documents:
- `model.go`: models, constants, schemas (simple domains may keep requests here)
- `model_request.go`: recommended for complex domains
- `model_response.go`: required for all domains

## Priority Justification

**Low** because:
- No functional issues
- Purely organizational improvement
- Current pattern is explicitly allowed for simple domains
- Can be done incrementally per domain
- Low risk (moving files, no logic changes)

## Related Instructions

- `.github/instructions/go-backend-patterns.instructions.md` (updated with this pattern)
- `.github/instructions/backend-core.instructions.md` (Kyora-specific structure)
