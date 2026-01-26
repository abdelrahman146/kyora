import { createFileRoute } from '@tanstack/react-router'
import { ChartsDemo } from '@/features/charts-demo/components/ChartsDemo'

/**
 * Charts Demo Route
 *
 * Showcases all chart variants with the revamped design system:
 * - Bar charts (vertical, horizontal, stacked, rounded)
 * - Line charts (simple, area, multi-series, smooth)
 * - Gauge charts with trend indicators
 * - Stat cards with sparklines
 * - Donut/pie charts
 * - Mixed charts
 */
export const Route = createFileRoute('/demo/charts')({
  component: ChartsDemo,
})
