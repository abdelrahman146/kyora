import { useNavigate, useParams, useSearch } from '@tanstack/react-router'
import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Edit, Eye, Trash2, Users } from 'lucide-react'
import { CustomerCard } from './CustomerCard'
import { CustomerListSkeleton } from './CustomerListSkeleton'
import { AddCustomerSheet } from './AddCustomerSheet'
import { EditCustomerSheet } from './EditCustomerSheet'
import { CountrySelect } from './CountrySelect'
import type { MouseEvent } from 'react'
import type { QueryClient } from '@tanstack/react-query'

import type { Customer, SocialPlatform } from '@/api/customer'
import type { TableColumn } from '@/components/organisms/Table'
import type { SortOption } from '@/components/organisms/SortButton'
import type { CustomersSearch } from '@/features/customers/schema/customersSearch'
import {
  customerQueries,
  useCustomersQuery,
  useDeleteCustomerMutation,
} from '@/api/customer'
import { Avatar } from '@/components/atoms/Avatar'
import { Dialog } from '@/components/atoms/Dialog'
import { ResourceListLayout } from '@/components/templates/ResourceListLayout'
import { showErrorFromException, showSuccessToast } from '@/lib/toast'
import { getSelectedBusiness } from '@/stores/businessStore'
import { useKyoraForm } from '@/lib/form/useKyoraForm'
import { getMetadata } from '@/stores/metadataStore'
import { formatDateShort } from '@/lib/formatDate'
import { CustomersSearchSchema } from '@/features/customers/schema/customersSearch'

