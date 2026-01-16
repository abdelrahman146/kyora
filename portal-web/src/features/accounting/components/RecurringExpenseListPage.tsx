/**
 * RecurringExpenseListPage Component
 *
 * List page for recurring expense templates with:
 * - Filterable list (status)
 * - Desktop: Table view with edit/delete/status actions
 * - Mobile: Card view with RecurringExpenseCard
 * - Add recurring expense navigates to expenses page with recurring toggle
 *
 * Note: Backend does NOT support search for recurring expenses.
 * The search input is hidden per UX spec.
 */

import { useNavigate, useParams, useSearch } from '@tanstack/react-router'
import { useEffect, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Receipt, Repeat } from 'lucide-react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-hot-toast'

import { RecurringExpensesSearchSchema } from '../schema/recurringExpensesSearch'
import { categoryColors, categoryIcons } from '../schema/options'
import { RecurringExpenseCard } from './RecurringExpenseCard'
import { RecurringExpenseQuickActions } from './RecurringExpenseQuickActions'
import { ExpensesTabs } from './ExpenseListPage'
import type { RecurringExpensesSearch } from '../schema/recurringExpensesSearch'
import type { RecurringExpense } from '@/api/accounting'
import type { RecurringExpenseStatus } from '@/api/types/accounting'
import type { TableColumn } from '@/components/organisms/Table'
import type { SortOption } from '@/components/molecules/SortButton'
import { accountingQueries, useRecurringExpensesQuery } from '@/api/accounting'
import { ResourceListLayout } from '@/components/templates/ResourceListLayout'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'
import { getSelectedBusiness } from '@/stores/businessStore'
import { useKyoraForm } from '@/lib/form'

const STATUS_OPTIONS: Array<{
  value: RecurringExpenseStatus
  labelKey: string
}> = [
  { value: 'active', labelKey: 'status.active' },
  { value: 'paused', labelKey: 'status.paused' },
  { value: 'ended', labelKey: 'status.ended' },
  { value: 'canceled', labelKey: 'status.canceled' },
]

const statusColors: Record<string, string> = {
  active: 'badge-success',
  paused: 'badge-warning',
  ended: 'badge-ghost',
  canceled: 'badge-error',
}

