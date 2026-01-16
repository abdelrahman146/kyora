/**
 * Accounting API Types & Zod Schemas
 *
 * Based on backend accounting domain models in:
 * backend/internal/domain/accounting/model.go
 *
 * All amounts are returned as strings (decimal representation).
 * Dates are RFC3339 strings.
 */
import { z } from 'zod'

// =============================================================================
// Expense Category Enum
// =============================================================================

export const expenseCategoryEnum = z.enum([
  'rent',
  'marketing',
  'salaries',
  'packaging',
  'software',
  'logistics',
  'transaction_fee',
  'travel',
  'supplies',
  'other',
])

export type ExpenseCategory = z.infer<typeof expenseCategoryEnum>

// =============================================================================
// Expense Type Enum
// =============================================================================

export const expenseTypeEnum = z.enum(['one_time', 'recurring'])

export type ExpenseType = z.infer<typeof expenseTypeEnum>

// =============================================================================
// Recurring Expense Frequency Enum
// =============================================================================

export const recurringExpenseFrequencyEnum = z.enum([
  'daily',
  'weekly',
  'monthly',
  'yearly',
])

export type RecurringExpenseFrequency = z.infer<
  typeof recurringExpenseFrequencyEnum
>

// =============================================================================
// Recurring Expense Status Enum
// =============================================================================

export const recurringExpenseStatusEnum = z.enum([
  'active',
  'paused',
  'ended',
  'canceled',
])

export type RecurringExpenseStatus = z.infer<typeof recurringExpenseStatusEnum>

// =============================================================================
// Expense Schema
// =============================================================================

