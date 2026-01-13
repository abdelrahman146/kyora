import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { CreditCard, MapPin, Package, Plus, Trash2, User } from 'lucide-react'

import type {
  CreateOrderRequest,
  DiscountType,
  OrderPaymentMethod,
  OrderPaymentStatus,
  OrderStatus,
} from '@/api/order'
import { useAddressesQuery } from '@/api/address'
import { useCustomersQuery } from '@/api/customer'
import { useCreateOrderMutation } from '@/api/order'
import { usePaymentMethodsQuery, useShippingZonesQuery } from '@/api/business'
import { BottomSheet } from '@/components'
import { useKyoraForm } from '@/lib/form'
import { showSuccessToast } from '@/lib/toast'
import { businessStore } from '@/stores/businessStore'

export interface CreateOrderSheetProps {
  isOpen: boolean
  onClose: () => void
  onCreated?: (orderId: string) => void | Promise<void>
}

const CHANNELS = [
  { value: 'instagram', labelKey: 'channel_instagram' },
  { value: 'whatsapp', labelKey: 'channel_whatsapp' },
  { value: 'facebook', labelKey: 'channel_facebook' },
  { value: 'tiktok', labelKey: 'channel_tiktok' },
  { value: 'in_person', labelKey: 'channel_in_person' },
  { value: 'other', labelKey: 'channel_other' },
]

const ORDER_STATUSES: Array<OrderStatus> = [
  'pending',
  'placed',
  'ready_for_shipment',
  'shipped',
  'fulfilled',
]

const PAYMENT_STATUSES: Array<OrderPaymentStatus> = [
  'pending',
  'paid',
  'failed',
]

