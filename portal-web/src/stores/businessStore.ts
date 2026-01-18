import { Store } from '@tanstack/react-store'

import type { Business } from '@/api/business'
import { createPersistencePlugin } from '@/lib/storePersistence'

/**
 * Business Store State
 *
 * Manages business list, UI preferences, and "last selected" business.
 *
 * IMPORTANT: The URL path param `$businessDescriptor` is the single source of
 * truth for the currently active business in business-scoped pages.
 * `selectedBusinessDescriptor` is persisted only as a convenience preference
 * for redirects (e.g., from homepage to last visited business).
 *
 * Only selectedBusinessDescriptor and sidebarCollapsed are persisted.
 */
interface BusinessState {
  businesses: Array<Business>
  selectedBusinessDescriptor: string | null // Convenience preference, NOT source of truth
  sidebarCollapsed: boolean // Desktop: collapsed (icon-only) state
  sidebarOpen: boolean // Mobile: drawer open state
}

const initialState: BusinessState = {
  businesses: [],
  selectedBusinessDescriptor: null,
  sidebarCollapsed: false,
  sidebarOpen: false,
}

/**
 * Business Store
 *
 * Manages business state using TanStack Store.
 * Persists selectedBusinessDescriptor (for convenience redirects) and
 * sidebarCollapsed (UI preference) to localStorage.
 *
 * IMPORTANT: Business-scoped routes use URL `$businessDescriptor` as the
 * single source of truth. This store's `selectedBusinessDescriptor` is synced
 * when navigating to business routes (via route loaders), but it is NOT the
 * authoritative source during business page renders.
 *
 * Dev-only devtools integration via conditional import.
 */
export const businessStore = new Store<BusinessState>(initialState)

// Set up persistence plugin for preferences only
const persistencePlugin = createPersistencePlugin({
  key: 'kyora_business_prefs',
  store: businessStore,
  select: (state: BusinessState) => ({
    selectedBusinessDescriptor: state.selectedBusinessDescriptor,
    sidebarCollapsed: state.sidebarCollapsed,
  }),
  restore: (persisted, currentState) => ({
    ...currentState,
    selectedBusinessDescriptor: persisted.selectedBusinessDescriptor,
    sidebarCollapsed: persisted.sidebarCollapsed,
  }),
  // No TTL - preferences persist until explicitly changed
})

// Initialize store from localStorage on app load
const persistedPrefs = persistencePlugin.loadState()
if (persistedPrefs) {
  businessStore.setState((state) => ({
    ...state,
    selectedBusinessDescriptor: persistedPrefs.selectedBusinessDescriptor,
    sidebarCollapsed: persistedPrefs.sidebarCollapsed,
  }))
}

/**
 * Business Store Actions
 */

/**
 * Set businesses list
 *
 * Called after fetching businesses from API.
 * Businesses are not persisted - always fetched fresh.
 */
export function setBusinesses(businesses: Array<Business>): void {
  businessStore.setState((state) => ({
    ...state,
    businesses,
  }))
}

/**
 * Select a business
 *
 * Updates selectedBusinessDescriptor and persists to localStorage.
 * This is called by business route loaders to sync the store preference with
 * the current URL `$businessDescriptor` param (the single source of truth).
 *
 * Note: Query invalidation should be handled by the caller
 * (typically in the business route loader after this action).
 *
 * @param descriptor - Business descriptor from URL to sync into store
 */
export function selectBusiness(descriptor: string): void {
  businessStore.setState((state) => ({
    ...state,
    selectedBusinessDescriptor: descriptor,
  }))
}

/**
 * Clear selected business
 *
 * Called on logout or when business is no longer accessible.
 */
export function clearSelectedBusiness(): void {
  businessStore.setState((state) => ({
    ...state,
    selectedBusinessDescriptor: null,
  }))
}

/**
 * Toggle sidebar collapsed state
 *
 * Persisted to localStorage for consistent UI across sessions.
 */
export function toggleSidebar(): void {
  businessStore.setState((state) => ({
    ...state,
    sidebarCollapsed: !state.sidebarCollapsed,
  }))
}

/**
 * Set sidebar collapsed state explicitly
 */
export function setSidebarCollapsed(collapsed: boolean): void {
  businessStore.setState((state) => ({
    ...state,
    sidebarCollapsed: collapsed,
  }))
}

/**
 * Open mobile sidebar drawer
 */
export function openSidebar(): void {
  businessStore.setState((state) => ({
    ...state,
    sidebarOpen: true,
  }))
}

/**
 * Close mobile sidebar drawer
 */
export function closeSidebar(): void {
  businessStore.setState((state) => ({
    ...state,
    sidebarOpen: false,
  }))
}

/**
 * Clear all business data
 *
 * Called on logout. Clears both state and persisted preferences.
 */
export function clearBusinessData(): void {
  businessStore.setState(() => initialState)
  persistencePlugin.clearState()
}

/**
 * Get currently selected business from store
 */
export function getSelectedBusiness(): Business | null {
  const { businesses, selectedBusinessDescriptor } = businessStore.state
  if (!selectedBusinessDescriptor) return null
  return (
    businesses.find((b) => b.descriptor === selectedBusinessDescriptor) ?? null
  )
}

// Re-export Business type for convenience
export type { Business } from '@/api/business'
