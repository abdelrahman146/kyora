import {
  keepPreviousData,
  queryOptions,
  useMutation,
  useQuery,
} from '@tanstack/react-query'

import { del, get, patch, post } from './client'
import type { UseMutationOptions } from '@tanstack/react-query'
import type { SocialPlatform } from '@/api/customer'
import { STALE_TIME } from '@/lib/queryKeys'

/**
 * Order API Types
 * Based on backend order domain model
 */

export type OrderStatus =
  | 'pending'
  | 'placed'
  | 'ready_for_shipment'
  | 'shipped'
  | 'fulfilled'
  | 'cancelled'
  | 'returned'

export type OrderPaymentStatus = 'pending' | 'paid' | 'failed' | 'refunded'

export type OrderPaymentMethod =
  | 'credit_card'
  | 'paypal'
  | 'bank_transfer'
  | 'cash_on_delivery'
  | 'tamara'
  | 'tabby'

export interface OrderCustomer {
  id: string
  name: string
  email: string | null
  phoneCode: string | null
  phoneNumber: string | null
  avatarUrl?: string
  whatsappNumber?: string | null
  instagramUsername?: string | null
  facebookUsername?: string | null
  tiktokUsername?: string | null
  xUsername?: string | null
  snapchatUsername?: string | null
}

export interface OrderAddress {
  id: string
  street: string | null
  city: string
  state: string
  zipCode: string | null
  countryCode: string
  phoneCode: string
  phoneNumber: string
}

export interface OrderShippingZone {
  id: string
  name: string
  countries: Array<string>
  currency: string
  shippingCost: string
  freeShippingThreshold: string
}

export interface OrderItem {
  id: string
  orderId: string
  productId: string
  variantId: string
  quantity: number
  currency: string
  unitPrice: string
  unitCost: string
  totalCost: string
  total: string
  product?: {
    id: string
    name: string
    photos: Array<{ url: string; thumbnailUrl?: string }>
  }
  variant?: {
    id: string
    name: string
    code: string
    sku: string
  }
}

export interface OrderNote {
  id: string
  orderId: string
  content: string
  createdAt: string
  updatedAt: string
}

export interface Order {
  id: string
  orderNumber: string
  businessId: string
  customerId: string
  customer?: OrderCustomer
  shippingAddressId: string
  shippingAddress?: OrderAddress
  shippingZone?: OrderShippingZone | null
  shippingZoneId: string | null
  channel: string
  subtotal: string
  vat: string
  vatRate: string
  shippingFee: string
  discount: string
  cogs: string
  total: string
  currency: string
  status: OrderStatus
  paymentStatus: OrderPaymentStatus
  paymentMethod: OrderPaymentMethod
  paymentReference: string | null
  placedAt: string | null
  readyForShipmentAt: string | null
  orderedAt: string
  shippedAt: string | null
  fulfilledAt: string | null
  cancelledAt: string | null
  returnedAt: string | null
  paidAt: string | null
  failedAt: string | null
  refundedAt: string | null
  items?: Array<OrderItem>
  notes?: Array<OrderNote>
  createdAt: string
  updatedAt: string
  deletedAt: string | null
}

export interface ListResponse<T> {
  hasMore: boolean
  items: Array<T>
  page: number
  pageSize: number
  totalCount: number
  totalPages: number
}

export type ListOrdersResponse = ListResponse<Order>

export interface ListOrdersFilters {
  status?: Array<OrderStatus>
  paymentStatus?: Array<OrderPaymentStatus>
  socialPlatforms?: Array<SocialPlatform>
  customerId?: string
  orderNumber?: string
  from?: string
  to?: string
}

export interface CreateOrderItemRequest {
  variantId: string
  quantity: number
  unitPrice: string
  unitCost?: string
}

export interface CreateOrderRequest {
  customerId: string
  channel: string
  shippingAddressId: string
  shippingZoneId?: string
  shippingFee?: string
  discount?: string
  paymentMethod?: OrderPaymentMethod
  paymentReference?: string
  orderedAt?: string
  items: Array<CreateOrderItemRequest>
}

export interface UpdateOrderRequest {
  shippingAddressId?: string
  shippingZoneId?: string
  shippingFee?: string
  channel?: string
  discount?: string
  orderedAt?: string
  items?: Array<CreateOrderItemRequest>
}

