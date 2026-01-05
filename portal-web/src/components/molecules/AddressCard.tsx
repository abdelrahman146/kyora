/**
 * AddressCard Component
 *
 * Displays a customer address in a card format with actions.
 * Mobile-first, RTL-aware, reusable component.
 *
 * Features:
 * - Clean card layout with address details
 * - Country flag and name display
 * - Edit and delete actions
 * - Responsive design
 */

import { Edit, MapPin, Phone, Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useStore } from '@tanstack/react-store'

import type { CustomerAddress } from '@/api/customer'
import { metadataStore } from '@/stores/metadataStore'

interface AddressCardProps {
  address: CustomerAddress
  onEdit?: () => void
  onDelete?: () => void
  isDeleting?: boolean
}

export function AddressCard({
  address,
  onEdit,
  onDelete,
  isDeleting,
}: AddressCardProps) {
  const { t, i18n } = useTranslation()
  const countries = useStore(metadataStore, (s) => s.countries)
  const isArabic = i18n.language.toLowerCase().startsWith('ar')

  const getCountryInfo = (countryCode: string) => {
    const country = countries.find((c) => c.code === countryCode)
    return {
      name: isArabic
        ? (country?.nameAr ?? countryCode)
        : (country?.name ?? countryCode),
      flag: country?.flag,
    }
  }

  const countryInfo = getCountryInfo(address.countryCode)

  const phoneDisplay =
    address.phoneCode && address.phoneNumber
      ? `${address.phoneCode} ${address.phoneNumber}`
      : null

  return (
    <div className="border border-base-300 rounded-xl p-4 bg-base-100">
      {/* Header with Location Icon and Actions */}
      <div className="flex items-start gap-3 mb-3">
        <div className="flex-shrink-0 mt-0.5">
          <div className="size-10 rounded-lg bg-primary/10 flex items-center justify-center">
            <MapPin size={20} className="text-primary" />
          </div>
        </div>

        <div className="flex-1 min-w-0">
          <h3 className="font-bold text-base leading-tight mb-1">
            {address.city}, {address.state}
          </h3>
          <div className="flex items-center gap-2 text-sm text-base-content/60">
            {countryInfo.flag && (
              <span className="text-lg leading-none">{countryInfo.flag}</span>
            )}
            <span>{countryInfo.name}</span>
          </div>
        </div>

        {/* Action Buttons */}
        {(onEdit ?? onDelete) && (
          <div className="flex-shrink-0 flex gap-1">
            {onEdit && (
              <button
                type="button"
                className="btn btn-ghost btn-sm btn-square"
                onClick={onEdit}
                disabled={isDeleting}
                aria-label={t('common.edit')}
              >
                <Edit size={16} />
              </button>
            )}
            {onDelete && (
              <button
                type="button"
                className="btn btn-ghost btn-sm btn-square text-error hover:bg-error/10"
                onClick={onDelete}
                disabled={isDeleting}
                aria-label={t('common.delete')}
              >
                {isDeleting ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  <Trash2 size={16} />
                )}
              </button>
            )}
          </div>
        )}
      </div>

      {/* Address Details */}
      <div className="space-y-2 border-t border-base-300 pt-3">
        {address.street && (
          <div className="flex items-start gap-2">
            <span className="text-sm text-base-content/60 min-w-[4rem]">
              {t('customers.form.street')}:
            </span>
            <span className="text-sm font-medium flex-1">{address.street}</span>
          </div>
        )}

        {address.zipCode && (
          <div className="flex items-start gap-2">
            <span className="text-sm text-base-content/60 min-w-[4rem]">
              {t('customers.form.zip_code')}:
            </span>
            <span className="text-sm font-medium flex-1">
              {address.zipCode}
            </span>
          </div>
        )}

        {phoneDisplay && (
          <div className="flex items-start gap-2">
            <Phone size={14} className="text-base-content/60 mt-0.5" />
            <span className="text-sm font-medium" dir="ltr">
              {phoneDisplay}
            </span>
          </div>
        )}
      </div>
    </div>
  )
}
