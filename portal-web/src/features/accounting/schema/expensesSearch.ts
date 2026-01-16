/**
 * Expenses Search Schema
 *
 * URL-driven state for the Expenses List page.
 * Validated via TanStack Router `validateSearch`.
 *
 * Note: Backend does NOT support search parameter for expenses.
 * We include it here for potential future use but the UI will disable the search input.
 */
import { z } from 'zod'

import { expenseCategoryEnum } from '@/api/types/accounting'

export const ExpensesSearchSchema = z.object({
  page: z.number().optional().default(1),
  pageSize: z.number().optional().default(20),
  sortBy: z.string().optional().default('occurredOn'),
  sortOrder: z.enum(['asc', 'desc']).optional().default('desc'),
  category: expenseCategoryEnum.optional(),
  from: z.string().optional(), // YYYY-MM-DD
  to: z.string().optional(), // YYYY-MM-DD
})

export type ExpensesSearch = z.infer<typeof ExpensesSearchSchema>
