---
title: "Portal Order API Missing Get-By-Number Endpoint"
date: 2026-01-18
category: enhancement
area: portal-web/api
impact: low
status: open
---

# Portal Order API Missing Get-By-Number Endpoint

## Summary

The portal-web order API client does not expose `GET /orders/by-number/:orderNumber`, even though the backend supports this endpoint.

## Current State

**Backend:**
- Endpoint exists: `GET /v1/businesses/:businessDescriptor/orders/by-number/:orderNumber`
- Route: `backend/internal/server/routes.go:351`
- Handler: `backend/internal/domain/order/handler_http.go:GetOrderByNumber`
- Permission: `role.ActionView` on `role.ResourceOrder`
- Returns: `OrderResponse` (same as `GetOrder`)

**Portal:**
- Missing from `portal-web/src/api/order.ts`
- No corresponding query options factory
- No React Query hook

## Expected State

Portal should expose this endpoint for cases where:
1. URL routing uses order number instead of ID (e.g., shareable links)
2. Quick order lookup by user-facing order number
3. Order search result navigation

## Suggested Implementation

### 1. Add API method in `portal-web/src/api/order.ts`:

```typescript
/**
 * Get order by order number (unique per business)
 */
async getOrderByNumber(
  businessDescriptor: string,
  orderNumber: string,
): Promise<Order> {
  return get<Order>(
    `v1/businesses/${businessDescriptor}/orders/by-number/${orderNumber}`
  )
}
```

### 2. Add query options factory:

```typescript
byNumber: (businessDescriptor: string, orderNumber: string) =>
  queryOptions({
    queryKey: [
      ...orderQueries.details(),
      businessDescriptor,
      'by-number',
      orderNumber,
    ] as const,
    queryFn: () =>
      orderApi.getOrderByNumber(businessDescriptor, orderNumber),
    staleTime: STALE_TIME.THIRTY_SECONDS,
  }),
```

### 3. Add React Query hook:

```typescript
export function useOrderByNumberQuery(
  businessDescriptor: string,
  orderNumber: string,
) {
  return useQuery(orderQueries.byNumber(businessDescriptor, orderNumber))
}
```

## Impact Assessment

- **User Experience:** No current impact; feature is not used
- **Consistency:** Low - creates gap between backend and portal API surface
- **Developer Experience:** Missing convenience method for order number lookups

## Related Files

- `backend/internal/server/routes.go`
- `backend/internal/domain/order/handler_http.go`
- `portal-web/src/api/order.ts`

## References

- Orders SSOT: `.github/instructions/orders.instructions.md`
