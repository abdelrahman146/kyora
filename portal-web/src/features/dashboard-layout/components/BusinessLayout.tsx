import { Suspense } from 'react'
import { Outlet } from '@tanstack/react-router'

import { DashboardLayout } from '@/features/dashboard-layout/components/DashboardLayout'

export function BusinessLayout() {
  return (
    <DashboardLayout>
      <Suspense
        fallback={
          <div className="flex min-h-[400px] items-center justify-center">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        }
      >
        <Outlet />
      </Suspense>
    </DashboardLayout>
  )
}
