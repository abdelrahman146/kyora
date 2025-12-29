import { createFileRoute } from '@tanstack/react-router'

/**
 * Business Dashboard Route
 *
 * Home page for a specific business showing:
 * - Overview cards (revenue, orders, customers, inventory)
 * - Recent activity feed
 * - Quick action buttons
 *
 * TODO: Implement dashboard analytics query and display
 */
export const Route = createFileRoute('/business/$businessDescriptor/')({
  component: BusinessDashboard,
})

/**
 * Business Dashboard Component
 *
 * Displays business overview and quick actions.
 */
function BusinessDashboard() {
  const { businessDescriptor } = Route.useParams()
  const { business } = Route.useRouteContext()

  return (
    <div className="space-y-6">
      {/* Overview Cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {/* Revenue Card */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="text-sm font-medium text-base-content/70">
              الإيرادات
            </h3>
            <div className="flex items-baseline gap-2">
              <span className="text-3xl font-bold">0</span>
              <span className="text-sm text-base-content/70">
                {business.currency}
              </span>
            </div>
            <p className="text-xs text-success">+0% من الشهر الماضي</p>
          </div>
        </div>

        {/* Orders Card */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="text-sm font-medium text-base-content/70">
              الطلبات
            </h3>
            <div className="flex items-baseline gap-2">
              <span className="text-3xl font-bold">0</span>
            </div>
            <p className="text-xs text-base-content/50">لا توجد طلبات بعد</p>
          </div>
        </div>

        {/* Customers Card */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="text-sm font-medium text-base-content/70">
              العملاء
            </h3>
            <div className="flex items-baseline gap-2">
              <span className="text-3xl font-bold">0</span>
            </div>
            <p className="text-xs text-base-content/50">لا يوجد عملاء بعد</p>
          </div>
        </div>

        {/* Inventory Card */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="text-sm font-medium text-base-content/70">
              المخزون
            </h3>
            <div className="flex items-baseline gap-2">
              <span className="text-3xl font-bold">0</span>
            </div>
            <p className="text-xs text-base-content/50">لا توجد منتجات بعد</p>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <h2 className="card-title">إجراءات سريعة</h2>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <a
              href={`/business/${businessDescriptor}/customers`}
              className="btn btn-outline"
            >
              إضافة عميل
            </a>
            <button className="btn btn-outline" disabled>
              إنشاء طلب
            </button>
            <button className="btn btn-outline" disabled>
              إضافة منتج
            </button>
            <button className="btn btn-outline" disabled>
              عرض التقارير
            </button>
          </div>
        </div>
      </div>

      {/* Recent Activity - Placeholder */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <h2 className="card-title">النشاط الأخير</h2>
          <div className="flex min-h-[200px] items-center justify-center">
            <p className="text-base-content/50">لا يوجد نشاط حتى الآن</p>
          </div>
        </div>
      </div>
    </div>
  )
}
