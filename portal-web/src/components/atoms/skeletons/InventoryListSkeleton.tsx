/**
 * InventoryListSkeleton Component
 *
 * Content-aware skeleton for inventory list page.
 * Matches actual layout: header + search bar + filters + table/cards
 */

import { useMediaQuery } from '@/hooks/useMediaQuery'

export function InventoryListSkeleton() {
  const isMobile = useMediaQuery('(max-width: 768px)')

  return (
    <div className="space-y-6 animate-pulse">
      {/* Header Skeleton */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-2">
          <div className="h-8 w-32 bg-base-300 rounded" />
          <div className="h-4 w-56 bg-base-300 rounded" />
        </div>
        <div className="h-[52px] w-36 bg-base-300 rounded-xl" />
      </div>

      {/* Search Bar Skeleton */}
      <div className="h-[50px] w-full bg-base-300 rounded-xl" />

      {/* Filter Button Skeleton */}
      <div className="flex gap-3">
        <div className="h-[50px] w-40 bg-base-300 rounded-xl" />
      </div>

      {/* Desktop Table Skeleton */}
      {!isMobile && (
        <div className="card bg-base-100 border border-base-300 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="table">
              {/* Header */}
              <thead>
                <tr>
                  <th>
                    <div className="h-4 w-32 bg-base-300 rounded" />
                  </th>
                  <th className="text-center">
                    <div className="h-4 w-20 bg-base-300 rounded mx-auto" />
                  </th>
                  <th className="text-center">
                    <div className="h-4 w-24 bg-base-300 rounded mx-auto" />
                  </th>
                  <th className="text-center">
                    <div className="h-4 w-16 bg-base-300 rounded mx-auto" />
                  </th>
                  <th className="text-center">
                    <div className="h-4 w-16 bg-base-300 rounded mx-auto" />
                  </th>
                  <th className="text-center">
                    <div className="h-4 w-20 bg-base-300 rounded mx-auto" />
                  </th>
                </tr>
              </thead>
              <tbody>
                {Array.from({ length: 10 }).map((_, i) => (
                  <tr key={i}>
                    <td>
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-base-300 rounded-full" />
                        <div className="space-y-2">
                          <div className="h-4 w-40 bg-base-300 rounded" />
                          <div className="h-3 w-24 bg-base-300 rounded" />
                        </div>
                      </div>
                    </td>
                    <td>
                      <div className="flex justify-center">
                        <div className="h-6 w-20 bg-base-300 rounded-full" />
                      </div>
                    </td>
                    <td>
                      <div className="flex justify-center">
                        <div className="h-4 w-20 bg-base-300 rounded" />
                      </div>
                    </td>
                    <td>
                      <div className="flex justify-center">
                        <div className="h-4 w-12 bg-base-300 rounded" />
                      </div>
                    </td>
                    <td>
                      <div className="flex justify-center">
                        <div className="h-6 w-16 bg-base-300 rounded-full" />
                      </div>
                    </td>
                    <td>
                      <div className="flex justify-center">
                        <div className="h-8 w-8 bg-base-300 rounded-lg" />
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Mobile Card Grid Skeleton */}
      {isMobile && (
        <div className="space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <div
              key={i}
              className="card bg-base-100 border border-base-300"
            >
              <div className="card-body p-4">
                <div className="flex items-center gap-3 mb-3">
                  <div className="w-10 h-10 bg-base-300 rounded-full" />
                  <div className="flex-1 space-y-2">
                    <div className="h-4 w-32 bg-base-300 rounded" />
                    <div className="h-3 w-24 bg-base-300 rounded" />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-2">
                    <div className="h-3 w-20 bg-base-300 rounded" />
                    <div className="h-4 w-16 bg-base-300 rounded" />
                  </div>
                  <div className="space-y-2">
                    <div className="h-3 w-16 bg-base-300 rounded" />
                    <div className="h-4 w-20 bg-base-300 rounded" />
                  </div>
                </div>
                <div className="mt-3">
                  <div className="h-6 w-full bg-base-300 rounded-full" />
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Pagination Skeleton */}
      <div className="flex items-center justify-between gap-4">
        <div className="h-10 w-32 bg-base-300 rounded" />
        <div className="flex gap-2">
          <div className="h-10 w-10 bg-base-300 rounded" />
          <div className="h-10 w-10 bg-base-300 rounded" />
          <div className="h-10 w-10 bg-base-300 rounded" />
        </div>
        <div className="h-10 w-32 bg-base-300 rounded" />
      </div>
    </div>
  )
}
