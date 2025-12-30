import { useEffect, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import { Check, Loader2 } from 'lucide-react'
import { z } from 'zod'
import type { Plan } from '@/api/types/onboarding'
import type { RouterContext } from '@/router'
import { onboardingApi, onboardingQueries } from '@/api/onboarding'
import { Modal } from '@/components/atoms/Modal'
import { ResumeSessionDialog } from '@/components/molecules/ResumeSessionDialog'
import {
  clearSession,
  loadSessionFromStorage,
  onboardingStore,
  setPlan,
} from '@/stores/onboardingStore'
import { translateErrorAsync } from '@/lib/translateError'

const PlanSearchSchema = z.object({
  plan: z.string().optional(),
})

export const Route = createFileRoute('/onboarding/plan')({
  validateSearch: (search): z.infer<typeof PlanSearchSchema> => {
    return PlanSearchSchema.parse(search)
  },

  loader: async ({ context }) => {
    const { queryClient } = context as RouterContext
    // Prefetch plans for faster page load
    await queryClient.prefetchQuery(onboardingQueries.plans())
  },

  component: PlanSelectionPage,
})

/**
 * Plan Selection Step - Step 1 of Onboarding
 *
 * Features:
 * - Displays all available billing plans
 * - Shows plan features and limits
 * - Highlights recommended plan
 * - Mobile-responsive card grid
 * - Proceeds to email entry after selection
 */
function PlanSelectionPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const { t: tTranslation } = useTranslation('translation')
  const navigate = useNavigate()
  const state = useStore(onboardingStore)
  const { plan: planParam } = Route.useSearch()

  const [plans, setPlans] = useState<Array<Plan>>([])
  const [selectedPlanId, setSelectedPlanId] = useState<string | null>(
    state.planId,
  )
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [viewingPlan, setViewingPlan] = useState<Plan | null>(null)
  const [showResumeDialog, setShowResumeDialog] = useState(false)
  const [isResuming, setIsResuming] = useState(false)
  const [reloadKey, setReloadKey] = useState(0)

  // Check for existing session on mount
  useEffect(() => {
    const checkExistingSession = async () => {
      const hasSession = await loadSessionFromStorage()
      if (hasSession) {
        setShowResumeDialog(true)
      }
    }

    void checkExistingSession()
  }, [])

  // Load plans on mount
  useEffect(() => {
    const loadPlans = async () => {
      try {
        setIsLoading(true)
        setError('')

        const fetchedPlans = await onboardingApi.listPlans()
        setPlans(fetchedPlans)

        // Pre-select plan from URL if provided
        if (planParam) {
          const matched = fetchedPlans.find((p) => p.descriptor === planParam)
          if (matched) {
            setSelectedPlanId(matched.id)
          }
        }
      } catch (err) {
        const message = await translateErrorAsync(err, tTranslation)
        setError(message)
      } finally {
        setIsLoading(false)
      }
    }

    void loadPlans()
  }, [planParam, reloadKey, tTranslation])

  const handleSelectPlan = (plan: Plan) => {
    setSelectedPlanId(plan.id)
  }

  const handleContinue = async () => {
    const plan = plans.find((p) => p.id === selectedPlanId)
    if (!plan) {
      setError(tOnboarding('plan.selectPlanRequired'))
      return
    }

    const isPaid = parseFloat(plan.price) > 0

    setPlan(plan.id, plan.descriptor, isPaid)
    await navigate({ to: '/onboarding/email' })
  }

  const getEnabledFeatures = (plan: Plan) => {
    const features: Array<string> = []
    if (plan.features.orderManagement)
      features.push(tOnboarding('plan.features.orderManagement'))
    if (plan.features.inventoryManagement)
      features.push(tOnboarding('plan.features.inventoryManagement'))
    if (plan.features.customerManagement)
      features.push(tOnboarding('plan.features.customerManagement'))
    if (plan.features.expenseManagement)
      features.push(tOnboarding('plan.features.expenseManagement'))
    if (plan.features.accounting)
      features.push(tOnboarding('plan.features.accounting'))
    if (plan.features.basicAnalytics)
      features.push(tOnboarding('plan.features.basicAnalytics'))
    if (plan.features.financialReports)
      features.push(tOnboarding('plan.features.financialReports'))
    if (plan.features.dataImport)
      features.push(tOnboarding('plan.features.dataImport'))
    if (plan.features.dataExport)
      features.push(tOnboarding('plan.features.dataExport'))
    if (plan.features.advancedAnalytics)
      features.push(tOnboarding('plan.features.advancedAnalytics'))
    if (plan.features.advancedFinancialReports)
      features.push(tOnboarding('plan.features.advancedFinancialReports'))
    if (plan.features.orderPaymentLinks)
      features.push(tOnboarding('plan.features.orderPaymentLinks'))
    if (plan.features.invoiceGeneration)
      features.push(tOnboarding('plan.features.invoiceGeneration'))
    if (plan.features.exportAnalyticsData)
      features.push(tOnboarding('plan.features.exportAnalyticsData'))
    if (plan.features.aiBusinessAssistant)
      features.push(tOnboarding('plan.features.aiBusinessAssistant'))
    return features
  }

  const navigateToCurrentStage = async () => {
    if (!state.stage) {
      await navigate({ to: '/onboarding/email' })
      return
    }

    switch (state.stage) {
      case 'plan_selected':
      case 'identity_pending':
        await navigate({ to: '/onboarding/email' })
        break
      case 'identity_verified':
        await navigate({ to: '/onboarding/business' })
        break
      case 'business_staged':
      case 'payment_pending':
        await navigate({ to: '/onboarding/payment' })
        break
      case 'payment_confirmed':
      case 'ready_to_commit':
        await navigate({ to: '/onboarding/complete' })
        break
      default:
        await navigate({ to: '/onboarding/email' })
    }
  }

  const handleResumeSession = async () => {
    try {
      setIsResuming(true)
      setShowResumeDialog(false)
      await navigateToCurrentStage()
    } finally {
      setIsResuming(false)
    }
  }

  const handleStartFresh = async () => {
    setShowResumeDialog(false)
    await clearSession()
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-12 h-12 text-primary animate-spin" />
          <p className="text-base-content/60">{tCommon('loading')}</p>
        </div>
      </div>
    )
  }

  if (error && plans.length === 0) {
    return (
      <div className="max-w-2xl mx-auto">
        <div className="alert alert-error">
          <span>{error}</span>
        </div>
        <button
          onClick={() => {
            setReloadKey((k) => k + 1)
          }}
          className="btn btn-primary mt-4"
        >
          {tCommon('retry')}
        </button>
      </div>
    )
  }

  return (
    <div className="max-w-8xl mx-auto px-4">
      {/* Header */}
      <div className="text-center mb-8">
        <h1 className="text-3xl md:text-4xl font-bold text-base-content mb-3">
          {tOnboarding('plan.title')}
        </h1>
        <p className="text-base md:text-lg text-base-content/70">
          {tOnboarding('plan.subtitle')}
        </p>
      </div>

      {/* Plans Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6 mb-8">
        {plans.map((plan) => {
          const isSelected = plan.id === selectedPlanId
          const isFree = plan.price === '0' || plan.price === '0.00'
          const isRecommended = plan.descriptor === 'starter'
          const enabledFeatures = getEnabledFeatures(plan)

          return (
            <div
              key={plan.id}
              onClick={() => {
                handleSelectPlan(plan)
              }}
              className={`card bg-base-100 border-2 cursor-pointer transition-all hover:shadow-xl relative ${
                isSelected
                  ? 'border-primary shadow-lg scale-105'
                  : 'border-base-300 hover:border-primary/50'
              } ${isRecommended ? 'ring-2 ring-secondary ring-offset-2' : ''}`}
            >
              <div className="card-body">
                {/* Recommended Badge */}
                {isRecommended && (
                  <div className="absolute -top-3 start-1/2 -translate-x-1/2 z-10">
                    <span className="badge badge-secondary badge-sm font-semibold">
                      {tOnboarding('plan.recommended')}
                    </span>
                  </div>
                )}

                {/* Plan Name */}
                <h3 className="card-title text-xl md:text-2xl line-clamp-1">
                  {plan.name}
                </h3>

                {/* Price */}
                <div className="my-3">
                  <div className="flex items-baseline gap-1 flex-wrap">
                    <span className="text-3xl md:text-4xl font-bold">
                      {isFree ? tCommon('free') : plan.price}
                    </span>
                    {!isFree && (
                      <>
                        <span className="text-base md:text-lg text-base-content/60">
                          {plan.currency.toUpperCase()}
                        </span>
                        <span className="text-xs md:text-sm text-base-content/60">
                          / {plan.billingCycle}
                        </span>
                      </>
                    )}
                  </div>
                </div>

                {/* Description */}
                {plan.description && (
                  <p className="text-base-content/70 text-xs md:text-sm mb-4 line-clamp-2">
                    {plan.description}
                  </p>
                )}

                <div className="divider my-2"></div>

                {/* Limits */}
                <ul className="space-y-2 mb-4 min-h-[120px]">
                  <li className="flex items-start gap-2">
                    <Check className="w-4 h-4 md:w-5 md:h-5 text-success mt-0.5 shrink-0" />
                    <span className="text-xs md:text-sm line-clamp-2">
                      {plan.limits.maxTeamMembers === -1
                        ? tOnboarding('plan.unlimitedTeamMembers')
                        : tOnboarding('plan.maxTeamMembers', {
                            count: plan.limits.maxTeamMembers,
                          })}
                    </span>
                  </li>
                  <li className="flex items-start gap-2">
                    <Check className="w-4 h-4 md:w-5 md:h-5 text-success mt-0.5 shrink-0" />
                    <span className="text-xs md:text-sm line-clamp-2">
                      {plan.limits.maxBusinesses === -1
                        ? tOnboarding('plan.unlimitedBusinesses')
                        : tOnboarding('plan.maxBusinesses', {
                            count: plan.limits.maxBusinesses,
                          })}
                    </span>
                  </li>
                  <li className="flex items-start gap-2">
                    <Check className="w-4 h-4 md:w-5 md:h-5 text-success mt-0.5 shrink-0" />
                    <span className="text-xs md:text-sm line-clamp-2">
                      {plan.limits.maxOrdersPerMonth === -1
                        ? tOnboarding('plan.unlimitedOrders')
                        : tOnboarding('plan.maxMonthlyOrders', {
                            count: plan.limits.maxOrdersPerMonth,
                          })}
                    </span>
                  </li>
                </ul>

                {/* Top Features (max 3) */}
                {enabledFeatures.length > 0 && (
                  <>
                    <div className="divider my-2 text-xs">
                      {tOnboarding('plan.featuresIncluded')}
                    </div>
                    <ul className="space-y-1.5">
                      {enabledFeatures.slice(0, 3).map((feature, idx) => (
                        <li key={idx} className="flex items-start gap-2">
                          <Check className="w-3 h-3 md:w-4 md:h-4 text-primary mt-0.5 shrink-0" />
                          <span className="text-xs line-clamp-1">{feature}</span>
                        </li>
                      ))}
                      {enabledFeatures.length > 3 && (
                        <li className="text-xs text-base-content/60 ps-6">
                          {tOnboarding('plan.andMore', {
                            count: enabledFeatures.length - 3,
                          })}
                        </li>
                      )}
                    </ul>
                  </>
                )}

                {/* View Details Button */}
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation()
                    setViewingPlan(plan)
                  }}
                  className="btn btn-ghost btn-sm btn-block mt-2 text-xs"
                >
                  {tOnboarding('plan.viewAllFeatures')}
                </button>

                {/* Select Button */}
                <button
                  type="button"
                  className={`btn btn-block mt-4 ${
                    isSelected ? 'btn-primary' : 'btn-outline'
                  }`}
                >
                  {isSelected ? tCommon('selected') : tCommon('select')}
                </button>
              </div>
            </div>
          )
        })}
      </div>

      {/* Continue Button */}
      <div className="max-w-md mx-auto">
        {error && (
          <div className="alert alert-error mb-4">
            <span className="text-sm">{error}</span>
          </div>
        )}

        <button
          onClick={handleContinue}
          disabled={!selectedPlanId}
          className="btn btn-primary btn-block btn-lg"
        >
          {tOnboarding('plan.continue')}
        </button>
      </div>

      {/* Plan Details Modal */}
      <Modal
        isOpen={viewingPlan !== null}
        onClose={() => {
          setViewingPlan(null)
        }}
        title={viewingPlan?.name}
        size="lg"
        footer={
          <>
            <button
              onClick={() => {
                setViewingPlan(null)
              }}
              className="btn btn-ghost"
            >
              {tCommon('close')}
            </button>
            <button
              onClick={() => {
                if (viewingPlan) {
                  setSelectedPlanId(viewingPlan.id)
                  setViewingPlan(null)
                }
              }}
              className="btn btn-primary"
            >
              {tOnboarding('plan.selectThisPlan')}
            </button>
          </>
        }
      >
        {viewingPlan && (
          <>
            <div className="flex items-baseline gap-2 mb-4">
              <span className="text-3xl font-bold">
                {viewingPlan.price === '0' || viewingPlan.price === '0.00'
                  ? tCommon('free')
                  : viewingPlan.price}
              </span>
              {viewingPlan.price !== '0' && viewingPlan.price !== '0.00' && (
                <>
                  <span className="text-lg text-base-content/60">
                    {viewingPlan.currency.toUpperCase()}
                  </span>
                  <span className="text-sm text-base-content/60">
                    / {viewingPlan.billingCycle}
                  </span>
                </>
              )}
            </div>

            {viewingPlan.description && (
              <p className="text-base-content/70 mb-6">
                {viewingPlan.description}
              </p>
            )}

            <div className="divider">{tOnboarding('plan.limitsTitle')}</div>

            <ul className="space-y-3 mb-6">
              <li className="flex items-start gap-3">
                <Check className="w-5 h-5 text-success mt-0.5 shrink-0" />
                <div>
                  <div className="font-semibold">
                    {tOnboarding('plan.teamMembers')}
                  </div>
                  <div className="text-sm text-base-content/70">
                    {viewingPlan.limits.maxTeamMembers === -1
                      ? tOnboarding('plan.unlimitedTeamMembers')
                      : tOnboarding('plan.maxTeamMembers', {
                          count: viewingPlan.limits.maxTeamMembers,
                        })}
                  </div>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <Check className="w-5 h-5 text-success mt-0.5 shrink-0" />
                <div>
                  <div className="font-semibold">
                    {tOnboarding('plan.businesses')}
                  </div>
                  <div className="text-sm text-base-content/70">
                    {viewingPlan.limits.maxBusinesses === -1
                      ? tOnboarding('plan.unlimitedBusinesses')
                      : tOnboarding('plan.maxBusinesses', {
                          count: viewingPlan.limits.maxBusinesses,
                        })}
                  </div>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <Check className="w-5 h-5 text-success mt-0.5 shrink-0" />
                <div>
                  <div className="font-semibold">
                    {tOnboarding('plan.orders')}
                  </div>
                  <div className="text-sm text-base-content/70">
                    {viewingPlan.limits.maxOrdersPerMonth === -1
                      ? tOnboarding('plan.unlimitedOrders')
                      : tOnboarding('plan.maxMonthlyOrders', {
                          count: viewingPlan.limits.maxOrdersPerMonth,
                        })}
                  </div>
                </div>
              </li>
            </ul>

            <div className="divider">{tOnboarding('plan.allFeatures')}</div>

            <ul className="space-y-2">
              {getEnabledFeatures(viewingPlan).map((feature, idx) => (
                <li key={idx} className="flex items-start gap-2">
                  <Check className="w-5 h-5 text-primary mt-0.5 shrink-0" />
                  <span className="text-sm">{feature}</span>
                </li>
              ))}
            </ul>
          </>
        )}
      </Modal>

      {/* Resume Session Dialog */}
      <ResumeSessionDialog
        open={showResumeDialog}
        onResume={handleResumeSession}
        onStartFresh={handleStartFresh}
        email={state.email ?? undefined}
        stage={state.stage ?? undefined}
        isLoading={isResuming}
      />
    </div>
  )
}
