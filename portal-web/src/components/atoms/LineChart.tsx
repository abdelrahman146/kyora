import { useMemo } from 'react'
import { Line } from 'react-chartjs-2'
import type { ChartData, ChartOptions } from 'chart.js'
import { useChartTheme } from '@/lib/charts'
import { cn } from '@/lib/utils'

export interface LineChartProps {
  data: ChartData<'line'>
  options?: ChartOptions<'line'>
  height?: number
  className?: string
  enableArea?: boolean
  enableDecimation?: boolean
}

/**
 * LineChart Component
 *
 * Generic line chart with automatic theming from daisyUI and RTL support.
 * Supports area fill and data decimation for large datasets.
 */
export function LineChart({
  data,
  options = {},
  height = 320,
  className,
  enableArea = false,
  enableDecimation = false,
}: LineChartProps) {
  const { themedOptions, backgroundPlugin } = useChartTheme()

  const mergedOptions = useMemo<ChartOptions<'line'>>(() => {
    const decimationConfig = enableDecimation
      ? {
          plugins: {
            decimation: {
              enabled: true,
              algorithm: 'min-max' as const,
            },
          },
        }
      : {}

    return {
      ...themedOptions,
      ...decimationConfig,
      ...options,
      plugins: {
        ...themedOptions.plugins,
        ...decimationConfig.plugins,
        ...options.plugins,
      },
    } as ChartOptions<'line'>
  }, [themedOptions, options, enableDecimation])

  const processedData = useMemo(() => {
    if (!enableArea) return data

    return {
      ...data,
      datasets: data.datasets.map((dataset) => ({
        ...dataset,
        fill: true,
        tension: 0.4,
      })),
    }
  }, [data, enableArea])

  return (
    <div className={cn('relative h-full w-full', className)} style={{ height }}>
      <Line
        data={processedData}
        options={mergedOptions}
        plugins={[backgroundPlugin] as any}
      />
    </div>
  )
}
