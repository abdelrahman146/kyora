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
    <div className="card bg-base-100 border border-base-300  transition-shadow">
      <div className="card-body p-4">
        <div className="flex items-start justify-between gap-3">
          {/* Address Icon */}
          <div className="flex-shrink-0 mt-1">
            <MapPin size={20} className="text-primary" />
          </div>

          {/* Address Details */}
          <div className="flex-1 min-w-0">
            {/* City, State */}
            <h3 className="font-semibold text-base mb-1">
              {address.city}, {address.state}
            </h3>

            {/* Country */}
            <div className="flex items-center gap-2 text-sm text-base-content/70 mb-1">
              {countryInfo.flag && (
                <span className="text-base">{countryInfo.flag}</span>
              )}
              <span>{countryInfo.name}</span>
            </div>

            {/* Street */}
            {address.street && (
              <p className="text-sm text-base-content/70 mb-1">
                {address.street}
              </p>
            )}

            {/* Zip Code */}
            {address.zipCode && (
              <p className="text-sm text-base-content/70 mb-2">
                {t('customers.form.zip_code')}: {address.zipCode}
              </p>
            )}

            {/* Phone */}
            {phoneDisplay && (
              <div className="flex items-center gap-2 text-sm text-base-content/70">
                <Phone size={14} />
                <span dir="ltr">{phoneDisplay}</span>
              </div>
            )}
          </div>

          {/* Actions */}
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
      </div>
    </div>
  )
}
