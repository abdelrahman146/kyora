import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { Edit, Package, Plus } from 'lucide-react'

import type { TableColumn } from '@/components/organisms/Table'
import type { Product } from '@/api/inventory'
import {
  Avatar,
  Badge,
  Button,
  InventoryCard,
  InventoryListSkeleton,
  Pagination,
  SearchInput,
  Table,
  Tooltip,
} from '@/components'
import { AddProductSheet } from '@/components/organisms/AddProductSheet'
import { EditProductSheet } from '@/components/organisms/EditProductSheet'
import { FilterButton } from '@/components/organisms/FilterButton'
import { ProductDetailsSheet } from '@/components/organisms/ProductDetailsSheet'
import { FormRadio } from '@/components/atoms/FormRadio'
import { FormSelect } from '@/components/atoms/FormSelect'
import {
  inventoryQueries,
  useCategoriesQuery,
  useProductsQuery,
} from '@/api/inventory'
import { useMediaQuery } from '@/hooks'
import { formatCurrency } from '@/lib/formatCurrency'
import {
  calculateAverageCostPrice,
  calculateTotalStock,
  getStockStatus,
  getStockStatusBadgeVariant,
  hasLowStock,
} from '@/lib/inventoryUtils'
import { getSelectedBusiness } from '@/stores/businessStore'

/**
 * Search schema for inventory list
 */
const InventorySearchSchema = z.object({
  search: z.string().optional(),
  page: z.number().optional().default(1),
  pageSize: z.number().optional().default(20),
  sortBy: z.string().optional(),
  sortOrder: z.enum(['asc', 'desc']).optional().default('asc'),
  categoryId: z.string().optional(),
  stockStatus: z.enum(['in_stock', 'low_stock', 'out_of_stock']).optional(),
})

type InventorySearch = z.infer<typeof InventorySearchSchema>

export const Route = createFileRoute(
  '/business/$businessDescriptor/inventory/',
)({
  staticData: {
    titleKey: 'inventory.title',
  },
  validateSearch: (search): InventorySearch => {
    return InventorySearchSchema.parse(search)
  },
  loader: async ({ context, params, location }) => {
    const { queryClient } = context as any

    // Parse search params
    const searchParams = InventorySearchSchema.parse(location.search)

    // Build orderBy
    let orderBy: Array<string> | undefined
    if (searchParams.sortBy) {
      const prefix = searchParams.sortOrder === 'desc' ? '-' : ''
      orderBy = [`${prefix}${searchParams.sortBy}`]
    }

    // Prefetch products
    await queryClient.prefetchQuery(
      inventoryQueries.list(params.businessDescriptor, {
        search: searchParams.search,
        page: searchParams.page,
        pageSize: searchParams.pageSize,
        orderBy,
        categoryId: searchParams.categoryId,
        stockStatus: searchParams.stockStatus,
      }),
    )

    // Prefetch categories
    await queryClient.prefetchQuery(
      inventoryQueries.categories(params.businessDescriptor),
    )
  },
  component: () => (
    <Suspense fallback={<InventoryListSkeleton />}>
      <InventoryListPage />
    </Suspense>
  ),
})

