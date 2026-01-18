---
name: implement-enhancement
description: "Implement an enhancement proposal with comprehensive validation and quality checks"
agent: "Feature Builder"
argument-hint: "Path to enhancement report (e.g., backlog/enhancements/2026-01-18-phone-validation-helper.md)"
tools: ["vscode", "execute", "read", "edit", "search", "agent", "todo"]
---

# Enhancement Implementation

You are about to implement the enhancement described in:

**Enhancement Report:** `${input:enhancementReport:Path to enhancement report (e.g., backlog/enhancements/2026-01-18-issue.md)}`

## Mission

Implement a complete, production-ready enhancement by:

1. Understanding the enhancement proposal thoroughly
2. Implementing the solution across all required layers
3. Validating against instruction files and patterns
4. Ensuring comprehensive testing
5. Verifying success criteria are met

## Workflow

### Phase 1: Understanding & Planning

1. **Read the enhancement report completely**
   - Understand problem statement and proposed solution
   - Review implementation plan
   - Note success criteria
   - Review alternatives considered

2. **Determine enhancement category**
   - Documentation (instruction file updates)
   - Pattern (new utilities/helpers/components)
   - Tooling (scripts/build improvements)
   - Testing (test utilities/coverage)
   - Performance (optimizations)
   - Developer Experience (DX improvements)

3. **Read relevant instruction files**
   - All files listed in enhancement's `related-instructions`
   - Component-specific patterns
   - Related domain instructions

4. **Create detailed implementation plan**
   - Break into concrete steps
   - Identify files to create
   - Identify files to modify
   - Plan test coverage
   - Plan validation approach

### Phase 2: Implementation

#### For Documentation Enhancements:

1. **Instruction file updates**
   - Add new sections with clear, actionable rules
   - Provide concrete examples
   - Link to related SSOT files
   - Update frontmatter if scope changes

2. **Quality checks**
   - [ ] Rule is testable and specific
   - [ ] Examples show real code patterns
   - [ ] No conflicts with other instructions
   - [ ] Proper markdown formatting
   - [ ] Valid YAML frontmatter

#### For Pattern Enhancements (Backend):

1. **Create reusable utilities**
   - Location: `backend/internal/platform/utils/` or domain-specific
   - Clear, descriptive function names
   - Proper error handling
   - Type-safe interfaces
   - Documented with godoc comments

2. **Integration patterns**
   - Service helpers for business logic
   - Repository scopes for queries
   - Middleware for cross-cutting concerns
   - Request/response utilities

3. **Example usage**
   - Add examples to instruction files
   - Update related documentation
   - Consider adding to test fixtures

#### For Pattern Enhancements (Portal-Web):

1. **Create reusable components/hooks**
   - Location per `portal-web-code-structure.instructions.md`
   - Components: `src/components/` or feature-specific
   - Hooks: `src/hooks/` or feature-specific
   - Utilities: `src/lib/`

2. **Component patterns**
   - Proper TypeScript types
   - Mobile-first, RTL-compatible
   - Use design tokens
   - daisyUI integration
   - Accessible markup

3. **Hook patterns**
   - Clear responsibility
   - Proper dependency arrays
   - Error handling
   - Loading states
   - Type-safe returns

4. **Integration**
   - Update existing code to use new pattern
   - Add examples to instruction files
   - Ensure consistency

#### For Tooling Enhancements:

1. **Create scripts**
   - Location: `scripts/`
   - Clear, documented purpose
   - Error handling
   - Input validation
   - Usage examples

2. **Build/development improvements**
   - Update `Makefile` or `package.json`
   - Document new commands
   - Ensure cross-platform compatibility
   - Test on clean environment

#### For Testing Enhancements:

1. **Test utilities** (Backend)
   - Location: `backend/internal/tests/testutils/`
   - Fixtures for common scenarios
   - Helper functions for setup/teardown
   - Reusable assertion patterns

2. **Test utilities** (Portal-Web)
   - Location: `portal-web/src/test/` (if not exists, create)
   - Testing library utilities
   - Mock data factories
   - Reusable test helpers

3. **Coverage improvements**
   - Add missing test scenarios
   - Create test patterns for common cases
   - Document testing approach

### Phase 3: Testing & Validation

#### Backend Implementation Testing:

1. **Create E2E tests for new functionality**
   - Test happy paths
   - Test error cases
   - Test edge cases
   - Test integration with existing features

