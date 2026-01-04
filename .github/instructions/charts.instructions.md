---
description: Chart.js Data Visualization - Portal Web
applyTo: "portal-web/**"
---

# Chart.js Data Visualization System

**SSOT Hierarchy:**

- Parent: copilot-instructions.md
- Peers: portal-web-architecture.instructions.md, ui-implementation.instructions.md
- Required Reading: design-tokens.instructions.md

**When to Read:**

- Creating charts/data visualizations
- Dashboard metrics
- Analytics pages
- Stat cards

---

## 1. Quick Start

```tsx
import { LineChart, ChartCard, useChartTheme } from "@/lib/charts";

function SalesChart({ data }) {
  const { tokens } = useChartTheme();

  const chartData = {
    labels: data.map((d) => d.date),
    datasets: [
      {
        label: "Revenue",
        data: data.map((d) => d.amount),
        borderColor: tokens.primary,
        backgroundColor: tokens.primaryLight,
      },
    ],
  };

  return (
    <ChartCard title="Sales Over Time" height={320}>
      <LineChart data={chartData} enableArea />
    </ChartCard>
  );
}
```

---

## 2. Chart Components

### LineChart

**Use for:** Time-series data, trends, continuous metrics

```tsx
<LineChart
  data={chartData}
  height={320}
  enableArea // Fill area under line
  enableDecimation // For large datasets (>100 points)
/>
```

### BarChart

**Use for:** Categorical comparisons, distributions

```tsx
<BarChart
  data={chartData}
  horizontal // Horizontal bars
  stacked // Stacked bars
/>
```

### PieChart / DoughnutChart

**Use for:** Proportions, percentages

```tsx
<DoughnutChart
  data={chartData}
  centerLabel="150 Orders" // Center label (doughnut only)
/>
```

### MixedChart

**Use for:** Combining bar + line (e.g., revenue vs profit)

```tsx
<MixedChart data={{
  labels: months,
  datasets: [
    { type: 'bar' as const, label: 'Revenue', data: [...] },
    { type: 'line' as const, label: 'Profit', data: [...] },
  ]
}} />
```

---

## 3. ChartCard Wrapper (REQUIRED)

**Always wrap charts in ChartCard for loading/error/empty states:**

```tsx
<ChartCard
  title="Revenue"
  subtitle="Last 30 days"
  isLoading={isLoading}
  isEmpty={!data || data.series.length === 0}
  error={error}
  onRetry={refetch}
  actions={<RefreshButton onClick={refetch} />}
  height={320}
  chartType="line"
>
  <LineChart data={chartData} />
</ChartCard>
```

---

## 4. Statistics Components

### StatCard

**Use for:** Key metrics with trend indicators

```tsx
import { StatCard, StatCardGroup } from "@/components";
import { DollarSign } from "lucide-react";

<StatCardGroup cols={3}>
  <StatCard
    label="Total Revenue"
    value="$12,450"
    icon={<DollarSign className="h-5 w-5" />}
    trend="up"
    trendValue="+12.5%"
    variant="success"
  />
</StatCardGroup>;
```

**Variants:** `default`, `success`, `warning`, `error`, `info`

### ComplexStatCard

**Use for:** Multiple metrics, status badges

```tsx
<ComplexStatCard
  label="Cash on Hand"
  value="$25,000"
  comparisonText="vs last month: +$2,500"
  statusBadge={{ label: "Healthy", variant: "success" }}
  secondaryMetrics={[
    { label: "Cash In", value: "$30,000" },
    { label: "Cash Out", value: "$5,000" },
  ]}
/>
```

---

## 5. Theme Integration

### useChartTheme Hook

**Automatically applies daisyUI colors and RTL:**

```tsx
const { tokens, themedOptions, backgroundPlugin } = useChartTheme();

// tokens.primary → daisyUI primary color
// tokens.success → daisyUI success color
// themedOptions → Pre-configured Chart.js options (RTL-aware)
// backgroundPlugin → Applies bg-base-100 background
```

**Token Map:**

