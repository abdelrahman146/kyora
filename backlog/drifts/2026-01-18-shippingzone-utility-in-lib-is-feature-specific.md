# shippingZone.ts Utility in lib/ is Feature-Specific

**Type:** Drift (Code Consistency)
**Severity:** Low
**Domain:** Portal Web / Code Structure
**Date:** 2026-01-18

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

### Option 1: Move to Orders Feature (Recommended if primarily used in orders)

```
portal-web/src/features/orders/utils/shippingZone.ts
```

Update all imports to `@/features/orders/utils/shippingZone`.

### Option 2: Move to Business Feature (If used in business settings)

If the utility is primarily for business/shipping zone management:

```
portal-web/src/features/business/utils/shippingZone.ts
```

or create a dedicated business feature if it doesn't exist.

### Option 3: Keep in lib/ Only If Truly Cross-Cutting

Only if the utility is demonstrably used across 3+ distinct features (orders, customers, business settings, analytics, etc.) with no single "owner" feature.

In that case:
- Rename to make it clear it's business-domain cross-cutting (e.g., `lib/business/shippingZone.ts`)
- Document why it's cross-cutting in the file header

## Implementation Steps (Option 1 - Recommended)

1. Create `features/orders/utils/` directory if it doesn't exist
2. Move `lib/shippingZone.ts` to `features/orders/utils/shippingZone.ts`
3. Search for all imports: `grep -r "from '@/lib/shippingZone'" portal-web/src/`
4. Update imports to `@/features/orders/utils/shippingZone`
5. Verify all imports still resolve correctly

## Related

- Code Structure SSOT: `.github/instructions/portal-web-code-structure.instructions.md`
- Existing precedent: `features/inventory/utils/inventoryUtils.ts`, `features/onboarding/utils/onboarding.ts`
