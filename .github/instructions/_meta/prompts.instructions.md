---
description: "Guidelines for creating strong, optimized GitHub Copilot prompt files (.prompt.md)"
applyTo: "**/*.prompt.md,**/.github/prompts/**"
---

# GitHub Copilot Prompt Files — Complete Guidelines

Rules for creating effective, reusable prompt files that define specific tasks with clear inputs, outputs, and workflows for GitHub Copilot.

**Artifact Selection**: See [ai-artifacts.instructions.md](./ai-artifacts.instructions.md) for decision matrix on when to create prompts vs agents vs skills vs instructions.

---

## ⚠️ Critical: Prompts Are User-Only (PO-Only)

**Prompts are triggered by the user (PO) via `/prompt-name` in chat. Agents CANNOT invoke prompts.**

This is a fundamental capability limitation:

- **Prompts**: User-triggered tasks with input variables
- **Skills**: Agent-accessible workflows with bundled resources
- **Agent tool**: How agents delegate to other agents

### Implications for Workflow Design

| Need                              | Use                       |
| --------------------------------- | ------------------------- |
| PO triggers a task                | Prompt (`.prompt.md`)     |
| Agent needs workflow instructions | Skill (`SKILL.md`)        |
| Agent delegates to another agent  | `agent` tool with handoff |

**Common Mistake**: Creating prompts for agent workflows. If an agent needs to execute a workflow (like creating a delegation packet), that workflow MUST be in a skill, not a prompt.

---

## 1) Prompt File Structure

### Required Location

```
.github/prompts/                    # Recommended
  create-component.prompt.md
  generate-tests.prompt.md
  review-code.prompt.md
```

Alternative locations (configurable in VS Code):

- User profile prompts (cross-workspace)
- Custom workspace folders via `chat.promptFilesLocations` setting

### Required Format

```markdown
---
description: "Clear description of WHAT it does and WHEN to use it"
agent: "agent"
tools: ["codebase", "editFiles", "search"]
model: "Claude Sonnet 4"
---

# Prompt Title

[Persona and context]

## Task

[Clear task description]

## Instructions

[Step-by-step workflow]

## Input

[Variable handling and context requirements]

## Output

[Expected format and structure]
```

---

## 2) Frontmatter Specification

### Required and Recommended Fields

| Field           | Required    | Description                                           |
| --------------- | ----------- | ----------------------------------------------------- |
| `description`   | Recommended | 1-500 chars describing WHAT + WHEN (discovery key)    |
| `name`          | Optional    | Display name after `/` in chat (defaults to filename) |
| `agent`         | Recommended | `'ask'`, `'edit'`, `'agent'`, or custom agent name    |
| `model`         | Optional    | AI model to use (defaults to current selection)       |
| `tools`         | Optional    | List of available tools for this prompt               |
| `argument-hint` | Optional    | Hint text shown in chat input field                   |

### Agent Values

| Agent    | Use Case                                      |
| -------- | --------------------------------------------- |
| `ask`    | Questions, explanations, analysis only        |
| `edit`   | Code modifications with edit capabilities     |
| `agent`  | Full agent mode with tools (default if tools) |
| `<name>` | Custom agent defined in `.agent.md`           |

### Tool Categories

**File Operations:**

- `codebase` / `search/codebase` — Search code semantically
- `editFiles` / `edit/editFiles` — Modify files
- `search` — Find files and text (grep, glob)
- `problems` — Get diagnostics/errors

**Execution:**

- `runCommands` — Run terminal commands
- `runTasks` — Execute VS Code tasks
- `runTests` — Run tests
- `terminalLastCommand` — Get last terminal command

**External:**

- `fetch` / `web/fetch` — Fetch web content
- `githubRepo` — Search GitHub repositories
- `openSimpleBrowser` — Open browser preview

**Specialized:**

- `playwright/*` — Browser automation
- `github/*` — GitHub MCP tools
- `usages` — Find code usages
- `vscodeAPI` — VS Code extension API
- `extensions` — Search extensions
- `changes` — Git changes

**Analysis:**

- `findTestFiles` — Locate test files
- `testFailure` — Get test failure info
- `searchResults` — Access search results

### Tool Selection Principles

1. **Principle of Least Privilege**: Only include tools necessary for the task
2. **If tools specified**: Agent mode is auto-enabled if current mode is `ask` or `edit`
3. **Tool order**: List in preferred execution sequence when order matters
4. **Unknown tools**: Silently ignored (enables environment-specific tools)

