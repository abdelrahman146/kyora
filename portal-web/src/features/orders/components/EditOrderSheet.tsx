import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
  CreditCard,
  Lock,
  MapPin,
  Package,
  Plus,
  Trash2,
  User,
} from 'lucide-react'

import { ItemPickerSheet } from './ItemPickerSheet'
import { LiveSummaryCard } from './LiveSummaryCard'
import { OrderPreviewManager } from './OrderPreviewManager'
import { PreviewStatusPill } from './PreviewStatusPill'
import type { OrderItem } from './ItemPickerSheet'
import type { OrderFormValues } from './OrderPreviewManager'
import type { DiscountType, Order, UpdateOrderRequest } from '@/api/order'
import { useUpdateOrderMutation } from '@/api/order'
import { useAddressesQuery } from '@/api/address'
import { usePaymentMethodsQuery, useShippingZonesQuery } from '@/api/business'
import { useVariantsQuery } from '@/api/inventory'
import { BottomSheet } from '@/components'
import { ShippingZoneInfo } from '@/components/organisms/ShippingZoneInfo'
import { StandaloneAddressSheet } from '@/features/customers/components/StandaloneAddressSheet'
import { showSuccessToast } from '@/lib/toast'
import { businessStore, getSelectedBusiness } from '@/stores/businessStore'
import { useKyoraForm } from '@/lib/form'
import { inferShippingZoneFromAddress } from '../utils/shippingZone'

export interface EditOrderSheetProps {
  isOpen: boolean
  onClose: () => void
  order: Order | null
  onUpdated?: () => void | Promise<void>
}

/**
 * Status-aware edit permissions
 * Based on order status, determines which fields can be edited
 */
function getEditPermissions(order: Order) {
  const status = order.status
  const paymentStatus = order.paymentStatus

  // Terminal or locked statuses for items/prices
  const itemsLockedStatuses = ['shipped', 'fulfilled', 'cancelled', 'returned']
  const isStatusLocked = itemsLockedStatuses.includes(status)

  // Address locked after shipped
  const addressLockedStatuses = ['shipped', 'fulfilled', 'returned']
  const isAddressLocked = addressLockedStatuses.includes(status)

  // Payment locked if already completed or in certain states
  const paymentLockedStatuses = ['paid', 'refunded']
  const isPaymentCompleted = paymentLockedStatuses.includes(paymentStatus)

  const isCancelled = status === 'cancelled'

  return {
    canEditItems: !isStatusLocked,
    canEditAddress: !isAddressLocked && !isCancelled,
    canEditDiscount: !isStatusLocked,
    canEditShipping: !isAddressLocked && !isCancelled,
    canEditPayment:
      !isPaymentCompleted && !isCancelled && status !== 'returned',
    isStatusLocked,
    isAddressLocked,
    isPaymentCompleted,
    isCancelled,
  }
}

/**
 * Get status badge color
 */
function getStatusBadgeClass(status: string): string {
  switch (status) {
    case 'fulfilled':
      return 'badge-success'
    case 'cancelled':
    case 'returned':
    case 'failed':
      return 'badge-error'
    case 'shipped':
    case 'ready_for_shipment':
      return 'badge-info'
    case 'placed':
      return 'badge-primary'
    default:
      return 'badge-warning'
  }
}

/**
 * Get payment status badge color
 */
function getPaymentBadgeClass(status: string): string {
  switch (status) {
    case 'paid':
      return 'badge-success'
    case 'failed':
    case 'refunded':
      return 'badge-error'
    default:
      return 'badge-warning'
  }
}

