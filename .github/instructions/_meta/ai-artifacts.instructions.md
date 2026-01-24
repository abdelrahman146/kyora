---
description: "GitHub Copilot AI Artifacts — Single Source of Truth for choosing between prompts, agents, skills, and instructions"
applyTo: "**/*.prompt.md,**/*.agent.md,**/*.instructions.md,**/.github/prompts/**,**/.github/agents/**,**/.github/instructions/**,**/.github/skills/**"
---

# GitHub Copilot AI Artifacts — Selection Guide

**SSOT for artifact selection, decision matrix, and orchestration.**

This file is the authoritative reference for determining which artifact type to create. Individual artifact files (prompts.instructions.md, agents.instructions.md, agent-skills.instructions.md) contain implementation details for each type.

---

## 1) Artifact Types Overview

GitHub Copilot supports four customization artifact types:

| Artifact           | Purpose              | Activation        | Scope                  |
| ------------------ | -------------------- | ----------------- | ---------------------- |
| `.instructions.md` | Coding standards     | **Always active** | File patterns          |
| `.prompt.md`       | Reusable tasks       | **On-demand** `/` | Single task            |
| `.agent.md`        | Specialized personas | **On-demand** `@` | Multi-task, role-based |
| `SKILL.md`         | Bundled workflows    | **On-demand**     | Portable, asset-rich   |

### Key Distinctions

| Dimension          | Instructions  | Prompts        | Agents          | Skills              |
| ------------------ | ------------- | -------------- | --------------- | ------------------- |
| **Identity**       | Rules         | Tasks          | Personas        | Workflows           |
| **Activation**     | Automatic     | User-triggered | User-triggered  | User-triggered      |
| **User Input**     | None          | Variables      | Conversation    | Context             |
| **Tool Scope**     | N/A           | Configurable   | Restricted      | Configurable        |
| **Bundled Assets** | None          | None           | None            | Scripts, templates  |
| **Portability**    | VS Code only  | VS Code        | VS Code, GitHub | VS Code, CLI, Agent |
| **Persistence**    | Always loaded | Per-invocation | Per-session     | Per-invocation      |

---

## 2) Master Decision Matrix

### Quick Reference Table

| Your Need                                          | Create This        | File Location            |
| -------------------------------------------------- | ------------------ | ------------------------ |
| Coding standards applied to all matching files     | `.instructions.md` | `.github/instructions/`  |
| Reusable single-purpose task with user variables   | `.prompt.md`       | `.github/prompts/`       |
| Specialized persona with expertise and tool limits | `.agent.md`        | `.github/agents/`        |
| Complex workflow with bundled scripts/templates    | `SKILL.md`         | `.github/skills/<name>/` |
| Project-wide guidance for all agents               | `AGENTS.md`        | Repository root          |

### Decision Flowchart

```
START: I need to customize Copilot behavior
         │
         ▼
┌────────────────────────────────────────────┐
│ Should this apply AUTOMATICALLY to all     │
│ files matching a pattern?                  │
└────────────────────────────────────────────┘
         │
    YES  │  NO
         │    │
         ▼    ▼
   ┌─────────┐  ┌────────────────────────────┐
   │INSTRUCTION│ │ Does it need BUNDLED ASSETS │
   │  .md    │  │ (scripts, templates, data)? │
   └─────────┘  └────────────────────────────┘
                         │
                    YES  │  NO
                         │    │
                         ▼    ▼
                   ┌─────────┐  ┌────────────────────────┐
                   │ SKILL.md│  │ Does it need a PERSONA │
                   └─────────┘  │ with tool restrictions │
                                │ or workflow handoffs?  │
                                └────────────────────────┘
                                         │
                                    YES  │  NO
                                         │    │
                                         ▼    ▼
                                   ┌─────────┐  ┌─────────┐
                                   │ AGENT   │  │ PROMPT  │
                                   │  .md    │  │  .md    │
                                   └─────────┘  └─────────┘
```

---

## 3) Detailed Comparison: When to Choose Each

### Instructions vs Prompts

| Criterion                      | Instructions           | Prompts                 |
| ------------------------------ | ---------------------- | ----------------------- |
| Activation trigger             | Automatic (file match) | Manual (`/prompt-name`) |
| User input needed              | No                     | Yes (variables)         |
| Applies to specific files      | Yes (glob patterns)    | No (workspace-wide)     |
| Coding standards/conventions   | ✅ Use this            | ❌ Wrong choice         |
| One-off task with clear output | ❌ Wrong choice        | ✅ Use this             |

**Rule**: If it's a **standard** → instruction. If it's a **task** → prompt.

### Prompts vs Agents

