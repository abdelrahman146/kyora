import { useMemo } from 'react'
import type { ChartData } from 'chart.js'
import {
  BarChart,
  ChartCard,
  DoughnutChart,
  GaugeChart,
  LineChart,
  MixedChart,
  PieChart,
  StatCard,
} from '@/components/charts'
import {
  AREA_FILL_OPACITY_NORMAL,
  colorWithOpacity,
  getMultiSeriesColors,
  useChartTheme,
} from '@/lib/charts'

/**
 * ChartsDemo Component
 *
 * Comprehensive showcase of all chart components with the revamped design system.
 * Demonstrates:
 * - Modern, elegant aesthetics
 * - Strict adherence to Kyora design system (no gradients, no shadows)
 * - Semantic color usage
 * - RTL compatibility
 * - Mobile responsiveness
 */
export function ChartsDemo() {
  const { tokens } = useChartTheme()
  const multiSeriesColors = useMemo(
    () => getMultiSeriesColors(tokens),
    [tokens],
  )

  // Bar Chart Data - Vertical
  const verticalBarData = useMemo<ChartData<'bar'>>(() => {
    return {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
      datasets: [
        {
          label: 'Revenue',
          data: [12000, 19000, 15000, 25000, 22000, 30000],
          backgroundColor: tokens.primary,
        },
      ],
    }
  }, [tokens.primary])

  // Bar Chart Data - Horizontal
  const horizontalBarData = useMemo<ChartData<'bar'>>(() => {
    return {
      labels: ['Product A', 'Product B', 'Product C', 'Product D', 'Product E'],
      datasets: [
        {
          label: 'Sales',
          data: [65, 59, 80, 81, 56],
          backgroundColor: tokens.info,
        },
      ],
    }
  }, [tokens.info])

  // Bar Chart Data - Stacked
  const stackedBarData = useMemo<ChartData<'bar'>>(() => {
    return {
      labels: ['Q1', 'Q2', 'Q3', 'Q4'],
      datasets: [
        {
          label: 'Online',
          data: [12, 19, 15, 25],
          backgroundColor: tokens.primary,
        },
        {
          label: 'In-Store',
          data: [8, 11, 13, 18],
          backgroundColor: tokens.info,
        },
        {
          label: 'Wholesale',
          data: [5, 7, 9, 12],
          backgroundColor: tokens.secondary,
        },
      ],
    }
  }, [tokens.primary, tokens.info, tokens.secondary])

  // Line Chart Data - Simple
  const simpleLineData = useMemo<ChartData<'line'>>(() => {
    return {
      labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Week 5', 'Week 6'],
      datasets: [
        {
          label: 'Orders',
          data: [45, 52, 49, 60, 58, 65],
          borderColor: tokens.success,
          backgroundColor: tokens.success,
          tension: 0.4,
        },
      ],
    }
  }, [tokens.success])

  // Line Chart Data - Area
  const areaLineData = useMemo<ChartData<'line'>>(() => {
    return {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug'],
      datasets: [
        {
          label: 'Visitors',
          data: [300, 450, 400, 550, 520, 600, 580, 650],
          borderColor: tokens.primary,
          backgroundColor: colorWithOpacity(
            tokens.primary,
            AREA_FILL_OPACITY_NORMAL,
          ),
          borderWidth: 3,
          fill: true,
          tension: 0.4,
        },
      ],
    }
  }, [tokens.primary])

  // Line Chart Data - Multi-Series
  const multiSeriesLineData = useMemo<ChartData<'line'>>(() => {
    return {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
      datasets: [
        {
          label: 'Revenue',
          data: [15000, 18000, 17000, 21000, 19000, 23000],
          borderColor: multiSeriesColors[0],
          backgroundColor: multiSeriesColors[0],
          tension: 0.4,
        },
        {
          label: 'Profit',
          data: [8000, 9500, 9000, 11000, 10500, 12500],
          borderColor: multiSeriesColors[1],
          backgroundColor: multiSeriesColors[1],
          tension: 0.4,
        },
        {
          label: 'Expenses',
          data: [7000, 8500, 8000, 10000, 8500, 10500],
          borderColor: multiSeriesColors[2],
          backgroundColor: multiSeriesColors[2],
          tension: 0.4,
        },
      ],
    }
  }, [multiSeriesColors])

  // Doughnut Chart Data
  const doughnutData = useMemo<ChartData<'doughnut'>>(() => {
    return {
      labels: ['Electronics', 'Clothing', 'Food', 'Books', 'Other'],
      datasets: [
        {
          data: [35, 25, 20, 12, 8],
          backgroundColor: multiSeriesColors,
          borderWidth: 0,
        },
      ],
    }
  }, [multiSeriesColors])

  // Pie Chart Data
  const pieData = useMemo<ChartData<'pie'>>(() => {
    return {
      labels: ['Completed', 'Processing', 'Cancelled'],
      datasets: [
        {
          data: [65, 25, 10],
          backgroundColor: [tokens.success, tokens.warning, tokens.error],
          borderWidth: 0,
        },
      ],
    }
  }, [tokens.success, tokens.warning, tokens.error])

  // Mixed Chart Data
  const mixedData = useMemo<ChartData<'bar' | 'line'>>(() => {
    return {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
      datasets: [
        {
          type: 'bar' as const,
          label: 'Sales',
          data: [12, 19, 15, 25, 22, 30],
          backgroundColor: tokens.primary,
        },
        {
          type: 'line' as const,
          label: 'Target',
          data: [18, 18, 20, 22, 24, 26],
          borderColor: tokens.warning,
          backgroundColor: tokens.warning,
          tension: 0.4,
        },
      ],
    }
  }, [tokens.primary, tokens.warning])

  // Sparkline data
  const sparklineData1 = [12, 15, 13, 18, 16, 21, 19, 25, 23, 28]
  const sparklineData2 = [30, 28, 32, 29, 35, 33, 38, 36, 40, 42]
  const sparklineData3 = [50, 48, 45, 42, 40, 38, 35, 33, 30, 28]

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8 flex flex-col gap-2">
        <h1 className="text-3xl font-bold text-base-content">Charts Demo</h1>
        <p className="text-base-content/60">
          Showcase of all chart components with the revamped design system
        </p>
      </div>

      {/* Stat Cards Section */}
      <section className="mb-12">
        <h2 className="mb-4 text-2xl font-semibold text-base-content">
          Stat Cards
        </h2>
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          <StatCard
            title="Total Revenue"
            value="$127,500"
            trend={{ value: 12.5, direction: 'up' }}
            sparklineData={sparklineData1}
            sparklineColor={tokens.primary}
            subtitle="Last 30 days"
          />
          <StatCard
            title="Active Orders"
            value="342"
            trend={{ value: 8.2, direction: 'up' }}
            sparklineData={sparklineData2}
            sparklineColor={tokens.success}
            subtitle="Current period"
          />
          <StatCard
            title="Conversion Rate"
            value="24.8%"
            trend={{ value: -3.1, direction: 'down' }}
            sparklineData={sparklineData3}
            sparklineColor={tokens.error}
            subtitle="Compared to last month"
          />
        </div>
      </section>

      {/* Gauge Charts Section */}
      <section className="mb-12">
        <h2 className="mb-4 text-2xl font-semibold text-base-content">
          Gauge Charts
        </h2>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          <GaugeChart
            value={75}
            label="Sales Target"
            trend={{ value: 5.2, direction: 'up' }}
            color={tokens.primary}
          />
          <GaugeChart
            value={92}
            label="Customer Satisfaction"
            trend={{ value: 2.1, direction: 'up' }}
            color={tokens.success}
          />
          <GaugeChart
            value={45}
            label="Inventory Health"
            trend={{ value: -1.5, direction: 'down' }}
            color={tokens.warning}
          />
          <GaugeChart
            value={88}
            label="On-Time Delivery"
            trend={{ value: 4.3, direction: 'up' }}
            color={tokens.info}
          />
        </div>
      </section>

      {/* Bar Charts Section */}
      <section className="mb-12">
        <h2 className="mb-4 text-2xl font-semibold text-base-content">
          Bar Charts
        </h2>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <ChartCard
            title="Monthly Revenue"
            subtitle="Vertical bar chart with rounded corners"
            chartType="bar"
            height={320}
          >
            <BarChart data={verticalBarData} height={320} />
          </ChartCard>

          <ChartCard
            title="Product Sales"
            subtitle="Horizontal bar chart"
            chartType="bar"
            height={320}
          >
            <BarChart data={horizontalBarData} height={320} horizontal />
          </ChartCard>

          <ChartCard
            title="Quarterly Sales by Channel"
            subtitle="Stacked bar chart with multiple series"
            chartType="bar"
            height={320}
            className="lg:col-span-2"
          >
            <BarChart data={stackedBarData} height={320} stacked />
          </ChartCard>
        </div>
      </section>

      {/* Line Charts Section */}
      <section className="mb-12">
        <h2 className="mb-4 text-2xl font-semibold text-base-content">
          Line Charts
        </h2>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <ChartCard
            title="Weekly Orders"
            subtitle="Simple line chart with smooth curves"
            chartType="line"
            height={320}
          >
            <LineChart data={simpleLineData} height={320} />
          </ChartCard>

          <ChartCard
            title="Website Traffic"
            subtitle="Area chart with 20% opacity fill"
            chartType="line"
            height={320}
          >
            <LineChart data={areaLineData} height={320} enableArea />
          </ChartCard>

          <ChartCard
            title="Financial Overview"
            subtitle="Multi-series line chart"
            chartType="line"
            height={400}
            className="lg:col-span-2"
          >
            <LineChart data={multiSeriesLineData} height={400} />
          </ChartCard>
        </div>
      </section>

      {/* Pie & Doughnut Charts Section */}
      <section className="mb-12">
        <h2 className="mb-4 text-2xl font-semibold text-base-content">
          Pie & Doughnut Charts
        </h2>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <ChartCard
            title="Sales by Category"
            subtitle="Doughnut chart with center label"
            chartType="doughnut"
            height={320}
          >
            <DoughnutChart
              data={doughnutData}
              height={320}
              centerLabel="100%"
              centerLabelColor={tokens.baseContent}
            />
          </ChartCard>

          <ChartCard
            title="Order Status"
            subtitle="Pie chart showing distribution"
            chartType="pie"
            height={320}
          >
            <PieChart data={pieData} height={320} />
          </ChartCard>
        </div>
      </section>

      {/* Mixed Chart Section */}
      <section className="mb-12">
        <h2 className="mb-4 text-2xl font-semibold text-base-content">
          Mixed Charts
        </h2>
        <ChartCard
          title="Sales vs Target"
          subtitle="Combined bar and line chart"
          chartType="mixed"
          height={400}
        >
          <MixedChart data={mixedData} height={400} />
        </ChartCard>
      </section>

      {/* Compact Section */}
      <section className="mb-12">
        <h2 className="mb-4 text-2xl font-semibold text-base-content">
          Compact Charts
        </h2>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <ChartCard title="Quick View" chartType="bar" height={240}>
            <BarChart data={verticalBarData} height={240} />
          </ChartCard>
          <ChartCard title="Trend" chartType="line" height={240}>
            <LineChart data={simpleLineData} height={240} />
          </ChartCard>
          <ChartCard title="Distribution" chartType="doughnut" height={240}>
            <DoughnutChart data={doughnutData} height={240} />
          </ChartCard>
        </div>
      </section>

      {/* Design System Notes */}
      <section className="rounded-box border border-base-300 bg-base-200 p-6">
        <h3 className="mb-3 text-lg font-semibold text-base-content">
          Design System Notes
        </h3>
        <ul className="space-y-2 text-sm text-base-content/70">
          <li>✅ Bar charts: 8px border radius</li>
          <li>✅ Line charts: 0.4 tension (smooth curves)</li>
          <li>✅ Grid lines: 30% opacity on base-300, Y-axis only</li>
          <li>✅ Axis ticks: 60% opacity, 12px font, 8px padding</li>
          <li>✅ Tooltips: 8px radius, 12px padding</li>
          <li>
            ✅ Area fills: 12% opacity for normal charts, 10% for sparklines (NO
            gradients)
          </li>
          <li>
            ✅ Multi-series: primary → info → secondary → success → warning →
            accent
          </li>
          <li>✅ All colors use daisyUI semantic tokens</li>
          <li>✅ RTL-compatible layouts</li>
          <li>✅ Mobile responsive</li>
        </ul>
      </section>
    </div>
  )
}
