import { useMemo } from 'react'
import { useLoaderData, useNavigate, useSearch } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { Building2 } from 'lucide-react'
import { z } from 'zod'

import type { CountryMetadata } from '@/api/types/metadata'
import { useCountriesQuery } from '@/api/metadata'
import { useSetBusinessMutation } from '@/api/onboarding'
import { useKyoraForm } from '@/lib/form'
import { useLanguage } from '@/hooks/useLanguage'

import { OnboardingLayout } from '@/features/onboarding/components/OnboardingLayout'

export function BusinessSetupPage() {
  const { t: tOnboarding } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const { t: tErrors } = useTranslation('errors')
  const navigate = useNavigate()

  const { session } = useLoaderData({ from: '/onboarding/business' })
  const { session: sessionToken } = useSearch({ from: '/onboarding/business' })

  const {
    data: countries = [],
    isLoading: isLoadingCountries,
    isError: isCountriesError,
  } = useCountriesQuery()

  const { isArabic, language } = useLanguage()

  const sortedCountries = useMemo(() => {
    return [...countries].sort((a, b) => {
      const nameA = isArabic ? a.nameAr : a.name
      const nameB = isArabic ? b.nameAr : b.name
      return nameA.localeCompare(nameB, language)
    })
  }, [countries, isArabic, language])

  const countryByCode = useMemo(() => {
    const map = new Map<string, CountryMetadata>()
    for (const c of sortedCountries) {
      map.set(c.code, c)
    }
    return map
  }, [sortedCountries])

  const countryOptions = useMemo(() => {
    return sortedCountries.map((c) => {
      const label = `${c.flag ? `${c.flag} ` : ''}${isArabic ? c.nameAr : c.name}`
      return {
        value: c.code,
        label,
      }
    })
  }, [sortedCountries, isArabic])

  const currencyOptions = useMemo(() => {
    const seen = new Set<string>()
    const options: Array<{ value: string; label: string }> = []
    for (const c of sortedCountries) {
      if (!c.currencyCode || seen.has(c.currencyCode)) continue
      seen.add(c.currencyCode)
      options.push({
        value: c.currencyCode,
        label: c.currencyLabel || c.currencyCode,
      })
    }
    return options
  }, [sortedCountries])

  const setBusinessMutation = useSetBusinessMutation({
    onSuccess: async (response) => {
      if (response.stage === 'ready_to_commit') {
        await navigate({
          to: '/onboarding/complete',
          search: { session: sessionToken },
        })
      } else if (response.stage === 'business_staged') {
        await navigate({
          to: '/onboarding/payment',
          search: { session: sessionToken },
        })
      } else {
        await navigate({
          to: '/onboarding/complete',
          search: { session: sessionToken },
        })
      }
    },
  })

  const form = useKyoraForm({
    defaultValues: {
      name: session.businessName ?? '',
      descriptor: session.businessDescriptor ?? '',
      country: '',
      currency: '',
    },
    listeners: {
      onChange: ({ formApi, fieldApi }) => {
        if (fieldApi.name === 'name') {
          const businessName = fieldApi.state.value as string
          const currentDescriptor = formApi.getFieldValue('descriptor')

          const expectedDescriptor = businessName
            .toLowerCase()
            .replace(/[^a-z0-9\s-]/g, '')
            .replace(/\s+/g, '-')
            .slice(0, 20)

          if (
            businessName &&
            (!currentDescriptor ||
              currentDescriptor ===
                expectedDescriptor.slice(0, currentDescriptor.length))
          ) {
            formApi.setFieldValue('descriptor', expectedDescriptor)
          }
        }

        if (fieldApi.name === 'country') {
          const country = fieldApi.state.value as string
          const currentCurrency = formApi.getFieldValue('currency')

          if (country && !currentCurrency) {
            const selected = countryByCode.get(country)
            if (selected?.currencyCode) {
              formApi.setFieldValue('currency', selected.currencyCode)
            }
          }
        }
      },
    },
    onSubmit: async ({
      value,
    }: {
      value: {
        name: string
        descriptor: string
        country: string
        currency: string
      }
    }) => {
      await setBusinessMutation.mutateAsync({
        sessionToken,
        businessName: value.name,
        businessDescriptor: value.descriptor,
        country: value.country,
        currency: value.currency,
      })
    },
  })

  return (
    <OnboardingLayout>
      <div className="max-w-2xl mx-auto">
        <div className="card bg-base-100 border border-base-300">
          <div className="card-body">
            <div className="text-center mb-6">
              <div className="flex justify-center mb-4">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                  <Building2 className="w-8 h-8 text-primary" />
                </div>
              </div>
              <h2 className="text-2xl font-bold">
                {tOnboarding('business.title')}
              </h2>
              <p className="text-base-content/70 mt-2">
                {tOnboarding('business.subtitle')}
              </p>
            </div>

            <form.AppForm>
              <form.FormRoot className="space-y-6">
                <form.AppField
                  name="name"
                  validators={{
                    onBlur: z.string().min(1, 'validation.required'),
                  }}
                >
                  {(field) => (
                    <field.TextField
                      type="text"
                      label={tOnboarding('business.name')}
                      placeholder={tOnboarding('business.namePlaceholder')}
                      autoFocus
                      hint={tOnboarding('business.nameHint')}
                      startIcon={<Building2 className="w-5 h-5" />}
                    />
                  )}
                </form.AppField>

                <form.AppField
                  name="descriptor"
                  validators={{
                    onBlur: z.string().min(1, 'validation.required'),
                  }}
                >
                  {(field) => (
                    <field.TextField
                      type="text"
                      label={tOnboarding('business.descriptor')}
                      placeholder={tOnboarding(
                        'business.descriptorPlaceholder',
                      )}
                      hint={tOnboarding('business.descriptorHint')}
                    />
                  )}
                </form.AppField>

                <form.AppField
                  name="country"
                  validators={{
                    onBlur: z.string().min(1, 'validation.required'),
                  }}
                >
                  {(field) => (
                    <field.SelectField
                      label={tOnboarding('business.country')}
                      placeholder={tOnboarding('business.selectCountry')}
                      options={countryOptions}
                      disabled={isLoadingCountries || isCountriesError}
                    />
                  )}
                </form.AppField>

                <form.AppField
                  name="currency"
                  validators={{
                    onBlur: z.string().min(1, 'validation.required'),
                  }}
                >
                  {(field) => (
                    <field.SelectField
                      label={tOnboarding('business.currency')}
                      placeholder={tOnboarding('business.selectCurrency')}
                      options={currencyOptions}
                      disabled={isLoadingCountries || isCountriesError}
                    />
                  )}
                </form.AppField>

                {isCountriesError && (
                  <div className="alert alert-error">
                    <span className="text-sm">
                      {tErrors('generic.unexpected')}
                    </span>
                  </div>
                )}

                {setBusinessMutation.error && (
                  <div className="alert alert-error">
                    <span className="text-sm">
                      {setBusinessMutation.error.message}
                    </span>
                  </div>
                )}

                <form.SubmitButton
                  variant="primary"
                  size="lg"
                  fullWidth
                  disabled={setBusinessMutation.isPending}
                >
                  {tCommon('continue')}
                </form.SubmitButton>
              </form.FormRoot>
            </form.AppForm>
          </div>
        </div>
      </div>
    </OnboardingLayout>
  )
}
