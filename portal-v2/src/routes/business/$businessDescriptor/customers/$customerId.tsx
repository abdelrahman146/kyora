import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useState } from 'react'

import {
  useCustomerQuery,
  useDeleteCustomerMutation,
} from '@/api/customer'

/**
 * Customer Detail Route
 *
 * Displays detailed customer information with:
 * - Customer profile card
 * - Order history
 * - Edit/delete actions
 * - Notes and activity
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/customers/$customerId',
)({
  component: () => (
    <Suspense fallback={<CustomerDetailSkeleton />}>
      <CustomerDetailPage />
    </Suspense>
  ),
})

/**
 * Customer Detail Page Component
 */
function CustomerDetailPage() {
  const { businessDescriptor, customerId } = Route.useParams()
  const { business } = Route.useRouteContext()
  const navigate = useNavigate()
  const [showDeleteModal, setShowDeleteModal] = useState(false)

  // Fetch customer data
  const { data: customer, isLoading, error } = useCustomerQuery(
    businessDescriptor,
    customerId,
  )

  // Delete mutation
  const deleteMutation = useDeleteCustomerMutation(businessDescriptor, {
    onSuccess: () => {
      // Navigate back to customers list
      void navigate({
        to: '/business/$businessDescriptor/customers',
        params: { businessDescriptor },
        search: { page: 1, limit: 20 },
      })
    },
  })

  const handleDelete = () => {
    deleteMutation.mutate(customerId)
  }

  if (isLoading) {
    return <CustomerDetailSkeleton />
  }

  if (error || !customer) {
    return (
      <div className="flex min-h-[400px] flex-col items-center justify-center gap-4">
        <p className="text-error">حدث خطأ في تحميل بيانات العميل</p>
        <button
          className="btn btn-sm"
          onClick={() => window.location.reload()}
        >
          إعادة المحاولة
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-4">
          <a
            href={`/business/${businessDescriptor}/customers`}
            className="btn btn-circle btn-ghost btn-sm"
          >
            ←
          </a>
          <div>
            <h1 className="text-2xl font-bold">{customer.fullName}</h1>
            <p className="text-sm text-base-content/70">تفاصيل العميل</p>
          </div>
        </div>
        <div className="flex gap-2">
          <button className="btn btn-outline btn-sm">تعديل</button>
          <button
            className="btn btn-error btn-outline btn-sm"
            onClick={() => setShowDeleteModal(true)}
          >
            حذف
          </button>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Customer Info Card */}
        <div className="card bg-base-100 shadow lg:col-span-2">
          <div className="card-body">
            <h2 className="card-title">معلومات العميل</h2>

            <div className="space-y-4">
              {/* Avatar and Name */}
              <div className="flex items-center gap-4">
                <div className="avatar placeholder">
                  <div className="w-20 rounded-full bg-neutral text-neutral-content">
                    <span className="text-3xl">
                      {customer.fullName.charAt(0)}
                    </span>
                  </div>
                </div>
                <div>
                  <h3 className="text-xl font-bold">{customer.fullName}</h3>
                  <p className="text-sm text-base-content/70">
                    عميل منذ{' '}
                    {new Date(customer.createdAt).toLocaleDateString('ar-EG')}
                  </p>
                </div>
              </div>

              {/* Contact Info */}
              <div className="divider"></div>
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="label label-text text-xs font-semibold">
                    البريد الإلكتروني
                  </label>
                  <p className="text-sm">
                    {customer.email || (
                      <span className="text-base-content/50">غير محدد</span>
                    )}
                  </p>
                </div>
                <div>
                  <label className="label label-text text-xs font-semibold">
                    رقم الهاتف
                  </label>
                  <p className="text-sm">
                    {customer.phoneNumber ? (
                      <span>
                        {customer.phonePrefix} {customer.phoneNumber}
                      </span>
                    ) : (
                      <span className="text-base-content/50">غير محدد</span>
                    )}
                  </p>
                </div>
              </div>

              {/* Address */}
              {(customer.address || customer.city || customer.country) && (
                <>
                  <div className="divider"></div>
                  <div>
                    <label className="label label-text text-xs font-semibold">
                      العنوان
                    </label>
                    <p className="text-sm">
                      {customer.address && <span>{customer.address}</span>}
                      {customer.city && <span>, {customer.city}</span>}
                      {customer.country && <span>, {customer.country}</span>}
                    </p>
                  </div>
                </>
              )}

              {/* Social Media */}
              {(customer.instagramHandle || customer.facebookHandle) && (
                <>
                  <div className="divider"></div>
                  <div className="grid gap-4 sm:grid-cols-2">
                    {customer.instagramHandle && (
                      <div>
                        <label className="label label-text text-xs font-semibold">
                          Instagram
                        </label>
                        <p className="text-sm">@{customer.instagramHandle}</p>
                      </div>
                    )}
                    {customer.facebookHandle && (
                      <div>
                        <label className="label label-text text-xs font-semibold">
                          Facebook
                        </label>
                        <p className="text-sm">{customer.facebookHandle}</p>
                      </div>
                    )}
                  </div>
                </>
              )}

              {/* Notes */}
              {customer.notes && (
                <>
                  <div className="divider"></div>
                  <div>
                    <label className="label label-text text-xs font-semibold">
                      ملاحظات
                    </label>
                    <p className="text-sm">{customer.notes}</p>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Stats Card */}
        <div className="space-y-6">
          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <h2 className="card-title">الإحصائيات</h2>
              <div className="space-y-4">
                <div>
                  <label className="label label-text text-xs">
                    إجمالي الطلبات
                  </label>
                  <p className="text-2xl font-bold">{customer.totalOrders}</p>
                </div>
                <div>
                  <label className="label label-text text-xs">
                    إجمالي المبيعات
                  </label>
                  <p className="text-2xl font-bold">
                    {customer.totalSpent.toFixed(2)} {business.currency}
                  </p>
                </div>
                <div>
                  <label className="label label-text text-xs">
                    متوسط قيمة الطلب
                  </label>
                  <p className="text-2xl font-bold">
                    {customer.totalOrders > 0
                      ? (customer.totalSpent / customer.totalOrders).toFixed(2)
                      : '0.00'}{' '}
                    {business.currency}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <h2 className="card-title">إجراءات</h2>
              <div className="space-y-2">
                <button className="btn btn-primary btn-block btn-sm" disabled>
                  إنشاء طلب جديد
                </button>
                <button className="btn btn-outline btn-block btn-sm">
                  تعديل المعلومات
                </button>
                <button
                  className="btn btn-error btn-outline btn-block btn-sm"
                  onClick={() => setShowDeleteModal(true)}
                >
                  حذف العميل
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Order History */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <h2 className="card-title">سجل الطلبات</h2>
          <div className="flex min-h-[200px] items-center justify-center">
            <p className="text-base-content/50">لا توجد طلبات بعد</p>
          </div>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      {showDeleteModal && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="text-lg font-bold">تأكيد الحذف</h3>
            <p className="py-4">
              هل أنت متأكد من حذف العميل "{customer.fullName}"؟ هذا الإجراء لا
              يمكن التراجع عنه.
            </p>
            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() => setShowDeleteModal(false)}
                disabled={deleteMutation.isPending}
              >
                إلغاء
              </button>
              <button
                className="btn btn-error"
                onClick={handleDelete}
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending ? (
                  <span className="loading loading-spinner loading-sm"></span>
                ) : (
                  'حذف'
                )}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button onClick={() => setShowDeleteModal(false)}>close</button>
          </form>
        </dialog>
      )}
    </div>
  )
}

