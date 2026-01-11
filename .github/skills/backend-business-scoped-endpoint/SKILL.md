---
name: backend-business-scoped-endpoint
description: Add a new authenticated business-scoped API endpoint in Kyora backend (Gin) with correct tenant isolation, RBAC/plan gates, error handling, and Swagger annotations.
---

## What this skill does

A safe, repeatable checklist for adding a new endpoint under:

- `GET|POST|PATCH|DELETE /v1/businesses/:businessDescriptor/...`

…matching Kyora’s existing patterns.

## Workflow

1. Decide the scope
   - Business-scoped endpoints belong under the `registerBusinessScopedRoutes` group in `backend/internal/server/routes.go`.
   - Prefer business scoping over accepting `workspaceId`/`businessId` from the client.

2. Add/extend domain logic
   - Domain code lives under `backend/internal/domain/<domain>/`.
   - Keep storage queries scoped by `biz.ID` (tenant isolation).

3. Implement the HTTP handler method
   - Get the authenticated actor with `account.ActorFromContext(c)`.
   - Resolve the business via `business.BusinessFromContext(c)` (or the domain’s helper).
   - Validate request bodies using `request.ValidBody(c, &req)`.
   - Validate query params using `ShouldBindQuery` + `binding` tags.
   - Return errors via `response.Error(c, err)` and success via `response.SuccessJSON(...)`.

4. Wire the route with RBAC and plan gates
   - Add route in `backend/internal/server/routes.go` under the business-scoped group.
   - Apply `account.EnforceActorPermissions(...)` using the correct resource/action.
   - If the endpoint should be plan-gated, apply the relevant billing middleware.

5. Swagger annotations
   - Follow the existing `@Summary/@Tags/@Param/@Success/@Failure/@Router/@Security` style.

6. Verify with E2E
   - Add/extend tests under `backend/internal/tests/e2e/`.
   - Prefer using existing `*_helpers_test.go` helpers for that domain.

## References (SSOT)

- Backend architecture: `.github/instructions/backend-core.instructions.md`
- Go patterns: `.github/instructions/go-backend-patterns.instructions.md`
- Errors & responses: `.github/instructions/errors-handling.instructions.md`
- DTOs/Swagger: `.github/instructions/responses-dtos-swagger.instructions.md`
- Domain SSOT: `.github/instructions/*.instructions.md` (pick the relevant domain file)

## Optional (manual) OpenAPI generation

If you changed Swagger annotations, regenerate docs:

- `make openapi`

(Requires local Go toolchain; run manually.)
