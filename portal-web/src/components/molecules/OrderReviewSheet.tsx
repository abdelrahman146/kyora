import { useTranslation } from 'react-i18next'
import { FileText, MapPin, Package, User } from 'lucide-react'

import { BottomSheet } from './BottomSheet'
import type { Order } from '@/api/order'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export interface OrderReviewSheetProps {
  order: Order | null
  isOpen: boolean
  onClose: () => void
}

export function OrderReviewSheet({
  order,
  isOpen,
  onClose,
}: OrderReviewSheetProps) {
  const { t: tOrders } = useTranslation('orders')
  const { t: tCommon } = useTranslation('common')

  if (!order) return null

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

  const latestTimestamp =
    order.fulfilledAt ||
    order.shippedAt ||
    order.readyForShipmentAt ||
    order.placedAt ||
    order.orderedAt ||
    order.createdAt

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={onClose}
      title={tOrders('quick_review')}
      size="lg"
    >
      <div className="space-y-6">
        <div className="bg-base-200 rounded-lg p-4 space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-sm text-base-content/70">
              {tOrders('order_number')}
            </span>
            <span className="font-bold text-lg">{order.orderNumber}</span>
          </div>
          <div className="flex gap-2">
            <span className={`badge ${getStatusBadgeClass(order.status)}`}>
              {tOrders(`status_${order.status}`)}
            </span>
            <span
              className={`badge ${getPaymentStatusBadgeClass(order.paymentStatus)}`}
            >
              {tOrders(`payment_status_${order.paymentStatus}`)}
            </span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-base-content/70">
              {tOrders('last_updated')}
            </span>
            <span>{formatDateShort(latestTimestamp)}</span>
          </div>
        </div>

        {order.customer && (
          <div>
            <div className="flex items-center gap-2 mb-3">
              <User size={18} className="text-base-content/70" />
              <h3 className="font-semibold">{tOrders('customer')}</h3>
            </div>
            <div className="bg-base-100 border border-base-300 rounded-lg p-3 space-y-2">
              <div className="font-medium">{order.customer.name}</div>
              {order.customer.phoneNumber && (
                <div className="text-sm text-base-content/70">
                  <span dir="ltr">
                    {order.customer.phoneCode} {order.customer.phoneNumber}
                  </span>
                </div>
              )}
              {order.customer.email && (
                <div className="text-sm text-base-content/70">
                  {order.customer.email}
                </div>
              )}
            </div>
          </div>
        )}

        {order.shippingAddress && (
          <div>
            <div className="flex items-center gap-2 mb-3">
              <MapPin size={18} className="text-base-content/70" />
              <h3 className="font-semibold">{tOrders('shipping_address')}</h3>
            </div>
            <div className="bg-base-100 border border-base-300 rounded-lg p-3 text-sm">
              {order.shippingAddress.street && (
                <div>{order.shippingAddress.street}</div>
              )}
              <div>
                {order.shippingAddress.city}, {order.shippingAddress.state}
              </div>
              {order.shippingAddress.zipCode && (
                <div>{order.shippingAddress.zipCode}</div>
              )}
              <div>{order.shippingAddress.countryCode}</div>
              <div className="mt-2 text-base-content/70">
                <span dir="ltr">
                  {order.shippingAddress.phoneCode}{' '}
                  {order.shippingAddress.phoneNumber}
                </span>
              </div>
            </div>
          </div>
        )}

        {order.items && order.items.length > 0 && (
          <div>
            <div className="flex items-center gap-2 mb-3">
              <Package size={18} className="text-base-content/70" />
              <h3 className="font-semibold">
                {tOrders('items')} ({order.items.length})
              </h3>
            </div>
            <div className="space-y-2">
              {order.items.map((item) => (
                <div
                  key={item.id}
                  className="flex items-center gap-3 bg-base-100 border border-base-300 rounded-lg p-3"
                >
                  {item.product?.photos && item.product.photos[0] && (
                    <img
                      src={
                        item.product.photos[0].thumbnailUrl ||
                        item.product.photos[0].url
                      }
                      alt={item.product.name}
                      className="w-12 h-12 rounded object-cover"
                    />
                  )}
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-sm truncate">
                      {item.product?.name || tCommon('unknown')}
                    </div>
                    {item.variant && (
                      <div className="text-xs text-base-content/70">
                        {item.variant.name}
                      </div>
                    )}
                    <div className="text-xs text-base-content/60">
                      {tOrders('quantity')}: {item.quantity}
                    </div>
                  </div>
                  <div className="text-end">
                    <div className="font-semibold text-sm">
                      {formatCurrency(parseFloat(item.total), order.currency)}
                    </div>
                    <div className="text-xs text-base-content/60">
                      {formatCurrency(
                        parseFloat(item.unitPrice),
                        order.currency,
                      )}{' '}
                      × {item.quantity}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */}
        {order.paymentMethod && (
          <div className="bg-base-100 border border-base-300 rounded-lg p-3">
            <div className="flex justify-between text-sm">
              <span className="text-base-content/70">
                {tOrders('payment_method')}
              </span>
              <span className="font-medium">
                {tOrders(`payment_method_${order.paymentMethod}`)}
              </span>
            </div>
            {order.paymentReference && (
              <div className="flex justify-between text-sm mt-2">
                <span className="text-base-content/70">
                  {tOrders('payment_reference')}
                </span>
                <span className="font-medium font-mono text-xs">
                  {order.paymentReference}
                </span>
              </div>
            )}
          </div>
        )}

        {order.notes && order.notes.length > 0 && (
          <div>
            <div className="flex items-center gap-2 mb-3">
              <FileText size={18} className="text-base-content/70" />
              <h3 className="font-semibold">{tOrders('notes')}</h3>
            </div>
            <div className="space-y-2">
              {order.notes.map((note) => (
                <div
                  key={note.id}
                  className="bg-base-100 border border-base-300 rounded-lg p-3"
                >
                  <div className="text-sm whitespace-pre-wrap">
                    {note.content}
                  </div>
                  <div className="text-xs text-base-content/60 mt-2">
                    {formatDateShort(note.createdAt)}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        <div>
          <h3 className="font-semibold mb-3">{tOrders('checkout_summary')}</h3>
          <div className="bg-base-200 rounded-lg p-4 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-base-content/70">
                {tOrders('subtotal')}
              </span>
              <span className="font-medium">
                {formatCurrency(parseFloat(order.subtotal), order.currency)}
              </span>
            </div>
            {parseFloat(order.shippingFee) > 0 && (
              <div className="flex justify-between text-sm">
                <span className="text-base-content/70">
                  {tOrders('shipping_fee')}
                </span>
                <span className="font-medium">
                  {formatCurrency(
                    parseFloat(order.shippingFee),
                    order.currency,
                  )}
                </span>
              </div>
            )}
            {parseFloat(order.discount) > 0 && (
              <div className="flex justify-between text-sm">
                <span className="text-base-content/70">
                  {tOrders('discount')}
                </span>
                <span className="font-medium text-success">
                  -{formatCurrency(parseFloat(order.discount), order.currency)}
                </span>
              </div>
            )}
            {parseFloat(order.vat) > 0 && (
              <div className="flex justify-between text-sm">
                <span className="text-base-content/70">
                  {tOrders('vat')} ({parseFloat(order.vatRate) * 100}%)
                </span>
                <span className="font-medium">
                  {formatCurrency(parseFloat(order.vat), order.currency)}
                </span>
              </div>
            )}
            <div className="divider my-2" />
            <div className="flex justify-between items-center">
              <span className="font-semibold text-base">
                {tOrders('total')}
              </span>
              <span className="font-bold text-xl text-primary">
                {formatCurrency(parseFloat(order.total), order.currency)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </BottomSheet>
  )
}
