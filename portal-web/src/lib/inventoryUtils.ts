import type { Variant } from '@/api/inventory'

/**
 * Get price display (single value or range)
 * @param variants Array of product variants
 * @param priceKey Either 'costPrice' or 'salePrice'
 * @returns Object with min, max, and isSame flag
 */
export function getPriceRange(
  variants: Array<Variant> | undefined,
  priceKey: 'costPrice' | 'salePrice',
): { min: number; max: number; isSame: boolean } {
  if (!variants || variants.length === 0) {
    return { min: 0, max: 0, isSame: true }
  }

  const prices = variants.map((v) => parseFloat(v[priceKey]))
  const min = Math.min(...prices)
  const max = Math.max(...prices)

  return {
    min,
    max,
    isSame: min === max,
  }
}

/**
 * Calculate average cost price across all variants
 * @param variants Array of product variants
 * @returns Average cost price as a number
 */
export function calculateAverageCostPrice(
  variants: Array<Variant> | undefined,
): number {
  if (!variants || variants.length === 0) return 0

  const total = variants.reduce((sum, variant) => {
    return sum + parseFloat(variant.costPrice)
  }, 0)

  return total / variants.length
}

/**
 * Calculate total stock quantity across all variants
 * @param variants Array of product variants
 * @returns Total stock quantity
 */
export function calculateTotalStock(
  variants: Array<Variant> | undefined,
): number {
  if (!variants || variants.length === 0) return 0

  return variants.reduce((sum, variant) => {
    return sum + variant.stockQuantity
  }, 0)
}

/**
 * Get stock status for a product based on its variants
 * @param variants Array of product variants
 * @returns Stock status: 'in_stock' | 'low_stock' | 'out_of_stock'
 */
export function getStockStatus(
  variants: Array<Variant> | undefined,
): 'in_stock' | 'low_stock' | 'out_of_stock' {
  if (!variants || variants.length === 0) return 'out_of_stock'

  const totalStock = calculateTotalStock(variants)

  // Out of stock if total is 0
  if (totalStock === 0) return 'out_of_stock'

  // Low stock if any variant is at or below its alert threshold
  const hasLowStockVariant = variants.some(
    (variant) =>
      variant.stockQuantity > 0 &&
      variant.stockQuantity <= variant.stockQuantityAlert,
  )

  if (hasLowStockVariant) return 'low_stock'

  return 'in_stock'
}

/**
 * Check if product has any variants with low stock
 * @param variants Array of product variants
 * @returns True if any variant has low stock
 */
export function hasLowStock(variants: Array<Variant> | undefined): boolean {
  if (!variants || variants.length === 0) return false

  return variants.some(
    (variant) =>
      variant.stockQuantity > 0 &&
      variant.stockQuantity <= variant.stockQuantityAlert,
  )
}

/**
 * Build orderBy array from search params
 * @param search Search params with sortBy and sortOrder
 * @returns Array of orderBy strings or undefined
 */
export function buildOrderBy(search: {
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}): Array<string> | undefined {
  if (!search.sortBy) return undefined

  const prefix = search.sortOrder === 'desc' ? '-' : ''
  return [`${prefix}${search.sortBy}`]
}

/**
 * Get badge variant based on stock status
 * @param status Stock status
 * @returns Badge variant string
 */
export function getStockStatusBadgeVariant(
  status: 'in_stock' | 'low_stock' | 'out_of_stock',
): 'success' | 'warning' | 'error' {
  switch (status) {
    case 'in_stock':
      return 'success'
    case 'low_stock':
      return 'warning'
    case 'out_of_stock':
      return 'error'
  }
}
