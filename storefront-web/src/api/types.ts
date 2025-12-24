export type ISO3166Alpha2 = string;

export interface Problem {
  type?: string;
  title?: string;
  status?: number;
  detail?: string;
  instance?: string;
  [key: string]: unknown;
}

export interface StorefrontTheme {
  primaryColor?: string;
  secondaryColor?: string;
  accentColor?: string;
  backgroundColor?: string;
  textColor?: string;
  fontFamily?: string;
  headingFontFamily?: string;
}

export interface PublicBusiness {
  id: string;
  name: string;
  descriptor: string;
  brand?: string;
  logoUrl?: string;
  countryCode: ISO3166Alpha2;
  currency: string;
  storefrontPublicId: string;
  storefrontEnabled: boolean;
  storefrontTheme: StorefrontTheme;
  supportEmail?: string;
  phoneNumber?: string;
  whatsappNumber?: string;
  address?: string;
  websiteUrl?: string;
  instagramUrl?: string;
  facebookUrl?: string;
  tiktokUrl?: string;
  xUrl?: string;
  snapchatUrl?: string;
}

export interface PublicCategory {
  id: string;
  name: string;
  descriptor: string;
}

export interface PublicVariant {
  id: string;
  productId: string;
  code: string;
  name: string;
  sku: string;
  salePrice: string; // decimal string
  currency: string;
  photos?: string[];
}

export interface PublicProduct {
  id: string;
  name: string;
  description?: string;
  categoryId: string;
  photos?: string[];
  variants: PublicVariant[];
}

export interface CatalogResponse {
  business: PublicBusiness;
  categories: PublicCategory[];
  products: PublicProduct[];
}

export interface PublicShippingZone {
  id: string;
  name: string;
  countries: ISO3166Alpha2[];
  currency: string;
  shippingCost: string; // decimal string
  freeShippingThreshold: string; // decimal string
}

export interface CreateOrderItem {
  variantId: string;
  quantity: number;
  specialRequest?: string;
}

export interface CreateOrderCustomer {
  email: string;
  name: string;
  phoneNumber?: string;
  instagramUsername?: string;
}

export interface CreateOrderShippingAddress {
  countryCode: ISO3166Alpha2;
  state: string;
  city: string;
  street?: string;
  zipCode?: string;
  phoneCode: string;
  phoneNumber: string;
}

export interface CreateOrderRequest {
  customer: CreateOrderCustomer;
  shippingAddress: CreateOrderShippingAddress;
  items: CreateOrderItem[];
}

export interface CreateOrderResponse {
  orderId: string;
  orderNumber: string;
  status: string;
  paymentStatus: string;
  total: string; // decimal string
  currency: string;
}
