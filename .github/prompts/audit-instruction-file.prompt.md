---
name: audit-instruction-file
description: "Audit and update instruction file, detect bugs/drifts, create actionable reports"
agent: "AI Architect"
argument-hint: "Path to the instruction file to audit (e.g., .github/instructions/backend-core.instructions.md)"
tools: ["read", "edit", "search"]
---

# Instruction File Audit & Update

You are about to perform a comprehensive audit and update of the instruction file:

**Target File:** `${input:instructionFile:Path to instruction file (e.g., .github/instructions/backend-core.instructions.md)}`

## Mission

Perform a deep audit of the specified instruction file to ensure it is:

1. **Accurate** — matches current repository implementation
2. **Complete** — covers all relevant patterns and conventions
3. **Consistent** — aligns with Kyora SSOT hierarchy and doesn't conflict with other instructions
4. **Actionable** — provides clear, testable rules
5. **Up-to-date** — reflects recent changes and patterns

## Workflow

### Phase 1: Discovery & Context Gathering

1. **Read the target instruction file** completely
2. **Identify the scope** from the `applyTo` frontmatter
3. **Map related files** using the scope pattern (search for actual code files)
4. **Read related SSOT files** mentioned in `.github/copilot-instructions.md`
5. **Gather recent changes** (search for patterns mentioned in the instruction file)

### Phase 2: Validation & Gap Analysis

For each rule/pattern in the instruction file:

1. **Verify it exists** in the codebase
   - Use grep/semantic search to find real examples
   - If a pattern is prescribed but not used, note it for removal or mark as "planned"

2. **Check for completeness**
   - Are there conventions used in code but not documented?
   - Are there new libraries/patterns added since the instruction was last updated?

3. **Test for consistency**
   - Does this rule conflict with other instruction files?
   - Does it align with `.github/copilot-instructions.md`?

4. **Assess clarity**
   - Is the rule testable and actionable?
   - Are examples provided when needed?

### Phase 3: Drift & Bug Detection

Search the codebase for:

1. **Known drifts** — code that doesn't follow the instruction file patterns
2. **Inconsistencies** — multiple competing patterns for the same problem
3. **Anti-patterns** — code that violates stated rules
4. **Missing patterns** — conventions used but not documented
5. **Outdated references** — mentions of deprecated libraries, removed files, or changed APIs

For each issue found, collect:

- **Exact location** (file path + line numbers)
- **Current state** (what the code does now)
- **Expected state** (what it should do per instructions/best practice)
- **Impact** (how this affects consistency/maintainability/correctness)
- **Suggested fix** (concrete steps to resolve)

### Phase 4: Check Existing Reports

Before creating new reports:

1. **Read existing open reports** in `backlog/bugs/`, `backlog/drifts/`, `backlog/enhancements/`
2. **For each new issue found**, check if:
   - A similar issue already exists (same file/pattern/problem)
   - A related issue exists (same domain but different aspect)
3. **Deduplication strategy**:
   - If exact duplicate: skip creating new report, note in summary
   - If partial overlap: enhance existing report with new findings instead of creating new
   - If related but distinct: create new report and reference the related one
4. **Enhancement strategy**:
   - Add newly discovered locations to existing reports
   - Update impact assessment if scope is larger
   - Add additional context or fix suggestions

### Phase 5: Report Generation

Create detailed reports in `backlog/` following this structure:

#### For Bugs (correctness issues):

- File: `backlog/bugs/YYYY-MM-DD-<short-slug>.md`
- Template: Use the bug report template (`.github/prompts/templates/bug-report.template.md`)
- Criteria: Behavior that breaks functionality or causes errors

#### For Drifts (consistency issues):

- File: `backlog/drifts/YYYY-MM-DD-<short-slug>.md`
- Template: Use the drift report template (`.github/prompts/templates/drift-report.template.md`)
- Criteria: Code that works but doesn't follow established patterns/instructions

#### For Enhancements (missing features/patterns):

- File: `backlog/enhancements/YYYY-MM-DD-<short-slug>.md`
- Template: Use the enhancement template (`.github/prompts/templates/enhancement.template.md`)
- Criteria: Gaps in coverage, missing documentation, or improvement opportunities

### Phase 6: Update the Instruction File

Based on your findings:

1. **Add missing patterns** found in Phase 2
2. **Remove invalid patterns** (those not actually used)
3. **Update outdated patterns** to match current implementation
4. **Clarify ambiguous rules** with examples from actual code
5. **Add references** to related SSOT files (link, don't duplicate)
6. **Update frontmatter** if scope has changed
7. **Fix any "Known Drifts" sections** — if drifts are resolved, mark them as resolved; if new drifts exist, document them briefly

### Phase 7: Provide Summary

Provide a concise summary in the chat (not as a separate file):

```markdown
## Audit Summary for [instruction-file-name]

**Date:** [current-date]
**Files Scanned:** [count]
**Patterns Validated:** [count]

### Changes Made to Instruction File

- [list each specific change made to the instruction file]
- [be specific: "Added pattern X", "Removed outdated Y", "Updated Z to match implementation"]

### Reports Created

- Bugs: [count] (see backlog/bugs/)
  - New: [count] | Enhanced: [count]
- Drifts: [count] (see backlog/drifts/)
  - New: [count] | Enhanced: [count]
- Enhancements: [count] (see backlog/enhancements/)
  - New: [count] | Enhanced: [count]
- Duplicates Skipped: [count]

### Top Priority Issues

1. [most critical bug/drift]
2. [second most critical]
3. [third most critical]

### Recommendations

- [actionable next steps]
```

## Quality Standards

### For Instruction File Updates

- Every rule must have a real example in the codebase (or be marked as "planned")
- No contradictions with other SSOT files
- Clear `applyTo` scope
- Links to related files instead of duplicating content
- Concrete, testable guidance (avoid "use best practices")

### For Bug/Drift/Enhancement Reports

- **Specific**: Exact file paths and line numbers
- **Reproducible**: Clear steps to see the issue
- **Actionable**: Concrete fix proposal
- **Contextual**: Links to relevant SSOT files
- **Prioritized**: Impact assessment (critical/high/medium/low)

## Safety Checks

- **Only modify instruction files** (`.github/instructions/*.instructions.md`)
- **Create bug/drift/enhancement reports** in `backlog/` subdirectories
- **Do NOT create**: audit reports, README files, documentation files, or summary markdown files
- Do not modify production code
- Do not add secrets or credentials to any file
- Ensure all created markdown files have valid frontmatter
- Link to existing files; verify paths before adding references
- Create `backlog/` folder structure if it doesn't exist
- Read existing reports before creating new ones to avoid duplicates

## Output Requirements

- **Update the instruction file** to match current codebase reality
- **Create actionable bug/drift/enhancement reports** for issues found
- **Provide summary in chat** (not as a separate markdown file)
- **Do NOT create**: audit reports (e.g., `AUDIT-YYYY-MM-DD-*.md`), README files, or general documentation

## References

- AI Infrastructure SSOT: `.github/instructions/ai-infrastructure.instructions.md`
- Repo Orchestration: `.github/copilot-instructions.md`
- Related Instructions: See `.github/instructions/` for domain-specific rules
