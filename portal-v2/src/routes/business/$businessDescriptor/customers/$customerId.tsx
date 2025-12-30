import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import {
  ArrowLeft,
  Edit2,
  Facebook,
  Instagram,
  Mail,
  MapPin,
  MessageCircle,
  Phone,
  Trash2,
} from 'lucide-react'

import { useCustomerQuery, useDeleteCustomerMutation } from '@/api/customer'
import { Dialog } from '@/components/atoms/Dialog'
import { EditCustomerSheet } from '@/components/organisms/customers'
import { CustomerDetailSkeleton } from '@/components/atoms/skeletons/CustomerDetailSkeleton'
import { queryKeys } from '@/lib/queryKeys'
import { showErrorFromException, showSuccessToast } from '@/lib/toast'

/**
 * Customer Detail Route
 *
 * Displays detailed customer information with:
 * - Customer profile card
 * - Order history
 * - Edit in BottomSheet with TanStack Form
 * - Delete confirmation Dialog
 * - Optimistic updates with toast on rollback only
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/customers/$customerId',
)({
  component: () => (
    <Suspense fallback={<CustomerDetailSkeleton />}>
      <CustomerDetailPage />
    </Suspense>
  ),
})

/**
 * Customer Detail Page Component
 */
function CustomerDetailPage() {
  const { t } = useTranslation()
  const { businessDescriptor, customerId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [showEditSheet, setShowEditSheet] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)

  // Fetch customer data
  const {
    data: customer,
    isLoading,
    error,
  } = useCustomerQuery(businessDescriptor, customerId)

  // Delete mutation
  const deleteMutation = useDeleteCustomerMutation(businessDescriptor, {
    onSuccess: () => {
      // Invalidate customers list
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.list(businessDescriptor),
      })

      showSuccessToast(t('customers.delete_success'))

      // Navigate back to customers list
      void navigate({
        to: '/business/$businessDescriptor/customers',
        params: { businessDescriptor },
        search: { page: 1, limit: 20 },
      })
    },
    onError: (err) => {
      void showErrorFromException(err, t)
    },
  })

  // Handle delete customer
  const handleDelete = () => {
    deleteMutation.mutate(customerId)
  }

  if (isLoading) {
    return <CustomerDetailSkeleton />
  }

  if (error || !customer) {
    return (
      <div className="flex min-h-[400px] flex-col items-center justify-center gap-4">
        <p className="text-error">{t('errors.generic.load_failed')}</p>
        <button className="btn btn-sm" onClick={() => window.location.reload()}>
          {t('common.retry')}
        </button>
      </div>
    )
  }

  // Format phone display
  const phoneDisplay =
    customer.phoneCode && customer.phoneNumber
      ? `${customer.phoneCode} ${customer.phoneNumber}`
      : null

  // Get primary address if exists
  const primaryAddress = customer.addresses?.[0]
  const addressDisplay = primaryAddress
    ? [primaryAddress.street, primaryAddress.city, primaryAddress.state]
        .filter(Boolean)
        .join(', ')
    : null

  return (
    <>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-4">
            <button
              className="btn btn-circle btn-ghost btn-sm"
              onClick={() => {
                void navigate({
                  to: '/business/$businessDescriptor/customers',
                  params: { businessDescriptor },
                  search: { page: 1, limit: 20 },
                })
              }}
            >
              <ArrowLeft size={20} />
            </button>
            <div>
              <h1 className="text-2xl font-bold">{customer.name}</h1>
              <p className="text-sm text-base-content/70">
                {t('customers.details_title')}
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            <button
              className="btn btn-outline btn-sm gap-2"
              onClick={() => setShowEditSheet(true)}
            >
              <Edit2 size={16} />
              {t('common.edit')}
            </button>
            <button
              className="btn btn-error btn-outline btn-sm gap-2"
              onClick={() => setShowDeleteDialog(true)}
            >
              <Trash2 size={16} />
              {t('common.delete')}
            </button>
          </div>
        </div>

        {/* Customer Profile Card */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <div className="flex flex-col md:flex-row gap-6">
              {/* Avatar */}
              <div className="avatar placeholder">
                <div className="w-24 h-24 bg-primary/10 text-primary rounded-full">
                  <span className="text-3xl font-bold">
                    {customer.name.charAt(0).toUpperCase()}
                  </span>
                </div>
              </div>

              {/* Customer Info */}
              <div className="flex-1 space-y-4">
                <div>
                  <h2 className="text-2xl font-bold">{customer.name}</h2>
                  <p className="text-sm text-base-content/70">
                    {t('customers.since', {
                      date: new Date(customer.createdAt).toLocaleDateString(),
                    })}
                  </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {/* Email */}
                  {customer.email && (
                    <div className="flex items-center gap-2">
                      <Mail size={16} className="text-base-content/60" />
                      <span className="text-sm">{customer.email}</span>
                    </div>
                  )}

                  {/* Phone */}
                  {phoneDisplay && (
                    <div className="flex items-center gap-2">
                      <Phone size={16} className="text-base-content/60" />
                      <span className="text-sm">{phoneDisplay}</span>
                    </div>
                  )}

                  {/* Address */}
                  {addressDisplay && (
                    <div className="flex items-center gap-2">
                      <MapPin size={16} className="text-base-content/60" />
                      <span className="text-sm">{addressDisplay}</span>
                    </div>
                  )}

                  {/* WhatsApp */}
                  {customer.whatsappNumber && (
                    <div className="flex items-center gap-2">
                      <MessageCircle
                        size={16}
                        className="text-base-content/60"
                      />
                      <span className="text-sm">{customer.whatsappNumber}</span>
                    </div>
                  )}

                  {/* Instagram */}
                  {customer.instagramUsername && (
                    <div className="flex items-center gap-2">
                      <Instagram size={16} className="text-base-content/60" />
                      <a
                        href={`https://instagram.com/${customer.instagramUsername.replace('@', '')}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-sm link link-primary"
                      >
                        @{customer.instagramUsername.replace('@', '')}
                      </a>
                    </div>
                  )}

                  {/* Facebook */}
                  {customer.facebookUsername && (
                    <div className="flex items-center gap-2">
                      <Facebook size={16} className="text-base-content/60" />
                      <a
                        href={`https://facebook.com/${customer.facebookUsername.replace('@', '')}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-sm link link-primary"
                      >
                        {customer.facebookUsername}
                      </a>
                    </div>
                  )}
                </div>

                {/* Notes */}
                {customer.notes && customer.notes.length > 0 && (
                  <div className="pt-4 border-t border-base-300">
                    <h3 className="text-sm font-semibold mb-2">
                      {t('customers.details.notes')}
                    </h3>
                    <p className="text-sm text-base-content/70 whitespace-pre-wrap">
                      {customer.notes[0].content}
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="text-sm text-base-content/70">
                {t('customers.total_orders')}
              </div>
              <div className="text-3xl font-bold">
                {customer.ordersCount ?? 0}
              </div>
            </div>
          </div>

          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="text-sm text-base-content/70">
                {t('customers.total_spent')}
              </div>
              <div className="text-3xl font-bold">
                {(customer.totalSpent ?? 0).toFixed(2)}
              </div>
            </div>
          </div>

          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="text-sm text-base-content/70">
                {t('customers.average_order')}
              </div>
              <div className="text-3xl font-bold">
                {(customer.ordersCount ?? 0) > 0
                  ? (
                      (customer.totalSpent ?? 0) / (customer.ordersCount ?? 1)
                    ).toFixed(2)
                  : '0.00'}
              </div>
            </div>
          </div>
        </div>

        {/* Recent Orders */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="text-lg font-semibold mb-4">
              {t('customers.recent_orders')}
            </h3>
            <div className="text-center text-base-content/70 py-8">
              {t('customers.no_orders_yet')}
            </div>
          </div>
        </div>
      </div>

      {/* Edit Customer Sheet */}
      <EditCustomerSheet
        isOpen={showEditSheet}
        onClose={() => setShowEditSheet(false)}
        businessDescriptor={businessDescriptor}
        customer={customer}
        onUpdated={() => {
          void queryClient.invalidateQueries({
            queryKey: queryKeys.customers.list(businessDescriptor),
          })
          void queryClient.invalidateQueries({
            queryKey: queryKeys.customers.detail(
              businessDescriptor,
              customerId,
            ),
          })
          setShowEditSheet(false)
        }}
      />

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
        title={t('customers.delete_confirm_title')}
        description={t('customers.delete_confirm_message', {
          name: customer.name,
        })}
      >
        <div className="flex justify-end gap-2 mt-4">
          <button
            className="btn btn-ghost"
            onClick={() => setShowDeleteDialog(false)}
            disabled={deleteMutation.isPending}
          >
            {t('common.cancel')}
          </button>
          <button
            className="btn btn-error"
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? (
              <span className="loading loading-spinner loading-sm"></span>
            ) : (
              t('common.delete')
            )}
          </button>
        </div>
      </Dialog>
    </>
  )
}
