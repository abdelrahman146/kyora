# ShippingZoneInfo and Notes Components Violate Organisms Resource-Agnostic Rule

**Type:** Drift (Code Consistency)
**Severity:** Medium
**Domain:** Portal Web / Code Structure
**Date:** 2026-01-18

## Description

`components/organisms/ShippingZoneInfo.tsx` and `components/organisms/Notes.tsx` are resource/feature-specific components that violate the portal-web code structure SSOT rule that organisms must be resource-agnostic.

## Current State

Both components are placed in `portal-web/src/components/organisms/`:

1. **ShippingZoneInfo.tsx** (85 lines)
   - Imports `ShippingZone` type from `@/api/business`
   - Used only in order forms (EditOrderSheet, CreateOrderSheet)
   - Component name and props explicitly reference `ShippingZone` resource

2. **Notes.tsx** (399 lines)
   - Designed for customer notes, order notes, and entity note-taking
   - Used in CustomerDetailPage
   - While more generic than ShippingZoneInfo, it's still coupled to resource-specific behavior patterns

## Expected State

Per `.github/instructions/portal-web-code-structure.instructions.md`:

```markdown
### 1.2 Shared UI components

All shared components must be resource-agnostic.

- `portal-web/src/components/organisms/`
  - Higher-level reusable composites.
  - Still generic. If it's feature/resource-specific, it does not belong here.
  - Allowed examples: `Table`, `Header`, `Sidebar` (app chrome), generic filters/sorting.

**Strict naming rule:** shared component names must not include resource names.
```

Components referencing business resources (`ShippingZone`, note-taking patterns specific to customers/orders) are feature-specific.

## Impact

- **Architectural Clarity**: Blurs the boundary between shared UI primitives and feature components
- **Maintainability**: Other developers may incorrectly assume organisms can be feature-specific
- **Reusability**: Limits ability to maintain truly generic component library

## Affected Files

```
portal-web/src/components/organisms/ShippingZoneInfo.tsx (85 lines)
portal-web/src/components/organisms/Notes.tsx (399 lines)
```

**Usage:**

- ShippingZoneInfo: imported by `features/orders/components/EditOrderSheet.tsx` and `CreateOrderSheet.tsx`
- Notes: imported by `features/customers/components/CustomerDetailPage.tsx`

## Suggested Fix

### Option 1: Move to Feature Modules (Recommended)

Move components to their primary feature usage location:

1. **ShippingZoneInfo** → `portal-web/src/features/orders/components/ShippingZoneInfo.tsx`
   - Update imports in EditOrderSheet and CreateOrderSheet
   - If needed by other features later, can move to a shared orders utility

2. **Notes** → `portal-web/src/features/customers/components/Notes.tsx` OR create a dedicated `features/notes/` cross-cutting feature
   - Since Notes is designed to be reusable across entities (customers, orders, etc.), consider:
     - Creating `features/notes/components/Notes.tsx` as a cross-cutting feature module
     - This follows the pattern of `features/auth/`, `features/business-switcher/` (cross-cutting but feature-specific)

### Option 2: Make Truly Generic (Alternative)

If Notes is intended to be a primitive:
- Remove all resource-specific vocabulary from documentation
- Ensure interface is 100% generic (no assumptions about customer/order patterns)
- Keep current location

ShippingZoneInfo cannot be made generic as it's inherently tied to the ShippingZone business concept.

## Implementation Steps

1. Move `ShippingZoneInfo.tsx` to `features/orders/components/`
2. Update imports in:
   - `features/orders/components/EditOrderSheet.tsx`
   - `features/orders/components/CreateOrderSheet.tsx`
3. Decide on Notes placement (feature-specific location)
4. Move `Notes.tsx` to chosen location
5. Update import in `features/customers/components/CustomerDetailPage.tsx`
6. Remove from `components/organisms/index.ts` exports

## Related

- Code Structure SSOT: `.github/instructions/portal-web-code-structure.instructions.md`
- Portal Architecture: `.github/instructions/portal-web-architecture.instructions.md`
