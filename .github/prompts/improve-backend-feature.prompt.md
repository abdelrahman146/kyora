---
description: Improve an existing feature in the backend (Go API)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Improve Backend Feature

You are improving an existing feature in the Kyora backend (Go monolith).

## Improvement Brief

${input:improvement:Describe what you want to improve (e.g., "Reduce N+1 queries on orders list", "Tighten validation for inventory adjustments", "Make order create endpoint return clearer problem details")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No DB migrations", "No new endpoints", "Keep response shape unchanged", "Only internal refactor", "Must keep backward compatibility")}

## Instructions

Before implementing, read the backend rules:

- Read [backend-core.instructions.md](../instructions/backend-core.instructions.md) for domain/service/storage/handler patterns
- If writing tests: [backend-testing.instructions.md](../instructions/backend-testing.instructions.md)
- If feature involves email: [resend.instructions.md](../instructions/resend.instructions.md)
- If feature involves billing: [stripe.instructions.md](../instructions/stripe.instructions.md)
- If feature involves file uploads: [asset_upload.instructions.md](../instructions/asset_upload.instructions.md)

## Improvement Standards

1. **Improve what exists**: Prefer refactors, performance, correctness, and better errors over expanding surface area.
2. **Architecture**: Follow `backend-core.instructions.md` conventions (model/storage/service/errors/handler_http).
3. **Multi-tenancy**: Workspace → Businesses (no cross-scope leaks).
   - Workspace-scoped resources: scope by `WorkspaceID`.
   - Business-owned resources: scope by `BusinessID`.
4. **Validation**: Use `request.ValidBody(c, &req)` for JSON binding + validation.
5. **Error handling**: Return RFC 7807 via `response.Error(c, err)`; use domain errors via `problem.*`.
6. **Database**: Use repository pattern; use `AtomicProcess.Exec` for multi-step writes.
7. **Security**: Prevent BOLA, validate permissions, sanitize inputs, avoid SQL injection.
8. **Docs**: If endpoints change, regenerate Swagger docs (`swag init`).

## Workflow

1. **Locate current behavior**

   - Find the handler/service/repository implementing the feature
   - Identify slow paths, unclear errors, duplicated logic, or missing validation

2. **Define “better” (acceptance criteria)**

   - Specify measurable outcome (fewer queries, less latency, clearer errors, fewer edge-case failures)
   - Confirm what must NOT change (API shape, DB schema, routes)

3. **Implement improvement**

   - Fix root cause (not symptoms)
   - Refactor shared logic into appropriate domain/platform utilities when reused
   - Keep tenancy filters explicit and correct in repositories

4. **Regression safety**

   - Add/extend tests if there is an existing testing pattern for this area

5. **Verify**
   - `cd backend && go test ./...`
   - If behavior is user-facing/API: validate via E2E or manual request against dev server

## Done

- Improvement shipped and constraints respected
- Tenancy rules preserved (workspace/business scoping)
- Tests pass
- No TODOs or FIXMEs
- Follows backend-core and backend-testing rules
