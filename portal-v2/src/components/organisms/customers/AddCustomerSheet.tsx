/**
 * AddCustomerSheet Component
 *
 * Bottom sheet form for creating new customers.
 * Features auto-linking of country to phone code on selection.
 *
 * Uses TanStack Form for form state management with Zod validation.
 */

import { useEffect, useId, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useForm } from '@tanstack/react-form'
import { z } from 'zod'

import { BottomSheet } from '../../molecules/BottomSheet'
import { CountrySelect } from '../../molecules/CountrySelect'
import { PhoneCodeSelect } from '../../molecules/PhoneCodeSelect'
import { SocialMediaInputs } from '../../molecules/SocialMediaInputs'
import type { Customer, CustomerGender } from '@/api/customer'
import { FormInput, FormSelect } from '@/components'
import { useCreateCustomerMutation } from '@/api/customer'
import { useCountriesQuery } from '@/api/metadata'
import { buildE164Phone } from '@/lib/phone'
import { showErrorToast, showSuccessToast } from '@/lib/toast'
import { translateErrorAsync } from '@/lib/translateError'

export interface AddCustomerSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  businessCountryCode: string
  onCreated?: (customer: Customer) => void | Promise<void>
}

const addCustomerSchema = z
  .object({
    name: z.string().trim().min(1, 'validation.required'),
    email: z
      .string()
      .trim()
      .min(1, 'validation.required')
      .email('validation.invalid_email'),
    gender: z.enum(['male', 'female', 'other'], {
      message: 'validation.required',
    }),
    countryCode: z
      .string()
      .trim()
      .min(1, 'validation.required')
      .refine((v) => /^[A-Za-z]{2}$/.test(v), 'validation.invalid_country'),
    phoneCode: z
      .string()
      .trim()
      .refine(
        (v) => v === '' || /^\+?\d{1,4}$/.test(v),
        'validation.invalid_phone_code',
      ),
    phoneNumber: z
      .string()
      .trim()
      .refine(
        (v) => v === '' || /^[0-9\-\s()]{6,20}$/.test(v),
        'validation.invalid_phone',
      ),
    instagramUsername: z.string().trim().optional(),
    facebookUsername: z.string().trim().optional(),
    tiktokUsername: z.string().trim().optional(),
    snapchatUsername: z.string().trim().optional(),
    xUsername: z.string().trim().optional(),
    whatsappNumber: z.string().trim().optional(),
  })
  .refine(
    (values) => {
      const hasPhoneNumber = values.phoneNumber.trim() !== ''
      return !hasPhoneNumber || values.phoneCode.trim() !== ''
    },
    { message: 'validation.required', path: ['phoneCode'] },
  )

export type AddCustomerFormValues = z.infer<typeof addCustomerSchema>

function getDefaultValues(businessCountryCode: string): AddCustomerFormValues {
  return {
    name: '',
    email: '',
    gender: 'other',
    countryCode: businessCountryCode,
    phoneCode: '',
    phoneNumber: '',
    instagramUsername: '',
    facebookUsername: '',
    tiktokUsername: '',
    snapchatUsername: '',
    xUsername: '',
    whatsappNumber: '',
  }
}

