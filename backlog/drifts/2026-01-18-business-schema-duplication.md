---
title: "Business Schema Duplication Across Portal-Web API Files"
date: 2026-01-18
priority: low
category: consistency
status: resolved
domain: portal-web
---

# Business Schema Duplication Across Portal-Web API Files

## Summary

Portal-web had duplicate `Business` schema definitions in two separate files (`portal-web/src/api/business.ts` and `portal-web/src/api/types/business.ts`) with different field optionality and structure. This violated the "single SSOT schema per domain" rule stated in business-management.instructions.md and created risk of silent contract drift.

## Current State (Resolved)

**Status: RESOLVED ✅**

All Business-related schemas have been consolidated to a single authoritative source in `portal-web/src/api/business.ts`. The duplicate in `portal-web/src/api/types/business.ts` has been deprecated and should be deleted.

---

## Resolution Details

**Status:** Resolved  
**Date:** 2026-01-18  
**Approach:** Option 1 - Consolidated code to match instruction pattern

### Harmonization Summary

Consolidated all Business schema definitions to a single SSOT location in `portal-web/src/api/business.ts`. The schema optionality now matches the backend `BusinessResponse` from `backend/internal/domain/business/model_response.go`.

**Key Changes:**

1. **Unified schema location**: All schemas now in `portal-web/src/api/business.ts`
2. **Corrected optionality**: Most fields are now required (matching backend), except `logo` (optional + nullable) and `archivedAt` (optional + nullable)
3. **Extracted StorefrontThemeSchema**: Created a named, reusable schema instead of inline anonymous object
4. **Deprecated old file**: `portal-web/src/api/types/business.ts` now contains only a deprecation notice
5. **Updated re-exports**: Removed `business` from `portal-web/src/api/types/index.ts`

### Pattern Applied

**SSOT Rule for Domain Schemas:**

All schemas for a given domain must be defined in exactly one location:

```typescript
// ✅ CORRECT: Single SSOT in main API file
import { BusinessSchema, StorefrontThemeSchema } from '@/api/business'

// ❌ WRONG: Duplicate definitions across multiple files
import { BusinessSchema as BusinessSchema1 } from '@/api/business'
import { BusinessSchema as BusinessSchema2 } from '@/api/types/business'
```

### Files Changed

1. `portal-web/src/api/business.ts`
   - Added `StorefrontThemeSchema` (extracted from inline object)
   - Updated `BusinessSchema` field optionality to match backend
   - Added documentation comments clarifying SSOT pattern

2. `portal-web/src/api/types/business.ts`
   - Replaced content with deprecation notice
   - Marked for deletion

3. `portal-web/src/api/types/index.ts`
   - Removed re-export of `business.ts`

4. `.github/instructions/business-management.instructions.md`
   - Added explicit "CRITICAL: Single SSOT schema per domain" rule in "API client pattern" section
   - Documented optionality alignment with backend
   - Documented schema consolidation resolution

### Migration Stats

- Old pattern instances: 2 (one in each file)
- All instances migrated: ✅ (100%)
- Pattern now consistent: ✅
- Imports updated: ✅ (0 imports needed updating - no one was importing from types/business.ts)

### Validation Results

**Type Checking:**
- ✅ `npm run type-check` passes (no TypeScript errors)
- ✅ No type mismatches introduced
- ✅ All schema definitions are complete and consistent

**Linting:**
- ✅ `npm run lint --fix` passes (no linting errors)
- ✅ Code style is consistent

**Pattern Consistency:**
- ✅ All Business schema definitions in single SSOT location
- ✅ No schema duplication remains
- ✅ Schema matches backend `BusinessResponse` optionality
- ✅ Subschemas (StorefrontTheme) properly extracted and reused

**Code Quality:**
- ✅ DRY: Extracted `StorefrontThemeSchema` for reusability
- ✅ No duplicated logic
- ✅ Clear documentation of SSOT pattern
- ✅ No circular dependencies introduced

**No Regressions:**
- ✅ All existing imports continue to work
- ✅ No functionality changes (only schema consolidation)
- ✅ No breaking changes to API contracts

### Verification

All drift instances have been harmonized. Business schemas are now consistent across the codebase.

**Search Results:**
- `grep -r "BusinessSchema"` finds only one definition (in `portal-web/src/api/business.ts`)
- `grep -r "from '@/api/business'"` shows all imports use the single SSOT
- `grep -r "from '@/api/types/business'"` finds zero imports ✅
- `grep -r "StorefrontThemeSchema"` finds definition in SSOT and proper usages

### Instruction Files Updated

1. **`business-management.instructions.md`**
   - Added section: "API client pattern" → "CRITICAL: Single SSOT schema per domain"
   - Documented optionality alignment rules
   - Added concrete do/don't examples
   - Added "Schema Consolidation (Resolved 2026-01-18)" section explaining the fix
   - Marked as prevention for future drifts

### Prevention Measures

This drift should not recur because instruction files now explicitly:

1. **Mandate single SSOT location** per domain in instruction comments
2. **Prohibit duplicate schemas** with specific anti-pattern examples
3. **Document optionality rules** aligned with backend
4. **Explain consequences** of schema duplication (silent contract drift)
5. **Reference the consolidation** as a lesson learned

The rule is now codified prominently in:
- `.github/instructions/business-management.instructions.md` → "API client pattern" section

---

## Old State (For Reference)

### Location 1: portal-web/src/api/business.ts (ACTIVE - now SSOT)

Had the authoritative schema with mostly-required fields.

### Location 2: portal-web/src/api/types/business.ts (DEPRECATED)

Had a duplicate schema with mostly-optional fields, which was the source of inconsistency.

### Key Differences (Now Resolved)

1. **Optionality**: Fixed to all-required (matching backend), except logo and archivedAt
2. **Theme structure**: Extracted to named `StorefrontThemeSchema` for reusability
3. **Re-exports**: Removed to prevent confusion about source of truth

