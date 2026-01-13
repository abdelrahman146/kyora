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
import { useCreateProductWithVariantsMutation } from '@/api/inventory'
import { queryKeys } from '@/lib/queryKeys'
import { getSelectedBusiness } from '@/stores/businessStore'

interface SingleVariantProductFormProps {
  onSuccess: () => void
  onCancel: () => void
}

export function SingleVariantProductForm({
  onSuccess,
  onCancel,
}: SingleVariantProductFormProps) {
  const { t: tInventory } = useTranslation('inventory')
  const { t: tCommon } = useTranslation('common')
  const { businessDescriptor } = useParams({ strict: false })
  const queryClient = useQueryClient()
  const selectedBusiness = getSelectedBusiness()
  const currencyCode = selectedBusiness?.currency ?? 'USD'

  const createMutation = useCreateProductWithVariantsMutation(
    businessDescriptor!,
    {
      onSuccess: () => {
        toast.success(tInventory('product_created'))
        queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all })
        onSuccess()
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
          },
        ],
      }

      try {
        await createMutation.mutateAsync(payload)
      } catch {
        // Global QueryClient mutation handler will toast.
        return
      }
    },
  })

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
                    businessDescriptor={businessDescriptor!}
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
              disabled={createMutation.isPending}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton variant="primary" className="flex-1">
              {tInventory('create_product')}
            </form.SubmitButton>
          </div>
        </form.FormRoot>
      </form.AppForm>
    </BusinessContext.Provider>
  )
}
