import { useEffect } from 'react'
import { useStore } from '@tanstack/react-store'
import { useMatches } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import type { ReactNode } from 'react'
import { Header } from '@/components/organisms/Header'
import { Sidebar } from '@/components/organisms/Sidebar'
import { BottomNav } from '@/components/organisms/BottomNav'
import { businessStore, closeSidebar } from '@/stores/businessStore'
import { useMediaQuery } from '@/hooks/useMediaQuery'
import { cn } from '@/lib/utils'

export interface DashboardLayoutProps {
  title?: string
  children: ReactNode
}

/**
 * DashboardLayout Component
 *
 * Mobile-first application layout with responsive behavior.
 *
 * Features:
 * - Mobile: Fixed header + Bottom nav + Drawer sidebar (overlay)
 * - Desktop: Fixed sidebar + Header (no bottom nav)
 * - Collapsible sidebar on desktop
 * - Overlay with backdrop on mobile when sidebar opens
 * - Safe area padding for mobile devices
 * - RTL support with logical properties
 * - Smooth transitions and animations
 * - Thumb-friendly touch targets (minimum 44px)
 *
 * Layout Structure:
 * ```
 * Desktop (≥768px):
 * ┌─────────────┬──────────────────────────┐
 * │   Sidebar   │         Header           │
 * │  (Fixed)    ├──────────────────────────┤
 * │             │                          │
 * │  Collapsible│        Content           │
 * │  64→20px    │      (Responsive)        │
 * │             │                          │
 * └─────────────┴──────────────────────────┘
 *
 * Mobile (<768px):
 * ┌──────────────────────────────────────┐
 * │         Header (Fixed)               │
 * ├──────────────────────────────────────┤
 * │                                      │
 * │          Content Area                │
 * │      (Scrollable, Safe)              │
 * │                                      │
 * ├──────────────────────────────────────┤
 * │       Bottom Nav (Fixed)             │
 * └──────────────────────────────────────┘
 *
 * Mobile Sidebar (Drawer Overlay):
 * ┌──────────────┬───────────────────────┐
 * │   Sidebar    │   Backdrop (Blur)     │
 * │   (Drawer)   │                       │
 * │   Slide-in   │   Tap to Close        │
 * │              │                       │
 * └──────────────┴───────────────────────┘
 * ```
 *
 * @example
 * ```tsx
 * <DashboardLayout title="Inventory">
 *   <InventoryPage />
 * </DashboardLayout>
 * ```
 */
export function DashboardLayout({
  title,
  children,
}: DashboardLayoutProps) {
  const { t } = useTranslation()
  const matches = useMatches()
  const isDesktop = useMediaQuery('(min-width: 768px)')
  const sidebarCollapsed = useStore(businessStore, (s) => s.sidebarCollapsed)
  const sidebarOpen = useStore(businessStore, (s) => s.sidebarOpen)

  const derivedTitleKey =
    matches.length > 0
      ? (matches[matches.length - 1] as unknown as { staticData?: unknown })
          .staticData &&
        typeof (matches[matches.length - 1] as unknown as { staticData?: any })
          .staticData?.titleKey === 'string'
        ? (matches[matches.length - 1] as unknown as { staticData?: any })
            .staticData.titleKey
        : undefined
      : undefined

  const resolvedTitle = title ?? (derivedTitleKey ? t(derivedTitleKey) : undefined)

  // Close sidebar when switching to desktop
  useEffect(() => {
    if (isDesktop && sidebarOpen) {
      closeSidebar()
    }
  }, [isDesktop, sidebarOpen])

  // Prevent body scroll when mobile sidebar is open
  useEffect(() => {
    if (!isDesktop && sidebarOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }

    return () => {
      document.body.style.overflow = ''
    }
  }, [isDesktop, sidebarOpen])

  return (
    <div className="min-h-screen bg-base-100">
      {/* Sidebar - Desktop: Fixed, Mobile: Drawer Overlay */}
      <Sidebar />

      {/* Mobile Sidebar Backdrop */}
      {!isDesktop && sidebarOpen && (
        <div
          className="fixed inset-0 bg-base-content/40 backdrop-blur-sm z-40 md:hidden transition-opacity duration-300"
          onClick={closeSidebar}
          onKeyDown={(e) => {
            if (e.key === 'Escape') closeSidebar()
          }}
          role="button"
          tabIndex={0}
          aria-label={t('dashboard.close_menu')}
        />
      )}

      {/* Header - Fixed at top, adjusts based on sidebar state */}
      <Header title={resolvedTitle} />

      {/* Main Content Area */}
      <main
        className={cn(
          // Base: Content below header
          'min-h-screen pt-16 transition-all duration-300',
          // Desktop: Adjust for sidebar width
          isDesktop && !sidebarCollapsed && 'md:ms-64',
          isDesktop && sidebarCollapsed && 'md:ms-20',
          // Mobile: Add bottom nav padding
          !isDesktop && 'pb-20',
        )}
      >
        {/* Content Container with max-width and padding */}
        <div className="container mx-auto px-4 py-6 max-w-7xl">{children}</div>
      </main>

      {/* Mobile Bottom Navigation - Only visible on mobile */}
      {!isDesktop && <BottomNav />}
    </div>
  )
}
