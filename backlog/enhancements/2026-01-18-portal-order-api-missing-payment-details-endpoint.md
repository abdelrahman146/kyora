---
title: "Portal Order API Missing Add-Payment-Details Endpoint"
date: 2026-01-18
category: enhancement
area: portal-web/api
impact: low
status: implemented
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

## Implementation

**Status:** ✅ Implemented
**Date:** 2026-01-19
**Implementation Summary:**
Added `addOrderPaymentDetails` API method and `useAddOrderPaymentDetailsMutation` React Query hook to portal-web order API client, following existing patterns for order mutations.

**Files Modified:**
- `portal-web/src/api/order.ts` - Added API method and mutation hook

**Implementation Details:**

1. **API Method Added** (line ~385):
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

2. **React Query Hook Added** (line ~577):
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

3. **Request Type** (already existed at line ~234):
```typescript
export interface AddOrderPaymentDetailsRequest {
  paymentMethod: OrderPaymentMethod
  paymentReference?: string
}
```

**Usage Example:**

```typescript
import { useAddOrderPaymentDetailsMutation } from '@/api/order'

function OrderPaymentDetailsForm({ businessDescriptor, orderId }) {
  const mutation = useAddOrderPaymentDetailsMutation(businessDescriptor, orderId, {
    onSuccess: (updatedOrder) => {
      showSuccessToast('Payment details updated')
      queryClient.invalidateQueries(orderQueries.detail(businessDescriptor, orderId))
    },
  })

  const handleSubmit = (data: AddOrderPaymentDetailsRequest) => {
    mutation.mutate(data)
  }

  return (
    // Form implementation
  )
}
```

**Success Criteria Met:**
- [x] API method added following existing patterns
- [x] React Query mutation hook added
- [x] TypeScript types properly defined
- [x] Follows http-tanstack-query.instructions.md patterns
- [x] Consistent with other order mutation hooks
- [x] Type checking passes
- [x] Linting passes with no issues

**Validation:**
- [x] Type check passes (100%)
- [x] Lint passes with no issues
- [x] Follows portal-web-architecture.instructions.md
- [x] Follows http-tanstack-query.instructions.md
- [x] DRY and reusable implementation
- [x] Well-documented with JSDoc comments
- [x] Consistent naming with existing hooks
- [x] Uses typed helpers (patch)
- [x] Proper return types

**Pattern Consistency:**
- ✅ Uses `patch<Order>` helper from client.ts
- ✅ Returns `Promise<Order>` matching backend response
- ✅ Follows business-scoped URL pattern
- ✅ Mutation hook structure matches other order mutations
- ✅ Options parameter allows customization
- ✅ TypeScript generics properly configured

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
- `portal-web/src/features/orders/components/OrderQuickActions.tsx`
- `portal-web/src/i18n/en/orders.json`
- `portal-web/src/i18n/ar/orders.json`

## UI Implementation

**Status:** ✅ Implemented
**Date:** 2026-01-19

The payment details functionality has been integrated into OrderQuickActions component with smart conditional visibility and dynamic labeling.

**Features Implemented:**

1. **Conditional Menu Item:**
   - Only visible when order is not in final states (cancelled, returned)
   - Only visible when payment is not in final states (paid, refunded)
   - Dynamic label: "Add Payment Details" vs "Update Payment Details"

2. **Smart Labeling:**
   - Shows "Add" when order has no payment method
   - Shows "Update" when payment method already exists
   - Applied to both menu item and sheet title

3. **Bottom Sheet Form:**
   - Payment method select field (all 6 payment methods)
   - Payment reference text field (optional)
   - Shows current payment details before editing
   - Proper loading states during submission
   - Success/error handling via global error handler

4. **i18n Support:**
   - Added English translations for all new UI labels
   - Added Arabic translations for all new UI labels
   - Reused existing payment method translations

**Files Modified for UI:**
- `portal-web/src/features/orders/components/OrderQuickActions.tsx`
  - Added state for payment details sheet
  - Added mutation hook integration
  - Added payment details form with validation
  - Added conditional menu item with dynamic label
  - Added bottom sheet with current details display
- `portal-web/src/i18n/en/orders.json` - Added 5 new keys
- `portal-web/src/i18n/ar/orders.json` - Added 5 new keys

**UX Flow:**
1. User opens order quick actions
2. "Add/Update Payment Details" appears (if eligible)
3. User clicks to open bottom sheet
4. Sheet shows current payment details (if any)
5. User selects/updates payment method and reference
6. User submits → API call → success toast → sheet closes → data refreshes

**Conditional Visibility Logic:**
```typescript
{!['cancelled', 'returned'].includes(order.status) &&
  !['paid', 'refunded'].includes(order.paymentStatus) && (
    <li>
      <button onClick={() => setShowPaymentDetailsSheet(true)}>
        <CreditCard size={18} />
        {order.paymentMethod
          ? tOrders('update_payment_details')
          : tOrders('add_payment_details')}
      </button>
    </li>
  )}
```

**Validation:**
- [x] Type check passes (100%)
- [x] Lint passes with no issues
- [x] Follows OrderQuickActions existing patterns
- [x] Mobile-first responsive layout
- [x] RTL-compatible
- [x] Proper loading/error states
- [x] i18n complete (en + ar)
- [x] Conditional visibility implemented
- [x] Dynamic labeling based on state
- [x] Follows forms.instructions.md patterns
- [x] Uses global error handler

**Component Integration:**
- Integrated seamlessly into existing OrderQuickActions
- Follows same pattern as status/payment/address updates
- Reuses existing form system and UI components
- Consistent with other quick action items

## References

- Orders SSOT: `.github/instructions/orders.instructions.md`
- Backend implementation: `backend/internal/domain/order/service.go:734-775`
