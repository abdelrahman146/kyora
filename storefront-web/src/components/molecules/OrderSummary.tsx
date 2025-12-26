import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import type Decimal from 'decimal.js-light';
import { formatMoney } from '../../utils/money';

interface OrderSummaryProps {
  subtotal: Decimal;
  shipping: Decimal;
  total: Decimal;
  currency: string;
}

export const OrderSummary = memo<OrderSummaryProps>(function OrderSummary({
  subtotal,
  shipping,
  total,
  currency,
}) {
  const { t, i18n } = useTranslation();

  return (
    <div className="bg-neutral-50 rounded-lg p-3 space-y-2">
      <div className="flex justify-between text-sm">
        <span className="text-neutral-600">{t('subtotal')}</span>
        <span className="font-medium text-neutral-900">
          {formatMoney(subtotal, currency, i18n.language)}
        </span>
      </div>
      <div className="flex justify-between text-sm">
        <span className="text-neutral-600">{t('shipping')}</span>
        <span className="font-medium text-neutral-900">
          {shipping.greaterThan(0) ? formatMoney(shipping, currency, i18n.language) : t('free')}
        </span>
      </div>
      <div className="h-px bg-neutral-200 my-2" />
      <div className="flex justify-between">
        <span className="font-semibold text-neutral-900">{t('total')}</span>
        <span className="font-bold text-lg text-primary-700">
          {formatMoney(total, currency, i18n.language)}
        </span>
      </div>
    </div>
  );
});
