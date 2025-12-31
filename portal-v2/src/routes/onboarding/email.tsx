import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries, useStartSessionMutation } from '@/api/onboarding'
import { authApi } from '@/api/auth'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { OnboardingLayout } from '@/components/templates/OnboardingLayout'

// Search params schema for URL-driven state
const EmailSearchSchema = z.object({
  plan: z.string().min(1),
})

export const Route = createFileRoute('/onboarding/email')({
  validateSearch: (search): z.infer<typeof EmailSearchSchema> => {
    return EmailSearchSchema.parse(search)
  },
  // Prefetch plan data before rendering
  loader: async ({ context, location }) => {
    const { queryClient } = context as RouterContext
    const parsed = EmailSearchSchema.parse(location.search)
    
    // Redirect if no plan selected
    if (!parsed.plan) {
      throw redirect({ to: '/onboarding/plan' })
    }

    // Prefetch plan details for summary card
    const plan = await queryClient.ensureQueryData(
      onboardingQueries.plan(parsed.plan)
    )

    return { plan }
  },
  component: EmailEntryPage,
})

/**
 * Email Entry Step - Step 2 of Onboarding
 *
 * User enters email to start onboarding session
 * - Plan descriptor is in URL (URL-driven state)
 * - TanStack Form for validation
 * - Mutation hook for API call
 * - Prefetched plan data from loader
 */
function EmailEntryPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const navigate = useNavigate()
  const { plan } = Route.useLoaderData()
  const { plan: planDescriptor } = Route.useSearch()

  // Start session mutation
  const startSessionMutation = useStartSessionMutation({
    onSuccess: async (response) => {
      // Navigate to verify with session token in URL
      await navigate({
        to: '/onboarding/verify',
        search: { session: response.sessionToken },
      })
    },
  })

  // TanStack Form with Zod validation
  const form = useKyoraForm({
    defaultValues: {
      email: '',
    },
    onSubmit: async ({ value }) => {
      await startSessionMutation.mutateAsync({
        email: value.email,
        planDescriptor,
      })
    },
  })

  const handleGoogleSignIn = async () => {
    try {
      const { url } = await authApi.getGoogleAuthUrl()
      // Store plan descriptor for OAuth callback
      sessionStorage.setItem('kyora_onboarding_plan', planDescriptor)
      window.location.href = url
    } catch (err) {
      startSessionMutation.reset()
      startSessionMutation.error = err as Error
    }
  }

  return (
    <OnboardingLayout>
      <div className="max-w-md mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-base-content mb-3">
            {tOnboarding('email.title')}
          </h1>
          <p className="text-lg text-base-content/70">
            {tOnboarding('email.subtitle')}
          </p>
        </div>

      {/* Selected Plan Summary */}
      <div className="card bg-base-200 mb-6">
        <div className="card-body">
          <div className="flex justify-between items-center">
            <div>
              <h3 className="font-semibold text-lg">{plan.name}</h3>
              <p className="text-sm text-base-content/70">
                {plan.price === '0'
                  ? tCommon('free')
                  : `${plan.price} ${plan.currency.toUpperCase()}`}
                {plan.price !== '0' && ` / ${plan.billingCycle}`}
              </p>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={async () => {
                await navigate({ to: '/onboarding/plan' })
              }}
            >
              {tCommon('change')}
            </Button>
          </div>
        </div>
      </div>

      {/* Email Form */}
      <div className="card bg-base-100 border border-base-300 shadow-lg">
        <div className="card-body">
          <form.AppForm>
            <form.FormRoot className="space-y-6">
              <form.AppField
                name="email"
                validators={{
                  onBlur: z.string().min(1, 'validation.required').email('validation.invalid_email'),
                }}
              >
                {(field) => (
                  <field.TextField
                    type="email"
                    label={tCommon('email')}
                    placeholder={tOnboarding('email.emailPlaceholder')}
                    autoFocus
                  />
                )}
              </form.AppField>

              {startSessionMutation.error && (
                <div className="alert alert-error">
                  <span className="text-sm">
                    {startSessionMutation.error.message}
                  </span>
                </div>
              )}

              <form.SubmitButton
                variant="primary"
                size="lg"
                fullWidth
                disabled={startSessionMutation.isPending}
              >
                {tOnboarding('email.continue')}
              </form.SubmitButton>
            </form.FormRoot>
          </form.AppForm>

          <div className="divider">{tCommon('or')}</div>

          <form.Subscribe selector={(state) => ({ isSubmitting: state.isSubmitting })}>
            {({ isSubmitting }) => (
              <Button
                type="button"
                variant="outline"
                size="lg"
                fullWidth
                onClick={() => void handleGoogleSignIn()}
                disabled={isSubmitting || startSessionMutation.isPending}
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
                {tOnboarding('email.continueWithGoogle')}
              </Button>
            )}
          </form.Subscribe>
        </div>
      </div>
      </div>
    </OnboardingLayout>
  )
}
