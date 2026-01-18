---
name: implement-bug-fix
description: "Implement a bug fix from a bug report with comprehensive validation and testing"
agent: "Feature Builder"
argument-hint: "Path to bug report (e.g., backlog/bugs/2026-01-18-orders-date-filter-500-error.md)"
tools: ["vscode", "execute", "read", "edit", "search", "agent", "todo"]
---

# Bug Fix Implementation

You are about to implement a fix for the bug report:

**Bug Report:** `${input:bugReport:Path to bug report (e.g., backlog/bugs/2026-01-18-issue.md)}`

## Mission

Implement a complete, production-ready bug fix by:

1. Understanding the bug report thoroughly
2. Implementing the fix across all affected layers
3. Validating against instruction files and code patterns
4. Ensuring tests cover the fix
5. Verifying no regressions introduced

## Workflow

### Phase 1: Understanding & Planning

1. **Read the bug report completely**
   - Understand the problem, root cause, and proposed fix
   - Identify all affected files
   - Note related instruction files
   - Understand impact assessment

2. **Determine component(s) affected**
   - Backend only
   - Portal-web only
   - Storefront-web only
   - Full-stack (multiple components)
   - Infrastructure

3. **Read relevant instruction files**
   - All files listed in bug report's `related-instructions`
   - Component-specific patterns (backend-core, portal-web-architecture, etc.)
   - Cross-cutting concerns (errors, forms, HTTP, state)

4. **Create implementation plan**
   - List all files to modify
   - Identify new files to create (if any)
   - Determine test coverage needed
   - Plan validation steps

### Phase 2: Implementation

Implement the fix following the component-specific patterns:

#### For Backend Fixes:

1. **Follow clean architecture**
   - Handler layer: Request validation, response formatting
   - Service layer: Business logic, orchestration
   - Storage layer: Database operations, caching
   - Keep layers properly separated

2. **Use established patterns**
   - Use `request.ValidBody` for JSON validation
   - Use `response.Error` for error responses
   - Use `AtomicProcess.Exec` for transactions
   - Use repository scopes for queries
   - Use GORM models with proper tags

3. **Maintain security invariants**
   - Workspace/business scoping
   - RBAC checks via middleware
   - Input validation
   - No SQL injection vectors

#### For Portal-Web Fixes:

1. **Follow component structure**
   - Place files in correct locations per `portal-web-code-structure.instructions.md`
   - Use proper feature organization
   - Keep route components thin

2. **Use established patterns**
   - TanStack Router for routing
   - TanStack Query for data fetching
   - TanStack Form for forms
   - TanStack Store for client state
   - Use `apiClient` from `src/api/client.ts`

3. **Maintain UI/UX standards**
   - Mobile-first, RTL-first design
   - Use daisyUI components
   - Follow design tokens
   - Arabic translations required
   - Accessible markup

#### For Full-Stack Fixes:

1. **Backend first**
   - Implement and validate backend changes
   - Ensure API contract is correct

2. **Frontend second**
   - Consume new backend contract
   - Update types/schemas
   - Update UI components

### Phase 3: Testing

#### Backend Testing:

1. **Create/update E2E tests**
   - Location: `backend/internal/tests/e2e/`
   - Use testcontainers pattern
   - Test the bug scenario explicitly
   - Test happy path still works
   - Test edge cases

2. **Test structure**
   - Use suite-based pattern
   - Proper setup/teardown (truncate tables)
   - Use domain storage layers (no raw SQL)
   - Assert all response fields
   - Test RBAC boundaries

3. **Run tests**
   - Execute: `cd backend && go test ./internal/tests/e2e -v`
   - Verify 100% pass
   - Check test covers the bug fix

#### Portal-Web Testing:

1. **Type checking**
   - Execute: `cd portal-web && npm run type-check`
   - Fix all TypeScript errors
   - Ensure no type assertions bypass real issues

2. **Linting**
   - Execute: `cd portal-web && npm run lint -- --fix`
   - Fix all linting errors
   - Ensure code follows style guide

3. **Manual validation checklist**
   - [ ] Component renders without errors
   - [ ] Forms validate correctly
   - [ ] API calls use proper error handling
   - [ ] Loading states work
   - [ ] i18n keys exist and are used
   - [ ] RTL layout works
   - [ ] Mobile responsive

### Phase 4: Validation & Consistency Checks

#### Backend Validation:

1. **Instruction file alignment**
   - [ ] Follows `backend-core.instructions.md` patterns
   - [ ] Follows domain-specific instructions (orders, inventory, etc.)
   - [ ] Error handling uses ProblemDetails
   - [ ] Transactions use AtomicProcess
   - [ ] Repository uses scopes correctly

