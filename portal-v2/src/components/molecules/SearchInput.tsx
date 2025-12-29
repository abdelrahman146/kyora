/**
 * SearchInput Component
 *
 * A production-grade debounced search input with clear functionality.
 *
 * Features:
 * - ✅ Automatic debouncing (configurable, default 300ms)
 * - ✅ Clear button with loading state indicator
 * - ✅ Search icon (start position, RTL-compatible)
 * - ✅ Consistent styling with form fields (focus states, colors, transitions)
 * - ✅ Full accessibility (ARIA labels, keyboard navigation)
 * - ✅ RTL-first design (start/end logical properties)
 * - ✅ Mobile-friendly (50px min-height, proper touch targets)
 * - ✅ Disabled state support
 *
 * Consistency Notes:
 * - Uses same color scheme as FormInput/Input components
 * - Focus ring: ring-2 ring-primary/20
 * - Border: border-base-300, focus:border-primary
 * - Icon color: text-base-content/50
 * - Placeholder: text-base-content/40
 * - Smooth transitions: transition-all duration-200
 *
 * @example
 * ```tsx
 * const [search, setSearch] = useState("");
 *
 * <SearchInput
 *   value={search}
 *   onChange={setSearch}
 *   placeholder="Search customers..."
 *   debounceMs={300}
 * />
 * ```
 */

import { useEffect, useRef, useState } from 'react'
import { Search, X } from 'lucide-react'

export interface SearchInputProps {
  /** Current search value (controlled) */
  value: string
  /** Callback fired after debounce when value changes */
  onChange: (value: string) => void
  /** Placeholder text */
  placeholder?: string
  /** Debounce delay in milliseconds */
  debounceMs?: number
  /** Whether the input is disabled */
  disabled?: boolean
  /** Additional CSS classes */
  className?: string
}

export function SearchInput({
  value,
  onChange,
  placeholder = 'Search...',
  debounceMs = 300,
  disabled = false,
  className = '',
}: SearchInputProps) {
  const [localValue, setLocalValue] = useState(value)
  const [isDebouncePending, setIsDebouncePending] = useState(false)
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)

  // Sync external value changes
  useEffect(() => {
    setLocalValue(value)
  }, [value])

  // Debounced onChange
  useEffect(() => {
    if (localValue === value) {
      setIsDebouncePending(false)
      return
    }

    setIsDebouncePending(true)

    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }

    timeoutRef.current = setTimeout(() => {
      onChange(localValue)
      setIsDebouncePending(false)
    }, debounceMs)

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [localValue, value, onChange, debounceMs])

  const handleClear = () => {
    setLocalValue('')
    onChange('')
  }

  return (
    <div className={`relative ${className}`}>
      {/* Search Icon - Consistent with form fields */}
      <span
        className="absolute top-1/2 -translate-y-1/2 start-3 text-base-content/50 pointer-events-none z-10"
        aria-hidden="true"
      >
        <Search size={20} />
      </span>

      {/* Input - Consistent styling with FormInput/Input */}
      <input
        type="text"
        value={localValue}
        onChange={(e) => {
          setLocalValue(e.target.value)
        }}
        placeholder={placeholder}
        disabled={disabled}
        className="
          input input-bordered w-full
          min-h-[50px]
          ps-10 pe-10
          text-base text-base-content
          placeholder:text-base-content/40
          border-base-300
          bg-base-100
          transition-all duration-200
          focus:border-primary
          focus:ring-2
          focus:ring-primary/20
          focus:outline-none
          disabled:opacity-50
          disabled:cursor-not-allowed
          z-0
        "
        aria-label={placeholder}
      />

      {/* Clear Button / Loading Indicator */}
      {localValue && (
        <button
          type="button"
          onClick={handleClear}
          disabled={disabled}
          className="
            btn btn-ghost btn-sm btn-circle
            absolute top-1/2 -translate-y-1/2 end-2
            z-10
            hover:bg-base-300
            transition-colors
            disabled:opacity-50
            disabled:cursor-not-allowed
          "
          aria-label="Clear search"
        >
          {isDebouncePending ? (
            <span className="loading loading-spinner loading-xs"></span>
          ) : (
            <X size={16} />
          )}
        </button>
      )}
    </div>
  )
}
