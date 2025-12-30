import { useEffect, useMemo } from 'react'
import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { Building2 } from 'lucide-react'
import { z } from 'zod'
import type { RouterContext } from '@/router'
import type { CountryMetadata } from '@/api/types/metadata'
import { useCountriesQuery } from '@/api/metadata'
import { onboardingQueries, useSetBusinessMutation } from '@/api/onboarding'

import { OnboardingLayout } from '@/components/templates/OnboardingLayout'
import { useKyoraForm } from '@/lib/form'
import { SelectField, TextField } from '@/lib/form/components'

// Search params schema
const BusinessSearchSchema = z.object({
  sessionToken: z.string().min(1),
})

export const Route = createFileRoute('/onboarding/business')({
  validateSearch: (search): z.infer<typeof BusinessSearchSchema> => {
    return BusinessSearchSchema.parse(search)
  },
  loader: async ({ context, location }) => {
    const parsed = BusinessSearchSchema.parse(location.search)
    const { queryClient } = context as RouterContext
    
    // Redirect if no session token
    if (!parsed.sessionToken) {
      throw redirect({ to: '/onboarding/plan' })
    }

    // Prefetch and validate session
    const session = await queryClient.ensureQueryData(
      onboardingQueries.session(parsed.sessionToken)
    )

    // Redirect if wrong stage
    if (
      session.stage !== 'identity_verified' &&
      session.stage !== 'business_staged'
    ) {
      if (session.stage === 'plan_selected') {
        throw redirect({
          to: '/onboarding/verify',
          search: { sessionToken: parsed.sessionToken },
        })
      }
    }

    return { session }
  },
  component: BusinessSetupPage,
})

/**
 * Business Setup Step - Step 4 of Onboarding
 *
 * Features:
 * - Business name input
 * - Business descriptor (slug) with auto-generation
 * - Country selection
 * - Currency selection (auto-selected from country)
 */
function BusinessSetupPage() {
  const { t: tOnboarding, i18n } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const { t: tTranslation } = useTranslation('translation')
  const navigate = useNavigate()
  const { session } = Route.useLoaderData()
  const { sessionToken } = Route.useSearch()

  // Fetch countries and currencies
  const {
    data: countries = [],
    isLoading: isLoadingCountries,
    isError: isCountriesError,
  } = useCountriesQuery()

  const isArabic = i18n.language.toLowerCase().startsWith('ar')
  const language = isArabic ? 'ar' : 'en'

  // Sort countries by name based on language
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

  // Set business mutation
  const setBusinessMutation = useSetBusinessMutation({
    onSuccess: async (response) => {
      // Navigate based on next stage
      if (response.stage === 'ready_to_commit') {
        await navigate({
          to: '/onboarding/complete',
          search: { session: sessionToken },
        })
      } else if (response.stage === 'business_staged') {
        // Check if plan requires payment
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

  // TanStack Form with Zod validation and field listeners
  const form = useKyoraForm({
    defaultValues: {
      name: session.businessName ?? '',
      descriptor: session.businessDescriptor ?? '',
      country: '',
      currency: '',
    },
    onSubmit: async ({ value }: { value: { name: string; descriptor: string; country: string; currency: string } }) => {
      await setBusinessMutation.mutateAsync({
        sessionToken,
        businessName: value.name,
        businessDescriptor: value.descriptor,
        country: value.country,
        currency: value.currency,
      })
    },
  })

  // Auto-generate descriptor from business name using field listener
  useEffect(() => {
    const unsubscribe = form.store.subscribe(() => {
      const businessName = form.getFieldValue('name')
      const descriptor = form.getFieldValue('descriptor')
      
      if (businessName && !descriptor) {
        const generated = businessName
          .toLowerCase()
          .replace(/[^a-z0-9\s-]/g, '')
          .replace(/\s+/g, '-')
          .slice(0, 20)

        form.setFieldValue('descriptor', generated)
      }
    })

    return unsubscribe
  }, [form])

  // Auto-select currency when country changes
  useEffect(() => {
    const unsubscribe = form.store.subscribe(() => {
      const country = form.getFieldValue('country')
      if (country) {
        const selected = countryByCode.get(country)
        const currency = form.getFieldValue('currency')
        if (selected?.currencyCode && !currency) {
          form.setFieldValue('currency', selected.currencyCode)
        }
      }
    })

    return unsubscribe
  }, [form, countryByCode])

  return (
    <OnboardingLayout>
      <div className="max-w-2xl mx-auto">
        <div className="card bg-base-100 border border-base-300 shadow-xl">
        <div className="card-body">
          {/* Header */}
          <div className="text-center mb-6">
            <div className="flex justify-center mb-4">
              <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                <Building2 className="w-8 h-8 text-primary" />
              </div>
            </div>
            <h2 className="text-2xl font-bold">{tOnboarding('business.title')}</h2>
            <p className="text-base-content/70 mt-2">
              {tOnboarding('business.subtitle')}
            </p>
          </div>

          <form.FormRoot className="space-y-6">
            {/* Business Name */}
            <form.Field 
              name="name"
              validators={{
                onBlur: z.string().min(1, 'validation.required'),
              }}
            >
              {() => (
                <TextField
                  type="text"
                  label={tOnboarding('business.name')}
                  placeholder={tOnboarding('business.namePlaceholder')}
                  autoFocus
                  hint={tOnboarding('business.nameHint')}
                  startIcon={<Building2 className="w-5 h-5" />}
                />
              )}
            </form.Field>

            {/* Business Descriptor */}
            <form.Field 
              name="descriptor"
              validators={{
                onBlur: z.string().min(1, 'validation.required'),
              }}
            >
              {() => (
                <TextField
                  type="text"
                  label={tOnboarding('business.descriptor')}
                  placeholder={tOnboarding('business.descriptorPlaceholder')}
                  hint={tOnboarding('business.descriptorHint')}
                />
              )}
            </form.Field>

            {/* Country */}
            <form.Field 
              name="country"
              validators={{
                onBlur: z.string().min(1, 'validation.required'),
              }}
            >
              {() => (
                <SelectField
                  label={tOnboarding('business.country')}
                  placeholder={tOnboarding('business.selectCountry')}
                  options={countryOptions}
                  disabled={isLoadingCountries || isCountriesError}
                />
              )}
            </form.Field>

            {/* Currency */}
            <form.Field 
              name="currency"
              validators={{
                onBlur: z.string().min(1, 'validation.required'),
              }}
            >
              {() => (
                <SelectField
                  label={tOnboarding('business.currency')}
                  placeholder={tOnboarding('business.selectCurrency')}
                  options={currencyOptions}
                  disabled={isLoadingCountries || isCountriesError}
                />
              )}
            </form.Field>

            {isCountriesError && (
              <div className="alert alert-error">
                <span className="text-sm">
                  {tTranslation('errors:generic.unexpected')}
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
        </div>
      </div>
      </div>
    </OnboardingLayout>
  )
}
