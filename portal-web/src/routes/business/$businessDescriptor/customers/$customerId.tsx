import { createFileRoute } from '@tanstack/react-router'
import { Suspense } from 'react'

import type { RouterContext } from '@/router'
import {
  CustomerDetailPage,
  customerDetailLoader,
} from '@/features/customers/components/CustomerDetailPage'
import { CustomerDetailSkeleton } from '@/features/customers/components/CustomerDetailSkeleton'
import { RouteErrorFallback } from '@/components/molecules/RouteErrorFallback'

/**
 * Customer Detail Route
 *
 * Portal-web parity:
 * - Profile card + stats
 * - Social handles
 * - Basic info + addresses management
 * - Notes section
 * - Edit + Delete actions
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/customers/$customerId',
)({
  staticData: {
    titleKey: 'customers.details_title',
  },

  loader: async ({ context, params }) => {
    const { queryClient } = context as unknown as RouterContext

    await customerDetailLoader({
      queryClient,
      businessDescriptor: params.businessDescriptor,
      customerId: params.customerId,
    })
  },

  errorComponent: RouteErrorFallback,

  component: () => (
    <Suspense fallback={<CustomerDetailSkeleton />}>
      <CustomerDetailPage />
    </Suspense>
  ),
})
