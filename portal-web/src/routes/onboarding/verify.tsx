import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Mail, Check } from "lucide-react";
import { OnboardingLayout } from "@/components/templates";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { authApi } from "@/api/auth";
import { translateErrorAsync } from "@/lib/translateError";

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
  const { t } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const {
    sessionToken,
    email,
    stage,
    markEmailVerified,
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
  const [canResend, setCanResend] = useState(false);
  const [resendCooldown, setResendCooldown] = useState(0);

  const otpInputRefs = useRef<(HTMLInputElement | null)[]>([]);

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

  // Start resend cooldown
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | undefined;

    if (resendCooldown > 0) {
      timer = setTimeout(() => {
        setResendCooldown(resendCooldown - 1);
      }, 1000);
    } else {
      setCanResend(true);
    }

    return () => {
      if (timer) clearTimeout(timer);
    };
  }, [resendCooldown]);

  const sendOTP = async () => {
    if (!sessionToken) return;

    try {
      setError("");
      setIsSubmitting(true);
      await onboardingApi.sendEmailOTP({ sessionToken });
      setSuccess(t("onboarding:verify.otpSent"));
      setCanResend(false);
      setResendCooldown(30);
    } catch (err) {
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

      markEmailVerified();
      void navigate("/onboarding/business");
    } catch (err) {
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
                onClick={handleVerifyOtp}
                disabled={otpCode.join("").length !== 6 || isSubmitting}
                className="btn btn-primary btn-block mb-4"
              >
                {t("onboarding:verify.verifyCode")}
              </button>

              {/* Resend OTP */}
              <div className="text-center">
                <button
                  onClick={() => {
                    void sendOTP();
                  }}
                  disabled={!canResend || isSubmitting}
                  className="btn btn-ghost btn-sm"
                >
                  {canResend
                    ? t("onboarding:verify.resendCode")
                    : t("onboarding:verify.resendIn", { seconds: resendCooldown })}
                </button>
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

              <form onSubmit={handleSubmitProfile} className="space-y-4">
                {/* First Name */}
                <div className="form-control">
                  <label htmlFor="firstName" className="label">
                    <span className="label-text font-medium">
                      {t("common:firstName")}
                    </span>
                  </label>
                    <input
                      type="text"
                      id="firstName"
                      value={firstName}
                      onChange={(e) => {
                        setFirstName(e.target.value);
                      }}
                    className="input input-bordered"
                    required
                    disabled={isSubmitting}
                  />
                </div>

                {/* Last Name */}
                <div className="form-control">
                  <label htmlFor="lastName" className="label">
                    <span className="label-text font-medium">
                      {t("common:lastName")}
                    </span>
                  </label>
                    <input
                      type="text"
                      id="lastName"
                      value={lastName}
                      onChange={(e) => {
                        setLastName(e.target.value);
                      }}
                    className="input input-bordered"
                    required
                    disabled={isSubmitting}
                  />
                </div>

                {/* Password */}
                <div className="form-control">
                  <label htmlFor="password" className="label">
                    <span className="label-text font-medium">
                      {t("common:password")}
                    </span>
                  </label>
                    <input
                      type="password"
                      id="password"
                      value={password}
                      onChange={(e) => {
                        setPassword(e.target.value);
                      }}
                    className="input input-bordered"
                    minLength={8}
                    required
                    disabled={isSubmitting}
                  />
                  <label className="label">
                    <span className="label-text-alt">
                      {t("onboarding:verify.passwordHint")}
                    </span>
                  </label>
                </div>

                {/* Confirm Password */}
                <div className="form-control">
                  <label htmlFor="confirmPassword" className="label">
                    <span className="label-text font-medium">
                      {t("common:confirmPassword")}
                    </span>
                  </label>
                    <input
                      type="password"
                      id="confirmPassword"
                      value={confirmPassword}
                      onChange={(e) => {
                        setConfirmPassword(e.target.value);
                      }}
                    className="input input-bordered"
                    minLength={8}
                    required
                    disabled={isSubmitting}
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
