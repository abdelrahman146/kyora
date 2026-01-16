import { forwardRef, useCallback } from 'react'
import type { ChangeEvent, FocusEvent, KeyboardEvent } from 'react'
import type { FormInputProps } from '@/components/form/FormInput'
import { FormInput } from '@/components/form/FormInput'
import { cn } from '@/lib/utils'

export interface PriceInputProps extends Omit<
  FormInputProps,
  'type' | 'inputMode' | 'pattern' | 'onBlur'
> {
  currencyCode?: string
  /**
   * Optional callback for when the input loses focus,
   * allowing you to sanitize the final value in the parent form
   */
  onBlur?: (event: FocusEvent<HTMLInputElement>) => void
}

export const PriceInput = forwardRef<HTMLInputElement, PriceInputProps>(
  (
    { currencyCode = 'AED', className, placeholder = '0.00', onBlur, ...props },
    ref,
  ) => {
    // 1. Prevent invalid keys (letters, symbols) before they register
    const handleKeyDown = useCallback((e: KeyboardEvent<HTMLInputElement>) => {
      const validKeys = [
        'Backspace',
        'Tab',
        'Enter',
        'Delete',
        'ArrowLeft',
        'ArrowRight',
        'Home',
        'End',
        '.',
        ',',
      ]

      // Allow standard control keys
      if (
        validKeys.includes(e.key) ||
        // Allow select all/copy/paste (Ctrl/Cmd + A/C/V)
        e.metaKey ||
        e.ctrlKey
      ) {
        return
      }

      // Block non-numeric keys
      if (!/^[0-9]$/.test(e.key)) {
        e.preventDefault()
      }
    }, [])

    // 2. Handle the value change
    const handleChange = useCallback(
      (event: ChangeEvent<HTMLInputElement>) => {
        let value = event.target.value

        // Standardize comma to dot
        value = value.replace(/,/g, '.')

        // Regex Explanation:
        // ^\d* -> Starts with zero or more digits
        // (\.\d{0,2})? -> Optionally followed by a dot and 0-2 digits
        // $          -> End of string
        const isValidMoneyFormat = /^\d*(\.\d{0,2})?$/.test(value)

        // If the value isn't valid (and isn't empty), ignore the change
        if (!isValidMoneyFormat && value !== '') return

        // Prevent multiple leading zeros (e.g., 005 -> 5)
        if (value.length > 1 && value.startsWith('0') && value[1] !== '.') {
          value = value.replace(/^0+/, '')
        }

        // Pass the sanitized value to the parent
        // We create a new event to respect the FormInputProps interface
        const syntheticEvent = {
          ...event,
          target: {
            ...event.target,
            value,
          },
        }

        props.onChange?.(syntheticEvent)
      },
      [props],
    )

    // 3. Format strictly on blur (e.g. "10" -> "10.00")
    const handleBlur = useCallback(
      (event: FocusEvent<HTMLInputElement>) => {
        let value = event.target.value

        if (value) {
          // Standardize comma to dot again just in case
          value = value.replace(/,/g, '.')

          // Parse float to remove dangling decimals like "12."
          const number = parseFloat(value)

          if (!isNaN(number)) {
            value = number.toFixed(2) // Force 2 decimal places
          } else {
            value = ''
          }
        }

        // Update the input value visually
        event.target.value = value

        // Trigger parent onChange to save the formatted "10.00"
        const syntheticEvent = {
          ...event,
          target: { ...event.target, value },
        } as ChangeEvent<HTMLInputElement>

        props.onChange?.(syntheticEvent)

        // Call original onBlur if provided
        onBlur?.(event)
      },
      [onBlur, props],
    )

    return (
      <FormInput
        ref={ref}
        type="text" // Keep text to handle custom formatting better than "number"
        inputMode="decimal" // triggers numeric keyboard on mobile
        dir="ltr"
        placeholder={placeholder}
        startIcon={
          <span
            className={cn(
              'uppercase text-sm font-medium text-base-content/80',
              'me-3 px-1 select-none', // select-none prevents highlighting currency
            )}
          >
            {currencyCode}
          </span>
        }
        className={cn('ps-12 font-mono', className)} // font-mono looks better for numbers
        {...props}
        onKeyDown={handleKeyDown}
        onChange={handleChange}
        onBlur={handleBlur}
        // Disable scroll wheel changing numbers (common annoyance)
        onWheel={(e) => e.currentTarget.blur()}
        autoComplete="off"
      />
    )
  },
)

PriceInput.displayName = 'PriceInput'
