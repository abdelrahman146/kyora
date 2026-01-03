/**
 * OrderCard Component
 *
 * A card component displaying order information for mobile list view.
 * Designed for social commerce entrepreneurs to quickly understand order status.
 *
 * Features:
 * - Clean layout with order number, customer, and total at top
 * - Order and payment status badges with color coding
 * - Item count and channel information
 * - Date ordered displayed
 * - RTL-compatible with logical properties
 * - Mobile-optimized (min 44px touch target)
 */

import { useTranslation } from 'react-i18next'
import { Calendar, Package, ShoppingBag, User } from 'lucide-react'

import { Avatar } from '../atoms/Avatar'
import type { Order } from '@/api/order'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export interface OrderCardProps {
  order: Order
  onClick?: (order: Order) => void
}

export function OrderCard({ order, onClick }: OrderCardProps) {
  const { t } = useTranslation()

  const getOrderStatusBadge = (status: Order['status']) => {
    const statusMap: Record<Order['status'], { class: string; label: string }> =
      {
        pending: { class: 'badge-warning', label: 'pending' },
        placed: { class: 'badge-info', label: 'placed' },
        ready_for_shipment: {
          class: 'badge-info',
          label: 'ready_for_shipment',
        },
        shipped: { class: 'badge-primary', label: 'shipped' },
        fulfilled: { class: 'badge-success', label: 'fulfilled' },
        cancelled: { class: 'badge-error', label: 'cancelled' },
        returned: { class: 'badge-error', label: 'returned' },
      }
    const config = statusMap[status]
    return (
      <span className={`badge badge-sm ${config.class}`}>
        {t(`orders:status_${config.label}`)}
      </span>
    )
  }

  const getPaymentStatusBadge = (status: Order['paymentStatus']) => {
    const statusMap: Record<
      Order['paymentStatus'],
      { class: string; label: string }
    > = {
      pending: { class: 'badge-warning', label: 'pending' },
      paid: { class: 'badge-success', label: 'paid' },
      failed: { class: 'badge-error', label: 'failed' },
      refunded: { class: 'badge-ghost', label: 'refunded' },
    }
    const config = statusMap[status]
    return (
      <span className={`badge badge-sm ${config.class}`}>
        {t(`orders:payment_status_${config.label}`)}
      </span>
    )
  }

  const getCustomerInitials = (name: string): string => {
    return name
      .split(' ')
      .map((word) => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  const itemsCount = order.items?.length ?? 0

  return (
    <div
      className={`bg-base-100 border border-base-300 rounded-xl p-4 hover:shadow-md transition-shadow ${
        onClick ? 'cursor-pointer active:scale-[0.98]' : ''
      }`}
      onClick={() => onClick?.(order)}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      onKeyDown={(e) => {
        if (onClick && (e.key === 'Enter' || e.key === ' ')) {
          e.preventDefault()
          onClick(order)
        }
      }}
    >
      {/* Header: Order Number + Total */}
      <div className="flex items-start justify-between gap-3 mb-3">
        <div className="flex-1 min-w-0">
          <h3 className="font-bold text-base text-base-content truncate mb-1">
            {order.orderNumber}
          </h3>
          <div className="flex items-center gap-2 flex-wrap">
            {getOrderStatusBadge(order.status)}
            {getPaymentStatusBadge(order.paymentStatus)}
          </div>
        </div>
        <div className="text-end">
          <div className="text-xs text-base-content/60 mb-1">
            {t('orders:total')}
          </div>
          <div className="font-bold text-lg text-primary">
            {formatCurrency(parseFloat(order.total), order.currency)}
          </div>
        </div>
      </div>

      {/* Customer Info */}
      {order.customer && (
        <div className="flex items-center gap-2 mb-3 p-2 bg-base-200 rounded-lg">
          <Avatar
            src={order.customer.avatarUrl}
            alt={order.customer.name}
            fallback={getCustomerInitials(order.customer.name)}
            size="sm"
          />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-1 text-sm text-base-content/70">
              <User size={14} />
              <span className="truncate font-medium">
                {order.customer.name}
              </span>
            </div>
          </div>
        </div>
      )}

      {/* Order Details Grid */}
      <div className="grid grid-cols-2 gap-3 mb-3">
        {/* Items Count */}
        <div className="flex items-center gap-2">
          <div className="p-1.5 bg-primary/10 rounded-lg">
            <Package size={16} className="text-primary" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="text-xs text-base-content/60">
              {t('orders:items')}
            </div>
            <div className="font-semibold text-sm">{itemsCount}</div>
          </div>
        </div>

        {/* Channel */}
        <div className="flex items-center gap-2">
          <div className="p-1.5 bg-secondary/10 rounded-lg">
            <ShoppingBag size={16} className="text-secondary" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="text-xs text-base-content/60">
              {t('orders:channel')}
            </div>
            <div className="font-semibold text-sm truncate">
              {order.channel}
            </div>
          </div>
        </div>
      </div>

      {/* Pricing Details */}
      <div className="space-y-1 mb-3 p-2 bg-base-200 rounded-lg">
        <div className="flex items-baseline justify-between text-sm">
          <span className="text-base-content/60">{t('orders:subtotal')}</span>
          <span className="font-medium">
            {formatCurrency(parseFloat(order.subtotal), order.currency)}
          </span>
        </div>
        {parseFloat(order.shippingFee) > 0 && (
          <div className="flex items-baseline justify-between text-sm">
            <span className="text-base-content/60">
              {t('orders:shipping_fee')}
            </span>
            <span className="font-medium">
              {formatCurrency(parseFloat(order.shippingFee), order.currency)}
            </span>
          </div>
        )}
        {parseFloat(order.discount) > 0 && (
          <div className="flex items-baseline justify-between text-sm">
            <span className="text-base-content/60">{t('orders:discount')}</span>
            <span className="font-medium text-success">
              -{formatCurrency(parseFloat(order.discount), order.currency)}
            </span>
          </div>
        )}
        {parseFloat(order.vat) > 0 && (
          <div className="flex items-baseline justify-between text-sm">
            <span className="text-base-content/60">{t('orders:vat')}</span>
            <span className="font-medium">
              {formatCurrency(parseFloat(order.vat), order.currency)}
            </span>
          </div>
        )}
      </div>

      {/* Date Ordered */}
      <div className="flex items-center gap-2 text-sm text-base-content/60 pt-2 border-t border-base-300">
        <Calendar size={14} />
        <span>{formatDateShort(order.orderedAt)}</span>
      </div>
    </div>
  )
}
