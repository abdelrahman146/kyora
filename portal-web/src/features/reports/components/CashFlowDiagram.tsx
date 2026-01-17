/**
 * CashFlowDiagram Component
 *
 * Visual flow representation for cash movement:
 * Start → Money In box → Money Out box → Net Change → End
 *
 * Shows a vertical flow with color-coded boxes for inflows (green) and
 * outflows (orange/red). Arrows connect each section.
 *
 * Mobile-first, RTL-compatible layout.
 */
import { ArrowDown, TrendingDown, TrendingUp } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { formatCurrency } from '@/lib/formatCurrency'
import { cn } from '@/lib/utils'

export interface CashFlowItem {
  /** Display label for this item */
  label: string
  /** Numeric value (positive) */
  value: number
}

export interface CashFlowDiagramProps {
  /** Starting cash amount */
  startAmount: number
  /** Ending cash amount */
  endAmount: number
  /** Array of cash inflow items */
  cashIn: Array<CashFlowItem>
  /** Array of cash outflow items */
  cashOut: Array<CashFlowItem>
  /** Currency code for formatting */
  currency: string
  /** Additional CSS classes */
  className?: string
}

export function CashFlowDiagram({
  startAmount,
  endAmount,
  cashIn,
  cashOut,
  currency,
  className,
}: CashFlowDiagramProps) {
  const { t } = useTranslation('reports')

  const totalIn = cashIn.reduce((sum, item) => sum + item.value, 0)
  const totalOut = cashOut.reduce((sum, item) => sum + item.value, 0)
  const netChange = totalIn - totalOut
  const isPositiveChange = netChange >= 0

  return (
    <div className={cn('space-y-4', className)}>
      {/* Start Node */}
      <div className="text-center">
        <span className="text-sm text-base-content/60">
          {t('cashflow.cash_start')}
        </span>
        <p className="text-lg font-semibold tabular-nums">
          {formatCurrency(startAmount, currency)}
        </p>
      </div>

      <div className="flex justify-center">
        <ArrowDown
          className="h-5 w-5 text-base-content/30"
          aria-hidden="true"
        />
      </div>

      {/* Money Coming In Box */}
      <div className="rounded-box border border-success/20 bg-success/5 p-4">
        <h3 className="mb-3 font-semibold text-success">
          {t('cashflow.money_in')}
        </h3>
        <div className="space-y-2">
          {cashIn.map((item, index) => (
            <div key={index} className="flex justify-between text-sm">
              <span className="text-base-content/80">{item.label}</span>
              <span className="font-medium tabular-nums text-success">
                +{formatCurrency(item.value, currency)}
              </span>
            </div>
          ))}
          <div className="mt-2 flex justify-between border-t border-success/20 pt-2 font-semibold">
            <span>{t('cashflow.total_in')}</span>
            <span className="tabular-nums text-success">
              {formatCurrency(totalIn, currency)}
            </span>
          </div>
        </div>
      </div>

      <div className="flex justify-center">
        <ArrowDown
          className="h-5 w-5 text-base-content/30"
          aria-hidden="true"
        />
      </div>

      {/* Money Going Out Box */}
      <div className="rounded-box border border-warning/20 bg-warning/5 p-4">
        <h3 className="mb-3 font-semibold text-warning">
          {t('cashflow.money_out')}
        </h3>
        <div className="space-y-2">
          {cashOut.map((item, index) => (
            <div key={index} className="flex justify-between text-sm">
              <span className="text-base-content/80">{item.label}</span>
              <span className="font-medium tabular-nums text-error">
                -{formatCurrency(item.value, currency)}
              </span>
            </div>
          ))}
          <div className="mt-2 flex justify-between border-t border-warning/20 pt-2 font-semibold">
            <span>{t('cashflow.total_out')}</span>
            <span className="tabular-nums text-error">
              {formatCurrency(totalOut, currency)}
            </span>
          </div>
        </div>
      </div>

      <div className="flex justify-center">
        <ArrowDown
          className="h-5 w-5 text-base-content/30"
          aria-hidden="true"
        />
      </div>

      {/* Net Change Indicator */}
      <div className="py-2 text-center">
        <div className="flex items-center justify-center gap-2">
          {isPositiveChange ? (
            <TrendingUp
              className="h-5 w-5 text-success"
              aria-label={t('cashflow.cash_increased', {
                amount: formatCurrency(netChange, currency),
              })}
            />
          ) : (
            <TrendingDown
              className="h-5 w-5 text-error"
              aria-label={t('cashflow.cash_decreased', {
                amount: formatCurrency(Math.abs(netChange), currency),
              })}
            />
          )}
          <span className="text-sm text-base-content/60">
            {t('cashflow.net_change')}
          </span>
        </div>
        <p
          className={cn(
            'text-xl font-bold tabular-nums',
            isPositiveChange ? 'text-success' : 'text-error',
          )}
        >
          {isPositiveChange ? '+' : ''}
          {formatCurrency(netChange, currency)}
        </p>
      </div>

      <div className="flex justify-center">
        <ArrowDown
          className="h-5 w-5 text-base-content/30"
          aria-hidden="true"
        />
      </div>

      {/* End Node */}
      <div className="rounded-box bg-base-200/50 p-4 text-center">
        <span className="text-sm text-base-content/60">
          {t('cashflow.cash_end')}
        </span>
        <p
          className={cn(
            'text-2xl font-bold tabular-nums',
            endAmount >= 0 ? 'text-success' : 'text-error',
          )}
        >
          {formatCurrency(endAmount, currency)}
        </p>
      </div>
    </div>
  )
}

CashFlowDiagram.displayName = 'CashFlowDiagram'
