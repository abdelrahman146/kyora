import { createFileRoute } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { LineChart, Package } from 'lucide-react'

export const Route = createFileRoute(
  '/business/$businessDescriptor/analytics/sales',
)({
  staticData: {
    titleKey: 'pages.analytics_sales',
  },
  component: SalesAnalyticsPage,
})

function SalesAnalyticsPage() {
  const { t } = useTranslation('common')

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <LineChart className="size-8 text-primary" />
        <h1 className="text-2xl font-bold">{t('pages.analytics_sales')}</h1>
      </div>

      <div className="card bg-base-200 border border-base-300">
        <div className="card-body text-center py-12">
          <Package className="size-16 mx-auto text-base-content/30 mb-4" />
          <h2 className="text-xl font-semibold text-base-content/70">
            {t('coming_soon')}
          </h2>
          <p className="text-base-content/60 max-w-md mx-auto">
            Sales analytics with insights on revenue trends, best-selling
            products, and growth opportunities.
          </p>
        </div>
      </div>
    </div>
  )
}
