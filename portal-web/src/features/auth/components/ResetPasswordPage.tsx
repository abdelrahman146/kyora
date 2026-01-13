import { useEffect, useState } from 'react'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { AlertCircle, CheckCircle, Loader2 } from 'lucide-react'
import { z } from 'zod'

import { useResetPasswordMutation } from '@/api/auth'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'

type PageStatus = 'loading' | 'ready' | 'success' | 'error'

export function ResetPasswordPage() {
  const navigate = useNavigate()
  const { token } = useSearch({ from: '/auth/reset-password' })
  const { t: tAuth } = useTranslation('auth')
  const [pageStatus, setPageStatus] = useState<PageStatus>('loading')
  const [errorMessage, setErrorMessage] = useState('')

  const resetPasswordMutation = useResetPasswordMutation()

  useEffect(() => {
    if (!token) {
      queueMicrotask(() => {
        setPageStatus('error')
        setErrorMessage(tAuth('reset_password_missing_token'))
      })
      return
    }

    queueMicrotask(() => {
      setPageStatus('ready')
    })
  }, [token, tAuth])

  const handleBackToLogin = async () => {
    await navigate({
      to: '/auth/login',
      search: { redirect: '/' },
      replace: true,
    })
  }

  if (pageStatus === 'loading') {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="text-center">
          <Loader2
            className="animate-spin text-primary mx-auto mb-4"
            size={48}
          />
          <p className="text-base-content/70">
            {tAuth('reset_password_validating')}
          </p>
        </div>
      </div>
    )
  }

  if (pageStatus === 'error') {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-200">
            <div className="card-body items-center text-center">
              <div className="w-16 h-16 rounded-full bg-error/20 flex items-center justify-center mb-4">
                <AlertCircle className="text-error" size={32} />
              </div>

              <h1 className="card-title text-2xl mb-2">
                {tAuth('reset_password_error_title')}
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
                  {tAuth('return_to_login')}
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="lg"
                  fullWidth
                  onClick={() => {
                    void navigate({
                      to: '/auth/forgot-password',
                      replace: true,
                    })
                  }}
                >
                  {tAuth('request_new_link')}
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
          <div className="card bg-base-200">
            <div className="card-body items-center text-center">
              <div className="w-16 h-16 rounded-full bg-success/20 flex items-center justify-center mb-4">
                <CheckCircle className="text-success" size={32} />
              </div>

              <h1 className="card-title text-2xl mb-2">
                {tAuth('password_reset_success_title')}
              </h1>

              <p className="text-base-content/70 mb-6">
                {tAuth('password_reset_success_description')}
              </p>

              <div className="flex items-center gap-2 text-sm text-base-content/60 mb-4">
                <Loader2 className="animate-spin" size={16} />
                <span>{tAuth('redirecting_to_login')}</span>
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
                {tAuth('return_to_login')}
              </Button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const form = useKyoraForm({
    defaultValues: {
      password: '',
      confirmPassword: '',
    },
    onSubmit: async ({
      value,
    }: {
      value: { password: string; confirmPassword: string }
    }) => {
      await resetPasswordMutation.mutateAsync({
        token,
        newPassword: value.password,
      })

      setPageStatus('success')
      void setTimeout(() => {
        void navigate({
          to: '/auth/login',
          search: { redirect: '/' },
          replace: true,
        })
      }, 2000)
    },
  })

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="card bg-base-200">
          <div className="card-body">
            <h1 className="card-title text-3xl mb-2">
              {tAuth('reset_password_title')}
            </h1>
            <p className="text-base-content/70 mb-6">
              {tAuth('reset_password_description')}
            </p>

            <form.AppForm>
              <form.FormRoot className="space-y-6">
                <form.FormError />

                <form.AppField
                  name="password"
                  validators={{
                    onBlur: z.string().min(8, 'validation.password_min_length'),
                    onChange: ({ value }: { value: string }) => {
                      if (value.length > 0 && value.length < 8) {
                        return 'validation.password_min_length'
                      }
                      return undefined
                    },
                  }}
                >
                  {(field) => (
                    <field.PasswordField
                      label={tAuth('new_password')}
                      placeholder={tAuth('new_password_placeholder')}
                      hint={tAuth('password_requirements')}
                      autoComplete="new-password"
                      autoFocus
                    />
                  )}
                </form.AppField>

                <form.AppField
                  name="confirmPassword"
                  validators={{
                    onBlur: z.string().min(1, 'validation.required'),
                    onChangeListenTo: ['password'],
                    onChange: ({
                      value,
                      fieldApi,
                    }: {
                      value: string
                      fieldApi: any
                    }) => {
                      if (value !== fieldApi.form.getFieldValue('password')) {
                        return 'validation.passwords_must_match'
                      }
                      return undefined
                    },
                  }}
                >
                  {(field) => (
                    <field.PasswordField
                      label={tAuth('confirm_password')}
                      placeholder={tAuth('confirm_password_placeholder')}
                      autoComplete="new-password"
                    />
                  )}
                </form.AppField>

                <form.SubmitButton variant="primary" size="lg" fullWidth>
                  {tAuth('reset_password_submit')}
                </form.SubmitButton>

                <div className="text-center">
                  <p className="text-sm text-base-content/60">
                    {tAuth('remember_password')}{' '}
                    <button
                      type="button"
                      onClick={() => {
                        void handleBackToLogin()
                      }}
                      className="text-primary hover:text-primary-focus hover:underline transition-colors font-medium"
                    >
                      {tAuth('login')}
                    </button>
                  </p>
                </div>
              </form.FormRoot>
            </form.AppForm>
          </div>
        </div>
      </div>
    </div>
  )
}
