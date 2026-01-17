/**
 * Business Health Route
 *
 * Route: /business/$businessDescriptor/reports/health
 * Shows financial position (what the business owns, owes, and is worth)
 *
 * Search params:
 * - asOf: optional ISO date string for historical data
 */
import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

import { financialPositionQueryOptions } from '@/api/accounting'
import { BusinessHealthPage } from '@/features/reports/components'

const searchSchema = z.object({
  asOf: z.string().optional(),
})

export const Route = createFileRoute(
  '/business/$businessDescriptor/reports/health',
)({
  validateSearch: searchSchema,
  staticData: {
    titleKey: 'pages.reports_health',
  },
  loader: async ({ params, context }) => {
    const { queryClient } = context as any

    // Prefetch financial position data
    await queryClient.ensureQueryData(
      financialPositionQueryOptions(params.businessDescriptor),
    )
  },
  component: BusinessHealthPage,
})
