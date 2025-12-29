/**
 * CustomerDetailSkeleton Component
 *
 * Content-aware skeleton for customer detail page.
 * Matches actual layout: header + profile card + sections
 */

export function CustomerDetailSkeleton() {
  return (
    <div className="space-y-6 animate-pulse">
      {/* Header Skeleton */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-4">
          <div className="h-10 w-10 bg-base-300 rounded-full" />
          <div className="space-y-2">
            <div className="h-8 w-48 bg-base-300 rounded" />
            <div className="h-4 w-32 bg-base-300 rounded" />
          </div>
        </div>
        <div className="flex gap-2">
          <div className="h-9 w-20 bg-base-300 rounded" />
          <div className="h-9 w-20 bg-base-300 rounded" />
        </div>
      </div>

      {/* Profile Card Skeleton */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="flex items-start gap-6">
            <div className="w-24 h-24 bg-base-300 rounded-full" />
            <div className="flex-1 space-y-4">
              <div className="space-y-2">
                <div className="h-6 w-40 bg-base-300 rounded" />
                <div className="h-4 w-32 bg-base-300 rounded" />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {Array.from({ length: 6 }).map((_, i) => (
                  <div key={i} className="space-y-1">
                    <div className="h-3 w-20 bg-base-300 rounded" />
                    <div className="h-4 w-full bg-base-300 rounded" />
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Stats Cards Skeleton */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="h-4 w-24 bg-base-300 rounded mb-2" />
              <div className="h-8 w-32 bg-base-300 rounded" />
            </div>
          </div>
        ))}
      </div>

      {/* Orders Section Skeleton */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="h-6 w-32 bg-base-300 rounded mb-4" />
          <div className="space-y-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="flex items-center justify-between p-3 bg-base-200 rounded">
                <div className="flex-1 space-y-2">
                  <div className="h-4 w-40 bg-base-300 rounded" />
                  <div className="h-3 w-24 bg-base-300 rounded" />
                </div>
                <div className="h-6 w-20 bg-base-300 rounded" />
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
