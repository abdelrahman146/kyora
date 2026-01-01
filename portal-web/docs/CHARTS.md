# Chart.js Data Visualization System

**Status:** Production Ready  
**Last Updated:** January 1, 2026

## Overview

This document provides comprehensive guidance on using the Chart.js data visualization system in Kyora Portal Web. The system is fully integrated with daisyUI theming, supports RTL/Arabic localization, and provides mobile-first responsive charts.

## Table of Contents

1. [Architecture](#architecture)
2. [Getting Started](#getting-started)
3. [Chart Components](#chart-components)
4. [Statistics Components](#statistics-components)
5. [Theme Integration](#theme-integration)
6. [Data Transformers](#data-transformers)
7. [RTL Support](#rtl-support)
8. [Integration with TanStack Query](#integration-with-tanstack-query)
9. [Best Practices](#best-practices)
10. [Examples](#examples)

---

## Architecture

### Layer Structure

```
┌─────────────────────────────────────────────┐
│         Analytics API (Backend)             │
│   Returns: TimeSeries, KeyValue, Metrics   │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          TanStack Query Layer               │
│   Caching, Refetching, Loading States      │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│         Chart Components Layer              │
│   LineChart, BarChart, PieChart, etc.      │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          Theme Integration Layer            │
│   useChartTheme(), daisyUI Token Resolver  │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│              Chart.js Core                  │
│   Rendering Engine, Canvas, Plugins        │
└─────────────────────────────────────────────┘
```

### File Organization

```
portal-web/src/
├── lib/charts/               # Theme integration utilities
│   ├── chartTheme.ts         # Theme token resolver & hook
│   ├── chartPlugins.ts       # Custom Chart.js plugins
│   ├── chartUtils.ts         # Data transformers & formatters
│   ├── rtlSupport.ts         # RTL configuration helpers
│   └── index.ts              # Public exports
├── components/atoms/         # Chart & stat components
│   ├── LineChart.tsx         # Line/area charts
│   ├── BarChart.tsx          # Bar charts (vertical/horizontal)
│   ├── PieChart.tsx          # Pie charts
│   ├── DoughnutChart.tsx     # Doughnut charts with center label
│   ├── MixedChart.tsx        # Combined bar + line charts
│   ├── ChartCard.tsx         # Chart wrapper with title, loading, error
│   ├── ChartSkeleton.tsx     # Loading skeletons
│   ├── ChartEmptyState.tsx   # Empty state with icons
│   ├── StatCard.tsx          # Simple stat cards
│   ├── StatCardGroup.tsx     # Responsive stat card grid
│   ├── ComplexStatCard.tsx   # Advanced stat cards
│   └── StatCardSkeleton.tsx  # Stat card loading state
└── i18n/
    ├── en/analytics.json     # English chart translations
    └── ar/analytics.json     # Arabic chart translations
```

---

## Getting Started

### Basic Usage

```tsx
import { LineChart, ChartCard } from '@/components'
import { useChartTheme, transformTimeSeriesToChartData } from '@/lib/charts'

function SalesChart({ data }) {
  const { tokens } = useChartTheme()
  
  const chartData = transformTimeSeriesToChartData(
    data,
    'Revenue',
    tokens.primary
  )

  return (
    <ChartCard title="Sales Over Time" height={320}>
      <LineChart data={chartData} enableArea />
    </ChartCard>
  )
}
```

---

## Chart Components

### LineChart

**Purpose:** Display time-series data, trends, and continuous metrics.

**Props:**
```typescript
interface LineChartProps {
  data: ChartData<'line'>         // Chart.js line chart data
  options?: ChartOptions<'line'>  // Custom Chart.js options
  height?: number                 // Chart height (default: 320px)
  className?: string              // Additional CSS classes
  enableArea?: boolean            // Fill area under line (default: false)
  enableDecimation?: boolean      // Enable data decimation for large datasets
}
```

**Example:**
```tsx
import { LineChart, ChartCard } from '@/components'
import { useQuery } from '@tanstack/react-query'
import { salesApi } from '@/api'

function RevenueChart() {
  const { data, isLoading } = useQuery({
    queryKey: ['analytics', 'revenue'],
    queryFn: () => salesApi.getRevenueOverTime()
  })

  const chartData = {
    labels: data?.series.map(s => s.label) || [],
    datasets: [{
      label: 'Revenue',
      data: data?.series.map(s => s.value) || [],
      borderColor: 'rgb(13, 148, 136)',
      backgroundColor: 'rgba(13, 148, 136, 0.1)',
    }]
  }

  return (
    <ChartCard title="Revenue Over Time" isLoading={isLoading} height={320}>
      <LineChart data={chartData} enableArea />
    </ChartCard>
  )
}
```

---

### BarChart

**Purpose:** Compare categorical data, show distributions.

**Props:**
```typescript
interface BarChartProps {
  data: ChartData<'bar'>
  options?: ChartOptions<'bar'>
  height?: number                // Default: 320px
  className?: string
  horizontal?: boolean           // Horizontal bars (default: false)
  stacked?: boolean             // Stacked bars (default: false)
}
```

**Example:**
```tsx
import { BarChart, ChartCard } from '@/components'
import { transformKeyValueToBarData, generateColorPalette } from '@/lib/charts'

function ExpensesByCategoryChart({ data }) {
  const colors = generateColorPalette('rgb(234, 179, 8)', data.length)
  const chartData = transformKeyValueToBarData(data, 'Expenses', colors)

  return (
    <ChartCard title="Expenses by Category">
      <BarChart data={chartData} />
    </ChartCard>
  )
}
```

---

### PieChart & DoughnutChart

**Purpose:** Show proportions and percentages.

**Props:**
```typescript
interface PieChartProps {
  data: ChartData<'pie'>
  options?: ChartOptions<'pie'>
  height?: number                // Default: 320px
  className?: string
}

interface DoughnutChartProps {
  data: ChartData<'doughnut'>
  options?: ChartOptions<'doughnut'>
  height?: number
  className?: string
  centerLabel?: string           // Display label in center (e.g., total count)
  centerLabelColor?: string      // Center label color
}
```

**Example:**
```tsx
import { DoughnutChart, ChartCard } from '@/components'
import { transformKeyValueToPieData } from '@/lib/charts'

function OrderStatusChart({ data }) {
  const colors = ['#10B981', '#EAB308', '#EF4444', '#64748b']
  const chartData = transformKeyValueToPieData(data, colors)
  const total = data.reduce((sum, kv) => sum + Number(kv.value), 0)

  return (
    <ChartCard title="Order Status Breakdown">
      <DoughnutChart 
        data={chartData} 
        centerLabel={`${total} Orders`}
      />
    </ChartCard>
  )
}
```

---

### MixedChart

**Purpose:** Combine bar and line charts for comparative analysis (e.g., revenue vs profit).

**Props:**
```typescript
interface MixedChartProps {
  data: ChartData<'bar' | 'line'>
  options?: ChartOptions
  height?: number
  className?: string
}
```

**Example:**
```tsx
import { MixedChart, ChartCard } from '@/components'

function ProfitLossChart({ data }) {
  const chartData = {
    labels: data.months,
    datasets: [
      {
        type: 'bar' as const,
        label: 'Revenue',
        data: data.revenue,
        backgroundColor: 'rgb(13, 148, 136)',
      },
      {
        type: 'line' as const,
        label: 'Profit',
        data: data.profit,
        borderColor: 'rgb(234, 179, 8)',
        backgroundColor: 'rgba(234, 179, 8, 0.1)',
      }
    ]
  }

  return (
    <ChartCard title="Profit & Loss Statement">
      <MixedChart data={chartData} height={400} />
    </ChartCard>
  )
}
```

---

### ChartCard

**Purpose:** Wrapper component for charts with title, loading states, error handling.

**Props:**
```typescript
interface ChartCardProps {
  title?: string                 // Card title
  subtitle?: string              // Card subtitle
  children: ReactNode            // Chart component
  isLoading?: boolean            // Show loading skeleton
  isEmpty?: boolean              // Show empty state
  error?: Error | null           // Show error state
  onRetry?: () => void          // Retry function for errors
  actions?: ReactNode            // Action buttons (top-right)
  className?: string
  height?: number                // Chart area height
  chartType?: 'line' | 'bar' | 'pie' | 'doughnut' | 'mixed'
}
```

**Example:**
```tsx
import { ChartCard, LineChart } from '@/components'
import { RefreshCw } from 'lucide-react'

function MyChart() {
  const { data, isLoading, error, refetch } = useQuery(...)

  return (
    <ChartCard
      title="Sales Performance"
      subtitle="Last 30 days"
      isLoading={isLoading}
      isEmpty={!data || data.series.length === 0}
      error={error}
      onRetry={refetch}
      actions={
        <button onClick={refetch}>
          <RefreshCw className="h-4 w-4" />
        </button>
      }
      chartType="line"
      height={320}
    >
      <LineChart data={chartData} />
    </ChartCard>
  )
}
```

---

## Statistics Components

### StatCard

**Purpose:** Display key metrics with optional trend indicators.

**Props:**
```typescript
interface StatCardProps {
  label: string                  // Metric label
  value: string | number         // Primary value
  icon?: ReactNode               // Icon component
  trend?: 'up' | 'down'         // Trend direction
  trendValue?: string            // Trend percentage (e.g., "+12%")
  variant?: 'default' | 'success' | 'warning' | 'error' | 'info'
  className?: string
}
```

**Example:**
```tsx
import { StatCard, StatCardGroup } from '@/components'
import { DollarSign, TrendingUp, ShoppingCart } from 'lucide-react'
import { formatCurrency } from '@/lib/formatCurrency'

function DashboardMetrics({ metrics }) {
  return (
    <StatCardGroup cols={3}>
      <StatCard
        label="Total Revenue"
        value={formatCurrency(metrics.revenue, 'USD')}
        icon={<DollarSign className="h-5 w-5 text-primary" />}
        trend="up"
        trendValue="+12.5%"
        variant="success"
      />
      
      <StatCard
        label="Open Orders"
        value={metrics.openOrders}
        icon={<ShoppingCart className="h-5 w-5 text-secondary" />}
      />
      
      <StatCard
        label="Low Stock Items"
        value={metrics.lowStockItems}
        icon={<TrendingUp className="h-5 w-5 text-warning" />}
        variant="warning"
      />
    </StatCardGroup>
  )
}
```

---

### ComplexStatCard

**Purpose:** Advanced stat card with multiple metrics and badges.

**Props:**
```typescript
interface ComplexStatCardProps {
  label: string
  value: string | number
  icon?: ReactNode
  secondaryMetrics?: SecondaryMetric[]  // Array of additional metrics
  comparisonText?: string               // "vs last month" text
  statusBadge?: {
    label: string
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info'
  }
  variant?: 'default' | 'success' | 'warning' | 'error' | 'info'
  className?: string
}

interface SecondaryMetric {
  label: string
  value: string | number
}
```

**Example:**
```tsx
import { ComplexStatCard } from '@/components'
import { formatCurrency } from '@/lib/formatCurrency'

function CashFlowCard({ data }) {
  return (
    <ComplexStatCard
      label="Cash on Hand"
      value={formatCurrency(data.cashOnHand, 'USD')}
      comparisonText="vs last month: +$2,500"
      statusBadge={{ label: "Healthy", variant: "success" }}
      secondaryMetrics={[
        { label: "Cash In", value: formatCurrency(data.cashIn, 'USD') },
        { label: "Cash Out", value: formatCurrency(data.cashOut, 'USD') }
      ]}
      variant="success"
    />
  )
}
```

---

## Theme Integration

### useChartTheme Hook

**Purpose:** Automatically apply daisyUI theme colors and RTL configuration to charts.

**Returns:**
```typescript
interface UseChartThemeResult {
  tokens: ChartTokens           // Resolved color tokens
  themedOptions: ChartOptions   // Pre-configured Chart.js options
  backgroundPlugin: Plugin      // Background plugin
}
```

**Example:**
```tsx
import { useChartTheme } from '@/lib/charts'

function ThemedChart() {
  const { tokens, themedOptions, backgroundPlugin } = useChartTheme()

  const chartData = {
    datasets: [{
      data: [10, 20, 30],
      borderColor: tokens.primary,       // Uses daisyUI primary color
      backgroundColor: tokens.success,   // Uses daisyUI success color
    }]
  }

  return (
    <Line 
      data={chartData} 
      options={themedOptions}           // Auto-themed options
      plugins={[backgroundPlugin]}      // Applies bg-base-100 background
    />
  )
}
```

---

## Data Transformers

### transformTimeSeriesToChartData

**Purpose:** Convert backend `TimeSeries` format to Chart.js line chart data.

**Signature:**
```typescript
function transformTimeSeriesToChartData(
  timeSeries: TimeSeries,
  datasetLabel: string,
  color: string
): ChartData<'line'>
```

**Example:**
```tsx
import { transformTimeSeriesToChartData } from '@/lib/charts'

const chartData = transformTimeSeriesToChartData(
  apiResponse.revenueOverTime,
  'Revenue',
  'rgb(13, 148, 136)'
)
```

---

### transformKeyValueToBarData

**Purpose:** Convert `KeyValue[]` to bar chart data.

**Signature:**
```typescript
function transformKeyValueToBarData(
  keyValues: KeyValue[],
  datasetLabel: string,
  colors: string[]
): ChartData<'bar'>
```

---

### transformKeyValueToPieData

**Purpose:** Convert `KeyValue[]` to pie/doughnut chart data.

**Signature:**
```typescript
function transformKeyValueToPieData(
  keyValues: KeyValue[],
  colors: string[]
): ChartData<'pie' | 'doughnut'>
```

---

### Utility Functions

```typescript
// Generate color palette with varying opacity
generateColorPalette(baseColor: string, count: number): string[]

// Format currency for tooltips
formatChartCurrency(value: number, currency: string): string

// Format large numbers (1000 → 1K, 1000000 → 1M)
formatChartNumber(value: number): string

// Check if decimation should be enabled
shouldEnableDecimation(dataPointCount: number): boolean

// Calculate percentages for pie charts
calculatePercentages(values: number[]): string[]
```

---

## RTL Support

All charts automatically support RTL when the language is set to Arabic:

- **X-axis reversal**: Latest data appears on the left (RTL) or right (LTR)
- **Legend positioning**: Aligned to the start (right in RTL, left in LTR)
- **Tooltip text direction**: Proper RTL text rendering
- **Number formatting**: Western Arabic numerals (0-9) used consistently

**No additional configuration required!** The `useChartTheme()` hook handles everything automatically based on the current `i18n` language.

---

## Integration with TanStack Query

### Pattern: Single API Call → Multiple Charts

Analytics APIs return all data at once. Handle loading/error at the page level:

```tsx
import { useQuery } from '@tanstack/react-query'
import { ChartCard, LineChart, BarChart, StatCardGroup, StatCard } from '@/components'

function SalesAnalyticsPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['analytics', 'sales', dateRange],
    queryFn: () => analyticsApi.getSalesAnalytics(dateRange),
  })

  if (isLoading) {
    return <PageLoadingSkeleton />
  }

  if (error) {
    return <PageErrorBoundary error={error} />
  }

  return (
    <div className="space-y-6">
      {/* Statistics Cards */}
      <StatCardGroup cols={4}>
        <StatCard label="Total Revenue" value={data.totalRevenue} />
        <StatCard label="Gross Profit" value={data.grossProfit} />
        <StatCard label="Total Orders" value={data.totalOrders} />
        <StatCard label="Avg Order Value" value={data.averageOrderValue} />
      </StatCardGroup>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ChartCard title="Revenue Over Time">
          <LineChart data={transformTimeSeriesToChartData(data.revenueOverTime)} />
        </ChartCard>

        <ChartCard title="Top Products">
          <BarChart data={transformTopProductsData(data.topSellingProducts)} />
        </ChartCard>
      </div>
    </div>
  )
}
```

---

## Best Practices

### 1. Always use ChartCard wrapper

```tsx
// ✅ Good
<ChartCard title="Sales" isLoading={isLoading}>
  <LineChart data={data} />
</ChartCard>

// ❌ Bad (no loading/error handling)
<LineChart data={data} />
```

---

### 2. Handle empty data gracefully

```tsx
<ChartCard 
  title="Revenue" 
  isEmpty={!data || data.series.length === 0}
>
  <LineChart data={data} />
</ChartCard>
```

---

### 3. Use semantic variants for stat cards

```tsx
// Low stock warning
<StatCard 
  label="Low Stock Items" 
  value={count} 
  variant="warning"  // Yellow background
/>

// Profit success
<StatCard 
  label="Net Profit" 
  value={profit} 
  variant="success"  // Green background
/>
```

---

### 4. Enable decimation for large datasets

```tsx
<LineChart 
  data={largeDataset} 
  enableDecimation={data.series.length > 100} 
/>
```

---

### 5. Use proper height for mobile

```tsx
// Default 320px is good for mobile
<ChartCard height={320}>
  <LineChart data={data} />
</ChartCard>

// Taller charts for desktop-only views
<ChartCard height={400} className="hidden lg:block">
  <LineChart data={data} />
</ChartCard>
```

---


## Troubleshooting

### Charts not rendering

1. Ensure Chart.js is registered in `main.tsx`
2. Check if `data` prop is correctly formatted
3. Verify chart height is set (default 320px)

### Colors not matching theme

- Use `useChartTheme()` hook
- Access colors via `tokens.primary`, `tokens.secondary`, etc.
- Don't hard-code hex values

### RTL not working

- Check `i18n.dir()` returns `'rtl'`
- Use `useChartTheme()` which handles RTL automatically
- Verify `document.documentElement.dir === 'rtl'`

### Performance issues with large datasets

- Enable decimation: `<LineChart enableDecimation />`
- Reduce data points before passing to chart
- Consider pagination or date range filters

---

## Related Documentation

- [Kyora Design System (KDS)](./branding.instructions.md)
- [Portal Web Architecture](./portal-web.instructions.md)
- [Best Practices](./bestpractice.instructions.md)
- [Chart.js Official Docs](https://www.chartjs.org)
- [react-chartjs-2 Guide](https://react-chartjs-2.js.org)

---

**Questions?** Refer to the codebase examples or consult the Chart.js documentation.