---

## 3) Description: The Discovery Mechanism

**CRITICAL**: The `description` is how Copilot discovers prompts. Users see only `name` and `description` when typing `/` in chat.

### Formula for Strong Descriptions

```
[WHAT it does] + [WHEN to use it] + [Keywords/triggers]
```

### ✅ Good Description Examples

```yaml
description: 'Generate comprehensive README.md files for projects. Use when creating documentation, setting up new repos, or improving project discoverability.'

description: 'Create unit tests with full coverage. Use when writing tests for functions, classes, or modules. Supports Jest, Pytest, xUnit frameworks.'

description: 'Review code for security vulnerabilities. Use for security audits, pre-deployment checks, or OWASP compliance verification.'
```

### ❌ Bad Description Examples

```yaml
description: 'Creates READMEs'  # Too vague, no triggers
description: 'Test generation'  # No capabilities, no when
description: 'Code review tool' # No specifics, no keywords
```

---

## 4) Body Structure

### Recommended Sections

| Section           | Purpose                                       | Required    |
| ----------------- | --------------------------------------------- | ----------- |
| `# Title`         | Matches prompt intent, surfaces in Quick Pick | Yes         |
| `## Task`         | Clear task description with requirements      | Yes         |
| `## Instructions` | Step-by-step workflow                         | Yes         |
| `## Input`        | Variable handling, context requirements       | Recommended |
| `## Output`       | Expected format, location, structure          | Recommended |
| `## Examples`     | Good/bad examples for guidance                | Recommended |
| `## Validation`   | Success criteria, verification steps          | Optional    |

### Logical Flow

Structure prompts following this logical progression:

```
WHY → CONTEXT → INPUTS → ACTIONS → OUTPUTS → VALIDATION
```

1. **Why**: Brief persona or purpose statement
2. **Context**: Scope, preconditions, constraints
3. **Inputs**: What the prompt needs from user/workspace
4. **Actions**: Step-by-step workflow
5. **Outputs**: Expected results and format
6. **Validation**: How to verify success

---

## 5) Variable System

### Built-in Variables

| Variable                     | Description                |
| ---------------------------- | -------------------------- |
| `${workspaceFolder}`         | Workspace root path        |
| `${workspaceFolderBasename}` | Workspace folder name      |
| `${selection}`               | Currently selected text    |
| `${selectedText}`            | Alias for selection        |
| `${file}`                    | Current file path          |
| `${fileBasename}`            | Current filename           |
| `${fileDirname}`             | Current file directory     |
| `${fileBasenameNoExtension}` | Filename without extension |

### Input Variables (User-Provided)

```markdown
${input:variableName}                    # Basic input
${input:variableName:placeholder text} # With placeholder hint
```

**Example:**

```markdown
## Input

- **Component Name**: ${input:componentName:Enter component name (e.g., UserProfile)}
- **Target Directory**: ${input:targetDir:src/components}

Create a new React component named **${input:componentName}** in `${input:targetDir}/`.
```

### Tool References in Body

Reference tools using `#tool:<tool-name>` syntax:

```markdown
Use #tool:githubRepo to search for similar implementations.
Run #tool:runTests after generating the code.
```

### File References

Use relative Markdown links to reference workspace files:

```markdown
Follow the patterns in [coding standards](../../instructions/typescript.instructions.md).
See [API reference](../docs/api.md) for endpoint details.
```

---

## 6) Prompt Writing Best Practices

### ✅ Do

| Practice                       | Example                                         |
| ------------------------------ | ----------------------------------------------- |
| Use imperative mood            | "Create", "Generate", "Analyze"                 |
| Be specific and actionable     | "Create a React component with props interface" |
| Include expected output format | "Output as Markdown with code blocks"           |
| Provide concrete examples      | Show good/bad patterns                          |
| State constraints clearly      | "Do not modify existing files"                  |
| Define error handling          | "If file not found, ask user for path"          |
| Use structured formatting      | Headers, bullets, tables                        |

### ❌ Don't

| Anti-Pattern         | Why                                      |
| -------------------- | ---------------------------------------- |
| Vague instructions   | "Make it better" — not actionable        |
| Missing context      | Assumes agent knows project specifics    |
| Ambiguous terms      | "should", "might", "possibly" — unclear  |
| Hardcoded paths      | Use variables instead                    |
| Excessive prose      | Keep scannable with bullets/tables       |
| Conflicting guidance | "Be brief" + "Include all details"       |
| Over 500 lines       | Split into multiple prompts or use skill |

