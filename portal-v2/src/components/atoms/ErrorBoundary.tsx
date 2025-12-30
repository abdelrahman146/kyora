import { Component } from 'react'
import { AlertTriangle } from 'lucide-react'
import { Button } from './Button'
import type { ErrorInfo, ReactNode } from 'react'

export interface ErrorBoundaryProps {
  children: ReactNode
  fallback?: (error: Error, reset: () => void) => ReactNode
  onError?: (error: Error, errorInfo: ErrorInfo) => void
  compact?: boolean
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

/**
 * ErrorBoundary Component
 *
 * Component-level error boundary with inline fallback UI.
 * Features:
 * - Smart retry that resets state
 * - Maintains component layout space
 * - Optional compact mode for smaller organisms
 * - Logs errors in development
 */
export class ErrorBoundary extends Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Log error in development
    if (import.meta.env.DEV) {
      console.error('ErrorBoundary caught an error:', error, errorInfo)
    }

    // Call optional error handler
    this.props.onError?.(error, errorInfo)
  }

  reset = () => {
    this.setState({ hasError: false, error: null })
  }

  render() {
    if (this.state.hasError && this.state.error) {
      // Custom fallback
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.reset)
      }

      // Default inline fallback
      const { compact } = this.props

      return (
        <div
          className={cn(
            'flex flex-col items-center justify-center rounded-lg border border-error/20 bg-error/5 p-6',
            compact ? 'gap-2' : 'gap-4',
          )}
        >
          <AlertTriangle className="text-error" size={compact ? 24 : 32} />
          <div className="text-center">
            <p
              className={cn(
                'font-semibold text-error',
                compact ? 'text-sm' : 'text-base',
              )}
            >
              حدث خطأ
            </p>
            {!compact && (
              <p className="mt-1 text-sm text-base-content/60">
                {this.state.error.message || 'حدث خطأ غير متوقع'}
              </p>
            )}
          </div>
          <Button
            variant="outline"
            size={compact ? 'xs' : 'sm'}
            onClick={this.reset}
          >
            إعادة المحاولة
          </Button>
        </div>
      )
    }

    return this.props.children
  }
}

// Helper to fix import
function cn(...inputs: Array<string | boolean | undefined | null>) {
  return inputs.filter(Boolean).join(' ')
}
