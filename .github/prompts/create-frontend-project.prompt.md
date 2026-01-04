---
description: Create a new frontend project structure (React feature module)
agent: agent
tools: ["vscode", "execute", "read", "edit", "search", "web", "agent", "todo"]
model: Claude Opus 4.5 (copilot)
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

## Feature Module Structure

Create the following structure in portal-web:

```
portal-web/src/
├── routes/${feature}/        # TanStack Router routes
│   ├── index.tsx            # List view
│   ├── create.tsx           # Create form
│   └── $id.tsx              # Detail/edit view
├── api/${feature}.ts         # API client methods
├── schemas/${feature}.ts     # Zod validation schemas
├── components/${feature}/    # Feature-specific components
│   ├── ${Feature}Form.tsx   # Main form component
│   ├── ${Feature}Card.tsx   # Display component
│   └── ${Feature}List.tsx   # List component
└── i18n/
    ├── en.json              # English translations
    └── ar.json              # Arabic translations
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

3. **Routes** (`routes/${feature}/`)

   - **index.tsx**: List view with TanStack Query
   - **create.tsx**: Create form with TanStack Form
   - **$id.tsx**: Detail/edit with params and loader

4. **Components** (`components/${feature}/`)

   - Reusable components for the feature
   - Use form field components from `components/forms/fields/`
   - Apply daisyUI styling + design tokens

5. **Translations** (`i18n/`)

   - Add keys to both `en.json` and `ar.json`
   - Use semantic keys: `${feature}.title`, `${feature}.create`, etc.

6. **Tests**
   - Component tests with React Testing Library
   - Test user interactions and edge cases
   - Mock API calls

## Architecture Patterns

- **Routing**: File-based TanStack Router, use loaders for data fetching
- **State**: TanStack Query for server state, Zustand for client state
- **Forms**: TanStack Form + Zod validation
- **HTTP**: Ky client with proper error handling
- **UI**: daisyUI components, RTL-first, responsive
- **i18n**: All text must have translation keys

## Example Reference

Look at existing features for patterns:

- `routes/orders/` - Full CRUD example
- `routes/customers/` - Simple resource
- `routes/inventory/` - Complex forms

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
