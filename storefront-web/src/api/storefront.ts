import { apiFetch } from "./client";
import type {
  CatalogResponse,
  CreateOrderRequest,
  CreateOrderResponse,
  PublicShippingZone,
} from "./types";

export const storefrontApi = {
  getCatalog: (storefrontPublicId: string, signal?: AbortSignal) =>
    apiFetch<CatalogResponse>(
      `/v1/storefront/${encodeURIComponent(storefrontPublicId)}/catalog`,
      {
        method: "GET",
        signal,
      }
    ),

  listShippingZones: (storefrontPublicId: string, signal?: AbortSignal) =>
    apiFetch<PublicShippingZone[]>(
      `/v1/storefront/${encodeURIComponent(storefrontPublicId)}/shipping-zones`,
      { method: "GET", signal }
    ),

  createOrder: (
    storefrontPublicId: string,
    req: CreateOrderRequest,
    idempotencyKey: string,
    signal?: AbortSignal
  ) =>
    apiFetch<CreateOrderResponse>(
      `/v1/storefront/${encodeURIComponent(storefrontPublicId)}/orders`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Idempotency-Key": idempotencyKey,
        },
        body: JSON.stringify(req),
        signal,
      }
    ),
};