export const expenseSchema = z.object({
  id: z.string(),
  businessId: z.string(),
  amount: z.string(),
  currency: z.string(),
  category: expenseCategoryEnum,
  type: expenseTypeEnum,
  occurredOn: z.string(), // RFC3339
  note: z.string().nullable(),
  orderId: z.string().nullable(),
  recurringExpenseId: z.string().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type Expense = z.infer<typeof expenseSchema>

// =============================================================================
// Recurring Expense Schema
// =============================================================================

export const recurringExpenseSchema = z.object({
  id: z.string(),
  businessId: z.string(),
  amount: z.string(),
  currency: z.string(),
  category: expenseCategoryEnum,
  frequency: recurringExpenseFrequencyEnum,
  recurringStartDate: z.string(), // RFC3339
  recurringEndDate: z.string().nullable(),
  nextRecurringDate: z.string().nullable(),
  note: z.string().nullable(),
  status: recurringExpenseStatusEnum,
  expenses: z.array(expenseSchema).optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type RecurringExpense = z.infer<typeof recurringExpenseSchema>

// =============================================================================
// Investment Schema (Owner Equity In)
// =============================================================================

export const investmentSchema = z.object({
  id: z.string(),
  businessId: z.string(),
  amount: z.string(),
  currency: z.string(),
  investedAt: z.string(), // RFC3339
  note: z.string().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type Investment = z.infer<typeof investmentSchema>

// =============================================================================
// Withdrawal Schema (Owner Equity Out)
// =============================================================================

export const withdrawalSchema = z.object({
  id: z.string(),
  businessId: z.string(),
  amount: z.string(),
  currency: z.string(),
  withdrawnAt: z.string(), // RFC3339
  note: z.string().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type Withdrawal = z.infer<typeof withdrawalSchema>

// =============================================================================
// Asset Type Enum
// =============================================================================

export const assetTypeEnum = z.enum([
  'software',
  'equipment',
  'vehicle',
  'furniture',
  'other',
])

export type AssetType = z.infer<typeof assetTypeEnum>

// =============================================================================
// Asset Schema (Fixed Assets)
// =============================================================================

export const assetSchema = z.object({
  id: z.string(),
  businessId: z.string(),
  name: z.string(),
  type: assetTypeEnum,
  value: z.string(),
  currency: z.string(),
  purchasedAt: z.string(), // RFC3339
  note: z.string().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type Asset = z.infer<typeof assetSchema>

// =============================================================================
// Accounting Summary Schema
// =============================================================================

export const accountingSummarySchema = z.object({
  totalAssetValue: z.string(),
  totalInvestments: z.string(),
  totalWithdrawals: z.string(),
  totalExpenses: z.string(),
  safeToDrawAmount: z.string(),
  currency: z.string(),
  from: z.string().nullable().optional(),
  to: z.string().nullable().optional(),
})

export type AccountingSummary = z.infer<typeof accountingSummarySchema>

// =============================================================================
// List Response Schema (Generic)
// =============================================================================

export const listResponseSchema = <T extends z.ZodTypeAny>(itemSchema: T) =>
  z.object({
    items: z.array(itemSchema),
    page: z.number(),
    pageSize: z.number(),
    totalCount: z.number(),
    hasMore: z.boolean(),
  })

export type ListResponse<T> = {
  items: Array<T>
  page: number
  pageSize: number
  totalCount: number
  hasMore: boolean
}

// =============================================================================
// Request DTOs
// =============================================================================

// Expense Create Request
export const createExpenseRequestSchema = z.object({
  amount: z.string().min(1),
  category: expenseCategoryEnum,
  occurredOn: z.string(), // YYYY-MM-DD or RFC3339
  note: z.string().optional(),
})

export type CreateExpenseRequest = z.infer<typeof createExpenseRequestSchema>

// Expense Update Request
export const updateExpenseRequestSchema = z.object({
  amount: z.string().optional(),
  category: expenseCategoryEnum.optional(),
  occurredOn: z.string().optional(),
  note: z.string().nullable().optional(),
})

export type UpdateExpenseRequest = z.infer<typeof updateExpenseRequestSchema>

// Recurring Expense Create Request
export const createRecurringExpenseRequestSchema = z.object({
  amount: z.string().min(1),
  category: expenseCategoryEnum,
  frequency: recurringExpenseFrequencyEnum,
  recurringStartDate: z.string(),
  recurringEndDate: z.string().nullable().optional(),
  note: z.string().optional(),
  autoCreateHistoricalExpenses: z.boolean().optional(),
})

export type CreateRecurringExpenseRequest = z.infer<
  typeof createRecurringExpenseRequestSchema
>

// Recurring Expense Update Request
export const updateRecurringExpenseRequestSchema = z.object({
  amount: z.string().optional(),
  category: expenseCategoryEnum.optional(),
  frequency: recurringExpenseFrequencyEnum.optional(),
  recurringEndDate: z.string().nullable().optional(),
  note: z.string().nullable().optional(),
})

export type UpdateRecurringExpenseRequest = z.infer<
  typeof updateRecurringExpenseRequestSchema
>

// Recurring Expense Status Update Request
export const updateRecurringExpenseStatusRequestSchema = z.object({
  status: recurringExpenseStatusEnum,
})

export type UpdateRecurringExpenseStatusRequest = z.infer<
  typeof updateRecurringExpenseStatusRequestSchema
>

// Investment Create Request
export const createInvestmentRequestSchema = z.object({
  investorId: z.string().min(1),
  amount: z.string().min(1),
  investedAt: z.string(), // RFC3339/ISO8601
  note: z.string().optional(),
})

export type CreateInvestmentRequest = z.infer<
  typeof createInvestmentRequestSchema
>

// Investment Update Request
export const updateInvestmentRequestSchema = z.object({
  investorId: z.string().optional(),
  amount: z.string().optional(),
  investedAt: z.string().optional(), // RFC3339/ISO8601
  note: z.string().nullable().optional(),
})

export type UpdateInvestmentRequest = z.infer<
  typeof updateInvestmentRequestSchema
>

// Withdrawal Create Request
export const createWithdrawalRequestSchema = z.object({
  withdrawerId: z.string().min(1),
  amount: z.string().min(1),
  withdrawnAt: z.string(), // RFC3339/ISO8601
  note: z.string().optional(),
})

export type CreateWithdrawalRequest = z.infer<
  typeof createWithdrawalRequestSchema
>

// Withdrawal Update Request
export const updateWithdrawalRequestSchema = z.object({
  withdrawerId: z.string().optional(),
  amount: z.string().optional(),
  withdrawnAt: z.string().optional(), // RFC3339/ISO8601
  note: z.string().nullable().optional(),
})

export type UpdateWithdrawalRequest = z.infer<
  typeof updateWithdrawalRequestSchema
>

// Asset Create Request
export const createAssetRequestSchema = z.object({
  name: z.string().min(1),
  type: assetTypeEnum,
  value: z.string().min(1),
  purchasedAt: z.string(),
  note: z.string().optional(),
})

export type CreateAssetRequest = z.infer<typeof createAssetRequestSchema>

// Asset Update Request
export const updateAssetRequestSchema = z.object({
  name: z.string().optional(),
  type: assetTypeEnum.optional(),
  value: z.string().optional(),
  purchasedAt: z.string().optional(),
  note: z.string().nullable().optional(),
})

export type UpdateAssetRequest = z.infer<typeof updateAssetRequestSchema>

// =============================================================================
// Recent Activity Types
// =============================================================================

export const recentActivityTypeEnum = z.enum([
  'expense',
  'investment',
  'withdrawal',
  'asset',
])

export type RecentActivityType = z.infer<typeof recentActivityTypeEnum>

export const recentActivitySchema = z.object({
  id: z.string(),
  type: recentActivityTypeEnum,
  amount: z.string(),
  currency: z.string(),
  description: z.string(),
  occurredAt: z.string(), // RFC3339
  createdAt: z.string(), // RFC3339
  category: expenseCategoryEnum.nullable().optional(), // Only for expenses
  expenseType: expenseTypeEnum.nullable().optional(), // Only for expenses
  personId: z.string().nullable().optional(), // InvestorID or WithdrawerID
})

export type RecentActivity = z.infer<typeof recentActivitySchema>

export const recentActivitiesResponseSchema = z.object({
  items: z.array(recentActivitySchema),
})

export type RecentActivitiesResponse = z.infer<
  typeof recentActivitiesResponseSchema
>
