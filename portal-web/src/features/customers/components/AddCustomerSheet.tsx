import { useEffect, useId, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { AdditionalDetailsInputs } from './AdditionalDetailsInputs'
import { CountrySelect } from './CountrySelect'
import { PhoneCodeSelect } from './PhoneCodeSelect'
import { SocialMediaInputs } from './SocialMediaInputs'
import type { Customer, CustomerGender } from '@/api/customer'
import { useCreateCustomerMutation } from '@/api/customer'
import { useCountriesQuery } from '@/api/metadata'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { useKyoraForm } from '@/lib/form'
import { buildE164Phone } from '@/lib/phone'
import { showSuccessToast } from '@/lib/toast'

export interface AddCustomerSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  businessCountryCode: string
  onCreated?: (customer: Customer) => void | Promise<void>
}

const customerSchema = z.object({
  name: z.string().trim().min(1, 'validation.required'),
  // Optional fields - email can be empty string or valid email
  email: z
    .string()
    .trim()
    .refine(
      (v) => v === '' || z.string().email().safeParse(v).success,
      'validation.invalid_email',
    ),
  // Optional - string field where empty means not set, validated values are 'male', 'female', 'other'
  gender: z
    .string()
    .refine(
      (v) => v === '' || ['male', 'female', 'other'].includes(v),
      'validation.invalid_gender',
    ),
  countryCode: z
    .string()
    .trim()
    .min(1, 'validation.required')
    .refine((v) => /^[A-Za-z]{2}$/.test(v), 'validation.invalid_country'),
  // Required - phone fields
  phoneCode: z
    .string()
    .trim()
    .min(1, 'validation.required')
    .refine((v) => /^\+?\d{1,4}$/.test(v), 'validation.invalid_phone_code'),
  phoneNumber: z
    .string()
    .trim()
    .min(1, 'validation.required')
    .refine((v) => /^[0-9\-\s()]{6,20}$/.test(v), 'validation.invalid_phone'),
  // Social media - all optional
  instagramUsername: z.string().trim().optional(),
  facebookUsername: z.string().trim().optional(),
  tiktokUsername: z.string().trim().optional(),
  snapchatUsername: z.string().trim().optional(),
  xUsername: z.string().trim().optional(),
  whatsappNumber: z.string().trim().optional(),
})

export type AddCustomerFormValues = z.infer<typeof customerSchema>

