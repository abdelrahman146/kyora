import { TrendingDown, TrendingUp } from 'lucide-react'
import { SparklineChart } from './SparklineChart'
import { cn } from '@/lib/utils'

export interface StatCardProps {
  title: string
  value: string | number
  trend?: {
    value: number
    direction: 'up' | 'down'
  }
  sparklineData?: Array<number>
  sparklineColor?: string
  className?: string
  subtitle?: string
}

/**
 * StatCard Component
 *
 * Metric card with inline sparkline showing trend.
 * Design specs: Header with metric, optional trend indicator, 48px sparkline height.
 */
export function StatCard({
  title,
  value,
  trend,
  sparklineData,
  sparklineColor,
  className,
  subtitle,
}: StatCardProps) {
  return (
    <div
      className={cn(
        'card rounded-lg border border-base-300 bg-base-100 p-5',
        className,
      )}
    >
      <div className="flex flex-col gap-4">
        {/* Header */}
        <div className="flex items-start justify-between gap-3">
          <div className="flex flex-col gap-1.5">
            <p className="text-sm font-medium text-base-content/60">{title}</p>
            <p className="text-3xl font-bold text-base-content">{value}</p>
            {subtitle && (
              <p className="text-xs text-base-content/50">{subtitle}</p>
            )}
          </div>
          {trend && (
            <div
              className={cn(
                'flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-sm font-semibold',
                trend.direction === 'up'
                  ? 'bg-success/10 text-success'
                  : 'bg-error/10 text-error',
              )}
            >
              {trend.direction === 'up' ? (
                <TrendingUp className="h-3.5 w-3.5" />
              ) : (
                <TrendingDown className="h-3.5 w-3.5" />
              )}
              <span>
                {trend.value > 0 ? '+' : ''}
                {trend.value}%
              </span>
            </div>
          )}
        </div>

        {/* Sparkline */}
        {sparklineData && sparklineData.length > 0 && (
          <div className="h-16 w-full">
            <SparklineChart
              data={sparklineData}
              color={sparklineColor}
              height={64}
            />
          </div>
        )}
      </div>
    </div>
  )
}
