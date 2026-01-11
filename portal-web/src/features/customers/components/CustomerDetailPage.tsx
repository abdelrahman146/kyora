import { Link, useNavigate, useParams } from '@tanstack/react-router'
import { useQueryClient } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
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
import { AddressSheet } from './AddressSheet'
import { CustomerDetailSkeleton } from './CustomerDetailSkeleton'
import { EditCustomerSheet } from './EditCustomerSheet'
import { AddressCard } from './AddressCard'
import type { QueryClient } from '@tanstack/react-query'

import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress, CustomerGender } from '@/api/customer'
import type { Order } from '@/api/order'
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
import { Dialog } from '@/components/molecules/Dialog'
import { StatCard } from '@/components/atoms/StatCard'
import { SocialMediaHandles } from '@/components/molecules/SocialMediaHandles'
import { Notes } from '@/components/organisms/Notes'
import { queryKeys } from '@/lib/queryKeys'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'
import { showErrorFromException, showSuccessToast } from '@/lib/toast'
import { getSelectedBusiness } from '@/stores/businessStore'

export async function customerDetailLoader({
  queryClient,
  businessDescriptor,
  customerId,
}: {
  queryClient: QueryClient
  businessDescriptor: string
  customerId: string
}) {
  await queryClient.ensureQueryData(
    customerQueries.detail(businessDescriptor, customerId),
  )

  void queryClient.prefetchQuery(
    addressQueries.list(businessDescriptor, customerId),
  )
  void queryClient.prefetchQuery(metadataQueries.countries())
}