2. **Run all tests**
   - Execute: `cd backend && go test ./... -v`
   - Verify 100% pass
   - Check new tests are actually running

3. **Validate patterns**
   - [ ] Follows clean architecture
   - [ ] Uses repository patterns
   - [ ] Proper error handling
   - [ ] Transaction safety
   - [ ] RBAC where applicable

#### Portal-Web Implementation Testing:

1. **Type checking**
   - Execute: `cd portal-web && npm run type-check`
   - Fix all TypeScript errors
   - Ensure new code is properly typed

2. **Linting**
   - Execute: `cd portal-web && npm run lint -- --fix`
   - Fix all linting errors
   - Ensure code style consistency

3. **Component testing** (if applicable)
   - Component renders correctly
   - Props validation
   - Event handlers work
   - Accessibility checks
   - RTL compatibility

4. **Integration testing**
   - Works with existing features
   - API integration correct
   - State management proper
   - No console errors

#### Documentation Testing:

1. **If instruction file updated**
   - [ ] All patterns described exist in code
   - [ ] Examples are accurate
   - [ ] No conflicts with other instructions
   - [ ] Frontmatter valid
   - [ ] Markdown well-formed

2. **If new pattern added**
   - [ ] Pattern documented in relevant instruction file
   - [ ] Examples provided
   - [ ] Usage guidance clear

### Phase 4: DRY & Reusability Validation

1. **Check for duplication**
   - Search for similar patterns in codebase
   - Refactor existing code to use new pattern
   - Ensure new pattern is reusable

2. **Validate reusability**
   - [ ] Pattern is generic enough
   - [ ] Not over-engineered
   - [ ] Clear API/interface
   - [ ] Well-documented
   - [ ] Easy to use correctly

3. **Update existing code** (if applicable)
   - Find locations that could use new pattern
   - Refactor to use new utility/component
   - Ensure consistency

### Phase 5: Success Criteria Validation

1. **Review success criteria from enhancement report**
   - Check each criterion
   - Verify measurable outcomes
   - Document evidence

2. **Validate against metrics** (if defined)
   - Measure before/after
   - Document improvements
   - Verify targets met

3. **Functional validation**
   - [ ] Problem is solved
   - [ ] No workarounds needed
   - [ ] User experience improved
   - [ ] Developer experience improved

### Phase 6: Comprehensive Validation

#### Backend Validation:

1. **Instruction file alignment**
   - [ ] Follows `backend-core.instructions.md`
   - [ ] Follows domain-specific instructions
   - [ ] Matches established patterns
   - [ ] Proper separation of concerns

2. **Code quality**
   - [ ] DRY: Reusable implementation
   - [ ] Clear, descriptive names
   - [ ] Proper error handling
   - [ ] Well-documented (godoc)
   - [ ] Type-safe interfaces

3. **Testing**
   - [ ] E2E tests pass 100%
   - [ ] New tests added/updated
   - [ ] Coverage maintained/improved
   - [ ] Tests lock in new behavior

#### Portal-Web Validation:

1. **Instruction file alignment**
   - [ ] Follows `portal-web-architecture.instructions.md`
   - [ ] Follows `portal-web-code-structure.instructions.md`
   - [ ] Follows relevant pattern instructions
   - [ ] Proper file placement

2. **Code quality**
   - [ ] DRY: Reusable components/hooks
   - [ ] Clear, descriptive names
   - [ ] Proper TypeScript types
   - [ ] Well-documented (JSDoc/comments)
   - [ ] Accessible markup

3. **UI/UX** (if applicable)
   - [ ] Mobile-first responsive
   - [ ] RTL-compatible
   - [ ] Uses design tokens
   - [ ] Follows UI patterns
   - [ ] Arabic translations added

4. **Testing**
   - [ ] Type check passes
   - [ ] Lint passes
   - [ ] Component tests (if applicable)
   - [ ] Integration validated

### Phase 7: Update Enhancement Report

After successful implementation:

1. **Update the enhancement report file**
   - Change status to `implemented`
   - Add implementation notes
   - List files created/modified
   - Document success criteria met

