import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { useParams } from '@tanstack/react-router'
import toast from 'react-hot-toast'
import { Plus, Trash2 } from 'lucide-react'

import type {
  AssetReference,
  CreateProductWithVariantsRequest,
} from '@/api/inventory'
import { useKyoraForm } from '@/lib/form'
import { BusinessContext } from '@/lib/form/components/FileUploadField'
import { multiVariantProductSchema } from '@/schemas/inventory'
import { useCreateProductWithVariantsMutation } from '@/api/inventory'
import { queryKeys } from '@/lib/queryKeys'
import { translateErrorAsync } from '@/lib/translateError'
import { getSelectedBusiness } from '@/stores/businessStore'

interface MultiVariantProductFormProps {
  onSuccess: () => void
  onCancel: () => void
}

export function MultiVariantProductForm({
  onSuccess,
  onCancel,
}: MultiVariantProductFormProps) {
  const { t } = useTranslation()
  const { businessDescriptor } = useParams({ strict: false })
  const queryClient = useQueryClient()
  const selectedBusiness = getSelectedBusiness()
  const currencyCode = selectedBusiness?.currency ?? 'USD'

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
      variants: [
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
      const payload: CreateProductWithVariantsRequest = {
        product: {
          name: value.name,
          description: value.description || undefined,
          photos: value.photos,
          categoryId: value.categoryId,
        },
        variants: value.variants.map((v) => ({
          code: v.code,
          sku: v.sku || undefined,
          photos: v.photos,
          costPrice: v.costPrice,
          salePrice: v.salePrice,
          stockQuantity:
            typeof v.stockQuantity === 'string'
              ? parseInt(v.stockQuantity, 10)
              : v.stockQuantity,
          stockQuantityAlert:
            typeof v.stockQuantityAlert === 'string'
              ? parseInt(v.stockQuantityAlert, 10)
              : v.stockQuantityAlert,
        })),
      }

      await createMutation.mutateAsync(payload)
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
              <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide border-b border-base-300 pb-2">
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
                  <field.CategorySelectField
                    businessDescriptor={businessDescriptor!}
                    label={t('category', { ns: 'inventory' })}
                    placeholder={t('select_category', { ns: 'inventory' })}
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
                <h3 className="text-sm font-semibold text-base-content/60 uppercase tracking-wide">
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
                    {field.state.value.map((_variant: any, index: number) => (
                      <div
                        key={index}
                        className="p-4 rounded-lg border border-base-300 bg-base-50 space-y-4"
                      >
                        <div className="flex items-center justify-between">
                          <span className="font-medium text-sm text-base-content">
                            {t('variant', { ns: 'inventory' })} {index + 1}
                          </span>
                          {field.state.value.length > 1 && (
                            <button
                              type="button"
                              onClick={() => {
                                const newVariants = field.state.value.filter(
                                  (_: any, i: number) => i !== index,
                                )
                                field.handleChange(newVariants)
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
                            {(costPriceField) => (
                              <costPriceField.PriceField
                                label={t('cost_price', { ns: 'inventory' })}
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
                            {(salePriceField) => (
                              <salePriceField.PriceField
                                label={t('sale_price', { ns: 'inventory' })}
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

                    <button
                      type="button"
                      className="btn btn-outline w-full"
                      onClick={() =>
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
                          },
                        ])
                      }
                    >
                      <Plus className="w-4 h-4" />
                      <span className="ms-2">
                        {t('add_variant', { ns: 'inventory' })}
                      </span>
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
