import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import {
  ArrowLeft,
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

import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress, CustomerGender } from '@/api/customer'
import type { RouterContext } from '@/router'
import { addressApi, addressQueries } from '@/api/address'
import {
  customerQueries,
  useCustomerQuery,
  useDeleteCustomerMutation,
} from '@/api/customer'
import { metadataQueries, useCountriesQuery } from '@/api/metadata'
import { Avatar } from '@/components/atoms/Avatar'
import { Dialog } from '@/components/atoms/Dialog'
import { CustomerDetailSkeleton } from '@/components/atoms/skeletons/CustomerDetailSkeleton'
import { AddressCard } from '@/components/molecules/AddressCard'
import { RouteErrorFallback } from '@/components/molecules/RouteErrorFallback'
import { SocialMediaHandles } from '@/components/molecules/SocialMediaHandles'
import { AddressSheet } from '@/components/organisms/customers/AddressSheet'
import { EditCustomerSheet } from '@/components/organisms/customers'
import { queryKeys } from '@/lib/queryKeys'
import { showErrorFromException, showSuccessToast } from '@/lib/toast'
import { getSelectedBusiness } from '@/stores/businessStore'

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
      <div className="space-y-4">
        {/* Header */}
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

          <div className="flex gap-2">
            <button
              type="button"
              className="btn btn-outline btn-sm gap-2"
              onClick={() => {
                setIsEditOpen(true)
              }}
            >
              <Edit size={16} />
              <span className="hidden sm:inline">{t('common.edit')}</span>
            </button>
            <button
              type="button"
              className="btn btn-error btn-outline btn-sm gap-2"
              onClick={() => {
                setIsDeleteDialogOpen(true)
              }}
            >
              <Trash2 size={16} />
              <span className="hidden sm:inline">{t('common.delete')}</span>
            </button>
          </div>
        </div>

        {/* Header Card with Avatar */}
        <div className="card bg-base-100 border border-base-300">
          <div className="card-body">
            <div className="flex flex-col sm:flex-row items-center sm:items-start gap-4">
              <Avatar
                src={customer.avatarUrl}
                alt={customer.name}
                fallback={getInitials(customer.name)}
                size="xl"
              />

              <div className="flex-1 text-center sm:text-start">
                <h2 className="text-2xl font-bold">{customer.name}</h2>

                {customer.email && (
                  <div className="flex items-center gap-2 justify-center sm:justify-start mt-2 text-base-content/70">
                    <Mail size={16} />
                    <span>{customer.email}</span>
                  </div>
                )}

                {formatPhone() && (
                  <div className="flex items-center gap-2 justify-center sm:justify-start mt-1 text-base-content/70">
                    <Phone size={16} />
                    <span dir="ltr">{formatPhone()}</span>
                  </div>
                )}

                {/* Stats */}
                <div className="grid grid-cols-2 gap-4 mt-6">
                  <div className="stat bg-base-200 rounded-box p-4">
                    <div className="stat-figure text-success">
                      <ShoppingBag size={28} />
                    </div>
                    <div className="stat-title text-sm">
                      {t('customers.orders_count')}
                    </div>
                    <div className="text-2xl font-bold text-success">
                      {customer.ordersCount ?? 0}
                    </div>
                  </div>
                  <div className="stat bg-base-200 rounded-box p-4">
                    <div className="stat-title text-sm">
                      {t('customers.total_spent')}
                    </div>
                    <div className="text-2xl font-bold text-primary">
                      {currency}{' '}
                      {(customer.totalSpent ?? 0).toLocaleString(
                        isArabic ? 'ar-AE' : 'en-US',
                        {
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        },
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Social Media Handles */}
            <div className="mt-6">
              <SocialMediaHandles
                instagramUsername={customer.instagramUsername}
                facebookUsername={customer.facebookUsername}
                tiktokUsername={customer.tiktokUsername}
                snapchatUsername={customer.snapchatUsername}
                xUsername={customer.xUsername}
                whatsappNumber={customer.whatsappNumber}
              />
            </div>

            {/* Details Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
              {/* Basic Information */}
              <div className="card bg-base-100 border border-base-300">
                <div className="card-body">
                  <h3 className="card-title text-lg">
                    {t('customers.details.basic_info')}
                  </h3>

                  <div className="space-y-3 mt-4">
                    <div className="flex items-start gap-3">
                      <User size={18} className="text-base-content/60 mt-0.5" />
                      <div className="flex-1">
                        <div className="text-xs text-base-content/60">
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
                        className="text-base-content/60 mt-0.5"
                      />
                      <div className="flex-1">
                        <div className="text-xs text-base-content/60">
                          {t('customers.form.country')}
                        </div>
                        <div className="flex items-center gap-2 font-medium">
                          {getCountryInfo(customer.countryCode).flag && (
                            <span className="text-lg">
                              {getCountryInfo(customer.countryCode).flag}
                            </span>
                          )}
                          <span>
                            {getCountryInfo(customer.countryCode).name}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Addresses */}
              <div className="card bg-base-100 border border-base-300">
                <div className="card-body">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="card-title text-lg">
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
                      <div className="skeleton h-24 rounded-box"></div>
                      <div className="skeleton h-24 rounded-box"></div>
                    </div>
                  ) : addresses.length > 0 ? (
                    <div className="space-y-3">
                      {addresses.map((address) => (
                        <AddressCard
                          key={address.id}
                          address={address}
                          onEdit={() => {
                            handleEditAddress(address)
                          }}
                          onDelete={() => {
                            handleDeleteAddressClick(address.id)
                          }}
                          isDeleting={
                            isDeletingAddress &&
                            deletingAddressId === address.id
                          }
                        />
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8 text-base-content/60">
                      <MapPin size={32} className="mx-auto mb-2 opacity-40" />
                      <p className="mb-3">
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

            {/* Notes Section */}
            {customer.notes && customer.notes.length > 0 && (
              <div className="card bg-base-100 border border-base-300 mt-6">
                <div className="card-body">
                  <h3 className="card-title text-lg">
                    {t('customers.details.notes')}
                  </h3>
                  <div className="space-y-2 mt-4">
                    {customer.notes.map((note) => (
                      <div key={note.id} className="p-3 bg-base-200 rounded-lg">
                        <p className="text-sm">{note.content}</p>
                        <div className="text-xs text-base-content/60 mt-2">
                          {new Date(note.createdAt).toLocaleDateString(
                            isArabic ? 'ar-AE' : 'en-US',
                            {
                              year: 'numeric',
                              month: 'long',
                              day: 'numeric',
                            },
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <EditCustomerSheet
        isOpen={isEditOpen}
        onClose={() => {
          setIsEditOpen(false)
        }}
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
        onClose={() => {
          setIsDeleteDialogOpen(false)
        }}
        title={t('customers.delete_confirm_title')}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => {
                setIsDeleteDialogOpen(false)
              }}
              disabled={isDeleting}
            >
              {t('common.cancel')}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={() => {
                void handleDeleteCustomer()
              }}
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
              onClick={() => {
                void handleDeleteAddressConfirm()
              }}
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