export function CustomerDetailPage() {
  const { i18n } = useTranslation()
  const { t: tCommon } = useTranslation('common')
  const { t: tCustomers } = useTranslation('customers')
  const { t: tDashboard } = useTranslation('dashboard')
  const { t: tErrors } = useTranslation('errors')
  const { t: tOrders } = useTranslation('orders')
  const isArabic = i18n.language.toLowerCase().startsWith('ar')

  const { businessDescriptor, customerId } = useParams({
    from: '/business/$businessDescriptor/customers/$customerId',
  })
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const selectedBusiness = getSelectedBusiness()
  const currency = selectedBusiness?.currency ?? 'AED'

  const { data: countries = [] } = useCountriesQuery()

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

  const [recentOrders, setRecentOrders] = useState<Array<Order>>([])
  const [isLoadingOrders, setIsLoadingOrders] = useState(true)

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

      showSuccessToast(tCustomers('delete_success'))

      void navigate({
        to: '/business/$businessDescriptor/customers',
        params: { businessDescriptor },
        search: { page: 1, pageSize: 20, sortOrder: 'desc' },
      })
    },
    onError: (err) => {
      void showErrorFromException(err, tCustomers)
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
        showSuccessToast(tCommon('notes.note_added'))
      },
      onError: (err) => {
        void showErrorFromException(err, tCustomers)
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
        showSuccessToast(tCommon('notes.note_deleted'))
      },
      onError: (err) => {
        void showErrorFromException(err, tCustomers)
      },
    },
  )

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
        void showErrorFromException(err, tCustomers)
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
  }, [businessDescriptor, customerId, i18n.language, tCustomers])

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
      } catch {
        if (!mounted) return
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
      male: tCustomers('form.gender_male'),
      female: tCustomers('form.gender_female'),
      other: tCustomers('form.gender_other'),
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
      showSuccessToast(tCustomers('address.delete_success'))
    } catch (err) {
      void showErrorFromException(err, tCustomers)
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
        <p className="text-error">{tErrors('generic.load_failed')}</p>
        <button className="btn btn-sm" onClick={() => window.location.reload()}>
          {tCommon('retry')}
        </button>
      </div>
    )
  }

  return (
    <>
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <button
            type="button"
            className="btn btn-ghost btn-sm gap-2"
            onClick={handleBack}
            aria-label={tCommon('back')}
          >
            <ArrowLeft size={18} className={isArabic ? 'rotate-180' : ''} />
            <span className="hidden sm:inline">{tCommon('back')}</span>
          </button>

          <h1 className="text-2xl font-bold flex-1 truncate">
            {customer.name}
          </h1>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <StatCard
            label={tCustomers('orders_count')}
            value={customer.ordersCount ?? 0}
            icon={<ShoppingBag size={24} className="text-success" />}
            variant="success"
          />
          <StatCard
            label={tCustomers('total_spent')}
            value={formatCurrency(customer.totalSpent ?? 0, currency)}
            icon={<ShoppingBag size={24} className="text-primary" />}
            variant="default"
          />
        </div>

        <div className="card bg-base-100 border border-base-300">
          <div className="card-body p-4">
            <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide mb-3">
              {tDashboard('quick_actions')}
            </h3>
            <div className="flex flex-col sm:flex-row gap-2">
              <button
                type="button"
                className="btn btn-primary btn-sm sm:flex-1 gap-2"
                onClick={() => setIsEditOpen(true)}
              >
                <Edit size={16} />
                <span>{tCommon('edit')}</span>
              </button>
              {customer.whatsappNumber && (
                <button
                  type="button"
                  className="btn btn-success btn-outline btn-sm sm:flex-1 gap-2"
                  onClick={handleWhatsAppClick}
                >
                  <FaWhatsapp size={16} />
                  <span>{tCustomers('talk_on_whatsapp')}</span>
                </button>
              )}
              <button
                type="button"
                className="btn btn-error btn-outline btn-sm sm:flex-1 gap-2"
                onClick={() => setIsDeleteDialogOpen(true)}
              >
                <Trash2 size={16} />
                <span>{tCommon('delete')}</span>
              </button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="space-y-6">
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body p-4">
                <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                  {tCustomers('details.title')}
                </h3>

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
                        {tCustomers('details.joined')}{' '}
                        {formatDateShort(customer.joinedAt)}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="space-y-4">
                  {customer.email && (
                    <div className="flex items-start gap-3">
                      <Mail
                        size={18}
                        className="text-base-content/40 mt-0.5 flex-shrink-0"
                      />
                      <div className="flex-1 min-w-0">
                        <div className="text-xs text-base-content/60 mb-0.5">
                          {tCustomers('form.email')}
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
                          {tCustomers('phone')}
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
                        {tCustomers('form.gender')}
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
                        {tCustomers('form.country')}
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

                {(customer.instagramUsername ||
                  customer.facebookUsername ||
                  customer.tiktokUsername ||
                  customer.snapchatUsername ||
                  customer.xUsername ||
                  customer.whatsappNumber) && (
                  <div className="border-t border-base-300 pt-4 mt-6">
                    <div className="text-xs text-base-content/60 uppercase tracking-wide font-semibold mb-3">
                      {tCustomers('details.social_media')}
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

            <div className="card bg-base-100 border border-base-300">
              <div className="card-body p-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                    {tCustomers('details.addresses')}
                  </h3>
                  <button
                    type="button"
                    className="btn btn-primary btn-sm gap-2"
                    onClick={handleAddAddress}
                  >
                    <Plus size={16} />
                    <span className="hidden sm:inline">
                      {tCustomers('address.add_button')}
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
                      {tCustomers('details.no_addresses')}
                    </p>
                    <button
                      type="button"
                      className="btn btn-outline btn-sm gap-2"
                      onClick={handleAddAddress}
                    >
                      <Plus size={16} />
                      {tCustomers('address.add_first')}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="space-y-6">
            <Notes
              notes={customer.notes ?? []}
              onAddNote={handleAddNote}
              onDeleteNote={handleDeleteNote}
              isAddingNote={isAddingNote}
              isDeletingNote={isDeletingNote}
              deletingNoteId={deletingNoteId}
            />

            <div className="card bg-base-100 border border-base-300">
              <div className="card-body p-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                    {tCustomers('details.recent_orders')}
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
                        {tCustomers('details.view_all_orders')}
                      </span>
                      <span className="sm:hidden">{tCommon('view')}</span>
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
                      {tCustomers('details.no_orders')}
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

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

      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => setIsDeleteDialogOpen(false)}
        title={tCustomers('delete_confirm_title')}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => setIsDeleteDialogOpen(false)}
              disabled={isDeleting}
            >
              {tCommon('cancel')}
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
              {tCommon('delete')}
            </button>
          </div>
        }
      >
        <p>{tCustomers('delete_confirm_message', { name: customer.name })}</p>
      </Dialog>

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

      <Dialog
        open={addressDeleteDialogOpen}
        onClose={() => {
          setAddressDeleteDialogOpen(false)
          setDeletingAddressId(null)
        }}
        title={tCustomers('address.delete_confirm_title')}
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
              {tCommon('cancel')}
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
              {tCommon('delete')}
            </button>
          </div>
        }
      >
        <p>{tCustomers('address.delete_confirm_message')}</p>
      </Dialog>
    </>
  )
}
