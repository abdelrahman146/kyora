import { useMemo } from 'react'
import { Bar } from 'react-chartjs-2'
import type { ChartData, ChartOptions } from 'chart.js'
import { useChartTheme } from '@/lib/charts'
import { cn } from '@/lib/utils'

export interface BarChartProps {
  data: ChartData<'bar'>
  options?: ChartOptions<'bar'>
  height?: number
  className?: string
  horizontal?: boolean
  stacked?: boolean
}

/**
 * BarChart Component
 *
 * Generic bar chart with automatic theming from daisyUI and RTL support.
 * Supports horizontal orientation and stacked mode.
 */
export function BarChart({
  data,
  options = {},
  height = 320,
  className,
  horizontal = false,
  stacked = false,
}: BarChartProps) {
  const { themedOptions, backgroundPlugin } = useChartTheme()

  const mergedOptions = useMemo<ChartOptions<'bar'>>(() => {
    const indexAxis = horizontal ? ('y' as const) : ('x' as const)

    const stackedConfig = stacked
      ? {
          scales: {
            x: { stacked: true },
            y: { stacked: true },
          },
        }
      : {}

    return {
      ...themedOptions,
      indexAxis,
      ...stackedConfig,
      ...options,
      scales: {
        ...themedOptions.scales,
        ...stackedConfig.scales,
        ...options.scales,
      },
      plugins: {
        ...themedOptions.plugins,
        ...options.plugins,
      },
    } as ChartOptions<'bar'>
  }, [themedOptions, options, horizontal, stacked])

  return (
    <div className={cn('relative h-full w-full', className)} style={{ height }}>
      <Bar
        data={data}
        options={mergedOptions}
        plugins={[backgroundPlugin] as any}
      />
    </div>
  )
}
