---
goal: Implement Kyora Agent OS operational model as enforceable Copilot artifacts (agents, prompts, skills) + lightweight validation automation
version: 1
date_created: 2026-01-22
last_updated: 2026-01-22
owner: Kyora Engineering
status: 'In Progress'
tags: [process, architecture, tooling]
---

# Introduction

![Status: Planned](https://img.shields.io/badge/status-Planned-blue)

This plan implements all requirements in [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md) by turning its operating model into concrete, repo-local Copilot artifacts (agents, prompts, skills) plus a small validation harness (Makefile target + scripts) that checks for required files and metadata. The goal is to make the OS “runnable” and consistently applied across sessions, lanes, and roles.

## 1. Requirements & Constraints

- **REQ-001**: Encode OS roles as custom agents with tool restrictions aligned to [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L201-L347) (Role specs).
- **REQ-002**: Implement mandatory handoff contract and packets: Delegation, Phase Handoff, Recovery/Resume, as first-class outputs in prompts/agents, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L349-L399) and [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1080-L1171).
- **REQ-003**: Implement delegation-by-inference triggers and mandatory reviewers as routing logic in Orchestrator + routing prompt, and require Phase 0 cross-stack contract agreement when Backend + Web both involved, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L360-L395).
- **REQ-004**: Implement routing algorithm artifacts: Minimum Task Brief template, classification, lane selection, owner assignment, artifact selection, tool selection, stop-and-ask rules, stop conditions, multi-phase rules, and Task Packet output contract per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L401-L575).
- **REQ-005**: Implement all lane definitions and default handoffs as reusable “lane playbooks” (skill references + prompt outputs), per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L577-L699).
- **REQ-006**: Implement Recovery Lane behavior (new session continuation) with a dedicated prompt that emits the Recovery Packet before any implementation continues, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L687-L699).
- **REQ-007**: Implement Quality Gates and Checklists as reusable checklists referenced by skills/prompts (not duplicated SSOT), per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L701-L789).
- **REQ-008**: Implement MCP + tooling policy: least-privilege tool lists by lane and safe MCP usage guidance embedded in Orchestrator + prompts/skills, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L791-L879).
- **REQ-009**: Implement governance: drift-sync workflow, drift ticket template, rule change ledger table maintenance, and artifact promotion/pruning process as prompts + skill reference material, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1030-L1079).
- **REQ-010**: Provide PO-executable prompts (not long docs) for repeatable tasks (routing, packets, drift ticket) per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L451-L474) and prompt spec SSOT [prompts.instructions.md](../.github/instructions/prompts.instructions.md).
- **REQ-011**: Implement “Operating Principles” explicitly in Orchestrator + lane skill: Correct+consistent over speed, reuse-first, token discipline, and ambiguity handling (“assumption-first” vs “clarify-first”), per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L10-L48) and [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L421-L444).
- **REQ-012**: Implement stop conditions and search/tool loop limits (3 failures / 3 searches) in Orchestrator + lane skill and include escalation paths, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L520-L545).
- **REQ-013**: Implement example scenarios as runnable reference examples inside skills (not always-on instructions) to improve determinism and reduce session drift, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L881-L1008).
- **REQ-015**: Treat any existing agents/prompts not created by this plan as out-of-scope and DO NOT modify them; OS artifacts must be additive and isolated. Validation automation MUST validate only OS artifacts declared by this plan (not the entire `.github/agents/` / `.github/prompts/` folders).
- **REQ-016**: Prompts intended to be PO-executable MUST use `${input:...}` variables for Minimum Task Brief fields and risk hints, and MUST include a “No surprise docs” clause, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L111-L139) and [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L121-L139).
- **REQ-017**: Every prompt and skill created by this plan MUST satisfy Appendix C (Prompt + Skill Design Checklist) before being considered complete, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1206-L1234).
- **SEC-001**: Enforce tenant isolation as a non-negotiable in all OS artifacts and checklists, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L33-L48) and [.github/copilot-instructions.md](../.github/copilot-instructions.md).
- **SEC-002**: Enforce “stop-and-ask” before auth/RBAC/payments/PII/schema/dependency/breaking-contract work proceeds, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L520-L545).
- **CON-001**: “No surprise docs” gate: do not create/expand docs unless requested; OS artifacts must not encourage long narrative outputs, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L701-L707).
- **CON-002**: Follow SSOT hierarchy: do not duplicate SSOT rules; link to SSOT instruction files instead, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L53-L103).
- **GUD-001**: Use least-privilege tool lists in prompts and agents, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L477-L489) and [agents.instructions.md](../.github/instructions/agents.instructions.md).
- **PAT-001**: Use Appendix A (handoff packet) and Appendix B (drift ticket) templates verbatim as outputs for the relevant prompts/agents, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1080-L1196).
- **PAT-003**: Use Appendix C (Prompt + Skill Design Checklist) as an authoring checklist for all prompts/skills created by this plan, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1206-L1234).
- **PAT-002**: Ensure skills are phase-sliced and require a Phase Handoff Packet between phases (especially for cross-stack workflows), per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L468-L476) and [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L547-L575).

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Establish OS foundations in repo (AGENTS.md + role agents + handoff wiring).

