---
description: Refactor a single portal-web feature
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Refactor Portal Web Feature

You are refactoring a single feature in the Kyora portal-web (React TanStack dashboard).

## Refactor Brief

${input:refactor:Describe the refactor (e.g., "Extract reusable OrderForm fields into shared components", "Unify query keys for inventory routes", "Replace ad-hoc fetches with Ky client + query hooks")}

## Feature/Area

${input:area:Which feature/area is being refactored? (e.g., "orders", "inventory", "customers", "analytics")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No UI changes", "No new routes", "No API changes", "Must keep existing i18n keys", "Only code motion + type fixes")}

## Instructions

Read relevant portal-web rules first:

- [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md)
- [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md)

If relevant:

- Forms: [forms.instructions.md](../instructions/forms.instructions.md)
- HTTP: [ky.instructions.md](../instructions/ky.instructions.md)
- UI/RTL/A11y: [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md)
- Design tokens: [design-tokens.instructions.md](../instructions/design-tokens.instructions.md)
- Charts: [charts.instructions.md](../instructions/charts.instructions.md)

## Refactor Standards

1. **Refactor ≠ feature**: Default to no behavior/UI changes unless explicitly requested.
2. **Architecture**: TanStack Router + TanStack Query + TanStack Store (no Redux, no Zustand).
3. **Tenancy & scoping (SSOT)**: Business-owned features remain scoped by `businessDescriptor`.
4. **Forms**: Use Kyora form system when forms are involved.
5. **No duplication**: Consolidate repeated logic into `portal-web/src/lib/` or shared components.
6. **i18n/RTL**: If UI text changes, keep Arabic + English keys consistent and verify RTL.

## Workflow

1. Locate the current route(s)/component(s)/query hooks for `${input:area}`
2. Define invariants (what must not change)
3. Apply mechanical refactor steps with compiler guidance (rename/move/extract)
4. Update TanStack Query keys/invalidation patterns if touched
5. Run checks:
   - `cd portal-web && npm run lint`
   - `cd portal-web && npm run type-check`
   - `cd portal-web && npm run test` (if present)
6. Smoke test the feature in both LTR + RTL if UI is involved

## Done

- Feature refactor completed and constraints respected
- Type-check + lint pass
- No TODOs or FIXMEs
- Tenancy scoping preserved (businessDescriptor)
