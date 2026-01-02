/**
 * ProductDetailsSheet Component
 *
 * A bottom sheet/drawer displaying detailed product information including:
 * - Product name, description, category
 * - Photo carousel with horizontal scroll
 * - Variants table with stock information
 *
 * Features:
 * - Mobile: Bottom sheet sliding up from bottom
 * - Desktop: Side drawer from end
 * - Photo carousel with snap scrolling
 * - Scrollable variants table (max-h-96)
 * - RTL-compatible
 * - Empty states for no photos/variants
 */

import { useTranslation } from 'react-i18next'
import { ImageOff, Package } from 'lucide-react'

import { BottomSheet } from '../molecules/BottomSheet'
import { Badge } from '../atoms/Badge'
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
      title={t('inventory.product_details')}
      size="lg"
      side="end"
    >
      <div className="space-y-6">
        {/* Product Header */}
        <div>
          <h2 className="text-xl font-bold mb-2">{displayProduct.name}</h2>
          {displayProduct.category && (
            <Badge variant="neutral" className="mb-3">
              {displayProduct.category.name}
            </Badge>
          )}
          {displayProduct.description && (
            <p className="text-base-content/70 leading-relaxed">
              {displayProduct.description}
            </p>
          )}
        </div>

        {/* Photos Carousel */}
        <div>
          <h3 className="text-lg font-semibold mb-3">
            {t('inventory.photos')}
          </h3>
          {displayProduct.photos.length > 0 ? (
            <div className="flex overflow-x-auto snap-x snap-mandatory gap-3 pb-3 scrollbar-hide">
              {displayProduct.photos.map((photo, index) => (
                <div
                  key={photo.asset_id || index}
                  className="snap-start shrink-0 w-64 aspect-square"
                >
                  <img
                    src={photo.url}
                    alt={`${displayProduct.name} - ${index + 1}`}
                    className="w-full h-full object-cover rounded-lg border border-base-300"
                    loading="lazy"
                  />
                </div>
              ))}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center text-base-content/50">
              <ImageOff size={48} className="mb-3" />
              <p>{t('inventory.no_photos')}</p>
            </div>
          )}
        </div>

        {/* Variants Section */}
        <div>
          <h3 className="text-lg font-semibold mb-3">
            {t('inventory.variants')}
          </h3>
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <span className="loading loading-spinner loading-md"></span>
            </div>
          ) : displayProduct.variants && displayProduct.variants.length > 0 ? (
            <div className="space-y-3 max-h-96 overflow-y-auto">
              {displayProduct.variants.map((variant) => {
                const isLowStock =
                  variant.stock_quantity > 0 &&
                  variant.stock_quantity <= variant.stock_quantity_alert
                const isOutOfStock = variant.stock_quantity === 0
                const variantPhoto =
                  variant.photos?.[0] || displayProduct.photos[0]

                return (
                  <div
                    key={variant.id}
                    className="border border-base-300 rounded-lg p-4 bg-base-100 hover:bg-base-200 transition-colors"
                  >
                    <div className="flex gap-3">
                      {/* Variant Thumbnail */}
                      {variantPhoto ? (
                        <img
                          src={variantPhoto.thumbnail_url || variantPhoto.url}
                          alt={variant.code}
                          className="w-16 h-16 object-cover rounded-lg border border-base-300 shrink-0"
                        />
                      ) : (
                        <div className="w-16 h-16 bg-base-200 rounded-lg border border-base-300 flex items-center justify-center shrink-0">
                          <Package size={24} className="text-base-content/40" />
                        </div>
                      )}

                      {/* Variant Info */}
                      <div className="flex-1 min-w-0">
                        {/* Header with code and status */}
                        <div className="flex items-start justify-between gap-2 mb-2">
                          <div className="min-w-0 flex-1">
                            <p className="font-semibold text-base truncate">
                              {variant.code}
                            </p>
                            {variant.sku && (
                              <p className="text-sm text-base-content/60 truncate">
                                {t('inventory.sku')}: {variant.sku}
                              </p>
                            )}
                          </div>
                          {isOutOfStock ? (
                            <Badge size="sm" variant="error">
                              {t('inventory.out_of_stock')}
                            </Badge>
                          ) : isLowStock ? (
                            <Badge size="sm" variant="warning">
                              {t('inventory.low_stock')}
                            </Badge>
                          ) : (
                            <Badge size="sm" variant="success">
                              {t('inventory.in_stock')}
                            </Badge>
                          )}
                        </div>

                        {/* Pricing and Stock Grid */}
                        <div className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                          <div>
                            <span className="text-base-content/60">
                              {t('inventory.cost_price')}:
                            </span>
                            <span className="ms-1 font-medium">
                              {formatCurrency(
                                parseFloat(variant.cost_price),
                                variant.currency,
                              )}
                            </span>
                          </div>
                          <div>
                            <span className="text-base-content/60">
                              {t('inventory.sale_price')}:
                            </span>
                            <span className="ms-1 font-medium">
                              {formatCurrency(
                                parseFloat(variant.sale_price),
                                variant.currency,
                              )}
                            </span>
                          </div>
                          <div>
                            <span className="text-base-content/60">
                              {t('inventory.stock_quantity')}:
                            </span>
                            <span className="ms-1 font-semibold">
                              {variant.stock_quantity}
                            </span>
                          </div>
                          <div>
                            <span className="text-base-content/60">
                              {t('inventory.stock_alert')}:
                            </span>
                            <span className="ms-1 text-base-content/70">
                              {variant.stock_quantity_alert}
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center text-base-content/50 border border-base-300 rounded-lg">
              <Package size={48} className="mb-3" />
              <p>{t('inventory.no_variants')}</p>
            </div>
          )}
        </div>
      </div>
    </BottomSheet>
  )
}
