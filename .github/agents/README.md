# Kyora Custom Agents

Custom GitHub Copilot agents tailored for Kyora's development workflows. These agents are specialized for building production-ready features across Kyora's full stack (Go backend + React frontend) with a focus on multi-tenancy, RTL/Arabic-first UX, and social commerce workflows.

## Available Agents

### 🏗️ **Feature Builder** (Most Common)

**Use when**: Building complete features from backend to frontend

The orchestrator agent that builds end-to-end features by coordinating between specialists. Handles everything from database models to UI components, ensuring consistency across the stack.

**Invokes**: Backend Specialist, Portal-Web Specialist  
**Best for**: "Add expense tracking feature", "Build customer loyalty program", "Implement discount codes"

---

### 🧰 **Maintenance Engineer** (Bugfix/Refactor)

**Use when**: Fixing fullstack bugs, refactoring, or reducing drift (without adding features)

Cross-stack maintenance agent focused on correctness, SSOT compliance, and small safe changes.

**Best for**: "Fix a regression across API + UI", "Refactor drifted code to match SSOT", "Improve performance without changing behavior"

**⚠️ Manual activation required**: Set `infer: false`

---

### 🧱 **Backend Maintenance**

**Use when**: Fixing backend bugs or refactoring backend code (no new endpoints)

Focuses on tenant isolation, RBAC/plan gates, error handling, and query performance.

**Best for**: "Fix a storage scoping bug", "Remove duplication in a service", "Optimize a slow query"

**⚠️ Manual activation required**: Set `infer: false`

---

### 🧼 **Portal Web Maintenance**

**Use when**: Fixing portal-web bugs or cleaning up drift (no new pages)

Focuses on TanStack Router/Query/Form correctness, state ownership, i18n parity, and RTL safety.

**Best for**: "Fix broken list filters", "Remove ad-hoc API calls", "Align a route with URL-state rules"

**⚠️ Manual activation required**: Set `infer: false`

---

### 🔧 **Backend Specialist**

**Use when**: Implementing backend-only features or fixing backend issues

Go backend expert specializing in domain-driven design, GORM/Postgres, multi-tenancy, Stripe/Resend integration, and E2E testing with testcontainers.

**Best for**: "Add API endpoint for bulk imports", "Fix order state machine", "Implement recurring billing"

---

### 🎨 **Portal-Web Specialist**

**Use when**: Building frontend-only features or fixing UI issues

React 19 + TanStack stack expert specializing in mobile-first, RTL/Arabic-first UI with daisyUI, Chart.js visualizations, and i18n.

**Best for**: "Create analytics dashboard", "Fix mobile layout issues", "Add Arabic date formatting"

---

### 🏛️ **Domain Architect** (Planning Phase)

**Use when**: Designing new domains before implementation

Plans domain architecture, data models, API contracts, state machines, and integration points. Creates comprehensive specs for other agents to implement.

**Best for**: "Design subscription management domain", "Plan multi-warehouse inventory", "Architect loyalty points system"

**⚠️ Manual activation required**: Set `infer: false`

---

### 🧪 **E2E Test Specialist**

**Use when**: Writing comprehensive backend tests

Testcontainers expert who writes E2E tests covering complete workflows, multi-tenancy isolation, business rules, and integration flows.

**Best for**: "Add tests for order workflow", "Test multi-tenancy isolation", "Cover error scenarios for payments"

---

### 🔍 **SSOT Auditor** (Code Quality)

**Use when**: Checking code compliance with instruction files

Audits codebase for pattern violations, instruction drift, missing documentation, and SSOT integrity. Generates comprehensive audit reports.

**Best for**: "Audit multi-tenancy scoping", "Check forms compliance", "Find undocumented patterns"

**⚠️ Manual activation required**: Set `infer: false`

---

### 🤖 **AI Architect** (Infrastructure)

**Use when**: Maintaining Copilot AI customization layer

Maintains `.github/` AI infrastructure (agents, instructions, prompts, skills). Keeps instructions synced with codebase reality.

**Best for**: "Update backend instructions", "Fix instruction drift", "Create new instruction file"

**⚠️ Manual activation required**: Set `infer: false`

---

## Quick Start Guide

### Option 1: Use in VS Code

1. Open VS Code with Copilot enabled
2. Start a chat session
3. Mention the agent you want (Copilot will auto-select based on context)
4. For manual agents (⚠️), explicitly select from agent dropdown

**Example prompts**:

