/**
 * ItemPickerSheet Component
 *
 * A mobile-first sheet for selecting/editing order items.
 * Handles variant selection, quantity input, and price override.
 *
 * Per UX spec:
 * - Searchable variant select with thumbnail + stock info
 * - Quantity stepper (min 1, max stock)
 * - Optional price override field
 * - Focused field on open
 */

import { useEffect, useId, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import type { Variant } from '@/api/inventory'
import type { FormSelectOption } from '@/components/form/FormSelect'
import { useVariantQuery, useVariantsQuery } from '@/api/inventory'
import { FormSelect } from '@/components/form/FormSelect'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { useKyoraForm } from '@/lib/form'
import { getThumbnailUrl } from '@/lib/assetUrl'

export interface OrderItem {
  variantId: string
  quantity: number
  unitPrice: string
  unitCost: string
}

export interface ItemPickerSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  /** If provided, sheet is in edit mode */
  existingItem?: OrderItem
  /** Items already added (to show stock deductions) */
  currentItems?: Array<OrderItem>
  onSave: (item: OrderItem) => void
  onRemove?: () => void
}

const itemSchema = z.object({
  variantId: z.string().min(1, 'validation.required'),
  quantity: z.coerce.number().int().min(1, 'validation.min_value'),
  unitPrice: z.string().refine(
    (v) => {
      const n = parseFloat(v)
      return !isNaN(n) && n > 0
    },
    { message: 'validation.invalid_price' },
  ),
  unitCost: z.string().optional(),
})

type FormData = z.infer<typeof itemSchema>

