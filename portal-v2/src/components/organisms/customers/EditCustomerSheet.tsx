/**
 * EditCustomerSheet Component
 *
 * Bottom sheet form for editing existing customers.
 * NOTE: Unlike AddCustomerSheet, this does NOT auto-link country to phone code
 * because user may have different phone code than their country.
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
import { useUpdateCustomerMutation } from '@/api/customer'
import { buildE164Phone } from '@/lib/phone'
import { showErrorToast, showSuccessToast } from '@/lib/toast'
import { translateErrorAsync } from '@/lib/translateError'

export interface EditCustomerSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  customer: Customer
  onUpdated?: (customer: Customer) => void | Promise<void>
}

const editCustomerSchema = z
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

export type EditCustomerFormValues = z.infer<typeof editCustomerSchema>

function getDefaultValues(customer: Customer): EditCustomerFormValues {
  return {
    name: customer.name,
    email: customer.email ?? '',
    gender: customer.gender,
    countryCode: customer.countryCode,
    phoneCode: customer.phoneCode ?? '',
    phoneNumber: customer.phoneNumber ?? '',
    instagramUsername: customer.instagramUsername ?? '',
    facebookUsername: customer.facebookUsername ?? '',
    tiktokUsername: customer.tiktokUsername ?? '',
    snapchatUsername: customer.snapchatUsername ?? '',
    xUsername: customer.xUsername ?? '',
    whatsappNumber: customer.whatsappNumber ?? '',
  }
}

export function EditCustomerSheet({
  isOpen,
  onClose,
  businessDescriptor,
  customer,
  onUpdated,
}: EditCustomerSheetProps) {
  const { t } = useTranslation()
  const { t: tErrors } = useTranslation('errors')
  const formId = useId()

  // Local state for social media values
  const [socialMediaValues, setSocialMediaValues] = useState(() => ({
    instagramUsername: customer.instagramUsername ?? '',
    facebookUsername: customer.facebookUsername ?? '',
    tiktokUsername: customer.tiktokUsername ?? '',
    snapchatUsername: customer.snapchatUsername ?? '',
    xUsername: customer.xUsername ?? '',
    whatsappNumber: customer.whatsappNumber ?? '',
  }))

  // Update customer mutation
  const updateMutation = useUpdateCustomerMutation(businessDescriptor, {
    onSuccess: async (updated) => {
      showSuccessToast(t('customers.update_success'))
      if (onUpdated) {
        await onUpdated(updated)
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
    defaultValues: getDefaultValues(customer),
    validators: {
      onBlur: editCustomerSchema,
    },
    onSubmit: async ({ value }) => {
      const phoneCode = value.phoneCode.trim()
      const phoneNumber = value.phoneNumber.trim()

      const normalizedPhone =
        phoneNumber !== '' && phoneCode !== ''
          ? buildE164Phone(phoneCode, phoneNumber)
          : undefined

      await updateMutation.mutateAsync({
        customerId: customer.id,
        data: {
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
        },
      })
    },
  })

  // Note: In TanStack Form v0.x, accessing form.state creates subscriptions.
  // The best practice is to use form.Subscribe component with selectors,
  // but for pragmatic inline state access, form.state is acceptable for boolean flags
  // when the component already has other reasons to re-render (e.g., mutation state)
  const isSubmitting = form.state.isSubmitting || updateMutation.isPending
  const isDirty = form.state.isDirty

  // Reset form and local state when sheet opens or customer changes
  useEffect(() => {
    if (isOpen) {
      form.reset()
      const defaults = getDefaultValues(customer)
      Object.entries(defaults).forEach(([key, value]) => {
        form.setFieldValue(key as keyof EditCustomerFormValues, value)
      })
      setSocialMediaValues({
        instagramUsername: customer.instagramUsername ?? '',
        facebookUsername: customer.facebookUsername ?? '',
        tiktokUsername: customer.tiktokUsername ?? '',
        snapchatUsername: customer.snapchatUsername ?? '',
        xUsername: customer.xUsername ?? '',
        whatsappNumber: customer.whatsappNumber ?? '',
      })
    }
  }, [isOpen, customer, form])

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
        form={`edit-customer-form-${formId}`}
        className="btn btn-primary flex-1"
        disabled={isSubmitting || !isDirty}
        aria-disabled={isSubmitting || !isDirty}
      >
        {isSubmitting
          ? t('customers.update_submitting')
          : t('customers.update_submit')}
      </button>
    </div>
  )

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={safeClose}
      title={t('customers.edit_title')}
      footer={footer}
      side="end"
      size="md"
      closeOnOverlayClick={!isSubmitting}
      closeOnEscape={!isSubmitting}
      contentClassName="space-y-4"
      ariaLabel={t('customers.edit_title')}
    >
      <form
        id={`edit-customer-form-${formId}`}
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
            onBlur: editCustomerSchema.shape.name,
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
            onBlur: editCustomerSchema.shape.email,
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
              onBlur: editCustomerSchema.shape.countryCode,
            }}
          >
            {(field) => (
              <CountrySelect
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
                required
              />
            )}
          </form.Field>

          <form.Field
            name="gender"
            validators={{
              onBlur: editCustomerSchema.shape.gender,
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
              onBlur: editCustomerSchema.shape.phoneCode,
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
                onBlur: editCustomerSchema.shape.phoneNumber,
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
          defaultExpanded={true}
        />
      </form>
    </BottomSheet>
  )
}
