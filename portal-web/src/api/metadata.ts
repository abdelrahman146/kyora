import { apiClient } from "./client";
import {
  ListCountriesResponseSchema,
  type ListCountriesResponse,
} from "./types/metadata";

export const metadataApi = {
  async listCountries(): Promise<ListCountriesResponse> {
    const data = await apiClient.get("v1/metadata/countries").json();
    return ListCountriesResponseSchema.parse(data);
  },
};
