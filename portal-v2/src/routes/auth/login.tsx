import { useState } from 'react'
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { Eye, EyeOff, Loader2 } from 'lucide-react'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { login } from '@/stores/authStore'
import { authApi } from '@/api/auth'
import { LoginSchema } from '@/schemas/auth'
import { translateErrorAsync } from '@/lib/translateError'
import { useLanguage } from '@/hooks/useLanguage'

export const Route = createFileRoute('/auth/login')({
  beforeLoad: redirectIfAuthenticated,
  component: LoginPage,
  validateSearch: (search: Record<string, unknown>) => ({
    redirect: (search.redirect as string) || '/',
  }),
})

function LoginPage() {
  const navigate = useNavigate()
  const { redirect } = useSearch({ from: '/auth/login' })
  const { t } = useTranslation()
  const { isRTL } = useLanguage()
  const [showPassword, setShowPassword] = useState(false)
  const [isGoogleLoading, setIsGoogleLoading] = useState(false)

  const form = useForm({
    defaultValues: {
      email: '',
      password: '',
    },
    onSubmit: async ({ value }) => {
      try {
        await login(value.email, value.password)
        await navigate({ to: redirect, search: { redirect } })
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
    validators: {
      onBlur: LoginSchema,
    },
  })

  const handleGoogleLogin = async () => {
    try {
      setIsGoogleLoading(true)

      // Get OAuth URL from backend
      const { url } = await authApi.getGoogleAuthUrl()

      // Append redirect to OAuth state
      const state = encodeURIComponent(JSON.stringify({ redirect }))
      const oauthUrl = url.includes('state=') ? url : `${url}&state=${state}`

      // Redirect to Google OAuth
      window.location.href = oauthUrl
    } catch (error) {
      setIsGoogleLoading(false)
      const message = await translateErrorAsync(error, t)
      toast.error(message)
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
            {t('auth:welcome_message')}
          </p>
        </div>
      </div>

      {/* Right Side - Login Form */}
      <div className="flex-1 flex items-center justify-center p-6 lg:p-12">
        <div className="w-full max-w-md">
          {/* Mobile Logo */}
          <div className="text-center mb-8 lg:hidden">
            <h1 className="text-4xl font-bold text-primary mb-2">Kyora</h1>
            <p className="text-base-content/60">{t('auth:login_subtitle')}</p>
          </div>

          {/* Page Title */}
          <div className="mb-8">
            <h2 className="text-3xl font-bold text-base-content mb-2">
              {t('auth:welcome_back')}
            </h2>
            <p className="text-base-content/60">
              {t('auth:login_description')}
            </p>
          </div>

          {/* Login Form */}
          <form
            onSubmit={(e) => {
              e.preventDefault()
              e.stopPropagation()
              void form.handleSubmit()
            }}
            className="space-y-6"
          >
            {/* Email Field */}
            <form.Field
              name="email"
              validators={{
                onBlur: LoginSchema.shape.email,
              }}
            >
              {(field) => (
                <div className="form-control">
                  <label htmlFor="email" className="label">
                    <span className="label-text font-medium">
                      {t('common:email')}
                    </span>
                  </label>
                  <input
                    id="email"
                    name="email"
                    type="email"
                    autoComplete="email"
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    className={`input input-bordered w-full ${
                      field.state.meta.errors.length > 0 ? 'input-error' : ''
                    }`}
                    placeholder={t('auth:email_placeholder')}
                  />
                  {field.state.meta.errors.length > 0 && (
                    <label className="label">
                      <span className="label-text-alt text-error">
                        {field.state.meta.errors[0]?.message || 'Invalid value'}
                      </span>
                    </label>
                  )}
                </div>
              )}
            </form.Field>

            {/* Password Field */}
            <form.Field
              name="password"
              validators={{
                onBlur: LoginSchema.shape.password,
              }}
            >
              {(field) => (
                <div className="form-control">
                  <label htmlFor="password" className="label">
                    <span className="label-text font-medium">
                      {t('common:password')}
                    </span>
                  </label>
                  <div className="relative">
                    <input
                      id="password"
                      name="password"
                      type={showPassword ? 'text' : 'password'}
                      autoComplete="current-password"
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      className={`input input-bordered w-full ${
                        isRTL ? 'ps-12' : 'pe-12'
                      } ${field.state.meta.errors.length > 0 ? 'input-error' : ''}`}
                      placeholder={t('auth:password_placeholder')}
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className={`absolute ${
                        isRTL ? 'start-3' : 'end-3'
                      } top-1/2 -translate-y-1/2 text-base-content/60 hover:text-base-content transition-colors`}
                    >
                      {showPassword ? (
                        <EyeOff className="w-5 h-5" />
                      ) : (
                        <Eye className="w-5 h-5" />
                      )}
                    </button>
                  </div>
                  {field.state.meta.errors.length > 0 && (
                    <label className="label">
                      <span className="label-text-alt text-error">
                        {field.state.meta.errors[0]?.message || 'Invalid value'}
                      </span>
                    </label>
                  )}
                </div>
              )}
            </form.Field>

            {/* Forgot Password Link */}
            <div className={`flex ${isRTL ? 'justify-start' : 'justify-end'}`}>
              <a
                href="/auth/forgot-password"
                className="text-sm text-primary hover:text-primary-focus transition-colors"
              >
                {t('auth:forgot_password')}
              </a>
            </div>

            {/* Submit Button */}
            <form.Subscribe
              selector={(state) => ({
                canSubmit: state.canSubmit,
                isSubmitting: state.isSubmitting,
              })}
            >
              {({ canSubmit, isSubmitting }) => (
                <button
                  type="submit"
                  disabled={!canSubmit || isSubmitting}
                  className="btn btn-primary w-full"
                >
                  {isSubmitting && <Loader2 className="w-4 h-4 animate-spin" />}
                  {t('auth:login')}
                </button>
              )}
            </form.Subscribe>

            {/* Divider */}
            <div className="divider">{t('common:or')}</div>

            {/* Google Login Button */}
            <button
              type="button"
              onClick={handleGoogleLogin}
              disabled={isGoogleLoading}
              className="btn btn-outline w-full"
            >
              {isGoogleLoading ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <svg className="w-5 h-5" viewBox="0 0 24 24">
                  <path
                    fill="currentColor"
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                  />
                  <path
                    fill="currentColor"
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                  />
                  <path
                    fill="currentColor"
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                  />
                  <path
                    fill="currentColor"
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                  />
                </svg>
              )}
              {t('auth:continue_with_google')}
            </button>
          </form>

          {/* Sign Up Link */}
          <div className="mt-8 text-center">
            <p className="text-base-content/60">
              {t('auth:no_account')}{' '}
              <a
                href="/onboarding"
                className="text-primary hover:text-primary-focus font-semibold hover:underline transition-colors"
              >
                {t('auth:sign_up')}
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