export function CreateOrderSheet({
  isOpen,
  onClose,
  onCreated,
}: CreateOrderSheetProps) {
  const { t: tOrders } = useTranslation('orders')
  const { t: tCommon } = useTranslation('common')

  const businessDescriptor =
    businessStore.state.selectedBusinessDescriptor || ''

  const [selectedCustomerId, setSelectedCustomerId] = useState<string>('')

  useCustomersQuery(
    businessDescriptor,
    {
      page: 1,
      pageSize: 100,
    },
    isOpen,
  )
  useAddressesQuery(businessDescriptor, selectedCustomerId, isOpen)
  const { data: paymentMethods = [] } = usePaymentMethodsQuery(
    businessDescriptor,
    isOpen,
  )
  const { data: shippingZones = [] } = useShippingZonesQuery(
    businessDescriptor,
    isOpen,
  )

  const createMutation = useCreateOrderMutation(businessDescriptor, {
    onSuccess: async (order) => {
      showSuccessToast(tOrders('create_success'))
      await onCreated?.(order.id)
      onClose()
    },
  })

  const form = useKyoraForm({
    defaultValues: {
      customerId: '',
      shippingAddressId: '',
      channel: 'instagram' as string,
      shippingZoneId: '' as string,
      shippingFee: '' as string,
      discountType: 'amount' as DiscountType,
      discountValue: '' as string,
      paymentMethod: '' as string,
      paymentReference: '' as string,
      status: '' as string,
      paymentStatus: '' as string,
      note: '' as string,
      items: [
        {
          variantId: '',
          quantity: 1,
          unitPrice: '',
          unitCost: '',
        },
      ],
    },
    onSubmit: async ({ value }) => {
      const payload: CreateOrderRequest = {
        customerId: value.customerId,
        shippingAddressId: value.shippingAddressId,
        channel: value.channel,
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
        paymentMethod:
          (value.paymentMethod as OrderPaymentMethod | undefined) || undefined,
        paymentReference: value.paymentReference || undefined,
        status: (value.status as OrderStatus | undefined) || undefined,
        paymentStatus:
          (value.paymentStatus as OrderPaymentStatus | undefined) || undefined,
        note: value.note || undefined,
        items: value.items,
      }
      await createMutation.mutateAsync(payload)
    },
  })

  const addItem = () => {
    form.pushFieldValue('items', {
      variantId: '',
      quantity: 1,
      unitPrice: '',
      unitCost: '',
    })
  }

  const removeItem = (index: number) => {
    form.removeFieldValue('items', index)
  }

  useEffect(() => {
    if (!isOpen) {
      form.reset()
      setSelectedCustomerId('')
    }
  }, [isOpen, form])

  useEffect(() => {
    const customerId = form.getFieldValue('customerId')
    if (customerId !== selectedCustomerId) {
      setSelectedCustomerId(customerId)
      // Reset address when customer changes
      if (customerId) {
        form.setFieldValue('shippingAddressId', '')
      }
    }
  }, [form.getFieldValue('customerId'), selectedCustomerId, form])

  const safeClose = () => {
    if (createMutation.isPending) return
    onClose()
  }

  const enabledPaymentMethods = paymentMethods.filter((pm) => pm.enabled)

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={safeClose}
        title={tOrders('create_order_title')}
        footer={
          <div className="flex gap-2">
            <button
              type="button"
              className="btn btn-ghost flex-1"
              onClick={safeClose}
              disabled={createMutation.isPending}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton variant="primary" className="flex-1">
              {createMutation.isPending
                ? tOrders('create_submitting')
                : tOrders('create_submit')}
            </form.SubmitButton>
          </div>
        }
        side="end"
        size="lg"
        closeOnOverlayClick={!createMutation.isPending}
        closeOnEscape={!createMutation.isPending}
        contentClassName="space-y-6"
      >
        <form.FormRoot className="space-y-6">
          <form.FormError />

          {/* Customer Selection */}
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <User size={18} className="text-base-content/70" />
              <h3 className="font-semibold">{tOrders('customer')}</h3>
            </div>

            <form.AppField
              name="customerId"
              validators={{
                onBlur: z.string().min(1, 'validation.required'),
              }}
            >
              {(field) => (
                <field.CustomerSelectField
                  label={tOrders('select_customer')}
                  businessDescriptor={businessDescriptor}
                  placeholder={tOrders('search_customer')}
                  required
                />
              )}
            </form.AppField>
          </div>

          {/* Address Selection */}
          {selectedCustomerId && (
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <MapPin size={18} className="text-base-content/70" />
                <h3 className="font-semibold">{tOrders('shipping_address')}</h3>
              </div>

              <form.AppField
                name="shippingAddressId"
                validators={{
                  onBlur: z.string().min(1, 'validation.required'),
                }}
              >
                {(field) => (
                  <field.AddressSelectField
                    label={tOrders('select_address')}
                    businessDescriptor={businessDescriptor}
                    customerId={selectedCustomerId}
                    placeholder={tOrders('search_address')}
                    required
                  />
                )}
              </form.AppField>
            </div>
          )}

          {/* Channel */}
          <form.AppField
            name="channel"
            validators={{
              onBlur: z.string().min(1, 'validation.required'),
            }}
          >
            {(field) => (
              <field.SelectField
                label={tOrders('channel')}
                options={CHANNELS.map((ch) => ({
                  value: ch.value,
                  label: tOrders(ch.labelKey),
                }))}
                required
              />
            )}
          </form.AppField>

          {/* Items */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Package size={18} className="text-base-content/70" />
                <h3 className="font-semibold">{tOrders('items')}</h3>
              </div>
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={addItem}
              >
                <Plus size={16} />
                {tOrders('add_item')}
              </button>
            </div>

            <form.Subscribe selector={(state) => state.values.items}>
              {(items) =>
                items.map((_, index: number) => (
                  <div
                    key={index}
                    className="rounded-lg bg-base-200 p-4 space-y-3"
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium">
                        {tOrders('item')} {index + 1}
                      </span>
                      {items.length > 1 && (
                        <button
                          type="button"
                          className="btn btn-ghost btn-sm btn-circle"
                          onClick={() => removeItem(index)}
                        >
                          <Trash2 size={16} />
                        </button>
                      )}
                    </div>

                    <form.AppField
                      name={`items[${index}].variantId`}
                      validators={{
                        onBlur: z.string().min(1, 'validation.required'),
                      }}
                    >
                      {(field) => (
                        <field.ProductVariantSelectField
                          label={tOrders('product')}
                          businessDescriptor={businessDescriptor}
                          placeholder={tOrders('search_product')}
                          required
                          showDetails={false}
                          onVariantSelect={(variant: {
                            id: string
                            productId: string
                            productName: string
                            variantName: string
                            salePrice: string
                            costPrice: string
                            stockQuantity: number
                          }) => {
                            // Auto-fill prices and update max quantity
                            form.setFieldValue(
                              `items[${index}].unitPrice`,
                              variant.salePrice,
                            )
                            form.setFieldValue(
                              `items[${index}].unitCost`,
                              variant.costPrice,
                            )
                            // Store stock quantity for reactive max validation
                            form.setFieldMeta(
                              `items[${index}].variantId`,
                              (prev) => ({
                                ...prev,
                                stockQuantity: variant.stockQuantity,
                              }),
                            )

                            // Clamp quantity if exceeds new stock
                            const currentQuantity =
                              form.getFieldValue(`items[${index}].quantity`) ||
                              1
                            if (currentQuantity > variant.stockQuantity) {
                              form.setFieldValue(
                                `items[${index}].quantity`,
                                variant.stockQuantity,
                              )
                            }
                          }}
                        />
                      )}
                    </form.AppField>

                    <form.AppField
                      name={`items[${index}].quantity`}
                      validators={{
                        onBlur: z.coerce
                          .number()
                          .int()
                          .min(1, 'validation.min_value')
                          .max(10000, 'validation.max_value'),
                      }}
                    >
                      {(field) => (
                        <form.Subscribe
                          selector={(state) => {
                            const variantField = state.fieldMeta[
                              `items[${index}].variantId`
                            ] as any
                            return variantField?.stockQuantity
                          }}
                        >
                          {(stockQuantity) => {
                            const maxQuantity = stockQuantity || 10000

                            return (
                              <field.QuantityField
                                label={tOrders('quantity')}
                                min={1}
                                max={maxQuantity}
                                required
                                helperText={
                                  stockQuantity !== undefined
                                    ? tOrders('available_stock', {
                                        count: stockQuantity,
                                      })
                                    : undefined
                                }
                              />
                            )
                          }}
                        </form.Subscribe>
                      )}
                    </form.AppField>
                  </div>
                ))
              }
            </form.Subscribe>
          </div>

          {/* Shipping Zone */}
          {shippingZones.length > 0 && (
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
          <div className="space-y-3">
            <h3 className="font-semibold">{tOrders('discount')}</h3>
            <div className="grid grid-cols-2 gap-3">
              <form.AppField name="discountType">
                {(field) => (
                  <field.SelectField
                    label={tOrders('discount_type')}
                    options={[
                      { value: 'amount', label: tOrders('discount_amount') },
                      { value: 'percent', label: tOrders('discount_percent') },
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

          {/* Payment Details */}
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <CreditCard size={18} className="text-base-content/70" />
              <h3 className="font-semibold">{tOrders('payment_details')}</h3>
            </div>

            {enabledPaymentMethods.length > 0 && (
              <form.AppField name="paymentMethod">
                {(field) => (
                  <field.SelectField
                    label={tOrders('payment_method')}
                    options={[
                      { value: '', label: tCommon('none') },
                      ...enabledPaymentMethods.map((pm) => ({
                        value: pm.descriptor,
                        label: pm.name,
                      })),
                    ]}
                    clearable
                  />
                )}
              </form.AppField>
            )}

            <form.AppField name="paymentReference">
              {(field) => (
                <field.TextField
                  label={tOrders('payment_reference')}
                  type="text"
                  placeholder={tOrders('payment_reference_placeholder')}
                />
              )}
            </form.AppField>
          </div>

          {/* Status & Payment Status */}
          <div className="grid grid-cols-2 gap-3">
            <form.AppField name="status">
              {(field) => (
                <field.SelectField
                  label={tOrders('order_status')}
                  options={[
                    { value: '', label: tOrders('status_pending') },
                    ...ORDER_STATUSES.filter((s) => s !== 'pending').map(
                      (status) => ({
                        value: status,
                        label: tOrders(`status_${status}`),
                      }),
                    ),
                  ]}
                  clearable
                />
              )}
            </form.AppField>

            <form.AppField name="paymentStatus">
              {(field) => (
                <field.SelectField
                  label={tOrders('payment_status_label')}
                  options={[
                    { value: '', label: tOrders('payment_status_pending') },
                    ...PAYMENT_STATUSES.filter((s) => s !== 'pending').map(
                      (status) => ({
                        value: status,
                        label: tOrders(`payment_status_${status}`),
                      }),
                    ),
                  ]}
                  clearable
                />
              )}
            </form.AppField>
          </div>

          {/* Note */}
          <form.AppField name="note">
            {(field) => (
              <field.TextareaField
                label={tOrders('note')}
                placeholder={tOrders('note_placeholder')}
                rows={3}
              />
            )}
          </form.AppField>
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
