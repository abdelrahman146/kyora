import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import {
  CreditCard,
  ExternalLink,
  MapPin,
  MessageCircle,
  Trash2,
  Truck,
  Zap,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { BottomSheet } from './BottomSheet'
import type { Order } from '@/api/order'
import { orderApi, orderQueries } from '@/api/order'
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
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [showStatusSheet, setShowStatusSheet] = useState(false)
  const [showPaymentSheet, setShowPaymentSheet] = useState(false)
  const [showAddressSheet, setShowAddressSheet] = useState(false)
  const [isUpdating, setIsUpdating] = useState(false)

  const statusForm = useKyoraForm({
    defaultValues: {
      status: order.status,
    },
    onSubmit: async ({ value }) => {
      setIsUpdating(true)
      try {
        await orderApi.updateOrderStatus(businessDescriptor, order.id, {
          status: value.status,
        })
        await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        toast.success(t('orders:status_updated'))
        setShowStatusSheet(false)
      } catch {
        toast.error(t('common:error_occurred'))
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
        await orderApi.updateOrderPaymentStatus(businessDescriptor, order.id, {
          paymentStatus: value.paymentStatus,
        })
        /* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */
        if (value.paymentStatus === 'paid' && value.paymentMethod) {
          await orderApi.updateOrder(businessDescriptor, order.id, {
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
        toast.success(t('orders:payment_status_updated'))
        setShowPaymentSheet(false)
      } catch {
        toast.error(t('common:error_occurred'))
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
        await orderApi.updateOrder(businessDescriptor, order.id, {
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
        toast.success(t('orders:shipping_address_updated'))
        setShowAddressSheet(false)
      } catch {
        toast.error(t('common:error_occurred'))
      } finally {
        setIsUpdating(false)
      }
    },
  })

  const handleDelete = async () => {
    if (
      window.confirm(
        t('orders:delete_confirm_message', { orderNumber: order.orderNumber }),
      )
    ) {
      try {
        await orderApi.deleteOrder(businessDescriptor, order.id)
        await queryClient.invalidateQueries({ queryKey: orderQueries.all })
        toast.success(t('orders:delete_success'))
        onDeleteSuccess?.()
      } catch {
        toast.error(t('common:error_occurred'))
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
          aria-label={t('actions')}
        >
          <Zap size={18} />
        </button>
        <ul
          tabIndex={0}
          className="dropdown-content menu bg-base-100 rounded-box z-[100] w-64 p-2 shadow-lg border border-base-300 mt-2"
        >
          <li>
            <button
              type="button"
              onClick={() => {
                setShowStatusSheet(true)
              }}
            >
              <Truck size={18} />
              {t('orders:update_status')}
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
              {t('orders:update_payment')}
            </button>
          </li>
          <li>
            <button
              type="button"
              onClick={() => {
                setShowAddressSheet(true)
              }}
            >
              <MapPin size={18} />
              {t('orders:update_address')}
            </button>
          </li>
          {whatsappUrl && (
            <li>
              <a href={whatsappUrl} target="_blank" rel="noopener noreferrer">
                <MessageCircle size={18} />
                {t('orders:open_whatsapp')}
              </a>
            </li>
          )}
          {platformUrl && (
            <li>
              <a href={platformUrl} target="_blank" rel="noopener noreferrer">
                <ExternalLink size={18} />
                {t('orders:view_on_platform', { platform: platformName })}
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
              {t('common:delete')}
            </button>
          </li>
        </ul>
      </div>

      <BottomSheet
        isOpen={showStatusSheet}
        onClose={() => setShowStatusSheet(false)}
        title={t('orders:update_status')}
        footer={
          <statusForm.SubmitButton
            variant="primary"
            disabled={isUpdating}
            className="w-full"
          >
            {isUpdating ? t('common:loading') : t('common:update')}
          </statusForm.SubmitButton>
        }
      >
        <statusForm.AppForm>
          <statusForm.FormRoot className="space-y-4">
            <statusForm.AppField name="status">
              {(field) => (
                <field.RadioField
                  label={t('orders:status')}
                  options={[
                    { value: 'pending', label: t('orders:status_pending') },
                    { value: 'placed', label: t('orders:status_placed') },
                    {
                      value: 'ready_for_shipment',
                      label: t('orders:status_ready_for_shipment'),
                    },
                    { value: 'shipped', label: t('orders:status_shipped') },
                    { value: 'fulfilled', label: t('orders:status_fulfilled') },
                    { value: 'cancelled', label: t('orders:status_cancelled') },
                    { value: 'returned', label: t('orders:status_returned') },
                  ]}
                />
              )}
            </statusForm.AppField>
          </statusForm.FormRoot>
        </statusForm.AppForm>
      </BottomSheet>

      <BottomSheet
        isOpen={showPaymentSheet}
        onClose={() => setShowPaymentSheet(false)}
        title={t('orders:update_payment')}
        footer={
          <paymentForm.SubmitButton
            variant="primary"
            disabled={isUpdating}
            className="w-full"
          >
            {isUpdating ? t('common:loading') : t('common:update')}
          </paymentForm.SubmitButton>
        }
      >
        <paymentForm.AppForm>
          <paymentForm.FormRoot className="space-y-4">
            <paymentForm.AppField name="paymentStatus">
              {(field) => (
                <field.RadioField
                  label={t('orders:payment_status')}
                  options={[
                    {
                      value: 'pending',
                      label: t('orders:payment_status_pending'),
                    },
                    { value: 'paid', label: t('orders:payment_status_paid') },
                    {
                      value: 'failed',
                      label: t('orders:payment_status_failed'),
                    },
                    {
                      value: 'refunded',
                      label: t('orders:payment_status_refunded'),
                    },
                  ]}
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
                          label={t('orders:payment_method')}
                          placeholder={t('orders:select_payment_method')}
                          options={[
                            {
                              value: 'cash_on_delivery',
                              label: t('orders:payment_method_cod'),
                            },
                            {
                              value: 'bank_transfer',
                              label: t('orders:payment_method_bank'),
                            },
                            {
                              value: 'credit_card',
                              label: t('orders:payment_method_card'),
                            },
                            {
                              value: 'paypal',
                              label: t('orders:payment_method_paypal'),
                            },
                            {
                              value: 'tamara',
                              label: t('orders:payment_method_tamara'),
                            },
                            {
                              value: 'tabby',
                              label: t('orders:payment_method_tabby'),
                            },
                          ]}
                        />
                      )}
                    </paymentForm.AppField>
                    <paymentForm.AppField name="paymentReference">
                      {(field) => (
                        <field.TextField
                          label={t('orders:payment_reference')}
                          placeholder={t(
                            'orders:payment_reference_placeholder',
                          )}
                        />
                      )}
                    </paymentForm.AppField>
                  </>
                )
              }
            </paymentForm.Subscribe>

            <div className="bg-base-200 rounded-lg p-3 space-y-1">
              <div className="flex justify-between text-sm">
                <span className="text-base-content/70">
                  {t('orders:total')}
                </span>
                <span className="font-semibold">
                  {formatCurrency(parseFloat(order.total), order.currency)}
                </span>
              </div>
            </div>
          </paymentForm.FormRoot>
        </paymentForm.AppForm>
      </BottomSheet>

      <BottomSheet
        isOpen={showAddressSheet}
        onClose={() => setShowAddressSheet(false)}
        title={t('orders:update_address')}
        footer={
          <addressForm.SubmitButton
            variant="primary"
            disabled={isUpdating}
            className="w-full"
          >
            {isUpdating ? t('common:loading') : t('common:update')}
          </addressForm.SubmitButton>
        }
      >
        <addressForm.AppForm>
          <addressForm.FormRoot className="space-y-4">
            <addressForm.AppField name="shippingAddressId">
              {(field) => (
                <field.AddressSelectField
                  label={t('orders:shipping_address')}
                  businessDescriptor={businessDescriptor}
                  customerId={order.customerId}
                  placeholder={t('orders:select_address')}
                />
              )}
            </addressForm.AppField>
          </addressForm.FormRoot>
        </addressForm.AppForm>
      </BottomSheet>
    </>
  )
}
