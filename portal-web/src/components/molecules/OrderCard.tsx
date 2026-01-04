/**
 * OrderCard Component
 *
 * Mobile-first card component for order list display.
 * Shows all essential order information with explicit action buttons.
 *
 * Features:
 * - Order number, customer, platform info
 * - Status badges with color coding
 * - Payment status and method
 * - Order timestamps (ordered + latest update)
 * - Explicit "Quick Review" and "Actions" buttons
 * - RTL-compatible with logical properties
 * - Mobile-optimized touch targets
 * - Professional card layout with clear visual hierarchy
 */

import { useTranslation } from 'react-i18next'
import { Link } from '@tanstack/react-router'
import { Calendar, Clock, Eye, Package, User } from 'lucide-react'
import {
  FaFacebook,
  FaInstagram,
  FaSnapchat,
  FaTiktok,
  FaWhatsapp,
} from 'react-icons/fa'
import { FaXTwitter } from 'react-icons/fa6'

import { OrderQuickActions } from './OrderQuickActions'
import type { Order } from '@/api/order'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export interface OrderCardProps {
  order: Order
  businessDescriptor: string
  onReviewClick?: (order: Order) => void
  onDeleteSuccess?: () => void
}

export function OrderCard({
  order,
  businessDescriptor,
  onReviewClick,
  onDeleteSuccess,
}: OrderCardProps) {
  const { t } = useTranslation('orders')

  const getStatusBadgeClass = (status: Order['status']) => {
    const statusMap: Record<Order['status'], string> = {
      pending: 'badge-warning',
      placed: 'badge-info',
      ready_for_shipment: 'badge-info',
      shipped: 'badge-primary',
      fulfilled: 'badge-success',
      cancelled: 'badge-error',
      returned: 'badge-error',
    }
    return statusMap[status]
  }

  const getPaymentStatusBadgeClass = (status: Order['paymentStatus']) => {
    const statusMap: Record<Order['paymentStatus'], string> = {
      pending: 'badge-warning',
      paid: 'badge-success',
      failed: 'badge-error',
      refunded: 'badge-ghost',
    }
    return statusMap[status]
  }

  const getPlatformIcon = (
    instagramUsername: string | null,
    tiktokUsername: string | null,
    facebookUsername: string | null,
    xUsername: string | null,
    snapchatUsername: string | null,
    whatsappNumber: string | null,
  ) => {
    if (instagramUsername)
      return <FaInstagram size={16} className="text-pink-600" />
    if (tiktokUsername) return <FaTiktok size={16} className="text-black" />
    if (facebookUsername)
      return <FaFacebook size={16} className="text-blue-600" />
    if (xUsername) return <FaXTwitter size={16} className="text-black" />
    if (snapchatUsername)
      return <FaSnapchat size={16} className="text-yellow-400" />
    if (whatsappNumber)
      return <FaWhatsapp size={16} className="text-green-600" />
    return null
  }

  const getPlatformHandle = (
    instagramUsername: string | null,
    tiktokUsername: string | null,
    facebookUsername: string | null,
    xUsername: string | null,
    snapchatUsername: string | null,
  ) => {
    if (instagramUsername) return `@${instagramUsername}`
    if (tiktokUsername) return `@${tiktokUsername}`
    if (facebookUsername) return `@${facebookUsername}`
    if (xUsername) return `@${xUsername}`
    if (snapchatUsername) return `@${snapchatUsername}`
    return null
  }

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

  const platformIcon = order.customer
    ? getPlatformIcon(
        order.customer.instagramUsername ?? null,
        order.customer.tiktokUsername ?? null,
        order.customer.facebookUsername ?? null,
        order.customer.xUsername ?? null,
        order.customer.snapchatUsername ?? null,
        order.customer.whatsappNumber ?? null,
      )
    : null

  const platformHandle = order.customer
    ? getPlatformHandle(
        order.customer.instagramUsername ?? null,
        order.customer.tiktokUsername ?? null,
        order.customer.facebookUsername ?? null,
        order.customer.xUsername ?? null,
        order.customer.snapchatUsername ?? null,
      )
    : null

  return (
    <div className="bg-base-100 border border-base-300 rounded-xl">
      {/* Header: Order Number + Total */}
      <div className="flex items-center justify-between gap-4 px-4 py-3 border-b border-base-300">
        <div className="flex-1">
          <div className="flex items-center gap-1.5 text-xs text-base-content/60 mb-1">
            <Package size={14} />
            <span>{t('order_number')}</span>
          </div>
          <h3 className="font-bold text-lg text-base-content">
            {order.orderNumber}
          </h3>
        </div>
        <div className="text-end">
          <div className="text-xs text-base-content/60 mb-0.5">
            {t('total')}
          </div>
          <div className="font-bold text-xl text-primary">
            {formatCurrency(parseFloat(order.total), order.currency)}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="px-4 py-3 space-y-3">
        {/* Status Badges */}
        <div className="flex items-center gap-2 flex-wrap">
          <span className={`badge ${getStatusBadgeClass(order.status)}`}>
            {t(`status_${order.status}`)}
          </span>
          <span
            className={`badge ${getPaymentStatusBadgeClass(order.paymentStatus)}`}
          >
            {t(`payment_status_${order.paymentStatus}`)}
          </span>
          {/* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */}
          {order.paymentMethod && order.paymentStatus === 'paid' && (
            <span className="badge badge-ghost">
              {t(`payment_method_${order.paymentMethod}`)}
            </span>
          )}
        </div>

        {/* Customer Info */}
        {order.customer && (
          <Link
            to="/business/$businessDescriptor/customers/$customerId"
            params={{
              businessDescriptor,
              customerId: order.customer.id,
            }}
            className="flex items-center gap-2.5 cursor-pointer hover:bg-base-200 rounded-lg p-2 -mx-2 transition-colors"
          >
            <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
              <User size={16} className="text-primary" />
            </div>
            <div className="flex-1 min-w-0">
              <div className="font-medium text-sm text-base-content truncate">
                {order.customer.name}
              </div>
              {platformHandle && (
                <div className="flex items-center gap-1.5 text-xs text-base-content/60">
                  {platformIcon}
                  <span className="truncate">{platformHandle}</span>
                </div>
              )}
            </div>
          </Link>
        )}

        {/* Timestamps */}
        <div className="flex items-center justify-between text-xs text-base-content/60">
          <div className="flex items-center gap-1.5">
            <Calendar size={14} className="flex-shrink-0" />
            <span>{formatDateShort(order.orderedAt)}</span>
          </div>
          {latestTimestamp !== order.orderedAt && (
            <div className="flex items-center gap-1.5">
              <Clock size={14} className="flex-shrink-0" />
              <span>{formatDateShort(latestTimestamp)}</span>
            </div>
          )}
        </div>
      </div>

      {/* Actions Footer */}
      <div className="flex items-center justify-end gap-2 px-4 py-3 bg-base-200/50 border-t border-base-300 rounded-b-xl">
        <button
          type="button"
          className="btn btn-ghost btn-sm btn-square"
          onClick={() => onReviewClick?.(order)}
          aria-label={t('quick_review')}
        >
          <Eye size={18} />
        </button>
        <OrderQuickActions
          order={order}
          businessDescriptor={businessDescriptor}
          onDeleteSuccess={onDeleteSuccess}
        />
      </div>
    </div>
  )
}
