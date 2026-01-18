---
type: drift
date: 2026-01-18
priority: medium
component: portal-web
affected-files:
  - portal-web/src/features/onboarding/components/PaymentPage.tsx
  - portal-web/src/features/onboarding/components/BusinessSetupPage.tsx
  - portal-web/src/features/onboarding/components/CompleteOnboardingPage.tsx
  - portal-web/src/features/onboarding/components/VerifyEmailPage.tsx
  - portal-web/src/features/onboarding/components/EmailEntryPage.tsx
  - portal-web/src/features/onboarding/components/OAuthCallbackPage.tsx
related-instructions:
  - .github/instructions/errors-handling.instructions.md
  - .github/instructions/http-tanstack-query.instructions.md
status: open
assignee: null
pattern-category: error-handling
---

# Onboarding components bypass global React Query error handler with manual mutation.error displays

## Summary

Several onboarding components manually display `mutation.error.message` in the UI, bypassing the global React Query error handling system configured in `portal-web/src/main.tsx`. This creates inconsistent UX where onboarding pages show inline error alerts while other features rely on toast notifications from the global handler.

## Current State

**What exists now:**
- Onboarding pages manually render mutation errors using inline alert components
- Direct access to `mutation.error` and `mutation.error.message` in JSX
- Mixed patterns: some mutations show errors inline, others rely on global toast handler

**Affected locations:**
```
portal-web/src/features/onboarding/components/PaymentPage.tsx:130
portal-web/src/features/onboarding/components/PaymentPage.tsx:198
portal-web/src/features/onboarding/components/BusinessSetupPage.tsx:252-255
portal-web/src/features/onboarding/components/CompleteOnboardingPage.tsx:65
portal-web/src/features/onboarding/components/VerifyEmailPage.tsx:168,223,385,389
portal-web/src/features/onboarding/components/EmailEntryPage.tsx:60
portal-web/src/features/onboarding/components/OAuthCallbackPage.tsx:96
```

**Example from BusinessSetupPage.tsx:**
```tsx
{setBusinessMutation.error && (
  <div className="alert alert-error">
    <span className="text-sm">
      {setBusinessMutation.error.message}
    </span>
  </div>
)}
```

**Example from PaymentPage.tsx:**
```tsx
{startPaymentMutation.error && (
  <div className="alert alert-error mb-4">
    <AlertCircle className="h-5 w-5" />
    <span>{startPaymentMutation.error.message}</span>
  </div>
)}
```

## Expected Pattern

**Per instruction file** (`.github/instructions/errors-handling.instructions.md` + `.github/instructions/http-tanstack-query.instructions.md`):

Portal-web has a global error handling system configured via `QueryClient` with `MutationCache.onError` and `QueryCache.onError` callbacks that:
1. Automatically translate backend errors via `translateErrorAsync`
2. Show user-friendly toast notifications
3. Support opt-out via `meta: { errorToast: 'off' }` for special cases

**Standard implementation** (from `portal-web/src/main.tsx`):
```typescript
const queryClient = new QueryClient({
  mutationCache: new MutationCache({
    onError: (error, _variables, _context, mutation) => {
      if (shouldIgnoreGlobalError(error)) return
      const meta = mutation.meta as undefined | { errorToast?: 'global' | 'off' }
      if (meta?.errorToast === 'off') return
      void showErrorFromException(error, i18n.t) // Translates + shows toast
    },
  }),
  // ...
})
```

**Current portal-web reality:**
- ✅ All other domains (customer, order, inventory, accounting, business, user) rely on global handler
- ❌ Only onboarding components bypass this pattern with manual `mutation.error.message` displays
- ❌ Manual displays show **untranslated** error messages (raw backend `error.message` instead of `translateErrorAsync`)

Components should:
- **Default:** Rely on global error handler (no manual error display) - followed by all domains except onboarding
- **Special case:** Use `meta: { errorToast: 'off' }` if custom error handling is needed (e.g., inline form validation)

## Pattern Deviation Analysis

**Type of drift:**
- [x] Error handling convention violation
- [ ] Naming convention violation
- [ ] File structure deviation
- [ ] API contract inconsistency
- [ ] State management pattern deviation

**Why this is drift:**
1. **Inconsistent UX:** Onboarding shows inline alerts, rest of app shows toasts
2. **Bypasses i18n:** Direct `.error.message` doesn't go through `translateErrorAsync`
3. **Duplicate logic:** Each component reimplements error display instead of reusing global system
4. **Maintenance burden:** Changes to error display format require updating multiple components

**Why it might have happened:**
- Onboarding was implemented before global error handler was established
- UX requirement for inline errors in forms wasn't communicated
- Lack of documentation about the global error handler

## Impact Assessment

**Severity:** Medium

**User Impact:**
- Inconsistent error experience (alerts vs toasts)
- Untranslated error messages in onboarding flow
- Less polished UX in critical user flow (onboarding)

**Developer Impact:**
- Confusion about correct error handling pattern
- Maintenance burden (updating 6 components for error display changes)
- Code duplication

**Maintenance Cost:** Low (straightforward refactor)

## Proposed Fix

### Option 1: Use global handler (recommended for most)

Remove manual error displays and rely on global toast handler:

```tsx
// BEFORE (BusinessSetupPage.tsx)
{setBusinessMutation.error && (
  <div className="alert alert-error">
    <span>{setBusinessMutation.error.message}</span>
  </div>
)}

// AFTER (rely on global handler - no changes needed)
// Global MutationCache.onError will show toast automatically
```

### Option 2: Opt-out with meta flag (if inline errors are required)

If UX specifically requires inline errors in onboarding forms:

```tsx
// In mutation hook definition (e.g., api/onboarding.ts)
export function useSetBusinessMutation(
  options?: UseMutationOptions<SetBusinessResponse, Error, SetBusinessRequest>,
) {
  return useMutation({
    mutationFn: (data: SetBusinessRequest) => onboardingApi.setBusiness(data),
    meta: { errorToast: 'off' }, // Opt out of global handler
    ...options,
  })
}

// In component (BusinessSetupPage.tsx)
{setBusinessMutation.error && (
  <div className="alert alert-error">
    <span>{await translateErrorAsync(setBusinessMutation.error, t)}</span>
  </div>
)}
```

### Implementation Steps

1. **Decision:** Determine if onboarding needs inline errors or can use toasts
2. **If toasts (Option 1):**
   - Remove all `{mutation.error && ...}` JSX blocks from affected files
   - Test onboarding flow to confirm errors appear as toasts
3. **If inline errors (Option 2):**
   - Add `meta: { errorToast: 'off' }` to mutation hook definitions
   - Update inline error displays to use `translateErrorAsync`
   - Consider creating shared `<InlineErrorAlert>` component
4. **Update tests:** Ensure error scenarios are covered

## Related Issues

- None found in existing backlog

## References

- Global error handler implementation: `portal-web/src/main.tsx:91-104`
- Error translation utility: `portal-web/src/lib/translateError.ts`
- Toast utility: `portal-web/src/lib/toast.ts`
- HTTP + TanStack Query SSOT: `.github/instructions/http-tanstack-query.instructions.md`
- Error handling SSOT: `.github/instructions/errors-handling.instructions.md`

## Notes

- This is a pattern alignment issue, not a functionality bug
- Current implementation works but creates inconsistent UX
- Fix should align with product decision on onboarding error UX (toasts vs inline)
