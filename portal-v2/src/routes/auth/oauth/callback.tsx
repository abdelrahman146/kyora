import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { AlertCircle, CheckCircle, Loader2 } from 'lucide-react'
import { authApi } from '@/api/auth'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/auth/oauth/callback')({
  component: OAuthCallbackPage,
  validateSearch: (search: Record<string, unknown>) => ({
    code: (search.code as string) || '',
    state: (search.state as string) || '',
    error: (search.error as string) || '',
    error_description: (search.error_description as string) || '',
  }),
})

function OAuthCallbackPage() {
  const navigate = useNavigate()
  const { code, state, error, error_description } = useSearch({
    from: '/auth/oauth/callback',
  })
  const { t } = useTranslation()

  const [status, setStatus] = useState<'loading' | 'success' | 'error'>(
    'loading',
  )
  const [errorMessage, setErrorMessage] = useState<string>('')

  useEffect(() => {
    const handleCallback = async () => {
      try {
        if (error) {
          setStatus('error')
          setErrorMessage(
            error_description || t('auth.oauth_error', { error }),
          )
          return
        }

        if (!code) {
          setStatus('error')
          setErrorMessage(t('auth.oauth_missing_code'))
          return
        }

        // Parse state to get redirect destination
        let redirectTo = '/'
        try {
          if (state) {
            const stateData = JSON.parse(decodeURIComponent(state)) as {
              from?: string
              redirect?: string
            }
            redirectTo = stateData.from ?? stateData.redirect ?? '/'
          }
        } catch {
          // Invalid state, use default redirect
        }

        // Exchange code for tokens (this automatically saves tokens and user)
        await authApi.loginWithGoogle({ code })

        setStatus('success')

        void setTimeout(() => {
          void navigate({ to: redirectTo as never, replace: true })
        }, 1000)
      } catch (err) {
        setStatus('error')
        const message = await translateErrorAsync(err, t)
        setErrorMessage(message)
      }
    }

    void handleCallback()
  }, [code, state, error, error_description, navigate, t])

  const handleRetry = async () => {
    await navigate({ to: '/auth/login', replace: true })
  }

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
      <div className="card w-full max-w-md bg-base-200 shadow-xl">
        <div className="card-body items-center text-center">
          {status === 'loading' && (
            <>
              <Loader2 className="w-16 h-16 text-primary animate-spin" />
              <h2 className="card-title text-2xl mt-4">
                {t('auth.oauth_processing')}
              </h2>
              <p className="text-base-content/60">
                {t('auth.oauth_processing_description')}
              </p>
            </>
          )}

          {status === 'success' && (
            <>
              <div className="w-16 h-16 rounded-full bg-success/20 flex items-center justify-center">
                <CheckCircle className="w-10 h-10 text-success" />
              </div>
              <h2 className="card-title text-2xl mt-4 text-success">
                {t('auth.oauth_success')}
              </h2>
              <p className="text-base-content/60">
                {t('auth.oauth_success_description')}
              </p>
            </>
          )}

          {status === 'error' && (
            <>
              <div className="w-16 h-16 rounded-full bg-error/20 flex items-center justify-center">
                <AlertCircle className="w-10 h-10 text-error" />
              </div>
              <h2 className="card-title text-2xl mt-4 text-error">
                {t('auth.oauth_failed')}
              </h2>
              <p className="text-base-content/60">{errorMessage}</p>
              <div className="card-actions mt-6">
                <button
                  onClick={() => {
                    void handleRetry()
                  }}
                  className="btn btn-primary"
                  type="button"
                >
                  {t('auth.return_to_login')}
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
