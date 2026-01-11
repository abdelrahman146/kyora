import { useEffect } from 'react'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { AlertCircle, CreditCard, Loader2 } from 'lucide-react'

import { onboardingQueries, useStartPaymentMutation } from '@/api/onboarding'
import { formatCurrency } from '@/lib/formatCurrency'

import { OnboardingLayout } from '@/features/onboarding/components/OnboardingLayout'

export function PaymentPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const navigate = useNavigate()

  const { session: sessionToken, status } = useSearch({
    from: '/onboarding/payment',
  })

  const { data: session } = useSuspenseQuery(
    onboardingQueries.session(sessionToken),
  )
  const { data: plans } = useSuspenseQuery(onboardingQueries.plans())

  const startPaymentMutation = useStartPaymentMutation()

  const selectedPlan = plans.find(
    (p) => p.descriptor === session.planDescriptor,
  )
  const isFree = selectedPlan ? parseFloat(selectedPlan.price) === 0 : false

  useEffect(() => {
    if (status === 'success') {
      void navigate({
        to: '/onboarding/complete',
        search: { session: sessionToken },
        replace: true,
      })
    } else if (status === 'cancelled') {
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
      window.location.href = result.checkoutUrl
    } catch (error) {
      console.error('[Payment] Failed to start payment:', error)
    }
  }

  if (status === 'cancelled') {
    return (
      <OnboardingLayout>
        <div className="max-w-2xl mx-auto">
          <div className="alert alert-warning mb-6">
            <AlertCircle className="w-5 h-5" />
            <div>
              <h3 className="font-semibold">
                {tOnboarding('payment.cancelled')}
              </h3>
              <p className="text-sm">{tOnboarding('payment.cancelledDesc')}</p>
            </div>
          </div>

          {selectedPlan && (
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">
                  {tOnboarding('payment.title')}
                </h2>

                <div className="bg-base-200 rounded-lg p-6 mb-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h3 className="text-xl font-semibold">
                        {selectedPlan.name}
                      </h3>
                      {selectedPlan.description && (
                        <p className="text-sm text-base-content/70 mt-1">
                          {selectedPlan.description}
                        </p>
                      )}
                    </div>
                    <div className="text-end">
                      <div className="text-3xl font-bold">
                        {formatCurrency(
                          parseFloat(selectedPlan.price),
                          selectedPlan.currency,
                        )}
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

  return (
    <OnboardingLayout>
      <div className="max-w-2xl mx-auto">
        {selectedPlan && (
          <div className="card bg-base-100 border border-base-300">
            <div className="card-body">
              <h2 className="card-title text-2xl mb-4">
                {tOnboarding('payment.title')}
              </h2>

              <div className="bg-base-200 rounded-lg p-6 mb-6">
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h3 className="text-xl font-semibold">
                      {selectedPlan.name}
                    </h3>
                    {selectedPlan.description && (
                      <p className="text-sm text-base-content/70 mt-1">
                        {selectedPlan.description}
                      </p>
                    )}
                  </div>
                  <div className="text-end">
                    <div className="text-3xl font-bold">
                      {formatCurrency(
                        parseFloat(selectedPlan.price),
                        selectedPlan.currency,
                      )}
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
