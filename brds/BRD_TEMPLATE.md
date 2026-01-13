---
status: draft
owner: product
created_at: 2026-01-13
updated_at: 2026-01-13
stakeholders:
  - name: ""
    role: ""
  - name: ""
    role: ""
areas:
  - portal-web
  - backend
  - storefront-web
kpis:
  - name: ""
    definition: ""
    baseline: ""
    target: ""
---

# BRD: <Title>

## 1) Problem (in plain language)

- What is the user trying to do?
- What is painful/confusing today?
- Why does it matter for Kyora’s customers (mobile-first, Arabic-first, DM commerce)?

## 2) Customer & Context

- Primary user persona(s):
- Where does this happen (Instagram/WhatsApp/etc.)?
- Device context (mobile-first):
- Language direction: RTL/Arabic-first requirements:

## 3) Goals (what success looks like)

- Goal 1:
- Goal 2:

## 4) Non-goals (explicitly out of scope)

- Non-goal 1:
- Non-goal 2:

## 5) User journey (happy path)

1.
2.
3.

## 6) Edge cases & failure handling

- Case:
  - Expected behavior:
  - What the user sees (plain language):

## 7) UX / IA (mobile-first)

### Pages / Surfaces

For each page/surface:

- Purpose:
- Primary action:
- Secondary actions:
- Content (what must be shown):
- Empty state:
- Loading state:
- Error state:
- i18n keys needed (en/ar parity):

### Copy principles

- Use plain language (avoid accounting jargon)
- Prefer actionable CTAs (“Add product”, “Send to WhatsApp”, “Mark as paid”)

## 8) Functional requirements

Write as testable bullets:

- FR-1:
- FR-2:

## 9) Data & permissions

- Tenant scoping (workspace + business):
- Roles (admin/member):
- What must never leak across tenants:

## 10) Analytics & KPIs

- Event(s) to track:
- KPI impact expectation:

## 11) Rollout & risks

- Rollout plan:
- Risks:
- Mitigations:

## 12) Open questions

- Q1:
- Q2:

## 13) Acceptance criteria (definition of done)

- [ ] Works end-to-end on mobile
- [ ] RTL/Arabic parity verified
- [ ] Clear empty/loading/error states
- [ ] No confusing jargon
- [ ] KPIs/events defined
- [ ] Multi-tenant safety respected
