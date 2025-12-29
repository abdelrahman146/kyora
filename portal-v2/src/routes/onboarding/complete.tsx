import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { CheckCircle2, Loader2, PartyPopper } from 'lucide-react'
import { clearSession, onboardingStore } from '@/stores/onboardingStore'
import { setTokens } from '@/api/client'
import { useCompleteOnboardingMutation } from '@/api/onboarding'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/onboarding/complete')({
  component: CompletePage,
})

/**
 * Completion Step - Final step of Onboarding
 *
 * Features:
 * - Finalizes onboarding session
 * - Retrieves authentication tokens
 * - Clears onboarding state
 * - Redirects to main application
 */
function CompletePage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const state = useStore(onboardingStore)
  const [isCompleting, setIsCompleting] = useState(false)
  const [completionSteps, setCompletionSteps] = useState({
    finalizingAccount: false,
    configuringWorkspace: false,
    redirecting: false,
  })

  const completeMutation = useCompleteOnboardingMutation()

  // Redirect if no session
  useEffect(() => {
    if (!state.sessionToken) {
      navigate({ to: '/onboarding/email', replace: true })
      return
    }

    // Check if business data is set
    if (!state.businessData) {
      navigate({ to: '/onboarding/business', replace: true })
      return
    }

    // Check if paid plan requires payment
    if (state.isPaidPlan && !state.paymentCompleted) {
      navigate({ to: '/onboarding/payment', replace: true })
      return
    }
  }, [state, navigate])

  // Auto-complete onboarding
  useEffect(() => {
    if (state.sessionToken && !isCompleting) {
      void handleComplete()
    }
  }, [state.sessionToken])

  const handleComplete = async () => {
    if (isCompleting) return

    setIsCompleting(true)

    try {
      // Step 1: Finalizing account
      setCompletionSteps((prev) => ({ ...prev, finalizingAccount: true }))
      await new Promise((resolve) => setTimeout(resolve, 800))

      // Step 2: Call complete API
      const response = await completeMutation.mutateAsync({
        sessionToken: state.sessionToken!,
      })

      // Step 3: Configuring workspace
      setCompletionSteps((prev) => ({ ...prev, configuringWorkspace: true }))
      await new Promise((resolve) => setTimeout(resolve, 800))

      // Step 4: Set authentication tokens
      setTokens(response.token, response.refreshToken)

      // Step 5: Clear onboarding session
      clearSession()

      // Step 6: Redirecting
      setCompletionSteps((prev) => ({ ...prev, redirecting: true }))
      await new Promise((resolve) => setTimeout(resolve, 500))

      toast.success(t('onboarding:welcome_to_kyora'))

      // Navigate to main app
      await navigate({ to: '/', replace: true })
    } catch (error) {
      setIsCompleting(false)
      const message = await translateErrorAsync(error, t)
      toast.error(message)
      
      // Allow retry
      setCompletionSteps({
        finalizingAccount: false,
        configuringWorkspace: false,
        redirecting: false,
      })
    }
  }

  const handleRetry = () => {
    void handleComplete()
  }

  return (
    <div className="max-w-lg mx-auto">
      <div className="text-center">
        {/* Success Icon */}
        <div className="flex justify-center mb-6">
          <div className="w-20 h-20 rounded-full bg-success/10 flex items-center justify-center">
            {isCompleting ? (
              <Loader2 className="w-10 h-10 text-success animate-spin" />
            ) : (
              <PartyPopper className="w-10 h-10 text-success" />
            )}
          </div>
        </div>

        {/* Header */}
        <h1 className="text-4xl font-bold text-base-content mb-3">
          {isCompleting
            ? t('onboarding:completing_setup')
            : t('onboarding:setup_complete')}
        </h1>
        <p className="text-base-content/70 text-lg mb-8">
          {isCompleting
            ? t('onboarding:completing_description')
            : t('onboarding:setup_complete_description')}
        </p>

        {/* Progress Steps */}
        {isCompleting && (
          <div className="space-y-4 mb-8">
            {/* Step 1: Finalizing Account */}
            <div className="flex items-center justify-between p-4 bg-base-200 rounded-lg">
              <div className="flex items-center gap-3">
                {completionSteps.finalizingAccount ? (
                  <CheckCircle2 className="w-5 h-5 text-success" />
                ) : (
                  <Loader2 className="w-5 h-5 text-base-content/40 animate-spin" />
                )}
                <span
                  className={
                    completionSteps.finalizingAccount
                      ? 'text-base-content'
                      : 'text-base-content/60'
                  }
                >
                  {t('onboarding:finalizing_account')}
                </span>
              </div>
            </div>

            {/* Step 2: Configuring Workspace */}
            <div className="flex items-center justify-between p-4 bg-base-200 rounded-lg">
              <div className="flex items-center gap-3">
                {completionSteps.configuringWorkspace ? (
                  <CheckCircle2 className="w-5 h-5 text-success" />
                ) : (
                  <Loader2
                    className={`w-5 h-5 ${
                      completionSteps.finalizingAccount
                        ? 'text-base-content/40 animate-spin'
                        : 'text-base-content/20'
                    }`}
                  />
                )}
                <span
                  className={
                    completionSteps.configuringWorkspace
                      ? 'text-base-content'
                      : 'text-base-content/60'
                  }
                >
                  {t('onboarding:configuring_workspace')}
                </span>
              </div>
            </div>

            {/* Step 3: Redirecting */}
            <div className="flex items-center justify-between p-4 bg-base-200 rounded-lg">
              <div className="flex items-center gap-3">
                {completionSteps.redirecting ? (
                  <CheckCircle2 className="w-5 h-5 text-success" />
                ) : (
                  <Loader2
                    className={`w-5 h-5 ${
                      completionSteps.configuringWorkspace
                        ? 'text-base-content/40 animate-spin'
                        : 'text-base-content/20'
                    }`}
                  />
                )}
                <span
                  className={
                    completionSteps.redirecting
                      ? 'text-base-content'
                      : 'text-base-content/60'
                  }
                >
                  {t('onboarding:redirecting_to_app')}
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Retry Button (shown on error) */}
        {!isCompleting && completeMutation.isError && (
          <button onClick={handleRetry} className="btn btn-primary">
            {t('common:retry')}
          </button>
        )}

        {/* Fun Message */}
        {!isCompleting && (
          <div className="alert alert-success mt-6">
            <PartyPopper className="w-5 h-5" />
            <span>{t('onboarding:ready_to_start')}</span>
          </div>
        )}
      </div>
    </div>
  )
}
