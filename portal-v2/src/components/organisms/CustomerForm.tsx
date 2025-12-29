/**
 * CustomerForm Organism
 *
 * TanStack Form for creating/editing customers with:
 * - Zod validation (onBlur mode)
 * - Phone and country selects
 * - Social media handles
 * - RTL support
 */

import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { FormInput } from '../atoms/FormInput'
import { FormTextarea } from '../atoms/FormTextarea'
import { PhoneCodeSelect } from '../molecules/PhoneCodeSelect'
import { CountrySelect } from '../molecules/CountrySelect'
import { Button } from '../atoms/Button'
import type { CreateCustomerRequest, Customer } from '@/api/customer'

/**
 * Customer Form Validation Schema
 */
const CustomerFormSchema = z.object({
  fullName: z.string().min(1, 'common:validation.required'),
  email: z.string().email('common:validation.invalid_email').optional().or(z.literal('')),
  phonePrefix: z.string().optional().or(z.literal('')),
  phoneNumber: z.string().optional().or(z.literal('')),
  address: z.string().optional().or(z.literal('')),
  city: z.string().optional().or(z.literal('')),
  country: z.string().optional().or(z.literal('')),
  instagramHandle: z.string().optional().or(z.literal('')),
  facebookHandle: z.string().optional().or(z.literal('')),
  notes: z.string().optional().or(z.literal('')),
})

type CustomerFormValues = z.infer<typeof CustomerFormSchema>

export interface CustomerFormProps {
  /**
   * Customer to edit (omit for create mode)
   */
  customer?: Customer
  /**
   * Form submission handler
   */
  onSubmit: (values: CreateCustomerRequest) => Promise<void>
  /**
   * Cancel handler
   */
  onCancel?: () => void
  /**
   * Loading state
   */
  isSubmitting?: boolean
}

export function CustomerForm({
  customer,
  onSubmit,
  onCancel,
  isSubmitting = false,
}: CustomerFormProps) {
  const { t } = useTranslation(['common', 'errors'])

  // Initialize form with TanStack Form
  const form = useForm({
    defaultValues: {
      fullName: customer?.fullName ?? '',
      email: customer?.email ?? '',
      phonePrefix: customer?.phonePrefix ?? '',
      phoneNumber: customer?.phoneNumber ?? '',
      address: customer?.address ?? '',
      city: customer?.city ?? '',
      country: customer?.country ?? '',
      instagramHandle: customer?.instagramHandle ?? '',
      facebookHandle: customer?.facebookHandle ?? '',
      notes: customer?.notes ?? '',
    } as CustomerFormValues,
    onSubmit: async ({ value }) => {
      // Clean up empty strings to undefined for optional fields
      const cleanedValue = {
        ...value,
        email: value.email || undefined,
        phonePrefix: value.phonePrefix || undefined,
        phoneNumber: value.phoneNumber || undefined,
        address: value.address || undefined,
        city: value.city || undefined,
        country: value.country || undefined,
        instagramHandle: value.instagramHandle || undefined,
        facebookHandle: value.facebookHandle || undefined,
        notes: value.notes || undefined,
      }
      await onSubmit(cleanedValue)
    },
    validators: {
      onBlur: CustomerFormSchema,
    },
  })

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        e.stopPropagation()
        void form.handleSubmit()
      }}
      className="space-y-6"
    >
      {/* Basic Information */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">{t('common:customer.basic_info')}</h3>

        {/* Full Name */}
        <form.Field
          name="fullName"
          validators={{
            onBlur: CustomerFormSchema.shape.fullName,
          }}
        >
          {(field) => (
            <FormInput
              label={t('common:customer.full_name')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
              required
              placeholder={t('common:customer.full_name_placeholder')}
            />
          )}
        </form.Field>

        {/* Email */}
        <form.Field
          name="email"
          validators={{
            onBlur: CustomerFormSchema.shape.email,
          }}
        >
          {(field) => (
            <FormInput
              type="email"
              label={t('common:customer.email')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
              placeholder={t('common:customer.email_placeholder')}
            />
          )}
        </form.Field>

        {/* Phone */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="md:col-span-1">
            <form.Field
              name="phonePrefix"
              validators={{
                onBlur: CustomerFormSchema.shape.phonePrefix,
              }}
            >
              {(field) => (
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('common:customer.phone_code')}</span>
                  </label>
                  <PhoneCodeSelect
                    value={field.state.value || ''}
                    onChange={field.handleChange}
                    error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
                  />
                </div>
              )}
            </form.Field>
          </div>
          <div className="md:col-span-3">
            <form.Field
              name="phoneNumber"
              validators={{
                onBlur: CustomerFormSchema.shape.phoneNumber,
              }}
            >
              {(field) => (
                <FormInput
                  type="tel"
                  label={t('common:customer.phone_number')}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                  error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
                  placeholder={t('common:customer.phone_number_placeholder')}
                />
              )}
            </form.Field>
          </div>
        </div>
      </div>

      {/* Address Information */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">{t('common:customer.address_info')}</h3>

        {/* Address */}
        <form.Field
          name="address"
          validators={{
            onBlur: CustomerFormSchema.shape.address,
          }}
        >
          {(field) => (
            <FormInput
              label={t('common:customer.address')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
              placeholder={t('common:customer.address_placeholder')}
            />
          )}
        </form.Field>

        {/* City and Country */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <form.Field
            name="city"
            validators={{
              onBlur: CustomerFormSchema.shape.city,
            }}
          >
            {(field) => (
              <FormInput
                label={t('common:customer.city')}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                onBlur={field.handleBlur}
                error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
                placeholder={t('common:customer.city_placeholder')}
              />
            )}
          </form.Field>

          <form.Field
            name="country"
            validators={{
              onBlur: CustomerFormSchema.shape.country,
            }}
          >
            {(field) => (
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('common:customer.country')}</span>
                </label>
                <CountrySelect
                  value={field.state.value || ''}
                  onChange={field.handleChange}
                  error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
                />
              </div>
            )}
          </form.Field>
        </div>
      </div>

      {/* Social Media */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">{t('common:customer.social_media')}</h3>

        <form.Field
          name="instagramHandle"
          validators={{
            onBlur: CustomerFormSchema.shape.instagramHandle,
          }}
        >
          {(field) => (
            <FormInput
              label={t('common:customer.instagram')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
              placeholder="@username"
            />
          )}
        </form.Field>

        <form.Field
          name="facebookHandle"
          validators={{
            onBlur: CustomerFormSchema.shape.facebookHandle,
          }}
        >
          {(field) => (
            <FormInput
              label={t('common:customer.facebook')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
              placeholder="@username"
            />
          )}
        </form.Field>
      </div>

      {/* Notes */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">{t('common:customer.notes')}</h3>

        <form.Field
          name="notes"
          validators={{
            onBlur: CustomerFormSchema.shape.notes,
          }}
        >
          {(field) => (
            <FormTextarea
              label={t('common:customer.notes')}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              onBlur={field.handleBlur}
              error={field.state.meta.errors.length > 0 ? String(field.state.meta.errors[0]) : undefined}
              placeholder={t('common:customer.notes_placeholder')}
              rows={4}
            />
          )}
        </form.Field>
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 justify-end pt-4 border-t border-base-300">
        {onCancel && (
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={isSubmitting}
          >
            {t('common:actions.cancel')}
          </Button>
        )}
        <Button
          type="submit"
          variant="primary"
          loading={isSubmitting}
        >
          {customer ? t('common:actions.save_changes') : t('common:actions.create')}
        </Button>
      </div>
    </form>
  )
}
