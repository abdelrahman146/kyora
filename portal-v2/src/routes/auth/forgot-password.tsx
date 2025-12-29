import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { Loader2 } from 'lucide-react'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { authApi } from '@/api/auth'
import { ForgotPasswordSchema } from '@/schemas/auth'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/auth/forgot-password')({
  beforeLoad: redirectIfAuthenticated,
  component: ForgotPasswordPage,
})

function ForgotPasswordPage() {
  const navigate = useNavigate()
  const { t } = useTranslation()

  const form = useForm({
    defaultValues: {
      email: '',
    },
    onSubmit: async ({ value }) => {
      try {
        await authApi.forgotPassword(value)
        toast.success(t('auth:reset_link_sent'))
        await navigate({ to: '/auth/login', search: { redirect: '/' } })
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
    validators: {
      onBlur: ForgotPasswordSchema,
    },
  })

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-6">
      <div className="w-full max-w-md">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary mb-2">Kyora</h1>
          <h2 className="text-2xl font-bold text-base-content mb-2">
            {t('auth:forgot_password')}
          </h2>
          <p className="text-base-content/60">
            {t('auth:forgot_password_description')}
          </p>
        </div>

        {/* Forgot Password Form */}
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
              onBlur: ForgotPasswordSchema.shape.email,
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
                {t('auth:send_reset_link')}
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
