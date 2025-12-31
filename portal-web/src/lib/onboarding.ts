import { redirect } from '@tanstack/react-router'
import type { SessionStage } from '@/api/types/onboarding'

/**
 * Maps onboarding session stages to their corresponding routes
 */
export const STAGE_ROUTES: Record<SessionStage, string> = {
  plan_selected: '/onboarding/email',
  identity_pending: '/onboarding/verify',
  identity_verified: '/onboarding/business',
  business_staged: '/onboarding/payment',
  payment_pending: '/onboarding/payment',
  payment_confirmed: '/onboarding/complete',
  ready_to_commit: '/onboarding/complete',
  committed: '/onboarding/complete',
}

/**
 * Redirects to the correct onboarding step based on session stage
 * If the current route matches the expected stage route, no redirect occurs
 *
 * @param currentPath - The current route path (e.g., '/onboarding/business')
 * @param sessionStage - The current session stage
 * @param sessionToken - The session token to include in search params
 * @returns A redirect response if stage mismatch, otherwise null
 */
export function redirectToCorrectStage(
  currentPath: string,
  sessionStage: SessionStage,
  sessionToken: string,
) {
  const expectedRoute = STAGE_ROUTES[sessionStage]

  // If we're already on the correct route, don't redirect
  if (currentPath === expectedRoute) {
    return null
  }

  // Redirect to the correct stage with the session token
  return redirect({
    to: expectedRoute,
    search: { session: sessionToken },
  })
}
