import { useId, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import {
  CreditCard,
  Edit,
  ExternalLink,
  MapPin,
  MessageCircle,
  Trash2,
  Truck,
  Zap,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { EditOrderSheet } from './EditOrderSheet'

import type { Order } from '@/api/order'
import {
  orderQueries,
  useAddOrderPaymentDetailsMutation,
  useDeleteOrderMutation,
  useUpdateOrderMutation,
  useUpdateOrderPaymentStatusMutation,
  useUpdateOrderStatusMutation,
} from '@/api/order'
import { BottomSheet } from '@/components'
import { formatCurrency } from '@/lib/formatCurrency'
import { useKyoraForm } from '@/lib/form'

export interface OrderQuickActionsProps {
  order: Order
  businessDescriptor: string
  onDeleteSuccess?: () => void
}

export function OrderQuickActions({
  order,
  businessDescriptor,
  onDeleteSuccess,
}: OrderQuickActionsProps) {
  const { t: tOrders } = useTranslation('orders')
  const { t: tCommon } = useTranslation('common')
  const queryClient = useQueryClient()
  const formId = useId()
  const [showStatusSheet, setShowStatusSheet] = useState(false)
  const [showPaymentSheet, setShowPaymentSheet] = useState(false)
  const [showPaymentDetailsSheet, setShowPaymentDetailsSheet] = useState(false)
  const [showAddressSheet, setShowAddressSheet] = useState(false)
  const [showEditSheet, setShowEditSheet] = useState(false)
  const [isUpdating, setIsUpdating] = useState(false)

  // State machine: Valid status transitions
  const getValidStatusTransitions = (
    currentStatus: Order['status'],
  ): Array<Order['status']> => {
    const transitions: Record<Order['status'], Array<Order['status']>> = {
      pending: ['placed', 'cancelled'],
      placed: ['ready_for_shipment', 'cancelled'],
      ready_for_shipment: ['shipped', 'cancelled'],
      shipped: ['fulfilled', 'returned'],
      fulfilled: ['returned'],
      cancelled: [],
      returned: [],
    }
    return transitions[currentStatus]
  }

  // State machine: Valid payment status transitions
  const getValidPaymentTransitions = (
    currentPaymentStatus: Order['paymentStatus'],
  ): Array<Order['paymentStatus']> => {
    const transitions: Record<
      Order['paymentStatus'],
      Array<Order['paymentStatus']>
    > = {
      pending: ['paid', 'failed'],
      paid: ['refunded'],
      failed: ['paid'],
      refunded: [],
    }
    return transitions[currentPaymentStatus]
  }

  const validStatusOptions = getValidStatusTransitions(order.status)
  const validPaymentOptions = getValidPaymentTransitions(order.paymentStatus)

  const updateStatusMutation = useUpdateOrderStatusMutation(
    businessDescriptor,
    order.id,
  )
  const updatePaymentStatusMutation = useUpdateOrderPaymentStatusMutation(
    businessDescriptor,
    order.id,
  )
  const addPaymentDetailsMutation = useAddOrderPaymentDetailsMutation(
    businessDescriptor,
    order.id,
  )
  const updateOrderMutation = useUpdateOrderMutation(
    businessDescriptor,
    order.id,
  )
  const deleteOrderMutation = useDeleteOrderMutation(
    businessDescriptor,
    order.id,
  )

  const statusFormId = `order-status-form-${formId}`
  const paymentFormId = `order-payment-form-${formId}`
  const paymentDetailsFormId = `order-payment-details-form-${formId}`
  const addressFormId = `order-address-form-${formId}`

  const statusForm = useKyoraForm({
    defaultValues: {
      status: order.status,
    },
    onSubmit: async ({ value }) => {
      setIsUpdating(true)
      try {
        await updateStatusMutation.mutateAsync({ status: value.status })
        await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        toast.success(tOrders('status_updated'))
        setShowStatusSheet(false)
      } catch {
        // Global QueryClient error handler shows the error toast.
      } finally {
        setIsUpdating(false)
      }
    },
  })

  const paymentForm = useKyoraForm({
    defaultValues: {
      paymentStatus: order.paymentStatus,
      /* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */
      paymentMethod: order.paymentMethod || 'cash_on_delivery',
      paymentReference: order.paymentReference || '',
    },
    onSubmit: async ({ value }) => {
      setIsUpdating(true)
      try {
        await updatePaymentStatusMutation.mutateAsync({
          paymentStatus: value.paymentStatus,
        })
        /* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */
        if (value.paymentStatus === 'paid' && value.paymentMethod) {
          await updateOrderMutation.mutateAsync({
            items:
              order.items?.map((item) => ({
                variantId: item.variantId,
                quantity: item.quantity,
                unitPrice: item.unitPrice,
                unitCost: item.unitCost,
              })) || [],
          })
        }
        await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        toast.success(tOrders('payment_status_updated'))
        setShowPaymentSheet(false)
      } catch {
        // Global QueryClient error handler shows the error toast.
      } finally {
        setIsUpdating(false)
      }
    },
  })

  const paymentDetailsForm = useKyoraForm({
    defaultValues: {
      /* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */
      paymentMethod: order.paymentMethod || 'cash_on_delivery',
      paymentReference: order.paymentReference || '',
    },
    onSubmit: async ({ value }) => {
      setIsUpdating(true)
      try {
        await addPaymentDetailsMutation.mutateAsync({
          paymentMethod: value.paymentMethod,
          paymentReference: value.paymentReference || undefined,
        })
        await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        toast.success(tOrders('payment_details_updated'))
        setShowPaymentDetailsSheet(false)
      } catch {
        // Global QueryClient error handler shows the error toast.
      } finally {
        setIsUpdating(false)
      }
    },
  })

  const addressForm = useKyoraForm({
    defaultValues: {
      shippingAddressId: order.shippingAddressId,
    },
    onSubmit: async ({ value }) => {
      setIsUpdating(true)
      try {
        await updateOrderMutation.mutateAsync({
          shippingAddressId: value.shippingAddressId,
          items:
            order.items?.map((item) => ({
              variantId: item.variantId,
              quantity: item.quantity,
              unitPrice: item.unitPrice,
              unitCost: item.unitCost,
            })) || [],
        })
        await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        toast.success(tOrders('shipping_address_updated'))
        setShowAddressSheet(false)
      } catch {
        // Global QueryClient error handler shows the error toast.
      } finally {
        setIsUpdating(false)
      }
    },
  })

  const handleDelete = async () => {
    if (
      window.confirm(
        tOrders('delete_confirm_message', { orderNumber: order.orderNumber }),
      )
    ) {
      try {
        await deleteOrderMutation.mutateAsync()
        await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        toast.success(tOrders('delete_success'))
        onDeleteSuccess?.()
      } catch {
        // Global QueryClient error handler shows the error toast.
      }
    }
  }

  const getPlatformUrl = (platform: string, handle: string | null) => {
    if (!handle) return null
    const cleanHandle = handle.replace('@', '')
    const urls: Record<string, string> = {
      instagram: `https://instagram.com/${cleanHandle}`,
      tiktok: `https://tiktok.com/@${cleanHandle}`,
      facebook: `https://facebook.com/${cleanHandle}`,
      x: `https://x.com/${cleanHandle}`,
      snapchat: `https://snapchat.com/add/${cleanHandle}`,
    }
    return urls[platform.toLowerCase()]
  }

  const getWhatsAppUrl = (
    phoneCode: string | null,
    phoneNumber: string | null,
  ) => {
    if (!phoneCode || !phoneNumber) return null
    const cleanPhone = phoneNumber.replace(/\D/g, '')
    const cleanCode = phoneCode.replace(/\D/g, '')
    return `https://wa.me/${cleanCode}${cleanPhone}`
  }

  const whatsappUrl = order.customer?.whatsappNumber
    ? getWhatsAppUrl(order.customer.phoneCode, order.customer.whatsappNumber)
    : null

  const platformHandle =
    order.customer?.instagramUsername ||
    order.customer?.tiktokUsername ||
    order.customer?.facebookUsername ||
    order.customer?.xUsername ||
    order.customer?.snapchatUsername

  let platformUrl = null
  let platformName = ''
  if (platformHandle) {
    if (order.customer?.instagramUsername) {
      platformUrl = getPlatformUrl(
        'instagram',
        order.customer.instagramUsername,
      )
      platformName = 'Instagram'
    } else if (order.customer?.tiktokUsername) {
      platformUrl = getPlatformUrl('tiktok', order.customer.tiktokUsername)
      platformName = 'TikTok'
    } else if (order.customer?.facebookUsername) {
      platformUrl = getPlatformUrl('facebook', order.customer.facebookUsername)
      platformName = 'Facebook'
    } else if (order.customer?.xUsername) {
      platformUrl = getPlatformUrl('x', order.customer.xUsername)
      platformName = 'X'
    } else if (order.customer?.snapchatUsername) {
      platformUrl = getPlatformUrl('snapchat', order.customer.snapchatUsername)
      platformName = 'Snapchat'
    }
  }

  return (
    <>
      <div className="dropdown dropdown-end">
        <button
          type="button"
          tabIndex={0}
          role="button"
          className="btn btn-ghost btn-sm btn-square"
          aria-label={tCommon('actionsLabel')}
        >
          <Zap size={18} />
        </button>
        <ul
          tabIndex={0}
          className="dropdown-content menu bg-base-100 rounded-box z-[100] w-64 p-2 border border-base-300 mt-2"
        >
          <li>
            <button
              type="button"
              onClick={() => {
                setShowEditSheet(true)
              }}
            >
              <Edit size={18} />
              {tCommon('edit')}
            </button>
          </li>
          <div className="divider my-1" />
          <li>
            <button
              type="button"
              onClick={() => {
                setShowStatusSheet(true)
              }}
            >
              <Truck size={18} />
              {tOrders('update_status')}
            </button>
          </li>
          <li>
            <button
              type="button"
              onClick={() => {
                setShowPaymentSheet(true)
              }}
            >
              <CreditCard size={18} />
              {tOrders('update_payment')}
            </button>
          </li>
          {/* Only show payment details action when order is not in final state */}
          {!['cancelled', 'returned'].includes(order.status) &&
            !['paid', 'refunded'].includes(order.paymentStatus) && (
              <li>
                <button
                  type="button"
                  onClick={() => {
                    setShowPaymentDetailsSheet(true)
                  }}
                >
                  <CreditCard size={18} />
                  {/* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */}
                  {order.paymentMethod
                    ? tOrders('update_payment_details')
                    : tOrders('add_payment_details')}
                </button>
              </li>
            )}
          <li>
            <button
              type="button"
              onClick={() => {
                setShowAddressSheet(true)
              }}
            >
              <MapPin size={18} />
              {tOrders('update_address')}
            </button>
          </li>
          {whatsappUrl && (
            <li>
              <a href={whatsappUrl} target="_blank" rel="noopener noreferrer">
                <MessageCircle size={18} />
                {tOrders('open_whatsapp')}
              </a>
            </li>
          )}
          {platformUrl && (
            <li>
              <a href={platformUrl} target="_blank" rel="noopener noreferrer">
                <ExternalLink size={18} />
                {tOrders('view_on_platform', { platform: platformName })}
              </a>
            </li>
          )}
          <div className="divider my-1" />
          <li>
            <button
              type="button"
              className="text-error hover:bg-error/10"
              onClick={() => {
                handleDelete()
              }}
            >
              <Trash2 size={18} />
              {tCommon('delete')}
            </button>
          </li>
        </ul>
      </div>

      <statusForm.AppForm>
        <BottomSheet
          isOpen={showStatusSheet}
          onClose={() => setShowStatusSheet(false)}
          title={tOrders('update_status')}
          footer={
            <statusForm.SubmitButton
              form={statusFormId}
              variant="primary"
              disabled={isUpdating || validStatusOptions.length === 0}
              className="w-full"
            >
              {isUpdating ? tCommon('loading') : tCommon('update')}
            </statusForm.SubmitButton>
          }
        >
          <statusForm.FormRoot id={statusFormId} className="space-y-4">
            {/* Current Status Info */}
            <div className="rounded-lg bg-base-200 p-3">
              <div className="text-sm">
                <span className="text-base-content/70">
                  {tOrders('current_status')}:{' '}
                </span>
                <span className="font-medium">
                  {tOrders(`status_${order.status}`)}
                </span>
              </div>
            </div>

            {validStatusOptions.length > 0 ? (
              <statusForm.AppField name="status">
                {(field) => (
                  <field.RadioField
                    label={tOrders('new_status')}
                    options={validStatusOptions.map((status) => ({
                      value: status,
                      label: tOrders(`status_${status}`),
                    }))}
                  />
                )}
              </statusForm.AppField>
            ) : (
              <div className="alert alert-info">
                <span className="text-sm">
                  {tOrders('no_status_transitions_available')}
                </span>
              </div>
            )}
          </statusForm.FormRoot>
        </BottomSheet>
      </statusForm.AppForm>

      <paymentForm.AppForm>
        <BottomSheet
          isOpen={showPaymentSheet}
          onClose={() => setShowPaymentSheet(false)}
          title={tOrders('update_payment')}
          footer={
            <paymentForm.SubmitButton
              form={paymentFormId}
              variant="primary"
              disabled={isUpdating || validPaymentOptions.length === 0}
              className="w-full"
            >
              {isUpdating ? tCommon('loading') : tCommon('update')}
            </paymentForm.SubmitButton>
          }
        >
          <paymentForm.FormRoot id={paymentFormId} className="space-y-4">
            {/* Current Payment Status Info */}
            <div className="rounded-lg bg-base-200 p-3">
              <div className="text-sm">
                <span className="text-base-content/70">
                  {tOrders('current_payment_status')}:{' '}
                </span>
                <span className="font-medium">
                  {tOrders(`payment_status_${order.paymentStatus}`)}
                </span>
              </div>
            </div>

            {validPaymentOptions.length > 0 ? (
              <>
                <paymentForm.AppField name="paymentStatus">
                  {(field) => (
                    <field.RadioField
                      label={tOrders('new_payment_status')}
                      options={validPaymentOptions.map((status) => ({
                        value: status,
                        label: tOrders(`payment_status_${status}`),
                      }))}
                    />
                  )}
                </paymentForm.AppField>

                <paymentForm.Subscribe
                  selector={(state) => state.values.paymentStatus}
                >
                  {(paymentStatus) =>
                    paymentStatus === 'paid' && (
                      <>
                        <paymentForm.AppField name="paymentMethod">
                          {(field) => (
                            <field.SelectField
                              label={tOrders('payment_method')}
                              placeholder={tOrders('select_payment_method')}
                              options={[
                                {
                                  value: 'cash_on_delivery',
                                  label: tOrders('payment_method_cod'),
                                },
                                {
                                  value: 'bank_transfer',
                                  label: tOrders('payment_method_bank'),
                                },
                                {
                                  value: 'credit_card',
                                  label: tOrders('payment_method_card'),
                                },
                                {
                                  value: 'paypal',
                                  label: tOrders('payment_method_paypal'),
                                },
                                {
                                  value: 'tamara',
                                  label: tOrders('payment_method_tamara'),
                                },
                                {
                                  value: 'tabby',
                                  label: tOrders('payment_method_tabby'),
                                },
                              ]}
                            />
                          )}
                        </paymentForm.AppField>
                        <paymentForm.AppField name="paymentReference">
                          {(field) => (
                            <field.TextField
                              label={tOrders('payment_reference')}
                              placeholder={tOrders(
                                'payment_reference_placeholder',
                              )}
                            />
                          )}
                        </paymentForm.AppField>
                      </>
                    )
                  }
                </paymentForm.Subscribe>
              </>
            ) : (
              <div className="alert alert-info">
                <span className="text-sm">
                  {tOrders('no_payment_transitions_available')}
                </span>
              </div>
            )}

            <div className="bg-base-200 rounded-lg p-3 space-y-1">
              <div className="flex justify-between text-sm">
                <span className="text-base-content/70">{tOrders('total')}</span>
                <span className="font-semibold">
                  {formatCurrency(parseFloat(order.total), order.currency)}
                </span>
              </div>
            </div>
          </paymentForm.FormRoot>
        </BottomSheet>
      </paymentForm.AppForm>

      <paymentDetailsForm.AppForm>
        <BottomSheet
          isOpen={showPaymentDetailsSheet}
          onClose={() => setShowPaymentDetailsSheet(false)}
          title={
            /* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */
            order.paymentMethod
              ? tOrders('update_payment_details')
              : tOrders('add_payment_details')
          }
          footer={
            <paymentDetailsForm.SubmitButton
              form={paymentDetailsFormId}
              variant="primary"
              disabled={isUpdating}
              className="w-full"
            >
              {isUpdating ? tCommon('loading') : tCommon('update')}
            </paymentDetailsForm.SubmitButton>
          }
        >
          <paymentDetailsForm.FormRoot
            id={paymentDetailsFormId}
            className="space-y-4"
          >
            <div className="rounded-lg bg-base-200 p-3">
              <div className="text-sm">
                <span className="text-base-content/70">
                  {tOrders('current_payment_details')}:{' '}
                </span>
                <span className="font-medium">
                  {/* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */}
                  {order.paymentMethod
                    ? tOrders(`payment_method_${order.paymentMethod}`)
                    : tOrders('no_payment_method')}
                </span>
                {order.paymentReference && (
                  <>
                    <br />
                    <span className="text-base-content/70">
                      {tOrders('payment_reference')}:{' '}
                    </span>
                    <span className="font-medium">
                      {order.paymentReference}
                    </span>
                  </>
                )}
              </div>
            </div>

            <paymentDetailsForm.AppField name="paymentMethod">
              {(field) => (
                <field.SelectField
                  label={tOrders('payment_method')}
                  placeholder={tOrders('select_payment_method')}
                  options={[
                    {
                      value: 'cash_on_delivery',
                      label: tOrders('payment_method_cod'),
                    },
                    {
                      value: 'bank_transfer',
                      label: tOrders('payment_method_bank'),
                    },
                    {
                      value: 'credit_card',
                      label: tOrders('payment_method_card'),
                    },
                    {
                      value: 'paypal',
                      label: tOrders('payment_method_paypal'),
                    },
                    {
                      value: 'tamara',
                      label: tOrders('payment_method_tamara'),
                    },
                    {
                      value: 'tabby',
                      label: tOrders('payment_method_tabby'),
                    },
                  ]}
                />
              )}
            </paymentDetailsForm.AppField>

            <paymentDetailsForm.AppField name="paymentReference">
              {(field) => (
                <field.TextField
                  label={tOrders('payment_reference')}
                  placeholder={tOrders('payment_reference_placeholder')}
                />
              )}
            </paymentDetailsForm.AppField>
          </paymentDetailsForm.FormRoot>
        </BottomSheet>
      </paymentDetailsForm.AppForm>

      <addressForm.AppForm>
        <BottomSheet
          isOpen={showAddressSheet}
          onClose={() => setShowAddressSheet(false)}
          title={tOrders('update_address')}
          footer={
            <addressForm.SubmitButton
              form={addressFormId}
              variant="primary"
              disabled={isUpdating}
              className="w-full"
            >
              {isUpdating ? tCommon('loading') : tCommon('update')}
            </addressForm.SubmitButton>
          }
        >
          <addressForm.FormRoot id={addressFormId} className="space-y-4">
            <addressForm.AppField name="shippingAddressId">
              {(field) => (
                <field.AddressSelectField
                  label={tOrders('shipping_address')}
                  businessDescriptor={businessDescriptor}
                  customerId={order.customerId}
                  placeholder={tOrders('select_address')}
                />
              )}
            </addressForm.AppField>
          </addressForm.FormRoot>
        </BottomSheet>
      </addressForm.AppForm>

      <EditOrderSheet
        isOpen={showEditSheet}
        onClose={() => setShowEditSheet(false)}
        order={order}
        onUpdated={async () => {
          await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        }}
      />
    </>
  )
}