function InventoryListPage() {
  const { t } = useTranslation()
  const { businessDescriptor } = Route.useParams()
  const navigate = useNavigate()
  const search = Route.useSearch()
  const business = getSelectedBusiness()
  const currency = business?.currency ?? 'USD'

  // Local state
  const [selectedProductId, setSelectedProductId] = useState<string | null>(
    null,
  )
  const [isAddSheetOpen, setIsAddSheetOpen] = useState(false)
  const [isEditSheetOpen, setIsEditSheetOpen] = useState(false)

  // Filter state
  const [categoryIdFilter, setCategoryIdFilter] = useState<string | undefined>(
    search.categoryId,
  )
  const [stockStatusFilter, setStockStatusFilter] = useState<
    'in_stock' | 'low_stock' | 'out_of_stock' | undefined
  >(search.stockStatus)

  // Media query for responsive layout
  const isMobile = useMediaQuery('(max-width: 768px)')

  // Build orderBy from search params
  const orderBy = useMemo<Array<string>>(() => {
    if (!search.sortBy) return ['-createdAt']
    return [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
  }, [search.sortBy, search.sortOrder])

  // Fetch products
  const productsResponse = useProductsQuery(businessDescriptor, {
    search: search.search,
    page: search.page,
    pageSize: search.pageSize,
    orderBy,
    categoryId: search.categoryId,
    stockStatus: search.stockStatus,
  })

  // Fetch categories
  const categoriesResponse = useCategoriesQuery(businessDescriptor)

  const products = productsResponse.data?.items ?? []
  const totalItems = productsResponse.data?.total_count ?? 0
  const totalPages = productsResponse.data?.total_pages ?? 0
  const categories = categoriesResponse.data ?? []

  // Build category options
  const categoryOptions: Array<{ value: string; label: string }> = [
    { value: '', label: t('inventory.all_categories') },
    ...categories.map((cat) => ({
      value: cat.id,
      label: cat.name,
    })),
  ]

  // Build stock status options
  const stockStatusOptions: Array<{ value: string; label: string }> = [
    { value: 'in_stock', label: t('inventory.in_stock') },
    { value: 'low_stock', label: t('inventory.low_stock') },
    { value: 'out_of_stock', label: t('inventory.out_of_stock') },
  ]

  // Calculate active filter count
  const activeFilterCount =
    (search.categoryId ? 1 : 0) + (search.stockStatus ? 1 : 0)

  // Handle search
  const handleSearch = (value: string) => {
    void navigate({
      to: '.',
      search: {
        ...search,
        search: value || undefined,
        page: 1,
      },
    })
  }

  // Handle pagination
  const handlePageChange = (newPage: number) => {
    void navigate({
      to: '.',
      search: { ...search, page: newPage },
    })
  }

  const handlePageSizeChange = (newPageSize: number) => {
    void navigate({
      to: '.',
      search: { ...search, pageSize: newPageSize, page: 1 },
    })
  }

  // Handle apply filters
  const handleApplyFilters = () => {
    void navigate({
      to: '.',
      search: {
        ...search,
        categoryId: categoryIdFilter,
        stockStatus: stockStatusFilter,
        page: 1,
      },
    })
  }

  // Handle reset filters
  const handleResetFilters = () => {
    setCategoryIdFilter(undefined)
    setStockStatusFilter(undefined)
    void navigate({
      to: '.',
      search: {
        ...search,
        categoryId: undefined,
        stockStatus: undefined,
        page: 1,
      },
    })
  }

  // Handle product click
  const handleProductClick = (product: Product) => {
    setSelectedProductId(product.id)
  }

  // Handle edit product
  const handleEditProduct = (product: Product) => {
    setSelectedProductId(product.id)
    setIsEditSheetOpen(true)
  }

  // Table columns (desktop)
  const columns: Array<TableColumn<Product>> = [
    {
      key: 'name',
      label: t('inventory.product_name'),
      sortable: true,
      render: (product: Product) => (
        <div className="flex items-center gap-3">
          <Avatar
            src={product.photos[0]?.thumbnail_url}
            alt={product.name}
            fallback={product.name.charAt(0).toUpperCase()}
            size="sm"
          />
          <div>
            <p className="font-medium text-base-content">{product.name}</p>
          </div>
        </div>
      ),
    },
    {
      key: 'category',
      label: t('inventory.category'),
      render: (product: Product) =>
        product.category ? (
          <Badge variant="neutral">{product.category.name}</Badge>
        ) : (
          <span className="text-base-content/40">
            {t('inventory.category')}
          </span>
        ),
    },
    {
      key: 'cost_price',
      label: t('inventory.cost_price_avg'),
      sortable: true,
      render: (product: Product) => {
        const avgCost = calculateAverageCostPrice(product.variants)
        return (
          <div className="flex items-center gap-2">
            <span>{formatCurrency(avgCost, currency)}</span>
            <Tooltip content={t('inventory.average_cost_tooltip')}>
              <span className="text-base-content/60 cursor-help">ⓘ</span>
            </Tooltip>
          </div>
        )
      },
    },
    {
      key: 'stock',
      label: t('inventory.stock_quantity'),
      render: (product: Product) => {
        const totalStock = calculateTotalStock(product.variants)
        const isLowStock = hasLowStock(product.variants)

        return (
          <div className="flex items-center gap-2">
            <span>{totalStock}</span>
            {isLowStock && (
              <Badge variant="warning" size="sm">
                {t('inventory.low_stock')}
              </Badge>
            )}
          </div>
        )
      },
    },
    {
      key: 'status',
      label: t('inventory.status'),
      render: (product: Product) => {
        const status = getStockStatus(product.variants)
        const badgeVariant = getStockStatusBadgeVariant(status)

        return <Badge variant={badgeVariant}>{t(`inventory.${status}`)}</Badge>
      },
    },
    {
      key: 'actions',
      label: t('common.actions'),
      render: (product: Product) => (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation()
            handleEditProduct(product)
          }}
          className="btn btn-ghost btn-sm"
          aria-label={t('common.edit')}
          title={t('common.edit')}
        >
          <Edit size={16} />
        </button>
      ),
    },
  ]

  return (
    <>
      <div className="container mx-auto p-4 space-y-6">
        {/* Header */}
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <h1 className="text-2xl font-bold text-base-content">
              {t('inventory.title')}
            </h1>
            <p className="text-base-content/60">{t('inventory.subtitle')}</p>
          </div>
          <Button
            variant="primary"
            size="md"
            onClick={() => {
              setIsAddSheetOpen(true)
            }}
          >
            <Plus size={20} />
            {t('inventory.add_product')}
          </Button>
        </div>

        {/* Search and Filters */}
        <div className="flex flex-col gap-3 md:flex-row md:items-center">
          <div className="flex-1">
            <SearchInput
              value={search.search ?? ''}
              onChange={handleSearch}
              placeholder={t('inventory.search_placeholder')}
              debounceMs={300}
            />
          </div>
          <FilterButton
            title={t('inventory.filters_title')}
            buttonText={t('common.filter')}
            activeCount={activeFilterCount}
            onApply={handleApplyFilters}
            onReset={handleResetFilters}
            applyLabel={t('common.apply')}
            resetLabel={t('common.reset')}
          >
            <div className="space-y-6 p-4">
              {/* Category filter */}
              <FormSelect
                label={t('inventory.filter_by_category')}
                options={categoryOptions}
                value={categoryIdFilter ?? ''}
                onChange={(value) => {
                  const val = Array.isArray(value) ? value[0] : value
                  setCategoryIdFilter(val === '' ? undefined : val)
                }}
                disabled={categoriesResponse.isLoading}
              />

              {/* Stock status filter */}
              <FormRadio
                label={t('inventory.filter_by_stock')}
                name="stockStatus"
                options={stockStatusOptions}
                value={stockStatusFilter ?? ''}
                onChange={(e) => {
                  const newValue = e.target.value as
                    | 'in_stock'
                    | 'low_stock'
                    | 'out_of_stock'
                    | ''
                  setStockStatusFilter(newValue === '' ? undefined : newValue)
                }}
                orientation="vertical"
              />
            </div>
          </FilterButton>
        </div>

        {/* Content - Table (desktop) or Cards (mobile) */}
        {isMobile ? (
          // Mobile: Card Grid
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {products.map((product) => {
              const avgCost = calculateAverageCostPrice(product.variants)
              const totalStock = calculateTotalStock(product.variants)
              const stockStatus = getStockStatus(product.variants)
              const lowStock = hasLowStock(product.variants)

              return (
                <InventoryCard
                  key={product.id}
                  product={product}
                  avgCostPrice={avgCost}
                  totalStock={totalStock}
                  stockStatus={stockStatus}
                  hasLowStock={lowStock}
                  currency={currency}
                  onClick={() => {
                    handleProductClick(product)
                  }}
                />
              )
            })}
            {products.length === 0 && (
              <div className="col-span-full flex flex-col items-center justify-center py-12">
                <Package size={48} className="text-base-content/20 mb-4" />
                <p className="text-lg font-medium text-base-content/60">
                  {t('inventory.no_products')}
                </p>
                <p className="text-sm text-base-content/40 mt-2">
                  {t('inventory.get_started_message')}
                </p>
              </div>
            )}
          </div>
        ) : (
          // Desktop: Table
          <Table<Product>
            columns={columns}
            data={products}
            keyExtractor={(product) => product.id}
            isLoading={productsResponse.isLoading}
            emptyMessage={t('inventory.no_products')}
            sortBy={search.sortBy}
            sortOrder={search.sortOrder}
            onSort={(key) => {
              void navigate({
                to: '.',
                search: {
                  ...search,
                  sortBy: key,
                  sortOrder:
                    search.sortBy === key && search.sortOrder === 'asc'
                      ? 'desc'
                      : 'asc',
                  page: 1,
                },
              })
            }}
            onRowClick={(product) => {
              handleProductClick(product)
            }}
          />
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center">
            <Pagination
              currentPage={search.page}
              totalPages={totalPages}
              pageSize={search.pageSize}
              totalItems={totalItems}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
              showPageSizeSelector={!isMobile}
              itemsName={t('inventory.title').toLowerCase()}
            />
          </div>
        )}
      </div>
      {/* Sheets */}
      <ProductDetailsSheet
        product={
          selectedProductId
            ? (products.find((p) => p.id === selectedProductId) ?? null)
            : null
        }
        businessDescriptor={businessDescriptor}
        isOpen={!!selectedProductId && !isEditSheetOpen}
        onClose={() => {
          setSelectedProductId(null)
        }}
      />
      <AddProductSheet
        isOpen={isAddSheetOpen}
        onClose={() => {
          setIsAddSheetOpen(false)
        }}
      />
      <EditProductSheet
        isOpen={isEditSheetOpen}
        onClose={() => {
          setIsEditSheetOpen(false)
          setSelectedProductId(null)
        }}
      />
    </>
  )
}
