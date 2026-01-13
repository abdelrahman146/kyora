import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useStore } from '@tanstack/react-store'
import { useNavigate } from '@tanstack/react-router'
import { Building2, Check, ChevronDown, Loader2 } from 'lucide-react'

import type { Business } from '@/stores/businessStore'
import { Avatar } from '@/components/atoms/Avatar'
import { Dropdown } from '@/components/molecules/Dropdown'
import { useBusinessesQuery } from '@/api/business'
import { useMediaQuery } from '@/hooks/useMediaQuery'
import { getThumbnailUrl } from '@/lib/assetUrl'
import { cn } from '@/lib/utils'
import { businessStore } from '@/stores/businessStore'

export function BusinessSwitcher() {
  const { t: tCommon } = useTranslation('common')
  const { t: tDashboard } = useTranslation('dashboard')
  const isMobile = useMediaQuery('(max-width: 640px)')
  const navigate = useNavigate()

  const businesses = useStore(businessStore, (state) => state.businesses)
  const selectedBusinessDescriptor = useStore(
    businessStore,
    (state) => state.selectedBusinessDescriptor,
  )
  const selectedBusiness = businesses.find(
    (b) => b.descriptor === selectedBusinessDescriptor,
  )

  const businessesQuery = useBusinessesQuery()

  useEffect(() => {
    if (businesses.length > 0) return
    const fetched = businessesQuery.data?.businesses ?? []
    if (fetched.length === 0) return

    businessStore.setState((state) => ({
      ...state,
      businesses: fetched,
      selectedBusinessDescriptor:
        state.selectedBusinessDescriptor ?? fetched[0]?.descriptor,
    }))
  }, [businesses.length, businessesQuery.data, selectedBusinessDescriptor])

  const isLoading = businesses.length === 0 && businessesQuery.isLoading

  const handleSelectBusiness = (businessDescriptor: string) => {
    businessStore.setState((state) => ({
      ...state,
      selectedBusinessDescriptor: businessDescriptor,
    }))

    void navigate({
      to: '/business/$businessDescriptor',
      params: { businessDescriptor },
    })
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
            {tCommon('loading')}
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
            isMobile ? 'gap-1 px-2' : 'gap-2 px-3',
          )}
        >
          <Avatar
            src={getThumbnailUrl(selectedBusiness?.logo)}
            fallback={selectedBusiness?.name}
            size="sm"
            shape="square"
          />
          {!isMobile && (
            <div className="flex flex-col items-start max-w-40">
              <span className="text-sm font-semibold text-base-content truncate">
                {selectedBusiness?.name ?? tDashboard('select_business')}
              </span>
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
        <div className="space-y-1">
          {businesses.map((business: Business) => (
            <button
              key={business.id}
              type="button"
              onClick={() => {
                handleSelectBusinessById(business.id)
              }}
              className={cn(
                'flex items-center gap-3 w-full px-4 py-3 text-start transition-colors hover:cursor-pointer rounded-md',
                'hover:bg-base-200',
                selectedBusiness?.id === business.id
                  ? 'bg-primary/10 text-primary'
                  : 'text-base-content',
              )}
            >
              <Avatar
                src={getThumbnailUrl(business.logo)}
                fallback={business.name}
                size="sm"
                shape="square"
                className="shrink-0"
              />
              <div className="flex-1 min-w-0">
                <div
                  className={cn(
                    'text-sm truncate',
                    selectedBusiness?.id === business.id && 'font-semibold',
                  )}
                >
                  {business.name}
                </div>
              </div>
              {selectedBusiness?.id === business.id && (
                <Check size={18} className="text-primary shrink-0" />
              )}
            </button>
          ))}
        </div>

        <div className="divider my-1"></div>

        <button
          type="button"
          className="flex items-center gap-3 w-full px-4 py-3 text-primary hover:bg-primary/10 transition-colors text-start hover:cursor-pointer"
        >
          <Building2 size={18} className="shrink-0" />
          <span className="text-sm font-semibold">
            {tDashboard('add_business')}
          </span>
        </button>
      </div>
    </Dropdown>
  )
}
