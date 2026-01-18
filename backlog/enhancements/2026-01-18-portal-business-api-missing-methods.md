---
title: "Portal-Web Missing Business Management API Client Methods"
date: 2026-01-18
priority: low
category: completeness
status: open
domain: portal-web
---

# Portal-Web Missing Business Management API Client Methods

## Summary

Portal-web's business API client (`portal-web/src/api/business.ts`) does not expose several backend endpoints that are already implemented and documented in business-management.instructions.md. This limits the portal's ability to implement full business management features.

## Current State

The portal-web business API client currently exposes:
- ✅ List businesses
- ✅ Get business by descriptor
- ✅ Create business
- ✅ Update business
- ✅ Delete business
- ✅ List shipping zones
- ✅ List payment methods

### Missing API Methods (Backend Exists, Portal Missing)

1. **Descriptor Availability Check**
   - Backend: `GET /v1/businesses/descriptor/availability?descriptor=...`
   - Handler: `backend/internal/domain/business/handler_http.go` (CheckDescriptorAvailability)
   - Use case: Validate descriptor before creation; real-time feedback in forms

2. **Archive Business**
   - Backend: `POST /v1/businesses/:businessDescriptor/archive`
   - Handler: `backend/internal/domain/business/handler_http.go` (ArchiveBusiness)
   - Use case: Soft-delete businesses without permanent removal

3. **Unarchive Business**
   - Backend: `POST /v1/businesses/:businessDescriptor/unarchive`
   - Handler: `backend/internal/domain/business/handler_http.go` (UnarchiveBusiness)
   - Use case: Restore archived businesses

4. **Get Shipping Zone by ID**
   - Backend: `GET /v1/businesses/:businessDescriptor/shipping-zones/:zoneId`
   - Handler: `backend/internal/domain/business/handler_http.go` (GetShippingZone)
   - Use case: Load individual zone details for editing

5. **Create Shipping Zone**
   - Backend: `POST /v1/businesses/:businessDescriptor/shipping-zones`
   - Handler: `backend/internal/domain/business/handler_http.go` (CreateShippingZone)
   - Plan gates: Requires active subscription
   - Use case: Add new shipping zones

6. **Update Shipping Zone**
   - Backend: `PATCH /v1/businesses/:businessDescriptor/shipping-zones/:zoneId`
   - Handler: `backend/internal/domain/business/handler_http.go` (UpdateShippingZone)
   - Plan gates: Requires active subscription
   - Use case: Modify existing shipping zones

7. **Delete Shipping Zone**
   - Backend: `DELETE /v1/businesses/:businessDescriptor/shipping-zones/:zoneId`
   - Handler: `backend/internal/domain/business/handler_http.go` (DeleteShippingZone)
   - Plan gates: Requires active subscription
   - Use case: Remove shipping zones

8. **Update Payment Method Override**
   - Backend: `PATCH /v1/businesses/:businessDescriptor/payment-methods/:descriptor`
   - Handler: `backend/internal/domain/business/handler_http.go` (UpdatePaymentMethod)
   - Plan gates: Requires active subscription
   - Use case: Enable/disable payment methods and set per-business fee overrides

## Expected State

Per business-management.instructions.md:

> ### Frontend features not implemented yet (but backend supports)
> 
> If/when building business settings UI, align to backend routes above:
> 
> - Descriptor availability check (`GET /v1/businesses/descriptor/availability`)
> - Archive/unarchive (`POST /v1/businesses/:descriptor/archive|unarchive`)
> - Shipping zone CRUD (create/update/delete are plan-gated)
> - Payment methods list + per-business override update (plan-gated)

## Impact

- **Low**: Features work correctly with current implementation; these are additional capabilities
- Portal cannot implement business archive/unarchive UI
- Portal cannot implement shipping zone management UI
- Portal cannot implement payment method configuration UI
- Portal cannot provide real-time descriptor validation during business creation

## Affected Files

- `portal-web/src/api/business.ts` (missing methods)
- Future UI components that will need these methods

## Suggested Fix

Add the missing API client methods to `portal-web/src/api/business.ts`:

```typescript
// Add to businessApi object

/**
 * Check if business descriptor is available
 */
async checkDescriptorAvailability(descriptor: string): Promise<boolean> {
  const response = await get<{ available: boolean }>(
    `v1/businesses/descriptor/availability?descriptor=${encodeURIComponent(descriptor)}`
  )
  return response.available
},

/**
 * Archive a business
 */
async archiveBusiness(descriptor: string): Promise<void> {
  await post<void>(`v1/businesses/${descriptor}/archive`)
},

/**
 * Unarchive a business
 */
async unarchiveBusiness(descriptor: string): Promise<void> {
  await post<void>(`v1/businesses/${descriptor}/unarchive`)
},

/**
 * Get shipping zone by ID
 */
async getShippingZone(
  businessDescriptor: string,
  zoneId: string,
): Promise<ShippingZone> {
  const response = await get<{ zone: ShippingZone }>(
    `v1/businesses/${businessDescriptor}/shipping-zones/${zoneId}`,
  )
  return response.zone
},

/**
 * Create shipping zone
 */
async createShippingZone(
  businessDescriptor: string,
  data: CreateShippingZoneRequest,
): Promise<ShippingZone> {
  const response = await post<{ zone: ShippingZone }>(
    `v1/businesses/${businessDescriptor}/shipping-zones`,
    { json: data },
  )
  return response.zone
},

/**
 * Update shipping zone
 */
async updateShippingZone(
  businessDescriptor: string,
  zoneId: string,
  data: UpdateShippingZoneRequest,
): Promise<ShippingZone> {
  const response = await patch<{ zone: ShippingZone }>(
    `v1/businesses/${businessDescriptor}/shipping-zones/${zoneId}`,
    { json: data },
  )
  return response.zone
},

/**
 * Delete shipping zone
 */
async deleteShippingZone(
  businessDescriptor: string,
  zoneId: string,
): Promise<void> {
  await del<void>(
    `v1/businesses/${businessDescriptor}/shipping-zones/${zoneId}`,
  )
},

/**
 * Update payment method override
 */
async updatePaymentMethod(
  businessDescriptor: string,
  methodDescriptor: string,
  data: UpdatePaymentMethodRequest,
): Promise<PaymentMethod> {
  const response = await patch<PaymentMethod>(
    `v1/businesses/${businessDescriptor}/payment-methods/${methodDescriptor}`,
    { json: data },
  )
  return response
},
```

Also add corresponding:
- Type definitions for request/response types
- Query options factories in `businessQueries`
- Mutation hooks (e.g., `useArchiveBusinessMutation`, `useCreateShippingZoneMutation`, etc.)

### Verification

After implementation:
1. All backend business management endpoints should have corresponding portal API client methods
2. Methods should follow existing patterns in `portal-web/src/api/business.ts`
3. Add query options to `businessQueries` for new GET endpoints
4. Add mutation hooks for new POST/PATCH/DELETE endpoints

## References

- business-management.instructions.md (Backend route surface, Frontend features not implemented section)
- Backend handlers: `backend/internal/domain/business/handler_http.go`
- Route registration: `backend/internal/server/routes.go` (lines 208-217, 315-347)
- Portal API client: `portal-web/src/api/business.ts`
