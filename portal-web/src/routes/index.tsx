import { createFileRoute } from '@tanstack/react-router'
import { businessQueries } from '@/api/business'
import { RouteErrorFallback } from '@/components/molecules/RouteErrorFallback'
import { HomePage, HomePending } from '@/features/home/components/HomePage'
import { requireAuth } from '@/lib/routeGuards'
import { setBusinesses } from '@/stores/businessStore'

/**
 * Home Route Configuration
 *
 * Redirects authenticated users to their selected business or shows business selection hub.
 * - If user has a previously selected business, redirect to `/business/:descriptor`
 * - Otherwise, show business selection interface
 */
export const Route = createFileRoute('/')({
  beforeLoad: requireAuth,

  pendingComponent: HomePending,

  errorComponent: RouteErrorFallback,

  loader: async ({ context }) => {
    const queryClient = (context as any).queryClient
    if (!queryClient) {
      throw new Error('QueryClient not found in router context')
    }

    // Use businessQueries.list() for type-safe data fetching
    const response = await queryClient.ensureQueryData(businessQueries.list())

    // Extract businesses array from response (ensureQueryData doesn't apply select)
    const businesses = response.businesses

    setBusinesses(businesses)

    return { businesses }
  },

  component: HomePage,
})
