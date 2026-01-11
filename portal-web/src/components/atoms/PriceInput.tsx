import { forwardRef, useCallback } from 'react'

import type { ChangeEvent } from 'react'
import type { FormInputProps } from '@/components/form/FormInput'
import { FormInput } from '@/components/form/FormInput'
import { cn } from '@/lib/utils'

export interface PriceInputProps extends Omit<
  FormInputProps,
  'type' | 'inputMode' | 'pattern'
> {
  currencyCode?: string
}

export const PriceInput = forwardRef<HTMLInputElement, PriceInputProps>(
  (
    { currencyCode = 'AED', className, placeholder = '0.00', ...props },
    ref,
  ) => {
    const handleChange = useCallback(
      (event: ChangeEvent<HTMLInputElement>) => {
        const raw = event.target.value.replace(',', '.')
        const normalized = raw.replace(/[^0-9.]/g, '')

        if (normalized === '') {
          props.onChange?.(event)
          return
        }

        const hasDecimal = normalized.includes('.')
        const [integerPartRaw = '', ...rest] = normalized.split('.')
        const decimalRaw = rest.join('')
        const decimalPart = decimalRaw.slice(0, 2)

        // Strip leading zeros but keep a single zero when the value is 0
        const integerPart = integerPartRaw.replace(/^0+(?=\d)/, '') || '0'

        const sanitized = hasDecimal
          ? `${integerPart}.${decimalPart}`
          : integerPart

        const syntheticEvent = {
          ...event,
          target: {
            ...event.target,
            value: sanitized,
          },
        }

        props.onChange?.(syntheticEvent as React.ChangeEvent<HTMLInputElement>)
      },
      [props],
    )

    return (
      <FormInput
        ref={ref}
        type="text"
        inputMode="decimal"
        pattern="^[0-9]*([.,][0-9]{0,2})?$"
        dir="ltr"
        placeholder={placeholder}
        startIcon={
          <span
            className={cn(
              'uppercase text-sm font-medium text-base-content/80',
              'me-3 px-1',
            )}
          >
            {currencyCode}
          </span>
        }
        className={cn('ps-12', className)}
        {...props}
        onChange={handleChange}
      />
    )
  },
)

PriceInput.displayName = 'PriceInput'
