import { memo, useMemo, useCallback } from 'react';
import { Plus, Minus } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import type { PublicProduct, PublicVariant } from '../../api/types';
import { formatMoney, money } from '../../utils/money';
import { ImageTile } from '../atoms';

interface ProductCardProps {
  product: PublicProduct;
  variant: PublicVariant;
  quantity: number;
  onAddToCart: () => void;
  onIncrement: () => void;
  onDecrement: () => void;
  onVariantChange?: (variantId: string) => void;
}

/**
 * ProductCard Molecule - Product display with cart controls
 * Memoized to prevent unnecessary re-renders
 * Optimized with useMemo for computed values and useCallback for event handlers
 * Follows KDS 5.2 with 2-column grid layout and floating add button
 */
export const ProductCard = memo<ProductCardProps>(function ProductCard({
  product,
  variant,
  quantity,
  onAddToCart,
  onIncrement,
  onDecrement,
  onVariantChange,
}) {
  const { t, i18n } = useTranslation();

  // Memoize computed values to prevent unnecessary recalculations
  const photoUrl = useMemo(
    () => variant.photos?.[0] || product.photos?.[0],
    [variant.photos, product.photos]
  );

  const hasMultipleVariants = useMemo(
    () => product.variants.length > 1,
    [product.variants.length]
  );

  const formattedPrice = useMemo(
    () => formatMoney(money(variant.salePrice), variant.currency, i18n.language),
    [variant.salePrice, variant.currency, i18n.language]
  );

  const isOutOfStock = false; // Stock tracking not implemented yet

  const isLowStock = false; // Stock tracking not implemented yet

  // Memoize event handler to prevent re-creating function on every render
  const handleVariantChange = useCallback(
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      onVariantChange?.(e.target.value);
    },
    [onVariantChange]
  );

  return (
    <div className="bg-white rounded-2xl border border-base-300/50 overflow-hidden">
      {/* Image Container with 2px spacing */}
      <div className="relative overflow-hidden bg-base-50 p-1">
        <div className="rounded-xl overflow-hidden">
          <ImageTile src={photoUrl} alt={product.name} aspectClassName="aspect-square" />
        </div>
        
        {/* Stock badges if needed */}
        {isOutOfStock && (
          <div className="absolute top-2 start-2 px-2 py-1 bg-error text-error-content text-xs font-bold rounded-full">
            {t('outOfStock')}
          </div>
        )}
        {isLowStock && !isOutOfStock && (
          <div className="absolute top-2 start-2 px-2 py-1 bg-warning text-warning-content text-xs font-bold rounded-full">
            {t('lowStock')}
          </div>
        )}
      </div>

      {/* Product Info with consistent structure */}
      <div className="p-3 flex flex-col">
        {/* Top section: Title and Price - flexible height */}
        <div className="flex-1">
          {/* Title (Max 2 lines) */}
          <h3 className="font-semibold text-sm md:text-base leading-tight line-clamp-2 text-base-content mb-2">
            {product.name}
          </h3>

          {/* Price */}
          <div className="text-primary font-bold text-base md:text-lg">
            {formattedPrice}
          </div>
        </div>

        {/* Bottom section: Always at bottom for alignment */}
        <div className="mt-3 space-y-2">
          {/* Variant Selector (if multiple variants) - Always at bottom */}
          {hasMultipleVariants && onVariantChange && (
            <select
              className="select select-sm w-full bg-base-100 border-base-300 focus:border-primary text-xs"
              value={variant.id}
              onChange={handleVariantChange}
              aria-label={t('selectVariant')}
            >
              {product.variants.map((v) => (
                <option key={v.id} value={v.id}>
                  {v.name}
                </option>
              ))}
            </select>
          )}

          {/* Full-width Action Button */}
          {quantity === 0 ? (
            <button
              type="button"
              onClick={onAddToCart}
              className="btn btn-primary btn-sm w-full active-scale focus-ring"
              aria-label={t('addToCart')}
            >
              <Plus className="w-4 h-4" strokeWidth={2.5} />
              <span className="text-xs md:text-sm font-semibold">{t('addToCart')}</span>
            </button>
          ) : (
            <div className="flex items-center justify-center gap-1 bg-base-200/50 rounded-full px-2 py-1.5 w-full">
              <button
                type="button"
                onClick={onDecrement}
                className="btn btn-xs btn-circle btn-ghost active-scale focus-ring"
                aria-label={t('decrease')}
              >
                <Minus className="w-3 h-3" strokeWidth={2.5} />
              </button>
              <span className="px-3 font-semibold text-sm min-w-[2rem] text-center">
                {quantity}
              </span>
              <button
                type="button"
                onClick={onIncrement}
                className="btn btn-xs btn-circle btn-primary active-scale focus-ring"
                aria-label={t('increase')}
              >
                <Plus className="w-3 h-3" strokeWidth={2.5} />
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
});

/**
 * ProductCardSkeleton Molecule - Loading state for ProductCard
 * Memoized to prevent unnecessary re-renders
 * Follows KDS 4.6 with pulsing animation
 */
export const ProductCardSkeleton = memo(function ProductCardSkeleton() {
  return (
    <div className="bg-white rounded-2xl shadow-sm border border-neutral-200 overflow-hidden">
      {/* Image Skeleton */}
      <div className="aspect-square bg-neutral-200 animate-pulse" />

      {/* Content Skeleton */}
      <div className="p-4 space-y-2">
        {/* Title Skeleton (2 lines) */}
        <div className="space-y-2">
          <div className="h-4 bg-neutral-200 rounded animate-pulse w-full" />
          <div className="h-4 bg-neutral-200 rounded animate-pulse w-3/4" />
        </div>

        {/* Price Skeleton */}
        <div className="h-5 bg-neutral-200 rounded animate-pulse w-1/2" />
      </div>
    </div>
  );
});
