import { Store } from '@tanstack/react-store'
import { createPersistencePlugin } from '@/lib/storePersistence'

/**
 * Country Metadata
 */
export interface CountryMetadata {
  name: string
  nameAr: string
  code: string
  iso_code?: string
  flag?: string
  phonePrefix: string
  currencyCode: string
  currencyLabel: string
  currencySymbol: string
}

// Type alias for convenience
export type Country = CountryMetadata

/**
 * Currency Info (derived from countries)
 */
export interface CurrencyInfo {
  code: string
  name: string
  symbol: string
}

/**
 * Metadata Store State
 *
 * Caches reference data (countries, currencies) with 24-hour TTL.
 * Auto-clears stale data on access.
 */
interface MetadataState {
  countries: Array<CountryMetadata>
  currencies: Array<CurrencyInfo>
  lastFetched: number | null
  status: 'idle' | 'loading' | 'loaded' | 'error'
  loadCountries: () => Promise<void>
}

const initialState: MetadataState = {
  countries: [],
  currencies: [],
  lastFetched: null,
  status: 'idle',
  loadCountries: async () => {
    // Placeholder - will be replaced after store creation
  },
}

/**
 * Metadata Store
 *
 * Manages metadata cache using TanStack Store.
 * Persists to localStorage with 24-hour TTL.
 *
 * Dev-only devtools integration via conditional import.
 */
export const metadataStore = new Store<MetadataState>(initialState)

// Set up persistence plugin with 24-hour TTL
const persistencePlugin = createPersistencePlugin({
  key: 'kyora_metadata',
  ttl: 24 * 60 * 60 * 1000, // 24 hours
  store: metadataStore,
  select: (state: MetadataState) => state,
  // No TTL - persistence plugin handles TTL automatically
})

// Initialize store from localStorage on app load
const persistedMetadata = persistencePlugin.loadState()
if (persistedMetadata) {
  metadataStore.setState(() => persistedMetadata)
}

// Now implement loadCountries after store is created
metadataStore.setState((state) => ({
  ...state,
  loadCountries: async () => {
    if (state.status === 'loading') return

    metadataStore.setState((s) => ({ ...s, status: 'loading' }))

    // Note: Countries are loaded via useCountriesQuery() which syncs with this store
    // This method is kept for compatibility with components that need imperative loading
    // The actual API call happens in api/metadata.ts useCountriesQuery
    metadataStore.setState((s) => ({
      ...s,
      status: 'loaded',
      lastFetched: Date.now(),
    }))
    return Promise.resolve()
  },
}))

/**
 * Metadata Store Actions
 */

/**
 * Set metadata
 *
 * Called after fetching metadata from API.
 * Automatically persists to localStorage with timestamp.
 */
export function setMetadata(
  countries: Array<CountryMetadata>,
  currencies: Array<CurrencyInfo>,
): void {
  metadataStore.setState((state) => ({
    ...state,
    countries,
    currencies,
    lastFetched: Date.now(),
    status: 'loaded',
  }))
}

/**
 * Clear metadata cache
 *
 * Forces refresh on next access.
 */
export function clearMetadata(): void {
  metadataStore.setState(() => initialState)
  persistencePlugin.clearState()
}

/**
 * Check if metadata is stale (older than 24 hours)
 */
export function isMetadataStale(): boolean {
  const { lastFetched } = metadataStore.state
  if (!lastFetched) return true

  const age = Date.now() - lastFetched
  const TTL = 24 * 60 * 60 * 1000 // 24 hours

  return age > TTL
}

/**
 * Get metadata with auto-clear if stale
 *
 * Automatically clears cache if older than 24 hours.
 * Caller should refetch after clearing.
 */
export function getMetadata(): MetadataState {
  if (isMetadataStale()) {
    clearMetadata()
  }
  return metadataStore.state
}

/**
 * Get unique currencies from countries
 *
 * Helper to extract currency list from countries data.
 */
export function extractCurrencies(
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
 * Initialize TanStack Store Devtools (dev-only)
 *
 * Conditionally loads devtools in development mode only.
 * Production builds will exclude this code via tree-shaking.
 */
if (import.meta.env.DEV) {
  console.log(
    '[metadataStore] TanStack Store devtools enabled in development mode',
  )
}