### Tone and Style

- **Direct**: Write for Copilot, not humans explaining to Copilot
- **Neutral**: Avoid idioms, humor, cultural references
- **Localizable**: Short sentences, no ambiguous pronouns
- **Scannable**: Use headers, bullets, numbered lists

---

## 7) Example Patterns

### Simple Task Prompt

```markdown
---
description: "Generate README.md for the current project based on structure and existing docs"
agent: "agent"
tools: ["codebase", "editFiles"]
---

# Create README

You are a technical writer specializing in developer documentation.

## Task

Create a comprehensive README.md for the current workspace.

## Instructions

1. Analyze the project structure using #tool:codebase
2. Identify: tech stack, main features, setup requirements
3. Generate README with: overview, installation, usage, contributing

## Output

Create `README.md` in workspace root with:

- Project title and description
- Installation steps
- Usage examples with code blocks
- Contributing guidelines
```

### Input-Driven Prompt

```markdown
---
description: "Create a new React component with TypeScript and tests. Use for scaffolding components."
agent: "agent"
tools: ["codebase", "editFiles", "search"]
---

# Create React Component

## Task

Generate a new React functional component with TypeScript.

## Input

- **Component Name**: ${input:name:ComponentName}
- **Directory**: ${input:dir:src/components}
- **Include Tests**: ${input:tests:yes}

## Instructions

1. Check if component already exists at `${input:dir}/${input:name}/`
2. Create component file with:
   - Proper TypeScript interface for props
   - Functional component using hooks
   - Default export
3. If tests requested, create test file with basic render test

## Output

Create:

- `${input:dir}/${input:name}/${input:name}.tsx`
- `${input:dir}/${input:name}/${input:name}.test.tsx` (if tests=yes)
- `${input:dir}/${input:name}/index.ts` (barrel export)
```

### Review/Analysis Prompt

```markdown
---
description: "Perform security review on selected code. Use for vulnerability scanning and OWASP checks."
agent: "ask"
tools: ["codebase", "search"]
---

# Security Review

You are a security specialist focused on application security.

## Task

Analyze the provided code for security vulnerabilities.

## Input

Analyze: ${selection}

If no selection, scan the current file: ${file}

## Instructions

1. Check for OWASP Top 10 vulnerabilities
2. Review input validation and sanitization
3. Identify authentication/authorization issues
4. Check for secrets or hardcoded credentials
5. Review data exposure risks

## Output

Format findings as:

\`\`\`markdown

## [SEVERITY] Finding Title

**Location**: file:line
**Issue**: Description of vulnerability
**Risk**: Impact if exploited
**Fix**: Recommended remediation
\`\`\`

Severity levels: Critical, High, Medium, Low, Info
```

---

## 8) Prompt vs Instruction Content Rules

For the full artifact selection matrix, see [ai-artifacts.instructions.md](./ai-artifacts.instructions.md). This section covers prompt-specific content placement.

### Goes in PROMPTS (`.prompt.md`)

- On-demand task workflows
- User-provided variable inputs
- Specific output format requirements
- Step-by-step procedures for single tasks
- Tool-specific workflows (generate tests, create docs)
- Code generation templates

### Goes in INSTRUCTIONS (`.instructions.md`)

- Coding standards applied to file patterns
- Naming conventions
- Framework-specific patterns
- Error handling patterns
- Always-active rules (no user trigger)
- Language/style conventions

### Referencing Instructions from Prompts

Prompts can reference instruction files:

```markdown
## Guidelines

Follow the standards defined in:

- [TypeScript conventions](../instructions/typescript.instructions.md)
- [Testing patterns](../instructions/testing.instructions.md)
```

Do NOT duplicate instruction content in prompts — link to the single source of truth.

---

## 9) Quality Assurance Checklist

### Frontmatter

- [ ] `description` present and describes WHAT + WHEN + KEYWORDS
- [ ] `description` is 10-500 characters
- [ ] `description` wrapped in single quotes
- [ ] `agent` specified appropriately for task type
- [ ] `tools` limited to necessary set (least privilege)
- [ ] All strings in frontmatter use single quotes

### Body Content

