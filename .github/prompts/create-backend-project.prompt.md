---
description: Create a new backend project structure (Go domain module)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
model: Claude Opus 4.5 (copilot)
---

# Create Backend Project

You are creating a new domain module in the Kyora backend (Go monolith).

## Project Requirements

${input:domain:What domain are you creating? (e.g., "subscription", "notification", "shipping")}

## Instructions

Read backend architecture rules:

- [backend-core.instructions.md](../instructions/backend-core.instructions.md) for complete architecture patterns
- [backend-testing.instructions.md](../instructions/backend-testing.instructions.md) for test suite structure

## Domain Module Structure

Create the following structure in `backend/internal/domain/${domain}/`:

```
backend/internal/domain/${domain}/
├── domain.go           # Domain models and types
├── service.go          # Business logic
├── storage.go          # Database access layer
├── handler_http.go     # HTTP handlers
```

## Implementation Workflow

1. **Domain Models** (`domain.go`)

   - Define struct types with proper tags (json, db)
   - Include validation rules
   - Add workspace_id for multi-tenancy

2. **Storage** (`storage.go`)

   - Database repository.go implementation
   - Use gorm for query execution and create methods there only if repository.go methods are not sufficient.
   - Proper error handling

3. **Service** (`service.go`)

   - Business logic implementation
   - Validation rules
   - Transaction handling where needed
   - No direct DB calls
   - Proper error messages

4. **Handler** (`handler.go`)

   - HTTP request parsing and validation
   - Call service methods
   - Return proper HTTP responses
   - Error handling with status codes

5. **Tests**

   - Integration tests to lock in the new service features properly and cover all the features with all options.
   - proper assertions to verify correctness.
   - similar structure to existing tests in `internal/tests/e2e/`
   - Cover success and error cases

6. **Routes** (`backend/internal/server/routes.go`)

   - Register endpoints with proper middleware
   - Apply reusable middleware such as auth, permission checks.
   - Group related routes

## Architecture Patterns

- **Service-Repository Pattern**: Separate business logic from data access
- **Multi-Tenancy**: All data scoped to workspace_id
- **Error Handling**: Return domain-specific errors, handle at handler level
- **Validation**: Validate at handler, enforce business rules in service
- **Testing**: Unit tests for all layers

## Example Reference

Look at existing domains for patterns:

- `backend/internal/domain/order/` - Full CRUD example
- `backend/internal/domain/customer/` - Simple resource
- `backend/internal/domain/inventory/` - Complex business logic

## Done

- Complete domain module structure created
- All layers implemented (models, repository, service, handler)
- Routes registered
- Tests written and passing
- Multi-tenancy enforced
- No TODOs or FIXMEs
- Production-ready code
