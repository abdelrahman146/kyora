# Portal Web Shadow Usage Violates Design Tokens

**Type:** Drift (Code Consistency)
**Severity:** Medium
**Domain:** Portal Web / UI
**Date:** 2026-01-18

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

## Related

- Design Tokens SSOT: `.github/instructions/design-tokens.instructions.md`
- UI Implementation: `.github/instructions/ui-implementation.instructions.md`
