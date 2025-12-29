/**
 * Customer Address API Client
 *
 * Handles CRUD operations for customer addresses
 */

import { get, post, patch, delVoid } from "./client";
import type { CustomerAddress } from "./types/customer";

// Request types
export interface CreateAddressRequest {
  countryCode: string;
  state: string;
  city: string;
  phoneCode: string;
  phone: string;
  street?: string;
  zipCode?: string;
}

export interface UpdateAddressRequest {
  street?: string;
  city?: string;
  state?: string;
  countryCode?: string;
  phoneCode?: string;
  phoneNumber?: string;
  zipCode?: string;
}

/**
 * List all addresses for a customer
 */
export async function listAddresses(
  businessDescriptor: string,
  customerId: string
): Promise<CustomerAddress[]> {
  return get<CustomerAddress[]>(
    `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses`
  );
}

/**
 * Create a new address for a customer
 */
export async function createAddress(
  businessDescriptor: string,
  customerId: string,
  data: CreateAddressRequest
): Promise<CustomerAddress> {
  return post<CustomerAddress>(
    `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses`,
    { json: data }
  );
}

/**
 * Update an existing address
 */
export async function updateAddress(
  businessDescriptor: string,
  customerId: string,
  addressId: string,
  data: UpdateAddressRequest
): Promise<CustomerAddress> {
  return patch<CustomerAddress>(
    `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses/${addressId}`,
    { json: data }
  );
}

/**
 * Delete an address
 */
export async function deleteAddress(
  businessDescriptor: string,
  customerId: string,
  addressId: string
): Promise<void> {
  return delVoid(
    `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses/${addressId}`
  );
}
