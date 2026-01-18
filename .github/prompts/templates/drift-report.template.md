---
type: drift
date: YYYY-MM-DD
priority: high|medium|low
component: backend|portal-web|storefront-web|infrastructure
affected-files:
  - path/to/file.ext
related-instructions:
  - .github/instructions/relevant.instructions.md
status: open
assignee: null
pattern-category: naming|structure|error-handling|state-management|api-contract|testing|other
---

# [Short descriptive title]

## Summary

A clear, one-paragraph description of the drift from established patterns.

## Current State

**What exists now:**
- [Description of current implementation]
- [How it differs from the standard pattern]

**Affected locations:**
```
path/to/file1.ext:123-145
path/to/file2.ext:67-89
path/to/file3.ext:234
```

## Expected Pattern

**Per instruction file** (`.github/instructions/[relevant].instructions.md`):

[Quote or reference the specific rule being violated]

**Standard implementation:**
```[language]
// How this should be implemented
[expected code pattern]
```

## Pattern Deviation Analysis

**Type of drift:**
- [ ] Naming convention violation
- [ ] File structure deviation
- [ ] API contract inconsistency
- [ ] State management pattern mismatch
- [ ] Error handling divergence
- [ ] Testing pattern deviation
- [ ] Other: [specify]

**Why this matters:**
- **Consistency impact:** [how this affects codebase uniformity]
- **Maintainability impact:** [how this affects future changes]
- **Onboarding impact:** [how this confuses new developers/AI agents]

## Scope

**Extent of drift:**
- Files affected: [count]
- Lines of code: [approximate]
- Domains affected: [list]

**Is this drift widespread?**
- [ ] Isolated incident (1-2 files)
- [ ] Localized pattern (one feature/domain)
- [ ] Systemic issue (multiple domains)

## Root Cause

**Why did this drift occur?**
- [ ] Instruction file didn't exist when code was written
- [ ] Pattern changed after initial implementation
- [ ] Developer/AI was unaware of the pattern
- [ ] Intentional deviation (document reason)
- [ ] Legacy code not yet refactored
- [ ] Other: [specify]

## Harmonization Plan

**Option 1: Update code to match instructions**

1. **File:** `path/to/file.ext`
   ```[language]
   // Current (drift)
   [current code]
   
   // Aligned pattern
   [corrected code]
   ```

2. **File:** `path/to/another.ext`
   [describe change]

**Effort estimate:** [hours/days]
**Risk level:** [low/medium/high]

**Option 2: Update instructions to match code** (if current code is actually better)

[Justify why the instruction should change]
[Proposed instruction file update]

**Recommended approach:** [Option 1 or 2, with reasoning]

## Impact Assessment

- **Priority:** [High/Medium/Low]
- **Urgency:** [Immediate/Soon/Eventually]
- **Breaking changes:** [Yes/No - explain]
- **Testing effort:** [Small/Medium/Large]

## Migration Checklist

- [ ] Update all affected files
- [ ] Run existing tests
- [ ] Add/update tests for the corrected pattern
- [ ] Update related documentation
- [ ] Verify no new drifts introduced
- [ ] Update instruction file (if needed)

## References

- **Instruction file:** [link to .github/instructions/...]
- **Pattern examples:** [links to good examples in codebase]
- **Related drifts:** [links to similar drift reports]
- **SSOT hierarchy:** [relevant section in copilot-instructions.md]

## Decision Log

**Date:** [YYYY-MM-DD]
**Decision:** [Keep as-is / Refactor code / Update instructions]
**Rationale:** [reasoning for decision]
**Decided by:** [person/role]

## Notes

[Any additional context, trade-offs, or historical information]
