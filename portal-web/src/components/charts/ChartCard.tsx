import { useTranslation } from 'react-i18next'
import { RefreshCw } from 'lucide-react'
import { ChartSkeleton } from './ChartSkeleton'
import { ChartEmptyState } from './ChartEmptyState'
import type { ReactNode } from 'react'
import { Button } from '@/components/atoms/Button'
import { cn } from '@/lib/utils'

export interface ChartCardProps {
  title?: string
  subtitle?: string
  children: ReactNode
  isLoading?: boolean
  isEmpty?: boolean
  error?: Error | null
  onRetry?: () => void
  actions?: ReactNode
  className?: string
  height?: number
  chartType?: 'line' | 'bar' | 'pie' | 'doughnut' | 'mixed'
}

export const ChartCard = ({
  title,
  subtitle,
  children,
  isLoading = false,
  isEmpty = false,
  error = null,
  onRetry,
  actions,
  className,
  height = 320,
  chartType = 'line',
}: ChartCardProps) => {
  const { t } = useTranslation('analytics')

  return (
    <div
      className={cn(
        'card rounded-box bg-base-100 border border-base-300',
        className,
      )}
    >
      {/* Card Header */}
      {(title || subtitle || actions) && (
        <div className="card-body p-4 pb-2">
          <div className="flex items-start justify-between gap-4">
            <div className="flex flex-col gap-1">
              {title && (
                <h3 className="text-lg font-semibold text-base-content">
                  {title}
                </h3>
              )}
              {subtitle && (
                <p className="text-sm text-base-content/60">{subtitle}</p>
              )}
            </div>
            {actions && (
              <div className="flex items-center gap-2">{actions}</div>
            )}
          </div>
        </div>
      )}

      {/* Card Content */}
      <div className="card-body p-4 pt-2">
        {isLoading ? (
          <ChartSkeleton variant={chartType} height={height} />
        ) : error ? (
          <div
            className="flex min-h-[280px] w-full flex-col items-center justify-center gap-4 rounded-lg bg-error/5 p-8 text-center"
            style={{ height }}
          >
            <div className="flex flex-col gap-2">
              <p className="text-sm font-medium text-error">
                {t('error.title', 'Failed to load chart')}
              </p>
              <p className="text-xs text-base-content/60">
                {error.message || t('error.subtitle', 'Something went wrong')}
              </p>
            </div>
            {onRetry && (
              <Button variant="outline" size="sm" onClick={onRetry}>
                <RefreshCw className="h-4 w-4" />
                {t('error.retry', 'Retry')}
              </Button>
            )}
          </div>
        ) : isEmpty ? (
          <ChartEmptyState chartType={chartType} />
        ) : (
          <div className="relative" style={{ height }}>
            {children}
          </div>
        )}
      </div>
    </div>
  )
}
