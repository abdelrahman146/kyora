/**
 * InventoryCard Component
 *
 * A card component displaying product information for mobile list view.
 * Redesigned for simplicity, clarity, and elegance following Kyora's Basata principle.
 *
 * Features:
 * - Clean single-column layout with logical information hierarchy
 * - Product photo, name, and category at top
 * - Price information clearly displayed
 * - Stock status with color coding and tooltips
 * - Generous white space and touch-friendly (min 44px tap target)
 * - RTL-compatible with proper logical properties
 */

import { useTranslation } from 'react-i18next'

import { Avatar } from '../atoms/Avatar'
import { Tooltip } from '../atoms/Tooltip'
import type { Category, Product } from '@/api/inventory'
import { formatCurrency } from '@/lib/formatCurrency'
import {
  calculateTotalStock,
  getPriceRange,
  hasLowStock,
} from '@/lib/inventoryUtils'

export interface InventoryCardProps {
  product: Product
  currency: string
  categories: Array<Category>
  onClick?: (product: Product) => void
}

export function InventoryCard({
  product,
  currency,
  categories,
  onClick,
}: InventoryCardProps) {
  const { t } = useTranslation()

  const totalStock = calculateTotalStock(product.variants)
  const isLowStock = hasLowStock(product.variants)
  const isOutOfStock = totalStock === 0
  const costPriceRange = getPriceRange(product.variants, 'costPrice')
  const salePriceRange = getPriceRange(product.variants, 'salePrice')
  const category = categories.find((c) => c.id === product.categoryId)

  const getStockColorClass = () => {
    if (isOutOfStock) return 'text-error'
    if (isLowStock) return 'text-warning'
    return 'text-success'
  }

  const getStockTooltip = () => {
    if (isOutOfStock) return t('inventory.out_of_stock')
    if (isLowStock) return t('inventory.low_stock')
    return undefined
  }

  const formatPriceDisplay = (range: {
    min: number
    max: number
    isSame: boolean
  }) => {
    if (range.isSame) {
      return formatCurrency(range.min, currency)
    }
    return `${formatCurrency(range.min, currency)} - ${formatCurrency(range.max, currency)}`
  }

  return (
    <div
      className={`bg-base-100 border border-base-300 rounded-xl p-4 hover:shadow-md transition-shadow ${
        onClick ? 'cursor-pointer active:scale-[0.98]' : ''
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
      {/* Header: Photo + Name + Category */}
      <div className="flex items-start gap-3 mb-4">
        <Avatar
          src={product.photos[0]?.thumbnail_url || product.photos[0]?.url}
          alt={product.name}
          fallback={product.name.charAt(0).toUpperCase()}
          size="md"
          shape="square"
        />
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-base text-base-content truncate mb-1">
            {product.name}
          </h3>
          {category && (
            <p className="text-sm text-base-content/60">{category.name}</p>
          )}
        </div>
      </div>

      {/* Pricing Info */}
      <div className="space-y-2 mb-4">
        <div className="flex items-baseline justify-between">
          <span className="text-sm text-base-content/60">
            {t('inventory.cost_price')}
          </span>
          <span className="text-base font-semibold text-base-content">
            {formatPriceDisplay(costPriceRange)}
          </span>
        </div>
        <div className="flex items-baseline justify-between">
          <span className="text-sm text-base-content/60">
            {t('inventory.sale_price')}
          </span>
          <span className="text-base font-semibold text-base-content">
            {formatPriceDisplay(salePriceRange)}
          </span>
        </div>
      </div>

      {/* Stock Status */}
      <div className="flex items-center justify-between pt-3 border-t border-base-300">
        <span className="text-sm text-base-content/60">
          {t('inventory.stock_quantity')}
        </span>
        {getStockTooltip() ? (
          <Tooltip content={getStockTooltip()}>
            <span className={`text-lg font-bold ${getStockColorClass()}`}>
              {totalStock}
            </span>
          </Tooltip>
        ) : (
          <span className={`text-lg font-bold ${getStockColorClass()}`}>
            {totalStock}
          </span>
        )}
      </div>
    </div>
  )
}
