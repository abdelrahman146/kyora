import { Link } from '@tanstack/react-router'
import { useEffect, useId, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { ShippingZoneSelect } from './ShippingZoneSelect'
import { CountrySelect } from './CountrySelect'
import { PhoneCodeSelect } from './PhoneCodeSelect'
import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress } from '@/api/customer'
import { useShippingZonesQuery } from '@/api/business'
import { useCountriesQuery } from '@/api/metadata'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { useKyoraForm } from '@/lib/form'
import { buildE164Phone, parseE164Phone } from '@/lib/phone'
import { showErrorToast, showSuccessToast } from '@/lib/toast'

export interface AddressSheetProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (
    data: CreateAddressRequest | UpdateAddressRequest,
  ) => Promise<CustomerAddress>
  address?: CustomerAddress
  submitLabel?: string
  businessDescriptor: string
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
  onSubmit,
  address,
  submitLabel,
  businessDescriptor,
}: AddressSheetProps) {
  const { t } = useTranslation()
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
      try {
        const phoneData = buildE164Phone(value.phoneCode, value.phoneNumber)

        if (address) {
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
          const createData: CreateAddressRequest = {
            shippingZoneId: value.shippingZoneId,
            countryCode: value.countryCode,
            state: value.state,
            city: value.city,
            phoneCode: value.phoneCode,
            phone: phoneData.e164,
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

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isDirty, setIsDirty] = useState(false)

  useEffect(() => {
    const unsubscribe = form.store.subscribe(() => {
      setIsSubmitting(form.store.state.isSubmitting)
      setIsDirty(form.store.state.isDirty)
    })
    return unsubscribe
  }, [form])

  useEffect(() => {
    if (isOpen) {
      form.reset()
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

  useEffect(() => {
    if (selectedCountryCode && countriesReady) {
      const country = countries.find((c) => c.code === selectedCountryCode)
      if (country?.phonePrefix) {
        form.setFieldValue('phoneCode', country.phonePrefix)
      }
    }
  }, [selectedCountryCode, countries, countriesReady, form])

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={onClose}
        title={
          address
            ? t('customers.address.edit_title')
            : t('customers.address.add_title')
        }
        footer={
          <div className="flex gap-2">
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
              form={`address-form-${formId}`}
            >
              {isSubmitting && (
                <span className="loading loading-spinner loading-sm" />
              )}
              {submitLabel ?? (address ? t('common.update') : t('common.add'))}
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot
          id={`address-form-${formId}`}
          className="space-y-4"
          aria-busy={isSubmitting}
        >
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

          <form.AppField name="street">
            {(field) => (
              <field.TextField
                label={t('customers.form.street')}
                placeholder={t('customers.form.street_placeholder')}
              />
            )}
          </form.AppField>

          <form.AppField name="zipCode">
            {(field) => (
              <field.TextField
                label={t('customers.form.zip_code')}
                placeholder={t('customers.form.zip_placeholder')}
              />
            )}
          </form.AppField>

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
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
