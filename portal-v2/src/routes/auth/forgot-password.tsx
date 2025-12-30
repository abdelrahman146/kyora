import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import { ArrowLeft, CheckCircle, Mail } from 'lucide-react'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { authApi } from '@/api/auth'
import { ForgotPasswordSchema } from '@/schemas/auth'
import { translateErrorAsync } from '@/lib/translateError'
import { useLanguage } from '@/hooks/useLanguage'
import { Button } from '@/components/atoms/Button'
import { FormInput } from '@/components/atoms/FormInput'
import { getErrorText } from '@/lib/formErrors'

export const Route = createFileRoute('/auth/forgot-password')({
  beforeLoad: redirectIfAuthenticated,
  component: ForgotPasswordPage,
})

function ForgotPasswordPage() {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const { t: tErrors } = useTranslation('errors')
  const { isRTL } = useLanguage()
  const [isSuccess, setIsSuccess] = useState(false)
  const [submittedEmail, setSubmittedEmail] = useState('')
  const [submitErrorMessage, setSubmitErrorMessage] =
    useState<string>('')

  const form = useForm({
    defaultValues: {
      email: '',
    },
    onSubmit: async ({ value }) => {
      try {
        setSubmitErrorMessage('')
        await authApi.forgotPassword(value)
        setSubmittedEmail(value.email)
        setIsSuccess(true)
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        setSubmitErrorMessage(message)
      }
    },
    validators: {
      onBlur: ForgotPasswordSchema,
    },
  })

  // Extract form.state.isSubmitting to minimize subscriptions
  // Accessing form.state multiple times in JSX creates multiple subscriptions
  const isSubmitting = form.state.isSubmitting

  const handleBackToLogin = async () => {
    await navigate({
      to: '/auth/login',
      search: { redirect: '/' },
      replace: true,
    })
  }

  if (isSuccess) {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-200 shadow-xl">
            <div className="card-body items-center text-center">
              <div className="w-16 h-16 rounded-full bg-success/20 flex items-center justify-center mb-4">
                <CheckCircle className="text-success" size={32} />
              </div>

              <h1 className="card-title text-2xl mb-2">
                {t('auth.password_reset_sent_title')}
              </h1>

              <p className="text-base-content/70 mb-6">
                {t('auth.password_reset_sent_description', {
                  email: submittedEmail,
                })}
              </p>

              <div className="alert alert-info mb-6">
                <div className="flex flex-col gap-2 text-start">
                  <p className="text-sm font-medium">
                    {t('auth.password_reset_email_tips_title')}
                  </p>
                  <ul className="text-xs list-disc list-inside space-y-1">
                    <li>{t('auth.password_reset_tip_check_spam')}</li>
                    <li>{t('auth.password_reset_tip_expires')}</li>
                    <li>{t('auth.password_reset_tip_no_account')}</li>
                  </ul>
                </div>
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

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="card bg-base-200 shadow-xl">
          <div className="card-body">
            <button
              type="button"
              onClick={() => {
                void handleBackToLogin()
              }}
              className="btn btn-ghost btn-sm w-fit mb-4 -ms-2"
              aria-label={t('auth.back_to_login')}
            >
              <ArrowLeft size={20} className={isRTL ? 'rotate-180' : ''} />
              {t('auth.back_to_login')}
            </button>

            <h1 className="card-title text-3xl mb-2">
              {t('auth.forgot_password_title')}
            </h1>
            <p className="text-base-content/70 mb-6">
              {t('auth.forgot_password_description')}
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
                name="email"
                validators={{
                  onBlur: ForgotPasswordSchema.shape.email,
                }}
              >
                {(field) => {
                  const errorKey = getErrorText(field.state.meta.errors)

                  return (
                    <FormInput
                      id="email"
                      type="email"
                      label={t('auth.email')}
                      placeholder={t('auth.email_placeholder')}
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      onBlur={field.handleBlur}
                      error={errorKey ? tErrors(errorKey) : undefined}
                      startIcon={<Mail size={20} />}
                      autoComplete="email"
                      autoFocus
                      disabled={isSubmitting}
                    />
                  )
                }}
              </form.Field>

              <Button
                type="submit"
                variant="primary"
                size="lg"
                fullWidth
                loading={isSubmitting}
                disabled={isSubmitting}
              >
                {t('auth.send_reset_link')}
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
