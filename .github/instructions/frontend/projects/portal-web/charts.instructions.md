---
description: Portal Web Chart.js patterns - Dashboard charts, ChartCard wrapper, StatCard components, theme integration (portal-web only)
applyTo: "portal-web/**"
---

# Portal Web Charts

Chart.js patterns for analytics and dashboards.

**Cross-refs:**

- General UI patterns: `../../_general/ui-patterns.instructions.md`
- Portal UI components: `./ui-components.instructions.md`

---

## 1. Chart.js Setup

### Dependencies

```json
{
  "dependencies": {
    "chart.js": "^4.4.1",
    "react-chartjs-2": "^5.2.0"
  }
}
```

### Global Chart.js Registration

```typescript
// src/lib/charts/config.ts
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
);

// Default options for all charts
export const defaultChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false, // Handle legends outside chart
    },
    tooltip: {
      backgroundColor: "rgba(0, 0, 0, 0.8)",
      padding: 12,
      cornerRadius: 8,
      titleFont: { size: 14, weight: "bold" as const },
      bodyFont: { size: 13 },
      displayColors: true,
    },
  },
  interaction: {
    mode: "index" as const,
    intersect: false,
  },
};
```

---

## 2. ChartCard Component

Reusable card wrapper for dashboard charts.

```tsx
// src/components/charts/ChartCard.tsx
import type { ReactNode } from "react";

interface ChartCardProps {
  title: string;
  description?: string;
  icon?: ReactNode;
  actions?: ReactNode; // Dropdown, export button, etc.
  children: ReactNode; // Chart component
  height?: string; // Default: "300px"
}

export function ChartCard({
  title,
  description,
  icon,
  actions,
  children,
  height = "300px",
}: ChartCardProps) {
  return (
    <div className="card bg-base-100 border border-base-300">
      <div className="card-body p-4">
        {/* Header */}
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center gap-2">
            {icon && <div className="text-primary">{icon}</div>}
            <div>
              <h3 className="font-semibold text-base">{title}</h3>
              {description && (
                <p className="text-sm text-base-content/70 mt-1">
                  {description}
                </p>
              )}
            </div>
          </div>
          {actions && <div>{actions}</div>}
        </div>

        {/* Chart */}
        <div style={{ height }}>{children}</div>
      </div>
    </div>
  );
}
```

**Usage:**

```tsx
<ChartCard
  title={t("analytics:revenue_trend")}
  description={t("analytics:last_30_days")}
  icon={<TrendingUp className="h-5 w-5" />}
  actions={
    <button className="btn btn-ghost btn-sm">
      <Download size={16} />
    </button>
  }
  height="350px"
>
  <LineChart data={revenueData} />
</ChartCard>
```

---

## 3. Line Chart

### Basic Line Chart

```tsx
// src/components/charts/LineChart.tsx
import { Line } from "react-chartjs-2";
import { defaultChartOptions } from "@/lib/charts/config";
import type { ChartData, ChartOptions } from "chart.js";

interface LineChartProps {
  data: ChartData<"line">;
  options?: ChartOptions<"line">;
}

export function LineChart({ data, options }: LineChartProps) {
  return (
    <Line
      data={data}
      options={{
        ...defaultChartOptions,
        scales: {
          x: {
            grid: { display: false },
            ticks: { color: "hsl(var(--bc) / 0.7)" },
          },
          y: {
            grid: { color: "hsl(var(--bc) / 0.1)" },
            ticks: { color: "hsl(var(--bc) / 0.7)" },
          },
        },
        ...options,
      }}
    />
  );
}
```

### Revenue Trend Example

```tsx
import { LineChart } from "@/components/charts/LineChart";
import { useTranslation } from "react-i18next";
import { TrendingUp } from "lucide-react";

function RevenueTrendChart() {
  const { t } = useTranslation(["analytics"]);

  const data = {
    labels: ["Week 1", "Week 2", "Week 3", "Week 4"],
    datasets: [
      {
        label: t("analytics:revenue"),
        data: [1200, 1900, 1500, 2200],
        borderColor: "hsl(var(--p))",
        backgroundColor: "hsl(var(--p) / 0.1)",
        tension: 0.4,
        fill: true,
      },
    ],
  };

  return (
    <ChartCard
      title={t("analytics:revenue_trend")}
      icon={<TrendingUp className="h-5 w-5" />}
    >
      <LineChart data={data} />
    </ChartCard>
  );
}
```

