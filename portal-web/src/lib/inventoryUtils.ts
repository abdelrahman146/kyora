import type { Variant } from '@/api/inventory'

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
    return sum + parseFloat(variant.cost_price)
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
    return sum + variant.stock_quantity
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
      variant.stock_quantity > 0 &&
      variant.stock_quantity <= variant.stock_quantity_alert,
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
      variant.stock_quantity > 0 &&
      variant.stock_quantity <= variant.stock_quantity_alert,
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