2. **Code quality**
   - [ ] DRY: Reused existing utilities/helpers
   - [ ] Created reusable helpers if needed
   - [ ] No duplicated business logic
   - [ ] Proper separation of concerns
   - [ ] Clear, descriptive names

3. **Security**
   - [ ] Workspace/business scoping enforced
   - [ ] RBAC checks in place
   - [ ] Input validation complete
   - [ ] No secrets in code

#### Portal-Web Validation:

1. **Instruction file alignment**
   - [ ] Follows `portal-web-architecture.instructions.md`
   - [ ] Follows `portal-web-code-structure.instructions.md`
   - [ ] Follows `forms.instructions.md` (if forms involved)
   - [ ] Follows `http-tanstack-query.instructions.md`
   - [ ] Follows `state-management.instructions.md`
   - [ ] Follows `ui-implementation.instructions.md`

2. **Code quality**
   - [ ] DRY: Reused existing components/hooks
   - [ ] Created reusable components if needed
   - [ ] No duplicated logic
   - [ ] Proper component composition
   - [ ] Clear, descriptive names

3. **UI/UX consistency**
   - [ ] Matches existing UI patterns
   - [ ] Uses design tokens
   - [ ] Mobile-first responsive
   - [ ] RTL-compatible
   - [ ] Arabic translations present

### Phase 5: Update Bug Report

After successful implementation:

1. **Update the bug report file**
   - Change status to `resolved`
   - Add resolution notes
   - List files changed
   - Link to implementation details

2. **Add resolution section**

   ```markdown
   ## Resolution

   **Status:** Resolved
   **Date:** YYYY-MM-DD
   **Implementation Summary:**
   [Brief description of what was changed]

   **Files Changed:**

   - path/to/file1.ext - [what changed]
   - path/to/file2.ext - [what changed]

   **Tests Added/Updated:**

   - path/to/test.ext - [coverage description]

   **Validation:**

   - [x] All tests pass
   - [x] Type check passes (if frontend)
   - [x] Lint passes (if frontend)
   - [x] Follows instruction files
   - [x] No regressions introduced
   ```

### Phase 6: Final Summary

Provide a concise implementation summary:

```markdown
## Bug Fix Implementation Summary

**Bug Report:** [path/to/report.md]
**Status:** ✅ Resolved

### Changes Made

- [List key changes]

### Files Modified

- [Count] files changed

### Tests

- Backend: [Pass/Fail + coverage description]
- Frontend: [Type check + Lint results]

### Validation Results

**Backend:**

- [x] E2E tests pass (100%)
- [x] Follows clean architecture
- [x] Security checks complete
- [x] DRY principles applied

**Frontend:**

- [x] Type check passes
- [x] Lint passes
- [x] UI/UX consistent
- [x] Mobile/RTL compatible
- [x] i18n complete

### Verification

The fix has been validated and is ready for review.
```

## Quality Standards

### Must-Have

- **Complete**: All files mentioned in bug report are fixed
- **Tested**: E2E tests (backend) or type-check + lint (frontend) pass
- **Validated**: Follows all relevant instruction files
- **DRY**: Reuses existing patterns, creates reusable code
- **Secure**: No security regressions
- **Consistent**: Matches existing code patterns

### Must-Not-Have

- **Partial fixes**: Don't fix only some affected files
- **Workarounds**: Don't add temporary hacks
- **Test skips**: Don't skip or comment out failing tests
- **Type assertions**: Don't use `any` or `@ts-ignore` to bypass issues
- **Duplicate code**: Don't copy-paste; refactor to reusable utilities

## Safety Checks

- Read the full bug report before starting
- Verify all file paths exist
- Run tests after implementation
- Check for regressions in related features
- Ensure no secrets or credentials added
- Verify instruction file alignment

## Common Pitfalls

**Backend:**

- ❌ Skipping transaction wrapping for multi-model updates
- ❌ Forgetting workspace/business scoping
- ❌ Not writing E2E tests
- ❌ Using raw SQL instead of repository scopes

**Portal-Web:**

- ❌ Skipping type checking
- ❌ Creating components in wrong folders
- ❌ Not adding i18n translations
- ❌ Forgetting RTL/mobile considerations
- ❌ Using fetch instead of apiClient

## References

- **Backend Core:** `.github/instructions/backend-core.instructions.md`
- **Backend Testing:** `.github/instructions/backend-testing.instructions.md`
- **Portal Architecture:** `.github/instructions/portal-web-architecture.instructions.md`
- **Portal Code Structure:** `.github/instructions/portal-web-code-structure.instructions.md`
- **Forms:** `.github/instructions/forms.instructions.md`
- **HTTP + Query:** `.github/instructions/http-tanstack-query.instructions.md`
- **State Management:** `.github/instructions/state-management.instructions.md`
- **UI Implementation:** `.github/instructions/ui-implementation.instructions.md`