export interface UpdateOrderStatusRequest {
  status: OrderStatus
}

export interface UpdateOrderPaymentStatusRequest {
  paymentStatus: OrderPaymentStatus
}

export interface AddOrderPaymentDetailsRequest {
  paymentMethod: OrderPaymentMethod
  paymentReference?: string
}

export interface CreateOrderNoteRequest {
  content: string
}

export interface UpdateOrderNoteRequest {
  content: string
}

/**
 * Order API Client
 */
export const orderApi = {
  /**
   * List orders for a business
   */
  async listOrders(
    businessDescriptor: string,
    params?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
      status?: Array<OrderStatus>
      paymentStatus?: Array<OrderPaymentStatus>
      socialPlatforms?: Array<SocialPlatform>
      customerId?: string
      orderNumber?: string
      from?: string
      to?: string
    },
  ): Promise<ListOrdersResponse> {
    const searchParams = new URLSearchParams()
    if (params?.search) searchParams.set('search', params.search)
    if (params?.page) searchParams.set('page', params.page.toString())
    if (params?.pageSize)
      searchParams.set('pageSize', params.pageSize.toString())
    if (params?.orderBy && params.orderBy.length > 0) {
      params.orderBy.forEach((o) => searchParams.append('orderBy', o))
    }
    if (params?.status && params.status.length > 0) {
      params.status.forEach((s) => searchParams.append('status', s))
    }
    if (params?.paymentStatus && params.paymentStatus.length > 0) {
      params.paymentStatus.forEach((ps) =>
        searchParams.append('paymentStatus', ps),
      )
    }
    if (params?.socialPlatforms && params.socialPlatforms.length > 0) {
      params.socialPlatforms.forEach((p) =>
        searchParams.append('socialPlatforms', p),
      )
    }
    if (params?.customerId) searchParams.set('customerId', params.customerId)
    if (params?.orderNumber) searchParams.set('orderNumber', params.orderNumber)
    if (params?.from) searchParams.set('from', params.from)
    if (params?.to) searchParams.set('to', params.to)

    const query = searchParams.toString() ? `?${searchParams.toString()}` : ''
    return get<ListOrdersResponse>(
      `v1/businesses/${businessDescriptor}/orders${query}`,
    )
  },

  /**
   * Get order by ID with items and notes
   */
  async getOrder(businessDescriptor: string, orderId: string): Promise<Order> {
    return get<Order>(`v1/businesses/${businessDescriptor}/orders/${orderId}`)
  },

  /**
   * Create a new order
   */
  async createOrder(
    businessDescriptor: string,
    data: CreateOrderRequest,
  ): Promise<Order> {
    return post<Order>(`v1/businesses/${businessDescriptor}/orders`, {
      json: data,
    })
  },

  /**
   * Update existing order
   */
  async updateOrder(
    businessDescriptor: string,
    orderId: string,
    data: UpdateOrderRequest,
  ): Promise<Order> {
    return patch<Order>(
      `v1/businesses/${businessDescriptor}/orders/${orderId}`,
      {
        json: data,
      },
    )
  },

  /**
   * Update order status
   */
  async updateOrderStatus(
    businessDescriptor: string,
    orderId: string,
    data: UpdateOrderStatusRequest,
  ): Promise<Order> {
    return patch<Order>(
      `v1/businesses/${businessDescriptor}/orders/${orderId}/status`,
      {
        json: data,
      },
    )
  },

  /**
   * Update order payment status
   */
  async updateOrderPaymentStatus(
    businessDescriptor: string,
    orderId: string,
    data: UpdateOrderPaymentStatusRequest,
  ): Promise<Order> {
    return patch<Order>(
      `v1/businesses/${businessDescriptor}/orders/${orderId}/payment-status`,
      {
        json: data,
      },
    )
  },

  /**
   * Delete an order
   */
  async deleteOrder(
    businessDescriptor: string,
    orderId: string,
  ): Promise<void> {
    return del(`v1/businesses/${businessDescriptor}/orders/${orderId}`)
  },

  /**
   * Create order note
   */
  async createOrderNote(
    businessDescriptor: string,
    orderId: string,
    data: CreateOrderNoteRequest,
  ): Promise<OrderNote> {
    return post<OrderNote>(
      `v1/businesses/${businessDescriptor}/orders/${orderId}/notes`,
      {
        json: data,
      },
    )
  },

  /**
   * Update order note
   */
  async updateOrderNote(
    businessDescriptor: string,
    orderId: string,
    noteId: string,
    data: UpdateOrderNoteRequest,
  ): Promise<OrderNote> {
    return patch<OrderNote>(
      `v1/businesses/${businessDescriptor}/orders/${orderId}/notes/${noteId}`,
      {
        json: data,
      },
    )
  },

  /**
   * Delete order note
   */
  async deleteOrderNote(
    businessDescriptor: string,
    orderId: string,
    noteId: string,
  ): Promise<void> {
    return del(
      `v1/businesses/${businessDescriptor}/orders/${orderId}/notes/${noteId}`,
    )
  },
}

