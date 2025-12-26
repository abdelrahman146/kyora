import { memo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Minus, Plus, Trash2 } from 'lucide-react';
import { ImageTile } from '../atoms';
import { useCartStore } from '../../cart/useCartStore';
import { formatMoney, money } from '../../utils/money';

interface CartItemProps {
  variantId: string;
  productId: string;
  productName: string;
  variantName: string;
  quantity: number;
  unitPrice: string;
  currency: string;
  photoUrl?: string;
  note?: string;
}

export const CartItem = memo<CartItemProps>(function CartItem({
  variantId,
  productId,
  productName,
  variantName,
  quantity,
  unitPrice,
  currency,
  photoUrl,
  note,
}) {
  const { t, i18n } = useTranslation();

  const handleDecrease = useCallback(() => {
    useCartStore.getState().removeItem(variantId, 1);
  }, [variantId]);

  const handleIncrease = useCallback(() => {
    useCartStore.getState().addItem(
      { variantId, productId, productName, variantName, unitPrice, currency, photoUrl },
      1
    );
  }, [variantId, productId, productName, variantName, unitPrice, currency, photoUrl]);

  const handleRemove = useCallback(() => {
    useCartStore.getState().setQuantity(variantId, 0);
  }, [variantId]);

  const handleNoteChange = useCallback(
    (value: string) => {
      useCartStore.getState().updateNote(variantId, value);
    },
    [variantId]
  );

  return (
    <div className="bg-white rounded-2xl border border-neutral-200 p-3">
      <div className="flex items-start gap-3">
        <div className="w-16 shrink-0">
          <ImageTile src={photoUrl} alt={productName} aspectClassName="aspect-square" />
        </div>

        <div className="min-w-0 flex-1 space-y-3">
          <div className="flex items-start justify-between gap-3">
            <div className="min-w-0 flex-1">
              <h4 className="font-semibold text-sm text-neutral-900 truncate">{productName}</h4>
              <p className="text-xs text-neutral-500 truncate mt-0.5">{variantName}</p>
            </div>
            <div className="text-sm font-bold text-primary-700 whitespace-nowrap">
              {formatMoney(money(unitPrice).mul(quantity), currency, i18n.language)}
            </div>
          </div>

          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-1 bg-base-200/50 rounded-full px-1 py-0.5">
              <button
                type="button"
                className="btn btn-xs btn-circle btn-ghost active-scale focus-ring min-h-0 h-5 w-5"
                onClick={handleDecrease}
                aria-label={t('decrease')}
              >
                <Minus className="w-2.5 h-2.5" strokeWidth={2.5} />
              </button>
              <span className="px-1 font-semibold text-xs min-w-[1rem] text-center">
                {quantity}
              </span>
              <button
                type="button"
                className="btn btn-xs btn-circle btn-primary active-scale focus-ring min-h-0 h-5 w-5"
                onClick={handleIncrease}
                aria-label={t('increase')}
              >
                <Plus className="w-2.5 h-2.5" strokeWidth={2.5} />
              </button>
            </div>

            <button
              type="button"
              className="btn btn-ghost btn-sm btn-square active-scale focus-ring text-error hover:bg-error hover:text-error-content"
              onClick={handleRemove}
              aria-label={t('remove')}
            >
              <Trash2 className="w-4 h-4" strokeWidth={2} />
            </button>
          </div>

          <input
            className="input input-sm w-full border-neutral-200 bg-neutral-50 focus:bg-white focus:border-primary-500 text-sm placeholder:text-neutral-400"
            value={note || ''}
            onChange={(e) => handleNoteChange(e.target.value)}
            placeholder={`${t('specialRequest')} (${t('optional')})`}
          />
        </div>
      </div>
    </div>
  );
});
