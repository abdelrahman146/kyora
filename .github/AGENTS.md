# .github/AGENTS.md

## Scope

AI artifacts — agents, prompts, skills, and instruction files for GitHub Copilot customization.

**Parent AGENTS.md**: [../AGENTS.md](../AGENTS.md) (read first for project context)
**Governing OS**: [../KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md) (operating model for agent collaboration)

## Structure

```
.github/
├── agents/                 # Custom Copilot agents (.agent.md)
│   ├── orchestrator.agent.md
│   ├── backend-lead.agent.md
│   ├── backend-implementer.agent.md
│   ├── web-lead.agent.md
│   ├── web-implementer.agent.md
│   └── ...
├── prompts/                # Reusable prompts (.prompt.md)
│   ├── continue-session.prompt.md
│   ├── resolve-drift.prompt.md
│   └── ...
├── skills/                 # Multi-step workflow skills
│   ├── agent-workflows/    # Delegation, handoffs, recovery
│   ├── lane-playbooks/     # Lane-specific guidance
│   └── ...
├── instructions/           # Always-on coding standards (.instructions.md)
│   ├── ai-artifacts.instructions.md      # Artifact selection matrix
│   ├── backend-core.instructions.md      # Backend patterns
│   ├── portal-web-architecture.instructions.md
│   └── ...
└── copilot-instructions.md # Repo-wide baseline
```

## Artifact Selection Matrix

| Need | Create | Location |
|------|--------|----------|
| Coding standards (always-on) | `.instructions.md` | `instructions/` |
| User-triggered task | `.prompt.md` | `prompts/` |
| Specialized persona with tools | `.agent.md` | `agents/` |
| Complex workflow with resources | `SKILL.md` | `skills/<name>/` |
| Project-wide context | `AGENTS.md` | Repo root or app root |

**Key Rule**: Prompts are **PO-only** (user-triggered). Agents cannot invoke prompts. For agent-to-agent workflows, use **skills**.

## Artifact Guidelines

### Instructions (`.instructions.md`)

- **When**: Coding standards that should always apply
- **Activation**: Automatic (via `applyTo` glob pattern)
- **Size**: Keep under 500 lines
- **SSOT**: Never duplicate content across instruction files

```yaml
---
description: "What standards this enforces"
applyTo: "backend/**"  # Glob pattern
---
```

### Prompts (`.prompt.md`)

- **When**: Reusable task for the user (PO) to trigger
- **Activation**: Manual via `/prompt-name` in chat
- **User**: Variables via `${input:varName}`
- **Agents cannot trigger prompts** — prompts are PO-only

```yaml
---
description: "What it does + when to use + keywords"
agent: "agent"
tools: ["codebase", "editFiles"]
---
```

### Agents (`.agent.md`)

- **When**: Specialized persona with tool restrictions
- **Activation**: Manual via `@agent-name` or via `agent` tool (delegation)
- **Delegation**: Add `agent` tool + `infer: true` for agents that delegate

```yaml
---
description: "Role expertise + when to invoke"
name: "Display Name"
tools: ["read", "search", "edit", "agent"]  # Include 'agent' for delegation
infer: true  # Enable autonomous selection
handoffs:
  - label: "Next Step"
    agent: "target-agent"
---
```

### Skills (`SKILL.md`)

- **When**: Complex workflow with bundled resources OR agent-runnable workflow
- **Activation**: Auto-discovered by agents based on description keywords
- **Resources**: Can include `scripts/`, `references/`, `templates/`
- **Agents CAN use skills** — this is how agents access workflow instructions

```yaml
---
name: skill-name
description: "What it does + when to use + keywords (for agent discovery)"
---
```

## Boundaries

### ✅ Always do

- Follow [ai-artifacts.instructions.md](./instructions/ai-artifacts.instructions.md) for artifact selection
- Use `description` field for discoverability (WHAT + WHEN + keywords)
- Reference SSOT files instead of duplicating content
- Version changes to KYORA_AGENT_OS.md via changelog
- Test prompts/agents before committing

### ⚠️ Ask first

- New agent definitions
- Changes to agent tool configurations
- New skills
- Changes to KYORA_AGENT_OS.md routing rules

### 🚫 Never do

- Create prompts for agent-to-agent workflows (use skills instead)
- Give agents more tools than needed (least privilege)
- Duplicate guidance across instruction files
- Create agents without clear role boundaries
- Skip the `description` field (breaks discovery)

## Agent Configuration Requirements

For agents that need to delegate work:

```yaml
tools: ["read", "search", "edit", "execute", "agent"]  # Must include 'agent'
infer: true  # Required for autonomous routing
```

For read-only agents (reviewers, auditors):

```yaml
tools: ["read", "search"]  # No edit, no execute
```

For implementation agents:

```yaml
tools: ["read", "search", "edit", "execute"]  # Full implementation access
```

## SSOT Entry Points

- [ai-artifacts.instructions.md](./instructions/ai-artifacts.instructions.md) — Artifact selection matrix (SSOT)
- [prompts.instructions.md](./instructions/prompts.instructions.md) — Prompt writing guidelines
- [agents.instructions.md](./instructions/agents.instructions.md) — Agent definition guidelines
- [agent-skills.instructions.md](./instructions/agent-skills.instructions.md) — Skill creation guidelines
- [writing-instructions.instructions.md](./instructions/writing-instructions.instructions.md) — Instruction file guidelines

## Agent OS Governance

The Kyora Agent OS ([../KYORA_AGENT_OS.md](../KYORA_AGENT_OS.md)) defines:

- **Lanes**: Who owns what (Backend, Web, Shared, QA, etc.)
- **Gates**: When PO approval is required
- **Routing**: How tasks get assigned to agents
- **Delegation**: How agents hand off work to each other
- **Recovery**: How to handle blocked or failed tasks

**To propose OS changes**: Use `/align-agent-os` prompt with the PO.
