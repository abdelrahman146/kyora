import { useTranslation } from 'react-i18next'
import type { ErrorComponentProps } from '@tanstack/react-router'

export function OnboardingRouteError({ error }: ErrorComponentProps) {
  const { t } = useTranslation('errors')

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <div className="card bg-base-100 border border-base-300 max-w-md">
        <div className="card-body">
          <h2 className="card-title text-error">{t('title')}</h2>
          <p className="text-base-content/70">
            {error instanceof Error
              ? error.message || t('generic.unexpected')
              : t('generic.unexpected')}
          </p>
        </div>
      </div>
    </div>
  )
}
