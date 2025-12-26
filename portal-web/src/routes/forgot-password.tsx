import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";
import { Mail, ArrowLeft, CheckCircle } from "lucide-react";
import { z } from "zod";
import toast from "react-hot-toast";
import { Input } from "../components/atoms/Input";
import { Button } from "../components/atoms/Button";
import { authApi } from "../api/auth";
import { useLanguage } from "../hooks/useLanguage";
import { translateErrorAsync } from "../lib/translateError";

// Validation schema for forgot password form
const forgotPasswordSchema = z.object({
  email: z
    .string()
    .min(1, "errors.validation.required")
    .pipe(z.email("errors.validation.invalid_email")),
});

type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

/**
 * Forgot Password Page
 *
 * Features:
 * - Email input with validation
 * - Success state with instructions
 * - Error handling with translation
 * - RTL support
 * - Rate limiting protection
 * - Accessible form controls
 */
export default function ForgotPasswordPage() {
  const { t } = useTranslation();
  const { t: tErrors } = useTranslation("errors");
  const { isRTL } = useLanguage();
  const navigate = useNavigate();
  const [isSuccess, setIsSuccess] = useState(false);
  const [submittedEmail, setSubmittedEmail] = useState("");

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ForgotPasswordFormData>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = async (data: ForgotPasswordFormData) => {
    try {
      await authApi.forgotPassword({ email: data.email });

      // Show success state (even if email doesn't exist - security best practice)
      setSubmittedEmail(data.email);
      setIsSuccess(true);

      void toast.success(t("auth.password_reset_email_sent"), {
        duration: 5000,
        position: isRTL ? "top-right" : "top-left",
      });
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      void toast.error(message, {
        duration: 4000,
        position: isRTL ? "top-right" : "top-left",
      });
    }
  };

  const handleBackToLogin = () => {
    void navigate("/login", { replace: true });
  };

  // Success State
  if (isSuccess) {
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
                {t("auth.password_reset_sent_title")}
              </h1>

              {/* Instructions */}
              <p className="text-base-content/70 mb-6">
                {t("auth.password_reset_sent_description", {
                  email: submittedEmail,
                })}
              </p>

              {/* Additional Information */}
              <div className="alert alert-info mb-6">
                <div className="flex flex-col gap-2 text-start">
                  <p className="text-sm font-medium">
                    {t("auth.password_reset_email_tips_title")}
                  </p>
                  <ul className="text-xs list-disc list-inside space-y-1">
                    <li>{t("auth.password_reset_tip_check_spam")}</li>
                    <li>{t("auth.password_reset_tip_expires")}</li>
                    <li>{t("auth.password_reset_tip_no_account")}</li>
                  </ul>
                </div>
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

  // Request Form State
  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="card bg-base-200 shadow-xl">
          <div className="card-body">
            {/* Back Button */}
            <button
              type="button"
              onClick={handleBackToLogin}
              className="btn btn-ghost btn-sm w-fit mb-4 -ms-2"
              aria-label={t("auth.back_to_login")}
            >
              <ArrowLeft size={20} className={isRTL ? "rotate-180" : ""} />
              {t("auth.back_to_login")}
            </button>

            {/* Header */}
            <h1 className="card-title text-3xl mb-2">
              {t("auth.forgot_password_title")}
            </h1>
            <p className="text-base-content/70 mb-6">
              {t("auth.forgot_password_description")}
            </p>

            {/* Form */}
            <form
              onSubmit={(e) => {
                void handleSubmit(onSubmit)(e);
              }}
              className="space-y-6"
              noValidate
            >
              {/* Email Input */}
              <Input
                {...register("email")}
                id="email"
                type="email"
                label={t("auth.email")}
                placeholder={t("auth.email_placeholder")}
                error={
                  errors.email?.message
                    ? tErrors(errors.email.message)
                    : undefined
                }
                startIcon={<Mail size={20} />}
                autoComplete="email"
                autoFocus
                disabled={isSubmitting}
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
                {t("auth.send_reset_link")}
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
