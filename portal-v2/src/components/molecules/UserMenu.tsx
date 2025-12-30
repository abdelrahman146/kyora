import { useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { ChevronDown, LogOut } from 'lucide-react'
import { Avatar } from '../atoms/Avatar'
import { Dropdown } from '../atoms/Dropdown'
import { useAuth } from '../../hooks/useAuth'

/**
 * UserMenu Component
 *
 * Dropdown menu in the header showing user info and actions.
 *
 * Features:
 * - User avatar and name
 * - Profile, Settings, Help links
 * - Logout action
 * - RTL support
 *
 * @example
 * ```tsx
 * <UserMenu />
 * ```
 */
export function UserMenu() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { user, logout } = useAuth()

  const handleLogout = () => {
    void logout()
      .then(() => {
        void navigate({
          to: '/auth/login',
          search: { redirect: '/' },
          replace: true,
        })
      })
      .catch(() => {
        // Silent fail - logout will clear local state anyway
        void navigate({
          to: '/auth/login',
          search: { redirect: '/' },
          replace: true,
        })
      })
  }

  if (!user) return null

  const userFullName = `${user.firstName} ${user.lastName}`
  const userInitials = `${user.firstName[0]}${user.lastName[0]}`

  return (
    <Dropdown
      align="end"
      width="14rem"
      trigger={
        <button
          type="button"
          className="btn btn-ghost gap-2 h-auto min-h-0 py-2 px-2 hover:bg-base-200"
        >
          <Avatar
            src={undefined}
            fallback={userInitials}
            size="sm"
            shape="circle"
          />
          <ChevronDown
            size={16}
            className="transition-transform text-base-content/60 hidden md:block"
          />
        </button>
      }
    >
      <div className="py-2">
        {/* User Info */}
        <div className="px-4 py-3 border-b border-base-300">
          <div className="flex flex-col gap-1">
            <span className="font-semibold text-sm text-base-content">
              {userFullName}
            </span>
            <span className="text-xs text-base-content/60">{user.email}</span>
          </div>
        </div>

        {/* Menu Items */}
        <div className="py-1">
          {/* Profile, Settings, Help routes commented out until implemented */}
          {/* <button
            type="button"
            onClick={() => {
              void navigate({ to: '/profile' })
            }}
            className="flex items-center gap-3 w-full px-4 py-2.5 hover:bg-base-200 transition-colors text-start"
          >
            <User size={16} className="shrink-0" />
            <span className="text-sm">{t('dashboard.profile')}</span>
          </button>

          <button
            type="button"
            onClick={() => {
              void navigate({ to: '/settings' })
            }}
            className="flex items-center gap-3 w-full px-4 py-2.5 hover:bg-base-200 transition-colors text-start"
          >
            <Settings size={16} className="shrink-0" />
            <span className="text-sm">{t('dashboard.settings')}</span>
          </button>

          <button
            type="button"
            onClick={() => {
              void navigate({ to: '/help' })
            }}
            className="flex items-center gap-3 w-full px-4 py-2.5 hover:bg-base-200 transition-colors text-start"
          >
            <HelpCircle size={16} className="shrink-0" />
            <span className="text-sm">{t('dashboard.help')}</span>
          </button> */}
        </div>

        <div className="divider my-1"></div>

        {/* Logout */}
        <div className="py-1">
          <button
            type="button"
            onClick={handleLogout}
            className="flex items-center gap-3 w-full px-4 py-2.5 text-error hover:bg-error/10 transition-colors text-start"
          >
            <LogOut size={16} className="shrink-0" />
            <span className="text-sm">{t('auth.logout')}</span>
          </button>
        </div>
      </div>
    </Dropdown>
  )
}
