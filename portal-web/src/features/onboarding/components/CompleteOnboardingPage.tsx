import { useEffect } from 'react'
import { useSearch } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { CheckCircle2, Loader2 } from 'lucide-react'

import {
  useCompleteOnboardingMutation,
  useDeleteSessionMutation,
} from '@/api/onboarding'
import { setTokens } from '@/lib/auth'
import { setUser } from '@/stores/authStore'

export function CompleteOnboardingPage() {
  const { t: tOnboarding } = useTranslation('onboarding')

  const { session: sessionToken } = useSearch({ from: '/onboarding/complete' })

  const completeMutation = useCompleteOnboardingMutation()
  const deleteSessionMutation = useDeleteSessionMutation()

  useEffect(() => {
    const complete = async () => {
      try {
        const result = await completeMutation.mutateAsync({ sessionToken })
        setTokens(result.token, result.refreshToken)
        setUser(result.user)
        deleteSessionMutation.mutate(sessionToken)
        window.location.href = '/'
      } catch (error) {
        console.error('[Complete] Failed to complete onboarding:', error)
      }
    }

    void complete()
  }, [sessionToken, completeMutation, deleteSessionMutation])

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <div className="card bg-base-100 border border-base-300 max-w-md">
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
            ) : null}
          </div>
        </div>
      </div>
    </div>
  )
}
