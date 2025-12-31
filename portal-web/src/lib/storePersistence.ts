/**
 * TanStack Store Persistence Plugin
 *
 * Provides localStorage persistence for TanStack Store with TTL support,
 * automatic serialization/deserialization, and type-safe storage keys.
 *
 * @module storePersistence
 */

/**
 * Storage keys used by the application
 */
export const STORAGE_KEYS = {
  BUSINESS_PREFS: 'kyora_business_prefs',
  METADATA: 'kyora_metadata',
  ONBOARDING_SESSION: 'kyora_onboarding_session',
} as const

type StorageKey = (typeof STORAGE_KEYS)[keyof typeof STORAGE_KEYS]

/**
 * Data structure for persisted state with optional TTL
 */
interface PersistedState<T> {
  data: T
  lastFetched?: number
}

/**
 * Options for persistence plugin
 */
interface PersistenceOptions<T> {
  /**
   * Storage key to use in localStorage
   */
  key: StorageKey
  /**
   * Time-to-live in milliseconds. If provided, data older than this will be cleared.
   */
  ttl?: number
  /**
   * Function to serialize state before saving
   */
  serialize?: (state: T) => string
  /**
   * Function to deserialize state when loading
   */
  deserialize?: (value: string) => T
}

/**
 * Creates a persistence plugin for TanStack Store
 *
 * @example
 * ```ts
 * import { Store } from '@tanstack/react-store'
 * import { createPersistencePlugin, STORAGE_KEYS } from '@/lib/storePersistence'
 *
 * const store = new Store({
 *   businesses: [],
 *   selectedBusinessDescriptor: null,
 * })
 *
 * // Persist selected business with no TTL
 * createPersistencePlugin({
 *   key: STORAGE_KEYS.BUSINESS_PREFS,
 *   store,
 *   select: (state) => ({ selectedBusinessDescriptor: state.selectedBusinessDescriptor })
 * })
 *
 * // Persist metadata with 24-hour TTL
 * createPersistencePlugin({
 *   key: STORAGE_KEYS.METADATA,
 *   ttl: 24 * 60 * 60 * 1000,
 *   store,
 *   select: (state) => state
 * })
 * ```
 */
export function createPersistencePlugin<TState, TPersisted>(
  options: PersistenceOptions<TPersisted> & {
    /**
     * The store instance to persist
     */
    store: any
    /**
     * Function to select which part of state to persist
     */
    select: (state: TState) => TPersisted
    /**
     * Function to restore state from persisted data
     */
    restore?: (persisted: TPersisted, currentState: TState) => TState
  },
) {
  const {
    key,
    ttl,
    store,
    select,
    restore,
    serialize = JSON.stringify,
    deserialize = JSON.parse,
  } = options

  /**
   * Load persisted state from localStorage
   */
  function loadState(): TPersisted | null {
    try {
      const item = localStorage.getItem(key)
      if (!item) return null

      const persisted: PersistedState<TPersisted> = deserialize(item)

      // Check TTL if provided
      if (ttl && persisted.lastFetched) {
        const now = Date.now()
        const age = now - persisted.lastFetched

        if (age > ttl) {
          // Data is stale, clear it
          localStorage.removeItem(key)
          return null
        }
      }

      return persisted.data
    } catch (error) {
      console.error(`Failed to load persisted state from ${key}:`, error)
      return null
    }
  }

  /**
   * Save state to localStorage
   */
  function saveState(state: TPersisted) {
    try {
      const persisted: PersistedState<TPersisted> = {
        data: state,
        lastFetched: ttl ? Date.now() : undefined,
      }
      localStorage.setItem(key, serialize(persisted))
    } catch (error) {
      console.error(`Failed to save state to ${key}:`, error)
    }
  }

  /**
   * Clear persisted state
   */
  function clearState() {
    try {
      localStorage.removeItem(key)
    } catch (error) {
      console.error(`Failed to clear state from ${key}:`, error)
    }
  }

  // Load initial state
  const persistedState = loadState()
  if (persistedState && restore) {
    const currentState = store.state
    const restoredState = restore(persistedState, currentState)
    store.setState(() => restoredState)
  }

  // Subscribe to store changes and persist
  const unsubscribe = store.subscribe(() => {
    const currentState = store.state
    const stateToPersist = select(currentState)
    saveState(stateToPersist)
  })

  // Return cleanup and utility functions
  return {
    unsubscribe,
    clearState,
    loadState,
  }
}

/**
 * Helper to check if persisted data is stale
 */
export function isStateStale(key: StorageKey, ttl: number): boolean {
  try {
    const item = localStorage.getItem(key)
    if (!item) return true

    const persisted: PersistedState<any> = JSON.parse(item)

    if (!persisted.lastFetched) return false

    const now = Date.now()
    const age = now - persisted.lastFetched

    return age > ttl
  } catch {
    return true
  }
}

/**
 * Helper to clear all Kyora-related localStorage data
 */
export function clearAllPersistedData() {
  Object.values(STORAGE_KEYS).forEach((key) => {
    try {
      localStorage.removeItem(key)
    } catch (error) {
      console.error(`Failed to clear ${key}:`, error)
    }
  })
}
