/**
 * Customer Address API Client
 *
 * Handles CRUD operations for customer addresses.
 * Includes TanStack Query hooks for data fetching and mutations.
 */

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { del, get, patch, post } from './client'
import type { CustomerAddress } from './customer'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'
import { showErrorToast, showSuccessToast } from '@/lib/toast'

// Request types
export interface CreateAddressRequest {
  countryCode: string
  state: string
  city: string
  phoneCode: string
  phone: string
  street?: string
  zipCode?: string
}

export interface UpdateAddressRequest {
  street?: string
  city?: string
  state?: string
  countryCode?: string
  phoneCode?: string
  phoneNumber?: string
  zipCode?: string
}

/**
 * Address API Client
 */
export const addressApi = {
  /**
   * List all addresses for a customer
   */
  async listAddresses(
    businessDescriptor: string,
    customerId: string,
  ): Promise<Array<CustomerAddress>> {
    return get<Array<CustomerAddress>>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses`,
    )
  },

  /**
   * Create a new address for a customer
   */
  async createAddress(
    businessDescriptor: string,
    customerId: string,
    data: CreateAddressRequest,
  ): Promise<CustomerAddress> {
    return post<CustomerAddress>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses`,
      { json: data },
    )
  },

  /**
   * Update an existing address
   */
  async updateAddress(
    businessDescriptor: string,
    customerId: string,
    addressId: string,
    data: UpdateAddressRequest,
  ): Promise<CustomerAddress> {
    return patch<CustomerAddress>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses/${addressId}`,
      { json: data },
    )
  },

  /**
   * Delete an address
   */
  async deleteAddress(
    businessDescriptor: string,
    customerId: string,
    addressId: string,
  ): Promise<void> {
    return del<void>(
      `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses/${addressId}`,
    )
  },
}

/**
 * Query Hooks
 */

/**
 * Query to fetch customer addresses
 */
export function useAddressesQuery(
  businessDescriptor: string,
  customerId: string,
) {
  return useQuery({
    queryKey: queryKeys.addresses.list(businessDescriptor, customerId),
    queryFn: () => addressApi.listAddresses(businessDescriptor, customerId),
    staleTime: STALE_TIME.THIRTY_SECONDS,
    enabled: !!businessDescriptor && !!customerId,
  })
}

/**
 * Mutation Hooks
 */

/**
 * Mutation to create a new address
 */
export function useCreateAddressMutation(
  businessDescriptor: string,
  customerId: string,
) {
  const queryClient = useQueryClient()
  const { t } = useTranslation()

  return useMutation({
    mutationFn: (data: CreateAddressRequest) =>
      addressApi.createAddress(businessDescriptor, customerId, data),
    onSuccess: () => {
      // Invalidate addresses list
      void queryClient.invalidateQueries({
        queryKey: queryKeys.addresses.list(businessDescriptor, customerId),
      })
      // Invalidate customer detail to refresh addresses
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
      })
      showSuccessToast(t('customers.address.create_success'))
    },
    onError: () => {
      showErrorToast(t('errors.generic.unexpected'))
    },
  })
}

/**
 * Mutation to update an address
 */
export function useUpdateAddressMutation(
  businessDescriptor: string,
  customerId: string,
  addressId: string,
) {
  const queryClient = useQueryClient()
  const { t } = useTranslation()

  return useMutation({
    mutationFn: (data: UpdateAddressRequest) =>
      addressApi.updateAddress(businessDescriptor, customerId, addressId, data),
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: queryKeys.addresses.list(businessDescriptor, customerId),
      })
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
      })
      showSuccessToast(t('customers.address.update_success'))
    },
    onError: () => {
      showErrorToast(t('errors.generic.unexpected'))
    },
  })
}

/**
 * Mutation to delete an address
 */
export function useDeleteAddressMutation(
  businessDescriptor: string,
  customerId: string,
) {
  const queryClient = useQueryClient()
  const { t } = useTranslation()

  return useMutation({
    mutationFn: (addressId: string) =>
      addressApi.deleteAddress(businessDescriptor, customerId, addressId),
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: queryKeys.addresses.list(businessDescriptor, customerId),
      })
      void queryClient.invalidateQueries({
        queryKey: queryKeys.customers.detail(businessDescriptor, customerId),
      })
      showSuccessToast(t('customers.address.delete_success'))
    },
    onError: () => {
      showErrorToast(t('errors.generic.unexpected'))
    },
  })
}
