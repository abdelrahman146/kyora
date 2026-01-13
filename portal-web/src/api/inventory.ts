import {
  keepPreviousData,
  queryOptions,
  useMutation,
  useQuery,
} from '@tanstack/react-query'

import { delVoid, get, patch, post } from './client'
import type { UseMutationOptions } from '@tanstack/react-query'
import type { AssetReference } from '@/types/asset'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'

export type { AssetReference }

/**
 * Inventory API Types
 * Based on backend inventory domain response DTOs
 */

export interface Category {
  id: string
  businessId: string
  name: string
  descriptor: string
  createdAt: string
  updatedAt: string
}

export interface Product {
  id: string
  businessId: string
  name: string
  description: string
  photos: Array<AssetReference>
  categoryId: string
  category?: Category
  variants?: Array<Variant>
  createdAt: string
  updatedAt: string
}

export interface Variant {
  id: string
  businessId: string
  productId: string
  name: string
  code: string
  sku: string
  costPrice: string
  salePrice: string
  currency: string
  photos: Array<AssetReference>
  stockQuantity: number
  stockQuantityAlert: number
  product?: Product
  createdAt: string
  updatedAt: string
}

export interface ListResponse<T> {
  items: Array<T>
  page: number
  page_size: number
  total_count: number
  total_pages: number
  has_more: boolean
}

export type ListProductsResponse = ListResponse<Product>
export type ListVariantsResponse = ListResponse<Variant>
export type ListCategoriesResponse = Array<Category>

export interface CreateProductRequest {
  name: string
  description?: string
  photos?: Array<AssetReference>
  categoryId: string
}

export interface UpdateProductRequest {
  name?: string
  description?: string
  photos?: Array<AssetReference>
  categoryId?: string
}

export interface CreateVariantRequest {
  product_id: string
  code: string
  sku?: string
  photos?: Array<AssetReference>
  costPrice: string
  salePrice: string
  stockQuantity: number
  stockQuantityAlert: number
}

export interface CreateProductVariantRequest {
  code: string
  sku?: string
  photos?: Array<AssetReference>
  costPrice: string
  salePrice: string
  stockQuantity: number
  stockQuantityAlert: number
}

export interface CreateProductWithVariantsRequest {
  product: CreateProductRequest
  variants: Array<CreateProductVariantRequest>
}

export interface UpdateVariantRequest {
  code?: string
  sku?: string
  photos?: Array<AssetReference>
  costPrice?: string
  salePrice?: string
  currency?: string
  stockQuantity?: number
  stockQuantityAlert?: number
}

export interface CreateCategoryRequest {
  name: string
  descriptor: string
}

export interface UpdateCategoryRequest {
  name?: string
  descriptor?: string
}

export interface InventorySummaryResponse {
  products_count: number
  variants_count: number
  categories_count: number
  low_stock_variants_count: number
  out_of_stock_variants_count: number
  total_stock_units: number
  inventory_value: string
  top_products_by_inventory_value: Array<Product>
}

/**
 * Inventory API Client
 */
