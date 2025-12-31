import { useEffect } from 'react'
import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { AlertCircle, CreditCard, Loader2 } from 'lucide-react'
import type { RouterContext } from '@/router'
import { onboardingQueries, useStartPaymentMutation } from '@/api/onboarding'
import { formatCurrency } from '@/lib/formatCurrency'
import { OnboardingLayout } from '@/components/templates/OnboardingLayout'
import { redirectToCorrectStage } from '@/lib/onboarding'

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
      onboardingQueries.session(parsed.session)
    )

    // Validate stage and plan requirements
    if (!session.isPaidPlan) {
      throw redirect({
        to: '/onboarding/complete',
        search: { session: parsed.session },
        replace: true,
      })
    }

    if (session.stage !== 'business_staged' && session.stage !== 'payment_pending' && session.stage !== 'payment_confirmed') {
      throw redirect({
        to: '/onboarding/plan',
        replace: true,
      })
    }

    // If payment already confirmed, go to complete
    if (session.stage === 'payment_confirmed' || session.paymentStatus === 'succeeded') {
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
      onboardingQueries.session(parsed.session)
    )

    // Automatically redirect to correct stage based on session
    const stageRedirect = redirectToCorrectStage(
      '/onboarding/payment',
      session.stage,
      parsed.session
    )
    if (stageRedirect) {
      throw stageRedirect
    }
    
    // Load plan details
    await queryClient.ensureQueryData(onboardingQueries.plans())
  },

  component: PaymentPage,
  
  errorComponent: ({ error }) => {
    const { t } = useTranslation('translation')
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="card bg-base-100 border border-base-300 shadow-xl max-w-md">
          <div className="card-body">
            <h2 className="card-title text-error">{t('error.title')}</h2>
            <p className="text-base-content/70">{error.message || t('error.generic')}</p>
          </div>
        </div>
      </div>
    )
  },
})

/**
 * Payment Step - Step 5 of Onboarding (paid plans only)
 *
 * Features:
 * - Stripe Checkout integration
 * - Payment status handling
 * - Automatic navigation after successful payment
 */
function PaymentPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const navigate = useNavigate()
  const { session: sessionToken, status } = Route.useSearch()
  
  const { data: session } = useSuspenseQuery(onboardingQueries.session(sessionToken))
  const { data: plans } = useSuspenseQuery(onboardingQueries.plans())
  
  const startPaymentMutation = useStartPaymentMutation()

  // Find the selected plan
  const selectedPlan = plans.find(p => p.descriptor === session.planDescriptor)
  const isFree = selectedPlan ? parseFloat(selectedPlan.price) === 0 : false

  // Handle payment status from URL
  useEffect(() => {
    if (status === 'success') {
      // Navigate to complete step
      void navigate({
        to: '/onboarding/complete',
        search: { session: sessionToken },
        replace: true,
      })
    } else if (status === 'cancelled') {
      // User cancelled payment, stay on this page
      // Remove status from URL
      void navigate({
        to: '/onboarding/payment',
        search: { session: sessionToken },
        replace: true,
      })
    }
  }, [status, sessionToken, navigate])

  const handleStartPayment = async () => {
    try {
      const result = await startPaymentMutation.mutateAsync({
        sessionToken,
        successUrl: `${window.location.origin}/onboarding/payment?session=${sessionToken}&status=success`,
        cancelUrl: `${window.location.origin}/onboarding/payment?session=${sessionToken}&status=cancelled`,
      })
      // Redirect to Stripe Checkout
      window.location.href = result.checkoutUrl
    } catch (error) {
      // Error handling is done by the mutation
      console.error('[Payment] Failed to start payment:', error)
    }
  }

  // If payment was cancelled
  if (status === 'cancelled') {
    return (
      <OnboardingLayout>
        <div className="max-w-2xl mx-auto">
          <div className="alert alert-warning mb-6">
            <AlertCircle className="w-5 h-5" />
            <div>
              <h3 className="font-semibold">{tOnboarding('payment.cancelled')}</h3>
              <p className="text-sm">{tOnboarding('payment.cancelledDesc')}</p>
          </div>
        </div>

        {selectedPlan && (
          <div className="card bg-base-100 border border-base-300 shadow-xl">
            <div className="card-body">
              <h2 className="card-title text-2xl mb-4">
                {tOnboarding('payment.title')}
              </h2>

              <div className="bg-base-200 rounded-lg p-6 mb-6">
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h3 className="text-xl font-semibold">{selectedPlan.name}</h3>
                    {selectedPlan.description && (
                      <p className="text-sm text-base-content/70 mt-1">
                        {selectedPlan.description}
                      </p>
                    )}
                  </div>
                  <div className="text-end">
                    <div className="text-3xl font-bold">
                      {formatCurrency(parseFloat(selectedPlan.price), selectedPlan.currency)}
                    </div>
                    <div className="text-sm text-base-content/60">
                      / {selectedPlan.billingCycle}
                    </div>
                  </div>
                </div>
              </div>

              <button
                onClick={handleStartPayment}
                disabled={startPaymentMutation.isPending || isFree}
                className="btn btn-primary btn-lg btn-block"
              >
                {startPaymentMutation.isPending ? (
                  <>
                    <Loader2 className="w-5 h-5 animate-spin" />
                    {tCommon('loading')}
                  </>
                ) : (
                  <>
                    <CreditCard className="w-5 h-5" />
                    {tOnboarding('payment.payNow')}
                  </>
                )}
              </button>

              {startPaymentMutation.isError && (
                <div className="alert alert-error mt-4">
                  <AlertCircle className="w-5 h-5" />
                  <span>{startPaymentMutation.error.message}</span>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
      </OnboardingLayout>
    )
  }

  // Regular payment page
  return (
    <OnboardingLayout>
      <div className="max-w-2xl mx-auto">
        {selectedPlan && (
        <div className="card bg-base-100 border border-base-300 shadow-xl">
          <div className="card-body">
            <h2 className="card-title text-2xl mb-4">
              {tOnboarding('payment.title')}
            </h2>

            <div className="bg-base-200 rounded-lg p-6 mb-6">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h3 className="text-xl font-semibold">{selectedPlan.name}</h3>
                  {selectedPlan.description && (
                    <p className="text-sm text-base-content/70 mt-1">
                      {selectedPlan.description}
                    </p>
                  )}
                </div>
                <div className="text-end">
                  <div className="text-3xl font-bold">
                    {formatCurrency(parseFloat(selectedPlan.price), selectedPlan.currency)}
                  </div>
                  <div className="text-sm text-base-content/60">
                    / {selectedPlan.billingCycle}
                  </div>
                </div>
              </div>
            </div>

            <button
              onClick={handleStartPayment}
              disabled={startPaymentMutation.isPending || isFree}
              className="btn btn-primary btn-lg btn-block"
            >
              {startPaymentMutation.isPending ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  {tCommon('loading')}
                </>
              ) : (
                <>
                  <CreditCard className="w-5 h-5" />
                  {tOnboarding('payment.payNow')}
                </>
              )}
            </button>

            {startPaymentMutation.isError && (
              <div className="alert alert-error mt-4">
                <AlertCircle className="w-5 h-5" />
                <span>{startPaymentMutation.error.message}</span>
              </div>
            )}
          </div>
        </div>
      )}
      </div>
    </OnboardingLayout>
  )
}
