---
description: Portal Web code structure SSOT (routes/components/features/lib), no-legacy + strict placement rules
applyTo: "portal-web/**"
---

# Portal Web Code Structure (SSOT)

This is the SSOT for how code must be organized in `portal-web/**`.

**SSOT Hierarchy**

- Parent orchestration: `.github/copilot-instructions.md`
- Portal architecture (stack + patterns): `.github/instructions/portal-web-architecture.instructions.md`
- State ownership: `.github/instructions/state-management.instructions.md`
- HTTP + TanStack Query: `.github/instructions/http-tanstack-query.instructions.md`
- Forms system: `.github/instructions/forms.instructions.md`
- UI/RTL rules: `.github/instructions/ui-implementation.instructions.md`

---

## 0) Non‑negotiables

1. **Routes are routing only.**

- `portal-web/src/routes/**` is the routing layer.
- Route files may define:
  - `Route = createFileRoute(...)({ ... })` config
  - `validateSearch` Zod schema
  - `loader` / `beforeLoad` / `loaderDeps`
  - `staticData` (e.g. title keys; route page titles must use `common.pages.*` — see `.github/instructions/i18n-translations.instructions.md`)
  - a thin `component` wrapper that renders a Feature Page component
- Route files must not:
  - implement “page UI” (tables/cards/sheets/layouts)
  - contain domain/resource business logic
  - define feature-specific components (OrdersListPage, InventoryPage, CustomerDetailsPage, etc.)

2. **Shared components must be generic.**

`portal-web/src/components/**` is for reusable UI that is not bound to a specific resource/feature.

If the component name or props mention a resource (e.g. `Order`, `Customer`, `Product`, `ShippingZone`, `Onboarding`) it is almost always **feature-specific** and must live under `features/<feature>/...`.

3. **Features own their code end-to-end.**

A **feature** is a cohesive slice of functionality. It can be:

- **Resource/domain feature**: a business-owned resource and all of its UI + forms + feature-specific components.

  - Examples: `orders`, `inventory`, `customers`, `analytics`, `accounting`.

- **Cross-cutting/global functionality feature**: global-ish client behavior/state plus components used across many pages.

  - Examples: `auth` (auth widgets + flows), `business-switcher` (selected business UX), `language` (language switcher + i18n helpers), `onboarding` (flow pages + steps).

- **Layout/template feature**: a layout that has its own complex UI composition, local context/state, and supporting components.
  - Examples: `dashboard-layout` (or `app-shell`), `resource-list-layout`.

Important rule: a feature can be _cross-functional_ (used by many other features) and still be a **feature**.
That does **not** make it a shared atomic component.

4. **No-legacy policy.**

- Never create `*.old.tsx`, `*.bak`, or duplicate “v2” implementations.
- If you replace something, you must migrate callers and delete the old code in the same change.

5. **DRY policy.**

- Before creating any new component/utility, search the repo and reuse existing ones.
- Duplicate meaning is drift: one concept should have one implementation.

---

## 1) Target folder layout (required going forward)

### 1.1 Routing layer

- `portal-web/src/routes/**`
  - Only routing config + thin wrappers.
  - Pages live inside `portal-web/src/features/**`.

Example pattern:

- Route file:
  - defines search schema + loader prefetch
  - exports `component: OrdersRoute`
- Feature page:
  - `portal-web/src/features/orders/components/OrdersListPage.tsx`

### 1.2 Shared UI components

All shared components must be resource-agnostic.

- `portal-web/src/components/atoms/`

  - Small UI primitives.
  - No TanStack Query, no route params, no feature logic.
  - Allowed: generic `Button`, `Badge`, `Input`, `Dialog`, `Tooltip`, `Skeleton`, etc.

- `portal-web/src/components/molecules/`

  - Compositions of atoms.
  - Still generic and reusable.
  - Allowed examples: `SearchInput`, `Pagination`, `ConfirmDialog`, `BottomSheet`.

