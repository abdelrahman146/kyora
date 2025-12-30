/**
 * Home Page Route (Authenticated)
 *
 * Landing page for authenticated users to select a business or access workspace features.
 *
 * Features:
 * - Business selection cards
 * - Quick links to billing, workspace, account
 * - Support and documentation links
 * - Mobile-first design
 * - Fully localized
 */

import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
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
import { businessApi } from '@/api/business'
import { Logo } from '@/components/atoms/Logo'
import { LanguageSwitcher } from '@/components/molecules/LanguageSwitcher'
import { useAuth } from '@/hooks/useAuth'
import { useLanguage } from '@/hooks/useLanguage'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'
import { requireAuth } from '@/lib/routeGuards'
import {
  businessStore,
  selectBusiness,
  setBusinesses,
} from '@/stores/businessStore'

/**
 * Home Route Configuration
 *
 * Redirects authenticated users to their selected business or shows business selection hub.
 * - If user has a previously selected business, redirect to `/business/:descriptor`
 * - Otherwise, show business selection interface
 */
export const Route = createFileRoute('/')({
  beforeLoad: () => {
    requireAuth()
  },

  loader: async ({ context }) => {
     
    const queryClient = (context as any).queryClient
    if (!queryClient) {
      throw new Error('QueryClient not found in router context')
    }

    const data = await queryClient.ensureQueryData({
      queryKey: queryKeys.businesses.list(),
      queryFn: () => businessApi.listBusinesses(),
      staleTime: STALE_TIME.FIVE_MINUTES,
    })

    setBusinesses(data.businesses)

    const state = businessStore.state
    if (
      state.selectedBusinessDescriptor &&
      data.businesses.some(
        (b: Business) => b.descriptor === state.selectedBusinessDescriptor,
      )
    ) {
      throw redirect({
        to: '/business/$businessDescriptor',
        params: { businessDescriptor: state.selectedBusinessDescriptor },
      })
    }

    if (data.businesses.length === 1) {
      throw redirect({
        to: '/business/$businessDescriptor',
        params: { businessDescriptor: data.businesses[0].descriptor },
      })
    }

    return { businesses: data.businesses }
  },

  component: HomePage,
})

/**
 * Home Page Component
 *
 * Shows business selection hub when user has multiple businesses and no selection.
 */
function HomePage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { user } = useAuth()
  const { isRTL } = useLanguage()
  const { businesses } = Route.useLoaderData()
  const state = useStore(businessStore)

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
      label: t('home.workspace_settings'),
      description: t('home.workspace_settings_desc'),
      path: '/workspace',
      disabled: true,
    },
    {
      key: 'account',
      icon: Settings,
      label: t('home.account_settings'),
      description: t('home.account_settings_desc'),
      path: '/account',
      disabled: true,
    },
    {
      key: 'billing',
      icon: CreditCard,
      label: t('home.billing'),
      description: t('home.billing_desc'),
      path: '/billing',
      disabled: true,
    },
  ]

  const supportLinks = [
    {
      key: 'help',
      icon: HelpCircle,
      label: t('home.help_center'),
      path: 'https://help.kyora.app',
      external: true,
    },
    {
      key: 'docs',
      icon: BookOpen,
      label: t('home.documentation'),
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
              <LanguageSwitcher variant="iconOnly" />
              <div className="divider divider-horizontal mx-0" />
              <span className="hidden text-sm text-base-content/70 sm:inline">
                {user?.firstName} {user?.lastName}
              </span>
              <a href="/auth/login" className="btn btn-ghost btn-sm">
                {t('auth.logout')}
              </a>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto max-w-6xl px-4 py-8">
        {/* Welcome Section */}
        <div className="mb-8">
          <h1 className="mb-2 text-3xl font-bold">
            {t('home.welcome')}, {user?.firstName}!
          </h1>
          <p className="text-base-content/70">
            {t('home.select_business_or_manage')}
          </p>
        </div>

        {/* Businesses Section */}
        <section className="mb-8">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-xl font-semibold">{t('home.your_businesses')}</h2>
            <button className="btn btn-primary btn-sm gap-2" disabled>
              <Plus size={18} />
              {t('home.add_business')}
            </button>
          </div>

          {businesses.length === 0 ? (
            <div className="card border border-base-300 bg-base-100">
              <div className="card-body items-center py-12 text-center">
                <Building2 size={48} className="mb-4 text-base-content/30" />
                <h3 className="mb-2 text-lg font-semibold">
                  {t('home.no_businesses')}
                </h3>
                <p className="mb-4 text-sm text-base-content/60">
                  {t('home.no_businesses_desc')}
                </p>
                <button className="btn btn-primary gap-2" disabled>
                  <Plus size={18} />
                  {t('home.create_first_business')}
                </button>
              </div>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {businesses.map((business: Business) => (
                <button
                  key={business.id}
                  onClick={() => handleBusinessSelect(business)}
                  className="card border border-base-300 bg-base-100 text-start transition-all hover:border-primary hover:shadow-md"
                >
                  <div className="card-body p-4">
                    <div className="flex items-start gap-3">
                      {business.logoUrl ? (
                        <img
                          src={business.logoUrl}
                          alt={business.name}
                          className="h-12 w-12 rounded-lg object-cover"
                        />
                      ) : (
                        <div className="avatar placeholder">
                          <div className="h-12 w-12 rounded-lg bg-primary text-primary-content">
                            <span className="text-lg font-bold">
                              {business.name.charAt(0).toUpperCase()}
                            </span>
                          </div>
                        </div>
                      )}
                      <div className="min-w-0 flex-1">
                        <h3 className="truncate font-semibold">{business.name}</h3>
                        <p className="truncate text-sm text-base-content/60">
                          @{business.descriptor}
                        </p>
                      </div>
                      <ChevronRight
                        size={20}
                        className={`text-base-content/40 ${isRTL ? 'rotate-180' : ''}`}
                      />
                    </div>
                    {state.selectedBusinessDescriptor === business.descriptor && (
                      <span className="badge badge-primary badge-sm mt-2">
                        {t('common.selected', 'Selected')}
                      </span>
                    )}
                  </div>
                </button>
              ))}
            </div>
          )}
        </section>

        {/* Quick Links Section */}
        <section className="mb-8">
          <h2 className="mb-4 text-xl font-semibold">{t('home.quick_links')}</h2>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {quickLinks.map((link) => (
              <button
                key={link.key}
                disabled={link.disabled}
                className="card border border-base-300 bg-base-100 text-start transition-all hover:border-primary hover:shadow-md disabled:cursor-not-allowed disabled:opacity-50"
              >
                <div className="card-body p-4">
                  <div className="mb-2 flex items-center gap-3">
                    <link.icon size={24} className="text-primary" />
                    <h3 className="font-semibold">{link.label}</h3>
                  </div>
                  <p className="text-sm text-base-content/60">{link.description}</p>
                  {link.disabled && (
                    <span className="badge badge-ghost badge-sm mt-2">
                      {t('common.coming_soon')}
                    </span>
                  )}
                </div>
              </button>
            ))}
          </div>
        </section>

        {/* Support Links Section */}
        <section>
          <h2 className="mb-4 text-xl font-semibold">{t('home.support')}</h2>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {supportLinks.map((link) => (
              <a
                key={link.key}
                href={link.path}
                target={link.external ? '_blank' : undefined}
                rel={link.external ? 'noopener noreferrer' : undefined}
                className="card border border-base-300 bg-base-100 transition-all hover:border-primary hover:shadow-md"
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
