import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries } from '@/api/onboarding'
import { redirectToCorrectStage } from '@/features/onboarding/utils/onboarding'
import { VerifyEmailPage } from '@/features/onboarding/components/VerifyEmailPage'

// Search params schema
const VerifySearchSchema = z.object({
  session: z.string().min(1),
})

export const Route = createFileRoute('/onboarding/verify')({
  validateSearch: (search): z.infer<typeof VerifySearchSchema> => {
    return VerifySearchSchema.parse(search)
  },
  loader: async ({ context, location }) => {
    const parsed = VerifySearchSchema.parse(location.search)
    const { queryClient } = context as unknown as RouterContext

    // Redirect if no session token
    if (!parsed.session) {
      throw redirect({ to: '/onboarding/plan' })
    }

    // Prefetch and validate session
    const session = await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.session),
    )

    // Automatically redirect to correct stage based on session
    const stageRedirect = redirectToCorrectStage(
      '/onboarding/verify',
      session.stage,
      parsed.session,
    )
    if (stageRedirect) {
      throw stageRedirect
    }

    return { session }
  },
  component: VerifyEmailPage,
})
