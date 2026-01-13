/**
 * QuantityInput - Form-Agnostic Quantity Selector
 *
 * A reusable quantity input component with increment/decrement buttons.
 * Follows the same UI/UX patterns as FormInput for consistency.
 *
 * Features:
 * - Increment/decrement buttons with configurable step
 * - Min/max validation (enforced on input)
 * - Consistent height and styling with FormInput (50px)
 * - RTL-compatible layout
 * - Accessible ARIA labels
 * - Keyboard input support
 *
 * Usage:
 * ```tsx
 * <QuantityInput
 *   value={quantity}
 *   onChange={setQuantity}
 *   min={1}
 *   max={100}
 *   incrementLabel="Increase"
 *   decrementLabel="Decrease"
 * />
 * ```
 */

import { Minus, Plus } from 'lucide-react'
import { cn } from '@/lib/utils'

export interface QuantityInputProps {
  /** Current value */
  value: number
  /** Change handler */
  onChange: (value: number) => void
  /** Blur handler */
  onBlur?: () => void
  /** Minimum allowed value (default: 0) */
  min?: number
  /** Maximum allowed value (default: Infinity) */
  max?: number
  /** Step increment/decrement value (default: 1) */
  step?: number
  /** Whether the input is disabled */
  disabled?: boolean
  /** ARIA label for increment button */
  incrementLabel: string
  /** ARIA label for decrement button */
  decrementLabel: string
  /** Show error state */
  error?: unknown
  /** HTML name attribute */
  name?: string
  /** HTML id attribute */
  id?: string
  /** Additional CSS classes */
  className?: string
  /** Required field */
  required?: boolean
}

export function QuantityInput({
  value,
  onChange,
  onBlur,
  min = 0,
  max = Infinity,
  step = 1,
  disabled = false,
  incrementLabel,
  decrementLabel,
  error,
  name,
  id,
  className,
  required,
}: QuantityInputProps) {
  const hasError = Boolean(error)

  const handleIncrement = () => {
    if (disabled || value >= max) return
    const newValue = Math.min(value + step, max)
    onChange(newValue)
  }

  const handleDecrement = () => {
    if (disabled || value <= min) return
    const newValue = Math.max(value - step, min)
    onChange(newValue)
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value.trim()

    // Allow empty input for better UX - will be clamped on blur
    if (inputValue === '') {
      onChange(min)
      return
    }

    const numValue = parseInt(inputValue, 10)
    if (!isNaN(numValue)) {
      // Clamp value between min and max
      const clampedValue = Math.max(min, Math.min(numValue, max))
      onChange(clampedValue)
    }
  }

  const handleBlur = () => {
    // Ensure value is within bounds on blur
    if (value < min) {
      onChange(min)
    } else if (value > max) {
      onChange(max)
    }
    onBlur?.()
  }

  const canDecrement = value > min
  const canIncrement = value < max

  return (
    <div className={cn('flex items-center gap-2', className)}>
      {/* Decrement Button */}
      <button
        type="button"
        className={cn(
          'btn btn-circle',
          'h-[44px] w-[44px] min-h-[44px]',
          !canDecrement || disabled ? 'btn-disabled' : 'btn-primary',
        )}
        onClick={handleDecrement}
        disabled={!canDecrement || disabled}
        aria-label={decrementLabel}
      >
        <Minus size={16} />
      </button>

      {/* Value Input - Matches FormInput height (50px) */}
      <input
        id={id}
        name={name}
        type="text"
        inputMode="numeric"
        pattern="[0-9]*"
        value={value}
        onChange={handleInputChange}
        onBlur={handleBlur}
        disabled={disabled}
        required={required}
        className={cn(
          'input input-bordered w-full text-center font-medium',
          'h-[50px] text-base',
          'transition-all duration-200',
          'focus:outline-none focus:border-primary focus:ring-2 focus:ring-primary/20',
          hasError &&
            'input-error border-error focus:border-error focus:ring-error/20',
          disabled && 'opacity-60 cursor-not-allowed',
        )}
        dir="ltr"
        aria-invalid={hasError}
        aria-required={required}
      />

      {/* Increment Button */}
      <button
        type="button"
        className={cn(
          'btn btn-circle',
          'h-[44px] w-[44px] min-h-[44px]',
          !canIncrement || disabled ? 'btn-disabled' : 'btn-primary',
        )}
        onClick={handleIncrement}
        disabled={!canIncrement || disabled}
        aria-label={incrementLabel}
      >
        <Plus size={16} />
      </button>
    </div>
  )
}
