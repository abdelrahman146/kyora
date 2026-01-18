---
type: drift
date: 2026-01-18
priority: low
component: backend
affected-files:
  - backend/internal/domain/onboarding/model.go
  - backend/internal/domain/onboarding/service.go
  - portal-web/src/api/types/onboarding.ts
  - portal-web/src/routes/onboarding/payment.tsx
  - portal-web/src/routes/onboarding/complete.tsx
  - portal-web/src/features/onboarding/utils/onboarding.ts
related-instructions:
  - .github/instructions/onboarding.instructions.md
status: open
assignee: null
pattern-category: dead-code
---

# Onboarding stage `payment_confirmed` is defined but never set

## Summary

The onboarding state machine defines a `payment_confirmed` stage in the enum, and portal-web supports it in routing and guards, but the backend never actually sets this stage. When Stripe payment succeeds, the backend transitions directly from `payment_pending` to `ready_to_commit`, skipping `payment_confirmed` entirely.

## Current State

**Backend enum definition (backend/internal/domain/onboarding/model.go:26):**
```go
StagePaymentConfirmed SessionStage = "payment_confirmed" // stripe checkout finished
```

**Backend payment success handler (backend/internal/domain/onboarding/service.go:380):**
```go
func (s *Service) MarkPaymentSucceeded(ctx context.Context, sessionID, stripeSubID string) error {
    // lookup by checkout session id
    rec, err := s.storage.session.FindOne(ctx, s.storage.session.ScopeEquals(schema.NewField("checkout_session_id", "checkoutSessionId"), sessionID))
    if err != nil || rec == nil {
        return ErrSessionNotFound(err)
    }
    rec.PaymentStatus = PaymentStatusSucceeded
    rec.StripeSubID = stripeSubID
    rec.Stage = StageReadyToCommit  // ← Goes directly to ready_to_commit, not payment_confirmed
    return s.storage.UpdateSession(ctx, rec)
}
```

**Portal routing support:**
- `portal-web/src/features/onboarding/utils/onboarding.ts:10` - Maps `payment_confirmed` to `/onboarding/complete`
- `portal-web/src/routes/onboarding/payment.tsx:42` - Checks for `payment_confirmed` in beforeLoad guard
- `portal-web/src/routes/onboarding/complete.tsx:27` - Allows `payment_confirmed` in beforeLoad guard

**Grep verification:**
```bash
# Search for assignments to payment_confirmed stage
grep -r "Stage = StagePaymentConfirmed" backend/internal/domain/onboarding/
# Result: No matches found
```

## Expected State

Either:
1. **Remove the unused stage** from both backend enum and portal types, or
2. **Use the stage** by transitioning to `payment_confirmed` after Stripe webhook, then require an explicit user action to move to `ready_to_commit`

## Impact

- **Code clarity:** Developers reading the code expect all enum values to be used
- **Maintenance:** Extra complexity maintaining unused code paths
- **Testing:** Portal guards/routing test for a state that never occurs
- **Documentation:** Instructions document behavior that doesn't exist

## Root Cause

The stage was likely planned to represent an intermediate state where payment was confirmed but the user hadn't yet clicked "Complete" in the UI. However, the implementation evolved to go directly to `ready_to_commit` after webhook confirmation.

## Verification Steps

1. Check backend enum: `payment_confirmed` exists ✓
2. Search for backend assignments: No code sets this stage ✓
3. Check portal support: Portal routing/guards check for it ✓
4. Review state transitions: Backend goes `payment_pending` → `ready_to_commit` directly ✓

## Suggested Fix

**Option 1 (Recommended): Remove the unused stage**

Backend:
```diff
// backend/internal/domain/onboarding/model.go
const (
    StagePlanSelected     SessionStage = "plan_selected"
    StageIdentityPending  SessionStage = "identity_pending"
    StageIdentityVerified SessionStage = "identity_verified"
    StageBusinessStaged   SessionStage = "business_staged"
    StagePaymentPending   SessionStage = "payment_pending"
-   StagePaymentConfirmed SessionStage = "payment_confirmed"
    StageReadyToCommit    SessionStage = "ready_to_commit"
    StageCommitted        SessionStage = "committed"
)
```

Portal:
```diff
// portal-web/src/api/types/onboarding.ts
export const SessionStageSchema = z.enum([
  'plan_selected',
  'identity_pending',
  'identity_verified',
  'business_staged',
  'payment_pending',
- 'payment_confirmed',
  'ready_to_commit',
  'committed',
])
```

```diff
// portal-web/src/features/onboarding/utils/onboarding.ts
export const STAGE_ROUTES: Record<SessionStage, string> = {
  plan_selected: '/onboarding/email',
  identity_pending: '/onboarding/verify',
  identity_verified: '/onboarding/business',
  business_staged: '/onboarding/payment',
  payment_pending: '/onboarding/payment',
- payment_confirmed: '/onboarding/complete',
  ready_to_commit: '/onboarding/complete',
  committed: '/onboarding/complete',
}
```

Remove checks from payment.tsx and complete.tsx route guards.

**Option 2: Implement the stage properly**

Add transition to `payment_confirmed` after webhook:
```go
func (s *Service) MarkPaymentSucceeded(ctx context.Context, sessionID, stripeSubID string) error {
    // ...
    rec.PaymentStatus = PaymentStatusSucceeded
    rec.StripeSubID = stripeSubID
    rec.Stage = StagePaymentConfirmed  // ← Use the intermediate stage
    return s.storage.UpdateSession(ctx, rec)
}
```

Then require explicit user action (Complete button click) to move to `ready_to_commit`.

## Recommendation

**Option 1 is strongly recommended** because:
- Simpler mental model (fewer states)
- Current UX already works smoothly without the intermediate stage
- Webhook → ready_to_commit is semantic (payment confirmed = ready to commit)
- Reduces test surface area

## Related Issues

- None

## Notes

- The payment route's polling mechanism (`refetchInterval: 3000`) is specifically designed to observe the transition from `payment_pending` to `ready_to_commit` after the webhook fires.
- Removing `payment_confirmed` won't break any existing functionality since the backend never sets it.
