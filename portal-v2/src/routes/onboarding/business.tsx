import { useEffect, useMemo, useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useTranslation } from 'react-i18next'
import { Building2, DollarSign, Globe } from 'lucide-react'
import type { ReactNode } from 'react'
import { useCountriesQuery } from '@/api/metadata'
import { onboardingApi } from '@/api/onboarding'
import { FormInput, FormSelect } from '@/components'
import type { CountryMetadata } from '@/api/types/metadata'
import {
  loadSessionFromStorage,
  onboardingStore,
  setBusiness,
  updateStage,
} from '@/stores/onboardingStore'
import { translateErrorAsync } from '@/lib/translateError'

export const Route = createFileRoute('/onboarding/business')({
  component: BusinessSetupPage,
})

/**
 * Business Setup Step - Step 4 of Onboarding
 *
 * Features:
 * - Business name input
 * - Business descriptor (slug) input with validation
 * - Country selection
 * - Currency selection
 */
function BusinessSetupPage() {
  const { t: tOnboarding, i18n } = useTranslation('onboarding')
  const { t: tCommon } = useTranslation('common')
  const { t: tTranslation } = useTranslation('translation')
  const navigate = useNavigate()
  const state = useStore(onboardingStore)

  const [didAttemptRestore, setDidAttemptRestore] = useState(false)
  const [businessName, setBusinessName] = useState(state.businessData?.name ?? '')
  const [descriptor, setDescriptor] = useState(
    state.businessData?.descriptor ?? '',
  )
  const [country, setCountry] = useState(state.businessData?.country ?? '')
  const [currency, setCurrency] = useState(state.businessData?.currency ?? '')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [descriptorError, setDescriptorError] = useState('')
  const [isDescriptorManuallyEdited, setIsDescriptorManuallyEdited] = useState(false)

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
        icon: <Globe className="w-4 h-4" />,
      }
    })
  }, [sortedCountries, isArabic])

  const currencyOptions = useMemo(() => {
    const seen = new Set<string>()
    const options: Array<{ value: string; label: string; icon: ReactNode }> = []
    for (const c of sortedCountries) {
      if (!c.currencyCode || seen.has(c.currencyCode)) continue
      seen.add(c.currencyCode)
      options.push({
        value: c.currencyCode,
        label: c.currencyLabel || c.currencyCode,
        icon: <DollarSign className="w-4 h-4" />,
      })
    }
    return options
  }, [sortedCountries])

  // Restore session from localStorage on mount
  useEffect(() => {
    const restoreSession = async () => {
      try {
        if (!state.sessionToken) {
          const hasSession = await loadSessionFromStorage()
          if (!hasSession) {
            void navigate({ to: '/onboarding/plan', replace: true })
          }
        }
      } finally {
        setDidAttemptRestore(true)
      }
    }

    void restoreSession()
  }, [navigate, state.sessionToken])

  // Redirect if not verified
  useEffect(() => {
    if (!didAttemptRestore) return

    if (!state.sessionToken) {
      void navigate({ to: '/onboarding/plan', replace: true })
      return
    }

    if (state.stage !== 'identity_verified') {
      void navigate({ to: '/onboarding/verify', replace: true })
    }
  }, [didAttemptRestore, navigate, state.sessionToken, state.stage])

  // Auto-generate descriptor from business name (only if not manually edited)
  useEffect(() => {
    if (businessName && !isDescriptorManuallyEdited) {
      const generated = businessName
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, '')
        .replace(/\s+/g, '-')
        .slice(0, 20)

      setDescriptor(generated)
    }
  }, [businessName, isDescriptorManuallyEdited])

  // Validate descriptor format
  useEffect(() => {
    if (descriptor) {
      if (descriptor.length < 3) {
        setDescriptorError(tOnboarding('business.descriptorTooShort'))
      } else if (!/^[a-z0-9-]+$/.test(descriptor)) {
        setDescriptorError(tOnboarding('business.descriptorInvalidFormat'))
      } else {
        setDescriptorError('')
      }
    } else {
      setDescriptorError('')
    }
  }, [descriptor, tOnboarding])

  const submitBusiness = async () => {
    setError('')

    if (!businessName.trim()) {
      setError(tOnboarding('business.nameRequired'))
      return
    }

    if (!descriptor.trim()) {
      setError(tOnboarding('business.descriptorRequired'))
      return
    }

    if (descriptorError) {
      setError(descriptorError)
      return
    }

    if (!country) {
      setError(tOnboarding('business.countryRequired'))
      return
    }

    if (!currency) {
      setError(tOnboarding('business.currencyRequired'))
      return
    }

    if (!state.sessionToken) return

    try {
      setIsSubmitting(true)

      const response = await onboardingApi.setBusiness({
        sessionToken: state.sessionToken,
        businessName: businessName.trim(),
        businessDescriptor: descriptor.trim(),
        country,
        currency,
      })

      setBusiness({
        name: businessName.trim(),
        descriptor: descriptor.trim(),
        country,
        currency,
      })
      updateStage(response.stage)

      if (response.stage === 'ready_to_commit') {
        void navigate({ to: '/onboarding/complete' })
      } else if (response.stage === 'business_staged' || state.isPaidPlan) {
        void navigate({ to: '/onboarding/payment' })
      } else {
        void navigate({ to: '/onboarding/complete' })
      }
    } catch (err) {
      const message = await translateErrorAsync(err, tTranslation)
      setError(message)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault()
    void submitBusiness()
  }

  return (
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

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Business Name */}
            <FormInput
              label={tOnboarding('business.name')}
              type="text"
              value={businessName}
              onChange={(e) => {
                setBusinessName(e.target.value)
              }}
              placeholder={tOnboarding('business.namePlaceholder')}
              required
              disabled={isSubmitting}
              startIcon={<Building2 className="w-5 h-5" />}
              helperText={tOnboarding('business.nameHint')}
              autoFocus
            />

            {/* Business Descriptor */}
            <FormInput
              label={tOnboarding('business.descriptor')}
              type="text"
              value={descriptor}
              onChange={(e) => {
                setDescriptor(e.target.value.toLowerCase())
                setIsDescriptorManuallyEdited(true)
              }}
              onBlur={() => {
                if (!descriptor.trim()) {
                  setIsDescriptorManuallyEdited(false)
                }
              }}
              placeholder={tOnboarding('business.descriptorPlaceholder')}
              pattern="[a-z0-9-]+"
              minLength={3}
              maxLength={20}
              required
              disabled={isSubmitting}
              helperText={tOnboarding('business.descriptorHint')}
              error={descriptorError}
            />

            {/* Country */}
            <FormSelect
              label={tOnboarding('business.country')}
              value={country}
              onChange={(value) => {
                const next = Array.isArray(value) ? value[0] ?? '' : value
                setCountry(next)
                const selected = countryByCode.get(next)
                if (selected?.currencyCode) {
                  setCurrency(selected.currencyCode)
                }
              }}
              options={countryOptions}
              placeholder={tOnboarding('business.selectCountry')}
              disabled={isSubmitting || isLoadingCountries || isCountriesError}
              required
            />

            {/* Currency */}
            <FormSelect
              label={tOnboarding('business.currency')}
              value={currency}
              onChange={(value) => {
                const next = Array.isArray(value) ? value[0] ?? '' : value
                setCurrency(next)
              }}
              options={currencyOptions}
              placeholder={tOnboarding('business.selectCurrency')}
              disabled={isSubmitting || isLoadingCountries || isCountriesError}
              required
            />

            {isCountriesError && (
              <div className="alert alert-error">
                <span className="text-sm">{tTranslation('errors:generic.unexpected')}</span>
              </div>
            )}

            {error && (
              <div className="alert alert-error">
                <span className="text-sm">{error}</span>
              </div>
            )}

            <button
              type="submit"
              className="btn btn-primary btn-block"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <>
                  <span className="loading loading-spinner loading-sm"></span>
                  {tCommon('loading')}
                </>
              ) : (
                tCommon('continue')
              )}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}