export const inventoryApi = {
  /**
   * List products for a business
   */
  async listProducts(
    businessDescriptor: string,
    params?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
      categoryId?: string
      stockStatus?: 'in_stock' | 'low_stock' | 'out_of_stock'
    },
  ): Promise<ListProductsResponse> {
    const searchParams = new URLSearchParams()
    if (params?.search) searchParams.set('search', params.search)
    if (params?.page) searchParams.set('page', params.page.toString())
    if (params?.pageSize)
      searchParams.set('pageSize', params.pageSize.toString())
    if (params?.orderBy && params.orderBy.length > 0) {
      params.orderBy.forEach((o) => searchParams.append('orderBy', o))
    }
    if (params?.categoryId) searchParams.set('categoryId', params.categoryId)
    if (params?.stockStatus) searchParams.set('stockStatus', params.stockStatus)

    const query = searchParams.toString() ? `?${searchParams.toString()}` : ''
    return get<ListProductsResponse>(
      `v1/businesses/${businessDescriptor}/inventory/products${query}`,
    )
  },

  /**
   * Get product by ID with variants
   */
  async getProduct(
    businessDescriptor: string,
    productId: string,
  ): Promise<Product> {
    return get<Product>(
      `v1/businesses/${businessDescriptor}/inventory/products/${productId}`,
    )
  },

  /**
   * Create a new product
   */
  async createProduct(
    businessDescriptor: string,
    data: CreateProductRequest,
  ): Promise<Product> {
    return post<Product>(
      `v1/businesses/${businessDescriptor}/inventory/products`,
      { json: data },
    )
  },

  /**
   * Create a product with variants atomically
   */
  async createProductWithVariants(
    businessDescriptor: string,
    data: CreateProductWithVariantsRequest,
  ): Promise<Product> {
    return post<Product>(
      `v1/businesses/${businessDescriptor}/inventory/products/with-variants`,
      { json: data },
    )
  },

  /**
   * Update a product
   */
  async updateProduct(
    businessDescriptor: string,
    productId: string,
    data: UpdateProductRequest,
  ): Promise<Product> {
    return patch<Product>(
      `v1/businesses/${businessDescriptor}/inventory/products/${productId}`,
      { json: data },
    )
  },

  /**
   * Delete a product (cascades to variants)
   */
  async deleteProduct(
    businessDescriptor: string,
    productId: string,
  ): Promise<void> {
    return delVoid(
      `v1/businesses/${businessDescriptor}/inventory/products/${productId}`,
    )
  },

  /**
   * List all categories for a business (no pagination)
   */
  async listCategories(
    businessDescriptor: string,
  ): Promise<ListCategoriesResponse> {
    return get<ListCategoriesResponse>(
      `v1/businesses/${businessDescriptor}/inventory/categories`,
    )
  },

  /**
   * Get category by ID
   */
  async getCategory(
    businessDescriptor: string,
    categoryId: string,
  ): Promise<Category> {
    return get<Category>(
      `v1/businesses/${businessDescriptor}/inventory/categories/${categoryId}`,
    )
  },

  /**
   * Create a new category
   */
  async createCategory(
    businessDescriptor: string,
    data: CreateCategoryRequest,
  ): Promise<Category> {
    return post<Category>(
      `v1/businesses/${businessDescriptor}/inventory/categories`,
      { json: data },
    )
  },

  /**
   * Update a category
   */
  async updateCategory(
    businessDescriptor: string,
    categoryId: string,
    data: UpdateCategoryRequest,
  ): Promise<Category> {
    return patch<Category>(
      `v1/businesses/${businessDescriptor}/inventory/categories/${categoryId}`,
      { json: data },
    )
  },

  /**
   * Delete a category
   */
  async deleteCategory(
    businessDescriptor: string,
    categoryId: string,
  ): Promise<void> {
    return delVoid(
      `v1/businesses/${businessDescriptor}/inventory/categories/${categoryId}`,
    )
  },

  /**
   * List variants for a product
   */
  async listVariants(
    businessDescriptor: string,
    productId: string,
    params?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
    },
  ): Promise<ListVariantsResponse> {
    const searchParams = new URLSearchParams()
    if (params?.search) searchParams.set('search', params.search)
    if (params?.page) searchParams.set('page', params.page.toString())
    if (params?.pageSize)
      searchParams.set('pageSize', params.pageSize.toString())
    if (params?.orderBy && params.orderBy.length > 0) {
      params.orderBy.forEach((o) => searchParams.append('orderBy', o))
    }

    const query = searchParams.toString() ? `?${searchParams.toString()}` : ''
    return get<ListVariantsResponse>(
      `v1/businesses/${businessDescriptor}/inventory/products/${productId}/variants${query}`,
    )
  },

  /**
   * Create a new variant for a product
   */
  async createVariant(
    businessDescriptor: string,
    data: CreateVariantRequest,
  ): Promise<Variant> {
    return post<Variant>(
      `v1/businesses/${businessDescriptor}/inventory/variants`,
      { json: data },
    )
  },

  /**
   * Update a variant
   */
  async updateVariant(
    businessDescriptor: string,
    variantId: string,
    data: UpdateVariantRequest,
  ): Promise<Variant> {
    return patch<Variant>(
      `v1/businesses/${businessDescriptor}/inventory/variants/${variantId}`,
      { json: data },
    )
  },

  /**
   * Delete a variant
   */
  async deleteVariant(
    businessDescriptor: string,
    variantId: string,
  ): Promise<void> {
    return delVoid(
      `v1/businesses/${businessDescriptor}/inventory/variants/${variantId}`,
    )
  },

  /**
   * Get inventory summary with metrics
   */
  async getInventorySummary(
    businessDescriptor: string,
    topLimit?: number,
  ): Promise<InventorySummaryResponse> {
    const searchParams = new URLSearchParams()
    if (topLimit !== undefined)
      searchParams.set('topLimit', topLimit.toString())

    const query = searchParams.toString() ? `?${searchParams.toString()}` : ''
    return get<InventorySummaryResponse>(
      `v1/businesses/${businessDescriptor}/inventory/summary${query}`,
    )
  },
}

