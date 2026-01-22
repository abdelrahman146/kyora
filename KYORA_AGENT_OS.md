# KYORA_AGENT_OS.md

Kyora Agent OS is the operating model for how humans + GitHub Copilot artifacts (instructions, prompt files, custom agents, skills, MCP tools) work together in this monorepo.

This repo is a monorepo (currently includes `backend/` and `portal-web/`), and **must remain future-proof** as more apps/services are added.

**Prime Directive:** Correct + consistent beats fast.

## Versioning

- **Version**: 2026-01-22-v2
- **Changelog**:
  - **BREAKING**: Corrected artifact capability model — prompts are PO-only (user-triggered), skills are agent-runnable workflows. Agents cannot trigger prompts.
  - Added autonomous agent delegation via `agent` tool with `infer: true` configuration.
  - Replaced manual handoff patterns with delegation-first + PO approval gates.
  - Added scoped AGENTS.md files (backend/, portal-web/, .github/) for token efficiency.
  - New PO prompts: session continuation, drift resolution, SSOT maintenance.
  - Merged agent-workflow prompts into skills (packets, routing logic now in skills).
  - Rebuilt as an operational spec: routing algorithm, lanes, gates, templates, examples.
  - Added "recovery lane" for new sessions/partial history.
  - Added MCP/tooling policy aligned with current VS Code prompt files + MCP behavior.
  - Strengthened delegation-by-inference, handoff packets, and drift-sync governance.

---

## Table of Contents

1. Operating Principles
2. Artifact Strategy (SSOT + Progressive Disclosure)
3. Autonomous Agent Delegation Model
4. Team Topology and Hierarchy (Roles)
5. Routing Algorithm (Task Intake → Lane → Owner → Artifacts → Tools)
6. Lanes (Continuous Cycle)
7. Quality Gates and Checklists
8. MCP + Tooling Policy
9. Examples (Concrete Scenarios)
10. Governance and Maintenance

Appendix A. Handoff Packet Templates
Appendix B. Drift Ticket Template
Appendix C. Prompt + Skill Design Checklist

---

## 1) Operating Principles

### What the OS governs

- Task routing: who does what, in which order.
- Lane discipline: discovery vs planning vs implementation vs validation.
- Tool policy: least-privilege, minimal steps, safe MCP usage.
- Gates: when PO approval is required.
- Consistency enforcement: reuse-first + checklists.
- Drift-sync: keeping SSOT instructions aligned with reality.

### What stays ad-hoc

- Product prioritization and scope negotiation (PO-owned).
- One-off work that isn’t worth a reusable artifact.
- Deep architectural tradeoffs that require human judgment (unless delegated).

### Non-negotiables (Kyora)

- **Tenant safety**: strict tenant isolation (workspace top-level; business second-level). Never cross-scope.
- **Arabic/RTL-first**: UI must be RTL-safe and Arabic-first; avoid left/right assumptions.
- **Plain money language**: avoid accounting jargon in UI copy.

### Prime rules

- **No surprise docs**: Don’t create/expand docs/README/summaries unless explicitly requested.
- **Reuse first**: Search and align with existing patterns/components/utils before creating new ones.
- **Token discipline**:
  - Prefer incremental phases with verification.
  - Ask early when requirements are ambiguous.
  - If stopping mid-task, write a handoff packet.

---

## 2) Artifact Strategy (SSOT + Progressive Disclosure)

### 2.1 SSOT rule

- **Never duplicate rules that already exist** in SSOT instruction files.
- Instead: **link to SSOT** and extract only the minimum "must remember" bullets.

Core SSOT entry points:

- Repo baseline: [.github/copilot-instructions.md](.github/copilot-instructions.md)
- Artifact selection: [.github/instructions/ai-artifacts.instructions.md](.github/instructions/ai-artifacts.instructions.md)
- Prompt spec: [.github/instructions/prompts.instructions.md](.github/instructions/prompts.instructions.md)
- Agents spec: [.github/instructions/agents.instructions.md](.github/instructions/agents.instructions.md)
- Skills spec: [.github/instructions/agent-skills.instructions.md](.github/instructions/agent-skills.instructions.md)

Kyora domain SSOTs (examples; not exhaustive):

- Backend core: [.github/instructions/backend-core.instructions.md](.github/instructions/backend-core.instructions.md)
- Backend patterns: [.github/instructions/go-backend-patterns.instructions.md](.github/instructions/go-backend-patterns.instructions.md)
- Errors: [.github/instructions/errors-handling.instructions.md](.github/instructions/errors-handling.instructions.md)
- DTOs/Swagger: [.github/instructions/responses-dtos-swagger.instructions.md](.github/instructions/responses-dtos-swagger.instructions.md)
- Portal architecture: [.github/instructions/portal-web-architecture.instructions.md](.github/instructions/portal-web-architecture.instructions.md)
- Portal structure: [.github/instructions/portal-web-code-structure.instructions.md](.github/instructions/portal-web-code-structure.instructions.md)
- UI/RTL: [.github/instructions/ui-implementation.instructions.md](.github/instructions/ui-implementation.instructions.md)
- Forms: [.github/instructions/forms.instructions.md](.github/instructions/forms.instructions.md)
- i18n: [.github/instructions/i18n-translations.instructions.md](.github/instructions/i18n-translations.instructions.md)
- HTTP/TanStack Query: [.github/instructions/http-tanstack-query.instructions.md](.github/instructions/http-tanstack-query.instructions.md)

### 2.2 Artifact decision table

