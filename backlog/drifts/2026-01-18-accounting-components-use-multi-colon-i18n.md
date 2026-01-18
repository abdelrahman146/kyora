---
title: "Accounting components use multi-colon i18n pattern (t('common:actions.delete'))"
date: 2026-01-18
status: open
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
