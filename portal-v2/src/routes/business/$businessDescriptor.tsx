import { Outlet, createFileRoute } from '@tanstack/react-router'
import { Suspense } from 'react'

import type { RouterContext } from '@/router'

import { businessApi } from '@/api/business'
import { invalidateBusinessScopedQueries } from '@/lib/queryInvalidation'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'
import { requireAuth } from '@/lib/routeGuards'
import { selectBusiness } from '@/stores/businessStore'
import { DashboardLayout } from '@/components/templates/DashboardLayout'

/**
 * Business Layout Route
 *
 * Parent layout for all business-scoped routes.
 * - Validates business access by descriptor
 * - Updates businessStore with selected business
 * - Invalidates all business-scoped queries on business switch
 * - Wraps children with DashboardLayout (Sidebar, Header, BottomNav)
 */
export const Route = createFileRoute('/business/$businessDescriptor')({
  beforeLoad: async ({ context, params }) => {
    // Require authentication
    requireAuth()

    // Cast context to RouterContext to access custom properties
    const { queryClient } = context as RouterContext

    // Validate business access
    const business = await queryClient.ensureQueryData({
      queryKey: queryKeys.businesses.detail(params.businessDescriptor),
      queryFn: () => businessApi.getBusiness(params.businessDescriptor),
      staleTime: STALE_TIME.FIVE_MINUTES,
    })

    // Update selected business in store
    selectBusiness(params.businessDescriptor)

    // Invalidate all business-scoped queries to ensure fresh data
    invalidateBusinessScopedQueries(queryClient)

    return { business }
  },

  component: BusinessLayout,
})

/**
 * Business Layout Component
 *
 * Wraps business routes with DashboardLayout template.
 */
function BusinessLayout() {
  const { businessDescriptor } = Route.useParams()

  return (
    <DashboardLayout businessDescriptor={businessDescriptor}>
      {/* Content outlet with Suspense boundary */}
      <Suspense
        fallback={
          <div className="flex min-h-[400px] items-center justify-center">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        }
      >
        <Outlet />
      </Suspense>
    </DashboardLayout>
  )
}
