import { useEffect } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { AlertCircle, Loader2 } from 'lucide-react'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries, useOAuthGoogleMutation } from '@/api/onboarding'

const OAuthCallbackSearchSchema = z.object({
  session: z.string().min(1),
  code: z.string().optional(),
  error: z.string().optional(),
  error_description: z.string().optional(),
})

export const Route = createFileRoute('/onboarding/oauth-callback')({
  validateSearch: (search): z.infer<typeof OAuthCallbackSearchSchema> => {
    return OAuthCallbackSearchSchema.parse(search)
  },

  loader: async ({ location, context }) => {
    const { queryClient } = context as unknown as RouterContext
    const parsed = OAuthCallbackSearchSchema.parse(location.search)

    // Preload session data
    await queryClient.ensureQueryData(onboardingQueries.session(parsed.session))
  },

  component: OAuthCallbackPage,

  errorComponent: ({ error }) => {
    const { t } = useTranslation('translation')
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="card bg-base-100 border border-base-300 max-w-md">
          <div className="card-body">
            <h2 className="card-title text-error">{t('error.title')}</h2>
            <p className="text-base-content/70">
              {error.message || t('error.generic')}
            </p>
          </div>
        </div>
      </div>
    )
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
  const navigate = useNavigate()
  const {
    session: sessionToken,
    code,
    error,
    error_description,
  } = Route.useSearch()

  useSuspenseQuery(onboardingQueries.session(sessionToken))
  const oauthMutation = useOAuthGoogleMutation()

  useEffect(() => {
    const handleCallback = async () => {
      const oauthError = error_description || error

      // Handle OAuth error
      if (oauthError) {
        setTimeout(() => {
          void navigate({
            to: '/onboarding/verify',
            search: { session: sessionToken },
            replace: true,
          })
        }, 3000)
        return
      }

      // Missing authorization code
      if (!code) {
        setTimeout(() => {
          void navigate({
            to: '/onboarding/verify',
            search: { session: sessionToken },
            replace: true,
          })
        }, 3000)
        return
      }

      // Process OAuth callback
      try {
        await oauthMutation.mutateAsync({
          sessionToken,
          code,
        })

        // Navigate to business step on success
        void navigate({
          to: '/onboarding/business',
          search: { session: sessionToken },
          replace: true,
        })
      } catch (err) {
        // Error is handled by mutation, redirect after delay
        setTimeout(() => {
          void navigate({
            to: '/onboarding/verify',
            search: { session: sessionToken },
            replace: true,
          })
        }, 3000)
      }
    }

    void handleCallback()
  }, [code, error, error_description, sessionToken, oauthMutation, navigate])

  const oauthError = error_description || error
  const hasError = !!oauthError || !code || oauthMutation.isError

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <div className="card bg-base-100 border border-base-300 max-w-md">
        <div className="card-body">
          <div className="text-center">
            {hasError ? (
              <>
                <div className="flex justify-center mb-4">
                  <div className="w-16 h-16 bg-error/10 rounded-full flex items-center justify-center">
                    <AlertCircle className="h-8 w-8 text-error" />
                  </div>
                </div>
                <h2 className="text-xl font-bold text-error mb-3">
                  {tOnboarding('oauth.errorTitle')}
                </h2>
                <p className="text-base-content/70 mb-4">
                  {oauthError ||
                    oauthMutation.error?.message ||
                    tOnboarding('oauth.noCode')}
                </p>
                <p className="text-sm text-base-content/50">
                  {tOnboarding('oauth.redirecting')}
                </p>
              </>
            ) : (
              <>
                <Loader2 className="w-16 h-16 text-primary animate-spin mx-auto mb-4" />
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
