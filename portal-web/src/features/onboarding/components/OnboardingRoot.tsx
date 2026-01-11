import { Navigate, useRouterState } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { authStore } from '@/stores/authStore'

export function OnboardingRoot() {
  const { t } = useTranslation('common')
  const authState = useStore(authStore)

  const pathname = useRouterState({
    select: (s) => s.location.pathname,
  })

  if (authState.isAuthenticated && pathname !== '/onboarding/complete') {
    return <Navigate to="/" replace />
  }

  return (
    <>
      <Navigate to="/onboarding/plan" replace />
      <div className="min-h-screen flex items-center justify-center bg-base-100">
        <div className="text-center">
          <Loader2 className="w-12 h-12 animate-spin text-primary mx-auto mb-4" />
          <p className="text-lg text-base-content/70">{t('loading')}...</p>
        </div>
      </div>
    </>
  )
}
