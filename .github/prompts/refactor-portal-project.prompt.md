---
description: Project-wide refactor in portal-web
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Refactor Portal Web Project

You are performing a project-wide refactor across Kyora portal-web.

## Refactor Brief

${input:refactor:Describe the refactor goal (e.g., "Convert POST/PUT/DELETE requests to TanStack Query mutations across the app", "Standardize error handling for Ky requests", "Unify query key factories for all resources")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No UI changes", "No route changes", "No backend/API changes", "No new dependencies", "Must be mechanical and low-risk")}

## Instructions

Read portal-web SSOT first:

- [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md)
- [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md)

If relevant:

- HTTP: [ky.instructions.md](../instructions/ky.instructions.md)
- Forms: [forms.instructions.md](../instructions/forms.instructions.md)
- UI/RTL/A11y: [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md)
- Design tokens: [design-tokens.instructions.md](../instructions/design-tokens.instructions.md)

## Refactor Standards

1. **No behavior drift**: Treat as “mechanical refactor” unless explicitly stated otherwise.
2. **Consistency**: Introduce shared utilities rather than repeating patterns.
3. **Tenancy & scoping (SSOT)**: Any business-owned route/API usage must keep `businessDescriptor`.
4. **No new architecture**: Stay within TanStack Router/Query/Store conventions.
5. **Low risk**: Prefer incremental, search-and-replace friendly steps.

## Workflow

1. Define an explicit migration rule (before/after pattern)
2. Create/adjust shared helpers in `portal-web/src/lib/` (query key factories, mutation wrappers, error mapping) if needed
3. Update usages across the project via search:
   - Replace ad-hoc POST/PUT/DELETE calls with mutation hooks
   - Standardize invalidation/refetch behavior
4. Run checks:
   - `cd portal-web && npm run lint`
   - `cd portal-web && npm run type-check`
   - `cd portal-web && npm run test` (if present)
5. Smoke-test a representative flow for each major domain (orders, inventory, customers) affected

## Done

- Project-wide refactor completed and constraints respected
- Lint + type-check pass
- No TODOs or FIXMEs
- Tenancy scoping preserved (businessDescriptor)
