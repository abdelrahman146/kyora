import { useTranslation } from 'react-i18next'
import { useStore } from '@tanstack/react-store'
import { useEffect } from 'react'
import { ChevronLeft, ChevronRight, X } from 'lucide-react'
import { Logo } from '../atoms/Logo'
import { IconButton } from '../atoms/IconButton'
import { NavGroup } from './NavGroup'
import {
  businessStore,
  closeSidebar,
  getSelectedBusiness,
  initNavGroups,
  toggleSidebar,
} from '@/stores/businessStore'
import { useLanguage } from '@/hooks/useLanguage'
import { useMediaQuery } from '@/hooks/useMediaQuery'
import { cn } from '@/lib/utils'
import { navigationConfig } from '@/config/navigation'

export type SidebarProps = Record<string, never>

/**
 * Sidebar Component
 *
 * Responsive navigation sidebar:
 * - Desktop: Fixed vertical sidebar with collapse functionality
 * - Mobile: Drawer overlay (slides from start)
 *
 * Features:
 * - Collapsible groups with nested items
 * - Collapsible on desktop (full width ↔ icon-only)
 * - Drawer overlay on mobile with slide animation
 * - Active route highlighting
 * - RTL support with proper directional logic
 * - Smooth transitions
 * - Touch-friendly targets (min 44px)
 * - Business descriptor from store for navigation
 */
export function Sidebar(_props: SidebarProps) {
  const { t: tDashboard } = useTranslation('dashboard')
  const { t: tCommon } = useTranslation('common')
  const { isRTL } = useLanguage()
  const isDesktop = useMediaQuery('(min-width: 768px)')

  const sidebarCollapsed = useStore(businessStore, (s) => s.sidebarCollapsed)
  const sidebarOpen = useStore(businessStore, (s) => s.sidebarOpen)
  const selectedBusinessDescriptor = useStore(
    businessStore,
    (s) => s.selectedBusinessDescriptor,
  )
  const selectedBusiness = getSelectedBusiness()

  // Initialize nav group expanded state on mount
  useEffect(() => {
    initNavGroups(navigationConfig)
  }, [])

  // Build base path for navigation
  const basePath = selectedBusinessDescriptor
    ? `/business/${selectedBusinessDescriptor}`
    : ''

  return (
    <aside
      className={cn(
        'fixed top-0 h-screen bg-base-100 border-base-300 z-50 flex flex-col',
        // Desktop: Always visible, collapsible width, smooth transition
        isDesktop && 'start-0 border-e transition-all duration-300',
        isDesktop && !sidebarCollapsed && 'w-64',
        isDesktop && sidebarCollapsed && 'w-20',
        // Mobile: Drawer from start with slide animation
        !isDesktop && 'start-0 w-64 transition-transform duration-300',
        // Mobile animation states
        !isDesktop && !isRTL && !sidebarOpen && '-translate-x-full',
        !isDesktop && isRTL && !sidebarOpen && 'translate-x-full',
      )}
    >
      {/* Header Section */}
      <div
        className={cn(
          'h-16 flex items-center border-b border-base-300 shrink-0',
          isDesktop && sidebarCollapsed
            ? 'justify-center'
            : 'justify-between px-4',
        )}
      >
        {isDesktop && sidebarCollapsed ? (
          // Collapsed desktop: Center logo and chevron with gap
          <>
            <Logo size="md" showText={false} />
            <IconButton
              icon={isRTL ? ChevronLeft : ChevronRight}
              size="sm"
              variant="ghost"
              onClick={toggleSidebar}
              aria-label={tDashboard('toggle_sidebar')}
            />
          </>
        ) : (
          // Expanded: Logo on start, toggle/close on end
          <>
            <div className="flex-1">
              <Logo size="md" showText />
            </div>
            {isDesktop ? (
              <IconButton
                icon={isRTL ? ChevronRight : ChevronLeft}
                size="sm"
                variant="ghost"
                onClick={toggleSidebar}
                aria-label={tDashboard('toggle_sidebar')}
              />
            ) : (
              <IconButton
                icon={X}
                size="sm"
                variant="ghost"
                onClick={closeSidebar}
                aria-label={tCommon('close')}
              />
            )}
          </>
        )}
      </div>

      {/* Main Navigation - Scrollable */}
      <nav
        className="flex-1 p-2 space-y-1 overflow-y-auto min-h-0"
        onWheel={(e) => {
          // Prevent scroll from bubbling to parent when nav is scrollable
          const target = e.currentTarget
          const isScrollable = target.scrollHeight > target.clientHeight
          const isAtTop = target.scrollTop === 0
          const isAtBottom =
            target.scrollTop + target.clientHeight >= target.scrollHeight

          if (isScrollable) {
            // Stop propagation if scrolling within bounds
            if ((e.deltaY < 0 && !isAtTop) || (e.deltaY > 0 && !isAtBottom)) {
              e.stopPropagation()
            }
          }
        }}
      >
        {navigationConfig
          .filter((item) => {
            // Exclude settings separator and settings items from main nav
            if ('type' in item && item.type === 'separator') return false
            if ('key' in item && item.key.startsWith('settings.')) return false
            return true
          })
          .map((config) => (
            <NavGroup
              key={config.key}
              config={config}
              basePath={basePath}
              sidebarCollapsed={sidebarCollapsed}
              onItemClick={() => {
                // Close drawer on mobile after navigation
                if (!isDesktop) closeSidebar()
              }}
            />
          ))}
      </nav>

      {/* Bottom Section - Fixed */}
      <div className="border-t border-base-300 bg-base-100 shrink-0">
        {(!isDesktop || !sidebarCollapsed) && (
          <div className="px-4 py-3">
            <span className="text-xs font-semibold uppercase text-base-content/50">
              {tCommon('settings_section')}
            </span>
          </div>
        )}
        <div className="p-2 space-y-1">
          {navigationConfig
            .filter((item) => 'key' in item && item.key.startsWith('settings.'))
            .map((config) => (
              <NavGroup
                key={config.key}
                config={config}
                basePath={basePath}
                sidebarCollapsed={sidebarCollapsed}
                onItemClick={() => {
                  // Close drawer on mobile after navigation
                  if (!isDesktop) closeSidebar()
                }}
              />
            ))}
        </div>
      </div>
    </aside>
  )
}
