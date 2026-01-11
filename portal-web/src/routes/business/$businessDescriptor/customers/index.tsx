import { createFileRoute } from '@tanstack/react-router'
import { Suspense } from 'react'

import type { CustomersSearch } from '@/features/customers/schema/customersSearch'
import {
  CustomersListPage,
  customersListLoader,
} from '@/features/customers/components/CustomersListPage'
import { CustomerListSkeleton } from '@/features/customers/components/CustomerListSkeleton'
import { CustomersSearchSchema } from '@/features/customers/schema/customersSearch'
import { RouteErrorFallback } from '@/components/molecules/RouteErrorFallback'

/**
 * Customers List Route
 *
 * Displays list of customers with:
 * - Search/filter functionality (debounced 300ms)
 * - Pagination with URL search params
 * - Responsive table/card views
 * - Empty states with CTA
 * - Create customer in BottomSheet/Modal
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/customers/',
)({
  staticData: {
    titleKey: 'customers.title',
  },
  validateSearch: (search): CustomersSearch => {
    return CustomersSearchSchema.parse(search)
  },

  // Prefetch customer list data based on search params
  loader: async ({ context, params, location }) => {
    const { queryClient } = context as any

    const searchParams = CustomersSearchSchema.parse(location.search)

    await customersListLoader({
      queryClient,
      businessDescriptor: params.businessDescriptor,
      search: searchParams,
    })
  },

  errorComponent: RouteErrorFallback,

  component: () => (
    <Suspense fallback={<CustomerListSkeleton />}>
      <CustomersListPage />
    </Suspense>
  ),
})