/**
 * Query Options Factories
 */
export const inventoryQueries = {
  /**
   * List products query
   */
  list: (
    businessDescriptor: string,
    params?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
      categoryId?: string
      stockStatus?: 'in_stock' | 'low_stock' | 'out_of_stock'
    },
  ) =>
    queryOptions({
      queryKey: queryKeys.inventory.list(businessDescriptor, {
        search: params?.search,
        page: params?.page,
        limit: params?.pageSize,
        categoryId: params?.categoryId,
        stockStatus: params?.stockStatus,
        orderBy: params?.orderBy,
      }),
      queryFn: () => inventoryApi.listProducts(businessDescriptor, params),
      staleTime: STALE_TIME.ONE_MINUTE,
      enabled: !!businessDescriptor,
      placeholderData: keepPreviousData,
    }),

  /**
   * Get product detail query
   */
  detail: (businessDescriptor: string, productId: string) =>
    queryOptions({
      queryKey: queryKeys.inventory.detail(businessDescriptor, productId),
      queryFn: () => inventoryApi.getProduct(businessDescriptor, productId),
      staleTime: STALE_TIME.ONE_MINUTE,
      enabled: !!businessDescriptor && !!productId,
    }),

  /**
   * List categories query
   */
  categories: (businessDescriptor: string) =>
    queryOptions({
      queryKey: [...queryKeys.inventory.all, 'categories', businessDescriptor],
      queryFn: () => inventoryApi.listCategories(businessDescriptor),
      staleTime: STALE_TIME.FIVE_MINUTES,
      enabled: !!businessDescriptor,
    }),

  /**
   * Get inventory summary query
   */
  summary: (businessDescriptor: string, topLimit?: number) =>
    queryOptions({
      queryKey: [
        ...queryKeys.inventory.all,
        'summary',
        businessDescriptor,
        topLimit,
      ],
      queryFn: () =>
        inventoryApi.getInventorySummary(businessDescriptor, topLimit),
      staleTime: STALE_TIME.ONE_MINUTE,
      enabled: !!businessDescriptor,
    }),
}

/**
 * Query Hooks
 */

/**
 * Hook to list products
 */
export function useProductsQuery(
  businessDescriptor: string,
  params?: {
    search?: string
    page?: number
    pageSize?: number
    orderBy?: Array<string>
    categoryId?: string
    stockStatus?: 'in_stock' | 'low_stock' | 'out_of_stock'
  },
) {
  return useQuery(inventoryQueries.list(businessDescriptor, params))
}

/**
 * Hook to get product detail
 */
export function useProductQuery(businessDescriptor: string, productId: string) {
  return useQuery(inventoryQueries.detail(businessDescriptor, productId))
}

/**
 * Hook to list categories
 */
export function useCategoriesQuery(businessDescriptor: string) {
  return useQuery(inventoryQueries.categories(businessDescriptor))
}

/**
 * Hook to get inventory summary
 */
export function useInventorySummaryQuery(
  businessDescriptor: string,
  topLimit?: number,
) {
  return useQuery(inventoryQueries.summary(businessDescriptor, topLimit))
}

