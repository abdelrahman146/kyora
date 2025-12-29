import { Menu, X } from 'lucide-react'
import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import { ErrorBoundary } from '@/components/atoms/ErrorBoundary'
import { cn } from '@/lib/utils'

export interface SidebarProps {
  businessDescriptor: string
  businessName: string
}

/**
 * Sidebar Component
 *
 * Main navigation sidebar with collapsible mobile drawer.
 * Wrapped in ErrorBoundary to prevent layout crashes.
 */
function SidebarContent({ businessDescriptor, businessName }: SidebarProps) {
  const [isOpen, setIsOpen] = useState(false)

  const menuItems = [
    { label: 'لوحة التحكم', href: `/business/${businessDescriptor}`, icon: '📊' },
    { label: 'العملاء', href: `/business/${businessDescriptor}/customers`, icon: '👥' },
    { label: 'الطلبات', href: `/business/${businessDescriptor}/orders`, icon: '📦', disabled: true },
    { label: 'المخزون', href: `/business/${businessDescriptor}/inventory`, icon: '📦', disabled: true },
    { label: 'التقارير', href: `/business/${businessDescriptor}/analytics`, icon: '📈', disabled: true },
  ]

  return (
    <>
      {/* Mobile Toggle */}
      <div className="lg:hidden">
        <button
          onClick={() => setIsOpen(!isOpen)}
          className="btn btn-circle btn-ghost"
          aria-label="Toggle menu"
        >
          {isOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </div>

      {/* Sidebar */}
      <aside
        className={cn(
          'fixed inset-y-0 start-0 z-40 w-64 transform bg-base-200 transition-transform duration-300 lg:static lg:translate-x-0',
          isOpen ? 'translate-x-0' : '-translate-x-full',
        )}
      >
        <div className="flex h-full flex-col">
          {/* Logo/Business Name */}
          <div className="border-b border-base-300 p-4">
            <h2 className="text-lg font-bold">{businessName}</h2>
          </div>

          {/* Navigation */}
          <nav className="flex-1 overflow-y-auto p-4">
            <ul className="menu gap-2">
              {menuItems.map((item) => (
                <li key={item.href}>
                  {item.disabled ? (
                    <span className="cursor-not-allowed opacity-50">
                      <span className="text-xl">{item.icon}</span>
                      {item.label}
                    </span>
                  ) : (
                    <Link
                      to={item.href}
                      className="flex items-center gap-3"
                      activeProps={{ className: 'active' }}
                    >
                      <span className="text-xl">{item.icon}</span>
                      {item.label}
                    </Link>
                  )}
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </aside>

      {/* Backdrop */}
      {isOpen && (
        <div
          className="fixed inset-0 z-30 bg-black/50 lg:hidden"
          onClick={() => setIsOpen(false)}
        />
      )}
    </>
  )
}

export function Sidebar(props: SidebarProps) {
  return (
    <ErrorBoundary compact>
      <SidebarContent {...props} />
    </ErrorBoundary>
  )
}
