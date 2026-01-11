import {
  Activity,
  BarChart3,
  LineChart,
  PieChart as PieChartIcon,
  TrendingUp,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { LucideIcon } from 'lucide-react'
import { cn } from '@/lib/utils'

export interface ChartEmptyStateProps {
  chartType?: 'line' | 'bar' | 'pie' | 'doughnut' | 'mixed'
  message?: string
  className?: string
}

const chartIcons: Record<string, LucideIcon> = {
  line: LineChart,
  bar: BarChart3,
  pie: PieChartIcon,
  doughnut: PieChartIcon,
  mixed: Activity,
}

export const ChartEmptyState = ({
  chartType = 'line',
  message,
  className,
}: ChartEmptyStateProps) => {
  const { t } = useTranslation('analytics')
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  const Icon = chartIcons[chartType] || TrendingUp

  return (
    <div
      className={cn(
        'flex min-h-[280px] w-full flex-col items-center justify-center gap-3 rounded-lg bg-base-100 p-8 text-center',
        className,
      )}
      role="status"
      aria-label={message || t('emptyState.noData', 'No data available')}
    >
      <Icon className="h-12 w-12 text-base-content/30" strokeWidth={1.5} />
      <div className="flex flex-col gap-1">
        <p className="text-sm font-medium text-base-content/70">
          {message || t('emptyState.noData', 'No data available')}
        </p>
        <p className="text-xs text-base-content/50">
          {t('emptyState.subtitle', 'Data will appear here once available')}
        </p>
      </div>
    </div>
  )
}
