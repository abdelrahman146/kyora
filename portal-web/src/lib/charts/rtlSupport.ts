import { useTranslation } from 'react-i18next'
import type { ChartOptions } from 'chart.js'

/**
 * Detects if the current language direction is RTL
 */
export function isRTL(): boolean {
  return document.documentElement.dir === 'rtl'
}

/**
 * Hook to get the current language direction
 */
export function useRTL(): boolean {
  const { i18n } = useTranslation()
  return i18n.dir() === 'rtl'
}

/**
 * Applies RTL transformations to Chart.js options
 */
export function applyRTLTransformations(
  options: ChartOptions,
  rtl: boolean,
): ChartOptions {
  if (!rtl) return options

  const rtlOptions: ChartOptions = {
    ...options,
    plugins: {
      ...options.plugins,
      legend: {
        ...options.plugins?.legend,
        align: 'end',
        rtl: true,
        textDirection: 'rtl',
      },
      tooltip: {
        ...options.plugins?.tooltip,
        rtl: true,
        textDirection: 'rtl',
      },
    },
  }

  // Reverse x-axis for RTL
  if (rtlOptions.scales?.x) {
    rtlOptions.scales.x = {
      ...rtlOptions.scales.x,
      reverse: true,
    }
  }

  return rtlOptions
}

/**
 * Gets the position for chart elements based on language direction
 */
export function getDirectionalPosition(rtl: boolean): 'left' | 'right' {
  return rtl ? 'right' : 'left'
}

/**
 * Gets the alignment for chart elements based on language direction
 */
export function getDirectionalAlignment(rtl: boolean): 'start' | 'end' {
  return rtl ? 'end' : 'start'
}

/**
 * Mirrors tooltip position for RTL
 */
export function mirrorTooltipPosition(
  position: 'left' | 'right' | 'top' | 'bottom',
  rtl: boolean,
): 'left' | 'right' | 'top' | 'bottom' {
  if (!rtl) return position

  if (position === 'left') return 'right'
  if (position === 'right') return 'left'
  return position
}
