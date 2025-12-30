import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import { AlertCircle, CheckCircle, Loader2 } from 'lucide-react'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { authApi } from '@/api/auth'
import { ResetPasswordSchema } from '@/schemas/auth'
import { translateErrorAsync } from '@/lib/translateError'
import { Button } from '@/components/atoms/Button'
import { PasswordInput } from '@/components/atoms/PasswordInput'

export const Route = createFileRoute('/auth/reset-password')({
  beforeLoad: redirectIfAuthenticated,
  component: ResetPasswordPage,
  validateSearch: (search: Record<string, unknown>) => ({
    token: (search.token as string) || '',
  }),
})

type PageStatus = 'loading' | 'ready' | 'success' | 'error'

function ResetPasswordPage() {
  const navigate = useNavigate()
  const { token } = useSearch({ from: '/auth/reset-password' })
  const { t } = useTranslation()
  const { t: tErrors } = useTranslation('errors')
  const [pageStatus, setPageStatus] = useState<PageStatus>('loading')
  const [errorMessage, setErrorMessage] = useState('')
  const [submitErrorMessage, setSubmitErrorMessage] = useState<string>('')

  useEffect(() => {
    if (!token) {
      queueMicrotask(() => {
        setPageStatus('error')
        setErrorMessage(t('auth.reset_password_missing_token'))
      })
      return
    }

    queueMicrotask(() => {
      setPageStatus('ready')
    })
  }, [token, t])

  const handleBackToLogin = async () => {
    await navigate({ to: '/auth/login', replace: true })
  }

  if (pageStatus === 'loading') {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="text-center">
          <Loader2 className="animate-spin text-primary mx-auto mb-4" size={48} />
          <p className="text-base-content/70">
            {t('auth.reset_password_validating')}
          </p>
        </div>
      </div>
    )
  }

  if (pageStatus === 'error') {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-200 shadow-xl">
            <div className="card-body items-center text-center">
              <div className="w-16 h-16 rounded-full bg-error/20 flex items-center justify-center mb-4">
                <AlertCircle className="text-error" size={32} />
              </div>

              <h1 className="card-title text-2xl mb-2">
                {t('auth.reset_password_error_title')}
              </h1>

              <p className="text-base-content/70 mb-6">{errorMessage}</p>

              <div className="w-full space-y-3">
                <Button
                  type="button"
                  variant="primary"
                  size="lg"
                  fullWidth
                  onClick={() => {
                    void handleBackToLogin()
                  }}
                >
                  {t('auth.return_to_login')}
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="lg"
                  fullWidth
                  onClick={() => {
                    void navigate({ to: '/auth/forgot-password', replace: true })
                  }}
                >
                  {t('auth.request_new_link')}
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (pageStatus === 'success') {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-200 shadow-xl">
            <div className="card-body items-center text-center">
              <div className="w-16 h-16 rounded-full bg-success/20 flex items-center justify-center mb-4">
                <CheckCircle className="text-success" size={32} />
              </div>

              <h1 className="card-title text-2xl mb-2">
                {t('auth.password_reset_success_title')}
              </h1>

              <p className="text-base-content/70 mb-6">
                {t('auth.password_reset_success_description')}
              </p>

              <div className="flex items-center gap-2 text-sm text-base-content/60 mb-4">
                <Loader2 className="animate-spin" size={16} />
                <span>{t('auth.redirecting_to_login')}</span>
              </div>

              <Button
                type="button"
                variant="primary"
                size="lg"
                fullWidth
                onClick={() => {
                  void handleBackToLogin()
                }}
              >
                {t('auth.return_to_login')}
              </Button>
            </div>
          </div>
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
        setSubmitErrorMessage('')
        await authApi.resetPassword({
          token,
          password: value.password,
        })

        setPageStatus('success')
        void setTimeout(() => {
          void navigate({ to: '/auth/login', replace: true })
        }, 2000)
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        setSubmitErrorMessage(message)
      }
    },
    validators: {
      onBlur: ResetPasswordSchema,
    },
  })

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="card bg-base-200 shadow-xl">
          <div className="card-body">
            <h1 className="card-title text-3xl mb-2">
              {t('auth.reset_password_title')}
            </h1>
            <p className="text-base-content/70 mb-6">
              {t('auth.reset_password_description')}
            </p>

            <form
              onSubmit={(e) => {
                e.preventDefault()
                e.stopPropagation()
                void form.handleSubmit()
              }}
              className="space-y-6"
              noValidate
            >
              {submitErrorMessage ? (
                <div role="alert" className="alert alert-error">
                  <span>{submitErrorMessage}</span>
                </div>
              ) : null}

              <form.Field
                name="password"
                validators={{
                  onBlur: ResetPasswordSchema.shape.password,
                }}
              >
                {(field) => (
                  <PasswordInput
                    id="newPassword"
                    label={t('auth.new_password')}
                    placeholder={t('auth.new_password_placeholder')}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                    onBlur={field.handleBlur}
                    error={
                      field.state.meta.errors[0]
                        ? tErrors(field.state.meta.errors[0] as unknown as string)
                        : undefined
                    }
                    helperText={t('auth.password_requirements')}
                    autoComplete="new-password"
                    autoFocus
                    disabled={form.state.isSubmitting}
                    showPasswordToggle
                  />
                )}
              </form.Field>

              <form.Field
                name="confirmPassword"
                validators={{
                  onBlur: ResetPasswordSchema.shape.confirmPassword,
                }}
              >
                {(field) => (
                  <PasswordInput
                    id="confirmPassword"
                    label={t('auth.confirm_password')}
                    placeholder={t('auth.confirm_password_placeholder')}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                    onBlur={field.handleBlur}
                    error={
                      field.state.meta.errors[0]
                        ? tErrors(field.state.meta.errors[0] as unknown as string)
                        : undefined
                    }
                    autoComplete="new-password"
                    disabled={form.state.isSubmitting}
                    showPasswordToggle
                  />
                )}
              </form.Field>

              <Button
                type="submit"
                variant="primary"
                size="lg"
                fullWidth
                loading={form.state.isSubmitting}
                disabled={form.state.isSubmitting}
              >
                {t('auth.reset_password_submit')}
              </Button>

              <div className="text-center">
                <p className="text-sm text-base-content/60">
                  {t('auth.remember_password')}{' '}
                  <button
                    type="button"
                    onClick={() => {
                      void handleBackToLogin()
                    }}
                    className="text-primary hover:text-primary-focus hover:underline transition-colors font-medium"
                  >
                    {t('auth.login')}
                  </button>
                </p>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