const CHANNELS = [
  { value: 'instagram', labelKey: 'channel_instagram' },
  { value: 'whatsapp', labelKey: 'channel_whatsapp' },
  { value: 'facebook', labelKey: 'channel_facebook' },
  { value: 'tiktok', labelKey: 'channel_tiktok' },
  { value: 'in_person', labelKey: 'channel_in_person' },
  { value: 'other', labelKey: 'channel_other' },
]

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
  const selectedBusiness = getSelectedBusiness()
  const businessCountryCode = selectedBusiness?.countryCode ?? 'SA'

  // Inline sheet states
  const [isInlineAddressOpen, setIsInlineAddressOpen] = useState(false)
  const [itemPickerState, setItemPickerState] = useState<{
    isOpen: boolean
    editIndex: number | null
  }>({ isOpen: false, editIndex: null })

  // Fetch data
  const { data: addressesData = [] } = useAddressesQuery(
    businessDescriptor,
    order?.customer?.id || '',
    isOpen && !!order?.customer?.id,
  )
  const { data: shippingZones = [] } = useShippingZonesQuery(
    businessDescriptor,
    isOpen,
  )
  const { data: paymentMethods = [] } = usePaymentMethodsQuery(
    businessDescriptor,
    isOpen,
  )
  const { data: variantsData } = useVariantsQuery(
    businessDescriptor,
    { page: 1, pageSize: 100 },
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

  // Parse discount from order - use stored discountType/discountValue if available
  const getDiscountType = (): DiscountType => {
    // Use the stored discount type from the order if available
    if (order?.discountType) return order.discountType
    // Fallback for legacy orders without discountType - assume amount
    return 'amount'
  }

  const getDiscountValue = (): string => {
    // Use the stored discount value from the order if available
    if (order?.discountValue !== undefined && order.discountValue !== null) {
      return order.discountValue
    }
    // Fallback for legacy orders - use the calculated discount amount
    if (!order?.discount) return ''
    return order.discount
  }

  // Initialize items from order
  const getInitialItems = (): Array<OrderItem> => {
    if (!order?.items || order.items.length === 0) {
      return [{ variantId: '', quantity: 1, unitPrice: '', unitCost: '' }]
    }
    return order.items.map((item) => ({
      variantId: item.variantId || item.variant?.id || '',
      quantity: item.quantity,
      unitPrice: item.unitPrice,
      unitCost: item.unitCost || '0',
    }))
  }

  const form = useKyoraForm({
    defaultValues: {
      customerId: order?.customer?.id || '',
      shippingAddressId: order?.shippingAddress?.id || '',
      channel: order?.channel || 'instagram',
      shippingZoneId: order?.shippingZone?.id || '',
      discountType: getDiscountType(),
      discountValue: getDiscountValue(),
      paymentMethod: order?.paymentMethod ?? '',
      paymentReference: order?.paymentReference || '',
      status: order?.status ?? '',
      paymentStatus: order?.paymentStatus ?? '',
      note: '',
      items: getInitialItems(),
    },
    onSubmit: async ({ value }) => {
      if (!order) return

      const permissions = getEditPermissions(order)

      const payload: UpdateOrderRequest = {}

      // Only include items if we can edit them
      if (permissions.canEditItems) {
        // Filter out empty placeholder items
        const validItems = value.items.filter(
          (item) => item.variantId && item.variantId.trim() !== '',
        )
        if (validItems.length > 0) {
          payload.items = validItems.map((item) => ({
            variantId: item.variantId,
            quantity: item.quantity,
            unitPrice: item.unitPrice,
            unitCost: item.unitCost,
          }))
        }
      }

      // Address and shipping zone
      if (permissions.canEditAddress || permissions.canEditShipping) {
        if (value.shippingAddressId) {
          payload.shippingAddressId = value.shippingAddressId
        }
        if (value.shippingZoneId) {
          payload.shippingZoneId = value.shippingZoneId
        }
      }

      // Channel
      if (value.channel) {
        payload.channel = value.channel
      }

      // Discount
      if (permissions.canEditDiscount) {
        if (value.discountValue && value.discountValue.trim() !== '') {
          payload.discountType = value.discountType
          payload.discountValue = value.discountValue
        }
      }

      await updateMutation.mutateAsync(payload)
    },
  })

  // Initialize form from order when sheet opens
  useEffect(() => {
    if (order && isOpen) {
      form.setFieldValue('customerId', order.customer?.id || '')
      form.setFieldValue('shippingAddressId', order.shippingAddress?.id || '')
      form.setFieldValue('channel', order.channel || 'instagram')
      form.setFieldValue('shippingZoneId', order.shippingZone?.id || '')
      form.setFieldValue('discountType', getDiscountType())
      form.setFieldValue('discountValue', getDiscountValue())
      form.setFieldValue('paymentMethod', order.paymentMethod)
      form.setFieldValue('paymentReference', order.paymentReference || '')
      form.setFieldValue('status', order.status)
      form.setFieldValue('paymentStatus', order.paymentStatus)

      // Set items
      const initialItems = getInitialItems()
      form.setFieldValue('items', initialItems)
    }
  }, [order, isOpen])

  // Reset form when sheet closes
  useEffect(() => {
    if (!isOpen) {
      form.reset()
      setIsInlineAddressOpen(false)
      setItemPickerState({ isOpen: false, editIndex: null })
    }
  }, [isOpen])

  const safeClose = () => {
    if (updateMutation.isPending) return
    onClose()
  }

  // Item picker handlers
  const openItemPicker = (editIndex: number | null = null) => {
    setItemPickerState({ isOpen: true, editIndex })
  }

  const closeItemPicker = () => {
    setItemPickerState({ isOpen: false, editIndex: null })
  }

  const handleItemSave = (item: OrderItem) => {
    if (itemPickerState.editIndex !== null) {
      // Update existing item
      form.setFieldValue(
        `items[${itemPickerState.editIndex}].variantId`,
        item.variantId,
      )
      form.setFieldValue(
        `items[${itemPickerState.editIndex}].quantity`,
        item.quantity,
      )
      form.setFieldValue(
        `items[${itemPickerState.editIndex}].unitPrice`,
        item.unitPrice,
      )
      form.setFieldValue(
        `items[${itemPickerState.editIndex}].unitCost`,
        item.unitCost,
      )
    } else {
      // Add new item
      const currentItems = form.state.values.items
      // If first item is empty placeholder, replace it
      if (currentItems.length === 1 && currentItems[0].variantId === '') {
        form.setFieldValue('items[0].variantId', item.variantId)
        form.setFieldValue('items[0].quantity', item.quantity)
        form.setFieldValue('items[0].unitPrice', item.unitPrice)
        form.setFieldValue('items[0].unitCost', item.unitCost)
      } else {
        form.pushFieldValue('items', item)
      }
    }
  }

  const handleItemRemove = () => {
    if (itemPickerState.editIndex !== null) {
      form.removeFieldValue('items', itemPickerState.editIndex)
    }
  }

  const removeItem = (index: number) => {
    form.removeFieldValue('items', index)
  }

  // Helper to get variant name for display
  const getVariantName = (variantId: string): string => {
    if (!variantId || !variantsData?.items) return ''
    const variant = variantsData.items.find((v) => v.id === variantId)
    return variant?.name ?? ''
  }

  // Memoized permissions
  const permissions = useMemo(() => {
    if (!order) return null
    return getEditPermissions(order)
  }, [order])

  if (!order || !permissions) return null

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
              const canSubmit = previewCanSubmit && !updateMutation.isPending

              return (
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
                        disabled={!canSubmit}
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

                    {/* Status Header Strip */}
                    <div className="flex items-center justify-between gap-2 rounded-lg bg-base-200 p-3">
                      <div className="flex items-center gap-2">
                        <span
                          className="text-sm text-base-content/70"
                          dir="ltr"
                        >
                          #{order.orderNumber}
                        </span>
                      </div>
                      <div className="flex items-center gap-2">
                        <span
                          className={`badge badge-sm ${getStatusBadgeClass(order.status)}`}
                        >
                          {tOrders(`status_${order.status}`)}
                        </span>
                        <span
                          className={`badge badge-sm ${getPaymentBadgeClass(order.paymentStatus)}`}
                        >
                          {tOrders(`payment_status_${order.paymentStatus}`)}
                        </span>
                      </div>
                    </div>

                    {/* Edit Restriction Notice */}
                    {permissions.isStatusLocked && (
                      <div className="alert alert-warning text-sm">
                        <Lock size={16} />
                        <span>
                          {tOrders('editing_limited')}{' '}
                          <span className="font-medium">
                            {tOrders('editing_limited_status', {
                              status: tOrders(`status_${order.status}`),
                            })}
                          </span>
                        </span>
                      </div>
                    )}

                    {/* Preview Status */}
                    <div className="flex items-start justify-between gap-3">
                      <div className="space-y-1">
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
                        {/* ─────────── Card 1: Customer & Channel ─────────── */}
                        <div className="rounded-xl border border-base-300 bg-base-100 p-4 space-y-4">
                          <div className="flex items-center gap-2">
                            <User size={18} className="text-base-content/70" />
                            <h3 className="font-semibold text-start">
                              {tOrders('customer')}
                            </h3>
                          </div>

                          {/* Customer Info (Read-only in edit) */}
                          {order.customer && (
                            <div className="rounded-lg bg-base-200 p-3">
                              <div className="font-medium">
                                {order.customer.name}
                              </div>
                              {order.customer.email && (
                                <div className="text-sm text-base-content/60">
                                  {order.customer.email}
                                </div>
                              )}
                            </div>
                          )}

                          {/* Address Section */}
                          {order.customer && (
                            <div className="space-y-3">
                              <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                  <MapPin
                                    size={18}
                                    className="text-base-content/70"
                                  />
                                  <h3 className="font-semibold text-start">
                                    {tOrders('shipping_address')}
                                  </h3>
                                  {!permissions.canEditAddress && (
                                    <Lock
                                      size={14}
                                      className="text-base-content/50"
                                    />
                                  )}
                                </div>
                                {permissions.canEditAddress && (
                                  <button
                                    type="button"
                                    className="btn btn-ghost btn-xs gap-1"
                                    onClick={() => setIsInlineAddressOpen(true)}
                                  >
                                    <Plus size={14} />
                                    {tOrders('add_new_address')}
                                  </button>
                                )}
                              </div>

                              {permissions.canEditAddress ? (
                                <form.AppField name="shippingAddressId">
                                  {(field) => (
                                    <field.AddressSelectField
                                      label={tOrders('select_address')}
                                      businessDescriptor={businessDescriptor}
                                      customerId={order.customer!.id}
                                      placeholder={tOrders('search_address')}
                                      onAddressChange={(
                                        addressId: string | null,
                                      ) => {
                                        // Auto-infer shipping zone when address changes
                                        if (
                                          addressId &&
                                          addressesData.length > 0
                                        ) {
                                          const address = addressesData.find(
                                            (a) => a.id === addressId,
                                          )
                                          if (
                                            address &&
                                            shippingZones.length > 0
                                          ) {
                                            const inferredZone =
                                              inferShippingZoneFromAddress(
                                                address,
                                                shippingZones,
                                              )
                                            form.setFieldValue(
                                              'shippingZoneId',
                                              inferredZone?.id || '',
                                            )
                                          } else {
                                            form.setFieldValue(
                                              'shippingZoneId',
                                              '',
                                            )
                                          }
                                        } else {
                                          form.setFieldValue(
                                            'shippingZoneId',
                                            '',
                                          )
                                        }
                                      }}
                                    />
                                  )}
                                </form.AppField>
                              ) : (
                                <div className="rounded-lg bg-base-200 p-3 opacity-60">
                                  {order.shippingAddress ? (
                                    <div>
                                      <div className="font-medium">
                                        {order.shippingAddress.city},{' '}
                                        {order.shippingAddress.state}
                                      </div>
                                      {order.shippingAddress.street && (
                                        <div className="text-sm text-base-content/60">
                                          {order.shippingAddress.street}
                                        </div>
                                      )}
                                    </div>
                                  ) : (
                                    <span className="text-base-content/50">
                                      {tOrders('no_address')}
                                    </span>
                                  )}
                                </div>
                              )}

                              {/* Shipping Zone Info */}
                              {shippingZones.length > 0 && (
                                <form.Subscribe
                                  selector={(state) =>
                                    state.values.shippingZoneId
                                  }
                                >
                                  {(shippingZoneId) => {
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
                          )}

                          {/* Channel */}
                          <form.AppField name="channel">
                            {(field) => (
                              <field.SelectField
                                label={tOrders('channel')}
                                options={CHANNELS.map((ch) => ({
                                  value: ch.value,
                                  label: tOrders(ch.labelKey),
                                }))}
                              />
                            )}
                          </form.AppField>
                        </div>

                        {/* ─────────── Card 2: Items & Discounts ─────────── */}
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
                              {!permissions.canEditItems && (
                                <Lock
                                  size={14}
                                  className="text-base-content/50"
                                />
                              )}
                            </div>
                            {permissions.canEditItems && (
                              <button
                                type="button"
                                className="btn btn-ghost btn-sm"
                                onClick={() => openItemPicker(null)}
                              >
                                <Plus size={16} />
                                {tOrders('add_item')}
                              </button>
                            )}
                          </div>

                          {/* Items Locked Notice */}
                          {!permissions.canEditItems && (
                            <div className="alert alert-info text-sm py-2">
                              <Package size={16} />
                              <span>{tOrders('items_locked_notice')}</span>
                            </div>
                          )}

                          {/* Item Cards */}
                          <form.Subscribe
                            selector={(state) => state.values.items}
                          >
                            {(items) =>
                              items.map((item, index: number) => {
                                const variantName = getVariantName(
                                  item.variantId,
                                )
                                const hasProduct = !!item.variantId

                                return (
                                  <div
                                    key={index}
                                    className={`rounded-lg bg-base-200 p-3 ${permissions.canEditItems ? 'cursor-pointer hover:bg-base-300 transition-colors' : 'opacity-60'}`}
                                    onClick={
                                      permissions.canEditItems
                                        ? () => openItemPicker(index)
                                        : undefined
                                    }
                                    role={
                                      permissions.canEditItems
                                        ? 'button'
                                        : undefined
                                    }
                                    tabIndex={
                                      permissions.canEditItems ? 0 : undefined
                                    }
                                    onKeyDown={
                                      permissions.canEditItems
                                        ? (e) => {
                                            if (
                                              e.key === 'Enter' ||
                                              e.key === ' '
                                            ) {
                                              e.preventDefault()
                                              openItemPicker(index)
                                            }
                                          }
                                        : undefined
                                    }
                                  >
                                    <div className="flex items-center justify-between">
                                      <div className="flex-1 min-w-0">
                                        {hasProduct ? (
                                          <>
                                            <div className="text-sm font-medium text-base-content truncate">
                                              {variantName ||
                                                `Variant ${item.variantId.substring(0, 8)}...`}
                                            </div>
                                            <div className="text-xs text-base-content/60">
                                              {tOrders('quantity')}:{' '}
                                              {item.quantity} × {item.unitPrice}
                                            </div>
                                          </>
                                        ) : (
                                          <div className="text-sm text-base-content/60 italic">
                                            {tOrders('item_picker_add')}
                                          </div>
                                        )}
                                      </div>
                                      {permissions.canEditItems &&
                                        items.length > 1 && (
                                          <button
                                            type="button"
                                            className="btn btn-ghost btn-sm btn-circle ms-2"
                                            onClick={(e) => {
                                              e.stopPropagation()
                                              removeItem(index)
                                            }}
                                          >
                                            <Trash2 size={16} />
                                          </button>
                                        )}
                                    </div>
                                  </div>
                                )
                              })
                            }
                          </form.Subscribe>

                          {/* Discount Section */}
                          {permissions.canEditDiscount ? (
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
                                          placeholder="0.00"
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
                          ) : (
                            <div className="rounded-lg bg-base-200 p-3 opacity-60">
                              <span className="text-sm font-medium">
                                {tOrders('discount')}:{' '}
                              </span>
                              <span className="text-sm">
                                {order.discount || tCommon('none')}
                              </span>
                            </div>
                          )}
                        </div>

                        {/* ─────────── Card 3: Payment Details ─────────── */}
                        <div className="rounded-xl border border-base-300 bg-base-100 p-4 space-y-4">
                          <div className="flex items-center gap-2">
                            <CreditCard
                              size={18}
                              className="text-base-content/70"
                            />
                            <h3 className="font-semibold text-start">
                              {tOrders('payment_details')}
                            </h3>
                            {!permissions.canEditPayment && (
                              <Lock
                                size={14}
                                className="text-base-content/50"
                              />
                            )}
                          </div>

                          {/* Payment Locked Notice */}
                          {!permissions.canEditPayment && (
                            <div className="alert alert-info text-sm py-2">
                              <Lock size={16} />
                              <span>{tOrders('payment_locked_notice')}</span>
                            </div>
                          )}

                          {/* Payment Method (Read-only in edit) */}
                          <div className="rounded-lg bg-base-200 p-3 opacity-60">
                            <span className="text-sm font-medium">
                              {tOrders('payment_method')}:{' '}
                            </span>
                            <span className="text-sm">
                              {enabledPaymentMethods.find(
                                (pm) => pm.descriptor === order.paymentMethod,
                              )?.name ?? order.paymentMethod}
                            </span>
                          </div>

                          {/* Payment Reference (Read-only in edit) */}
                          {order.paymentReference && (
                            <div className="rounded-lg bg-base-200 p-3 opacity-60">
                              <span className="text-sm font-medium">
                                {tOrders('payment_reference')}:{' '}
                              </span>
                              <span className="text-sm">
                                {order.paymentReference}
                              </span>
                            </div>
                          )}
                        </div>
                      </div>

                      {/* Live Summary Sidebar */}
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

      {/* Inline Address Creation Sheet */}
      <form.Subscribe selector={(state) => state.values.customerId}>
        {(customerId) => (
          <StandaloneAddressSheet
            isOpen={isInlineAddressOpen}
            onClose={() => setIsInlineAddressOpen(false)}
            businessDescriptor={businessDescriptor}
            customerId={customerId || order.customer?.id || ''}
            businessCountryCode={businessCountryCode}
            onCreated={(address) => {
              form.setFieldValue('shippingAddressId', address.id)
              // Auto-infer shipping zone
              if (shippingZones.length > 0) {
                const inferredZone = inferShippingZoneFromAddress(
                  address,
                  shippingZones,
                )
                form.setFieldValue('shippingZoneId', inferredZone?.id || '')
              }
            }}
          />
        )}
      </form.Subscribe>

      {/* Item Picker Sheet */}
      <form.Subscribe selector={(state) => state.values.items}>
        {(items) => (
          <ItemPickerSheet
            isOpen={itemPickerState.isOpen}
            onClose={closeItemPicker}
            businessDescriptor={businessDescriptor}
            existingItem={
              itemPickerState.editIndex !== null
                ? items[itemPickerState.editIndex]
                : undefined
            }
            currentItems={items}
            onSave={handleItemSave}
            onRemove={
              itemPickerState.editIndex !== null && items.length > 1
                ? handleItemRemove
                : undefined
            }
          />
        )}
      </form.Subscribe>
    </form.AppForm>
  )
}
