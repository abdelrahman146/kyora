# Form Field Components Import Resource-Specific Types Violating Generic Component Rule

**Type:** Drift (Code Consistency)
**Severity:** Medium
**Domain:** Portal Web / Code Structure / Forms
**Date:** 2026-01-18

## Description

Form field components in `portal-web/src/lib/form/components/` import and depend on resource-specific types and API hooks from `@/api/inventory` and `@/api/customer`, violating the code structure SSOT rule that `lib/**` components must be cross-cutting and generic.

## Current State

Three form field components in `lib/form/components/` are tightly coupled to specific business resources:

1. **CategorySelectField.tsx**
   - Imports: `type { Category } from '@/api/inventory'`
   - Imports: `useCategoriesQuery, useCreateCategoryMutation from '@/api/inventory'`
   - Uses inventory-specific translation namespace: `t('category_created', { ns: 'inventory' })`
   - Uses inventory-specific query keys: `queryKeys.inventory.all`

2. **CustomerSelectField.tsx**
   - Imports: `type { Customer } from '@/api/customer'`
   - Imports: `useCustomerQuery, useCustomersQuery from '@/api/customer'`

3. **ProductVariantSelectField.tsx**
   - Imports: `type { Variant } from '@/api/inventory'`
   - Imports: `useVariantQuery, useVariantsQuery from '@/api/inventory'`

## Expected State

Per `.github/instructions/portal-web-code-structure.instructions.md`:

```markdown
## 2) `lib/` rules (strict)

`portal-web/src/lib/**` is reserved for cross-cutting utilities only.

A utility belongs in `lib/` only if all are true:
- Used by 2+ distinct features OR clearly cross-cutting
- Not tied to a single resource vocabulary (inventory/orders/customers/etc.)
- Has stable API and is expected to be shared
```

And per the forms instruction section:

```markdown
- `portal-web/src/components/form/`
  - Generic UI form controls **not** tied to TanStack Form.
  - Controls should be usable with any form library (TanStack Form, RHF, uncontrolled).
```

These field components are feature-specific because they:
- Reference domain-specific types (Category, Customer, Variant)
- Use domain-specific API hooks
- Use domain-specific i18n namespaces
- Are tightly coupled to specific business resources

## Impact

- **Architectural Clarity**: Blurs the distinction between generic form controls and feature-specific form fields
- **Reusability**: Cannot reuse these components for different resources without inventory/customer domain knowledge
- **Dependency Management**: Creates unexpected coupling from `lib/` to domain-specific `api/` modules
- **Maintainability**: Other developers may incorrectly assume `lib/form/components/` can contain resource-specific logic

## Affected Files

```
portal-web/src/lib/form/components/CategorySelectField.tsx
portal-web/src/lib/form/components/CustomerSelectField.tsx
portal-web/src/lib/form/components/ProductVariantSelectField.tsx
```

## Suggested Fix

### Option 1: Move to Feature Modules (Recommended)

Move these field components to their respective feature modules:

1. **CategorySelectField** → `features/inventory/components/fields/CategorySelectField.tsx`
2. **ProductVariantSelectField** → `features/inventory/components/fields/ProductVariantSelectField.tsx`
3. **CustomerSelectField** → `features/customers/components/fields/CustomerSelectField.tsx`

Update all imports across the codebase.

### Option 2: Create Generic Wrapper + Feature-Specific Implementations

If the underlying select pattern is reusable:

1. Keep a **truly generic** `AsyncSelectField` in `lib/form/components/`
   - Accepts: `fetchOptions`, `createOption`, generic types
   - No domain-specific imports

2. Create feature-specific wrappers:
   - `features/inventory/components/fields/CategorySelectField.tsx` wraps `AsyncSelectField` with inventory-specific behavior
   - `features/inventory/components/fields/ProductVariantSelectField.tsx` wraps `AsyncSelectField`
   - `features/customers/components/fields/CustomerSelectField.tsx` wraps `AsyncSelectField`

### Option 3: Create Cross-Cutting Feature Module (Alternative)

If these fields are truly used across many features (orders, customers, inventory):
- Create `features/resource-select-fields/` as a cross-cutting feature module
- Move all three components there
- This acknowledges they are feature-specific while maintaining shared location

## Implementation Steps (Option 1 - Recommended)

1. Create `features/inventory/components/fields/` directory
2. Move `CategorySelectField.tsx` and `ProductVariantSelectField.tsx` to `features/inventory/components/fields/`
3. Create `features/customers/components/fields/` directory
4. Move `CustomerSelectField.tsx` to `features/customers/components/fields/`
5. Update all imports (search for `@/lib/form/components/CategorySelectField` etc.)
6. Remove from `lib/form/components/index.ts` exports
7. Add to respective feature index exports

## Related

- Code Structure SSOT: `.github/instructions/portal-web-code-structure.instructions.md`
- Forms System: `.github/instructions/forms.instructions.md`
- Portal Architecture: `.github/instructions/portal-web-architecture.instructions.md`
