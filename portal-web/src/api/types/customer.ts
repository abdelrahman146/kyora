/**
 * Customer API Types
 *
 * TypeScript types for customer-related API operations.
 * Based on backend swagger.json definitions.
 */

export type CustomerGender = "male" | "female" | "other";

export interface CustomerAddress {
  id: string;
  customerId: string;
  street: string | null;
  city: string;
  state: string;
  zipCode: string | null;
  countryCode: string;
  phoneCode: string;
  phoneNumber: string;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
}

export interface CustomerNote {
  id: string;
  customerId: string;
  content: string;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
}

export interface Customer {
  id: string;
  businessId: string;
  name: string;
  email: string | null;
  phoneCode: string | null;
  phoneNumber: string | null;
  countryCode: string;
  whatsappNumber: string | null;
  gender: CustomerGender;
  joinedAt: string;
  instagramUsername: string | null;
  facebookUsername: string | null;
  tiktokUsername: string | null;
  snapchatUsername: string | null;
  xUsername: string | null;
  addresses: CustomerAddress[];
  notes: CustomerNote[];
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;

  // Computed fields from backend aggregation
  ordersCount?: number;
  totalSpent?: number;
  avatarUrl?: string;
}

export interface CreateCustomerRequest {
  name: string;
  email?: string;
  phoneCode?: string;
  phoneNumber?: string;
  countryCode: string;
  whatsappNumber?: string;
  gender?: CustomerGender;
  joinedAt?: string;
  instagramUsername?: string;
  facebookUsername?: string;
  tiktokUsername?: string;
  snapchatUsername?: string;
  xUsername?: string;
}

export interface UpdateCustomerRequest {
  name?: string;
  email?: string;
  phoneCode?: string;
  phoneNumber?: string;
  countryCode?: string;
  whatsappNumber?: string;
  gender?: CustomerGender;
  joinedAt?: string;
  instagramUsername?: string;
  facebookUsername?: string;
  tiktokUsername?: string;
  snapchatUsername?: string;
  xUsername?: string;
}

export interface ListCustomersParams {
  businessDescriptor: string;
  page?: number;
  pageSize?: number;
  orderBy?: string[];
  search?: string;
}

export interface ListCustomersResponse {
  items: Customer[];
  totalCount: number;
  page: number;
  pageSize: number;
  totalPages: number;
  hasMore: boolean;
}
