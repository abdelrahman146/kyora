---
type: drift
date: 2026-01-18
priority: low
component: portal-web
affected-files:
  - portal-web/src/features/onboarding/components/BusinessSetupPage.tsx
  - backend/internal/domain/onboarding/service.go
related-instructions:
  - .github/instructions/onboarding.instructions.md
status: open
assignee: null
pattern-category: api-contract-mismatch
---

# Portal business step checks for `business_staged` but backend sets `payment_pending`

## Summary

The portal-web business setup page navigation logic checks for `business_staged` stage after calling the backend, but the backend never returns this stage for paid plans. Instead, the backend immediately sets `payment_pending` for paid plans after business details are submitted.

## Current State

**Backend behavior (backend/internal/domain/onboarding/service.go:302-307):**
```go
sess.Stage = StageBusinessStaged
if !sess.IsPaidPlan {
    sess.Stage = StageReadyToCommit
    sess.PaymentStatus = PaymentStatusSkipped
} else {
    sess.Stage = StagePaymentPending
    sess.PaymentStatus = PaymentStatusPending
}
```

The backend sets `StageBusinessStaged` initially but immediately overwrites it with either `StageReadyToCommit` (free) or `StagePaymentPending` (paid).

**Portal behavior (portal-web/src/features/onboarding/components/BusinessSetupPage.tsx:79):**
```tsx
} else if (response.stage === 'business_staged') {
    await navigate({
        to: '/onboarding/payment',
        search: { session: sessionToken },
    })
}
```

The portal checks for `business_staged` but this branch will never execute for paid plans since the backend returns `payment_pending` instead.

## Expected State

The navigation logic should check for `payment_pending` instead of `business_staged` to match the backend's actual behavior:

```tsx
} else if (response.stage === 'payment_pending') {
    await navigate({
        to: '/onboarding/payment',
        search: { session: sessionToken },
    })
}
```

## Impact

- **Functional:** The current code still works because the final `else` branch catches all other cases and navigates to complete. However, paid plan users should go to payment first, not complete.
- **Maintenance:** This mismatch makes the code harder to understand and could cause bugs if the else branch is removed or modified.
- **User Experience:** Potential for incorrect routing if the fallback else branch behavior changes.

## Root Cause

The backend initially sets `business_staged` but immediately overwrites it for both free and paid plans. The portal was written expecting the intermediate `business_staged` value to be returned.

## Verification Steps

1. Start onboarding with a paid plan
2. Complete identity verification
3. Submit business details
4. Observe that backend returns `stage: "payment_pending"` not `stage: "business_staged"`
5. Observe that portal navigation goes through the final else branch, not the business_staged branch

## Suggested Fix

**Option 1 (Recommended): Update portal to match backend reality**
```tsx
if (response.stage === 'ready_to_commit') {
    await navigate({
        to: '/onboarding/complete',
        search: { session: sessionToken },
    })
} else if (response.stage === 'payment_pending') {
    await navigate({
        to: '/onboarding/payment',
        search: { session: sessionToken },
    })
} else {
    // Fallback for unexpected states
    await navigate({
        to: '/onboarding/complete',
        search: { session: sessionToken },
    })
}
```

**Option 2: Change backend to return business_staged**

This would require keeping the stage as `business_staged` for paid plans and only transitioning to `payment_pending` after the payment flow starts. However, this is less semantically correct since the payment is pending as soon as business details are submitted.

## Related Issues

- None

## Notes

- The `STAGE_ROUTES` mapping in `portal-web/src/features/onboarding/utils/onboarding.ts` correctly maps both `business_staged` and `payment_pending` to `/onboarding/payment`, so automatic stage redirects work correctly.
- The issue is only in the explicit navigation decision logic in BusinessSetupPage.tsx.
