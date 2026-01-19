# shippingZone.ts Utility in lib/ is Feature-Specific

**Type:** Drift (Code Consistency)
**Severity:** Low
**Domain:** Portal Web / Code Structure
**Date:** 2026-01-18
**Status:** ✅ Resolved (2026-01-19)

## Description

`portal-web/src/lib/shippingZone.ts` is a feature-specific utility file that violates the code structure SSOT rule that `lib/**` is reserved for cross-cutting utilities only.

## Current State

The file `lib/shippingZone.ts` (108 lines) contains utilities for:
- Inferring shipping zones from customer addresses
- Checking if countries are in zones
- Formatting shipping zone information
- Getting zone countries

All functions are tightly coupled to the `ShippingZone` business concept:

```typescript
import type { ShippingZone } from '@/api/business'
import type { CustomerAddress } from '@/api/customer'

export function inferShippingZoneFromAddress(
  address: CustomerAddress | Pick<CustomerAddress, 'countryCode'> | null | undefined,
  zones: Array<ShippingZone>,
): ShippingZone | undefined { ... }

export function isCountryInZone(
  countryCode: string,
  zone: ShippingZone | null | undefined,
): boolean { ... }

export function formatShippingZoneInfo(zone: ShippingZone, currency?: string) { ... }

export function getZoneCountries(zone: ShippingZone | null | undefined): Array<string> { ... }
```

## Expected State

Per `.github/instructions/portal-web-code-structure.instructions.md`:

```markdown
## 2) `lib/` rules (strict)

`portal-web/src/lib/**` is reserved for cross-cutting utilities only.

A utility belongs in `lib/` only if all are true:
- Used by 2+ distinct features OR clearly cross-cutting (e.g. error parsing, query keys, routing guards)
- Not tied to a single resource vocabulary (inventory/orders/customers/etc.)
- Has stable API and is expected to be shared

If a file is feature-specific (e.g. `inventoryUtils`, `onboarding`, `customers*`), it must live under `features/<feature>/utils/`.
```

`shippingZone.ts` is tied to shipping zone vocabulary, which is primarily an orders/business feature concern, not a generic cross-cutting utility.

## Impact

- **Architectural Clarity**: Blurs the boundary between generic utilities and feature-specific helpers
- **Discoverability**: Developers looking for orders/shipping utilities may not expect them in `lib/`
- **Maintainability**: Inconsistent with other feature utilities (e.g., `features/inventory/utils/inventoryUtils.ts`, `features/onboarding/utils/onboarding.ts`)

## Affected Files

```
portal-web/src/lib/shippingZone.ts (108 lines)
```

**Potential usage:**
- Likely used in order forms (CreateOrderSheet, EditOrderSheet)
- May be used in business/shipping zone management components

## Suggested Fix

###  Move to Business Feature (If used in business settings)

```
portal-web/src/features/business/utils/shippingZone.ts
```

## Related

- Code Structure SSOT: `.github/instructions/portal-web-code-structure.instructions.md`
- Existing precedent: `features/inventory/utils/inventoryUtils.ts`, `features/onboarding/utils/onboarding.ts`

---

## Resolution

**Status:** ✅ Harmonized  
**Date:** 2026-01-19  
**Approach:** Option 1 - Updated code to match instructions

### Harmonization Summary

The utility was moved from the generic `lib/` folder to the orders feature where it is actually used. This eliminates confusion about library vs. feature-specific code organization.

### Pattern Applied

All feature-specific utilities must live under `features/<feature>/utils/` per the code structure SSOT. The shipping zone utilities are tightly coupled to the orders domain (used exclusively by CreateOrderSheet and EditOrderSheet), so they belong in the orders feature.

### Files Changed

- ✅ Created: `portal-web/src/features/orders/utils/shippingZone.ts` (moved from `lib/`)
- ✅ Updated: `portal-web/src/features/orders/components/EditOrderSheet.tsx` (import path)
- ✅ Updated: `portal-web/src/features/orders/components/CreateOrderSheet.tsx` (import path)
- ✅ Updated: `.github/instructions/portal-web-code-structure.instructions.md` (strengthened rule with explicit anti-pattern example)

### Migration Completeness

- Total instances found: 2 (both in orders feature components)
- All instances migrated: ✅
- Pattern now consistent: ✅

### Validation

- [x] Imports updated in all affected components
- [x] Old import path no longer exists in codebase
- [x] New path follows code structure guidelines
- [x] Pattern matches instruction file requirements
- [x] No regressions (utilities still work the same way)
- [x] All related code references the new location

### Instruction Files Updated

- `portal-web-code-structure.instructions.md`
  - **Added:** Explicit anti-pattern showing wrong placement in `lib/`
  - **Added:** Correct pattern showing feature-specific utility placement
  - **Added:** Examples with code showing before/after

### Prevention

This drift should not recur because instruction files now explicitly:

- Show **anti-pattern example** of feature-specific utility in `lib/` (was implicit)
- Show **correct pattern** with actual code examples (was unclear)
- Reference this specific case in "Known drifts" section
- Document that utilities tied to resource vocabulary belong in features

The rule is now **explicit with concrete examples** rather than implicit guidance.
