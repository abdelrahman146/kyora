/**
 * AddressSheet Component
 *
 * Reusable bottom sheet for adding/editing customer addresses.
 * Handles form validation, submission, and RTL support.
 *
 * Features:
 * - Mobile-first responsive design
 * - Country and phone code selection from metadata
 * - Bilingual support (Arabic/English)
 * - Form validation with Zod via TanStack Form
 * - Auto-linking country to phone code
 */

import { useEffect, useMemo, useState } from 'react'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { BottomSheet } from '../../molecules/BottomSheet'
import { CountrySelect } from '../../molecules/CountrySelect'
import { PhoneCodeSelect } from '../../molecules/PhoneCodeSelect'
import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress } from '@/api/customer'
import { FormInput } from '@/components'
import { useCountriesQuery } from '@/api/metadata'
import { buildE164Phone, parseE164Phone } from '@/lib/phone'
import { showErrorToast, showSuccessToast } from '@/lib/toast'
import { translateErrorAsync } from '@/lib/translateError'

export interface AddressSheetProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (
    data: CreateAddressRequest | UpdateAddressRequest,
  ) => Promise<CustomerAddress>
  address?: CustomerAddress // If provided, we're editing
  submitLabel?: string
}

// Zod schema
const addressSchema = z.object({
  countryCode: z.string().length(2, 'validation.country_required'),
  state: z.string().min(1, 'validation.state_required'),
  city: z.string().min(1, 'validation.city_required'),
  phoneCode: z.string().min(1, 'validation.phone_code_required'),
  phoneNumber: z.string().min(1, 'validation.phone_required'),
  street: z.string().optional(),
  zipCode: z.string().optional(),
})

type FormData = z.infer<typeof addressSchema>

