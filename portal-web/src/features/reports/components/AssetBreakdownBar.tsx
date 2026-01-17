/**
 * AssetBreakdownBar Component
 *
 * Horizontal stacked bar showing asset distribution with legend.
 * Used on Business Health page to visualize "What You Own" breakdown.
 *
 * Features:
 * - Proportional segments based on value
 * - Color-coded categories (cash, inventory, fixed assets)
 * - Accessible with aria-labels for each segment
 * - Optional legend showing values and labels
 *
 * @example
 * ```tsx
 * <AssetBreakdownBar
 *   segments={[
 *     { label: "Cash", value: 8000, color: "success" },
 *     { label: "Inventory", value: 15000, color: "info" },
 *     { label: "Equipment", value: 2000, color: "secondary" },
 *   ]}
 *   total={25000}
 *   currency="USD"
 * />
 * ```
 */
import { cn } from '@/lib/utils'
import { formatCurrency } from '@/lib/formatCurrency'

export interface AssetSegment {
  /** Display label for the segment */
  label: string
  /** Monetary value (numeric) */
  value: number
  /** Color variant for visual distinction */
  color: 'success' | 'info' | 'secondary' | 'warning' | 'primary'
}

export interface AssetBreakdownBarProps {
  /** Array of asset segments to display */
  segments: Array<AssetSegment>
  /** Total value for calculating proportions */
  total: number
  /** Currency code for formatting */
  currency: string
  /** Whether to show the legend below the bar (default: true) */
  showLegend?: boolean
  /** Additional CSS classes */
  className?: string
}

const colorClasses = {
  success: { bg: 'bg-success', dot: 'bg-success' },
  info: { bg: 'bg-info', dot: 'bg-info' },
  secondary: { bg: 'bg-secondary', dot: 'bg-secondary' },
  warning: { bg: 'bg-warning', dot: 'bg-warning' },
  primary: { bg: 'bg-primary', dot: 'bg-primary' },
} as const

export function AssetBreakdownBar({
  segments,
  total,
  currency,
  showLegend = true,
  className,
}: AssetBreakdownBarProps) {
  // Filter out zero-value segments for the bar
  const nonZeroSegments = segments.filter((s) => s.value > 0)

  return (
    <div className={cn('space-y-3', className)}>
      {/* Stacked Bar */}
      <div
        className="h-4 w-full rounded-full bg-base-200 overflow-hidden flex"
        role="img"
        aria-label={`Asset breakdown: ${segments.map((s) => `${s.label}: ${formatCurrency(s.value, currency)}`).join(', ')}`}
      >
        {nonZeroSegments.map((segment, index) => {
          const percentage = total > 0 ? (segment.value / total) * 100 : 0
          return (
            <div
              key={index}
              className={cn(
                colorClasses[segment.color].bg,
                'h-full transition-all',
              )}
              style={{ width: `${percentage}%` }}
              aria-hidden="true"
            />
          )
        })}
      </div>

      {/* Legend */}
      {showLegend && (
        <div className="space-y-2">
          {segments.map((segment, index) => (
            <div
              key={index}
              className="flex items-center justify-between text-sm"
            >
              <div className="flex items-center gap-2">
                <span
                  className={cn(
                    'h-3 w-3 rounded-full shrink-0',
                    colorClasses[segment.color].dot,
                  )}
                  aria-hidden="true"
                />
                <span className="text-base-content/80">{segment.label}</span>
              </div>
              <span className="font-medium tabular-nums">
                {formatCurrency(segment.value, currency)}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

AssetBreakdownBar.displayName = 'AssetBreakdownBar'
