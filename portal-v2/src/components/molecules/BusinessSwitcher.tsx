import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useStore } from '@tanstack/react-store'
import { Building2, Check, ChevronDown, Loader2 } from 'lucide-react'
import { Avatar } from '../atoms/Avatar'
import { Dropdown } from '../atoms/Dropdown'
import { businessStore } from '../../stores/businessStore'
import { businessApi } from '../../api/business'
import { useMediaQuery } from '../../hooks/useMediaQuery'
import { cn } from '../../lib/utils'
import type { Business } from '../../stores/businessStore'

/**
 * BusinessSwitcher Component
 *
 * Dropdown in the header that allows users to switch between businesses.
 * Fetches businesses from API and updates global selectedBusinessId state on selection.
 *
 * Features:
 * - Lists all available businesses
 * - Shows current selected business
 * - Updates global state (selectedBusiness & selectedBusinessId) on selection
 * - Loading and error states
 * - Mock data fallback for development
 * - Responsive: Compact on mobile, full info on desktop
 * - RTL support
 * - Touch-friendly dropdown (min 44px targets)
 *
 * @example
 * ```tsx
 * <BusinessSwitcher />
 * ```
 */
export function BusinessSwitcher() {
  const { t } = useTranslation()
  const isMobile = useMediaQuery('(max-width: 640px)')
  const businesses = useStore(businessStore, (state) => state.businesses)
  const selectedBusinessDescriptor = useStore(businessStore, (state) => state.selectedBusinessDescriptor)
  const selectedBusiness = businesses.find(b => b.descriptor === selectedBusinessDescriptor)
  const [isLoading, setIsLoading] = useState(false)

  // Fetch businesses on mount if not already loaded
  useEffect(() => {
    const loadBusinesses = async () => {
      // If already loaded, skip
      if (businesses.length > 0) return

      setIsLoading(true)
      try {
        const fetchedBusinesses = await businessApi.listBusinesses()
        businessStore.setState((state) => ({ ...state, businesses: fetchedBusinesses.businesses }))

        // Auto-select first business if none selected
        if (!selectedBusinessDescriptor && fetchedBusinesses.businesses.length > 0) {
          businessStore.setState((state) => ({ ...state, selectedBusinessDescriptor: fetchedBusinesses.businesses[0].descriptor }))
        }
      } catch {
        // Silent fail - component will handle empty state
      } finally {
        setIsLoading(false)
      }
    }

    void loadBusinesses()
  }, [businesses.length, selectedBusinessDescriptor])

  const handleSelectBusiness = (businessDescriptor: string) => {
    businessStore.setState((state) => ({ ...state, selectedBusinessDescriptor: businessDescriptor }))
  }

  const handleSelectBusinessById = (businessId: string) => {
    const business = businesses.find((b) => b.id === businessId)
    if (business) {
      handleSelectBusiness(business.descriptor)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 px-3 py-2 bg-base-200 rounded-lg">
        <Loader2 size={16} className="animate-spin text-base-content/60" />
        {!isMobile && (
          <span className="text-sm text-base-content/60">
            {t('common.loading')}
          </span>
        )}
      </div>
    )
  }

  if (businesses.length === 0) {
    return null
  }

  return (
    <Dropdown
      align="end"
      width="18rem"
      trigger={
        <button
          type="button"
          className={cn(
            'btn btn-ghost h-auto min-h-0 py-2 hover:bg-base-200 transition-colors',
            isMobile ? 'gap-1 px-2' : 'gap-2 px-3'
          )}
        >
          <Avatar
            src={selectedBusiness?.logoUrl}
            fallback={selectedBusiness?.name}
            size="sm"
            shape="square"
          />
          {!isMobile && (
            <div className="flex flex-col items-start max-w-40">
              <span className="text-sm font-semibold text-base-content truncate">
                {selectedBusiness?.name ?? t('dashboard.select_business')}
              </span>
              {/* Brand display commented out until Business interface has brand property */}
              {/* {selectedBusiness?.brand && (
                <span className="text-xs text-base-content/60 truncate">
                  {selectedBusiness.brand}
                </span>
              )} */}
            </div>
          )}
          <ChevronDown
            size={16}
            className="transition-transform text-base-content/60 shrink-0"
          />
        </button>
      }
    >
      <div className="py-2 max-h-96 overflow-y-auto">
        {/* Business List */}
        <div className="space-y-1">
          {businesses.map((business: Business) => (
            <button
              key={business.id}
              type="button"
              onClick={() => {
                handleSelectBusinessById(business.id)
              }}
              className={cn(
                'flex items-center gap-3 w-full px-4 py-3 text-start transition-colors',
                'hover:bg-base-200',
                selectedBusiness?.id === business.id
                  ? 'bg-primary/10 text-primary'
                  : 'text-base-content'
              )}
            >
              <Avatar
                src={business.logoUrl}
                fallback={business.name}
                size="sm"
                shape="square"
                className="shrink-0"
              />
              <div className="flex-1 min-w-0">
                <div
                  className={cn(
                    'text-sm truncate',
                    selectedBusiness?.id === business.id && 'font-semibold'
                  )}
                >
                  {business.name}
                </div>
                {/* Brand display commented out until Business interface has brand property */}
                {/* {business.brand && (
                  <div className="text-xs text-base-content/60 truncate">
                    {business.brand}
                  </div>
                )} */}
              </div>
              {selectedBusiness?.id === business.id && (
                <Check size={18} className="text-primary shrink-0" />
              )}
            </button>
          ))}
        </div>

        {/* Divider */}
        <div className="divider my-1"></div>

        {/* Add Business Option */}
        <button
          type="button"
          className="flex items-center gap-3 w-full px-4 py-3 text-primary hover:bg-primary/10 transition-colors text-start"
        >
          <Building2 size={18} className="shrink-0" />
          <span className="text-sm font-semibold">
            {t('dashboard.add_business')}
          </span>
        </button>
      </div>
    </Dropdown>
  )
}
