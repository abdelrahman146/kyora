import { AlertCircle, RefreshCw, Sparkles } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import type { OrderPreview } from '@/api/order'
import { formatCurrency } from '@/lib/formatCurrency'

interface LiveSummaryCardProps {
  preview: OrderPreview | null
  isLoading: boolean
  isStale: boolean
  errorMessage?: string | null
  onRetry?: () => void
}

export function LiveSummaryCard({
  preview,
  isLoading,
  isStale,
  errorMessage,
  onRetry,
}: LiveSummaryCardProps) {
  const { t: tOrders } = useTranslation('orders')

  const currency = preview?.currency || 'USD'
  const showFreeShipping = preview && parseFloat(preview.shippingFee) === 0

  const formatAmount = (value?: string) => {
    if (!value) return '—'
    const parsed = Number(value)
    if (Number.isNaN(parsed)) return '—'
    return formatCurrency(parsed, currency)
  }

  return (
    <div className="rounded-xl border border-base-300 bg-base-100 p-4 space-y-3">
      <div className="flex items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <Sparkles size={16} className="text-primary" />
          <h3 className="font-semibold">{tOrders('live_summary')}</h3>
        </div>
        {isStale && (
          <span className="badge badge-ghost badge-sm text-xs">
            {tOrders('preview_status_stale')}
          </span>
        )}
      </div>

      {errorMessage && (
        <div className="alert alert-error gap-2">
          <AlertCircle size={16} />
          <span className="text-sm flex-1">{errorMessage}</span>
          {onRetry && (
            <button
              type="button"
              className="btn btn-ghost btn-sm btn-circle"
              onClick={onRetry}
            >
              <RefreshCw size={14} />
            </button>
          )}
        </div>
      )}

      {isLoading && !preview ? (
        <div className="space-y-2 animate-pulse">
          <div className="h-4 w-32 rounded bg-base-300" />
          <div className="h-4 w-28 rounded bg-base-300" />
          <div className="h-4 w-24 rounded bg-base-300" />
          <div className="h-6 w-36 rounded bg-base-300" />
        </div>
      ) : preview ? (
        <div className="space-y-2 text-sm">
          <SummaryRow
            label={tOrders('subtotal')}
            value={formatAmount(preview.subtotal)}
          />
          <SummaryRow
            label={tOrders('discount')}
            value={formatAmount(preview.discount)}
            emphasize="negative"
          />
          <SummaryRow
            label={tOrders('shipping_fee')}
            value={formatAmount(preview.shippingFee)}
          />
          <SummaryRow
            label={tOrders('vat')}
            value={formatAmount(preview.vat)}
          />
          <div className="divider my-2" />
          <SummaryRow
            label={tOrders('total')}
            value={formatAmount(preview.total)}
            emphasize="positive"
          />
          {showFreeShipping && (
            <div className="text-xs text-success">
              {tOrders('free_shipping')}
            </div>
          )}
        </div>
      ) : (
        <div className="text-sm text-base-content/70">
          {tOrders('preview_missing_requirements')}
        </div>
      )}
    </div>
  )
}

function SummaryRow({
  label,
  value,
  emphasize,
}: {
  label: string
  value: string
  emphasize?: 'positive' | 'negative'
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-base-content/70">{label}</span>
      <span
        className={`font-semibold ${
          emphasize === 'positive'
            ? 'text-primary'
            : emphasize === 'negative'
              ? 'text-error'
              : ''
        }`}
      >
        {emphasize === 'negative' ? `-${value}` : value}
      </span>
    </div>
  )
}
