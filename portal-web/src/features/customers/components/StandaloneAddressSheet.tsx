/**
 * StandaloneAddressSheet Component
 *
 * A wrapper around AddressSheet that handles mutations internally.
 * Used for creating addresses in contexts like order creation where
 * we need to handle the mutation without a parent form managing it.
 *
 * This reuses AddressSheet (callback-based) to maintain DRY principle.
 */

import { useMemo } from 'react'

import { AddressSheet } from './AddressSheet'
import type { CreateAddressRequest, UpdateAddressRequest } from '@/api/address'
import type { CustomerAddress } from '@/api/customer'
import { useCreateAddressMutation } from '@/api/address'
import { useShippingZonesQuery } from '@/api/business'
import { useCountriesQuery } from '@/api/metadata'

export interface StandaloneAddressSheetProps {
  isOpen: boolean
  onClose: () => void
  businessDescriptor: string
  customerId: string
  businessCountryCode: string
  onCreated: (address: CustomerAddress) => void
}

export function StandaloneAddressSheet({
  isOpen,
  onClose,
  businessDescriptor,
  customerId,
  businessCountryCode,
  onCreated,
}: StandaloneAddressSheetProps) {
  const createMutation = useCreateAddressMutation(
    businessDescriptor,
    customerId,
  )

  const { data: countries = [], isLoading: countriesLoading } =
    useCountriesQuery()
  const { data: shippingZones = [], isLoading: zonesLoading } =
    useShippingZonesQuery(businessDescriptor, isOpen)
  const countriesReady = countries.length > 0 && !countriesLoading
  const zonesReady = shippingZones.length > 0 && !zonesLoading

  // Find default shipping zone for business country
  const defaultZone = useMemo(() => {
    if (!zonesReady || !businessCountryCode) return null
    return shippingZones.find((zone) =>
      zone.countries.includes(businessCountryCode),
    )
  }, [shippingZones, zonesReady, businessCountryCode])

  // Build a "default" address to pre-populate the form
  const defaultAddress = useMemo<CustomerAddress | undefined>(() => {
    if (!zonesReady || !countriesReady) return undefined

    const country = countries.find((c) => c.code === businessCountryCode)

    return {
      id: '',
      customerId: '',
      shippingZoneId: defaultZone?.id ?? '',
      countryCode: businessCountryCode,
      state: '',
      city: '',
      phoneCode: country?.phonePrefix ?? '',
      phoneNumber: '',
      street: '',
      zipCode: '',
      createdAt: '',
      updatedAt: '',
    }
  }, [zonesReady, countriesReady, defaultZone, businessCountryCode, countries])

  const handleSubmit = async (
    data: CreateAddressRequest | UpdateAddressRequest,
  ) => {
    const requestData = data as CreateAddressRequest

    const created = await createMutation.mutateAsync(requestData)
    onCreated(created)
    return created
  }

  return (
    <AddressSheet
      isOpen={isOpen}
      onClose={onClose}
      businessDescriptor={businessDescriptor}
      onSubmit={handleSubmit}
      address={defaultAddress}
      isSubmitting={createMutation.isPending}
    />
  )
}
