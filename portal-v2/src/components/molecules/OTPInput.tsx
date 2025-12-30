import { useEffect, useRef } from 'react'
import type { ChangeEvent, ClipboardEvent, KeyboardEvent } from 'react'

import { cn } from '@/lib/utils'

/**
 * OTPInput Props
 */
export interface OTPInputProps {
  /**
   * Number of OTP digits (typically 4 or 6)
   * @default 6
   */
  length?: number

  /**
   * Current OTP value as array of digits
   */
  value: Array<string>

  /**
   * Callback when OTP value changes
   */
  onChange: (value: Array<string>) => void

  /**
   * Callback when OTP is complete (all digits filled)
   */
  onComplete?: (code: string) => void

  /**
   * Whether the inputs are disabled
   */
  disabled?: boolean

  /**
   * Whether to show error state
   */
  error?: boolean

  /**
   * Whether to auto-focus first input on mount
   */
  autoFocus?: boolean

  /**
   * Custom class name for container
   */
  className?: string

  /**
   * Custom class name for individual inputs
   */
  inputClassName?: string
}

/**
 * OTPInput Component
 *
 * A production-ready OTP (One-Time Password) input component with:
 * - Auto-focus next/previous inputs
 * - Paste support (auto-fills all digits)
 * - Keyboard navigation (Backspace, Arrow keys)
 * - Mobile-friendly (numeric keyboard)
 * - Accessible (proper ARIA labels)
 * - Flexible styling
 * - RTL support
 *
 * @example
 * ```tsx
 * const [otp, setOtp] = useState(['', '', '', '', '', ''])
 * 
 * <OTPInput
 *   value={otp}
 *   onChange={setOtp}
 *   onComplete={(code) => verifyOtp(code)}
 *   autoFocus
 * />
 * ```
 */
export function OTPInput({
  length = 6,
  value,
  onChange,
  onComplete,
  disabled = false,
  error = false,
  autoFocus = false,
  className,
  inputClassName,
}: OTPInputProps) {
  const inputRefs = useRef<Array<HTMLInputElement | null>>([])

  // Initialize refs array
  useEffect(() => {
    inputRefs.current = inputRefs.current.slice(0, length)
  }, [length])

  // Auto-focus first input on mount
  useEffect(() => {
    if (autoFocus && !disabled) {
      inputRefs.current[0]?.focus()
    }
  }, [autoFocus, disabled])

  // Check if OTP is complete and trigger onComplete
  useEffect(() => {
    const isComplete = value.every((digit) => digit !== '')
    if (isComplete && onComplete && value.length === length) {
      onComplete(value.join(''))
    }
  }, [value, onComplete, length])

  /**
   * Handle input change
   */
  const handleChange = (index: number, e: ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value

    // Only allow digits
    if (inputValue && !/^\d$/.test(inputValue)) {
      return
    }

    const newValue = [...value]
    newValue[index] = inputValue.slice(-1) // Take only last digit
    onChange(newValue)

    // Auto-focus next input if value was entered
    if (inputValue && index < length - 1) {
      inputRefs.current[index + 1]?.focus()
    }
  }

  /**
   * Handle paste event - auto-fill all digits
   */
  const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault()

    const pastedData = e.clipboardData.getData('text').trim()

    // Only accept digits with correct length
    const digitRegex = new RegExp(`^\\d{${length}}$`)
    if (!digitRegex.test(pastedData)) {
      return
    }

    const newValue = pastedData.split('').slice(0, length)
    onChange(newValue)

    // Focus last input after paste
    inputRefs.current[length - 1]?.focus()
  }

  /**
   * Handle keyboard navigation
   */
  const handleKeyDown = (index: number, e: KeyboardEvent<HTMLInputElement>) => {
    const currentValue = value[index]

    // Backspace: clear current or move to previous
    if (e.key === 'Backspace') {
      e.preventDefault()

      if (currentValue) {
        // Clear current digit
        const newValue = [...value]
        newValue[index] = ''
        onChange(newValue)
      } else if (index > 0) {
        // Move to previous input and clear it
        const newValue = [...value]
        newValue[index - 1] = ''
        onChange(newValue)
        inputRefs.current[index - 1]?.focus()
      }

      return
    }

    // Delete: clear current digit
    if (e.key === 'Delete') {
      e.preventDefault()
      const newValue = [...value]
      newValue[index] = ''
      onChange(newValue)
      return
    }

    // Arrow Left: move to previous input
    if (e.key === 'ArrowLeft' && index > 0) {
      e.preventDefault()
      inputRefs.current[index - 1]?.focus()
      return
    }

    // Arrow Right: move to next input
    if (e.key === 'ArrowRight' && index < length - 1) {
      e.preventDefault()
      inputRefs.current[index + 1]?.focus()
      return
    }

    // Home: focus first input
    if (e.key === 'Home') {
      e.preventDefault()
      inputRefs.current[0]?.focus()
      return
    }

    // End: focus last input
    if (e.key === 'End') {
      e.preventDefault()
      inputRefs.current[length - 1]?.focus()
      return
    }
  }

  /**
   * Handle focus - select all text for easy replacement
   */
  const handleFocus = (e: React.FocusEvent<HTMLInputElement>) => {
    e.target.select()
  }

  return (
    <div
      className={cn('flex justify-center gap-2', className)}
      role="group"
      aria-label="One-time password input"
    >
      {Array.from({ length }, (_, index) => (
        <input
          key={index}
          ref={(el) => {
            inputRefs.current[index] = el
          }}
          type="text"
          inputMode="numeric"
          pattern="\d*"
          maxLength={1}
          value={value[index] || ''}
          onChange={(e) => {
            handleChange(index, e)
          }}
          onKeyDown={(e) => {
            handleKeyDown(index, e)
          }}
          onPaste={handlePaste}
          onFocus={handleFocus}
          disabled={disabled}
          aria-label={`Digit ${index + 1} of ${length}`}
          className={cn(
            'input input-bordered w-12 h-14 text-center text-xl font-bold',
            'transition-all duration-200',
            'focus:ring-2 focus:ring-primary focus:ring-offset-2',
            error && 'input-error',
            disabled && 'opacity-50 cursor-not-allowed',
            inputClassName,
          )}
          autoComplete="one-time-code"
        />
      ))}
    </div>
  )
}
