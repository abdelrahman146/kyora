/**
 * FilterButton Component
 *
 * A production-grade, unified component that combines a filter trigger button
 * with a filter drawer/sheet. Designed for maximum reusability across resources.
 *
 * Features:
 * - ✅ Consistent button styling (matches form fields and design system)
 * - ✅ Integrated drawer (mobile: bottom sheet, desktop: side drawer)
 * - ✅ Built-in state management (open/close)
 * - ✅ Apply and Reset callbacks
 * - ✅ Active filter indicator (badge count)
 * - ✅ Full accessibility (ARIA labels, keyboard navigation)
 * - ✅ RTL-first design (start/end logical properties)
 * - ✅ Mobile-friendly (50px min-height, proper touch targets)
 * - ✅ Responsive design (adapts button size and drawer position)
 *
 * Design Consistency:
 * - Button uses same styling as SearchInput and form fields
 * - Focus ring: ring-2 ring-primary/20
 * - Border: border-base-300, hover:border-primary
 * - Icon color: text-base-content/50
 * - Smooth transitions: transition-all duration-200
 * - Badge indicator when filters are active
 *
 * Usage Patterns:
 * 1. Simple filter (no active count):
 *    <FilterButton title="Filter Customers">
 *      <YourFilterContent />
 *    </FilterButton>
 *
 * 2. With active filter count:
 *    <FilterButton
 *      title="Filter Orders"
 *      activeCount={3}
 *      onApply={handleApplyFilters}
 *      onReset={handleResetFilters}
 *    >
 *      <YourFilterContent />
 *    </FilterButton>
 *
 * 3. Custom labels and button text:
 *    <FilterButton
 *      title="Advanced Filters"
 *      buttonText="Filters"
 *      applyLabel="Apply"
 *      resetLabel="Clear All"
 *    >
 *      <YourFilterContent />
 *    </FilterButton>
 *
 * @example
 * ```tsx
 * // In your page component
 * const [activeFilterCount, setActiveFilterCount] = useState(0);
 *
 * const handleApplyFilters = () => {
 *   // Apply filter logic
 *   setActiveFilterCount(3); // Update active count
 * };
 *
 * const handleResetFilters = () => {
 *   // Reset filter logic
 *   setActiveFilterCount(0); // Clear count
 * };
 *
 * <FilterButton
 *   title="Filter Customers"
 *   activeCount={activeFilterCount}
 *   onApply={handleApplyFilters}
 *   onReset={handleResetFilters}
 * >
 *   <div className="space-y-4">
 *     <FormSelect label="Status" options={statusOptions} />
 *     <FormInput label="Min Amount" type="number" />
 *   </div>
 * </FilterButton>
 * ```
 */

import { useState } from 'react'
import { Filter } from 'lucide-react'
import { BottomSheet } from '../molecules/BottomSheet'
import type { ReactNode } from 'react'

export interface FilterButtonProps {
  /** Content to display in the filter drawer */
  children: ReactNode
  /** Title displayed in the drawer header */
  title: string
  /** Text displayed on the trigger button */
  buttonText?: string
  /** Number of active filters (displays badge indicator) */
  activeCount?: number
  /** Callback when Apply button is clicked */
  onApply?: () => void
  /** Callback when Reset button is clicked */
  onReset?: () => void
  /** Label for Apply button */
  applyLabel?: string
  /** Label for Reset button */
  resetLabel?: string
  /** Whether the button is disabled */
  disabled?: boolean
  /** Additional CSS classes for the trigger button */
  className?: string
}

/**
 * FilterButton Component
 *
 * A unified filter component with button trigger and drawer/sheet.
 * Manages its own open/close state internally for simplicity.
 */
export function FilterButton({
  children,
  title,
  buttonText = 'Filter',
  activeCount = 0,
  onApply,
  onReset,
  applyLabel = 'Apply Filters',
  resetLabel = 'Reset',
  disabled = false,
  className = '',
}: FilterButtonProps) {
  const [isOpen, setIsOpen] = useState(false)

  const handleApply = () => {
    onApply?.()
    setIsOpen(false)
  }

  const handleReset = () => {
    onReset?.()
  }

  // Build footer with action buttons
  const footer =
    onApply !== undefined || onReset !== undefined ? (
      <div className="flex gap-2">
        {onReset && (
          <button
            type="button"
            onClick={handleReset}
            className="btn btn-ghost flex-1"
          >
            {resetLabel}
          </button>
        )}
        {onApply && (
          <button
            type="button"
            onClick={handleApply}
            className="btn btn-primary flex-1"
          >
            {applyLabel}
          </button>
        )}
      </div>
    ) : undefined

  return (
    <>
      {/* Filter Trigger Button */}
      <button
        type="button"
        onClick={() => {
          setIsOpen(true)
        }}
        disabled={disabled}
        className={`
          btn
          min-h-[50px]
          gap-2
          border-base-300
          bg-base-100
          hover:border-primary
          hover:bg-base-100
          transition-all duration-200
          focus:border-primary
          focus:ring-2
          focus:ring-primary/20
          focus:outline-none
          disabled:opacity-50
          disabled:cursor-not-allowed
          ${className}
        `}
        aria-label={`${buttonText}${activeCount > 0 ? ` (${String(activeCount)} active)` : ''}`}
      >
        {/* Filter Icon */}
        <span className="text-base-content/50" aria-hidden="true">
          <Filter size={18} />
        </span>

        {/* Button Text */}
        <span className="text-base text-base-content">{buttonText}</span>

        {/* Active Filter Count Badge */}
        {activeCount > 0 && (
          <span
            className="badge badge-primary badge-sm"
            aria-label={`${String(activeCount)} filters active`}
          >
            {activeCount}
          </span>
        )}
      </button>

      {/* Filter Drawer */}
      <BottomSheet
        isOpen={isOpen}
        onClose={() => {
          setIsOpen(false)
        }}
        title={title}
        footer={footer}
        side="end"
        size="md"
        contentClassName="space-y-6"
        ariaLabel="Filter options"
      >
        {children}
      </BottomSheet>
    </>
  )
}
