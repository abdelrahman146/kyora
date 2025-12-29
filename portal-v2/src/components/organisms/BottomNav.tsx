import { BarChart3, Home, Package, Users } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { ErrorBoundary } from '@/components/atoms/ErrorBoundary'
import { cn } from '@/lib/utils'

export interface BottomNavProps {
  businessDescriptor: string
}

/**
 * BottomNav Component
 *
 * Mobile bottom navigation bar.
 * Wrapped in ErrorBoundary to prevent layout crashes.
 */
function BottomNavContent({ businessDescriptor }: BottomNavProps) {
  const navItems = [
    {
      label: 'الرئيسية',
      href: `/business/${businessDescriptor}`,
      icon: Home,
    },
    {
      label: 'العملاء',
      href: `/business/${businessDescriptor}/customers`,
      icon: Users,
    },
    {
      label: 'المخزون',
      href: `/business/${businessDescriptor}/inventory`,
      icon: Package,
      disabled: true,
    },
    {
      label: 'التقارير',
      href: `/business/${businessDescriptor}/analytics`,
      icon: BarChart3,
      disabled: true,
    },
  ]

  return (
    <nav className="btm-nav btm-nav-lg border-t border-base-300 bg-base-100 lg:hidden">
      {navItems.map((item) => {
        const Icon = item.icon
        return (
          <Link
            key={item.href}
            to={item.href}
            disabled={item.disabled}
            className={cn(item.disabled && 'cursor-not-allowed opacity-50')}
            activeProps={{ className: 'active text-primary' }}
          >
            <Icon size={20} />
            <span className="btm-nav-label text-xs">{item.label}</span>
          </Link>
        )
      })}
    </nav>
  )
}

export function BottomNav(props: BottomNavProps) {
  return (
    <ErrorBoundary compact>
      <BottomNavContent {...props} />
    </ErrorBoundary>
  )
}