```typescript
{
  primary: string; // oklch(55% 0.3 180) → Teal
  secondary: string; // oklch(70% 0.15 30) → Gold accent
  success: string; // oklch(70% 0.2 150) → Green
  warning: string; // oklch(80% 0.2 85) → Yellow
  error: string; // oklch(65% 0.25 25) → Red
  info: string; // oklch(70% 0.15 230) → Blue
  primaryLight: string; // With 10% opacity
  // ... more
}
```

---

## 6. Data Transformers

### Backend Format → Chart.js Format

**TimeSeries → Line Chart:**

```typescript
import { transformTimeSeriesToChartData } from "@/lib/charts";

// Backend: { series: [{ label: '2024-01-01', value: 1500 }, ...] }
const chartData = transformTimeSeriesToChartData(
  apiResponse.revenueOverTime,
  "Revenue",
  tokens.primary
);
```

**KeyValue → Bar Chart:**

```typescript
import { transformKeyValueToBarData } from "@/lib/charts";

// Backend: [{ key: 'Electronics', value: 1200 }, ...]
const chartData = transformKeyValueToBarData(
  apiResponse.categorySales,
  "Sales",
  [tokens.primary, tokens.secondary, tokens.success]
);
```

**KeyValue → Pie Chart:**

```typescript
import { transformKeyValueToPieData } from "@/lib/charts";

const chartData = transformKeyValueToPieData(apiResponse.orderStatus, [
  tokens.success,
  tokens.warning,
  tokens.error,
]);
```

---

## 7. Utility Functions

```typescript
// Generate color palette with varying opacity
generateColorPalette(tokens.primary, 5); // Returns 5 colors

// Format currency in tooltips
formatChartCurrency(1500, "USD"); // "$1,500"

// Format large numbers
formatChartNumber(1500000); // "1.5M"

// Check if decimation needed
shouldEnableDecimation(dataPoints.length); // true if >100 points

// Calculate percentages
calculatePercentages([100, 200, 300]); // ["16.67%", "33.33%", "50%"]
```

---

## 8. RTL Support (Automatic)

**All charts auto-support RTL when language is Arabic:**

- X-axis reversal (latest data on left in RTL)
- Legend positioning (right in RTL)
- Tooltip text direction
- Number formatting (Western Arabic numerals 0-9)

**No configuration needed** - `useChartTheme()` handles it.

---

## 9. TanStack Query Integration

**Pattern:** Load all analytics data at once, distribute to multiple charts

```tsx
import { useQuery } from "@tanstack/react-query";

function AnalyticsPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["analytics", "sales", dateRange],
    queryFn: () => analyticsApi.getSalesAnalytics(dateRange),
  });

  if (isLoading) return <PageLoadingSkeleton />;
  if (error) return <ErrorBoundary error={error} />;

  return (
    <div className="space-y-6">
      {/* Stats */}
      <StatCardGroup cols={4}>
        <StatCard label="Total Revenue" value={data.totalRevenue} />
        <StatCard label="Orders" value={data.totalOrders} />
        {/* ... */}
      </StatCardGroup>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ChartCard title="Revenue Over Time">
          <LineChart data={transformTimeSeriesToChartData(data.revenue)} />
        </ChartCard>

        <ChartCard title="Top Products">
          <BarChart data={transformKeyValueToBarData(data.topProducts)} />
        </ChartCard>
      </div>
    </div>
  );
}
```

---

## 10. Best Practices

### ✅ Do

```tsx
// 1. Always wrap in ChartCard
<ChartCard title="Sales" isLoading={isLoading}>
  <LineChart data={data} />
</ChartCard>

// 2. Handle empty data
<ChartCard isEmpty={!data || data.series.length === 0}>
  <LineChart data={data} />
</ChartCard>

// 3. Use theme tokens
const { tokens } = useChartTheme()
borderColor: tokens.primary

// 4. Enable decimation for large datasets
<LineChart data={data} enableDecimation={data.length > 100} />

// 5. Use semantic variants
<StatCard variant="warning" />  // Yellow for warnings
<StatCard variant="success" />  // Green for positive metrics
```

### ❌ Don't