/**
 * Query Options Factory
 */
export const orderQueries = {
  all: ['orders'] as const,
  lists: () => [...orderQueries.all, 'list'] as const,
  list: (
    businessDescriptor: string,
    params?: {
      search?: string
      page?: number
      pageSize?: number
      orderBy?: Array<string>
      status?: Array<OrderStatus>
      paymentStatus?: Array<OrderPaymentStatus>
      socialPlatforms?: Array<SocialPlatform>
      customerId?: string
      orderNumber?: string
      from?: string
      to?: string
    },
  ) =>
    queryOptions({
      queryKey: [...orderQueries.lists(), businessDescriptor, params] as const,
      queryFn: () => orderApi.listOrders(businessDescriptor, params),
      staleTime: STALE_TIME.FIFTEEN_SECONDS,
      placeholderData: keepPreviousData,
    }),
  details: () => [...orderQueries.all, 'detail'] as const,
  detail: (businessDescriptor: string, orderId: string) =>
    queryOptions({
      queryKey: [
        ...orderQueries.details(),
        businessDescriptor,
        orderId,
      ] as const,
      queryFn: () => orderApi.getOrder(businessDescriptor, orderId),
      staleTime: STALE_TIME.THIRTY_SECONDS,
    }),
}

/**
 * React Query Hooks
 */

export function useOrdersQuery(
  businessDescriptor: string,
  params?: {
    search?: string
    page?: number
    pageSize?: number
    orderBy?: Array<string>
    status?: Array<OrderStatus>
    paymentStatus?: Array<OrderPaymentStatus>
    socialPlatforms?: Array<SocialPlatform>
    customerId?: string
    orderNumber?: string
    from?: string
    to?: string
  },
) {
  return useQuery(orderQueries.list(businessDescriptor, params))
}

export function useOrderQuery(businessDescriptor: string, orderId: string) {
  return useQuery(orderQueries.detail(businessDescriptor, orderId))
}

export function useCreateOrderMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Order, Error, CreateOrderRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateOrderRequest) =>
      orderApi.createOrder(businessDescriptor, data),
    ...options,
  })
}

export function useUpdateOrderMutation(
  businessDescriptor: string,
  orderId: string,
  options?: UseMutationOptions<Order, Error, UpdateOrderRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateOrderRequest) =>
      orderApi.updateOrder(businessDescriptor, orderId, data),
    ...options,
  })
}

export function useUpdateOrderStatusMutation(
  businessDescriptor: string,
  orderId: string,
  options?: UseMutationOptions<Order, Error, UpdateOrderStatusRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateOrderStatusRequest) =>
      orderApi.updateOrderStatus(businessDescriptor, orderId, data),
    ...options,
  })
}

export function useUpdateOrderPaymentStatusMutation(
  businessDescriptor: string,
  orderId: string,
  options?: UseMutationOptions<Order, Error, UpdateOrderPaymentStatusRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateOrderPaymentStatusRequest) =>
      orderApi.updateOrderPaymentStatus(businessDescriptor, orderId, data),
    ...options,
  })
}

export function useDeleteOrderMutation(
  businessDescriptor: string,
  orderId: string,
  options?: UseMutationOptions<void, Error, void>,
) {
  return useMutation({
    mutationFn: () => orderApi.deleteOrder(businessDescriptor, orderId),
    ...options,
  })
}
