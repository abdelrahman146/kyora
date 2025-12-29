import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Edit2, Instagram, Mail, MapPin, MessageCircle, Phone, Trash2 } from 'lucide-react'
import type {UpdateCustomerRequest} from '@/api/customer';
import {
  
  useCustomerQuery,
  useDeleteCustomerMutation,
  useUpdateCustomerMutation
} from '@/api/customer'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Dialog } from '@/components/atoms/Dialog'
import { CustomerForm } from '@/components/organisms/CustomerForm'
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
  const { t } = useTranslation(['common', 'errors'])
  const { businessDescriptor, customerId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [showEditSheet, setShowEditSheet] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)

  // Fetch customer data
  const { data: customer, isLoading, error } = useCustomerQuery(
    businessDescriptor,
    customerId,
  )

  // Update mutation with optimistic updates
  const updateMutation = useUpdateCustomerMutation(businessDescriptor, {
    onMutate: async (variables: { customerId: string; data: UpdateCustomerRequest }) => {
      const { data: updatedData } = variables
      // Cancel outgoing refetches
      await queryClient.cancelQueries({
        queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
      })

      // Snapshot previous value
      const previousCustomer = queryClient.getQueryData(
        queryKeys.customers.detail(businessDescriptor, customerId)
      )

      // Optimistically update cache
      queryClient.setQueryData(
        queryKeys.customers.detail(businessDescriptor, customerId),
        (old: any) => ({ ...old, ...updatedData })
      )

      return { previousCustomer }
    },
    onSuccess: () => {
      // Invalidate and refetch
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.list(businessDescriptor),
      })
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
      })
      
      showSuccessToast(t('common:customer.updated_success'))
      setShowEditSheet(false)
    },
    onError: (err: Error, _variables: { customerId: string; data: UpdateCustomerRequest }, context: unknown) => {
      // Rollback on error
      const ctx = context as { previousCustomer: any } | undefined
      if (ctx?.previousCustomer) {
        queryClient.setQueryData(
          queryKeys.customers.detail(businessDescriptor, customerId),
          ctx.previousCustomer
        )
      }
      void showErrorFromException(err, t)
    },
  })

  // Delete mutation
  const deleteMutation = useDeleteCustomerMutation(businessDescriptor, {
    onSuccess: () => {
      // Invalidate customers list
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.list(businessDescriptor),
      })
      
      showSuccessToast(t('common:customer.deleted_success'))
      
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

  // Handle edit customer form submission
  const handleUpdateCustomer = async (values: UpdateCustomerRequest) => {
    await updateMutation.mutateAsync({ customerId, data: values })
  }

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
        <p className="text-error">{t('errors:generic.load_failed')}</p>
        <button
          className="btn btn-sm"
          onClick={() => window.location.reload()}
        >
          {t('common:actions.retry')}
        </button>
      </div>
    )
  }

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
              <h1 className="text-2xl font-bold">{customer.fullName}</h1>
              <p className="text-sm text-base-content/70">
                {t('common:customer.customer_details')}
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            <button
              className="btn btn-outline btn-sm gap-2"
              onClick={() => setShowEditSheet(true)}
            >
              <Edit2 size={16} />
              {t('common:actions.edit')}
            </button>
            <button
              className="btn btn-error btn-outline btn-sm gap-2"
              onClick={() => setShowDeleteDialog(true)}
            >
              <Trash2 size={16} />
              {t('common:actions.delete')}
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
                    {customer.fullName.charAt(0).toUpperCase()}
                  </span>
                </div>
              </div>

              {/* Customer Info */}
              <div className="flex-1 space-y-4">
                <div>
                  <h2 className="text-2xl font-bold">{customer.fullName}</h2>
                  <p className="text-sm text-base-content/70">
                    {t('common:customer.since', {
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
                  {customer.phonePrefix && customer.phoneNumber && (
                    <div className="flex items-center gap-2">
                      <Phone size={16} className="text-base-content/60" />
                      <span className="text-sm">
                        {customer.phonePrefix} {customer.phoneNumber}
                      </span>
                    </div>
                  )}

                  {/* Address */}
                  {(customer.address || customer.city || customer.country) && (
                    <div className="flex items-center gap-2">
                      <MapPin size={16} className="text-base-content/60" />
                      <span className="text-sm">
                        {[customer.address, customer.city, customer.country]
                          .filter(Boolean)
                          .join(', ')}
                      </span>
                    </div>
                  )}

                  {/* Instagram */}
                  {customer.instagramHandle && (
                    <div className="flex items-center gap-2">
                      <Instagram size={16} className="text-base-content/60" />
                      <a
                        href={`https://instagram.com/${customer.instagramHandle.replace('@', '')}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-sm link link-primary"
                      >
                        {customer.instagramHandle}
                      </a>
                    </div>
                  )}

                  {/* Facebook */}
                  {customer.facebookHandle && (
                    <div className="flex items-center gap-2">
                      <MessageCircle size={16} className="text-base-content/60" />
                      <a
                        href={`https://facebook.com/${customer.facebookHandle.replace('@', '')}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-sm link link-primary"
                      >
                        {customer.facebookHandle}
                      </a>
                    </div>
                  )}
                </div>

                {/* Notes */}
                {customer.notes && (
                  <div className="pt-4 border-t border-base-300">
                    <h3 className="text-sm font-semibold mb-2">
                      {t('common:customer.notes')}
                    </h3>
                    <p className="text-sm text-base-content/70 whitespace-pre-wrap">
                      {customer.notes}
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
                {t('common:customer.total_orders')}
              </div>
              <div className="text-3xl font-bold">{customer.totalOrders}</div>
            </div>
          </div>

          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="text-sm text-base-content/70">
                {t('common:customer.total_spent')}
              </div>
              <div className="text-3xl font-bold">
                {customer.totalSpent.toFixed(2)}
              </div>
            </div>
          </div>

          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="text-sm text-base-content/70">
                {t('common:customer.average_order')}
              </div>
              <div className="text-3xl font-bold">
                {customer.totalOrders > 0
                  ? (customer.totalSpent / customer.totalOrders).toFixed(2)
                  : '0.00'}
              </div>
            </div>
          </div>
        </div>

        {/* Recent Orders */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="text-lg font-semibold mb-4">
              {t('common:customer.recent_orders')}
            </h3>
            <div className="text-center text-base-content/70 py-8">
              {t('common:customer.no_orders_yet')}
            </div>
          </div>
        </div>
      </div>

      {/* Edit Customer BottomSheet */}
      <BottomSheet
        isOpen={showEditSheet}
        onClose={() => setShowEditSheet(false)}
        title={t('common:customer.edit_customer')}
      >
        <CustomerForm
          customer={customer}
          onSubmit={handleUpdateCustomer}
          onCancel={() => setShowEditSheet(false)}
          isSubmitting={updateMutation.isPending}
        />
      </BottomSheet>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
        title={t('common:customer.confirm_delete')}
        description={t('common:customer.confirm_delete_message', {
          name: customer.fullName,
        })}
        footer={
          <>
            <button
              className="btn btn-ghost"
              onClick={() => setShowDeleteDialog(false)}
              disabled={deleteMutation.isPending}
            >
              {t('common:actions.cancel')}
            </button>
            <button
              className="btn btn-error"
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? (
                <span className="loading loading-spinner loading-sm"></span>
              ) : (
                t('common:actions.delete')
              )}
            </button>
          </>
        }
      >
        <p>{t('common:customer.confirm_delete_warning')}</p>
      </Dialog>
    </>
  )
}
