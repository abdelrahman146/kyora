/**
 * SortButton Component
 *
 * A production-grade component that provides sorting functionality
 * for list views. Similar to FilterButton but specifically for sorting columns.
 *
 * Features:
 * - ✅ Consistent button styling (matches filter button and design system)
 * - ✅ Integrated bottom sheet (mobile-only component)
 * - ✅ Built-in state management (open/close)
 * - ✅ Radio-based sort field and direction selection
 * - ✅ Apply sorting callback
 * - ✅ Full accessibility (ARIA labels, keyboard navigation)
 * - ✅ RTL-first design (start/end logical properties)
 * - ✅ Mobile-friendly (50px min-height, proper touch targets)
 *
 * Usage:
 * ```tsx
 * <SortButton
 *   title="Sort Customers"
 *   sortOptions={[
 *     { value: 'name', label: 'Name' },
 *     { value: 'ordersCount', label: 'Orders Count' },
 *     { value: 'totalSpent', label: 'Total Spent' },
 *   ]}
 *   currentSortBy="name"
 *   currentSortOrder="asc"
 *   onApply={(sortBy, sortOrder) => {
 *     navigate({ search: { ...search, sortBy, sortOrder, page: 1 } })
 *   }}
 * />
 * ```
 */

import { useState } from 'react'
import { ArrowDownWideNarrow } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { BottomSheet } from '../molecules/BottomSheet'

export interface SortOption {
  value: string
  label: string
}

export interface SortButtonProps {
  /** Title displayed in the sort sheet header */
  title: string
  /** Text displayed on the trigger button */
  buttonText?: string
  /** Available sort field options */
  sortOptions: Array<SortOption>
  /** Current active sort field */
  currentSortBy?: string
  /** Current active sort order */
  currentSortOrder?: 'asc' | 'desc'
  /** Callback when Apply button is clicked */
  onApply: (sortBy: string, sortOrder: 'asc' | 'desc') => void
  /** Whether the button is disabled */
  disabled?: boolean
  /** Additional CSS classes for the trigger button */
  className?: string
}

/**
 * SortButton Component
 *
 * A mobile-focused sort component with button trigger and bottom sheet.
 * Manages its own open/close state and temporary sort selections.
 */
export function SortButton({
  title,
  buttonText,
  sortOptions,
  currentSortBy,
  currentSortOrder = 'asc',
  onApply,
  disabled = false,
  className = '',
}: SortButtonProps) {
  const { t } = useTranslation()
  const [isOpen, setIsOpen] = useState(false)

  // Internal state for temporary selections before applying
  const [tempSortBy, setTempSortBy] = useState(
    currentSortBy || sortOptions[0]?.value || '',
  )
  const [tempSortOrder, setTempSortOrder] = useState<'asc' | 'desc'>(
    currentSortOrder,
  )

  // Sync internal state when props change
  const handleOpen = () => {
    setTempSortBy(currentSortBy || sortOptions[0]?.value || '')
    setTempSortOrder(currentSortOrder)
    setIsOpen(true)
  }

  const handleApply = () => {
    onApply(tempSortBy, tempSortOrder)
    setIsOpen(false)
  }

  const defaultButtonText = buttonText || t('common:sort')
  const isActive = !!currentSortBy

  const footer = (
    <div className="flex gap-2">
      <button
        type="button"
        onClick={() => {
          setIsOpen(false)
        }}
        className="btn btn-ghost flex-1"
      >
        {t('common.cancel')}
      </button>
      <button
        type="button"
        onClick={handleApply}
        className="btn btn-primary flex-1"
      >
        {t('common.apply')}
      </button>
    </div>
  )

  return (
    <>
      {/* Sort Trigger Button */}
      <button
        type="button"
        onClick={handleOpen}
        disabled={disabled}
        className={`
          btn
          min-h-[50px]
          gap-2
          w-full
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
          ${isActive ? 'border-primary' : ''}
          ${className}
        `}
        aria-label={defaultButtonText}
      >
        {/* Sort Icon */}
        <span
          className={isActive ? 'text-primary' : 'text-base-content/50'}
          aria-hidden="true"
        >
          <ArrowDownWideNarrow size={18} />
        </span>

        {/* Button Text */}
        <span className="text-base text-base-content">{defaultButtonText}</span>
      </button>

      {/* Sort Bottom Sheet */}
      <BottomSheet
        isOpen={isOpen}
        onClose={() => {
          setIsOpen(false)
        }}
        title={title}
        footer={footer}
      >
        <div className="flex flex-row gap-4 p-4">
          {/* Sort Field Selection */}
          <div className="form-control flex-1">
            <label className="label">
              <span className="label-text font-semibold">
                {t('common:sort_by')}
              </span>
            </label>
            <div className="space-y-2">
              {sortOptions.map((option) => (
                <label
                  key={option.value}
                  className="flex items-center gap-3 p-3 rounded-lg border border-base-300 cursor-pointer hover:border-primary transition-colors"
                >
                  <input
                    type="radio"
                    name="sortBy"
                    value={option.value}
                    checked={tempSortBy === option.value}
                    onChange={(e) => {
                      setTempSortBy(e.target.value)
                    }}
                    className="radio radio-primary"
                  />
                  <span className="flex-1">{option.label}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Sort Order Selection */}
          <div className="form-control flex-1">
            <label className="label">
              <span className="label-text font-semibold">
                {t('common:sort_order')}
              </span>
            </label>
            <div className="space-y-2">
              <label className="flex items-center gap-3 p-3 rounded-lg border border-base-300 cursor-pointer hover:border-primary transition-colors">
                <input
                  type="radio"
                  name="sortOrder"
                  value="asc"
                  checked={tempSortOrder === 'asc'}
                  onChange={() => {
                    setTempSortOrder('asc')
                  }}
                  className="radio radio-primary"
                />
                <span className="flex-1">{t('common:ascending')}</span>
              </label>
              <label className="flex items-center gap-3 p-3 rounded-lg border border-base-300 cursor-pointer hover:border-primary transition-colors">
                <input
                  type="radio"
                  name="sortOrder"
                  value="desc"
                  checked={tempSortOrder === 'desc'}
                  onChange={() => {
                    setTempSortOrder('desc')
                  }}
                  className="radio radio-primary"
                />
                <span className="flex-1">{t('common:descending')}</span>
              </label>
            </div>
          </div>
        </div>
      </BottomSheet>
    </>
  )
}