| Need | Create/Use | Why |
|---|---|---|
| Always-on coding standard for a file glob | `.instructions.md` | Enforced automatically, scoped |
| Reusable single-purpose task | `.prompt.md` | On-demand `/`, consistent outputs |
| Persona + tool restrictions + handoffs | `.agent.md` | Role specialization + safety |
| Repeatable multi-step workflow + bundled references/templates/scripts | Skill (`.github/skills/<name>/SKILL.md`) | Progressive disclosure |
| Project-wide guidance for all agents | `AGENTS.md` | Quick onboarding + boundaries |

### 2.3 Prompt / agent mechanics (minimal reminders)

- Prompt files are `.prompt.md` and can specify `name`, `description`, `argument-hint`, `agent`, `tools`, and `model`.
- Tool lists are prioritized: **prompt tools → referenced agent tools → default agent tools**.
- Use the smallest tool list possible; tool overload hurts reliability.

### 2.4 Prompts = Product Owner (PO) executables ONLY

**Critical**: Prompts are **user-triggered only** via `/prompt-name`. Agents **cannot** trigger prompts.

Prompts are for the PO to reliably "launch" repeatable work with clean inputs.

Use prompts when:

- The PO keeps asking the same type of request (e.g., "add endpoint + wire portal query + add i18n keys").
- The task benefits from consistent framing: scope, risks, gates, validation commands.

Prompt rules:

- Output must be a task packet + lane/owner routing (not a long narrative).
- Include `${input:...}` fields for the Minimum Task Brief (section 5.1) and risk hints.
- Default to the Orchestrator or the relevant Lead; do not hardcode a specific implementation.
- Tools: least-privilege. Only include MCP tools when the prompt’s job actually needs them.
- "No surprise docs" clause must be included.
- **Never create prompts expecting agents to call them** — convert to skills instead.

### 2.5 Skills = repeatable multi-step workflows (agent-runnable)

Skills are the reliability engine for agentic work in a token-limited world.

**Critical**: Skills are **agent-discoverable**. Agents can read and execute skills when the task matches the skill's description. This is how agents access repeatable workflows.

Use skills when:

- The workflow is multi-step and repeatable (especially cross-stack).
- The workflow benefits from bundled templates/checklists/scripts.
- Agents need to autonomously execute the workflow (not prompts!).

Skill rules:

- Must be phase-sliced (see section 5.8) with a required handoff packet between phases.
- Must reference SSOT instructions instead of copying their rules.
- Must define a compact success checklist and exact validation commands.

Optional: keep a short "skill registry" table inside the skill itself, not in this OS.

### 2.6 Artifact maintenance protocol

Keep artifacts tidy and predictable:

- Locations:
  - Prompts: `.github/prompts/`
  - Agents: `.github/agents/`
  - Skills: `.github/skills/<skill-name>/`
  - Instructions: `.github/instructions/`
- Naming:
  - Prompts: `<verb>-<object>.prompt.md`
  - Agents: `<role>.agent.md` (role-focused, not person-focused)
  - Skills: folder name is the skill name; `SKILL.md` inside
- Deprecation:
  - Mark deprecated artifacts in their description.
  - Update any prompts/agents that referenced them.
  - Remove after 1–2 releases if unused.

### 2.7 Drift detection + sync protocol

Trigger events:

- PO decides a convention change (e.g., translation key casing, naming, folder placement).
- A repeated review comment indicates a "hidden rule".
- The team starts doing something different than the instructions say.

Procedure:

1. Identify SSOT target (which `.instructions.md` owns this rule).
2. Update SSOT file(s) only (no scattered duplicates).
3. Add a single bullet to this file’s changelog summarizing the shift.
4. Run the smallest validation loop relevant to the change.

See also: "Drift-Sync Governance" in section 10 for mandatory templates + approval gates.

---

## 3) Autonomous Agent Delegation Model

This section clarifies how Copilot artifacts interact with each other and with the human Product Owner (PO).

### 3.1 Capability boundaries

| Artifact | Triggered by | Can be triggered by agents? | Primary use |
|---|---|---|---|
| **Instructions** (`.instructions.md`) | Automatic (file pattern match) | N/A (always-on) | Coding standards, conventions |
| **Prompts** (`.prompt.md`) | Human via `/prompt-name` | **No** | PO-initiated tasks with variable inputs |
| **Agents** (`.agent.md`) | Human via `@agent-name` or agent via `agent` tool | **Yes** (with `agent` tool + `infer: true`) | Specialized personas with tool restrictions |
| **Skills** (`SKILL.md`) | Agent reads skill when description matches | **Yes** (description-based discovery) | Repeatable multi-step workflows with bundled resources |

**Key insight**: Prompts are **user-facing only**. Agents cannot trigger prompts. Instead, agents can:
1. Delegate to other agents using the `agent` tool
2. Consume skills by reading their SKILL.md when the task matches the skill description

### 3.2 Agent delegation configuration

For an agent to autonomously delegate to other agents, it must have:

1. **`agent` tool** in its tools list:
   ```yaml
   tools: ['read', 'search', 'edit', 'agent']
   ```

2. **`infer: true`** in frontmatter (enables context-based agent selection):
   ```yaml
   ---
   description: 'Routes tasks and delegates to specialized agents'
   infer: true
   tools: ['read', 'search', 'agent']
   ---
   ```

3. **Delegation awareness** in the agent prompt body explaining when/how to delegate.

### 3.3 PO approval gates

When an agent needs PO approval (e.g., schema changes, new dependencies, auth changes), it must:

1. **Pause execution** and emit an "APPROVAL NEEDED" block
2. **Provide context**: what decision is needed, why, risks
3. **Wait for PO response** before continuing

**Approval Needed template**:

