import { useState } from 'react'
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { Eye, EyeOff, Loader2 } from 'lucide-react'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { authApi } from '@/api/auth'
import { ResetPasswordSchema } from '@/schemas/auth'
import { translateErrorAsync } from '@/lib/translateError'
import { useLanguage } from '@/hooks/useLanguage'

export const Route = createFileRoute('/auth/reset-password')({
  beforeLoad: redirectIfAuthenticated,
  component: ResetPasswordPage,
  validateSearch: (search: Record<string, unknown>) => ({
    token: (search.token as string) || '',
  }),
})

function ResetPasswordPage() {
  const navigate = useNavigate()
  const { token } = useSearch({ from: '/auth/reset-password' })
  const { t } = useTranslation()
  const { isRTL } = useLanguage()
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)

  // Early return if no token provided in URL

  if (!token) {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-6">
        <div className="w-full max-w-md text-center">
          <div className="alert alert-error">
            <span>{t('auth:invalid_reset_link')}</span>
          </div>
          <a href="/auth/login" className="btn btn-primary mt-4">
            {t('auth:back_to_login')}
          </a>
        </div>
      </div>
    )
  }

  const form = useForm({
    defaultValues: {
      password: '',
      confirmPassword: '',
    },
    onSubmit: async ({ value }) => {
      try {
        await authApi.resetPassword({
          token,
          password: value.password,
        })

        toast.success(t('auth:password_reset_success'))
        await navigate({ to: '/auth/login', search: { redirect: '/' } })
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
    validators: {
      onBlur: ResetPasswordSchema,
    },
  })

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-6">
      <div className="w-full max-w-md">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary mb-2">Kyora</h1>
          <h2 className="text-2xl font-bold text-base-content mb-2">
            {t('auth:reset_password')}
          </h2>
          <p className="text-base-content/60">
            {t('auth:reset_password_description')}
          </p>
        </div>

        {/* Reset Password Form */}
        <form
          onSubmit={(e) => {
            e.preventDefault()
            e.stopPropagation()
            void form.handleSubmit()
          }}
          className="space-y-6"
        >
          {/* New Password Field */}
          <form.Field
            name="password"
            validators={{
              onBlur: ResetPasswordSchema.shape.password,
            }}
          >
            {(field) => (
              <div className="form-control">
                <label htmlFor="password" className="label">
                  <span className="label-text font-medium">
                    {t('auth:new_password')}
                  </span>
                </label>
                <div className="relative">
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete="new-password"
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

          {/* Confirm Password Field */}
          <form.Field
            name="confirmPassword"
            validators={{
              onBlur: ResetPasswordSchema.shape.confirmPassword,
            }}
          >
            {(field) => (
              <div className="form-control">
                <label htmlFor="confirmPassword" className="label">
                  <span className="label-text font-medium">
                    {t('auth:confirm_password')}
                  </span>
                </label>
                <div className="relative">
                  <input
                    id="confirmPassword"
                    name="confirmPassword"
                    type={showConfirmPassword ? 'text' : 'password'}
                    autoComplete="new-password"
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    className={`input input-bordered w-full ${
                      isRTL ? 'ps-12' : 'pe-12'
                    } ${field.state.meta.errors.length > 0 ? 'input-error' : ''}`}
                    placeholder={t('auth:confirm_password_placeholder')}
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className={`absolute ${
                      isRTL ? 'start-3' : 'end-3'
                    } top-1/2 -translate-y-1/2 text-base-content/60 hover:text-base-content transition-colors`}
                  >
                    {showConfirmPassword ? (
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
                {t('auth:reset_password')}
              </button>
            )}
          </form.Subscribe>

          {/* Back to Login Link */}
          <div className="text-center">
            <a
              href="/auth/login"
              className="text-sm text-base-content/60 hover:text-base-content transition-colors"
            >
              {t('auth:back_to_login')}
            </a>
          </div>
        </form>
      </div>
    </div>
  )
}