---

## 4. Bar Chart

### Basic Bar Chart

```tsx
// src/components/charts/BarChart.tsx
import { Bar } from "react-chartjs-2";
import { defaultChartOptions } from "@/lib/charts/config";
import type { ChartData, ChartOptions } from "chart.js";

interface BarChartProps {
  data: ChartData<"bar">;
  options?: ChartOptions<"bar">;
}

export function BarChart({ data, options }: BarChartProps) {
  return (
    <Bar
      data={data}
      options={{
        ...defaultChartOptions,
        scales: {
          x: {
            grid: { display: false },
            ticks: { color: "hsl(var(--bc) / 0.7)" },
          },
          y: {
            grid: { color: "hsl(var(--bc) / 0.1)" },
            ticks: { color: "hsl(var(--bc) / 0.7)" },
            beginAtZero: true,
          },
        },
        ...options,
      }}
    />
  );
}
```

### Orders by Status Example

```tsx
function OrdersByStatusChart() {
  const { t } = useTranslation(["analytics"]);

  const data = {
    labels: [
      t("analytics:pending"),
      t("analytics:processing"),
      t("analytics:shipped"),
      t("analytics:delivered"),
    ],
    datasets: [
      {
        label: t("analytics:orders"),
        data: [12, 19, 25, 30],
        backgroundColor: [
          "hsl(var(--wa) / 0.7)", // Warning (pending)
          "hsl(var(--in) / 0.7)", // Info (processing)
          "hsl(var(--p) / 0.7)", // Primary (shipped)
          "hsl(var(--su) / 0.7)", // Success (delivered)
        ],
      },
    ],
  };

  return (
    <ChartCard title={t("analytics:orders_by_status")}>
      <BarChart data={data} />
    </ChartCard>
  );
}
```

---

## 5. Doughnut Chart

### Basic Doughnut Chart

```tsx
// src/components/charts/DoughnutChart.tsx
import { Doughnut } from "react-chartjs-2";
import { defaultChartOptions } from "@/lib/charts/config";
import type { ChartData, ChartOptions } from "chart.js";

interface DoughnutChartProps {
  data: ChartData<"doughnut">;
  options?: ChartOptions<"doughnut">;
}

export function DoughnutChart({ data, options }: DoughnutChartProps) {
  return (
    <Doughnut
      data={data}
      options={{
        ...defaultChartOptions,
        ...options,
      }}
    />
  );
}
```

### Payment Methods Example

```tsx
function PaymentMethodsChart() {
  const { t } = useTranslation(["analytics"]);

  const data = {
    labels: [
      t("analytics:cash"),
      t("analytics:card"),
      t("analytics:bank_transfer"),
    ],
    datasets: [
      {
        data: [45, 35, 20],
        backgroundColor: [
          "hsl(var(--p))", // Primary
          "hsl(var(--s))", // Secondary
          "hsl(var(--a))", // Accent
        ],
        borderWidth: 2,
        borderColor: "hsl(var(--b1))",
      },
    ],
  };

  return (
    <ChartCard title={t("analytics:payment_methods_breakdown")}>
      <DoughnutChart data={data} />
    </ChartCard>
  );
}
```

---

## 6. RTL Support

Chart.js needs RTL configuration:

```typescript
// src/lib/charts/config.ts
import { useLanguage } from "@/hooks/useLanguage";

export function getChartOptions(isRTL: boolean) {
  return {
    ...defaultChartOptions,
    indexAxis: "x" as const, // Change to 'y' for horizontal bars
    scales: {
      x: {
        position: isRTL ? "right" : "left",
        grid: { display: false },
        ticks: { color: "hsl(var(--bc) / 0.7)" },
      },
      y: {
        grid: { color: "hsl(var(--bc) / 0.1)" },
        ticks: { color: "hsl(var(--bc) / 0.7)" },
        beginAtZero: true,
      },
    },
  };
}
```

