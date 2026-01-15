/**
 * AddressSheet Component
 *
 * A callback-based sheet for creating/editing customer addresses.
 * Uses `onSubmit` callback for create/update operations.
 * The parent component handles the mutation and success behavior.
 *
 * For standalone usage (e.g., in order creation), use StandaloneAddressSheet
 * which wraps this component and handles mutations internally.
 */

import { Link } from '@tanstack/react-router'
import { useEffect, useId, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { CountrySelect } from './CountrySelect'
import { PhoneCodeSelect } from './PhoneCodeSelect'
import { ShippingZoneSelect } from './ShippingZoneSelect'
import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress } from '@/api/customer'
import { useShippingZonesQuery } from '@/api/business'
import { useCountriesQuery } from '@/api/metadata'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { useKyoraForm } from '@/lib/form'
import { parseE164Phone } from '@/lib/phone'

export interface AddressSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  onSubmit: (
    data: CreateAddressRequest | UpdateAddressRequest,
  ) => Promise<CustomerAddress>
  /** Existing address for edit mode. If undefined, creates a new address */
  address?: CustomerAddress
  /** Custom submit button label */
  submitLabel?: string
  /** Whether the sheet is currently submitting (controlled externally) */
  isSubmitting?: boolean
}

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
  businessDescriptor,
  onSubmit,
  address,
  submitLabel,
  isSubmitting: externalIsSubmitting,
}: AddressSheetProps) {
  const { t: tCustomers } = useTranslation('customers')
  const { t: tCommon } = useTranslation('common')
  const formId = useId()

  const { data: countries = [], isLoading: countriesLoading } =
    useCountriesQuery()
  const { data: shippingZones = [], isLoading: zonesLoading } =
    useShippingZonesQuery(businessDescriptor, isOpen)
  const countriesReady = countries.length > 0 && !countriesLoading
  const zonesReady = shippingZones.length > 0 && !zonesLoading

  const [selectedZoneId, setSelectedZoneId] = useState(
    address?.shippingZoneId ?? '',
  )
  const [selectedCountryCode, setSelectedCountryCode] = useState(
    address?.countryCode ?? '',
  )

  const initialPhoneData = useMemo(() => {
    if (address) {
      return parseE164Phone(address.phoneCode, address.phoneNumber)
    }
    return { phoneCode: '', phoneNumber: '' }
  }, [address])

  const availableCountries = useMemo(() => {
    if (!selectedZoneId || !zonesReady) return countries

    const selectedZone = shippingZones.find(
      (zone) => zone.id === selectedZoneId,
    )
    if (!selectedZone) return countries

    return countries.filter((c) => selectedZone.countries.includes(c.code))
  }, [selectedZoneId, shippingZones, countries, zonesReady])

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

  const form = useKyoraForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      if (address) {
        // Update mode
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
      } else {
        // Create mode
        const createData: CreateAddressRequest = {
          countryCode: value.countryCode,
          state: value.state,
          city: value.city,
          phoneCode: value.phoneCode,
          phoneNumber: value.phoneNumber,
          street: value.street,
          zipCode: value.zipCode,
        }
        await onSubmit(createData)
      }
      onClose()
    },
  })

  const [formIsSubmitting, setFormIsSubmitting] = useState(false)
  const [isDirty, setIsDirty] = useState(false)

  useEffect(() => {
    const unsubscribe = form.store.subscribe(() => {
      setFormIsSubmitting(form.store.state.isSubmitting)
      setIsDirty(form.store.state.isDirty)
    })
    return unsubscribe
  }, [form])

  // Reset when sheet closes
  useEffect(() => {
    if (!isOpen) {
      form.reset()
      setSelectedZoneId(address?.shippingZoneId ?? '')
      setSelectedCountryCode(address?.countryCode ?? '')
    }
  }, [isOpen, address, form])

  // Re-populate when sheet opens with existing address
  useEffect(() => {
    if (!isOpen) return

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
  }, [isOpen, address, initialPhoneData, form])

  // Clear country when zone changes and country not in new zone
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

  // Sync phone code when country changes
  useEffect(() => {
    if (selectedCountryCode && countriesReady) {
      const country = countries.find((c) => c.code === selectedCountryCode)
      if (country?.phonePrefix) {
        form.setFieldValue('phoneCode', country.phonePrefix)
      }
    }
  }, [selectedCountryCode, countries, countriesReady, form])

  const isPending = externalIsSubmitting || formIsSubmitting

  const safeClose = () => {
    if (isPending) return
    onClose()
  }

  const getTitle = () => {
    return address
      ? tCustomers('address.edit_title')
      : tCustomers('address.add_title')
  }

  const getSubmitLabel = () => {
    if (isPending) {
      return <span className="loading loading-spinner loading-sm" />
    }
    if (submitLabel) return submitLabel
    return address ? tCommon('update') : tCommon('add')
  }

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={safeClose}
        title={getTitle()}
        footer={
          <div className="flex gap-2">
            <button
              type="button"
              className="btn btn-ghost flex-1"
              onClick={safeClose}
              disabled={isPending}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton
              variant="primary"
              className="flex-1"
              disabled={address ? !isDirty : false}
              form={`address-form-${formId}`}
            >
              {getSubmitLabel()}
            </form.SubmitButton>
          </div>
        }
        side="end"
        size="md"
        closeOnOverlayClick={!isPending}
        closeOnEscape={!isPending}
        contentClassName="space-y-4"
        ariaLabel={getTitle()}
      >
        <form.FormRoot
          id={`address-form-${formId}`}
          className="space-y-4"
          aria-busy={isPending}
        >
          {/* Shipping Zone */}
          <form.AppField
            name="shippingZoneId"
            validators={{
              onBlur: addressSchema.shape.shippingZoneId,
            }}
          >
            {(field) => (
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

          {!zonesLoading && shippingZones.length === 0 && (
            <div className="text-sm text-base-content/70 -mt-2">
              {tCustomers('address.no_zones_message')}{' '}
              <Link
                to="/business/$businessDescriptor"
                params={{ businessDescriptor }}
                className="link link-primary"
              >
                {tCustomers('address.configure_zones_link')}
              </Link>
            </div>
          )}

          {/* Country */}
          <form.AppField
            name="countryCode"
            validators={{
              onBlur: z.string().length(2, 'validation.country_required'),
            }}
          >
            {(field) => (
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
                    {tCustomers('address.filtered_by_zone', {
                      count: availableCountries.length,
                    })}
                  </div>
                )}
              </>
            )}
          </form.AppField>

          {/* City and State */}
          <form.AppField
            name="city"
            validators={{
              onBlur: z.string().min(1, 'validation.city_required'),
            }}
          >
            {(field) => (
              <field.TextField
                label={tCustomers('form.city')}
                placeholder={tCustomers('form.city_placeholder')}
                required
              />
            )}
          </form.AppField>

          <form.AppField
            name="state"
            validators={{
              onBlur: z.string().min(1, 'validation.state_required'),
            }}
          >
            {(field) => (
              <field.TextField
                label={tCustomers('form.state')}
                placeholder={tCustomers('form.state_placeholder')}
                required
              />
            )}
          </form.AppField>

          {/* Street and Zip Code */}
          <form.AppField name="street">
            {(field) => (
              <field.TextField
                label={tCustomers('form.street')}
                placeholder={tCustomers('form.street_placeholder')}
              />
            )}
          </form.AppField>

          <form.AppField name="zipCode">
            {(field) => (
              <field.TextField
                label={tCustomers('form.zip_code')}
                placeholder={tCustomers('form.zip_placeholder')}
              />
            )}
          </form.AppField>

          {/* Phone */}
          <div className="space-y-4">
            <form.AppField
              name="phoneCode"
              validators={{
                onBlur: z.string().min(1, 'validation.phone_code_required'),
              }}
            >
              {(field) => (
                <PhoneCodeSelect
                  value={field.state.value}
                  onChange={(value: string) => field.handleChange(value)}
                  disabled
                  required
                />
              )}
            </form.AppField>

            <form.AppField
              name="phoneNumber"
              validators={{
                onBlur: z.string().min(1, 'validation.phone_required'),
              }}
            >
              {(field) => (
                <field.TextField
                  type="tel"
                  label={tCustomers('form.phone_number')}
                  placeholder={tCustomers('form.phone_placeholder')}
                  dir="ltr"
                  inputMode="tel"
                  required
                />
              )}
            </form.AppField>
          </div>
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
