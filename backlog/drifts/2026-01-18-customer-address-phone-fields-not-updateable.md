---
title: "Customer Address Phone Fields Not Updateable Despite DTO Including Them"
date: 2026-01-18
priority: medium
category: consistency
status: resolved
domain: backend
---

# Customer Address Phone Fields Not Updateable Despite DTO Including Them

## Summary

Backend `UpdateCustomerAddressRequest` DTO includes `phoneCode` and `phoneNumber` fields with validation bindings, but the service implementation (`UpdateCustomerAddress`) silently ignores these fields and does not update them. This creates a misleading API contract.

## Current State

### Request DTO (`backend/internal/domain/customer/model.go`)

```go
type UpdateCustomerAddressRequest struct {
    Street      string `json:"street" binding:"omitempty"`
    City        string `json:"city" binding:"omitempty"`
    State       string `json:"state" binding:"omitempty"`
    CountryCode string `json:"countryCode" binding:"omitempty,len=2"`
    PhoneCode   string `json:"phoneCode" binding:"omitempty"`     // ← Present in DTO
    PhoneNumber string `json:"phoneNumber" binding:"omitempty"`   // ← Present in DTO
    ZipCode     string `json:"zipCode" binding:"omitempty"`
}
```

### Service Implementation (`backend/internal/domain/customer/service.go`, lines 534-570)

```go
func (s *Service) UpdateCustomerAddress(..., req *UpdateCustomerAddressRequest) (*CustomerAddress, error) {
    // ...
    if req.Street != "" {
        address.Street = transformer.ToNullableString(req.Street)
    }
    if req.City != "" {
        address.City = req.City
    }
    if req.State != "" {
        address.State = req.State
    }
    if req.ZipCode != "" {
        address.ZipCode = transformer.ToNullableString(req.ZipCode)
    }
    if req.CountryCode != "" {
        address.CountryCode = strings.ToUpper(strings.TrimSpace(req.CountryCode))
    }
    // ← phoneCode and phoneNumber are NOT updated
    err = s.storage.customerAddress.UpdateOne(ctx, address)
    // ...
}
```

## Expected State

Two valid approaches:

**Option 1 (Recommended): Remove phone fields from update DTO**

```go
type UpdateCustomerAddressRequest struct {
    Street      string `json:"street" binding:"omitempty"`
    City        string `json:"city" binding:"omitempty"`
    State       string `json:"state" binding:"omitempty"`
    CountryCode string `json:"countryCode" binding:"omitempty,len=2"`
    ZipCode     string `json:"zipCode" binding:"omitempty"`
    // phoneCode and phoneNumber removed - not updateable
}
```

**Option 2: Implement phone field updates in service**

Add to service implementation:
```go
if req.PhoneCode != "" {
    address.PhoneCode = req.PhoneCode
}
if req.PhoneNumber != "" {
    address.PhoneNumber = req.PhoneNumber
}
```

## Impact

- **Medium**: API contract is misleading - clients may send phone fields expecting them to be updated.
- Portal-web currently includes `phoneCode` and `phoneNumber` in `UpdateAddressRequest` type, but they are silently ignored by backend.
- Swagger/OpenAPI documentation includes these fields as updateable, which is incorrect.

## Affected Files

- `backend/internal/domain/customer/model.go` (line 206-214: `UpdateCustomerAddressRequest`)
- `backend/internal/domain/customer/service.go` (lines 534-570: `UpdateCustomerAddress`)
- `portal-web/src/api/address.ts` (lines 29-38: `UpdateAddressRequest`)

## Suggested Fix

Option 2: Implement phone field updates in service

## Related

- `.github/instructions/customer.instructions.md` (documents current behavior)
- `.github/instructions/responses-dtos-swagger.instructions.md` (DTO contract standards)

## Resolution

**Status:** Resolved

**Date:** 2026-01-19

**Approach Taken:** Option 2 — implement phone field updates in service

**Harmonization Summary:**

- `UpdateCustomerAddress` now applies `phoneCode` and `phoneNumber` updates (trimmed) alongside other address fields.
- Customer instructions updated to document that phone fields are updateable.

**Files Changed:**

- backend/internal/domain/customer/service.go — add phoneCode/phoneNumber updates in address update flow.
- .github/instructions/customer.instructions.md — reflect updateable phone fields in address updates.

**Migration Completeness:**

- Total instances found: 1 service method; all updated.
- Remaining drift: 0.

**Validation:**

- Tests not run (manual verification only).

**Instruction Files Updated:**

- .github/instructions/customer.instructions.md — clarified address update behavior to include phone fields.
