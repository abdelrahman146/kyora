import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries } from '@/api/onboarding'
import { OnboardingRouteError } from '@/features/onboarding/components/OnboardingRouteError'
import { PlanSelectionPage } from '@/features/onboarding/components/PlanSelectionPage'

const PlanSearchSchema = z.object({
  plan: z.string().optional(),
})

export const Route = createFileRoute('/onboarding/plan')({
  validateSearch: (search): z.infer<typeof PlanSearchSchema> => {
    return PlanSearchSchema.parse(search)
  },

  loader: async ({ context }) => {
    const { queryClient } = context as RouterContext
    await queryClient.ensureQueryData(onboardingQueries.plans())
  },

  component: PlanSelectionPage,

  errorComponent: OnboardingRouteError,
})
