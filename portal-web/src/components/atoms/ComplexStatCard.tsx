import { Badge } from './Badge'
import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

export interface SecondaryMetric {
  label: string
  value: string | number
}

export interface ComplexStatCardProps {
  label: string
  value: string | number
  icon?: ReactNode
  secondaryMetrics?: Array<SecondaryMetric>
  comparisonText?: string
  statusBadge?: {
    label: string
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info'
  }
  variant?: 'default' | 'success' | 'warning' | 'error' | 'info'
  className?: string
}

const variantClasses = {
  default: 'bg-base-100 border-base-300',
  success: 'bg-success/5 border-success/20',
  warning: 'bg-warning/5 border-warning/20',
  error: 'bg-error/5 border-error/20',
  info: 'bg-info/5 border-info/20',
}

/**
 * ComplexStatCard Component
 *
 * Advanced statistics card with primary value, secondary metrics,
 * comparison text, status badge, and semantic variants for KPI thresholds.
 */
export const ComplexStatCard = ({
  label,
  value,
  icon,
  secondaryMetrics = [],
  comparisonText,
  statusBadge,
  variant = 'default',
  className,
}: ComplexStatCardProps) => {
  return (
    <div
      className={cn(
        'card rounded-box border transition-all ',
        variantClasses[variant],
        className,
      )}
    >
      <div className="card-body p-4">
        <div className="flex flex-col gap-4">
          {/* Primary Metric */}
          <div className="flex items-start justify-between gap-4">
            <div className="flex flex-col gap-1 flex-1 min-w-0">
              <p className="text-sm text-base-content/60 font-medium truncate">
                {label}
              </p>
              <p className="text-3xl font-bold text-base-content tabular-nums">
                {value}
              </p>
              {comparisonText && (
                <p className="text-xs text-base-content/50 mt-1">
                  {comparisonText}
                </p>
              )}
            </div>
            {statusBadge && (
              <Badge
                variant={(statusBadge.variant as any) || 'primary'}
                size="sm"
              >
                {statusBadge.label}
              </Badge>
            )}
            {icon && !statusBadge && (
              <div className="flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-base-200/50">
                {icon}
              </div>
            )}
          </div>

          {/* Secondary Metrics */}
          {secondaryMetrics.length > 0 && (
            <div className="flex items-center gap-6 pt-3 border-t border-base-300/50">
              {secondaryMetrics.map((metric, index) => (
                <div key={index} className="flex flex-col gap-1 flex-1 min-w-0">
                  <p className="text-xs text-base-content/50 truncate">
                    {metric.label}
                  </p>
                  <p className="text-base font-semibold text-base-content tabular-nums">
                    {metric.value}
                  </p>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
