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
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { useFieldContext } from '../contexts'
import type { CustomerSelectFieldProps } from '../types'
import type { Customer } from '@/api/customer'
import type { FormSelectOption } from '@/components/atoms/FormSelect'
import { customerApi } from '@/api/customer'
import { FormSelect } from '@/components/atoms/FormSelect'

export function CustomerSelectField(props: CustomerSelectFieldProps) {
  const field = useFieldContext<string>()
  const { t } = useTranslation('errors')

  const [customerSearchQuery, setCustomerSearchQuery] = useState('')
  const [debouncedCustomerSearch, setDebouncedCustomerSearch] = useState('')
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(
    null,
  )

  // Debounce customer search query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedCustomerSearch(customerSearchQuery)
    }, 300)
    return () => clearTimeout(timer)
  }, [customerSearchQuery])

  // Fetch customers for autocomplete
  const { data: customersData } = useQuery({
    queryKey: [
      'customers',
      'search',
      props.businessDescriptor,
      debouncedCustomerSearch,
    ],
    queryFn: async () => {
      return customerApi.listCustomers(props.businessDescriptor, {
        search: debouncedCustomerSearch || undefined,
        pageSize: 10,
      })
    },
    enabled: true, // Always fetch customers (even with empty search)
  })

  // Load selected customer by ID when field value exists
  useEffect(() => {
    if (field.state.value && !selectedCustomer) {
      customerApi
        .getCustomer(props.businessDescriptor, field.state.value)
        .then((customer) => {
          setSelectedCustomer(customer)
          setCustomerSearchQuery(customer.name)
        })
        .catch(() => {
          setSelectedCustomer(null)
        })
    }
  }, [field.state.value, props.businessDescriptor, selectedCustomer])

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

  // Show error only after field has been touched
  const showError = field.state.meta.isTouched && error

  const handleCustomerSelect = (customerId: string | Array<string>) => {
    const id = Array.isArray(customerId) ? customerId[0] : customerId
    field.handleChange(id)
    const customer = customersData?.items.find((c) => c.id === id)
    if (customer) {
      setSelectedCustomer(customer)
      setCustomerSearchQuery(customer.name)
    }
  }

  const handleClear = () => {
    field.handleChange('')
    setSelectedCustomer(null)
    setCustomerSearchQuery('')
  }

  return (
    <div className="form-control w-full">
      {/* Label */}
      {props.label && (
        <label className="label pb-2">
          <span className="label-text font-medium">{props.label}</span>
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
        placeholder={
          selectedCustomer
            ? selectedCustomer.name
            : props.placeholder || t('common:select')
        }
        clearable
        onClear={handleClear}
        disabled={props.disabled || field.state.meta.isValidating}
        error={showError ? error : undefined}
        onOpen={() => {
          setCustomerSearchQuery('')
        }}
        onClose={() => {
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
