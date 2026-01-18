---
title: "Business Schema Duplication Across Portal-Web API Files"
date: 2026-01-18
priority: low
category: consistency
status: open
domain: portal-web
---

# Business Schema Duplication Across Portal-Web API Files

## Summary

Portal-web has duplicate `Business` schema definitions in two separate files (`portal-web/src/api/business.ts` and `portal-web/src/api/types/business.ts`) with different field optionality and structure. This violates the "single SSOT schema per domain" rule stated in business-management.instructions.md and creates risk of silent contract drift.

## Current State

### Location 1: portal-web/src/api/business.ts (lines 13-47)

```typescript
export const BusinessSchema = z.object({
  id: z.string(),
  workspaceId: z.string(),
  descriptor: z.string(),
  name: z.string(),
  brand: z.string(),
  logo: AssetReferenceSchema.optional().nullable(),
  countryCode: z.string(),
  currency: z.string(),
  storefrontPublicId: z.string(),
  storefrontEnabled: z.boolean(),
  storefrontTheme: z.object({
    primaryColor: z.string(),
    secondaryColor: z.string(),
    accentColor: z.string(),
    backgroundColor: z.string(),
    textColor: z.string(),
    fontFamily: z.string(),
    headingFontFamily: z.string(),
  }),
  supportEmail: z.string(),
  phoneNumber: z.string(),
  whatsappNumber: z.string(),
  address: z.string(),
  websiteUrl: z.string(),
  instagramUrl: z.string(),
  facebookUrl: z.string(),
  tiktokUrl: z.string(),
  xUrl: z.string(),
  snapchatUrl: z.string(),
  vatRate: z.string(),
  safetyBuffer: z.string(),
  establishedAt: z.string(),
  archivedAt: z.string().nullable().optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
})
```

All fields are **required** (non-optional).

### Location 2: portal-web/src/api/types/business.ts (lines 21-58)

```typescript
export const BusinessSchema = z.object({
  id: z.string(),
  workspaceId: z.string(),
  name: z.string(),
  descriptor: z.string(),
  brand: z.string().optional(),
  logo: AssetReferenceSchema.optional().nullable(),
  phoneNumber: z.string().optional(),
  whatsappNumber: z.string().optional(),
  supportEmail: z.string().optional(),
  websiteUrl: z.string().optional(),
  facebookUrl: z.string().optional(),
  instagramUrl: z.string().optional(),
  xUrl: z.string().optional(),
  tiktokUrl: z.string().optional(),
  snapchatUrl: z.string().optional(),
  address: z.string().optional(),
  countryCode: z.string().optional(),
  currency: z.string().optional(),
  vatRate: z.string().optional(),
  safetyBuffer: z.string().optional(),
  establishedAt: z.string().optional(),
  storefrontEnabled: z.boolean().optional(),
  storefrontPublicId: z.string().optional(),
  storefrontTheme: StorefrontThemeSchema.optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
  archivedAt: z.string().optional().nullable(),
})
```

Most fields are **optional**.

### Key Differences

1. **Optionality**: Location 1 treats most fields as required; Location 2 treats most as optional
2. **Theme structure**: Location 1 has inline anonymous object; Location 2 references `StorefrontThemeSchema`
3. **Which is used**: `portal-web/src/api/business.ts` is imported and used throughout the portal

## Expected State

Per business-management.instructions.md:

> - Co-locate query options factories in `portal-web/src/api/business.ts` (`businessQueries.*`) and use them in route loaders and components.
> - Prefer aligning Zod schemas with backend JSON shapes and keeping a **single SSOT schema** per domain.

## Impact

- **Low**: The duplication exists but `portal-web/src/api/business.ts` schema is consistently used across the portal
- Risk of confusion when maintaining schemas
- Risk of importing wrong schema if file structure changes
- `portal-web/src/api/types/business.ts` schema appears to be unused

## Affected Files

- `portal-web/src/api/business.ts` (active schema)
- `portal-web/src/api/types/business.ts` (duplicate/possibly unused)

## Suggested Fix

### Option 1: Consolidate to Single File (Recommended)

1. Keep `BusinessSchema` in `portal-web/src/api/business.ts`
2. Align field optionality with backend `BusinessResponse` in `backend/internal/domain/business/model_response.go`
3. Delete or repurpose `portal-web/src/api/types/business.ts` (it also contains asset type re-exports which should live elsewhere)

### Option 2: Use Dedicated Types File

1. Move schema to `portal-web/src/api/types/business.ts`
2. Import in `portal-web/src/api/business.ts`
3. Ensure single definition

### Verification

After fix:
1. Search for all `BusinessSchema` definitions (should find only one)
2. Verify all imports resolve to the single SSOT
3. Align optionality with backend response (most fields are non-optional in `BusinessResponse`)

## References

- business-management.instructions.md (Known portal drift section)
- Backend SSOT: `backend/internal/domain/business/model_response.go` (BusinessResponse)
- Active schema: `portal-web/src/api/business.ts`
- Duplicate: `portal-web/src/api/types/business.ts`
