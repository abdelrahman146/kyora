import { redirect } from '@tanstack/react-router'
import type { SessionStage } from '@/api/types/onboarding'

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

export function redirectToCorrectStage(
  currentPath: string,
  sessionStage: SessionStage,
  sessionToken: string,
) {
  const expectedRoute = STAGE_ROUTES[sessionStage]

  if (currentPath === expectedRoute) {
    return null
  }

  return redirect({
    to: expectedRoute,
    search: { session: sessionToken },
  })
}
