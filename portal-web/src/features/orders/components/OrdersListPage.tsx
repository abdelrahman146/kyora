import { Link, useNavigate, useParams, useSearch } from '@tanstack/react-router'
import { useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { Calendar, Clock, Eye, ShoppingCart, User } from 'lucide-react'
import {
  FaFacebook,
  FaInstagram,
  FaSnapchat,
  FaTiktok,
  FaWhatsapp,
} from 'react-icons/fa'
import { FaXTwitter } from 'react-icons/fa6'
import { useEffect, useState } from 'react'
import { formatISO } from 'date-fns'
import { OrdersSearchSchema } from '../schema/ordersSearch'
import { OrderCard } from './OrderCard'
import { OrderQuickActions } from './OrderQuickActions'
import { OrderReviewSheet } from './OrderReviewSheet'
import { CreateOrderSheet } from './CreateOrderSheet'
import type { OrdersSearch } from '../schema/ordersSearch'
import type { DateRange } from 'react-day-picker'

import type { Order } from '@/api/order'
import type { SocialPlatform } from '@/api/customer'
import type { SortOption } from '@/components/molecules/SortButton'
import type { TableColumn } from '@/components/organisms/Table'
import { orderQueries, useOrdersQuery } from '@/api/order'
import { ResourceListLayout } from '@/components/templates/ResourceListLayout'
import { DateRangePicker } from '@/components/form/DateRangePicker'
import { useKyoraForm } from '@/lib/form'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export function OrdersListPage() {
  const { t: tOrders } = useTranslation('orders')
  const { t: tCommon } = useTranslation('common')

  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/orders/',
  })

  const rawSearch = useSearch({ from: '/business/$businessDescriptor/orders/' })
  const search = OrdersSearchSchema.parse(rawSearch)

  const navigate = useNavigate({ from: '/business/$businessDescriptor/orders' })

  const queryClient = useQueryClient()
  const [reviewOrder, setReviewOrder] = useState<Order | null>(null)
  const [isCreateOrderOpen, setIsCreateOrderOpen] = useState(false)

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
    socialPlatforms: search.socialPlatforms as Array<SocialPlatform>,
    customerId: search.customerId,
    from: search.from,
    to: search.to,
  })

  const form = useKyoraForm({
    defaultValues: {
      status: search.status || [],
      paymentStatus: search.paymentStatus || [],
      socialPlatforms: (search.socialPlatforms || []) as Array<SocialPlatform>,
      customerId: search.customerId || '',
      dateRange:
        search.from && search.to
          ? ({
              from: new Date(search.from),
              to: new Date(search.to),
            } as DateRange)
          : undefined,
    },
    onSubmit: async ({ value }) => {
      await navigate({
        search: (prev) => ({
          ...prev,
          status: value.status.length > 0 ? value.status : undefined,
          paymentStatus:
            value.paymentStatus.length > 0 ? value.paymentStatus : undefined,
          socialPlatforms:
            value.socialPlatforms.length > 0
              ? value.socialPlatforms
              : undefined,
          customerId: value.customerId || undefined,
          from: value.dateRange?.from
            ? formatISO(value.dateRange.from)
            : undefined,
          to: value.dateRange?.to ? formatISO(value.dateRange.to) : undefined,
          page: 1,
        }),
      })
    },
  })

  useEffect(() => {
    form.setFieldValue('status', search.status || [])
    form.setFieldValue('paymentStatus', search.paymentStatus || [])
    form.setFieldValue(
      'socialPlatforms',
      (search.socialPlatforms || []) as Array<SocialPlatform>,
    )
    form.setFieldValue('customerId', search.customerId || '')
    form.setFieldValue(
      'dateRange',
      search.from && search.to
        ? ({
            from: new Date(search.from),
            to: new Date(search.to),
          } as DateRange)
        : undefined,
    )
  }, [
    search.status,
    search.paymentStatus,
    search.socialPlatforms,
    search.customerId,
    search.from,
    search.to,
  ])

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
    setIsCreateOrderOpen(true)
  }

  const handleViewClick = (order: Order) => {
    setReviewOrder(order)
  }

  const columns: Array<TableColumn<Order>> = [
    {
      key: 'orderNumber',
      label: tOrders('order_number'),
      sortable: true,
      render: (order: Order) => (
        <span className="font-medium">{order.orderNumber}</span>
      ),
    },
    {
      key: 'customer',
      label: tOrders('customer'),
      render: (order: Order) => {
        const getPlatformIcon = () => {
          if (order.customer?.instagramUsername)
            return (
              <FaInstagram
                size={14}
                className="text-pink-600"
                aria-label="Instagram"
              />
            )
          if (order.customer?.tiktokUsername)
            return (
              <FaTiktok size={14} className="text-black" aria-label="TikTok" />
            )
          if (order.customer?.facebookUsername)
            return (
              <FaFacebook
                size={14}
                className="text-blue-600"
                aria-label="Facebook"
              />
            )
          if (order.customer?.xUsername)
            return (
              <FaXTwitter size={14} className="text-black" aria-label="X" />
            )
          if (order.customer?.snapchatUsername)
            return (
              <FaSnapchat
                size={14}
                className="text-yellow-400"
                aria-label="Snapchat"
              />
            )
          if (order.customer?.whatsappNumber)
            return (
              <FaWhatsapp
                size={14}
                className="text-green-600"
                aria-label="WhatsApp"
              />
            )
          return null
        }

        const getPlatformHandle = () => {
          if (order.customer?.instagramUsername)
            return `@${order.customer.instagramUsername}`
          if (order.customer?.tiktokUsername)
            return `@${order.customer.tiktokUsername}`
          if (order.customer?.facebookUsername)
            return `@${order.customer.facebookUsername}`
          if (order.customer?.xUsername) return `@${order.customer.xUsername}`
          if (order.customer?.snapchatUsername)
            return `@${order.customer.snapchatUsername}`
          return null
        }

        const platformIcon = getPlatformIcon()
        const platformHandle = getPlatformHandle()

        return (
          <Link
            to="/business/$businessDescriptor/customers/$customerId"
            params={{
              businessDescriptor,
              customerId: order.customer?.id || '',
            }}
            className="flex cursor-pointer items-center gap-3 hover:bg-base-200 rounded-lg p-2 -m-2 transition-colors"
          >
            <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
              <User size={14} className="text-primary" />
            </div>
            <div className="min-w-0">
              <div className="font-medium truncate">
                {order.customer?.name || tCommon('unknown')}
              </div>
              {platformHandle && (
                <div className="flex items-center gap-1.5 text-xs text-base-content/60">
                  {platformIcon}
                  <span className="truncate">{platformHandle}</span>
                </div>
              )}
            </div>
          </Link>
        )
      },
    },
    {
      key: 'status',
      label: tOrders('status'),
      sortable: true,
      render: (order: Order) => {
        const statusMap: Record<
          Order['status'],
          { class: string; label: string }
        > = {
          pending: { class: 'badge-warning', label: tOrders('status_pending') },
          placed: { class: 'badge-info', label: tOrders('status_placed') },
          ready_for_shipment: {
            class: 'badge-info',
            label: tOrders('status_ready_for_shipment'),
          },
          shipped: { class: 'badge-primary', label: tOrders('status_shipped') },
          fulfilled: {
            class: 'badge-success',
            label: tOrders('status_fulfilled'),
          },
          cancelled: {
            class: 'badge-error',
            label: tOrders('status_cancelled'),
          },
          returned: { class: 'badge-error', label: tOrders('status_returned') },
        }

        const paymentStatusMap: Record<
          Order['paymentStatus'],
          { class: string; label: string }
        > = {
          pending: {
            class: 'badge-warning',
            label: tOrders('payment_status_pending'),
          },
          paid: {
            class: 'badge-success',
            label: tOrders('payment_status_paid'),
          },
          failed: {
            class: 'badge-error',
            label: tOrders('payment_status_failed'),
          },
          refunded: {
            class: 'badge-ghost',
            label: tOrders('payment_status_refunded'),
          },
        }

        const statusConfig = statusMap[order.status]
        const paymentConfig = paymentStatusMap[order.paymentStatus]

        return (
          <div className="flex flex-col gap-1.5">
            <span className={`badge badge-sm ${statusConfig.class}`}>
              {statusConfig.label}
            </span>
            <span className={`badge badge-sm ${paymentConfig.class}`}>
              {paymentConfig.label}
            </span>
            {/* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */}
            {order.paymentMethod && order.paymentStatus === 'paid' && (
              <span className="badge badge-sm badge-ghost">
                {tOrders(`payment_method_${order.paymentMethod}`)}
              </span>
            )}
          </div>
        )
      },
    },
    {
      key: 'total',
      label: tOrders('total'),
      sortable: true,
      render: (order: Order) => (
        <span className="font-medium">
          {formatCurrency(parseFloat(order.total), order.currency || 'USD')}
        </span>
      ),
    },
    {
      key: 'orderedAt',
      label: tOrders('ordered_date'),
      sortable: true,
      render: (order: Order) => {
        const latestTimestamp =
          order.fulfilledAt ||
          order.shippedAt ||
          order.readyForShipmentAt ||
          order.placedAt ||
          order.paidAt ||
          order.failedAt ||
          order.refundedAt ||
          order.cancelledAt ||
          order.returnedAt ||
          order.updatedAt

        return (
          <div className="flex flex-col gap-1 text-xs">
            <div className="flex items-center gap-1.5 text-base-content/70">
              <Calendar size={12} />
              <span>{formatDateShort(order.orderedAt)}</span>
            </div>
            {latestTimestamp !== order.orderedAt && (
              <div className="flex items-center gap-1.5 text-base-content/60">
                <Clock size={12} />
                <span>{formatDateShort(latestTimestamp)}</span>
              </div>
            )}
          </div>
        )
      },
    },
    {
      key: 'actions',
      label: tCommon('actionsLabel'),
      align: 'center',
      width: '200px',
      render: (order: Order) => (
        <div className="flex items-center gap-2 justify-end">
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-square"
            onClick={() => handleViewClick(order)}
            aria-label={tOrders('quick_review')}
          >
            <Eye size={16} />
          </button>
          <OrderQuickActions
            order={order}
            businessDescriptor={businessDescriptor}
            onDeleteSuccess={() => {
              void queryClient.invalidateQueries({ queryKey: orderQueries.all })
            }}
          />
        </div>
      ),
    },
  ]

  const sortOptions: Array<SortOption> = [
    { value: 'orderNumber', label: tOrders('order_number') },
    { value: 'total', label: tOrders('total') },
    { value: 'status', label: tOrders('status') },
    { value: 'paymentStatus', label: tOrders('payment_status') },
    { value: 'orderedAt', label: tOrders('ordered_date') },
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
        socialPlatforms: undefined,
        customerId: undefined,
        from: undefined,
        to: undefined,
        page: 1,
      }),
    })
  }

  const activeFilterCount =
    (search.status?.length || 0) +
    (search.paymentStatus?.length || 0) +
    (search.socialPlatforms?.length || 0) +
    (search.customerId ? 1 : 0) +
    (search.from && search.to ? 1 : 0)

  const filterContent = (
    <form.AppForm>
      <form.FormRoot className="space-y-6 p-4">
        <form.AppField name="customerId">
          {(field) => (
            <field.CustomerSelectField
              label={tOrders('filter_by_customer')}
              businessDescriptor={businessDescriptor}
              placeholder={tOrders('search_customer_placeholder')}
            />
          )}
        </form.AppField>

        <form.AppField name="dateRange">
          {(field) => (
            <div className="form-control">
              <label className="label pb-2">
                <span className="label-text font-medium">
                  {tOrders('filter_by_date_range')}
                </span>
              </label>
              <DateRangePicker
                value={field.state.value}
                onChange={(range) => field.handleChange(range)}
                placeholder={tOrders('select_date_range')}
              />
            </div>
          )}
        </form.AppField>

        <form.AppField name="status">
          {(field) => (
            <field.CheckboxGroupField
              label={tOrders('filter_by_status')}
              options={[
                { value: 'pending', label: tOrders('status_pending') },
                { value: 'placed', label: tOrders('status_placed') },
                {
                  value: 'ready_for_shipment',
                  label: tOrders('status_ready_for_shipment'),
                },
                { value: 'shipped', label: tOrders('status_shipped') },
                { value: 'fulfilled', label: tOrders('status_fulfilled') },
                { value: 'cancelled', label: tOrders('status_cancelled') },
                { value: 'returned', label: tOrders('status_returned') },
              ]}
            />
          )}
        </form.AppField>

        <form.AppField name="paymentStatus">
          {(field) => (
            <field.CheckboxGroupField
              label={tOrders('filter_by_payment_status')}
              options={[
                { value: 'pending', label: tOrders('payment_status_pending') },
                { value: 'paid', label: tOrders('payment_status_paid') },
                { value: 'failed', label: tOrders('payment_status_failed') },
                {
                  value: 'refunded',
                  label: tOrders('payment_status_refunded'),
                },
              ]}
            />
          )}
        </form.AppField>

        <form.AppField name="socialPlatforms">
          {(field) => (
            <field.CheckboxGroupField
              label={tOrders('filter_by_platform')}
              description={tOrders('filter_by_platform_desc')}
              options={[
                {
                  value: 'instagram' as const,
                  label: tOrders('platform_instagram'),
                },
                { value: 'tiktok' as const, label: tOrders('platform_tiktok') },
                {
                  value: 'facebook' as const,
                  label: tOrders('platform_facebook'),
                },
                { value: 'x' as const, label: tOrders('platform_x') },
                {
                  value: 'snapchat' as const,
                  label: tOrders('platform_snapchat'),
                },
                {
                  value: 'whatsapp' as const,
                  label: tOrders('platform_whatsapp'),
                },
              ]}
            />
          )}
        </form.AppField>
      </form.FormRoot>
    </form.AppForm>
  )

  return (
    <>
      <ResourceListLayout
        title={tOrders('title')}
        subtitle={tOrders('subtitle')}
        addButtonText={tOrders('add_order')}
        onAddClick={handleAddOrder}
        searchPlaceholder={tOrders('search_placeholder')}
        searchValue={search.search ?? ''}
        onSearchChange={handleSearchChange}
        filterTitle={tOrders('filters')}
        filterButtonText={tCommon('filter')}
        filterButton={filterContent}
        activeFilterCount={activeFilterCount}
        onApplyFilters={handleApplyFilters}
        onResetFilters={handleResetFilters}
        applyLabel={tCommon('apply')}
        resetLabel={tCommon('reset')}
        sortTitle={tOrders('sort_orders')}
        sortOptions={sortOptions}
        onSortApply={handleSortApply}
        emptyIcon={<ShoppingCart size={48} />}
        emptyTitle={
          search.search ? tOrders('no_results') : tOrders('no_orders')
        }
        emptyMessage={
          search.search
            ? tOrders('try_different_search')
            : tOrders('get_started_message')
        }
        emptyActionText={!search.search ? tOrders('add_order') : undefined}
        onEmptyAction={!search.search ? handleAddOrder : undefined}
        tableColumns={columns}
        tableData={data?.items || []}
        tableKeyExtractor={(order) => order.id}
        tableSortBy={search.sortBy}
        tableSortOrder={search.sortOrder}
        onTableSort={(column: string) =>
          handleSortChange(column as OrdersSearch['sortBy'])
        }
        mobileCard={(order: Order) => (
          <OrderCard
            order={order}
            businessDescriptor={businessDescriptor}
            onReviewClick={setReviewOrder}
            onDeleteSuccess={() => {
              void queryClient.invalidateQueries({ queryKey: orderQueries.all })
            }}
          />
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
        itemsName={tOrders('title')}
      />

      <OrderReviewSheet
        order={reviewOrder}
        isOpen={reviewOrder !== null}
        onClose={() => setReviewOrder(null)}
      />

      <CreateOrderSheet
        isOpen={isCreateOrderOpen}
        onClose={() => setIsCreateOrderOpen(false)}
        onCreated={async () => {
          await queryClient.invalidateQueries({
            queryKey: orderQueries.all,
          })
        }}
      />
    </>
  )
}
