import { Navigate, createFileRoute, useRouterState } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { authStore } from '@/stores/authStore'
import { OnboardingLayout } from '@/components/templates/OnboardingLayout'

export const Route = createFileRoute('/onboarding/')({
  component: OnboardingRoot,
})

/**
 * Onboarding Root Layout
 *
 * Wraps all onboarding steps with:
 * - OnboardingLayout for consistent UI
 * - Authentication guard (redirect if already logged in)
 * - Outlet for nested routes
 *
 * Route Structure:
 * /onboarding
 *   /plan - Plan selection
 *   /email - Email entry
 *   /verify - Email verification
 *   /business - Business setup
 *   /payment - Payment checkout (for paid plans)
 *   /complete - Finalization and welcome
 *   /oauth-callback - OAuth callback handler
 */
function OnboardingRoot() {
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

  return (
    <OnboardingLayout>
      {/* Nested routes will render here */}
      <div>Redirecting...</div>
    </OnboardingLayout>
  )
}
