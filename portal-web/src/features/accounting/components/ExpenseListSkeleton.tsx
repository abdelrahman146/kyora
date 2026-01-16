/**
 * ExpenseListSkeleton Component
 *
 * Loading skeleton for the Expenses List page.
 * Matches the visual structure of ExpenseCard for consistent loading UX.
 */

export function ExpenseListSkeleton() {
  return (
    <div className="space-y-4">
      {/* Header skeleton */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <div className="skeleton h-8 w-32 mb-2" />
          <div className="skeleton h-4 w-48" />
        </div>
        <div className="skeleton h-10 w-36" />
      </div>

      {/* Toolbar skeleton */}
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="skeleton h-10 flex-1" />
        <div className="flex gap-2">
          <div className="skeleton h-10 w-24" />
          <div className="skeleton h-10 w-24" />
        </div>
      </div>

      {/* Cards skeleton */}
      <div className="space-y-3">
        {Array.from({ length: 5 }).map((_, index) => (
          <ExpenseCardSkeleton key={index} />
        ))}
      </div>

      {/* Pagination skeleton */}
      <div className="flex items-center justify-between">
        <div className="skeleton h-4 w-32" />
        <div className="flex gap-2">
          <div className="skeleton h-8 w-8" />
          <div className="skeleton h-8 w-8" />
          <div className="skeleton h-8 w-8" />
        </div>
      </div>
    </div>
  )
}

function ExpenseCardSkeleton() {
  return (
    <div className="bg-base-100 border border-base-300 rounded-xl p-4">
      <div className="flex items-start gap-3">
        {/* Icon skeleton */}
        <div className="skeleton w-10 h-10 rounded-lg flex-shrink-0" />

        {/* Content skeleton */}
        <div className="flex-1 min-w-0">
          <div className="skeleton h-5 w-24 mb-2" />
          <div className="skeleton h-4 w-48 mb-2" />
          <div className="skeleton h-3 w-20" />
        </div>

        {/* Amount skeleton */}
        <div className="skeleton h-6 w-20 flex-shrink-0" />
      </div>
    </div>
  )
}
