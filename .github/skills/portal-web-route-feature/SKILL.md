---
name: portal-web-route-feature
description: Implement a new Portal Web route/page in Kyora (TanStack Router + TanStack Query/Form + ky HTTP + i18n + RTL). Use when adding a business-scoped screen.
---

## What this skill does

A repeatable workflow for building a Portal Web route that matches Kyora’s current stack and SSOT.

## Prerequisites

- You’re working under `portal-web/`.
- You have a backend endpoint (or a clear contract stub) to call.

## Workflow

1. Pick the correct route location
   - Business-scoped pages live under `portal-web/src/routes/business/$businessDescriptor/...`.

2. Define URL state with `validateSearch`
   - Use Zod in the route file (see existing list pages) to keep URL state typed and validated.

3. Prefetch using the route `loader`
   - Use `queryClient.ensureQueryData(...)` so navigation feels instant and data is cached.

4. Use TanStack Query hooks for data
   - Add/extend query factories and hooks under `portal-web/src/api/...`.
   - Avoid calling `fetch` directly; use the existing HTTP layer patterns.

5. Forms and filters
   - Use `useKyoraForm` for filter forms that sync to URL search params.
   - Keep long-lived state in the URL when it affects list results (search/page/sort/filters).

6. i18n and RTL
   - Add copy to the correct namespace JSON files under `portal-web/src/i18n/{en,ar}/`.
   - Prefer existing namespaces (`common`, `orders`, `inventory`, etc.) instead of creating new ones.

7. Error and empty states
   - Use the project’s error translation/parsing utilities rather than ad-hoc strings.

## References (SSOT)

- Portal architecture: `.github/instructions/portal-web-architecture.instructions.md`
- Code placement rules: `.github/instructions/portal-web-code-structure.instructions.md`
- HTTP + Query rules: `.github/instructions/http-tanstack-query.instructions.md`
- State ownership rules: `.github/instructions/state-management.instructions.md`
- i18n rules: `.github/instructions/i18n-translations.instructions.md`
- RTL/UI rules: `.github/instructions/ui-implementation.instructions.md` and `.github/instructions/portal-web-ui-guidelines.instructions.md`

## Quick self-check

- Route uses `createFileRoute(...)` and Zod `validateSearch`.
- Data access goes through `portal-web/src/api/*` (Query + ky client), not ad-hoc calls.
- Any new UI copy exists in both `en` and `ar`.
