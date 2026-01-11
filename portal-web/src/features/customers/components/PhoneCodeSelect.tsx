import { useEffect, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useStore } from '@tanstack/react-store'
import type { FormSelectOption } from '@/components/form/FormSelect'
import { metadataStore } from '@/stores/metadataStore'
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
  const { t, i18n } = useTranslation()
  const countries = useStore(metadataStore, (state) => state.countries)
  const countriesStatus = useStore(metadataStore, (state) => state.status)

  const isArabic = i18n.language.toLowerCase().startsWith('ar')
  const countriesReady = countries.length > 0 || countriesStatus === 'loaded'

  // Load countries on mount
  useEffect(() => {
    if (countriesStatus === 'idle') {
      void metadataStore.state.loadCountries()
    }
  }, [countriesStatus])

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
      label={t('customers.form.phone_code')}
      options={phoneCodeOptions}
      value={value}
      onChange={(val) => {
        onChange(val as string)
      }}
      required={required}
      disabled={disabled ?? !countriesReady}
      placeholder={placeholder ?? t('customers.form.select_phone_code')}
      searchable={searchable}
      error={error}
    />
  )
}
