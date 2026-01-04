---
description: Add a new feature to the backend (Go API)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
model: Claude Opus 4.5 (copilot)
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

1. **Architecture**: Follow service-storage-handler pattern that is already implemented in the codebase with proper seperation of concerns for each layer.
2. **Multi-Tenancy**: All queries must filter by `workspaceID` for workspace scoped resources and `workspaceID / businessID` for business scoped resources.
3. **Validation**: Validate all inputs at handler level
4. **Error Handling**: Return proper HTTP status codes and error messages
5. **Testing**: Include integration tests in `*_test.go` files inside `internal/tests/e2e/` following patterns from backend-testing.instructions.md
6. **Security**: Prevent SQL injection, validate permissions, sanitize inputs, ensure BOLA compliance and abuse prevention
7. **Database**: Use proper indexing, connection pooling, transactions where needed using the `atomic` processor we already have.
8. **Documentation**: Update Swagger docs if adding/modifying endpoints by generating new docs using `swag init`

## Workflow

1. Search for similar features to understand patterns
2. Create/modify domain models in `backend/internal/domain/{domain}/`
3. Implement service logic with business rules
4. Create storage methods for database access incase new queries are needed otherwise use the storage defined methods that were implemented already by applying the repository struct.
5. Add HTTP handlers with proper validation
6. Register routes in `backend/internal/server/routes.go`
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
