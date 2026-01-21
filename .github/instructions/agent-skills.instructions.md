---
description: "Guidelines for creating, maintaining, and using GitHub Copilot Agent Skills — SKILL.md format, bundled resources, and progressive disclosure patterns"
applyTo: "**/.github/skills/**/SKILL.md,**/.claude/skills/**/SKILL.md,**/skills/**/SKILL.md"
---

# Agent Skills Development Guide

Rules for creating effective, portable Agent Skills that enhance GitHub Copilot with specialized capabilities, workflows, and bundled resources.

**Artifact Selection**: See [ai-artifacts.instructions.md](./ai-artifacts.instructions.md) for decision matrix on when to create skills vs prompts vs agents vs instructions.

---

## 1) Directory Structure

Skills must be stored in specific locations:

| Location                          | Scope                | Use When                                    |
| --------------------------------- | -------------------- | ------------------------------------------- |
| `.github/skills/<skill-name>/`    | Project/repository   | **Recommended** for project-specific skills |
| `.claude/skills/<skill-name>/`    | Project/repository   | Legacy (backward compatibility only)        |
| `~/.copilot/skills/<skill-name>/` | Personal (user-wide) | Skills used across all projects             |
| `~/.claude/skills/<skill-name>/`  | Personal (user-wide) | Legacy (backward compatibility only)        |

### Folder Structure Example

```
.github/skills/my-skill/
├── SKILL.md              # Required: Main instructions
├── LICENSE.txt           # Recommended: Apache 2.0 or MIT
├── scripts/              # Optional: Executable automation
│   ├── helper.py         # Python script
│   └── helper.ps1        # PowerShell script
├── references/           # Optional: Documentation (loaded into context)
│   ├── api-reference.md
│   └── workflow-guide.md
├── assets/               # Optional: Static files used AS-IS
│   └── template.html
└── templates/            # Optional: Code scaffolds agent modifies
    └── starter.py
```

---

## 2) SKILL.md Format Specification

### Required Frontmatter

```yaml
---
name: skill-name
description: Clear description of WHAT it does AND WHEN to use it. Include keywords users might mention.
---
```

| Field           | Required | Constraints                                                                                                                     |
| --------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `name`          | **Yes**  | 1-64 chars, lowercase a-z, numbers, hyphens only. Must match folder name. No consecutive hyphens. Cannot start/end with hyphen. |
| `description`   | **Yes**  | 1-1024 chars. Must describe capabilities AND triggers.                                                                          |
| `license`       | No       | License name or reference to LICENSE.txt                                                                                        |
| `compatibility` | No       | 1-500 chars. Environment requirements if needed.                                                                                |
| `metadata`      | No       | Key-value pairs for additional properties                                                                                       |
| `allowed-tools` | No       | Space-delimited pre-approved tools (experimental)                                                                               |

### Description: The Discovery Mechanism

**CRITICAL**: The `description` is how Copilot discovers skills. Copilot reads ONLY `name` and `description` to decide whether to activate a skill.

**Formula for Good Descriptions:**

```
[WHAT it does] + [WHEN to use it] + [Keywords/triggers]
```

**✅ Good Description:**

```yaml
description: "Toolkit for testing local web applications using Playwright. Use when asked to verify frontend functionality, debug UI behavior, capture browser screenshots, check for visual regressions, or view browser console logs. Supports Chrome, Firefox, and WebKit browsers."
```

**❌ Bad Description:**

```yaml
description: "Web testing helpers"
```

Bad because:

- No capabilities (what can it do?)
- No triggers (when should Copilot load it?)
- No keywords (what prompts would match?)

### Body Content Sections

| Section                     | Purpose                                    | Required?                |
| --------------------------- | ------------------------------------------ | ------------------------ |
| `# Title`                   | Brief overview                             | Yes                      |
| `## When to Use This Skill` | List of scenarios (reinforces description) | Recommended              |
| `## Prerequisites`          | Required tools, dependencies               | Recommended              |
| `## Step-by-Step Workflows` | Numbered procedures                        | Yes (for complex skills) |
| `## Troubleshooting`        | Common issues table                        | Recommended              |
| `## References`             | Links to bundled docs                      | If resources exist       |

---

## 3) Progressive Disclosure Architecture

Skills use three-level loading for context efficiency:

| Level             | What Loads                     | When                             | Token Cost               |
| ----------------- | ------------------------------ | -------------------------------- | ------------------------ |
| **1. Discovery**  | `name` + `description` only    | Always (at startup)              | ~50-100 tokens           |
| **2. Activation** | Full `SKILL.md` body           | When request matches description | <5000 tokens recommended |
| **3. Resources**  | Scripts, references, templates | Only when explicitly referenced  | Variable                 |

