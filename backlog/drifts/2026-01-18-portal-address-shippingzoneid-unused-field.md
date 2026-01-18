---
title: "Portal-web Address Forms Include shippingZoneId Field Not Supported by Backend"
date: 2026-01-18
priority: medium
category: consistency
status: open
domain: portal-web
---

# Portal-web Address Forms Include shippingZoneId Field Not Supported by Backend

## Summary

Portal-web customer address forms (`AddressSheet` component) include a `shippingZoneId` field with validation, but backend `CustomerAddress` model does not have this field. The field is used for UI pre-selection but is never sent to the backend, creating confusion and dead code.

## Current State

### Portal-web Form Schema (`portal-web/src/features/customers/components/AddressSheet.tsx`)

```typescript
const addressSchema = z.object({
  shippingZoneId: z.string().min(1, 'validation.shipping_zone_required'), // ← Not in backend
  countryCode: z.string().length(2, 'validation.country_required'),
  state: z.string().min(1, 'validation.state_required'),
  city: z.string().min(1, 'validation.city_required'),
  phoneCode: z.string().min(1, 'validation.phone_code_required'),
  phoneNumber: z.string().min(1, 'validation.phone_required'),
  street: z.string().optional(),
  zipCode: z.string().optional(),
})
```

### Portal-web API Type (`portal-web/src/api/customer.ts`)

```typescript
export interface CustomerAddress {
  id: string
  customerId: string
  shippingZoneId?: string  // ← Not in backend model
  countryCode: string
  state: string
  city: string
  // ...
}
```

### Backend Model (`backend/internal/domain/customer/model.go`)

```go
type CustomerAddress struct {
    gorm.Model
    ID          string          `gorm:"column:id;primaryKey;type:text" json:"id"`
    CustomerID  string          `gorm:"column:customer_id;type:text;not null;index" json:"customerId"`
    CountryCode string          `gorm:"column:country_code;type:text; not null" json:"countryCode"`
    State       string          `gorm:"column:state;type:text; not null" json:"state"`
    City        string          `gorm:"column:city;type:text; not null" json:"city"`
    Street      nullable.String `gorm:"column:street;type:text" json:"street"`
    PhoneCode   string          `gorm:"column:phone_code;type:text; not null" json:"phoneCode"`
    PhoneNumber string          `gorm:"column:phone_number;type:text; not null" json:"phoneNumber"`
    ZipCode     nullable.String `gorm:"column:zip_code;type:text" json:"zipCode"`
    // No shippingZoneId field
}
```

### Form Implementation Reality

Portal-web form collects `shippingZoneId` but strips it before API submission:

```typescript
// Create mode
const createData: CreateAddressRequest = {
  countryCode: value.countryCode,
  state: value.state,
  city: value.city,
  phoneCode: value.phoneCode,
  phoneNumber: value.phoneNumber,
  street: value.street,
  zipCode: value.zipCode,
  // shippingZoneId is NOT included
}
```

## Expected State

**Option 1 (Quick Fix): Remove unused field from portal-web**

Remove `shippingZoneId` from:
1. `addressSchema` in `AddressSheet.tsx`
2. `CustomerAddress` interface in `portal-web/src/api/customer.ts`
3. All form state and validation logic
4. `ShippingZoneSelect` component usage in address forms

**Option 2 (Feature Enhancement): Add shipping zone support to backend**

If shipping zones per address are a planned feature:
1. Add `ShippingZoneID` field to backend `CustomerAddress` model
2. Add foreign key relationship to `ShippingZone`
3. Update request/response DTOs
4. Keep portal-web implementation as-is

## Impact

- **Medium**: Portal-web has unnecessary validation and UI for a field that doesn't exist in backend.
- Adds ~100 lines of dead code in address forms and components.
- Creates confusion for developers about whether addresses are actually linked to shipping zones.
- UI shows shipping zone selector but selection has no effect.

## Affected Files

- `portal-web/src/features/customers/components/AddressSheet.tsx` (lines 44, 77, 102, 159, 168, 176, 268, 270)
- `portal-web/src/features/customers/components/StandaloneAddressSheet.tsx` (line 66)
- `portal-web/src/features/customers/components/ShippingZoneSelect.tsx` (used in address forms)
- `portal-web/src/api/customer.ts` (line 36: `CustomerAddress` interface)
- `portal-web/src/api/address.ts` (`CreateAddressRequest` and `UpdateAddressRequest` do NOT include it, correctly)

## Suggested Fix

**Recommended: Remove from portal-web (Option 1)**

1. Remove `shippingZoneId` from address form schema
2. Remove `ShippingZoneSelect` component usage from address forms
3. Update `CustomerAddress` interface to remove `shippingZoneId?` field
4. Clean up any related state management code

**Justification:** Backend does not support this field, and there's no indication in backlog or BRDs that it's a planned feature. Shipping zones are business-level configuration, not per-address.

## Related

- `.github/instructions/customer.instructions.md` (documents this as known drift)
- `backend/internal/domain/business/model.go` (ShippingZone is business-level, not address-level)
