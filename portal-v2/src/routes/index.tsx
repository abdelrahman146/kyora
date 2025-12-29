import { createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'

import { businessApi } from '@/api/business'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'
import { requireAuth } from '@/lib/routeGuards'
import { businessStore, setBusinesses } from '@/stores/businessStore'

/**
 * Home Route
 *
 * Redirects authenticated users to their selected business or shows business selection hub.
 * - If user has a previously selected business, redirect to `/business/:descriptor`
 * - Otherwise, show business selection interface
 */
export const Route = createFileRoute('/')({
  beforeLoad: () => {
    // Require authentication
    requireAuth()
  },

  loader: async ({ context }) => {
    // Fetch user's businesses
    const data = await context.queryClient.ensureQueryData({
      queryKey: queryKeys.businesses.list(),
      queryFn: () => businessApi.listBusinesses(),
      staleTime: STALE_TIME.FIVE_MINUTES,
    })

    // Update businessStore with fetched businesses
    setBusinesses(data.businesses)

    // Check if user has a previously selected business
    const state = businessStore.state
    if (
      state.selectedBusinessDescriptor &&
      data.businesses.some((b) => b.descriptor === state.selectedBusinessDescriptor)
    ) {
      // Redirect to last selected business
      throw redirect({
        to: '/business/$businessDescriptor',
        params: { businessDescriptor: state.selectedBusinessDescriptor },
      })
    }

    // If user has only one business, auto-select it
    if (data.businesses.length === 1) {
      throw redirect({
        to: '/business/$businessDescriptor',
        params: { businessDescriptor: data.businesses[0].descriptor },
      })
    }

    return { businesses: data.businesses }
  },

  component: HomePage,
})

/**
 * Home Page Component
 *
 * Shows business selection hub when user has multiple businesses and no selection.
 */
function HomePage() {
  const { businesses } = Route.useLoaderData()
  const state = useStore(businessStore)

  return (
    <div className="min-h-screen bg-base-200">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8 text-center">
          <h1 className="mb-2 text-3xl font-bold">مرحباً بك في Kyora</h1>
          <p className="text-base-content/70">اختر مشروعك للمتابعة</p>
        </div>

        {businesses.length === 0 ? (
          // No businesses - show onboarding prompt
          <div className="card mx-auto max-w-md bg-base-100 shadow-xl">
            <div className="card-body text-center">
              <h2 className="card-title justify-center">لا توجد مشاريع</h2>
              <p className="text-base-content/70">
                ابدأ بإنشاء مشروعك الأول لإدارة أعمالك
              </p>
              <div className="card-actions justify-center">
                <a href="/onboarding/plan" className="btn btn-primary">
                  إنشاء مشروع جديد
                </a>
              </div>
            </div>
          </div>
        ) : (
          // Show business cards grid
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {businesses.map((business) => (
              <a
                key={business.id}
                href={`/business/${business.descriptor}`}
                className="card bg-base-100 shadow-xl transition-shadow hover:shadow-2xl"
              >
                <div className="card-body">
                  <h2 className="card-title">{business.name}</h2>
                  <div className="flex items-center gap-2 text-sm text-base-content/70">
                    <span>{business.country}</span>
                    <span>•</span>
                    <span>{business.currency}</span>
                  </div>
                  {state.selectedBusinessDescriptor === business.descriptor && (
                    <div className="badge badge-primary badge-sm">النشط</div>
                  )}
                </div>
              </a>
            ))}
          </div>
        )}

        {businesses.length > 0 && (
          <div className="mt-8 text-center">
            <a href="/onboarding/plan" className="btn btn-outline btn-sm">
              + إضافة مشروع جديد
            </a>
          </div>
        )}
      </div>
    </div>
  )
}
