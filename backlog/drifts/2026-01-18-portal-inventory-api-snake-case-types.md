# Portal Inventory API Uses snake_case Types (Backend Returns camelCase)

**Date:** 2026-01-18  
**Domain:** portal-web, inventory  
**Type:** Code Drift (API Contract Mismatch)  
**Impact:** High — Type definitions don't match backend responses, potential runtime errors  
**Status:** Open  

---

## Problem

`portal-web/src/api/inventory.ts` defines TypeScript interfaces using **snake_case** keys, but the backend inventory API returns **camelCase** responses.

This violates the Kyora API contract standard:
- Backend responses use camelCase (verified in `backend/internal/domain/inventory/model_response.go`)
- E2E tests verify camelCase responses (e.g., `backend/internal/tests/e2e/inventory_products_test.go` checks `pageSize`, `totalCount`)
- Other portal API clients (e.g., `portal-web/src/api/order.ts`) use camelCase correctly

---

## Current State (Incorrect)

**File:** `portal-web/src/api/inventory.ts`

### ListResponse interface (lines 60-67):

```typescript
export interface ListResponse<T> {
  items: Array<T>
  page: number
  page_size: number       // ❌ snake_case
  total_count: number     // ❌ snake_case
  total_pages: number     // ❌ snake_case
  has_more: boolean       // ❌ snake_case
}
```

### CreateVariantRequest interface (line 88):

```typescript
export interface CreateVariantRequest {
  product_id: string      // ❌ snake_case
  code: string
  // ... other fields
}
```

---

## Expected State (Backend Reality)

Backend returns:

```typescript
export interface ListResponse<T> {
  items: Array<T>
  page: number
  pageSize: number        // ✅ camelCase
  totalCount: number      // ✅ camelCase
  totalPages: number      // ✅ camelCase
  hasMore: boolean        // ✅ camelCase
}
```

Request payloads should use:

```typescript
export interface CreateVariantRequest {
  productId: string       // ✅ camelCase
  code: string
  // ... other fields
}
```

---

## Evidence

1. **Backend Response DTOs** (`backend/internal/domain/inventory/model_response.go`):
   - All fields use camelCase JSON tags: `json:"businessId"`, `json:"categoryId"`, etc.

2. **Backend E2E Tests** (`backend/internal/tests/e2e/inventory_products_test.go`, lines 104-105):
   ```go
   s.Equal(float64(2), page1["pageSize"])
   s.Equal(float64(3), page1["totalCount"])
   ```

3. **List response implementation** (`backend/internal/platform/types/list/list.go`):
   - Uses camelCase field names in JSON serialization

4. **Other portal API clients** (`portal-web/src/api/order.ts`):
   - Correctly use camelCase for all types

---

## Impact Assessment

**Severity:** High

- **Type Safety Broken:** TypeScript types don't match runtime data
- **Runtime Errors:** Accessing `response.page_size` returns `undefined`; actual data is at `response.pageSize`
- **Silent Failures:** Code may work if it doesn't access these specific fields
- **Consistency Violation:** Creates confusion about Kyora's API standards

**Affected Code:**
- All inventory list operations (products, variants, categories)
- Any code that reads pagination metadata from list responses
- Variant creation requests

---

## Root Cause

Portal inventory API types were defined without referencing backend response DTOs or OpenAPI schema. This is a code drift where portal types don't follow the established backend contract.

---

## Suggested Fix

**Update types in `portal-web/src/api/inventory.ts`:**

```typescript
export interface ListResponse<T> {
  items: Array<T>
  page: number
  pageSize: number        // Changed from page_size
  totalCount: number      // Changed from total_count
  totalPages: number      // Changed from total_pages
  hasMore: boolean        // Changed from has_more
}

export interface CreateVariantRequest {
  productId: string       // Changed from product_id
  code: string
  sku?: string
  photos?: Array<AssetReference>
  costPrice: string
  salePrice: string
  stockQuantity: number
  stockQuantityAlert: number
}
```

**Search for usages** and update any code that accesses the old snake_case properties.

**Verification:**
1. Update all type definitions in `portal-web/src/api/inventory.ts`
2. Search for usages of `page_size`, `total_count`, `total_pages`, `has_more`, `product_id`
3. Update any code that references these properties
4. Test inventory list operations in the UI
5. Verify pagination controls work correctly

---

## Related Issues

- This is part of a broader pattern where some portal API types don't match backend contracts
- See `.github/instructions/responses-dtos-swagger.instructions.md` for API contract standards
- Similar issues may exist in other domains (should audit systematically)

---

## References

- **SSOT:** `.github/instructions/inventory.instructions.md` (updated to document this drift)
- **Backend Types:** `backend/internal/domain/inventory/model_response.go`
- **Backend E2E Tests:** `backend/internal/tests/e2e/inventory_products_test.go`
- **List Response:** `backend/internal/platform/types/list/list.go`
- **Correct Pattern:** `portal-web/src/api/order.ts` (uses camelCase)
