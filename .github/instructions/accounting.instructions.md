---
description: "Kyora accounting SSOT (backend + portal-web): assets, investments, withdrawals, expenses, recurring expenses, summary, automation"
applyTo: "backend/internal/domain/accounting/**,backend/internal/server/routes.go,backend/internal/tests/e2e/accounting_*_test.go,portal-web/src/components/organisms/Sidebar.tsx,portal-web/src/routes/onboarding/plan.tsx,portal-web/src/api/types/onboarding.ts"
---

# Kyora Accounting — Single Source of Truth (SSOT)

This file documents the **accounting behavior implemented today** in Kyora.

Scope:

- Backend: `backend/internal/domain/accounting/**` (source of truth for API + domain behavior)
- Portal Web: **not implemented yet** as a feature module (there is a sidebar nav entry + plan feature flag only)

If you change accounting behavior, keep backend contract + any portal-web consumers consistent.

## Non-negotiables

- **Business-scoped always:** all accounting data is scoped to a business under `/v1/businesses/:businessDescriptor/accounting/...`.
- **No cross-tenant leaks:** never allow access across workspaces/businesses.
- **RBAC on every route:** accounting endpoints are guarded by `role.ResourceAccounting` with `ActionView` vs `ActionManage`.
- **ProblemDetails errors:** backend uses RFC7807 ProblemDetails.
- **Treat notes as plain text:** accounting `note` fields are stored and returned as-is; portal-web must not render them as HTML.

## Backend: API surface (authoritative)

All routes are business-scoped under:

- `/v1/businesses/:businessDescriptor/accounting`

### Assets

- `GET /assets` → `list.ListResponse<Asset>`
  - Query: `page`, `pageSize`, `orderBy[]`
- `GET /assets/:assetId` → `Asset`
- `POST /assets` → `Asset`
- `PATCH /assets/:assetId` → `Asset`
- `DELETE /assets/:assetId` → `204`

### Investments (owner injections)

- `GET /investments` → `list.ListResponse<Investment>`
  - Query: `page`, `pageSize`, `orderBy[]`
- `GET /investments/:investmentId` → `Investment`
- `POST /investments` → `Investment`
- `PATCH /investments/:investmentId` → `Investment`
- `DELETE /investments/:investmentId` → `204`

### Withdrawals (owner draws)

- `GET /withdrawals` → `list.ListResponse<Withdrawal>`
  - Query: `page`, `pageSize`, `orderBy[]`
- `GET /withdrawals/:withdrawalId` → `Withdrawal`
- `POST /withdrawals` → `Withdrawal`
- `PATCH /withdrawals/:withdrawalId` → `Withdrawal`
- `DELETE /withdrawals/:withdrawalId` → `204`

### Expenses

- `GET /expenses` → `list.ListResponse<Expense>`
  - Query: `page`, `pageSize`, `orderBy[]`
  - Default sort: `-occurredOn`
- `GET /expenses/:expenseId` → `Expense` (preloads `RecurringExpense`)
- `POST /expenses` → `Expense`
- `PATCH /expenses/:expenseId` → `Expense`
- `DELETE /expenses/:expenseId` → `204`

### Recurring expenses (templates + occurrences)

- `GET /recurring-expenses` → `list.ListResponse<RecurringExpense>`
  - Query: `page`, `pageSize`, `orderBy[]`
- `GET /recurring-expenses/:recurringExpenseId` → `RecurringExpense` (preloads `Expenses[]`)
- `POST /recurring-expenses` → `RecurringExpense`
- `PATCH /recurring-expenses/:recurringExpenseId` → `RecurringExpense`
- `DELETE /recurring-expenses/:recurringExpenseId` → `204`
- `PATCH /recurring-expenses/:recurringExpenseId/status` → `RecurringExpense`
  - Body: `{ "status": "active"|"paused"|"ended"|"canceled" }`
  - Returns `409` for invalid transitions
- `GET /recurring-expenses/:recurringExpenseId/occurrences` → `Expense[]`
  - Occurrences are `Expense` rows linked via `recurringExpenseId`.

### Summary

- `GET /summary?from=YYYY-MM-DD&to=YYYY-MM-DD` → summary
  - Computes totals and `safeToDrawAmount`.
  - `from/to` are optional; invalid date format returns `400`.

## Backend: list response contract

All list endpoints return `list.ListResponse<T>` with **camelCase** metadata:

- `items`
- `page`, `pageSize`
- `totalCount`, `hasMore`

Pagination is implemented by `offset = (page-1)*pageSize` and `hasMore = page*pageSize < totalCount`.

## Backend: RBAC (repo reality)

- View endpoints use `ActionView`.
- Create/update/delete endpoints use `ActionManage`.

