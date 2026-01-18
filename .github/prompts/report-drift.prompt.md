---
name: report-drift
description: "Analyze pattern inconsistency and create a comprehensive drift report with harmonization plan"
agent: "AI Architect"
argument-hint: "Brief drift description (e.g., 'Customer API uses snake_case but instruction says camelCase')"
tools: ["vscode", "execute", "read", "edit", "search", "todo"]
---

# Drift Report Generator

You are about to analyze and document a pattern drift reported by the user:

**Drift Description:** `${input:driftDescription:Describe the pattern inconsistency}`

## Mission

Create a comprehensive drift report by:

1. Understanding the inconsistency described
2. Investigating the extent of the drift
3. Identifying the authoritative pattern
4. Proposing harmonization strategy
5. Checking for duplicates in existing reports

## Workflow

### Phase 1: Understanding & Pattern Discovery

1. **Parse the user's description**
   - Identify the pattern category (naming, structure, API, state, etc.)
   - Extract what exists vs. what's expected
   - Determine affected component(s)

2. **Check for existing reports**
   - Search `backlog/drifts/` for similar drift reports
   - If duplicate: enhance existing report instead of creating new
   - If related: reference in the new report

3. **Find the authoritative pattern**
   - Search instruction files (`.github/instructions/*.instructions.md`)
   - Look for SSOT definitions in `copilot-instructions.md`
   - Check if the "expected" pattern is actually documented

### Phase 2: Scope Analysis

1. **Map the extent of drift**
   - Use grep/semantic search to find all instances
   - Count affected files and lines
   - Identify if it's isolated, localized, or systemic

2. **Categorize the drift**
   - [ ] Naming convention (camelCase vs snake_case, etc.)
   - [ ] File/folder structure
   - [ ] API contract inconsistency
   - [ ] State management pattern
   - [ ] Error handling approach
   - [ ] Testing pattern
   - [ ] Import/module organization
   - [ ] Other: [specify]

3. **Analyze pattern distribution**
   - How many files follow the "correct" pattern?
   - How many files follow the "drift" pattern?
   - Are there multiple competing patterns?

### Phase 3: Root Cause Investigation

Determine why the drift exists:

- [ ] Pattern was introduced after initial code
- [ ] Instruction file didn't exist when code was written
- [ ] Developer/AI was unaware of the pattern
- [ ] Pattern changed mid-project
- [ ] Intentional deviation (find reasoning)
- [ ] Legacy code not yet refactored
- [ ] Copy-paste from external source

### Phase 4: Impact Assessment

Evaluate consequences:

- **Consistency**: How does this hurt uniformity?
- **Maintainability**: Does this make code harder to change?
- **Onboarding**: Does this confuse new developers/AI?
- **Interoperability**: Does this break integrations?

Priority guidelines:

- **High**: Systemic drift affecting multiple domains, blocks new work
- **Medium**: Localized drift, causes confusion, should fix soon
- **Low**: Isolated incident, minimal impact, fix eventually

### Phase 5: Harmonization Strategy

Design the fix:

**Option 1: Update code to match instructions** (usual choice)

- Identify all locations needing change
- Estimate effort (hours/days)
- Assess risk (breaking changes?)
- Plan migration steps

**Option 2: Update instructions to match code** (rare, when code is better)

- Justify why the "drift" pattern is actually superior
- Check for conflicts with other SSOT files
- Propose instruction file update

**Recommendation**: Choose the option that:

- Aligns with broader Kyora patterns
- Minimizes breaking changes
- Improves consistency
- Is easier to maintain long-term

### Phase 6: Report Generation

Create drift report at: `backlog/drifts/YYYY-MM-DD-<slug>.md`

Use the drift report template with:

- **Specific locations**: All affected files + line ranges
- **Pattern comparison**: Clear before/after examples
- **Scope assessment**: Isolated/localized/systemic
- **Harmonization plan**: Concrete steps with effort estimate
- **References**: Links to instruction files

**Slug naming**: Use kebab-case, max 50 chars, descriptive
Examples: `customer-api-snake-case-drift`, `order-state-pattern-mismatch`

### Phase 7: Verification

Before finalizing:

- [ ] Checked for duplicates in `backlog/drifts/`
- [ ] Verified instruction file references
- [ ] Found all drift instances (not just examples)
- [ ] Proposed harmonization is concrete
- [ ] Effort estimate is realistic
- [ ] Priority reflects actual impact
- [ ] Frontmatter is valid YAML
- [ ] Pattern category is set correctly

## Output Format

Provide:

1. **The generated drift report** (full path)
2. **Brief summary** of findings:

   ```
   Drift Report Created: backlog/drifts/2026-01-18-[slug].md

   Priority: [high|medium|low]
   Component: [component]
   Category: [pattern-category]
   Scope: [isolated|localized|systemic]
   Files Affected: [count]
   Harmonization Effort: [estimate]
   Recommended Approach: [Option 1 or 2]

   Related Reports: [if any]
   ```

## Safety Checks

- Do not modify production code (only create report in `backlog/drifts/`)
- Verify all file references are valid
- Ensure frontmatter is complete and valid
- Create `backlog/drifts/` directory if needed
- If proposing instruction changes, clearly justify

## Investigation Tips

**Pattern sources to check:**

- Backend: `.github/instructions/backend-core.instructions.md`
- Portal: `.github/instructions/portal-web-architecture.instructions.md`
- Forms: `.github/instructions/forms.instructions.md`
- HTTP: `.github/instructions/ky.instructions.md`
- State: `.github/instructions/state-management.instructions.md`
- i18n: `.github/instructions/i18n-translations.instructions.md`

**Common drift patterns:**

- **Naming**: camelCase vs snake_case, PascalCase inconsistency
- **Imports**: relative vs absolute paths
- **Errors**: different error response shapes
- **State**: multiple ways to manage same type of state
- **Files**: components in wrong folders
- **API**: different request/response formats

**Scope indicators:**

- Isolated: 1-2 files, one feature
- Localized: One domain/module (5-15 files)
- Systemic: Multiple domains (15+ files)

## References

- Drift Template: `.github/prompts/templates/drift-report.template.md`
- AI Infrastructure: `.github/instructions/ai-infrastructure.instructions.md`
- All Instruction Files: `.github/instructions/`
