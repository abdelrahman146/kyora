import { createFileRoute } from '@tanstack/react-router'

import {
  RecurringExpenseListPage,
  RecurringExpensesSearchSchema,
  recurringExpenseListLoader,
} from '@/features/accounting'

/**
 * Recurring Expenses List Route
 *
 * Displays paginated list of recurring expense templates with:
 * - Desktop: Table view
 * - Mobile: Card view
 * - Filters: Status
 * - Status management (pause, resume, end, cancel)
 * - Edit/Delete actions
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/accounting/expenses/recurring/',
)({
  staticData: {
    titleKey: 'pages.recurring_expenses',
  },
  validateSearch: RecurringExpensesSearchSchema.parse,
  loaderDeps: ({ search }) => ({
    page: search.page,
    pageSize: search.pageSize,
    sortBy: search.sortBy,
    sortOrder: search.sortOrder,
    status: search.status,
  }),
  loader: async ({ context, params, deps }) => {
    const { queryClient } = context as any
    const { businessDescriptor } = params

    await recurringExpenseListLoader({
      queryClient,
      businessDescriptor,
      search: {
        page: deps.page,
        pageSize: deps.pageSize,
        sortBy: deps.sortBy,
        sortOrder: deps.sortOrder,
        status: deps.status,
      },
    })
  },
  component: RecurringExpenseListPage,
})