```markdown
---
## ⚠️ APPROVAL NEEDED

**Decision Required**: [Brief description of what needs approval]
**Context**: [Why this decision point was reached]
**Options**:
1. [Option A] — [pros/cons]
2. [Option B] — [pros/cons]

**Risks if proceeding without approval**: [List risks]

**To continue**: Reply with your decision or ask for more context.
---
```

### 3.4 Delegation patterns

#### Lead → Implementer delegation

```
PO → @orchestrator → (routes to) → @backend-lead → (delegates via agent tool) → @backend-implementer
```

The Backend Lead:
1. Creates a clear task spec
2. Uses `agent` tool to invoke Backend Implementer
3. Reviews Implementer output before completing

#### Multi-domain coordination

```
PO → @orchestrator → classifies as cross-stack →
  → @backend-lead (Phase 1: API contract) → APPROVAL NEEDED (schema gate)
  → @web-lead (Phase 2: UI integration plan)
  → @backend-implementer + @web-implementer (Phase 3: implementation)
```

### 3.5 Skills vs prompts for agent work

**Use a Skill** when:
- Agents need to execute a repeatable workflow autonomously
- The workflow has multiple steps with bundled resources
- You want progressive disclosure (only load when needed)

**Use a Prompt** when:
- The PO needs to initiate work with specific inputs
- The task requires human judgment on parameters
- You want consistent task framing visible to the PO

**Never do**:
- Create prompts expecting agents to trigger them
- Put agent-only workflows in prompts (convert to skills)

### 3.6 Scoped AGENTS.md files

To optimize token usage and provide domain-specific context:

| Location | Purpose |
|---|---|
| Root `AGENTS.md` | General project context, boundaries, tech stack overview |
| `backend/AGENTS.md` | Go patterns, API conventions, testing requirements |
| `portal-web/AGENTS.md` | React/TypeScript patterns, i18n rules, UI conventions |
| `.github/AGENTS.md` | Artifact creation rules, OS governance, prompt/skill patterns |

Nested `AGENTS.md` files override parent context (closest file wins).

---

## 4) Team Topology and Hierarchy (Roles)

### Hierarchy

- **Product Owner (human)**: decision authority, scope, acceptance criteria, gates.
- **Chief-of-Staff / Orchestrator (agent)**: routes tasks, enforces lanes + gates + handoffs.
- **Domain Leads (agents)**: Backend Lead, Web Lead, Design/UX Lead, Data/Analytics Lead, DevOps/Platform Lead, Content/Marketing Lead, i18n/Localization Lead.
- **Implementers (agents)**: Backend Implementer, Web Implementer, Shared/Platform Implementer.
- **Quality/Safety (agents)**: QA/Test Specialist, Security/Privacy Reviewer (read-only).

Tool shorthand:

- Workspace tools: read/search/edit/execute/tests
- MCP tools: GitHub MCP, Context7, Playwright, Chrome DevTools

### Role specs

#### Product Owner (Human)

- When: Always
- Outputs: goal, constraints, acceptance criteria, approvals
- Allowed tools: N/A
- Forbidden: N/A
- DoD: accepts outcome; confirms gates
- Escalation: N/A

#### Chief-of-Staff / Orchestrator

- When: default entry for ambiguous/cross-stack/large work
- Outputs: classification + risk + scope; lane selection; owners; task packet; gates; handoff packets
- Allowed tools: read/search/todo (edit for plans/notes only)
- MCP: GitHub MCP optional (context only)
- Forbidden: implementing production code; adding deps; running destructive commands
- DoD: task packet complete or clarifications requested; owners/handoffs defined
- Escalation: missing acceptance criteria; schema/auth/payments/PII; multi-project scope

#### Backend Lead

- When: API contracts, domain modeling, backend architecture decisions
- Outputs: endpoint shapes, DTO decisions, error semantics, migration approach, gates
- Allowed tools: read/search/edit (specs), execute optional (validation)
- MCP: Context7 only when dependency/library usage must be verified
- Forbidden: schema changes without PO gate; large refactors without phased plan
- DoD: contract decisions explicit + testable; compatibility noted
- Escalation: auth/RBAC/tenant safety; migrations; payments

#### Web Lead

- When: portal architecture, routing/state, UI patterns, cross-stack UI integration plan
- Outputs: UI approach, component placement, i18n plan, query/mutation plan, UI states checklist
- Allowed tools: read/search/edit (spec), execute optional
- MCP: Playwright/Chrome DevTools optional for audits
- Forbidden: new design primitives without Design/UX sign-off
- DoD: plan matches portal SSOT; RTL/i18n considerations explicit
- Escalation: major UX; accessibility; contract changes

#### Design/UX Lead

- When: redesign/revamp, theming, UI consistency
- Outputs: UX spec (states/variants), acceptance criteria, RTL notes
- Allowed tools: read/search (edit only for specs)
- MCP: Playwright/Chrome DevTools for audits
- Forbidden: production code edits unless explicitly delegated
- DoD: spec is actionable and reviewable
- Escalation: brand voice; accessibility; redesign scope

#### Data/Analytics Lead

- When: dashboards, metric definitions, reporting semantics
- Outputs: metric definitions, date-range semantics, query shapes
- Allowed tools: read/search (DB read-only optional)
- MCP: Postgres MCP (read-only) optional if available/approved
- Forbidden: migrations without PO gate
- DoD: metrics unambiguous + testable
- Escalation: privacy/PII; financial semantics

#### DevOps/Platform Lead

- When: CI/CD, infra, env configs, tooling/MCP setup
- Outputs: reproducible steps, env var matrix, rollout/rollback plan
- Allowed tools: read/search/edit/execute
- MCP: GitHub MCP optional
- Forbidden: destructive infra changes without PO gate
- DoD: reproducible + rollback clear
- Escalation: prod-impact, secrets, policy

