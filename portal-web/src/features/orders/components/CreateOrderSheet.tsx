import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { CreditCard, MapPin, Package, Plus, Trash2, User } from 'lucide-react'

import { PreviewStatusPill } from './PreviewStatusPill'
import { LiveSummaryCard } from './LiveSummaryCard'
import { OrderPreviewManager } from './OrderPreviewManager'
import type { OrderFormValues } from './OrderPreviewManager'
import type {
  CreateOrderRequest,
  DiscountType,
  OrderPaymentMethod,
  OrderPaymentStatus,
  OrderStatus,
} from '@/api/order'
import type { ShippingZone } from '@/api/business'
import { useAddressesQuery } from '@/api/address'
import { useCustomersQuery } from '@/api/customer'
import { useCreateOrderMutation } from '@/api/order'
import { usePaymentMethodsQuery, useShippingZonesQuery } from '@/api/business'
import { BottomSheet } from '@/components'
import { ShippingZoneInfo } from '@/components/organisms/ShippingZoneInfo'
import { useKyoraForm } from '@/lib/form'
import { showSuccessToast } from '@/lib/toast'
import { businessStore } from '@/stores/businessStore'
import { inferShippingZoneFromAddress } from '@/lib/shippingZone'

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

interface AddressSectionContentProps {
  form: any
  businessDescriptor: string
  customerId: string
  shippingZones: Array<ShippingZone>
}

/**
 * Address selection content rendered inside form.Subscribe when customer is selected.
 * Handles address field + shipping zone inference.
 */