/**
 * Mutation Hooks
 */

/**
 * Hook to create a product
 */
export function useCreateProductMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Product, Error, CreateProductRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateProductRequest) =>
      inventoryApi.createProduct(businessDescriptor, data),
    ...options,
  })
}

/**
 * Hook to create a product with variants
 */
export function useCreateProductWithVariantsMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<
    Product,
    Error,
    CreateProductWithVariantsRequest
  >,
) {
  return useMutation({
    mutationFn: (data: CreateProductWithVariantsRequest) =>
      inventoryApi.createProductWithVariants(businessDescriptor, data),
    ...options,
  })
}

/**
 * Hook to update a product
 */
export function useUpdateProductMutation(
  businessDescriptor: string,
  productId: string,
  options?: UseMutationOptions<Product, Error, UpdateProductRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateProductRequest) =>
      inventoryApi.updateProduct(businessDescriptor, productId, data),
    ...options,
  })
}

/**
 * Hook to delete a product
 */
export function useDeleteProductMutation(
  businessDescriptor: string,
  productId: string,
  options?: UseMutationOptions<void, Error, void>,
) {
  return useMutation({
    mutationFn: () => inventoryApi.deleteProduct(businessDescriptor, productId),
    ...options,
  })
}

/**
 * Hook to create a category
 */
export function useCreateCategoryMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Category, Error, CreateCategoryRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateCategoryRequest) =>
      inventoryApi.createCategory(businessDescriptor, data),
    ...options,
  })
}

/**
 * Hook to update a category
 */
export function useUpdateCategoryMutation(
  businessDescriptor: string,
  categoryId: string,
  options?: UseMutationOptions<Category, Error, UpdateCategoryRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateCategoryRequest) =>
      inventoryApi.updateCategory(businessDescriptor, categoryId, data),
    ...options,
  })
}

/**
 * Hook to delete a category
 */
export function useDeleteCategoryMutation(
  businessDescriptor: string,
  categoryId: string,
  options?: UseMutationOptions<void, Error, void>,
) {
  return useMutation({
    mutationFn: () =>
      inventoryApi.deleteCategory(businessDescriptor, categoryId),
    ...options,
  })
}

/**
 * Hook to create a variant
 */
export function useCreateVariantMutation(
  businessDescriptor: string,
  productId: string,
  options?: UseMutationOptions<Variant, Error, CreateVariantRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateVariantRequest) =>
      inventoryApi.createVariant(businessDescriptor, {
        ...data,
        product_id: productId,
      }),
    ...options,
  })
}

/**
 * Hook to update a variant
 */
export function useUpdateVariantMutation(
  businessDescriptor: string,
  variantId: string,
  options?: UseMutationOptions<Variant, Error, UpdateVariantRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateVariantRequest) =>
      inventoryApi.updateVariant(businessDescriptor, variantId, data),
    ...options,
  })
}

export interface UpdateVariantByIdRequest {
  variantId: string
  data: UpdateVariantRequest
}

/**
 * Hook to update a variant when the variantId is only known at call time.
 */
export function useUpdateVariantByIdMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Variant, Error, UpdateVariantByIdRequest>,
) {
  return useMutation({
    mutationFn: ({ variantId, data }: UpdateVariantByIdRequest) =>
      inventoryApi.updateVariant(businessDescriptor, variantId, data),
    ...options,
  })
}

/**
 * Hook to delete a variant
 */
export function useDeleteVariantMutation(
  businessDescriptor: string,
  variantId: string,
  options?: UseMutationOptions<void, Error, void>,
) {
  return useMutation({
    mutationFn: () => inventoryApi.deleteVariant(businessDescriptor, variantId),
    ...options,
  })
}

/**
 * Hook to delete a variant when the variantId is only known at call time.
 */
export function useDeleteVariantByIdMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<void, Error, string>,
) {
  return useMutation({
    mutationFn: (variantId: string) =>
      inventoryApi.deleteVariant(businessDescriptor, variantId),
    ...options,
  })
}
