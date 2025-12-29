import type { ReactNode } from 'react'
import { Header } from '@/components/organisms/Header'
import { Sidebar } from '@/components/organisms/Sidebar'
import { BottomNav } from '@/components/organisms/BottomNav'

export interface DashboardLayoutProps {
  businessDescriptor: string
  businessName: string
  children: ReactNode
}

/**
 * DashboardLayout Template
 *
 * Main layout composing Sidebar, Header, and BottomNav.
 * Each organism wrapped in individual ErrorBoundaries to prevent cascading failures.
 */
export function DashboardLayout({
  businessDescriptor,
  businessName,
  children,
}: DashboardLayoutProps) {
  return (
    <div className="flex min-h-screen bg-base-200">
      {/* Sidebar - Desktop */}
      <Sidebar businessDescriptor={businessDescriptor} businessName={businessName} />

      {/* Main Content Area */}
      <div className="flex flex-1 flex-col">
        {/* Header */}
        <Header />

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-4 lg:p-8">
          {children}
        </main>

        {/* Bottom Navigation - Mobile */}
        <BottomNav businessDescriptor={businessDescriptor} />
      </div>
    </div>
  )
}