export function RecurringExpenseListPage() {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')

  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/accounting/expenses/recurring/',
  })

  const rawSearch = useSearch({
    from: '/business/$businessDescriptor/accounting/expenses/recurring/',
  })
  const search = RecurringExpensesSearchSchema.parse(rawSearch)

  const navigate = useNavigate({
    from: '/business/$businessDescriptor/accounting/expenses/recurring',
  })

  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'
  const queryClient = useQueryClient()

  // Build orderBy array from search params
  const orderBy = useMemo<Array<string>>(() => {
    if (!search.sortBy) return ['-createdAt']
    return [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
  }, [search.sortBy, search.sortOrder])

  // Filter form for drawer
  const filterForm = useKyoraForm({
    defaultValues: {
      status: search.status ?? '',
    },
    onSubmit: async ({ value }) => {
      const status =
        value.status === ''
          ? undefined
          : (value.status as RecurringExpenseStatus | undefined)

      await navigate({
        to: '.',
        search: (prev) => ({
          ...prev,
          status,
          page: 1,
        }),
      })
    },
  })

  // Sync filter form with URL state
  useEffect(() => {
    filterForm.setFieldValue('status', search.status ?? '')
  }, [search.status])

  // Fetch recurring expenses
  const { data: recurringExpensesResponse, isLoading } =
    useRecurringExpensesQuery(businessDescriptor, {
      page: search.page,
      pageSize: search.pageSize,
      orderBy,
    })

  // Client-side filter by status (backend doesn't support status filter yet)
  const allRecurringExpenses = recurringExpensesResponse?.items ?? []
  const recurringExpenses = search.status
    ? allRecurringExpenses.filter((re) => re.status === search.status)
    : allRecurringExpenses

  const totalItems = search.status
    ? recurringExpenses.length
    : (recurringExpensesResponse?.totalCount ?? 0)
  const totalPages = Math.ceil(totalItems / search.pageSize)

  // Build status select options for filter
  const statusOptions: Array<{ value: string; label: string }> = [
    { value: '', label: t('filters.all_statuses') },
    ...STATUS_OPTIONS.map((opt) => ({
      value: opt.value,
      label: t(opt.labelKey),
    })),
  ]

  // Sort options
  const sortOptions = useMemo<Array<SortOption>>(
    () => [
      { value: 'createdAt', label: tCommon('date_added') },
      { value: 'amount', label: t('form.amount') },
      { value: 'category', label: t('category.label') },
    ],
    [t, tCommon],
  )

  // Calculate active filter count
  const activeFilterCount = search.status ? 1 : 0

  // Handlers
  const handlePageChange = (newPage: number) => {
    void navigate({
      to: '.',
      search: (prev) => ({ ...prev, page: newPage }),
    })
  }

  const handlePageSizeChange = (newPageSize: number) => {
    void navigate({
      to: '.',
      search: (prev) => ({ ...prev, pageSize: newPageSize, page: 1 }),
    })
  }

  const handleSortApply = (sortBy: string, sortOrder: 'asc' | 'desc') => {
    void navigate({
      to: '.',
      search: (prev) => ({ ...prev, sortBy, sortOrder, page: 1 }),
    })
  }

  const handleTableSort = (key: string) => {
    const newOrder =
      search.sortBy === key && search.sortOrder === 'asc' ? 'desc' : 'asc'
    void navigate({
      to: '.',
      search: (prev) => ({
        ...prev,
        sortBy: key,
        sortOrder: newOrder,
        page: 1,
      }),
    })
  }

  const handleResetFilters = async () => {
    filterForm.reset()
    await navigate({
      to: '.',
      search: (prev) => ({
        ...prev,
        status: undefined,
        page: 1,
      }),
    })
  }

  const handleAddClick = () => {
    // Navigate to expenses page and show a toast to use the recurring toggle
    toast(t('recurring.add_hint'), { icon: '💡' })
    void navigate({
      to: '/business/$businessDescriptor/accounting/expenses',
      params: { businessDescriptor },
      search: {
        page: 1,
        pageSize: 20,
        sortBy: 'occurredOn',
        sortOrder: 'desc',
      },
    })
  }

  const handleActionComplete = async () => {
    await invalidateQueries()
  }

  const invalidateQueries = async () => {
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.recurringExpenses(),
    })
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.expenses(),
    })
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.summaries(),
    })
    await queryClient.invalidateQueries({
      queryKey: accountingQueries.recentActivitiesKey(),
    })
  }

  // Table columns for desktop view
  const tableColumns: Array<TableColumn<RecurringExpense>> = useMemo(
    () => [
      {
        key: 'category',
        label: t('category.label'),
        sortable: true,
        width: 'w-48',
        render: (item) => {
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          const CategoryIcon = categoryIcons[item.category] ?? Receipt
          const colorClass =
            // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
            categoryColors[item.category] ?? 'bg-base-200 text-base-content/70'
          return (
            <div className="flex items-center gap-3">
              <div
                className={`w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0 ${colorClass}`}
              >
                <CategoryIcon className="w-4 h-4" />
              </div>
              <span className="font-medium text-base-content">
                {t(`category.${item.category}`)}
              </span>
            </div>
          )
        },
      },
      {
        key: 'frequency',
        label: t('frequency.label'),
        width: 'w-28',
        render: (item) => (
          <span className="badge badge-ghost badge-sm gap-1">
            <Repeat className="w-3 h-3" />
            {t(`frequency.${item.frequency}`)}
          </span>
        ),
      },
      {
        key: 'status',
        label: t('status.label'),
        width: 'w-24',
        render: (item) => (
          <span className={`badge badge-sm ${statusColors[item.status]}`}>
            {t(`status.${item.status}`)}
          </span>
        ),
      },
      {
        key: 'nextRecurringDate',
        label: t('recurring.next_date'),
        width: 'w-32',
        render: (item) =>
          item.nextRecurringDate && item.status === 'active' ? (
            <span className="text-base-content/80">
              {formatDateShort(item.nextRecurringDate)}
            </span>
          ) : (
            <span className="text-base-content/40">—</span>
          ),
      },
      {
        key: 'amount',
        label: t('form.amount'),
        sortable: true,
        align: 'start',
        width: 'w-38',
        render: (item) => (
          <span className="font-bold text-error tabular-nums">
            -{formatCurrency(parseFloat(item.amount), currency)}
          </span>
        ),
      },
      {
        key: 'actions',
        label: tCommon('actionsLabel'),
        width: 'w-12',
        align: 'center',
        render: (item) => (
          <RecurringExpenseQuickActions
            recurringExpense={item}
            businessDescriptor={businessDescriptor}
            onActionComplete={handleActionComplete}
          />
        ),
      },
    ],
    [t, tCommon, currency, businessDescriptor],
  )

  // Skeleton component for loading state
  const RecurringExpenseListSkeleton = () => (
    <div className="space-y-4">
      {[...Array(5)].map((_, i) => (
        <div
          key={i}
          className="skeleton h-28 w-full rounded-xl bg-base-200/50"
        />
      ))}
    </div>
  )

  return (
    <>
      {/* Tab Navigation */}
      <ExpensesTabs
        businessDescriptor={businessDescriptor}
        activeTab="recurring"
      />

      <ResourceListLayout<RecurringExpense>
        // Header
        title={t('header.recurring_templates')}
        addButtonText={t('actions.add_recurring')}
        onAddClick={handleAddClick}
        // Search - disabled since backend doesn't support it
        hideSearch
        // Filters
        filterTitle={t('filters.title')}
        activeFilterCount={activeFilterCount}
        onApplyFilters={() => void filterForm.handleSubmit()}
        onResetFilters={handleResetFilters}
        filterButton={
          <filterForm.AppForm>
            <div className="space-y-4 p-4">
              {/* Status Filter */}
              <filterForm.AppField name="status">
                {(field) => (
                  <field.SelectField
                    label={t('status.label')}
                    options={statusOptions}
                  />
                )}
              </filterForm.AppField>
            </div>
          </filterForm.AppForm>
        }
        // Sort
        sortTitle={tCommon('sort')}
        sortOptions={sortOptions}
        onSortApply={handleSortApply}
        // Empty State
        emptyIcon={<Repeat size={48} />}
        emptyTitle={t('empty.recurring_title')}
        emptyMessage={t('empty.recurring_description')}
        emptyActionText={t('actions.add_first_recurring')}
        onEmptyAction={handleAddClick}
        // Table (Desktop)
        tableColumns={tableColumns}
        tableData={recurringExpenses}
        tableKeyExtractor={(item) => item.id}
        tableSortBy={search.sortBy}
        tableSortOrder={search.sortOrder}
        onTableSort={handleTableSort}
        // Mobile Card
        mobileCard={(item) => (
          <RecurringExpenseCard
            recurringExpense={item}
            currency={currency}
            businessDescriptor={businessDescriptor}
            onActionComplete={handleActionComplete}
          />
        )}
        // Loading & Pagination
        isLoading={isLoading}
        hasSearchQuery={false}
        currentPage={search.page}
        totalPages={totalPages}
        pageSize={search.pageSize}
        totalItems={totalItems}
        onPageChange={handlePageChange}
        onPageSizeChange={handlePageSizeChange}
        itemsName={t('header.recurring_templates').toLowerCase()}
        skeleton={<RecurringExpenseListSkeleton />}
      />
    </>
  )
}

/**
 * Route loader for prefetching recurring expenses data
 */
export async function recurringExpenseListLoader({
  queryClient,
  businessDescriptor,
  search,
}: {
  queryClient: any
  businessDescriptor: string
  search: RecurringExpensesSearch
}) {
  const orderBy = search.sortBy
    ? [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
    : ['-createdAt']

  await queryClient.ensureQueryData(
    accountingQueries.recurringExpenseList(businessDescriptor, {
      page: search.page,
      pageSize: search.pageSize,
      orderBy,
    }),
  )
}