export function AddCustomerSheet({
  isOpen,
  onClose,
  businessDescriptor,
  businessCountryCode,
  onCreated,
}: AddCustomerSheetProps) {
  const { t } = useTranslation()
  const { t: tErrors } = useTranslation('errors')
  const formId = useId()

  // Local state for social media values (since form.useStore doesn't exist)
  const [socialMediaValues, setSocialMediaValues] = useState({
    instagramUsername: '',
    facebookUsername: '',
    tiktokUsername: '',
    snapchatUsername: '',
    xUsername: '',
    whatsappNumber: '',
  })

  // Track selected country code for auto-linking phone code
  const [selectedCountryCode, setSelectedCountryCode] =
    useState(businessCountryCode)

  // Fetch countries using TanStack Query
  const { data: countries = [], isLoading: countriesLoading } =
    useCountriesQuery()
  const countriesReady = countries.length > 0 && !countriesLoading

  // Create customer mutation
  const createMutation = useCreateCustomerMutation(businessDescriptor, {
    onSuccess: async (created) => {
      showSuccessToast(t('customers.create_success'))
      if (onCreated) {
        await onCreated(created)
      }
      onClose()
    },
    onError: async (error) => {
      const message = await translateErrorAsync(error, t)
      showErrorToast(message)
    },
  })

  // TanStack Form setup
  const form = useForm({
    defaultValues: getDefaultValues(businessCountryCode),
    validators: {
      onBlur: addCustomerSchema,
    },
    onSubmit: async ({ value }) => {
      const phoneCode = value.phoneCode.trim()
      const phoneNumber = value.phoneNumber.trim()

      const normalizedPhone =
        phoneNumber !== '' && phoneCode !== ''
          ? buildE164Phone(phoneCode, phoneNumber)
          : undefined

      await createMutation.mutateAsync({
        name: value.name.trim(),
        email: value.email.trim(),
        gender: value.gender as CustomerGender,
        countryCode: value.countryCode.trim().toUpperCase(),
        phoneCode: normalizedPhone?.phoneCode,
        phoneNumber: normalizedPhone?.phoneNumber,
        instagramUsername: value.instagramUsername?.trim() || undefined,
        facebookUsername: value.facebookUsername?.trim() || undefined,
        tiktokUsername: value.tiktokUsername?.trim() || undefined,
        snapchatUsername: value.snapchatUsername?.trim() || undefined,
        xUsername: value.xUsername?.trim() || undefined,
        whatsappNumber: value.whatsappNumber?.trim() || undefined,
      })
    },
  })

  const isSubmitting = form.state.isSubmitting || createMutation.isPending

  // Auto-link country to phone code when country changes
  useEffect(() => {
    if (!isOpen) return
    if (!countriesReady) return
    const selected = countries.find((c) => c.code === selectedCountryCode)
    if (!selected?.phonePrefix) return
    form.setFieldValue('phoneCode', selected.phonePrefix)
  }, [isOpen, countriesReady, countries, selectedCountryCode, form])

  // Reset form and local state when sheet closes
  useEffect(() => {
    if (!isOpen) {
      form.reset()
      setSelectedCountryCode(businessCountryCode)
      setSocialMediaValues({
        instagramUsername: '',
        facebookUsername: '',
        tiktokUsername: '',
        snapchatUsername: '',
        xUsername: '',
        whatsappNumber: '',
      })
    }
  }, [isOpen, form, businessCountryCode])

  const safeClose = () => {
    if (isSubmitting) return
    onClose()
  }

  const footer = (
    <div className="flex gap-2">
      <button
        type="button"
        className="btn btn-ghost flex-1"
        onClick={safeClose}
        disabled={isSubmitting}
        aria-disabled={isSubmitting}
      >
        {t('common.cancel')}
      </button>
      <button
        type="submit"
        form={`add-customer-form-${formId}`}
        className="btn btn-primary flex-1"
        disabled={isSubmitting}
        aria-disabled={isSubmitting}
      >
        {isSubmitting
          ? t('customers.create_submitting')
          : t('customers.create_submit')}
      </button>
    </div>
  )

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={safeClose}
      title={t('customers.create_title')}
      footer={footer}
      side="end"
      size="md"
      closeOnOverlayClick={!isSubmitting}
      closeOnEscape={!isSubmitting}
      contentClassName="space-y-4"
      ariaLabel={t('customers.create_title')}
    >
      <form
        id={`add-customer-form-${formId}`}
        onSubmit={(e) => {
          e.preventDefault()
          e.stopPropagation()
          void form.handleSubmit()
        }}
        className="space-y-4"
        aria-busy={isSubmitting}
      >
        <form.Field
          name="name"
          validators={{
            onBlur: addCustomerSchema.shape.name,
          }}
        >
          {(field) => (
            <FormInput
              label={t('customers.form.name')}
              placeholder={t('customers.form.name_placeholder')}
              autoComplete="name"
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

        <form.Field
          name="email"
          validators={{
            onBlur: addCustomerSchema.shape.email,
          }}
        >
          {(field) => (
            <FormInput
              label={t('customers.form.email')}
              type="email"
              placeholder={t('customers.form.email_placeholder')}
              autoComplete="email"
              inputMode="email"
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

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <form.Field
            name="countryCode"
            validators={{
              onBlur: addCustomerSchema.shape.countryCode,
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
                disabled={isSubmitting}
                required
              />
            )}
          </form.Field>

          <form.Field
            name="gender"
            validators={{
              onBlur: addCustomerSchema.shape.gender,
            }}
          >
            {(field) => (
              <FormSelect<CustomerGender>
                label={t('customers.form.gender')}
                options={[
                  { value: 'male', label: t('customers.form.gender_male') },
                  { value: 'female', label: t('customers.form.gender_female') },
                  { value: 'other', label: t('customers.form.gender_other') },
                ]}
                value={field.state.value}
                onChange={(value) => {
                  // FormSelect can return array for multiSelect, but we use single select
                  const singleValue = Array.isArray(value) ? value[0] : value
                  field.handleChange(singleValue)
                }}
                required
                disabled={isSubmitting}
                placeholder={t('customers.form.select_gender')}
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
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <form.Field
            name="phoneCode"
            validators={{
              onBlur: addCustomerSchema.shape.phoneCode,
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
                disabled={isSubmitting}
              />
            )}
          </form.Field>

          <div className="sm:col-span-2">
            <form.Field
              name="phoneNumber"
              validators={{
                onBlur: addCustomerSchema.shape.phoneNumber,
              }}
            >
              {(field) => (
                <FormInput
                  label={t('customers.form.phone_number')}
                  placeholder={t('customers.form.phone_placeholder')}
                  autoComplete="tel"
                  inputMode="tel"
                  dir="ltr"
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
          </div>
        </div>

        <SocialMediaInputs
          instagramUsername={socialMediaValues.instagramUsername}
          onInstagramChange={(value) => {
            setSocialMediaValues((prev) => ({
              ...prev,
              instagramUsername: value,
            }))
            form.setFieldValue('instagramUsername', value)
          }}
          facebookUsername={socialMediaValues.facebookUsername}
          onFacebookChange={(value) => {
            setSocialMediaValues((prev) => ({
              ...prev,
              facebookUsername: value,
            }))
            form.setFieldValue('facebookUsername', value)
          }}
          tiktokUsername={socialMediaValues.tiktokUsername}
          onTiktokChange={(value) => {
            setSocialMediaValues((prev) => ({ ...prev, tiktokUsername: value }))
            form.setFieldValue('tiktokUsername', value)
          }}
          snapchatUsername={socialMediaValues.snapchatUsername}
          onSnapchatChange={(value) => {
            setSocialMediaValues((prev) => ({
              ...prev,
              snapchatUsername: value,
            }))
            form.setFieldValue('snapchatUsername', value)
          }}
          xUsername={socialMediaValues.xUsername}
          onXChange={(value) => {
            setSocialMediaValues((prev) => ({ ...prev, xUsername: value }))
            form.setFieldValue('xUsername', value)
          }}
          whatsappNumber={socialMediaValues.whatsappNumber}
          onWhatsappChange={(value) => {
            setSocialMediaValues((prev) => ({ ...prev, whatsappNumber: value }))
            form.setFieldValue('whatsappNumber', value)
          }}
          disabled={isSubmitting}
          defaultExpanded={false}
        />
      </form>
    </BottomSheet>
  )
}
