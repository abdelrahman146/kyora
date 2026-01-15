/**
 * AdditionalDetailsInputs Component
 *
 * A collapsible section for optional customer details (email, gender).
 * Uses the same progressive disclosure pattern as SocialMediaInputs.
 *
 * Features:
 * - Collapsible section (collapsed by default in create, expanded when data exists in edit)
 * - Shows count of filled fields in header
 * - Mobile-first responsive grid (1-2 columns)
 * - RTL/LTR support
 * - All fields optional
 *
 * @example
 * ```tsx
 * <AdditionalDetailsInputs
 *   email={emailField.value}
 *   onEmailChange={emailField.onChange}
 *   gender={genderField.value}
 *   onGenderChange={genderField.onChange}
 *   disabled={isSubmitting}
 * />
 * ```
 */

import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { ChevronDown, ChevronUp, Mail } from 'lucide-react'
import type { CustomerGender } from '@/api/customer'
import { cn } from '@/lib/utils'
import { FormInput } from '@/components/form/FormInput'
import { FormSelect } from '@/components/form/FormSelect'

export interface AdditionalDetailsInputsProps {
  // Email
  email?: string
  onEmailChange?: (value: string) => void
  emailError?: string

  // Gender
  gender?: CustomerGender | ''
  onGenderChange?: (value: CustomerGender | '') => void
  genderError?: string

  // State
  disabled?: boolean
  defaultExpanded?: boolean
}

export function AdditionalDetailsInputs({
  email = '',
  onEmailChange,
  emailError,
  gender = '',
  onGenderChange,
  genderError,
  disabled = false,
  defaultExpanded = false,
}: AdditionalDetailsInputsProps) {
  const { t: tCustomers } = useTranslation('customers')

  // Count filled fields for summary
  const filledCount = useMemo(() => {
    let count = 0
    if (email.trim()) count++
    if (gender) count++
    return count
  }, [email, gender])

  // Auto-expand if data exists, otherwise respect defaultExpanded
  const [isExpanded, setIsExpanded] = useState(
    defaultExpanded || filledCount > 0,
  )

  const genderOptions: Array<{ value: CustomerGender; label: string }> = [
    { value: 'male', label: tCustomers('form.gender_male') },
    { value: 'female', label: tCustomers('form.gender_female') },
    { value: 'other', label: tCustomers('form.gender_other') },
  ]

  return (
    <div className="space-y-3">
      {/* Header - Collapsible */}
      <button
        type="button"
        onClick={() => {
          setIsExpanded(!isExpanded)
        }}
        className={cn(
          'w-full flex items-center justify-between',
          'px-4 py-3 rounded-lg',
          'bg-base-200/50 hover:bg-base-200',
          'transition-colors duration-200',
          'focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2',
          disabled && 'opacity-60 cursor-not-allowed',
        )}
        disabled={disabled}
        aria-expanded={isExpanded}
        aria-controls="additional-details-inputs"
      >
        <div className="flex items-center gap-3">
          <span className="font-medium text-base-content">
            {tCustomers('form.additional_details_section')}
          </span>
          {filledCount > 0 && (
            <span className="badge badge-primary badge-sm">{filledCount}</span>
          )}
        </div>
        {isExpanded ? (
          <ChevronUp size={20} className="text-base-content/60" />
        ) : (
          <ChevronDown size={20} className="text-base-content/60" />
        )}
      </button>

      {/* Content - Expandable */}
      {isExpanded && (
        <div
          id="additional-details-inputs"
          className={cn(
            'grid gap-3',
            'grid-cols-1 sm:grid-cols-2',
            'animate-in fade-in slide-in-from-top-2 duration-200',
          )}
        >
          {/* Email */}
          <FormInput
            type="email"
            label={tCustomers('form.email')}
            placeholder={tCustomers('form.email_placeholder')}
            value={email}
            onChange={(e) => {
              onEmailChange?.(e.target.value)
            }}
            error={emailError}
            disabled={disabled}
            autoComplete="email"
            startIcon={<Mail className="w-4 h-4 text-base-content/60" />}
          />

          {/* Gender */}
          <FormSelect<CustomerGender | ''>
            label={tCustomers('form.gender')}
            options={genderOptions}
            value={gender}
            onChange={(value) => {
              const singleValue = Array.isArray(value) ? value[0] : value
              onGenderChange?.(singleValue)
            }}
            disabled={disabled}
            placeholder={tCustomers('form.select_gender')}
            error={genderError}
          />
        </div>
      )}
    </div>
  )
}
