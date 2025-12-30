import { Outlet, createFileRoute } from '@tanstack/react-router'
import { Suspense } from 'react'

import type { RouterContext } from '@/router'

import { businessQueries } from '@/api/business'
import { RouteErrorFallback } from '@/components/molecules/RouteErrorFallback'
import { requireAuth } from '@/lib/routeGuards'
import { selectBusiness } from '@/stores/businessStore'
import { DashboardLayout } from '@/components/templates/DashboardLayout'

/**
 * Business Layout Route
 *
 * Parent layout for all business-scoped routes.
 * - Validates business access by descriptor
 * - Updates businessStore with selected business
 * - Preloads business details into Query cache
 * - Wraps children with DashboardLayout (Sidebar, Header, BottomNav)
 *
 * Note: Query invalidation happens in mutations, not in loaders.
 * This prevents unnecessary refetches of data we just loaded.
 */
export const Route = createFileRoute('/business/$businessDescriptor')({
  beforeLoad: async ({ context, params }) => {
    // Require authentication
    await requireAuth()

    // Cast context to RouterContext to access custom properties
    const { queryClient } = context as RouterContext

    // Validate business access and prefetch business details
    const business = await queryClient.ensureQueryData(
      businessQueries.detail(params.businessDescriptor),
    )

    // Update selected business in store
    selectBusiness(params.businessDescriptor)

    return { business }
  },

  errorComponent: RouteErrorFallback,

  component: BusinessLayout,
})

/**
 * Business Layout Component
 *
 * Wraps business routes with DashboardLayout template.
 */
function BusinessLayout() {
  return (
    <DashboardLayout>
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
