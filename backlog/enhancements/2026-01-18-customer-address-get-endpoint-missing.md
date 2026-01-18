---
title: "Customer API Missing GET Endpoint for Individual Customer Addresses"
date: 2026-01-18
priority: low
category: api-completeness
status: open
domain: backend
---

# Customer API Missing GET Endpoint for Individual Customer Addresses

## Summary

Backend customer API provides `GET /customers/:customerId/addresses` to list all addresses but does not provide `GET /customers/:customerId/addresses/:addressId` to retrieve a single address. This creates an inconsistent REST API pattern.

## Current State

### Existing Address Endpoints (`backend/internal/server/routes.go`)

```go
addressGroup := customers.Group("/:customerId/addresses")
{
    addressGroup.GET("", ..., customerHandler.ListCustomerAddresses)     // List all
    addressGroup.POST("", ..., customerHandler.CreateCustomerAddress)    // Create
    addressGroup.PATCH("/:addressId", ..., customerHandler.UpdateCustomerAddress)  // Update
    addressGroup.DELETE("/:addressId", ..., customerHandler.DeleteCustomerAddress) // Delete
}
```

Missing: `addressGroup.GET("/:addressId", ..., customerHandler.GetCustomerAddress)` // Get single address

### Comparison with Other Resources

**Customer Notes** (same resource structure):
```go
noteGroup.GET("", ..., customerHandler.ListCustomerNotes)     // List all
// ← No GET /:noteId either
noteGroup.POST("", ..., customerHandler.CreateCustomerNote)   // Create
noteGroup.DELETE("/:noteId", ..., customerHandler.DeleteCustomerNote) // Delete
```

**Orders** (for reference):
```go
orders.GET("", ..., orderHandler.ListOrders)                  // List all
orders.GET("/:orderId", ..., orderHandler.GetOrder)           // Get single ← This exists
orders.POST("", ..., orderHandler.CreateOrder)                // Create
orders.PATCH("/:orderId", ..., orderHandler.UpdateOrder)      // Update
orders.DELETE("/:orderId", ..., orderHandler.DeleteOrder)     // Delete
```

## Expected State

Two perspectives:

**Perspective 1: Standard REST pattern**
- Add `GET /:addressId` endpoint for consistency with other resources
- Add `GET /:noteId` endpoint for consistency

**Perspective 2: Current pattern is intentional**
- Addresses are always fetched as part of customer detail (preloaded)
- Individual address retrieval is not needed because forms pre-populate from customer detail data
- Notes are also fetched as part of customer detail

## Impact

- **Low**: Portal-web doesn't need individual address GET (it prefetches customer detail with addresses).
- API pattern is slightly inconsistent with resources like Orders/Products/Variants.
- External API consumers may expect individual resource GET endpoints.

## Affected Files

- `backend/internal/domain/customer/handler_http.go` (would need `GetCustomerAddress` handler)
- `backend/internal/server/routes.go` (lines 263-267: address routes)
- `backend/internal/domain/customer/service.go` (would need corresponding service method)

## Suggested Fix

**Option 1: Add GET endpoints (if needed for API completeness)**

```go
// backend/internal/domain/customer/handler_http.go
func (h *HttpHandler) GetCustomerAddress(c *gin.Context) {
    // Implementation similar to GetCustomer
    // Enforces customer exists in business
    // Returns single address
}
```

**Option 2: Document as intentional design (recommended)**

If individual GET is not needed:
1. Document in `.github/instructions/customer.instructions.md` that addresses are always accessed via customer detail
2. Add comment in routes.go explaining why GET /:addressId is omitted
3. Consider same for notes if following same pattern

## Rationale for Current Behavior

The current pattern makes sense because:
- Addresses are lightweight and few per customer
- Customer detail always preloads addresses and notes
- Update/delete operations use the address/note ID directly without needing to fetch first
- This avoids unnecessary roundtrips

## Related

- `.github/instructions/customer.instructions.md` (documents available endpoints)
- `.github/instructions/backend-core.instructions.md` (REST patterns)
