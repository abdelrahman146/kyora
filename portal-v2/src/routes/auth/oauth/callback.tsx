import { useEffect } from 'react'
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { Loader2 } from 'lucide-react'
import { authApi } from '@/api/auth'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/auth/oauth/callback')({
  component: OAuthCallbackPage,
  validateSearch: (search: Record<string, unknown>) => ({
    code: (search.code as string) || '',
    state: (search.state as string) || '',
  }),
})

function OAuthCallbackPage() {
  const navigate = useNavigate()
  const { code, state } = useSearch({ from: '/auth/oauth/callback' })
  const { t } = useTranslation()

  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        if (!code) {
          toast.error(t('auth:oauth_error'))
          await navigate({ to: '/auth/login', search: { redirect: '/' } })
          return
        }

        // Parse state to get redirect destination
        let redirect = '/'
        try {
          if (state) {
            const stateData = JSON.parse(decodeURIComponent(state))
            redirect = stateData.redirect || '/'
          }
        } catch {
          // Invalid state, use default redirect
        }

        // Exchange code for tokens - this will set tokens in memory
        await authApi.loginWithGoogle({ code })

        toast.success(t('auth:login_success'))
        await navigate({ to: redirect, search: { redirect } })
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
        await navigate({ to: '/auth/login', search: { redirect: '/' } })
      }
    }

    void handleOAuthCallback()
  }, [code, state, navigate, t])

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-6">
      <div className="flex flex-col items-center gap-4">
        <Loader2 className="w-12 h-12 text-primary animate-spin" />
        <div className="text-center">
          <h2 className="text-xl font-semibold text-base-content mb-2">
            {t('auth:completing_authentication')}
          </h2>
          <p className="text-base-content/60">{t('common:please_wait')}</p>
        </div>
      </div>
    </div>
  )
}
