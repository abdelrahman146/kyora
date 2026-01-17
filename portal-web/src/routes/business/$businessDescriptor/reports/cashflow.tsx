/**
 * Cash Movement Route
 *
 * Route: /business/$businessDescriptor/reports/cashflow
 * Shows cash inflows, outflows, and net cash position
 *
 * Search params:
 * - asOf: optional ISO date string for historical data
 */
import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

import { cashFlowQueryOptions } from '@/api/accounting'
import { CashMovementPage } from '@/features/reports/components'

const searchSchema = z.object({
  asOf: z.string().optional(),
})

export const Route = createFileRoute(
  '/business/$businessDescriptor/reports/cashflow',
)({
  validateSearch: searchSchema,
  staticData: {
    titleKey: 'pages.reports_cashflow',
  },
  loader: async ({ params, context }) => {
    const { queryClient } = context as any

    // Prefetch cash flow data
    await queryClient.ensureQueryData(
      cashFlowQueryOptions(params.businessDescriptor),
    )
  },
  component: CashMovementPage,
})
