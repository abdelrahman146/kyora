/**
 * RouteErrorFallback Component
 *
 * Generic error boundary fallback for TanStack Router routes.
 * Displays user-friendly error message with retry action.
 *
 * Used as `errorComponent` in route definitions.
 */

import { useRouter } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { AlertCircle, RefreshCcw } from 'lucide-react'

import type { ErrorComponentProps } from '@tanstack/react-router'

export function RouteErrorFallback({ error, reset }: ErrorComponentProps) {
  const { t: tErrors } = useTranslation('errors')
  const router = useRouter()

  const handleRetry = () => {
    // Reset the error boundary
    reset()
    // Invalidate all queries to refetch
    router.invalidate()
  }

  const handleGoHome = () => {
    router.navigate({ to: '/' })
  }

  // Parse error message
  const errorMessage =
    error instanceof Error ? error.message : tErrors('route.unknown_error')

  return (
    <div className="flex min-h-[400px] items-center justify-center p-4">
      <div className="card bg-base-100 shadow-lg max-w-md w-full">
        <div className="card-body items-center text-center">
          <AlertCircle className="w-16 h-16 text-error mb-4" />

          <h2 className="card-title text-error mb-2">
            {tErrors('route.route_error_title')}
          </h2>

          <p className="text-base-content/70 mb-6">{errorMessage}</p>

          <div className="card-actions justify-center gap-2 flex-col sm:flex-row w-full">
            <button
              type="button"
              onClick={handleRetry}
              className="btn btn-primary gap-2"
            >
              <RefreshCcw className="w-4 h-4" />
              {tErrors('route.retry')}
            </button>

            <button
              type="button"
              onClick={handleGoHome}
              className="btn btn-ghost"
            >
              {tErrors('route.go_home')}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
