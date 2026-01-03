import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Suspense, useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { Edit, Layers, Package, Trash2 } from 'lucide-react'
import toast from 'react-hot-toast'
import { useQueryClient } from '@tanstack/react-query'

import type { TableColumn } from '@/components/organisms/Table'
import type { Product } from '@/api/inventory'
import type { SortOption } from '@/components/organisms/SortButton'
import {
  Avatar,
  Dialog,
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
  useDeleteProductMutation,
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
import { queryKeys } from '@/lib/queryKeys'
import { translateErrorAsync } from '@/lib/translateError'
import { formatDateShort } from '@/lib/formatDate'

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
  const queryClient = useQueryClient()

  const [selectedProductId, setSelectedProductId] = useState<string | null>(
    null,
  )
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [isAddSheetOpen, setIsAddSheetOpen] = useState(false)
  const [isEditSheetOpen, setIsEditSheetOpen] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)

  const deleteProductMutation = useDeleteProductMutation(
    businessDescriptor,
    selectedProductId || '',
    {
      onSuccess: () => {
        toast.success(t('product_deleted', { ns: 'inventory' }))
        queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all })
        setSelectedProductId(null)
        setIsEditSheetOpen(false)
      },
      onError: async (error) => {
        const message = await translateErrorAsync(error, t)
        toast.error(message)
      },
    },
  )

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

  const products = productsResponse.data?.items ?? []
  const totalItems = productsResponse.data?.total_count ?? 0
  const totalPages = productsResponse.data?.total_pages ?? 0
  const categories = categoriesResponse.data ?? []

  const categoryOptions: Array<{ value: string; label: string }> = [
    { value: '', label: t('all_categories', { ns: 'inventory' }) },
    ...categories.map((cat) => ({
      value: cat.id,
      label: cat.name,
    })),
  ]

  const stockStatusOptions: Array<{ value: string; label: string }> = [
    { value: 'in_stock', label: t('in_stock', { ns: 'inventory' }) },
    { value: 'low_stock', label: t('low_stock', { ns: 'inventory' }) },
    { value: 'out_of_stock', label: t('out_of_stock', { ns: 'inventory' }) },
  ]

  const activeFilterCount =
    (search.categoryId ? 1 : 0) + (search.stockStatus ? 1 : 0)

  const sortOptions = useMemo<Array<SortOption>>(
    () => [
      { value: 'name', label: t('product_name', { ns: 'inventory' }) },
      { value: 'variantsCount', label: t('variants', { ns: 'inventory' }) },
      { value: 'stock', label: t('stock_quantity', { ns: 'inventory' }) },
      { value: 'costPrice', label: t('cost_price', { ns: 'inventory' }) },
      { value: 'createdAt', label: t('date_added', { ns: 'inventory' }) },
    ],
    [t],
  )

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

  const handleDeleteClick = (product: Product) => {
    setSelectedProduct(product)
    setIsDeleteDialogOpen(true)
  }

  const handleDeleteProduct = () => {
    if (selectedProductId) {
      deleteProductMutation.mutate()
    }
  }

  const handleSuccessAdd = () => {
    setIsAddSheetOpen(false)
    queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all })
  }

  const handleSuccessEdit = () => {
    setIsEditSheetOpen(false)
    setSelectedProductId(null)
    queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all })
  }

  const columns: Array<TableColumn<Product>> = [
    {
      key: 'name',
      label: t('product_name', { ns: 'inventory' }),
      sortable: true,
      render: (product: Product) => (
        <div className="flex items-center gap-3">
          <Avatar
            src={product.photos[0]?.thumbnailUrl || product.photos[0]?.url}
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
      label: t('category', { ns: 'inventory' }),
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
      key: 'costPrice',
      label: t('cost_price', { ns: 'inventory' }),
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
      label: t('sale_price', { ns: 'inventory' }),
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
      key: 'variantsCount',
      label: t('variants', { ns: 'inventory' }),
      sortable: true,
      align: 'center',
      render: (product: Product) => {
        const variantsCount = product.variants?.length ?? 0
        if (variantsCount > 1) {
          return (
            <div className="flex items-center justify-center gap-1">
              <Layers size={16} className="text-primary" />
              <span className="font-medium">{variantsCount}</span>
            </div>
          )
        }
        return <span className="text-base-content/60">{variantsCount}</span>
      },
    },
    {
      key: 'stock',
      label: t('stock_quantity', { ns: 'inventory' }),
      sortable: true,
      render: (product: Product) => {
        const totalStock = calculateTotalStock(product.variants)
        const isLowStock = hasLowStock(product.variants)
        const isOutOfStock = totalStock === 0

        let colorClass = ''
        let tooltipText = ''

        if (isOutOfStock) {
          colorClass = 'text-error font-semibold'
          tooltipText = t('out_of_stock', { ns: 'inventory' })
        } else if (isLowStock) {
          colorClass = 'text-warning font-semibold'
          tooltipText = t('low_stock', { ns: 'inventory' })
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
      key: 'createdAt',
      label: t('date_added', { ns: 'inventory' }),
      sortable: true,
      render: (product: Product) => (
        <span className="text-sm text-base-content/70">
          {formatDateShort(product.createdAt)}
        </span>
      ),
    },
    {
      key: 'actions',
      label: t('common.actions'),
      render: (product: Product) => (
        <div className="flex gap-1">
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
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation()
              handleDeleteClick(product)
            }}
            className="btn btn-ghost btn-sm text-error hover:bg-error/10"
            aria-label={t('common.delete')}
            title={t('common.delete')}
          >
            <Trash2 size={16} />
          </button>
        </div>
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

  const handleSortApply = (sortBy: string, sortOrder: 'asc' | 'desc') => {
    void navigate({
      to: '.',
      search: {
        ...search,
        sortBy,
        sortOrder,
        page: 1,
      },
    })
  }

  return (
    <>
      <ResourceListLayout
        title={t('title', { ns: 'inventory' })}
        subtitle={t('subtitle', { ns: 'inventory' })}
        addButtonText={t('add_product', { ns: 'inventory' })}
        onAddClick={() => {
          setIsAddSheetOpen(true)
        }}
        searchPlaceholder={t('search_placeholder', { ns: 'inventory' })}
        searchValue={search.search ?? ''}
        onSearchChange={handleSearch}
        filterTitle={t('filters_title', { ns: 'inventory' })}
        filterButtonText={t('common.filter')}
        filterButton={
          <filterForm.AppForm>
            <div className="space-y-6 p-4">
              <filterForm.AppField name="categoryId">
                {(field) => (
                  <field.SelectField
                    label={t('filter_by_category', { ns: 'inventory' })}
                    options={categoryOptions}
                    disabled={categoriesResponse.isLoading}
                    clearable
                  />
                )}
              </filterForm.AppField>
              <filterForm.AppField name="stockStatus">
                {(field) => (
                  <field.RadioField
                    label={t('filter_by_stock', { ns: 'inventory' })}
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
        sortTitle={t('sort_products', { ns: 'inventory' })}
        sortOptions={sortOptions}
        onSortApply={handleSortApply}
        emptyIcon={<Package size={48} />}
        emptyTitle={
          search.search
            ? t('no_results', { ns: 'inventory' })
            : t('no_products', { ns: 'inventory' })
        }
        emptyMessage={
          search.search
            ? t('try_different_search', { ns: 'inventory' })
            : t('get_started_message', { ns: 'inventory' })
        }
        emptyActionText={
          !search.search
            ? t('add_first_product', { ns: 'inventory' })
            : undefined
        }
        onEmptyAction={
          !search.search
            ? () => {
                setIsAddSheetOpen(true)
              }
            : undefined
        }
        noResultsTitle={t('no_results', { ns: 'inventory' })}
        noResultsMessage={t('try_different_search', { ns: 'inventory' })}
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
        itemsName={t('title', { ns: 'inventory' }).toLowerCase()}
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
        onSuccess={handleSuccessAdd}
      />

      <EditProductSheet
        productId={selectedProductId}
        businessDescriptor={businessDescriptor}
        isOpen={isEditSheetOpen}
        onClose={() => {
          setIsEditSheetOpen(false)
          setSelectedProductId(null)
        }}
        onSuccess={handleSuccessEdit}
        onDelete={handleDeleteProduct}
      />

      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => {
          setIsDeleteDialogOpen(false)
          setSelectedProduct(null)
          setSelectedProductId(null)
        }}
        title={t('delete_confirm_title', { ns: 'inventory' })}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => {
                setIsDeleteDialogOpen(false)
                setSelectedProduct(null)
                setSelectedProductId(null)
              }}
              disabled={deleteProductMutation.isPending}
            >
              {t('common.cancel')}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={() => {
                handleDeleteProduct()
                setIsDeleteDialogOpen(false)
              }}
              disabled={deleteProductMutation.isPending}
            >
              {deleteProductMutation.isPending ? (
                <>
                  <span className="loading loading-spinner loading-sm"></span>
                  {t('common.deleting')}
                </>
              ) : (
                t('common.delete')
              )}
            </button>
          </div>
        }
      >
        <p className="text-base-content/70">
          {selectedProduct && (
            <span
              dangerouslySetInnerHTML={{
                __html: t('delete_confirm_message', {
                  ns: 'inventory',
                  name: `<strong>${selectedProduct.name}</strong>`,
                }),
              }}
            />
          )}
        </p>
      </Dialog>
    </>
  )
}
