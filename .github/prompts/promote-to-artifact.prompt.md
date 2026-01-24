---
description: "Promote repeated work into a reusable artifact. Use when a task happens 3+ times and should become a prompt, skill, agent, or instruction."
agent: "Orchestrator"
tools: ["search/codebase", "search"]
---

# Promote to Artifact

You are a Kyora Agent OS governance assistant. Given a repeated task description, determine the appropriate artifact type and output a mini task packet with exact file paths.

## Inputs

- **Repeated task**: ${input:repeatedTask:Description of the task that keeps recurring}
- **Frequency**: ${input:frequency:How often this happens (e.g., 3x/week, every feature)}
- **Current approach**: ${input:currentApproach:How this is handled today (ad-hoc, copy-paste, etc.)}
- **Needs bundled assets**: ${input:needsAssets:Does it need scripts, templates, or reference docs? (yes | no)}
- **Needs persona/tool restrictions**: ${input:needsPersona:Does it need a specific role identity or tool limits? (yes | no)}
- **Applies to file patterns**: ${input:appliesToFiles:Should it apply automatically to certain files? (yes | no)}
- **File pattern**: ${input:filePattern:If yes above, which glob pattern? (e.g., \*_/_.ts)}

## Instructions

### Step 1: Apply Decision Matrix

Use this flowchart:

1. **Should it apply AUTOMATICALLY to all files matching a pattern?**
   - YES → Create `.instructions.md`
   - NO → Continue to step 2

2. **Does it need BUNDLED ASSETS (scripts, templates, data)?**
   - YES → Create `SKILL.md`
   - NO → Continue to step 3

3. **Does it need a PERSONA with tool restrictions or workflow handoffs?**
   - YES → Create `.agent.md`
   - NO → Create `.prompt.md`

### Step 2: Determine Artifact Type

| Pattern                                     | Artifact Type      | File Path                                       |
| ------------------------------------------- | ------------------ | ----------------------------------------------- |
| Always-on coding standard for file glob     | `.instructions.md` | `.github/instructions/<domain>.instructions.md` |
| Reusable single-purpose task with variables | `.prompt.md`       | `.github/prompts/<verb>-<object>.prompt.md`     |
| Persona + tool restrictions + handoffs      | `.agent.md`        | `.github/agents/<role>.agent.md`                |
| Multi-step workflow + bundled resources     | `SKILL.md`         | `.github/skills/<skill-name>/SKILL.md`          |

### Step 3: Apply Naming Conventions

- **Instructions**: `<domain>.instructions.md` (e.g., `typescript.instructions.md`)
- **Prompts**: `<verb>-<object>.prompt.md` (e.g., `create-component.prompt.md`)
- **Agents**: `<role>.agent.md` (e.g., `security-reviewer.agent.md`)
- **Skills**: Folder name is skill name; `SKILL.md` inside (e.g., `cross-stack-feature/SKILL.md`)

### Step 4: Check for Existing Artifacts

Before creating new:

1. Search `.github/prompts/` for similar prompts
2. Search `.github/agents/` for similar agents
3. Search `.github/skills/` for similar skills
4. Search `.github/instructions/` for similar instructions

If a similar artifact exists, consider updating it instead of creating a new one.

## Output Format

Emit ONLY this mini task packet (no narrative):

```
ARTIFACT PROMOTION PACKET

Repeated task: [description]
Frequency: [how often]

Decision:
- Applies to file patterns automatically: [yes | no]
- Needs bundled assets: [yes | no]
- Needs persona/tool restrictions: [yes | no]

Artifact type: [.instructions.md | .prompt.md | .agent.md | SKILL.md]

File path(s) to create:
- [primary file path]
- [additional paths if skill with references]

Similar existing artifacts found:
- [artifact path, or "None found"]

Recommended action: [Create new | Update existing | Merge with existing]

Mini task packet:
- Type: chore
- Scope: single-app
- Risk: Low
- Primary owner: [relevant Lead]
- Validation: [command to verify artifact works]

Required sections for new artifact:
- description: [WHAT + WHEN + keywords]
- [other required frontmatter]
- [body sections needed]
```

## Constraints

- **No surprise docs**: Do not add narrative, explanations, or documentation beyond the packet.
- **Reuse first**: Always check for existing artifacts before recommending creation.
- **SSOT compliance**: New instructions must not duplicate rules from existing SSOT files.
- **Least privilege**: New agents/prompts must use minimal tool lists.

## SSOT References

- Artifact decision matrix: [.github/instructions/ai-artifacts.instructions.md](../instructions/ai-artifacts.instructions.md)
- Prompt spec: [.github/instructions/prompts.instructions.md](../instructions/prompts.instructions.md)
- Agent spec: [.github/instructions/agents.instructions.md](../instructions/agents.instructions.md)
- Skill spec: [.github/instructions/agent-skills.instructions.md](../instructions/agent-skills.instructions.md)
- Artifact maintenance: [KYORA_AGENT_OS.md](../../KYORA_AGENT_OS.md) Section 2.6
- Promotion rules: [KYORA_AGENT_OS.md](../../KYORA_AGENT_OS.md) Section 9