| Task     | Description | Completed | Date |
| -------- | ----------- | --------- | ---- |
| TASK-001 | Create root AGENTS.md that mirrors OS intent and provides validated commands from .github/copilot-instructions.md. Include required sections (Project Overview, Tech Stack, Setup Commands, Project Structure, Code Style, Boundaries) per agents SSOT. File: `AGENTS.md`. Content MUST reference (not duplicate) [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md) and key SSOT entry points. Boundaries MUST include: no secrets, tenant isolation, RTL-first, plain money language, no surprise docs. | ✅ | 2026-01-22 |
| TASK-002 | Add Orchestrator agent file implementing routing + lane enforcement only (no production code edits). File: `.github/agents/orchestrator.agent.md`. Frontmatter tools MUST include `todo` (for lane tracking / packets) and MUST be least-privilege: `['read','search','todo','edit']`, where `edit` is explicitly restricted to plans/notes only (no production code edits), per Orchestrator role spec [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L180-L199). It MUST explicitly forbid editing production code, adding deps, or running destructive commands. Body MUST: (1) require Minimum Task Brief, (2) emit TASK PACKET format per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L547-L575), (3) enforce inference triggers per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L360-L395), (4) enforce stop-and-ask + stop conditions per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L520-L545), (5) output Delegation/Phase/Recovery packets when switching owner/lane per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L349-L399), (6) enforce token discipline/reuse-first per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L39-L48), (7) when scope spans multiple surfaces (e.g., `backend/` + `portal-web/`), require a Phase 0 “Contract Agreement” between the relevant domain leads before any implementation lane begins, and define contract as endpoint, DTO, error semantics, and required i18n copy, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L382-L395) and [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L547-L575). | ✅ | 2026-01-22 |
| TASK-003 | Add Lead agents with OS-aligned responsibilities + tool limits: `.github/agents/backend-lead.agent.md`, `.github/agents/web-lead.agent.md`, `.github/agents/design-ux-lead.agent.md`, `.github/agents/data-analytics-lead.agent.md`, `.github/agents/devops-platform-lead.agent.md`, `.github/agents/i18n-localization-lead.agent.md`, `.github/agents/content-marketing-lead.agent.md`, `.github/agents/security-privacy-reviewer.agent.md`. Each MUST match Role specs in [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L201-L311) including forbidden actions and DoD. Security/Privacy MUST be read/search only. Lead agent bodies MUST also: (1) apply the Delegation-by-inference triggers when reviewing/scoping work (not only the Orchestrator), per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L360-L395), and (2) enforce the cross-stack coordination rule: if Backend + Web both involved, Leads must agree Phase 0 contract BEFORE implementation starts, and the contract includes endpoint, DTO, error semantics, and required i18n copy, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L382-L395). | ✅ | 2026-01-22 |
| TASK-004 | Add Implementer agents: `.github/agents/backend-implementer.agent.md`, `.github/agents/web-implementer.agent.md`, `.github/agents/shared-platform-implementer.agent.md`, `.github/agents/qa-test-specialist.agent.md`. Each MUST: (1) require a TASK PACKET before implementation (unless tiny/low-risk), (2) require a Delegation Packet when switching owners/scopes, and require a Recovery Packet when resuming in a new session, (3) include validation commands section with repo Makefile targets, (4) enforce no dead code + no TODO/FIXME gate, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L701-L739) and repo rules in [.github/copilot-instructions.md](../.github/copilot-instructions.md). | ✅ | 2026-01-22 |
| TASK-005 | Wire agent handoffs from `.github/agents/orchestrator.agent.md` to leads/implementers using VS Code `handoffs:` frontmatter. Handoffs MUST include labels that match lane transitions (e.g., "Start Discovery", "Start Planning", "Start Implementation", "Start Review", "Start Validation", "Resume (Recovery)"). | ✅ | 2026-01-22 |