E2E tests confirm **members can view but cannot manage** investments/withdrawals/expenses/recurring-expenses.

## Backend: core data semantics

### Currency

All accounting records store `currency` as the business currency (`biz.Currency`) at creation time.

### Notes (security)

- `note` fields are persisted and returned as raw strings.
- E2E tests intentionally store strings containing HTML/JS. Portal-web must escape when rendering.

### Validation (what is enforced today)

- Recurring expense create enforces:
  - `frequency` ∈ `daily|weekly|monthly|yearly`
  - `category` is one of the allowed enum values
  - `recurringEndDate > recurringStartDate` (if provided)
  - `amount > 0` (service-level check)
- Expenses enforce `recurringExpenseId` when `type=recurring` via request binding (`required_if=Type recurring`).

Important repo reality:

- Investments/withdrawals/assets/one-time expenses are **not** currently enforcing `amount > 0` in service/handler (despite existing error helpers). If you need that rule, add it explicitly and cover with E2E tests.

## Backend: recurring expense status machine

Status values:

- `active`, `paused`, `ended`, `canceled`

Allowed transitions:

- `active` → `paused|ended|canceled`
- `paused` → `active|ended|canceled`
- `ended` → `active|canceled`
- `canceled` → `active`

Invalid transitions return a `409` ProblemDetails (`ErrRecurringExpenseInvalidTransition`).

## Backend: recurring expense occurrences

There are two ways occurrences (expense rows) appear:

1. **Backfill on create** (optional)

- `POST /recurring-expenses` supports `autoCreateHistoricalExpenses=true`.
- Backend creates past `Expense` rows from `recurringStartDate` up to “today”, stepping by frequency.

2. **Ongoing creation** (internal automation)

- `CreateNewRecurringExpenseOccurrence(...)` creates an `Expense` and updates `nextRecurringDate`.

If you add a scheduler/cron worker, it should call service-level helpers and remain business-scoped.

## Backend: transaction fee automation (event-driven)

Accounting listens to `bus.OrderPaymentSucceededTopic`.

When an order payment succeeds:

- Resolve the effective payment-method fee from the business service.
- Compute fee as `orderTotal * feePercent + feeFixed` and round to 2 decimals.
- If fee is positive, call `UpsertTransactionFeeExpenseForOrder(...)`.

Idempotency:

- Transaction fee is represented as an `Expense` with:
  - `category = transaction_fee`
  - `type = one_time`
  - `orderId` set
- A unique constraint on `(business_id, order_id, category)` makes the upsert safe.

This method is intentionally **internal** and does not do actor permission checks.

## Backend: accounting summary and “safe to draw”

`GET /summary` returns:

- `totalAssetValue`
- `totalInvestments`
- `totalWithdrawals`
- `totalExpenses`
- `safeToDrawAmount`
- `currency`
- optional echo of `from`, `to`

Safe-to-draw computation:

- Uses **order revenue and order COGS** (not investments).
- Formula:

$$
\text{safeToDraw} = \text{income} - \text{COGS} - \text{expenses} - \text{withdrawals} - \text{safetyBuffer}
$$

Safety buffer:

- If `biz.SafetyBuffer` is set (non-zero): use it.
- If it is zero: default to **sum of expenses in the last 30 days**, anchored to `to` (if provided) or `now`.
- If computed safe-to-draw is negative: returns `0`.

E2E tests confirm date ranges apply to totals and safe-to-draw.

## Portal Web: repo reality (current)

- There is a sidebar nav item pointing to `/accounting`.
- There are translations and onboarding plan features referencing “accounting”.
- There is **no implemented accounting route tree or API client** in portal-web today.

If you implement portal accounting:

- Add a business-scoped API client (follow patterns in `portal-web/src/api/inventory.ts`).
- Keep list state URL-driven (TanStack Router search params) like customers/inventory.
- Add query keys under `portal-web/src/lib/queryKeys.ts` and invalidate correctly after mutations.
- Use the standard form system (`.github/instructions/forms.instructions.md`) and RTL UI rules (`.github/instructions/ui-implementation.instructions.md`).

## Change checklist (when touching accounting)

Backend:

- Keep all queries business-scoped (`ScopeBusinessID(biz.ID)` or equivalent).
- Keep list responses as `list.ListResponse<T>` with camelCase metadata.
- When adding new invariants (e.g., amount > 0), add/extend E2E tests under `backend/internal/tests/e2e/accounting_*_test.go`.
- For automation (bus/cron), use idempotent patterns (unique constraints + atomic upserts).

Portal Web:

- Treat `note` fields as plain text.
- Reuse existing HTTP client patterns (`.github/instructions/ky.instructions.md`).
- Do not invent new UI patterns; follow existing sheets/tables conventions.
