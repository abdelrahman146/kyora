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
import { Link } from '@tanstack/react-router'

import { BottomSheet } from '../../molecules/BottomSheet'
import { CountrySelect } from '../../molecules/CountrySelect'
import { PhoneCodeSelect } from '../../molecules/PhoneCodeSelect'
import { ShippingZoneSelect } from '../../molecules/ShippingZoneSelect'
import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress } from '@/api/customer'
import { useKyoraForm } from '@/lib/form'
import { useCountriesQuery } from '@/api/metadata'
import { useShippingZonesQuery } from '@/api/business'
import { buildE164Phone, parseE164Phone } from '@/lib/phone'
import { showErrorToast, showSuccessToast } from '@/lib/toast'

export interface AddressSheetProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (
    data: CreateAddressRequest | UpdateAddressRequest,
  ) => Promise<CustomerAddress>
  address?: CustomerAddress // If provided, we're editing
  submitLabel?: string
  businessDescriptor: string // Required for fetching shipping zones
}

// Zod schema
const addressSchema = z.object({
  shippingZoneId: z.string().min(1, 'validation.shipping_zone_required'),
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
  businessDescriptor,
}: AddressSheetProps) {
  const { t } = useTranslation()

  // Fetch countries and shipping zones
  const { data: countries = [], isLoading: countriesLoading } =
    useCountriesQuery()
  const { data: shippingZones = [], isLoading: zonesLoading } =
    useShippingZonesQuery(businessDescriptor, isOpen)
  const countriesReady = countries.length > 0 && !countriesLoading
  const zonesReady = shippingZones.length > 0 && !zonesLoading

  // Track selected shipping zone and country for filtering/auto-linking
  const [selectedZoneId, setSelectedZoneId] = useState(
    address?.shippingZoneId ?? '',
  )
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

  // Filter countries by selected shipping zone
  const availableCountries = useMemo(() => {
    if (!selectedZoneId || !zonesReady) return countries

    const selectedZone = shippingZones.find(
      (zone) => zone.id === selectedZoneId,
    )
    if (!selectedZone) return countries

    // Filter countries to only those in the selected zone
    return countries.filter((c) => selectedZone.countries.includes(c.code))
  }, [selectedZoneId, shippingZones, countries, zonesReady])

  // Default values
  const defaultValues: FormData = {
    shippingZoneId: address?.shippingZoneId ?? '',
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
    onSubmit: async ({ value }) => {
      try {
        // Build E.164 phone
        const phoneData = buildE164Phone(value.phoneCode, value.phoneNumber)

        if (address) {
          // Update
          const updateData: UpdateAddressRequest = {
            shippingZoneId: value.shippingZoneId,
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
            shippingZoneId: value.shippingZoneId,
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

  // Get form state
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isDirty, setIsDirty] = useState(false)

  useEffect(() => {
    const unsubscribe = form.store.subscribe(() => {
      setIsSubmitting(form.store.state.isSubmitting)
      setIsDirty(form.store.state.isDirty)
    })
    return unsubscribe
  }, [form])

  // Reset form when sheet opens or address changes
  useEffect(() => {
    if (isOpen) {
      form.reset()
      // Update default values when address changes
      form.setFieldValue('shippingZoneId', address?.shippingZoneId ?? '')
      form.setFieldValue('countryCode', address?.countryCode ?? '')
      form.setFieldValue('state', address?.state ?? '')
      form.setFieldValue('city', address?.city ?? '')
      form.setFieldValue('phoneCode', initialPhoneData.phoneCode)
      form.setFieldValue('phoneNumber', initialPhoneData.phoneNumber)
      form.setFieldValue('street', address?.street ?? '')
      form.setFieldValue('zipCode', address?.zipCode ?? '')
      setSelectedZoneId(address?.shippingZoneId ?? '')
      setSelectedCountryCode(address?.countryCode ?? '')
    }
  }, [isOpen, address, initialPhoneData, form])

  // Reset country when zone changes and country not in new zone
  useEffect(() => {
    if (
      selectedZoneId &&
      selectedCountryCode &&
      availableCountries.length > 0
    ) {
      const isCountryInZone = availableCountries.some(
        (c) => c.code === selectedCountryCode,
      )
      if (!isCountryInZone) {
        form.setFieldValue('countryCode', '')
        setSelectedCountryCode('')
      }
    }
  }, [selectedZoneId, selectedCountryCode, availableCountries, form])

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
      <form.AppForm>
        <form.FormRoot className="space-y-4" aria-busy={isSubmitting}>
          {/* Shipping Zone */}
          <form.AppField
            name="shippingZoneId"
            validators={{
              onBlur: addressSchema.shape.shippingZoneId,
            }}
          >
            {(field: any) => (
              <ShippingZoneSelect
                value={field.state.value}
                onChange={(value: string) => {
                  field.handleChange(value)
                  setSelectedZoneId(value)
                }}
                zones={shippingZones}
                isLoading={zonesLoading}
                required
              />
            )}
          </form.AppField>

          {/* Help text for no zones */}
          {!zonesLoading && shippingZones.length === 0 && (
            <div className="text-sm text-base-content/70 -mt-2">
              {t('customers.address.no_zones_message')}{' '}
              <Link
                to="/business/$businessDescriptor"
                params={{ businessDescriptor }}
                className="link link-primary"
              >
                {t('customers.address.configure_zones_link')}
              </Link>
            </div>
          )}

          {/* Country - filtered by zone */}
          <form.AppField
            name="countryCode"
            validators={{
              onBlur: z.string().length(2, 'validation.country_required'),
            }}
          >
            {(field: any) => (
              <>
                <CountrySelect
                  value={field.state.value}
                  onChange={(value: string) => {
                    field.handleChange(value)
                    setSelectedCountryCode(value)
                  }}
                  availableCountries={availableCountries}
                  disabled={!selectedZoneId}
                  required
                />
                {selectedZoneId && availableCountries.length > 0 && (
                  <div className="text-xs text-base-content/60 mt-1">
                    {t('customers.address.filtered_by_zone', {
                      count: availableCountries.length,
                    })}
                  </div>
                )}
              </>
            )}
          </form.AppField>

          {/* State */}
          <form.AppField
            name="state"
            validators={{
              onBlur: z.string().min(1, 'validation.state_required'),
            }}
          >
            {(field) => (
              <field.TextField
                label={t('customers.form.state')}
                placeholder={t('customers.form.state_placeholder')}
                required
              />
            )}
          </form.AppField>

          {/* City */}
          <form.AppField
            name="city"
            validators={{
              onBlur: z.string().min(1, 'validation.city_required'),
            }}
          >
            {(field) => (
              <field.TextField
                label={t('customers.form.city')}
                placeholder={t('customers.form.city_placeholder')}
                required
              />
            )}
          </form.AppField>

          {/* Street (Optional) */}
          <form.AppField name="street">
            {(field) => (
              <field.TextField
                label={t('customers.form.street')}
                placeholder={t('customers.form.street_placeholder')}
              />
            )}
          </form.AppField>

          {/* Zip Code (Optional) */}
          <form.AppField name="zipCode">
            {(field) => (
              <field.TextField
                label={t('customers.form.zip_code')}
                placeholder={t('customers.form.zip_placeholder')}
              />
            )}
          </form.AppField>

          {/* Phone Code - Auto-updated from country, disabled */}
          <form.AppField
            name="phoneCode"
            validators={{
              onBlur: z.string().min(1, 'validation.phone_code_required'),
            }}
          >
            {(field: any) => (
              <PhoneCodeSelect
                value={field.state.value}
                onChange={(value: string) => field.handleChange(value)}
                disabled
                required
              />
            )}
          </form.AppField>

          {/* Phone Number */}
          <form.AppField
            name="phoneNumber"
            validators={{
              onBlur: z.string().min(1, 'validation.phone_required'),
            }}
          >
            {(field) => (
              <field.TextField
                type="tel"
                label={t('customers.form.phone_number')}
                placeholder={t('customers.form.phone_placeholder')}
                required
              />
            )}
          </form.AppField>

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
      </form.AppForm>
    </BottomSheet>
  )
}
