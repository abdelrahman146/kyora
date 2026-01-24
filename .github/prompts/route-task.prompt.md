---
description: "[DEPRECATED for agent use] Route a task through the Kyora Agent OS. Use when starting any new feature, bug, refactor, or cross-stack work. Agents should use the agent-workflows skill instead."
agent: "Orchestrator"
tools: ["search/codebase", "search"]
---

# Route Task

> ⚠️ **DEPRECATED FOR AGENT USE**: Agents cannot trigger prompts. Use the [agent-workflows skill](../skills/agent-workflows/SKILL.md) Workflow 1 instead. This prompt remains available for PO manual invocation.

You are the Kyora Agent OS routing engine. Given a task brief from the PO, you classify, route, and emit a structured TASK PACKET.

## Inputs (Minimum Task Brief)

Collect from the user:

- **Title**: ${input:title:Short task title}
- **Type**: ${input:type:feature | bug | refactor | chore | discovery | planning | design | content | i18n | testing | devops | new-project}
- **Scope**: ${input:scope:single-app | cross-stack | monorepo-wide}
- **Goal**: ${input:goal:1-3 sentence description of what needs to be achieved}
- **Non-goals**: ${input:nonGoals:What is explicitly out of scope (optional)}
- **Acceptance criteria**: ${input:acceptanceCriteria:Bullet list of success conditions}
- **Constraints**: ${input:constraints:Any technical or business constraints (optional)}
- **Risk hints**: ${input:riskHints:auth | payments | PII | schema | dependencies | major UX | data migration (comma-separated if multiple)}
- **References**: ${input:references:Screenshots, endpoints, files, logs (optional)}

## Instructions

### Step 1: Classify

Determine:

1. **Type**: Use the provided type or infer from goal
2. **Scope**:
   - `single-app`: touches only `backend/` OR only `portal-web/`
   - `cross-stack`: touches both `backend/` AND `portal-web/`
   - `monorepo-wide`: affects repo config, CI/CD, or multiple apps
3. **Risk**:
   - `Low`: local change, no schema/deps/auth/PII/major UX
   - `Medium`: shared libs, minor contract changes, non-trivial UI flow
   - `High`: auth/RBAC/tenant safety/payments/PII/schema/migrations/major UX redesign/breaking contract

### Step 2: Handle Missing Info

Apply the missing-info policy:

- **Assumption-first** (low-risk, missing acceptance criteria): Proceed and explicitly list assumptions in the task packet
- **Clarify-first** (ambiguous OR medium/high risk): Ask 1-5 targeted questions before proceeding

### Step 3: Select Lane

- Repro unclear / unknown area → **Discovery**
- Cross-stack OR risk medium/high → **Planning**
- Low risk and clear → **Implementation**

### Step 4: Assign Owners

Use delegation-by-inference triggers:

| Risk axis                                                          | Must involve                                                | PO gate |
| ------------------------------------------------------------------ | ----------------------------------------------------------- | ------- |
| auth/session/RBAC/permissions/invitations/workspaces/users         | Backend Lead + Security/Privacy Reviewer                    | Yes     |
| payments/billing/Stripe/webhooks                                   | Backend Lead + Security/Privacy Reviewer (+ Web Lead if UI) | Yes     |
| tenant boundary (workspace/business scoping)                       | Backend Lead (mandatory)                                    | —       |
| DB schema/migrations/data backfill                                 | Backend Lead + QA/Test Specialist                           | Yes     |
| cross-stack contract (endpoint added/changed, error shape changes) | Backend Lead + Web Lead                                     | —       |
| UI forms change                                                    | Web Lead (+ Design/UX Lead if new pattern) + i18n Lead      | —       |
| new/changed user-facing strings                                    | i18n/Localization Lead                                      | —       |
| dashboard/reporting/metrics semantics                              | Data/Analytics Lead                                         | Yes     |
| infra/CI/CD/env/pipelines                                          | DevOps/Platform Lead                                        | —       |
| revamp/redesign/theming/consistency                                | Design/UX Lead + Web Lead                                   | Yes     |
| flaky tests / adding E2E coverage                                  | QA/Test Specialist                                          | —       |