```
@Feature-Builder Add expense tracking with recurring expenses support

@Maintenance-Engineer Fix a bug without adding features; include the smallest safe refactor if needed

@Backend-Maintenance Fix backend bug and add/adjust E2E coverage

@Portal-Web-Maintenance Refactor drifted portal-web code to match TanStack + i18n/RTL rules

@Backend-Specialist Fix order state machine to handle partial refunds

@Portal-Web-Specialist Create mobile-first analytics dashboard with charts

@Domain-Architect Design loyalty points system with tier-based rewards

@E2E-Test-Specialist Add tests for complete order workflow including payments

@SSOT-Auditor Audit multi-tenancy scoping across all domains
```

### Option 2: Use on GitHub.com (Copilot Coding Agent)

1. Go to [github.com/copilot/agents](https://github.com/copilot/agents)
2. Select repository and branch
3. Choose agent from dropdown
4. Enter task description
5. Copilot creates PR with implementation

---

## Agent Skills

Reusable, on-demand workflows live under `.github/skills/`.

- Skills inventory: [.github/skills/README.md](../skills/README.md)

---

## Agent Capabilities Matrix

| Agent                      | Backend           | Frontend          | Testing              | Design           | Audit             |
| -------------------------- | ----------------- | ----------------- | -------------------- | ---------------- | ----------------- |
| **Feature Builder**        | ✅ (orchestrates) | ✅ (orchestrates) | ✅ (via specialists) | ❌               | ❌                |
| **Maintenance Engineer**   | ✅                | ✅                | ✅ (targeted)        | ❌               | ✅ (drift fixes)  |
| **Backend Maintenance**    | ✅✅✅            | ❌                | ✅ (E2E)             | ❌               | ✅ (SSOT fixes)   |
| **Portal Web Maintenance** | ❌                | ✅✅✅            | ✅ (typecheck/build) | ❌               | ✅ (SSOT fixes)   |
| **Backend Specialist**     | ✅✅✅            | ❌                | ✅ (E2E)             | ❌               | ❌                |
| **Portal-Web Specialist**  | ❌                | ✅✅✅            | ❌                   | ✅ (UX patterns) | ❌                |
| **Domain Architect**       | ✅ (design only)  | ❌                | ✅ (test planning)   | ✅✅✅           | ❌                |
| **E2E Test Specialist**    | ✅ (tests only)   | ❌                | ✅✅✅               | ❌               | ❌                |
| **SSOT Auditor**           | ✅ (audit)        | ✅ (audit)        | ✅ (audit)           | ❌               | ✅✅✅            |
| **AI Architect**           | ❌                | ❌                | ❌                   | ✅ (docs)        | ✅ (instructions) |

Legend: ✅✅✅ = Primary expertise | ✅ = Can do | ❌ = Not in scope

---

## Decision Tree: Which Agent to Use?

```
Start here
    ↓
Fixing bugs, refactoring, or reducing drift?
    YES → Maintenance Engineer (or Backend Maintenance / Portal Web Maintenance)
    NO → Continue
        ↓
Building a complete feature?
    YES → Feature Builder (orchestrates everything)
    NO → Continue
        ↓
Backend work only?
    YES → Backend Specialist
    NO → Continue
        ↓
Frontend work only?
    YES → Portal-Web Specialist
    NO → Continue
        ↓
Need to plan/design first?
    YES → Domain Architect (creates spec)
    NO → Continue
        ↓
Writing tests?
    YES → E2E Test Specialist
    NO → Continue
        ↓
Checking code quality?
    YES → SSOT Auditor
    NO → Continue
        ↓
Maintaining AI layer (.github/)?
    YES → AI Architect
    NO → Use default Copilot

```

---

## Agent Design Philosophy

### 1. Specialization Over Generalization

Each agent has a narrow, well-defined scope. This prevents confusion and ensures expertise.

### 2. Orchestration Pattern

Feature Builder delegates to specialists rather than doing everything itself. This maintains separation of concerns.

### 3. Instruction-First

All agents reference `.github/instructions/` files as SSOT. They implement what's documented, not what they "think" is best.

### 4. Production-Ready Code

No TODOs, no FIXMEs, no "example" code. Everything is complete and ready for production.

### 5. Context-Aware

Agents understand Kyora's unique context:

- Arabic-first, RTL-native UX
- Multi-tenancy (workspace + business)
- Social commerce workflows (DM-driven orders)
- Mobile-heavy users with low tech literacy

---

## How Agents Work Together

### Example: Building "Discount Codes" Feature

**Step 1: Planning (optional)**

```
@Domain-Architect Design discount codes domain with percentage/fixed amount types,
usage limits, expiry dates, and order integration
```

Output: Comprehensive design doc

**Step 2: Implementation**

```
@Feature-Builder Implement discount codes feature based on design doc
```

Feature Builder automatically:

1. Invokes `@Backend-Specialist` to create domain module + API
2. Invokes `@Portal-Web-Specialist` to create UI + forms
3. Coordinates E2E testing
4. Verifies multi-tenancy + plan gates
5. Ensures frontend/backend consistency

**Step 3: Testing (if more coverage needed)**

```
@E2E-Test-Specialist Add edge case tests for discount code stacking and expiry
```

**Step 4: Audit (before merging)**

```
@SSOT-Auditor Audit discount codes implementation for multi-tenancy and pattern compliance
```

---

## Tool Scoping

Each agent has carefully scoped tool access:

- **Read-only agents** (Domain Architect, SSOT Auditor): `read`, `search`, `grep_search`, `semantic_search`, `usages`
- **Implementation agents** (Backend/Portal-Web Specialist): + `edit`, `multi_replace`, `get_errors`
- **Orchestrator** (Feature Builder): + `agent` (for runSubagent)
- **AI Infrastructure** (AI Architect): Limited to `.github/**` paths

This prevents agents from accidentally modifying code outside their scope.

---

## Best Practices

### ✅ Do

- **Use Feature Builder for new features** - It orchestrates everything correctly
- **Let agents reference instructions** - They know where to look
- **Provide business context** - Explain the "why", not just the "what"
- **Review agent output** - Especially for complex features
- **Use Domain Architect for complex domains** - Design before implementing
- **Run SSOT Auditor before big PRs** - Catch violations early

### ❌ Don't

- **Don't micro-manage agents** - Let them follow their patterns
- **Don't mix agent scopes** - Use the right agent for the job
- **Don't skip testing** - E2E Test Specialist exists for a reason
- **Don't ignore audit reports** - SSOT violations accumulate into tech debt
- **Don't modify AI infrastructure without AI Architect** - Let it maintain consistency

---

## Troubleshooting

### Agent not activating automatically

**Cause**: `infer: false` or unclear prompt

**Solution**: Manually select agent from dropdown, or rephrase prompt with more context

---

### Agent doing too much / too little

**Cause**: Wrong agent selected for the task

**Solution**: Check decision tree above and use appropriate agent

---

### Agent output violates patterns

**Cause**: Instruction drift or missing instruction file

**Solution**:

1. Run `@SSOT-Auditor` to find violations
2. Use `@AI-Architect` to update instructions
3. Re-run implementation agent

---

### Backend and frontend not aligned

**Cause**: Using specialists separately instead of Feature Builder

**Solution**: Always use Feature Builder for full-stack features

---

## Contributing

### Adding New Agents

1. Read [.github/instructions/ai-infrastructure.instructions.md](../instructions/ai-infrastructure.instructions.md)
2. Use `@AI-Architect` to create new agent file
3. Test with representative tasks
4. Document in this README
5. Update decision tree if needed

### Modifying Existing Agents

1. Discuss changes (why is the current scope insufficient?)
2. Use `@AI-Architect` to update agent frontmatter/body
3. Verify no conflicts with other agents
4. Test with existing task examples
5. Update documentation

---

## Agent Metrics

Track agent effectiveness (update quarterly):

| Agent                 | Avg. Task Success Rate | Avg. Iterations | Common Issues |
| --------------------- | ---------------------- | --------------- | ------------- |
| Feature Builder       | TBD                    | TBD             | TBD           |
| Backend Specialist    | TBD                    | TBD             | TBD           |
| Portal-Web Specialist | TBD                    | TBD             | TBD           |
| Domain Architect      | TBD                    | TBD             | TBD           |
| E2E Test Specialist   | TBD                    | TBD             | TBD           |
| SSOT Auditor          | TBD                    | TBD             | TBD           |
| AI Architect          | TBD                    | TBD             | TBD           |

---

## Resources

- **Instruction Files**: `.github/instructions/` - SSOT for all patterns
- **AI Infrastructure Guide**: `.github/instructions/ai-infrastructure.instructions.md`
- **Copilot Instructions**: `.github/copilot-instructions.md` - Always-on context
- **GitHub Copilot Docs**: [docs.github.com/copilot](https://docs.github.com/en/copilot)
- **VS Code Agent Docs**: [code.visualstudio.com/docs/copilot/customization/custom-agents](https://code.visualstudio.com/docs/copilot/customization/custom-agents)

---

## Support

Having issues with agents?

1. Check this README's troubleshooting section
2. Review relevant instruction files
3. Use `@SSOT-Auditor` to check for pattern violations
4. Ask in team chat with agent output + expected behavior

---

**Last Updated**: January 2026  
**Agent Count**: 7 specialized agents  
**Maintenance**: Use `@AI-Architect` for updates
