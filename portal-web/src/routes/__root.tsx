import { useEffect } from 'react'
import { Outlet, createRootRoute } from '@tanstack/react-router'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { Toaster } from 'react-hot-toast'
import type { RouterContext } from '@/router'
import { useLanguage } from '@/hooks/useLanguage'
import { initializeAuth } from '@/stores/authStore'

export const Route = createRootRoute<RouterContext>({
  component: RootComponent,
})

function RootComponent() {
  const { isRTL } = useLanguage()

  // Restore session on mount
  useEffect(() => {
    void initializeAuth()
  }, [])

  // Toast position based on RTL and screen size
  const getToastPosition = () => {
    const isDesktop = window.matchMedia('(min-width: 768px)').matches
    if (!isDesktop) return 'top-center'
    return isRTL ? 'top-right' : 'top-left'
  }

  return (
    <>
      {/* Global Toast Notifications */}
      <Toaster
        position={getToastPosition()}
        toastOptions={{
          duration: 4000,
          style: {
            fontFamily:
              'IBM Plex Sans Arabic, -apple-system, BlinkMacSystemFont, sans-serif',
            fontSize: '14px',
            lineHeight: '1.5',
            borderRadius: '12px',
            padding: '12px 16px',
            boxShadow:
              '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
            maxWidth: '400px',
          },
          success: {
            iconTheme: {
              primary: '#10B981',
              secondary: '#FFFFFF',
            },
            style: {
              background: '#FFFFFF',
              color: '#0F172A',
              border: '1px solid #10B981',
            },
          },
          error: {
            iconTheme: {
              primary: '#EF4444',
              secondary: '#FFFFFF',
            },
            style: {
              background: '#FFFFFF',
              color: '#0F172A',
              border: '1px solid #EF4444',
            },
          },
          loading: {
            iconTheme: {
              primary: '#0D9488',
              secondary: '#FFFFFF',
            },
            style: {
              background: '#FFFFFF',
              color: '#0F172A',
              border: '1px solid #E2E8F0',
            },
          },
        }}
      />

      {/* Main App Content */}
      <Outlet />

      {/* Development Tools */}
      {import.meta.env.DEV && (
        <>
          <TanStackRouterDevtools />
          <ReactQueryDevtools />
        </>
      )}
    </>
  )
}
