---
name: implement-drift-fix
description: "Harmonize code to fix pattern drift with comprehensive validation and consistency checks"
agent: "Feature Builder"
argument-hint: "Path to drift report (e.g., backlog/drifts/2026-01-18-customer-api-snake-case-drift.md)"
tools: ["vscode", "execute", "read", "edit", "search", "agent", "todo"]
---

# Drift Fix Implementation

You are about to harmonize code to fix the drift described in:

**Drift Report:** `${input:driftReport:Path to drift report (e.g., backlog/drifts/2026-01-18-issue.md)}`

## Mission

Harmonize code to eliminate pattern drift by:

1. Understanding the drift and authoritative pattern
2. Implementing consistent pattern across all affected code
3. Validating against instruction files
4. Ensuring no functionality breaks
5. Verifying complete harmonization

## Workflow

### Phase 1: Understanding & Planning

1. **Read the drift report completely**
   - Understand current state vs. expected pattern
   - Identify all affected locations
   - Review harmonization plan
   - Note the recommended approach (Option 1 or 2)

2. **Confirm authoritative pattern**
   - Read instruction files referenced
   - Find examples of correct pattern in codebase
   - Understand why this pattern is preferred

3. **Assess scope**
   - Isolated (1-2 files)
   - Localized (one domain/feature)
   - Systemic (multiple domains)

4. **Create implementation plan**
   - List all files to modify
   - Determine if refactoring needed
   - Plan for breaking changes (if any)
   - Identify test updates needed

### Phase 2: Implementation

#### Strategy Selection

**Option 1: Update code to match instructions** (most common)

- Apply the authoritative pattern from instruction files
- Update all affected locations
- Ensure consistency across the codebase

**Option 2: Update instructions to match code** (rare)

- Only if code pattern is demonstrably better
- Update instruction file
- Verify no conflicts with other SSOT files
- Document reasoning

#### For Backend Harmonization:

1. **Common drift types:**

   **Naming conventions:**
   - snake_case vs camelCase in JSON
   - Inconsistent function/variable names
   - Model field naming

   **Error handling:**
   - Different error response shapes
   - Inconsistent ProblemDetails usage
   - Mixed error handling patterns

   **Data access:**
   - Raw SQL vs repository scopes
   - Inconsistent transaction wrapping
   - Mixed query patterns

2. **Harmonization patterns:**
   - Create reusable helpers if pattern is complex
   - Use bulk rename/refactor for widespread changes
   - Update DTOs for API contract changes
   - Regenerate Swagger/OpenAPI after API changes

3. **Breaking change handling:**
   - If changing API contracts, update both backend + frontend
   - Update E2E tests for new contracts
   - Consider versioning if external consumers exist

#### For Portal-Web Harmonization:

1. **Common drift types:**

   **Component structure:**
   - Components in wrong folders
   - Inconsistent file organization
   - Mixed feature structure

   **State management:**
   - Multiple patterns for same state type
   - Mixed use of URL/Query/Store/Form
   - Inconsistent hook usage

   **API calls:**
   - Direct fetch vs apiClient
   - Inconsistent error handling
   - Mixed response typing

   **UI patterns:**
   - Inconsistent form components
   - Mixed button/input styles
   - Different modal patterns

2. **Harmonization patterns:**
   - Move files to correct locations
   - Extract reusable components
   - Standardize on TanStack patterns
   - Create shared hooks for common logic

3. **Type safety:**
   - Update types after API changes
   - Ensure Zod schemas match backend
   - Fix any type assertions

### Phase 3: Comprehensive Migration

1. **Find all instances**
   - Use grep/semantic search extensively
   - Don't miss edge cases
   - Check related files

2. **Apply pattern consistently**
   - Same pattern everywhere
   - No half-migrations
   - Update all variations

3. **Update related code**
   - Tests using the old pattern
   - Documentation mentioning it
   - Related utilities/helpers

