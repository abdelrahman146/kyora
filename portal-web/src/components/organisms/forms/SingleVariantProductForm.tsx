import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { useParams } from '@tanstack/react-router'
import toast from 'react-hot-toast'

import type {
  AssetReference,
  CreateProductWithVariantsRequest,
} from '@/api/inventory'
import { useKyoraForm } from '@/lib/form'
import { BusinessContext } from '@/lib/form/components/FileUploadField'
import { singleVariantProductSchema } from '@/schemas/inventory'
import {
  useCategoriesQuery,
  useCreateProductWithVariantsMutation,
} from '@/api/inventory'
import { queryKeys } from '@/lib/queryKeys'
import { translateErrorAsync } from '@/lib/translateError'

interface SingleVariantProductFormProps {
  onSuccess: () => void
  onCancel: () => void
}

export function SingleVariantProductForm({
  onSuccess,
  onCancel,
}: SingleVariantProductFormProps) {
  const { t } = useTranslation()
  const { businessDescriptor } = useParams({ strict: false })
  const queryClient = useQueryClient()

  const categoriesQuery = useCategoriesQuery(businessDescriptor!)
  const categories = categoriesQuery.data ?? []

  const createMutation = useCreateProductWithVariantsMutation(
    businessDescriptor!,
    {
      onSuccess: () => {
        toast.success(t('product_created', { ns: 'inventory' }))
        queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all })
        onSuccess()
      },
      onError: async (error) => {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      },
    },
  )

  const form = useKyoraForm({
    defaultValues: {
      name: '',
      description: '',
      photos: [] as Array<AssetReference>,
      categoryId: '',
      sku: '',
      costPrice: '',
      salePrice: '',
      stockQuantity: 0,
      stockQuantityAlert: 0,
    },
    onSubmit: async ({ value }) => {
      const payload: CreateProductWithVariantsRequest = {
        product: {
          name: value.name,
          description: value.description || undefined,
          photos: value.photos,
          categoryId: value.categoryId,
        },
        variants: [
          {
            code: 'STANDARD',
            sku: value.sku || undefined,
            photos: value.photos, // Use same photos as product
            costPrice: value.costPrice,
            salePrice: value.salePrice,
            stockQuantity:
              typeof value.stockQuantity === 'string'
                ? parseInt(value.stockQuantity, 10)
                : value.stockQuantity,
            stockQuantityAlert:
              typeof value.stockQuantityAlert === 'string'
                ? parseInt(value.stockQuantityAlert, 10)
                : value.stockQuantityAlert,
          },
        ],
      }

      await createMutation.mutateAsync(payload)
    },
  })

  const categoryOptions = categories.map((cat) => ({
    value: cat.id,
    label: cat.name,
  }))

  return (
    <BusinessContext.Provider
      value={{ businessDescriptor: businessDescriptor! }}
    >
      <form.AppForm>
        <form.FormRoot className="flex flex-col h-full">
          <div className="flex-1 overflow-y-auto p-4 space-y-6">
            <form.FormError />

            <div className="space-y-4">
              <form.AppField
                name="name"
                validators={{
                  onBlur: singleVariantProductSchema.shape.name,
                }}
              >
                {(field) => (
                  <field.TextField
                    label={t('product_name', { ns: 'inventory' })}
                    placeholder={t('product_name_placeholder', {
                      ns: 'inventory',
                    })}
                    required
                  />
                )}
              </form.AppField>

              <form.AppField name="description">
                {(field) => (
                  <field.TextareaField
                    label={t('description', { ns: 'inventory' })}
                    placeholder={t('description_placeholder', {
                      ns: 'inventory',
                    })}
                    rows={3}
                  />
                )}
              </form.AppField>

              <form.AppField
                name="categoryId"
                validators={{
                  onBlur: singleVariantProductSchema.shape.categoryId,
                }}
              >
                {(field) => (
                  <field.SelectField
                    label={t('category', { ns: 'inventory' })}
                    placeholder={t('select_category', { ns: 'inventory' })}
                    options={categoryOptions}
                    required
                  />
                )}
              </form.AppField>

              <form.AppField name="photos">
                {(field) => (
                  <field.FileUploadField
                    label={t('product_photos', { ns: 'inventory' })}
                    hint={t('product_photos_hint', { ns: 'inventory' })}
                    accept="image/*"
                    maxFiles={10}
                    multiple
                  />
                )}
              </form.AppField>

              <form.AppField name="sku">
                {(field) => (
                  <field.TextField
                    label={t('sku', { ns: 'inventory' })}
                    placeholder={t('sku_placeholder', { ns: 'inventory' })}
                    hint={t('sku_hint', { ns: 'inventory' })}
                  />
                )}
              </form.AppField>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <form.AppField
                  name="costPrice"
                  validators={{
                    onBlur: singleVariantProductSchema.shape.costPrice,
                  }}
                >
                  {(field) => (
                    <field.TextField
                      type="text"
                      label={t('cost_price', { ns: 'inventory' })}
                      placeholder="0.00"
                      required
                    />
                  )}
                </form.AppField>

                <form.AppField
                  name="salePrice"
                  validators={{
                    onBlur: singleVariantProductSchema.shape.salePrice,
                  }}
                >
                  {(field) => (
                    <field.TextField
                      type="text"
                      label={t('sale_price', { ns: 'inventory' })}
                      placeholder="0.00"
                      required
                    />
                  )}
                </form.AppField>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <form.AppField name="stockQuantity">
                  {(field) => (
                    <field.TextField
                      type="text"
                      label={t('stock_quantity', { ns: 'inventory' })}
                      placeholder="0"
                      required
                    />
                  )}
                </form.AppField>

                <form.AppField name="stockQuantityAlert">
                  {(field) => (
                    <field.TextField
                      type="text"
                      label={t('stock_alert', { ns: 'inventory' })}
                      placeholder="0"
                      hint={t('stock_alert_hint', { ns: 'inventory' })}
                    />
                  )}
                </form.AppField>
              </div>
            </div>
          </div>

          <div className="flex gap-2 p-4 border-t border-base-300 bg-base-100">
            <button
              type="button"
              onClick={onCancel}
              className="btn btn-ghost flex-1"
              disabled={createMutation.isPending}
            >
              {t('common.cancel')}
            </button>
            <form.SubmitButton variant="primary" className="flex-1">
              {t('create_product', { ns: 'inventory' })}
            </form.SubmitButton>
          </div>
        </form.FormRoot>
      </form.AppForm>
    </BusinessContext.Provider>
  )
}
