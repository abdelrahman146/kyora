import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { ArrowLeft, Check, Eye, EyeOff, Loader2 } from 'lucide-react'
import { onboardingStore, updateStage } from '@/stores/onboardingStore'
import { useSendOTPMutation, useVerifyEmailMutation } from '@/api/onboarding'
import { OTPVerificationSchema } from '@/schemas/onboarding'
import { translateErrorAsync } from '@/lib/translateError'
import { useLanguage } from '@/hooks/useLanguage'

export const Route = createFileRoute('/onboarding/verify')({
  component: VerifyEmailPage,
})

/**
 * Email Verification Step - Step 3 of Onboarding
 *
 * Features:
 * - 6-digit OTP verification
 * - Profile information (firstName, lastName)
 * - Password setup
 * - Resend OTP with rate limiting
 */
function VerifyEmailPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { isRTL } = useLanguage()
  const state = useStore(onboardingStore)
  const [showPassword, setShowPassword] = useState(false)
  const [resendCooldown, setResendCooldown] = useState(0)

  const verifyMutation = useVerifyEmailMutation()
  const sendOTPMutation = useSendOTPMutation()

  // Redirect if no session
  useEffect(() => {
    if (!state.sessionToken || !state.email) {
      navigate({ to: '/onboarding/email', replace: true })
    }
  }, [state.sessionToken, state.email, navigate])

  // Cooldown timer
  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setInterval(() => {
        setResendCooldown((prev) => Math.max(0, prev - 1))
      }, 1000)
      return () => clearInterval(timer)
    }
  }, [resendCooldown])

  const form = useForm({
    defaultValues: {
      code: '',
      fullName: '',
      password: '',
    },
    onSubmit: async ({ value }) => {
      try {
        const response = await verifyMutation.mutateAsync({
          sessionToken: state.sessionToken!,
          otp: value.code,
          fullName: value.fullName,
          password: value.password,
        })

        updateStage(response.stage)
        toast.success(t('onboarding:verification_success'))
        await navigate({ to: '/onboarding/business' })
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
    validators: {
      onBlur: OTPVerificationSchema,
    },
  })

  const handleResendOTP = async () => {
    try {
      await sendOTPMutation.mutateAsync({
        sessionToken: state.sessionToken!,
      })
      // Set cooldown to 60 seconds (default)
      setResendCooldown(60)
      toast.success(t('onboarding:otp_sent'))
    } catch (error) {
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    }
  }

  return (
    <div className="max-w-lg mx-auto">
      {/* Header */}
      <div className="text-center mb-8">
        <div className="flex justify-center mb-4">
          <div className="w-16 h-16 rounded-full bg-success/10 flex items-center justify-center">
            <Check className="w-8 h-8 text-success" />
          </div>
        </div>
        <h1 className="text-3xl font-bold text-base-content mb-2">
          {t('onboarding:verify_email')}
        </h1>
        <p className="text-base-content/70">
          {t('onboarding:verification_code_sent')} <strong>{state.email}</strong>
        </p>
      </div>

      {/* Form */}
      <form
        onSubmit={(e) => {
          e.preventDefault()
          e.stopPropagation()
          void form.handleSubmit()
        }}
        className="space-y-6"
      >
        {/* OTP Code */}
        <form.Field
          name="code"
          validators={{
            onBlur: OTPVerificationSchema.shape.code,
          }}
        >
          {(field) => (
            <div className="form-control">
              <label htmlFor="code" className="label">
                <span className="label-text font-medium">
                  {t('onboarding:verification_code')}
                </span>
              </label>
              <input
                id="code"
                name="code"
                type="text"
                maxLength={6}
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => {
                  const value = e.target.value.replace(/\D/g, '').slice(0, 6)
                  field.handleChange(value)
                }}
                className={`input input-bordered w-full text-center text-2xl tracking-widest ${
                  field.state.meta.errors.length > 0 ? 'input-error' : ''
                }`}
                placeholder="000000"
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

        {/* Resend OTP */}
        <div className="text-center">
          <button
            type="button"
            onClick={handleResendOTP}
            disabled={resendCooldown > 0 || sendOTPMutation.isPending}
            className="btn btn-ghost btn-sm"
          >
            {sendOTPMutation.isPending ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              t('onboarding:resend_code')
            )}
            {resendCooldown > 0 && ` (${resendCooldown}s)`}
          </button>
        </div>

        <div className="divider">{t('onboarding:profile_information')}</div>

        {/* Full Name */}
        <form.Field
          name="fullName"
          validators={{
            onBlur: OTPVerificationSchema.shape.fullName,
          }}
        >
          {(field) => (
            <div className="form-control">
              <label htmlFor="fullName" className="label">
                <span className="label-text font-medium">
                  {t('common:full_name')}
                </span>
              </label>
              <input
                id="fullName"
                name="fullName"
                type="text"
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
                className={`input input-bordered w-full ${
                  field.state.meta.errors.length > 0 ? 'input-error' : ''
                }`}
                placeholder={t('common:full_name_placeholder')}
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

        {/* Password */}
        <form.Field
          name="password"
          validators={{
            onBlur: OTPVerificationSchema.shape.password,
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

        {/* Submit Button */}
        <form.Subscribe
          selector={(formState) => ({
            canSubmit: formState.canSubmit,
            isSubmitting: formState.isSubmitting,
          })}
        >
          {({ canSubmit, isSubmitting }) => (
            <button
              type="submit"
              disabled={!canSubmit || isSubmitting || verifyMutation.isPending}
              className="btn btn-primary w-full"
            >
              {(isSubmitting || verifyMutation.isPending) && (
                <Loader2 className="w-4 h-4 animate-spin" />
              )}
              {t('onboarding:verify_and_continue')}
            </button>
          )}
        </form.Subscribe>

        {/* Back Button */}
        <button
          type="button"
          onClick={() => navigate({ to: '/onboarding/email' })}
          className="btn btn-ghost w-full"
        >
          <ArrowLeft className="w-4 h-4" />
          {t('common:back')}
        </button>
      </form>
    </div>
  )
}
