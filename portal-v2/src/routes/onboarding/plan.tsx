import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import { Check, Loader2 } from 'lucide-react'
import toast from 'react-hot-toast'
import type { Plan } from '@/api/types/onboarding'
import { usePlansQuery } from '@/api/onboarding'
import { clearSession, onboardingStore, setPlan } from '@/stores/onboardingStore'

export const Route = createFileRoute('/onboarding/plan')({
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
  const { t } = useTranslation()
  const navigate = useNavigate()
  const state = useStore(onboardingStore)
  const [selectedPlanId, setSelectedPlanId] = useState<string | null>(
    state.planId
  )

  const { data: plans, isLoading, error } = usePlansQuery()

  const handleSelectPlan = (plan: Plan) => {
    setSelectedPlanId(plan.id)
  }

  const handleContinue = async () => {
    const plan = plans?.find((p) => p.id === selectedPlanId)
    if (!plan) {
      toast.error(t('onboarding:plan_selection_required'))
      return
    }

    const isPaid = parseFloat(plan.price) > 0

    setPlan(plan.id, plan.descriptor, isPaid)
    await navigate({ to: '/onboarding/email' })
  }

  const getEnabledFeatures = (plan: Plan) => {
    const features: Array<string> = []
    if (plan.features.orderManagement) features.push(t('onboarding:feature_order_management'))
    if (plan.features.inventoryManagement) features.push(t('onboarding:feature_inventory_management'))
    if (plan.features.customerManagement) features.push(t('onboarding:feature_customer_management'))
    if (plan.features.expenseManagement) features.push(t('onboarding:feature_expense_management'))
    if (plan.features.accounting) features.push(t('onboarding:feature_accounting'))
    if (plan.features.basicAnalytics) features.push(t('onboarding:feature_basic_analytics'))
    if (plan.features.financialReports) features.push(t('onboarding:feature_financial_reports'))
    if (plan.features.advancedAnalytics) features.push(t('onboarding:feature_advanced_analytics'))
    if (plan.features.orderPaymentLinks) features.push(t('onboarding:feature_order_payment_links'))
    if (plan.features.invoiceGeneration) features.push(t('onboarding:feature_invoice_generation'))
    if (plan.features.aiBusinessAssistant) features.push(t('onboarding:feature_ai_assistant'))
    return features
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-12 h-12 text-primary animate-spin" />
          <p className="text-base-content/60">{t('common:loading')}</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto">
        <div className="alert alert-error">
          <span>{error.message || t('errors:something_went_wrong')}</span>
        </div>
        <button
          onClick={() => clearSession()}
          className="btn btn-primary mt-4"
        >
          {t('common:try_again')}
        </button>
      </div>
    )
  }

  return (
    <div className="max-w-6xl mx-auto">
      {/* Header */}
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold text-base-content mb-4">
          {t('onboarding:choose_your_plan')}
        </h1>
        <p className="text-xl text-base-content/70">
          {t('onboarding:plan_subtitle')}
        </p>
      </div>

      {/* Plans Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        {plans?.map((plan) => {
          const isSelected = plan.id === selectedPlanId
          const isPaid = parseFloat(plan.price) > 0
          const features = getEnabledFeatures(plan)
          const isRecommended = plan.descriptor === 'starter'

          return (
            <div
              key={plan.id}
              className={`card bg-base-100 shadow-xl border-2 transition-all cursor-pointer hover:shadow-2xl ${
                isSelected
                  ? 'border-primary'
                  : 'border-base-300 hover:border-primary/50'
              } ${isRecommended ? 'ring-2 ring-secondary ring-offset-2' : ''}`}
              onClick={() => handleSelectPlan(plan)}
            >
              <div className="card-body">
                {/* Recommended Badge */}
                {isRecommended && (
                  <div className="badge badge-secondary absolute -top-3 left-1/2 -translate-x-1/2">
                    {t('onboarding:recommended')}
                  </div>
                )}

                {/* Plan Name */}
                <h2 className="card-title text-2xl">{plan.name}</h2>

                {/* Price */}
                <div className="my-4">
                  <span className="text-4xl font-bold text-primary">
                    {plan.currency === 'USD' ? '$' : plan.currency}
                    {plan.price}
                  </span>
                  <span className="text-base-content/60">
                    /{plan.billingCycle}
                  </span>
                </div>

                {/* Description */}
                {plan.description && (
                  <p className="text-base-content/70 mb-4">{plan.description}</p>
                )}

                {/* Features */}
                <div className="space-y-2">
                  {features.slice(0, 6).map((feature, index) => (
                    <div key={index} className="flex items-start gap-2">
                      <Check className="w-5 h-5 text-success flex-shrink-0 mt-0.5" />
                      <span className="text-sm">{feature}</span>
                    </div>
                  ))}
                  {features.length > 6 && (
                    <p className="text-sm text-base-content/60 ps-7">
                      +{features.length - 6} {t('onboarding:more_features')}
                    </p>
                  )}
                </div>

                {/* Limits */}
                <div className="divider"></div>
                <div className="text-sm space-y-1 text-base-content/70">
                  <p>
                    • {plan.limits.maxOrdersPerMonth === -1 ? t('onboarding:unlimited') : plan.limits.maxOrdersPerMonth} {t('onboarding:orders_per_month')}
                  </p>
                  <p>
                    • {plan.limits.maxTeamMembers === -1 ? t('onboarding:unlimited') : plan.limits.maxTeamMembers} {t('onboarding:team_members')}
                  </p>
                  <p>
                    • {plan.limits.maxBusinesses === -1 ? t('onboarding:unlimited') : plan.limits.maxBusinesses} {t('onboarding:businesses')}
                  </p>
                </div>

                {/* Select Button */}
                <div className="card-actions mt-4">
                  <button
                    className={`btn w-full ${
                      isSelected ? 'btn-primary' : 'btn-outline btn-primary'
                    }`}
                  >
                    {isSelected ? (
                      <>
                        <Check className="w-5 h-5" />
                        {t('onboarding:selected')}
                      </>
                    ) : (
                      t('onboarding:select_plan')
                    )}
                  </button>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Continue Button */}
      <div className="flex justify-center">
        <button
          onClick={handleContinue}
          disabled={!selectedPlanId}
          className="btn btn-primary btn-lg"
        >
          {t('common:continue')}
        </button>
      </div>
    </div>
  )
}
