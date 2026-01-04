---
description: Project-wide refactor in the backend
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Refactor Backend Project

You are performing a project-wide refactor across the Kyora backend (Go monolith).

## Refactor Brief

${input:refactor:Describe the refactor goal (e.g., "Standardize request validation + problem errors across handlers", "Introduce shared query helper to remove duplication", "Unify pagination patterns across list endpoints")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No endpoint changes", "No DB migrations", "No behavior changes", "No new dependencies", "Must be mechanical and low-risk")}

## Instructions

Read backend SSOT first:

- [backend-core.instructions.md](../instructions/backend-core.instructions.md)
- If writing tests: [backend-testing.instructions.md](../instructions/backend-testing.instructions.md)

## Refactor Standards

1. **No behavior drift**: Treat as a mechanical refactor unless explicitly stated otherwise.
2. **Consistency**: Prefer shared helpers over copy/paste.
3. **Multi-tenancy**: Ensure all repositories/queries keep correct workspace/business scoping.
4. **Error/response contract**: Keep RFC 7807 (`response.Error`) and existing problem types unless the refactor explicitly changes them.
5. **Low risk**: Execute in small, verifiable steps with compiler/test guidance.

## Workflow

1. Define the migration rule (before/after pattern)
2. Add/adjust shared utilities in `backend/internal/platform/utils/` (or appropriate platform module) only when it reduces duplication
3. Apply changes across the codebase via search; keep each change mechanical
4. Run: `cd backend && go test ./...`
5. If handlers changed: spot-check a representative endpoint for each affected domain

## Done

- Project-wide refactor completed and constraints respected
- Tenancy scoping preserved (workspace/business)
- Tests pass
- No TODOs or FIXMEs
