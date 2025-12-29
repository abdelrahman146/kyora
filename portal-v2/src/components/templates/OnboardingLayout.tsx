import { Outlet } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitcher } from '../molecules/LanguageSwitcher'
import { getCurrentStageNumber } from '@/stores/onboardingStore'

/**
 * OnboardingLayout Template
 *
 * Minimal layout for onboarding flow with:
 * - Progress bar showing completion percentage
 * - Language switcher (icon-only variant)
 * - Outlet for step content
 * - Clean, distraction-free design matching portal-web
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
  const { t } = useTranslation(['onboarding', 'common'])

  // Calculate progress percentage based on current stage
  const currentStage = getCurrentStageNumber()
  const totalStages = 6
  const progressPercentage = Math.round((currentStage / totalStages) * 100)

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
          <LanguageSwitcher variant="icon" />
        </div>
      </header>

      {/* Progress Bar */}
      {currentStage > 0 && (
        <div className="bg-base-100/80 backdrop-blur-md">
          <div className="container mx-auto px-4 pb-4 max-w-7xl">
            <div className="flex items-center gap-3">
              <span className="text-sm text-base-content/70 font-medium whitespace-nowrap">
                {t('onboarding:progress.step', { current: currentStage, total: totalStages })}
              </span>
              <div className="flex-1 h-2 bg-base-300 rounded-full overflow-hidden">
                <div
                  className="h-full bg-linear-to-r from-primary to-secondary transition-all duration-500 ease-out rounded-full"
                  style={{ width: `${progressPercentage.toString()}%` }}
                  role="progressbar"
                  aria-valuenow={currentStage}
                  aria-valuemin={0}
                  aria-valuemax={totalStages}
                  aria-label={t('onboarding:progress.label', { percentage: progressPercentage })}
                />
              </div>
              <span className="text-sm text-base-content/70 font-medium whitespace-nowrap">
                {progressPercentage}%
              </span>
            </div>
          </div>
        </div>
      )}

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 py-8 max-w-5xl">
        <div className="flex items-center justify-center min-h-[calc(100vh-200px)]">
          <div className="w-full">
            <Outlet />
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="py-6 border-t border-base-300 bg-base-100">
        <div className="container mx-auto px-4 max-w-7xl">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-base-content/60">
            <p>
              {t('common:copyright', { year: new Date().getFullYear(), company: 'Kyora' })}
            </p>
            <div className="flex gap-4">
              <a
                href="/privacy"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {t('common:privacy')}
              </a>
              <a
                href="/terms"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {t('common:terms')}
              </a>
              <a
                href="/support"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {t('common:support')}
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
