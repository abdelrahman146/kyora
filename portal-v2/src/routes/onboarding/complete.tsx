import { useCallback, useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import { ArrowRight, Sparkles } from 'lucide-react'
import { onboardingApi } from '@/api/onboarding'
import { setTokens } from '@/api/client'
import { authStore, setUser } from '@/stores/authStore'
import { clearSession, loadSessionFromStorage, onboardingStore } from '@/stores/onboardingStore'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/onboarding/complete')({
  component: CompletePage,
})

function CompletePage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const navigate = useNavigate()
  const onboardingState = useStore(onboardingStore)
  const authState = useStore(authStore)

  const [isCompleting, setIsCompleting] = useState(false)
  const [isComplete, setIsComplete] = useState(false)
  const [error, setError] = useState('')
  const [hasValidatedSession, setHasValidatedSession] = useState(false)

  // Restore onboarding session from storage on mount (if needed)
  useEffect(() => {
    if (isComplete) return

    let isCancelled = false

    const restoreSession = async () => {
      if (onboardingState.sessionToken) {
        setHasValidatedSession(true)
        return
      }

      const hasSession = await loadSessionFromStorage()
      if (isCancelled) return

      setHasValidatedSession(true)

      if (!hasSession) {
        if (authState.isAuthenticated) {
          void navigate({ to: '/', replace: true })
          return
        }

        void navigate({ to: '/onboarding/plan', replace: true })
      }
    }

    void restoreSession()
    return () => {
      isCancelled = true
    }
  }, [authState.isAuthenticated, isComplete, navigate, onboardingState.sessionToken])

  // Redirect if prerequisites not met
  useEffect(() => {
    if (!hasValidatedSession) return
    if (isCompleting || isComplete) return

    if (!onboardingState.sessionToken) {
      void navigate({ to: '/onboarding/plan', replace: true })
      return
    }

    if (
      onboardingState.isPaidPlan &&
      !onboardingState.paymentCompleted &&
      onboardingState.stage !== 'ready_to_commit'
    ) {
      void navigate({ to: '/onboarding/payment', replace: true })
      return
    }

    if (
      !onboardingState.isPaidPlan &&
      onboardingState.stage !== 'business_staged' &&
      onboardingState.stage !== 'ready_to_commit'
    ) {
      void navigate({ to: '/onboarding/business', replace: true })
    }
  }, [hasValidatedSession, isComplete, isCompleting, navigate, onboardingState.isPaidPlan, onboardingState.paymentCompleted, onboardingState.sessionToken, onboardingState.stage])

  const completeOnboarding = useCallback(async () => {
    if (!onboardingState.sessionToken) return

    try {
      setIsCompleting(true)
      setError('')

      const response = await onboardingApi.complete({
        sessionToken: onboardingState.sessionToken,
      })

      // Hydrate auth immediately to avoid protected-route flicker
      setTokens(response.token, response.refreshToken)
      setUser(response.user)

      setIsComplete(true)

      setTimeout(() => {
        void (async () => {
          await clearSession()
          await navigate({ to: '/', replace: true })
        })()
      }, 3000)
    } catch (err) {
      const message = await translateErrorAsync(err, tOnboarding)
      setError(message)
    } finally {
      setIsCompleting(false)
    }
  }, [navigate, onboardingState.sessionToken, tOnboarding])

  // Auto-complete on mount if ready
  useEffect(() => {
    if (!hasValidatedSession) return
    if (!onboardingState.sessionToken) return
    if (isComplete || isCompleting || error) return

    const isReady =
      (onboardingState.isPaidPlan && onboardingState.paymentCompleted) ||
      !onboardingState.isPaidPlan

    if (!isReady) return

    void completeOnboarding()
  }, [completeOnboarding, error, hasValidatedSession, isComplete, isCompleting, onboardingState.isPaidPlan, onboardingState.paymentCompleted, onboardingState.sessionToken])

  if (error) {
    return (
      <div className="max-w-lg mx-auto">
        <div className="card bg-base-100 border border-error shadow-xl">
          <div className="card-body">
            <div className="text-center">
              <h2 className="text-2xl font-bold text-error mb-3">
                {tOnboarding('complete.errorTitle')}
              </h2>
              <p className="text-base-content/70 mb-6">{error}</p>
              <button
                onClick={() => {
                  void completeOnboarding()
                }}
                className="btn btn-primary"
                disabled={isCompleting}
              >
                {isCompleting ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {tCommon('loading')}
                  </>
                ) : (
                  tCommon('retry')
                )}
              </button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (isComplete) {
    const businessName = onboardingState.businessData?.name ?? ''

    return (
      <div className="max-w-2xl mx-auto">
        <div className="card bg-linear-to-br from-primary/10 to-secondary/10 border-2 border-primary shadow-2xl">
          <div className="card-body">
            <div className="text-center">
              <div className="flex justify-center mb-6">
                <div className="relative">
                  <div className="w-24 h-24 bg-linear-to-br from-primary to-secondary rounded-full flex items-center justify-center animate-bounce">
                    <Sparkles className="w-12 h-12 text-primary-content" />
                  </div>
                  <div className="absolute inset-0 w-24 h-24 bg-linear-to-br from-primary to-secondary rounded-full animate-ping opacity-20"></div>
                </div>
              </div>

              <h1 className="text-4xl font-bold mb-3">
                {tOnboarding('complete.welcomeTitle')}
              </h1>
              <p className="text-xl text-base-content/80 mb-6">
                {tOnboarding('complete.welcomeMessage', { businessName })}
              </p>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 my-8">
                <div className="bg-base-100/50 backdrop-blur rounded-lg p-4">
                  <div className="text-2xl mb-2">📦</div>
                  <h3 className="font-semibold mb-1">
                    {tOnboarding('complete.feature1Title')}
                  </h3>
                  <p className="text-sm text-base-content/70">
                    {tOnboarding('complete.feature1Desc')}
                  </p>
                </div>
                <div className="bg-base-100/50 backdrop-blur rounded-lg p-4">
                  <div className="text-2xl mb-2">📊</div>
                  <h3 className="font-semibold mb-1">
                    {tOnboarding('complete.feature2Title')}
                  </h3>
                  <p className="text-sm text-base-content/70">
                    {tOnboarding('complete.feature2Desc')}
                  </p>
                </div>
                <div className="bg-base-100/50 backdrop-blur rounded-lg p-4">
                  <div className="text-2xl mb-2">🚀</div>
                  <h3 className="font-semibold mb-1">
                    {tOnboarding('complete.feature3Title')}
                  </h3>
                  <p className="text-sm text-base-content/70">
                    {tOnboarding('complete.feature3Desc')}
                  </p>
                </div>
              </div>

              <div className="flex items-center justify-center gap-2 text-base-content/60">
                <span className="loading loading-spinner loading-sm"></span>
                <span>{tOnboarding('complete.redirecting')}</span>
                <ArrowRight className="w-4 h-4 animate-pulse" />
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
          <div className="text-center">
            <div className="flex justify-center mb-6">
              <span className="loading loading-spinner loading-lg text-primary"></span>
            </div>
            <h2 className="text-2xl font-bold mb-3">
              {tOnboarding('complete.processingTitle')}
            </h2>
            <p className="text-base-content/70">
              {tOnboarding('complete.processingMessage')}
            </p>
            <div className="mt-6 space-y-2">
              <div className="flex items-center gap-3 justify-center text-sm">
                <span className="loading loading-ring loading-sm text-success"></span>
                <span>{tOnboarding('complete.creatingWorkspace')}</span>
              </div>
              <div className="flex items-center gap-3 justify-center text-sm">
                <span className="loading loading-ring loading-sm text-success"></span>
                <span>{tOnboarding('complete.settingUpBusiness')}</span>
              </div>
              <div className="flex items-center gap-3 justify-center text-sm">
                <span className="loading loading-ring loading-sm text-success"></span>
                <span>{tOnboarding('complete.preparingDashboard')}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
