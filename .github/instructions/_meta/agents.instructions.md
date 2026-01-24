---
description: "Guidelines for creating strong, optimized GitHub Copilot custom agents and AGENTS.md files"
applyTo: "**/*.agent.md,**/AGENTS.md,**/.github/agents/**"
---

# GitHub Copilot Custom Agents — Complete Guidelines

Rules for creating effective custom agents (`.agent.md` files) and `AGENTS.md` project guidance files.

**Artifact Selection**: See [ai-artifacts.instructions.md](./ai-artifacts.instructions.md) for decision matrix on when to create agents vs prompts vs skills vs instructions.

---

## 1) AGENTS.md Specification

`AGENTS.md` is the README for coding agents. It provides project context agents need to work effectively.

### Location and Discovery

```
# Single project
project/
  AGENTS.md              # Root guidance

# Monorepo (nested overrides closer files)
monorepo/
  AGENTS.md              # Base guidance
  packages/
    api/
      AGENTS.md          # API-specific overrides
    web/
      AGENTS.md          # Web-specific overrides
```

Agents read the **nearest** `AGENTS.md` in the directory tree. Files closer to the current directory override earlier guidance.

### Required Sections (Cover These Six Areas)

```markdown
# AGENTS.md

## Project Overview

[2-3 sentences: what this project does, who it's for]

## Tech Stack

- **Language**: [Go 1.22+ / TypeScript 5.5+ / etc.]
- **Framework**: [Gin, React, etc. with versions]
- **Database**: [PostgreSQL 16, etc.]
- **Key Dependencies**: [List critical packages]

## Setup Commands

\`\`\`bash

# Install dependencies

npm ci

# Build

npm run build

# Test

npm test

# Lint

npm run lint
\`\`\`

## Project Structure

\`\`\`
src/
api/ # REST endpoints
services/ # Business logic
models/ # Data models
tests/ # Test files
\`\`\`

## Code Style

[Include 1-2 concrete examples showing correct patterns]

\`\`\`typescript
// ✅ Good: Descriptive names, proper error handling
async function fetchUserById(id: string): Promise<User> {
if (!id) throw new Error('User ID required');
return await api.get(`/users/${id}`);
}
\`\`\`

## Boundaries

- ✅ **Always**: Run tests before commits, follow naming conventions
- ⚠️ **Ask first**: Database schema changes, new dependencies
- 🚫 **Never**: Commit secrets, modify `node_modules/`, edit production configs
```

### AGENTS.md Best Practices

| ✅ Do                                       | ❌ Don't                      |
| ------------------------------------------- | ----------------------------- |
| Put executable commands early               | Bury commands in prose        |
| Include exact flags: `pytest -v --cov`      | Just tool names: "run pytest" |
| Show code examples of style                 | Describe style abstractly     |
| State versions: "React 18 + TypeScript 5.5" | "Modern React stack"          |
| Set explicit boundaries                     | Assume agent knows limits     |
| Keep under 32KB total                       | Write novel-length docs       |

### Three-Tier Boundaries (Required)

Every AGENTS.md must define boundaries in three tiers:

```markdown
## Boundaries

- ✅ **Always do**: [Safe actions agent should take without asking]
- ⚠️ **Ask first**: [Actions requiring user confirmation]
- 🚫 **Never do**: [Forbidden actions]
```

---

## 2) Custom Agent Files (`.agent.md`)

Custom agents are specialized Copilot personas with defined expertise, tools, and boundaries.

### File Location

| Scope              | Location                                 |
| ------------------ | ---------------------------------------- |
| Repository-level   | `.github/agents/<name>.agent.md`         |
| Organization-level | `.github-private/agents/<name>.agent.md` |
| Enterprise-level   | `.github-private/agents/<name>.agent.md` |

### Required Frontmatter

```yaml
---
description: "Expert test engineer for comprehensive test coverage" # REQUIRED
name: "Test Specialist" # Optional
tools: ["read", "edit", "search", "execute"] # Optional
model: "Claude Sonnet 4" # Recommended
target: "vscode" # Optional
infer: true # Optional
---
```

