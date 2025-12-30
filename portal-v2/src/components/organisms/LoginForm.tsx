import { Link } from '@tanstack/react-router'
import { Loader2, Mail } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '../atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { createLoginValidators } from '@/schemas/auth'
import type { LoginFormData } from '@/schemas/auth'

interface LoginFormProps {
  onSubmit: (data: LoginFormData) => Promise<void>
  onGoogleLogin?: () => void
  isGoogleLoading?: boolean
}

/**
 * Login Form Component
 *
 * Features:
 * - Email and password validation with Zod
 * - TanStack Form integration with useKyoraForm
 * - Progressive validation (submit mode → blur mode after first submission)
 * - Automatic focus management on validation errors
 * - Server error injection for invalid credentials
 * - Loading states during submission
 * - RTL support
 * - Accessible form controls with ARIA attributes
 * - Google OAuth button
 *
 * @example
 * ```tsx
 * <LoginForm
 *   onSubmit={handleLogin}
 *   onGoogleLogin={handleGoogleLogin}
 * />
 * ```
 */
export function LoginForm({
  onSubmit,
  onGoogleLogin,
  isGoogleLoading = false,
}: LoginFormProps) {
  const { t } = useTranslation()

  const form = useKyoraForm({
    defaultValues: {
      email: '',
      password: '',
    },
    validators: createLoginValidators(),
    onSubmit: async ({ value }) => {
      await onSubmit(value)
    },
  })

  return (
    <form.FormRoot className="space-y-6">
      <form.FormError />

      {/* Email Input */}
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
          />
        )}
      </form.Field>

      {/* Password Input */}
      <form.Field name="password">
        {(field) => (
          <form.PasswordField
            {...field}
            id="password"
            label={t('auth.password')}
            placeholder={t('auth.password_placeholder')}
            autoComplete="current-password"
          />
        )}
      </form.Field>

      {/* Forgot Password Link */}
      <div className="text-end">
        <Link
          to="/auth/forgot-password"
          className="text-sm text-primary hover:text-primary-focus hover:underline transition-colors"
        >
          {t('auth.forgot_password')}
        </Link>
      </div>

      {/* Submit Button */}
      <form.SubmitButton variant="primary" size="lg" fullWidth loadingText={t('auth.logging_in')}>
        {t('auth.login')}
      </form.SubmitButton>

      {/* Divider */}
      {onGoogleLogin && (
        <>
          <div className="divider text-neutral-500 text-sm">
            {t('auth.or_continue_with')}
          </div>

          {/* Google Login Button */}
          <form.Subscribe selector={(state) => ({ isSubmitting: state.isSubmitting })}>
            {({ isSubmitting }) => (
              <button
                type="button"
                onClick={onGoogleLogin}
                disabled={isSubmitting || isGoogleLoading}
                className="btn btn-outline btn-lg w-full h-13 rounded-xl border-2 border-neutral-200 hover:border-neutral-300 hover:bg-base-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                {isGoogleLoading ? (
                  <>
                    <Loader2 size={20} className="animate-spin" />
                    <span className="font-semibold">
                      {t('auth.connecting_google')}
                    </span>
                  </>
                ) : (
                  <>
                    <svg
                      viewBox="0 0 24 24"
                      className="w-5 h-5"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                        fill="#4285F4"
                      />
                      <path
                        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                        fill="#34A853"
                      />
                      <path
                        d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                        fill="#FBBC05"
                      />
                      <path
                        d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                        fill="#EA4335"
                      />
                    </svg>
                    <span className="font-semibold">
                      {t('auth.continue_with_google')}
                    </span>
                  </>
                )}
              </button>
            )}
          </form.Subscribe>
        </>
      )}
    </form.FormRoot>
  )
}
