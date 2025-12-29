import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { AlertCircle, Loader2 } from 'lucide-react'
import { z } from 'zod'
import { onboardingStore, updateStage } from '@/stores/onboardingStore'
import { useOAuthGoogleMutation } from '@/api/onboarding'
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
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { code, error, error_description } = Route.useSearch()
  const storeState = useStore(onboardingStore)
  const [isProcessing, setIsProcessing] = useState(true)

  const oauthMutation = useOAuthGoogleMutation()

  useEffect(() => {
    void handleOAuthCallback()
  }, [])

  const handleOAuthCallback = async () => {
    // Handle OAuth errors
    if (error) {
      const errorMessage = error_description || error
      toast.error(t('auth:oauth_error', { error: errorMessage }))
      await navigate({ to: '/onboarding/verify', replace: true })
      return
    }

    // Validate required parameters
    if (!code) {
      toast.error(t('auth:oauth_invalid_callback'))
      await navigate({ to: '/onboarding/verify', replace: true })
      return
    }

    // Validate session token
    if (!storeState.sessionToken) {
      toast.error(t('onboarding:session_expired'))
      await navigate({ to: '/onboarding/email', replace: true })
      return
    }

    try {
      setIsProcessing(true)

      // Exchange authorization code for tokens
      const response = await oauthMutation.mutateAsync({
        sessionToken: storeState.sessionToken,
        code,
      })

      // Update onboarding stage
      updateStage(response.stage)

      toast.success(t('auth:oauth_success'))

      // Navigate to next step based on stage
      if (response.stage === 'identity_verified') {
        await navigate({ to: '/onboarding/business', replace: true })
      } else if (response.stage === 'business_staged') {
        // Business already set, go to payment or complete
        if (storeState.isPaidPlan) {
          await navigate({ to: '/onboarding/payment', replace: true })
        } else {
          await navigate({ to: '/onboarding/complete', replace: true })
        }
      } else {
        // Fallback to verify step
        await navigate({ to: '/onboarding/verify', replace: true })
      }
    } catch (err) {
      const message = await translateErrorAsync(err, t)
      toast.error(message)
      await navigate({ to: '/onboarding/verify', replace: true })
    } finally {
      setIsProcessing(false)
    }
  }

  return (
    <div className="max-w-lg mx-auto">
      <div className="text-center">
        {/* Loading State */}
        {isProcessing && (
          <>
            <div className="flex justify-center mb-6">
              <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
                <Loader2 className="w-8 h-8 text-primary animate-spin" />
              </div>
            </div>
            <h1 className="text-3xl font-bold text-base-content mb-2">
              {t('auth:processing_oauth')}
            </h1>
            <p className="text-base-content/70">
              {t('auth:processing_oauth_description')}
            </p>
          </>
        )}

        {/* Error State */}
        {!isProcessing && (error || oauthMutation.isError) && (
          <>
            <div className="flex justify-center mb-6">
              <div className="w-16 h-16 rounded-full bg-error/10 flex items-center justify-center">
                <AlertCircle className="w-8 h-8 text-error" />
              </div>
            </div>
            <h1 className="text-3xl font-bold text-base-content mb-2">
              {t('auth:oauth_failed')}
            </h1>
            <p className="text-base-content/70 mb-6">
              {error_description || t('auth:oauth_failed_description')}
            </p>
            <button
              onClick={() => navigate({ to: '/onboarding/verify' })}
              className="btn btn-primary"
            >
              {t('common:back_to_verification')}
            </button>
          </>
        )}
      </div>
    </div>
  )
}
