/**
 * CustomerSelectField Component - Form Composition Layer
 *
 * Pre-bound customer selector with autocomplete search functionality.
 * Automatically wires to TanStack Form field context and handles customer fetching.
 *
 * DESIGN GUIDELINES:
 * ==================
 * This component provides a searchable customer selector with autocomplete,
 * displaying customer details (phone/email) for better identification.
 *
 * Key Features:
 * 1. Autocomplete search - Real-time search with debouncing
 * 2. Customer details - Shows phone/email below selection
 * 3. Persistent state - Loads selected customer on mount (survives refresh)
 * 4. Mobile-first - Uses bottom sheet on mobile via FormSelect
 * 5. RTL support - Full right-to-left layout support
 * 6. Accessibility - WCAG AA compliant with ARIA attributes
 *
 * Usage within form:
 * ```tsx
 * <form.AppField name="customerId">
 *   {(field) => (
 *     <field.CustomerSelectField
 *       label="Select Customer"
 *       businessDescriptor="my-business"
 *       placeholder="Search by name, phone, or email..."
 *     />
 *   )}
 * </form.AppField>
 * ```
 */

import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { useFieldContext } from '../contexts'
import type { CustomerSelectFieldProps } from '../types'
import type { Customer } from '@/api/customer'
import type { FormSelectOption } from '@/components/form/FormSelect'
import { useCustomerQuery, useCustomersQuery } from '@/api/customer'
import { FormSelect } from '@/components/form/FormSelect'

export function CustomerSelectField(props: CustomerSelectFieldProps) {
  const field = useFieldContext<string>()
  const { t: tErrors } = useTranslation('errors')
  const { t: tCommon } = useTranslation('errors')

  const [customerSearchQuery, setCustomerSearchQuery] = useState('')
  const [debouncedCustomerSearch, setDebouncedCustomerSearch] = useState('')
  const [isOpen, setIsOpen] = useState(false)

  // Debounce customer search query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedCustomerSearch(customerSearchQuery)
    }, 300)
    return () => clearTimeout(timer)
  }, [customerSearchQuery])

  const listParams = useMemo(
    () => ({
      search: debouncedCustomerSearch || undefined,
      page: 1,
      pageSize: 10,
    }),
    [debouncedCustomerSearch],
  )

  const customersQuery = useCustomersQuery(props.businessDescriptor, listParams)
  const customersData = customersQuery.data

  const selectedCustomerQuery = useCustomerQuery(
    props.businessDescriptor,
    field.state.value,
  )

  const selectedCustomer = useMemo<Customer | null>(() => {
    const fromList = customersData?.items.find(
      (c) => c.id === field.state.value,
    )
    return fromList ?? selectedCustomerQuery.data ?? null
  }, [customersData?.items, field.state.value, selectedCustomerQuery.data])

  useEffect(() => {
    if (isOpen) return
    if (!field.state.value) return
    if (!selectedCustomer) return
    setCustomerSearchQuery(selectedCustomer.name)
  }, [field.state.value, isOpen, selectedCustomer])

  // Transform customers to select options
  const customerOptions = useMemo<Array<FormSelectOption<string>>>(() => {
    if (!customersData?.items) return []
    return customersData.items.map((customer) => ({
      value: customer.id,
      label: customer.name,
      description: customer.phoneNumber
        ? `${customer.phoneCode} ${customer.phoneNumber}`
        : customer.email || undefined,
    }))
  }, [customersData])

  // Extract error from field state and translate
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return tErrors(firstError)
    }

    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return tErrors(errorObj.message)
    }

    return undefined
  }, [field.state.meta.errors, tErrors])

  // Show error only after field has been touched
  const showError = field.state.meta.isTouched && error

  const handleCustomerSelect = (customerId: string | Array<string>) => {
    const id = Array.isArray(customerId) ? customerId[0] : customerId
    field.handleChange(id)
    const customer = customersData?.items.find((c) => c.id === id)
    if (customer) {
      setCustomerSearchQuery(customer.name)
    }
  }

  const handleClear = () => {
    field.handleChange('')
    setCustomerSearchQuery('')
  }

  return (
    <div className="form-control w-full">
      {/* Label - Matches TextField pattern */}
      {props.label && (
        <label className="label">
          <span className="label-text text-base-content/70 font-medium">
            {props.label}
            {props.required && <span className="text-error ms-1">*</span>}
          </span>
        </label>
      )}

      {/* FormSelect with customer search */}
      <FormSelect
        options={customerOptions}
        value={field.state.value}
        onChange={handleCustomerSelect}
        searchable
        searchValue={customerSearchQuery}
        onSearchChange={setCustomerSearchQuery}
        placeholder={props.placeholder || tCommon('select')}
        clearable
        onClear={handleClear}
        disabled={props.disabled || field.state.meta.isValidating}
        error={showError ? error : undefined}
        onOpen={() => {
          setIsOpen(true)
          setCustomerSearchQuery('')
        }}
        onClose={() => {
          setIsOpen(false)
          if (selectedCustomer) {
            setCustomerSearchQuery(selectedCustomer.name)
          }
        }}
      />

      {/* Selected customer details */}
      {selectedCustomer && !showError && (
        <label className="label pt-1">
          <span className="label-text-alt text-base-content/60">
            {selectedCustomer.phoneNumber
              ? `${selectedCustomer.phoneCode} ${selectedCustomer.phoneNumber}`
              : selectedCustomer.email || ''}
          </span>
        </label>
      )}

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
