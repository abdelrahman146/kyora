/**
 * InventoryCard Component
 *
 * A card component displaying product information for mobile list view.
 *
 * Features:
 * - Displays product name, category, cost price, stock quantity
 * - Stock status badge with visual indicators
 * - Low stock warning badge
 * - Click to view details
 * - Photo thumbnail with fallback icon
 * - RTL-compatible
 * - Mobile-optimized (min 44px touch target)
 */

import { useTranslation } from 'react-i18next'
import { Box, DollarSign } from 'lucide-react'

import { Avatar } from '../atoms/Avatar'
import { Badge } from '../atoms/Badge'
import { Tooltip } from '../atoms/Tooltip'
import type { Product } from '@/api/inventory'
import { formatCurrency } from '@/lib/formatCurrency'

export interface InventoryCardProps {
  product: Product
  avgCostPrice: number
  totalStock: number
  stockStatus: 'in_stock' | 'low_stock' | 'out_of_stock'
  hasLowStock: boolean
  currency: string
  onClick?: (product: Product) => void
}

export function InventoryCard({
  product,
  avgCostPrice,
  totalStock,
  stockStatus,
  hasLowStock,
  currency,
  onClick,
}: InventoryCardProps) {
  const { t } = useTranslation()

  const getStockStatusBadgeClass = () => {
    switch (stockStatus) {
      case 'in_stock':
        return 'badge-success'
      case 'low_stock':
        return 'badge-warning'
      case 'out_of_stock':
        return 'badge-error'
      default:
        return 'badge-neutral'
    }
  }

  return (
    <div
      className={`card bg-base-100 border border-base-300  transition-shadow ${
        onClick ? 'cursor-pointer' : ''
      }`}
      onClick={() => onClick?.(product)}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      onKeyDown={(e) => {
        if (onClick && (e.key === 'Enter' || e.key === ' ')) {
          e.preventDefault()
          onClick(product)
        }
      }}
    >
      <div className="card-body p-4">
        {/* Header: Product Photo + Name + Category */}
        <div className="flex items-center gap-3 mb-3">
          <Avatar
            src={product.photos[0]?.thumbnailUrl}
            alt={product.name}
            fallback={product.name.charAt(0).toUpperCase()}
            size="sm"
          />
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-base truncate">{product.name}</h3>
            {product.category && (
              <p className="text-sm text-base-content/70 truncate">
                {product.category.name}
              </p>
            )}
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 gap-3">
          {/* Average Cost Price */}
          <div className="stat-item">
            <div className="flex items-center gap-1 text-sm text-base-content/60 mb-1">
              <DollarSign size={14} />
              <Tooltip content={t('average_cost_tooltip', { ns: 'inventory' })}>
                <span className="cursor-help">
                  {t('cost_price_avg', { ns: 'inventory' })}
                </span>
              </Tooltip>
            </div>
            <div className="text-base font-semibold">
              {formatCurrency(avgCostPrice, currency)}
            </div>
          </div>

          {/* Stock Quantity */}
          <div className="stat-item">
            <div className="flex items-center gap-1 text-sm text-base-content/60 mb-1">
              <Box size={14} />
              <span>{t('stock_quantity', { ns: 'inventory' })}</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-base font-semibold">{totalStock}</span>
              {hasLowStock && (
                <Badge size="sm" variant="warning">
                  {t('low_stock', { ns: 'inventory' })}
                </Badge>
              )}
            </div>
          </div>
        </div>

        {/* Status Badge */}
        <div className="mt-3">
          <div
            className={`badge badge-sm w-full ${getStockStatusBadgeClass()}`}
          >
            {t(`inventory.${stockStatus}`)}
          </div>
        </div>
      </div>
    </div>
  )
}