### Phase 4: Testing

#### Backend Testing:

1. **Run existing E2E tests**
   - Execute: `cd backend && go test ./internal/tests/e2e -v`
   - Verify 100% pass
   - Fix any tests broken by harmonization

2. **Update tests if needed**
   - Update assertions for new patterns
   - Add tests for harmonized behavior
   - Ensure coverage is maintained

3. **Regression testing**
   - Verify related features still work
   - Check for unexpected side effects
   - Test edge cases

#### Portal-Web Testing:

1. **Type checking**
   - Execute: `cd portal-web && npm run type-check`
   - Fix all TypeScript errors
   - Ensure types reflect new patterns

2. **Linting**
   - Execute: `cd portal-web && npm run lint -- --fix`
   - Fix all linting errors
   - Ensure consistent code style

3. **Functional validation**
   - [ ] All features work as before
   - [ ] No console errors
   - [ ] No network errors
   - [ ] Forms submit correctly
   - [ ] Navigation works

### Phase 5: Validation & Consistency Checks

#### Backend Validation:

1. **Pattern consistency**
   - [ ] All affected files use same pattern
   - [ ] No instances of old pattern remain
   - [ ] New pattern matches instruction files
   - [ ] Follows clean architecture layers

2. **Code quality**
   - [ ] DRY: Created reusable utilities if needed
   - [ ] No duplicated logic
   - [ ] Clear separation of concerns
   - [ ] Consistent naming throughout

3. **API contracts** (if changed)
   - [ ] Swagger regenerated (`make openapi`)
   - [ ] Response DTOs documented
   - [ ] Frontend updated to match
   - [ ] No breaking changes to external consumers

#### Portal-Web Validation:

1. **Pattern consistency**
   - [ ] All affected files use same pattern
   - [ ] No instances of old pattern remain
   - [ ] Matches instruction file patterns
   - [ ] Follows code structure guidelines

2. **Code quality**
   - [ ] DRY: Extracted reusable components/hooks
   - [ ] No duplicated markup
   - [ ] Proper component composition
   - [ ] Consistent naming

3. **File structure** (if changed)
   - [ ] Files in correct folders per code-structure instructions
   - [ ] Route files properly organized
   - [ ] Features properly isolated
   - [ ] No circular dependencies

#### Cross-Cutting Validation:

1. **Instruction file alignment**
   - [ ] Read all related instruction files
   - [ ] Verified pattern matches SSOT
   - [ ] No conflicts with other patterns
   - [ ] Harmonization is complete

2. **Migration completeness**
   - [ ] Searched for all instances (grep + semantic)
   - [ ] Updated all occurrences
   - [ ] Updated tests
   - [ ] Updated related docs

3. **No regressions**
   - [ ] Existing features work
   - [ ] Tests pass
   - [ ] No new errors introduced

### Phase 5.5: Update Instruction Files (Critical)

**Purpose:** Ensure this drift doesn't happen again by codifying the correct pattern.

#### Identify Relevant Instruction Files

1. **Check drift report**
   - Look at `related-instructions` field in frontmatter
   - These are the files that need updates

2. **Find additional relevant files**
   - If backend pattern: likely `backend-core.instructions.md` or domain-specific file
   - If frontend pattern: likely `portal-web-architecture.instructions.md` or component-specific file
   - If API contract: likely `responses-dtos-swagger.instructions.md`

#### Update Instructions

**For Option 1 (Updated code to match instructions):**

If the drift occurred because instructions were unclear or incomplete:

1. **Strengthen the existing rule**

   ```markdown
   # Before (weak/implicit)

   Use camelCase for API responses.

   # After (strong/explicit)

   **API Response Casing (CRITICAL):**

   - All JSON fields MUST be camelCase
   - Never use snake_case or PascalCase
   - Example: `createdAt`, `userId`, `businessDescriptor`
   - Common mistake: Exposing GORM models directly (they use PascalCase)
   - Solution: Create explicit response DTOs with `json:` tags
   ```

