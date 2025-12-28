import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Mail, Check, Lock, User } from "lucide-react";
import { OnboardingLayout } from "@/components/templates";
import { ResendCountdownButton, FormInput } from "@/components";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { authApi } from "@/api/auth";
import { translateErrorAsync } from "@/lib/translateError";
import { isHTTPError } from "@/lib/errorParser";
import { formatCountdownDuration } from "@/lib/utils";

/**
 * Email Verification Step - Step 2 of Onboarding
 *
 * Features:
 * - Email OTP verification flow
 * - Google OAuth alternative
 * - Resend OTP with rate limiting
 * - Auto-focus OTP input fields
 * - Profile information (firstName, lastName)
 * - Password setup
 *
 * Flow:
 * 1. Send OTP to email
 * 2. User enters 6-digit code
 * 3. User provides first name, last name, password
 * 4. POST /v1/onboarding/email/verify
 * 5. Navigate to /onboarding/business
 *
 * Alternative Flow (Google OAuth):
 * 1. User clicks "Continue with Google"
 * 2. Redirect to Google OAuth
 * 3. Callback with code
 * 4. POST /v1/onboarding/oauth/google
 * 5. Navigate to /onboarding/business
 */
export default function VerifyEmailPage() {
  const { t } = useTranslation(["onboarding", "common", "translation"]);
  const navigate = useNavigate();
  const {
    sessionToken,
    email,
    stage,
    loadSessionFromStorage,
    loadSession,
  } = useOnboarding();

  const [step, setStep] = useState<"otp" | "profile">("otp");
  const [otpCode, setOtpCode] = useState(["", "", "", "", "", ""]);
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [resendCooldownSeconds, setResendCooldownSeconds] = useState(0);
  const [showLoginCta, setShowLoginCta] = useState(false);

  const otpInputRefs = useRef<(HTMLInputElement | null)[]>([]);

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

  // Redirect if no session token
  useEffect(() => {
    if (!sessionToken || !email) {
      void navigate("/onboarding/plan", { replace: true });
    }
  }, [sessionToken, email, navigate]);

  // Send initial OTP
  useEffect(() => {
    if (sessionToken && stage === "plan_selected") {
      void sendOTP();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const extractRetryAfterSeconds = async (err: unknown): Promise<number | null> => {
    if (
      typeof err === "object" &&
      err !== null &&
      "response" in err &&
      (err as { response?: unknown }).response instanceof Response
    ) {
      const resp = (err as { response: Response }).response;
      try {
        const body = (await resp.clone().json()) as unknown;
        if (
          body &&
          typeof body === "object" &&
          "extensions" in body &&
          (body as { extensions?: unknown }).extensions &&
          typeof (body as { extensions: unknown }).extensions === "object"
        ) {
          const ext = (body as { extensions: Record<string, unknown> }).extensions;
          const v = (ext as { retryAfterSeconds?: unknown }).retryAfterSeconds;
          if (typeof v === "number" && Number.isFinite(v)) {
            return Math.max(0, Math.floor(v));
          }
        }
      } catch {
        // ignore
      }
    }
    return null;
  };

  const sendOTP = async () => {
    if (!sessionToken) return;

    try {
      setError("");
      setIsSubmitting(true);
      const { retryAfterSeconds } = await onboardingApi.sendEmailOTP({
        sessionToken,
      });
      setSuccess(t("onboarding:verify.otpSent"));
      setResendCooldownSeconds(retryAfterSeconds);
    } catch (err) {
      const retryAfterSeconds = await extractRetryAfterSeconds(err);
      if (retryAfterSeconds !== null && retryAfterSeconds > 0) {
        setResendCooldownSeconds(retryAfterSeconds);
      }
      const message = await translateErrorAsync(err, t);
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleOtpChange = (index: number, value: string) => {
    if (!/^\d*$/.test(value)) return;

    const newOtp = [...otpCode];
    newOtp[index] = value.slice(-1);
    setOtpCode(newOtp);

    // Auto-focus next input
    if (value && index < 5) {
      otpInputRefs.current[index + 1]?.focus();
    }
  };

  const handleOtpKeyDown = (
    index: number,
    e: React.KeyboardEvent<HTMLInputElement>
  ) => {
    if (e.key === "Backspace" && !otpCode[index] && index > 0) {
      otpInputRefs.current[index - 1]?.focus();
    }
  };

  const handleOtpPaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData("text").trim();
    if (/^\d{6}$/.test(pastedData)) {
      setOtpCode(pastedData.split(""));
      otpInputRefs.current[5]?.focus();
    }
  };

  const handleVerifyOtp = () => {
    const code = otpCode.join("");
    if (code.length !== 6) {
      setError(t("onboarding:verify.invalidCode"));
      return;
    }

    setStep("profile");
    setError("");
    setSuccess("");
  };

  const submitProfile = async () => {
    setError("");
    setShowLoginCta(false);

    if (password !== confirmPassword) {
      setError(t("onboarding:verify.passwordMismatch"));
      return;
    }

    if (password.length < 8) {
      setError(t("onboarding:verify.passwordTooShort"));
      return;
    }

    if (!sessionToken) return;

    try {
      setIsSubmitting(true);

      await onboardingApi.verifyEmail({
        sessionToken,
        code: otpCode.join(""),
        firstName,
        lastName,
        password,
      });

      // Reload session from backend to get updated state
      await loadSession(sessionToken);
      
      // Navigate to business setup
      void navigate("/onboarding/business");
    } catch (err) {
      if (
        isHTTPError(err) &&
        err.response.status === 409 &&
        typeof err.response.url === "string" &&
        err.response.url.endsWith("/v1/onboarding/email/verify")
      ) {
        setShowLoginCta(true);
      }
      const message = await translateErrorAsync(err, t);
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSubmitProfile: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    void submitProfile();
  };

  const handleGoogleOAuth = async () => {
    try {
      setIsSubmitting(true);
      const { url } = await authApi.getGoogleAuthUrl();
      
      // Store session token in sessionStorage for callback
      sessionStorage.setItem("kyora_onboarding_google_session", sessionToken ?? "");
      
      // Redirect to Google OAuth
      window.location.href = url;
    } catch (err) {
      const message = await translateErrorAsync(err, t);
      setError(message);
      setIsSubmitting(false);
    }
  };

  return (
    <OnboardingLayout currentStep={2} totalSteps={5}>
      <div className="max-w-lg mx-auto">
        {step === "otp" ? (
          <div className="card bg-base-100 border border-base-300 shadow-xl">
            <div className="card-body">
              {/* Header */}
              <div className="text-center mb-6">
                <div className="flex justify-center mb-4">
                  <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                    <Mail className="w-8 h-8 text-primary" />
                  </div>
                </div>
                <h2 className="text-2xl font-bold">
                  {t("onboarding:verify.title")}
                </h2>
                <p className="text-base-content/70 mt-2">
                  {t("onboarding:verify.subtitle", { email })}
                </p>
              </div>

              {/* Success Message */}
              {success && (
                <div className="alert alert-success mb-4">
                  <Check className="w-5 h-5" />
                  <span>{success}</span>
                </div>
              )}

              {/* Error Message */}
              {error && (
                <div className="alert alert-error mb-4">
                  <span>{error}</span>
                </div>
              )}

              {/* OTP Input */}
              <div className="flex justify-center gap-2 mb-6">
                {otpCode.map((digit, index) => (
                  <input
                    key={index}
                    ref={(el) => {
                      otpInputRefs.current[index] = el;
                    }}
                    type="text"
                    inputMode="numeric"
                    maxLength={1}
                    value={digit}
                    onChange={(e) => {
                      handleOtpChange(index, e.target.value);
                    }}
                    onKeyDown={(e) => {
                      handleOtpKeyDown(index, e);
                    }}
                    onPaste={index === 0 ? handleOtpPaste : undefined}
                    className="input input-bordered w-12 h-14 text-center text-xl font-bold"
                    disabled={isSubmitting}
                  />
                ))}
              </div>

              {/* Verify Button */}
              <button
                type="button"
                onClick={handleVerifyOtp}
                disabled={isSubmitting}
                className="btn btn-primary btn-block"
              >
                {t("onboarding:verify.verifyCode")}
              </button>

              {/* Resend Button */}
              <div className="text-center mt-4">
                <ResendCountdownButton
                  cooldownSeconds={resendCooldownSeconds}
                  isBusy={isSubmitting}
                  onResend={sendOTP}
                  className="btn btn-ghost btn-sm"
                  renderLabel={({ remainingSeconds, canResend }) =>
                    canResend
                      ? t("onboarding:verify.resendCode")
                      : t("onboarding:verify.resendIn", {
                          time: formatCountdownDuration(remainingSeconds),
                        })
                  }
                />
              </div>

              {/* Divider */}
              <div className="divider">{t("common:or")}</div>

              {/* Google OAuth */}
              <button
                onClick={() => {
                  void handleGoogleOAuth();
                }}
                disabled={isSubmitting}
                className="btn btn-outline btn-block gap-2"
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
                {t("onboarding:verify.continueWithGoogle")}
              </button>
            </div>
          </div>
        ) : (
          <div className="card bg-base-100 border border-base-300 shadow-xl">
            <div className="card-body">
              <h2 className="text-2xl font-bold mb-4">
                {t("onboarding:verify.completeProfile")}
              </h2>

              <form onSubmit={handleSubmitProfile} className="space-y-6">
                {/* First Name */}
                <FormInput
                  label={t("common:firstName")}
                  type="text"
                  value={firstName}
                  onChange={(e) => {
                    setFirstName(e.target.value);
                  }}
                  required
                  disabled={isSubmitting}
                  startIcon={<User className="w-5 h-5" />}
                  placeholder={t("common:firstName")}
                />

                {/* Last Name */}
                <FormInput
                  label={t("common:lastName")}
                  type="text"
                  value={lastName}
                  onChange={(e) => {
                    setLastName(e.target.value);
                  }}
                  required
                  disabled={isSubmitting}
                  startIcon={<User className="w-5 h-5" />}
                  placeholder={t("common:lastName")}
                />

                {/* Password */}
                <FormInput
                  label={t("common:password")}
                  type="password"
                  value={password}
                  onChange={(e) => {
                    setPassword(e.target.value);
                  }}
                  minLength={8}
                  required
                  disabled={isSubmitting}
                  startIcon={<Lock className="w-5 h-5" />}
                  placeholder={t("common:password")}
                  helperText={t("onboarding:verify.passwordHint")}
                />

                {/* Confirm Password */}
                <FormInput
                  label={t("common:confirmPassword")}
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => {
                    setConfirmPassword(e.target.value);
                  }}
                  minLength={8}
                  required
                  disabled={isSubmitting}
                  startIcon={<Lock className="w-5 h-5" />}
                  placeholder={t("common:confirmPassword")}
                />

                {error && (
                  <div className="alert alert-error">
                    <div className="flex flex-col gap-2">
                      <span className="text-sm">{error}</span>
                      {showLoginCta && (
                        <button
                          type="button"
                          className="btn btn-outline btn-sm self-start"
                          onClick={() => {
                            void navigate("/auth/login");
                          }}
                        >
                          {t("translation:login")}
                        </button>
                      )}
                    </div>
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
                    t("common:continue")
                  )}
                </button>
              </form>
            </div>
          </div>
        )}
      </div>
    </OnboardingLayout>
  );
}
