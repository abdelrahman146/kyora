---
description: Add a new feature to portal-web (React dashboard)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Add Portal Web Feature

You are implementing a new feature in the Kyora portal-web (React TanStack dashboard).

## Feature Requirements

${input:feature:Describe the feature you want to add (e.g., "Add order filtering by status with date range")}

## Instructions

Before implementing, read the frontend architecture rules:

- Read [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md) for tech stack, routing, state management
- Read [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md) for development workflow
- If feature involves forms: [forms.instructions.md](../instructions/forms.instructions.md)
- If feature involves HTTP requests: [ky.instructions.md](../instructions/ky.instructions.md)
- If feature involves UI components: [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md)
- If feature involves charts: [charts.instructions.md](../instructions/charts.instructions.md)
- If feature involves design tokens: [design-tokens.instructions.md](../instructions/design-tokens.instructions.md)
- If feature involves billing: [stripe.instructions.md](../instructions/stripe.instructions.md)
- If feature involves file uploads: [asset_upload.instructions.md](../instructions/asset_upload.instructions.md)

## Implementation Standards

1. **Architecture**: TanStack Router + TanStack Query + TanStack Store (no Redux, no Zustand)
2. **Routing**: File-based routing in `portal-web/src/routes/`
3. **Tenancy & Scoping (SSOT)**: Business-owned features must be scoped by `businessDescriptor`
   - UI routes under `/business/$businessDescriptor/...`
   - API calls to `v1/businesses/${businessDescriptor}/...`
   - Read `businessDescriptor` via `Route.useParams()` and pass it into API/query hooks
4. **Forms**: Use Kyora form system (`useKyoraForm` + `<form.AppField>` pattern). Follow `forms.instructions.md`
5. **HTTP**: Ky client with proper error handling
6. **UI**: daisyUI components, RTL-first, responsive design
7. **i18n**: All text strings must support Arabic + English
8. **Accessibility**: Proper ARIA labels, keyboard navigation
9. **Type Safety**: Zod schemas for validation + TypeScript types

## Workflow

1. Search for similar features to understand patterns
2. Define Zod schema in `portal-web/src/schemas/`
3. Create/update API client in `portal-web/src/api/`
4. Build route component in `portal-web/src/routes/`
5. Create/update/reuse reusable components in `portal-web/src/components/` that are shared across other resources.
6. Add translations in `portal-web/src/i18n/`
7. Style with daisyUI classes + design tokens
8. Test locally with `npm run lint` and `npm run type-check` in portal-web directory

## Done

- Implementation complete, production-ready
- RTL support verified
- i18n keys added for both locales
- Responsive design verified (with smooth UX for mobile/tablet/desktop)
- Accessibility checked
- No TODOs or FIXMEs
