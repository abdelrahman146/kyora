---
description: Create a new backend project structure (Go domain module)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
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
├── model.go            # GORM model(s) + schema mapping
├── storage.go          # Typed repositories + cache (repo pattern)
├── service.go          # Business logic (workspace-scoped)
├── errors.go           # Domain problems (RFC 7807)
├── handler_http.go     # HTTP handlers (request.ValidBody + response.Success/Error)
```

## Implementation Workflow

1. **Models** (`model.go`)

   - Follow `backend-core.instructions.md` model conventions (ID generation, schema mapping, gorm tags)
   - Decide scope first:
     - Workspace-scoped domain (e.g., account/billing/workspace): include `WorkspaceID` where applicable
     - Business-owned domain (most commerce domains): include `BusinessID` and enforce business scoping

2. **Storage** (`storage.go`)

   - Use `database.Repository[T]` and scopes (especially `ScopeWorkspaceID`)
   - Keep caching (if any) in storage only

3. **Service** (`service.go`)

   - Business logic implementation
   - Validation rules
   - Multi-entity writes must use `AtomicProcess.Exec`
   - No direct DB calls
   - Proper error messages

4. **HTTP Handler** (`handler_http.go`)

   - Extract actor/workspace from context
   - Parse + validate via `request.ValidBody`
   - Call service methods
   - Return via `response.Success*` / `response.Error`

5. **Tests**

   - If you add/modify endpoints: add E2E suite(s) under `backend/internal/tests/e2e/` per `backend-testing.instructions.md`
   - If you add service-only logic: add unit tests next to the code (e.g., `service_test.go`)
   - Cover success + error cases; keep tests isolated

6. **Routes** (`backend/internal/server/routes.go`)

   - Register endpoints with proper middleware
   - Apply reusable middleware such as auth, permission checks.
   - Group related routes

## Architecture Patterns

- **Service + Storage**: Business logic in service; persistence + cache in storage
- **Multi-Tenancy**: Two-layer scoped: workspace → businesses. Scope data by the correct key(s).
- **Errors**: Domain problems in `errors.go`, emitted via `response.Error`
- **Validation**: Handler validates request; service enforces business rules
- **Testing**: Follow `backend-testing.instructions.md` patterns

## Example Reference

Look at existing domains for patterns:

- `backend/internal/domain/order/` - Full CRUD example
- `backend/internal/domain/customer/` - Simple resource
- `backend/internal/domain/inventory/` - Complex business logic

## Done

- Complete domain module structure created
- All layers implemented (model, storage, service, handler)
- Routes registered
- Tests written and passing
- Multi-tenancy enforced
- No TODOs or FIXMEs
- Production-ready code
