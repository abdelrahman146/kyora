/**
 * FormSkeleton Component
 *
 * Content-aware skeleton for forms.
 * Matches actual form structure: labels + inputs + buttons
 */

interface FormSkeletonProps {
  /**
   * Number of form fields to show
   * @default 6
   */
  fields?: number
  /**
   * Show cancel button alongside submit
   * @default false
   */
  showCancel?: boolean
}

export function FormSkeleton({
  fields = 6,
  showCancel = false,
}: FormSkeletonProps) {
  return (
    <div className="space-y-6 animate-pulse">
      {/* Form Fields */}
      <div className="space-y-4">
        {Array.from({ length: fields }).map((_, i) => (
          <div key={i} className="form-control">
            <div className="h-4 w-24 bg-base-300 rounded mb-2" />
            <div className="h-12 w-full bg-base-300 rounded" />
          </div>
        ))}
      </div>

      {/* Form Actions */}
      <div
        className={`flex gap-3 ${showCancel ? 'justify-end' : 'justify-start'}`}
      >
        {showCancel && <div className="h-12 w-24 bg-base-300 rounded" />}
        <div className="h-12 w-32 bg-base-300 rounded" />
      </div>
    </div>
  )
}
