import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";
import { CheckCircle, AlertCircle, Loader2 } from "lucide-react";
import { z } from "zod";
import { PasswordInput } from "../components/atoms/PasswordInput";
import { Button } from "../components/atoms/Button";
import { authApi } from "../api/auth";
import { useLanguage } from "../hooks/useLanguage";
import { translateErrorAsync } from "../lib/translateError";

// Validation schema for reset password form
const resetPasswordSchema = z
  .object({
    newPassword: z
      .string()
      .min(1, "validation.required")
      .min(8, "validation.invalid_password"),
    confirmPassword: z.string().min(1, "validation.required"),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    path: ["confirmPassword"],
    message: "validation.password_mismatch",
  });

type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;

type PageStatus = "loading" | "ready" | "success" | "error";

/**
 * Reset Password Page
 *
 * Features:
 * - Token extraction from URL
 * - New password and confirmation with validation
 * - Password strength requirements
 * - Success and error states
 * - Token validation
 * - RTL support
 * - Accessible form controls
 *
 * Flow:
 * 1. Extract token from URL query params
 * 2. Validate token exists
 * 3. User enters new password and confirmation
 * 4. Submit to backend
 * 5. Show success and redirect to login
 */
export default function ResetPasswordPage() {
  const { t } = useTranslation();
  const { t: tErrors } = useTranslation("errors");
  useLanguage();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [pageStatus, setPageStatus] = useState<PageStatus>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [token, setToken] = useState("");
  const [submitErrorMessage, setSubmitErrorMessage] = useState<string>("");

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      newPassword: "",
      confirmPassword: "",
    },
  });

  // Extract and validate token on mount
  useEffect(() => {
    const tokenFromUrl = searchParams.get("token");

    if (!tokenFromUrl) {
      queueMicrotask(() => {
        setPageStatus("error");
        setErrorMessage(t("auth.reset_password_missing_token"));
      });
      return;
    }

    queueMicrotask(() => {
      setToken(tokenFromUrl);
      setPageStatus("ready");
    });
  }, [searchParams, t]);

  const onSubmit = async (data: ResetPasswordFormData) => {
    try {
      setSubmitErrorMessage("");
      await authApi.resetPassword({
        token,
        password: data.newPassword,
      });

      setPageStatus("success");

      // Redirect to login after 2 seconds
      void setTimeout(() => {
        void navigate("/auth/login", { replace: true });
      }, 2000);
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      setSubmitErrorMessage(message);
    }
  };

  const handleBackToLogin = () => {
    void navigate("/login", { replace: true });
  };

  // Loading State
  if (pageStatus === "loading") {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="text-center">
          <Loader2 className="animate-spin text-primary mx-auto mb-4" size={48} />
          <p className="text-base-content/70">
            {t("auth.reset_password_validating")}
          </p>
        </div>
      </div>
    );
  }

  // Error State (Invalid/Missing Token)
  if (pageStatus === "error") {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-200 shadow-xl">
            <div className="card-body items-center text-center">
              {/* Error Icon */}
              <div className="w-16 h-16 rounded-full bg-error/20 flex items-center justify-center mb-4">
                <AlertCircle className="text-error" size={32} />
              </div>

              {/* Error Title */}
              <h1 className="card-title text-2xl mb-2">
                {t("auth.reset_password_error_title")}
              </h1>

              {/* Error Message */}
              <p className="text-base-content/70 mb-6">{errorMessage}</p>

              {/* Actions */}
              <div className="w-full space-y-3">
                <Button
                  type="button"
                  variant="primary"
                  size="lg"
                  fullWidth
                  onClick={handleBackToLogin}
                >
                  {t("auth.return_to_login")}
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="lg"
                  fullWidth
                  onClick={() => {
                    void navigate("/forgot-password");
                  }}
                >
                  {t("auth.request_new_link")}
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Success State
  if (pageStatus === "success") {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-200 shadow-xl">
            <div className="card-body items-center text-center">
              {/* Success Icon */}
              <div className="w-16 h-16 rounded-full bg-success/20 flex items-center justify-center mb-4">
                <CheckCircle className="text-success" size={32} />
              </div>

              {/* Success Title */}
              <h1 className="card-title text-2xl mb-2">
                {t("auth.password_reset_success_title")}
              </h1>

              {/* Success Description */}
              <p className="text-base-content/70 mb-6">
                {t("auth.password_reset_success_description")}
              </p>

              {/* Redirecting Message */}
              <div className="flex items-center gap-2 text-sm text-base-content/60 mb-4">
                <Loader2 className="animate-spin" size={16} />
                <span>{t("auth.redirecting_to_login")}</span>
              </div>

              {/* Back to Login Button */}
              <Button
                type="button"
                variant="primary"
                size="lg"
                fullWidth
                onClick={handleBackToLogin}
              >
                {t("auth.return_to_login")}
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Reset Password Form (Ready State)
  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="card bg-base-200 shadow-xl">
          <div className="card-body">
            {/* Header */}
            <h1 className="card-title text-3xl mb-2">
              {t("auth.reset_password_title")}
            </h1>
            <p className="text-base-content/70 mb-6">
              {t("auth.reset_password_description")}
            </p>

            {/* Form */}
            <form
              onSubmit={(e) => {
                void handleSubmit(onSubmit)(e);
              }}
              className="space-y-6"
              noValidate
            >
              {submitErrorMessage ? (
                <div role="alert" className="alert alert-error">
                  <span>{submitErrorMessage}</span>
                </div>
              ) : null}

              {/* New Password Input */}
              <PasswordInput
                {...register("newPassword")}
                id="newPassword"
                label={t("auth.new_password")}
                placeholder={t("auth.new_password_placeholder")}
                error={
                  errors.newPassword?.message
                    ? tErrors(errors.newPassword.message)
                    : undefined
                }
                helperText={t("auth.password_requirements")}
                autoComplete="new-password"
                autoFocus
                disabled={isSubmitting}
                showPasswordToggle
              />

              {/* Confirm Password Input */}
              <PasswordInput
                {...register("confirmPassword")}
                id="confirmPassword"
                label={t("auth.confirm_password")}
                placeholder={t("auth.confirm_password_placeholder")}
                error={
                  errors.confirmPassword?.message
                    ? tErrors(errors.confirmPassword.message)
                    : undefined
                }
                autoComplete="new-password"
                disabled={isSubmitting}
                showPasswordToggle
              />

              {/* Submit Button */}
              <Button
                type="submit"
                variant="primary"
                size="lg"
                fullWidth
                loading={isSubmitting}
                disabled={isSubmitting}
              >
                {t("auth.reset_password_submit")}
              </Button>

              {/* Back to Login Link */}
              <div className="text-center">
                <p className="text-sm text-base-content/60">
                  {t("auth.remember_password")}{" "}
                  <button
                    type="button"
                    onClick={handleBackToLogin}
                    className="text-primary hover:text-primary-focus hover:underline transition-colors font-medium"
                  >
                    {t("auth.login")}
                  </button>
                </p>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}
