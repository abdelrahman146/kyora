---
status: draft
created_at: 2026-01-14
updated_at: 2026-01-14
brd_ref: ""
owners:
  - area: portal-web
    agent: UI/UX Designer
stakeholders:
  - name: ""
    role: "PM"
  - name: ""
    role: "Engineering Manager"
areas:
  - portal-web
---

# UX Spec: <Title>

## 0) Inputs & scope

- BRD: <link/path>
- Goals (what should feel better for the user):
  -
- Non-goals:
  -
- Assumptions:
  -

## 1) Reuse Map (Evidence)

List concrete existing files/components to reuse.

- Existing routes/pages:
  -
- Existing feature components:
  -
- Existing shared components (components/lib):
  -

**Do-not-duplicate list** (must reuse these if applicable):
- BottomSheet / sheet patterns:
- Form fields/selects:
- Empty/loading/error state patterns:
- Resource list layout patterns:

## 2) IA + Surfaces

List every UI surface you’re specifying.

- Page: <route>
- Sheet: <name>
- Modal: <name>

## 3) User flows (step-by-step)

Write flows as numbered steps, mobile-first.

### Flow A — <name>

1.
2.
3.

### Flow B — <name>

1.
2.

## 4) Per-surface specification (implementation-ready)

For each surface, specify layout, actions, states, and copy guidance.

### Surface: <name>

- Entry points:
- Layout structure (sections; collapsed/expanded; sticky actions):
- Primary CTA:
  - label:
  - enabled/disabled rules:
  - loading behavior:
- Secondary actions:
- Fields:
  - <field>: required? validation? helper text? input mode? `dir="ltr"`?
- Empty states:
- Loading states:
- Error states (mapped to user-friendly copy categories):
- Success behavior:
- Accessibility notes:

## 5) Responsiveness + RTL rules

- Mobile-first layout rules:
- Tablet/desktop layout rules:
- RTL rules:
- `dir="ltr"` fields:

## 6) Copy & i18n keys inventory

- Namespaces to use (prefer existing):
  -
- New keys needed (en/ar parity required):
  -

## 7) Component gaps / enhancements

### Missing components (if any)

- Component:
  - Why it’s needed:
  - Where it should live (feature vs shared):
  - Reuse candidates (1–2):

### Enhancements to existing components

- Component/file:
  - Enhancement:
  - Call sites to update (to avoid divergence):

## 8) Acceptance checklist (for Engineering Manager)

- [ ] Reuses existing Kyora UI patterns; no parallel “second versions” of flows
- [ ] Calm UI: safe defaults, advanced options collapsed
- [ ] Mobile-first verified (one-handed, sticky primary CTA where appropriate)
- [ ] RTL-first verified; `dir="ltr"` applied where needed
- [ ] All empty/loading/error states specified
- [ ] i18n keys listed; en/ar parity planned
- [ ] Any missing components/enhancements explicitly documented