2. **Add anti-patterns section** (if not present)

   ```markdown
   ## Anti-Patterns (Avoid)

   - ❌ [Describe the drift pattern that just occurred]
   - ✅ [Describe the correct pattern applied]
   ```

3. **Add code examples** (if helpful)

   ```markdown
   ## Example: [Pattern Name]

   ❌ **Wrong (caused drift):**
   [Show the pattern that led to drift]

   ✅ **Correct:**
   [Show the harmonized pattern]
   ```

**For Option 2 (Updated instructions to match code):**

If you determined the code pattern is better:

1. **Update the instruction file** with the new pattern
2. **Document why** the change was made
3. **Add date** of the pattern change
4. **Verify no conflicts** with other instruction files

#### Specific Updates by Drift Type

**Naming Conventions Drift:**

- Update relevant instruction file's naming section
- Add explicit "Do/Don't" examples
- Clarify scope (when does this rule apply?)

**Architecture Pattern Drift:**

- Update architecture section with clear layer rules
- Add sequence diagram if complex
- Document when to use each pattern

**API Contract Drift:**

- Update responses-dtos-swagger.instructions.md
- Document the required DTO structure
- Add validation checklist

**Error Handling Drift:**

- Update errors-handling.instructions.md
- Document the standard error shape
- Add examples of common scenarios

#### Validation Checklist

Before proceeding:

- [ ] Identified all instruction files that need updates
- [ ] Updated instructions with explicit rules
- [ ] Added anti-pattern examples showing the drift that occurred
- [ ] Added correct pattern examples
- [ ] Verified updated instructions don't conflict with other SSOT files
- [ ] Instructions are now clear enough that AI/human won't repeat this drift

#### Example Instruction File Update

```markdown
# What to add to the instruction file:

## [Pattern Name] - Updated 2026-01-18

**Context:** A drift was found where [describe what happened].

**Rule (Explicit):**

- MUST: [Specific requirement]
- MUST NOT: [Specific anti-pattern]
- WHY: [Brief reasoning]

**Common Mistake:**
❌ [Show the exact mistake that led to drift]

**Correct Pattern:**
✅ [Show the correct pattern]

**Affected Areas:**

- [List where this pattern applies]

**Related Drift Reports:**

- `backlog/drifts/2026-01-18-[slug].md` - [Brief description]
```

### Phase 6: Update Drift Report

After successful harmonization:

1. **Update the drift report file**
   - Change status to `resolved`
   - Add resolution notes
   - Document approach taken
   - List files changed

2. **Add resolution section**

   ```markdown
   ## Resolution

   **Status:** Resolved
   **Date:** YYYY-MM-DD
   **Approach Taken:** [Option 1 or 2]

   **Harmonization Summary:**
   [Brief description of changes made]

   **Pattern Applied:**
   [Description of the consistent pattern now used]

   **Files Changed:**

   - path/to/file1.ext - [what was harmonized]
   - path/to/file2.ext - [what was harmonized]

   **Migration Completeness:**

   - Total instances found: [count]
   - Instances harmonized: [count]
   - Remaining drift: 0

   **Validation:**

   - [x] All tests pass
   - [x] Type check passes (if frontend)
   - [x] Lint passes (if frontend)
   - [x] Pattern applied consistently
   - [x] No regressions introduced
   - [x] Instruction files aligned

   **Instruction Files Updated:**

   - path/to/instruction.md - [what was added/clarified]
   - [Brief description of how instructions were strengthened]

   **Prevention:**
   This drift should not recur because instruction files now explicitly:

   - [What specific rule was added]
   - [What anti-pattern was documented]
   ```

### Phase 7: Final Summary

Provide a concise harmonization summary:

