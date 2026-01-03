/**
 * ResourceListLayout Template Component
 *
 * A highly flexible and reusable template for list pages that share common patterns:
 * - Header with title, subtitle, and action button
 * - Search and filter toolbar
 * - Responsive views (Desktop: Table, Mobile: Cards)
 * - Empty states with CTAs
 * - Pagination
 * - Loading states
 *
 * This component uses composition and render props to maximize flexibility
 * while maintaining consistency across different resource types (customers, inventory, etc.).
 *
 * @example Customer List
 * ```tsx
 * <ResourceListLayout
 *   title={t('customers.title')}
 *   subtitle={t('customers.subtitle')}
 *   addButtonText={t('customers.add_customer')}
 *   onAddClick={() => setIsAddSheetOpen(true)}
 *   searchPlaceholder={t('customers.search_placeholder')}
 *   searchValue={search.search ?? ''}
 *   onSearchChange={handleSearch}
 *   filterTitle={t('customers.filters')}
 *   filterButton={<FilterContent />}
 *   activeFilterCount={filterCount}
 *   onApplyFilters={handleApply}
 *   onResetFilters={handleReset}
 *   emptyIcon={<Users size={48} />}
 *   emptyTitle={t('customers.no_customers')}
 *   emptyMessage={t('customers.get_started_message')}
 *   emptyActionText={t('customers.add_first_customer')}
 *   onEmptyAction={() => setIsAddSheetOpen(true)}
 *   tableColumns={columns}
 *   tableData={customers}
 *   tableKeyExtractor={(c) => c.id}
 *   tableSortBy={search.sortBy}
 *   tableSortOrder={search.sortOrder}
 *   onTableSort={handleSort}
 *   onTableRowClick={handleCustomerClick}
 *   mobileCard={(customer) => <CustomerCard customer={customer} />}
 *   isLoading={isLoading}
 *   hasSearchQuery={!!search.search}
 *   currentPage={search.page}
 *   totalPages={totalPages}
 *   pageSize={search.pageSize}
 *   totalItems={totalItems}
 *   onPageChange={handlePageChange}
 *   onPageSizeChange={handlePageSizeChange}
 *   itemsName={t('customers.customers').toLowerCase()}
 * />
 * ```
 */

import { Plus } from 'lucide-react'

import { Pagination } from '../molecules/Pagination'
import { SearchInput } from '../molecules/SearchInput'
import { FilterButton } from '../organisms/FilterButton'
import { SortButton } from '../organisms/SortButton'
import { Table } from '../organisms/Table'
import type { ReactNode } from 'react'
import type { TableColumn } from '../organisms/Table'
import type { SortOption } from '../organisms/SortButton'
import { useMediaQuery } from '@/hooks/useMediaQuery'

export interface ResourceListLayoutProps<T> {
  /** Page title */
  title: string
  /** Page subtitle */
  subtitle?: string
  /** Add button text */
  addButtonText: string
  /** Add button click handler */
  onAddClick: () => void
  /** Whether add button should be disabled */
  addButtonDisabled?: boolean
  /** Search placeholder text */
  searchPlaceholder: string
  /** Current search value */
  searchValue: string
  /** Search change handler */
  onSearchChange: (value: string) => void
  /** Filter drawer title */
  filterTitle: string
  /** Filter button text (default: "Filter") */
  filterButtonText?: string
  /** Filter content (rendered inside drawer) */
  filterButton: ReactNode
  /** Active filter count */
  activeFilterCount?: number
  /** Apply filters handler */
  onApplyFilters?: () => void
  /** Reset filters handler */
  onResetFilters?: () => void
  /** Apply button label */
  applyLabel?: string
  /** Reset button label */
  resetLabel?: string
  /** Sort sheet title */
  sortTitle?: string
  /** Sort button text */
  sortButtonText?: string
  /** Available sort options */
  sortOptions?: Array<SortOption>
  /** Callback when sort is applied */
  onSortApply?: (sortBy: string, sortOrder: 'asc' | 'desc') => void
  /** Empty state icon */
  emptyIcon: ReactNode
  /** Empty state title */
  emptyTitle: string
  /** Empty state message */
  emptyMessage?: string
  /** Empty state action text (for CTA button) */
  emptyActionText?: string
  /** Empty state action handler */
  onEmptyAction?: () => void
  /** No results title (when search query exists) */
  noResultsTitle?: string
  /** No results message (when search query exists) */
  noResultsMessage?: string
  /** Table columns for desktop view */
  tableColumns: Array<TableColumn<T>>
  /** Table data */
  tableData: Array<T>
  /** Table key extractor */
  tableKeyExtractor: (item: T) => string
  /** Table sort by column key */
  tableSortBy?: string
  /** Table sort order */
  tableSortOrder?: 'asc' | 'desc'
  /** Table sort handler */
  onTableSort?: (key: string) => void
  /** Table row click handler */
  onTableRowClick?: (item: T) => void
  /** Mobile card renderer */
  mobileCard: (item: T) => ReactNode
  /** Loading state */
  isLoading: boolean
  /** Whether there's an active search query */
  hasSearchQuery: boolean
  /** Current page number */
  currentPage: number
  /** Total number of pages */
  totalPages: number
  /** Page size */
  pageSize: number
  /** Total items count */
  totalItems: number
  /** Page change handler */
  onPageChange: (page: number) => void
  /** Page size change handler */
  onPageSizeChange?: (pageSize: number) => void
  /** Items name for pagination (e.g., "customers", "products") */
  itemsName: string
  /** Additional toolbar content (optional) */
  toolbarExtra?: ReactNode
  /** Custom skeleton component for loading state (optional) */
  skeleton?: ReactNode
}

