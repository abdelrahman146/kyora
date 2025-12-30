import { useTranslation } from 'react-i18next'
import { Link, useLocation } from '@tanstack/react-router'
import {
  LayoutDashboard,
  Menu,
  Package,
  ShoppingCart,
  Users,
} from 'lucide-react'
import { openSidebar } from '@/stores/businessStore'
import { cn } from '@/lib/utils'

interface NavItem {
  key: string
  icon: typeof LayoutDashboard
  path: string
}

const navItems: Array<NavItem> = [
  { key: 'dashboard', icon: LayoutDashboard, path: '' },
  { key: 'inventory', icon: Package, path: '/inventory' },
  { key: 'orders', icon: ShoppingCart, path: '/orders' },
  { key: 'customers', icon: Users, path: '/customers' },
]

export interface BottomNavProps {
  businessDescriptor: string
}

/**
 * BottomNav Component
 *
 * Fixed bottom navigation bar for mobile devices.
 *
 * Features:
 * - Shows 4 core navigation items + More button (opens drawer)
 * - Active route highlighting
 * - Touch-friendly targets (min 44px height)
 * - Safe area padding for devices with notches
 * - RTL support with logical properties
 * - Smooth press feedback
 *
 * Design Pattern:
 * - Primary actions (Dashboard, Inventory, Orders, Customers)
 * - "More" button opens the full sidebar drawer
 */
export function BottomNav({ businessDescriptor }: BottomNavProps) {
  const { t } = useTranslation()
  const location = useLocation()

  return (
    <nav className="fixed bottom-0 start-0 end-0 h-16 bg-base-100 border-t border-base-300 z-40 md:hidden safe-area-pb">
      <div className="h-full flex items-center justify-around px-2">
        {/* Core Navigation Items */}
        {navItems.map((item) => {
          // Build full path using props business descriptor
          const itemPath = `/business/${businessDescriptor}${item.path}`
          const isActive =
            item.path === ''
              ? location.pathname === `/business/${businessDescriptor}` ||
                location.pathname === `/business/${businessDescriptor}/`
              : location.pathname.startsWith(itemPath)
          const Icon = item.icon

          return (
            <Link
              key={item.key}
              to={itemPath}
              className={cn(
                'flex flex-col items-center justify-center gap-1 flex-1 h-full min-w-0',
                'transition-all active:scale-95 rounded-lg',
                // Min touch target
                'min-h-11',
                // Active state
                isActive && 'text-primary',
                !isActive && 'text-base-content/70',
              )}
            >
              <Icon size={22} className="shrink-0" />
              <span className="text-xs font-medium truncate">
                {t(`dashboard.${item.key}`)}
              </span>
            </Link>
          )
        })}

        {/* More Button - Opens Sidebar Drawer */}
        <button
          type="button"
          onClick={openSidebar}
          className={cn(
            'flex flex-col items-center justify-center gap-1 flex-1 h-full min-w-0',
            'transition-all active:scale-95 rounded-lg',
            'min-h-11',
            'text-base-content/70',
          )}
        >
          <Menu size={22} className="shrink-0" />
          <span className="text-xs font-medium truncate">
            {t('dashboard.more')}
          </span>
        </button>
      </div>
    </nav>
  )
}
