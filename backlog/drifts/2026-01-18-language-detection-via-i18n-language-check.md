# Drift: Components checking `i18n.language` directly instead of using `useLanguage` hook

**Date**: 2026-01-18  
**Category**: i18n/translations  
**Severity**: Low (code duplication, but functionally correct)  
**Status**: Open

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
