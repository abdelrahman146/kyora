---
description: "Kyora analytics SSOT (backend + portal-web): dashboards, sales/inventory/customer analytics, financial reports, date-range semantics"
applyTo: "**/*"
---

# Kyora Analytics — Single Source of Truth (SSOT)

This file documents **analytics** behavior implemented today across:

- Backend (source of truth): `backend/internal/domain/analytics/**` + route wiring in `backend/internal/server/routes.go`
- Portal Web (current status): analytics UI/API is mostly not implemented yet; portal currently only has nav + i18n scaffolding.

If you change analytics semantics (metric definitions, date range defaults, response JSON shapes, or RBAC), update backend + portal-web together.

## Non-negotiables

- **Business-scoped analytics:** analytics routes are under `/v1/businesses/:businessDescriptor/...` and must apply `business.EnforceBusinessValidity` (business is loaded from workspace + descriptor).
- **Workspace is the tenant:** analytics must not accept `workspaceId` from clients.
- **RBAC is required:** analytics endpoints enforce `role.ActionView` on either `role.ResourceBasicAnalytics` or `role.ResourceFinancialReports`.

## Backend: route surface (authoritative)

All analytics endpoints are business-scoped under:

- `/v1/businesses/:businessDescriptor/analytics/*`

These routes are registered in `backend/internal/server/routes.go` under the business-scoped group.

### Dashboard

- `GET /v1/businesses/:businessDescriptor/analytics/dashboard`
  - Permission: `role.ActionView` on `role.ResourceBasicAnalytics`
  - Returns: `DashboardMetrics`

### Sales analytics

- `GET /v1/businesses/:businessDescriptor/analytics/sales`
  - Permission: `role.ActionView` on `role.ResourceBasicAnalytics`
  - Query params:
    - `from` (optional) date string `YYYY-MM-DD`
    - `to` (optional) date string `YYYY-MM-DD`
  - Defaults:
    - if `to` missing → now (UTC)
    - if `from` missing → `to - 30 days`
    - validation: `to >= from`
  - Returns: `SalesAnalytics`

### Inventory analytics

- `GET /v1/businesses/:businessDescriptor/analytics/inventory`
  - Permission: `role.ActionView` on `role.ResourceBasicAnalytics`
  - Query params: `from`, `to` (`YYYY-MM-DD`, same defaulting rules)
  - Returns: `InventoryAnalytics`

### Customer analytics

- `GET /v1/businesses/:businessDescriptor/analytics/customers`
  - Permission: `role.ActionView` on `role.ResourceBasicAnalytics`
  - Query params: `from`, `to` (`YYYY-MM-DD`, same defaulting rules)
  - Returns: `CustomerAnalytics`

### Financial reports

These are nested under `/analytics/reports` and require the stronger permission:

- Permission: `role.ActionView` on `role.ResourceFinancialReports`

Routes:

- `GET /v1/businesses/:businessDescriptor/analytics/reports/financial-position`
- `GET /v1/businesses/:businessDescriptor/analytics/reports/profit-and-loss`
- `GET /v1/businesses/:businessDescriptor/analytics/reports/cash-flow`

Query params:

- `asOf` (optional) date string `YYYY-MM-DD`
  - default: today (UTC)

## Backend: date parsing and range semantics

Analytics uses **date-only** query parameters (not RFC3339).

- Expected format: `YYYY-MM-DD` (Go layout `2006-01-02`)
- If the date string is invalid, backend returns `400` with message: `invalid <field> date format, use YYYY-MM-DD`.

Defaulting for `from/to`-based endpoints:

- `to` defaults to `time.Now().UTC()`
- `from` defaults to `to.AddDate(0, 0, -30)`

## Backend: metric definitions (as implemented)

The analytics service composes metrics from other domains:

- Orders domain: revenue/COGS, order counts, funnels, top products, per-channel/country breakdowns.
- Inventory domain: stock totals, low/out-of-stock counts, inventory value, top products by inventory value.
- Customer domain: customer counts/time series and retrieving top customers.
- Accounting domain: safe-to-draw, expenses sums, assets/investments/withdrawals.

### DashboardMetrics

- `revenueLast30Days`: sum of order totals for the last 30 days.
- `grossProfitLast30Days`: last-30-day revenue minus last-30-day COGS.
- `openOrdersCount`: count of open orders.
- `lowStockItemsCount`: count of low-stock variants.
- `allTimeRevenue`: sum of order totals across all time.
- `safeToDrawAmount`: `accounting.ComputeSafeToDrawAmount(...)` using all-time revenue and all-time COGS.
- `salesPerformanceLast30Days`: revenue time series over last 30 days.
- `liveOrderFunnel`: distribution of live (non-completed) orders by stage.
- `topSellingProducts`: top 5 products by sales.
- `newCustomersTimeSeries`: new customers per day over last 30 days.

