import { useState } from 'react'
import { Link, createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import type { LoginFormData } from '@/schemas/auth'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { authStore, login } from '@/stores/authStore'
import { authApi } from '@/api/auth'
import { translateErrorAsync } from '@/lib/translateError'
import { LoginForm } from '@/components/organisms/LoginForm'
import { LanguageSwitcher } from '@/components/molecules/LanguageSwitcher'

export const Route = createFileRoute('/auth/login')({
  beforeLoad: redirectIfAuthenticated,
  component: LoginPage,
  validateSearch: (search: Record<string, unknown>) => ({
    redirect: (search.redirect as string) || '/',
  }),
})

function LoginPage() {
  const authState = useStore(authStore)
  const navigate = useNavigate()
  const { redirect } = useSearch({ from: '/auth/login' })
  const { t } = useTranslation()
  const [isGoogleLoading, setIsGoogleLoading] = useState(false)
  const [googleErrorMessage, setGoogleErrorMessage] = useState<string>('')

  // Don't show loading spinner during login - LoginForm handles its own loading state
  // Only show loading when initializing auth on page load
  const isInitializing = authState.isLoading && !authState.user && !authState.isAuthenticated

  if (isInitializing) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-base-100">
        <div className="flex flex-col items-center gap-4">
          <span className="loading loading-spinner loading-lg text-primary"></span>
          <p className="text-base-content/60 text-sm">{t('common.loading')}</p>
        </div>
      </div>
    )
  }

  const handleLogin = async (data: LoginFormData) => {
    setGoogleErrorMessage('')

    // The error will be handled by LoginForm component
    // We only navigate on success
    await login(data.email, data.password)
    await navigate({ to: redirect as never, replace: true })
  }

  const handleGoogleLogin = async () => {
    try {
      setIsGoogleLoading(true)
      setGoogleErrorMessage('')

      // Get OAuth URL from backend
      const { url } = await authApi.getGoogleAuthUrl()

      // Prepare state parameter with redirect destination
      const state = encodeURIComponent(JSON.stringify({ from: redirect }))
      const oauthUrl = url.includes('state=') ? url : `${url}&state=${state}`

      // Redirect to Google OAuth
      window.location.href = oauthUrl
    } catch (error) {
      setIsGoogleLoading(false)
      const message = await translateErrorAsync(error, t)
      setGoogleErrorMessage(message)
    }
  }

  return (
    <div className="min-h-screen bg-base-100 flex flex-col lg:flex-row">
      {/* Left Side - Branding (Hidden on mobile) */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary to-primary-focus items-center justify-center p-12">
        <div className="max-w-md text-center">
          <h1 className="text-5xl font-bold text-primary-content mb-6">
            Kyora
          </h1>
          <p className="text-xl text-primary-content/80 leading-relaxed">
            {t('auth.login_welcome_message')}
          </p>
        </div>
      </div>

      {/* Right Side - Login Form */}
      <div className="flex-1 flex items-center justify-center p-6 lg:p-12">
        <div className="w-full max-w-md">
          {/* Mobile Logo */}
          <div className="text-center mb-8 lg:hidden">
            <h1 className="text-4xl font-bold text-primary mb-2">Kyora</h1>
            <p className="text-base-content/60">{t('auth.login_subtitle')}</p>
          </div>

          {/* Page Title */}
          <div className="mb-8">
            <h2 className="text-3xl font-bold text-base-content mb-2">
              {t('auth.welcome_back')}
            </h2>
            <p className="text-base-content/60">
              {t('auth.login_description')}
            </p>
          </div>

          {googleErrorMessage ? (
            <div role="alert" className="alert alert-error mb-6">
              <span>{googleErrorMessage}</span>
            </div>
          ) : null}

          {/* Login Form */}
          <LoginForm
            onSubmit={handleLogin}
            onGoogleLogin={() => {
              void handleGoogleLogin()
            }}
            isGoogleLoading={isGoogleLoading}
          />

          {/* Sign Up Link */}
          <div className="mt-8 text-center">
            <p className="text-base-content/60">
              {t('auth.no_account')}{' '}
              <Link
                to="/onboarding"
                className="text-primary hover:text-primary-focus font-semibold hover:underline transition-colors"
              >
                {t('auth.sign_up')}
              </Link>
            </p>
          </div>

          {/* Language Switcher */}
          <div className="mt-8 text-center">
            <LanguageSwitcher variant="toggle" showLabel />
          </div>
        </div>
      </div>
    </div>
  )
}
