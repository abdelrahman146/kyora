import { useMemo } from 'react'
import { Doughnut } from 'react-chartjs-2'
import type { ChartData, ChartOptions } from 'chart.js'
import { createCenterLabelPlugin, useChartTheme } from '@/lib/charts'
import { cn } from '@/lib/utils'

export interface DoughnutChartProps {
  data: ChartData<'doughnut'>
  options?: ChartOptions<'doughnut'>
  height?: number
  className?: string
  centerLabel?: string
  centerLabelColor?: string
}

/**
 * DoughnutChart Component
 *
 * Generic doughnut chart with automatic theming from daisyUI and RTL support.
 * Supports center label for displaying totals or primary metrics.
 */
export function DoughnutChart({
  data,
  options = {},
  height = 320,
  className,
  centerLabel,
  centerLabelColor,
}: DoughnutChartProps) {
  const { themedOptions, backgroundPlugin, tokens } = useChartTheme()

  const mergedOptions = useMemo<ChartOptions<'doughnut'>>(() => {
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
      cutout: '70%',
    } as ChartOptions<'doughnut'>
  }, [themedOptions, options])

  const plugins = useMemo(() => {
    const basePlugins = [backgroundPlugin]

    if (centerLabel) {
      basePlugins.push(
        createCenterLabelPlugin(
          centerLabel,
          centerLabelColor || tokens.baseContent,
        ),
      )
    }

    return basePlugins
  }, [backgroundPlugin, centerLabel, centerLabelColor, tokens.baseContent])

  return (
    <div className={cn('relative h-full w-full', className)} style={{ height }}>
      <Doughnut data={data} options={mergedOptions} plugins={plugins as any} />
    </div>
  )
}