function getDefaultValues(businessCountryCode: string): AddCustomerFormValues {
  return {
    name: '',
    email: '',
    gender: '', // Optional - empty means not set
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
  const { t: tCustomers } = useTranslation('customers')
  const { t: tCommon } = useTranslation('common')
  const formId = useId()

  const [selectedCountryCode, setSelectedCountryCode] =
    useState(businessCountryCode)

  const { data: countries = [], isLoading: countriesLoading } =
    useCountriesQuery()
  const countriesReady = countries.length > 0 && !countriesLoading

  const createMutation = useCreateCustomerMutation(businessDescriptor, {
    onSuccess: async (created) => {
      showSuccessToast(tCustomers('create_success'))
      if (onCreated) {
        await onCreated(created)
      }
      onClose()
    },
  })

  const form = useKyoraForm({
    defaultValues: getDefaultValues(businessCountryCode),
    onSubmit: async ({ value }) => {
      const phoneCode = value.phoneCode.trim()
      const phoneNumber = value.phoneNumber.trim()

      const normalizedPhone = buildE164Phone(phoneCode, phoneNumber)

      // Handle optional email - only send if not empty
      const email = value.email.trim() || undefined

      // Handle optional gender - transform '' to undefined
      const gender =
        value.gender && value.gender !== ''
          ? (value.gender as CustomerGender)
          : undefined

      await createMutation.mutateAsync({
        name: value.name.trim(),
        email,
        gender,
        countryCode: value.countryCode.trim().toUpperCase(),
        phoneCode: normalizedPhone.phoneCode,
        phoneNumber: normalizedPhone.phoneNumber,
        instagramUsername: value.instagramUsername?.trim() || undefined,
        facebookUsername: value.facebookUsername?.trim() || undefined,
        tiktokUsername: value.tiktokUsername?.trim() || undefined,
        snapchatUsername: value.snapchatUsername?.trim() || undefined,
        xUsername: value.xUsername?.trim() || undefined,
        whatsappNumber: value.whatsappNumber?.trim() || undefined,
      })
    },
  })

  // Auto-populate phone code based on selected country
  useEffect(() => {
    if (!isOpen) return
    if (!countriesReady) return
    const selected = countries.find((c) => c.code === selectedCountryCode)
    if (!selected?.phonePrefix) return
    form.setFieldValue('phoneCode', selected.phonePrefix)
  }, [isOpen, countriesReady, countries, selectedCountryCode, form])

  // Reset form when sheet closes
  useEffect(() => {
    if (!isOpen) {
      form.reset()
      setSelectedCountryCode(businessCountryCode)
    }
  }, [isOpen, form, businessCountryCode])

  const safeClose = () => {
    if (createMutation.isPending) return
    onClose()
  }

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={safeClose}
        title={tCustomers('create_title')}
        footer={
          <div className="flex gap-2">
            <button
              type="button"
              className="btn btn-ghost flex-1"
              onClick={safeClose}
              disabled={createMutation.isPending}
              aria-disabled={createMutation.isPending}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton
              form={`add-customer-form-${formId}`}
              variant="primary"
              className="flex-1"
            >
              {createMutation.isPending
                ? tCustomers('create_submitting')
                : tCustomers('create_submit')}
            </form.SubmitButton>
          </div>
        }
        side="end"
        size="md"
        closeOnOverlayClick={!createMutation.isPending}
        closeOnEscape={!createMutation.isPending}
        contentClassName="space-y-4"
        ariaLabel={tCustomers('create_title')}
      >
        <form.FormRoot
          id={`add-customer-form-${formId}`}
          className="space-y-4"
          aria-busy={createMutation.isPending}
        >
          {/* Name field */}
          <form.AppField
            name="name"
            validators={{
              onBlur: z.string().trim().min(1, 'validation.required'),
            }}
          >
            {(field) => (
              <field.TextField
                label={tCustomers('form.name')}
                placeholder={tCustomers('form.name_placeholder')}
                autoComplete="name"
                required
              />
            )}
          </form.AppField>

          {/* Phone fields - Required */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <form.AppField
              name="phoneCode"
              validators={{
                onBlur: z
                  .string()
                  .trim()
                  .min(1, 'validation.required')
                  .refine(
                    (v) => /^\+?\d{1,4}$/.test(v),
                    'validation.invalid_phone_code',
                  ),
              }}
            >
              {(field) => (
                <PhoneCodeSelect
                  value={field.state.value}
                  onChange={(value: string) => field.handleChange(value)}
                  disabled={createMutation.isPending}
                  required
                />
              )}
            </form.AppField>

            <div className="sm:col-span-2">
              <form.AppField
                name="phoneNumber"
                validators={{
                  onBlur: z
                    .string()
                    .trim()
                    .min(1, 'validation.required')
                    .refine(
                      (v) => /^[0-9\-\s()]{6,20}$/.test(v),
                      'validation.invalid_phone',
                    ),
                }}
              >
                {(field) => (
                  <field.TextField
                    type="tel"
                    label={tCustomers('form.phone_number')}
                    placeholder={tCustomers('form.phone_placeholder')}
                    autoComplete="tel"
                    required
                  />
                )}
              </form.AppField>
            </div>
          </div>

          {/* Country - Required */}
          <form.AppField
            name="countryCode"
            validators={{
              onBlur: z
                .string()
                .trim()
                .min(1, 'validation.required')
                .refine(
                  (v) => /^[A-Za-z]{2}$/.test(v),
                  'validation.invalid_country',
                ),
            }}
          >
            {(field) => (
              <CountrySelect
                value={field.state.value}
                onChange={(value: string) => {
                  field.handleChange(value)
                  setSelectedCountryCode(value)
                }}
                disabled={createMutation.isPending}
                required
              />
            )}
          </form.AppField>

          {/* Additional Details (Email, Gender) - Optional collapsible */}
          <form.Subscribe selector={(state) => state.values}>
            {(values) => (
              <AdditionalDetailsInputs
                email={values.email || ''}
                onEmailChange={(value: string) =>
                  form.setFieldValue('email', value)
                }
                gender={(values.gender as CustomerGender | '') || ''}
                onGenderChange={(value: CustomerGender | '') =>
                  form.setFieldValue('gender', value as string)
                }
                disabled={createMutation.isPending}
                defaultExpanded={false}
              />
            )}
          </form.Subscribe>

          {/* Social Media */}
          <form.Subscribe selector={(state) => state.values}>
            {(values) => (
              <SocialMediaInputs
                instagramUsername={values.instagramUsername || ''}
                onInstagramChange={(value: string) =>
                  form.setFieldValue('instagramUsername', value)
                }
                facebookUsername={values.facebookUsername || ''}
                onFacebookChange={(value: string) =>
                  form.setFieldValue('facebookUsername', value)
                }
                tiktokUsername={values.tiktokUsername || ''}
                onTiktokChange={(value: string) =>
                  form.setFieldValue('tiktokUsername', value)
                }
                snapchatUsername={values.snapchatUsername || ''}
                onSnapchatChange={(value: string) =>
                  form.setFieldValue('snapchatUsername', value)
                }
                xUsername={values.xUsername || ''}
                onXChange={(value: string) =>
                  form.setFieldValue('xUsername', value)
                }
                whatsappNumber={values.whatsappNumber || ''}
                onWhatsappChange={(value: string) =>
                  form.setFieldValue('whatsappNumber', value)
                }
                disabled={createMutation.isPending}
                defaultExpanded={false}
              />
            )}
          </form.Subscribe>
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
