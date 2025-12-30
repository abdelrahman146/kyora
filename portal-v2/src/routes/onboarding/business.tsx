import { useEffect, useMemo } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useForm } from '@tanstack/react-form'
import { useTranslation } from 'react-i18next'
import toast from 'react-hot-toast'
import { ArrowLeft, Building2, Loader2 } from 'lucide-react'
import {
  onboardingStore,
  setBusiness,
  updateStage,
} from '@/stores/onboardingStore'
import { useSetBusinessMutation } from '@/api/onboarding'
import { BusinessSetupSchema } from '@/schemas/onboarding'
import { translateErrorAsync } from '@/lib/translateError'
import { getUniqueCurrencies, useCountriesQuery } from '@/api/metadata'
import { useLanguage } from '@/hooks/useLanguage'

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
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { language } = useLanguage()
  const state = useStore(onboardingStore)
  const setBusinessMutation = useSetBusinessMutation()

  // Fetch countries and currencies
  const {
    data: countries = [],
    isLoading: isLoadingCountries,
    isError: isCountriesError,
  } = useCountriesQuery()

  // Extract unique currencies from countries
  const currencies = useMemo(() => {
    return countries.length > 0 ? getUniqueCurrencies(countries) : []
  }, [countries])

  // Sort countries by name based on language
  const sortedCountries = useMemo(() => {
    return [...countries].sort((a, b) => {
      const nameA = language === 'ar' ? a.nameAr : a.name
      const nameB = language === 'ar' ? b.nameAr : b.name
      return nameA.localeCompare(nameB, language)
    })
  }, [countries, language])

  // Redirect if no session or not verified
  useEffect(() => {
    if (!state.sessionToken) {
      navigate({ to: '/onboarding/email', replace: true })
    } else if (state.stage === 'identity_pending') {
      navigate({ to: '/onboarding/verify', replace: true })
    }
  }, [state.sessionToken, state.stage, navigate])

  const form = useForm({
    defaultValues: {
      name: state.businessData?.name || '',
      descriptor: state.businessData?.descriptor || '',
      country: state.businessData?.country || '',
      currency: state.businessData?.currency || '',
    },
    onSubmit: async ({ value }) => {
      try {
        const response = await setBusinessMutation.mutateAsync({
          sessionToken: state.sessionToken!,
          businessName: value.name,
          businessDescriptor: value.descriptor,
          country: value.country,
          currency: value.currency,
        })

        setBusiness(value)
        updateStage(response.stage)
        toast.success(t('onboarding:business_setup_success'))

        // Navigate based on plan type
        if (state.isPaidPlan) {
          await navigate({ to: '/onboarding/payment' })
        } else {
          await navigate({ to: '/onboarding/complete' })
        }
      } catch (error) {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      }
    },
    validators: {
      onBlur: BusinessSetupSchema,
    },
  })

  return (
    <div className="max-w-lg mx-auto">
      {/* Header */}
      <div className="text-center mb-8">
        <div className="flex justify-center mb-4">
          <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
            <Building2 className="w-8 h-8 text-primary" />
          </div>
        </div>
        <h1 className="text-3xl font-bold text-base-content mb-2">
          {t('onboarding:setup_business')}
        </h1>
        <p className="text-base-content/70">
          {t('onboarding:business_setup_description')}
        </p>
      </div>

      {/* Form */}
      <form
        onSubmit={(e) => {
          e.preventDefault()
          e.stopPropagation()
          void form.handleSubmit()
        }}
        className="space-y-6"
      >
        {/* Business Name */}
        <form.Field
          name="name"
          validators={{
            onBlur: BusinessSetupSchema.shape.name,
          }}
        >
          {(field) => (
            <div className="form-control">
              <label htmlFor="name" className="label">
                <span className="label-text font-medium">
                  {t('onboarding:business_name')}
                </span>
              </label>
              <input
                id="name"
                name="name"
                type="text"
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
                className={`input input-bordered w-full ${
                  field.state.meta.errors.length > 0 ? 'input-error' : ''
                }`}
                placeholder={t('onboarding:business_name_placeholder')}
              />
              {field.state.meta.errors.length > 0 && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {field.state.meta.errors[0]?.message || 'Invalid value'}
                  </span>
                </label>
              )}
            </div>
          )}
        </form.Field>

        {/* Business Descriptor */}
        <form.Field
          name="descriptor"
          validators={{
            onBlur: BusinessSetupSchema.shape.descriptor,
          }}
        >
          {(field) => (
            <div className="form-control">
              <label htmlFor="descriptor" className="label">
                <span className="label-text font-medium">
                  {t('onboarding:business_descriptor')}
                </span>
              </label>
              <input
                id="descriptor"
                name="descriptor"
                type="text"
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => {
                  // Convert to lowercase and replace spaces with hyphens
                  const value = e.target.value
                    .toLowerCase()
                    .replace(/\s+/g, '-')
                    .replace(/[^a-z0-9-]/g, '')
                  field.handleChange(value)
                }}
                className={`input input-bordered w-full font-mono ${
                  field.state.meta.errors.length > 0 ? 'input-error' : ''
                }`}
                placeholder={t('onboarding:business_descriptor_placeholder')}
              />
              <label className="label">
                <span className="label-text-alt text-base-content/60">
                  {t('onboarding:business_descriptor_help')}
                </span>
              </label>
              {field.state.meta.errors.length > 0 && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {field.state.meta.errors[0]?.message || 'Invalid value'}
                  </span>
                </label>
              )}
            </div>
          )}
        </form.Field>

        {/* Country */}
        <form.Field
          name="country"
          validators={{
            onBlur: BusinessSetupSchema.shape.country,
          }}
        >
          {(field) => (
            <div className="form-control">
              <label htmlFor="country" className="label">
                <span className="label-text font-medium">
                  {t('onboarding:country')}
                </span>
              </label>
              <select
                id="country"
                name="country"
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
                disabled={isLoadingCountries}
                className={`select select-bordered w-full ${
                  field.state.meta.errors.length > 0 ? 'select-error' : ''
                }`}
              >
                <option value="">
                  {isLoadingCountries
                    ? t('common:loading')
                    : t('onboarding:select_country')}
                </option>
                {sortedCountries.map((country) => (
                  <option key={country.code} value={country.code}>
                    {language === 'ar' ? country.nameAr : country.name}
                  </option>
                ))}
              </select>
              {isCountriesError && (
                <label className="label">
                  <span className="label-text-alt text-warning">
                    {t('common:load_error')}
                  </span>
                </label>
              )}
              {field.state.meta.errors.length > 0 && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {field.state.meta.errors[0]?.message || 'Invalid value'}
                  </span>
                </label>
              )}
            </div>
          )}
        </form.Field>

        {/* Currency */}
        <form.Field
          name="currency"
          validators={{
            onBlur: BusinessSetupSchema.shape.currency,
          }}
        >
          {(field) => (
            <div className="form-control">
              <label htmlFor="currency" className="label">
                <span className="label-text font-medium">
                  {t('onboarding:currency')}
                </span>
              </label>
              <select
                id="currency"
                name="currency"
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
                disabled={isLoadingCountries}
                className={`select select-bordered w-full ${
                  field.state.meta.errors.length > 0 ? 'select-error' : ''
                }`}
              >
                <option value="">
                  {isLoadingCountries
                    ? t('common:loading')
                    : t('onboarding:select_currency')}
                </option>
                {currencies.map((currency) => (
                  <option key={currency.code} value={currency.code}>
                    {currency.code} - {currency.name} ({currency.symbol})
                  </option>
                ))}
              </select>
              {field.state.meta.errors.length > 0 && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {field.state.meta.errors[0]?.message || 'Invalid value'}
                  </span>
                </label>
              )}
            </div>
          )}
        </form.Field>

        {/* Submit Button */}
        <form.Subscribe
          selector={(formState) => ({
            canSubmit: formState.canSubmit,
            isSubmitting: formState.isSubmitting,
          })}
        >
          {({ canSubmit, isSubmitting }) => (
            <button
              type="submit"
              disabled={
                !canSubmit || isSubmitting || setBusinessMutation.isPending
              }
              className="btn btn-primary w-full"
            >
              {(isSubmitting || setBusinessMutation.isPending) && (
                <Loader2 className="w-4 h-4 animate-spin" />
              )}
              {t('common:continue')}
            </button>
          )}
        </form.Subscribe>

        {/* Back Button */}
        <button
          type="button"
          onClick={() => navigate({ to: '/onboarding/verify' })}
          className="btn btn-ghost w-full"
        >
          <ArrowLeft className="w-4 h-4" />
          {t('common:back')}
        </button>
      </form>
    </div>
  )
}
