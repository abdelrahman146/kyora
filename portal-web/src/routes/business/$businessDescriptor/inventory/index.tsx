import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { Edit, Package } from 'lucide-react'

import type { TableColumn } from '@/components/organisms/Table'
import type { Product } from '@/api/inventory'
import {
  Avatar,
  InventoryCard,
  InventoryListSkeleton,
  ResourceListLayout,
  Tooltip,
} from '@/components'
import { AddProductSheet } from '@/components/organisms/AddProductSheet'
import { EditProductSheet } from '@/components/organisms/EditProductSheet'
import { ProductDetailsSheet } from '@/components/organisms/ProductDetailsSheet'
import {
  inventoryQueries,
  useCategoriesQuery,
  useProductsQuery,
} from '@/api/inventory'
import { formatCurrency } from '@/lib/formatCurrency'
import {
  calculateTotalStock,
  getPriceRange,
  hasLowStock,
} from '@/lib/inventoryUtils'
import { getSelectedBusiness } from '@/stores/businessStore'
import { useKyoraForm } from '@/lib/form/useKyoraForm'

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

  const [selectedProductId, setSelectedProductId] = useState<string | null>(
    null,
  )
  const [isAddSheetOpen, setIsAddSheetOpen] = useState(false)
  const [isEditSheetOpen, setIsEditSheetOpen] = useState(false)

  const orderBy = useMemo<Array<string>>(() => {
    if (!search.sortBy) return ['-createdAt']
    return [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
  }, [search.sortBy, search.sortOrder])

  const filterForm = useKyoraForm({
    defaultValues: {
      categoryId: search.categoryId ?? '',
      stockStatus: search.stockStatus ?? '',
    },
    onSubmit: async ({ value }) => {
      const stockStatus =
        value.stockStatus === '' ? undefined : value.stockStatus
      await navigate({
        to: '.',
        search: {
          ...search,
          categoryId: value.categoryId || undefined,
          stockStatus: stockStatus as
            | 'in_stock'
            | 'low_stock'
            | 'out_of_stock'
            | undefined,
          page: 1,
        },
      })
    },
  })

  // Sync form with URL search params when they change
  useEffect(() => {
    filterForm.setFieldValue('categoryId', search.categoryId ?? '')
    filterForm.setFieldValue('stockStatus', search.stockStatus ?? '')
  }, [search.categoryId, search.stockStatus])

  const productsResponse = useProductsQuery(businessDescriptor, {
    search: search.search,
    page: search.page,
    pageSize: search.pageSize,
    orderBy,
    categoryId: search.categoryId,
    stockStatus: search.stockStatus,
  })

  const categoriesResponse = useCategoriesQuery(businessDescriptor)

  // Clear previous data when filters change to fix the bug
  const products = useMemo(() => {
    if (productsResponse.isFetching && productsResponse.data) {
      // If we're fetching new data and the query key has changed, return empty array
      // to prevent showing stale data
      return []
    }
    return productsResponse.data?.items ?? []
  }, [productsResponse.data, productsResponse.isFetching])

  const totalItems = productsResponse.data?.total_count ?? 0
  const totalPages = productsResponse.data?.total_pages ?? 0
  const categories = categoriesResponse.data ?? []

  const categoryOptions: Array<{ value: string; label: string }> = [
    { value: '', label: t('inventory.all_categories') },
    ...categories.map((cat) => ({
      value: cat.id,
      label: cat.name,
    })),
  ]

  const stockStatusOptions: Array<{ value: string; label: string }> = [
    { value: 'in_stock', label: t('inventory.in_stock') },
    { value: 'low_stock', label: t('inventory.low_stock') },
    { value: 'out_of_stock', label: t('inventory.out_of_stock') },
  ]

  const activeFilterCount =
    (search.categoryId ? 1 : 0) + (search.stockStatus ? 1 : 0)

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

  const handleResetFilters = async () => {
    filterForm.reset()
    await navigate({
      to: '.',
      search: {
        ...search,
        categoryId: undefined,
        stockStatus: undefined,
        page: 1,
      },
    })
  }

  const handleProductClick = (product: Product) => {
    setSelectedProductId(product.id)
  }

  const handleEditProduct = (product: Product) => {
    setSelectedProductId(product.id)
    setIsEditSheetOpen(true)
  }
  const columns: Array<TableColumn<Product>> = [
    {
      key: 'name',
      label: t('inventory.product_name'),
      sortable: true,
      render: (product: Product) => (
        <div className="flex items-center gap-3">
          <Avatar
            src={product.photos[0]?.thumbnail_url || product.photos[0]?.url}
            alt={product.name}
            fallback={product.name.charAt(0).toUpperCase()}
            size="sm"
            shape="square"
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
      render: (product: Product) => {
        const category = categories.find((c) => c.id === product.categoryId)
        return (
          <span className="text-sm text-base-content/70">
            {category ? category.name : '-'}
          </span>
        )
      },
    },
    {
      key: 'cost_price',
      label: t('inventory.cost_price'),
      sortable: true,
      render: (product: Product) => {
        const priceRange = getPriceRange(product.variants, 'costPrice')
        if (priceRange.isSame) {
          return <span>{formatCurrency(priceRange.min, currency)}</span>
        }
        return (
          <span>
            {formatCurrency(priceRange.min, currency)} -{' '}
            {formatCurrency(priceRange.max, currency)}
          </span>
        )
      },
    },
    {
      key: 'sale_price',
      label: t('inventory.sale_price'),
      render: (product: Product) => {
        const priceRange = getPriceRange(product.variants, 'salePrice')
        if (priceRange.isSame) {
          return <span>{formatCurrency(priceRange.min, currency)}</span>
        }
        return (
          <span>
            {formatCurrency(priceRange.min, currency)} -{' '}
            {formatCurrency(priceRange.max, currency)}
          </span>
        )
      },
    },
    {
      key: 'stock',
      label: t('inventory.stock_quantity'),
      render: (product: Product) => {
        const totalStock = calculateTotalStock(product.variants)
        const isLowStock = hasLowStock(product.variants)
        const isOutOfStock = totalStock === 0

        let colorClass = ''
        let tooltipText = ''

        if (isOutOfStock) {
          colorClass = 'text-error font-semibold'
          tooltipText = t('inventory.out_of_stock')
        } else if (isLowStock) {
          colorClass = 'text-warning font-semibold'
          tooltipText = t('inventory.low_stock')
        }

        if (tooltipText) {
          return (
            <Tooltip content={tooltipText}>
              <span className={colorClass}>{totalStock}</span>
            </Tooltip>
          )
        }

        return <span>{totalStock}</span>
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

  const handleSort = (key: string) => {
    void navigate({
      to: '.',
      search: {
        ...search,
        sortBy: key,
        sortOrder:
          search.sortBy === key && search.sortOrder === 'asc' ? 'desc' : 'asc',
        page: 1,
      },
    })
  }

  return (
    <>
      <ResourceListLayout
        title={t('inventory.title')}
        subtitle={t('inventory.subtitle')}
        addButtonText={t('inventory.add_product')}
        onAddClick={() => {
          setIsAddSheetOpen(true)
        }}
        searchPlaceholder={t('inventory.search_placeholder')}
        searchValue={search.search ?? ''}
        onSearchChange={handleSearch}
        filterTitle={t('inventory.filters_title')}
        filterButtonText={t('common.filter')}
        filterButton={
          <filterForm.AppForm>
            <div className="space-y-6 p-4">
              <filterForm.AppField name="categoryId">
                {(field) => (
                  <field.SelectField
                    label={t('inventory.filter_by_category')}
                    options={categoryOptions}
                    disabled={categoriesResponse.isLoading}
                    clearable
                  />
                )}
              </filterForm.AppField>
              <filterForm.AppField name="stockStatus">
                {(field) => (
                  <field.RadioField
                    label={t('inventory.filter_by_stock')}
                    options={stockStatusOptions}
                    orientation="vertical"
                  />
                )}
              </filterForm.AppField>
            </div>
          </filterForm.AppForm>
        }
        activeFilterCount={activeFilterCount}
        onApplyFilters={() => {
          filterForm.handleSubmit()
        }}
        onResetFilters={handleResetFilters}
        applyLabel={t('common.apply')}
        resetLabel={t('common.reset')}
        emptyIcon={<Package size={48} />}
        emptyTitle={
          search.search ? t('inventory.no_results') : t('inventory.no_products')
        }
        emptyMessage={
          search.search
            ? t('inventory.try_different_search')
            : t('inventory.get_started_message')
        }
        emptyActionText={
          !search.search ? t('inventory.add_first_product') : undefined
        }
        onEmptyAction={
          !search.search
            ? () => {
                setIsAddSheetOpen(true)
              }
            : undefined
        }
        noResultsTitle={t('inventory.no_results')}
        noResultsMessage={t('inventory.try_different_search')}
        tableColumns={columns}
        tableData={products}
        tableKeyExtractor={(product) => product.id}
        tableSortBy={search.sortBy}
        tableSortOrder={search.sortOrder}
        onTableSort={handleSort}
        onTableRowClick={handleProductClick}
        mobileCard={(product) => (
          <InventoryCard
            product={product}
            currency={currency}
            categories={categories}
            onClick={() => {
              handleProductClick(product)
            }}
          />
        )}
        isLoading={productsResponse.isLoading}
        hasSearchQuery={!!search.search}
        currentPage={search.page}
        totalPages={totalPages}
        pageSize={search.pageSize}
        totalItems={totalItems}
        onPageChange={handlePageChange}
        onPageSizeChange={handlePageSizeChange}
        itemsName={t('inventory.title').toLowerCase()}
        skeleton={<InventoryListSkeleton />}
      />

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
