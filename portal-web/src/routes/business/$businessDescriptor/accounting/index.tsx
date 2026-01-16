import { createFileRoute } from '@tanstack/react-router'

import { accountingQueries } from '@/api/accounting'
import { AccountingDashboard } from '@/features/accounting/components/AccountingDashboard'

/**
 * Accounting Dashboard Route
 *
 * Landing page for the Accounting module.
 * Shows "Safe to Draw" hero stat, summary metrics, and recent activity.
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/accounting/',
)({
  staticData: {
    titleKey: 'pages.accounting',
  },
  loader: async ({ context, params }) => {
    const { queryClient } = context as any

    // Prefetch accounting summary (includes Safe to Draw)
    await queryClient.ensureQueryData(
      accountingQueries.summary(params.businessDescriptor),
    )

    // Prefetch recent expenses (last 5)
    await queryClient.ensureQueryData(
      accountingQueries.expenseList(params.businessDescriptor, {
        page: 1,
        pageSize: 5,
        orderBy: ['-occurredOn'],
      }),
    )
  },
  component: AccountingDashboard,
})
