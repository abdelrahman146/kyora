import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Link } from "react-router-dom";
import { Mail, Lock, Loader2 } from "lucide-react";
import { Input } from "../atoms/Input";
import { Button } from "../atoms/Button";
import { loginSchema, type LoginFormData } from "../../schemas/auth";
import { useTranslation } from "react-i18next";

interface LoginFormProps {
  onSubmit: (data: LoginFormData) => Promise<void>;
  onGoogleLogin?: () => void;
  isGoogleLoading?: boolean;
}

/**
 * Login Form Component
 *
 * Features:
 * - Email and password validation with Zod
 * - React Hook Form integration
 * - Loading states during submission
 * - RTL support
 * - Accessible form controls
 * - Google OAuth button
 *
 * @example
 * ```tsx
 * <LoginForm
 *   onSubmit={handleLogin}
 *   onGoogleLogin={handleGoogleLogin}
 * />
 * ```
 */
export function LoginForm({ onSubmit, onGoogleLogin, isGoogleLoading = false }: LoginFormProps) {
  const { t } = useTranslation();
  const { t: tErrors } = useTranslation("errors");

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  return (
    <form onSubmit={(e) => { void handleSubmit(onSubmit)(e); }} className="space-y-6" noValidate>
      {/* Email Input */}
      <Input
        {...register("email")}
        id="email"
        type="email"
        label={t("auth.email")}
        placeholder={t("auth.email_placeholder")}
        error={errors.email?.message ? tErrors(errors.email.message) : undefined}
        startIcon={<Mail size={20} />}
        autoComplete="email"
        disabled={isSubmitting}
      />

      {/* Password Input */}
      <Input
        {...register("password")}
        id="password"
        type="password"
        label={t("auth.password")}
        placeholder={t("auth.password_placeholder")}
        error={errors.password?.message ? tErrors(errors.password.message) : undefined}
        startIcon={<Lock size={20} />}
        autoComplete="current-password"
        disabled={isSubmitting}
      />

      {/* Forgot Password Link */}
      <div className="text-end">
        <Link
          to="/forgot-password"
          className="text-sm text-primary hover:text-primary-focus hover:underline transition-colors"
        >
          {t("auth.forgot_password")}
        </Link>
      </div>

      {/* Submit Button */}
      <Button
        type="submit"
        variant="primary"
        size="lg"
        fullWidth
        loading={isSubmitting}
        disabled={isSubmitting}
      >
        {isSubmitting ? (
          <>
            <Loader2 className="animate-spin" size={20} />
            {t("auth.logging_in")}
          </>
        ) : (
          t("auth.login")
        )}
      </Button>

      {/* Divider */}
      {onGoogleLogin && (
        <>
          <div className="divider text-neutral-500 text-sm">
            {t("auth.or_continue_with")}
          </div>

          {/* Google Login Button */}
          <button
            type="button"
            onClick={onGoogleLogin}
            disabled={isSubmitting || isGoogleLoading}
            className="btn btn-outline btn-lg w-full h-[52px] rounded-xl border-2 border-neutral-200 hover:border-neutral-300 hover:bg-base-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            {isGoogleLoading ? (
              <>
                <Loader2 size={20} className="animate-spin" />
                <span className="font-semibold">{t("auth.connecting_google")}</span>
              </>
            ) : (
              <>
                <svg
                  viewBox="0 0 24 24"
                  className="w-5 h-5"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                    fill="#4285F4"
                  />
                  <path
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                    fill="#34A853"
                  />
                  <path
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                    fill="#FBBC05"
                  />
                  <path
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                    fill="#EA4335"
                  />
                </svg>
                <span className="font-semibold">{t("auth.continue_with_google")}</span>
              </>
            )}
          </button>
        </>
      )}
    </form>
  );
}
