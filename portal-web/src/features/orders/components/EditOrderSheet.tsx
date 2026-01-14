import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { MapPin } from 'lucide-react'

import type { CustomerAddress } from '@/api/customer'
import type { DiscountType, Order, UpdateOrderRequest } from '@/api/order'
import { useUpdateOrderMutation } from '@/api/order'
import { useAddressesQuery } from '@/api/address'
import { useShippingZonesQuery } from '@/api/business'
import { BottomSheet } from '@/components'
import { ShippingZoneInfo } from '@/components/organisms/ShippingZoneInfo'
import { showSuccessToast } from '@/lib/toast'
import { businessStore } from '@/stores/businessStore'
import { useKyoraForm } from '@/lib/form'
import { inferShippingZoneFromAddress } from '@/lib/shippingZone'

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

  const [selectedAddressId, setSelectedAddressId] = useState<string>('')
  const [selectedAddress, setSelectedAddress] =
    useState<CustomerAddress | null>(null)

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
      discountType: getDiscountType(),
      discountValue: getDiscountValue(),
    },
    onSubmit: async ({ value }) => {
      if (!order) return

      const payload: UpdateOrderRequest = {
        shippingAddressId: value.shippingAddressId || undefined,
        shippingZoneId: value.shippingZoneId || undefined,
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
      const addressId = order.shippingAddress?.id || ''
      form.setFieldValue('shippingAddressId', addressId)
      form.setFieldValue('shippingZoneId', order.shippingZone?.id || '')
      form.setFieldValue('discountType', getDiscountType())
      form.setFieldValue('discountValue', getDiscountValue())

      // Initialize selected address state
      setSelectedAddressId(addressId)
      if (addressId && addressesData) {
        const address = addressesData.find(
          (a: CustomerAddress) => a.id === addressId,
        )
        setSelectedAddress(address || null)
      }
    }
  }, [order, isOpen, addressesData])

  // Handle address selection changes and auto-infer shipping zone
  useEffect(() => {
    const addressId = form.getFieldValue('shippingAddressId')

    if (addressId !== selectedAddressId) {
      setSelectedAddressId(addressId)

      if (addressId && addressesData && addressesData.length > 0) {
        // Find the selected address
        const address = addressesData.find(
          (a: CustomerAddress) => a.id === addressId,
        )
        setSelectedAddress(address || null)

        if (address && shippingZones.length > 0) {
          // Auto-infer shipping zone from address
          const inferredZone = inferShippingZoneFromAddress(
            address,
            shippingZones,
          )
          if (inferredZone) {
            form.setFieldValue('shippingZoneId', inferredZone.id)
          } else {
            form.setFieldValue('shippingZoneId', '')
          }
        }
      } else {
        // Address cleared - reset zone
        setSelectedAddress(null)
        form.setFieldValue('shippingZoneId', '')
      }
    }
  }, [
    form.getFieldValue('shippingAddressId'),
    selectedAddressId,
    addressesData,
    shippingZones,
    form,
  ])

  useEffect(() => {
    if (!isOpen) {
      form.reset()
      setSelectedAddressId('')
      setSelectedAddress(null)
    }
  }, [isOpen])

  const safeClose = () => {
    if (updateMutation.isPending) return
    onClose()
  }

  if (!order) return null

  // Determine which fields are editable based on order status
  const isShipped = ['shipped', 'fulfilled', 'returned'].includes(order.status)
  const isCancelled = order.status === 'cancelled'
  const canEditAddress = !isShipped && !isCancelled
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
            <form.SubmitButton
              form="edit-order-form"
              variant="primary"
              className="flex-1"
            >
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
        <form.FormRoot id="edit-order-form" className="space-y-6">
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
          {canEditAddress && order.customer && (
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <MapPin size={18} className="text-base-content/70" />
                <h3 className="font-semibold text-start">{tOrders('shipping_address')}</h3>
              </div>

              <form.AppField name="shippingAddressId">
                {(field) => (
                  <field.AddressSelectField
                    label={tOrders('select_address')}
                    businessDescriptor={businessDescriptor}
                    customerId={order.customer!.id}
                    placeholder={tOrders('search_address')}
                  />
                )}
              </form.AppField>

              {/* Inferred Shipping Zone (Read-only Display) */}
              {selectedAddress && shippingZones.length > 0 && (
                <form.Subscribe
                  selector={(state) => state.values.shippingZoneId}
                >
                  {(shippingZoneId) => {
                    const inferredZone = shippingZones.find(
                      (z) => z.id === shippingZoneId,
                    )
                    return (
                      <ShippingZoneInfo
                        zone={inferredZone}
                        currency={inferredZone?.currency}
                      />
                    )
                  }}
                </form.Subscribe>
              )}
            </div>
          )}

          {/* Discount */}
          {canEditDiscount && (
            <div className="space-y-3">
              <h3 className="font-semibold text-start">{tOrders('discount')}</h3>
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

                <form.Subscribe selector={(state) => state.values.discountType}>
                  {(discountType) => (
                    <form.AppField name="discountValue">
                      {(field) =>
                        discountType === 'amount' ? (
                          <field.PriceField
                            label={tOrders('discount_value')}
                            placeholder="0"
                          />
                        ) : (
                          <field.TextField
                            label={tOrders('discount_value')}
                            type="text"
                            inputMode="decimal"
                            placeholder="0"
                            dir="ltr"
                            endIcon={
                              <span className="text-base-content/70">%</span>
                            }
                          />
                        )
                      }
                    </form.AppField>
                  )}
                </form.Subscribe>
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
