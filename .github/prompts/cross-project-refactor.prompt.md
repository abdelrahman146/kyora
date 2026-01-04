---
description: Cross-project refactor (backend + portal-web)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Cross-Project Refactor

You are performing a refactor that spans `backend/` and `portal-web/` to improve maintainability, consistency, and correctness.

## Refactor Brief

${input:refactor:Describe the refactor goal (e.g., "Unify error codes for order creation and update portal handling", "Refactor backend DTOs and portal schemas to remove duplication", "Rename a field across API + UI without breaking behavior")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No behavior changes", "No endpoint changes", "No DB migrations", "Must preserve API compatibility", "Only internal code motion")}

## Instructions (SSOT)

### Backend

- [backend-core.instructions.md](../instructions/backend-core.instructions.md)
- If writing tests: [backend-testing.instructions.md](../instructions/backend-testing.instructions.md)

### Portal Web

- [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md)
- [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md)
- [portal-web-ui-guidelines.instructions.md](../instructions/portal-web-ui-guidelines.instructions.md) (if any UI is involved)
- If forms/HTTP/UI are affected: [forms.instructions.md](../instructions/forms.instructions.md), [ky.instructions.md](../instructions/ky.instructions.md), [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md), [design-tokens.instructions.md](../instructions/design-tokens.instructions.md)

## Refactor Standards

1. **Refactor ≠ feature**: Default to no behavior changes unless explicitly requested.
2. **No contract drift**: If API types or shapes change, update portal schemas/types in the same change-set.
3. **Keep tenancy rules intact**:
   - Backend scoping by WorkspaceID/BusinessID.
   - Portal business scoping by `businessDescriptor` in both routes and API paths.
4. **No duplication**: If you find repeated logic, consolidate into shared utilities:
   - Backend: `backend/internal/platform/utils/` (or appropriate domain/platform module)
   - Portal: `portal-web/src/lib/`
5. **Safety**: Add regression coverage when there’s an existing testing pattern for the refactored area.

## Workflow

1. **Map the surface area**

   - List files/functions/endpoints affected
   - Identify portal consumers (routes, queries, schemas)

2. **Define invariants**

   - What must remain identical (responses, UX, routing, errors, translations)?

3. **Execute the refactor in small, verifiable steps**

   - Rename/move code with compiler guidance
   - Keep changes mechanical when possible
   - Update imports/usages across the monorepo

4. **Update schemas/types**

   - Backend: request/response structs, domain errors
   - Portal: Zod schemas, TS types, API client typings

5. **Verify**
   - Backend: `cd backend && go test ./...`
   - Portal: `cd portal-web && npm run lint` and `cd portal-web && npm run type-check`
   - Manual smoke test for affected flows

## Done

- Refactor completed and constraints respected
- No cross-scope data access regressions
- Portal contract consumption matches backend contract
- Tests/lint/type-check pass for touched projects
- No TODOs or FIXMEs
