/**
 * Format currency with proper formatting
 * Shows currency code (e.g., USD 10.00) instead of symbol for consistency
 */
export function formatCurrency(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency.toUpperCase(),
    currencyDisplay: 'code',
  }).format(amount)
}
