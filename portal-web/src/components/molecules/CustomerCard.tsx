/**
 * CustomerCard Component
 *
 * A card component displaying customer information for mobile list view.
 *
 * Features:
 * - Displays customer name, phone, orders count, total spent
 * - Click to view details
 * - Avatar/initials display
 * - RTL-compatible
 * - Mobile-optimized (min 44px touch target)
 */

import { useTranslation } from 'react-i18next'
import { Calendar, DollarSign, MapPin, Phone, ShoppingBag } from 'lucide-react'

import { Avatar } from '../atoms/Avatar'
import type { Customer } from '@/api/customer'
import { getMetadata } from '@/stores/metadataStore'
import { formatDateShort } from '@/lib/formatDate'

export interface CustomerCardProps {
  customer: Customer
  onClick?: (customer: Customer) => void
  ordersCount?: number
  totalSpent?: number
  currency?: string
}

export function CustomerCard({
  customer,
  onClick,
  ordersCount = 0,
  totalSpent = 0,
  currency = 'AED',
}: CustomerCardProps) {
  const { t } = useTranslation()

  const getInitials = (name: string): string => {
    return name
      .split(' ')
      .map((word) => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  const formatPhone = (): string | null => {
    if (customer.phoneCode && customer.phoneNumber) {
      return `${customer.phoneCode} ${customer.phoneNumber}`
    }
    return null
  }

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'decimal',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(amount)
  }

  const metadata = getMetadata()
  const country = metadata.countries.find(
    (c) =>
      c.code === customer.countryCode || c.iso_code === customer.countryCode,
  )

  return (
    <div
      className={`card bg-base-100 border border-base-300 shadow-sm hover:shadow-md transition-shadow ${
        onClick ? 'cursor-pointer' : ''
      }`}
      onClick={() => onClick?.(customer)}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      onKeyDown={(e) => {
        if (onClick && (e.key === 'Enter' || e.key === ' ')) {
          e.preventDefault()
          onClick(customer)
        }
      }}
    >
      <div className="card-body p-4">
        {/* Header: Avatar + Name */}
        <div className="flex items-center gap-3 mb-3">
          <Avatar
            src={customer.avatarUrl}
            alt={customer.name}
            fallback={getInitials(customer.name)}
            size="lg"
          />
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-base truncate">
              {customer.name}
            </h3>
            {formatPhone() && (
              <div className="flex items-center gap-1 text-sm text-base-content/70">
                <Phone size={14} />
                <span className="truncate">{formatPhone()}</span>
              </div>
            )}
            {country && (
              <div className="flex items-center gap-1 text-sm text-base-content/70">
                <MapPin size={14} />
                {country.flag && <span>{country.flag}</span>}
                <span className="truncate">{country.name}</span>
              </div>
            )}
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 gap-3">
          {/* Orders Count */}
          <div className="flex items-center gap-2 bg-base-200 rounded-lg p-2">
            <div className="p-1.5 bg-success/10 rounded-lg">
              <ShoppingBag size={16} className="text-success" />
            </div>
            <div className="flex-1 min-w-0">
              <div className="text-xs text-base-content/60 truncate">
                {t('customers.orders_count')}
              </div>
              <div className="font-semibold">{ordersCount}</div>
            </div>
          </div>

          {/* Total Spent */}
          <div className="flex items-center gap-2 bg-base-200 rounded-lg p-2">
            <div className="p-1.5 bg-primary/10 rounded-lg">
              <DollarSign size={16} className="text-primary" />
            </div>
            <div className="flex-1 min-w-0">
              <div className="text-xs text-base-content/60 truncate">
                {t('customers.total_spent')}
              </div>
              <div className="font-semibold truncate">
                {currency} {formatCurrency(totalSpent)}
              </div>
            </div>
          </div>
        </div>

        {/* Date Added */}
        <div className="flex items-center gap-2 text-sm text-base-content/60 pt-2 border-t border-base-300">
          <Calendar size={14} />
          <span>{formatDateShort(customer.joinedAt)}</span>
        </div>
      </div>
    </div>
  )
}