| Criterion               | Prompts         | Agents                 |
| ----------------------- | --------------- | ---------------------- |
| Identity                | Task definition | Persona with expertise |
| Invocation              | `/prompt-name`  | `@agent-name`          |
| Tool restrictions       | Optional        | Core feature           |
| Workflow handoffs       | Not supported   | Built-in (`handoffs:`) |
| Multi-request expertise | No (stateless)  | Yes (session-based)    |
| Sub-agent orchestration | No              | Yes (`agent` tool)     |

**Rule**: If it needs a **persona** or **restricted tools** → agent. If it's a **single task** → prompt.

### Prompts vs Skills

| Criterion          | Prompts           | Skills                     |
| ------------------ | ----------------- | -------------------------- |
| Bundled assets     | None              | Scripts, templates, data   |
| Complexity         | Single file       | Multi-file structure       |
| Portability        | VS Code only      | CLI, VS Code, coding agent |
| External resources | Inline only       | `references/`, `scripts/`  |
| Execution          | Instructions only | Can run scripts            |

**Rule**: If it needs **bundled files** → skill. If it's **self-contained** → prompt.

### Agents vs Skills

| Criterion         | Agents          | Skills              |
| ----------------- | --------------- | ------------------- |
| Identity          | Persona         | Workflow            |
| Tool restrictions | Core feature    | Not applicable      |
| Bundled assets    | None            | Scripts, templates  |
| Handoffs          | Yes             | No                  |
| Cross-platform    | VS Code, GitHub | VS Code, CLI, Agent |

**Rule**: If it needs **persona + tool limits** → agent. If it needs **bundled resources** → skill.

### Instructions vs Skills

| Criterion      | Instructions           | Skills            |
| -------------- | ---------------------- | ----------------- |
| Activation     | Automatic              | On-demand         |
| Bundled assets | None                   | Yes               |
| Scope          | File patterns          | Task-specific     |
| Size limit     | ~200 lines practical   | Unlimited         |
| Use case       | Standards, conventions | Complex workflows |

**Rule**: If it's **always-on standards** → instruction. If it's **on-demand workflow with assets** → skill.

---

## 4) Artifact Anatomy Quick Reference

### Instructions (`.instructions.md`)

```yaml
---
description: "What standards this enforces"
applyTo: "**/*.ts,**/*.tsx" # Glob patterns
---
# Title

[Coding standards, patterns, conventions]
```

**Location**: `.github/instructions/`

### Prompts (`.prompt.md`)

```yaml
---
description: "What task + when to use"
agent: "agent"                         # ask | edit | agent | <custom>
tools: ["codebase", "editFiles"]
---

# Title

[Task definition with ${input:variables}]
```

**Location**: `.github/prompts/`

### Agents (`.agent.md`)

```yaml
---
description: "Persona expertise"
name: "Display Name"
tools: ["read", "search"] # Restricted tools
handoffs:
  - label: "Next Step"
    agent: other-agent
---
# Persona Name

[Identity, responsibilities, constraints]
```

**Location**: `.github/agents/`

### Skills (`SKILL.md`)

```yaml
---
description: "Workflow purpose"
version: "1.0.0"
triggers: ["keyword", "phrase"]
---
# Skill Name

[Workflow instructions referencing bundled assets]
```

**Location**: `.github/skills/<skill-name>/SKILL.md`

### Project Guidance (`AGENTS.md`)

```markdown
# AGENTS.md

## Project Overview

## Tech Stack

## Setup Commands

## Project Structure

## Code Style

## Boundaries
```

**Location**: Repository root (or nested for monorepos)

---

## 5) Common Scenarios Mapped

| Scenario                                          | Artifact     | Why                                   |
| ------------------------------------------------- | ------------ | ------------------------------------- |
| TypeScript naming conventions for all `.ts` files | Instructions | Always-on, pattern-based              |
| "Create a React component" task                   | Prompt       | On-demand, user inputs component name |
| Security auditor that only reads code             | Agent        | Persona + tool restriction (no edit)  |
| API integration with SDK templates                | Skill        | Needs bundled templates and scripts   |
| "Document how to build this project"              | AGENTS.md    | Project-wide guidance for all agents  |
| Code review guidelines                            | Instructions | Standards, always apply to reviews    |
| Generate unit tests for a function                | Prompt       | Single task, user selects function    |
| Planning agent that hands off to implementer      | Agent        | Multi-step workflow with handoffs     |
| Database migration workflow with SQL templates    | Skill        | Bundled SQL scripts, multi-step       |

---

## 6) Anti-Patterns

### ❌ Wrong Artifact Choices

