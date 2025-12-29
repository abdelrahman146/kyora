import { Outlet, createFileRoute } from '@tanstack/react-router'
import { Suspense } from 'react'

import type { RouterContext } from '@/router'

import { businessApi } from '@/api/business'
import { invalidateBusinessScopedQueries } from '@/lib/queryInvalidation'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'
import { requireAuth } from '@/lib/routeGuards'
import { selectBusiness } from '@/stores/businessStore'

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
 * Wraps business routes with dashboard layout.
 * TODO: Implement DashboardLayout template (Sidebar, Header, BottomNav)
 */
function BusinessLayout() {
  const { business } = Route.useRouteContext()

  return (
    <div className="min-h-screen bg-base-200">
      {/* TODO: Replace with DashboardLayout template from Step 8 */}
      <div className="container mx-auto px-4 py-8">
        {/* Temporary header - will be replaced with DashboardLayout */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold">{business.name}</h1>
          <p className="text-sm text-base-content/70">
            {business.country} • {business.currency}
          </p>
        </div>

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
      </div>
    </div>
  )
}
