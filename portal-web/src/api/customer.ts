/**
 * Customer API Client
 *
 * API client for customer-related operations.
 * Provides CRUD operations and search functionality for customers.
 */

import { apiClient } from "./client";
import type {
  Customer,
  CreateCustomerRequest,
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

  return apiClient
    .get(
      `v1/businesses/${businessDescriptor}/customers?${searchParams.toString()}`
    )
    .json<ListCustomersResponse>();
}

/**
 * Fetch a single customer by ID
 */
export async function getCustomer(
  businessDescriptor: string,
  customerId: string
): Promise<Customer> {
  return apiClient
    .get(`v1/businesses/${businessDescriptor}/customers/${customerId}`)
    .json<Customer>();
}

/**
 * Create a new customer
 */
export async function createCustomer(
  businessDescriptor: string,
  data: CreateCustomerRequest
): Promise<Customer> {
  return apiClient
    .post(`v1/businesses/${businessDescriptor}/customers`, {
      json: data,
    })
    .json<Customer>();
}

/**
 * Update an existing customer
 */
export async function updateCustomer(
  businessDescriptor: string,
  customerId: string,
  data: UpdateCustomerRequest
): Promise<Customer> {
  return apiClient
    .put(`v1/businesses/${businessDescriptor}/customers/${customerId}`, {
      json: data,
    })
    .json<Customer>();
}

/**
 * Delete a customer (soft delete)
 */
export async function deleteCustomer(
  businessDescriptor: string,
  customerId: string
): Promise<void> {
  await apiClient.delete(
    `v1/businesses/${businessDescriptor}/customers/${customerId}`
  );
}
