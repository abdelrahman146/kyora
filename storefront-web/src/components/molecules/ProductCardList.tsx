import { memo, useMemo, useCallback } from 'react';
import { Plus, Minus } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import type { PublicProduct, PublicVariant } from '../../api/types';
import { formatMoney, money } from '../../utils/money';
import { ImageTile } from '../atoms';

interface ProductCardListProps {
  product: PublicProduct;
  variant: PublicVariant;
  quantity: number;
  onAddToCart: () => void;
  onIncrement: () => void;
  onDecrement: () => void;
  onVariantChange?: (variantId: string) => void;
}

/**
 * ProductCardList Molecule - Horizontal product card for list view
 * Optimized layout for list display across all screen sizes
 * Memoized to prevent unnecessary re-renders
 */
export const ProductCardList = memo<ProductCardListProps>(function ProductCardList({
  product,
  variant,
  quantity,
  onAddToCart,
  onIncrement,
  onDecrement,
  onVariantChange,
}) {
  const { t, i18n } = useTranslation();

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

  const handleVariantChange = useCallback(
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      onVariantChange?.(e.target.value);
    },
    [onVariantChange]
  );

  return (
    <div className="bg-white rounded-2xl border border-base-300/50 overflow-hidden">
      <div className="flex gap-3 p-3">
        {/* Fixed Image on Left with 2px spacing */}
        <div className="shrink-0 bg-base-50 p-0.5 rounded-xl">
          <div className="w-16 h-16 md:w-20 md:h-20 rounded-lg overflow-hidden">
            <ImageTile src={photoUrl} alt={product.name} aspectClassName="aspect-square" />
          </div>
        </div>

        {/* Product Info on Right */}
        <div className="flex-1 flex flex-col justify-between min-w-0">
          {/* Top Section: Name and Price */}
          <div className="space-y-1">
            <h3 className="font-semibold text-sm md:text-base line-clamp-1">
              {product.name}
            </h3>
            {product.description && (
              <p className="text-xs text-base-content/60 line-clamp-1">
                {product.description}
              </p>
            )}
            <div className="font-bold text-sm md:text-base text-primary">
              {formattedPrice}
            </div>
          </div>

          {/* Bottom Section: Controls */}
          <div className="flex items-center justify-between gap-2 mt-1">
            {/* Variant Selector */}
            {hasMultipleVariants && onVariantChange ? (
              <select
                className="select select-xs bg-base-100 border-base-300 focus:border-primary text-xs flex-1 min-w-0"
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
            ) : (
              <div className="flex-1" />
            )}

            {/* Action Buttons */}
            <div className="flex items-center gap-1 shrink-0">
              {quantity === 0 ? (
                <button
                  type="button"
                  onClick={onAddToCart}
                  className="btn btn-primary btn-xs btn-circle active-scale focus-ring"
                  aria-label={t('add')}
                >
                  <Plus className="w-3 h-3" strokeWidth={2.5} />
                </button>
              ) : (
                <div className="flex items-center gap-0.5 bg-base-200/50 rounded-full px-1 py-0.5">
                  <button
                    type="button"
                    onClick={onDecrement}
                    className="btn btn-xs btn-circle btn-ghost active-scale focus-ring min-h-0 h-5 w-5"
                    aria-label={t('decrease')}
                  >
                    <Minus className="w-2.5 h-2.5" strokeWidth={2.5} />
                  </button>
                  <span className="font-semibold text-xs min-w-[1rem] text-center px-1">
                    {quantity}
                  </span>
                  <button
                    type="button"
                    onClick={onIncrement}
                    className="btn btn-xs btn-circle btn-primary active-scale focus-ring min-h-0 h-5 w-5"
                    aria-label={t('increase')}
                  >
                    <Plus className="w-2.5 h-2.5" strokeWidth={2.5} />
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
});
