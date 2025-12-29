import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { z } from 'zod'
import { Plus, Search, Users } from 'lucide-react'
import type {CreateCustomerRequest} from '@/api/customer';
import {
  
  useCreateCustomerMutation,
  useCustomersQuery
} from '@/api/customer'
import { CustomerCard } from '@/components/molecules/CustomerCard'
import { Pagination } from '@/components/molecules/Pagination'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { CustomerForm } from '@/components/organisms/CustomerForm'
import { CustomerListSkeleton } from '@/components/atoms/skeletons/CustomerListSkeleton'
import { useMediaQuery } from '@/hooks/useMediaQuery'
import { queryKeys } from '@/lib/queryKeys'
import { showErrorFromException, showSuccessToast } from '@/lib/toast'

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
 * - Search/filter functionality (debounced 300ms)
 * - Pagination with URL search params
 * - Responsive table/card views
 * - Empty states with CTA
 * - Create customer in BottomSheet/Modal
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/customers/',
)({
  validateSearch: (search): CustomersSearch => {
    return CustomersSearchSchema.parse(search)
  },

  component: () => (
    <Suspense fallback={<CustomerListSkeleton />}>
      <CustomersListPage />
    </Suspense>
  ),
})

/**
 * Customers List Page Component
 */
function CustomersListPage() {
  const { t } = useTranslation(['common', 'errors'])
  const { businessDescriptor } = Route.useParams()
  const navigate = useNavigate()
  const search = Route.useSearch()
  const queryClient = useQueryClient()
  const isMobile = useMediaQuery('(max-width: 768px)')

  // Local search state for debouncing
  const [searchInput, setSearchInput] = useState(search.search || '')
  const [showCreateSheet, setShowCreateSheet] = useState(false)

  // Fetch customers with search params
  const { data, isLoading, error } = useCustomersQuery(businessDescriptor, {
    search: search.search,
    page: search.page,
    limit: search.limit,
  })

  // Create customer mutation
  const createMutation = useCreateCustomerMutation(businessDescriptor, {
    onSuccess: (newCustomer) => {
      // Invalidate customers list to refetch with new customer
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.list(businessDescriptor),
      })
      
      showSuccessToast(t('common:customer.created_success'))
      setShowCreateSheet(false)
      
      // Navigate to new customer detail page
      void navigate({
        to: '/business/$businessDescriptor/customers/$customerId',
        params: { businessDescriptor, customerId: newCustomer.id },
      })
    },
    onError: (err) => {
      void showErrorFromException(err, t)
    },
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

  // Handle create customer form submission
  const handleCreateCustomer = async (values: CreateCustomerRequest) => {
    await createMutation.mutateAsync(values)
  }

  // Handle pagination
  const handlePageChange = (newPage: number) => {
    void navigate({
      to: '.',
      search: { ...search, page: newPage },
    })
  }

  if (error) {
    return (
      <div className="flex min-h-[400px] flex-col items-center justify-center gap-4">
        <p className="text-error">{t('errors:generic.load_failed')}</p>
        <button
          className="btn btn-sm"
          onClick={() => window.location.reload()}
        >
          {t('common:actions.retry')}
        </button>
      </div>
    )
  }

  const customers = data?.customers ?? []
  const pagination = data?.pagination

  return (
    <>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-2xl font-bold">{t('common:customer.customers')}</h1>
            <p className="text-sm text-base-content/70">
              {t('common:customer.manage_database')}
            </p>
          </div>
          <button
            className="btn btn-primary btn-sm sm:btn-md gap-2"
            onClick={() => setShowCreateSheet(true)}
          >
            <Plus size={18} />
            {t('common:customer.add_customer')}
          </button>
        </div>

        {/* Search Bar */}
        <div className="card bg-base-100 shadow">
          <div className="card-body p-4">
            <div className="relative">
              <Search
                size={20}
                className="absolute start-3 top-1/2 -translate-y-1/2 text-base-content/40"
              />
              <input
                type="text"
                placeholder={t('common:customer.search_placeholder')}
                className="input input-bordered w-full ps-10"
                value={searchInput}
                onChange={(e) => handleSearchChange(e.target.value)}
              />
            </div>
          </div>
        </div>

        {/* Empty State */}
        {!isLoading && customers.length === 0 && (
          <div className="card bg-base-100 shadow">
            <div className="card-body items-center text-center py-12">
              <Users size={48} className="text-base-content/20 mb-4" />
              <h3 className="text-lg font-semibold mb-2">
                {search.search
                  ? t('common:customer.no_results')
                  : t('common:customer.no_customers')}
              </h3>
              <p className="text-sm text-base-content/70 mb-4">
                {search.search
                  ? t('common:customer.try_different_search')
                  : t('common:customer.get_started_message')}
              </p>
              {!search.search && (
                <button
                  className="btn btn-primary btn-sm gap-2"
                  onClick={() => setShowCreateSheet(true)}
                >
                  <Plus size={18} />
                  {t('common:customer.add_first_customer')}
                </button>
              )}
            </div>
          </div>
        )}

        {/* Desktop Table View */}
        {!isMobile && customers.length > 0 && (
          <div className="card bg-base-100 shadow overflow-hidden">
            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>{t('common:customer.name')}</th>
                    <th>{t('common:customer.email')}</th>
                    <th>{t('common:customer.phone')}</th>
                    <th>{t('common:customer.orders')}</th>
                    <th>{t('common:customer.total_spent')}</th>
                    <th>{t('common:actions.actions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {customers.map((customer) => (
                    <tr key={customer.id} className="hover">
                      <td>
                        <div className="flex items-center gap-3">
                          <div className="avatar placeholder">
                            <div className="w-10 h-10 bg-primary/10 text-primary rounded-full">
                              <span className="text-sm font-medium">
                                {customer.fullName.charAt(0).toUpperCase()}
                              </span>
                            </div>
                          </div>
                          <span className="font-medium">{customer.fullName}</span>
                        </div>
                      </td>
                      <td>
                        <span className="text-sm text-base-content/70">
                          {customer.email || '—'}
                        </span>
                      </td>
                      <td>
                        <span className="text-sm text-base-content/70">
                          {customer.phonePrefix && customer.phoneNumber
                            ? `${customer.phonePrefix} ${customer.phoneNumber}`
                            : '—'}
                        </span>
                      </td>
                      <td>
                        <span className="badge badge-ghost">{customer.totalOrders}</span>
                      </td>
                      <td>
                        <span className="font-medium">
                          {customer.totalSpent.toFixed(2)}
                        </span>
                      </td>
                      <td>
                        <button
                          className="btn btn-ghost btn-xs"
                          onClick={() => {
                            void navigate({
                              to: '/business/$businessDescriptor/customers/$customerId',
                              params: { businessDescriptor, customerId: customer.id },
                            })
                          }}
                        >
                          {t('common:actions.view')}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Mobile Card View */}
        {isMobile && customers.length > 0 && (
          <div className="grid gap-4">
            {customers.map((customer) => (
              <CustomerCard
                key={customer.id}
                customer={{
                  id: customer.id,
                  name: customer.fullName,
                  email: customer.email ?? undefined,
                  phone:
                    customer.phonePrefix && customer.phoneNumber
                      ? `${customer.phonePrefix} ${customer.phoneNumber}`
                      : undefined,
                  phoneCode: customer.phonePrefix ?? undefined,
                  phoneNumber: customer.phoneNumber ?? undefined,
                  totalOrders: customer.totalOrders,
                  totalSpent: customer.totalSpent,
                }}
                onClick={(cust) => {
                  void navigate({
                    to: '/business/$businessDescriptor/customers/$customerId',
                    params: { businessDescriptor, customerId: cust.id },
                  })
                }}
              />
            ))}
          </div>
        )}

        {/* Pagination */}
        {pagination && pagination.totalPages > 1 && (
          <div className="flex justify-center">
            <Pagination
              currentPage={pagination.page}
              totalPages={pagination.totalPages}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </div>

      {/* Create Customer BottomSheet/Modal */}
      <BottomSheet
        isOpen={showCreateSheet}
        onClose={() => setShowCreateSheet(false)}
        title={t('common:customer.add_customer')}
      >
        <CustomerForm
          onSubmit={handleCreateCustomer}
          onCancel={() => setShowCreateSheet(false)}
          isSubmitting={createMutation.isPending}
        />
      </BottomSheet>
    </>
  )
}
