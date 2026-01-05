/**
 * ProductDetailsSheet Component
 *
 * A bottom sheet/drawer displaying detailed product information.
 * Redesigned for simplicity and clarity following Kyora's Basata principle.
 *
 * Features:
 * - Clean, single-column layout with clear information hierarchy
 * - Product header with name, category, and description
 * - Photo carousel with horizontal snap scroll
 * - Variant cards with clean vertical layout (no grids)
 * - Mobile: Bottom sheet sliding up from bottom
 * - Desktop: Side drawer from end
 * - RTL-compatible with proper logical properties
 * - Stock status color coding with tooltips
 */

import { useTranslation } from 'react-i18next'
import { ImageOff, Package } from 'lucide-react'

import { BottomSheet } from '../molecules/BottomSheet'
import { Tooltip } from '../atoms/Tooltip'
import type { Product } from '@/api/inventory'
import { formatCurrency } from '@/lib/formatCurrency'
import { useProductQuery } from '@/api/inventory'

export interface ProductDetailsSheetProps {
  isOpen: boolean
  onClose: () => void
  product: Product | null
  businessDescriptor: string
}

export function ProductDetailsSheet({
  isOpen,
  onClose,
  product,
  businessDescriptor,
}: ProductDetailsSheetProps) {
  const { t } = useTranslation()

  // Fetch full product details with variants
  const { data: fullProduct, isLoading } = useProductQuery(
    businessDescriptor,
    product?.id || '',
  )

  const displayProduct = fullProduct || product

  if (!displayProduct) return null

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={onClose}
      title={t('product_details', { ns: 'inventory' })}
      size="lg"
      side="end"
    >
      <div className="space-y-6">
        {/* Product Header */}
        <div>
          <h2 className="text-2xl font-bold mb-2 text-base-content">
            {displayProduct.name}
          </h2>
          {displayProduct.category && (
            <p className="text-sm text-base-content/60 mb-3">
              {displayProduct.category.name}
            </p>
          )}
          {displayProduct.description && (
            <p className="text-base text-base-content/70 leading-relaxed">
              {displayProduct.description}
            </p>
          )}
        </div>

        {/* Photos Carousel */}
        <div>
          <h3 className="text-lg font-semibold mb-3 text-base-content">
            {t('photos', { ns: 'inventory' })}
          </h3>
          {displayProduct.photos.length > 0 ? (
            <div className="flex overflow-x-auto snap-x snap-mandatory gap-3 pb-3 scrollbar-hide -mx-1 px-1">
              {displayProduct.photos.map((photo, index) => (
                <div
                  key={photo.assetId || index}
                  className="snap-start shrink-0 w-64 aspect-square"
                >
                  <img
                    src={photo.url}
                    alt={`${displayProduct.name} - ${index + 1}`}
                    className="w-full h-full object-cover rounded-xl border border-base-300"
                    loading="lazy"
                  />
                </div>
              ))}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center text-base-content/50 border border-base-300 rounded-xl bg-base-100">
              <ImageOff size={48} className="mb-3" />
              <p>{t('no_photos', { ns: 'inventory' })}</p>
            </div>
          )}
        </div>

        {/* Variants Section */}
        <div>
          <h3 className="text-lg font-semibold mb-3 text-base-content">
            {t('variants', { ns: 'inventory' })}
          </h3>
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <span className="loading loading-spinner loading-md"></span>
            </div>
          ) : displayProduct.variants && displayProduct.variants.length > 0 ? (
            <div className="space-y-3">
              {displayProduct.variants.map((variant) => {
                const isLowStock =
                  variant.stockQuantity > 0 &&
                  variant.stockQuantity <= variant.stockQuantityAlert
                const isOutOfStock = variant.stockQuantity === 0
                const variantPhoto =
                  variant.photos[0] || displayProduct.photos[0]

                const getStockColorClass = () => {
                  if (isOutOfStock) return 'text-error'
                  if (isLowStock) return 'text-warning'
                  return 'text-success'
                }

                const getStockTooltip = () => {
                  if (isOutOfStock)
                    return t('out_of_stock', { ns: 'inventory' })
                  if (isLowStock) return t('low_stock', { ns: 'inventory' })
                  return undefined
                }

                return (
                  <div
                    key={variant.id}
                    className="bg-base-100 border border-base-300 rounded-xl p-4  transition-shadow"
                  >
                    {/* Header: Photo + Code + SKU */}
                    <div className="flex items-start gap-3 mb-4">
                      <img
                        src={variantPhoto.thumbnailUrl || variantPhoto.url}
                        alt={variant.code}
                        className="w-16 h-16 object-cover rounded-lg border border-base-300 shrink-0"
                      />
                      <div className="flex-1 min-w-0">
                        <h4 className="font-semibold text-base text-base-content truncate">
                          {variant.code}
                        </h4>
                        {variant.sku && (
                          <p className="text-sm text-base-content/60 truncate">
                            {t('sku', { ns: 'inventory' })}: {variant.sku}
                          </p>
                        )}
                      </div>
                    </div>

                    {/* Pricing Info */}
                    <div className="space-y-2 mb-4">
                      <div className="flex items-baseline justify-between">
                        <span className="text-sm text-base-content/60">
                          {t('cost_price', { ns: 'inventory' })}
                        </span>
                        <span className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                          {formatCurrency(
                            parseFloat(variant.costPrice),
                            variant.currency,
                          )}
                        </span>
                      </div>
                      <div className="flex items-baseline justify-between">
                        <span className="text-sm text-base-content/60">
                          {t('sale_price', { ns: 'inventory' })}
                        </span>
                        <span className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                          {formatCurrency(
                            parseFloat(variant.salePrice),
                            variant.currency,
                          )}
                        </span>
                      </div>
                    </div>

                    {/* Stock Status */}
                    <div className="flex items-center justify-between pt-3 border-t border-base-300">
                      <span className="text-sm text-base-content/60">
                        {t('stock_quantity', { ns: 'inventory' })}
                      </span>
                      {getStockTooltip() ? (
                        <Tooltip content={getStockTooltip()}>
                          <span
                            className={`text-lg font-bold ${getStockColorClass()}`}
                          >
                            {variant.stockQuantity}
                          </span>
                        </Tooltip>
                      ) : (
                        <span
                          className={`text-lg font-bold ${getStockColorClass()}`}
                        >
                          {variant.stockQuantity}
                        </span>
                      )}
                    </div>
                  </div>
                )
              })}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center text-base-content/50 border border-base-300 rounded-xl bg-base-100">
              <Package size={48} className="mb-3" />
              <p>{t('no_variants', { ns: 'inventory' })}</p>
            </div>
          )}
        </div>
      </div>
    </BottomSheet>
  )
}
