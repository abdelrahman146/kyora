import { useEffect } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { ArrowLeft, Loader2, Mail } from 'lucide-react'
import { onboardingStore, startSession } from '@/stores/onboardingStore'
import { useSendOTPMutation, useStartSessionMutation } from '@/api/onboarding'
import { EmailFormSchema } from '@/schemas/onboarding'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/onboarding/email')({
  component: EmailEntryPage,
})

/**
 * Email Entry Step - Step 2 of Onboarding
 *
 * User enters email and receives OTP for verification
 */
function EmailEntryPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const state = useStore(onboardingStore)

  const startSessionMutation = useStartSessionMutation()
  const sendOTPMutation = useSendOTPMutation()

  // Redirect if no plan selected
  useEffect(() => {
    if (!state.planDescriptor) {
      navigate({ to: '/onboarding/plan', replace: true })
    }
  }, [state.planDescriptor, navigate])

  const form = useForm({
    defaultValues: {
      email: state.email || '',
    },
    onSubmit: async ({ value }) => {
      try {
        // Start session with email and plan
        const response = await startSessionMutation.mutateAsync({
          email: value.email,
          planDescriptor: state.planDescriptor!,
        })

        // Store session in onboardingStore
        startSession(
          response.sessionToken,
          response.stage,
          value.email,
          state.planId!,
          state.planDescriptor!,
          response.isPaid,
        )

        // Send OTP
        await sendOTPMutation.mutateAsync({
          sessionToken: response.sessionToken,
        })

        toast.success(t('onboarding:otp_sent'))
        await navigate({ to: '/onboarding/verify' })
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
    validators: {
      onBlur: EmailFormSchema,
    },
  })

  return (
    <div className="max-w-md mx-auto">
      {/* Header */}
      <div className="text-center mb-8">
        <div className="flex justify-center mb-4">
          <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
            <Mail className="w-8 h-8 text-primary" />
          </div>
        </div>
        <h1 className="text-3xl font-bold text-base-content mb-2">
          {t('onboarding:enter_email')}
        </h1>
        <p className="text-base-content/70">{t('onboarding:email_subtitle')}</p>
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
        {/* Email Field */}
        <form.Field
          name="email"
          validators={{
            onBlur: EmailFormSchema.shape.email,
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
          selector={(formState) => ({
            canSubmit: formState.canSubmit,
            isSubmitting: formState.isSubmitting,
          })}
        >
          {({ canSubmit, isSubmitting }) => (
            <button
              type="submit"
              disabled={
                !canSubmit ||
                isSubmitting ||
                startSessionMutation.isPending ||
                sendOTPMutation.isPending
              }
              className="btn btn-primary w-full"
            >
              {(isSubmitting ||
                startSessionMutation.isPending ||
                sendOTPMutation.isPending) && (
                <Loader2 className="w-4 h-4 animate-spin" />
              )}
              {t('onboarding:send_verification_code')}
            </button>
          )}
        </form.Subscribe>

        {/* Back Button */}
        <button
          type="button"
          onClick={() => navigate({ to: '/onboarding/plan' })}
          className="btn btn-ghost w-full"
        >
          <ArrowLeft className="w-4 h-4" />
          {t('common:back')}
        </button>
      </form>
    </div>
  )
}
