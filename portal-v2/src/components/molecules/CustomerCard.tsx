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

import { DollarSign, Phone, ShoppingBag } from 'lucide-react'
import { Avatar } from '../atoms/Avatar'

// Define Customer interface inline (until customer types are properly set up)
interface Customer {
  id: string
  name: string
  email?: string
  phone?: string
  phoneCode?: string
  phoneNumber?: string
  avatarUrl?: string
  totalOrders?: number
  totalSpent?: number
}

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
                Orders
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
              <div className="text-xs text-base-content/60 truncate">Total</div>
              <div className="font-semibold truncate">
                {currency} {formatCurrency(totalSpent)}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
