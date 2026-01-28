import { createFileRoute } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { Building2, Package } from 'lucide-react'

export const Route = createFileRoute(
  '/business/$businessDescriptor/settings/business',
)({
  staticData: {
    titleKey: 'pages.business_management',
  },
  component: BusinessManagementPage,
})

function BusinessManagementPage() {
  const { t } = useTranslation('common')

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Building2 className="size-8 text-primary" />
        <h1 className="text-2xl font-bold">{t('pages.business_management')}</h1>
      </div>

      <div className="card bg-base-200 border border-base-300">
        <div className="card-body text-center py-12">
          <Package className="size-16 mx-auto text-base-content/30 mb-4" />
          <h2 className="text-xl font-semibold text-base-content/70">
            {t('coming_soon')}
          </h2>
          <p className="text-base-content/60 max-w-md mx-auto">
            Manage your business details, storefront settings, and operational
            preferences.
          </p>
        </div>
      </div>
    </div>
  )
}
