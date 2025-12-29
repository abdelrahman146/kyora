import { useQuery } from '@tanstack/react-query'
import { get } from './client'
import { ListCountriesResponseSchema } from './types/metadata'
import type { CountryMetadata, ListCountriesResponse } from './types/metadata'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'
import {
  extractCurrencies,
  getMetadata,
  setMetadata,
} from '@/stores/metadataStore'

/**
 * Metadata API Client
 *
 * Endpoints for fetching countries, currencies, and other metadata
 */
export const metadataApi = {
  /**
   * List all supported countries
   */
  async listCountries(): Promise<ListCountriesResponse> {
    const response = await get<unknown>('v1/metadata/countries')
    return ListCountriesResponseSchema.parse(response)
  },
}

/**
 * Query hook for countries metadata with dual-layer caching
 *
 * Fetches list of supported countries with currency info.
 * - TanStack Query cache: 24 hours
 * - localStorage persistence: 24 hours (via metadataStore)
 *
 * Automatically syncs with metadataStore for additional persistence layer.
 */
export function useCountriesQuery() {
  return useQuery({
    queryKey: queryKeys.metadata.countries(),
    queryFn: async () => {
      const response = await metadataApi.listCountries()
      const currencies = extractCurrencies(response.countries)

      // Sync with metadataStore for additional persistence layer
      setMetadata(response.countries, currencies)

      return response
    },
    staleTime: STALE_TIME.TWENTY_FOUR_HOURS,
    gcTime: STALE_TIME.TWENTY_FOUR_HOURS,
    select: (data) => data.countries,
    // Try to use persisted data on mount
    initialData: () => {
      const metadata = getMetadata()
      if (metadata.countries.length > 0) {
        return { countries: metadata.countries }
      }
      return undefined
    },
  })
}

/**
 * Query hook for currencies metadata
 *
 * Extracts unique currencies from countries data.
 * Uses same caching strategy as countries query.
 */
export function useCurrenciesQuery() {
  return useQuery({
    queryKey: queryKeys.metadata.currencies(),
    queryFn: async () => {
      const response = await metadataApi.listCountries()
      const currencies = extractCurrencies(response.countries)

      // Sync with metadataStore
      setMetadata(response.countries, currencies)

      return currencies
    },
    staleTime: STALE_TIME.TWENTY_FOUR_HOURS,
    gcTime: STALE_TIME.TWENTY_FOUR_HOURS,
    // Try to use persisted data on mount
    initialData: () => {
      const metadata = getMetadata()
      if (metadata.currencies.length > 0) {
        return metadata.currencies
      }
      return undefined
    },
  })
}

/**
 * Helper to get unique currencies from countries
 */
export function getUniqueCurrencies(countries: Array<CountryMetadata>) {
  return extractCurrencies(countries)
}
