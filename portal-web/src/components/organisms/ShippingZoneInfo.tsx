/**
 * ShippingZoneInfo Component
 *
 * Read-only display of inferred shipping zone information based on address.
 * Used in order forms to show applicable shipping details.
 */

import { useTranslation } from 'react-i18next'
import { Package } from 'lucide-react'

import type { ShippingZone } from '@/api/business'
import { formatCurrency } from '@/lib/formatCurrency'

export interface ShippingZoneInfoProps {
  zone: ShippingZone | null | undefined
  currency?: string
  className?: string
}

export function ShippingZoneInfo({
  zone,
  currency,
  className = '',
}: ShippingZoneInfoProps) {
  const { t: tOrders } = useTranslation('orders')

  if (!zone) {
    return (
      <div className={`alert alert-warning ${className}`}>
        <Package size={20} />
        <span className="text-sm">
          {tOrders('no_shipping_zone_applicable')}
        </span>
      </div>
    )
  }

  const zoneCurrency = currency || zone.currency
  const shippingCost = parseFloat(zone.shippingCost)
  const freeThreshold = parseFloat(zone.freeShippingThreshold)

  return (
    <div className={`rounded-lg bg-base-200 p-4 space-y-2 ${className}`}>
      <div className="flex items-center gap-2">
        <Package size={18} className="text-base-content/70" />
        <h4 className="font-semibold text-sm">{tOrders('shipping_zone')}</h4>
      </div>

      <div className="space-y-1 text-sm">
        <div className="flex justify-between">
          <span className="text-base-content/70">{tOrders('zone_name')}:</span>
          <span className="font-medium">{zone.name}</span>
        </div>

        <div className="flex justify-between">
          <span className="text-base-content/70">
            {tOrders('shipping_cost')}:
          </span>
          <span className="font-medium" dir="ltr">
            {formatCurrency(shippingCost, zoneCurrency)}
          </span>
        </div>

        {freeThreshold > 0 && (
          <div className="flex justify-between">
            <span className="text-base-content/70">
              {tOrders('free_shipping_threshold')}:
            </span>
            <span className="font-medium text-success" dir="ltr">
              {formatCurrency(freeThreshold, zoneCurrency)}
            </span>
          </div>
        )}

        <div className="flex justify-between">
          <span className="text-base-content/70">
            {tOrders('countries_covered')}:
          </span>
          <span className="font-medium">{zone.countries.length}</span>
        </div>
      </div>
    </div>
  )
}
