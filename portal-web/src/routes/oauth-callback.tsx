import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Loader2, AlertCircle, CheckCircle } from "lucide-react";
import { authApi } from "../api/auth";
import { translateErrorAsync } from "../lib/translateError";

/**
 * OAuth Callback Handler
 *
 * Handles the OAuth callback from Google after user authorization.
 * Exchanges the authorization code for access tokens.
 *
 * URL: /oauth/callback?code=...&state=...
 *
 * Flow:
 * 1. Extract code and state from URL params
 * 2. Exchange code for tokens via backend API
 * 3. Save tokens and user data
 * 4. Redirect to intended destination (from state) or dashboard
 *
 * Error Handling:
 * - Missing/invalid code
 * - Backend API errors
 * - Network errors
 * - Invalid state parameter
 */
export default function OAuthCallbackPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { t } = useTranslation();

  const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
  const [errorMessage, setErrorMessage] = useState<string>("");

  useEffect(() => {
    const handleCallback = async () => {
      try {
        // Extract OAuth code from URL
        const code = searchParams.get("code");
        const error = searchParams.get("error");
        const errorDescription = searchParams.get("error_description");

        // Handle OAuth errors (user denied access, etc.)
        if (error) {
          setStatus("error");
          setErrorMessage(
            errorDescription ?? t("auth.oauth_error", { error })
          );
          return;
        }

        // Validate code exists
        if (!code) {
          setStatus("error");
          setErrorMessage(t("auth.oauth_missing_code"));
          return;
        }

        // Exchange code for tokens (this automatically saves tokens and user)
        await authApi.loginWithGoogle({ code });

        // Mark as success
        setStatus("success");

        // Get redirect destination from state or default to dashboard
        const state = searchParams.get("state");
        let redirectTo = "/dashboard";

        if (state) {
          try {
            const stateData = JSON.parse(decodeURIComponent(state)) as { from?: string };
            redirectTo = stateData.from ?? "/dashboard";
          } catch {
            // Invalid state, use default
            redirectTo = "/dashboard";
          }
        }

        // Redirect after short delay to show success message
        void setTimeout(() => {
          void navigate(redirectTo, { replace: true });
        }, 1000);
      } catch (error) {
        setStatus("error");
        const message = await translateErrorAsync(error, t);
        setErrorMessage(message);
      }
    };

    void handleCallback().catch((error: unknown) => {
      console.error("OAuth callback error:", error);
    });
  }, [searchParams, navigate, t]);

  // Retry button handler
  const handleRetry = () => {
    void navigate("/login", { replace: true });
  };

  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
      <div className="card w-full max-w-md bg-base-200 shadow-xl">
        <div className="card-body items-center text-center">
          {/* Loading State */}
          {status === "loading" && (
            <>
              <Loader2 className="w-16 h-16 text-primary animate-spin" />
              <h2 className="card-title text-2xl mt-4">
                {t("auth.oauth_processing")}
              </h2>
              <p className="text-base-content/60">
                {t("auth.oauth_processing_description")}
              </p>
            </>
          )}

          {/* Success State */}
          {status === "success" && (
            <>
              <div className="w-16 h-16 rounded-full bg-success/20 flex items-center justify-center">
                <CheckCircle className="w-10 h-10 text-success" />
              </div>
              <h2 className="card-title text-2xl mt-4 text-success">
                {t("auth.oauth_success")}
              </h2>
              <p className="text-base-content/60">
                {t("auth.oauth_success_description")}
              </p>
            </>
          )}

          {/* Error State */}
          {status === "error" && (
            <>
              <div className="w-16 h-16 rounded-full bg-error/20 flex items-center justify-center">
                <AlertCircle className="w-10 h-10 text-error" />
              </div>
              <h2 className="card-title text-2xl mt-4 text-error">
                {t("auth.oauth_failed")}
              </h2>
              <p className="text-base-content/60">{errorMessage}</p>
              <div className="card-actions mt-6">
                <button onClick={handleRetry} className="btn btn-primary">
                  {t("auth.return_to_login")}
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
