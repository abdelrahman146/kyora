import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { FormSelect } from '../atoms/FormSelect'
import type { FormSelectOption } from '../atoms/FormSelect'
import type { ShippingZone } from '@/api/business'

export interface ShippingZoneSelectProps {
  value: string
  onChange: (value: string) => void
  zones: Array<ShippingZone>
  isLoading?: boolean
  error?: string
  disabled?: boolean
  required?: boolean
  placeholder?: string
  searchable?: boolean
}

/**
 * Reusable shipping zone select component
 * Displays shipping zones with country count
 * Supports search functionality
 */
export function ShippingZoneSelect({
  value,
  onChange,
  zones,
  isLoading,
  error,
  disabled,
  required,
  placeholder,
  searchable = true,
}: ShippingZoneSelectProps) {
  const { t } = useTranslation()

  const zoneOptions: Array<FormSelectOption> = useMemo(() => {
    return zones.map((zone) => {
      const countryCount = zone.countries.length
      const label = `${zone.name} (${countryCount} ${t('customers.address.countries', { count: countryCount })})`
      return { value: zone.id, label }
    })
  }, [zones, t])

  return (
    <FormSelect<string>
      label={t('customers.address.shipping_zone')}
      options={zoneOptions}
      value={value}
      onChange={(val) => {
        onChange(val as string)
      }}
      required={required}
      disabled={disabled ?? isLoading}
      placeholder={placeholder ?? t('customers.address.select_shipping_zone')}
      searchable={searchable}
      error={error}
    />
  )
}
