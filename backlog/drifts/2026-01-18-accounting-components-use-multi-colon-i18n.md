---
title: "Accounting components use multi-colon i18n pattern (t('common:actions.delete'))"
date: 2026-01-18
status: resolved
resolved-date: 2026-01-18
impact: medium
area: portal-web
tags: [i18n, code-quality, consistency]
---

# Drift: Accounting Components Use Multi-Colon i18n Pattern

## Summary

Accounting quick action components use the forbidden multi-colon namespace pattern `t('common:actions.delete')` instead of the canonical namespace-bound pattern.

## Current State

**Affected Files:**
- `portal-web/src/features/accounting/components/TransactionQuickActions.tsx` (lines 127-128)
- `portal-web/src/features/accounting/components/AssetQuickActions.tsx` (lines 116-117)

**Current Pattern:**
```tsx
confirmText={t('common:actions.delete')}
cancelText={t('common:actions.cancel')}
```

## Expected State

Per `.github/instructions/portal-web-architecture.instructions.md` and `.github/instructions/i18n-translations.instructions.md`, components must bind a translator to a namespace and call keys without any `ns:` prefix.

**Expected Pattern:**
```tsx
// At component level
const { t: tCommon } = useTranslation("common")

// In JSX
confirmText={tCommon('actions.delete')}
cancelText={tCommon('actions.cancel')}
```

## Impact

- **Code Quality**: Violates established i18n patterns documented in SSOT
- **Consistency**: Creates inconsistency with other portal components that follow the canonical pattern
- **Maintainability**: Multi-colon pattern is explicitly forbidden in instructions

## Why This Matters

The portal-web architecture instructions explicitly state:

> ### Forbidden Patterns
> - Do not use `t("ns:key")` anywhere (e.g. `t("orders:order_number")`).
> - Do not use multi-colon strings like `t("errors:route:retry")`.

This pattern was likely copied from an older implementation before the i18n guidelines were standardized.

## Suggested Fix

Update both components to use namespace-bound translators:

```tsx
// TransactionQuickActions.tsx
import { useTranslation } from 'react-i18next'

export function TransactionQuickActions() {
  const { t: tAccounting } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')
  
  // Use tCommon('actions.delete') and tCommon('actions.cancel')
}

// AssetQuickActions.tsx
import { useTranslation } from 'react-i18next'

export function AssetQuickActions() {
  const { t: tAccounting } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')
  
  // Use tCommon('actions.delete') and tCommon('actions.cancel')
}
```

## Files to Change

1. `portal-web/src/features/accounting/components/TransactionQuickActions.tsx`
2. `portal-web/src/features/accounting/components/AssetQuickActions.tsx`

## Related

- SSOT: `.github/instructions/portal-web-architecture.instructions.md` (Section 8: Internationalization)
- SSOT: `.github/instructions/i18n-translations.instructions.md`

## Priority

**Medium** - Does not affect functionality but violates documented patterns and creates technical debt.

---

## Resolution

**Status:** ✅ Resolved  
**Date:** 2026-01-18  
**Approach Taken:** Option 1 (Updated code to match instructions)

### Harmonization Summary

Updated both accounting quick action components to use the canonical namespace-bound translator pattern instead of the forbidden multi-colon pattern.

### Pattern Applied

Both components now:
1. Bind translators to their respective namespaces: `const { t } = useTranslation('accounting')` and `const { t: tCommon } = useTranslation('common')`
2. Call translation keys without namespace prefixes: `tCommon('actions.delete')` instead of `t('common:actions.delete')`

### Files Changed

- `portal-web/src/features/accounting/components/TransactionQuickActions.tsx` - Added `tCommon` translator and replaced 2 multi-colon patterns
- `portal-web/src/features/accounting/components/AssetQuickActions.tsx` - Added `tCommon` translator and replaced 2 multi-colon patterns

### Migration Completeness

- Total instances found: 4 (2 files × 2 patterns each)
- Instances harmonized: 4
- Remaining drift: 0

### Validation

- [x] All tests pass (type-check: ✅, lint: ✅)
- [x] Type check passes
- [x] Lint passes  
- [x] Pattern applied consistently across both files
- [x] No regressions introduced
- [x] Instruction files aligned (existing instructions were clear)

### Instruction Files Reviewed

- `.github/instructions/portal-web-architecture.instructions.md` - Section 8 clearly documents the forbidden pattern
- `.github/instructions/i18n-translations.instructions.md` - SSOT for i18n patterns

### Prevention

This drift should not recur because instruction files already explicitly forbid the multi-colon pattern:

- Portal-web architecture instructions list `t("ns:key")` as a forbidden pattern
- i18n translations instructions document the canonical namespace-bound approach
- Both files now follow the documented pattern correctly
