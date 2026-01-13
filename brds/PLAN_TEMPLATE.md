---
status: draft
created_at: 2026-01-13
updated_at: 2026-01-13
brd_ref: ""
owners:
  - area: backend
    agent: Feature Builder
  - area: portal-web
    agent: Feature Builder
  - area: tests
    agent: Feature Builder
risk_level: medium
---

# Engineering Plan: <Title>

## 0) Inputs

- BRD: <link/path>
- Assumptions:
  -

## 1) Confirmation Gate (must be approved before implementation)

List anything that requires explicit approval:

- New dependency/library?
- New project/app?
- Breaking change?
- Migration?
- Data model change with customer impact?

## 2) Architecture summary (high level)

- Backend: 
- Portal-web:
- Data model:
- Security/tenancy:

## 3) Work breakdown (handoff-ready)

### Milestone 1 (shippable)

- Goal:
- Backend tasks:
- Portal-web tasks:
- Tests:
- Rollout notes:

### Milestone 2 (shippable)

- ...

## 4) API contracts (high level)

- Endpoints:
- DTOs:
- Error cases:

## 5) Data model & migrations

- Tables/models:
- Indexing:
- Migration plan:

## 6) Security & privacy

- Tenant scoping:
- RBAC:
- Abuse prevention:

## 7) Observability & KPIs

- Events/metrics:
- Dashboards/alerts (if any):

## 8) Test strategy

- What is covered by E2E:
- What is covered by integration tests:
- Edge cases:

## 9) Risks & mitigations

- Risk:
  - Mitigation:

## 10) Definition of done

- [ ] Meets BRD acceptance criteria
- [ ] Mobile-first UX verified
- [ ] RTL/i18n parity verified
- [ ] Multi-tenancy verified
- [ ] Error handling + empty/loading states complete
- [ ] No TODO/FIXME
