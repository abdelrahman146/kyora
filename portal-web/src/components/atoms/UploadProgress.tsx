import { X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { formatSize } from '@/lib/upload/fileValidator'

export interface UploadProgressProps {
  fileName: string
  fileSize: number
  progress: number
  status?: string
  onCancel?: () => void
  className?: string
}

/**
 * UploadProgress - Progress indicator for file uploads
 *
 * Features:
 * - Progress bar with percentage
 * - File name and size display
 * - Cancel button
 * - Status text
 * - Mobile-optimized layout
 * - RTL support
 *
 * @example
 * <UploadProgress
 *   fileName="photo.jpg"
 *   fileSize={1024000}
 *   progress={65}
 *   status="Uploading..."
 *   onCancel={() => {}}
 * />
 */
export function UploadProgress({
  fileName,
  fileSize,
  progress,
  status = 'Uploading...',
  onCancel,
  className,
}: UploadProgressProps) {
  return (
    <div
      className={cn(
        'flex flex-col gap-2 p-4 rounded-lg',
        'bg-base-200 border border-base-300',
        className,
      )}
    >
      {/* Header: File info and cancel button */}
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <p
            className="text-sm font-medium text-base-content truncate"
            title={fileName}
          >
            {fileName}
          </p>
          <p className="text-xs text-base-content/60">{formatSize(fileSize)}</p>
        </div>

        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="btn btn-sm btn-circle btn-ghost flex-shrink-0"
            aria-label="Cancel upload"
          >
            <X className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* Progress bar */}
      <div className="relative">
        <progress
          className="progress progress-primary w-full"
          value={progress}
          max={100}
          aria-label={`Upload progress: ${progress}%`}
        />
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <span className="text-xs font-semibold text-base-content/70 mix-blend-difference">
            {progress}%
          </span>
        </div>
      </div>

      {/* Status text */}
      {status && (
        <p
          className="text-xs text-base-content/60"
          role="status"
          aria-live="polite"
        >
          {status}
        </p>
      )}
    </div>
  )
}
