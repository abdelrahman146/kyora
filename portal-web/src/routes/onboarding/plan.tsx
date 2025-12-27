import { useState, useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Check } from "lucide-react";
import { OnboardingLayout } from "@/components/templates";
import { Modal } from "@/components/atoms";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { translateErrorAsync } from "@/lib/translateError";
import type { Plan } from "@/api/types/onboarding";

/**
 * Plan Selection Step - Step 1 of Onboarding
 * 
 * Features:
 * - Displays all available billing plans
 * - Shows real plan features and limits from API
 * - Highlights recommended plan (starter)
 * - Mobile-responsive card grid
 * - Proceeds to email entry after selection
 */
export default function PlanSelectionPage() {
  const { t } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { setSelectedPlan } = useOnboarding();

  const [plans, setPlans] = useState<Plan[]>([]);
  const [selectedPlanId, setSelectedPlanId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [viewingPlan, setViewingPlan] = useState<Plan | null>(null);

  // Load plans on mount
  useEffect(() => {
    const loadPlans = async () => {
      try {
        setIsLoading(true);
        const fetchedPlans = await onboardingApi.listPlans();
        setPlans(fetchedPlans);
        
        // Pre-select plan from URL if provided
        const planParam = searchParams.get("plan");
        if (planParam) {
          const plan = fetchedPlans.find(p => p.descriptor === planParam);
          if (plan) {
            setSelectedPlanId(plan.id);
          }
        }
      } catch (err) {
        const message = await translateErrorAsync(err, t);
        setError(message);
      } finally {
        setIsLoading(false);
      }
    };

    void loadPlans();
  }, [searchParams, t]);

  const handleContinue = () => {
    const plan = plans.find(p => p.id === selectedPlanId);
    if (!plan) {
      setError(t("onboarding:plan.selectPlanRequired"));
      return;
    }

    setSelectedPlan(plan);
    void navigate("/onboarding/email");
  };

  // Helper to get enabled features from a plan
  const getEnabledFeatures = (plan: Plan) => {
    const features = [];
    if (plan.features.orderManagement) features.push(t("onboarding:plan.features.orderManagement"));
    if (plan.features.inventoryManagement) features.push(t("onboarding:plan.features.inventoryManagement"));
    if (plan.features.customerManagement) features.push(t("onboarding:plan.features.customerManagement"));
    if (plan.features.expenseManagement) features.push(t("onboarding:plan.features.expenseManagement"));
    if (plan.features.accounting) features.push(t("onboarding:plan.features.accounting"));
    if (plan.features.basicAnalytics) features.push(t("onboarding:plan.features.basicAnalytics"));
    if (plan.features.financialReports) features.push(t("onboarding:plan.features.financialReports"));
    if (plan.features.dataImport) features.push(t("onboarding:plan.features.dataImport"));
    if (plan.features.dataExport) features.push(t("onboarding:plan.features.dataExport"));
    if (plan.features.advancedAnalytics) features.push(t("onboarding:plan.features.advancedAnalytics"));
    if (plan.features.advancedFinancialReports) features.push(t("onboarding:plan.features.advancedFinancialReports"));
    if (plan.features.orderPaymentLinks) features.push(t("onboarding:plan.features.orderPaymentLinks"));
    if (plan.features.invoiceGeneration) features.push(t("onboarding:plan.features.invoiceGeneration"));
    if (plan.features.exportAnalyticsData) features.push(t("onboarding:plan.features.exportAnalyticsData"));
    if (plan.features.aiBusinessAssistant) features.push(t("onboarding:plan.features.aiBusinessAssistant"));
    return features;
  };

  if (isLoading) {
    return (
      <OnboardingLayout currentStep={0} totalSteps={5} showProgress={false}>
        <div className="flex flex-col items-center justify-center gap-4 py-12">
          <span className="loading loading-spinner loading-lg text-primary"></span>
          <p className="text-base-content/60">{t("common:loading")}</p>
        </div>
      </OnboardingLayout>
    );
  }

  if (error && plans.length === 0) {
    return (
      <OnboardingLayout currentStep={0} totalSteps={5} showProgress={false}>
        <div className="max-w-2xl mx-auto text-center">
          <div className="alert alert-error">
            <span>{error}</span>
          </div>
        </div>
      </OnboardingLayout>
    );
  }

  return (
    <OnboardingLayout currentStep={0} totalSteps={5} showProgress={false}>
      <div className="max-w-8xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-base-content mb-3">
            {t("onboarding:plan.title")}
          </h1>
          <p className="text-base md:text-lg text-base-content/70">
            {t("onboarding:plan.subtitle")}
          </p>
        </div>

        {/* Plans Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6 mb-8">
          {plans.map((plan) => {
            const isSelected = selectedPlanId === plan.id;
            const isFree = plan.price === "0" || plan.price === "0.00";
            const isRecommended = plan.descriptor === "starter";
            const enabledFeatures = getEnabledFeatures(plan);

            return (
              <div
                key={plan.id}
                onClick={() => {
                  setSelectedPlanId(plan.id);
                }}
                className={`card bg-base-100 border-2 cursor-pointer transition-all hover:shadow-xl relative ${
                  isSelected
                    ? "border-primary shadow-lg scale-105"
                    : "border-base-300 hover:border-primary/50"
                } ${isRecommended ? "ring-2 ring-secondary ring-offset-2" : ""}`}
              >
                {isRecommended && (
                  <div className="absolute -top-3 start-1/2 -translate-x-1/2 z-10">
                    <span className="badge badge-secondary badge-sm font-semibold">
                      {t("onboarding:plan.recommended")}
                    </span>
                  </div>
                )}

                <div className="card-body p-4 md:p-6">
                  {/* Plan Name */}
                  <h3 className="card-title text-xl md:text-2xl line-clamp-1">{plan.name}</h3>
                  
                  {/* Price */}
                  <div className="my-3">
                    <div className="flex items-baseline gap-1 flex-wrap">
                      <span className="text-3xl md:text-4xl font-bold">
                        {isFree ? t("common:free") : plan.price}
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
                  <ul className="space-y-2 mb-4 min-h-30">
                    <li className="flex items-start gap-2">
                      <Check className="w-4 h-4 md:w-5 md:h-5 text-success mt-0.5 shrink-0" />
                      <span className="text-xs md:text-sm line-clamp-2">
                        {plan.limits.maxTeamMembers === -1
                          ? t("onboarding:plan.unlimitedTeamMembers")
                          : t("onboarding:plan.maxTeamMembers", {
                              count: plan.limits.maxTeamMembers,
                            })}
                      </span>
                    </li>
                    <li className="flex items-start gap-2">
                      <Check className="w-4 h-4 md:w-5 md:h-5 text-success mt-0.5 shrink-0" />
                      <span className="text-xs md:text-sm line-clamp-2">
                        {plan.limits.maxBusinesses === -1
                          ? t("onboarding:plan.unlimitedBusinesses")
                          : t("onboarding:plan.maxBusinesses", {
                              count: plan.limits.maxBusinesses,
                            })}
                      </span>
                    </li>
                    <li className="flex items-start gap-2">
                      <Check className="w-4 h-4 md:w-5 md:h-5 text-success mt-0.5 shrink-0" />
                      <span className="text-xs md:text-sm line-clamp-2">
                        {plan.limits.maxOrdersPerMonth === -1
                          ? t("onboarding:plan.unlimitedOrders")
                          : t("onboarding:plan.maxMonthlyOrders", {
                              count: plan.limits.maxOrdersPerMonth,
                            })}
                      </span>
                    </li>
                  </ul>

                  {/* Top Features (show max 3) */}
                  {enabledFeatures.length > 0 && (
                    <>
                      <div className="divider my-2 text-xs">{t("onboarding:plan.featuresIncluded")}</div>
                      <ul className="space-y-1.5">
                        {enabledFeatures.slice(0, 3).map((feature, idx) => (
                          <li key={idx} className="flex items-start gap-2">
                            <Check className="w-3 h-3 md:w-4 md:h-4 text-primary mt-0.5 shrink-0" />
                            <span className="text-xs line-clamp-1">{feature}</span>
                          </li>
                        ))}
                        {enabledFeatures.length > 3 && (
                          <li className="text-xs text-base-content/60 ps-6">
                            {t("onboarding:plan.andMore", { count: enabledFeatures.length - 3 })}
                          </li>
                        )}
                      </ul>
                    </>
                  )}

                  {/* View Details Button */}
                  <button
                    type="button"
                    onClick={(e) => {
                      e.stopPropagation();
                      setViewingPlan(plan);
                    }}
                    className="btn btn-ghost btn-sm btn-block mt-2 text-xs"
                  >
                    {t("onboarding:plan.viewAllFeatures")}
                  </button>

                  {/* Select Button */}
                  <button
                    type="button"
                    className={`btn btn-block mt-4 ${
                      isSelected ? "btn-primary" : "btn-outline"
                    }`}
                  >
                    {isSelected ? t("common:selected") : t("common:select")}
                  </button>
                </div>
              </div>
            );
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
            className="btn btn-primary btn-block btn-lg"
            disabled={!selectedPlanId}
          >
            {t("onboarding:plan.continue")}
          </button>
        </div>

        {/* Plan Details Modal */}
        <Modal
          isOpen={viewingPlan !== null}
          onClose={() => {
            setViewingPlan(null);
          }}
          title={viewingPlan?.name}
          size="lg"
          footer={
            <>
              <button
                onClick={() => {
                  setViewingPlan(null);
                }}
                className="btn btn-ghost"
              >
                {t("common:close")}
              </button>
              <button
                onClick={() => {
                  if (viewingPlan) {
                    setSelectedPlanId(viewingPlan.id);
                    setViewingPlan(null);
                  }
                }}
                className="btn btn-primary"
              >
                {t("onboarding:plan.selectThisPlan")}
              </button>
            </>
          }
        >
          {viewingPlan && (
            <>
              {/* Pricing */}
              <div className="flex items-baseline gap-2 mb-4">
                <span className="text-3xl font-bold">
                  {viewingPlan.price === "0" || viewingPlan.price === "0.00"
                    ? t("common:free")
                    : viewingPlan.price}
                </span>
                {viewingPlan.price !== "0" && viewingPlan.price !== "0.00" && (
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
              <div className="divider">{t("onboarding:plan.limitsTitle")}</div>

              <ul className="space-y-3 mb-6">
                <li className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-success mt-0.5 shrink-0" />
                  <div>
                    <div className="font-semibold">{t("onboarding:plan.teamMembers")}</div>
                    <div className="text-sm text-base-content/70">
                      {viewingPlan.limits.maxTeamMembers === -1
                        ? t("onboarding:plan.unlimitedTeamMembers")
                        : t("onboarding:plan.maxTeamMembers", {
                            count: viewingPlan.limits.maxTeamMembers,
                          })}
                    </div>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-success mt-0.5 shrink-0" />
                  <div>
                    <div className="font-semibold">{t("onboarding:plan.businesses")}</div>
                    <div className="text-sm text-base-content/70">
                      {viewingPlan.limits.maxBusinesses === -1
                        ? t("onboarding:plan.unlimitedBusinesses")
                        : t("onboarding:plan.maxBusinesses", {
                            count: viewingPlan.limits.maxBusinesses,
                          })}
                    </div>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <Check className="w-5 h-5 text-success mt-0.5 shrink-0" />
                  <div>
                    <div className="font-semibold">{t("onboarding:plan.orders")}</div>
                    <div className="text-sm text-base-content/70">
                      {viewingPlan.limits.maxOrdersPerMonth === -1
                        ? t("onboarding:plan.unlimitedOrders")
                        : t("onboarding:plan.maxMonthlyOrders", {
                            count: viewingPlan.limits.maxOrdersPerMonth,
                          })}
                    </div>
                  </div>
                </li>
              </ul>

              {/* All Features Section */}
              <div className="divider">{t("onboarding:plan.allFeatures")}</div>

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
  );
}
