import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { Plus, Trash2 } from 'lucide-react'

import type { AssetReference, Product } from '@/api/inventory'
import { useKyoraForm } from '@/lib/form'
import { BusinessContext } from '@/lib/form/components/FileUploadField'
import { multiVariantProductSchema } from '@/schemas/inventory'
import {
  inventoryApi,
  useCreateVariantMutation,
  useUpdateProductMutation,
} from '@/api/inventory'
import { queryKeys } from '@/lib/queryKeys'
import { translateErrorAsync } from '@/lib/translateError'
import { getSelectedBusiness } from '@/stores/businessStore'

interface UpdateProductWithVariantsFormProps {
  product: Product
  businessDescriptor: string
  onSuccess: () => void
  onCancel: () => void
}

type VariantFormValue = {
  id?: string
  code: string
  sku?: string
  photos: Array<AssetReference>
  costPrice: string
  salePrice: string
  stockQuantity: number | string
  stockQuantityAlert: number | string
}

export function UpdateProductWithVariantsForm({
  product,
  businessDescriptor,
  onSuccess,
  onCancel,
}: UpdateProductWithVariantsFormProps) {
  const { t: tInventory } = useTranslation('inventory')
  const { t: tCommon } = useTranslation('common')
  const queryClient = useQueryClient()
  const selectedBusiness = getSelectedBusiness()
  const currencyCode = selectedBusiness?.currency ?? 'USD'
  const [variantsToDelete, setVariantsToDelete] = useState<Array<string>>([])

  const updateProductMutation = useUpdateProductMutation(
    businessDescriptor,
    product.id,
    {
      onError: async (error) => {
        const message = await translateErrorAsync(error, tInventory)
        toast.error(message)
      },
    },
  )

  const createVariantMutation = useCreateVariantMutation(
    businessDescriptor,
    product.id,
    {
      onError: async (error) => {
        const message = await translateErrorAsync(error, tInventory)
        toast.error(message)
      },
    },
  )

  const form = useKyoraForm({
    defaultValues: {
      name: product.name || '',
      description: product.description || '',
      photos: product.photos,
      categoryId: product.categoryId || '',
      variants: product.variants?.map((v) => ({
        id: v.id,
        code: v.code || '',
        sku: v.sku || '',
        photos: v.photos,
        costPrice: v.costPrice || '',
        salePrice: v.salePrice || '',
        stockQuantity: v.stockQuantity || 0,
        stockQuantityAlert: v.stockQuantityAlert || 0,
      })) || [
        {
          code: '',
          sku: '',
          photos: [] as Array<AssetReference>,
          costPrice: '',
          salePrice: '',
          stockQuantity: 0,
          stockQuantityAlert: 0,
        },
      ],
    },
    onSubmit: async ({ value }) => {
      try {
        await updateProductMutation.mutateAsync({
          name: value.name,
          description: value.description || undefined,
          photos: value.photos,
          categoryId: value.categoryId,
        })

        const variants = value.variants as Array<VariantFormValue>

        const existingVariants = variants.filter(
          (variant): variant is VariantFormValue & { id: string } =>
            Boolean(variant.id),
        )

        const newVariants = variants.filter((variant) => !variant.id)

        for (const variant of existingVariants) {
          await inventoryApi.updateVariant(businessDescriptor, variant.id, {
            code: variant.code,
            sku: variant.sku || undefined,
            photos: variant.photos,
            costPrice: variant.costPrice,
            salePrice: variant.salePrice,
            stockQuantity:
              typeof variant.stockQuantity === 'string'
                ? parseInt(variant.stockQuantity, 10)
                : variant.stockQuantity,
            stockQuantityAlert:
              typeof variant.stockQuantityAlert === 'string'
                ? parseInt(variant.stockQuantityAlert, 10)
                : variant.stockQuantityAlert,
          })
        }

        for (const variant of newVariants) {
          await createVariantMutation.mutateAsync({
            product_id: product.id,
            code: variant.code,
            sku: variant.sku || undefined,
            photos: variant.photos,
            costPrice: variant.costPrice,
            salePrice: variant.salePrice,
            stockQuantity:
              typeof variant.stockQuantity === 'string'
                ? parseInt(variant.stockQuantity, 10)
                : variant.stockQuantity,
            stockQuantityAlert:
              typeof variant.stockQuantityAlert === 'string'
                ? parseInt(variant.stockQuantityAlert, 10)
                : variant.stockQuantityAlert,
          })
        }

        for (const variantId of variantsToDelete) {
          await inventoryApi.deleteVariant(businessDescriptor, variantId)
        }

        toast.success(tInventory('product_updated'))
        queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all })
        onSuccess()
      } catch (error) {
        const message = await translateErrorAsync(error, tInventory)
        toast.error(message)
      }
    },
  })

  const isLoading =
    updateProductMutation.isPending || createVariantMutation.isPending

  const handleDeleteVariant = (variantId: string) => {
    const currentVariants = form.getFieldValue('variants') as Array<{
      id?: string
    }>
    if (currentVariants.length <= 1) {
      toast.error(tInventory('cannot_delete_last_variant'))
      return
    }

    setVariantsToDelete((prev) => [...prev, variantId])

    const updatedVariants = currentVariants.filter((v) => v.id !== variantId)
    form.setFieldValue('variants', updatedVariants as any)
  }

  return (
    <BusinessContext.Provider value={{ businessDescriptor }}>
      <form.AppForm>
        <form.FormRoot className="flex flex-col h-full">
          <div className="flex-1 overflow-y-auto p-4 space-y-6">
            <form.FormError />

            <div className="space-y-4">
              <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide border-b border-base-300 pb-2">
                {tInventory('product_information')}
              </h3>

              <form.AppField
                name="name"
                validators={{
                  onBlur: multiVariantProductSchema.shape.name,
                }}
              >
                {(field) => (
                  <field.TextField
                    label={tInventory('product_name')}
                    placeholder={tInventory('product_name_placeholder')}
                    required
                  />
                )}
              </form.AppField>

              <form.AppField name="description">
                {(field) => (
                  <field.TextareaField
                    label={tInventory('description')}
                    placeholder={tInventory('description_placeholder')}
                    rows={3}
                  />
                )}
              </form.AppField>

              <form.AppField
                name="categoryId"
                validators={{
                  onBlur: multiVariantProductSchema.shape.categoryId,
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
            </div>

            <div className="space-y-4">
              <div className="flex items-center justify-between border-b border-base-300 pb-2">
                <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
                  {tInventory('variants')}
                </h3>
                <form.Subscribe selector={(state) => state.values.variants}>
                  {(variants) => (
                    <span className="text-sm text-base-content/70">
                      {variants.length} {tInventory('variant_count')}
                    </span>
                  )}
                </form.Subscribe>
              </div>

              <form.AppField name="variants">
                {(field) => (
                  <div className="space-y-4">
                    {field.state.value.map((variant: any, index: number) => (
                      <div
                        key={variant.id || index}
                        className="p-4 rounded-lg border border-base-300 bg-base-50 space-y-4"
                      >
                        <div className="flex items-center justify-between">
                          <span className="font-medium text-sm text-base-content">
                            {variant.id
                              ? tInventory('variant') + ' ' + (index + 1)
                              : tInventory('new_variant')}
                          </span>
                          {field.state.value.length > 1 && (
                            <button
                              type="button"
                              onClick={() => {
                                if (variant.id) {
                                  handleDeleteVariant(variant.id)
                                } else {
                                  const newVariants = field.state.value.filter(
                                    (_: any, i: number) => i !== index,
                                  )
                                  field.handleChange(newVariants)
                                }
                              }}
                              className="btn btn-ghost btn-sm btn-circle text-error"
                            >
                              <Trash2 className="w-4 h-4" />
                            </button>
                          )}
                        </div>

                        <form.AppField
                          name={`variants[${index}].code`}
                          validators={{
                            onBlur:
                              multiVariantProductSchema.shape.variants.element
                                .shape.code,
                          }}
                        >
                          {(codeField) => (
                            <codeField.TextField
                              label={tInventory('variant_code')}
                              placeholder={tInventory(
                                'variant_code_placeholder',
                              )}
                              required
                            />
                          )}
                        </form.AppField>

                        <form.AppField name={`variants[${index}].sku`}>
                          {(skuField) => (
                            <skuField.TextField
                              label={tInventory('sku')}
                              placeholder={tInventory('sku_placeholder')}
                              hint={tInventory('sku_hint')}
                            />
                          )}
                        </form.AppField>

                        <form.AppField name={`variants[${index}].photos`}>
                          {(photosField) => (
                            <photosField.FileUploadField
                              label={tInventory('variant_photos')}
                              accept="image/*"
                              maxFiles={10}
                              multiple
                            />
                          )}
                        </form.AppField>

                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                          <form.AppField
                            name={`variants[${index}].costPrice`}
                            validators={{
                              onBlur:
                                multiVariantProductSchema.shape.variants.element
                                  .shape.costPrice,
                            }}
                          >
                            {(priceField) => (
                              <priceField.PriceField
                                label={tInventory('cost_price')}
                                placeholder="0.00"
                                currencyCode={currencyCode}
                                required
                              />
                            )}
                          </form.AppField>

                          <form.AppField
                            name={`variants[${index}].salePrice`}
                            validators={{
                              onBlur:
                                multiVariantProductSchema.shape.variants.element
                                  .shape.salePrice,
                            }}
                          >
                            {(priceField) => (
                              <priceField.PriceField
                                label={tInventory('sale_price')}
                                placeholder="0.00"
                                currencyCode={currencyCode}
                                required
                              />
                            )}
                          </form.AppField>
                        </div>

                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                          <form.AppField
                            name={`variants[${index}].stockQuantity`}
                          >
                            {(stockField) => (
                              <stockField.TextField
                                type="text"
                                label={tInventory('stock_quantity')}
                                placeholder="0"
                                required
                              />
                            )}
                          </form.AppField>

                          <form.AppField
                            name={`variants[${index}].stockQuantityAlert`}
                          >
                            {(alertField) => (
                              <alertField.TextField
                                type="text"
                                label={tInventory('stock_alert')}
                                placeholder="0"
                                hint={tInventory('stock_alert_hint')}
                              />
                            )}
                          </form.AppField>
                        </div>
                      </div>
                    ))}

                    {variantsToDelete.map((deletedId) => {
                      const deletedVariant = product.variants?.find(
                        (v) => v.id === deletedId,
                      )
                      return (
                        <div key={deletedId} className="alert alert-warning">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="stroke-current shrink-0 h-6 w-6"
                            fill="none"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth="2"
                              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                            />
                          </svg>
                          <div className="flex-1">
                            <span className="font-medium">
                              {deletedVariant?.code || tInventory('variant')}
                            </span>
                            {' - '}
                            <span>
                              {tInventory('variant_marked_for_deletion')}
                            </span>
                          </div>
                        </div>
                      )
                    })}

                    <button
                      type="button"
                      onClick={() => {
                        field.handleChange([
                          ...field.state.value,
                          {
                            code: '',
                            sku: '',
                            photos: [],
                            costPrice: '',
                            salePrice: '',
                            stockQuantity: 0,
                            stockQuantityAlert: 0,
                          } as any,
                        ])
                      }}
                      className="btn btn-outline btn-sm w-full"
                      disabled={field.state.value.length >= 50}
                    >
                      <Plus className="w-4 h-4" />
                      {tInventory('add_variant')}
                    </button>
                  </div>
                )}
              </form.AppField>
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
