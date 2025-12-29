import { useQuery } from '@tanstack/react-query'
import { client } from './client'
import {
  ListCountriesResponseSchema,
  type ListCountriesResponse,
  type CountryMetadata,
} from './types/metadata'

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
    const response = await client.get('v1/metadata/countries').json()
    return ListCountriesResponseSchema.parse(response)
  },
}

/**
 * Query hook for countries metadata
 *
 * Fetches list of supported countries with currency info.
 * Data is cached for 1 hour as it rarely changes.
 */
export function useCountriesQuery() {
  return useQuery({
    queryKey: ['metadata', 'countries'],
    queryFn: () => metadataApi.listCountries(),
    staleTime: 60 * 60 * 1000, // 1 hour
    gcTime: 24 * 60 * 60 * 1000, // 24 hours
    select: (data) => data.countries,
  })
}

/**
 * Helper to get unique currencies from countries
 */
export function getUniqueCurrencies(countries: CountryMetadata[]) {
  const currencyMap = new Map<
    string,
    { code: string; name: string; symbol: string }
  >()

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
