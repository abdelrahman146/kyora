/**
 * CreateAssetSheet Component
 *
 * Bottom sheet form for creating fixed assets.
 *
 * Features:
 * - Name field
 * - Type selector (software, equipment, vehicle, furniture, other)
 * - Value field with currency
 * - Purchase date picker
 * - Optional note field
 * - Validation with translated error messages
 */

import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { z } from 'zod'

import type { AssetType } from '@/api/accounting'
import { assetTypeEnum, useCreateAssetMutation } from '@/api/accounting'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { Button } from '@/components/atoms/Button'
import { useKyoraForm } from '@/lib/form'
import { getSelectedBusiness } from '@/stores/businessStore'

interface CreateAssetSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  onCreated?: () => void | Promise<void>
}

/**
 * Asset type options for the select field
 */
const assetTypeOptions: Array<{ value: AssetType; labelKey: string }> = [
  { value: 'software', labelKey: 'asset_type.software' },
  { value: 'equipment', labelKey: 'asset_type.equipment' },
  { value: 'vehicle', labelKey: 'asset_type.vehicle' },
  { value: 'furniture', labelKey: 'asset_type.furniture' },
  { value: 'other', labelKey: 'asset_type.other' },
]

export function CreateAssetSheet({
  isOpen,
  onClose,
  businessDescriptor,
  onCreated,
}: CreateAssetSheetProps) {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')

  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'

  const createAssetMutation = useCreateAssetMutation(businessDescriptor)

  const form = useKyoraForm({
    defaultValues: {
      name: '',
      type: 'equipment' as AssetType,
      value: '',
      purchaseDate: new Date(),
      note: '',
    },
    onSubmit: async ({ value }) => {
      // Let global error handler catch errors - don't manually toast them
      await createAssetMutation.mutateAsync({
        name: value.name,
        type: value.type,
        value: value.value,
        purchasedAt: value.purchaseDate.toISOString(),
        note: value.note || undefined,
      })
      toast.success(t('toast.asset_created'))

      // Reset and close
      form.reset()
      await onCreated?.()
      onClose()
    },
  })

  const isSubmitting = createAssetMutation.isPending

  const handleClose = () => {
    if (isSubmitting) return
    form.reset()
    onClose()
  }

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={handleClose}
        title={t('actions.add_asset')}
        closeOnOverlayClick={!isSubmitting}
        closeOnEscape={!isSubmitting}
        footer={
          <div className="flex gap-2">
            <Button
              type="button"
              variant="ghost"
              className="flex-1"
              onClick={handleClose}
              disabled={isSubmitting}
            >
              {tCommon('cancel')}
            </Button>
            <form.SubmitButton
              form="create-asset-form"
              variant="primary"
              className="flex-1"
              disabled={isSubmitting}
            >
              {isSubmitting ? tCommon('saving') : t('actions.save_asset')}
            </form.SubmitButton>
          </div>
        }
      >
        <form.FormRoot id="create-asset-form" className="space-y-4">
          {/* Name Field */}
          <form.AppField
            name="name"
            validators={{
              onChange: z
                .string()
                .min(1, 'accounting:validation.name_required'),
            }}
          >
            {(field) => (
              <field.TextField
                label={t('form.name')}
                placeholder={t('form.name_placeholder')}
                required
              />
            )}
          </form.AppField>

          {/* Type Field */}
          <form.AppField
            name="type"
            validators={{
              onChange: assetTypeEnum,
            }}
          >
            {(field) => (
              <field.SelectField
                label={t('form.type')}
                required
                options={assetTypeOptions.map((option) => ({
                  value: option.value,
                  label: t(option.labelKey),
                }))}
              />
            )}
          </form.AppField>

          {/* Value Field */}
          <form.AppField
            name="value"
            validators={{
              onChange: z
                .string()
                .min(1, 'accounting:validation.value_required')
                .refine(
                  (val) => {
                    const num = parseFloat(val)
                    return !isNaN(num) && num > 0
                  },
                  { message: 'accounting:validation.value_positive' },
                ),
            }}
          >
            {(field) => (
              <field.PriceField
                label={t('form.value')}
                placeholder={t('form.value_placeholder')}
                currencyCode={currency}
                required
              />
            )}
          </form.AppField>

          {/* Purchase Date Field */}
          <form.AppField
            name="purchaseDate"
            validators={{
              onChange: z
                .date({
                  required_error: 'accounting:validation.date_required',
                  invalid_type_error: 'accounting:validation.date_required',
                })
                .max(new Date(), 'accounting:validation.date_not_future'),
            }}
          >
            {(field) => (
              <field.DateField
                label={t('form.purchase_date')}
                maxDate={new Date()}
                required
              />
            )}
          </form.AppField>

          {/* Note Field */}
          <form.AppField name="note">
            {(field) => (
              <field.TextareaField
                label={t('form.note')}
                placeholder={t('form.note_placeholder')}
              />
            )}
          </form.AppField>
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