### Implementation Phase 2

- GOAL-002: Create PO-executable prompts that produce deterministic packets (task intake, routing, handoffs, drift tickets) and reference SSOT.

| Task     | Description | Completed | Date |
| -------- | ----------- | --------- | ---- |
| TASK-006 | Create routing prompt: `.github/prompts/route-task.prompt.md`. Frontmatter MUST include `description`, `agent: 'agent'`, and minimal tools `['codebase','search']` (no `editFiles`). Body MUST: (1) collect Minimum Task Brief fields using `${input:...}` variables (PO-executable) matching [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L405-L419), including risk hints, (2) encode the missing-info policy (“assumption-first” for low-risk missing acceptance criteria; “clarify-first” for ambiguous or medium/high risk) per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L421-L444), (3) include a “No surprise docs” clause, (4) output TASK PACKET exactly as [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L547-L575), (5) apply stop-and-ask + stop conditions per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L520-L545), (6) infer supporting roles per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L360-L395), (7) route by default to the Orchestrator or the relevant Domain Lead (not directly to an Implementer) per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L111-L139), (8) when the task spans multiple surfaces (e.g., backend + web), include a required Phase 0 “Contract Agreement” section in the packet that must be signed off by the relevant domain leads before Phase 1 starts, and explicitly define the contract as endpoint, DTO, error semantics, and required i18n copy, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L382-L395) and [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L547-L575). | ✅ | 2026-01-22 |
| TASK-007 | Create handoff prompts that emit Appendix A formats verbatim: `.github/prompts/create-delegation-packet.prompt.md`, `.github/prompts/create-phase-handoff-packet.prompt.md`, `.github/prompts/create-recovery-packet.prompt.md`. Each MUST be `agent: 'ask'` or `agent: 'agent'` with read/search tools only, MUST include a “No surprise docs” clause, and MUST output only the packet block (no narrative). Templates MUST match [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1080-L1171). | ✅ | 2026-01-22 |
| TASK-008 | Create drift ticket prompt: `.github/prompts/create-drift-ticket.prompt.md`. Must include a “No surprise docs” clause, must output Appendix B Drift Ticket format exactly as [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1173-L1196), and must require an explicit “PO gate required?” answer. | ✅ | 2026-01-22 |
| TASK-009 | Create “artifact promotion” prompt: `.github/prompts/promote-to-artifact.prompt.md` that, given a repeated task description, chooses prompt vs skill vs agent vs instruction using the decision matrix in [ai-artifacts.instructions.md](../.github/instructions/ai-artifacts.instructions.md) and OS guidance [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1056-L1064). Prompt MUST include a “No surprise docs” clause. Output MUST be a mini task packet plus exact file paths to create/update. | ✅ | 2026-01-22 |

### Implementation Phase 3

- GOAL-003: Add core skills that operationalize lanes, quality gates, and cross-stack phase slicing with bundled checklists.

