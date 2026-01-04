# Kyora Custom Agents

VS Code Copilot custom agents for specialized development tasks in the Kyora monorepo.

## Agent Directory

| Agent                  | File                    | Purpose                                                                    | Tools (Key)                            | Target |
| ---------------------- | ----------------------- | -------------------------------------------------------------------------- | -------------------------------------- | ------ |
| **AI Architect**       | `ai-architect.agent.md` | Meta-layer optimization, agent/prompt management, instruction optimization | readFile, editFiles, textSearch, fetch | vscode |
| **Backend Engineer**   | `backend.agent.md`      | Go API development, clean architecture, database design                    | + runInTerminal, runTests              | vscode |
| **Frontend Engineer**  | `frontend.agent.md`     | React TanStack, mobile-first UI, RTL support                               | + openSimpleBrowser                    | vscode |
| **Fullstack Engineer** | `fullstack.agent.md`    | End-to-end features, API contract alignment                                | All tools (full stack)                 | vscode |
| **UI Designer**        | `ui-designer.agent.md`  | Design system, RTL/accessibility, visual language                          | readFile, editFiles (no terminal)      | vscode |

## Usage

### VS Code

1. Open GitHub Copilot Chat
2. Click agent dropdown at bottom of chat view
3. Select agent matching your task:
   - **Backend** work → `@backend-engineer`
   - **Frontend** work → `@frontend-engineer`
   - **Full-stack** work → `@fullstack-engineer`
   - **UI/Design** work → `@ui-designer`
   - **Instructions/Docs** work → `@ai-architect`

### Task Examples

**Backend Engineer**

```
@backend-engineer Add revenue recognition endpoint
@backend-engineer Optimize workspace data query performance
@backend-engineer Implement Stripe webhook handler
```

**Frontend Engineer**

```
@frontend-engineer Build order entry form with validation
@frontend-engineer Add RTL support to dashboard layout
@frontend-engineer Create accessible data table component
```

**Fullstack Engineer**

```
@fullstack-engineer Add product inventory tracking feature
@fullstack-engineer Build customer management system
@fullstack-engineer Implement workspace invitation flow
```

**UI Designer**

```
@ui-designer Design mobile-first order summary card
@ui-designer Audit dashboard accessibility
@ui-designer Create loading state patterns
```

**AI Architect**

```
@ai-architect Optimize backend.instructions.md for token efficiency
@ai-architect Review instruction file consistency
@ai-architect Refactor agent decision tree
```

## Agent Structure

Each agent follows this structure:

```markdown
---
name: Agent Display Name
description: Brief description of agent purpose and capabilities
tools: ["read", "edit", "search", "execute", "grep"]
target: vscode
---

# Agent Name — Specialist Title

## Role

Brief role description

## Technical Expertise

- Bullet points of technical skills

## Coding Standards (Non-Negotiable)

Key principles (KISS, DRY, etc.)

## Domain: Kyora Context

Product/business context

## Definition of Done

Success criteria

## Key References

- Links to relevant instruction files

## Workflow

Step-by-step process

## [Domain-Specific] Principles

Rules specific to this agent's domain
```

## Design Principles

**High Signal/Low Noise**: Every line provides deterministic value. No fluff.

**Single Source of Truth (SSOT)**: Each rule exists in exactly one place. Agents reference, never duplicate.

**KISS (Keep It Simple)**: Simple instructions reduce agent error probability.

**DRY (Don't Repeat Yourself)**: Reference instruction files instead of repeating rules.

**Token Efficiency**: Maximum information in minimum tokens.

**Determinism**: Same instructions → same output every time.

**LLM Readability**: Self-documenting, semantically clear, unambiguous.

## Maintenance

### Adding New Agent

1. Create `<agent-name>.agent.md` in this directory
2. Follow structure above
3. Add YAML frontmatter with required fields
4. Keep description under 100 characters
5. Reference existing instruction files (don't duplicate)
6. Update this README's agent directory table

### Updating Agent

1. Edit agent file directly
2. Test in VS Code Copilot Chat
3. Verify references to instruction files are correct
4. Ensure no broken links
5. Maintain token efficiency (compress, don't expand)

### Deleting Agent

1. Remove agent file
2. Update this README
3. Check for references in other agents
4. Update `.github/copilot-instructions.md` if mentioned

## Related Documentation

- [`.github/copilot-instructions.md`](../copilot-instructions.md) — Agent orchestration layer
- [`.github/instructions/`](../instructions/) — Specialized instruction files
- [VS Code Custom Agents Docs](https://code.visualstudio.com/docs/copilot/customization/custom-agents)
- [VS Code Chat Tools Docs](https://code.visualstudio.com/docs/copilot/chat/chat-tools)
- [GitHub Copilot Agents Docs](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/create-custom-agents)

## Support

Issues with agents? Check:

1. Agent file has valid YAML frontmatter
2. File uses `.agent.md` extension
3. Tools list uses valid tool names
4. References point to existing instruction files
5. VS Code recognizes file in Configure Custom Agents dialog

For questions about agent behavior or instruction optimization, use `@ai-architect`.