2. **Add implementation section**

   ````markdown
   ## Implementation

   **Status:** Implemented
   **Date:** YYYY-MM-DD
   **Implementation Summary:**
   [Brief description of what was built]

   **Files Created:**

   - path/to/new/file.ext - [purpose]

   **Files Modified:**

   - path/to/existing/file.ext - [what changed]

   **Usage Example:**

   ```[language]
   // How to use the new pattern
   [code example]
   ```
   ````

   **Documentation:**
   - Instruction file updated: [path if applicable]
   - Examples added: [where]

   **Success Criteria Met:**
   - [x] [Criterion 1]
   - [x] [Criterion 2]
   - [x] [Criterion 3]

   **Validation:**
   - [x] All tests pass
   - [x] Type check passes (if frontend)
   - [x] Lint passes (if frontend)
   - [x] DRY principles applied
   - [x] Follows instruction files
   - [x] Reusable implementation

   ```

   ```

### Phase 8: Final Summary

Provide a concise implementation summary:

```markdown
## Enhancement Implementation Summary

**Enhancement Report:** [path/to/report.md]
**Status:** ✅ Implemented

### What Was Built

**Category:** [documentation|pattern|tooling|testing|performance|dx]
**Solution:** [Brief description]

### Implementation Details

- [List key components/changes]

### Files Created/Modified

- Created: [count] new files
- Modified: [count] existing files

### Reusability

- [Description of how pattern is reusable]
- [Existing code refactored to use it: count]

### Success Criteria

- [x] [Each criterion with checkmark]

### Validation Results

**Backend:**

- [x] E2E tests pass (100%)
- [x] Follows clean architecture
- [x] DRY and reusable
- [x] Proper separation of concerns
- [x] Well-documented

**Frontend:**

- [x] Type check passes
- [x] Lint passes
- [x] UI/UX consistent
- [x] Mobile/RTL compatible
- [x] DRY and reusable
- [x] Well-documented

### Impact

[Description of improvement/benefit delivered]

### Verification

Enhancement is complete, tested, and ready for use.
```

## Quality Standards

### Must-Have

- **Complete**: All aspects of enhancement implemented
- **Reusable**: Pattern can be used elsewhere
- **Tested**: Comprehensive test coverage
- **Documented**: Clear usage guidance
- **DRY**: No duplication introduced
- **Consistent**: Matches existing patterns
- **Validated**: Success criteria met

### Must-Not-Have

- **Partial implementation**: Some parts missing
- **Over-engineering**: More complex than needed
- **Duplication**: Similar code exists elsewhere
- **Poor naming**: Unclear purpose
- **Undocumented**: No usage guidance
- **Inconsistent**: Doesn't match codebase patterns

## Safety Checks

- Read the full enhancement report before starting
- Follow implementation plan from report
- Create reusable, not one-off solutions
- Run all tests after implementation
- Verify success criteria met
- Update instruction files if new pattern
- Ensure no regressions introduced

## Common Pitfalls

**Backend:**

- ❌ Creating domain-specific utility in platform layer
- ❌ Not adding E2E tests
- ❌ Forgetting error handling
- ❌ Not documenting with godoc

**Portal-Web:**

- ❌ Creating non-reusable "utility" components
- ❌ Placing files in wrong locations
- ❌ Not adding TypeScript types
- ❌ Forgetting i18n translations
- ❌ Not considering mobile/RTL

**Both:**

- ❌ Not checking for existing similar code
- ❌ Over-engineering simple solutions
- ❌ Not updating instruction files
- ❌ Skipping documentation

## Enhancement-Specific Guidance

### Documentation Enhancements

- Add concrete, testable rules
- Provide real code examples
- Link to related SSOT files
- Keep instructions actionable

### Pattern Enhancements

- Make it reusable from day one
- Update existing code to use it
- Add to relevant instruction file
- Provide clear usage examples

### Tooling Enhancements

- Test on clean environment
- Document all commands/options
- Handle errors gracefully
- Make cross-platform compatible

### Testing Enhancements

- Make test utilities truly reusable
- Document testing patterns
- Provide clear examples
- Ensure tests are maintainable

## References

- **Backend Core:** `.github/instructions/backend-core.instructions.md`
- **Backend Testing:** `.github/instructions/backend-testing.instructions.md`
- **Portal Architecture:** `.github/instructions/portal-web-architecture.instructions.md`
- **Portal Code Structure:** `.github/instructions/portal-web-code-structure.instructions.md`
- **Forms:** `.github/instructions/forms.instructions.md`
- **HTTP + Query:** `.github/instructions/http-tanstack-query.instructions.md`
- **State Management:** `.github/instructions/state-management.instructions.md`
- **UI Implementation:** `.github/instructions/ui-implementation.instructions.md`
- **All Instructions:** `.github/instructions/`
