import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Sparkles, ArrowRight } from "lucide-react";
import { OnboardingLayout } from "@/components/templates";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { setTokens } from "@/api/client";
import { translateErrorAsync } from "@/lib/translateError";

/**
 * Complete Onboarding Step - Final Step
 *
 * Features:
 * - Finalizes onboarding session
 * - Creates workspace, user, business, subscription
 * - Issues JWT tokens and logs user in
 * - Success animation and welcome message
 * - Auto-redirects to dashboard
 *
 * Flow:
 * 1. POST /v1/onboarding/complete
 * 2. Receive user data and tokens
 * 3. Set tokens in auth context
 * 4. Clear onboarding state
 * 5. Navigate to /dashboard
 */
export default function CompletePage() {
  const { t } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const {
    sessionToken,
    stage,
    isPaidPlan,
    isPaymentComplete,
    businessName,
    resetOnboarding,
    loadSessionFromStorage,
  } = useOnboarding();

  const [isCompleting, setIsCompleting] = useState(false);
  const [isComplete, setIsComplete] = useState(false);
  const [error, setError] = useState("");

  // Restore session from localStorage on mount
  useEffect(() => {
    const restoreSession = async () => {
      if (!sessionToken) {
        const hasSession = await loadSessionFromStorage();
        if (!hasSession) {
          await navigate("/onboarding/plan", { replace: true });
        }
      }
    };
    void restoreSession();
  }, [sessionToken, loadSessionFromStorage, navigate]);

  // Redirect if prerequisites not met
  useEffect(() => {
    if (!sessionToken) {
      void navigate("/onboarding/plan", { replace: true });
      return;
    }

    // For paid plans, ensure payment is complete
    if (isPaidPlan && !isPaymentComplete && stage !== "ready_to_commit") {
      void navigate("/onboarding/payment", { replace: true });
      return;
    }

    // For free plans, ensure business is set up
    if (!isPaidPlan && stage !== "business_staged" && stage !== "ready_to_commit") {
      void navigate("/onboarding/business", { replace: true });
    }
  }, [sessionToken, isPaidPlan, isPaymentComplete, stage, navigate]);

  const completeOnboarding = async () => {
    if (!sessionToken) return;

    try {
      setIsCompleting(true);
      setError("");

      const response = await onboardingApi.complete({
        sessionToken,
      });

      // Store tokens directly
      setTokens(response.token, response.refreshToken);

      // Mark as complete
      setIsComplete(true);

      // Clear onboarding state
      resetOnboarding();

      // Redirect to dashboard after animation
      setTimeout(() => {
        void navigate("/dashboard", { replace: true });
      }, 3000);
    } catch (err) {
      const message = await translateErrorAsync(err, t);
      setError(message);
    } finally {
      setIsCompleting(false);
    }
  };

  // Auto-complete on mount if ready
  useEffect(() => {
    if (
      sessionToken &&
      !isComplete &&
      !isCompleting &&
      !error &&
      ((isPaidPlan && isPaymentComplete) || !isPaidPlan)
    ) {
      void completeOnboarding();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (error) {
    return (
      <OnboardingLayout currentStep={5} totalSteps={5}>
        <div className="max-w-lg mx-auto">
          <div className="card bg-base-100 border border-error shadow-xl">
            <div className="card-body">
              <div className="text-center">
                <h2 className="text-2xl font-bold text-error mb-3">
                  {t("onboarding:complete.errorTitle")}
                </h2>
                <p className="text-base-content/70 mb-6">{error}</p>
                <button
                  onClick={() => void completeOnboarding()}
                  className="btn btn-primary"
                  disabled={isCompleting}
                >
                  {isCompleting ? (
                    <>
                      <span className="loading loading-spinner loading-sm"></span>
                      {t("common:loading")}
                    </>
                  ) : (
                    t("common:retry")
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>
      </OnboardingLayout>
    );
  }

  if (isComplete) {
    return (
      <OnboardingLayout currentStep={5} totalSteps={5} showProgress={false}>
        <div className="max-w-2xl mx-auto">
          <div className="card bg-linear-to-br from-primary/10 to-secondary/10 border-2 border-primary shadow-2xl">
            <div className="card-body">
              <div className="text-center">
                {/* Success Animation */}
                <div className="flex justify-center mb-6">
                  <div className="relative">
                    <div className="w-24 h-24 bg-linear-to-br from-primary to-secondary rounded-full flex items-center justify-center animate-bounce">
                      <Sparkles className="w-12 h-12 text-primary-content" />
                    </div>
                    <div className="absolute inset-0 w-24 h-24 bg-linear-to-br from-primary to-secondary rounded-full animate-ping opacity-20"></div>
                  </div>
                </div>

                {/* Welcome Message */}
                <h1 className="text-4xl font-bold mb-3">
                  {t("onboarding:complete.welcomeTitle")}
                </h1>
                <p className="text-xl text-base-content/80 mb-6">
                  {t("onboarding:complete.welcomeMessage", { businessName })}
                </p>

                {/* Features Highlight */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 my-8">
                  <div className="bg-base-100/50 backdrop-blur rounded-lg p-4">
                    <div className="text-2xl mb-2">📦</div>
                    <h3 className="font-semibold mb-1">
                      {t("onboarding:complete.feature1Title")}
                    </h3>
                    <p className="text-sm text-base-content/70">
                      {t("onboarding:complete.feature1Desc")}
                    </p>
                  </div>
                  <div className="bg-base-100/50 backdrop-blur rounded-lg p-4">
                    <div className="text-2xl mb-2">📊</div>
                    <h3 className="font-semibold mb-1">
                      {t("onboarding:complete.feature2Title")}
                    </h3>
                    <p className="text-sm text-base-content/70">
                      {t("onboarding:complete.feature2Desc")}
                    </p>
                  </div>
                  <div className="bg-base-100/50 backdrop-blur rounded-lg p-4">
                    <div className="text-2xl mb-2">🚀</div>
                    <h3 className="font-semibold mb-1">
                      {t("onboarding:complete.feature3Title")}
                    </h3>
                    <p className="text-sm text-base-content/70">
                      {t("onboarding:complete.feature3Desc")}
                    </p>
                  </div>
                </div>

                {/* Redirect Notice */}
                <div className="flex items-center justify-center gap-2 text-base-content/60">
                  <span className="loading loading-spinner loading-sm"></span>
                  <span>{t("onboarding:complete.redirecting")}</span>
                  <ArrowRight className="w-4 h-4 animate-pulse" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </OnboardingLayout>
    );
  }

  return (
    <OnboardingLayout currentStep={5} totalSteps={5}>
      <div className="max-w-lg mx-auto">
        <div className="card bg-base-100 border border-base-300 shadow-xl">
          <div className="card-body">
            <div className="text-center">
              <div className="flex justify-center mb-6">
                <span className="loading loading-spinner loading-lg text-primary"></span>
              </div>
              <h2 className="text-2xl font-bold mb-3">
                {t("onboarding:complete.processingTitle")}
              </h2>
              <p className="text-base-content/70">
                {t("onboarding:complete.processingMessage")}
              </p>
              <div className="mt-6 space-y-2">
                <div className="flex items-center gap-3 justify-center text-sm">
                  <span className="loading loading-ring loading-sm text-success"></span>
                  <span>{t("onboarding:complete.creatingWorkspace")}</span>
                </div>
                <div className="flex items-center gap-3 justify-center text-sm">
                  <span className="loading loading-ring loading-sm text-success"></span>
                  <span>{t("onboarding:complete.settingUpBusiness")}</span>
                </div>
                <div className="flex items-center gap-3 justify-center text-sm">
                  <span className="loading loading-ring loading-sm text-success"></span>
                  <span>{t("onboarding:complete.preparingDashboard")}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </OnboardingLayout>
  );
}
