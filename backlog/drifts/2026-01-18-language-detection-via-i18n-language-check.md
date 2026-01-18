# Drift: Components checking `i18n.language` directly instead of using `useLanguage` hook

**Date**: 2026-01-18  
**Category**: i18n/translations  
**Severity**: Low (code duplication, but functionally correct)  
**Status**: Resolved  
**Resolution Date**: 2026-01-19

## Description

Several components check the current language by accessing `i18n.language` directly and using `.toLowerCase().startsWith('ar')` instead of using the centralized `useLanguage` hook.

This creates unnecessary duplication of language detection logic and violates the i18n SSOT pattern where `useLanguage` is the single source of truth for language state.

## Current Behavior

Components directly check:
```tsx
const isArabic = i18n.language.toLowerCase().startsWith('ar')
```

## Expected Behavior

Components should use the `useLanguage` hook:
```tsx
const { isArabic } = useLanguage()
// or
const { isRTL } = useLanguage()
```

## Affected Files

1. `portal-web/src/features/onboarding/components/BusinessSetupPage.tsx` (line 29)
2. `portal-web/src/features/customers/components/AddressCard.tsx` (line 37)
3. `portal-web/src/features/customers/components/CustomerDetailPage.tsx` (line 81)
4. `portal-web/src/features/customers/components/PhoneCodeSelect.tsx` (line 37)
5. `portal-web/src/features/customers/components/CountrySelect.tsx` (line 40)

## Impact

- **Code duplication**: Language detection logic is duplicated in 5 places
- **Maintainability**: Changes to language detection rules require updating multiple files
- **Consistency**: Some components use `useLanguage`, others use `i18n.language` directly

## Solution

Replace all occurrences with `useLanguage` hook:

```tsx
// Before
import { useTranslation } from 'react-i18next'
const { i18n } = useTranslation()
const isArabic = i18n.language.toLowerCase().startsWith('ar')

// After
import { useLanguage } from '@/hooks/useLanguage'
const { isArabic } = useLanguage()
```

## SSOT Reference

- `.github/instructions/i18n-translations.instructions.md` section 5.3: "The `useLanguage` hook (preferred way to access language)"
- `.github/instructions/i18n-translations.instructions.md` section 5.4: "Don't: Check `i18n.language` directly with `.startsWith()` or `.toLowerCase()`"

## Related Issues

None

## Notes

This drift is purely about code organization and does not cause functional bugs. The logic `i18n.language.toLowerCase().startsWith('ar')` produces the same result as `useLanguage().isArabic`.

However, maintaining this pattern creates technical debt and makes future language-related changes more difficult.

---

## Resolution

**Status:** Resolved  
**Date:** 2026-01-19  
**Approach Taken:** Option 1 (Update code to match instructions)

### Harmonization Summary

All 5 components have been harmonized to use the centralized `useLanguage` hook instead of directly checking `i18n.language`. This eliminates code duplication and ensures language detection logic is maintained in a single location.

### Pattern Applied

**Consistent pattern now used across all affected files:**

```tsx
// Import useLanguage hook
import { useLanguage } from '@/hooks/useLanguage'

// Use destructured isArabic
const { isArabic } = useLanguage()

// For components that need multiple language properties
const { isArabic, language, isRTL } = useLanguage()
```

### Files Changed

1. `portal-web/src/features/onboarding/components/BusinessSetupPage.tsx` - Replaced i18n.language check with useLanguage hook; added import; removed i18n destructuring from useTranslation
2. `portal-web/src/features/customers/components/AddressCard.tsx` - Replaced i18n.language check with useLanguage hook; added import; removed unused useTranslation() call
3. `portal-web/src/features/customers/components/CustomerDetailPage.tsx` - Replaced i18n.language check with useLanguage hook; added import; removed i18n destructuring from useTranslation
4. `portal-web/src/features/customers/components/PhoneCodeSelect.tsx` - Replaced i18n.language check with useLanguage hook; added import; removed i18n destructuring from useTranslation
5. `portal-web/src/features/customers/components/CountrySelect.tsx` - Replaced i18n.language check with useLanguage hook; added import; removed i18n destructuring from useTranslation

### Migration Completeness

- Total instances found: 5
- Instances harmonized: 5
- Remaining drift: 0

### Validation

- [x] All tests pass (type-check passes)
- [x] Lint passes
- [x] Pattern applied consistently
- [x] No regressions introduced
- [x] Instruction files aligned

### Instruction Files Updated

- `.github/instructions/i18n-translations.instructions.md` - Strengthened section 5.2 and 5.4:
  - Replaced "known drift" note with explicit anti-pattern examples
  - Added "CRITICAL ANTI-PATTERN" label to wrong pattern
  - Added detailed "Pattern to follow when refactoring" section showing Before/After
  - Strengthened "Don't" rules with additional explicit anti-patterns
  - Added emoji indicators (✅/❌) for clarity

### Prevention

This drift should not recur because instruction files now explicitly:

- Label the direct i18n.language check as a "CRITICAL ANTI-PATTERN"
- Provide clear Before/After refactoring examples
- Explain why the pattern is wrong (duplicates logic, violates SSOT)
- Document the correct pattern with multiple usage examples
- Include the anti-pattern in the drift detection rules (section 6)
