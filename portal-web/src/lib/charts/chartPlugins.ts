import type { Plugin } from 'chart.js'

/**
 * Plugin that adds a touch-optimized tooltip with larger hit areas for mobile devices
 */
export const touchOptimizedTooltipPlugin: Plugin<any> = {
  id: 'touchOptimizedTooltip',
  defaults: {},
  beforeEvent(chart, args) {
    const event = args.event

    // Increase hit radius for touch events
    if (
      event.native &&
      (event.native.type === 'touchstart' || event.native.type === 'touchmove')
    ) {
      chart.options.plugins = chart.options.plugins || {}
      chart.options.plugins.tooltip = chart.options.plugins.tooltip || {}
      chart.options.elements = chart.options.elements || {}
      chart.options.elements.point = chart.options.elements.point || {}
      chart.options.elements.point.hitRadius = 20
    }
  },
}

/**
 * Plugin that reverses the x-axis for RTL languages
 */
export function createRTLPlugin(isRTL: boolean): Plugin<any> {
  return {
    id: 'rtlSupport',
    beforeInit(chart) {
      if (!isRTL) return

      const originalUpdate = chart.update
      chart.update = function (mode?: any) {
        if (chart.options.scales?.x) {
          chart.options.scales.x.reverse = true
        }
        originalUpdate.call(this, mode)
      }
    },
  }
}

/**
 * Plugin that adds gradient backgrounds to area charts
 */
export function createGradientPlugin(
  startColor: string,
  endColor: string,
  opacity: number = 0.2,
): Plugin<any> {
  return {
    id: 'gradientBackground',
    beforeDraw(chart) {
      const { ctx, chartArea } = chart
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      if (!chartArea) return

      const gradient = ctx.createLinearGradient(
        0,
        chartArea.top,
        0,
        chartArea.bottom,
      )
      gradient.addColorStop(
        0,
        startColor.replace(')', `, ${opacity})`).replace('rgb', 'rgba'),
      )
      gradient.addColorStop(
        1,
        endColor.replace(')', `, 0)`).replace('rgb', 'rgba'),
      )

      chart.data.datasets.forEach((dataset: any, index) => {
        if (dataset.fill) {
          const meta = chart.getDatasetMeta(index)
          if (meta.type === 'line') {
            dataset.backgroundColor = gradient
          }
        }
      })
    },
  }
}

/**
 * Plugin that adds a center label to doughnut/pie charts
 */
export function createCenterLabelPlugin(
  text: string,
  color: string,
): Plugin<any> {
  return {
    id: 'centerLabel',
    afterDraw(chart) {
      const { ctx, chartArea } = chart
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      if (!chartArea) return

      const centerX = (chartArea.left + chartArea.right) / 2
      const centerY = (chartArea.top + chartArea.bottom) / 2

      ctx.save()
      ctx.font = "600 24px 'IBM Plex Sans Arabic', -apple-system, sans-serif"
      ctx.fillStyle = color
      ctx.textAlign = 'center'
      ctx.textBaseline = 'middle'
      ctx.fillText(text, centerX, centerY)
      ctx.restore()
    },
  }
}
