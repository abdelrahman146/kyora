import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries } from '@/api/onboarding'
import { OnboardingRouteError } from '@/features/onboarding/components/OnboardingRouteError'
import { redirectToCorrectStage } from '@/features/onboarding/utils/onboarding'
import { PaymentPage } from '@/features/onboarding/components/PaymentPage'

const PaymentSearchSchema = z.object({
  session: z.string().min(1),
  status: z.enum(['success', 'cancelled']).optional(),
})

export const Route = createFileRoute('/onboarding/payment')({
  validateSearch: (search): z.infer<typeof PaymentSearchSchema> => {
    return PaymentSearchSchema.parse(search)
  },

  beforeLoad: async ({ location, context }) => {
    const { queryClient } = context as unknown as RouterContext
    const parsed = PaymentSearchSchema.parse(location.search)

    // Ensure session data is loaded
    const session = await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.session),
    )

    // Validate stage and plan requirements
    if (!session.isPaidPlan) {
      throw redirect({
        to: '/onboarding/complete',
        search: { session: parsed.session },
        replace: true,
      })
    }

    if (
      session.stage !== 'business_staged' &&
      session.stage !== 'payment_pending' &&
      session.stage !== 'payment_confirmed'
    ) {
      throw redirect({
        to: '/onboarding/plan',
        replace: true,
      })
    }

    // If payment already confirmed, go to complete
    if (
      session.stage === 'payment_confirmed' ||
      session.paymentStatus === 'succeeded'
    ) {
      throw redirect({
        to: '/onboarding/complete',
        search: { session: parsed.session },
        replace: true,
      })
    }
  },

  loader: async ({ location, context }) => {
    const { queryClient } = context as unknown as RouterContext
    const parsed = PaymentSearchSchema.parse(location.search)

    // Load session data
    const session = await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.session),
    )

    // Automatically redirect to correct stage based on session
    const stageRedirect = redirectToCorrectStage(
      '/onboarding/payment',
      session.stage,
      parsed.session,
    )
    if (stageRedirect) {
      throw stageRedirect
    }

    // Load plan details
    await queryClient.ensureQueryData(onboardingQueries.plans())
  },

  component: PaymentPage,

  errorComponent: OnboardingRouteError,
})
