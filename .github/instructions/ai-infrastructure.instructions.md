---
description: "Guidelines for maintaining Kyora Copilot AI infrastructure (.github agents, instructions, prompts, skills)"
applyTo: ".github/**"
---

# Kyora AI Infrastructure Rules (.github)

This file governs changes under `.github/**` only (agents, instructions, prompts, skills). It exists to keep Kyora’s Copilot customization layer consistent, auditable, and in sync with the codebase.

## What Goes Where

- `.github/copilot-instructions.md`
  - Always-on, repo-wide guidance.
  - Keep it broad, short, and stable.
- `.github/instructions/*.instructions.md`
  - Path-specific or file-type-specific guidance.
  - Prefer narrow `applyTo` patterns over global rules.
- `.github/agents/*.agent.md`
  - Personas and workflows with scoped tool access.
  - Use when you want a consistent “mode” (planner, reviewer, AI infra maintainer).
- `.github/prompts/*.prompt.md` (if present)
  - On-demand reusable prompts (invoked via `/...`).
  - Use variables (e.g., `${input:...}`, `${selection}`) to keep prompts generic.
- `.github/skills/<skill-name>/SKILL.md` (if present)
  - Agent Skills: on-demand instructions + optional bundled assets/scripts.
  - Use for complex repeatable workflows requiring extra resources.

## Precedence (when multiple instructions apply)

When writing instructions, avoid conflicts. On GitHub.com, precedence is generally:

1. Path-specific instructions in `.github/instructions/**/NAME.instructions.md`
2. Repo-wide `.github/copilot-instructions.md`
3. Agent instruction files like `AGENTS.md`

In VS Code, multiple instruction sources can be combined and no strict order is guaranteed. Design files to be non-conflicting.

## When To Use Each Primitive

- Use **instructions** for coding standards and repo conventions.
- Use **prompt files** for repeatable tasks with a strict output format or reusable scaffolding.
- Use **custom agents** for reusable roles with stable behavior and tool scoping.
- Use **skills** for progressive disclosure and workflows that may include templates/scripts.

## File Naming and Frontmatter Requirements

### Instructions files (`*.instructions.md`)

- Must live in `.github/instructions/`.
- Must include YAML frontmatter.
- Must include `applyTo` unless it is intended to be attached manually.

Frontmatter fields you may use:

- `description` (recommended)
- `name` (optional)
- `applyTo` (optional; required for auto-application)

Minimal template:

```md
---
description: "What this instructions file does"
applyTo: "**/*.go"
---

# Title

- Rules…
```

### Agent files (`*.agent.md`)

- Must live in `.github/agents/`.
- Must include YAML frontmatter with a non-empty `description`.
- Keep `tools` minimal; avoid over-permission.
- Prefer `infer: false` for specialized agents that should not auto-activate.

Minimal template:

```md
---
name: Agent Name
description: "What this agent does and does not do"
target: vscode
infer: false
tools: ["read", "search", "edit"]
---

# Instructions

...
```

### Prompt files (`*.prompt.md`)

- Must live in `.github/prompts/`.
- Must include YAML frontmatter with a non-empty `description`.
- Use `agent: 'agent'` or a custom agent name when you need a specific mode.

Minimal template:

```md
---
agent: "agent"
description: "What this prompt does"
---

# Prompt

...
```

### Skills (`.github/skills/<skill-name>/SKILL.md`)

Follow the Agent Skills specification:

- `SKILL.md` must contain YAML frontmatter and Markdown body.
- Required fields:
  - `name`: lowercase letters/numbers/hyphens only, 1–64 chars, must match directory name, must not start/end with `-`, must not contain `--`.
  - `description`: 1–1024 chars and must say what it does and when to use it.
- Optional: `license`, `compatibility`, `metadata`, `allowed-tools`.
- Keep `SKILL.md` concise; move deep content to `references/`.

Recommended skill structure:

- `scripts/` executable helpers (document dependencies and safe execution)
- `references/` long-form docs loaded on demand
- `assets/` templates/data

## Prompt Variables and Tool References

When authoring `.prompt.md` files, prefer variables over hardcoding values:

- Selection: `${selection}`, `${selectedText}`
- Workspace: `${workspaceFolder}`, `${workspaceFolderBasename}`
- File context: `${file}`, `${fileBasename}`, `${fileDirname}`, `${fileBasenameNoExtension}`
- Inputs: `${input:variableName}` and `${input:variableName:placeholder}`

When referencing tools inside agent/prompt/instructions bodies, use `#tool:<tool-name>`.

## SSOT and Drift Control

- Prefer **linking** to existing SSOT instruction files instead of copying their content.
- Only write rules that are demonstrably true in the codebase.
- If the codebase is inconsistent, avoid prescribing new standards unless explicitly asked; instead, document the dominant pattern and propose a consolidation plan.

## Security and Safety

- Never add secrets or tokens to instructions, prompts, agents, or skills.
- Treat third-party agents/prompts/skills as untrusted until reviewed.
- If a skill includes scripts, document:
  - What the script does
  - Expected inputs/outputs
  - Dependencies
  - Safe execution notes (and require explicit user approval before running)

## Minimal Review Checklist (for .github changes)

- Every new `*.instructions.md` has correct `applyTo` and narrow scope.
- No duplicated rules across files; SSOT is respected.
- No contradictions with `.github/copilot-instructions.md`.
- No vague directives (“use best practices”, “optimize performance”) without concrete, repo-specific guidance.
- No TODO/FIXME placeholders.
