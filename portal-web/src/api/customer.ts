/**
 * Customer API Client
 *
 * API client for customer-related operations.
 * Provides CRUD operations and search functionality for customers.
 */

import { get, post, put, delVoid } from "./client";
import type {
  Customer,
  CustomerAddress,
  CreateCustomerRequest,
  CreateCustomerAddressRequest,
  UpdateCustomerRequest,
  ListCustomersParams,
  ListCustomersResponse,
} from "./types/customer";

/**
 * Fetch paginated list of customers
 */
export async function listCustomers(
  params: ListCustomersParams
): Promise<ListCustomersResponse> {
  const {
    businessDescriptor,
    page = 1,
    pageSize = 20,
    orderBy,
    search,
  } = params;

  const searchParams = new URLSearchParams({
    page: page.toString(),
    pageSize: pageSize.toString(),
  });

  if (orderBy && orderBy.length > 0) {
    searchParams.append("orderBy", orderBy.join(","));
  }

  if (search) {
    searchParams.append("search", search);
  }

  return get<ListCustomersResponse>(
    `v1/businesses/${businessDescriptor}/customers?${searchParams.toString()}`
  );
}

/**
 * Fetch a single customer by ID
 */
export async function getCustomer(
  businessDescriptor: string,
  customerId: string
): Promise<Customer> {
  return get<Customer>(
    `v1/businesses/${businessDescriptor}/customers/${customerId}`
  );
}

/**
 * Create a new customer
 */
export async function createCustomer(
  businessDescriptor: string,
  data: CreateCustomerRequest
): Promise<Customer> {
  return post<Customer>(`v1/businesses/${businessDescriptor}/customers`, {
    json: data,
  });
}

/**
 * Create a new address for a customer
 */
export async function createCustomerAddress(
  businessDescriptor: string,
  customerId: string,
  data: CreateCustomerAddressRequest
): Promise<CustomerAddress> {
  return post<CustomerAddress>(
    `v1/businesses/${businessDescriptor}/customers/${customerId}/addresses`,
    {
      json: data,
    }
  );
}

/**
 * Update an existing customer
 */
export async function updateCustomer(
  businessDescriptor: string,
  customerId: string,
  data: UpdateCustomerRequest
): Promise<Customer> {
  return put<Customer>(
    `v1/businesses/${businessDescriptor}/customers/${customerId}`,
    {
      json: data,
    }
  );
}

/**
 * Delete a customer (soft delete)
 */
export async function deleteCustomer(
  businessDescriptor: string,
  customerId: string
): Promise<void> {
  await delVoid(`v1/businesses/${businessDescriptor}/customers/${customerId}`);
}
