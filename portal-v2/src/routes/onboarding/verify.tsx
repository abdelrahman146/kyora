import { useEffect, useRef, useState } from 'react'
import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import { Check, Mail, User } from 'lucide-react'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import {
  onboardingQueries,
  useSendOTPMutation,
  useVerifyEmailMutation,
} from '@/api/onboarding'
import { authApi } from '@/api/auth'
import { FormInput } from '@/components/atoms/FormInput'
import { PasswordInput } from '@/components/atoms/PasswordInput'
import { Button } from '@/components/atoms/Button'
import { ResendCountdownButton } from '@/components'
import { isHTTPError } from '@/lib/errorParser'
import { translateErrorAsync } from '@/lib/translateError'
import { formatCountdownDuration } from '@/lib/utils'
import { OnboardingLayout } from '@/components/templates/OnboardingLayout'

// Search params schema
const VerifySearchSchema = z.object({
  sessionToken: z.string().min(1),
})

// OTP form schema
const OTPFormSchema = z.object({
  code: z.array(z.string().regex(/^\d$/)).length(6),
})

// Profile form schema
const ProfileFormSchema = z
  .object({
    firstName: z.string().trim().min(1, 'validation.required'),
    lastName: z.string().trim().min(1, 'validation.required'),
    password: z.string().min(8, 'validation.password_min_length'),
    confirmPassword: z.string().min(1, 'validation.required'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'validation.passwords_must_match',
    path: ['confirmPassword'],
  })

export const Route = createFileRoute('/onboarding/verify')({
  validateSearch: (search): z.infer<typeof VerifySearchSchema> => {
    return VerifySearchSchema.parse(search)
  },
  loader: async ({ context, location }) => {
    const parsed = VerifySearchSchema.parse(location.search)
    const { queryClient } = context as unknown as RouterContext
    
    // Redirect if no session token
    if (!parsed.sessionToken) {
      throw redirect({ to: '/onboarding/plan' })
    }

    // Prefetch and validate session
    const session = await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.sessionToken)
    )

    // Redirect if wrong stage
    if (session.stage !== 'plan_selected') {
      if (session.stage === 'identity_verified' || session.stage === 'business_staged') {
        throw redirect({ 
          to: '/onboarding/business', 
          search: { sessionToken: parsed.sessionToken } 
        })
      }
    }

    return { session }
  },
  component: VerifyEmailPage,
})

/**
 * Email Verification Step - Step 3 of Onboarding
 *
 * Two-step process:
 * 1. OTP Verification: User enters 6-digit code
 * 2. Profile Setup: User provides name and password
 *
 * Features:
 * - Separate TanStack Forms for each step
 * - Auto-send OTP on mount
 * - Resend OTP with rate limiting
 * - OAuth option
 */
function VerifyEmailPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const { t: tTranslation } = useTranslation('translation')
  const navigate = useNavigate()
  const { session } = Route.useLoaderData()
  const { sessionToken } = Route.useSearch()

  const [step, setStep] = useState<'otp' | 'profile'>('otp')
  const [resendCooldownSeconds, setResendCooldownSeconds] = useState(0)
  const [showLoginCta, setShowLoginCta] = useState(false)
  const [didSendInitialOtp, setDidSendInitialOtp] = useState(false)

  const otpInputRefs = useRef<Array<HTMLInputElement | null>>([])

  // Send OTP mutation
  const sendOTPMutation = useSendOTPMutation({
    onSuccess: (response) => {
      setResendCooldownSeconds(response.retryAfterSeconds ?? 0)
    },
    onError: async (error) => {
      // Extract retry-after from error response
      const retryAfter = await extractRetryAfterSeconds(error)
      if (retryAfter !== null && retryAfter > 0) {
        setResendCooldownSeconds(retryAfter)
      }
    },
  })

  // Verify email mutation
  const verifyEmailMutation = useVerifyEmailMutation({
    onSuccess: async () => {
      await navigate({
        to: '/onboarding/business',
        search: { sessionToken },
      })
    },
    onError: (error) => {
      // Check for user already exists error (409)
      if (
        isHTTPError(error) &&
        error.response.status === 409
      ) {
        setShowLoginCta(true)
      }
    },
  })

  // OTP form
  const otpForm = useForm({
    defaultValues: {
      code: ['', '', '', '', '', ''],
    },
    onSubmit: ({ value }) => {
      const codeString = value.code.join('')
      if (codeString.length === 6) {
        setStep('profile')
      }
    },
    validators: {
      onBlur: OTPFormSchema,
    },
  })

  // Profile form  
  const profileForm = useForm({
    defaultValues: {
      firstName: '',
      lastName: '',
      password: '',
      confirmPassword: '',
    },
    onSubmit: async ({ value }) => {
      setShowLoginCta(false)

      await verifyEmailMutation.mutateAsync({
        sessionToken,
        code: otpForm.state.values.code.join(''),
        firstName: value.firstName,
        lastName: value.lastName,
        password: value.password,
      })
    },
    validators: {
      onBlur: ProfileFormSchema,
    },
  })

  const extractRetryAfterSeconds = async (err: unknown): Promise<number | null> => {
    if (isHTTPError(err)) {
      const resp = err.response
      try {
        const body = (await resp.clone().json()) as unknown
        if (
          body &&
          typeof body === 'object' &&
          'extensions' in body &&
          (body as { extensions?: unknown }).extensions &&
          typeof (body as { extensions: unknown }).extensions === 'object'
        ) {
          const ext = (body as { extensions: Record<string, unknown> }).extensions
          const v = (ext as { retryAfterSeconds?: unknown }).retryAfterSeconds
          if (typeof v === 'number' && Number.isFinite(v)) {
            return Math.max(0, Math.floor(v))
          }
        }
      } catch {
        // ignore
      }
    }
    return null
  }

  const sendOTP = async () => {
    await sendOTPMutation.mutateAsync({ sessionToken })
  }

  // Send initial OTP on mount
  useEffect(() => {
    if (didSendInitialOtp) return
    if (session.stage === 'plan_selected') {
      setDidSendInitialOtp(true)
      void sendOTP()
    }
  }, [didSendInitialOtp, session.stage])

  const handleOtpChange = (index: number, value: string) => {
    if (!/^\d*$/.test(value)) return

    const currentCode = [...otpForm.state.values.code]
    currentCode[index] = value.slice(-1)
    
    otpForm.setFieldValue('code', currentCode)

    if (value && index < 5) {
      otpInputRefs.current[index + 1]?.focus()
    }
  }

  const handleOtpKeyDown = (
    index: number,
    e: React.KeyboardEvent<HTMLInputElement>,
  ) => {
    if (e.key === 'Backspace' && !otpForm.state.values.code[index] && index > 0) {
      otpInputRefs.current[index - 1]?.focus()
    }
  }

  const handleOtpPaste = (e: React.ClipboardEvent) => {
    e.preventDefault()
    const pasted = e.clipboardData.getData('text').trim()
    if (/^\d{6}$/.test(pasted)) {
      otpForm.setFieldValue('code', pasted.split(''))
      otpInputRefs.current[5]?.focus()
    }
  }

  const handleGoogleOAuth = async () => {
    try {
      const { url } = await authApi.getGoogleAuthUrl()
      sessionStorage.setItem('kyora_onboarding_google_session', sessionToken)
      window.location.href = url
    } catch (err) {
      sendOTPMutation.error = err as Error
    }
  }

  const isSubmitting =
    sendOTPMutation.isPending ||
    verifyEmailMutation.isPending ||
    otpForm.state.isSubmitting ||
    profileForm.state.isSubmitting

  return (
    <OnboardingLayout>
      <div className="max-w-lg mx-auto">
        {step === 'otp' ? (
        <div className="card bg-base-100 border border-base-300 shadow-xl">
          <div className="card-body">
            <div className="text-center mb-6">
              <div className="flex justify-center mb-4">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                  <Mail className="w-8 h-8 text-primary" />
                </div>
              </div>
              <h2 className="text-2xl font-bold">{tOnboarding('verify.title')}</h2>
              <p className="text-base-content/70 mt-2">
                {tOnboarding('verify.subtitle', { email: session.email })}
              </p>
            </div>

            {sendOTPMutation.isSuccess && (
              <div className="alert alert-success mb-4">
                <Check className="w-5 h-5" />
                <span>{tOnboarding('verify.otpSent')}</span>
              </div>
            )}

            {sendOTPMutation.error && (
              <div className="alert alert-error mb-4">
                <span>
                  {translateErrorAsync(sendOTPMutation.error, tTranslation)}
                </span>
              </div>
            )}

            <form
              onSubmit={(e) => {
                e.preventDefault()
                e.stopPropagation()
                void otpForm.handleSubmit()
              }}
            >
              <div className="flex justify-center gap-2 mb-6">
                {otpForm.state.values.code.map((digit, index) => (
                  <input
                    key={index}
                    ref={(el) => {
                      otpInputRefs.current[index] = el
                    }}
                    type="text"
                    inputMode="numeric"
                    maxLength={1}
                    value={digit}
                    onChange={(e) => {
                      handleOtpChange(index, e.target.value)
                    }}
                    onKeyDown={(e) => {
                      handleOtpKeyDown(index, e)
                    }}
                    onPaste={index === 0 ? handleOtpPaste : undefined}
                    className="input input-bordered w-12 h-14 text-center text-xl font-bold"
                    disabled={isSubmitting}
                  />
                ))}
              </div>

              <Button
                type="submit"
                variant="primary"
                size="lg"
                fullWidth
                disabled={isSubmitting}
              >
                {tOnboarding('verify.verifyCode')}
              </Button>
            </form>

            <div className="text-center mt-4">
              <ResendCountdownButton
                cooldownSeconds={resendCooldownSeconds}
                isBusy={isSubmitting}
                onResend={sendOTP}
                className="btn btn-ghost btn-sm"
                renderLabel={({ remainingSeconds, canResend }) =>
                  canResend
                    ? tOnboarding('verify.resendCode')
                    : tOnboarding('verify.resendIn', {
                        time: formatCountdownDuration(remainingSeconds),
                      })
                }
              />
            </div>

            <div className="divider">{tCommon('or')}</div>

            <Button
              type="button"
              variant="outline"
              size="lg"
              fullWidth
              onClick={() => void handleGoogleOAuth()}
              disabled={isSubmitting}
            >
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
              {tOnboarding('verify.continueWithGoogle')}
            </Button>
          </div>
        </div>
      ) : (
        <div className="card bg-base-100 border border-base-300 shadow-xl">
          <div className="card-body">
            <div className="text-center mb-6">
              <div className="flex justify-center mb-4">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                  <User className="w-8 h-8 text-primary" />
                </div>
              </div>
              <h2 className="text-2xl font-bold">
                {tOnboarding('verify.completeProfile')}
              </h2>
              <p className="text-base-content/70 mt-2">
                {tOnboarding('verify.subtitle', { email: session.email })}
              </p>
            </div>

            <form
              onSubmit={(e) => {
                e.preventDefault()
                e.stopPropagation()
                void profileForm.handleSubmit()
              }}
              className="space-y-5"
            >
              <profileForm.Field name="firstName">
                {(field) => (
                    <FormInput
                      id={field.name}
                      type="text"
                      label={tCommon('firstName')}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder={tCommon('firstName')}
                      disabled={isSubmitting}
                      error={field.state.meta.errors[0]}
                    />
                )}
              </profileForm.Field>

              <profileForm.Field name="lastName">
                {(field) => (
                    <FormInput
                      id={field.name}
                      type="text"
                      label={tCommon('lastName')}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder={tCommon('lastName')}
                      disabled={isSubmitting}
                      error={field.state.meta.errors[0]}
                    />
                )}
              </profileForm.Field>

              <profileForm.Field name="password">
                {(field) => (
                    <PasswordInput
                      id={field.name}
                      label={tCommon('password')}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder={tCommon('password')}
                      disabled={isSubmitting}
                      error={field.state.meta.errors[0]}
                      helperText={tOnboarding('verify.passwordHint')}
                    />
                )}
              </profileForm.Field>

              <profileForm.Field name="confirmPassword">
                {(field) => (
                    <PasswordInput
                      id={field.name}
                      label={tCommon('confirmPassword')}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder={tCommon('confirmPassword')}
                      disabled={isSubmitting}
                      error={field.state.meta.errors[0]}
                    />
                )}
              </profileForm.Field>

              {verifyEmailMutation.error && (
                <div className="alert alert-error">
                  <div className="flex flex-col gap-2">
                    <span className="text-sm">
                      {translateErrorAsync(
                        verifyEmailMutation.error,
                        tTranslation
                      )}
                    </span>
                    {showLoginCta && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={async () => {
                          await navigate({
                            to: '/auth/login',
                            search: { redirect: '/' },
                          })
                        }}
                      >
                        {tTranslation('auth.login')}
                      </Button>
                    )}
                  </div>
                </div>
              )}

              <Button
                type="submit"
                variant="primary"
                size="lg"
                fullWidth
                disabled={isSubmitting}
                loading={isSubmitting}
              >
                {tCommon('continue')}
              </Button>
            </form>
          </div>
        </div>
      )}
      </div>
    </OnboardingLayout>
  )
}
