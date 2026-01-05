import { Link, createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import {
  ArrowLeft,
  Calendar,
  Edit,
  Globe,
  Mail,
  MapPin,
  Phone,
  Plus,
  ShoppingBag,
  Trash2,
  User,
} from 'lucide-react'
import { FaWhatsapp } from 'react-icons/fa'

import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress, CustomerGender } from '@/api/customer'
import type { Order } from '@/api/order'
import type { RouterContext } from '@/router'
import { addressApi, addressQueries } from '@/api/address'
import {
  customerQueries,
  useCreateCustomerNoteMutation,
  useCustomerQuery,
  useDeleteCustomerMutation,
  useDeleteCustomerNoteMutation,
} from '@/api/customer'
import { orderApi } from '@/api/order'
import { metadataQueries, useCountriesQuery } from '@/api/metadata'
import { Avatar } from '@/components/atoms/Avatar'
import { StatCard } from '@/components/atoms/StatCard'
import { Dialog } from '@/components/atoms/Dialog'
import { CustomerDetailSkeleton } from '@/components/atoms/skeletons/CustomerDetailSkeleton'
import { AddressCard } from '@/components/molecules/AddressCard'
import { RouteErrorFallback } from '@/components/molecules/RouteErrorFallback'
import { SocialMediaHandles } from '@/components/molecules/SocialMediaHandles'
import { AddressSheet } from '@/components/organisms/customers/AddressSheet'
import { EditCustomerSheet } from '@/components/organisms/customers'
import { Notes } from '@/components/organisms/Notes'
import { queryKeys } from '@/lib/queryKeys'
import { showErrorFromException, showSuccessToast } from '@/lib/toast'
import { getSelectedBusiness } from '@/stores/businessStore'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

