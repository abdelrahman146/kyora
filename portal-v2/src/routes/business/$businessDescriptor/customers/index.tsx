import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useState } from 'react'
import { z } from 'zod'

import { useCustomersQuery } from '@/api/customer'

/**
 * Customers List Route Search Params Schema
 */
const CustomersSearchSchema = z.object({
  search: z.string().optional(),
  page: z.number().optional().default(1),
  limit: z.number().optional().default(20),
})

type CustomersSearch = z.infer<typeof CustomersSearchSchema>

/**
 * Customers List Route
 *
 * Displays list of customers with:
 * - Search/filter functionality (debounced)
 * - Pagination
 * - Responsive table/card views
 * - Empty states
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/customers/',
)({
  validateSearch: (search): CustomersSearch => {
    return CustomersSearchSchema.parse(search)
  },

  component: () => (
    <Suspense fallback={<CustomersListSkeleton />}>
      <CustomersListPage />
    </Suspense>
  ),
})

/**
 * Customers List Page Component
 */
function CustomersListPage() {
  const { businessDescriptor } = Route.useParams()
  const { business } = Route.useRouteContext()
  const navigate = useNavigate()
  const search = Route.useSearch()

  // Local search state for debouncing
  const [searchInput, setSearchInput] = useState(search.search || '')

  // Fetch customers with search params
  const { data, isLoading, error } = useCustomersQuery(businessDescriptor, {
    search: search.search,
    page: search.page,
    limit: search.limit,
  })

  // Handle search input change (debounced)
  const handleSearchChange = (value: string) => {
    setSearchInput(value)

    // Debounce search by 300ms
    const timer = setTimeout(() => {
      void navigate({
        to: '.',
        search: { search: value || undefined, page: 1, limit: search.limit },
      })
    }, 300)

    return () => clearTimeout(timer)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">العملاء</h1>
          <p className="text-sm text-base-content/70">
            إدارة قاعدة بيانات العملاء
          </p>
        </div>
        <button className="btn btn-primary btn-sm sm:btn-md">
          + إضافة عميل
        </button>
      </div>

      {/* Search and Filters */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="flex flex-col gap-4 sm:flex-row">
            <div className="form-control flex-1">
              <input
                type="text"
                placeholder="البحث عن عميل..."
                className="input input-bordered w-full"
                value={searchInput}
                onChange={(e) => handleSearchChange(e.target.value)}
              />
            </div>
            <button className="btn btn-outline btn-sm sm:btn-md">
              تصفية
            </button>
          </div>
        </div>
      </div>

      {/* Customers List */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          {isLoading ? (
            // Loading skeleton
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="flex items-center gap-4">
                  <div className="skeleton h-12 w-12 shrink-0 rounded-full"></div>
                  <div className="flex-1 space-y-2">
                    <div className="skeleton h-4 w-full"></div>
                    <div className="skeleton h-3 w-2/3"></div>
                  </div>
                </div>
              ))}
            </div>
          ) : error ? (
            // Error state
            <div className="flex min-h-[300px] flex-col items-center justify-center gap-4">
              <p className="text-error">حدث خطأ في تحميل البيانات</p>
              <button className="btn btn-sm" onClick={() => window.location.reload()}>
                إعادة المحاولة
              </button>
            </div>
          ) : !data || data.customers.length === 0 ? (
            // Empty state
            <div className="flex min-h-[300px] flex-col items-center justify-center gap-4">
              <div className="text-center">
                <h3 className="mb-2 text-lg font-semibold">لا يوجد عملاء</h3>
                <p className="mb-4 text-base-content/70">
                  ابدأ بإضافة أول عميل لك
                </p>
                <button className="btn btn-primary btn-sm">+ إضافة عميل</button>
              </div>
            </div>
          ) : (
            // Customers table
            <>
              <div className="overflow-x-auto">
                <table className="table">
                  <thead>
                    <tr>
                      <th>الاسم</th>
                      <th>البريد الإلكتروني</th>
                      <th>الهاتف</th>
                      <th>الطلبات</th>
                      <th>الإجمالي</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.customers.map((customer) => (
                      <tr key={customer.id} className="hover">
                        <td>
                          <div className="flex items-center gap-3">
                            <div className="avatar placeholder">
                              <div className="w-12 rounded-full bg-neutral text-neutral-content">
                                <span className="text-xl">
                                  {customer.fullName.charAt(0)}
                                </span>
                              </div>
                            </div>
                            <div>
                              <div className="font-bold">{customer.fullName}</div>
                            </div>
                          </div>
                        </td>
                        <td>
                          {customer.email || (
                            <span className="text-base-content/50">-</span>
                          )}
                        </td>
                        <td>
                          {customer.phoneNumber ? (
                            <span>
                              {customer.phonePrefix} {customer.phoneNumber}
                            </span>
                          ) : (
                            <span className="text-base-content/50">-</span>
                          )}
                        </td>
                        <td>{customer.totalOrders}</td>
                        <td>
                          {customer.totalSpent.toFixed(2)} {business.currency}
                        </td>
                        <td>
                          <a
                            href={`/business/${businessDescriptor}/customers/${customer.id}`}
                            className="btn btn-ghost btn-xs"
                          >
                            عرض
                          </a>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {data.pagination.totalPages > 1 && (
                <div className="mt-4 flex justify-center">
                  <div className="join">
                    <button
                      className="btn join-item btn-sm"
                      disabled={search.page === 1}
                      onClick={() =>
                        navigate({
                          to: '.',
                          search: {
                            search: search.search,
                            page: search.page - 1,
                            limit: search.limit,
                          },
                        })
                      }
                    >
                      السابق
                    </button>
                    <button className="btn join-item btn-sm">
                      صفحة {search.page} من {data.pagination.totalPages}
                    </button>
                    <button
                      className="btn join-item btn-sm"
                      disabled={search.page >= data.pagination.totalPages}
                      onClick={() =>
                        navigate({
                          to: '.',
                          search: {
                            search: search.search,
                            page: search.page + 1,
                            limit: search.limit,
                          },
                        })
                      }
                    >
                      التالي
                    </button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}

/**
 * Customers List Skeleton
 *
 * Content-aware skeleton matching customers list structure
 */
function CustomersListSkeleton() {
  return (
    <div className="space-y-6">
      {/* Header Skeleton */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-2">
          <div className="skeleton h-8 w-32"></div>
          <div className="skeleton h-4 w-48"></div>
        </div>
        <div className="skeleton h-10 w-32"></div>
      </div>

      {/* Search and Filters Skeleton */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="flex flex-col gap-4 sm:flex-row">
            <div className="skeleton h-12 flex-1"></div>
            <div className="skeleton h-12 w-24"></div>
          </div>
        </div>
      </div>

      {/* Table Skeleton */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="space-y-4">
            {[...Array(8)].map((_, i) => (
              <div key={i} className="flex items-center gap-4">
                <div className="skeleton h-12 w-12 shrink-0 rounded-full"></div>
                <div className="flex-1 space-y-2">
                  <div className="skeleton h-4 w-full"></div>
                  <div className="skeleton h-3 w-2/3"></div>
                </div>
                <div className="skeleton h-8 w-16"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
