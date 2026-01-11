import { useMemo } from 'react'
import { Chart } from 'react-chartjs-2'
import type { ChartData, ChartOptions } from 'chart.js'
import { useChartTheme } from '@/lib/charts'
import { cn } from '@/lib/utils'

export interface MixedChartProps {
  data: ChartData<'bar' | 'line'>
  options?: ChartOptions
  height?: number
  className?: string
}

/**
 * MixedChart Component
 *
 * Combined bar + line chart with automatic theming from daisyUI and RTL support.
 * Perfect for P&L statements and comparative analytics (e.g., revenue vs profit).
 * Supports dual Y-axes for different scales.
 */
export function MixedChart({
  data,
  options = {},
  height = 320,
  className,
}: MixedChartProps) {
  const { themedOptions, backgroundPlugin } = useChartTheme()

  const mergedOptions = useMemo<ChartOptions>(() => {
    return {
      ...themedOptions,
      ...options,
      plugins: {
        ...themedOptions.plugins,
        ...options.plugins,
      },
      scales: {
        ...themedOptions.scales,
        ...options.scales,
      },
    }
  }, [themedOptions, options])

  return (
    <div className={cn('relative h-full w-full', className)} style={{ height }}>
      <Chart
        type="bar"
        data={data}
        options={mergedOptions}
        plugins={[backgroundPlugin] as any}
      />
    </div>
  )
}
