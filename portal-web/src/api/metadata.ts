import { queryOptions, useQuery } from '@tanstack/react-query'
import { get } from './client'
import { ListCountriesResponseSchema } from './types/metadata'
import type { CountryMetadata, ListCountriesResponse } from './types/metadata'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'

/**
 * Currency Info (derived from countries)
 */
export interface CurrencyInfo {
  code: string
  name: string
  symbol: string
}

/**
 * Extract unique currencies from countries
 */
function extractCurrencies(
  countries: Array<CountryMetadata>,
): Array<CurrencyInfo> {
  const currencyMap = new Map<string, CurrencyInfo>()

  for (const country of countries) {
    if (!currencyMap.has(country.currencyCode)) {
      currencyMap.set(country.currencyCode, {
        code: country.currencyCode,
        name: country.currencyLabel,
        symbol: country.currencySymbol,
      })
    }
  }

  return Array.from(currencyMap.values()).sort((a, b) =>
    a.code.localeCompare(b.code),
  )
}

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
 * Query Options Factories
 *
 * Co-locate query configuration (key + fn + staleTime) for type-safe reuse
 * in components, route loaders, and prefetching.
 */
export const metadataQueries = {
  /**
   * Query options for fetching countries
   *
   * Uses TanStack Query as the single source of truth with 24-hour cache.
   * No localStorage persistence - Query cache handles everything.
   */
  countries: () =>
    queryOptions({
      queryKey: queryKeys.metadata.countries(),
      queryFn: async () => {
        const response = await metadataApi.listCountries()
        return response
      },
      staleTime: STALE_TIME.TWENTY_FOUR_HOURS,
      gcTime: STALE_TIME.TWENTY_FOUR_HOURS,
      select: (data) => data.countries,
    }),

  /**
   * Query options for fetching currencies from countries data
   *
   * Derives currencies from countries endpoint.
   * Uses TanStack Query as the single source of truth.
   */
  currencies: () =>
    queryOptions({
      queryKey: queryKeys.metadata.currencies(),
      queryFn: async () => {
        const response = await metadataApi.listCountries()
        return extractCurrencies(response.countries)
      },
      staleTime: STALE_TIME.TWENTY_FOUR_HOURS,
      gcTime: STALE_TIME.TWENTY_FOUR_HOURS,
    }),
}

/**
 * Query hook for countries metadata
 *
 * Fetches list of supported countries with currency info.
 * TanStack Query cache: 24 hours (stale time + gc time).
 *
 * This is the single source of truth for countries data.
 */
export function useCountriesQuery() {
  return useQuery(metadataQueries.countries())
}

/**
 * Query hook for currencies metadata
 *
 * Derives unique currencies from countries endpoint.
 * TanStack Query cache: 24 hours (stale time + gc time).
 *
 * This is the single source of truth for currencies data.
 */
export function useCurrenciesQuery() {
  return useQuery(metadataQueries.currencies())
}

/**
 * Helper to get unique currencies from countries
 */
export function getUniqueCurrencies(countries: Array<CountryMetadata>) {
  return extractCurrencies(countries)
}
