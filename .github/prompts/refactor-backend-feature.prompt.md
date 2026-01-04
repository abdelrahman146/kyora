---
description: Refactor a single backend feature/domain
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Refactor Backend Feature

You are refactoring a single feature/domain in the Kyora backend (Go monolith).

## Refactor Brief

${input:refactor:Describe the refactor (e.g., "Extract shared validation into a utility", "Split a large service method into smaller ones", "Standardize problem errors for a domain")}

## Domain/Area

${input:area:Which domain/area is being refactored? (e.g., "order", "inventory", "customer", "analytics")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No endpoint changes", "No DB migrations", "No behavior changes", "Only internal refactor", "Must keep backward compatibility")}

## Instructions

Read backend SSOT first:

- [backend-core.instructions.md](../instructions/backend-core.instructions.md)
- If writing tests: [backend-testing.instructions.md](../instructions/backend-testing.instructions.md)

If relevant:

- Email: [resend.instructions.md](../instructions/resend.instructions.md)
- Billing: [stripe.instructions.md](../instructions/stripe.instructions.md)
- Uploads: [asset_upload.instructions.md](../instructions/asset_upload.instructions.md)

## Refactor Standards

1. **Refactor ≠ feature**: Default to no behavior changes unless explicitly requested.
2. **Architecture**: Preserve the domain conventions (model/storage/service/errors/handler_http).
3. **Multi-tenancy**: Keep correct scoping (WorkspaceID / BusinessID) with zero leaks.
4. **Validation/errors**: Keep `request.ValidBody` and RFC 7807 via `response.Error`.
5. **No duplication**: Consolidate reusable logic into the appropriate domain/platform utilities.

## Workflow

1. Locate the handler/service/storage code for `${input:area}`
2. Define invariants (contract, behavior, tenancy) and constraints
3. Apply mechanical refactor steps with compiler guidance (rename/move/extract)
4. Update tests or add minimal regression coverage where patterns exist
5. Run: `cd backend && go test ./...`

## Done

- Feature/domain refactor completed and constraints respected
- Tenancy scoping preserved (workspace/business)
- Tests pass
- No TODOs or FIXMEs
