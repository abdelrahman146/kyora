import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { OnboardingLayout } from "@/components/templates";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { translateErrorAsync } from "@/lib/translateError";
import { authApi } from "@/api/auth";

/**
 * Email Entry Step - Step 2 of Onboarding
 * 
 * User enters email and chooses verification method:
 * - Email OTP verification
 * - Google OAuth
 */
export default function EmailEntryPage() {
  const { t } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const { selectedPlan, sessionToken, startSession } = useOnboarding();

  const [email, setEmail] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  // If no plan selected, redirect back
  if (!selectedPlan) {
    void navigate("/onboarding/plan", { replace: true });
    return null;
  }

  const handleEmailSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setError(t("onboarding:plan.invalidEmail"));
      return;
    }

    try {
      setIsSubmitting(true);

      const response = await onboardingApi.startSession({
        email,
        planDescriptor: selectedPlan.descriptor,
      });

      startSession(response, email, selectedPlan);
      void navigate("/onboarding/verify");
    } catch (err) {
      const message = await translateErrorAsync(err, t);
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleGoogleSignIn = async () => {
    if (!sessionToken) {
      setError(t("onboarding:oauth.noSession"));
      return;
    }

    try {
      const { url } = await authApi.getGoogleAuthUrl();
      
      // Store session token for OAuth callback
      sessionStorage.setItem("kyora_onboarding_google_session", sessionToken);
      
      // Redirect to Google OAuth
      window.location.href = url;
    } catch (err) {
      const message = await translateErrorAsync(err, t);
      setError(message);
    }
  };

  return (
    <OnboardingLayout currentStep={1} totalSteps={5} showProgress>
      <div className="max-w-md mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-base-content mb-3">
            {t("onboarding:email.title")}
          </h1>
          <p className="text-lg text-base-content/70">
            {t("onboarding:email.subtitle")}
          </p>
        </div>

        {/* Selected Plan Summary */}
        <div className="card bg-base-200 mb-6">
          <div className="card-body">
            <div className="flex justify-between items-center">
              <div>
                <h3 className="font-semibold text-lg">{selectedPlan.name}</h3>
                <p className="text-sm text-base-content/70">
                  {selectedPlan.price === "0" 
                    ? t("common:free") 
                    : `${selectedPlan.price} ${selectedPlan.currency.toUpperCase()}`}
                  {selectedPlan.price !== "0" && ` / ${selectedPlan.billingCycle}`}
                </p>
              </div>
              <button
                type="button"
                onClick={() => void navigate("/onboarding/plan")}
                className="btn btn-ghost btn-sm"
              >
                {t("common:change")}
              </button>
            </div>
          </div>
        </div>

        {/* Email Form */}
        <div className="card bg-base-100 border border-base-300 shadow-lg">
          <div className="card-body">
            <form onSubmit={(e) => void handleEmailSubmit(e)} className="space-y-4">
              <div className="form-control">
                <label htmlFor="email" className="label">
                  <span className="label-text font-medium">
                    {t("common:email")}
                  </span>
                </label>
                <input
                  type="email"
                  id="email"
                  value={email}
                  onChange={(e) => {
                    setEmail(e.target.value);
                  }}
                  placeholder={t("onboarding:email.emailPlaceholder")}
                  className="input input-bordered w-full"
                  required
                  disabled={isSubmitting}
                  autoFocus
                />
              </div>

              {error && (
                <div className="alert alert-error">
                  <span className="text-sm">{error}</span>
                </div>
              )}

              <button
                type="submit"
                className="btn btn-primary btn-block"
                disabled={isSubmitting}
              >
                {isSubmitting ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {t("common:loading")}
                  </>
                ) : (
                  t("onboarding:email.continue")
                )}
              </button>
            </form>

            <div className="divider">{t("common:or")}</div>

            <button
              type="button"
              onClick={() => void handleGoogleSignIn()}
              className="btn btn-outline btn-block"
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
              {t("onboarding:email.continueWithGoogle")}
            </button>
          </div>
        </div>
      </div>
    </OnboardingLayout>
  );
}