| Task     | Description | Completed | Date |
| -------- | ----------- | --------- | ---- |
| TASK-011 | Create skill folder `.github/skills/lane-playbooks/` with `SKILL.md` describing lane workflows: Discovery, Planning, Implementation, Review, Validation, Recovery. Skill MUST keep body <500 lines and move long checklists into `references/`. References MUST include: `references/lane-discovery.md`, `references/lane-planning.md`, `references/lane-recovery.md` anchored to [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L577-L699). Skill MUST also include (in `SKILL.md`): a compact success checklist and exact validation commands (or explicit "no commands" if lane is non-executable), per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L140-L163). | ✅ | 2026-01-22 |
| TASK-012 | Create skill `.github/skills/cross-stack-feature/` with `SKILL.md` implementing the Phase 0→3 slicing rule per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L547-L575) and example pattern [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L803-L846). Skill MUST require a Phase Handoff Packet after each phase and must include “Never mix new behavior + large refactor in same phase unless PO approves”. Include `references/quality-gates.md` that links (not copies) portal/backend SSOT instruction files named in [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L719-L777). Skill MUST also include: a compact success checklist and exact validation commands per phase, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L140-L163). | ✅ | 2026-01-22 |
| TASK-013 | Create skill `.github/skills/drift-sync/` with `SKILL.md` that runs the drift-sync protocol steps 1–6 per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1046-L1065), and includes `references/rule-change-ledger.md` describing how to append a row to the ledger table in [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1067-L1071). Skill MUST explicitly include the OS “change log discipline” step to add 1–3 bullets under KYORA Agent OS Versioning when behavior changes, per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1028-L1065). | ✅ | 2026-01-22 |
| TASK-014 | Create skill `.github/skills/tool-minimization/` with `SKILL.md` that encodes MCP/tool selection rules and lane-default toolchains per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L791-L879). Include a `references/mcp-safety-checklist.md` for secrets handling and “when NOT to use MCP”, and include the operational note about tool list caching + when to reset cached tools per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L808-L823). Skill MUST also include a compact success checklist and exact validation commands (where applicable), per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L140-L163). | ✅ | 2026-01-22 |
| TASK-014A | Add skills reference for OS examples to reduce drift: `.github/skills/lane-playbooks/references/examples.md` containing short, structured copies of Examples 1–6 from [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L881-L1008) focusing on classification + owners + phases + gates (no narrative). | ✅ | 2026-01-22 |

### Implementation Phase 4

- GOAL-004: Add lightweight automation to validate OS artifacts exist and conform to metadata/tooling constraints.

| Task     | Description | Completed | Date |
| -------- | ----------- | --------- | ---- |
| TASK-015 | Add OS artifact manifest `scripts/agent-os/manifest.json` listing ONLY the OS artifacts created by this plan (agents, prompts, skills, references). This manifest is the source of truth for validation and MUST NOT include existing non-OS artifacts. | ✅ | 2026-01-22 |
| TASK-016 | Add validation script. | ✅ | 2026-01-22 |
| TASK-017 | Add Makefile targets. | ✅ | 2026-01-22 |
| TASK-018 | Add MCP template config. | ✅ | 2026-01-22 |
| TASK-019 | Add monthly artifact prune checklist as a reference doc under `.github/skills/drift-sync/references/artifact-prune.md` implementing [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1066-L1079). This stays inside a skill to avoid “always-on docs”. | ✅ | 2026-01-22 |

## 3. Alternatives

- **ALT-001**: Keep KYORA_AGENT_OS as “guidance only” with no artifacts. Rejected because it does not satisfy routing/tooling/packet determinism requirements and fails in new sessions.
- **ALT-002**: Encode everything into always-on `.instructions.md` files. Rejected because it violates the OS’s SSOT/scanability principle and increases token cost for every chat.
- **ALT-003**: Add external tooling (Node/Yarn packages) to validate artifacts. Rejected because it introduces new dependencies and requires PO gating; Phase 4 uses no-new-deps scripts.

## 4. Dependencies

