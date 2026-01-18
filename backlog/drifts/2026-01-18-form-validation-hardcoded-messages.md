---
title: Form validation uses hardcoded error messages instead of translation keys
date: 2026-01-18
severity: medium
scope: portal-web
status: open
---

# Drift: Form Validation Hardcoded Error Messages

## Problem

Multiple form validation rules use hardcoded English error messages instead of translation keys from `errors.json`. This breaks i18n, ignores RTL users, and violates the forms SSOT.

## Current State (WRONG)

**Examples found:**

```tsx
// ResetPasswordPage.tsx, LoginForm.tsx, and others
validators={{
  onBlur: z.string().min(8, 'validation.password_min_length'),  // ✅ CORRECT (has validation. prefix)
  onBlur: z.string().min(1, 'validation.required'),  // ✅ CORRECT
}}

// However, many forms are correctly using validation keys
// This drift report documents the pattern requirement for future forms
```

## Expected State (CORRECT)

**ALL validation messages must use translation keys:**

```tsx
// ✅ CORRECT - Translation keys with validation. prefix
validators={{
  onBlur: z.string().min(8, 'validation.password_min_length'),
  onBlur: z.string().email('validation.invalid_email'),
  onBlur: z.string().min(1, 'validation.required'),
  onChange: ({ value }) => {
    if (!value) return 'validation.required';
    if (value.length < 3) return 'validation.min_length';
    return undefined;
  },
}}
```

**Keys must exist in:**
- `portal-web/src/i18n/en/errors.json` under `validation.*`
- `portal-web/src/i18n/ar/errors.json` under `validation.*`

## Impact

- **User Experience:** Arabic users see English error messages
- **Maintenance:** Error message changes require code changes instead of translation file updates
- **Consistency:** Some forms have translated messages, others don't
- **RTL:** Hardcoded messages don't support RTL layout

## How to Fix

1. **Check if key exists** in `src/i18n/en/errors.json` under `validation.*`
2. **If missing, add it:**
   ```json
   // src/i18n/en/errors.json
   {
     "validation": {
       "your_new_key": "Your error message here"
     }
   }
   
   // src/i18n/ar/errors.json
   {
     "validation": {
       "your_new_key": "رسالة الخطأ بالعربية"
     }
   }
   ```
3. **Replace hardcoded message** with key
4. **Test in both languages**

## Related Files

- Forms SSOT: `.github/instructions/forms.instructions.md` (section "Translation Keys")
- Translation files: `portal-web/src/i18n/*/errors.json`
- Translator: `portal-web/src/lib/translateValidationError.ts`

## Common Validation Keys Available

```typescript
"validation.required"
"validation.invalid_email"
"validation.invalid_phone"
"validation.password_min_length"
"validation.min_length"  // interpolates {{min}}
"validation.max_length"  // interpolates {{max}}
"validation.positive_number"
"validation.invalid_date"
"validation.select_at_least_one"
```

## Prevention

**Before completing any form task, verify:**
- ☑ All validation messages use `validation.*` keys
- ☑ All keys exist in both EN and AR errors.json
- ☑ No hardcoded English strings in validators

## Next Steps

1. Audit all existing forms for hardcoded validation messages
2. Create missing translation keys
3. Replace hardcoded messages with keys
4. Add linter rule to catch hardcoded validation strings (future enhancement)