### SalesAnalytics

Computed for the selected `[from, to]` range:

- `totalRevenue`: sum of order totals.
- `grossProfit`: revenue minus COGS.
- `totalOrders`, `averageOrderValue`, `itemsSold`.
- Time series:
  - `numberOfSalesOverTime` (orders count)
  - `revenueOverTime`
- Breakdowns:
  - `orderStatusBreakdown`
  - `salesByCountry`
  - `salesByChannel`
- `topSellingProducts`: top 5 products.

### InventoryAnalytics

- `totalInventoryValue`: current inventory value (valued at cost).
- `totalInStock`: sum of variant stock quantities.
- `lowStockItems`, `outOfStockItems`.
- Ratios:
  - `inventoryTurnoverRatio`: `COGS(from,to) / avgInventory`, where `avgInventory` is approximated as current inventory value.
  - `sellThroughRate`: `itemsSold / (itemsSold + currentStockUnits)`.
- `topProductsByInventoryValue`: top 5.

### CustomerAnalytics

- `newCustomers`: customers created in the period.
- `returningCustomers`: number of returning customers in the period (as computed by orders domain).
- `repeatCustomerRate`: `returningCustomers / uniquePurchasers` where `uniquePurchasers = len(ordersByCustomer)`.
- `averageRevenuePerCustomer`: `totalRevenue / uniquePurchasers`.
- `customerAcquisitionCost`: marketing spend / new customers.
  - marketing spend is computed via `accounting.SumExpensesAmountByCategory(..., ExpenseCategoryMarketing, from, to)`.
- `customerLifetimeValue`: `averageOrderValue(all-time) * averagePurchaseFrequency(period)`.
- `averageCustomerPurchaseFrequency`: `totalOrders(period) / uniquePurchasers`.
- `newCustomersOverTime`: customers time series.
- `topCustomersByRevenue`: top 5 customers by revenue.

### Financial reports (current approximations)

These reports are currently **inception-to-date** (from `time.Time{}` to `asOf`), not period reports.

- Financial position (`ComputeFinancialPosition`):

  - Retained earnings: `revenue - cogs - expenses`.
  - Cash on hand approximation:
    - `cash = (revenue + ownerInvestment) - (expenses + ownerDraws + fixedAssets + inventoryValue)`.
  - Liabilities are not tracked yet; set to zero.

- Cash flow (`ComputeCashFlow`):
  - Uses the same cash approximation inputs as financial position.
  - Assumes `cashAtStart = 0`, and `cashAtEnd = netCashFlow` for inception-to-date.

## Backend: time series JSON shape

Time series values are returned as:

- `{ granularity: <enum>, series: [{ timestamp, label, value }] }`

Where:

- `granularity` is one of: `hourly|daily|weekly|monthly|quarterly|yearly`.
- `timestamp` is a timestamp.
- `label` is a backend-generated human label for chart ticks.
- `value` is a number.

Portal should prefer using `label` for axis ticks and `timestamp` for sorting.

## Portal Web: alignment guidance (even if not implemented yet)

### Routing (business-scoped)

Implement analytics routes as children under the business layout route:

- `portal-web/src/routes/business/$businessDescriptor/analytics/index.tsx` (overview)
- Optional deeper pages for parity with backend:
  - `sales.tsx`, `inventory.tsx`, `customers.tsx`, `reports/*`

Avoid a global `/analytics` route that is not business-scoped.

### API client

Implement a dedicated analytics API module:

- `portal-web/src/api/analytics.ts`
- `portal-web/src/api/types/analytics.ts`

Rules:

- Use `get()` from `portal-web/src/api/client.ts`.
- Use Zod schemas that match backend JSON exactly.
- Encode date filters as `YYYY-MM-DD` strings.
- Co-locate `queryOptions` factories (like other domains) for route loaders + components reuse.

### Chart rendering

For chart UI and Chart.js patterns, follow the SSOT in:

- `.github/instructions/charts.instructions.md`

Do not re-implement separate chart conventions inside analytics.

### i18n

Portal already has an `analytics` translation namespace and should keep analytics UI labels there.

## Known portal drift / gaps (track before building on it)

- Portal currently has an Analytics nav item but no matching analytics routes or API clients.
- Ensure date filter UX aligns to analytics backend date-only contract (`YYYY-MM-DD`).
