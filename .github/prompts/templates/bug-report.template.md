---
type: bug
date: YYYY-MM-DD
priority: critical|high|medium|low
component: backend|portal-web|storefront-web|infrastructure
affected-files:
  - path/to/file.ext
related-instructions:
  - .github/instructions/relevant.instructions.md
status: open
assignee: null
---

# [Short descriptive title]

## Summary

A clear, one-paragraph description of the bug.

## Current Behavior

What currently happens (be specific):
- Exact behavior observed
- Error messages (if any)
- Unexpected output

## Expected Behavior

What should happen according to:
- Instruction files
- API contracts
- Functional requirements
- Security policies

## Location

**Primary affected files:**
```
path/to/file.ext:123-145
path/to/another.ext:67
```

**Related files:**
- path/to/related.ext (explanation why)

## Root Cause

**Analysis:**
[Explanation of why this bug exists]

**Pattern violated:**
[Which instruction/convention is being broken]

## Impact

- **Severity:** [Critical/High/Medium/Low]
- **Affected users:** [Who is impacted]
- **Frequency:** [How often this occurs]
- **Workaround:** [Temporary mitigation, if any]

## Reproduction Steps

1. [Step 1]
2. [Step 2]
3. [Step 3]
4. Observe: [what happens]

**Environment:**
- OS: [if relevant]
- Browser: [if relevant]
- Backend version: [if relevant]

## Proposed Fix

**Changes required:**

1. **File:** `path/to/file.ext`
   ```[language]
   // Current code (simplified)
   [problematic code]
   
   // Proposed fix
   [corrected code]
   ```

2. **File:** `path/to/another.ext`
   [describe change]

**Testing checklist:**
- [ ] Unit tests updated/added
- [ ] E2E tests cover this scenario
- [ ] No regressions in [related feature]
- [ ] Follows [relevant instruction file]

## References

- Instruction file: [link to .github/instructions/...]
- Related SSOT: [links]
- API contract: [if applicable]
- Original implementation: [commit/PR if known]

## Dependencies

- Blocked by: [other tickets, if any]
- Blocks: [other work that depends on this fix]
- Related to: [similar issues]

## Notes

[Any additional context, edge cases, or considerations]