#### Content/Marketing Lead

- When: marketing copy, product copy (non-UI), website content
- Outputs: drafts aligned to Kyora voice; translation-ready text
- Allowed tools: read/search
- MCP: none by default
- Forbidden: code edits unless requested
- DoD: calm, simple, non-technical copy
- Escalation: legal/privacy claims

#### i18n/Localization Lead

- When: translation keys, Arabic-first phrasing, glossary consistency
- Outputs: keys + Arabic/English copy; glossary notes
- Allowed tools: read/search/edit
- MCP: none by default
- Forbidden: large unrelated refactors
- DoD: keys present; no hardcoded UI strings; Arabic phrasing natural
- Escalation: ambiguous meaning; domain terminology conflicts

#### Backend Implementer

- When: backend code changes
- Outputs: code + tests; OpenAPI updates if required by repo norms
- Allowed tools: read/search/edit/execute/tests
- MCP: Context7 only when a dependency/library requires validation
- Forbidden: schema changes without PO gate; cross-tenant access; new deps without gate
- DoD: acceptance criteria met; relevant tests pass; no dead code
- Escalation: unclear contracts; failing unrelated tests

#### Web Implementer

- When: portal UI + API integration + i18n
- Outputs: UI with loading/empty/error; API integration; i18n keys
- Allowed tools: read/search/edit/execute/tests
- MCP: Playwright/Chrome DevTools optional
- Forbidden: new UI primitives without lead sign-off; new deps without gate
- DoD: RTL verified; i18n complete; consistent components/tokens
- Escalation: uncertain translations; contract mismatch

#### Shared/Platform Implementer

- When: shared libs/utilities across apps/services
- Outputs: reusable utilities; minimal API surface; adoption plan
- Allowed tools: read/search/edit/execute/tests
- MCP: Context7 optional
- Forbidden: widespread refactors without plan
- DoD: reuse improves consistency; no duplication
- Escalation: breaking changes across projects

#### QA/Test Specialist

- When: adding tests, E2E, triage flakiness
- Outputs: test plan; stable tests; failure triage
- Allowed tools: read/search/edit/execute/tests
- MCP: Playwright preferred for UI flows
- Forbidden: changing production code unless asked
- DoD: tests cover change; evidence captured (commands + results)
- Escalation: non-determinism; unclear acceptance criteria

#### Security/Privacy Reviewer (read-only)

- When: auth/session/RBAC/tenant isolation, payments, PII
- Outputs: findings with severity + remediation steps
- Allowed tools: read/search
- MCP: GitHub MCP optional
- Forbidden: code edits, running commands
- DoD: findings actionable and triaged
- Escalation: critical/high issues

### Handoff contract (required)

Handoffs are mandatory when:

- Switching lane owners (e.g., Lead → Implementer)
- Switching scope (single-app → cross-stack)
- Pausing mid-task (token/time)
- Starting a new session for unfinished work

Use Appendix A templates (designed to survive token limits and "new session amnesia"):

- Delegation Packet
- Phase Handoff Packet
- Recovery/Resume Packet

Minimum rule: no role starts Implementation without a task packet + a handoff packet (unless it’s a tiny, low-risk change).

### Delegation-by-inference (required)

The Orchestrator and Domain Leads must infer supporting roles using the triggers below. This is not optional—these triggers exist to prevent partial completion, cross-stack mismatch, and rule drift.

Inference triggers (auto-involve):

- Auth/session/RBAC/permissions/invitations/workspaces/users → Backend Lead + Security/Privacy Reviewer
- Payments/billing/Stripe/webhooks → Backend Lead + Security/Privacy Reviewer (+ Web Lead if UI)
- Tenant boundary (workspace/business scoping) → Backend Lead (mandatory)
- DB schema/migrations/data backfill → Backend Lead (PO gate) + QA/Test Specialist
- Any cross-stack contract touch (endpoint added/changed, error shape changes) → Backend Lead + Web Lead
- Any UI forms change → Web Lead (+ Design/UX Lead if new pattern) + i18n Lead if user-facing copy
- Any new or changed user-facing strings → i18n/Localization Lead
- Any dashboard/reporting/metrics semantics → Data/Analytics Lead
- Any infra/CI/CD/env/pipelines → DevOps/Platform Lead
- Any "revamp/redesign/theming/consistency" request → Design/UX Lead + Web Lead
- Any flaky tests / adding E2E coverage → QA/Test Specialist

Cross-stack coordination rule:

- If Backend + Web both involved, the Leads must agree on the contract BEFORE implementation starts (Phase 0). Contract means: endpoint, DTO, error semantics, and required i18n copy.

Token-limit guardrails:

- If the request is medium/high risk OR cross-stack, do not attempt full end-to-end completion in one pass; slice into phases with verifiable outputs and handoff packets.
- If an instruction file is long, read only the relevant sections, then anchor the change using the quality gates in section 6 (don’t rely on memory).

### Decision authority (when PO approval is required)

| Decision | Owner | Requires PO gate |
|---|---:|---:|
| Acceptance criteria / scope | PO | Yes |
| New dependency (any app) | Lead + Implementer | Yes |
| DB schema/migrations | Backend Lead | Yes |
| Auth/session/RBAC changes | Backend Lead + Security | Yes |
| Breaking API contract | Backend/Web Leads | Yes |
| Major UX redesign / new primitives | Design/UX Lead | Yes |
| Minor UI changes within existing patterns | Web Lead | No |
| Refactor that changes behavior | Lead | Yes |
| Chore that does not change behavior | Implementer | No |

