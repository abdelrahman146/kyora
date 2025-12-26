import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { ShoppingBag } from 'lucide-react';
import { LanguageSwitcher } from '../molecules';

interface StorefrontHeaderProps {
  cartItemCount: number;
  onCartClick: () => void;
}

/**
 * StorefrontHeader Organism - Minimal transparent header
 * Only contains cart and language switcher on opposite sides
 * Not sticky, short height
 * Memoized to prevent unnecessary re-renders
 */
export const StorefrontHeader = memo<StorefrontHeaderProps>(function StorefrontHeader({
  cartItemCount,
  onCartClick,
}) {
  const { t } = useTranslation();

  return (
    <header className="bg-transparent safe-top">
      <div className="mx-auto max-w-5xl">
        <div className="flex items-center justify-between px-4 py-3 min-h-14">
          {/* Left: Language Switcher */}
          <LanguageSwitcher />

          {/* Right: Cart Button */}
          <button
            type="button"
            onClick={onCartClick}
            className="btn btn-sm btn-square bg-base-100 border border-base-300 hover:border-primary hover:bg-base-100 active-scale focus-ring relative shadow-sm"
            aria-label={t('cart')}
          >
            <ShoppingBag className="w-5 h-5" strokeWidth={2} />
            {cartItemCount > 0 && (
              <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-primary-content text-[10px] font-bold">
                {cartItemCount > 9 ? '9+' : cartItemCount}
              </span>
            )}
          </button>
        </div>
      </div>
    </header>
  );
});