export async function customersListLoader({
  queryClient,
  businessDescriptor,
  search,
}: {
  queryClient: QueryClient
  businessDescriptor: string
  search: CustomersSearch
}) {
  const orderBy = search.sortBy
    ? [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
    : ['-joinedAt']

  await queryClient.prefetchQuery(
    customerQueries.list(businessDescriptor, {
      search: search.search,
      page: search.page,
      pageSize: search.pageSize,
      orderBy,
      countryCode: search.countryCode,
      hasOrders: search.hasOrders,
      socialPlatforms: search.socialPlatforms,
    }),
  )
}

export function CustomersListPage() {
  const { t } = useTranslation()

  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/customers/',
  })

  const rawSearch = useSearch({
    from: '/business/$businessDescriptor/customers/',
  })
  const search = CustomersSearchSchema.parse(rawSearch)

  const navigate = useNavigate({
    from: '/business/$businessDescriptor/customers',
  })

  const [isAddCustomerOpen, setIsAddCustomerOpen] = useState(false)
  const [isEditCustomerOpen, setIsEditCustomerOpen] = useState(false)
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(
    null,
  )
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)

  const orderBy = useMemo<Array<string>>(() => {
    if (!search.sortBy) return ['-joinedAt']
    return [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
  }, [search.sortBy, search.sortOrder])

  const filterForm = useKyoraForm({
    defaultValues: {
      countryCode: search.countryCode ?? '',
      hasOrders: search.hasOrders ?? false,
      socialPlatforms: (search.socialPlatforms ?? []) as Array<SocialPlatform>,
    },
    onSubmit: async ({ value }) => {
      await navigate({
        to: '.',
        search: (prev) => ({
          ...prev,
          countryCode: value.countryCode || undefined,
          hasOrders: value.hasOrders ? true : undefined,
          socialPlatforms:
            value.socialPlatforms.length > 0
              ? value.socialPlatforms
              : undefined,
          page: 1,
        }),
      })
    },
  })

  useEffect(() => {
    filterForm.setFieldValue('countryCode', search.countryCode ?? '')
    filterForm.setFieldValue('hasOrders', search.hasOrders ?? false)
    filterForm.setFieldValue(
      'socialPlatforms',
      (search.socialPlatforms ?? []) as Array<SocialPlatform>,
    )
  }, [search.countryCode, search.hasOrders, search.socialPlatforms])

  const { data: customersResponse, isLoading } = useCustomersQuery(
    businessDescriptor,
    {
      search: search.search,
      page: search.page,
      pageSize: search.pageSize,
      orderBy,
      countryCode: search.countryCode,
      hasOrders: search.hasOrders,
      socialPlatforms: search.socialPlatforms,
    },
  )

  const customers = customersResponse?.items ?? []
  const totalItems = customersResponse?.totalCount ?? 0
  const totalPages = customersResponse?.totalPages ?? 0

  const activeFilterCount =
    (search.countryCode ? 1 : 0) +
    (search.hasOrders ? 1 : 0) +
    (search.socialPlatforms && search.socialPlatforms.length > 0 ? 1 : 0)

  const deleteMutation = useDeleteCustomerMutation(businessDescriptor, {
    onSuccess: () => {
      showSuccessToast(t('customers.delete_success'))
      setIsDeleteDialogOpen(false)
      setSelectedCustomer(null)
    },
    onError: (error) => {
      void showErrorFromException(error, t)
    },
  })

  const handleSearch = (value: string) => {
    void navigate({
      to: '.',
      search: (prev) => ({
        ...prev,
        search: value || undefined,
        page: 1,
      }),
    })
  }

  const handlePageChange = (newPage: number) => {
    void navigate({
      to: '.',
      search: (prev) => ({ ...prev, page: newPage }),
    })
  }

  const handlePageSizeChange = (newPageSize: number) => {
    void navigate({
      to: '.',
      search: (prev) => ({ ...prev, pageSize: newPageSize, page: 1 }),
    })
  }

  const handleResetFilters = async () => {
    filterForm.reset()
    await navigate({
      to: '.',
      search: (prev) => ({
        ...prev,
        countryCode: undefined,
        hasOrders: undefined,
        socialPlatforms: undefined,
        page: 1,
      }),
    })
  }

  const handleSort = (key: string) => {
    const newSortOrder =
      search.sortBy === key && search.sortOrder === 'asc' ? 'desc' : 'asc'

    void navigate({
      to: '.',
      search: (prev) => ({
        ...prev,
        sortBy: key,
        sortOrder: newSortOrder,
        page: 1,
      }),
    })
  }

  const handleSortApply = (sortBy: string, sortOrder: 'asc' | 'desc') => {
    void navigate({
      to: '.',
      search: (prev) => ({
        ...prev,
        sortBy,
        sortOrder,
        page: 1,
      }),
    })
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

  const handleDelete = () => {
    if (!selectedCustomer) return
    deleteMutation.mutate(selectedCustomer.id)
  }

  const handleCustomerCreated = () => {
    setIsAddCustomerOpen(false)
  }

  const handleCustomerUpdated = () => {
    setIsEditCustomerOpen(false)
  }

  const selectedBusiness = getSelectedBusiness()
  const businessCountryCode = selectedBusiness?.countryCode ?? 'AE'
  const currency = selectedBusiness?.currency ?? 'AED'

  const sortOptions = useMemo<Array<SortOption>>(
    () => [
      { value: 'name', label: t('customers.name') },
      { value: 'countryCode', label: t('customers.country') },
      { value: 'ordersCount', label: t('customers.orders_count') },
      { value: 'totalSpent', label: t('customers.total_spent') },
      { value: 'joinedAt', label: t('customers.joined_date') },
    ],
    [t],
  )

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
        key: 'countryCode',
        label: t('customers.country'),
        sortable: true,
        render: (customer) => {
          const metadata = getMetadata()
          const country = metadata.countries.find(
            (c) =>
              c.code === customer.countryCode ||
              c.iso_code === customer.countryCode,
          )
          if (!country) {
            return <span className="text-base-content/40">—</span>
          }
          return (
            <div className="flex items-center gap-2">
              {country.flag && <span className="text-lg">{country.flag}</span>}
              <span className="text-sm">{country.name}</span>
            </div>
          )
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
        key: 'joinedAt',
        label: t('customers.joined_date'),
        sortable: true,
        render: (customer) => (
          <span className="text-sm text-base-content/70">
            {formatDateShort(customer.joinedAt)}
          </span>
        ),
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
    [currency, handleCustomerClick, handleDeleteClick, handleEditClick, t],
  )

  return (
    <>
      <ResourceListLayout
        title={t('customers.title')}
        subtitle={t('customers.subtitle')}
        addButtonText={t('customers.add_customer')}
        onAddClick={() => {
          setIsAddCustomerOpen(true)
        }}
        addButtonDisabled={!businessDescriptor}
        searchPlaceholder={t('customers.search_placeholder')}
        searchValue={search.search ?? ''}
        onSearchChange={handleSearch}
        filterTitle={t('customers.filters')}
        filterButtonText={t('common.filter')}
        filterButton={
          <filterForm.AppForm>
            <div className="space-y-6 p-4">
              <filterForm.AppField name="countryCode">
                {(field) => (
                  <div className="form-control">
                    <CountrySelect
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      placeholder={t('customers.all_countries')}
                      searchable
                    />
                  </div>
                )}
              </filterForm.AppField>

              <filterForm.AppField name="hasOrders">
                {(field) => (
                  <field.ToggleField
                    label={t('customers.filter_only_with_orders')}
                    description={t('customers.filter_only_with_orders_desc')}
                  />
                )}
              </filterForm.AppField>

              <filterForm.AppField name="socialPlatforms">
                {(field) => (
                  <field.CheckboxGroupField
                    label={t('customers.filter_by_social_platform')}
                    options={[
                      {
                        value: 'instagram' as const,
                        label: t('customers.instagram'),
                      },
                      {
                        value: 'tiktok' as const,
                        label: t('customers.tiktok'),
                      },
                      {
                        value: 'facebook' as const,
                        label: t('customers.facebook'),
                      },
                      {
                        value: 'x' as const,
                        label: t('customers.x'),
                      },
                      {
                        value: 'snapchat' as const,
                        label: t('customers.snapchat'),
                      },
                      {
                        value: 'whatsapp' as const,
                        label: t('customers.whatsapp'),
                      },
                    ]}
                  />
                )}
              </filterForm.AppField>
            </div>
          </filterForm.AppForm>
        }
        activeFilterCount={activeFilterCount}
        applyLabel={t('common.apply')}
        resetLabel={t('common.reset')}
        onApplyFilters={() => {
          filterForm.handleSubmit()
        }}
        onResetFilters={handleResetFilters}
        sortTitle={t('customers.sort_customers')}
        sortOptions={sortOptions}
        onSortApply={handleSortApply}
        emptyIcon={<Users size={48} />}
        emptyTitle={
          search.search
            ? t('customers.no_results')
            : t('customers.no_customers')
        }
        emptyMessage={
          search.search
            ? t('customers.try_different_search')
            : t('customers.get_started_message')
        }
        emptyActionText={
          !search.search ? t('customers.add_first_customer') : undefined
        }
        onEmptyAction={
          !search.search
            ? () => {
                setIsAddCustomerOpen(true)
              }
            : undefined
        }
        noResultsTitle={t('customers.no_results')}
        noResultsMessage={t('customers.try_different_search')}
        tableColumns={tableColumns}
        tableData={customers}
        tableKeyExtractor={(customer) => customer.id}
        tableSortBy={search.sortBy}
        tableSortOrder={search.sortOrder}
        onTableSort={handleSort}
        onTableRowClick={handleCustomerClick}
        mobileCard={(customer) => (
          <div className="relative group">
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
                className="btn btn-sm btn-circle btn-ghost bg-base-100"
                onClick={(e) => {
                  handleEditClick(customer, e)
                }}
                aria-label={t('common.edit')}
              >
                <Edit size={16} />
              </button>
              <button
                type="button"
                className="btn btn-sm btn-circle btn-ghost bg-base-100 text-error"
                onClick={(e) => {
                  handleDeleteClick(customer, e)
                }}
                aria-label={t('common.delete')}
              >
                <Trash2 size={16} />
              </button>
            </div>
          </div>
        )}
        isLoading={isLoading}
        hasSearchQuery={!!search.search}
        currentPage={search.page}
        totalPages={totalPages}
        pageSize={search.pageSize}
        totalItems={totalItems}
        onPageChange={handlePageChange}
        onPageSizeChange={handlePageSizeChange}
        itemsName={t('customers.customers').toLowerCase()}
        skeleton={<CustomerListSkeleton />}
      />

      <AddCustomerSheet
        isOpen={isAddCustomerOpen}
        onClose={() => {
          setIsAddCustomerOpen(false)
        }}
        businessDescriptor={businessDescriptor}
        businessCountryCode={businessCountryCode}
        onCreated={handleCustomerCreated}
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
              disabled={deleteMutation.isPending}
            >
              {t('common.cancel')}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? (
                <>
                  <span className="loading loading-spinner loading-sm"></span>
                  {t('common.deleting')}
                </>
              ) : (
                t('common.delete')
              )}
            </button>
          </div>
        }
      >
        <p className="text-base-content/70">
          {t('customers.delete_confirm_message', {
            name: selectedCustomer?.name,
          })}
        </p>
      </Dialog>
    </>
  )
}