---

## 5) Routing Algorithm (Task Intake → Lane → Owner → Artifacts → Tools)

### 4.1 Minimum Task Brief (PO → Agent)

Copy/paste:

```
TITLE:
TYPE: feature | bug | refactor | chore | discovery | planning | design | content | i18n | testing | devops | new-project
SCOPE: single-app | cross-stack | monorepo-wide
GOAL (1–3 sentences):
NON-GOALS:
ACCEPTANCE CRITERIA (bullets):
CONSTRAINTS (bullets):
RISK HINTS: auth | payments | PII | schema | dependencies | major UX | data migration
REFERENCES: screenshots | endpoints | files | logs
```

Missing info policy (agents):

- If acceptance criteria are missing but the work is low-risk: proceed with "Assumption-first" Phase 0 and explicitly list assumptions in the task packet.
- If behavior is ambiguous or risk is medium/high: "Clarify-first" (ask 1–5 targeted questions) before implementation.

### 4.2 Classify: type + scope + risk

- **Type**: feature | bug | refactor | chore | discovery | planning | design | content | i18n | testing | devops | new-project
- **Scope**:
  - single-app
  - cross-stack (backend + portal-web)
  - monorepo-wide
- **Risk**:
  - Low: local change, no schema/deps/auth/PII/major UX
  - Medium: shared libs, minor contract changes, non-trivial UI flow
  - High: auth/RBAC/tenant safety/payments/PII/schema/migrations/major UX redesign/breaking contract

### 4.3 Lane selection (default)

- Repro unclear/unknown area → **Discovery**
- Cross-stack or risk medium/high → **Planning**
- Low risk and clear → **Implementation**

### 4.4 Assign owners (roles)

| Task type | Primary roles | Typical MCP/tools |
|---|---|---|
| feature (backend) | Backend Lead → Backend Implementer | execute/tests; Context7 only if needed |
| feature (portal) | Web Lead → Web Implementer | execute/tests; Playwright optional |
| feature (cross-stack) | Orchestrator + Backend/Web Leads + Implementers | execute/tests; Playwright optional |
| bug (clear repro) | Implementer | execute/tests |
| bug (unclear) | Orchestrator → Discovery | search/read; minimal execute |
| refactor (large) | Orchestrator + Leads | phased; tests early |
| chore | Implementer | minimal tools |
| design/UX | Design/UX Lead + Web Lead | Playwright/DevTools for audits |
| i18n | i18n Lead | read/edit only unless UI needs change |
| content | Content Lead | read-only by default |
| testing/e2e | QA/Test | Playwright |
| devops | DevOps Lead | execute; GitHub MCP optional |
| new-project | Orchestrator + relevant Lead | scaffolding + repo conventions |

Mandatory reviewers by risk axis (inference):

| Risk axis | Must involve | Requires PO gate |
|---|---|---:|
| auth/session/RBAC | Backend Lead + Security | Yes |
| payments/billing | Backend Lead + Security | Yes |
| PII/privacy | Security + relevant Lead | Yes |
| schema/migrations/backfill | Backend Lead + QA | Yes |
| breaking API contract | Backend Lead + Web Lead | Yes |
| major UX redesign/new primitives | Design/UX Lead + Web Lead | Yes |
| analytics semantics | Data/Analytics Lead | Yes |

### 4.5 Artifact selection

- Prefer existing prompts/skills.
- Create a **prompt** when the task is repeatable and single-purpose.
- Create a **skill** when the workflow is multi-step and repeatable.
- Create/update **instructions** only for always-on standards.

Artifact intent reminders (to prevent misuse):

- Prompts are PO-executable "task launchers" (standardize input, constraints, acceptance criteria, gates). Prompts are not for generating long docs.
- Skills are repeatable multi-step workflows (often cross-stack), optimized for consistency and phase slicing.
- Agents define who does what and which tools are allowed.

### 4.6 Tool selection (minimize)

- Start with workspace read/search.
- Run only the smallest relevant checks.
- Use MCP only when it clearly reduces work or increases certainty.

MCP tool selection:

- Context7: only for up-to-date library usage.
- Playwright: UI flows, RTL verification, E2E smoke.
- Chrome DevTools: perf/layout/network debugging.
- GitHub MCP: issues/PRs/policies; avoid for local-only tasks.

### 4.7 Stop-and-ask rules

Ask the PO before proceeding if any are true:

- Acceptance criteria are missing and behavior is ambiguous.
- Schema changes or migrations are needed.
- New dependency is needed.
- Breaking API contract or major UX redesign is implied.
- Auth/RBAC/tenant boundary is touched.

Safe assumptions (only for low risk):

- Minimal change; no new deps; no schema changes.
- Match existing patterns/components/tokens.
- Prefer backwards-compatible API changes.

Stop conditions (avoid tool loops / token burn):

- If you hit the same error 3 times (build/test/tool), stop and escalate with a small repro summary.
- If you can’t find a pattern after 3 searches, ask for a hint (file/feature name) OR switch to a Lead for guidance.

### 4.8 Multi-phase plan rule

Split into phases if any apply:

- Cross-stack work
- >2 feature areas touched
- High-risk axis involved
- Work likely exceeds a single session

Each phase must be independently verifiable.

Phase slicing rules (to beat token limits):

- Phase 0 (always for cross-stack): agree contract + locate patterns to reuse.
- Phase 1: backend or portal (one side) with tests.
- Phase 2: other side integration + i18n + UI states.
- Phase 3: delete dead code + consistency gates + focused E2E (if applicable).

Never mix "new behavior" + "large refactor" in the same phase unless PO explicitly approves.

