/**
 * ExpenseListPage Component
 *
 * List page for expenses with:
 * - Filterable list (date range, category)
 * - Desktop: Table view with edit/delete actions
 * - Mobile: Card view with ExpenseCard
 * - Add/Edit Expense sheets
 * - URL-driven pagination and filters
 *
 * Note: Backend does NOT support search for expenses.
 * The search input is hidden per UX spec.
 */

import { useNavigate, useParams, useSearch } from '@tanstack/react-router'
import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Receipt, Repeat } from 'lucide-react'
import { useQueryClient } from '@tanstack/react-query'
import { format, parseISO } from 'date-fns'

import { ExpensesSearchSchema } from '../schema/expensesSearch'
import {
  CATEGORY_OPTIONS,
  categoryColors,
  categoryIcons,
} from '../schema/options'
import { ExpenseCard } from './ExpenseCard'
import { ExpenseListSkeleton } from './ExpenseListSkeleton'
import { ExpenseQuickActions } from './ExpenseQuickActions'
import { CreateExpenseSheet } from './sheets/CreateExpenseSheet'
import type { ExpensesSearch } from '../schema/expensesSearch'
import type { Expense, ExpenseCategory } from '@/api/accounting'
import type { TableColumn } from '@/components/organisms/Table'
import type { SortOption } from '@/components/molecules/SortButton'
import type { DateRange } from 'react-day-picker'
import { accountingQueries, useExpensesQuery } from '@/api/accounting'
import { ResourceListLayout } from '@/components/templates/ResourceListLayout'
import { DateRangePicker } from '@/components/form/DateRangePicker'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'
import { getSelectedBusiness } from '@/stores/businessStore'
import { useKyoraForm } from '@/lib/form'

