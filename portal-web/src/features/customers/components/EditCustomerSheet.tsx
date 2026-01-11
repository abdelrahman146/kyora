import { useEffect, useId, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { CountrySelect } from './CountrySelect'
import { PhoneCodeSelect } from './PhoneCodeSelect'
import { SocialMediaInputs } from './SocialMediaInputs'
import type { Customer, CustomerGender } from '@/api/customer'
import { useUpdateCustomerMutation } from '@/api/customer'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { FormSelect } from '@/components'
import { useKyoraForm } from '@/lib/form'
import { buildE164Phone } from '@/lib/phone'
import { showErrorToast, showSuccessToast } from '@/lib/toast'

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
  const { t: tCustomers } = useTranslation('customers')
  const { t: tCommon } = useTranslation('common')
  const formId = useId()

  const updateMutation = useUpdateCustomerMutation(businessDescriptor, {
    onSuccess: async (updated) => {
      showSuccessToast(tCustomers('update_success'))
      if (onUpdated) {
        await onUpdated(updated)
      }
      onClose()
    },
    onError: (error) => {
      showErrorToast(error.message)
    },
  })

  const form = useKyoraForm({
    defaultValues: getDefaultValues(customer),
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

  const [isDirty, setIsDirty] = useState(false)
  useEffect(() => {
    const unsubscribe = form.store.subscribe(() => {
      setIsDirty(form.store.state.isDirty)
    })
    return unsubscribe
  }, [form])

  useEffect(() => {
    if (isOpen) {
      form.reset()
      const defaults = getDefaultValues(customer)
      Object.entries(defaults).forEach(([key, value]) => {
        form.setFieldValue(key as keyof EditCustomerFormValues, value)
      })
    }
  }, [isOpen, customer, form])

  const safeClose = () => {
    if (updateMutation.isPending) return
    onClose()
  }

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={safeClose}
        title={tCustomers('edit_title')}
        footer={
          <div className="flex gap-2">
            <button
              type="button"
              className="btn btn-ghost flex-1"
              onClick={safeClose}
              disabled={updateMutation.isPending}
              aria-disabled={updateMutation.isPending}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton
              form={`edit-customer-form-${formId}`}
              variant="primary"
              className="flex-1"
              disabled={!isDirty}
            >
              {updateMutation.isPending
                ? tCustomers('update_submitting')
                : tCustomers('update_submit')}
            </form.SubmitButton>
          </div>
        }
        side="end"
        size="md"
        closeOnOverlayClick={!updateMutation.isPending}
        closeOnEscape={!updateMutation.isPending}
        contentClassName="space-y-4"
        ariaLabel={tCustomers('edit_title')}
      >
        <form.FormRoot
          id={`edit-customer-form-${formId}`}
          className="space-y-4"
          aria-busy={updateMutation.isPending}
        >
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

          <form.AppField
            name="email"
            validators={{
              onBlur: z
                .string()
                .trim()
                .min(1, 'validation.required')
                .email('validation.invalid_email'),
            }}
          >
            {(field) => (
              <field.TextField
                type="email"
                label={tCustomers('form.email')}
                placeholder={tCustomers('form.email_placeholder')}
                autoComplete="email"
                required
              />
            )}
          </form.AppField>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
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
              {(field: any) => (
                <CountrySelect
                  value={field.state.value}
                  onChange={(value: string) => field.handleChange(value)}
                  disabled={updateMutation.isPending}
                  required
                />
              )}
            </form.AppField>

            <form.AppField
              name="gender"
              validators={{
                onBlur: z.enum(['male', 'female', 'other'], {
                  message: 'validation.required',
                }),
              }}
            >
              {(field: any) => (
                <FormSelect<CustomerGender>
                  label={tCustomers('form.gender')}
                  options={[
                    { value: 'male', label: tCustomers('form.gender_male') },
                    {
                      value: 'female',
                      label: tCustomers('form.gender_female'),
                    },
                    { value: 'other', label: tCustomers('form.gender_other') },
                  ]}
                  value={field.state.value}
                  onChange={(value: CustomerGender | Array<CustomerGender>) => {
                    const singleValue = Array.isArray(value) ? value[0] : value
                    field.handleChange(singleValue)
                  }}
                  required
                  disabled={updateMutation.isPending}
                  placeholder={tCustomers('form.select_gender')}
                />
              )}
            </form.AppField>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <form.AppField
              name="phoneCode"
              validators={{
                onBlur: z
                  .string()
                  .trim()
                  .refine(
                    (v) => v === '' || /^\+?\d{1,4}$/.test(v),
                    'validation.invalid_phone_code',
                  ),
              }}
            >
              {(field: any) => (
                <PhoneCodeSelect
                  value={field.state.value}
                  onChange={(value: string) => field.handleChange(value)}
                  disabled={updateMutation.isPending}
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
                    .refine(
                      (v) => v === '' || /^[0-9\-\s()]{6,20}$/.test(v),
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
                disabled={updateMutation.isPending}
                defaultExpanded={true}
              />
            )}
          </form.Subscribe>
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