### 4.9 Task Packet (output contract)

Every non-trivial request produces a task packet:

```
TASK PACKET

Classification:
- Type:
- Scope:
- Risk:

Owners:
- Primary:
- Supporting:

Lanes:
- Sequence:

Artifacts:
- Use:
- Create/update:

Plan (phases):
- Phase 0:
- Phase 1:
- Phase 2:

Gates:
- Ask PO for:

Tool plan:
- Workspace tools:
- MCP/tools:

Validation:
- Commands:
- Evidence expected:

Definition of done:
-
```

---

## 6) Lanes (Continuous Cycle)

### Default lane handoffs

| Lane | Typical owner | Default next lane |
|---|---|---|
| Discovery | Orchestrator / Lead | Planning or Deferred |
| Planning | Lead | Implementation |
| Implementation | Implementer | Review |
| Review | Lead / QA | Validation |
| Validation | QA / Implementer | Release/Follow-up (optional) |

### 5.1 Discovery

- Entry: unclear bug; unknown area; cross-stack unknown
- Output: findings + hypotheses + repro steps + next plan
- DoD: can explain the smallest next step with high confidence
- Token playbook: search first; read only the relevant slices

### 5.2 Planning

- Entry: medium/high risk; cross-stack; large refactor; UX changes
- Output: phased plan + acceptance checks + gates
- DoD: phases are verifiable; contracts and boundaries explicit
- Token playbook: keep plan compact; link SSOT files

### 5.3 Implementation

- Entry: plan approved or low-risk change
- Output: code changes + tests; minimal notes
- DoD: acceptance criteria met; touched-area checks pass
- Token playbook: small change → validate immediately

### 5.4 Review

- Entry: before declaring "done"
- Output: checklist results + fix list
- DoD: gates passed; consistency verified

### 5.5 Validation

- Entry: after changes land
- Output: command evidence (tests/build/e2e) + any follow-ups
- DoD: relevant checks green; unrelated failures triaged, not "fixed opportunistically"

### 5.6 Release / Follow-up (optional)

- Entry: deploy/pipeline/infrastructure change
- Output: rollout/rollback plan; monitoring notes

### 5.7 Recovery Lane (new session / partial completion)

Trigger: a new session continues an unfinished task.

Steps:

1. Reconstruct objective + current lane.
2. Identify what’s done from git changes + notes + test history.
3. Re-assert acceptance criteria and pending gates.
4. Continue with the next smallest verifiable phase.

Required output: a Recovery/Resume handoff packet (Appendix A) before continuing implementation.

---

## 7) Quality Gates and Checklists

These gates prevent the common failure modes: inconsistency, wrong assumptions, partial completion, and drift.

### 6.1 "No surprise docs" gate

- If user did not ask for docs: do not create or expand docs/README/long summaries.

### 6.2 "No dead code" gate

- No commented-out blocks, unused exports, unused files.
- No TODO/FIXME placeholders.

### 6.3 UI consistency gate (portal-web)

Use SSOT:

- [.github/instructions/ui-implementation.instructions.md](.github/instructions/ui-implementation.instructions.md)
- [.github/instructions/design-tokens.instructions.md](.github/instructions/design-tokens.instructions.md)
- [.github/instructions/portal-web-ui-guidelines.instructions.md](.github/instructions/portal-web-ui-guidelines.instructions.md)

Checklist:

- Uses existing components/patterns (no new primitives by default)
- RTL-safe layout and spacing
- Loading/empty/error states exist
- Copy is simple and non-technical

Reuse-first verification (required):

- Before building a new component/pattern, search for an existing one in `portal-web/src/components/` and `portal-web/src/features/`.
- Before adding a new API call, search for similar calls in `portal-web/src/api/` and reuse shared client/query utilities.

### 6.4 Forms gate

Use SSOT:

- [.github/instructions/forms.instructions.md](.github/instructions/forms.instructions.md)

Checklist:

- Uses the project form system
- Validation errors shown consistently
- Submit/disabled/server errors handled

### 6.5 i18n gate

Use SSOT:

- [.github/instructions/i18n-translations.instructions.md](.github/instructions/i18n-translations.instructions.md)

Checklist:

- No hardcoded UI strings (unless SSOT allows)
- Keys exist for all supported locales
- Arabic phrasing natural + consistent with domain language

### 6.6 Backend API contract gate

Use SSOT:

- [.github/instructions/backend-core.instructions.md](.github/instructions/backend-core.instructions.md)
- [.github/instructions/errors-handling.instructions.md](.github/instructions/errors-handling.instructions.md)
- [.github/instructions/responses-dtos-swagger.instructions.md](.github/instructions/responses-dtos-swagger.instructions.md)

Checklist:

- Inputs validated
- Tenant isolation enforced
- Errors follow Kyora Problem/RFC7807 patterns (per SSOT)
- DTOs/OpenAPI aligned (per repo norms)

Reuse-first verification (required):

- Before adding a new pattern/util, search `backend/internal/platform/utils/` and related domain modules.
- Prefer existing domain boundaries: domain logic in `backend/internal/domain/**`, infra in `backend/internal/platform/**`.

### 6.7 Cross-stack alignment gate

Checklist:

- Backend endpoint + portal API client agree on request/response
- Error semantics handled in UI (per HTTP layer SSOT)
- i18n keys added for new user-facing text

### 6.8 Testing gate

Use repo commands from `.github/copilot-instructions.md` and the Makefile.

Checklist:

- Run the smallest relevant test suite(s)
- Add/adjust tests where it’s natural
- Don’t fix unrelated failures

---

## 8) MCP + Tooling Policy

This repo uses VS Code built-in tools plus MCP servers. Keep tool usage minimal and safe.

