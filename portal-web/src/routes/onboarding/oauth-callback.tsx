import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { translateErrorAsync } from "@/lib/translateError";

/**
 * Google OAuth Callback Handler for Onboarding
 *
 * Handles the OAuth callback during onboarding flow:
 * 1. Extracts OAuth code from URL
 * 2. Retrieves session token from sessionStorage
 * 3. Calls POST /v1/onboarding/oauth/google
 * 4. Updates onboarding state
 * 5. Redirects to /onboarding/business
 */
export default function OnboardingOAuthCallbackPage() {
  const { t } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { sessionToken, loadSession } = useOnboarding();

  const [error, setError] = useState("");

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get("code");
      const errorParam = searchParams.get("error");

      if (errorParam) {
        setError(t("onboarding:oauth.error", { error: errorParam }));
        setTimeout(() => {
          void navigate("/onboarding/verify", { replace: true });
        }, 3000);
        return;
      }

      if (!code) {
        setError(t("onboarding:oauth.noCode"));
        setTimeout(() => {
          void navigate("/onboarding/verify", { replace: true });
        }, 3000);
        return;
      }

      // Retrieve session token from sessionStorage
      const storedToken =
        sessionToken ??
        sessionStorage.getItem("kyora_onboarding_google_session");

      if (!storedToken) {
        setError(t("onboarding:oauth.noSession"));
        setTimeout(() => {
          void navigate("/onboarding/plan", { replace: true });
        }, 3000);
        return;
      }

      try {
        await onboardingApi.oauthGoogle({
          sessionToken: storedToken,
          code,
        });

        // Clear stored token
        sessionStorage.removeItem("kyora_onboarding_google_session");

        // Reload session from backend to get updated state
        await loadSession(storedToken);

        // Redirect to business setup
        void navigate("/onboarding/business", { replace: true });
      } catch (err) {
        const message = await translateErrorAsync(err, t);
        setError(message);
        setTimeout(() => {
          void navigate("/onboarding/verify", { replace: true });
        }, 3000);
      }
    };

    void handleCallback();
  }, [searchParams, sessionToken, markEmailVerified, navigate, t]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-base-100">
      <div className="card bg-base-100 border border-base-300 shadow-xl max-w-md">
        <div className="card-body">
          <div className="text-center">
            {error ? (
              <>
                <div className="flex justify-center mb-4">
                  <div className="w-16 h-16 bg-error/10 rounded-full flex items-center justify-center">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-8 w-8 text-error"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M6 18L18 6M6 6l12 12"
                      />
                    </svg>
                  </div>
                </div>
                <h2 className="text-xl font-bold text-error mb-3">
                  {t("onboarding:oauth.errorTitle")}
                </h2>
                <p className="text-base-content/70 mb-4">{error}</p>
                <p className="text-sm text-base-content/50">
                  {t("onboarding:oauth.redirecting")}
                </p>
              </>
            ) : (
              <>
                <span className="loading loading-spinner loading-lg text-primary mb-4"></span>
                <h2 className="text-xl font-bold mb-2">
                  {t("onboarding:oauth.processing")}
                </h2>
                <p className="text-base-content/70">
                  {t("onboarding:oauth.pleaseWait")}
                </p>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