### Frontmatter Properties Reference

| Property      | Type    | Required    | Description                                    |
| ------------- | ------- | ----------- | ---------------------------------------------- |
| `description` | string  | **Yes**     | Agent purpose (50-150 chars, single-quoted)    |
| `name`        | string  | No          | Display name (defaults to filename)            |
| `tools`       | list    | No          | Available tools (defaults to all if omitted)   |
| `model`       | string  | Recommended | AI model: `'Claude Sonnet 4'`, `'gpt-4o'`      |
| `target`      | string  | No          | `'vscode'` or `'github-copilot'` or both       |
| `infer`       | boolean | No          | Auto-select based on context (default: `true`) |
| `handoffs`    | list    | No          | Workflow transitions (VS Code only)            |
| `mcp-servers` | object  | No          | MCP servers (org/enterprise only)              |

### Tool Configuration

**Tool Aliases (case-insensitive):**

| Alias     | Alternatives                   | Purpose                  |
| --------- | ------------------------------ | ------------------------ |
| `execute` | shell, bash, powershell        | Run terminal commands    |
| `read`    | view, NotebookRead             | Read file contents       |
| `edit`    | Write, MultiEdit, NotebookEdit | Modify files             |
| `search`  | Grep, Glob                     | Search codebase          |
| `agent`   | custom-agent, Task             | Invoke other agents      |
| `web`     | WebSearch, WebFetch            | Fetch web content        |
| `todo`    | TodoWrite                      | Task list (VS Code only) |

**Tool Strategies:**

```yaml
# Read-only agent (planning, review)
tools: ['read', 'search', 'web']

# Full implementation agent
tools: ['read', 'edit', 'search', 'execute']

# MCP server tools
tools: ['read', 'edit', 'github/*', 'playwright/navigate']

# All tools (default if omitted)
tools: ['*']

# No tools (pure chat)
tools: []
```

**Principle of Least Privilege**: Only enable tools necessary for the agent's purpose. Fewer tools = clearer agent behavior.

### Agent Prompt Structure

The markdown body below frontmatter defines agent behavior. Include:

1. **Identity and Role**: Who the agent is
2. **Core Responsibilities**: What specific tasks it performs
3. **Approach/Methodology**: How it accomplishes tasks
4. **Guidelines and Constraints**: What to do/avoid
5. **Output Expectations**: Expected format and quality

```markdown
---
description: "Security auditor for vulnerability detection"
name: "Security Auditor"
tools: ["read", "search", "web"]
---

# Security Auditor

You are a security specialist focused on identifying vulnerabilities and security issues.

## Your Role

- Analyze code for security vulnerabilities
- Check against OWASP Top 10
- Review authentication and authorization logic
- Identify injection risks and data exposure

## Approach

1. Scan for common vulnerability patterns
2. Review input validation and sanitization
3. Check secrets management
4. Analyze access control logic
5. Document findings with severity ratings

## Guidelines

- Focus on security issues only, not code style
- Provide actionable remediation steps
- Rate severity: Critical, High, Medium, Low
- Never modify code directly - only report findings

## Output Format

Present findings as:
\`\`\`markdown

## [SEVERITY] Finding Title

**Location**: file:line
**Issue**: Description
**Risk**: Impact explanation
**Fix**: Remediation steps
\`\`\`
```

### Prompt Writing Rules

| ✅ Do                                 | ❌ Don't                   |
| ------------------------------------- | -------------------------- |
| Use imperative: "Analyze", "Generate" | "You should consider"      |
| Define clear boundaries               | Leave scope ambiguous      |
| Include output format examples        | Describe format abstractly |
| State what to avoid                   | Assume agent knows limits  |
| Keep under 30,000 characters          | Write endless instructions |

---

## 3) Handoffs (VS Code Only)

Handoffs create guided workflows that transition between agents.

### Configuration