### 7.1 Safety rules

- Treat local MCP servers as code execution: only add trusted servers.
- Never paste secrets; use VS Code input variables for sensitive values.
- Prefer tool sets (grouped tools) and keep enabled tools under the model/tool limits.

Operational notes:

- MCP servers are configured via `mcp.json` (workspace or user).
- Tool lists are cached; reset cached tools when server tools change.

### 7.2 Approved MCP servers (current)

- GitHub MCP (remote): issues/PRs/policies; keep it off for local-only tasks.
- Context7: up-to-date library docs/snippets; use sparingly.
- Playwright: UI automation, screenshots, smoke flows.
- Chrome DevTools: layout/perf/network debugging.

### 7.3 When NOT to use MCP

- Simple local refactors
- Tasks solvable via workspace search + reading existing code
- Anything that would require secrets or sensitive logs

### 7.4 Tool minimization protocol

1. Start with read/search.
2. Enable only the tools you need.
3. Disable entire servers if not relevant.
4. Use tool sets to keep the picker tidy.

Default toolchain by lane:

| Lane | Default tools | Add only if needed |
|---|---|---|
| Discovery | workspace search/read | GitHub MCP (repo context), Context7 (library uncertainty) |
| Planning | search/read + small targeted reads | Context7 (API usage), diagrams only if asked |
| Implementation | edit + execute + focused tests | Playwright/DevTools (UI validation), Context7 (new dep/library) |
| Validation | execute/tests | Playwright (UI smoke), DevTools (perf/layout) |
| Recovery | changes + search/read | none unless blocked |

Tool correctness rules:

- Prefer workspace tools for local codebase truth.
- Use Context7 only to confirm third-party API usage; stop after the minimum needed.
- Use Playwright/DevTools when visual/RTL/layout correctness matters or when UI bugs are hard to reproduce.

---

## 9) Examples (Concrete Scenarios)

Each example shows classification, routing, and a task packet.

### Example 1: Cross-stack feature (backend + portal)

User request:

> "Add a new ‘Low stock’ widget to the dashboard and expose a backend endpoint for it."

Task packet (example):

```
TASK PACKET

Classification:
- Type: feature
- Scope: cross-stack
- Risk: Medium

Owners:
- Primary: Orchestrator
- Supporting: Backend Lead, Web Lead, Backend Implementer, Web Implementer, QA

Lanes:
- Planning → Implementation → Review → Validation

Artifacts:
- Use: existing cross-stack feature skill/prompt (if present)
- Create/update: cross-stack task packet prompt (only if recurring)

Plan (phases):
- Phase 0: define endpoint + response shape + error behavior
- Phase 1: backend endpoint + tests (+ OpenAPI if required)
- Phase 2: portal API client + query + widget UI + i18n keys + states
- Phase 3: Playwright smoke for dashboard + RTL check

Gates:
- Ask PO if schema changes or new deps are required

Tool plan:
- Workspace: search/read/edit + focused tests
- MCP/tools: Playwright optional

Definition of done:
- Widget correct; RTL-safe; no hardcoded strings; relevant tests green
```

### Example 2: UI/UX redesign request

User request:

> "Revamp the Orders list page to feel cleaner and more mobile-friendly."

Task packet (example):

```
TASK PACKET

Classification:
- Type: design
- Scope: single-app (portal-web)
- Risk: High

Owners:
- Primary: Design/UX Lead
- Supporting: Web Lead, i18n Lead, Web Implementer, QA

Lanes:
- Discovery → Planning → Implementation → Review → Validation

Gates:
- PO approval for redesign scope and new UI primitives

Tool plan:
- Workspace: read existing list patterns + components
- MCP/tools: Playwright + Chrome DevTools (as needed)

Definition of done:
- Page is simpler, mobile-first, RTL-safe, consistent components/tokens
```

### Example 3: Large refactor (phased)

User request:

> "Refactor customer search/filter logic across backend and portal; it’s getting messy."

Task packet (example):

```
TASK PACKET

Classification:
- Type: refactor
- Scope: cross-stack
- Risk: High

Owners:
- Primary: Orchestrator
- Supporting: Backend Lead, Web Lead, Shared/Platform Implementer, QA

Lanes:
- Planning (phased) → Implementation → Review → Validation

Plan (phases):
- Phase 0: map current behavior + add/adjust tests
- Phase 1: backend refactor behind compatible API
- Phase 2: portal refactor (reuse query keys/patterns)
- Phase 3: delete dead code + consistency pass

Gates:
- PO approval if behavior changes or contract breaks

Definition of done:
- Behavior preserved (unless approved); tests cover; no dead code
```

### Example 4: Bug report deferred into tech debt

User request:

> "Sometimes the dashboard numbers look wrong. Not urgent."

Task packet (example):

```
TASK PACKET

Classification:
- Type: bug
- Scope: unknown
- Risk: Medium

Owners:
- Primary: Orchestrator
- Supporting: Data/Analytics Lead (if metrics), QA

Lanes:
- Discovery → Deferred/Backlog

Output:
- Repro attempts + what data/logs are needed
- Suspected causes
- Structured backlog item PO can prioritize
```

### Example 5: Translation rewrite

User request:

> "Rewrite the Arabic translations for onboarding screens to sound more natural."

Task packet (example):

```
TASK PACKET

Classification:
- Type: i18n
- Scope: portal-web
- Risk: Medium

Owners:
- Primary: i18n/Localization Lead
- Supporting: Web Lead

Lanes:
- Planning-lite → Implementation → Review

Gates:
- PO approval if meaning changes (not just phrasing)

Definition of done:
- Natural Arabic; keys consistent; no hardcoded strings
```

