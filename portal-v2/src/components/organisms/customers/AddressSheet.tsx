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
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { BottomSheet } from '../../molecules/BottomSheet'
import { CountrySelect } from '../../molecules/CountrySelect'
import { PhoneCodeSelect } from '../../molecules/PhoneCodeSelect'
import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress } from '@/api/customer'
import { useKyoraForm } from '@/lib/form'
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

  // TanStack Form setup with useKyoraForm
  const form = useKyoraForm({
    defaultValues,
    validators: {
      countryCode: { onBlur: z.string().length(2, 'country_required') },
      state: { onBlur: z.string().min(1, 'state_required') },
      city: { onBlur: z.string().min(1, 'city_required') },
      phoneCode: { onBlur: z.string().min(1, 'phone_code_required') },
      phoneNumber: { onBlur: z.string().min(1, 'phone_required') },
      street: { onBlur: z.string().optional() },
      zipCode: { onBlur: z.string().optional() },
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
        showErrorToast((error as Error).message)
      }
    },
  })

  // Subscribe to specific form state to minimize re-renders
  const isSubmitting = form.useStore((state) => state.isSubmitting)
  const isDirty = form.useStore((state) => state.isDirty)

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
      <form.FormRoot
        className="space-y-4"
        aria-busy={isSubmitting}
      >
        {/* Country */}
        <form.Field name="countryCode">
          {(field: any) => (
            <CountrySelect
              value={field.state.value}
              onChange={(value: string) => {
                field.handleChange(value)
                setSelectedCountryCode(value)
              }}
              required
            />
          )}
        </form.Field>

        {/* State */}
        <form.TextField
          name="state"
          label={t('customers.form.state')}
          placeholder={t('customers.form.state_placeholder')}
          required
        />

        {/* City */}
        <form.TextField
          name="city"
          label={t('customers.form.city')}
          placeholder={t('customers.form.city_placeholder')}
          required
        />

        {/* Street (Optional) */}
        <form.TextField
          name="street"
          label={t('customers.form.street')}
          placeholder={t('customers.form.street_placeholder')}
        />

        {/* Zip Code (Optional) */}
        <form.TextField
          name="zipCode"
          label={t('customers.form.zip_code')}
          placeholder={t('customers.form.zip_placeholder')}
        />

        {/* Phone Code - Auto-updated from country, disabled */}
        <form.Field name="phoneCode">
          {(field: any) => (
            <PhoneCodeSelect
              value={field.state.value}
              onChange={(value: string) => field.handleChange(value)}
              countryCode={selectedCountryCode}
              disabled
              required
            />
          )}
        </form.Field>

        {/* Phone Number */}
        <form.TextField
          name="phoneNumber"
          type="tel"
          label={t('customers.form.phone_number')}
          placeholder={t('customers.form.phone_placeholder')}
          inputMode="tel"
          required
          dir="ltr"
        />

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
          <form.SubmitButton
            variant="primary"
            className="flex-1"
            disabled={address ? !isDirty : false}
          >
            {isSubmitting && (
              <span className="loading loading-spinner loading-sm" />
            )}
            {submitLabel ?? (address ? t('common.update') : t('common.add'))}
          </form.SubmitButton>
        </div>
      </form.FormRoot>
    </BottomSheet>
  )
}