export function ItemPickerSheet({
  isOpen,
  onClose,
  businessDescriptor,
  existingItem,
  currentItems = [],
  onSave,
  onRemove,
}: ItemPickerSheetProps) {
  const { t: tOrders } = useTranslation('orders')
  const { t: tCommon } = useTranslation('common')
  const formId = useId()

  const isEditMode = !!existingItem

  // Variant search state
  const [variantSearch, setVariantSearch] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')
  const [isSelectOpen, setIsSelectOpen] = useState(false)

  // Debounce search query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(variantSearch)
    }, 300)
    return () => clearTimeout(timer)
  }, [variantSearch])

  // Fetch variants
  const listParams = useMemo(
    () => ({
      search: debouncedSearch || undefined,
      page: 1,
      pageSize: 15,
    }),
    [debouncedSearch],
  )

  const variantsQuery = useVariantsQuery(businessDescriptor, listParams)
  const variantsData = variantsQuery.data

  // Local state for selected variant (for display details)
  const [selectedVariantId, setSelectedVariantId] = useState<string>(
    existingItem?.variantId ?? '',
  )

  // Query selected variant details if not in list
  const selectedVariantQuery = useVariantQuery(
    businessDescriptor,
    selectedVariantId,
  )

  const selectedVariant = useMemo<Variant | null>(() => {
    if (!selectedVariantId) return null
    const fromList = variantsData?.items.find((v) => v.id === selectedVariantId)
    return fromList ?? selectedVariantQuery.data ?? null
  }, [selectedVariantId, variantsData?.items, selectedVariantQuery.data])

  // Calculate available stock (original stock minus items in form, plus back what's being edited)
  const availableStock = useMemo(() => {
    if (!selectedVariant) return 0

    const baseStock = selectedVariant.stockQuantity
    const usedByOtherItems = currentItems
      .filter((item) => {
        // Don't count the item being edited
        if (existingItem && item.variantId === existingItem.variantId) {
          return false
        }
        return item.variantId === selectedVariantId
      })
      .reduce((sum, item) => sum + item.quantity, 0)

    return Math.max(0, baseStock - usedByOtherItems)
  }, [selectedVariant, currentItems, selectedVariantId, existingItem])

  // Build variant options for select
  const variantOptions = useMemo<Array<FormSelectOption<string>>>(() => {
    if (!variantsData?.items) return []
    return variantsData.items.map((variant) => {
      const firstPhoto =
        variant.photos.length > 0 ? variant.photos[0] : undefined
      return {
        value: variant.id,
        label: variant.name,
        description: `${tOrders('price')}: ${variant.salePrice} ${variant.currency} | ${tOrders('stock')}: ${variant.stockQuantity}`,
        thumbnail: firstPhoto ? getThumbnailUrl(firstPhoto) : undefined,
      }
    })
  }, [variantsData?.items, tOrders])

  const defaultValues: FormData = {
    variantId: existingItem?.variantId ?? '',
    quantity: existingItem?.quantity ?? 1,
    unitPrice: existingItem?.unitPrice ?? '',
    unitCost: existingItem?.unitCost ?? '',
  }

  const form = useKyoraForm({
    defaultValues,
    onSubmit: ({ value }) => {
      onSave({
        variantId: value.variantId,
        quantity: value.quantity,
        unitPrice: value.unitPrice,
        unitCost: value.unitCost ?? '',
      })
      onClose()
    },
  })

  // Sync selectedVariantId with form value
  useEffect(() => {
    const unsubscribe = form.store.subscribe(() => {
      const variantId = form.state.values.variantId
      if (variantId !== selectedVariantId) {
        setSelectedVariantId(variantId)
      }
    })
    return unsubscribe
  }, [form, selectedVariantId])

  // Reset form when sheet opens/closes
  useEffect(() => {
    if (isOpen) {
      form.reset()
      if (existingItem) {
        form.setFieldValue('variantId', existingItem.variantId)
        form.setFieldValue('quantity', existingItem.quantity)
        form.setFieldValue('unitPrice', existingItem.unitPrice)
        form.setFieldValue('unitCost', existingItem.unitCost)
        setSelectedVariantId(existingItem.variantId)
      } else {
        setSelectedVariantId('')
      }
      setVariantSearch('')
      setDebouncedSearch('')
    } else {
      setVariantSearch('')
      setDebouncedSearch('')
    }
  }, [isOpen, existingItem, form])

  // Populate search field with selected variant name when closing dropdown
  useEffect(() => {
    if (!isSelectOpen && selectedVariant) {
      setVariantSearch(selectedVariant.name)
    }
  }, [isSelectOpen, selectedVariant])

  const handleVariantSelect = (variantId: string | Array<string>) => {
    const id = Array.isArray(variantId) ? variantId[0] : variantId
    form.setFieldValue('variantId', id)

    const variant = variantsData?.items.find((v) => v.id === id)
    if (variant) {
      setSelectedVariantId(id)
      setVariantSearch(variant.name)
      // Auto-fill price and cost from variant
      form.setFieldValue('unitPrice', variant.salePrice)
      form.setFieldValue('unitCost', variant.costPrice)
      // Reset quantity to 1 or max stock
      const maxQty = Math.min(variant.stockQuantity, 1)
      form.setFieldValue('quantity', Math.max(1, maxQty))
    }
  }

  const handleClearVariant = () => {
    form.setFieldValue('variantId', '')
    form.setFieldValue('unitPrice', '')
    form.setFieldValue('unitCost', '')
    form.setFieldValue('quantity', 1)
    setSelectedVariantId('')
    setVariantSearch('')
  }

  const safeClose = () => {
    onClose()
  }

  return (
    <form.AppForm>
      <BottomSheet
        isOpen={isOpen}
        onClose={safeClose}
        title={
          isEditMode ? tOrders('item_picker_edit') : tOrders('item_picker_add')
        }
        footer={
          <div className="flex gap-2">
            {isEditMode && onRemove && (
              <button
                type="button"
                className="btn btn-error btn-outline"
                onClick={() => {
                  onRemove()
                  onClose()
                }}
              >
                {tCommon('remove')}
              </button>
            )}
            <button
              type="button"
              className="btn btn-ghost flex-1"
              onClick={safeClose}
            >
              {tCommon('cancel')}
            </button>
            <form.SubmitButton
              form={`item-picker-form-${formId}`}
              variant="primary"
              className="flex-1"
            >
              {isEditMode
                ? tOrders('item_picker_update')
                : tOrders('item_picker_save')}
            </form.SubmitButton>
          </div>
        }
        side="end"
        size="md"
        contentClassName="space-y-4"
        ariaLabel={
          isEditMode ? tOrders('item_picker_edit') : tOrders('item_picker_add')
        }
      >
        <form.FormRoot id={`item-picker-form-${formId}`} className="space-y-4">
          {/* Variant Selector */}
          <div className="form-control w-full">
            <label className="label">
              <span className="label-text text-base-content/70 font-medium">
                {tOrders('product')}
                <span className="text-error ms-1">*</span>
              </span>
            </label>

            <FormSelect
              options={variantOptions}
              value={form.state.values.variantId}
              onChange={handleVariantSelect}
              searchable
              searchValue={variantSearch}
              onSearchChange={setVariantSearch}
              placeholder={tOrders('search_product')}
              clearable
              onClear={handleClearVariant}
              isLoading={variantsQuery.isPending && !variantsData}
              onOpen={() => {
                setIsSelectOpen(true)
                setVariantSearch('')
                setDebouncedSearch('')
              }}
              onClose={() => {
                setIsSelectOpen(false)
              }}
            />

            {/* Selected variant preview */}
            {selectedVariant && !isSelectOpen && (
              <div className="px-3 py-2 mt-2 bg-base-200 rounded-lg">
                <div className="flex items-center gap-3">
                  {selectedVariant.photos.length > 0 &&
                    selectedVariant.photos[0] && (
                      <div className="avatar">
                        <div className="w-10 h-10 rounded">
                          <img
                            src={getThumbnailUrl(selectedVariant.photos[0])}
                            alt={selectedVariant.name}
                          />
                        </div>
                      </div>
                    )}
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-base-content truncate">
                      {selectedVariant.name}
                    </div>
                    <div className="text-xs text-base-content/60">
                      {tOrders('available_stock', { count: availableStock })}
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Quantity Field */}
          <form.AppField
            name="quantity"
            validators={{
              onBlur: z.coerce.number().int().min(1, 'validation.min_value'),
            }}
          >
            {(field) => (
              <field.QuantityField
                label={tOrders('quantity')}
                min={1}
                max={availableStock || 10000}
                required
                disabled={!selectedVariant}
                suffix={
                  selectedVariant && availableStock > 0
                    ? `/ ${availableStock}`
                    : undefined
                }
              />
            )}
          </form.AppField>

          {/* Price Override */}
          <form.AppField
            name="unitPrice"
            validators={{
              onBlur: z.string().refine(
                (v) => {
                  const n = parseFloat(v)
                  return !isNaN(n) && n > 0
                },
                { message: 'validation.invalid_price' },
              ),
            }}
          >
            {(field) => (
              <field.PriceField
                label={tOrders('item_unit_price')}
                required
                hint={
                  selectedVariant
                    ? `${tOrders('original_price')}: ${selectedVariant.salePrice} ${selectedVariant.currency}`
                    : undefined
                }
              />
            )}
          </form.AppField>

          {/* Cost (hidden from user but stored) */}
          <input
            type="hidden"
            name="unitCost"
            value={form.state.values.unitCost || ''}
          />
        </form.FormRoot>
      </BottomSheet>
    </form.AppForm>
  )
}