/**
 * Customer Detail Skeleton
 *
 * Content-aware skeleton matching customer detail page structure
 */
function CustomerDetailSkeleton() {
  return (
    <div className="space-y-6">
      {/* Header Skeleton */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-4">
          <div className="skeleton h-10 w-10 rounded-full"></div>
          <div className="space-y-2">
            <div className="skeleton h-8 w-48"></div>
            <div className="skeleton h-4 w-32"></div>
          </div>
        </div>
        <div className="flex gap-2">
          <div className="skeleton h-10 w-20"></div>
          <div className="skeleton h-10 w-16"></div>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Customer Info Skeleton */}
        <div className="card bg-base-100 shadow lg:col-span-2">
          <div className="card-body">
            <div className="skeleton h-6 w-32 mb-4"></div>
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <div className="skeleton h-20 w-20 shrink-0 rounded-full"></div>
                <div className="space-y-2">
                  <div className="skeleton h-6 w-40"></div>
                  <div className="skeleton h-4 w-32"></div>
                </div>
              </div>
              <div className="divider"></div>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="skeleton h-16 w-full"></div>
                <div className="skeleton h-16 w-full"></div>
              </div>
            </div>
          </div>
        </div>

        {/* Stats Skeleton */}
        <div className="space-y-6">
          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="skeleton h-6 w-24 mb-4"></div>
              <div className="space-y-4">
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="space-y-2">
                    <div className="skeleton h-4 w-32"></div>
                    <div className="skeleton h-8 w-24"></div>
                  </div>
                ))}
              </div>
            </div>
          </div>
          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="skeleton h-6 w-20 mb-4"></div>
              <div className="space-y-2">
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="skeleton h-10 w-full"></div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Order History Skeleton */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="skeleton h-6 w-32 mb-4"></div>
          <div className="skeleton h-40 w-full"></div>
        </div>
      </div>
    </div>
  )
}
