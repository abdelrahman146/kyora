import { createFileRoute } from '@tanstack/react-router'
import { BusinessDashboardPage } from '@/features/dashboard/components/BusinessDashboardPage'

export const Route = createFileRoute('/business/$businessDescriptor/')({
  staticData: {
    titleKey: 'pages.dashboard',
  },
  component: BusinessDashboardPage,
})
