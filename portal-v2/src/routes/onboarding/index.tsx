import { Navigate, createFileRoute } from '@tanstack/react-router'
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

  // Redirect to home if already authenticated
  if (authState.isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return <OnboardingLayout />
}
