import { useCallback, useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { AlertCircle, CheckCircle2, CreditCard } from 'lucide-react'
import type { Plan } from '@/api/types/onboarding'
import { onboardingApi } from '@/api/onboarding'
import { onboardingStore, setCheckoutUrl } from '@/stores/onboardingStore'
import { translateErrorAsync } from '@/lib/translateError'

const PaymentSearchSchema = z.object({
  status: z.enum(['success', 'cancelled']).optional(),
})

export const Route = createFileRoute('/onboarding/payment')({
  validateSearch: (search): z.infer<typeof PaymentSearchSchema> => {
    return PaymentSearchSchema.parse(search)
  },
  component: PaymentPage,
})

/**
 * Payment Step - Step 5 of Onboarding (paid plans only)
 *
 * Features:
 * - Stripe Checkout integration
 * - Payment status polling
 * - Automatic navigation after successful payment
 */
function PaymentPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const { t: tTranslation } = useTranslation('translation')
  const navigate = useNavigate()
  const state = useStore(onboardingStore)
  const { status } = Route.useSearch()

  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [paymentStatus, setPaymentStatus] = useState<
    'pending' | 'success' | 'cancelled' | null
  >(null)

  // Load selected plan details
  useEffect(() => {
    const loadSelectedPlan = async () => {
      if (!state.planDescriptor) return
      try {
        const plan = await onboardingApi.getPlan(state.planDescriptor)
        setSelectedPlan(plan)
      } catch {
        // Non-blocking
      }
    }

    void loadSelectedPlan()
  }, [state.planDescriptor])

  // Check URL params for payment status
  useEffect(() => {
    if (status === 'success') {
      setPaymentStatus('success')
      setTimeout(() => {
        void navigate({ to: '/onboarding/complete' })
      }, 2000)
      return
    }

    if (status === 'cancelled') {
      setPaymentStatus('cancelled')
    }
  }, [navigate, status])

  // Redirect if not at business_staged or payment_pending stage
  useEffect(() => {
    if (!state.sessionToken) {
      void navigate({ to: '/onboarding/plan', replace: true })
      return
    }

    if (!state.isPaidPlan) {
      void navigate({ to: '/onboarding/complete', replace: true })
      return
    }

    if (
      state.stage !== 'business_staged' &&
      state.stage !== 'payment_pending' &&
      state.stage !== 'ready_to_commit'
    ) {
      void navigate({ to: '/onboarding/business', replace: true })
    }
  }, [navigate, state.isPaidPlan, state.sessionToken, state.stage])

  const initiatePayment = useCallback(async () => {
    if (!state.sessionToken) return

    try {
      setIsLoading(true)
      setError('')

      const successUrl = `${window.location.origin}/onboarding/payment?status=success`
      const cancelUrl = `${window.location.origin}/onboarding/payment?status=cancelled`

      const response = await onboardingApi.startPayment({
        sessionToken: state.sessionToken,
        successUrl,
        cancelUrl,
      })

      setCheckoutUrl(response.checkoutUrl)
      window.location.href = response.checkoutUrl
    } catch (err) {
      const message = await translateErrorAsync(err, tTranslation)
      setError(message)
    } finally {
      setIsLoading(false)
    }
  }, [state.sessionToken, tTranslation])

  // Auto-initiate payment if no checkout URL exists
  useEffect(() => {
    if (
      !state.checkoutUrl &&
      !paymentStatus &&
      !error &&
      state.sessionToken &&
      state.isPaidPlan
    ) {
      void initiatePayment()
    }
  }, [error, initiatePayment, paymentStatus, state.checkoutUrl, state.isPaidPlan, state.sessionToken])

  if (paymentStatus === 'success') {
    return (
      <div className="max-w-lg mx-auto">
        <div className="card bg-base-100 border border-success shadow-xl">
          <div className="card-body">
            <div className="text-center">
              <div className="flex justify-center mb-4">
                <div className="w-16 h-16 bg-success/10 rounded-full flex items-center justify-center">
                  <CheckCircle2 className="w-8 h-8 text-success" />
                </div>
              </div>
              <h2 className="text-2xl font-bold text-success mb-3">
                  {tOnboarding('payment.successTitle')}
              </h2>
              <p className="text-base-content/70">
                  {tOnboarding('payment.successMessage')}
              </p>
              <div className="mt-6">
                <span className="loading loading-spinner loading-sm"></span>
                <span className="ms-2 text-sm text-base-content/60">
                    {tOnboarding('payment.redirecting')}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (paymentStatus === 'cancelled') {
    return (
      <div className="max-w-lg mx-auto">
        <div className="card bg-base-100 border border-warning shadow-xl">
          <div className="card-body">
            <div className="text-center">
              <div className="flex justify-center mb-4">
                <div className="w-16 h-16 bg-warning/10 rounded-full flex items-center justify-center">
                  <AlertCircle className="w-8 h-8 text-warning" />
                </div>
              </div>
              <h2 className="text-2xl font-bold text-warning mb-3">
                {tOnboarding('payment.cancelledTitle')}
              </h2>
              <p className="text-base-content/70 mb-6">
                {tOnboarding('payment.cancelledMessage')}
              </p>
              <div className="flex gap-3 justify-center">
                <button
                  onClick={() => {
                    void navigate({ to: '/onboarding/plan' })
                  }}
                  className="btn btn-ghost"
                >
                  {tOnboarding('payment.changePlan')}
                </button>
                <button
                  onClick={() => {
                    void initiatePayment()
                  }}
                  className="btn btn-primary"
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <>
                      <span className="loading loading-spinner loading-sm"></span>
                      {tCommon('loading')}
                    </>
                  ) : (
                    tOnboarding('payment.tryAgain')
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-lg mx-auto">
      <div className="card bg-base-100 border border-base-300 shadow-xl">
        <div className="card-body">
          <div className="text-center mb-6">
            <div className="flex justify-center mb-4">
              <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                <CreditCard className="w-8 h-8 text-primary" />
              </div>
            </div>
            <h2 className="text-2xl font-bold">{tOnboarding('payment.title')}</h2>
            <p className="text-base-content/70 mt-2">
              {tOnboarding('payment.subtitle')}
            </p>
          </div>

          {selectedPlan && (
            <div className="bg-base-200 rounded-lg p-4 mb-6">
              <h3 className="font-semibold text-lg mb-2">{selectedPlan.name}</h3>
              <div className="flex items-baseline gap-1">
                <span className="text-3xl font-bold text-primary">
                  {selectedPlan.price}
                </span>
                <span className="text-base-content/60">
                  {selectedPlan.currency.toUpperCase()}
                </span>
                <span className="text-base-content/60">/ {selectedPlan.billingCycle}</span>
              </div>
              {selectedPlan.description && (
                <p className="text-sm text-base-content/70 mt-2">
                  {selectedPlan.description}
                </p>
              )}
            </div>
          )}

          {error && (
            <div className="alert alert-error mb-4">
              <span className="text-sm">{error}</span>
            </div>
          )}

          <div className="space-y-3 mb-6">
            <div className="flex items-center gap-3 text-sm">
              <CheckCircle2 className="w-5 h-5 text-success shrink-0" />
              <span>{tOnboarding('payment.securePayment')}</span>
            </div>
            <div className="flex items-center gap-3 text-sm">
              <CheckCircle2 className="w-5 h-5 text-success shrink-0" />
              <span>{tOnboarding('payment.cancelAnytime')}</span>
            </div>
            <div className="flex items-center gap-3 text-sm">
              <CheckCircle2 className="w-5 h-5 text-success shrink-0" />
              <span>{tOnboarding('payment.instantAccess')}</span>
            </div>
          </div>

          <button
            onClick={() => {
              void initiatePayment()
            }}
            className="btn btn-primary w-full"
            disabled={isLoading}
          >
            {isLoading ? (
              <>
                <span className="loading loading-spinner loading-sm"></span>
                {tCommon('loading')}
              </>
            ) : (
              tOnboarding('payment.proceedToPayment')
            )}
          </button>

          <p className="text-xs text-center text-base-content/50 mt-4">
            {tOnboarding('payment.poweredByStripe')}
          </p>
        </div>
      </div>
    </div>
  )
}
