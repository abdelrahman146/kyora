/**
 * Recurring Expenses Search Schema
 *
 * URL-driven state for the Recurring Expenses List page.
 * Validated via TanStack Router `validateSearch`.
 *
 * Note: Backend does NOT support search parameter for recurring expenses.
 * We include it here for potential future use but the UI will disable the search input.
 */
import { z } from 'zod'

import { recurringExpenseStatusEnum } from '@/api/types/accounting'

export const RecurringExpensesSearchSchema = z.object({
  page: z.number().optional().default(1),
  pageSize: z.number().optional().default(20),
  sortBy: z.string().optional().default('createdAt'),
  sortOrder: z.enum(['asc', 'desc']).optional().default('desc'),
  status: recurringExpenseStatusEnum.optional(),
})

export type RecurringExpensesSearch = z.infer<
  typeof RecurringExpensesSearchSchema
>
