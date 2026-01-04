---
description: Improve an existing feature in portal-web (React dashboard)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Improve Portal Web Feature

You are improving an existing feature in the Kyora portal-web (React TanStack dashboard).

## Improvement Brief

${input:improvement:Describe what you want to improve (e.g., "Make order creation faster by reducing steps", "Improve inventory adjustment UX on mobile", "Fix slow dashboard chart loading without changing UI")}

## Constraints

${input:constraints:List constraints (optional) (e.g., "No new routes", "No new API endpoints", "Keep existing UI layout", "Must keep backward compatibility")}

## Instructions

Before implementing, read the portal-web rules:

- Read [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md) for tech stack, routing, state management
- Read [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md) for development workflow
- Read [portal-web-ui-guidelines.instructions.md](../instructions/portal-web-ui-guidelines.instructions.md) for portal UX/UI SSOT (mobile-first, Arabic/RTL-first, minimal)

If relevant to the improvement:

- Forms: [forms.instructions.md](../instructions/forms.instructions.md)
- HTTP requests: [ky.instructions.md](../instructions/ky.instructions.md)
- UI components / RTL / accessibility: [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md)
- Charts: [charts.instructions.md](../instructions/charts.instructions.md)
- Design tokens: [design-tokens.instructions.md](../instructions/design-tokens.instructions.md)
- Billing: [stripe.instructions.md](../instructions/stripe.instructions.md)
- File uploads: [asset_upload.instructions.md](../instructions/asset_upload.instructions.md)

## Improvement Standards

1. **Improve what exists**: Prefer refactors, UX polish, and performance improvements over adding new screens or flows.
2. **Architecture**: TanStack Router + TanStack Query + TanStack Store (no Redux, no Zustand)
3. **Routing**: File-based routing in `portal-web/src/routes/` (only add routes if explicitly required)
4. **Tenancy & scoping (SSOT)**: Business-owned functionality must be scoped by `businessDescriptor`
   - UI routes under `/business/$businessDescriptor/...`
   - API calls to `v1/businesses/${businessDescriptor}/...`
   - Read `businessDescriptor` via `Route.useParams()` and pass it into API/query hooks
5. **Forms**: Use Kyora form system (`useKyoraForm` + `<form.AppField>` pattern)
6. **HTTP**: Use Ky client with proper error handling
7. **UI (Portal SSOT)**: Follow `portal-web-ui-guidelines.instructions.md` (mobile-first, Arabic/RTL-first, minimal: no shadows, no gradients)
8. **i18n**: All user-visible strings must support Arabic + English
9. **Accessibility**: Keyboard navigation + appropriate ARIA
10. **Type safety**: Zod schemas + TypeScript types

## Workflow

1. **Locate current behavior**

   - Find the route/component(s) implementing the feature
   - Identify current UX, loading/error states, and data flow
   - Note any existing patterns used elsewhere in portal-web

2. **Define “better” (acceptance criteria)**

   - Specify the measurable outcome (fewer steps, faster load, clearer errors, better RTL, etc.)
   - Confirm what must NOT change (layout, API contract, route, etc.)

3. **Implement improvement**

   - Refactor for readability and reuse (move shared logic to `portal-web/src/lib/` or shared components)
   - Reduce unnecessary re-renders, fix cache invalidation, and tighten query keys
   - Improve error messages and loading skeletons/spinners using existing UI primitives

4. **Regression safety**

   - Add/extend tests when there’s an existing testing pattern for that area
   - Ensure Zod schemas match the API

5. **Verify**
   - `cd portal-web && npm run lint`
   - `cd portal-web && npm run type-check`
   - `cd portal-web && npm run test` (if present in this repo)
   - Manual verification in both English (LTR) and Arabic (RTL)
   - Responsive checks (mobile/tablet/desktop)

## Done

- Improvement shipped and statisfied 100%
- All constraints respected
- RTL + LTR verified
- i18n keys added for both locales (if any text changed)
- Type-check and lint pass
- No TODOs or FIXMEs