```markdown
## Drift Fix Implementation Summary

**Drift Report:** [path/to/report.md]
**Status:** ✅ Harmonized

### Harmonization Details

**Scope:** [Isolated/Localized/Systemic]
**Pattern Applied:** [Description]
**Approach:** [Option 1 or 2]

### Changes Made

- [List key changes]

### Files Modified

- [Count] files harmonized
- [Count] tests updated

### Migration Stats

- Old pattern instances: [count]
- All instances migrated: ✅
- Pattern now consistent: ✅

### Validation Results

**Backend:**

- [x] E2E tests pass (100%)
- [x] Pattern matches instruction files
- [x] Clean architecture maintained
- [x] No API breaking changes (or: documented breaking changes)

**Frontend:**

- [x] Type check passes
- [x] Lint passes
- [x] Pattern matches instruction files
- [x] Code structure correct
- [x] No regressions

### Verification

All drift instances have been harmonized. Pattern is now consistent across the codebase.

**Instruction Files Updated:**

- [List instruction files updated]
- [Brief description of clarifications added]

**Prevention Measures:**

- Explicit rules added to prevent recurrence
- Anti-pattern examples documented
- Correct pattern examples provided

Pattern is now codified in instruction files to prevent future drift.
```

## Quality Standards

### Must-Have

- **Complete**: All drift instances fixed, not just examples
- **Consistent**: Same pattern everywhere
- **Tested**: All tests pass after harmonization
- **Validated**: Matches instruction file patterns exactly
- **DRY**: Extracted reusable utilities for complex patterns
- **Documented**: Clear what changed and why
- **Instructions Updated**: Relevant instruction files strengthened to prevent recurrence

### Must-Not-Have

- **Partial migration**: Some files use old pattern, some new
- **New variations**: Creating yet another pattern variant
- **Test skips**: Commenting out tests that break
- **Shortcuts**: Ignoring some instances "because they work"
- **Undocumented changes**: Not explaining why pattern chosen

## Safety Checks

- Read the full drift report before starting
- Confirm the authoritative pattern from instruction files
- Search comprehensively for all instances
- Run tests after each major change
- Verify no functionality changes (except in breaking-change scenarios)
- Update relevant instruction files to prevent recurrence
- Update drift report with resolution

## Common Pitfalls

**Backend:**

- ❌ Forgetting to update DTOs after model changes
- ❌ Missing some instances in tests
- ❌ Not regenerating Swagger
- ❌ Breaking API contracts without updating frontend

**Portal-Web:**

- ❌ Moving files but not updating imports
- ❌ Extracting components but not using them everywhere
- ❌ Fixing types in some files but not others
- ❌ Harmonizing code but not updating tests

**Both:**

- ❌ Stopping at the first few instances
- ❌ Not searching for variations of the pattern
- ❌ Forgetting related documentation
- ❌ Not validating against instruction files
- ❌ **Fixing drift but not updating instruction files (drift WILL recur!)**
- ❌ **Adding vague instructions instead of explicit rules with examples**

## When to Update Instructions Instead

Choose Option 2 (update instructions) only when:

- Current code pattern is demonstrably superior
- Pattern is already widespread and working well
- Instruction was based on outdated assumptions
- New pattern better aligns with framework/library best practices

If updating instructions:

- Clearly document why in the instruction file
- Verify no conflicts with other SSOT files
- Get consensus if this is a major pattern change
- Update related instruction files

## References

- **Backend Core:** `.github/instructions/backend-core.instructions.md`
- **Backend Testing:** `.github/instructions/backend-testing.instructions.md`
- **Portal Architecture:** `.github/instructions/portal-web-architecture.instructions.md`
- **Portal Code Structure:** `.github/instructions/portal-web-code-structure.instructions.md`
- **Responses/DTOs:** `.github/instructions/responses-dtos-swagger.instructions.md`
- **Error Handling:** `.github/instructions/errors-handling.instructions.md`
- **State Management:** `.github/instructions/state-management.instructions.md`
- **All Instructions:** `.github/instructions/`
