import { get, post } from "./client";
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
  /**
   * Fetches all businesses for the authenticated workspace
   * @returns Array of businesses
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async listBusinesses(): Promise<Business[]> {
    const response = await get<unknown>("v1/businesses");

    // Backend returns { businesses: [...] }
    const validatedResponse = ListBusinessesResponseSchema.parse(response);

    return validatedResponse.businesses;
  },

  /**
   * Fetches a single business by descriptor
   * @param descriptor - Business descriptor (unique identifier)
   * @returns Business details
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async getBusiness(descriptor: string): Promise<Business> {
    const response = await get<unknown>(`v1/businesses/${descriptor}`);

    // Backend returns { business: {...} }
    const validated = z.object({ business: BusinessSchema }).parse(response);

    return validated.business;
  },

  /**
   * Creates a new business in the authenticated workspace
   * @param input - Business creation data
   * @returns Created business
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async createBusiness(input: CreateBusinessInput): Promise<Business> {
    // Validate input
    const validatedInput = CreateBusinessInputSchema.parse(input);

    const response = await post<unknown>("v1/businesses", {
      json: validatedInput,
    });

    // Backend returns { business: {...} }
    const validated = z.object({ business: BusinessSchema }).parse(response);

    return validated.business;
  },
};
