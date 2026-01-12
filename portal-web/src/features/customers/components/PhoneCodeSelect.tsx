import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import type { FormSelectOption } from '@/components/form/FormSelect'
import { useCountriesQuery } from '@/api/metadata'
import { FormSelect } from '@/components/form/FormSelect'

export interface PhoneCodeSelectProps {
  value: string
  onChange: (value: string) => void
  error?: string
  disabled?: boolean
  required?: boolean
  placeholder?: string
  searchable?: boolean
}

/**
 * Reusable phone code select component
 * Uses FormSelect with phone codes from metadata store
 * Displays phone prefix with country name
 * Automatically deduplicates phone codes
 * Supports search functionality
 */
export function PhoneCodeSelect({
  value,
  onChange,
  error,
  disabled,
  required,
  placeholder,
  searchable = true,
}: PhoneCodeSelectProps) {
  const { i18n } = useTranslation()
  const { t: tCustomers } = useTranslation('customers')
  const { data: countries = [], isSuccess } = useCountriesQuery()

  const isArabic = i18n.language.toLowerCase().startsWith('ar')
  const countriesReady = isSuccess && countries.length > 0

  const phoneCodeOptions: Array<FormSelectOption> = useMemo(() => {
    const seen = new Set<string>()
    const options: Array<FormSelectOption> = []

    for (const c of countries) {
      if (!c.phonePrefix) continue
      if (seen.has(c.phonePrefix)) continue
      seen.add(c.phonePrefix)

      const countryLabel = isArabic ? c.nameAr : c.name
      const label = `\u200E${c.phonePrefix} — ${countryLabel}`
      options.push({ value: c.phonePrefix, label })
    }

    return options
  }, [countries, isArabic])

  return (
    <FormSelect<string>
      label={tCustomers('form.phone_code')}
      options={phoneCodeOptions}
      value={value}
      onChange={(val) => {
        onChange(val as string)
      }}
      required={required}
      disabled={disabled ?? !countriesReady}
      placeholder={placeholder ?? tCustomers('form.select_phone_code')}
      searchable={searchable}
      error={error}
    />
  )
}
