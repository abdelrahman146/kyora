import { useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useMetadataStore } from "@/stores/metadataStore";
import { FormSelect } from "../atoms/FormSelect";
import type { FormSelectOption } from "../atoms/FormSelect";

export interface CountrySelectProps {
  value: string;
  onChange: (value: string) => void;
  error?: string;
  disabled?: boolean;
  required?: boolean;
  placeholder?: string;
  searchable?: boolean;
}

/**
 * Reusable country select component
 * Uses FormSelect with country metadata from store
 * Displays country flag and localized name
 * Supports search functionality
 */
export function CountrySelect({
  value,
  onChange,
  error,
  disabled,
  required,
  placeholder,
  searchable = true,
}: CountrySelectProps) {
  const { t, i18n } = useTranslation();

  const countries = useMetadataStore((s) => s.countries);
  const countriesStatus = useMetadataStore((s) => s.status);
  const loadCountries = useMetadataStore((s) => s.loadCountries);

  const isArabic = i18n.language.toLowerCase().startsWith("ar");
  const countriesReady = countries.length > 0 || countriesStatus === "loaded";

  // Load countries on mount
  useEffect(() => {
    void loadCountries();
  }, [loadCountries]);

  const countryOptions: FormSelectOption[] = useMemo(() => {
    return countries.map((c) => {
      const label = `${c.flag ? `${c.flag} ` : ""}${isArabic ? c.nameAr : c.name}`;
      return { value: c.code, label };
    });
  }, [countries, isArabic]);

  return (
    <FormSelect<string>
      label={t("customers.form.country")}
      options={countryOptions}
      value={value}
      onChange={(val) => {
        onChange(val as string);
      }}
      required={required}
      disabled={disabled ?? !countriesReady}
      placeholder={placeholder ?? t("customers.form.select_country")}
      searchable={searchable}
      error={error}
    />
  );
}
