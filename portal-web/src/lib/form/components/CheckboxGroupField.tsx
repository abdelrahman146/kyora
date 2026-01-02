/**
 * CheckboxGroupField Component - Form Composition Layer
 *
 * Pre-bound checkbox group that automatically wires to TanStack Form field context.
 * Designed for selecting multiple values from a list of options.
 *
 * DESIGN GUIDELINES:
 * ==================
 * This component uses a grid layout with visual feedback for better UX,
 * especially for options like social media platforms.
 *
 * Key Design Principles:
 * 1. Grid layout - Shows options in a responsive grid (2 columns on mobile, 3 on larger screens)
 * 2. Visual icons - Uses emoji or icons to make options scannable
 * 3. Bordered cards - Bordered containers that highlight when selected
 * 4. Consistent colors - Uses design system colors (primary for selection)
 * 5. Proper spacing - Matches the spacing used in other form fields
 *
 * Visual Feedback:
 * - Unselected: Gray border with hover effect
 * - Selected: Primary border with subtle primary background
 * - Always: Smooth transitions for state changes
 *
 * Usage within form:
 * ```tsx
 * <form.Field name="socialPlatforms">
 *   {(field) => (
 *     <field.CheckboxGroupField
 *       label="Social Media Platforms"
 *       options={[
 *         { value: 'instagram', label: 'Instagram', icon: '📷' },
 *         { value: 'tiktok', label: 'TikTok', icon: '🎵' },
 *       ]}
 *     />
 *   )}
 * </form.Field>
 * ```
 */

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useFieldContext } from '../contexts'
import type { CheckboxGroupFieldProps } from '../types'
import type { SocialPlatform } from '@/components/icons/social'
import { SocialIcon } from '@/components/icons/social'

export function CheckboxGroupField<T extends string = string>(
  props: CheckboxGroupFieldProps<T>,
) {
  const field = useFieldContext<Array<T>>()
  const { t } = useTranslation('errors')

  // Extract error from field state and translate
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return t(firstError)
    }

    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message)
    }

    return undefined
  }, [field.state.meta.errors, t])

  // Show error only after field has been touched (better UX)
  const showError = field.state.meta.isTouched && error

  return (
    <div className="form-control w-full">
      {/* Label and Description */}
      {(props.label || props.description) && (
        <label className="label pb-2">
          <div className="flex flex-col gap-1">
            {props.label && (
              <span className="label-text font-medium">{props.label}</span>
            )}
            {props.description && (
              <span className="label-text-alt text-base-content/60">
                {props.description}
              </span>
            )}
          </div>
        </label>
      )}

      {/* Checkbox Options in Grid Layout with Borders */}
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
        {props.options.map((option) => {
          const isChecked = field.state.value.includes(option.value)
          const optionId = `${field.name}-${option.value}`
          const isSocialPlatform = [
            'instagram',
            'facebook',
            'tiktok',
            'snapchat',
            'x',
            'whatsapp',
          ].includes(option.value)

          return (
            <label
              key={option.value}
              htmlFor={optionId}
              className={`
                flex items-center gap-2 p-3 rounded-lg border cursor-pointer transition-all
                ${
                  isChecked
                    ? 'border-primary bg-primary/10 ring-2 ring-primary/20'
                    : 'border-base-300 hover:border-base-content/20'
                }
                ${props.disabled || field.state.meta.isValidating ? 'opacity-60 cursor-not-allowed' : ''}
              `}
            >
              <input
                type="checkbox"
                id={optionId}
                className="checkbox checkbox-primary checkbox-sm"
                checked={isChecked}
                disabled={props.disabled || field.state.meta.isValidating}
                onChange={(e) => {
                  const newValue = e.target.checked
                    ? [...field.state.value, option.value]
                    : field.state.value.filter((v) => v !== option.value)
                  field.handleChange(newValue)
                }}
                onBlur={field.handleBlur}
              />
              {isSocialPlatform ? (
                <SocialIcon
                  platform={option.value as SocialPlatform}
                  className="w-5 h-5 text-base-content"
                  aria-hidden="true"
                />
              ) : (
                option.icon && (
                  <span className="text-lg" role="img" aria-hidden="true">
                    {option.icon}
                  </span>
                )
              )}
              <span className="text-sm font-medium flex-1">{option.label}</span>
            </label>
          )
        })}
      </div>

      {/* Error Message */}
      {showError && (
        <label className="label pt-1">
          <span className="label-text-alt text-error" role="alert">
            {error}
          </span>
        </label>
      )}
    </div>
  )
}
