import {
  createFileRoute,
  useParams,
  useRouteContext,
} from '@tanstack/react-router'

import { CapitalListPage } from '@/features/accounting/components/CapitalListPage'

/**
 * Capital Management Route
 *
 * Manage capital transactions: Investments & Withdrawals.
 * Tracks owner equity movements in/out of the business.
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/accounting/capital',
)({
  staticData: {
    titleKey: 'accounting:header.capital',
  },
  component: CapitalRoute,
})

function CapitalRoute() {
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/accounting/capital',
  })
  const { business } = useRouteContext({
    from: '/business/$businessDescriptor',
  })

  return (
    <CapitalListPage
      businessDescriptor={businessDescriptor}
      currency={business.currency}
    />
  )
}
