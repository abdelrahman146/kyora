import { createFileRoute } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { Briefcase, Package } from 'lucide-react'

export const Route = createFileRoute('/settings/workspace')({
  staticData: {
    titleKey: 'pages.workspace_settings',
  },
  component: WorkspaceSettingsPage,
})

function WorkspaceSettingsPage() {
  const { t } = useTranslation('common')

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Briefcase className="size-8 text-primary" />
        <h1 className="text-2xl font-bold">{t('pages.workspace_settings')}</h1>
      </div>

      <div className="card bg-base-200 border border-base-300">
        <div className="card-body text-center py-12">
          <Package className="size-16 mx-auto text-base-content/30 mb-4" />
          <h2 className="text-xl font-semibold text-base-content/70">
            {t('coming_soon')}
          </h2>
          <p className="text-base-content/60 max-w-md mx-auto">
            Manage workspace members, permissions, and collaboration settings.
          </p>
        </div>
      </div>
    </div>
  )
}