### Optimization Rules

1. **Keep SKILL.md body under 500 lines** — Split large content into `references/` folder
2. **Install many skills safely** — Only metadata loads until needed
3. **Split long workflows (>5 steps)** — Move to `references/workflow-name.md`
4. **Reference bundled files explicitly** — Agent loads on-demand

---

## 4) Bundling Resources

### Resource Type Distinction

| Folder        | Purpose                           | Loaded into Context?     | Example Files                      |
| ------------- | --------------------------------- | ------------------------ | ---------------------------------- |
| `scripts/`    | Executable automation             | When **executed**        | `helper.py`, `validate.sh`         |
| `references/` | Documentation agent reads         | Yes, when **referenced** | `api-reference.md`, `schema.md`    |
| `assets/`     | Static files used **AS-IS**       | No (used in output)      | `logo.png`, `report-template.html` |
| `templates/`  | Code scaffolds agent **modifies** | Yes, when **referenced** | `scaffold.py`, `hello-world/`      |

### Assets vs Templates

| Type          | Agent Action              | Example Use Case                      |
| ------------- | ------------------------- | ------------------------------------- |
| **Assets**    | Copy without modification | Brand images, fonts, fixed configs    |
| **Templates** | Read, modify, extend      | Starter code, scaffolds, boilerplates |

### Referencing Resources

Use relative paths from SKILL.md:

```markdown
## Available Scripts

Run the [helper script](./scripts/helper.py) to automate common tasks.

See [API reference](./references/api-reference.md) for detailed documentation.

Use the [scaffold](./templates/scaffold.py) as a starting point.
```

### When to Bundle Scripts

Include scripts when:

- Same code would be rewritten repeatedly by agent
- Deterministic reliability is critical
- Complex logic benefits from being pre-tested
- Operation has self-contained purpose
- Testability matters
- Predictable behavior preferred over dynamic generation

---

## 5) Script Requirements

### Cross-Platform Languages

| Language               | Use Case                            |
| ---------------------- | ----------------------------------- |
| Python                 | Complex automation, data processing |
| PowerShell Core (pwsh) | Cross-platform PowerShell           |
| Node.js                | JavaScript-based tooling            |
| Bash/Shell             | Simple Unix automation              |

### Script Best Practices

```python
#!/usr/bin/env python3
"""
Helper script for [skill-name].

Usage:
    python helper.py --input <file> --action <action>

Arguments:
    --input   Input file or URL to process
    --action  Action to perform (extract|convert|validate)
    --help    Show this help message
"""

import argparse
import sys

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--input', required=True, help='Input file')
    parser.add_argument('--action', required=True, choices=['extract', 'convert', 'validate'])
    args = parser.parse_args()

    try:
        # Implementation here
        pass
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()
```

**Requirements:**

- Include `--help` documentation
- Handle errors gracefully with clear messages
- Use exit codes (0 success, non-zero failure)
- No hardcoded credentials
- Use relative paths where possible
- Warn before irreversible actions

---

## 6) Writing Style

### Do

- Use imperative mood: "Run", "Create", "Configure" (not "You should run")
- Be specific and actionable
- Include exact commands with parameters
- Show expected outputs where helpful
- Keep sections focused and scannable
- Use tables for parameter documentation

### Don't

- Vague instructions ("consider optimizing")
- Missing context (assume agent knows project)
- Hardcoded paths (use placeholders)
- Credentials in examples
- Excessive prose (prefer bullets/tables)

---

## 7) Common Patterns

### Parameter Table Pattern

```markdown
| Parameter   | Required | Default     | Description                  |
| ----------- | -------- | ----------- | ---------------------------- |
| `--input`   | Yes      | -           | Input file or URL to process |
| `--output`  | No       | `./output/` | Output directory             |
| `--verbose` | No       | `false`     | Enable verbose logging       |
```

### Workflow Execution Pattern

````markdown
## Workflow: [Task Name]

### Prerequisites

- [ ] [Prerequisite 1]
- [ ] [Prerequisite 2]

### Steps

1. **[Step Name]**: [Action]
   ```bash
   command --flag value
   ```
````

2. **[Step Name]**: [Action]
   - Expected output: [description]

### Verification

- [ ] [Verification step 1]
- [ ] [Verification step 2]

````

### Troubleshooting Table Pattern

