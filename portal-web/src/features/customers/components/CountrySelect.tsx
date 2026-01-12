import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import type { FormSelectOption } from '@/components/form/FormSelect'
import type { CountryMetadata } from '@/api/types/metadata'
import { useCountriesQuery } from '@/api/metadata'
import { FormSelect } from '@/components/form/FormSelect'

export interface CountrySelectProps {
  value: string
  onChange: (value: string) => void
  error?: string
  disabled?: boolean
  required?: boolean
  placeholder?: string
  searchable?: boolean
  availableCountries?: Array<CountryMetadata> // Optional filter: only show these countries
}

/**
 * Reusable country select component
 * Uses FormSelect with country metadata from store
 * Displays country flag and localized name
 * Supports search functionality
 * Can be filtered to show only specific countries via availableCountries prop
 */
export function CountrySelect({
  value,
  onChange,
  error,
  disabled,
  required,
  placeholder,
  searchable = true,
  availableCountries,
}: CountrySelectProps) {
  const { i18n } = useTranslation()
  const { t: tCustomers } = useTranslation('customers')
  const { data: countries = [], isSuccess } = useCountriesQuery()

  const isArabic = i18n.language.toLowerCase().startsWith('ar')
  const countriesReady = isSuccess && countries.length > 0

  const countryOptions: Array<FormSelectOption> = useMemo(() => {
    // Use filtered list if provided, otherwise use all countries
    const countriesToShow = availableCountries ?? countries
    return countriesToShow.map((c: CountryMetadata) => {
      const label = `${c.flag ? `${c.flag} ` : ''}${isArabic ? c.nameAr : c.name}`
      return { value: c.code, label }
    })
  }, [countries, availableCountries, isArabic])

  return (
    <FormSelect<string>
      label={tCustomers('form.country')}
      options={countryOptions}
      value={value}
      onChange={(val) => {
        onChange(val as string)
      }}
      required={required}
      disabled={disabled ?? !countriesReady}
      placeholder={placeholder ?? tCustomers('form.select_country')}
      searchable={searchable}
      error={error}
    />
  )
}
