import { Outlet } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import { Languages } from 'lucide-react'
import { onboardingStore, getCurrentStageNumber } from '@/stores/onboardingStore'
import { useLanguage } from '@/hooks/useLanguage'

/**
 * OnboardingLayout Template
 *
 * Minimal layout for onboarding flow with:
 * - Progress bar showing completion percentage
 * - Language switcher (icon-only variant)
 * - Outlet for step content
 *
 * Progress stages:
 * 1. Plan selection (16%)
 * 2. Email verification (33%)
 * 3. Identity verified (50%)
 * 4. Business setup (66%)
 * 5. Payment (83%)
 * 6. Complete (100%)
 */
export function OnboardingLayout() {
  const { t } = useTranslation()
  const { language, toggleLanguage, isRTL } = useLanguage()
  const state = useStore(onboardingStore)

  // Calculate progress percentage based on current stage
  const currentStage = getCurrentStageNumber()
  const totalStages = 6
  const progressPercentage = Math.round((currentStage / totalStages) * 100)

  return (
    <div className="min-h-screen bg-base-100" dir={isRTL ? 'rtl' : 'ltr'}>
      {/* Header with progress bar and language switcher */}
      <header className="border-b border-base-300 bg-base-100">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between gap-4">
            {/* Logo */}
            <div className="flex items-center gap-3">
              <h1 className="text-2xl font-bold text-primary">Kyora</h1>
              {state.stage && (
                <span className="text-sm text-base-content/60">
                  {t('onboarding:setup')}
                </span>
              )}
            </div>

            {/* Language Switcher - Icon only */}
            <button
              onClick={toggleLanguage}
              className="btn btn-ghost btn-sm btn-circle"
              aria-label={t('common:change_language')}
              title={t('common:change_language')}
            >
              <Languages className="w-5 h-5" />
            </button>
          </div>

          {/* Progress Bar */}
          {currentStage > 0 && (
            <div className="mt-4">
              <div className="flex items-center justify-between text-sm mb-2">
                <span className="text-base-content/60">
                  {t('onboarding:progress')}
                </span>
                <span className="font-medium text-primary">
                  {progressPercentage}%
                </span>
              </div>
              <progress
                className="progress progress-primary w-full"
                value={progressPercentage}
                max="100"
              ></progress>
            </div>
          )}
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="border-t border-base-300 bg-base-100 py-6 mt-auto">
        <div className="container mx-auto px-4 text-center">
          <p className="text-sm text-base-content/60">
            © 2025 Kyora. {t('common:all_rights_reserved')}
          </p>
        </div>
      </footer>
    </div>
  )
}
