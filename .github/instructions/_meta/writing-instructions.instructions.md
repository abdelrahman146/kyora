---
description: "Guidelines for creating high-quality, token-efficient custom instruction files for GitHub Copilot"
applyTo: "**/*.instructions.md"
---

# Writing Effective Custom Instructions

Rules for creating and maintaining `.instructions.md` files that are token-efficient, SSOT-compliant, and actionable.

---

## 1) Required Frontmatter

Every instruction file must include YAML frontmatter:

```yaml
---
description: "Brief purpose (1-500 chars)"
applyTo: "glob pattern for target files"
---
```

### `applyTo` Patterns

| Pattern              | Matches                   |
| -------------------- | ------------------------- |
| `**/*.ts`            | All TypeScript files      |
| `**/*.ts,**/*.tsx`   | Multiple extensions       |
| `src/**/*.py`        | Specific directory        |
| `**/*`               | All files (use sparingly) |
| `**/tests/*.spec.ts` | Pattern matching          |

---

## 2) Token Economy

Instructions are sent with every chat message. Every word costs tokens.

### Compression Strategies

| ❌ Wasteful                           | ✅ Efficient |
| ------------------------------------- | ------------ |
| "You should always make sure to use"  | "Use"        |
| "It is recommended that you consider" | "Prefer"     |
| "In order to ensure that"             | "For"        |

### Format Efficiency

1. **Tables over prose** — 3x more compact for comparisons
2. **Bullet lists over paragraphs** — Scannable, fewer filler words
3. **Code examples over descriptions** — Show, don't tell
4. **Remove redundant context** — Agent has codebase access

### Size Guidelines

| File Type                    | Target         | Max     |
| ---------------------------- | -------------- | ------- |
| Repository-wide instructions | 500-1000 words | 2 pages |
| Path-specific instructions   | 200-500 words  | 1 page  |

---

## 3) Writing Style

### Do

- Use imperative mood: "Use", "Implement", "Avoid"
- Be specific: "Max 3 levels of nesting"
- Provide code examples
- Use tables for comparisons
- Reference versions: "TypeScript 5.0+"

### Don't

- "You should", "It's recommended"
- "Avoid deep nesting" (vague)
- Abstract descriptions without code
- Long prose for options
- Unversioned references

---

## 4) SSOT Compliance

### Never Duplicate

If a rule exists elsewhere, reference it:

```markdown
<!-- ❌ BAD: Duplicates forms.instructions.md -->

## Form Handling

Use TanStack Form with Zod validation...
[copies 50 lines]

<!-- ✅ GOOD: References SSOT -->

## Form Handling

See [forms.instructions.md](forms.instructions.md). Key rule: all forms use `useAppForm`.
```

### Cross-Reference Pattern

```markdown
## [Topic]

See [source-file.instructions.md](source-file.instructions.md) for complete rules.

Key points for this context:

- Essential rule 1
- Essential rule 2
```

---

## 5) Content Structure

### Recommended Sections

```markdown
---
description: "[Technology] coding standards for [scope]"
applyTo: "**/*.ext"
---

# [Technology] Development

[One sentence: what this file covers]

## Tech Stack

- **Runtime**: [language] [version]
- **Framework**: [name] [version]

## Core Rules

- [Imperative rule 1]
- [Imperative rule 2]

## Naming Conventions

| Entity    | Convention | Example       |
| --------- | ---------- | ------------- |
| Variables | camelCase  | `userName`    |
| Functions | camelCase  | `getUserById` |

## Common Patterns

### [Pattern Name]

\`\`\`language
// Correct implementation
code here
\`\`\`

## Validation

- Build: `command`
- Test: `command`
```

---

## 6) What to Include

1. **Project overview** — Elevator pitch (2-3 sentences max)
2. **Tech stack** — Languages, frameworks with versions
3. **Coding guidelines** — Naming, formatting, patterns
4. **Project structure** — Key directories and purpose
5. **Build/test commands** — How to validate changes
6. **Available resources** — Scripts, tools (NOT EXTERNAL LINKS)

---

## 7) What to Avoid

| Anti-Pattern              | Why                    |
| ------------------------- | ---------------------- |
| Verbose explanations      | Wastes tokens          |
| Outdated information      | Incorrect suggestions  |
| Ambiguous guidelines      | Agent confusion        |
| Missing examples          | Abstract rules fail    |
| Contradictory advice      | Conflicts cause errors |
| Copy-pasted docs          | No added value         |
| "See styleguide.md" alone | No context provided    |

---

## 8) Examples

### ❌ Bad: Vague

```markdown
Follow good coding practices and write maintainable code.
```

### ✅ Good: Specific

```markdown
- Max function length: 50 lines
- Max nesting depth: 3 levels
- Extract functions over 20 lines
```

### ❌ Bad: External Reference Only

```markdown
Follow the patterns in our style guide at /docs/styleguide.md
```

### ✅ Good: Key Points Extracted

```markdown
Style rules (from /docs/styleguide.md):

- camelCase for variables
- PascalCase for components
- 2-space indentation
```

---

## 9) Testing Instructions

### Validation Process

1. **Syntax**: Verify frontmatter parses correctly
2. **Pattern**: Test `applyTo` matches intended files
3. **Behavior**: Use Copilot to verify rules are followed

### Debug Steps

If instructions aren't followed:

1. Check `References` section in chat response
2. Simplify rules — complex conditions may be ignored
3. Add explicit examples
4. Check for conflicting rules in other files

---

## 10) Quality Checklist

Before committing:

- [ ] `description` present in frontmatter
- [ ] `applyTo` targets specific files (not `**` unless intentional)
- [ ] No rules duplicated from other instruction files
- [ ] No conflicts with `copilot-instructions.md`
- [ ] Every rule is specific and verifiable
- [ ] No TODO/FIXME placeholders
- [ ] Code examples for non-obvious patterns
- [ ] Tested with Copilot

---

## 11) Maintenance

### When to Update

- Framework/library version upgrades
- New patterns emerge in codebase
- Existing patterns deprecated
- Instructions not being followed

### Update Protocol

1. Identify the authoritative file (SSOT)
2. Update only that file
3. Verify `applyTo` patterns still match
4. Test with Copilot
5. Remove obsolete rules
