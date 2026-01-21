---
name: copilot-instructions-authoring
description: Rules for authoring and maintaining .github/copilot-instructions.md to maximize Copilot effectiveness under context limits.
applyTo: ".github/copilot-instructions.md"
---

# Authoring Rules for .github/copilot-instructions.md

## Goal

Maintain a short, high-signal repository onboarding file that improves Copilot Chat, coding agent tasks, and Copilot code review.

## Hard constraints (token optimization)

- Keep this file concise and scannable: prefer headings + bullets; avoid long prose.
- Put the highest-impact details in the first ~15–25 lines.
- Prefer concrete, testable instructions (exact commands, file paths, invariants) over general advice.
- Do not repeat language/framework-specific rules here; place them in scoped `*.instructions.md` files to avoid conflicts and reduce noise.

## Content requirements (use this structure)

1. **Project overview (2–5 lines)**
   - Elevator pitch: what the repo is, who it serves, what “success” looks like.

2. **Tech stack (bullets)**
   - Primary languages/frameworks and key infrastructure dependencies (DB/cache/queues/etc.).
   - Testing tools and where tests live.

3. **Repo map (bullets)**
   - Top-level folders and what they contain.
   - Where to find configs, docs, scripts, CI, and common entry points.

4. **Build / test / validate (must be actionable)**
   - Minimal “happy path” commands to: set up, run, test, lint, format, and build.
   - Include common failure fixes only if they are frequent and deterministic.

5. **Coding rules (10–20 bullets max, MUST/SHOULD/MAY)**
   - Only include rules that measurably reduce bugs or rework:
     - required tests
     - error handling invariants
     - API/contract expectations
     - security-critical guardrails
   - Prefer “Do X, not Y” where ambiguity is common.

6. **Resources (bullets)**
   - Point to _in-repo_ scripts and automation (Make targets, package scripts, task runners).
   - Mention available MCP tools only in terms of what they are used for (e.g., Playwright for UI verification).

## What NOT to include

- Vague directives: “be more accurate”, “don’t miss issues”, “be consistent”.
- Instructions that require following external links; copy critical content into the repo instead.
- Requests to change Copilot’s UI/formatting behavior (e.g., “use bold”, “add emojis”).
- Large reference material, style guides, or tutorial content (move to dedicated docs/skills).

## Consistency + maintenance rules

- Every command must be verified against the repository (do not guess).
- Only pin versions if the repo tooling is pinned; otherwise omit versions.
- If a rule conflicts with a scoped `*.instructions.md` file, remove it here and keep it scoped.
- After edits, validate effectiveness by running a small change/PR review and iterating on missed/ignored instructions.

## Suggested minimal template (copy/paste inside copilot-instructions.md)

- Title + 2–5 line overview
- Tech stack (bullets)
- Repo map (bullets)
- Build/test/validate commands (bullets)
- Coding rules (10–20 bullets)
- Resources (scripts/tools) (bullets)
