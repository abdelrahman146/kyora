import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { ChevronRight } from 'lucide-react';
import { formatMoney } from '../../utils/money';
import type Decimal from 'decimal.js-light';

interface StickyCartBarProps {
  itemCount: number;
  totalAmount: Decimal;
  currency: string;
  onViewCart: () => void;
}

/**
 * StickyCartBar Molecule - Bottom cart summary
 * Memoized to prevent unnecessary re-renders
 * Optimized with useMemo for formatted price
 * Follows KDS 5.3 with high contrast primary-900 background
 */
export const StickyCartBar = memo<StickyCartBarProps>(function StickyCartBar({
  itemCount,
  totalAmount,
  currency,
  onViewCart,
}) {
  const { t, i18n } = useTranslation();

  // Memoize formatted price to prevent recalculation on every render
  const formattedPrice = useMemo(
    () => formatMoney(totalAmount, currency, i18n.language),
    [totalAmount, currency, i18n.language]
  );

  // Memoize item text to prevent recalculation
  const itemText = useMemo(
    () => `${itemCount} ${itemCount === 1 ? t('item') : t('items')}`,
    [itemCount, t]
  );

  if (itemCount <= 0) return null;

  return (
    <div className="fixed bottom-20 inset-x-0 z-40 safe-bottom">
      <div className="mx-auto max-w-5xl px-4 md:px-6 pb-4">
        <button
          type="button"
          onClick={onViewCart}
          className="w-full h-14 rounded-2xl bg-primary text-primary-content shadow-xl flex items-center justify-between px-5 md:px-6 active-scale focus-ring transition-all"
          aria-label={t('viewCart')}
        >          {/* Left Side: Item Count and Price */}
          <div className="flex items-center gap-2 font-semibold text-sm md:text-base">
            <span>{itemText}</span>
            <span className="opacity-60">•</span>
            <span className="font-bold">{formattedPrice}</span>
          </div>

          {/* Right Side: "View Cart" with Chevron */}
          <div className="flex items-center gap-1 font-semibold text-sm md:text-base">
            <span>{t('viewCart')}</span>
            <ChevronRight className="w-5 h-5 rtl-mirror" strokeWidth={2.5} />
          </div>
        </button>
      </div>
    </div>
  );
});
