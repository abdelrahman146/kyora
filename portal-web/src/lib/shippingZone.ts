/**
 * Shipping Zone Utilities
 *
 * Production-ready utilities for inferring and managing shipping zones
 * based on customer addresses.
 */

import type { ShippingZone } from '@/api/business'
import type { CustomerAddress } from '@/api/customer'

/**
 * Infers the appropriate shipping zone for a given address based on country code.
 *
 * @param address - The customer address to match
 * @param zones - Available shipping zones
 * @returns The matched shipping zone, or undefined if no match found
 *
 * @example
 * ```ts
 * const zone = inferShippingZoneFromAddress(address, allZones)
 * if (zone) {
 *   console.log(`Matched zone: ${zone.name}`)
 * }
 * ```
 */
export function inferShippingZoneFromAddress(
  address:
    | CustomerAddress
    | Pick<CustomerAddress, 'countryCode'>
    | null
    | undefined,
  zones: Array<ShippingZone>,
): ShippingZone | undefined {
  if (!address?.countryCode || zones.length === 0) {
    return undefined
  }

  // Normalize country code to uppercase for comparison
  const countryCode = address.countryCode.toUpperCase()

  // Find the first zone that includes this country
  return zones.find((zone) =>
    zone.countries.some(
      (zoneCountry) => zoneCountry.toUpperCase() === countryCode,
    ),
  )
}

/**
 * Checks if a country code is valid for a given shipping zone.
 *
 * @param countryCode - The ISO 3166-1 alpha-2 country code
 * @param zone - The shipping zone to check against
 * @returns true if the country is in the zone, false otherwise
 */
export function isCountryInZone(
  countryCode: string,
  zone: ShippingZone | null | undefined,
): boolean {
  if (!countryCode || !zone) {
    return false
  }

  const normalizedCode = countryCode.toUpperCase()
  return zone.countries.some(
    (zoneCountry) => zoneCountry.toUpperCase() === normalizedCode,
  )
}

/**
 * Formats shipping zone information for display.
 *
 * @param zone - The shipping zone
 * @param currency - Optional currency override
 * @returns Formatted display object
 */
export function formatShippingZoneInfo(
  zone: ShippingZone,
  currency?: string,
): {
  name: string
  cost: string
  freeThreshold: string
  countryCount: number
} {
  return {
    name: zone.name,
    cost: `${parseFloat(zone.shippingCost).toFixed(2)} ${currency || zone.currency}`,
    freeThreshold: `${parseFloat(zone.freeShippingThreshold).toFixed(2)} ${currency || zone.currency}`,
    countryCount: zone.countries.length,
  }
}

/**
 * Gets available countries for a shipping zone.
 *
 * @param zone - The shipping zone
 * @returns Array of uppercase country codes
 */
export function getZoneCountries(
  zone: ShippingZone | null | undefined,
): string[] {
  if (!zone) {
    return []
  }
  return zone.countries.map((c) => c.toUpperCase())
}
