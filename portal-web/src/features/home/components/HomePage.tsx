import { useLoaderData, useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import {
  BookOpen,
  Building2,
  ChevronRight,
  CreditCard,
  HelpCircle,
  Plus,
  Settings,
  Users,
} from 'lucide-react'

import type { Business } from '@/stores/businessStore'
import { Avatar } from '@/components/atoms/Avatar'
import { Logo } from '@/components/atoms/Logo'
import { LanguageSwitcher } from '@/features/language/components/LanguageSwitcher'
import { useAuth } from '@/hooks/useAuth'
import { useLanguage } from '@/hooks/useLanguage'
import { getThumbnailUrl } from '@/lib/assetUrl'
import { selectBusiness } from '@/stores/businessStore'

export function HomePending() {
  return (
    <div className="min-h-screen bg-base-200">
      <header className="sticky top-0 z-30 border-b border-base-300 bg-base-100">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Logo size="md" showText />
            <div className="flex items-center gap-2">
              <LanguageSwitcher variant="compact" />
            </div>
          </div>
        </div>
      </header>

      <main className="container mx-auto max-w-6xl px-4 py-8">
        <div className="mb-8">
          <div className="skeleton h-9 w-64 mb-2" />
          <div className="skeleton h-5 w-96" />
        </div>

        <section className="mb-8">
          <div className="mb-4 flex items-center justify-between">
            <div className="skeleton h-7 w-48" />
            <div className="skeleton h-9 w-36" />
          </div>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="skeleton h-32 rounded-box" />
            ))}
          </div>
        </section>
      </main>
    </div>
  )
}

