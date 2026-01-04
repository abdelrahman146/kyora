---
description: Cross-project enhancement (backend + portal-web)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Cross-Project Enhancement

You are implementing an enhancement that requires coordinated changes across:

- `backend/` (Go monolith API)
- `portal-web/` (React TanStack dashboard)

## Enhancement Brief

${input:enhancement:Describe the enhancement end-to-end (e.g., "Add a new computed profit field to order details and display it in the portal", "Add bulk inventory adjustments API + UI", "Add stronger validation errors surfaced in the form")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No DB migrations", "No new routes", "No breaking API changes", "Portal-only UI changes", "Backend-only refactor + keep UI unchanged")}

## Instructions (SSOT)

Read the relevant instruction files before making changes:

### Backend

- [backend-core.instructions.md](../instructions/backend-core.instructions.md)
- If writing tests: [backend-testing.instructions.md](../instructions/backend-testing.instructions.md)
- If billing/email/uploads are involved: [stripe.instructions.md](../instructions/stripe.instructions.md), [resend.instructions.md](../instructions/resend.instructions.md), [asset_upload.instructions.md](../instructions/asset_upload.instructions.md)

### Portal Web

- [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md)
- [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md)
- [portal-web-ui-guidelines.instructions.md](../instructions/portal-web-ui-guidelines.instructions.md) (portal UX/UI SSOT)
- If forms/HTTP/UI/charts are involved: [forms.instructions.md](../instructions/forms.instructions.md), [ky.instructions.md](../instructions/ky.instructions.md), [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md), [charts.instructions.md](../instructions/charts.instructions.md), [design-tokens.instructions.md](../instructions/design-tokens.instructions.md)

## Cross-Project Standards

1. **Single change-set, coherent contract**: Backend API shape and portal types/schemas must match.
2. **No tenancy leaks**:
   - Backend: workspace → business scoping enforced (WorkspaceID / BusinessID as appropriate).
   - Portal: business-owned features scoped via `businessDescriptor`.
     - Routes under `/business/$businessDescriptor/...`
     - API paths under `v1/businesses/${businessDescriptor}/...`
3. **Validation & errors**:
   - Backend returns RFC 7807 via `response.Error(c, err)` and domain errors via `problem.*`.
   - Portal surfaces errors in forms and UX states (no raw stack traces).
4. **i18n + RTL**: Any new/changed portal UI text must have Arabic + English keys and work in RTL.
5. **Keep scope tight**: Do not add unrelated features.

## Workflow

1. **Locate current contract + UI**

   - Find existing endpoint(s), handler/service/repository
   - Find portal route(s) and API client methods consuming them

2. **Define acceptance criteria**

   - What user-visible behavior changes?
   - What stays the same?

3. **Implement backend changes**

   - Update domain/service/storage/handler following backend-core conventions
   - Add/adjust routes in `backend/internal/server/routes.go` if needed
   - Update swagger if endpoints change

4. **Implement portal changes**

   - Update schema(s) in `portal-web/src/schemas/`
   - Update API client in `portal-web/src/api/`
   - Update route/component in `portal-web/src/routes/`
   - Add i18n keys in both locales

5. **Test/verify**
   - Backend: `cd backend && go test ./...`
   - Portal: `cd portal-web && npm run lint` and `cd portal-web && npm run type-check`
   - Manual verification for the specific flow end-to-end

## Done

- Backend + portal ship together with matching contract
- Constraints respected
- Tenancy scoping verified (backend + portal)
- Portal: RTL/LTR + i18n verified if UI changed
- Lint/type-check/tests pass for touched projects
- No TODOs or FIXMEs
