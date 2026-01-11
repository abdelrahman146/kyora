import type { Variant } from '@/api/inventory'

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

export function calculateAverageCostPrice(
  variants: Array<Variant> | undefined,
): number {
  if (!variants || variants.length === 0) return 0

  const total = variants.reduce((sum, variant) => {
    return sum + parseFloat(variant.costPrice)
  }, 0)

  return total / variants.length
}

export function calculateTotalStock(
  variants: Array<Variant> | undefined,
): number {
  if (!variants || variants.length === 0) return 0

  return variants.reduce((sum, variant) => {
    return sum + variant.stockQuantity
  }, 0)
}

export function getStockStatus(
  variants: Array<Variant> | undefined,
): 'in_stock' | 'low_stock' | 'out_of_stock' {
  if (!variants || variants.length === 0) return 'out_of_stock'

  const totalStock = calculateTotalStock(variants)

  if (totalStock === 0) return 'out_of_stock'

  const hasLowStockVariant = variants.some(
    (variant) =>
      variant.stockQuantity > 0 &&
      variant.stockQuantity <= variant.stockQuantityAlert,
  )

  if (hasLowStockVariant) return 'low_stock'

  return 'in_stock'
}

export function hasLowStock(variants: Array<Variant> | undefined): boolean {
  if (!variants || variants.length === 0) return false

  return variants.some(
    (variant) =>
      variant.stockQuantity > 0 &&
      variant.stockQuantity <= variant.stockQuantityAlert,
  )
}

export function buildOrderBy(search: {
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}): Array<string> | undefined {
  if (!search.sortBy) return undefined

  const prefix = search.sortOrder === 'desc' ? '-' : ''
  return [`${prefix}${search.sortBy}`]
}

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