export function HomePage() {
  const { t: tHome } = useTranslation('home')
  const { t: tAuth } = useTranslation('auth')
  const { t: tCommon } = useTranslation('common')
  const navigate = useNavigate()
  const { user, logout } = useAuth()
  const { isRTL } = useLanguage()
  const { businesses } = useLoaderData({ from: '/' })

  const handleLogout = () => {
    void logout()
      .catch(() => {
        // Silent fail - auth state will be cleared locally anyway
      })
      .finally(() => {
        void navigate({
          to: '/auth/login',
          search: { redirect: '/' },
          replace: true,
        })
      })
  }

  const handleBusinessSelect = (business: Business) => {
    selectBusiness(business.descriptor)
    void navigate({
      to: '/business/$businessDescriptor',
      params: { businessDescriptor: business.descriptor },
    })
  }

  const quickLinks = [
    {
      key: 'workspace',
      icon: Users,
      label: tHome('workspace_settings'),
      description: tHome('workspace_settings_desc'),
      path: '/workspace',
      disabled: true,
    },
    {
      key: 'account',
      icon: Settings,
      label: tHome('account_settings'),
      description: tHome('account_settings_desc'),
      path: '/account',
      disabled: true,
    },
    {
      key: 'billing',
      icon: CreditCard,
      label: tHome('billing'),
      description: tHome('billing_desc'),
      path: '/billing',
      disabled: true,
    },
  ]

  const supportLinks = [
    {
      key: 'help',
      icon: HelpCircle,
      label: tHome('help_center'),
      path: 'https://help.kyora.app',
      external: true,
    },
    {
      key: 'docs',
      icon: BookOpen,
      label: tHome('documentation'),
      path: 'https://docs.kyora.app',
      external: true,
    },
  ]

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <header className="sticky top-0 z-30 border-b border-base-300 bg-base-100">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Logo size="md" showText />
            <div className="flex items-center gap-2">
              <LanguageSwitcher variant="compact" />
              <div className="divider divider-horizontal mx-0" />
              <span className="hidden text-sm text-base-content/70 sm:inline">
                {user?.firstName} {user?.lastName}
              </span>
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={handleLogout}
              >
                {tAuth('logout')}
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto max-w-6xl px-4 py-8">
        {/* Welcome Section */}
        <div className="mb-8">
          <h1 className="mb-2 text-3xl font-bold">
            {tHome('welcome')}, {user?.firstName}!
          </h1>
          <p className="text-base-content/70">
            {tHome('select_business_or_manage')}
          </p>
        </div>

        {/* Businesses Section */}
        <section className="mb-8">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-xl font-semibold">
              {tHome('your_businesses')}
            </h2>
            <button className="btn btn-primary btn-sm gap-2" disabled>
              <Plus size={18} />
              {tHome('add_business')}
            </button>
          </div>

          {businesses.length === 0 ? (
            <div className="card border border-base-300 bg-base-100">
              <div className="card-body items-center py-12 text-center">
                <Building2 size={48} className="mb-4 text-base-content/30" />
                <h3 className="mb-2 text-lg font-semibold">
                  {tHome('no_businesses')}
                </h3>
                <p className="mb-4 text-sm text-base-content/60">
                  {tHome('no_businesses_desc')}
                </p>
                <button className="btn btn-primary gap-2" disabled>
                  <Plus size={18} />
                  {tHome('create_first_business')}
                </button>
              </div>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {businesses.map((business: Business) => (
                <button
                  key={business.id}
                  onClick={() => handleBusinessSelect(business)}
                  className="card border border-base-300 bg-base-100 text-start transition-all hover:border-primary  hover:cursor-pointer"
                >
                  <div className="card-body p-4">
                    <div className="flex items-start gap-3">
                      <Avatar
                        src={getThumbnailUrl(business.logo)}
                        fallback={business.name}
                        size="md"
                        shape="square"
                        className="shrink-0"
                      />
                      <div className="min-w-0 flex-1">
                        <h3 className="truncate font-semibold">
                          {business.name}
                        </h3>
                        <p className="truncate text-sm text-base-content/60">
                          @{business.descriptor}
                        </p>
                      </div>
                      <ChevronRight
                        size={20}
                        className={`text-base-content/40 ${isRTL ? 'rotate-180' : ''}`}
                      />
                    </div>
                  </div>
                </button>
              ))}
            </div>
          )}
        </section>

        {/* Quick Links Section */}
        <section className="mb-8">
          <h2 className="mb-4 text-xl font-semibold">{tHome('quick_links')}</h2>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {quickLinks.map((link) => (
              <button
                key={link.key}
                disabled={link.disabled}
                className="card border border-base-300 bg-base-100 text-start transition-all hover:border-primary  hover:cursor-pointer disabled:cursor-not-allowed disabled:opacity-50"
              >
                <div className="card-body p-4">
                  <div className="mb-2 flex items-center gap-3">
                    <link.icon size={24} className="text-primary" />
                    <h3 className="font-semibold">{link.label}</h3>
                  </div>
                  <p className="text-sm text-base-content/60">
                    {link.description}
                  </p>
                  {link.disabled && (
                    <span className="badge badge-ghost badge-sm mt-2">
                      {tCommon('coming_soon')}
                    </span>
                  )}
                </div>
              </button>
            ))}
          </div>
        </section>

        {/* Support Links Section */}
        <section>
          <h2 className="mb-4 text-xl font-semibold">{tHome('support')}</h2>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {supportLinks.map((link) => (
              <a
                key={link.key}
                href={link.path}
                target={link.external ? '_blank' : undefined}
                rel={link.external ? 'noopener noreferrer' : undefined}
                className="card border border-base-300 bg-base-100 transition-all hover:border-primary "
              >
                <div className="card-body p-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <link.icon size={20} className="text-primary" />
                      <span className="font-medium">{link.label}</span>
                    </div>
                    <ChevronRight
                      size={20}
                      className={`text-base-content/40 ${isRTL ? 'rotate-180' : ''}`}
                    />
                  </div>
                </div>
              </a>
            ))}
          </div>
        </section>
      </main>
    </div>
  )
}