```markdown
## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Skill not discovered | Description too vague | Add keywords and triggers |
| Script fails | Missing dependency | Run `pip install -r requirements.txt` |
| Assets not found | Wrong path | Use relative paths from skill root |
````

---

## 8) Validation Checklist

Before publishing a skill:

**Frontmatter:**

- [ ] `name` is lowercase with hyphens, ≤64 characters
- [ ] `name` matches folder name exactly
- [ ] No consecutive hyphens in `name`
- [ ] `description` is 10-1024 characters
- [ ] `description` states WHAT, WHEN, and KEYWORDS

**Body:**

- [ ] Includes "When to Use" section
- [ ] Includes "Prerequisites" if applicable
- [ ] Step-by-step workflows have numbered steps
- [ ] Body under 500 lines (split if larger)

**Resources:**

- [ ] Scripts include help documentation
- [ ] Scripts handle errors gracefully
- [ ] Relative paths used for all references
- [ ] No hardcoded credentials or secrets
- [ ] Assets under 5MB each
- [ ] LICENSE.txt included (Apache 2.0 recommended)

---

## 9) Security Considerations

### Script Execution

- Scripts run in user's environment with user privileges
- VS Code provides terminal tool controls with auto-approve options
- Include `--force` flags only for destructive operations
- Document any network operations or external calls

### Audit Before Use

When using shared skills:

1. Read all files in the skill directory
2. Check code dependencies
3. Review bundled assets
4. Watch for instructions connecting to external sources
5. Verify scripts don't exfiltrate data

---

## 10) Complete Skill Example

```
.github/skills/webapp-testing/
├── SKILL.md
├── LICENSE.txt
└── test-helper.js
```

**SKILL.md:**

````markdown
---
name: webapp-testing
description: "Toolkit for testing local web applications using Playwright. Use when asked to verify frontend functionality, debug UI behavior, capture browser screenshots, check for visual regressions, or view browser console logs. Supports Chrome, Firefox, and WebKit."
license: Complete terms in LICENSE.txt
---

# Web Application Testing

Comprehensive testing and debugging of local web applications using Playwright automation.

## When to Use This Skill

- Test frontend functionality in a real browser
- Verify UI behavior and interactions
- Debug web application issues
- Capture screenshots for documentation
- Inspect browser console logs
- Validate form submissions and user flows

## Prerequisites

- Node.js installed on the system
- A locally running web application
- Playwright (installed automatically if not present)

## Core Capabilities

### 1. Browser Automation

- Navigate to URLs
- Click buttons and links
- Fill form fields
- Handle dialogs

### 2. Verification

- Assert element presence
- Verify text content
- Test responsive behavior

### 3. Debugging

- Capture screenshots
- View console logs

## Usage Examples

### Basic Navigation Test

```javascript
await page.goto("http://localhost:3000");
const title = await page.title();
console.log("Page title:", title);
```
````

### Form Interaction

```javascript
await page.fill("#username", "testuser");
await page.fill("#password", "password123");
await page.click('button[type="submit"]');
await page.waitForURL("**/dashboard");
```

## Guidelines

1. **Verify app is running** before tests
2. **Use explicit waits** for elements
3. **Capture screenshots** on failure
4. **Clean up resources** — always close browser
5. **Handle timeouts gracefully**

## Troubleshooting

| Issue             | Solution                                   |
| ----------------- | ------------------------------------------ |
| Page not loading  | Verify server is running on specified port |
| Element not found | Check selector, use `waitForSelector`      |
| Timeout errors    | Increase timeout, verify network           |

````

---

## 11) Skill Maintenance

### Version Control

- Track skills in `.github/skills/` with version control
- Use meaningful commit messages for skill changes
- Consider semantic versioning in `metadata` field

### Documentation Updates

When skill behavior changes:
1. Update `description` if triggers change
2. Update prerequisites if dependencies change
3. Update workflows if procedures change
4. Bump version in metadata

### Deprecation

When deprecating a skill:
1. Add deprecation notice to SKILL.md body
2. Point to replacement skill if applicable
3. Keep functional for transition period
4. Remove after communication period

---

## 12) Agent Awareness

For agents to effectively use skills:

### In System Prompts

Skills are automatically discovered from configured directories. Agents receive skill metadata as:

```xml
<available_skills>
  <skill>
    <name>webapp-testing</name>
    <description>Toolkit for testing local web applications...</description>
    <location>/path/to/skills/webapp-testing/SKILL.md</location>
  </skill>
</available_skills>
````

### Activation Flow

1. Agent receives user prompt
2. Agent matches prompt keywords against skill descriptions
3. If matched, agent reads full SKILL.md into context
4. Agent follows instructions, accessing resources as needed

### Multi-Skill Usage

Agents can:

- Activate multiple skills in one session
- Chain skills for complex workflows
- Reference resources across activated skills
