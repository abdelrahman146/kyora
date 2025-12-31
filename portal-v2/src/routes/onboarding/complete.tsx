import { useEffect } from 'react'
import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { CheckCircle2, Loader2 } from 'lucide-react'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import { onboardingQueries, useCompleteOnboardingMutation, useDeleteSessionMutation } from '@/api/onboarding'
import { setTokens } from '@/api/client'
import { setUser } from '@/stores/authStore'

const CompleteSearchSchema = z.object({
  session: z.string().min(1),
})

export const Route = createFileRoute('/onboarding/complete')({
  validateSearch: (search): z.infer<typeof CompleteSearchSchema> => {
    return CompleteSearchSchema.parse(search)
  },

  beforeLoad: async ({ location, context }) => {
    const { queryClient } = context as RouterContext
    const parsed = CompleteSearchSchema.parse(location.search)
    
    const session = await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.session)
    )

    if (session.stage !== 'ready_to_commit' && session.stage !== 'payment_confirmed') {
      throw redirect({
        to: '/onboarding/plan',
        replace: true,
      })
    }
  },

  loader: async ({ location, context }) => {
    const { queryClient } = context as unknown as RouterContext
    const parsed = CompleteSearchSchema.parse(location.search)
    
    await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.session)
    )
  },

  component: CompleteOnboardingPage,
  
  errorComponent: ({ error }) => {
    const { t } = useTranslation('translation')
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="card bg-base-100 border border-base-300 shadow-xl max-w-md">
          <div className="card-body">
            <h2 className="card-title text-error">{t('error.title')}</h2>
            <p className="text-base-content/70">{error.message || t('error.generic')}</p>
          </div>
        </div>
      </div>
    )
  },
})

function CompleteOnboardingPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const navigate = useNavigate()
  const { session: sessionToken } = Route.useSearch()
  
  const completeMutation = useCompleteOnboardingMutation()
  const deleteSessionMutation = useDeleteSessionMutation()

  useEffect(() => {
    const complete = async () => {
      try {
        const result = await completeMutation.mutateAsync({ sessionToken })
        setTokens(result.token, result.refreshToken)
        setUser(result.user)
        deleteSessionMutation.mutate(sessionToken)
        // Success - navigate to home page
        window.location.href = '/'
      } catch (error) {
        console.error('[Complete] Failed to complete onboarding:', error)
      }
    }

    void complete()
  }, [sessionToken, completeMutation, deleteSessionMutation, navigate])

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <div className="card bg-base-100 border border-base-300 shadow-xl max-w-md">
        <div className="card-body">
          <div className="text-center">
            {completeMutation.isPending ? (
              <>
                <Loader2 className="w-16 h-16 text-primary animate-spin mx-auto mb-4" />
                <h2 className="text-2xl font-bold mb-3">
                  {tOnboarding('complete.settingUp')}
                </h2>
                <p className="text-base-content/70">
                  {tOnboarding('complete.pleaseWait')}
                </p>
              </>
            ) : completeMutation.isSuccess ? (
              <>
                <CheckCircle2 className="w-16 h-16 text-success mx-auto mb-4" />
                <h2 className="text-2xl font-bold text-success mb-3">
                  {tOnboarding('complete.success')}
                </h2>
                <p className="text-base-content/70">
                  {tOnboarding('complete.redirecting')}
                </p>
              </>
            ) : completeMutation.isError ? (
              <>
                <div className="alert alert-error mb-4">
                  <span>{completeMutation.error.message}</span>
                </div>
                <button
                  onClick={() => completeMutation.reset()}
                  className="btn btn-primary"
                >
                  {tOnboarding('complete.retry')}
                </button>
              </>
            ) : null}
          </div>
        </div>
      </div>
    </div>
  )
}