```tsx
// Don't render charts without ChartCard
<LineChart data={data} />  // Missing loading/error handling

// Don't hardcode colors
borderColor: '#0d9488'  // Use tokens.primary instead

// Don't skip empty states
// ChartCard handles it automatically

// Don't use large datasets without decimation
<LineChart data={thousandsOfPoints} />  // Enable decimation!
```

---

## 11. Mobile Considerations

**Height:**

- Default 320px (good for mobile)
- 400px for desktop-only views

**Responsive Grids:**

```tsx
<div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
  <ChartCard>...</ChartCard>
  <ChartCard>...</ChartCard>
</div>
```

**Touch Interactions:**

- Chart.js tooltips work with touch
- Use larger hit areas for legends
- Test on real mobile devices

---

## 12. Performance

### Large Datasets (>100 points)

```tsx
<LineChart
  data={largeData}
  enableDecimation // Automatically reduces points
/>
```

### Canvas vs SVG

Chart.js uses Canvas (better performance for complex charts).

### Lazy Loading

```tsx
import { lazy, Suspense } from 'react'

const HeavyChart = lazy(() => import('./HeavyChart'))

<Suspense fallback={<ChartSkeleton />}>
  <HeavyChart />
</Suspense>
```

---

## 13. Troubleshooting

### Charts not rendering

**Cause:** Chart.js not registered

**Solution:** Ensure registered in `main.tsx`:

```typescript
import { Chart, registerables } from "chart.js";
Chart.register(...registerables);
```

### Colors not matching theme

**Cause:** Hardcoded colors

**Solution:** Use `useChartTheme()` tokens

### RTL not working

**Cause:** Not using `useChartTheme()`

**Solution:**

```tsx
const { themedOptions } = useChartTheme()
<Line options={themedOptions} />
```

### Performance issues

**Cause:** Too many data points

**Solution:** Enable decimation or reduce data

---

## 14. File Structure

```
portal-web/src/
├── lib/charts/                  # Chart utilities
│   ├── chartTheme.ts            # useChartTheme hook
│   ├── chartUtils.ts            # Data transformers
│   ├── chartPlugins.ts          # Custom plugins
│   ├── rtlSupport.ts            # RTL helpers
│   └── index.ts                 # Public exports
├── components/atoms/
│   ├── LineChart.tsx
│   ├── BarChart.tsx
│   ├── PieChart.tsx
│   ├── DoughnutChart.tsx
│   ├── MixedChart.tsx
│   ├── ChartCard.tsx
│   ├── ChartSkeleton.tsx
│   ├── ChartEmptyState.tsx
│   ├── StatCard.tsx
│   ├── StatCardGroup.tsx
│   └── ComplexStatCard.tsx
└── i18n/
    ├── ar/analytics.json        # Arabic chart labels
    └── en/analytics.json        # English chart labels
```

---

## Agent Validation Checklist

Before completing chart task:

- ☑ Charts wrapped in ChartCard
- ☑ Using `useChartTheme()` for colors
- ☑ Empty/loading/error states handled
- ☑ Data transformers used (no manual Chart.js data construction)
- ☑ Decimation enabled for large datasets
- ☑ Mobile-friendly height (320px default)
- ☑ RTL: No hardcoded left/right positioning
- ☑ Stat cards use semantic variants
- ☑ Translation keys, not hardcoded labels
- ☑ Query keys from `queryKeys` factory

---

## See Also

- **Design Tokens:** `.github/instructions/design-tokens.instructions.md` → Color palette
- **UI Patterns:** `.github/instructions/ui-implementation.instructions.md` → Layout patterns
- **Architecture:** `.github/instructions/portal-web-architecture.instructions.md` → TanStack Query
- **Forms:** `.github/instructions/forms.instructions.md` → Date range pickers for filtering

---

## Resources

- Chart.js Docs: https://www.chartjs.org
- react-chartjs-2: https://react-chartjs-2.js.org
- Implementation: `portal-web/src/lib/charts/`, `portal-web/src/components/atoms/`
- Examples: `portal-web/src/routes/_app/dashboard.tsx`
