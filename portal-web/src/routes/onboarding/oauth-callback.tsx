import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries } from '@/api/onboarding'
import { OnboardingRouteError } from '@/features/onboarding/components/OnboardingRouteError'
import { OAuthCallbackPage } from '@/features/onboarding/components/OAuthCallbackPage'

const OAuthCallbackSearchSchema = z.object({
  session: z.string().min(1),
  code: z.string().optional(),
  error: z.string().optional(),
  error_description: z.string().optional(),
})

export const Route = createFileRoute('/onboarding/oauth-callback')({
  validateSearch: (search): z.infer<typeof OAuthCallbackSearchSchema> => {
    return OAuthCallbackSearchSchema.parse(search)
  },

  loader: async ({ location, context }) => {
    const { queryClient } = context as unknown as RouterContext
    const parsed = OAuthCallbackSearchSchema.parse(location.search)

    // Preload session data
    await queryClient.ensureQueryData(onboardingQueries.session(parsed.session))
  },

  component: OAuthCallbackPage,

  errorComponent: OnboardingRouteError,
})
