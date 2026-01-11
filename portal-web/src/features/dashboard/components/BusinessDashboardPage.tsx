import { Link, useParams, useRouteContext } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { useAuth } from '@/hooks/useAuth'

export function BusinessDashboardPage() {
  const { t: tDashboard } = useTranslation('dashboard')
  const { t: tCustomers } = useTranslation('customers')
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/',
  })
  const { business } = useRouteContext({
    from: '/business/$businessDescriptor/',
  })
  const { user } = useAuth()

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div className="card bg-base-200">
        <div className="card-body">
          <h2 className="card-title text-2xl">
            {tDashboard('welcome')}, {user?.firstName}!
          </h2>
          <p className="text-base-content/70">
            {tDashboard('managing')}: {business.name}
          </p>
        </div>
      </div>

      {/* Overview Cards */}
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {/* Revenue Card */}
        <div className="card border border-base-300 bg-base-100">
          <div className="card-body">
            <h3 className="card-title text-lg">{tDashboard('revenue')}</h3>
            <p className="text-3xl font-bold text-primary">
              {business.currency} 0
            </p>
            <p className="text-sm text-base-content/60">
              {tDashboard('this_month')}
            </p>
          </div>
        </div>

        {/* Orders Card */}
        <div className="card border border-base-300 bg-base-100">
          <div className="card-body">
            <h3 className="card-title text-lg">{tDashboard('orders')}</h3>
            <p className="text-3xl font-bold text-success">0</p>
            <p className="text-sm text-base-content/60">
              {tDashboard('pending')}
            </p>
          </div>
        </div>

        {/* Inventory Card */}
        <div className="card border border-base-300 bg-base-100">
          <div className="card-body">
            <h3 className="card-title text-lg">{tDashboard('inventory')}</h3>
            <p className="text-3xl font-bold text-warning">0</p>
            <p className="text-sm text-base-content/60">
              {tDashboard('low_stock')}
            </p>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <h2 className="card-title">{tDashboard('quick_actions')}</h2>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Link
              to="/business/$businessDescriptor/customers"
              params={{ businessDescriptor }}
              search={{ page: 1, pageSize: 20, sortOrder: 'desc' }}
              className="btn btn-outline"
            >
              {tCustomers('add_customer')}
            </Link>
            <button className="btn btn-outline" disabled>
              {tDashboard('orders')}
            </button>
            <button className="btn btn-outline" disabled>
              {tDashboard('inventory')}
            </button>
            <button className="btn btn-outline" disabled>
              {tDashboard('analytics')}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
