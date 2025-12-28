import { z } from "zod";

/**
 * Metadata API types
 */

export const CountryMetadataSchema = z.object({
  name: z.string(),
  nameAr: z.string(),
  code: z.string(),
  iso_code: z.string().optional(),
  flag: z.string().optional(),
  phonePrefix: z.string(),
  currencyCode: z.string(),
  currencyLabel: z.string(),
  currencySymbol: z.string(),
});

export type CountryMetadata = z.infer<typeof CountryMetadataSchema>;

export const ListCountriesResponseSchema = z.object({
  countries: z.array(CountryMetadataSchema),
});

export type ListCountriesResponse = z.infer<typeof ListCountriesResponseSchema>;