```yaml
---
description: "Generate implementation plans"
name: "Planner"
tools: ["read", "search"]
handoffs:
  - label: Start Implementation # Button text
    agent: implementer # Target agent
    prompt: "Implement the plan above." # Pre-filled prompt
    send: false # Manual send (true = auto-submit)
---
```

### Common Handoff Patterns

| Flow                          | Use Case                          |
| ----------------------------- | --------------------------------- |
| Planning → Implementation     | Plan first, then code             |
| Implementation → Review       | Code, then quality check          |
| Review → Planning             | Issues found, re-plan             |
| Failing Tests → Passing Tests | Write tests first, then implement |

### Handoff Best Practices

- **Clear labels**: "Start Implementation" not "Next"
- **Context-aware prompts**: Reference completed work
- **Limit to 2-3**: Most relevant next steps only
- **Verify targets exist**: Handoffs to missing agents are silently ignored

---

## 4) Sub-Agent Orchestration

Agents can invoke other agents using the `agent` tool for multi-step workflows.

### Enable Orchestration

```yaml
tools: ["read", "edit", "search", "agent"] # Include 'agent' tool
```

### Invocation Pattern

```text
This phase must be performed as the agent "data-processor" defined in ".github/agents/data-processor.agent.md".

IMPORTANT:
- Read and apply the entire .agent.md spec.
- Project: "${projectName}"
- Base path: "${basePath}"

Task:
1. Process input from ${basePath}/input/
2. Write results to ${basePath}/output/
3. Return summary of actions taken.
```

### Orchestrator Structure

Document these elements:

- **Dynamic parameters**: Values extracted from user (`projectName`, `basePath`)
- **Sub-agent registry**: Step → agent mapping
- **Step ordering**: Explicit sequence
- **Trigger conditions**: When steps run/skip
- **Logging strategy**: Single log file updated per step

### Limitations

- **NOT for large-scale processing**: Each invocation adds latency
- **Max ~5-10 steps**: Beyond this, implement directly in one agent
- **Tool ceiling**: Sub-agents can only use tools available to orchestrator

---

## 5) Variables in Agents

Use template variables for dynamic, context-aware agents.

### Declaration Pattern

```markdown
## Dynamic Parameters

- **projectName**: Name of the project (string, required)
- **basePath**: Root directory (path, required)
- **outputDir**: Defaults to ${basePath}/output

## Your Mission

Process the **${projectName}** project at `${basePath}`.
```

### Extraction Methods

1. **Explicit**: Ask user if not provided
2. **Implicit**: Extract from user prompt
3. **Contextual**: Derive from file/workspace context

### Variable Best Practices

- Document all expected variables
- Use consistent naming (`projectName` not `name`)
- Specify types and constraints
- Provide defaults where sensible

---

## 6) Common Agent Patterns

### Testing Specialist

```yaml
description: "Test coverage and quality"
tools: ["*"] # Needs full access to write tests
```

Focus: Write tests, identify gaps, never modify production code.

### Implementation Planner

```yaml
description: "Technical planning and specifications"
tools: ["read", "search", "edit"] # No execute
```

Focus: Create plans, not implementations.

### Code Reviewer

```yaml
description: "Code quality analysis"
tools: ["read", "search"] # Read-only
```

Focus: Analyze and suggest, never modify.

### Security Auditor

```yaml
description: "Security vulnerability detection"
tools: ["read", "search", "web"]
```

Focus: Find vulnerabilities, report findings.

### Documentation Writer

```yaml
description: "Technical documentation"
tools: ["read", "search", "edit"]
```

Focus: Read code, write docs. Never modify source.

---

## 7) File Organization

### Naming Conventions

- **Format**: `lowercase-with-hyphens.agent.md`
- **Characters**: `.`, `-`, `_`, `a-z`, `A-Z`, `0-9` only
- **Purpose**: Name should reflect agent purpose
- **Examples**: `test-specialist.agent.md`, `security-auditor.agent.md`

### Directory Structure