export function ResourceListLayout<T>({
  title,
  subtitle,
  addButtonText,
  onAddClick,
  addButtonDisabled = false,
  searchPlaceholder,
  searchValue,
  onSearchChange,
  filterTitle,
  filterButtonText,
  filterButton,
  activeFilterCount,
  onApplyFilters,
  onResetFilters,
  applyLabel,
  resetLabel,
  sortTitle,
  sortButtonText,
  sortOptions,
  onSortApply,
  emptyIcon,
  emptyTitle,
  emptyMessage,
  emptyActionText,
  onEmptyAction,
  noResultsTitle,
  noResultsMessage,
  tableColumns,
  tableData,
  tableKeyExtractor,
  tableSortBy,
  tableSortOrder,
  onTableSort,
  onTableRowClick,
  mobileCard,
  isLoading,
  hasSearchQuery,
  currentPage,
  totalPages,
  pageSize,
  totalItems,
  onPageChange,
  onPageSizeChange,
  itemsName,
  toolbarExtra,
  skeleton,
}: ResourceListLayoutProps<T>) {
  const isMobile = useMediaQuery('(max-width: 768px)')

  const isEmpty = !isLoading && tableData.length === 0
  const shouldShowEmptyState = isEmpty && !hasSearchQuery
  const shouldShowNoResults = isEmpty && hasSearchQuery
  const shouldShowContent = !isEmpty || isLoading

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">{title}</h1>
          {subtitle && (
            <p className="text-sm text-base-content/60 mt-1">{subtitle}</p>
          )}
        </div>
        <button
          type="button"
          className="btn btn-primary gap-2"
          onClick={onAddClick}
          disabled={addButtonDisabled}
          aria-disabled={addButtonDisabled}
        >
          <Plus size={20} />
          {addButtonText}
        </button>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="flex-1">
          <SearchInput
            value={searchValue}
            onChange={onSearchChange}
            placeholder={searchPlaceholder}
          />
        </div>
        <div className="flex gap-2">
          {/* Sort Button (Mobile Only) */}
          {isMobile && sortOptions && sortOptions.length > 0 && onSortApply && (
            <div className="flex-1">
              <SortButton
                title={sortTitle || 'Sort'}
                buttonText={sortButtonText}
                sortOptions={sortOptions}
                currentSortBy={tableSortBy}
                currentSortOrder={tableSortOrder}
                onApply={onSortApply}
              />
            </div>
          )}
          <div className="flex-1">
            <FilterButton
              title={filterTitle}
              buttonText={filterButtonText}
              activeCount={activeFilterCount}
              onApply={onApplyFilters}
              onReset={onResetFilters}
              applyLabel={applyLabel}
              resetLabel={resetLabel}
            >
              {filterButton}
            </FilterButton>
          </div>
          {toolbarExtra}
        </div>
      </div>

      {/* Empty State (No Items) */}
      {shouldShowEmptyState && (
        <div className="card bg-base-100 shadow">
          <div className="card-body items-center text-center py-12">
            <div className="text-base-content/20 mb-4">{emptyIcon}</div>
            <h3 className="text-lg font-semibold mb-2">{emptyTitle}</h3>
            {emptyMessage && (
              <p className="text-sm text-base-content/70 mb-4">
                {emptyMessage}
              </p>
            )}
            {emptyActionText && onEmptyAction && (
              <button
                type="button"
                className="btn btn-primary btn-sm gap-2"
                onClick={onEmptyAction}
              >
                <Plus size={18} />
                {emptyActionText}
              </button>
            )}
          </div>
        </div>
      )}

      {/* No Results State (Search Query Exists) */}
      {shouldShowNoResults && (
        <div className="card bg-base-100 shadow">
          <div className="card-body items-center text-center py-12">
            <div className="text-base-content/20 mb-4">{emptyIcon}</div>
            <h3 className="text-lg font-semibold mb-2">
              {noResultsTitle ?? emptyTitle}
            </h3>
            {noResultsMessage && (
              <p className="text-sm text-base-content/70 mb-4">
                {noResultsMessage}
              </p>
            )}
          </div>
        </div>
      )}

      {/* Desktop: Table View */}
      {!isMobile && shouldShowContent && (
        <>
          {isLoading && tableData.length === 0 && skeleton ? (
            skeleton
          ) : (
            <div className="overflow-x-auto">
              <Table
                columns={tableColumns}
                data={tableData}
                keyExtractor={tableKeyExtractor}
                isLoading={isLoading}
                emptyMessage={emptyTitle}
                sortBy={tableSortBy}
                sortOrder={tableSortOrder}
                onSort={onTableSort}
                onRowClick={onTableRowClick}
              />
            </div>
          )}
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            pageSize={pageSize}
            totalItems={totalItems}
            onPageChange={onPageChange}
            onPageSizeChange={onPageSizeChange}
            itemsName={itemsName}
          />
        </>
      )}

      {/* Mobile: Card View with Pagination */}
      {isMobile && shouldShowContent && (
        <>
          <div className="space-y-3">
            {isLoading && tableData.length === 0 ? (
              (skeleton ??
              Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="skeleton h-32 rounded-box" />
              )))
            ) : tableData.length === 0 ? (
              <div className="text-center py-12 text-base-content/60">
                {emptyTitle}
              </div>
            ) : (
              tableData.map((item) => (
                <div key={tableKeyExtractor(item)}>{mobileCard(item)}</div>
              ))
            )}
          </div>
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            pageSize={pageSize}
            totalItems={totalItems}
            onPageChange={onPageChange}
            onPageSizeChange={onPageSizeChange}
            itemsName={itemsName}
          />
        </>
      )}
    </div>
  )
}
