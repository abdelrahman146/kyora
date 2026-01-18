---
title: "Customer Email Field Optional in Backend but Portal Treats as Required"
date: 2026-01-18
priority: low
category: consistency
status: open
domain: backend, portal-web
---

# Customer Email Field Optional in Backend but Portal Treats as Required

## Summary

Backend `CreateCustomerRequest.Email` is optional (`binding:"omitempty,email"`), but portal-web's `CreateCustomerRequest` interface types `email` as optional without documenting that phone is the primary identifier. This creates ambiguity about which fields are truly required for customer creation.

## Current State

### Backend (`backend/internal/domain/customer/model.go`)

```go
type CreateCustomerRequest struct {
    Name              string         `json:"name" binding:"required"`
    Email             string         `json:"email" binding:"omitempty,email"`
    PhoneNumber       string         `json:"phoneNumber" binding:"required"`
    PhoneCode         string         `json:"phoneCode" binding:"required"`
    // ...
}
```

Email is **optional**, but PhoneNumber + PhoneCode are **required**.

### Portal-web (`portal-web/src/api/customer.ts`)

```typescript
export interface CreateCustomerRequest {
  name: string
  email?: string
  phoneCode: string
  phoneNumber: string
  countryCode: string
  // ...
}
```

Portal correctly types `email` as optional, but forms/UI may not clearly communicate that phone is the primary identifier.

## Expected State

Backend behavior is correct (phone-first for social commerce). Portal-web should:

1. Ensure customer creation forms clearly indicate that phone is required and email is optional
2. Consider adding helper text: "Email is optional. Phone number is the primary identifier for social commerce."

## Impact

- **Low**: The typing is technically correct on both sides.
- Users may be confused about which fields are truly required.
- Social commerce businesses often don't collect emails upfront (DM orders start with WhatsApp/Instagram handles).

## Affected Files

- `backend/internal/domain/customer/model.go` (lines 102-118)
- `portal-web/src/api/customer.ts` (lines 90-106)
- `portal-web/src/features/customers/components/AddCustomerSheet.tsx` (form UI)

## Suggested Fix

**Portal-web UI enhancement:**

In customer create forms, add helper text to the email field:

```tsx
<TextInputField
  name="email"
  label={t('email')}
  placeholder={t('email_placeholder')}
  helperText={t('email_optional_helper')} // "Email is optional. Use phone for identification."
/>
```

**i18n key:**
- Add `customers.email_optional_helper` translation

## Related

- `.github/instructions/customer.instructions.md` (documents this as expected behavior)
