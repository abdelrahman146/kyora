import { createFileRoute } from '@tanstack/react-router'

import {
  ExpenseListPage,
  ExpensesSearchSchema,
  expenseListLoader,
} from '@/features/accounting'

/**
 * Expenses List Route
 *
 * Displays paginated list of expenses with:
 * - Desktop: Table view
 * - Mobile: Card view
 * - Filters: Category, Date Range
 * - Add Expense sheet
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/accounting/expenses/',
)({
  staticData: {
    titleKey: 'pages.expenses',
  },
  validateSearch: ExpensesSearchSchema.parse,
  loaderDeps: ({ search }) => ({
    page: search.page,
    pageSize: search.pageSize,
    sortBy: search.sortBy,
    sortOrder: search.sortOrder,
    category: search.category,
    from: search.from,
    to: search.to,
  }),
  loader: async ({ context, params, deps }) => {
    const { queryClient } = context as any
    const { businessDescriptor } = params

    await expenseListLoader({
      queryClient,
      businessDescriptor,
      search: {
        page: deps.page,
        pageSize: deps.pageSize,
        sortBy: deps.sortBy,
        sortOrder: deps.sortOrder,
        category: deps.category,
        from: deps.from,
        to: deps.to,
      },
    })
  },
  component: ExpenseListPage,
})
