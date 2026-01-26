import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'
import type { ChartOptions, Plugin } from 'chart.js'

/**
 * Area chart fill opacity constants
 * Design spec: Subtle fills with no gradients
 */
export const AREA_FILL_OPACITY_SPARKLINE = 0.1 // Ultra-compact sparklines
export const AREA_FILL_OPACITY_NORMAL = 0.08 // Standard area charts - very light translucent fill, bold line on top

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
 * Returns the multi-series color palette sequence following design specs
 * Sequence: primary → info → secondary → success → warning → accent
 */
export function getMultiSeriesColors(tokens: ChartTokens): Array<string> {
  return [
    tokens.primary,
    tokens.info,
    tokens.secondary,
    tokens.success,
    tokens.warning,
    tokens.accent,
  ]
}

/**
 * Checks if a color string is already translucent (alpha < 1)
 * @param color - Color string to check
 * @returns true if the color has alpha < 1
 */
export function isTranslucentColor(color: string): boolean {
  const trimmed = color.trim()

  // Match rgba(r, g, b, a) with comma or space separator
  const rgbaMatch = trimmed.match(
    /rgba?\([^,]+,\s*[^,]+,\s*[^,]+(?:[,/]\s*0?\.\d+|[,/]\s*0)\s*\)/,
  )
  if (rgbaMatch) return true

  // Match rgb(r g b / alpha) modern format
  const modernMatch = trimmed.match(/rgb\([^)]+\/\s*0?\.\d+\s*\)/)
  if (modernMatch) return true

  return false
}

/**
 * Converts a color to RGBA with specified opacity
 * Used for area fills (subtly tinted surfaces as per design specs)
 *
 * Handles modern CSS Color Level 4 formats:
 * - rgb(r, g, b) and rgba(r, g, b, a) (comma-separated)
 * - rgb(r g b) and rgb(r g b / a) (space-separated)
 * - #rgb and #rrggbb (hex)
 * - Other formats (hsl, oklch, css vars): resolved via hidden element
 *
 * Always returns comma-separated rgba(r, g, b, opacity)
 */
export function colorWithOpacity(color: string, opacity: number): string {
  const trimmed = color.trim()

  // Handle hex colors
  if (trimmed.startsWith('#')) {
    const hex = trimmed.slice(1)
    let r: number, g: number, b: number

    if (hex.length === 3) {
      // #rgb → #rrggbb
      r = parseInt(hex[0] + hex[0], 16)
      g = parseInt(hex[1] + hex[1], 16)
      b = parseInt(hex[2] + hex[2], 16)
    } else if (hex.length === 6) {
      r = parseInt(hex.slice(0, 2), 16)
      g = parseInt(hex.slice(2, 4), 16)
      b = parseInt(hex.slice(4, 6), 16)
    } else {
      // Invalid hex, fallback to original
      return color
    }

    return `rgba(${r}, ${g}, ${b}, ${opacity})`
  }

  // Handle rgb/rgba formats
  if (trimmed.startsWith('rgb')) {
    // Extract content inside parentheses
    const match = trimmed.match(/rgba?\(([^)]+)\)/)
    if (!match) return color

    const content = match[1]
    // Replace forward slash with space, then split on comma or whitespace
    const normalized = content.replace(/\//g, ' ')
    const parts = normalized
      .split(/[\s,]+/)
      .map((s) => s.trim())
      .filter(Boolean)

    if (parts.length >= 3) {
      const r = parseFloat(parts[0])
      const g = parseFloat(parts[1])
      const b = parseFloat(parts[2])

      if (!isNaN(r) && !isNaN(g) && !isNaN(b)) {
        return `rgba(${r}, ${g}, ${b}, ${opacity})`
      }
    }

    // Failed to parse, return original
    return color
  }

  // Fallback: resolve other formats (hsl, oklch, css vars) via hidden element
  try {
    const el = document.createElement('div')
    el.style.position = 'absolute'
    el.style.left = '-9999px'
    el.style.visibility = 'hidden'
    el.style.color = trimmed
    document.body.appendChild(el)

    const resolvedColor = getComputedStyle(el).color
    document.body.removeChild(el)

    // Recursively parse the resolved color
    if (resolvedColor && resolvedColor !== trimmed) {
      return colorWithOpacity(resolvedColor, opacity)
    }
  } catch {
    // If DOM manipulation fails, return original
  }

  // Ultimate fallback
  return color
}

/**
 * Builds themed Chart.js options that automatically use colors from the active daisyUI theme
 * Following design specs:
 * - Bar charts: 8px border radius
 * - Line charts: 0.4 tension (smooth curves)
 * - Grid lines: 30% opacity on base-300, Y-axis only
 * - Axis ticks: 60% opacity on base-content, 12px font, 8px padding
 * - Tooltips: 8px radius, 12px padding, base-100 bg with base-300 border
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
    datasets: {
      line: {
        pointRadius: 0,
        pointHoverRadius: 4,
        pointHitRadius: 10,
      },
    },
    elements: {
      bar: {
        borderRadius: 8,
      },
      line: {
        tension: 0.4,
      },
      point: {
        radius: 0,
        hitRadius: 10,
        hoverRadius: 4,
      },
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
        cornerRadius: 8,
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
          color: colorWithOpacity(tokens.baseContent, 0.6),
          font: {
            family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
            size: 12,
          },
          padding: 8,
        },
        grid: {
          display: false,
        },
        border: {
          display: false,
        },
      },
      y: {
        ticks: {
          color: colorWithOpacity(tokens.baseContent, 0.6),
          font: {
            family: "'IBM Plex Sans Arabic', -apple-system, sans-serif",
            size: 12,
          },
          padding: 8,
        },
        grid: {
          color: colorWithOpacity(tokens.base300, 0.2),
          lineWidth: 1,
        },
        border: {
          display: false,
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