- `portal-web/src/components/organisms/`

  - Higher-level reusable composites.
  - Still generic. If it’s feature/resource-specific, it does not belong here.
  - Allowed examples: `Table`, `Header`, `Sidebar` (app chrome), generic filters/sorting.

- `portal-web/src/components/templates/`

  - Layout-level wrappers with slots.
  - Must not bake in resource-specific behavior.

- `portal-web/src/components/form/`

  - Generic UI form controls **not** tied to TanStack Form.
  - Controls should be usable with any form library (TanStack Form, RHF, uncontrolled).

- `portal-web/src/components/charts/`

  - Generic, reusable Chart.js components.
  - Must not be analytics-specific; receive `data`/`options` via props.

- `portal-web/src/components/icons/`
  - Custom icons only when `lucide-react` doesn’t provide what we need.

**Strict naming rule:** shared component names must not include resource names.

### 1.3 Feature modules

Create features under:

- `portal-web/src/features/<feature>/...`

#### Feature types (allowed)

- Resource features (business-scoped): `features/orders`, `features/inventory`, ...
- Global/cross-cutting features: `features/auth`, `features/business-switcher`, `features/language`, ...
- Layout/template features: `features/dashboard-layout`, `features/app-shell`, ...

If a module has its own state + components + forms, it is a strong sign it should be a feature.

Allowed subfolders (all optional):

- `features/<feature>/state/` feature-local store (route-bound preferred)
- `features/<feature>/components/` feature-specific components and pages
- `features/<feature>/forms/` feature-specific forms
- `features/<feature>/schema/` feature-specific zod schemas
- `features/<feature>/utils/` feature-specific helpers
- `features/<feature>/types/` feature-specific types

Rules:

- Feature components may compose shared components.
- Feature code may use Query hooks from `portal-web/src/api/**`.
- Feature code may be reused across the app by importing from `features/<feature>/*`.
  - This is the preferred home for cross-cutting functionality (auth/business-switcher/language).
- Feature code should not export project-wide “misc utilities”.
  - If it is truly cross-cutting and not feature-vocabulary-specific, it belongs in `lib/`.
  - If it is cross-cutting but still clearly “this feature’s vocabulary”, keep it in that feature.

#### Cross-cutting feature rule (strict)

If a component is reused across many pages but is still feature-specific (e.g. BusinessSwitcher, LanguageSwitcher, Auth widgets), it must live under:

- `features/<feature>/components/**`

It must **not** live under `components/atoms|molecules|organisms`.

Shared `components/**` is reserved for resource-agnostic UI primitives and compositions.

---

## 2) `lib/` rules (strict)

`portal-web/src/lib/**` is reserved for cross-cutting utilities only.

A utility belongs in `lib/` only if all are true:

- Used by 2+ distinct features OR clearly cross-cutting (e.g. error parsing, query keys, routing guards)
- Not tied to a single resource vocabulary (inventory/orders/customers/etc.)
- Has stable API and is expected to be shared

If a file is feature-specific (e.g. `inventoryUtils`, `onboarding`, `customers*`), it must live under `features/<feature>/utils/`.

---

## 3) How to decide: shared vs feature vs route

Use this decision tree:

1. Does it talk to the backend or represent server data?

- Put it in `api/` (already SSOT).

2. Is it required to render a route and should be shareable/bookmarkable?

- Put state in URL/search params (Router).
- Put the page implementation in `features/<feature>/components/`.

3. Is it a UI component used across multiple features and not resource-specific?

- Put it in `components/**` (atoms/molecules/organisms/templates).

4. Is it specific to one resource/workflow?

- Put it under `features/<feature>/...`.

---

## 4) Enforcement notes for AI agents

- When adding a new page:

  - Create a route wrapper in `routes/**`.
  - Implement the page in `features/<feature>/components/**`.

- When moving code:

  - Move + update imports in the same change.
  - Delete the old location immediately (no transitional duplicates).

- If current code violates this SSOT, do not copy the pattern.
  - Follow this SSOT and log drift in `DRIFT_TODO.md`.
