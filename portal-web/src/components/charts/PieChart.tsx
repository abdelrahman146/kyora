import { useMemo } from 'react'
import { Pie } from 'react-chartjs-2'
import type { ChartData, ChartOptions } from 'chart.js'
import { useChartTheme } from '@/lib/charts'
import { cn } from '@/lib/utils'

export interface PieChartProps {
  data: ChartData<'pie'>
  options?: ChartOptions<'pie'>
  height?: number
  className?: string
}

/**
 * PieChart Component
 *
 * Generic pie chart with automatic theming from daisyUI and RTL support.
 */
export function PieChart({
  data,
  options = {},
  height = 320,
  className,
}: PieChartProps) {
  const { themedOptions, backgroundPlugin } = useChartTheme()

  const mergedOptions = useMemo<ChartOptions<'pie'>>(() => {
    return {
      ...themedOptions,
      ...options,
      plugins: {
        ...themedOptions.plugins,
        ...options.plugins,
        legend: {
          ...themedOptions.plugins?.legend,
          ...options.plugins?.legend,
          position: 'bottom',
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
    } as ChartOptions<'pie'>
  }, [themedOptions, options])

  return (
    <div className={cn('relative h-full w-full', className)} style={{ height }}>
      <Pie
        data={data}
        options={mergedOptions}
        plugins={[backgroundPlugin] as any}
      />
    </div>
  )
}
