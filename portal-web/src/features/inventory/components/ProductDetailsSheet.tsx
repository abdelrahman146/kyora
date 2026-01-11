import { useTranslation } from 'react-i18next'
import { ImageOff, Package } from 'lucide-react'

import type { Product } from '@/api/inventory'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Tooltip } from '@/components/atoms/Tooltip'
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
  const { t: tInventory } = useTranslation('inventory')

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
      title={tInventory('product_details')}
      size="lg"
      side="end"
    >
      <div className="space-y-6">
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

        <div>
          <h3 className="text-lg font-semibold mb-3 text-base-content">
            {tInventory('photos')}
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
              <p>{tInventory('no_photos')}</p>
            </div>
          )}
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-3 text-base-content">
            {tInventory('variants')}
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
                  if (isOutOfStock) return tInventory('out_of_stock')
                  if (isLowStock) return tInventory('low_stock')
                  return undefined
                }

                return (
                  <div
                    key={variant.id}
                    className="bg-base-100 border border-base-300 rounded-xl p-4  transition-shadow"
                  >
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
                            {tInventory('sku')}: {variant.sku}
                          </p>
                        )}
                      </div>
                    </div>

                    <div className="space-y-2 mb-4">
                      <div className="flex items-baseline justify-between">
                        <span className="text-sm text-base-content/60">
                          {tInventory('cost_price')}
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
                          {tInventory('sale_price')}
                        </span>
                        <span className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                          {formatCurrency(
                            parseFloat(variant.salePrice),
                            variant.currency,
                          )}
                        </span>
                      </div>
                    </div>

                    <div className="flex items-center justify-between pt-3 border-t border-base-300">
                      <span className="text-sm text-base-content/60">
                        {tInventory('stock_quantity')}
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
              <p>{tInventory('no_variants')}</p>
            </div>
          )}
        </div>
      </div>
    </BottomSheet>
  )
}
