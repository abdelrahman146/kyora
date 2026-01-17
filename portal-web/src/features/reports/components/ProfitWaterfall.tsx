/**
 * ProfitWaterfall Component
 *
 * Professional P&L flow visualization with prominent operator badges:
 * Revenue → COGS → Gross Profit → Operating Expenses → Net Profit
 *
 * Features:
 * - Large, colored operator badges (+, -, =) for clarity
 * - Color-coded bars (green for revenue/results, red for costs)
 * - Mobile-first, RTL-compatible layout
 * - Professional financial statement appearance
 */
import { Minus, Plus, Equal } from 'lucide-react'

import { formatCurrency } from '@/lib/formatCurrency'
import { cn } from '@/lib/utils'

export interface WaterfallStep {
  /** Display label for this step */
  label: string
  /** Numeric value (positive, even for subtractions) */
  value: number
  /** Type determines visual treatment: total (starting), subtract (deduction), result (subtotal/total) */
  type: 'total' | 'subtract' | 'result'
}

export interface ProfitWaterfallProps {
  /** Array of steps in the waterfall */
  steps: Array<WaterfallStep>
  /** Currency code for formatting */
  currency: string
  /** Additional CSS classes */
  className?: string
}

export function ProfitWaterfall({
  steps,
  currency,
  className,
}: ProfitWaterfallProps) {
  // Find max value for bar scaling
  const maxValue = Math.max(...steps.map((s) => Math.abs(s.value)))

  return (
    <div className={cn('space-y-4', className)}>
      {steps.map((step, index) => {
        const percentage =
          maxValue > 0 ? (Math.abs(step.value) / maxValue) * 100 : 0
        const isSubtract = step.type === 'subtract'
        const isResult = step.type === 'result'
        const isTotal = step.type === 'total'
        const isLast = index === steps.length - 1

        // Determine operator icon and colors
        let OperatorIcon = null
        let operatorBadgeClass = ''

        if (isTotal) {
          OperatorIcon = Plus
          operatorBadgeClass = 'bg-success/20 text-success'
        } else if (isSubtract) {
          OperatorIcon = Minus
          operatorBadgeClass = 'bg-error/20 text-error'
        } else if (isResult) {
          OperatorIcon = Equal
          operatorBadgeClass =
            step.value >= 0
              ? 'bg-primary/20 text-primary'
              : 'bg-error/20 text-error'
        }

        return (
          <div key={index}>
            {/* Step Row with Operator Badge */}
            <div className="flex items-start gap-4">
              {/* Large Operator Badge */}
              {OperatorIcon && (
                <div
                  className={cn(
                    'flex h-10 w-10 shrink-0 items-center justify-center rounded-lg font-bold',
                    operatorBadgeClass,
                  )}
                  aria-hidden="true"
                >
                  <OperatorIcon className="h-5 w-5" strokeWidth={2.5} />
                </div>
              )}

              {/* Label and Bar */}
              <div className="flex-1 pt-1">
                <div className="mb-2 flex items-center justify-between gap-3">
                  <span
                    className={cn(
                      'text-sm',
                      isResult && 'font-semibold text-base-content',
                      !isResult && 'text-base-content/80',
                    )}
                  >
                    {step.label}
                  </span>
                  <span
                    className={cn(
                      'font-semibold tabular-nums',
                      isSubtract && 'text-error',
                      isTotal && 'text-base-content',
                      isResult && step.value >= 0 && 'text-success',
                      isResult && step.value < 0 && 'text-error',
                    )}
                  >
                    {isSubtract && '-'}
                    {formatCurrency(Math.abs(step.value), currency)}
                  </span>
                </div>
                <div className="h-2.5 w-full overflow-hidden rounded-full bg-base-200">
                  <div
                    className={cn(
                      'h-full rounded-full transition-all duration-500',
                      isTotal && 'bg-success',
                      isSubtract && 'bg-error',
                      isResult && step.value >= 0 && 'bg-primary',
                      isResult && step.value < 0 && 'bg-error',
                    )}
                    style={{ width: `${percentage}%` }}
                  />
                </div>
              </div>
            </div>

            {/* Divider line between steps (not after last) */}
            {!isLast && (
              <div
                className="ms-5 border-s-2 border-dashed border-base-300 py-2"
                aria-hidden="true"
              />
            )}
          </div>
        )
      })}
    </div>
  )
}

ProfitWaterfall.displayName = 'ProfitWaterfall'
