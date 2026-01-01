import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'
import type { ChartOptions, Plugin } from 'chart.js'

/**
 * Resolves a background color from a daisyUI className by rendering a hidden element
 * and reading its computed style. This ensures colors are always derived from the active theme.
 */
export function resolveBgColorFromClass(className: string): string {
  const el = document.createElement('div')
  el.className = className
  el.style.position = 'absolute'
  el.style.left = '-9999px'
  el.style.visibility = 'hidden'
  document.body.appendChild(el)

  const color = getComputedStyle(el).backgroundColor
  document.body.removeChild(el)
  return color
}

/**
 * Resolves a text color from a daisyUI className by rendering a hidden element
 * and reading its computed style.
 */
export function resolveTextColorFromClass(className: string): string {
  const el = document.createElement('div')
  el.className = className
  el.style.position = 'absolute'
  el.style.left = '-9999px'
  el.style.visibility = 'hidden'
  document.body.appendChild(el)

  const color = getComputedStyle(el).color
  document.body.removeChild(el)
  return color
}

/**
 * Extracts all semantic color tokens from the active daisyUI theme
 * by resolving colors from semantic class names.
 */
export interface ChartTokens {
  primary: string
  secondary: string
  accent: string
  success: string
  warning: string
  error: string
  info: string
  base100: string
  base200: string
  base300: string
  baseContent: string
  neutral: string
}

export function getChartTokens(): ChartTokens {
  return {
    primary: resolveBgColorFromClass('bg-primary'),
    secondary: resolveBgColorFromClass('bg-secondary'),
    accent: resolveBgColorFromClass('bg-accent'),
    success: resolveBgColorFromClass('bg-success'),
    warning: resolveBgColorFromClass('bg-warning'),
    error: resolveBgColorFromClass('bg-error'),
    info: resolveBgColorFromClass('bg-info'),
    base100: resolveBgColorFromClass('bg-base-100'),
    base200: resolveBgColorFromClass('bg-base-200'),
    base300: resolveBgColorFromClass('bg-base-300'),
    baseContent: resolveTextColorFromClass('text-base-content'),
    neutral: resolveBgColorFromClass('bg-neutral'),
  }
}

/**
 * Builds themed Chart.js options that automatically use colors from the active daisyUI theme
 */
export function buildThemedOptions(tokens: ChartTokens): ChartOptions {
  return {
    responsive: true,
    maintainAspectRatio: false,
    color: tokens.baseContent,
    font: {
      family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
      size: 12,
    },
    plugins: {
      legend: {
        labels: {
          color: tokens.baseContent,
          font: {
            family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
            size: 14,
          },
          padding: 12,
          usePointStyle: true,
        },
      },
      tooltip: {
        titleColor: tokens.baseContent,
        bodyColor: tokens.baseContent,
        backgroundColor: tokens.base100,
        borderColor: tokens.base300,
        borderWidth: 1,
        padding: 12,
        boxPadding: 6,
        usePointStyle: true,
        titleFont: {
          family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
          size: 14,
          weight: 600,
        },
        bodyFont: {
          family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
          size: 13,
        },
      },
    },
    scales: {
      x: {
        ticks: {
          color: tokens.baseContent,
          font: {
            family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
            size: 12,
          },
        },
        grid: {
          color: tokens.base300,
          lineWidth: 1,
        },
        border: {
          color: tokens.base300,
        },
      },
      y: {
        ticks: {
          color: tokens.baseContent,
          font: {
            family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
            size: 12,
          },
        },
        grid: {
          color: tokens.base300,
          lineWidth: 1,
        },
        border: {
          color: tokens.base300,
        },
      },
    },
  }
}

/**
 * Creates a canvas background plugin that fills the chart area with a background color
 * matching the daisyUI theme
 */
export function createCanvasBackgroundPlugin(backgroundColor: string): Plugin {
  return {
    id: 'canvasBackground',
    beforeDraw(chart) {
      const { ctx, chartArea } = chart
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      if (!chartArea) return

      ctx.save()
      ctx.fillStyle = backgroundColor
      ctx.fillRect(
        chartArea.left,
        chartArea.top,
        chartArea.right - chartArea.left,
        chartArea.bottom - chartArea.top,
      )
      ctx.restore()
    },
  }
}

/**
 * React hook that provides themed Chart.js options and plugins
 * based on the current daisyUI theme and language direction
 */
export interface UseChartThemeResult {
  tokens: ChartTokens
  themedOptions: ChartOptions
  backgroundPlugin: Plugin
}

export function useChartTheme(): UseChartThemeResult {
  const { i18n } = useTranslation()
  const isRTL = i18n.dir() === 'rtl'

  const result = useMemo(() => {
    const tokens = getChartTokens()
    const baseOptions = buildThemedOptions(tokens)
    const backgroundPlugin = createCanvasBackgroundPlugin(tokens.base100)

    // Apply RTL transformations if needed
    const themedOptions: ChartOptions = {
      ...baseOptions,
      indexAxis:
        isRTL && baseOptions.indexAxis === 'x' ? 'y' : baseOptions.indexAxis,
      plugins: {
        ...baseOptions.plugins,
        legend: {
          ...baseOptions.plugins?.legend,
          align: isRTL ? 'end' : 'start',
          rtl: isRTL,
        },
        tooltip: {
          ...baseOptions.plugins?.tooltip,
          rtl: isRTL,
          textDirection: isRTL ? 'rtl' : 'ltr',
        },
      },
    }

    return {
      tokens,
      themedOptions,
      backgroundPlugin,
    }
  }, [isRTL])

  return result
}