```
.github/
  agents/                    # Repository agents
    test-specialist.agent.md
    code-reviewer.agent.md
  instructions/              # Coding standards
    typescript.instructions.md
  prompts/                   # Task prompts
    create-component.prompt.md
  skills/                    # Skills with assets
    playwright-testing/
      SKILL.md
      scripts/helper.py
AGENTS.md                    # Project guidance
```

---

## 8) Content Placement: Agent vs Instruction vs AGENTS.md

For the full artifact selection matrix, see [ai-artifacts.instructions.md](./ai-artifacts.instructions.md). This section covers agent-specific content placement.

### Content Placement Rules

| Content Type                             | Goes In             |
| ---------------------------------------- | ------------------- |
| Build/test commands                      | `AGENTS.md`         |
| Project structure overview               | `AGENTS.md`         |
| Global boundaries (never commit secrets) | `AGENTS.md`         |
| Tech stack versions                      | `AGENTS.md`         |
| Specialized persona definition           | `*.agent.md`        |
| Tool restrictions for safety             | `*.agent.md`        |
| Workflow handoffs                        | `*.agent.md`        |
| Role-specific expertise                  | `*.agent.md`        |
| Coding standards per file type           | `*.instructions.md` |
| Language conventions                     | `*.instructions.md` |
| Naming patterns                          | `*.instructions.md` |
| Framework-specific rules                 | `*.instructions.md` |

### Key Distinction

- **AGENTS.md**: Context about the project (what exists, how to build)
- **Custom agents**: Persona with restricted capabilities (who to be, what tools)
- **Instructions**: Rules for code patterns (how to write specific file types)

---

## 9) Common Mistakes to Avoid

### Frontmatter Errors

- ❌ Missing `description` field
- ❌ Description not in quotes
- ❌ Invalid YAML syntax
- ❌ Invalid tool names

### Tool Configuration

- ❌ Granting excessive tools unnecessarily
- ❌ Missing required tools for purpose
- ❌ Forgetting MCP namespace: `server-name/tool`

### Prompt Problems

- ❌ Vague instructions: "Write good code"
- ❌ Conflicting guidelines
- ❌ No scope definition
- ❌ Missing output expectations
- ❌ Over 30,000 characters

### Organizational

- ❌ Filename doesn't reflect purpose
- ❌ Wrong directory (repo vs org)
- ❌ Spaces in filename
- ❌ Duplicate names

---

## 10) Validation Checklist

### Agent File Checklist

**Frontmatter:**

- [ ] `description` present and descriptive (50-150 chars)
- [ ] `description` wrapped in single quotes
- [ ] `tools` configured appropriately
- [ ] `model` specified (recommended)

**Prompt Content:**

- [ ] Clear identity and role
- [ ] Core responsibilities listed
- [ ] Guidelines and constraints specified
- [ ] Output format documented
- [ ] Under 30,000 characters

**File Structure:**

- [ ] Lowercase-with-hyphens filename
- [ ] Correct directory (`.github/agents/`)
- [ ] `.agent.md` extension

### AGENTS.md Checklist

- [ ] Project overview (2-3 sentences)
- [ ] Tech stack with versions
- [ ] Setup/build/test commands with flags
- [ ] Project structure description
- [ ] Code style examples (not descriptions)
- [ ] Three-tier boundaries (Always/Ask/Never)
- [ ] Under 32KB total

---

## 11) Version Compatibility

| Feature                | GitHub.com   | VS Code | JetBrains/Eclipse/Xcode |
| ---------------------- | ------------ | ------- | ----------------------- |
| Standard frontmatter   | ✅           | ✅      | ✅                      |
| `model` property       | ❌           | ✅      | ✅                      |
| `handoffs` property    | ❌           | ✅      | ⚠️ Limited              |
| `mcp-servers` in agent | Org/Ent only | ❌      | ❌                      |
| Repository MCP         | ✅           | ✅      | ⚠️                      |

Use `target` property for environment-specific agents when needed.
