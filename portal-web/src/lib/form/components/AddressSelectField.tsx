/**
 * AddressSelectField Component - Form Composition Layer
 *
 * Pre-bound address selector with customer's addresses.
 * Automatically wires to TanStack Form field context and handles address fetching.
 *
 * DESIGN GUIDELINES:
 * ==================
 * This component provides a searchable address selector displaying customer addresses,
 * with full address details for better identification.
 *
 * Key Features:
 * 1. Customer addresses - Loads all addresses for a specific customer
 * 2. Address details - Shows full address with phone for identification
 * 3. Mobile-first - Uses bottom sheet on mobile via FormSelect
 * 4. RTL support - Full right-to-left layout support
 * 5. Accessibility - WCAG AA compliant with ARIA attributes
 *
 * Usage within form:
 * ```tsx
 * <form.AppField name="shippingAddressId">
 *   {(field) => (
 *     <field.AddressSelectField
 *       label="Shipping Address"
 *       businessDescriptor="my-business"
 *       customerId="customer-id"
 *       placeholder="Select address..."
 *     />
 *   )}
 * </form.AppField>
 * ```
 */

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { useFieldContext } from '../contexts'
import type { FormSelectOption } from '@/components/atoms/FormSelect'
import { FormSelect } from '@/components/atoms/FormSelect'
import { useAddressesQuery } from '@/api/address'

export interface AddressSelectFieldProps {
  /** Field label */
  label?: string
  /** Placeholder text */
  placeholder?: string
  /** Business descriptor for fetching addresses */
  businessDescriptor: string
  /** Customer ID for fetching addresses */
  customerId: string
  /** Whether the field is disabled */
  disabled?: boolean
  /** Whether the field is required (visual indicator) */
  required?: boolean
  /** Additional CSS classes */
  className?: string
}

export function AddressSelectField(props: AddressSelectFieldProps) {
  const field = useFieldContext<string>()
  const { t } = useTranslation(['common', 'errors'])

  const {
    label = t('common:address'),
    placeholder = t('common:select_address'),
    businessDescriptor,
    customerId,
    disabled = false,
    required = false,
    className,
  } = props

  const { data: addresses = [], isLoading } = useAddressesQuery(
    businessDescriptor,
    customerId,
  )

  const addressOptions: Array<FormSelectOption<string>> = useMemo(() => {
    return addresses.map((address) => ({
      value: address.id,
      label: `${address.city}, ${address.state}`,
      description: [
        address.street,
        address.zipCode,
        `${address.phoneCode} ${address.phoneNumber}`,
      ]
        .filter(Boolean)
        .join(' • '),
    }))
  }, [addresses])

  const errorMessage = field.state.meta.errors[0]
  const translatedError =
    typeof errorMessage === 'string'
      ? t(`errors:${errorMessage}`, { defaultValue: errorMessage })
      : errorMessage?.message

  return (
    <div className={className}>
      <FormSelect
        label={label}
        placeholder={placeholder}
        options={addressOptions}
        value={field.state.value}
        onChange={(value) => field.handleChange(value as string)}
        error={translatedError}
        disabled={disabled || isLoading}
        required={required}
        searchable
        helperText={isLoading ? t('common:loading') : undefined}
      />
    </div>
  )
}
