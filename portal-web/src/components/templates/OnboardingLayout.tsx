import { type ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { useLanguage } from "@/hooks/useLanguage";
import { cn } from "@/lib/utils";

interface OnboardingLayoutProps {
  children: ReactNode;
  currentStep?: number;
  totalSteps?: number;
  showProgress?: boolean;
}

/**
 * OnboardingLayout Component
 *
 * Simple, focused layout for the onboarding flow without sidebar navigation.
 *
 * Features:
 * - Clean, distraction-free layout
 * - Progress indicator (optional)
 * - Language switcher in header
 * - Mobile-first responsive design
 * - RTL support with logical properties
 * - Smooth animations and transitions
 * - Safe area padding for mobile devices
 *
 * Layout Structure:
 * ```
 * ┌────────────────────────────────────┐
 * │    Header (Logo + Lang Switch)    │
 * ├────────────────────────────────────┤
 * │    Progress Bar (if enabled)      │
 * ├────────────────────────────────────┤
 * │                                    │
 * │         Content Area               │
 * │      (Centered, Scrollable)        │
 * │                                    │
 * └────────────────────────────────────┘
 * ```
 *
 * @example
 * ```tsx
 * <OnboardingLayout currentStep={2} totalSteps={5} showProgress>
 *   <PlanSelectionStep />
 * </OnboardingLayout>
 * ```
 */
export function OnboardingLayout({
  children,
  currentStep = 0,
  totalSteps = 5,
  showProgress = true,
}: OnboardingLayoutProps) {
  const { t } = useTranslation(["onboarding", "common"]);
  const { language, toggleLanguage } = useLanguage();

  const progressPercentage = totalSteps > 0 ? (currentStep / totalSteps) * 100 : 0;

  return (
    <div className="min-h-screen bg-linear-to-br from-base-100 to-base-200 flex flex-col">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-base-100/80 backdrop-blur-md border-b border-base-300">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between max-w-7xl">
          {/* Logo */}
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
              <span className="text-primary-content font-bold text-lg">K</span>
            </div>
            <h1 className="text-xl font-bold text-base-content">Kyora</h1>
          </div>

          {/* Language Switcher */}
          <button
            onClick={toggleLanguage}
            className="btn btn-ghost btn-sm gap-2"
            aria-label={t("common:switchLanguage")}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"
              />
            </svg>
            <span className="font-medium">{language.toUpperCase()}</span>
          </button>
        </div>
      </header>

      {/* Progress Bar */}
      {showProgress && totalSteps > 0 && (
        <div className="sticky top-18 z-40 bg-base-100/80 backdrop-blur-md border-b border-base-300">
          <div className="container mx-auto px-4 py-3 max-w-7xl">
            <div className="flex items-center gap-3">
              <span className="text-sm text-base-content/70 font-medium whitespace-nowrap">
                {t("onboarding.progress.step", { current: currentStep, total: totalSteps })}
              </span>
              <div className="flex-1 h-2 bg-base-300 rounded-full overflow-hidden">
                <div
                  className="h-full bg-linear-to-r from-primary to-secondary transition-all duration-500 ease-out rounded-full"
                  style={{ width: `${progressPercentage}%` }}
                  role="progressbar"
                  aria-valuenow={currentStep}
                  aria-valuemin={0}
                  aria-valuemax={totalSteps}
                  aria-label={t("onboarding.progress.label", { percentage: Math.round(progressPercentage) })}
                />
              </div>
              <span className="text-sm text-base-content/70 font-medium whitespace-nowrap">
                {Math.round(progressPercentage)}%
              </span>
            </div>
          </div>
        </div>
      )}

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 py-8 max-w-4xl">
        <div className="flex items-center justify-center min-h-[calc(100vh-200px)]">
          <div className="w-full">
            {children}
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="py-6 border-t border-base-300 bg-base-100">
        <div className="container mx-auto px-4 max-w-7xl">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-base-content/60">
            <p>
              {t("common:copyright", { year: new Date().getFullYear(), company: "Kyora" })}
            </p>
            <div className="flex gap-4">
              <a
                href="/privacy"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {t("common:privacy")}
              </a>
              <a
                href="/terms"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {t("common:terms")}
              </a>
              <a
                href="/support"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {t("common:support")}
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
