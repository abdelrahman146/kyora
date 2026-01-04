import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { useSuspenseQuery } from '@tanstack/react-query'
import { Check } from 'lucide-react'
import { useMemo, useState } from 'react'
import { z } from 'zod'
import type { Plan } from '@/api/types/onboarding'
import type { RouterContext } from '@/router'
import { onboardingQueries } from '@/api/onboarding'
import { OnboardingLayout } from '@/components/templates/OnboardingLayout'
import { Modal } from '@/components/atoms/Modal'

const PlanSearchSchema = z.object({
  plan: z.string().optional(),
})

export const Route = createFileRoute('/onboarding/plan')({
  validateSearch: (search): z.infer<typeof PlanSearchSchema> => {
    return PlanSearchSchema.parse(search)
  },

  loader: async ({ context }) => {
    const { queryClient } = context as RouterContext
    await queryClient.ensureQueryData(onboardingQueries.plans())
  },

  component: PlanSelectionPage,
  
  errorComponent: ({ error }) => {
    const { t } = useTranslation('translation')
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="card bg-base-100 border border-base-300 max-w-md">
          <div className="card-body">
            <h2 className="card-title text-error">{t('error.title')}</h2>
            <p className="text-base-content/70">{error.message || t('error.generic')}</p>
          </div>
        </div>
      </div>
    )
  },
})

/**
 * Plan Selection Step - Step 1 of Onboarding
 *
 * Features:
 * - Displays all available billing plans (sorted by price)
 * - Shows plan features and limits
 * - Highlights recommended plan
 * - Modal to view all plan details
 * - Mobile-responsive card grid
 * - Proceeds to email entry after selection
 * - Plan selection is URL-driven via search params
 */
function PlanSelectionPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const navigate = useNavigate()
  const { plan: planParam } = Route.useSearch()
  const [viewingPlan, setViewingPlan] = useState<Plan | null>(null)

  const { data: plans } = useSuspenseQuery(onboardingQueries.plans())

  // Sort plans by price (cheapest first)
  const sortedPlans = useMemo(() => {
    return [...plans].sort((a, b) => {
      const priceA = parseFloat(a.price)
      const priceB = parseFloat(b.price)
      return priceA - priceB
    })
  }, [plans])

  const handleSelectPlan = (plan: Plan) => {
    // Navigate to email step with plan in URL
    void navigate({ 
      to: '/onboarding/email', 
      search: { plan: plan.descriptor } 
    })
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

  return (
    <OnboardingLayout>
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
          {sortedPlans.map((plan) => {
            const isSelected = plan.descriptor === planParam
            const isFree = plan.price === '0' || plan.price === '0.00'
            const isRecommended = plan.descriptor === 'starter'
            const enabledFeatures = getEnabledFeatures(plan)

            return (
              <div
                key={plan.id}
                onClick={() => {
                  handleSelectPlan(plan)
                }}
                className={`card bg-base-100 border-2 cursor-pointer transition-all  relative ${
                  isSelected
                    ? 'border-primary scale-105'
                    : 'border-base-300 hover:border-primary/50'
                } ${isRecommended ? 'ring-2 ring-secondary ring-offset-2' : ''}`}
              >
                {/* Recommended Badge - Fixed centering with start-1/2 */}
                {isRecommended && (
                  <div className="absolute -top-3 start-1/2 -translate-x-1/2 z-10">
                    <span className="badge badge-secondary badge-sm font-semibold">
                      {tOnboarding('plan.recommended')}
                    </span>
                  </div>
                )}

                <div className="card-body p-4 md:p-6">
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
                type="button"
                onClick={() => {
                  setViewingPlan(null)
                }}
                className="btn btn-ghost"
              >
                {tCommon('close')}
              </button>
              <button
                type="button"
                onClick={() => {
                  if (viewingPlan) {
                    handleSelectPlan(viewingPlan)
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
              {/* Pricing */}
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

              {/* Description */}
              {viewingPlan.description && (
                <p className="text-base-content/70 mb-6">{viewingPlan.description}</p>
              )}

              {/* Limits Section */}
              <div className="divider">{tOnboarding('plan.limitsTitle')}</div>

              <ul className="space-y-3 mb-6">
                <li className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-success mt-0.5 shrink-0" />
                  <div>
                    <div className="font-semibold">{tOnboarding('plan.teamMembers')}</div>
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
                    <div className="font-semibold">{tOnboarding('plan.businesses')}</div>
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
                    <div className="font-semibold">{tOnboarding('plan.orders')}</div>
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

              {/* All Features Section */}
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
      </div>
    </OnboardingLayout>
  )
}
