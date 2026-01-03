import { createFileRoute } from '@tanstack/react-router'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { useTranslation } from 'react-i18next'
import { Edit, Eye, ShoppingCart, Trash2 } from 'lucide-react'
import { z } from 'zod'

import type { Order } from '@/api/order'
import type { SortOption } from '@/components/organisms/SortButton'
import type { TableColumn } from '@/components/organisms/Table'

import { OrderCard } from '@/components'
import { ResourceListLayout } from '@/components/templates/ResourceListLayout'
import { useKyoraForm } from '@/lib/form'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'
import { orderApi, orderQueries, useOrdersQuery } from '@/api/order'

const OrdersSearchSchema = z.object({
  search: z.string().optional(),
  page: z.number().optional(),
  pageSize: z.number().optional(),
  sortBy: z
    .enum(['orderNumber', 'total', 'status', 'paymentStatus', 'orderedAt'])
    .optional(),
  sortOrder: z.enum(['asc', 'desc']).optional(),
  status: z.array(z.string()).optional(),
  paymentStatus: z.array(z.string()).optional(),
  customerId: z.string().optional(),
  from: z.string().optional(),
  to: z.string().optional(),
})

type OrdersSearch = z.infer<typeof OrdersSearchSchema>

export const Route = createFileRoute('/business/$businessDescriptor/orders/')({
  validateSearch: (search) => OrdersSearchSchema.parse(search),
  loaderDeps: ({ search }) => search,
  loader: async ({ context, deps: search, params }) => {
    // Access queryClient from context
    const { queryClient } = context as any

    const orderByArray = search.sortBy
      ? [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
      : ['-orderedAt']

    await queryClient.ensureQueryData(
      orderQueries.list(params.businessDescriptor, {
        search: search.search,
        page: search.page,
        pageSize: search.pageSize,
        orderBy: orderByArray,
        status: search.status as Array<Order['status']>,
        paymentStatus: search.paymentStatus as Array<Order['paymentStatus']>,
        customerId: search.customerId,
        from: search.from,
        to: search.to,
      }),
    )
  },
  component: OrdersListPage,
})

function OrdersListPage() {
  const { t } = useTranslation()
  const { businessDescriptor } = Route.useParams()
  const search = Route.useSearch()
  const navigate = Route.useNavigate()
  const queryClient = useQueryClient()

  const orderByArray = search.sortBy
    ? [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
    : ['-orderedAt']

  const { data, isLoading } = useOrdersQuery(businessDescriptor, {
    search: search.search,
    page: search.page,
    pageSize: search.pageSize,
    orderBy: orderByArray,
    status: search.status as Array<Order['status']>,
    paymentStatus: search.paymentStatus as Array<Order['paymentStatus']>,
    customerId: search.customerId,
    from: search.from,
    to: search.to,
  })

  const deleteMutation = useMutation({
    mutationFn: (params: { businessDescriptor: string; orderId: string }) =>
      orderApi.deleteOrder(params.businessDescriptor, params.orderId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: orderQueries.all })
      toast.success(t('orders:delete_success'))
    },
    onError: () => {
      toast.error(t('common:error_occurred'))
    },
  })

  const form = useKyoraForm({
    defaultValues: {
      status: search.status || [],
      paymentStatus: search.paymentStatus || [],
    },
    onSubmit: async ({ value }) => {
      await navigate({
        search: (prev) => ({
          ...prev,
          status: value.status.length > 0 ? value.status : undefined,
          paymentStatus:
            value.paymentStatus.length > 0 ? value.paymentStatus : undefined,
          page: 1,
        }),
      })
    },
  })

  const handleSearchChange = (value: string) => {
    navigate({
      search: (prev) => ({
        ...prev,
        search: value || undefined,
        page: 1,
      }),
    })
  }

  const handleSortChange = (sortBy: OrdersSearch['sortBy']) => {
    navigate({
      search: (prev) => ({
        ...prev,
        sortBy,
        sortOrder:
          prev.sortBy === sortBy && prev.sortOrder === 'asc' ? 'desc' : 'asc',
      }),
    })
  }

  const handleAddOrder = () => {
    console.log('TODO: Open add order sheet')
  }

  const handleViewClick = (order: Order) => {
    console.log('TODO: Navigate to order details', order.id)
  }

  const handleEditClick = (order: Order) => {
    console.log('TODO: Open edit order sheet', order.id)
  }

  const handleDelete = async (order: Order) => {
    if (
      window.confirm(
        t('orders:delete_confirm_message', { orderNumber: order.orderNumber }),
      )
    ) {
      await deleteMutation.mutateAsync({
        businessDescriptor,
        orderId: order.id,
      })
    }
  }

  const columns: Array<TableColumn<Order>> = [
    {
      key: 'orderNumber',
      label: t('orders:order_number'),
      sortable: true,
      render: (order: Order) => (
        <span className="font-medium">{order.orderNumber}</span>
      ),
    },
    {
      key: 'customer',
      label: t('orders:customer'),
      render: (order: Order) => (
        <div className="flex items-center gap-3">
          <div className="avatar placeholder">
            <div className="w-10 rounded-full bg-primary/10 text-primary">
              <span className="text-sm">
                {order.customer?.name.charAt(0).toUpperCase() || '?'}
              </span>
            </div>
          </div>
          <div>
            <div className="font-medium">
              {order.customer?.name || t('common:unknown')}
            </div>
            {order.customer?.phoneNumber && (
              <div className="text-sm text-base-content/70" dir="ltr">
                {order.customer.phoneCode} {order.customer.phoneNumber}
              </div>
            )}
          </div>
        </div>
      ),
    },
    {
      key: 'status',
      label: t('orders:status'),
      sortable: true,
      render: (order: Order) => {
        const statusMap: Record<
          Order['status'],
          { class: string; label: string }
        > = {
          pending: {
            class: 'badge-warning',
            label: t('orders:status_pending'),
          },
          placed: { class: 'badge-info', label: t('orders:status_placed') },
          ready_for_shipment: {
            class: 'badge-info',
            label: t('orders:status_ready_for_shipment'),
          },
          shipped: {
            class: 'badge-primary',
            label: t('orders:status_shipped'),
          },
          fulfilled: {
            class: 'badge-success',
            label: t('orders:status_fulfilled'),
          },
          cancelled: {
            class: 'badge-error',
            label: t('orders:status_cancelled'),
          },
          returned: {
            class: 'badge-error',
            label: t('orders:status_returned'),
          },
        }

        const config = statusMap[order.status]
        return <span className={`badge ${config.class}`}>{config.label}</span>
      },
    },
    {
      key: 'paymentStatus',
      label: t('orders:payment_status'),
      sortable: true,
      render: (order: Order) => {
        const statusMap: Record<
          Order['paymentStatus'],
          { class: string; label: string }
        > = {
          pending: {
            class: 'badge-warning',
            label: t('orders:payment_status_pending'),
          },
          paid: {
            class: 'badge-success',
            label: t('orders:payment_status_paid'),
          },
          failed: {
            class: 'badge-error',
            label: t('orders:payment_status_failed'),
          },
          refunded: {
            class: 'badge-ghost',
            label: t('orders:payment_status_refunded'),
          },
        }

        const config = statusMap[order.paymentStatus]
        return <span className={`badge ${config.class}`}>{config.label}</span>
      },
    },
    {
      key: 'total',
      label: t('orders:total'),
      sortable: true,
      render: (order: Order) => (
        <span className="font-medium">
          {formatCurrency(parseFloat(order.total), order.currency || 'USD')}
        </span>
      ),
    },
    {
      key: 'orderedAt',
      label: t('orders:ordered_date'),
      sortable: true,
      render: (order: Order) => (
        <span className="text-sm">{formatDateShort(order.orderedAt)}</span>
      ),
    },
    {
      key: 'actions',
      label: t('common.actions'),
      align: 'center',
      width: '120px',
      render: (order: Order) => (
        <div className="flex gap-2 justify-end">
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-square"
            onClick={() => handleViewClick(order)}
            aria-label={t('common:view')}
          >
            <Eye size={18} />
          </button>
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-square"
            onClick={() => handleEditClick(order)}
            aria-label={t('common:edit')}
          >
            <Edit size={18} />
          </button>
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-square text-error hover:bg-error/10"
            onClick={() => handleDelete(order)}
            aria-label={t('common:delete')}
          >
            <Trash2 size={18} />
          </button>
        </div>
      ),
    },
  ]

  const sortOptions: Array<SortOption> = [
    { value: 'orderNumber', label: t('orders:order_number') },
    { value: 'total', label: t('orders:total') },
    { value: 'status', label: t('orders:status') },
    { value: 'paymentStatus', label: t('orders:payment_status') },
    { value: 'orderedAt', label: t('orders:ordered_date') },
  ]

  const handleSortApply = (sortBy: string) => {
    navigate({
      search: (prev) => ({
        ...prev,
        sortBy: sortBy as OrdersSearch['sortBy'],
        sortOrder:
          prev.sortBy === sortBy && prev.sortOrder === 'asc' ? 'desc' : 'asc',
      }),
    })
  }

  const handleApplyFilters = () => {
    form.handleSubmit()
  }

  const handleResetFilters = () => {
    form.reset()
    navigate({
      search: (prev) => ({
        ...prev,
        status: undefined,
        paymentStatus: undefined,
        page: 1,
      }),
    })
  }

  const activeFilterCount =
    (search.status?.length || 0) + (search.paymentStatus?.length || 0)

  const filterContent = (
    <form.AppForm>
      <form.FormRoot className="space-y-4">
        <div>
          <h3 className="font-medium mb-3">{t('orders:filter_by_status')}</h3>
          <form.AppField name="status">
            {(field) => (
              <div className="space-y-2">
                {[
                  { value: 'pending', label: t('orders:status_pending') },
                  { value: 'placed', label: t('orders:status_placed') },
                  {
                    value: 'ready_for_shipment',
                    label: t('orders:status_ready_for_shipment'),
                  },
                  { value: 'shipped', label: t('orders:status_shipped') },
                  { value: 'fulfilled', label: t('orders:status_fulfilled') },
                  { value: 'cancelled', label: t('orders:status_cancelled') },
                  { value: 'returned', label: t('orders:status_returned') },
                ].map((option) => (
                  <label
                    key={option.value}
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      className="checkbox checkbox-sm"
                      checked={field.state.value.includes(option.value)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          field.handleChange([
                            ...field.state.value,
                            option.value,
                          ])
                        } else {
                          field.handleChange(
                            field.state.value.filter((v) => v !== option.value),
                          )
                        }
                      }}
                    />
                    <span className="text-sm">{option.label}</span>
                  </label>
                ))}
              </div>
            )}
          </form.AppField>
        </div>

        <div>
          <h3 className="font-medium mb-3">
            {t('orders:filter_by_payment_status')}
          </h3>
          <form.AppField name="paymentStatus">
            {(field) => (
              <div className="space-y-2">
                {[
                  {
                    value: 'pending',
                    label: t('orders:payment_status_pending'),
                  },
                  { value: 'paid', label: t('orders:payment_status_paid') },
                  { value: 'failed', label: t('orders:payment_status_failed') },
                  {
                    value: 'refunded',
                    label: t('orders:payment_status_refunded'),
                  },
                ].map((option) => (
                  <label
                    key={option.value}
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      className="checkbox checkbox-sm"
                      checked={field.state.value.includes(option.value)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          field.handleChange([
                            ...field.state.value,
                            option.value,
                          ])
                        } else {
                          field.handleChange(
                            field.state.value.filter((v) => v !== option.value),
                          )
                        }
                      }}
                    />
                    <span className="text-sm">{option.label}</span>
                  </label>
                ))}
              </div>
            )}
          </form.AppField>
        </div>
      </form.FormRoot>
    </form.AppForm>
  )

  return (
    <ResourceListLayout
      title={t('orders:title')}
      subtitle={t('orders:subtitle')}
      addButtonText={t('orders:add_order')}
      onAddClick={handleAddOrder}
      searchPlaceholder={t('orders:search_placeholder')}
      searchValue={search.search ?? ''}
      onSearchChange={handleSearchChange}
      filterTitle={t('orders:filters')}
      filterButtonText={t('common:filter')}
      filterButton={filterContent}
      activeFilterCount={activeFilterCount}
      onApplyFilters={handleApplyFilters}
      onResetFilters={handleResetFilters}
      applyLabel={t('common:apply')}
      resetLabel={t('common:reset')}
      sortTitle={t('orders:sort_orders')}
      sortOptions={sortOptions}
      onSortApply={handleSortApply}
      emptyIcon={<ShoppingCart size={48} />}
      emptyTitle={search.search ? t('no_results') : t('no_orders')}
      emptyMessage={
        search.search ? t('try_different_search') : t('get_started_message')
      }
      emptyActionText={!search.search ? t('add_order') : undefined}
      onEmptyAction={!search.search ? handleAddOrder : undefined}
      tableColumns={columns}
      tableData={data?.items || []}
      tableKeyExtractor={(order) => order.id}
      tableSortBy={search.sortBy}
      tableSortOrder={search.sortOrder}
      onTableSort={(column: string) =>
        handleSortChange(column as OrdersSearch['sortBy'])
      }
      onTableRowClick={handleViewClick}
      mobileCard={(order: Order) => (
        <div className="relative group">
          <OrderCard order={order} onClick={() => handleViewClick(order)} />
          <div className="absolute top-2 end-2 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              type="button"
              className="btn btn-sm btn-circle btn-ghost bg-base-100"
              onClick={(e) => {
                e.stopPropagation()
                handleEditClick(order)
              }}
              aria-label={t('common:edit')}
            >
              <Edit size={16} />
            </button>
            <button
              type="button"
              className="btn btn-sm btn-circle btn-ghost bg-base-100 text-error"
              onClick={(e) => {
                e.stopPropagation()
                handleDelete(order)
              }}
              aria-label={t('common:delete')}
            >
              <Trash2 size={16} />
            </button>
          </div>
        </div>
      )}
      isLoading={isLoading}
      hasSearchQuery={!!search.search}
      currentPage={search.page || 1}
      totalPages={data?.totalPages || 1}
      pageSize={search.pageSize || 20}
      totalItems={data?.totalCount || 0}
      onPageChange={(page: number) => {
        navigate({
          search: (prev) => ({
            ...prev,
            page,
          }),
        })
      }}
      itemsName={t('orders:orders')}
    />
  )
}
