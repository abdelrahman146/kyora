import { RotateCw, X } from 'lucide-react'
import { cn } from '@/lib/utils'

export interface FilePreviewProps {
  src: string
  alt?: string
  onRemove?: () => void
  onRetry?: () => void
  isLoading?: boolean
  progress?: number
  error?: string
  disabled?: boolean
  className?: string
}

/**
 * FilePreview - Image preview with loading/error states
 *
 * Features:
 * - Image preview with object-fit cover
 * - Loading state with progress ring
 * - Error state with retry button
 * - Remove button (top-end positioned)
 * - Optimistic rendering
 * - Touch-optimized (40px+ touch targets)
 * - RTL support
 *
 * @example
 * <FilePreview
 *   src={objectURL}
 *   alt="Product photo"
 *   isLoading={true}
 *   progress={45}
 *   onRemove={() => {}}
 * />
 */
export function FilePreview({
  src,
  alt = 'Preview',
  onRemove,
  onRetry,
  isLoading = false,
  progress = 0,
  error,
  disabled = false,
  className,
}: FilePreviewProps) {
  const hasError = Boolean(error)

  return (
    <div
      className={cn(
        'relative group',
        'aspect-square rounded-lg overflow-hidden',
        'border-2 transition-all duration-200',
        hasError ? 'border-error' : 'border-base-300',
        className,
      )}
    >
      {/* Image */}
      <img
        src={src}
        alt={alt}
        className={cn(
          'w-full h-full object-cover',
          (isLoading || hasError) && 'opacity-50',
        )}
        loading="lazy"
      />

      {/* Loading overlay */}
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-base-300/50 backdrop-blur-sm">
          <div className="relative">
            <div
              className="radial-progress text-primary"
              style={
                { '--value': progress, '--size': '3rem' } as React.CSSProperties
              }
              role="progressbar"
              aria-valuenow={progress}
              aria-valuemin={0}
              aria-valuemax={100}
              aria-label={`Uploading ${progress}%`}
            >
              <span className="text-xs font-semibold">{progress}%</span>
            </div>
          </div>
        </div>
      )}

      {/* Error overlay */}
      {hasError && (
        <div className="absolute inset-0 flex flex-col items-center justify-center bg-error/10 backdrop-blur-sm">
          <p className="text-error text-xs font-medium mb-2 px-2 text-center">
            {error}
          </p>
          {onRetry && (
            <button
              type="button"
              onClick={onRetry}
              className="btn btn-sm btn-error btn-outline"
              aria-label="Retry upload"
            >
              <RotateCw className="w-4 h-4" />
              Retry
            </button>
          )}
        </div>
      )}

      {/* Remove button */}
      {onRemove && (
        <button
          type="button"
          onClick={onRemove}
          disabled={disabled}
          className={cn(
            'absolute top-2 end-2 z-10',
            'btn btn-sm btn-circle btn-error',
            'opacity-0 group-hover:opacity-100 focus:opacity-100',
            'transition-opacity duration-200',
            'touch-manipulation',
            disabled && 'btn-disabled',
          )}
          aria-label={`Remove ${alt}`}
        >
          <X className="w-4 h-4" />
        </button>
      )}

      {/* Success checkmark */}
      {!isLoading && !hasError && (
        <div className="absolute top-2 start-2 bg-success text-success-content rounded-full p-1">
          <svg
            className="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M5 13l4 4L19 7"
            />
          </svg>
        </div>
      )}
    </div>
  )
}