| Mistake                                    | Why Wrong                             | Correct Choice |
| ------------------------------------------ | ------------------------------------- | -------------- |
| Prompt for "follow TypeScript conventions" | Should always apply, not on-demand    | Instructions   |
| Instruction for "create a component"       | Needs user input, task-based          | Prompt         |
| Agent for simple README generation         | No persona/tool restrictions needed   | Prompt         |
| Prompt with 500+ lines of templates        | Too complex, needs bundled assets     | Skill          |
| Skill for coding standards                 | Should always apply automatically     | Instructions   |
| Agent without tool restrictions            | No benefit over prompt, adds overhead | Prompt         |

### ❌ Duplication Anti-Patterns

| Anti-Pattern                          | Fix                                   |
| ------------------------------------- | ------------------------------------- |
| Same decision logic in multiple files | Reference this SSOT                   |
| Prompt + Instruction for same concern | Pick one based on activation model    |
| Agent + Skill for same workflow       | Agent if persona needed, Skill if not |
| Multiple prompts doing similar tasks  | Parameterize with variables           |

---

## 7) Artifact Composition Patterns

### Layered System (Recommended)

```
AGENTS.md                    # Project-wide context (all agents read this)
    │
    ├── .instructions.md     # Always-on standards (auto-applied)
    │
    ├── .agent.md            # Personas invoke prompts/skills
    │       │
    │       └── .prompt.md   # Agents can suggest prompts
    │
    └── SKILL.md             # Complex workflows with assets
```

### Composition Rules

1. **AGENTS.md** provides **context** all artifacts can reference
2. **Instructions** set **baselines** that prompts/agents must follow
3. **Agents** can **invoke** prompts as sub-tasks
4. **Skills** can be **referenced** by agents for complex workflows
5. **Prompts** are **atomic** — one task, one output

### Harmony Example

```
.github/
├── AGENTS.md                           # "Use Go 1.22, follow DDD"
├── instructions/
│   └── go-backend.instructions.md      # Go coding standards (always-on)
├── agents/
│   ├── planner.agent.md                # Plans features (read-only)
│   └── implementer.agent.md            # Implements plans (full tools)
├── prompts/
│   ├── create-handler.prompt.md        # Scaffold a new handler
│   └── generate-tests.prompt.md        # Generate tests for selection
└── skills/
    └── api-endpoint/
        ├── SKILL.md                    # Full endpoint workflow
        ├── templates/handler.go.tmpl
        └── scripts/generate.sh
```

**Flow**:

1. Developer asks `@planner` to design a feature
2. Planner reads AGENTS.md, applies go-backend.instructions.md
3. Planner hands off to `@implementer`
4. Implementer uses `/create-handler` prompt for scaffolding
5. Complex workflows use `api-endpoint` skill with templates

---

## 8) Governance Guidelines

### Naming Conventions

| Artifact     | Naming Pattern             | Example                      |
| ------------ | -------------------------- | ---------------------------- |
| Instructions | `<domain>.instructions.md` | `go-backend.instructions.md` |
| Prompts      | `<verb>-<noun>.prompt.md`  | `create-component.prompt.md` |
| Agents       | `<role>.agent.md`          | `security-auditor.agent.md`  |
| Skills       | `<workflow-name>/SKILL.md` | `api-endpoint/SKILL.md`      |

### Description Requirements

All artifacts MUST have a `description` in frontmatter:

| Artifact     | Description Formula                           |
| ------------ | --------------------------------------------- |
| Instructions | "Coding standards for [domain/pattern]"       |
| Prompts      | "[What it does] + [When to use] + [Keywords]" |
| Agents       | "[Role/expertise] for [domain/tasks]"         |
| Skills       | "[Workflow purpose] with [key capabilities]"  |

### Review Checklist

Before creating any artifact:

- [ ] Consulted this decision matrix
- [ ] Verified no existing artifact serves the need
- [ ] Chose correct artifact type based on flowchart
- [ ] Placed in correct location
- [ ] Wrote actionable description
- [ ] Tested in isolation

---

## 9) Cross-References

For implementation details of each artifact type:

| Artifact     | Implementation Guide                                                           |
| ------------ | ------------------------------------------------------------------------------ |
| Instructions | [writing-instructions.instructions.md](./writing-instructions.instructions.md) |
| Prompts      | [prompts.instructions.md](./prompts.instructions.md)                           |
| Agents       | [agents.instructions.md](./agents.instructions.md)                             |
| Skills       | [agent-skills.instructions.md](./agent-skills.instructions.md)                 |

**This file (ai-artifacts.instructions.md)**: Selection and decision logic only.  
**Specialized files**: Implementation details, frontmatter specs, examples.
