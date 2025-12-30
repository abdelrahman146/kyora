import { useRouterState } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitcher } from '../molecules/LanguageSwitcher'
import type { ReactNode } from 'react'


/**
 * OnboardingLayout Template
 *
 * Minimal layout for onboarding flow with:
 * - Progress bar showing completion percentage
 * - Language switcher (icon-only variant)
 * - Children content for step content
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

interface OnboardingLayoutProps {
  children: ReactNode
}

export function OnboardingLayout({ children }: OnboardingLayoutProps) {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')

  const pathname = useRouterState({
    select: (s: { location: { pathname: string } }) => s.location.pathname,
  })

  const totalSteps = 5
  const currentStep = (() => {
    switch (pathname) {
      case '/onboarding/email':
        return 1
      case '/onboarding/verify':
        return 2
      case '/onboarding/business':
        return 3
      case '/onboarding/payment':
        return 4
      case '/onboarding/complete':
        return 5
      default:
        return 0
    }
  })()

  const showProgress = pathname !== '/onboarding/plan' && currentStep > 0

  const progressPercentage = Math.round((currentStep / totalSteps) * 100)

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
          <LanguageSwitcher variant="iconOnly" />
        </div>
      </header>

      {/* Progress Bar */}
      {showProgress && (
        <div className="bg-base-100/80 backdrop-blur-md">
          <div className="container mx-auto px-4 pb-4 max-w-7xl">
            <div className="flex items-center gap-3">
              <span className="text-sm text-base-content/70 font-medium whitespace-nowrap">
                {tOnboarding('progress.step', {
                  current: currentStep,
                  total: totalSteps,
                })}
              </span>
              <div className="flex-1 h-2 bg-base-300 rounded-full overflow-hidden">
                <div
                  className="h-full bg-linear-to-r from-primary to-secondary transition-all duration-500 ease-out rounded-full"
                  style={{ width: `${progressPercentage.toString()}%` }}
                  role="progressbar"
                  aria-valuenow={currentStep}
                  aria-valuemin={0}
                  aria-valuemax={totalSteps}
                  aria-label={tOnboarding('progress.label', {
                    percentage: progressPercentage,
                  })}
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
            {children}
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="py-6 border-t border-base-300 bg-base-100">
        <div className="container mx-auto px-4 max-w-7xl">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-base-content/60">
            <p>
              {tCommon('copyright', {
                year: new Date().getFullYear(),
                company: 'Kyora',
              })}
            </p>
            <div className="flex gap-4">
              <a
                href="/privacy"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {tCommon('privacy')}
              </a>
              <a
                href="/terms"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {tCommon('terms')}
              </a>
              <a
                href="/support"
                className="hover:text-primary transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {tCommon('support')}
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
