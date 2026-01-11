import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import type { ShippingZone } from '@/api/business'
import type { FormSelectOption } from '@/components/form/FormSelect'
import { FormSelect } from '@/components/form/FormSelect'

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
  const { t: tCustomers } = useTranslation('customers')

  const zoneOptions: Array<FormSelectOption> = useMemo(() => {
    return zones.map((zone) => {
      const countryCount = zone.countries.length
      const label = `${zone.name} (${countryCount} ${tCustomers('address.countries', { count: countryCount })})`
      return { value: zone.id, label }
    })
  }, [zones, tCustomers])

  return (
    <FormSelect<string>
      label={tCustomers('address.shipping_zone')}
      options={zoneOptions}
      value={value}
      onChange={(val) => {
        onChange(val as string)
      }}
      required={required}
      disabled={disabled ?? isLoading}
      placeholder={placeholder ?? tCustomers('address.select_shipping_zone')}
      searchable={searchable}
      error={error}
    />
  )
}
