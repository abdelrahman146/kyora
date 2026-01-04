import { TrendingDown, TrendingUp } from 'lucide-react'
import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

export interface StatCardProps {
  label: string
  value: string | number
  icon?: ReactNode
  trend?: 'up' | 'down'
  trendValue?: string
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

const trendClasses = {
  up: 'text-success',
  down: 'text-error',
}

/**
 * StatCard Component
 *
 * Simple statistics card with value, label, optional icon, trend indicator,
 * and semantic variants for KPI thresholds (success, warning, error).
 */
export const StatCard = ({
  label,
  value,
  icon,
  trend,
  trendValue,
  variant = 'default',
  className,
}: StatCardProps) => {
  const TrendIcon = trend === 'up' ? TrendingUp : TrendingDown

  return (
    <div
      className={cn(
        'card rounded-box border transition-all ',
        variantClasses[variant],
        className,
      )}
    >
      <div className="card-body p-4 min-h-[52px]">
        <div className="flex items-start justify-between gap-4">
          <div className="flex flex-col gap-1 flex-1 min-w-0">
            <p className="text-sm text-base-content/60 font-medium truncate">
              {label}
            </p>
            <p className="text-2xl font-bold text-base-content tabular-nums">
              {value}
            </p>
            {trend && trendValue && (
              <div
                className={cn(
                  'flex items-center gap-1 text-xs font-medium',
                  trendClasses[trend],
                )}
              >
                <TrendIcon className="h-3 w-3" />
                <span>{trendValue}</span>
              </div>
            )}
          </div>
          {icon && (
            <div className="flex-shrink-0 flex items-center justify-center h-10 w-10 rounded-full bg-base-200/50">
              {icon}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