export function ExpenseListPage() {
  const { t } = useTranslation('accounting')
  const { t: tCommon } = useTranslation('common')

  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/accounting/expenses/',
  })

  const rawSearch = useSearch({
    from: '/business/$businessDescriptor/accounting/expenses/',
  })
  const search = ExpensesSearchSchema.parse(rawSearch)

  const navigate = useNavigate({
    from: '/business/$businessDescriptor/accounting/expenses',
  })

  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'
  const queryClient = useQueryClient()

  // Sheet state
  const [isCreateSheetOpen, setIsCreateSheetOpen] = useState(false)

  // Build orderBy array from search params
  const orderBy = useMemo<Array<string>>(() => {
    if (!search.sortBy) return ['-occurredOn']
    return [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
  }, [search.sortBy, search.sortOrder])

  // Helper to parse URL date string to Date object
  const parseDateParam = (dateStr: string | undefined): Date | undefined => {
    if (!dateStr) return undefined
    try {
      return parseISO(dateStr)
    } catch {
      return undefined
    }
  }

  // Filter form for drawer
  const filterForm = useKyoraForm({
    defaultValues: {
      category: search.category ?? '',
      dateRange:
        search.from && search.to
          ? ({
              from: parseDateParam(search.from),
              to: parseDateParam(search.to),
            } as DateRange)
          : undefined,
    },
    onSubmit: async ({ value }) => {
      // Parse category - empty string means no filter
      const category =
        value.category === ''
          ? undefined
          : (value.category as ExpenseCategory | undefined)

      await navigate({
        to: '.',
        search: (prev) => ({
          ...prev,
          category,
          from: value.dateRange?.from
            ? format(value.dateRange.from, 'yyyy-MM-dd')
            : undefined,
          to: value.dateRange?.to
            ? format(value.dateRange.to, 'yyyy-MM-dd')
            : undefined,
          page: 1,
        }),
      })
    },
  })

  // Sync filter form with URL state
  useEffect(() => {
    filterForm.setFieldValue('category', search.category ?? '')
    filterForm.setFieldValue(
      'dateRange',
      search.from && search.to
        ? ({
            from: parseDateParam(search.from),
            to: parseDateParam(search.to),
          } as DateRange)
        : undefined,
    )
  }, [search.category, search.from, search.to])

  // Fetch expenses
  const { data: expensesResponse, isLoading } = useExpensesQuery(
    businessDescriptor,
    {
      page: search.page,
      pageSize: search.pageSize,
      orderBy,
      category: search.category,
      from: search.from,
      to: search.to,
    },
  )

  const expenses = expensesResponse?.items ?? []
  const totalItems = expensesResponse?.totalCount ?? 0
  const totalPages = Math.ceil(totalItems / search.pageSize)

  // Build category select options for filter
  const categoryOptions: Array<{ value: string; label: string }> = [
    { value: '', label: t('filters.all_categories') },
    ...CATEGORY_OPTIONS.map((opt) => ({
      value: opt.value,
      label: t(opt.labelKey),
    })),
  ]

  // Sort options
  const sortOptions = useMemo<Array<SortOption>>(
    () => [
      { value: 'occurredOn', label: t('form.date') },
      { value: 'amount', label: t('form.amount') },
      { value: 'category', label: t('category.label') },
      { value: 'createdAt', label: tCommon('date_added') },
    ],
    [t, tCommon],
  )

  // Calculate active filter count
  const activeFilterCount =
    (search.category ? 1 : 0) + (search.from && search.to ? 1 : 0)

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
        category: undefined,
        from: undefined,
        to: undefined,
        page: 1,
      }),
    })
  }

  const handleExpenseCreated = async () => {
    setIsCreateSheetOpen(false)
    await invalidateQueries()
  }

  const handleActionComplete = async () => {
    await invalidateQueries()
  }

  const invalidateQueries = async () => {
    // Invalidate expenses list and summary queries
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
  const tableColumns: Array<TableColumn<Expense>> = useMemo(
    () => [
      {
        key: 'category',
        label: t('category.label'),
        sortable: true,
        width: 'w-48',
        render: (expense) => {
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          const CategoryIcon = categoryIcons[expense.category] ?? Receipt
          const colorClass =
            // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
            categoryColors[expense.category] ??
            'bg-base-200 text-base-content/70'
          return (
            <div className="flex items-center gap-3">
              {/* Category Icon */}
              <div
                className={`w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0 ${colorClass}`}
              >
                <CategoryIcon className="w-4 h-4" />
              </div>
              <span className="font-medium text-base-content">
                {t(`category.${expense.category}`)}
              </span>
              {expense.recurringExpenseId && (
                <span className="badge badge-ghost badge-sm gap-1 text-secondary">
                  <Repeat className="w-3 h-3" />
                </span>
              )}
            </div>
          )
        },
      },
      {
        key: 'note',
        label: t('form.note'),
        render: (expense) => (
          <span className="text-base-content/70 line-clamp-2 max-w-md">
            {expense.note || (
              <span className="text-base-content/40 italic">
                {t('list.no_note')}
              </span>
            )}
          </span>
        ),
      },
      {
        key: 'occurredOn',
        label: t('form.date'),
        sortable: true,
        width: 'w-32',
        render: (expense) => (
          <span className="text-base-content/80">
            {formatDateShort(expense.occurredOn)}
          </span>
        ),
      },
      {
        key: 'amount',
        label: t('form.amount'),
        sortable: true,
        align: 'start',
        width: 'w-38',
        render: (expense) => (
          <span className="font-bold text-error tabular-nums">
            -{formatCurrency(parseFloat(expense.amount), currency)}
          </span>
        ),
      },
      {
        key: 'actions',
        label: tCommon('actionsLabel'),
        width: 'w-12',
        align: 'center',
        render: (expense) => (
          <ExpenseQuickActions
            expense={expense}
            businessDescriptor={businessDescriptor}
            currency={currency}
            onActionComplete={handleActionComplete}
          />
        ),
      },
    ],
    [t, tCommon, currency, businessDescriptor],
  )

  return (
    <>
      <ResourceListLayout<Expense>
        // Header
        title={t('header.expenses')}
        addButtonText={t('actions.add_expense')}
        onAddClick={() => setIsCreateSheetOpen(true)}
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
              {/* Category Filter */}
              <filterForm.AppField name="category">
                {(field) => (
                  <field.SelectField
                    label={t('filters.category')}
                    options={categoryOptions}
                  />
                )}
              </filterForm.AppField>

              {/* Date Range Filter */}
              <filterForm.AppField name="dateRange">
                {(field) => (
                  <div className="form-control">
                    <label className="label pb-2">
                      <span className="label-text font-medium">
                        {t('filters.date_range')}
                      </span>
                    </label>
                    <DateRangePicker
                      value={field.state.value}
                      onChange={(range) => field.handleChange(range)}
                      placeholder={t('filters.select_date_range')}
                    />
                  </div>
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
        emptyIcon={<Receipt size={48} />}
        emptyTitle={t('empty.expenses_title')}
        emptyMessage={t('empty.expenses_description')}
        emptyActionText={t('actions.add_first_expense')}
        onEmptyAction={() => setIsCreateSheetOpen(true)}
        // Table (Desktop)
        tableColumns={tableColumns}
        tableData={expenses}
        tableKeyExtractor={(expense) => expense.id}
        tableSortBy={search.sortBy}
        tableSortOrder={search.sortOrder}
        onTableSort={handleTableSort}
        // Mobile Card
        mobileCard={(expense) => (
          <ExpenseCard
            expense={expense}
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
        itemsName={t('header.expenses').toLowerCase()}
        skeleton={<ExpenseListSkeleton />}
      />

      {/* Create Expense Sheet */}
      <CreateExpenseSheet
        isOpen={isCreateSheetOpen}
        onClose={() => setIsCreateSheetOpen(false)}
        businessDescriptor={businessDescriptor}
        onCreated={handleExpenseCreated}
      />
    </>
  )
}

/**
 * Route loader for prefetching expenses data
 */
export async function expenseListLoader({
  queryClient,
  businessDescriptor,
  search,
}: {
  queryClient: any
  businessDescriptor: string
  search: ExpensesSearch
}) {
  const orderBy = search.sortBy
    ? [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
    : ['-occurredOn']

  await queryClient.ensureQueryData(
    accountingQueries.expenseList(businessDescriptor, {
      page: search.page,
      pageSize: search.pageSize,
      orderBy,
      category: search.category,
      from: search.from,
      to: search.to,
    }),
  )
}