function AddressSectionContent({
  form,
  businessDescriptor,
  customerId,
  shippingZones,
}: AddressSectionContentProps) {
  const { t: tOrders } = useTranslation('orders')

  // Fetch addresses for the selected customer
  const { data: addresses = [] } = useAddressesQuery(
    businessDescriptor,
    customerId,
    true,
  )

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2">
        <MapPin size={18} className="text-base-content/70" />
        <h3 className="font-semibold text-start">
          {tOrders('shipping_address')}
        </h3>
      </div>

      <form.AppField
        name="shippingAddressId"
        validators={{
          onBlur: z.string().min(1, 'validation.required'),
        }}
      >
        {(field: {
          AddressSelectField: React.FC<{
            label: string
            businessDescriptor: string
            customerId: string
            placeholder: string
            required: boolean
            onAddressChange: (id: string | null) => void
          }>
        }) => (
          <field.AddressSelectField
            label={tOrders('select_address')}
            businessDescriptor={businessDescriptor}
            customerId={customerId}
            placeholder={tOrders('search_address')}
            required
            onAddressChange={(addressId: string | null) => {
              // Auto-infer shipping zone when address changes
              if (addressId && addresses.length > 0) {
                const address = addresses.find((a) => a.id === addressId)
                if (address && shippingZones.length > 0) {
                  const inferredZone = inferShippingZoneFromAddress(
                    address,
                    shippingZones,
                  )
                  form.setFieldValue('shippingZoneId', inferredZone?.id || '')
                } else {
                  form.setFieldValue('shippingZoneId', '')
                }
              } else {
                form.setFieldValue('shippingZoneId', '')
              }
            }}
          />
        )}
      </form.AppField>

      {/* Show inferred shipping zone */}
      {shippingZones.length > 0 && (
        <form.Subscribe
          selector={(state: { values: OrderFormValues }) =>
            state.values.shippingZoneId
          }
        >
          {(shippingZoneId: string) => {
            if (!shippingZoneId) return null
            const inferredZone = shippingZones.find(
              (zone) => zone.id === shippingZoneId,
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
  )
}

export function CreateOrderSheet({
  isOpen,
  onClose,
  onCreated,
}: CreateOrderSheetProps) {
  const { t: tOrders } = useTranslation('orders')
  const { t: tCommon } = useTranslation('common')

  const businessDescriptor =
    businessStore.state.selectedBusinessDescriptor || ''

  // Fetch lookup data
  useCustomersQuery(
    businessDescriptor,
    {
      page: 1,
      pageSize: 100,
    },
    isOpen,
  )
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

  // Reset form when sheet closes - only react to isOpen changes
  useEffect(() => {
    if (!isOpen) {
      form.reset()
    }
  }, [isOpen])

  const safeClose = () => {
    if (createMutation.isPending) return
    onClose()
  }

  const enabledPaymentMethods = paymentMethods.filter((pm) => pm.enabled)

  return (
    <form.AppForm>
      {/* Subscribe to form values to pass to OrderPreviewManager */}
      <form.Subscribe selector={(state) => state.values}>
        {(formValues) => (
          <OrderPreviewManager
            businessDescriptor={businessDescriptor}
            isOpen={isOpen}
            formValues={formValues as OrderFormValues}
          >
            {({
              previewData,
              isLoading,
              isStale,
              errorMessage,
              previewState,
              lastPreviewAt,
              canSubmit: previewCanSubmit,
              triggerPreview,
            }) => {
              const canSubmit = previewCanSubmit && !createMutation.isPending

              return (
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
                      <form.SubmitButton
                        form="create-order-form"
                        variant="primary"
                        className="flex-1"
                        disabled={!canSubmit}
                      >
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
                  <form.FormRoot id="create-order-form" className="space-y-6">
                    <form.FormError />

                    <div className="flex items-start justify-between gap-3">
                      <div className="space-y-1">
                        <h3 className="text-lg font-semibold text-start">
                          {tOrders('create_order_title')}
                        </h3>
                        <p className="text-sm text-base-content/60">
                          {tOrders('preview_totals_hint')}
                        </p>
                      </div>
                      <PreviewStatusPill
                        state={previewState}
                        lastUpdated={lastPreviewAt}
                        message={errorMessage}
                      />
                    </div>

                    <div className="grid gap-4 lg:grid-cols-[1.6fr,1fr]">
                      <div className="space-y-4">
                        <div className="rounded-xl border border-base-300 bg-base-100 p-4 space-y-4">
                          <div className="flex items-center gap-2">
                            <User size={18} className="text-base-content/70" />
                            <h3 className="font-semibold text-start">
                              {tOrders('customer')}
                            </h3>
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
                                onCustomerChange={() => {
                                  // Reset address and shipping zone when customer changes
                                  form.setFieldValue('shippingAddressId', '')
                                  form.setFieldValue('shippingZoneId', '')
                                }}
                              />
                            )}
                          </form.AppField>

                          {/* Address section - only shows when customer is selected */}
                          <form.Subscribe
                            selector={(state) => state.values.customerId}
                          >
                            {(customerId) =>
                              customerId ? (
                                <AddressSectionContent
                                  form={form}
                                  businessDescriptor={businessDescriptor}
                                  customerId={customerId}
                                  shippingZones={shippingZones}
                                />
                              ) : null
                            }
                          </form.Subscribe>

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
                        </div>

                        <div className="rounded-xl border border-base-300 bg-base-100 p-4 space-y-4">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <Package
                                size={18}
                                className="text-base-content/70"
                              />
                              <h3 className="font-semibold text-start">
                                {tOrders('items')}
                              </h3>
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

                          <form.Subscribe
                            selector={(state) => state.values.items}
                          >
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
                                      onBlur: z
                                        .string()
                                        .min(1, 'validation.required'),
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
                                          form.setFieldValue(
                                            `items[${index}].unitPrice`,
                                            variant.salePrice,
                                          )
                                          form.setFieldValue(
                                            `items[${index}].unitCost`,
                                            variant.costPrice,
                                          )
                                          form.setFieldMeta(
                                            `items[${index}].variantId`,
                                            (prev) => ({
                                              ...prev,
                                              stockQuantity:
                                                variant.stockQuantity,
                                            }),
                                          )

                                          const currentQuantity =
                                            form.getFieldValue(
                                              `items[${index}].quantity`,
                                            ) || 1
                                          if (
                                            currentQuantity >
                                            variant.stockQuantity
                                          ) {
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
                                          const maxQuantity =
                                            stockQuantity || 10000

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

                          <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
                            <form.AppField name="discountType">
                              {(field) => (
                                <field.SelectField
                                  label={tOrders('discount_type')}
                                  options={[
                                    {
                                      value: 'amount',
                                      label: tOrders('discount_amount'),
                                    },
                                    {
                                      value: 'percent',
                                      label: tOrders('discount_percent'),
                                    },
                                  ]}
                                />
                              )}
                            </form.AppField>

                            <form.Subscribe
                              selector={(state) => state.values.discountType}
                            >
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
                                          <span className="text-base-content/70">
                                            %
                                          </span>
                                        }
                                      />
                                    )
                                  }
                                </form.AppField>
                              )}
                            </form.Subscribe>
                          </div>
                        </div>

                        <div className="rounded-xl border border-base-300 bg-base-100 p-4 space-y-4">
                          <div className="flex items-center gap-2">
                            <CreditCard
                              size={18}
                              className="text-base-content/70"
                            />
                            <h3 className="font-semibold text-start">
                              {tOrders('payment_details')}
                            </h3>
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
                                placeholder={tOrders(
                                  'payment_reference_placeholder',
                                )}
                              />
                            )}
                          </form.AppField>

                          <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
                            <form.AppField name="status">
                              {(field) => (
                                <field.SelectField
                                  label={tOrders('order_status')}
                                  options={[
                                    {
                                      value: '',
                                      label: tOrders('status_pending'),
                                    },
                                    ...ORDER_STATUSES.filter(
                                      (s) => s !== 'pending',
                                    ).map((status) => ({
                                      value: status,
                                      label: tOrders(`status_${status}`),
                                    })),
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
                                    {
                                      value: '',
                                      label: tOrders('payment_status_pending'),
                                    },
                                    ...PAYMENT_STATUSES.filter(
                                      (s) => s !== 'pending',
                                    ).map((status) => ({
                                      value: status,
                                      label: tOrders(
                                        `payment_status_${status}`,
                                      ),
                                    })),
                                  ]}
                                  clearable
                                />
                              )}
                            </form.AppField>
                          </div>

                          <form.AppField name="note">
                            {(field) => (
                              <field.TextareaField
                                label={tOrders('note')}
                                placeholder={tOrders('note_placeholder')}
                                rows={3}
                              />
                            )}
                          </form.AppField>
                        </div>
                      </div>

                      <div className="space-y-3">
                        <LiveSummaryCard
                          preview={previewData}
                          isLoading={isLoading}
                          isStale={isStale}
                          errorMessage={errorMessage}
                          onRetry={triggerPreview}
                        />
                      </div>
                    </div>
                  </form.FormRoot>
                </BottomSheet>
              )
            }}
          </OrderPreviewManager>
        )}
      </form.Subscribe>
    </form.AppForm>
  )
}
