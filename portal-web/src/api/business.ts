import apiClient from "./client";
import {
  ListBusinessesResponseSchema,
  BusinessSchema,
  CreateBusinessInputSchema,
  type Business,
  type CreateBusinessInput,
} from "./types/business";
import { z } from "zod";

/**
 * Business API Service
 *
 * Provides methods to interact with business endpoints.
 * All methods validate request/response data using Zod schemas.
 */

export const businessApi = {
  // ==========================================================================
  // List Businesses - GET /v1/businesses
  // ==========================================================================

  /**
   * Fetches all businesses for the authenticated workspace
   * @returns Array of businesses
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async listBusinesses(): Promise<Business[]> {
    const response = await apiClient.get("v1/businesses").json();

    // Backend returns { businesses: [...] }
    const validatedResponse = ListBusinessesResponseSchema.parse(response);

    return validatedResponse.businesses;
  },

  // ==========================================================================
  // Get Business - GET /v1/businesses/{businessDescriptor}
  // ==========================================================================

  /**
   * Fetches a single business by descriptor
   * @param descriptor - Business descriptor (unique identifier)
   * @returns Business details
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async getBusiness(descriptor: string): Promise<Business> {
    const response = await apiClient.get(`v1/businesses/${descriptor}`).json();

    // Backend returns { business: {...} }
    const validated = z.object({ business: BusinessSchema }).parse(response);

    return validated.business;
  },

  // ==========================================================================
  // Create Business - POST /v1/businesses
  // ==========================================================================

  /**
   * Creates a new business in the authenticated workspace
   * @param input - Business creation data
   * @returns Created business
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async createBusiness(input: CreateBusinessInput): Promise<Business> {
    // Validate input
    const validatedInput = CreateBusinessInputSchema.parse(input);

    const response = await apiClient
      .post("v1/businesses", { json: validatedInput })
      .json();

    // Backend returns { business: {...} }
    const validated = z.object({ business: BusinessSchema }).parse(response);

    return validated.business;
  },
};