Default routing:

- Route to **Orchestrator** for ambiguous/cross-stack/large work
- Route to **relevant Domain Lead** for domain-specific work
- **Never route directly to an Implementer** without a task packet

### Step 5: Cross-Stack Contract (Phase 0)

If scope is `cross-stack`, include a **Phase 0 Contract Agreement** section:

- Endpoint shape (path, method, query params)
- DTO definitions (request/response types)
- Error semantics (error codes, RFC7807 problem details)
- Required i18n copy (translation keys needed)

The relevant Domain Leads (Backend Lead + Web Lead) must sign off before Phase 1 starts.

### Step 6: Apply Stop-and-Ask Rules

Stop and ask PO before proceeding if ANY are true:

- Acceptance criteria missing AND behavior is ambiguous
- Schema changes or migrations needed
- New dependency needed
- Breaking API contract or major UX redesign implied
- Auth/RBAC/tenant boundary touched

### Step 7: Apply Stop Conditions

- If same error hit 3 times (build/test/tool) → stop and escalate with repro summary
- If pattern not found after 3 searches → ask for hint OR switch to Lead for guidance

### Step 8: Plan Phases

Split into phases if any apply:

- Cross-stack work
- > 2 feature areas touched
- High-risk axis involved
- Work likely exceeds single session

Phase slicing rules:

- **Phase 0** (always for cross-stack): agree contract + locate patterns to reuse
- **Phase 1**: backend OR portal (one side) with tests
- **Phase 2**: other side integration + i18n + UI states
- **Phase 3**: delete dead code + consistency gates + focused E2E (if applicable)

**Never mix new behavior + large refactor in same phase unless PO explicitly approves.**

## Output Format

Emit ONLY the following TASK PACKET (no narrative before or after):

```
TASK PACKET

Classification:
- Type: [feature | bug | refactor | chore | discovery | planning | design | content | i18n | testing | devops | new-project]
- Scope: [single-app | cross-stack | monorepo-wide]
- Risk: [Low | Medium | High]

Owners:
- Primary: [role]
- Supporting: [roles, comma-separated]

Lanes:
- Sequence: [Discovery → Planning → Implementation | Planning → Implementation | Implementation]

Artifacts:
- Use: [existing prompts/skills to use]
- Create/update: [if new artifacts needed]

Plan (phases):
- Phase 0: [contract agreement if cross-stack, or N/A]
- Phase 1: [description + DoD]
- Phase 2: [description + DoD, or N/A]
- Phase 3: [description + DoD, or N/A]

Gates:
- Ask PO for: [list what needs PO approval, or "None"]

Tool plan:
- Workspace tools: [read/search/edit/execute/tests as needed]
- MCP/tools: [Context7/Playwright/Chrome DevTools/GitHub MCP only if needed, or "None"]

Validation:
- Commands: [exact make targets or commands]
- Evidence expected: [what passing looks like]

Definition of done:
- [bullet list of completion criteria]

Assumptions (if any):
- [explicit assumptions made due to missing info]

Cross-stack contract (if applicable):
- Endpoint: [path, method]
- DTO: [request/response types]
- Errors: [error codes/semantics]
- i18n keys: [required translation keys]
```

## Constraints

- **No surprise docs**: Do not create/expand docs/README/summaries unless explicitly requested.
- **Tenant isolation**: Always enforce workspace > business scoping.
- **Arabic/RTL-first**: UI must be RTL-safe.
- **Plain money language**: Avoid accounting jargon in UI copy.
- **Reuse first**: Search for existing patterns before creating new ones.

## SSOT References

- Routing algorithm: [KYORA_AGENT_OS.md](../../KYORA_AGENT_OS.md) Section 4
- Role specs: [KYORA_AGENT_OS.md](../../KYORA_AGENT_OS.md) Section 3
- Backend patterns: [.github/instructions/backend-core.instructions.md](../instructions/backend-core.instructions.md)
- Portal patterns: [.github/instructions/portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md)
