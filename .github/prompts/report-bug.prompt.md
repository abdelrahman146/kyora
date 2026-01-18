---
name: report-bug
description: "Analyze user-reported bug description and create a comprehensive bug report with proper investigation"
agent: "AI Architect"
argument-hint: "Brief bug description (e.g., 'Orders list returns 500 error when filtering by date')"
tools: ["vscode", "execute", "read", "edit", "search", "todo"]
---

# Bug Report Generator

You are about to analyze and document a bug reported by the user:

**Bug Description:** `${input:bugDescription:Describe the bug you encountered}`

## Mission

Create a comprehensive, actionable bug report by:

1. Understanding the user's description
2. Investigating the codebase to find affected areas
3. Analyzing the root cause
4. Proposing concrete fixes
5. Checking for duplicates in existing reports

## Workflow

### Phase 1: Understanding & Scope Discovery

1. **Parse the user's description**
   - Identify affected component (backend/portal-web/storefront-web/infrastructure)
   - Extract symptoms (error messages, unexpected behavior)
   - Determine user-facing vs internal issue

2. **Check for existing reports**
   - Search `backlog/bugs/` for similar issues
   - If duplicate found: enhance the existing report instead of creating new
   - If related: reference in the new report

3. **Locate affected code**
   - Use semantic search for relevant files
   - Use grep for error messages/stack traces
   - Identify the likely problematic code sections

### Phase 2: Deep Investigation

1. **Reproduce the issue mentally**
   - Trace the code path from user action to error
   - Identify preconditions and triggers
   - Determine if it's environment-specific

2. **Find related code**
   - Check related handlers/services/components
   - Look for similar patterns that might have the same bug
   - Review tests (or lack thereof) for this functionality

3. **Identify root cause**
   - What specific code/logic causes this?
   - Why was it implemented this way?
   - Is it violating any instruction file patterns?

### Phase 3: Impact Assessment

Determine:

- **Severity**: Does this break core functionality? Data integrity risk?
- **Frequency**: Always, sometimes, edge case?
- **Affected users**: Who encounters this? How often?
- **Workaround**: Is there a temporary mitigation?

Priority guidelines:

- **Critical**: Data loss, security issue, core feature broken for all users
- **High**: Important feature broken, affects many users, no workaround
- **Medium**: Feature partially broken, affects some users, workaround exists
- **Low**: Edge case, cosmetic, minor inconvenience

### Phase 4: Solution Design

1. **Propose fix**
   - Provide specific code changes needed
   - Show before/after snippets
   - Consider edge cases

2. **Identify related files**
   - What else needs updating? (tests, docs, related code)
   - Are there similar bugs elsewhere?

3. **Plan testing**
   - How to verify the fix?
   - What tests need adding/updating?
   - Regression risk assessment

### Phase 5: Report Generation

Create bug report at: `backlog/bugs/YYYY-MM-DD-<slug>.md`

Use the bug report template with:

- **Specific locations**: File paths + line numbers
- **Concrete examples**: Actual code snippets
- **Clear reproduction**: Step-by-step
- **Actionable fix**: Detailed implementation guidance
- **References**: Links to related SSOT files

**Slug naming**: Use kebab-case, max 50 chars, descriptive
Examples: `orders-date-filter-500-error`, `customer-search-null-pointer`

### Phase 6: Verification

Before finalizing:

- [ ] Checked for duplicates in `backlog/bugs/`
- [ ] Verified all file paths exist
- [ ] Included exact line numbers
- [ ] Linked to relevant instruction files
- [ ] Proposed fix is concrete and actionable
- [ ] Priority accurately reflects impact
- [ ] Frontmatter is valid YAML

## Output Format

Provide:

1. **The generated bug report** (full path)
2. **Brief summary** of findings:

   ```
   Bug Report Created: backlog/bugs/2026-01-18-[slug].md

   Priority: [critical|high|medium|low]
   Component: [component]
   Root Cause: [one-line explanation]
   Affected Files: [count]
   Fix Effort: [estimate]

   Related Reports: [if any]
   ```

## Safety Checks

- Do not modify production code (only create report in `backlog/bugs/`)
- Verify all file references are valid
- Ensure frontmatter is complete and valid
- Create `backlog/bugs/` directory if needed
- If unsure about root cause, state assumptions clearly

## Investigation Tips

**For backend bugs:**

- Check `backend/internal/domain/` for business logic
- Review middleware chains in `backend/internal/server/routes.go`
- Look for error handling patterns in `backend/internal/platform/response/`
- Check E2E tests in `backend/internal/tests/e2e/`

**For portal-web bugs:**

- Check route files in `portal-web/src/routes/`
- Review API clients in `portal-web/src/api/`
- Look for state management in stores/hooks
- Check form validation in schemas

**Common bug patterns:**

- Missing null checks
- Incorrect error handling
- Race conditions in async code
- Missing RBAC checks
- Cross-tenant data leaks
- Unhandled edge cases

## References

- Bug Template: `.github/prompts/templates/bug-report.template.md`
- AI Infrastructure: `.github/instructions/ai-infrastructure.instructions.md`
- Backend Patterns: `.github/instructions/backend-core.instructions.md`
- Portal Patterns: `.github/instructions/portal-web-architecture.instructions.md`
