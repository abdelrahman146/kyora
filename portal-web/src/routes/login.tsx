import { useNavigate, useLocation, Navigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useState } from "react";
import { LoginForm } from "../components/organisms/LoginForm";
import { useAuth } from "../hooks/useAuth";
import { useLanguage } from "../hooks/useLanguage";
import { translateErrorAsync } from "../lib/translateError";
import { authApi } from "../api/auth";
import type { LoginFormData } from "../schemas/auth";

/**
 * Login Page
 *
 * Features:
 * - Email/password login with validation
 * - Google OAuth integration
 * - Redirect to intended destination after login
 * - Inline error messages for auth flows
 * - RTL support
 * - Responsive layout (mobile-first)
 *
 * Flow:
 * 1. User enters credentials
 * 2. Form validates with Zod
 * 3. API call to /v1/auth/login
 * 4. On success: Save tokens, redirect to dashboard or intended destination
 * 5. On error: Show translated error message within the form
 */
export default function LoginPage() {
  const { login, isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();
  const { language, toggleLanguage } = useLanguage();
  const [isGoogleLoading, setIsGoogleLoading] = useState(false);
  const [googleErrorMessage, setGoogleErrorMessage] = useState<string>("");

  // Redirect if already authenticated
  if (!isLoading && isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  // Handle login submission
  const handleLogin = async (data: LoginFormData) => {
    setGoogleErrorMessage("");

    await login(data);

    // Redirect to intended destination or dashboard
    const from =
      (location.state as { from?: { pathname?: string } } | null)?.from
        ?.pathname ?? "/dashboard";
    void navigate(from, { replace: true });
  };

  // Handle Google OAuth
  const handleGoogleLogin = async () => {
    try {
      setIsGoogleLoading(true);
      setGoogleErrorMessage("");

      // Get OAuth URL from backend
      const { url } = await authApi.getGoogleAuthUrl();

      // Prepare state parameter with redirect destination
      const from = (location.state as { from?: { pathname?: string } } | null)?.from?.pathname ?? "/dashboard";
      const state = encodeURIComponent(JSON.stringify({ from }));

      // Append state to OAuth URL if not already present
      const oauthUrl = url.includes("state=") ? url : `${url}&state=${state}`;

      // Redirect to Google OAuth
      window.location.href = oauthUrl;
    } catch (error) {
      setIsGoogleLoading(false);
      const message = await translateErrorAsync(error, t);
      setGoogleErrorMessage(message);
    }
  };

  // Show loading while checking auth status
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-base-100">
        <div className="flex flex-col items-center gap-4">
          <span className="loading loading-spinner loading-lg text-primary"></span>
          <p className="text-base-content/60 text-sm">{t("common.loading")}</p>
        </div>
      </div>
    );
  }

  return (
    <>
      {/* Main Login Layout */}
      <div className="min-h-screen bg-base-100 flex flex-col lg:flex-row">
        {/* Left Side - Branding (Hidden on mobile) */}
        <div className="hidden lg:flex lg:w-1/2 bg-linear-to-br from-primary-600 to-primary-800 items-center justify-center p-12">
          <div className="max-w-md text-center">
            <h1 className="text-5xl font-bold text-white mb-6">Kyora</h1>
            <p className="text-xl text-primary-100 leading-relaxed">
              {t("auth.login_welcome_message")}
            </p>
          </div>
        </div>

        {/* Right Side - Login Form */}
        <div className="flex-1 flex items-center justify-center p-6 lg:p-12">
          <div className="w-full max-w-md">
            {/* Mobile Logo */}
            <div className="text-center mb-8 lg:hidden">
              <h1 className="text-4xl font-bold text-primary-600 mb-2">
                Kyora
              </h1>
              <p className="text-base-content/60">
                {t("auth.login_subtitle")}
              </p>
            </div>

            {/* Page Title */}
            <div className="mb-8">
              <h2 className="text-3xl font-bold text-base-content mb-2">
                {t("auth.welcome_back")}
              </h2>
              <p className="text-base-content/60">
                {t("auth.login_description")}
              </p>
            </div>

            {googleErrorMessage ? (
              <div role="alert" className="alert alert-error mb-6">
                <span>{googleErrorMessage}</span>
              </div>
            ) : null}

            {/* Login Form */}
            <LoginForm 
              onSubmit={handleLogin} 
              onGoogleLogin={() => { void handleGoogleLogin(); }}
              isGoogleLoading={isGoogleLoading}
            />

            {/* Sign Up Link */}
            <div className="mt-8 text-center">
              <p className="text-base-content/60">
                {t("auth.no_account")}{" "}
                <a
                  href="/register"
                  className="text-primary-600 hover:text-primary-700 font-semibold hover:underline transition-colors"
                >
                  {t("auth.sign_up")}
                </a>
              </p>
            </div>

            {/* Language Switcher */}
            <div className="mt-8 text-center">
              <button
                onClick={toggleLanguage}
                className="btn btn-ghost btn-sm"
              >
                {language === "ar" ? "English" : "العربية"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
