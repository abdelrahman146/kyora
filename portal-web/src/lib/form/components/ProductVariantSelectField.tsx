/**
 * ProductVariantSelectField Component - Form Composition Layer
 *
 * Pre-bound product variant selector with autocomplete search functionality.
 * Automatically wires to TanStack Form field context and handles variant fetching.
 *
 * DESIGN GUIDELINES:
 * ==================
 * This component provides a searchable variant selector with autocomplete,
 * displaying variant details (price, stock) and product thumbnail for better identification.
 *
 * Key Features:
 * 1. Autocomplete search - Real-time search with debouncing
 * 2. Variant details - Shows price and available stock below selection
 * 3. Product thumbnail - Displays variant/product image next to name
 * 4. Persistent state - Loads selected variant on mount (survives refresh)
 * 5. Mobile-first - Uses bottom sheet on mobile via FormSelect
 * 6. RTL support - Full right-to-left layout support
 * 7. Accessibility - WCAG AA compliant with ARIA attributes
 *
 * Usage within form:
 * ```tsx
 * <form.AppField name="variantId">
 *   {(field) => (
 *     <field.ProductVariantSelectField
 *       label="Select Product"
 *       businessDescriptor="my-business"
 *       placeholder="Search products..."
 *     />
 *   )}
 * </form.AppField>
 * ```
 */

import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { useFieldContext } from '../contexts'
import type { ProductVariantSelectFieldProps } from '../types'
import type { Variant } from '@/api/inventory'
import type { FormSelectOption } from '@/components/form/FormSelect'
import { useVariantQuery, useVariantsQuery } from '@/api/inventory'
import { FormSelect } from '@/components/form/FormSelect'
import { getThumbnailUrl } from '@/lib/assetUrl'

export function ProductVariantSelectField(
  props: ProductVariantSelectFieldProps,
) {
  const field = useFieldContext<string>()
  const { t: tErrors } = useTranslation('errors')
  const { t: tOrders } = useTranslation('orders')

  const [variantSearchQuery, setVariantSearchQuery] = useState('')
  const [debouncedVariantSearch, setDebouncedVariantSearch] = useState('')
  const [isOpen, setIsOpen] = useState(false)

  // Debounce variant search query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedVariantSearch(variantSearchQuery)
    }, 300)
    return () => clearTimeout(timer)
  }, [variantSearchQuery])

  const listParams = useMemo(
    () => ({
      search: debouncedVariantSearch || undefined,
      page: 1,
      pageSize: 10,
    }),
    [debouncedVariantSearch],
  )

  const variantsQuery = useVariantsQuery(props.businessDescriptor, listParams)
  const variantsData = variantsQuery.data

  const selectedVariantQuery = useVariantQuery(
    props.businessDescriptor,
    field.state.value,
  )

  const selectedVariant = useMemo<Variant | null>(() => {
    const fromList = variantsData?.items.find((v) => v.id === field.state.value)
    return fromList ?? selectedVariantQuery.data ?? null
  }, [variantsData?.items, field.state.value, selectedVariantQuery.data])

  useEffect(() => {
    if (isOpen) return
    if (!field.state.value) return
    if (!selectedVariant) return
    setVariantSearchQuery(selectedVariant.name)
  }, [field.state.value, isOpen, selectedVariant])

  // Transform variants to select options with thumbnail
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
  }, [variantsData, tOrders])

  // Extract error from field state and translate
  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') {
      return tErrors(firstError)
    }

    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return tErrors(errorObj.message)
    }

    return undefined
  }, [field.state.meta.errors, tErrors])

  // Show error only after field has been touched
  const showError = field.state.meta.isTouched && error

  const handleVariantSelect = (variantId: string | Array<string>) => {
    const id = Array.isArray(variantId) ? variantId[0] : variantId
    field.handleChange(id)
    const variant = variantsData?.items.find((v) => v.id === id)
    if (variant) {
      setVariantSearchQuery(variant.name)
      // Call onVariantSelect callback if provided
      if (props.onVariantSelect) {
        props.onVariantSelect({
          id: variant.id,
          productId: variant.productId,
          productName: variant.name, // Variant name includes product name
          variantName: variant.name,
          salePrice: variant.salePrice,
          costPrice: variant.costPrice,
          stockQuantity: variant.stockQuantity,
        })
      }
    }
  }

  const handleClear = () => {
    field.handleChange('')
    setVariantSearchQuery('')
    if (props.onClear) {
      props.onClear()
    }
  }

  return (
    <div className="form-control w-full">
      {/* Label - Matches TextField pattern */}
      {props.label && (
        <label className="label">
          <span className="label-text text-base-content/70 font-medium">
            {props.label}
            {props.required && <span className="text-error ms-1">*</span>}
          </span>
        </label>
      )}

      {/* FormSelect with variant search */}
      <FormSelect
        options={variantOptions}
        value={field.state.value}
        onChange={handleVariantSelect}
        searchable
        searchValue={variantSearchQuery}
        onSearchChange={setVariantSearchQuery}
        placeholder={props.placeholder || tOrders('search_product')}
        clearable
        onClear={handleClear}
        disabled={props.disabled || field.state.meta.isValidating}
        error={showError ? error : undefined}
        isLoading={variantsQuery.isFetching}
        onOpen={() => {
          setIsOpen(true)
          setVariantSearchQuery('')
        }}
        onClose={() => {
          setIsOpen(false)
          if (selectedVariant) {
            setVariantSearchQuery(selectedVariant.name)
          }
        }}
      />

      {/* Selected variant details */}
      {selectedVariant && !isOpen && (
        <div className="px-3 py-2 mt-2 bg-base-200 rounded-lg">
          <div className="flex items-center gap-3">
            {/* Variant thumbnail */}
            {selectedVariant.photos.length > 0 && selectedVariant.photos[0] && (
              <div className="avatar">
                <div className="w-12 h-12 rounded">
                  <img
                    src={getThumbnailUrl(selectedVariant.photos[0])}
                    alt={selectedVariant.name}
                  />
                </div>
              </div>
            )}

            {/* Variant details */}
            <div className="flex-1 min-w-0">
              <div className="text-sm font-medium text-base-content truncate">
                {selectedVariant.name}
              </div>
              <div className="text-xs text-base-content/60 flex gap-3 mt-1">
                <span>
                  {tOrders('price')}: {selectedVariant.salePrice}{' '}
                  {selectedVariant.currency}
                </span>
                <span className="opacity-40">•</span>
                <span>
                  {tOrders('available_stock')}: {selectedVariant.stockQuantity}
                </span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