**Usage:**

```tsx
function MyChart() {
  const { isRTL } = useLanguage();
  const options = getChartOptions(isRTL);

  return <BarChart data={data} options={options} />;
}
```

---

## 7. Responsive Patterns

### Mobile: Show Fewer Labels

```tsx
function ResponsiveChart() {
  const isMobile = window.innerWidth < 768;

  const data = {
    labels: isMobile
      ? ["Jan", "Feb", "Mar", "Apr"] // Short labels
      : ["January", "February", "March", "April"], // Full labels
    datasets: [
      {
        label: "Revenue",
        data: [1200, 1900, 1500, 2200],
        borderColor: "hsl(var(--p))",
      },
    ],
  };

  return <LineChart data={data} />;
}
```

### Adjust Height on Mobile

```tsx
<ChartCard
  title={t("analytics:revenue_trend")}
  height={isMobile ? "250px" : "350px"}
>
  <LineChart data={data} />
</ChartCard>
```

---

## 8. Loading & Empty States

### Loading State

```tsx
<ChartCard title={t("analytics:revenue_trend")}>
  {isLoading ? (
    <div className="flex items-center justify-center h-full">
      <span className="loading loading-spinner loading-lg"></span>
    </div>
  ) : (
    <LineChart data={data} />
  )}
</ChartCard>
```

### Empty State

```tsx
<ChartCard title={t("analytics:revenue_trend")}>
  {data.length === 0 ? (
    <div className="flex flex-col items-center justify-center h-full text-center">
      <BarChart3 className="h-12 w-12 text-base-content/30 mb-3" />
      <p className="text-base-content/70">{t("analytics:no_data_yet")}</p>
    </div>
  ) : (
    <LineChart data={data} />
  )}
</ChartCard>
```

---

## 9. Color System (daisyUI Integration)

Use daisyUI semantic colors for consistency:

| Purpose   | Color Variable   | Usage                 |
| --------- | ---------------- | --------------------- |
| Primary   | `hsl(var(--p))`  | Main brand color      |
| Secondary | `hsl(var(--s))`  | Secondary brand color |
| Accent    | `hsl(var(--a))`  | Accent highlights     |
| Success   | `hsl(var(--su))` | Positive values       |
| Warning   | `hsl(var(--wa))` | Cautionary values     |
| Error     | `hsl(var(--er))` | Negative values       |
| Info      | `hsl(var(--in))` | Informational values  |

**Example:**

```tsx
const data = {
  datasets: [
    {
      label: "Revenue",
      data: [1200, 1900, 1500],
      backgroundColor: [
        "hsl(var(--su) / 0.7)", // Success
        "hsl(var(--p) / 0.7)", // Primary
        "hsl(var(--wa) / 0.7)", // Warning
      ],
    },
  ],
};
```

---

## 10. Best Practices

### ✅ Do

- Use `ChartCard` for consistent layout
- Set explicit `height` prop (default 300px)
- Use daisyUI color variables for theme consistency
- Handle loading/empty/error states
- Show fewer labels on mobile
- Use translated labels (`t('analytics:...')`)
- Use RTL-aware axis positioning

### ❌ Don't

- Don't hardcode colors (use daisyUI variables)
- Don't skip loading/empty states
- Don't use tiny fonts (<11px) on mobile
- Don't overflow x-axis labels (rotate if needed)
- Don't render charts server-side (use `'use client'` in Next.js)

---

## Agent Validation

Before completing charts task:

- ☑ Using `ChartCard` wrapper
- ☑ Using daisyUI color variables (`hsl(var(--p))`)
- ☑ RTL-aware axis positioning
- ☑ Loading/empty/error states handled
- ☑ Mobile-friendly height and labels
- ☑ Translated chart labels
- ☑ Chart.js registered globally

---

## Resources

- Chart.js Docs: https://www.chartjs.org/docs/latest/
- react-chartjs-2: https://react-chartjs-2.js.org/
- Portal UI components: `./ui-components.instructions.md`
