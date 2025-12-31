import { Navigate, createFileRoute, useRouterState } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { authStore } from '@/stores/authStore'

export const Route = createFileRoute('/onboarding/')({
  component: OnboardingRoot,
})

/**
 * Onboarding Root - Redirect Handler
 *
 * Handles the base /onboarding route by:
 * - Redirecting authenticated users to home
 * - Redirecting unauthenticated users to /onboarding/plan (first step)
 *
 * Route Structure:
 * /onboarding → redirects to /onboarding/plan
 *   /plan - Plan selection (Step 1)
 *   /email - Email entry (Step 2)
 *   /verify - Email verification (Step 3)
 *   /business - Business setup (Step 4)
 *   /payment - Payment checkout (Step 5 - for paid plans)
 *   /complete - Finalization and welcome (Final step)
 *   /oauth-callback - OAuth callback handler
 */
function OnboardingRoot() {
  const { t } = useTranslation('common')
  const authState = useStore(authStore)

  const pathname = useRouterState({
    select: (s) => s.location.pathname,
  })

  // Redirect to home if already authenticated
  // Exception: allow /onboarding/complete to render its success screen
  // after we set auth tokens, matching portal-web UX.
  if (authState.isAuthenticated && pathname !== '/onboarding/complete') {
    return <Navigate to="/" replace />
  }

  // Redirect to first step of onboarding
  return (
    <>
      <Navigate to="/onboarding/plan" replace />
      {/* Show loading state during redirect */}
      <div className="min-h-screen flex items-center justify-center bg-base-100">
        <div className="text-center">
          <Loader2 className="w-12 h-12 animate-spin text-primary mx-auto mb-4" />
          <p className="text-lg text-base-content/70">{t('loading')}...</p>
        </div>
      </div>
    </>
  )
}