export function AddressSheet({
  isOpen,
  onClose,
  onSubmit,
  address,
  submitLabel,
}: AddressSheetProps) {
  const { t } = useTranslation()
  const { t: tErrors } = useTranslation('errors')

  // Fetch countries using TanStack Query
  const { data: countries = [], isLoading: countriesLoading } =
    useCountriesQuery()
  const countriesReady = countries.length > 0 && !countriesLoading

  // Track selected country code for auto-linking phone code
  const [selectedCountryCode, setSelectedCountryCode] = useState(
    address?.countryCode ?? '',
  )

  // Parse address phone if editing
  const initialPhoneData = useMemo(() => {
    if (address) {
      return parseE164Phone(address.phoneCode, address.phoneNumber)
    }
    return { phoneCode: '', phoneNumber: '' }
  }, [address])

  // Default values
  const defaultValues: FormData = {
    countryCode: address?.countryCode ?? '',
    state: address?.state ?? '',
    city: address?.city ?? '',
    phoneCode: initialPhoneData.phoneCode,
    phoneNumber: initialPhoneData.phoneNumber,
    street: address?.street ?? '',
    zipCode: address?.zipCode ?? '',
  }

  // TanStack Form setup
  const form = useForm({
    defaultValues,
    validators: {
      onBlur: addressSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        // Build E.164 phone
        const phoneData = buildE164Phone(value.phoneCode, value.phoneNumber)

        if (address) {
          // Update
          const updateData: UpdateAddressRequest = {
            countryCode: value.countryCode,
            state: value.state,
            city: value.city,
            phoneCode: value.phoneCode,
            phoneNumber: value.phoneNumber,
            street: value.street,
            zipCode: value.zipCode,
          }
          await onSubmit(updateData)
          showSuccessToast(t('customers.address.update_success'))
        } else {
          // Create
          const createData: CreateAddressRequest = {
            countryCode: value.countryCode,
            state: value.state,
            city: value.city,
            phoneCode: value.phoneCode,
            phone: phoneData.e164, // Backend expects 'phone' field with E.164 format
            street: value.street,
            zipCode: value.zipCode,
          }
          await onSubmit(createData)
          showSuccessToast(t('customers.address.create_success'))
        }
        onClose()
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        showErrorToast(message)
      }
    },
  })

  // Subscribe to specific form state to minimize re-renders
  // TanStack Form pattern: extract form state subscriptions early to avoid
  // repeated subscriptions in JSX
  const isSubmitting = form.state.isSubmitting
  const isDirty = form.state.isDirty

  // Reset form when sheet opens or address changes
  useEffect(() => {
    if (isOpen) {
      form.reset()
      // Update default values when address changes
      form.setFieldValue('countryCode', address?.countryCode ?? '')
      form.setFieldValue('state', address?.state ?? '')
      form.setFieldValue('city', address?.city ?? '')
      form.setFieldValue('phoneCode', initialPhoneData.phoneCode)
      form.setFieldValue('phoneNumber', initialPhoneData.phoneNumber)
      form.setFieldValue('street', address?.street ?? '')
      form.setFieldValue('zipCode', address?.zipCode ?? '')
      setSelectedCountryCode(address?.countryCode ?? '')
    }
  }, [isOpen, address, initialPhoneData, form])

  // Auto-link country to phone code when country changes
  useEffect(() => {
    if (selectedCountryCode && countriesReady) {
      const country = countries.find((c) => c.code === selectedCountryCode)
      if (country?.phonePrefix) {
        form.setFieldValue('phoneCode', country.phonePrefix)
      }
    }
  }, [selectedCountryCode, countries, countriesReady, form])

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={onClose}
      title={
        address
          ? t('customers.address.edit_title')
          : t('customers.address.add_title')
      }
    >
      <form
        onSubmit={(e) => {
          e.preventDefault()
          e.stopPropagation()
          void form.handleSubmit()
        }}
        className="space-y-4"
      >
        {/* Country */}
        <form.Field
          name="countryCode"
          validators={{
            onBlur: addressSchema.shape.countryCode,
          }}
        >
          {(field) => (
            <CountrySelect
              value={field.state.value}
              onChange={(value) => {
                field.handleChange(value)
                setSelectedCountryCode(value)
              }}
              error={
                field.state.meta.errors.length > 0
                  ? tErrors(
                      field.state.meta.errors[0]?.message ??
                        'validation.invalid',
                    )
                  : undefined
              }
              required
            />
          )}
        </form.Field>

        {/* State */}
        <form.Field
          name="state"
          validators={{
            onBlur: addressSchema.shape.state,
          }}
        >
          {(field) => (
            <FormInput
              label={t('customers.form.state')}
              placeholder={t('customers.form.state_placeholder')}
              required
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={
                field.state.meta.errors.length > 0
                  ? tErrors(
                      field.state.meta.errors[0]?.message ??
                        'validation.invalid',
                    )
                  : undefined
              }
            />
          )}
        </form.Field>

        {/* City */}
        <form.Field
          name="city"
          validators={{
            onBlur: addressSchema.shape.city,
          }}
        >
          {(field) => (
            <FormInput
              label={t('customers.form.city')}
              placeholder={t('customers.form.city_placeholder')}
              required
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={
                field.state.meta.errors.length > 0
                  ? tErrors(
                      field.state.meta.errors[0]?.message ??
                        'validation.invalid',
                    )
                  : undefined
              }
            />
          )}
        </form.Field>

        {/* Street (Optional) */}
        <form.Field name="street">
          {(field) => (
            <FormInput
              label={t('customers.form.street')}
              placeholder={t('customers.form.street_placeholder')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
            />
          )}
        </form.Field>

        {/* Zip Code (Optional) */}
        <form.Field name="zipCode">
          {(field) => (
            <FormInput
              label={t('customers.form.zip_code')}
              placeholder={t('customers.form.zip_placeholder')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
            />
          )}
        </form.Field>

        {/* Phone Code - Auto-updated from country, disabled */}
        <form.Field
          name="phoneCode"
          validators={{
            onBlur: addressSchema.shape.phoneCode,
          }}
        >
          {(field) => (
            <PhoneCodeSelect
              value={field.state.value}
              onChange={(value) => field.handleChange(value)}
              error={
                field.state.meta.errors.length > 0
                  ? tErrors(
                      field.state.meta.errors[0]?.message ??
                        'validation.invalid',
                    )
                  : undefined
              }
              disabled
              required
            />
          )}
        </form.Field>

        {/* Phone Number */}
        <form.Field
          name="phoneNumber"
          validators={{
            onBlur: addressSchema.shape.phoneNumber,
          }}
        >
          {(field) => (
            <FormInput
              label={t('customers.form.phone_number')}
              placeholder={t('customers.form.phone_placeholder')}
              type="tel"
              inputMode="tel"
              required
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={
                field.state.meta.errors.length > 0
                  ? tErrors(
                      field.state.meta.errors[0]?.message ??
                        'validation.invalid',
                    )
                  : undefined
              }
            />
          )}
        </form.Field>

        {/* Footer Actions */}
        <div className="flex gap-2 pt-4">
          <button
            type="button"
            className="btn btn-ghost flex-1"
            onClick={onClose}
            disabled={isSubmitting}
          >
            {t('common.cancel')}
          </button>
          <button
            type="submit"
            className="btn btn-primary flex-1"
            disabled={isSubmitting || (address ? !isDirty : false)}
          >
            {isSubmitting && (
              <span className="loading loading-spinner loading-sm" />
            )}
            {submitLabel ?? (address ? t('common.update') : t('common.add'))}
          </button>
        </div>
      </form>
    </BottomSheet>
  )
}
