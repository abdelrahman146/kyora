import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { Plus, Trash2 } from 'lucide-react'
import { useState } from 'react'

import type { Product } from '@/api/inventory'
import { useKyoraForm } from '@/lib/form'
import { BusinessContext } from '@/lib/form/components/FileUploadField'
import { multiVariantProductSchema } from '@/schemas/inventory'
import {
  inventoryApi,
  useCategoriesQuery,
  useCreateVariantMutation,
  useUpdateProductMutation,
} from '@/api/inventory'
import { queryKeys } from '@/lib/queryKeys'
import { translateErrorAsync } from '@/lib/translateError'

interface UpdateProductWithVariantsFormProps {
  product: Product
  businessDescriptor: string
  onSuccess: () => void
  onCancel: () => void
}

export function UpdateProductWithVariantsForm({
  product,
  businessDescriptor,
  onSuccess,
  onCancel,
}: UpdateProductWithVariantsFormProps) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [variantsToDelete, setVariantsToDelete] = useState<Array<string>>([])

  const categoriesQuery = useCategoriesQuery(businessDescriptor)
  const categories = categoriesQuery.data ?? []

  const updateProductMutation = useUpdateProductMutation(
    businessDescriptor,
    product.id,
    {
      onError: async (error) => {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      },
    },
  )

  const createVariantMutation = useCreateVariantMutation(
    businessDescriptor,
    product.id,
    {
      onError: async (error) => {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      },
    },
  )

  // Deletion happens on form submit, not immediately

  const form = useKyoraForm({
    defaultValues: {
      name: product.name || '',
      description: product.description || '',
      photos: product.photos,
      categoryId: product.categoryId || '',
      variants:
        product.variants?.map((v) => ({
          id: v.id,
          code: v.code || '',
          sku: v.sku || '',
          photos: v.photos,
          costPrice: v.costPrice || '',
          salePrice: v.salePrice || '',
          stockQuantity: v.stockQuantity || 0,
          stockQuantityAlert: v.stockQuantityAlert || 0,
        })) || [],
    },
    onSubmit: async ({ value }) => {
      try {
        // First, update the product
        await updateProductMutation.mutateAsync({
          name: value.name,
          description: value.description || undefined,
          photos: value.photos,
          categoryId: value.categoryId,
        })

        // Then handle variants
        const existingVariants = value.variants.filter((v) => 'id' in v && v.id)
        const newVariants = value.variants.filter((v) => !('id' in v) || !v.id)

        // Update existing variants using inventoryApi directly
        for (const variant of existingVariants) {
          if ('id' in variant && variant.id) {
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
        }

        // Create new variants
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

        // Delete variants marked for deletion
        for (const variantId of variantsToDelete) {
          await inventoryApi.deleteVariant(businessDescriptor, variantId)
        }

        toast.success(t('product_updated', { ns: 'inventory' }))
        queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all })
        onSuccess()
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
  })

  const categoryOptions = categories.map((cat) => ({
    value: cat.id,
    label: cat.name,
  }))

  const isLoading =
    updateProductMutation.isPending || createVariantMutation.isPending

  const handleDeleteVariant = (variantId: string) => {
    const currentVariants = form.getFieldValue('variants') as Array<{
      id?: string
    }>
    if (currentVariants.length <= 1) {
      toast.error(t('cannot_delete_last_variant', { ns: 'inventory' }))
      return
    }

    // Mark variant for deletion
    setVariantsToDelete((prev) => [...prev, variantId])

    // Remove from form UI
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
              <h3 className="text-base font-semibold text-base-content border-b border-base-300 pb-2">
                {t('product_information', { ns: 'inventory' })}
              </h3>

              <form.AppField
                name="name"
                validators={{
                  onBlur: multiVariantProductSchema.shape.name,
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
                  onBlur: multiVariantProductSchema.shape.categoryId,
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
            </div>

            <div className="space-y-4">
              <div className="flex items-center justify-between border-b border-base-300 pb-2">
                <h3 className="text-base font-semibold text-base-content">
                  {t('variants', { ns: 'inventory' })}
                </h3>
                <form.Subscribe selector={(state) => state.values.variants}>
                  {(variants) => (
                    <span className="text-sm text-base-content/70">
                      {variants.length}{' '}
                      {t('variant_count', { ns: 'inventory' })}
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
                              ? t('variant', { ns: 'inventory' }) +
                                ' ' +
                                (index + 1)
                              : t('new_variant', { ns: 'inventory' })}
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
                              label={t('variant_code', { ns: 'inventory' })}
                              placeholder={t(
                                'inventory.variant_code_placeholder',
                              )}
                              required
                            />
                          )}
                        </form.AppField>

                        <form.AppField name={`variants[${index}].sku`}>
                          {(skuField) => (
                            <skuField.TextField
                              label={t('sku', { ns: 'inventory' })}
                              placeholder={t('sku_placeholder', {
                                ns: 'inventory',
                              })}
                              hint={t('sku_hint', { ns: 'inventory' })}
                            />
                          )}
                        </form.AppField>

                        <form.AppField name={`variants[${index}].photos`}>
                          {(photosField) => (
                            <photosField.FileUploadField
                              label={t('variant_photos', { ns: 'inventory' })}
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
                              <priceField.TextField
                                type="text"
                                label={t('cost_price', { ns: 'inventory' })}
                                placeholder="0.00"
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
                              <priceField.TextField
                                type="text"
                                label={t('sale_price', { ns: 'inventory' })}
                                placeholder="0.00"
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
                                label={t('stock_quantity', { ns: 'inventory' })}
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
                                label={t('stock_alert', { ns: 'inventory' })}
                                placeholder="0"
                                hint={t('stock_alert_hint', {
                                  ns: 'inventory',
                                })}
                              />
                            )}
                          </form.AppField>
                        </div>
                      </div>
                    ))}

                    {/* Show alerts for deleted variants */}
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
                              {deletedVariant?.code ||
                                t('variant', { ns: 'inventory' })}
                            </span>
                            {' - '}
                            <span>
                              {t('variant_marked_for_deletion', {
                                ns: 'inventory',
                              })}
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
                            // Don't include id for new variants
                            code: '',
                            sku: '',
                            photos: [],
                            costPrice: '',
                            salePrice: '',
                            stockQuantity: 0,
                            stockQuantityAlert: 0,
                          } as any, // Type assertion since new variants don't have id
                        ])
                      }}
                      className="btn btn-outline btn-sm w-full"
                      disabled={field.state.value.length >= 50}
                    >
                      <Plus className="w-4 h-4" />
                      {t('add_variant', { ns: 'inventory' })}
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
              {t('common.cancel')}
            </button>
            <form.SubmitButton variant="primary" className="flex-1">
              {t('update_product', { ns: 'inventory' })}
            </form.SubmitButton>
          </div>
        </form.FormRoot>
      </form.AppForm>
    </BusinessContext.Provider>
  )
}
