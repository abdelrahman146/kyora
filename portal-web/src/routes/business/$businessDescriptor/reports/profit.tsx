/**
 * Profit & Earnings Route
 *
 * Route: /business/$businessDescriptor/reports/profit
 * Shows revenue, costs, and profitability (P&L statement)
 *
 * Search params:
 * - asOf: optional ISO date string for historical data
 */
import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

import { profitAndLossQueryOptions } from '@/api/accounting'
import { ProfitEarningsPage } from '@/features/reports/components'

const searchSchema = z.object({
  asOf: z.string().optional(),
})

export const Route = createFileRoute(
  '/business/$businessDescriptor/reports/profit',
)({
  validateSearch: searchSchema,
  staticData: {
    titleKey: 'pages.reports_profit',
  },
  loader: async ({ params, context }) => {
    const { queryClient } = context as any

    // Prefetch profit and loss data
    await queryClient.ensureQueryData(
      profitAndLossQueryOptions(params.businessDescriptor),
    )
  },
  component: ProfitEarningsPage,
})
