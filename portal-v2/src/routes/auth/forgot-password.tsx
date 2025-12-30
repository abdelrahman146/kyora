import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { ArrowLeft, CheckCircle, Mail } from 'lucide-react'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { authApi } from '@/api/auth'
import { useLanguage } from '@/hooks/useLanguage'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { z } from 'zod'

export const Route = createFileRoute('/auth/forgot-password')({
  beforeLoad: redirectIfAuthenticated,
  component: ForgotPasswordPage,
})

function ForgotPasswordPage() {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const { isRTL } = useLanguage()
  const [isSuccess, setIsSuccess] = useState(false)
  const [submittedEmail, setSubmittedEmail] = useState('')

  const form = useKyoraForm({
    defaultValues: {
      email: '',
    },
    validators: {
      email: {
        onBlur: z.string().email('invalid_email'),
      },
    },
    onSubmit: async ({ value }) => {
      await authApi.forgotPassword(value)
      setSubmittedEmail(value.email)
      setIsSuccess(true)
    },
  })

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

            <form.FormRoot className="space-y-6">
              <form.FormError />

              <form.Field name="email">
                {(field) => (
                  <form.TextField
                    {...field}
                    id="email"
                    type="email"
                    label={t('auth.email')}
                    placeholder={t('auth.email_placeholder')}
                    startIcon={<Mail size={20} />}
                    autoComplete="email"
                    autoFocus
                  />
                )}
              </form.Field>

              <form.SubmitButton variant="primary" size="lg" fullWidth>
                {t('auth.send_reset_link')}
              </form.SubmitButton>

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
            </form.FormRoot>
          </div>
        </div>
      </div>
    </div>
  )
}
