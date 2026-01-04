---
description: Add a new feature to the backend (Go API)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Add Backend Feature

You are implementing a new feature in the Kyora backend (Go monolith).

## Feature Requirements

${input:feature:Describe the feature you want to add (e.g., "Add endpoint to export orders as CSV")}

## Instructions

Before implementing, read the backend architecture rules:

- Read [backend-core.instructions.md](../instructions/backend-core.instructions.md) for architecture patterns, service layer structure, and database conventions
- If writing tests: [backend-testing.instructions.md](../instructions/backend-testing.instructions.md)
- If feature involves email: [resend.instructions.md](../instructions/resend.instructions.md)
- If feature involves billing: [stripe.instructions.md](../instructions/stripe.instructions.md)
- If feature involves file uploads: [asset_upload.instructions.md](../instructions/asset_upload.instructions.md)

## Implementation Standards

1. **Architecture**: Follow the domain conventions from `backend-core.instructions.md` (model/storage/service/errors/handler_http).
2. **Multi-Tenancy**: Kyora is two-layer scoped: workspace → businesses.
   - Workspace-scoped resources: scope by `WorkspaceID`.
   - Business-owned resources (inventory/accounting/orders/assets/analytics/customers/storefront): scope by `BusinessID`.
3. **Validation**: Use `request.ValidBody(c, &req)` for JSON binding + validation.
4. **Error Handling**: Return RFC 7807 via `response.Error(c, err)` (domain errors via `problem.*`).
5. **Testing**: Add tests following `backend-testing.instructions.md` (E2E suites under `backend/internal/tests/e2e/`; unit tests live next to code in the domain/platform folders).
6. **Security**: Prevent SQL injection, validate permissions, sanitize inputs, ensure BOLA compliance and abuse prevention
7. **Database**: Use the repository pattern and `AtomicProcess.Exec` for multi-step writes.
8. **Documentation**: Update Swagger docs if adding/modifying endpoints by generating new docs using `swag init`

## Workflow

1. Search for similar features to understand patterns
2. Create/modify domain models in `backend/internal/domain/{domain}/`
3. Implement service logic with business rules
4. Use storage repositories + scopes; add new storage helpers only when repository/scopes are insufficient.
5. Add HTTP handlers with proper validation
6. Register routes in `backend/internal/server/routes.go`
   - If business-scoped: add under `registerBusinessScopedRoutes` (`/v1/businesses/:businessDescriptor/...`) and apply `business.EnforceBusinessValidity`.
7. Write tests covering success and error cases
8. Update Swagger documentation
9. Test by running `cd backend && go test ./...`.

## Done

- Implementation complete, production-ready
- All tests pass
- No TODOs or FIXMEs
- Follows backend-core.instructions.md patterns
- Tests follow backend-testing.instructions.md patterns
- Workspace-scoped for workspace resources and business scoped for business resources (multi-tenancy compliant)