- [ ] Title matches prompt intent (surfaces in Quick Pick)
- [ ] Task section clearly states objective
- [ ] Instructions are step-by-step and actionable
- [ ] Input variables documented with placeholders
- [ ] Output format specified (structure, location, format)
- [ ] Body under 500 lines (split if larger)

### Variables

- [ ] All input variables have descriptive names
- [ ] Placeholders provide guidance to users
- [ ] Default behaviors documented for optional inputs
- [ ] Fallback handling for missing context

### Testing

- [ ] Prompt executes successfully in VS Code (`Chat: Run Prompt`)
- [ ] Output matches expected format
- [ ] Error cases handled gracefully
- [ ] Tool permissions sufficient but minimal

---

## 10) Maintenance Guidance

### Version Control

- Store prompts in `.github/prompts/` with version control
- Use meaningful commit messages for prompt changes
- Test prompts after dependency/framework updates

### When to Update Prompts

- Tool lists change (new tools available, tools deprecated)
- Linked instruction files change
- Output format requirements evolve
- User feedback indicates confusion
- Framework versions change (update examples)

### Deprecation Process

1. Add deprecation notice at top of prompt
2. Point to replacement prompt if applicable
3. Keep functional for transition period
4. Remove after communication to team

### Extracting Shared Guidance

When patterns emerge across multiple prompts:

1. Extract common guidance to `.instructions.md`
2. Update prompts to reference instruction file
3. Remove duplicated content from prompts

---

## 11) Advanced Patterns

### Prompt Chaining (Manual)

For multi-step workflows, create separate prompts that reference each other:

```markdown
## Next Steps

After running this prompt, consider:

- `/review-code` — Review the generated code
- `/generate-tests` — Create tests for new code
```

### Conditional Sections

Use Markdown conditions for dynamic behavior:

```markdown
## Additional Steps

If generating for production:

1. Add error boundaries
2. Include logging
3. Add performance monitoring

If generating for development:

1. Include console.log statements
2. Add TODO comments for incomplete sections
```

### Referencing External Documentation

```markdown
## References

- [Official API docs](https://example.com/api)
- [Framework guide](https://framework.io/guide)

Use these as authoritative sources for implementation details.
```

---

## 12) Naming Conventions

### File Naming

- **Format**: `kebab-case.prompt.md`
- **Characters**: `a-z`, `0-9`, `-`, `.` only
- **Length**: Descriptive but concise
- **Action-oriented**: Start with verb when possible

### ✅ Good Names

```
create-component.prompt.md
generate-tests.prompt.md
review-security.prompt.md
document-api.prompt.md
```

### ❌ Bad Names

```
prompt1.prompt.md           # Not descriptive
myPrompt.prompt.md          # Not kebab-case
CreateComponent.prompt.md   # Wrong case
component creation.prompt.md # Spaces not allowed
```

---

## 13) Common Mistakes to Avoid

### Description Mistakes

- ❌ Too short: "Creates tests"
- ❌ No triggers: "Test generator" (when do I use it?)
- ❌ Not wrapped in quotes in frontmatter

### Tool Mistakes

- ❌ Including all tools when only `codebase` needed
- ❌ Missing `editFiles` when prompt creates files
- ❌ Wrong tool names (check exact spelling)

### Body Mistakes

- ❌ Vague instructions: "Generate good code"
- ❌ No output format: User doesn't know what to expect
- ❌ Conflicting guidance in different sections
- ❌ Over 500 lines without splitting

### Variable Mistakes

- ❌ No placeholder text for inputs
- ❌ Using variables without documenting them
- ❌ Hardcoded paths instead of variables

---

## 14) Additional Resources

### Official Documentation

- [VS Code Prompt Files](https://code.visualstudio.com/docs/copilot/customization/prompt-files)
- [GitHub Copilot Customization Library](https://docs.github.com/en/copilot/tutorials/customization-library)
- [Prompt Engineering Guide](https://docs.github.com/en/copilot/concepts/prompting/prompt-engineering)

### Community Resources

- [Awesome Copilot Prompts](https://github.com/github/awesome-copilot/tree/main/prompts)
- [Copilot Chat Cookbook](https://docs.github.com/en/copilot/tutorials/copilot-chat-cookbook)

### Related Instructions

- [Agent Guidelines](./agents.instructions.md) — For creating custom agents
- [Skills Guidelines](./agent-skills.instructions.md) — For creating skills with bundled assets
- [Instructions Guidelines](./writing-instructions.instructions.md) — For creating instruction files
