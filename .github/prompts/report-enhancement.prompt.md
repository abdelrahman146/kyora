---
name: report-enhancement
description: "Analyze improvement idea and create a comprehensive enhancement proposal with implementation plan"
agent: "AI Architect"
argument-hint: "Brief enhancement description (e.g., 'Add form validation helper for phone numbers')"
tools: ["vscode", "execute", "read", "edit", "search", "todo"]
---

# Enhancement Report Generator

You are about to analyze and document an enhancement proposed by the user:

**Enhancement Description:** `${input:enhancementDescription:Describe the improvement or missing feature}`

## Mission

Create a comprehensive enhancement proposal by:

1. Understanding the improvement idea
2. Analyzing current gaps and pain points
3. Researching existing patterns and solutions
4. Designing the enhancement
5. Checking for duplicates in existing reports

## Workflow

### Phase 1: Understanding & Context Gathering

1. **Parse the user's description**
   - Identify the category (documentation/pattern/tooling/testing/performance/dx)
   - Extract the problem being solved
   - Determine affected component(s)
   - Understand who benefits (developers/AI/users/infrastructure)

2. **Check for existing reports**
   - Search `backlog/enhancements/` for similar proposals
   - If duplicate: enhance existing report instead of creating new
   - If related: reference in the new report

3. **Gather current state**
   - What exists today (if anything)?
   - What workarounds are people using?
   - Are there partial solutions scattered around?

### Phase 2: Problem Analysis

1. **Identify pain points**
   - What friction does this solve?
   - What errors does this prevent?
   - What time does this save?
   - What confusion does this eliminate?

2. **Determine scope**
   - Is this a localized improvement or systemic change?
   - Does this affect one component or multiple?
   - Is this backward compatible?

3. **Research existing solutions**
   - Check if this pattern exists elsewhere in the codebase
   - Look for similar solutions in instruction files
   - Consider industry best practices (but prefer Kyora patterns)

### Phase 3: Solution Design

1. **Propose the enhancement**

   **For documentation enhancement:**
   - Which instruction file(s) need updating?
   - What sections to add/modify?
   - What examples to include?

   **For pattern enhancement:**
   - Where should the new pattern live?
   - What's the API/interface?
   - How does it integrate with existing code?

   **For tooling enhancement:**
   - What script/tool is needed?
   - What dependencies are required?
   - How is it invoked?

   **For testing enhancement:**
   - What test coverage is missing?
   - What test utilities are needed?
   - What test patterns should be standardized?

2. **Design implementation**
   - Break into phases/milestones
   - Identify files to create/modify
   - Estimate effort
   - Assess risk

3. **Consider alternatives**
   - Are there other ways to solve this?
   - What are the tradeoffs?
   - Why is the proposed approach best?

### Phase 4: Impact Assessment

Evaluate benefits:

- **Consistency**: Improves uniformity how?
- **Productivity**: Time saved per use?
- **Quality**: Errors prevented?
- **Maintainability**: Future work made easier how?
- **Onboarding**: Friction reduced for new developers/AI?

Priority guidelines:

- **High**: Solves recurring pain, affects many, quick win
- **Medium**: Valuable improvement, moderate scope
- **Low**: Nice-to-have, niche benefit, large effort

### Phase 5: Success Criteria

Define measurability:

- How do we know this is working?
- What metrics indicate success?
- What does "done" look like?

Examples:

- "No more manual X needed"
- "Test coverage increases by Y%"
- "Documentation mentions pattern Z"
- "Build time reduces by N seconds"

### Phase 6: Report Generation

Create enhancement report at: `backlog/enhancements/YYYY-MM-DD-<slug>.md`

Use the enhancement template with:

- **Clear problem statement**: Why this matters
- **Concrete solution**: Specific implementation details
- **Alternatives considered**: With pros/cons
- **Implementation plan**: Phases with effort estimates
- **Success criteria**: Measurable outcomes
- **References**: Links to instruction files

**Slug naming**: Use kebab-case, max 50 chars, descriptive
Examples: `phone-validation-helper`, `test-coverage-orders`, `form-error-display-pattern`

### Phase 7: Verification

Before finalizing:

- [ ] Checked for duplicates in `backlog/enhancements/`
- [ ] Solution is concrete and actionable
- [ ] Implementation plan is realistic
- [ ] Success criteria are measurable
- [ ] Alternatives are considered
- [ ] Priority reflects value vs. effort
- [ ] Category is set correctly
- [ ] Frontmatter is valid YAML

## Output Format

Provide:

1. **The generated enhancement report** (full path)
2. **Brief summary** of proposal:

   ```
   Enhancement Created: backlog/enhancements/2026-01-18-[slug].md

   Priority: [high|medium|low]
   Category: [category]
   Component: [component]
   Problem: [one-line problem statement]
   Solution: [one-line solution summary]
   Effort: [estimate]
   Impact: [who benefits and how]

   Related Reports: [if any]
   ```

## Safety Checks

- Do not modify production code (only create report in `backlog/enhancements/`)
- Verify all file references are valid
- Ensure frontmatter is complete and valid
- Create `backlog/enhancements/` directory if needed
- Keep proposals scoped and actionable

## Investigation Tips

**Enhancement categories:**

**Documentation:**

- Missing instruction files
- Incomplete coverage of patterns
- Need for examples
- Outdated references

**Pattern:**

- Repetitive code that should be abstracted
- Missing utilities/helpers
- Inconsistent approaches to same problem
- New architectural pattern needed

**Tooling:**

- Manual process that should be automated
- Missing developer scripts
- Build/deploy improvements
- Code generation opportunities

**Testing:**

- Missing test coverage
- Test utilities needed
- E2E scenarios not covered
- Test patterns to standardize

**Performance:**

- Slow queries/operations
- Bundle size reductions
- Caching opportunities
- Database optimization

**Developer Experience:**

- Friction in common workflows
- Confusing patterns
- Missing guardrails
- Better error messages

## Common Enhancement Patterns

**Backend:**

- Service helpers for common operations
- Middleware for cross-cutting concerns
- Repository scopes for complex queries
- Test fixtures for E2E tests

**Portal-web:**

- Form field components
- Reusable hooks
- API client patterns
- State management utilities

**Cross-cutting:**

- Instruction file updates
- Shared type definitions
- Error handling patterns
- Logging/monitoring improvements

## References

- Enhancement Template: `.github/prompts/templates/enhancement.template.md`
- AI Infrastructure: `.github/instructions/ai-infrastructure.instructions.md`
- Backend Patterns: `.github/instructions/backend-core.instructions.md`
- Portal Patterns: `.github/instructions/portal-web-architecture.instructions.md`
