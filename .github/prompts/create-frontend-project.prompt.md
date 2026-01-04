---
description: Create a new frontend project structure (React feature module)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
---

# Create Frontend Project

You are creating a new feature module in the Kyora portal-web (React TanStack dashboard).

## Project Requirements

${input:feature:What feature module are you creating? (e.g., "subscriptions management", "notification center")}

## Instructions

Read frontend architecture rules:

- [portal-web-architecture.instructions.md](../instructions/portal-web-architecture.instructions.md) for complete architecture
- [portal-web-development.instructions.md](../instructions/portal-web-development.instructions.md) for workflow
- [forms.instructions.md](../instructions/forms.instructions.md) for form patterns
- [ui-implementation.instructions.md](../instructions/ui-implementation.instructions.md) for UI patterns

## Tenancy & Scoping (SSOT)

- Workspace can contain multiple businesses.
- Business-owned domains (Orders, Inventory, Customers, Analytics, Accounting, Assets, Storefront, Onboarding) must be scoped by `businessDescriptor`.
- Route convention: `/business/$businessDescriptor/...`.
- API convention: `v1/businesses/${businessDescriptor}/...`.

## Feature Module Structure

Create/extend the following areas in `portal-web/` (follow existing route group conventions and avoid inventing new folder systems):

```
portal-web/src/
├── routes/_app/${feature}/         # App routes (typical location; match existing layout groups)
│   ├── index.tsx                  # List view
│   ├── create.tsx                 # Create form (if needed)
│   └── $id.tsx                    # Detail/edit view (if needed)
├── api/${feature}.ts               # API client methods
├── api/types/${feature}.ts         # API types
├── schemas/${feature}.ts           # Zod validation schemas
├── components/{atoms|molecules|organisms|templates}/  # Reuse existing component taxonomy
└── i18n/{ar|en}/*.json             # Add keys to existing namespace files (don’t create new root JSONs)
```

## Implementation Workflow

1. **Define Schema** (`schemas/${feature}.ts`)

   - Zod schema for validation
   - TypeScript types inferred from schema
   - Match backend API contract

2. **API Client** (`api/${feature}.ts`)

   - Ky HTTP client methods (GET, POST, PUT, DELETE)
   - Proper error handling
   - Type-safe request/response

3. **Routes** (`routes/_app/${feature}/`)

   - **index.tsx**: List view with TanStack Query
   - **create.tsx**: Create form with TanStack Form
   - **$id.tsx**: Detail/edit with params and loader

4. **Components** (`components/{atoms|molecules|organisms|templates}/`)

   - Prefer reusing existing atoms/molecules/organisms/templates; add new ones only when necessary
   - Forms must follow `forms.instructions.md` (Kyora form system + `<form.AppField>` pattern)
   - Apply daisyUI semantics + design tokens (no custom colors)

5. **Translations** (`i18n/`)

   - Add keys to both `portal-web/src/i18n/ar/*.json` and `portal-web/src/i18n/en/*.json` in the appropriate namespace
   - Use semantic keys: `${feature}.title`, `${feature}.create`, etc.

6. **Tests**
   - Component tests with React Testing Library
   - Test user interactions and edge cases
   - Mock API calls

## Architecture Patterns

- **Routing**: File-based TanStack Router, use loaders for data fetching
- **State**: TanStack Query for server state; TanStack Store for client state (no Redux, no Zustand)
- **Forms**: TanStack Form + Zod validation
- **HTTP**: Ky client with proper error handling
- **UI**: daisyUI components, RTL-first, responsive
- **i18n**: All text must have translation keys

## Example Reference

Look at existing features for patterns:

- `routes/_app/orders/` - Full CRUD example
- `routes/_app/customers/` - Simple resource
- `routes/_app/inventory/` - Complex forms

## Done

- Complete feature module structure created
- Routes implemented with proper loaders
- API client created with type safety
- Forms implemented with validation
- Components styled with daisyUI + design tokens
- Translations added for both locales
- Tests written and passing
- RTL support verified
- Responsive design verified
- No TODOs or FIXMEs
- Production-ready code
