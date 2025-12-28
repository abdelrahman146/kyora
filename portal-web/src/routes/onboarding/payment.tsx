import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { CreditCard, AlertCircle, CheckCircle2 } from "lucide-react";
import { OnboardingLayout } from "@/components/templates";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { translateErrorAsync } from "@/lib/translateError";

/**
 * Payment Step - Step 4 of Onboarding (for paid plans only)
 *
 * Features:
 * - Creates Stripe checkout session
 * - Redirects to Stripe for payment
 * - Handles success/cancel callbacks
 * - Polls payment status after return
 *
 * Flow:
 * 1. POST /v1/onboarding/payment/start
 * 2. Redirect to Stripe Checkout
 * 3. User completes payment
 * 4. Stripe redirects to /onboarding/payment?status=success
 * 5. Navigate to /onboarding/complete
 *
 * Note: For free plans, this step is skipped
 */
export default function PaymentPage() {
  const { t } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const {
    sessionToken,
    stage,
    isPaidPlan,
    selectedPlan,
    checkoutUrl,
    setCheckoutUrl,
    markPaymentComplete,
  } = useOnboarding();

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [paymentStatus, setPaymentStatus] = useState<"pending" | "success" | "cancelled" | null>(null);

  // Check URL params for payment status
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const status = params.get("status");
    
    if (status === "success") {
      setPaymentStatus("success");
      markPaymentComplete();
      
      // Auto-redirect after 2 seconds
      setTimeout(() => {
        void navigate("/onboarding/complete");
      }, 2000);
    } else if (status === "cancelled") {
      setPaymentStatus("cancelled");
    }
  }, [markPaymentComplete, navigate]);

  // Redirect if not at business_staged or payment_pending stage
  useEffect(() => {
    if (!sessionToken) {
      void navigate("/onboarding/plan", { replace: true });
      return;
    }

    if (!isPaidPlan) {
      void navigate("/onboarding/complete", { replace: true });
      return;
    }

    if (stage !== "business_staged" && stage !== "payment_pending" && stage !== "ready_to_commit") {
      void navigate("/onboarding/business", { replace: true });
    }
  }, [sessionToken, isPaidPlan, stage, navigate]);

  const initiatePayment = async () => {
    if (!sessionToken) return;

    try {
      setIsLoading(true);
      setError("");

      const successUrl = `${window.location.origin}/onboarding/payment?status=success`;
      const cancelUrl = `${window.location.origin}/onboarding/payment?status=cancelled`;

      const response = await onboardingApi.startPayment({
        sessionToken,
        successUrl,
        cancelUrl,
      });

      if (response.checkoutUrl) {
        setCheckoutUrl(response.checkoutUrl);
        
        // Redirect to Stripe Checkout
        window.location.href = response.checkoutUrl;
      } else {
        setError(t("onboarding:payment.noCheckoutUrl"));
      }
    } catch (err) {
      const message = await translateErrorAsync(err, t);
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  // Auto-initiate payment if no checkout URL exists
  useEffect(() => {
    if (!checkoutUrl && !paymentStatus && !error && sessionToken && isPaidPlan) {
      void initiatePayment();
    }
  }, []);

  if (paymentStatus === "success") {
    return (
      <OnboardingLayout currentStep={4} totalSteps={5}>
        <div className="max-w-lg mx-auto">
          <div className="card bg-base-100 border border-success shadow-xl">
            <div className="card-body">
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <div className="w-16 h-16 bg-success/10 rounded-full flex items-center justify-center">
                    <CheckCircle2 className="w-8 h-8 text-success" />
                  </div>
                </div>
                <h2 className="text-2xl font-bold text-success mb-3">
                  {t("onboarding:payment.successTitle")}
                </h2>
                <p className="text-base-content/70">
                  {t("onboarding:payment.successMessage")}
                </p>
                <div className="mt-6">
                  <span className="loading loading-spinner loading-sm"></span>
                  <span className="ms-2 text-sm text-base-content/60">
                    {t("onboarding:payment.redirecting")}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </OnboardingLayout>
    );
  }

  if (paymentStatus === "cancelled") {
    return (
      <OnboardingLayout currentStep={4} totalSteps={5}>
        <div className="max-w-lg mx-auto">
          <div className="card bg-base-100 border border-warning shadow-xl">
            <div className="card-body">
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <div className="w-16 h-16 bg-warning/10 rounded-full flex items-center justify-center">
                    <AlertCircle className="w-8 h-8 text-warning" />
                  </div>
                </div>
                <h2 className="text-2xl font-bold text-warning mb-3">
                  {t("onboarding:payment.cancelledTitle")}
                </h2>
                <p className="text-base-content/70 mb-6">
                  {t("onboarding:payment.cancelledMessage")}
                </p>
                <div className="flex gap-3 justify-center">
                  <button
                    onClick={() => {
                      void navigate("/onboarding/plan");
                    }}
                    className="btn btn-ghost"
                  >
                    {t("onboarding:payment.changePlan")}
                  </button>
                  <button
                    onClick={() => {
                      void initiatePayment();
                    }}
                    className="btn btn-primary"
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <>
                        <span className="loading loading-spinner loading-sm"></span>
                        {t("common:loading")}
                      </>
                    ) : (
                      t("onboarding:payment.tryAgain")
                    )}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </OnboardingLayout>
    );
  }

  return (
    <OnboardingLayout currentStep={4} totalSteps={5}>
      <div className="max-w-lg mx-auto">
        <div className="card bg-base-100 border border-base-300 shadow-xl">
          <div className="card-body">
            {/* Header */}
            <div className="text-center mb-6">
              <div className="flex justify-center mb-4">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                  <CreditCard className="w-8 h-8 text-primary" />
                </div>
              </div>
              <h2 className="text-2xl font-bold">
                {t("onboarding:payment.title")}
              </h2>
              <p className="text-base-content/70 mt-2">
                {t("onboarding:payment.subtitle")}
              </p>
            </div>

            {/* Plan Summary */}
            {selectedPlan && (
              <div className="bg-base-200 rounded-lg p-4 mb-6">
                <h3 className="font-semibold text-lg mb-2">{selectedPlan.name}</h3>
                <div className="flex items-baseline gap-1">
                  <span className="text-3xl font-bold text-primary">
                    {selectedPlan.price}
                  </span>
                  <span className="text-base-content/60">
                    {selectedPlan.currency.toUpperCase()}
                  </span>
                  <span className="text-base-content/60">
                    / {selectedPlan.billingCycle}
                  </span>
                </div>
                {selectedPlan.description && (
                  <p className="text-sm text-base-content/70 mt-2">
                    {selectedPlan.description}
                  </p>
                )}
              </div>
            )}

            {error && (
              <div className="alert alert-error mb-4">
                <span className="text-sm">{error}</span>
              </div>
            )}

            {/* Payment Info */}
            <div className="space-y-3 mb-6">
              <div className="flex items-center gap-3 text-sm">
                <CheckCircle2 className="w-5 h-5 text-success shrink-0" />
                <span>{t("onboarding:payment.securePayment")}</span>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <CheckCircle2 className="w-5 h-5 text-success shrink-0" />
                <span>{t("onboarding:payment.cancelAnytime")}</span>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <CheckCircle2 className="w-5 h-5 text-success shrink-0" />
                <span>{t("onboarding:payment.instantAccess")}</span>
              </div>
            </div>

            {/* Actions */}
            <div>
              <button
                onClick={() => {
                  void initiatePayment();
                }}
                className="btn btn-primary w-full"
                disabled={isLoading}
              >
                {isLoading ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {t("common:loading")}
                  </>
                ) : (
                  t("onboarding:payment.proceedToPayment")
                )}
              </button>
            </div>

            <p className="text-xs text-center text-base-content/50 mt-4">
              {t("onboarding:payment.poweredByStripe")}
            </p>
          </div>
        </div>
      </div>
    </OnboardingLayout>
  );
}