### Example 6: Content writing (marketing)

User request:

> "Write a short landing page section explaining Kyora’s ‘Cash in hand’ benefit."

Task packet (example):

```
TASK PACKET

Classification:
- Type: content
- Scope: monorepo-wide (copy only)
- Risk: Low

Owners:
- Primary: Content/Marketing Lead
- Supporting: i18n Lead (if Arabic version needed)

Lanes:
- Discovery (tone/claims) → Implementation (copy draft) → Review

Gates:
- PO approval for any claims touching legal/tax promises

Definition of done:
- Simple, calm copy; no accounting jargon; ready for translation
```

---

## 10) Governance and Maintenance

### OS maintenance rules

- Keep this file operational and scannable; don’t stuff SSOT details here.
- Any new always-on rule must live in an SSOT `.instructions.md` file.

### Drift-Sync governance (required)

Drift is any mismatch between:

- How the codebase actually behaves / is implemented today
- What SSOT instructions or this OS say should happen

Common drift examples:

- Translation key casing convention changes
- New form/error-handling pattern becomes the de-facto standard
- New folder placement rules emerge
- UI tokens/components change and older guidance becomes wrong

Rules:

- Do not "quietly diverge" in implementation. If you change a convention, you must sync SSOT.
- Any convention change requires a PO gate (even if code already drifted).

Sync protocol:

1. Identify the SSOT owner file (which `.instructions.md` governs the rule).
2. Prepare a Drift Ticket (Appendix B) with: what changed, why, blast radius, and proposed new rule.
3. Get PO approval if it changes conventions.
4. Update only the SSOT file(s) (do not scatter the same rule elsewhere).
5. Add 1–3 bullets to this OS changelog under Versioning.
6. Run the smallest relevant validation loop (lint/tests/build, or a focused E2E smoke).

Emergency patch path (allowed):

- If a mismatch is actively breaking work today, apply the minimal fix in code, then immediately open a Drift Ticket and schedule SSOT sync in the next cycle.

Rule change ledger (inside this OS):

- Use the Versioning changelog for human-readable notes.
- Use this table when a rule/convention changes (append-only; keep entries short).

| Date | Rule changed | SSOT owner file | PO approved | Validation |
|---|---|---|---:|---|

### Promote repeated work into artifacts

If something happens 3+ times, promote it:

- Repeated single-purpose task → create a prompt
- Repeated multi-step workflow → create a skill
- Repeated persona/tool boundary need → create/refine an agent

### Prompt + skill design checklist (use Appendix C)

When creating/updating prompts or skills, keep them PO-usable, token-efficient, and strongly typed around inputs/outputs. If the artifact would require copying SSOT rules, reference the SSOT instead.

### Prune stale artifacts

Monthly (or per release):

- Remove prompts that aren’t used
- Deprecate skills that no longer match current patterns
- Reduce agent tool lists to avoid tool overload

### Change log discipline

- When Agent OS behavior changes, add a 1–3 bullet entry under the version header.

---

## Appendix A. Handoff Packet Templates

These templates are intentionally verbose in structure (not prose) to survive token limits and reduce "new session drift". Keep bullets tight.

### A1) Delegation Packet (Orchestrator/Lead → Lead/Implementer)

```
DELEGATION PACKET

From:
To:
Date:

Objective (1 sentence):

Classification:
- Type:
- Scope:
- Risk:

Acceptance criteria (copy from PO):
-

Constraints / non-negotiables:
- Tenant isolation (workspace > business)
- Arabic/RTL-first
- Plain money language
- No surprise docs

Gates (PO approvals needed):
-

Key decisions:
-

Assumptions (explicit):
-

Reuse targets (what to reuse / where to look):
-

SSOT references used:
-

Plan (phases + DoD per phase):
- Phase 0:
- Phase 1:

Validation plan (exact commands):
-

Files/areas likely touched:
-

Risks / watch-outs:
-
```

### A2) Phase Handoff Packet (end of a phase)

```
PHASE HANDOFF PACKET

Phase just completed:
Current lane:
Next lane:

What changed (facts only):
-

What’s verified (commands run + result):
-

What remains (next 3–7 steps, ordered):
-

Open questions / pending gates:
-

Cross-stack state snapshot (if applicable):
- Backend contract: stable | changed | pending
- Portal integration: not started | WIP | done
- i18n keys: not started | WIP | done
- E2E/RTL validation: not started | WIP | done
```

### A3) Recovery/Resume Packet (new session continues unfinished work)

```
RECOVERY PACKET

Goal (1 sentence):
Last known lane:

What’s already done (based on git changes/tests):
-

What’s broken / failing (if any):
-

Next smallest verifiable step:
-

Commands to run first:
-

Pending PO gates:
-

Assumptions to re-confirm:
-
```

## Appendix B. Drift Ticket Template

```
DRIFT TICKET

What drifted:
-

Current reality (what code does today):
-

Current SSOT guidance (what instructions say):
-

Proposed new rule:
-

Why (benefit + risk):
-

Blast radius (what areas/files/users affected):
-

PO gate required?: yes | no

Validation plan:
-
```

## Appendix C. Prompt + Skill Design Checklist

Use this when authoring prompts/skills (do not duplicate SSOT rules—link them).

- `description` includes WHAT + WHEN + trigger keywords.
- Tool list is least-privilege; avoid enabling MCP by default.
- Inputs are explicit (goal, scope, constraints, acceptance criteria, risk flags).
- Output is explicit (task packet + phases + validation commands).
- Includes a "Stop-and-ask" clause for high-risk axes.
- Includes a short validation section (exact commands).