import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { ArrowLeft, Check, CreditCard, Loader2 } from 'lucide-react'
import { onboardingStore, setCheckoutUrl, setPaymentCompleted, updateStage } from '@/stores/onboardingStore'
import { useOnboardingSessionQuery, useStartPaymentMutation } from '@/api/onboarding'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/onboarding/payment')({
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
  const { t } = useTranslation()
  const navigate = useNavigate()
  const state = useStore(onboardingStore)
  const [isRedirecting, setIsRedirecting] = useState(false)
  
  const startPaymentMutation = useStartPaymentMutation()
  
  // Poll session to check payment status
  const { data: session } = useOnboardingSessionQuery(state.sessionToken)

  // Redirect if no session or not business staged
  useEffect(() => {
    if (!state.sessionToken) {
      navigate({ to: '/onboarding/email', replace: true })
    } else if (!state.businessData) {
      navigate({ to: '/onboarding/business', replace: true })
    } else if (!state.isPaidPlan) {
      // Free plan shouldn't reach payment step
      navigate({ to: '/onboarding/complete', replace: true })
    }
  }, [state.sessionToken, state.businessData, state.isPaidPlan, navigate])

  // Handle payment status changes
  useEffect(() => {
    if (session) {
      updateStage(session.stage)
      
      if (session.stage === 'payment_confirmed' || session.stage === 'ready_to_commit') {
        setPaymentCompleted(true)
        toast.success(t('onboarding:payment_success'))
        void navigate({ to: '/onboarding/complete' })
      }
    }
  }, [session, navigate, t])

  // Start payment flow
  const handleStartPayment = async () => {
    try {
      setIsRedirecting(true)
      const response = await startPaymentMutation.mutateAsync({
        sessionToken: state.sessionToken!,
      })
      
      setCheckoutUrl(response.checkoutUrl)
      
      // Redirect to Stripe Checkout
      window.location.href = response.checkoutUrl
    } catch (error) {
      setIsRedirecting(false)
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    }
  }

  // If already has checkout URL and payment is pending, show polling UI
  const isPolling = state.checkoutUrl && state.stage === 'payment_pending'
  const isPaymentConfirmed = state.paymentCompleted || session?.stage === 'payment_confirmed'

  if (isPolling || isRedirecting) {
    return (
      <div className="max-w-lg mx-auto">
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <div className="w-16 h-16 rounded-full bg-warning/10 flex items-center justify-center">
              <Loader2 className="w-8 h-8 text-warning animate-spin" />
            </div>
          </div>
          <h1 className="text-3xl font-bold text-base-content mb-2">
            {isRedirecting ? t('onboarding:redirecting_to_payment') : t('onboarding:waiting_for_payment')}
          </h1>
          <p className="text-base-content/70 mb-6">
            {isRedirecting 
              ? t('onboarding:redirecting_description')
              : t('onboarding:payment_polling_description')
            }
          </p>
          {!isRedirecting && state.checkoutUrl && (
            <a
              href={state.checkoutUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="btn btn-primary"
            >
              {t('onboarding:return_to_payment')}
            </a>
          )}
        </div>
      </div>
    )
  }

  if (isPaymentConfirmed) {
    return (
      <div className="max-w-lg mx-auto">
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <div className="w-16 h-16 rounded-full bg-success/10 flex items-center justify-center">
              <Check className="w-8 h-8 text-success" />
            </div>
          </div>
          <h1 className="text-3xl font-bold text-base-content mb-2">
            {t('onboarding:payment_confirmed')}
          </h1>
          <p className="text-base-content/70 mb-6">
            {t('onboarding:payment_confirmed_description')}
          </p>
          <button
            onClick={() => navigate({ to: '/onboarding/complete' })}
            className="btn btn-primary"
          >
            {t('common:continue')}
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-lg mx-auto">
      {/* Header */}
      <div className="text-center mb-8">
        <div className="flex justify-center mb-4">
          <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
            <CreditCard className="w-8 h-8 text-primary" />
          </div>
        </div>
        <h1 className="text-3xl font-bold text-base-content mb-2">
          {t('onboarding:complete_payment')}
        </h1>
        <p className="text-base-content/70">
          {t('onboarding:payment_description')}
        </p>
      </div>

      {/* Plan Summary */}
      <div className="card bg-base-200 mb-6">
        <div className="card-body">
          <h2 className="card-title text-lg">{t('onboarding:selected_plan')}</h2>
          <div className="flex justify-between items-center">
            <span className="text-base-content/70">{state.planDescriptor}</span>
            <span className="text-2xl font-bold text-primary">
              {/* Price would come from plan data */}
              {t('onboarding:view_pricing')}
            </span>
          </div>
        </div>
      </div>

      {/* Actions */}
      <div className="space-y-4">
        <button
          onClick={handleStartPayment}
          disabled={startPaymentMutation.isPending || isRedirecting}
          className="btn btn-primary w-full"
        >
          {startPaymentMutation.isPending && (
            <Loader2 className="w-4 h-4 animate-spin" />
          )}
          {t('onboarding:proceed_to_payment')}
        </button>

        <button
          type="button"
          onClick={() => navigate({ to: '/onboarding/business' })}
          className="btn btn-ghost w-full"
        >
          <ArrowLeft className="w-4 h-4" />
          {t('common:back')}
        </button>
      </div>

      {/* Payment Info */}
      <div className="alert alert-info mt-6">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          className="stroke-current shrink-0 w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <span>{t('onboarding:secure_payment_info')}</span>
      </div>
    </div>
  )
}
