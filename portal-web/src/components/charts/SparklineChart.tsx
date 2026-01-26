import { useMemo } from 'react'
import { Line } from 'react-chartjs-2'
import type { ChartData, ChartOptions } from 'chart.js'
import {
  AREA_FILL_OPACITY_SPARKLINE,
  colorWithOpacity,
  useChartTheme,
} from '@/lib/charts'
import { cn } from '@/lib/utils'

export interface SparklineChartProps {
  data: Array<number>
  color?: string
  className?: string
  height?: number
}

/**
 * SparklineChart Component
 *
 * Ultra-compact trend indicator for inline display.
 * Design specs: 60px × 24px, minimal styling, no axes/labels/grid.
 */
export function SparklineChart({
  data,
  color,
  className,
  height = 24,
}: SparklineChartProps) {
  const { tokens } = useChartTheme()

  const lineColor = color || tokens.primary

  const chartData = useMemo<ChartData<'line'>>(() => {
    return {
      labels: data.map((_, i) => i.toString()),
      datasets: [
        {
          data,
          borderColor: lineColor,
          backgroundColor: colorWithOpacity(
            lineColor,
            AREA_FILL_OPACITY_SPARKLINE,
          ),
          borderWidth: 2,
          fill: true,
          tension: 0.4,
          pointRadius: 0,
          pointHoverRadius: 0,
        },
      ],
    }
  }, [data, lineColor])

  const options = useMemo<ChartOptions<'line'>>(() => {
    return {
      responsive: true,
      maintainAspectRatio: false,
      interaction: {
        intersect: false,
      },
      scales: {
        x: {
          display: false,
        },
        y: {
          display: false,
        },
      },
      plugins: {
        legend: {
          display: false,
        },
        tooltip: {
          enabled: false,
        },
      },
      elements: {
        line: {
          tension: 0.4,
          borderWidth: 2,
        },
        point: {
          radius: 0,
        },
      },
    } as ChartOptions<'line'>
  }, [])

  return (
    <div className={cn('relative', className)} style={{ height }}>
      <Line data={chartData} options={options} />
    </div>
  )
}
