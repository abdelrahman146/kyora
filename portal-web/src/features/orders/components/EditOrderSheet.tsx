import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { MapPin } from 'lucide-react'

import type { CustomerAddress } from '@/api/customer'
import type { DiscountType, Order, UpdateOrderRequest } from '@/api/order'
import { useUpdateOrderMutation } from '@/api/order'
import { useAddressesQuery } from '@/api/address'
import { useShippingZonesQuery } from '@/api/business'
import { BottomSheet } from '@/components'
import { showSuccessToast } from '@/lib/toast'
import { businessStore } from '@/stores/businessStore'
import { useKyoraForm } from '@/lib/form'

export interface EditOrderSheetProps {
  isOpen: boolean
  onClose: () => void
  order: Order | null
  onUpdated?: () => void | Promise<void>
}

export function EditOrderSheet({
  isOpen,
  onClose,
  order,
  onUpdated,
}: EditOrderSheetProps) {
  const { t: tOrders } = useTranslation('orders')
  const { t: tCommon } = useTranslation('common')

  const businessDescriptor =
    businessStore.state.selectedBusinessDescriptor || ''

  const { data: addressesData } = useAddressesQuery(
    businessDescriptor,
    order?.customer?.id || '',
    isOpen,
  )
  const { data: shippingZones = [] } = useShippingZonesQuery(
    businessDescriptor,
    isOpen,
  )

  const updateMutation = useUpdateOrderMutation(
    businessDescriptor,
    order?.id || '',
    {
      onSuccess: async () => {
        showSuccessToast(tOrders('order_updated_success'))
        await onUpdated?.()
        onClose()
      },
    },
  )

  // Parse discount from order
  const getDiscountType = (): DiscountType => {
    // If order has discount string like "10%", it's percent
    if (order?.discount && order.discount.includes('%')) return 'percent'
    return 'amount'
  }

  const getDiscountValue = (): string => {
    if (!order?.discount) return ''
    // Remove % sign if present and extract numeric value
    return order.discount.replace('%', '').replace(/[^0-9.]/g, '')
  }

  const form = useKyoraForm({
    defaultValues: {
      shippingAddressId: order?.shippingAddress?.id || '',
      shippingZoneId: order?.shippingZone?.id || '',
      shippingFee: order?.shippingFee || '',
      discountType: getDiscountType(),
      discountValue: getDiscountValue(),
    },
    onSubmit: async ({ value }) => {
      if (!order) return

      const payload: UpdateOrderRequest = {
        shippingAddressId: value.shippingAddressId || undefined,
        shippingZoneId: value.shippingZoneId || undefined,
        shippingFee: value.shippingFee || undefined,
        discountType:
          value.discountValue && value.discountValue.trim() !== ''
            ? value.discountType
            : undefined,
        discountValue:
          value.discountValue && value.discountValue.trim() !== ''
            ? value.discountValue
            : undefined,
      }
      await updateMutation.mutateAsync(payload)
    },
  })

  useEffect(() => {
    if (order && isOpen) {
      form.setFieldValue('shippingAddressId', order.shippingAddress?.id || '')
      form.setFieldValue('shippingZoneId', order.shippingZone?.id || '')
      form.setFieldValue('shippingFee', order.shippingFee || '')
      form.setFieldValue('discountType', getDiscountType())
      form.setFieldValue('discountValue', getDiscountValue())
    }
  }, [order, isOpen])

  useEffect(() => {
    if (!isOpen) {
      form.reset()
    }
  }, [isOpen])

  const safeClose = () => {
    if (updateMutation.isPending) return
    onClose()
  }

  if (!order) return null

  const addressOptions =
    addressesData?.map((addr: CustomerAddress) => ({
      value: addr.id,
      label: `${addr.street}, ${addr.city}`,
    })) || []

  // Determine which fields are editable based on order status
  const isShipped = ['shipped', 'fulfilled', 'returned'].includes(order.status)
  const isCancelled = order.status === 'cancelled'
  const canEditAddress = !isShipped && !isCancelled
  const canEditShipping = !isCancelled
  const canEditDiscount = !isCancelled

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={safeClose}
        title={tOrders('edit_order_title')}
        footer={
          <div className="flex gap-2">
            <button
              type="button"
              className="btn btn-ghost flex-1"
              onClick={safeClose}
              disabled={updateMutation.isPending}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton variant="primary" className="flex-1">
              {updateMutation.isPending
                ? tOrders('update_submitting')
                : tOrders('update_submit')}
            </form.SubmitButton>
          </div>
        }
        side="end"
        size="lg"
        closeOnOverlayClick={!updateMutation.isPending}
        closeOnEscape={!updateMutation.isPending}
        contentClassName="space-y-6"
      >
        <form.FormRoot className="space-y-6">
          <form.FormError />

          {/* Order Info (Read-only) */}
          <div className="rounded-lg bg-base-200 p-4 space-y-2">
            <div className="text-sm">
              <span className="text-base-content/70">
                {tOrders('order_number')}:{' '}
              </span>
              <span className="font-medium" dir="ltr">
                #{order.orderNumber}
              </span>
            </div>
            {order.customer && (
              <div className="text-sm">
                <span className="text-base-content/70">
                  {tOrders('customer')}:{' '}
                </span>
                <span className="font-medium">{order.customer.name}</span>
              </div>
            )}
            <div className="text-sm">
              <span className="text-base-content/70">
                {tOrders('status')}:{' '}
              </span>
              <span
                className={`badge ${
                  order.status === 'fulfilled'
                    ? 'badge-success'
                    : order.status === 'cancelled' ||
                        order.status === 'returned'
                      ? 'badge-error'
                      : order.status === 'shipped'
                        ? 'badge-info'
                        : 'badge-warning'
                }`}
              >
                {tOrders(`status_${order.status}`)}
              </span>
            </div>
          </div>

          {/* Shipping Address (Editable pre-shipped) */}
          {canEditAddress && addressOptions.length > 0 && (
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <MapPin size={18} className="text-base-content/70" />
                <h3 className="font-semibold">{tOrders('shipping_address')}</h3>
              </div>

              <form.AppField name="shippingAddressId">
                {(field) => (
                  <field.SelectField
                    label={tOrders('select_address')}
                    options={addressOptions}
                    searchable
                    clearable
                  />
                )}
              </form.AppField>
            </div>
          )}

          {/* Shipping Zone */}
          {canEditShipping && shippingZones.length > 0 && (
            <form.AppField name="shippingZoneId">
              {(field) => (
                <field.SelectField
                  label={tOrders('shipping_zone')}
                  options={[
                    { value: '', label: tCommon('none') },
                    ...shippingZones.map((zone) => ({
                      value: zone.id,
                      label: zone.name,
                    })),
                  ]}
                  clearable
                />
              )}
            </form.AppField>
          )}

          {/* Discount */}
          {canEditDiscount && (
            <div className="space-y-3">
              <h3 className="font-semibold">{tOrders('discount')}</h3>
              <div className="grid grid-cols-2 gap-3">
                <form.AppField name="discountType">
                  {(field) => (
                    <field.SelectField
                      label={tOrders('discount_type')}
                      options={[
                        { value: 'amount', label: tOrders('discount_amount') },
                        {
                          value: 'percent',
                          label: tOrders('discount_percent'),
                        },
                      ]}
                    />
                  )}
                </form.AppField>

                <form.AppField name="discountValue">
                  {(field) => (
                    <field.TextField
                      label={tOrders('discount_value')}
                      type="text"
                      inputMode="decimal"
                      dir="ltr"
                      placeholder="0"
                    />
                  )}
                </form.AppField>
              </div>
            </div>
          )}

          {/* Info about restrictions */}
          {(isShipped || isCancelled) && (
            <div className="alert alert-info">
              <span className="text-sm">
                {isShipped
                  ? tOrders('edit_restriction_shipped')
                  : tOrders('edit_restriction_cancelled')}
              </span>
            </div>
          )}
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
