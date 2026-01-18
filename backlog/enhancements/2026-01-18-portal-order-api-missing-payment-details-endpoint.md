---
title: "Portal Order API Missing Add-Payment-Details Endpoint"
date: 2026-01-18
category: enhancement
area: portal-web/api
impact: low
status: open
---

# Portal Order API Missing Add-Payment-Details Endpoint

## Summary

The portal-web order API client does not expose `PATCH /orders/:orderId/payment-details`, even though the backend supports this endpoint for updating payment method/reference without changing payment status.

## Current State

**Backend:**
- Endpoint exists: `PATCH /v1/businesses/:businessDescriptor/orders/:orderId/payment-details`
- Route: `backend/internal/server/routes.go:370`
- Handler: `backend/internal/domain/order/handler_http.go:AddOrderPaymentDetails`
- Permission: `role.ActionManage` on `role.ResourceOrder` + plan gates
- Request: `{ paymentMethod, paymentReference? }`
- Returns: `OrderResponse`
- Behavior: Updates payment method/reference without changing payment status; validates payment method is enabled for business

**Portal:**
- Missing from `portal-web/src/api/order.ts`
- No corresponding mutation hook
- Portal only has `updateOrderPaymentStatus` which updates the payment status itself

## Expected State

Portal should expose this endpoint to allow operators to:
1. Update payment method after order creation (e.g., customer switched from COD to bank transfer)
2. Add payment reference for tracking (e.g., transaction ID after payment received)
3. Correct payment method selection without affecting payment status workflow

## Suggested Implementation

### 1. Add request type in `portal-web/src/api/order.ts`:

```typescript
export interface AddOrderPaymentDetailsRequest {
  paymentMethod: OrderPaymentMethod
  paymentReference?: string | null
}
```

### 2. Add API method:

```typescript
/**
 * Add or update payment details (method/reference) without changing status
 */
async addOrderPaymentDetails(
  businessDescriptor: string,
  orderId: string,
  data: AddOrderPaymentDetailsRequest,
): Promise<Order> {
  return patch<Order>(
    `v1/businesses/${businessDescriptor}/orders/${orderId}/payment-details`,
    {
      json: data,
    },
  )
}
```

### 3. Add React Query mutation hook:

```typescript
export function useAddOrderPaymentDetailsMutation(
  businessDescriptor: string,
  orderId: string,
  options?: UseMutationOptions<Order, Error, AddOrderPaymentDetailsRequest>,
) {
  return useMutation({
    mutationFn: (data: AddOrderPaymentDetailsRequest) =>
      orderApi.addOrderPaymentDetails(businessDescriptor, orderId, data),
    ...options,
  })
}
```

## Use Cases

1. **Payment method correction:**
   - Order created with COD, customer wants to pay via bank transfer
   - Update payment method without triggering status change

2. **Reference tracking:**
   - Customer paid via bank transfer
   - Add transaction reference for reconciliation
   - Payment status updated separately via `updateOrderPaymentStatus`

3. **Payment provider changes:**
   - Customer wants to switch from Tamara to Tabby installments
   - Update method before payment status transitions

## Impact Assessment

- **User Experience:** No current impact; feature workflow is not exposed in UI
- **Functionality:** Low - operators currently must update payment status to track payment details
- **Data Quality:** Medium - missing dedicated endpoint for payment details tracking

## Implementation Notes

- Backend validates payment method is enabled for the business
- Backend prevents updates for finalized order statuses (cancelled, returned)
- This is intentionally separate from `UpdateOrderPaymentStatus` to avoid accidental status changes

## Related Files

- `backend/internal/server/routes.go`
- `backend/internal/domain/order/handler_http.go`
- `backend/internal/domain/order/service.go` (AddOrderPaymentDetails)
- `portal-web/src/api/order.ts`

## References

- Orders SSOT: `.github/instructions/orders.instructions.md`
- Backend implementation: `backend/internal/domain/order/service.go:734-775`
