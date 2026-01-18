# Portal Web Shadow Usage Violates Design Tokens

**Type:** Drift (Code Consistency)
**Severity:** Medium
**Domain:** Portal Web / UI
**Date:** 2026-01-18
**Status:** ✅ Resolved
**Resolution Date:** 2026-01-19

## Description

Portal Web codebase uses shadow utilities in multiple locations, violating the design tokens rule "No Shadows" which explicitly states "Do not use shadow utilities for elevation."

## Current State

Shadows are used in 4+ locations:

1. **AssetListPage.tsx** (line 151): `className="btn-circle shadow-lg"`
2. **CapitalListPage.tsx** (line 253): `className="btn-circle shadow-lg"`
3. **AdvisorPanel.tsx** (line 116): `hover:shadow-sm`
4. **DatePicker.tsx** (line 459): `shadow-2xl`

## Expected State

Per `.github/instructions/design-tokens.instructions.md`:

```
## Shadows (No Shadows)

- Do not use shadow utilities for elevation.
- Do not introduce "elevation" via `shadow-*` or `drop-shadow-*`.
```

All elevation and separation should come from **spacing + typography + borders**.

## Impact

- **Consistency**: Visual inconsistency across portal UI
- **Design System Integrity**: Undermines the minimal, calm aesthetic
- **Maintainability**: Other developers may assume shadows are acceptable

## Affected Files

```
portal-web/src/features/accounting/components/AssetListPage.tsx:151
portal-web/src/features/accounting/components/CapitalListPage.tsx:253
portal-web/src/components/molecules/AdvisorPanel.tsx:116
portal-web/src/components/form/DatePicker.tsx:459
```

## Suggested Fix

1. **Remove all shadow utilities** from the affected files
2. **Replace with borders** where visual separation is needed:
   - `shadow-lg` → `border border-base-300`
   - `hover:shadow-sm` → `hover:border-primary/50`
   - `shadow-2xl` → `border-2 border-base-300`

Example fix for AssetListPage.tsx:
```tsx
// Before
className="btn-circle shadow-lg"

// After
className="btn-circle border-2 border-base-300"
```

## Resolution

**Status:** Resolved
**Date:** 2026-01-19
**Approach Taken:** Option 1 (Update code to match instructions)

### Harmonization Summary

All shadow utilities have been removed from portal-web and replaced with appropriate border utilities according to the design tokens SSOT.

### Pattern Applied

**Elevation and separation now use borders exclusively:**
- FAB buttons: `border-2 border-base-300` instead of `shadow-lg`
- Interactive hover states: `hover:border-primary/50` instead of `hover:shadow-sm`
- Modals/overlays: `border-2 border-base-300` instead of `shadow-2xl`
- Cards with borders: `border border-base-300` for base separation

### Files Changed

- `portal-web/src/features/accounting/components/AssetListPage.tsx` - Replaced `shadow-lg` with `border-2 border-base-300` on FAB button
- `portal-web/src/features/accounting/components/CapitalListPage.tsx` - Replaced `shadow-lg` with `border-2 border-base-300` on FAB button
- `portal-web/src/components/molecules/AdvisorPanel.tsx` - Replaced `hover:shadow-sm` with `border border-base-300 hover:border-primary/50` and adjusted transition class
- `portal-web/src/components/form/DatePicker.tsx` - Replaced `shadow-2xl` with `border-2 border-base-300` on modal popup

### Migration Completeness

- Total shadow instances found: 4
- Instances harmonized: 4
- Remaining drift: 0

### Validation

- [x] All shadow utilities removed (grep search confirms 0 matches)
- [x] Type check passes (`npm run type-check`)
- [x] Lint passes (`npm run lint`)
- [x] Pattern applied consistently across all files
- [x] No regressions introduced
- [x] Instruction files updated with strengthened rules

### Instruction Files Updated

- `.github/instructions/design-tokens.instructions.md` - Enhanced "Shadows (No Shadows)" section with:
  - Explicit MUST/MUST NOT rules
  - Common mistake examples (❌ Wrong patterns)
  - Correct pattern examples (✅ Right patterns)
  - Affected areas documentation
  - Reference to this resolved drift report

### Prevention

This drift should not recur because design-tokens.instructions.md now explicitly:

- Documents the exact anti-patterns that caused this drift (using `shadow-*` utilities)
- Provides concrete before/after examples for common use cases (FAB buttons, hover states, modals)
- Lists affected areas where the rule applies
- Explains WHY shadows are prohibited (minimal aesthetic)
- Specifies the correct alternative (borders with appropriate classes)

## Related

- Design Tokens SSOT: `.github/instructions/design-tokens.instructions.md`
- UI Implementation: `.github/instructions/ui-implementation.instructions.md`
