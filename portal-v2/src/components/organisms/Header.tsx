import { Bell, Settings } from 'lucide-react'
import { ErrorBoundary } from '@/components/atoms/ErrorBoundary'
import { LanguageSwitcher } from '@/components/molecules/LanguageSwitcher'

/**
 * Header Component
 *
 * Top navigation bar with notifications and settings.
 * Wrapped in ErrorBoundary to prevent layout crashes.
 */
function HeaderContent() {
  return (
    <header className="border-b border-base-300 bg-base-100">
      <div className="navbar">
        <div className="flex-1">
          {/* Empty - can add breadcrumbs later */}
        </div>

        <div className="flex-none gap-2">
          {/* Language Switcher */}
          <LanguageSwitcher variant="icon" />

          {/* Notifications */}
          <button className="btn btn-circle btn-ghost" aria-label="Notifications">
            <Bell size={20} />
          </button>

          {/* Settings */}
          <button className="btn btn-circle btn-ghost" aria-label="Settings">
            <Settings size={20} />
          </button>

          {/* User Menu */}
          <div className="dropdown dropdown-end">
            <div
              tabIndex={0}
              role="button"
              className="avatar btn btn-circle btn-ghost"
            >
              <div className="w-10 rounded-full bg-neutral text-neutral-content">
                <span className="text-lg">👤</span>
              </div>
            </div>
            <ul
              tabIndex={0}
              className="menu dropdown-content menu-sm z-[1] mt-3 w-52 rounded-box bg-base-100 p-2 shadow"
            >
              <li>
                <a>الملف الشخصي</a>
              </li>
              <li>
                <a>الإعدادات</a>
              </li>
              <li>
                <a>تسجيل الخروج</a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </header>
  )
}

export function Header() {
  return (
    <ErrorBoundary compact>
      <HeaderContent />
    </ErrorBoundary>
  )
}
