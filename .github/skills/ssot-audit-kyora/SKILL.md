---
name: ssot-audit-kyora
description: Audit a change (backend or portal-web) for Kyora SSOT compliance: tenant isolation, RBAC/plan gates, HTTP/query patterns, i18n/RTL, and instruction-file drift.
---

## What this skill does

A lightweight, repeatable review checklist to catch SSOT violations early.

## Workflow

1. Identify the touched surface area
   - Backend only, portal-web only, or both.

2. Tenant isolation (backend)
   - Confirm no endpoint accepts `workspaceId`/`businessId` from the client unless explicitly intended.
   - Confirm all queries are scoped to the workspace/business derived from middleware/context.

3. RBAC and plan gates
   - Confirm new/changed routes enforce permissions (`admin`/`member`) and any plan limits where applicable.

4. Error handling
   - Confirm handlers return structured problems via the existing response/problem utilities.

5. Portal-web state ownership
   - Confirm list/search/filter state lives in the URL.
   - Confirm server state uses TanStack Query; avoid ad-hoc request logic.

6. i18n/RTL
   - Confirm any new UI copy exists in both `en` and `ar`.
   - Confirm RTL isn’t broken by hard-coded left/right assumptions.

7. SSOT drift prevention
   - If you introduced a new pattern, either:
     - align it to the existing SSOT, or
     - update the relevant instruction file (narrow scope) so SSOT matches reality.

## References (SSOT)

- Repo SSOT: `.github/copilot-instructions.md`
- Backend: `.github/instructions/backend-core.instructions.md`
- Portal web: `.github/instructions/portal-web-architecture.instructions.md`
- Errors: `.github/instructions/errors-handling.instructions.md`
- HTTP/Query: `.github/instructions/http-tanstack-query.instructions.md`
- State: `.github/instructions/state-management.instructions.md`
- i18n: `.github/instructions/i18n-translations.instructions.md`