/**
 * Customer Detail Route
 *
 * Portal-web parity:
 * - Profile card + stats
 * - Social handles
 * - Basic info + addresses management
 * - Notes section
 * - Edit + Delete actions
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/customers/$customerId',
)({
  staticData: {
    titleKey: 'customers.details_title',
  },

  loader: async ({ context, params }) => {
    const { queryClient } = context as unknown as RouterContext

    // Prefetch customer details (critical data - await)
    await queryClient.ensureQueryData(
      customerQueries.detail(params.businessDescriptor, params.customerId),
    )

    // Prefetch addresses and metadata (non-critical - parallel, no await)
    void queryClient.prefetchQuery(
      addressQueries.list(params.businessDescriptor, params.customerId),
    )
    void queryClient.prefetchQuery(metadataQueries.countries())
  },

  errorComponent: RouteErrorFallback,

  component: () => (
    <Suspense fallback={<CustomerDetailSkeleton />}>
      <CustomerDetailPage />
    </Suspense>
  ),
})

function CustomerDetailPage() {
  const { t, i18n } = useTranslation()
  const { t: tOrders } = useTranslation('orders')
  const isArabic = i18n.language.toLowerCase().startsWith('ar')

  const { businessDescriptor, customerId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const selectedBusiness = getSelectedBusiness()
  const currency = selectedBusiness?.currency ?? 'AED'

  // Ensure countries metadata is available for display
  const { data: countries = [] } = useCountriesQuery()

  // State
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

  const [addresses, setAddresses] = useState<Array<CustomerAddress>>([])
  const [isLoadingAddresses, setIsLoadingAddresses] = useState(true)
  const [isAddressSheetOpen, setIsAddressSheetOpen] = useState(false)
  const [editingAddress, setEditingAddress] = useState<
    CustomerAddress | undefined
  >(undefined)
  const [deletingAddressId, setDeletingAddressId] = useState<string | null>(
    null,
  )
  const [isDeletingAddress, setIsDeletingAddress] = useState(false)
  const [addressDeleteDialogOpen, setAddressDeleteDialogOpen] = useState(false)

  // Recent orders state
  const [recentOrders, setRecentOrders] = useState<Array<Order>>([])
  const [isLoadingOrders, setIsLoadingOrders] = useState(true)

  // Customer notes state
  const [isAddingNote, setIsAddingNote] = useState(false)
  const [isDeletingNote, setIsDeletingNote] = useState(false)
  const [deletingNoteId, setDeletingNoteId] = useState<string | null>(null)

  const {
    data: customer,
    isLoading,
    error,
  } = useCustomerQuery(businessDescriptor, customerId)

  const deleteMutation = useDeleteCustomerMutation(businessDescriptor, {
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.list(businessDescriptor),
      })

      showSuccessToast(t('customers.delete_success'))

      void navigate({
        to: '/business/$businessDescriptor/customers',
        params: { businessDescriptor },
        search: { page: 1, pageSize: 20, sortOrder: 'desc' },
      })
    },
    onError: (err) => {
      void showErrorFromException(err, t)
    },
  })

  const createNoteMutation = useCreateCustomerNoteMutation(
    businessDescriptor,
    customerId,
    {
      onSuccess: () => {
        void queryClient.invalidateQueries({
          queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
        })
        showSuccessToast(t('notes.note_added', { ns: 'common' }))
      },
      onError: (err) => {
        void showErrorFromException(err, t)
      },
    },
  )

  const deleteNoteMutation = useDeleteCustomerNoteMutation(
    businessDescriptor,
    customerId,
    {
      onSuccess: () => {
        void queryClient.invalidateQueries({
          queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
        })
        showSuccessToast(t('notes.note_deleted', { ns: 'common' }))
      },
      onError: (err) => {
        void showErrorFromException(err, t)
      },
    },
  )

  // Load addresses
  useEffect(() => {
    let mounted = true

    const fetchAddresses = async () => {
      try {
        setIsLoadingAddresses(true)
        const data = await addressApi.listAddresses(
          businessDescriptor,
          customerId,
        )
        if (!mounted) return
        setAddresses(data)
      } catch (err) {
        if (!mounted) return
        void showErrorFromException(err, t)
      } finally {
        if (mounted) {
          setIsLoadingAddresses(false)
        }
      }
    }

    void fetchAddresses()

    return () => {
      mounted = false
    }
  }, [businessDescriptor, customerId, i18n.language, t])

  // Load recent orders
  useEffect(() => {
    let mounted = true

    const fetchOrders = async () => {
      try {
        setIsLoadingOrders(true)
        const data = await orderApi.listOrders(businessDescriptor, {
          customerId,
          page: 1,
          pageSize: 5,
          orderBy: ['-orderedAt'],
        })
        if (!mounted) return
        setRecentOrders(data.items)
      } catch (err) {
        if (!mounted) return
        // Silent fail for orders - not critical
      } finally {
        if (mounted) {
          setIsLoadingOrders(false)
        }
      }
    }

    void fetchOrders()

    return () => {
      mounted = false
    }
  }, [businessDescriptor, customerId])

  const getInitials = (name: string): string => {
    return name
      .split(' ')
      .map((w) => w[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  const formatPhone = (): string | null => {
    if (customer?.phoneCode && customer.phoneNumber) {
      return `${customer.phoneCode} ${customer.phoneNumber}`
    }
    return null
  }

  const getCountryInfo = (countryCode: string) => {
    const country = countries.find((c) => c.code === countryCode)
    return {
      name: isArabic
        ? (country?.nameAr ?? countryCode)
        : (country?.name ?? countryCode),
      flag: country?.flag,
    }
  }

  const getGenderLabel = (gender: CustomerGender): string => {
    const genderMap: Record<CustomerGender, string> = {
      male: t('customers.form.gender_male'),
      female: t('customers.form.gender_female'),
      other: t('customers.form.gender_other'),
    }
    return genderMap[gender]
  }

  const handleBack = () => {
    void navigate({
      to: '/business/$businessDescriptor/customers',
      params: { businessDescriptor },
      search: { page: 1, pageSize: 20, sortOrder: 'desc' },
    })
  }

  const handleDeleteCustomer = async () => {
    try {
      setIsDeleting(true)
      await deleteMutation.mutateAsync(customerId)
    } finally {
      setIsDeleting(false)
      setIsDeleteDialogOpen(false)
    }
  }

  const handleAddAddress = () => {
    setEditingAddress(undefined)
    setIsAddressSheetOpen(true)
  }

  const handleEditAddress = (address: CustomerAddress) => {
    setEditingAddress(address)
    setIsAddressSheetOpen(true)
  }

  const handleDeleteAddressClick = (addressId: string) => {
    setDeletingAddressId(addressId)
    setAddressDeleteDialogOpen(true)
  }

  const handleDeleteAddressConfirm = async () => {
    if (!deletingAddressId) return

    try {
      setIsDeletingAddress(true)
      await addressApi.deleteAddress(
        businessDescriptor,
        customerId,
        deletingAddressId,
      )
      setAddresses((prev) => prev.filter((a) => a.id !== deletingAddressId))
      showSuccessToast(t('customers.address.delete_success'))
    } catch (err) {
      void showErrorFromException(err, t)
    } finally {
      setIsDeletingAddress(false)
      setDeletingAddressId(null)
      setAddressDeleteDialogOpen(false)
    }
  }

  const handleAddNote = async (content: string) => {
    try {
      setIsAddingNote(true)
      await createNoteMutation.mutateAsync(content)
    } finally {
      setIsAddingNote(false)
    }
  }

  const handleDeleteNote = async (noteId: string) => {
    try {
      setIsDeletingNote(true)
      setDeletingNoteId(noteId)
      await deleteNoteMutation.mutateAsync(noteId)
    } finally {
      setIsDeletingNote(false)
      setDeletingNoteId(null)
    }
  }

  const handleAddressSubmit = async (
    data: CreateAddressRequest | UpdateAddressRequest,
  ): Promise<CustomerAddress> => {
    if (editingAddress) {
      const updated = await addressApi.updateAddress(
        businessDescriptor,
        customerId,
        editingAddress.id,
        data as UpdateAddressRequest,
      )
      setAddresses((prev) =>
        prev.map((a) => (a.id === updated.id ? updated : a)),
      )
      return updated
    }

    const created = await addressApi.createAddress(
      businessDescriptor,
      customerId,
      data as CreateAddressRequest,
    )
    setAddresses((prev) => [...prev, created])
    return created
  }

  const handleWhatsAppClick = () => {
    if (!customer?.whatsappNumber) return
    const cleanNumber = customer.whatsappNumber.replace(/[^\d+]/g, '')
    window.open(`https://wa.me/${cleanNumber}`, '_blank', 'noopener,noreferrer')
  }

  const getOrderStatusBadgeClass = (status: Order['status']): string => {
    const statusMap: Record<Order['status'], string> = {
      pending: 'badge-warning',
      placed: 'badge-info',
      ready_for_shipment: 'badge-info',
      shipped: 'badge-primary',
      fulfilled: 'badge-success',
      cancelled: 'badge-error',
      returned: 'badge-error',
    }
    return statusMap[status] || 'badge-ghost'
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

  return (
    <>
      <div className="space-y-6">
        {/* Header with Back Button */}
        <div className="flex items-center gap-4">
          <button
            type="button"
            className="btn btn-ghost btn-sm gap-2"
            onClick={handleBack}
            aria-label={t('common.back')}
          >
            <ArrowLeft size={18} className={isArabic ? 'rotate-180' : ''} />
            <span className="hidden sm:inline">{t('common.back')}</span>
          </button>

          <h1 className="text-2xl font-bold flex-1 truncate">
            {customer.name}
          </h1>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <StatCard
            label={t('customers.orders_count')}
            value={customer.ordersCount ?? 0}
            icon={<ShoppingBag size={24} className="text-success" />}
            variant="success"
          />
          <StatCard
            label={t('customers.total_spent')}
            value={formatCurrency(customer.totalSpent ?? 0, currency)}
            icon={<ShoppingBag size={24} className="text-primary" />}
            variant="default"
          />
        </div>

        {/* Quick Actions */}
        <div className="card bg-base-100 border border-base-300">
          <div className="card-body p-4">
            <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide mb-3">
              {t('dashboard.quick_actions')}
            </h3>
            <div className="flex flex-col sm:flex-row gap-2">
              <button
                type="button"
                className="btn btn-primary btn-sm sm:flex-1 gap-2"
                onClick={() => setIsEditOpen(true)}
              >
                <Edit size={16} />
                <span>{t('common.edit')}</span>
              </button>
              {customer.whatsappNumber && (
                <button
                  type="button"
                  className="btn btn-success btn-outline btn-sm sm:flex-1 gap-2"
                  onClick={handleWhatsAppClick}
                >
                  <FaWhatsapp size={16} />
                  <span>{t('customers.talk_on_whatsapp')}</span>
                </button>
              )}
              <button
                type="button"
                className="btn btn-error btn-outline btn-sm sm:flex-1 gap-2"
                onClick={() => setIsDeleteDialogOpen(true)}
              >
                <Trash2 size={16} />
                <span>{t('common.delete')}</span>
              </button>
            </div>
          </div>
        </div>

        {/* Desktop: 2-column layout, Mobile: stacked */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Left Column */}
          <div className="space-y-6">
            {/* Customer Details Card */}
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body p-4">
                <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                  {t('customers.details.title')}
                </h3>

                {/* Avatar + Basic Info */}
                <div className="flex items-start gap-2 mb-6 pb-6 border-b border-base-300">
                  <Avatar
                    src={customer.avatarUrl}
                    alt={customer.name}
                    fallback={getInitials(customer.name)}
                    size="md"
                  />
                  <div className="flex-1 min-w-0">
                    <h2 className="text-xl font-bold truncate mb-1">
                      {customer.name}
                    </h2>
                    <div className="flex items-center gap-1 text-xs text-base-content/60">
                      <Calendar size={14} />
                      <span>
                        {t('customers.details.joined')}{' '}
                        {formatDateShort(customer.joinedAt)}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Contact Details */}
                <div className="space-y-4">
                  {customer.email && (
                    <div className="flex items-start gap-3">
                      <Mail
                        size={18}
                        className="text-base-content/40 mt-0.5 flex-shrink-0"
                      />
                      <div className="flex-1 min-w-0">
                        <div className="text-xs text-base-content/60 mb-0.5">
                          {t('customers.form.email')}
                        </div>
                        <div className="font-medium break-all">
                          {customer.email}
                        </div>
                      </div>
                    </div>
                  )}

                  {formatPhone() && (
                    <div className="flex items-start gap-3">
                      <Phone
                        size={18}
                        className="text-base-content/40 mt-0.5 flex-shrink-0"
                      />
                      <div className="flex-1 min-w-0">
                        <div className="text-xs text-base-content/60 mb-0.5">
                          {t('customers.phone')}
                        </div>
                        <span className="font-medium" dir="ltr">
                          {formatPhone()}
                        </span>
                      </div>
                    </div>
                  )}

                  <div className="flex items-start gap-3">
                    <User
                      size={18}
                      className="text-base-content/40 mt-0.5 flex-shrink-0"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="text-xs text-base-content/60 mb-0.5">
                        {t('customers.form.gender')}
                      </div>
                      <div className="font-medium">
                        {getGenderLabel(customer.gender)}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-start gap-3">
                    <Globe
                      size={18}
                      className="text-base-content/40 mt-0.5 flex-shrink-0"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="text-xs text-base-content/60 mb-0.5">
                        {t('customers.form.country')}
                      </div>
                      <div className="flex items-center gap-2 font-medium">
                        {getCountryInfo(customer.countryCode).flag && (
                          <span className="text-lg leading-none">
                            {getCountryInfo(customer.countryCode).flag}
                          </span>
                        )}
                        <span>{getCountryInfo(customer.countryCode).name}</span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Social Media Handles */}
                {(customer.instagramUsername ||
                  customer.facebookUsername ||
                  customer.tiktokUsername ||
                  customer.snapchatUsername ||
                  customer.xUsername ||
                  customer.whatsappNumber) && (
                  <div className="border-t border-base-300 pt-4 mt-6">
                    <div className="text-xs text-base-content/60 uppercase tracking-wide font-semibold mb-3">
                      {t('customers.details.social_media')}
                    </div>
                    <SocialMediaHandles
                      instagramUsername={customer.instagramUsername}
                      facebookUsername={customer.facebookUsername}
                      tiktokUsername={customer.tiktokUsername}
                      snapchatUsername={customer.snapchatUsername}
                      xUsername={customer.xUsername}
                      whatsappNumber={customer.whatsappNumber}
                      size="md"
                    />
                  </div>
                )}
              </div>
            </div>

            {/* Addresses Card */}
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body p-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                    {t('customers.details.addresses')}
                  </h3>
                  <button
                    type="button"
                    className="btn btn-primary btn-sm gap-2"
                    onClick={handleAddAddress}
                  >
                    <Plus size={16} />
                    <span className="hidden sm:inline">
                      {t('customers.address.add_button')}
                    </span>
                  </button>
                </div>

                {isLoadingAddresses ? (
                  <div className="space-y-3">
                    <div className="skeleton h-32 rounded-xl"></div>
                    <div className="skeleton h-32 rounded-xl"></div>
                  </div>
                ) : addresses.length > 0 ? (
                  <div className="space-y-3">
                    {addresses.map((address) => (
                      <AddressCard
                        key={address.id}
                        address={address}
                        onEdit={() => handleEditAddress(address)}
                        onDelete={() => handleDeleteAddressClick(address.id)}
                        isDeleting={
                          isDeletingAddress && deletingAddressId === address.id
                        }
                      />
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12 text-base-content/60">
                    <div className="size-16 rounded-full bg-base-200 flex items-center justify-center mx-auto mb-4">
                      <MapPin size={32} className="opacity-40" />
                    </div>
                    <p className="font-medium mb-4">
                      {t('customers.details.no_addresses')}
                    </p>
                    <button
                      type="button"
                      className="btn btn-outline btn-sm gap-2"
                      onClick={handleAddAddress}
                    >
                      <Plus size={16} />
                      {t('customers.address.add_first')}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Right Column */}
          <div className="space-y-6">
            {/* Customer Notes - Always show */}
            <Notes
              notes={customer.notes ?? []}
              onAddNote={handleAddNote}
              onDeleteNote={handleDeleteNote}
              isAddingNote={isAddingNote}
              isDeletingNote={isDeletingNote}
              deletingNoteId={deletingNoteId}
            />

            {/* Recent Orders */}
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body p-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                    {t('customers.details.recent_orders')}
                  </h3>
                  {recentOrders.length > 0 && (
                    <Link
                      to="/business/$businessDescriptor/orders"
                      params={{ businessDescriptor }}
                      search={{
                        customerId,
                        page: 1,
                        pageSize: 20,
                      }}
                      className="btn btn-ghost btn-sm gap-2"
                    >
                      <span className="hidden sm:inline">
                        {t('customers.details.view_all_orders')}
                      </span>
                      <span className="sm:hidden">{t('common.view')}</span>
                      <ArrowLeft
                        size={16}
                        className={isArabic ? '' : 'rotate-180'}
                      />
                    </Link>
                  )}
                </div>

                {isLoadingOrders ? (
                  <div className="space-y-3">
                    {Array.from({ length: 3 }).map((_, i) => (
                      <div
                        key={i}
                        className="skeleton h-20 rounded-lg animate-pulse"
                      ></div>
                    ))}
                  </div>
                ) : recentOrders.length > 0 ? (
                  <div className="space-y-2">
                    {recentOrders.map((order) => (
                      <div
                        key={order.id}
                        className="flex flex-col sm:flex-row sm:items-center justify-between p-3 rounded-lg border border-base-300 gap-3"
                      >
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-1 flex-wrap">
                            <span className="font-semibold">
                              #{order.orderNumber}
                            </span>
                            <span
                              className={`badge badge-sm ${getOrderStatusBadgeClass(order.status)}`}
                            >
                              {tOrders(`status_${order.status}`)}
                            </span>
                          </div>
                          <div className="flex items-center gap-2 text-xs text-base-content/60">
                            <Calendar size={12} />
                            <span>{formatDateShort(order.orderedAt)}</span>
                          </div>
                        </div>
                        <div className="text-start sm:text-end flex-shrink-0">
                          <div className="font-bold text-primary text-lg">
                            {formatCurrency(
                              parseFloat(order.total),
                              order.currency,
                            )}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12 text-base-content/60">
                    <div className="size-16 rounded-full bg-base-200 flex items-center justify-center mx-auto mb-4">
                      <ShoppingBag size={32} className="opacity-40" />
                    </div>
                    <p className="font-medium">
                      {t('customers.details.no_orders')}
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Edit Customer Sheet */}
      <EditCustomerSheet
        isOpen={isEditOpen}
        onClose={() => setIsEditOpen(false)}
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
          setIsEditOpen(false)
        }}
      />

      {/* Delete Customer Dialog */}
      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => setIsDeleteDialogOpen(false)}
        title={t('customers.delete_confirm_title')}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => setIsDeleteDialogOpen(false)}
              disabled={isDeleting}
            >
              {t('common.cancel')}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={() => void handleDeleteCustomer()}
              disabled={isDeleting}
            >
              {isDeleting && (
                <span className="loading loading-spinner loading-sm" />
              )}
              {t('common.delete')}
            </button>
          </div>
        }
      >
        <p>{t('customers.delete_confirm_message', { name: customer.name })}</p>
      </Dialog>

      {/* Address Sheet */}
      <AddressSheet
        isOpen={isAddressSheetOpen}
        onClose={() => {
          setIsAddressSheetOpen(false)
          setEditingAddress(undefined)
        }}
        onSubmit={handleAddressSubmit}
        address={editingAddress}
        businessDescriptor={businessDescriptor}
      />

      {/* Delete Address Dialog */}
      <Dialog
        open={addressDeleteDialogOpen}
        onClose={() => {
          setAddressDeleteDialogOpen(false)
          setDeletingAddressId(null)
        }}
        title={t('customers.address.delete_confirm_title')}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => {
                setAddressDeleteDialogOpen(false)
                setDeletingAddressId(null)
              }}
              disabled={isDeletingAddress}
            >
              {t('common.cancel')}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={() => void handleDeleteAddressConfirm()}
              disabled={isDeletingAddress}
            >
              {isDeletingAddress && (
                <span className="loading loading-spinner loading-sm" />
              )}
              {t('common.delete')}
            </button>
          </div>
        }
      >
        <p>{t('customers.address.delete_confirm_message')}</p>
      </Dialog>
    </>
  )
}
