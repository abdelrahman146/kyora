/**
 * CustomerListSkeleton Component
 *
 * Content-aware skeleton for customer list page.
 */

export function CustomerListSkeleton() {
  return (
    <div className="space-y-6 animate-pulse">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-2">
          <div className="h-8 w-32 bg-base-300 rounded" />
          <div className="h-4 w-48 bg-base-300 rounded" />
        </div>
        <div className="h-10 w-32 bg-base-300 rounded" />
      </div>

      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="h-10 w-full bg-base-300 rounded" />
        </div>
      </div>

      <div className="hidden md:block card bg-base-100 shadow overflow-hidden">
        <div className="overflow-x-auto">
          <table className="table">
            <thead>
              <tr>
                <th>
                  <div className="h-4 w-24 bg-base-300 rounded" />
                </th>
                <th>
                  <div className="h-4 w-20 bg-base-300 rounded" />
                </th>
                <th>
                  <div className="h-4 w-24 bg-base-300 rounded" />
                </th>
                <th>
                  <div className="h-4 w-20 bg-base-300 rounded" />
                </th>
                <th>
                  <div className="h-4 w-20 bg-base-300 rounded" />
                </th>
                <th>
                  <div className="h-4 w-16 bg-base-300 rounded" />
                </th>
              </tr>
            </thead>
            <tbody>
              {Array.from({ length: 5 }).map((_, i) => (
                <tr key={i}>
                  <td>
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-base-300 rounded-full" />
                      <div className="h-4 w-32 bg-base-300 rounded" />
                    </div>
                  </td>
                  <td>
                    <div className="h-4 w-28 bg-base-300 rounded" />
                  </td>
                  <td>
                    <div className="h-4 w-24 bg-base-300 rounded" />
                  </td>
                  <td>
                    <div className="h-4 w-16 bg-base-300 rounded" />
                  </td>
                  <td>
                    <div className="h-4 w-20 bg-base-300 rounded" />
                  </td>
                  <td>
                    <div className="h-8 w-20 bg-base-300 rounded" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="md:hidden grid gap-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-12 h-12 bg-base-300 rounded-full" />
                <div className="flex-1 space-y-2">
                  <div className="h-5 w-32 bg-base-300 rounded" />
                  <div className="h-3 w-24 bg-base-300 rounded" />
                </div>
              </div>
              <div className="space-y-2">
                <div className="h-4 w-full bg-base-300 rounded" />
                <div className="h-4 w-3/4 bg-base-300 rounded" />
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="flex justify-center gap-2">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="h-10 w-10 bg-base-300 rounded" />
        ))}
      </div>
    </div>
  )
}
