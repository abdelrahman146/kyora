import type { ChartData, ChartDataset } from 'chart.js'

/**
 * TimeSeries data structure from backend
 */
export interface TimeSeriesRow {
  timestamp: string
  label: string
  value: number
}

export interface TimeSeries {
  granularity:
    | 'hourly'
    | 'daily'
    | 'weekly'
    | 'monthly'
    | 'quarterly'
    | 'yearly'
  series: Array<TimeSeriesRow>
}

/**
 * KeyValue data structure from backend
 */
export interface KeyValue {
  key: string | number
  value: string | number
}

/**
 * Transforms TimeSeries data from backend into Chart.js compatible format
 */
export function transformTimeSeriesToChartData(
  timeSeries: TimeSeries,
  datasetLabel: string,
  color: string,
): ChartData<'line'> {
  return {
    labels: timeSeries.series.map((row) => row.label),
    datasets: [
      {
        label: datasetLabel,
        data: timeSeries.series.map((row) => row.value),
        borderColor: color,
        backgroundColor: color.replace(')', ', 0.1)').replace('rgb', 'rgba'),
        tension: 0.4,
        fill: true,
      },
    ],
  }
}

/**
 * Transforms KeyValue array into Chart.js bar chart data
 */
export function transformKeyValueToBarData(
  keyValues: Array<KeyValue>,
  datasetLabel: string,
  colors: Array<string>,
): ChartData<'bar'> {
  return {
    labels: keyValues.map((kv) => String(kv.key)),
    datasets: [
      {
        label: datasetLabel,
        data: keyValues.map((kv) => Number(kv.value)),
        backgroundColor: colors,
        borderWidth: 0,
      },
    ],
  }
}

/**
 * Transforms KeyValue array into Chart.js pie/doughnut chart data
 */
export function transformKeyValueToPieData(
  keyValues: Array<KeyValue>,
  colors: Array<string>,
): ChartData<'pie' | 'doughnut'> {
  return {
    labels: keyValues.map((kv) => String(kv.key)),
    datasets: [
      {
        data: keyValues.map((kv) => Number(kv.value)),
        backgroundColor: colors,
        borderWidth: 2,
        borderColor: '#ffffff',
      },
    ],
  }
}

/**
 * Generates a color palette from a base color with varying opacity
 */
export function generateColorPalette(
  baseColor: string,
  count: number,
): Array<string> {
  const colors: Array<string> = []
  const opacities = [1, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3]

  for (let i = 0; i < count; i++) {
    const opacity = opacities[i % opacities.length]
    colors.push(baseColor.replace(')', `, ${opacity})`).replace('rgb', 'rgba'))
  }

  return colors
}

/**
 * Formats currency values for chart labels and tooltips
 */
export function formatChartCurrency(value: number, currency: string): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(value)
}

/**
 * Formats large numbers with K/M/B suffixes for better readability
 */
export function formatChartNumber(value: number): string {
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(1)}B`
  }
  if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(1)}M`
  }
  if (value >= 1_000) {
    return `${(value / 1_000).toFixed(1)}K`
  }
  return value.toFixed(0)
}

/**
 * Determines if decimation should be enabled based on data point count
 */
export function shouldEnableDecimation(dataPointCount: number): boolean {
  return dataPointCount > 100
}

/**
 * Gets optimal animation configuration based on dataset size
 */
export function getOptimalAnimationConfig(
  dataPointCount: number,
): boolean | { duration: number } {
  if (dataPointCount > 500) {
    return false // Disable animations for very large datasets
  }
  if (dataPointCount > 200) {
    return { duration: 300 } // Faster animations for large datasets
  }
  return { duration: 750 } // Default smooth animations
}

/**
 * Merges multiple chart datasets for mixed charts
 */
export function mergeChartDatasets<T extends 'line' | 'bar'>(
  ...datasets: Array<ChartDataset<T>>
): Array<ChartDataset<T>> {
  return datasets
}

/**
 * Calculates percentage for pie chart labels
 */
export function calculatePercentages(values: Array<number>): Array<string> {
  const total = values.reduce((sum, val) => sum + val, 0)
  return values.map((val) => `${((val / total) * 100).toFixed(1)}%`)
}