- **DEP-001**: Existing SSOT instruction files must remain the authoritative rule sources; skills/prompts must link to them: [.github/instructions/*](../.github/instructions/).
- **DEP-002**: Repo validation commands must remain accurate as listed in [.github/copilot-instructions.md](../.github/copilot-instructions.md) (Makefile targets).

## 5. Files

- **FILE-001**: AGENTS.md
- **FILE-002**: .github/agents/orchestrator.agent.md
- **FILE-003**: .github/agents/backend-lead.agent.md
- **FILE-004**: .github/agents/web-lead.agent.md
- **FILE-005**: .github/agents/design-ux-lead.agent.md
- **FILE-006**: .github/agents/data-analytics-lead.agent.md
- **FILE-007**: .github/agents/devops-platform-lead.agent.md
- **FILE-008**: .github/agents/i18n-localization-lead.agent.md
- **FILE-009**: .github/agents/content-marketing-lead.agent.md
- **FILE-010**: .github/agents/security-privacy-reviewer.agent.md
- **FILE-011**: .github/agents/backend-implementer.agent.md
- **FILE-012**: .github/agents/web-implementer.agent.md
- **FILE-013**: .github/agents/shared-platform-implementer.agent.md
- **FILE-014**: .github/agents/qa-test-specialist.agent.md
- **FILE-015**: .github/prompts/route-task.prompt.md
- **FILE-016**: .github/prompts/create-delegation-packet.prompt.md
- **FILE-017**: .github/prompts/create-phase-handoff-packet.prompt.md
- **FILE-018**: .github/prompts/create-recovery-packet.prompt.md
- **FILE-019**: .github/prompts/create-drift-ticket.prompt.md
- **FILE-020**: .github/prompts/promote-to-artifact.prompt.md
- **FILE-021**: .github/skills/lane-playbooks/SKILL.md
- **FILE-022**: .github/skills/lane-playbooks/references/lane-discovery.md
- **FILE-023**: .github/skills/lane-playbooks/references/lane-planning.md
- **FILE-024**: .github/skills/lane-playbooks/references/lane-recovery.md
- **FILE-025**: .github/skills/cross-stack-feature/SKILL.md
- **FILE-026**: .github/skills/cross-stack-feature/references/quality-gates.md
- **FILE-027**: .github/skills/drift-sync/SKILL.md
- **FILE-028**: .github/skills/drift-sync/references/rule-change-ledger.md
- **FILE-029**: .github/skills/drift-sync/references/artifact-prune.md
- **FILE-030**: .github/skills/tool-minimization/SKILL.md
- **FILE-031**: .github/skills/tool-minimization/references/mcp-safety-checklist.md
- **FILE-031A**: .github/skills/lane-playbooks/references/examples.md
- **FILE-032**: scripts/agent-os/manifest.json
- **FILE-033**: scripts/agent-os/validate.sh
- **FILE-034**: Makefile
- **FILE-035**: .vscode/mcp.template.json

## 6. Testing

- **TEST-001**: Validate artifact metadata and presence: `make agent.os.check`.
- **TEST-002**: Validate repo baseline tools remain healthy (precondition): `make doctor`.
- **TEST-003**: Validate backend CI-relevant tests still pass (spot check): `make test.quick`.
- **TEST-004**: Validate portal checks still pass (spot check): `make portal.check`.

## 7. Risks & Assumptions

- **RISK-001**: Too many artifacts reduce discoverability; mitigate by strong `description` fields and pruning cadence per [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md#L1066-L1079).
- **RISK-003**: Tool restrictions may block needed execution; mitigate by keeping Orchestrator limited to read/search/todo (+ edit for plans/notes only) and routing implementation to implementers.
- **RISK-004**: Drift between OS and SSOT instructions; mitigate by enforcing Drift Ticket prompt and adding a `make agent.os.check` requirement in PR checklist (outside scope of this plan unless requested).
- **RISK-005**: Existing non-OS agents/prompts may conflict semantically with the OS; mitigate by OS-only validation and by clearly naming OS artifacts (e.g., `orchestrator.agent.md`, `route-task.prompt.md`) while leaving non-OS artifacts untouched.
- **ASSUMPTION-001**: Team uses VS Code Copilot prompt files and custom agents; if not, prompts/agents still serve as documentation but automated enforcement is reduced.
- **ASSUMPTION-002**: No new dependencies required to validate YAML frontmatter; Phase 4 uses a minimal frontmatter parser.

## 8. Related Specifications / Further Reading

- [KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md)
- [.github/copilot-instructions.md](../.github/copilot-instructions.md)
- [.github/instructions/ai-artifacts.instructions.md](../.github/instructions/ai-artifacts.instructions.md)
- [.github/instructions/agents.instructions.md](../.github/instructions/agents.instructions.md)
- [.github/instructions/prompts.instructions.md](../.github/instructions/prompts.instructions.md)
- [.github/instructions/agent-skills.instructions.md](../.github/instructions/agent-skills.instructions.md)