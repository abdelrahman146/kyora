/**
 * OrderListSkeleton Component
 *
 * Loading skeleton for orders list page
 */

export function OrderListSkeleton() {
  return (
    <div className="space-y-4">
      {/* Mobile view skeleton */}
      <div className="block md:hidden space-y-4">
        {Array.from({ length: 5 }).map((_, i) => (
          <div
            key={i}
            className="bg-base-100 border border-base-300 rounded-xl p-4 animate-pulse"
          >
            <div className="flex items-start justify-between gap-3 mb-3">
              <div className="flex-1 space-y-2">
                <div className="h-5 bg-base-300 rounded w-32"></div>
                <div className="flex gap-2">
                  <div className="h-6 bg-base-300 rounded-full w-20"></div>
                  <div className="h-6 bg-base-300 rounded-full w-16"></div>
                </div>
              </div>
              <div className="space-y-1">
                <div className="h-3 bg-base-300 rounded w-12"></div>
                <div className="h-6 bg-base-300 rounded w-20"></div>
              </div>
            </div>
            <div className="flex items-center gap-2 mb-3 p-2 bg-base-200 rounded-lg">
              <div className="w-8 h-8 bg-base-300 rounded-full"></div>
              <div className="h-4 bg-base-300 rounded flex-1"></div>
            </div>
            <div className="grid grid-cols-2 gap-3 mb-3">
              <div className="h-12 bg-base-300 rounded-lg"></div>
              <div className="h-12 bg-base-300 rounded-lg"></div>
            </div>
            <div className="h-16 bg-base-300 rounded-lg mb-3"></div>
            <div className="h-4 bg-base-300 rounded w-24"></div>
          </div>
        ))}
      </div>

      {/* Desktop view skeleton */}
      <div className="hidden md:block">
        <div className="overflow-x-auto">
          <table className="table">
            <thead>
              <tr>
                {Array.from({ length: 7 }).map((_, i) => (
                  <th key={i}>
                    <div className="h-4 bg-base-300 rounded w-full"></div>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {Array.from({ length: 10 }).map((_, i) => (
                <tr key={i}>
                  {Array.from({ length: 7 }).map((_cell, j) => (
                    <td key={j}>
                      <div className="h-4 bg-base-300 rounded w-full"></div>
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
