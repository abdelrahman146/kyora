import { createFileRoute, useNavigate } from '@tanstack/react-router'
import {
  Suspense,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { Edit, Eye, Plus, Trash2, Users } from 'lucide-react'
import type { MouseEvent } from 'react'

import type { Customer } from '@/api/customer'
import type { TableColumn } from '@/components/organisms/Table'
import { customerApi } from '@/api/customer'
import { Avatar } from '@/components/atoms/Avatar'
import { Dialog } from '@/components/atoms/Dialog'
import { CustomerListSkeleton } from '@/components/atoms/skeletons/CustomerListSkeleton'
import { CustomerCard } from '@/components/molecules/CustomerCard'
import { InfiniteScroll } from '@/components/molecules/InfiniteScroll'
import { Pagination } from '@/components/molecules/Pagination'
import { SearchInput } from '@/components/molecules/SearchInput'
import { FilterButton } from '@/components/organisms/FilterButton'
import { Table } from '@/components/organisms/Table'
import { AddCustomerSheet, EditCustomerSheet } from '@/components/organisms/customers'
import { useMediaQuery } from '@/hooks/useMediaQuery'
import { showErrorFromException, showSuccessToast } from '@/lib/toast'
import { getSelectedBusiness } from '@/stores/businessStore'

/**
 * Customers List Route Search Params Schema
 */
const CustomersSearchSchema = z.object({
  search: z.string().optional(),
  page: z.number().optional().default(1),
  pageSize: z.number().optional().default(20),
  sortBy: z.string().optional(),
  sortOrder: z.enum(['asc', 'desc']).optional().default('desc'),
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
  staticData: {
    titleKey: 'customers.title',
  },
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
  const { t } = useTranslation()
  const { businessDescriptor } = Route.useParams()
  const navigate = useNavigate()
  const search = Route.useSearch()
  const isMobile = useMediaQuery('(max-width: 768px)')

  const [customers, setCustomers] = useState<Array<Customer>>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isLoadingMore, setIsLoadingMore] = useState(false)
  const [totalItems, setTotalItems] = useState(0)
  const [totalPages, setTotalPages] = useState(0)

  const [isAddCustomerOpen, setIsAddCustomerOpen] = useState(false)
  const [isEditCustomerOpen, setIsEditCustomerOpen] = useState(false)
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(
    null,
  )
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

  const orderBy = useMemo<Array<string>>(() => {
    if (!search.sortBy) return ['-createdAt']
    return [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
  }, [search.sortBy, search.sortOrder])

  const fetchCustomers = useCallback(
    async (append?: boolean) => {
      try {
        if (append) {
          setIsLoadingMore(true)
        } else {
          setIsLoading(true)
        }

        const pageToFetch = append ? search.page + 1 : search.page

        const response = await customerApi.listCustomers(businessDescriptor, {
          search: search.search,
          page: pageToFetch,
          pageSize: search.pageSize,
          orderBy,
        })

        if (append) {
          setCustomers((prev) => [...prev, ...response.items])
          void navigate({
            to: '.',
            search: { ...search, page: pageToFetch },
          })
        } else {
          setCustomers(response.items)
        }

        setTotalItems(response.totalCount)
        setTotalPages(response.totalPages)
      } catch (error) {
        void showErrorFromException(error, t)
      } finally {
        setIsLoading(false)
        setIsLoadingMore(false)
      }
    },
    [
      businessDescriptor,
      navigate,
      orderBy,
      search,
      t,
      setCustomers,
      setTotalItems,
      setTotalPages,
    ],
  )

  useEffect(() => {
    void fetchCustomers(false)
  }, [fetchCustomers])

  const handleSearch = (value: string) => {
    void navigate({
      to: '.',
      search: {
        ...search,
        search: value || undefined,
        page: 1,
      },
    })
  }

  const handlePageChange = (newPage: number) => {
    void navigate({
      to: '.',
      search: { ...search, page: newPage },
    })
  }

  const handlePageSizeChange = (newPageSize: number) => {
    void navigate({
      to: '.',
      search: { ...search, pageSize: newPageSize, page: 1 },
    })
  }

  const handleSort = (key: string) => {
    const newSortOrder =
      search.sortBy === key && search.sortOrder === 'asc' ? 'desc' : 'asc'

    void navigate({
      to: '.',
      search: {
        ...search,
        sortBy: key,
        sortOrder: newSortOrder,
        page: 1,
      },
    })
  }

  const handleLoadMore = () => {
    void fetchCustomers(true)
  }

  const handleCustomerClick = (customer: Customer) => {
    void navigate({
      to: '/business/$businessDescriptor/customers/$customerId',
      params: { businessDescriptor, customerId: customer.id },
    })
  }

  const handleEditClick = (customer: Customer, event?: MouseEvent) => {
    if (event) event.stopPropagation()
    setSelectedCustomer(customer)
    setIsEditCustomerOpen(true)
  }

  const handleDeleteClick = (customer: Customer, event?: MouseEvent) => {
    if (event) event.stopPropagation()
    setSelectedCustomer(customer)
    setIsDeleteDialogOpen(true)
  }

  const handleDelete = async () => {
    if (!selectedCustomer) return

    try {
      setIsDeleting(true)
      await customerApi.deleteCustomer(businessDescriptor, selectedCustomer.id)
      showSuccessToast(t('customers.delete_success'))
      await fetchCustomers(false)
    } catch (error) {
      void showErrorFromException(error, t)
    } finally {
      setIsDeleting(false)
      setIsDeleteDialogOpen(false)
      setSelectedCustomer(null)
    }
  }

  const handleCustomerUpdated = () => {
    void fetchCustomers(false)
  }

  // Get selected business country code for form defaults
  const selectedBusiness = getSelectedBusiness()
  const businessCountryCode = selectedBusiness?.country ?? 'AE'
  const currency = selectedBusiness?.currency ?? 'AED'

  const tableColumns = useMemo<Array<TableColumn<Customer>>>(
    () => [
      {
        key: 'name',
        label: t('customers.name'),
        sortable: true,
        render: (customer) => (
          <div className="flex items-center gap-3">
            <Avatar
              src={customer.avatarUrl}
              alt={customer.name}
              fallback={customer.name
                .split(' ')
                .map((w) => w[0])
                .join('')
                .toUpperCase()
                .slice(0, 2)}
              size="sm"
            />
            <span className="font-medium">{customer.name}</span>
          </div>
        ),
      },
      {
        key: 'phone',
        label: t('customers.phone'),
        render: (customer) => {
          if (customer.phoneCode && customer.phoneNumber) {
            return `${customer.phoneCode} ${customer.phoneNumber}`
          }
          return <span className="text-base-content/40">—</span>
        },
      },
      {
        key: 'ordersCount',
        label: t('customers.orders_count'),
        sortable: true,
        align: 'center',
        render: (customer) => (
          <div className="badge badge-success badge-sm">
            {customer.ordersCount ?? 0}
          </div>
        ),
      },
      {
        key: 'totalSpent',
        label: t('customers.total_spent'),
        sortable: true,
        align: 'end',
        render: (customer) => {
          const spent = customer.totalSpent ?? 0
          return (
            <span className="font-semibold">
              {currency} {spent.toFixed(2)}
            </span>
          )
        },
      },
      {
        key: 'actions',
        label: t('common.actions'),
        align: 'center',
        width: '120px',
        render: (customer) => (
          <div className="flex items-center justify-center gap-2">
            <button
              type="button"
              className="btn btn-ghost btn-sm btn-square"
              onClick={(e) => {
                e.stopPropagation()
                handleCustomerClick(customer)
              }}
              aria-label={t('common.view')}
              title={t('common.view')}
            >
              <Eye size={16} />
            </button>
            <button
              type="button"
              className="btn btn-ghost btn-sm btn-square"
              onClick={(e) => {
                handleEditClick(customer, e)
              }}
              aria-label={t('common.edit')}
              title={t('common.edit')}
            >
              <Edit size={16} />
            </button>
            <button
              type="button"
              className="btn btn-ghost btn-sm btn-square text-error"
              onClick={(e) => {
                handleDeleteClick(customer, e)
              }}
              aria-label={t('common.delete')}
              title={t('common.delete')}
            >
              <Trash2 size={16} />
            </button>
          </div>
        ),
      },
    ],
    [
      currency,
      handleCustomerClick,
      handleDeleteClick,
      handleEditClick,
      t,
    ],
  )

  return (
    <>
      <div className="space-y-4">
        {/* Header */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold">{t('customers.title')}</h1>
            <p className="text-sm text-base-content/60 mt-1">
              {t('customers.subtitle')}
            </p>
          </div>
          <button
            type="button"
            className="btn btn-primary gap-2"
            onClick={() => {
              setIsAddCustomerOpen(true)
            }}
            disabled={!businessDescriptor}
            aria-disabled={!businessDescriptor}
          >
            <Plus size={20} />
            {t('customers.add_customer')}
          </button>
        </div>

        {/* Toolbar */}
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="flex-1">
            <SearchInput
              value={search.search ?? ''}
              onChange={handleSearch}
              placeholder={t('customers.search_placeholder')}
            />
          </div>
          <FilterButton
            title={t('customers.filters')}
            buttonText={t('common.filter')}
            applyLabel={t('common.apply')}
            resetLabel={t('common.reset')}
            onApply={() => {}}
            onReset={() => {}}
          >
            <div className="space-y-4">
              <p className="text-sm text-base-content/60">
                {t('customers.filters_coming_soon')}
              </p>
            </div>
          </FilterButton>
        </div>

        {/* Empty State */}
        {!isLoading && customers.length === 0 && (
          <div className="card bg-base-100 shadow">
            <div className="card-body items-center text-center py-12">
              <Users size={48} className="text-base-content/20 mb-4" />
              <h3 className="text-lg font-semibold mb-2">
                {search.search ? t('customers.no_results') : t('customers.no_customers')}
              </h3>
              <p className="text-sm text-base-content/70 mb-4">
                {search.search
                  ? t('customers.try_different_search')
                  : t('customers.get_started_message')}
              </p>
              {!search.search && (
                <button
                  type="button"
                  className="btn btn-primary btn-sm gap-2"
                  onClick={() => {
                    setIsAddCustomerOpen(true)
                  }}
                >
                  <Plus size={18} />
                  {t('customers.add_first_customer')}
                </button>
              )}
            </div>
          </div>
        )}

        {/* Desktop: Table View */}
        {!isMobile && (
          <>
            <div className="overflow-x-auto">
              <Table
                columns={tableColumns}
                data={customers}
                keyExtractor={(customer) => customer.id}
                isLoading={isLoading}
                emptyMessage={t('customers.no_customers')}
                sortBy={search.sortBy}
                sortOrder={search.sortOrder}
                onSort={handleSort}
              />
            </div>
            <Pagination
              currentPage={search.page}
              totalPages={totalPages}
              pageSize={search.pageSize}
              totalItems={totalItems}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
              itemsName={t('customers.customers').toLowerCase()}
            />
          </>
        )}

        {/* Mobile: Card View with Infinite Scroll */}
        {isMobile && (
          <InfiniteScroll
            hasMore={search.page < totalPages}
            isLoading={isLoadingMore}
            onLoadMore={handleLoadMore}
            loadingMessage={t('common.loading_more')}
            endMessage={t('customers.no_more_customers')}
          >
            <div className="space-y-3">
              {isLoading && customers.length === 0 ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <div key={i} className="skeleton h-32 rounded-box" />
                ))
              ) : customers.length === 0 ? (
                <div className="text-center py-12 text-base-content/60">
                  {t('customers.no_customers')}
                </div>
              ) : (
                customers.map((customer) => (
                  <div key={customer.id} className="relative group">
                    <CustomerCard
                      customer={customer}
                      onClick={handleCustomerClick}
                      ordersCount={customer.ordersCount ?? 0}
                      totalSpent={customer.totalSpent ?? 0}
                      currency={currency}
                    />
                    <div className="absolute top-2 end-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        type="button"
                        className="btn btn-sm btn-circle btn-ghost bg-base-100 shadow-md"
                        onClick={(e) => {
                          handleEditClick(customer, e)
                        }}
                        aria-label={t('common.edit')}
                      >
                        <Edit size={14} />
                      </button>
                      <button
                        type="button"
                        className="btn btn-sm btn-circle btn-ghost bg-base-100 shadow-md text-error"
                        onClick={(e) => {
                          handleDeleteClick(customer, e)
                        }}
                        aria-label={t('common.delete')}
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </InfiniteScroll>
        )}
      </div>

      <AddCustomerSheet
        isOpen={isAddCustomerOpen}
        onClose={() => {
          setIsAddCustomerOpen(false)
        }}
        businessDescriptor={businessDescriptor}
        businessCountryCode={businessCountryCode}
        onCreated={async () => {
          await fetchCustomers(false)
        }}
      />

      {selectedCustomer && (
        <EditCustomerSheet
          isOpen={isEditCustomerOpen}
          onClose={() => {
            setIsEditCustomerOpen(false)
            setSelectedCustomer(null)
          }}
          businessDescriptor={businessDescriptor}
          customer={selectedCustomer}
          onUpdated={handleCustomerUpdated}
        />
      )}

      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => {
          setIsDeleteDialogOpen(false)
          setSelectedCustomer(null)
        }}
        title={t('customers.delete_confirm_title')}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => {
                setIsDeleteDialogOpen(false)
                setSelectedCustomer(null)
              }}
              disabled={isDeleting}
            >
              {t('common.cancel')}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={() => {
                void handleDelete()
              }}
              disabled={isDeleting}
            >
              {isDeleting && (
                <span className="loading loading-spinner loading-sm" />
              )}
              {t('common.delete')}
            </button>
          </div>
        }
      >
        <p>
          {t('customers.delete_confirm_message', {
            name: selectedCustomer?.name,
          })}
        </p>
      </Dialog>
    </>
  )
}
