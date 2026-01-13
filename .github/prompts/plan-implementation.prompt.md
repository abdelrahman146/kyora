---
name: plan-implementation
description: "Create a staged implementation/refactor plan and generate handoff prompts for Kyora agents. Planning-only; no code changes."
argument-hint: "Describe the feature/refactor, affected areas (backend/portal-web), and any constraints"
agent: Implementation Planner
---

You are planning work for the Kyora monorepo.

Input:

- Request: ${input:request}
- Constraints: ${input:constraints:optional (time, scope, backwards-compat, etc.)}
- Target areas: ${input:areas:backend | portal-web | full-stack | refactor}

Requirements:

- Produce output exactly using the agent’s **Output Format (Always Use)**.
- Include a **Handoff Package (Prompts)** with the right specialist agents.
- If critical details are missing, ask only the minimum blocking questions.
