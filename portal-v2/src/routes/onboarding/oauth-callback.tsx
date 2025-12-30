import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import { AlertCircle, Loader2 } from 'lucide-react'
import { z } from 'zod'
import { onboardingApi } from '@/api/onboarding'
import { loadSession, onboardingStore } from '@/stores/onboardingStore'
import { translateErrorAsync } from '@/lib/translateError'

const OAuthCallbackSearchSchema = z.object({
  code: z.string().optional(),
  error: z.string().optional(),
  error_description: z.string().optional(),
})

export const Route = createFileRoute('/onboarding/oauth-callback')({
  component: OAuthCallbackPage,
  validateSearch: (search): z.infer<typeof OAuthCallbackSearchSchema> => {
    return OAuthCallbackSearchSchema.parse(search)
  },
})

/**
 * OAuth Callback Handler - Handles OAuth redirects during onboarding
 *
 * Features:
 * - Processes OAuth authorization code
 * - Updates onboarding session with OAuth data
 * - Handles OAuth errors
 * - Redirects to appropriate onboarding step
 */
function OAuthCallbackPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tTranslation } = useTranslation('translation')
  const navigate = useNavigate()
  const { code, error, error_description } = Route.useSearch()
  const storeState = useStore(onboardingStore)
  const [isProcessing, setIsProcessing] = useState(true)
  const [errorMessage, setErrorMessage] = useState('')

  useEffect(() => {
    const handleCallback = async () => {
      const oauthError = error_description || error

      if (oauthError) {
        setErrorMessage(tOnboarding('oauth.error', { error: oauthError }))
        setTimeout(() => {
          void navigate({ to: '/onboarding/verify', replace: true })
        }, 3000)
        return
      }

      if (!code) {
        setErrorMessage(tOnboarding('oauth.noCode'))
        setTimeout(() => {
          void navigate({ to: '/onboarding/verify', replace: true })
        }, 3000)
        return
      }

      const storedToken =
        storeState.sessionToken ??
        sessionStorage.getItem('kyora_onboarding_google_session')

      if (!storedToken) {
        setErrorMessage(tOnboarding('oauth.noSession'))
        setTimeout(() => {
          void navigate({ to: '/onboarding/plan', replace: true })
        }, 3000)
        return
      }

      try {
        await onboardingApi.oauthGoogle({
          sessionToken: storedToken,
          code,
        })

        sessionStorage.removeItem('kyora_onboarding_google_session')
        await loadSession(storedToken)

        void navigate({ to: '/onboarding/business', replace: true })
      } catch (err) {
        const message = await translateErrorAsync(err, tTranslation)
        setErrorMessage(message)
        setTimeout(() => {
          void navigate({ to: '/onboarding/verify', replace: true })
        }, 3000)
      } finally {
        setIsProcessing(false)
      }
    }

    void handleCallback()
  }, [code, error, error_description, navigate, storeState.sessionToken, tOnboarding, tTranslation])

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <div className="card bg-base-100 border border-base-300 shadow-xl max-w-md">
        <div className="card-body">
          <div className="text-center">
            {errorMessage ? (
              <>
                <div className="flex justify-center mb-4">
                  <div className="w-16 h-16 bg-error/10 rounded-full flex items-center justify-center">
                    <AlertCircle className="h-8 w-8 text-error" />
                  </div>
                </div>
                <h2 className="text-xl font-bold text-error mb-3">
                  {tOnboarding('oauth.errorTitle')}
                </h2>
                <p className="text-base-content/70 mb-4">{errorMessage}</p>
                <p className="text-sm text-base-content/50">
                  {tOnboarding('oauth.redirecting')}
                </p>
              </>
            ) : (
              <>
                {isProcessing && (
                  <span className="loading loading-spinner loading-lg text-primary mb-4"></span>
                )}
                <h2 className="text-xl font-bold mb-2">
                  {tOnboarding('oauth.processing')}
                </h2>
                <p className="text-base-content/70">
                  {tOnboarding('oauth.pleaseWait')}
                </p>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
