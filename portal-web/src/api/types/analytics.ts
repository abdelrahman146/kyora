/**
 * Analytics API Types
 *
 * TypeScript types for analytics and financial report API responses.
 * Based on backend models at: backend/internal/domain/analytics/model.go
 *
 * See: .github/instructions/analytics.instructions.md
 */

// =============================================================================
// Financial Position (Balance Sheet)
// =============================================================================

/**
 * Financial Position represents a business's assets, liabilities, and equity.
 * All monetary values are returned as decimal strings to avoid floating-point issues.
 */
export interface FinancialPosition {
  businessID: string
  asOf: string // ISO date string
  totalAssets: string // decimal as string
  totalLiabilities: string
  totalEquity: string
  cashOnHand: string
  totalInventoryValue: string
  currentAssets: string
  fixedAssets: string
  ownerInvestment: string
  retainedEarnings: string
  ownerDraws: string
}

// =============================================================================
// Profit and Loss Statement
// =============================================================================

/**
 * Key-value pair for expense category breakdown
 * Note: Backend returns Key/Value with capital letters (Go struct without json tags)
 */
export interface ExpenseCategoryValue {
  Key: string // expense category identifier
  Value: string // decimal amount as string
}

/**
 * Profit and Loss Statement shows revenue, costs, and profitability.
 * All monetary values are returned as decimal strings.
 */
export interface ProfitAndLossStatement {
  businessID: string
  asOf: string // ISO date string
  grossProfit: string
  totalExpenses: string
  netProfit: string
  revenue: string
  cogs: string // Cost of Goods Sold
  expensesByCategory: Array<ExpenseCategoryValue>
}

// =============================================================================
// Cash Flow Statement
// =============================================================================

/**
 * Cash Flow Statement tracks cash movement through the business.
 * All monetary values are returned as decimal strings.
 */
export interface CashFlowStatement {
  businessID: string
  asOf: string // ISO date string
  cashAtStart: string
  cashAtEnd: string
  cashFromCustomers: string
  cashFromOwner: string
  totalCashIn: string
  inventoryPurchases: string
  operatingExpenses: string
  totalBusinessOperation: string
  businessInvestments: string
  ownerDraws: string
  totalCashOut: string
  netCashFlow: string
}
