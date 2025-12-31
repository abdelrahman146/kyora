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
import { z } from 'zod'

import { BottomSheet } from '../../molecules/BottomSheet'
import { CountrySelect } from '../../molecules/CountrySelect'
import { PhoneCodeSelect } from '../../molecules/PhoneCodeSelect'
import { SocialMediaInputs } from '../../molecules/SocialMediaInputs'
import type { Customer, CustomerGender } from '@/api/customer'
import { FormSelect } from '@/components'
import { useCreateCustomerMutation } from '@/api/customer'
import { useCountriesQuery } from '@/api/metadata'
import { useKyoraForm } from '@/lib/form'
import { buildE164Phone } from '@/lib/phone'
import { showErrorToast, showSuccessToast } from '@/lib/toast'

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
  const formId = useId()

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
    onError: (error) => {
      showErrorToast(error.message)
    },
  })

  // TanStack Form setup with useKyoraForm
  const form = useKyoraForm({
    defaultValues: getDefaultValues(businessCountryCode),
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

  // Auto-link country to phone code when country changes
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
        title={t('customers.create_title')}
        footer={
          <div className="flex gap-2">
            <button
              type="button"
              className="btn btn-ghost flex-1"
              onClick={safeClose}
              disabled={createMutation.isPending}
              aria-disabled={createMutation.isPending}
            >
              {t('common.cancel')}
            </button>
            <form.SubmitButton
              form={`add-customer-form-${formId}`}
              variant="primary"
              className="flex-1"
            >
              {createMutation.isPending
                ? t('customers.create_submitting')
                : t('customers.create_submit')}
            </form.SubmitButton>
          </div>
        }
        side="end"
        size="md"
        closeOnOverlayClick={!createMutation.isPending}
        closeOnEscape={!createMutation.isPending}
        contentClassName="space-y-4"
        ariaLabel={t('customers.create_title')}
      >
        <form.FormRoot
          id={`add-customer-form-${formId}`}
          className="space-y-4"
          aria-busy={createMutation.isPending}
        >
        <form.AppField
          name="name"
          validators={{
            onBlur: z.string().trim().min(1, 'validation.required'),
          }}
        >
          {(field) => (
            <field.TextField
              label={t('customers.form.name')}
              placeholder={t('customers.form.name_placeholder')}
              autoComplete="name"
              required
            />
          )}
        </form.AppField>

        <form.AppField
          name="email"
          validators={{
            onBlur: z.string().trim().min(1, 'validation.required').email('validation.invalid_email'),
          }}
        >
          {(field) => (
            <field.TextField
              type="email"
              label={t('customers.form.email')}
              placeholder={t('customers.form.email_placeholder')}
              autoComplete="email"
              required
            />
          )}
        </form.AppField>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <form.AppField
            name="countryCode"
            validators={{
              onBlur: z.string().trim().min(1, 'validation.required').refine((v) => /^[A-Za-z]{2}$/.test(v), 'validation.invalid_country'),
            }}
          >
            {(field: any) => (
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

          <form.AppField
            name="gender"
            validators={{
              onBlur: z.enum(['male', 'female', 'other'], { message: 'validation.required' }),
            }}
          >
            {(field: any) => (
              <FormSelect<CustomerGender>
                label={t('customers.form.gender')}
                options={[
                  { value: 'male', label: t('customers.form.gender_male') },
                  { value: 'female', label: t('customers.form.gender_female') },
                  { value: 'other', label: t('customers.form.gender_other') },
                ]}
                value={field.state.value}
                onChange={(value: CustomerGender | Array<CustomerGender>) => {
                  // FormSelect can return array for multiSelect, but we use single select
                  const singleValue = Array.isArray(value) ? value[0] : value
                  field.handleChange(singleValue)
                }}
                required
                disabled={createMutation.isPending}
                placeholder={t('customers.form.select_gender')}
              />
            )}
          </form.AppField>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <form.AppField
            name="phoneCode"
            validators={{
              onBlur: z.string().trim().refine((v) => v === '' || /^\+?\d{1,4}$/.test(v), 'validation.invalid_phone_code'),
            }}
          >
            {(field: any) => (
              <PhoneCodeSelect
                value={field.state.value}
                onChange={(value: string) => field.handleChange(value)}
                disabled={createMutation.isPending}
              />
            )}
          </form.AppField>

          <div className="sm:col-span-2">
            <form.AppField
              name="phoneNumber"
              validators={{
                onBlur: z.string().trim().refine((v) => v === '' || /^[0-9\-\s()]{6,20}$/.test(v), 'validation.invalid_phone'),
              }}
            >
              {(field) => (
                <field.TextField
                  type="tel"
                  label={t('customers.form.phone_number')}
                  placeholder={t('customers.form.phone_placeholder')}
                  autoComplete="tel"
                />
              )}
            </form.AppField>
          </div>
        </div>

        <form.Subscribe selector={(state: any) => state.values}>
          {(values: any) => (
            <SocialMediaInputs
              instagramUsername={values.instagramUsername}
              onInstagramChange={(value: string) =>
                form.setFieldValue('instagramUsername', value)
              }
              facebookUsername={values.facebookUsername}
              onFacebookChange={(value: string) =>
                form.setFieldValue('facebookUsername', value)
              }
              tiktokUsername={values.tiktokUsername}
              onTiktokChange={(value: string) =>
                form.setFieldValue('tiktokUsername', value)
              }
              snapchatUsername={values.snapchatUsername}
              onSnapchatChange={(value: string) =>
                form.setFieldValue('snapchatUsername', value)
              }
              xUsername={values.xUsername}
              onXChange={(value: string) =>
                form.setFieldValue('xUsername', value)
              }
              whatsappNumber={values.whatsappNumber}
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
