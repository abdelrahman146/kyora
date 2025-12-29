import React from 'react';
import { AlertCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';

interface Props {
  children: React.ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

/**
 * ErrorBoundary Component
 *
 * Catches render-time errors in child components and displays a friendly fallback UI.
 * Prevents the entire app from crashing due to component errors.
 *
 * Features:
 * - Friendly error message with consistent styling
 * - Retry button to attempt re-rendering
 * - Follows design system (DaisyUI + branding)
 * - Mobile-first and RTL-aware
 * - Accessible with ARIA labels
 *
 * Usage:
 * ```tsx
 * <ErrorBoundary>
 *   <YourComponent />
 * </ErrorBoundary>
 * ```
 */
export class ErrorBoundary extends React.Component<Props, State> {
  public override state: State = { hasError: false, error: null };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public override componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // In production, you could log this to an error reporting service
    if (import.meta.env.DEV) {
      console.error('Error caught by ErrorBoundary:', error, errorInfo);
    }
  }

  private handleRetry = () => {
    this.setState({ hasError: false, error: null });
  };

  public override render() {
    if (!this.state.hasError) return this.props.children;

    return <ErrorBoundaryFallback onRetry={this.handleRetry} />;
  }
}

/**
 * ErrorBoundaryFallback Component
 *
 * Displays the fallback UI when an error is caught by ErrorBoundary.
 * Separated as a functional component to support hooks (useTranslation).
 */
function ErrorBoundaryFallback({ onRetry }: { onRetry: () => void }) {
  const { t } = useTranslation();

  return (
    <div className="flex min-h-[400px] items-center justify-center p-4">
      <div className="card bg-base-100 border border-base-300 shadow-sm w-full max-w-md">
        <div className="card-body text-center">
          {/* Error Icon */}
          <div className="flex justify-center mb-4">
            <div className="w-16 h-16 bg-error/10 rounded-full flex items-center justify-center">
              <AlertCircle className="text-error" size={32} />
            </div>
          </div>

          {/* Error Message */}
          <h2 className="card-title text-xl justify-center mb-2">
            {t('errors:generic.unexpected_title', { defaultValue: 'Something went wrong' })}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('errors:generic.unexpected_description', {
              defaultValue: 'An unexpected error occurred. Please try again.',
            })}
          </p>

          {/* Retry Button */}
          <div className="card-actions justify-center">
            <button
              type="button"
              onClick={onRetry}
              className="btn btn-primary"
              aria-label={t('common.retry', { defaultValue: 'Retry' })}
            >
              {t('common.retry', { defaultValue: 'Retry' })}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
