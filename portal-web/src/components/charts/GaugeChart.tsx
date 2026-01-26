import { useMemo } from 'react'
import { TrendingDown, TrendingUp } from 'lucide-react'
import { Doughnut } from 'react-chartjs-2'
import type { ChartData, ChartOptions } from 'chart.js'
import {
  colorWithOpacity,
  createGaugeCenterLabelPlugin,
  useChartTheme,
} from '@/lib/charts'
import { cn } from '@/lib/utils'

export interface GaugeChartProps {
  value: number
  max?: number
  label: string
  trend?: {
    value: number
    direction: 'up' | 'down'
  }
  className?: string
  color?: string
}

/**
 * GaugeChart Component
 *
 * Progress circle with center label and optional trend indicator.
 * Design specs: 180px × 180px, semantic colors only, no gradients.
 */
export function GaugeChart({
  value,
  max = 100,
  label,
  trend,
  className,
  color,
}: GaugeChartProps) {
  const { themedOptions, backgroundPlugin, tokens } = useChartTheme()

  const percentage = Math.min((value / max) * 100, 100)
  const displayValue = `${Math.round(percentage)}%`

  const primaryColor = color || tokens.primary
  const emptyColor = colorWithOpacity(tokens.base300, 0.3)

  const data = useMemo<ChartData<'doughnut'>>(() => {
    return {
      labels: [label, ''],
      datasets: [
        {
          data: [percentage, 100 - percentage],
          backgroundColor: [primaryColor, emptyColor],
          borderWidth: 0,
          circumference: 360,
          rotation: -90,
        },
      ],
    }
  }, [percentage, primaryColor, emptyColor, label])

  const options = useMemo<ChartOptions<'doughnut'>>(() => {
    return {
      ...themedOptions,
      plugins: {
        ...themedOptions.plugins,
        legend: {
          display: false,
        },
        tooltip: {
          enabled: false,
        },
      },
      scales: {
        x: {
          display: false,
        },
        y: {
          display: false,
        },
      },
      cutout: '75%',
    } as ChartOptions<'doughnut'>
  }, [themedOptions])

  const centerLabelPlugin = useMemo(() => {
    return createGaugeCenterLabelPlugin(
      displayValue,
      label,
      tokens.baseContent,
      colorWithOpacity(tokens.baseContent, 0.6),
    )
  }, [displayValue, label, tokens.baseContent])

  return (
    <div className={cn('relative flex flex-col items-center gap-2', className)}>
      <div className="relative h-[180px] w-[180px]">
        <Doughnut
          data={data}
          options={options}
          plugins={[backgroundPlugin, centerLabelPlugin] as any}
        />
      </div>
      {trend && (
        <div
          className={cn(
            'flex items-center gap-1 text-sm font-medium',
            trend.direction === 'up' ? 'text-success' : 'text-error',
          )}
        >
          {trend.direction === 'up' ? (
            <TrendingUp className="h-4 w-4" />
          ) : (
            <TrendingDown className="h-4 w-4" />
          )}
          <span>
            {trend.value > 0 ? '+' : ''}
            {trend.value}%
          </span>
        </div>
      )}
    </div>
  )
}
