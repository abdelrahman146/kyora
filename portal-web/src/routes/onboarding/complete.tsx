import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries } from '@/api/onboarding'
import { CompleteOnboardingPage } from '@/features/onboarding/components/CompleteOnboardingPage'
import { OnboardingRouteError } from '@/features/onboarding/components/OnboardingRouteError'

const CompleteSearchSchema = z.object({
  session: z.string().min(1),
})

export const Route = createFileRoute('/onboarding/complete')({
  validateSearch: (search): z.infer<typeof CompleteSearchSchema> => {
    return CompleteSearchSchema.parse(search)
  },

  beforeLoad: async ({ location, context }) => {
    const { queryClient } = context as RouterContext
    const parsed = CompleteSearchSchema.parse(location.search)

    const session = await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.session),
    )

    if (session.stage !== 'ready_to_commit') {
      throw redirect({
        to: '/onboarding/plan',
        replace: true,
      })
    }
  },

  loader: async ({ location, context }) => {
    const { queryClient } = context as unknown as RouterContext
    const parsed = CompleteSearchSchema.parse(location.search)

    await queryClient.ensureQueryData(onboardingQueries.session(parsed.session))
  },

  component: CompleteOnboardingPage,

  errorComponent: OnboardingRouteError,
})
