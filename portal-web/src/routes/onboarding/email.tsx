import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries } from '@/api/onboarding'
import { EmailEntryPage } from '@/features/onboarding/components/EmailEntryPage'

// Search params schema for URL-driven state
const EmailSearchSchema = z.object({
  plan: z.string().min(1),
})

export const Route = createFileRoute('/onboarding/email')({
  validateSearch: (search): z.infer<typeof EmailSearchSchema> => {
    return EmailSearchSchema.parse(search)
  },
  // Prefetch plan data before rendering
  loader: async ({ context, location }) => {
    const { queryClient } = context as RouterContext
    const parsed = EmailSearchSchema.parse(location.search)

    // Redirect if no plan selected
    if (!parsed.plan) {
      throw redirect({ to: '/onboarding/plan' })
    }

    // Prefetch plan details for summary card
    const plan = await queryClient.ensureQueryData(
      onboardingQueries.plan(parsed.plan),
    )

    return { plan }
  },
  component: EmailEntryPage,
})
