/**
 * Reports Hub Route
 *
 * Landing page for financial reports showing:
 * - Safe to Draw hero metric (from accounting summary)
 * - 3 report cards (Business Health, Profit & Earnings, Cash Movement)
 *
 * URL: /business/:businessDescriptor/reports
 * Search params: asOf (optional date filter, YYYY-MM-DD)
 *
 * Prefetches all required data in parallel for instant page loads.
 * Uses existing query options from accounting API.
 */
import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

import {
  accountingQueries,
  cashFlowQueryOptions,
  financialPositionQueryOptions,
  profitAndLossQueryOptions,
} from '@/api/accounting'
import { ReportsHubPage } from '@/features/reports/components/ReportsHubPage'

const searchSchema = z.object({
  asOf: z.string().optional(),
})

export const Route = createFileRoute('/business/$businessDescriptor/reports/')({
  validateSearch: searchSchema,
  staticData: {
    titleKey: 'pages.reports',
  },
  loader: async ({ context, params }) => {
    const { businessDescriptor } = params

    const { queryClient } = context as any

    // Prefetch all required data in parallel for instant page loads
    // Note: asOf param is optional; when omitted, backend uses current date
    await Promise.all([
      queryClient.ensureQueryData(
        accountingQueries.summary(businessDescriptor),
      ),
      queryClient.ensureQueryData(
        financialPositionQueryOptions(businessDescriptor),
      ),
      queryClient.ensureQueryData(
        profitAndLossQueryOptions(businessDescriptor),
      ),
      queryClient.ensureQueryData(cashFlowQueryOptions(businessDescriptor)),
    ])
  },
  component: ReportsHubPage,
})
