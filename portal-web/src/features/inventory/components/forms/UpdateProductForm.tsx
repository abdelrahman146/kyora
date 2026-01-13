import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'

import type { Product } from '@/api/inventory'
import { useKyoraForm } from '@/lib/form'
import { BusinessContext } from '@/lib/form/components/FileUploadField'
import { singleVariantProductSchema } from '@/schemas/inventory'
import {
  useUpdateProductMutation,
  useUpdateVariantMutation,
} from '@/api/inventory'
import { queryKeys } from '@/lib/queryKeys'
import { getSelectedBusiness } from '@/stores/businessStore'

interface UpdateProductFormProps {
  product: Product
  businessDescriptor: string
  onSuccess: () => void
  onCancel: () => void
}

export function UpdateProductForm({
  product,
  businessDescriptor,
  onSuccess,
  onCancel,
}: UpdateProductFormProps) {
  const { t: tInventory } = useTranslation('inventory')
  const { t: tCommon } = useTranslation('common')
  const queryClient = useQueryClient()
  const selectedBusiness = getSelectedBusiness()
  const currencyCode = selectedBusiness?.currency ?? 'USD'

  const variant = product.variants?.[0]

  const updateProductMutation = useUpdateProductMutation(
    businessDescriptor,
    product.id,
  )

  const updateVariantMutation = useUpdateVariantMutation(
    businessDescriptor,
    variant?.id || '',
  )

  const form = useKyoraForm({
    defaultValues: {
      name: product.name || '',
      description: product.description || '',
      photos: product.photos,
      categoryId: product.categoryId || '',
      sku: variant?.sku || '',
      costPrice: variant?.costPrice || '',
      salePrice: variant?.salePrice || '',
      stockQuantity: variant?.stockQuantity || 0,
      stockQuantityAlert: variant?.stockQuantityAlert || 0,
    },
    onSubmit: async ({ value }) => {
      try {
        await updateProductMutation.mutateAsync({
          name: value.name,
          description: value.description || undefined,
          photos: value.photos,
          categoryId: value.categoryId,
        })

        if (variant) {
          await updateVariantMutation.mutateAsync({
            code: variant.code,
            sku: value.sku || undefined,
            photos: value.photos,
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
          })
        }

        toast.success(tInventory('product_updated'))
        void queryClient.invalidateQueries({
          queryKey: queryKeys.inventory.all,
        })
        onSuccess()
      } catch {
        // Global QueryClient error handler shows the error toast.
      }
    },
  })

  const isLoading =
    updateProductMutation.isPending || updateVariantMutation.isPending

  return (
    <BusinessContext.Provider value={{ businessDescriptor }}>
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
                    label={tInventory('product_name')}
                    placeholder={tInventory('product_name_placeholder', {
                      ns: 'inventory',
                    })}
                    required
                  />
                )}
              </form.AppField>

              <form.AppField name="description">
                {(field) => (
                  <field.TextareaField
                    label={tInventory('description')}
                    placeholder={tInventory('description_placeholder', {
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
                  <field.CategorySelectField
                    businessDescriptor={businessDescriptor}
                    label={tInventory('category')}
                    placeholder={tInventory('select_category')}
                    required
                  />
                )}
              </form.AppField>

              <form.AppField name="photos">
                {(field) => (
                  <field.FileUploadField
                    label={tInventory('product_photos')}
                    hint={tInventory('product_photos_hint')}
                    accept="image/*"
                    maxFiles={10}
                    multiple
                  />
                )}
              </form.AppField>

              <form.AppField name="sku">
                {(field) => (
                  <field.TextField
                    label={tInventory('sku')}
                    placeholder={tInventory('sku_placeholder')}
                    hint={tInventory('sku_hint')}
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
                    <field.PriceField
                      label={tInventory('cost_price')}
                      placeholder="0.00"
                      currencyCode={currencyCode}
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
                    <field.PriceField
                      label={tInventory('sale_price')}
                      placeholder="0.00"
                      currencyCode={currencyCode}
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
                      label={tInventory('stock_quantity')}
                      placeholder="0"
                      required
                    />
                  )}
                </form.AppField>

                <form.AppField name="stockQuantityAlert">
                  {(field) => (
                    <field.TextField
                      type="text"
                      label={tInventory('stock_alert')}
                      placeholder="0"
                      hint={tInventory('stock_alert_hint')}
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
              disabled={isLoading}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton variant="primary" className="flex-1">
              {tInventory('update_product')}
            </form.SubmitButton>
          </div>
        </form.FormRoot>
      </form.AppForm>
    </BusinessContext.Provider>
  )
}
