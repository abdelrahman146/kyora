import { createFileRoute } from '@tanstack/react-router'
import { Suspense } from 'react'

import type {InventorySearch} from '@/features/inventory/schema/inventorySearch';
import {
  InventoryListPage,
  inventoryListLoader,
} from '@/features/inventory/components/InventoryListPage'
import { InventoryListSkeleton } from '@/features/inventory/components/InventoryListSkeleton'
import {
  
  InventorySearchSchema
} from '@/features/inventory/schema/inventorySearch'

export const Route = createFileRoute(
  '/business/$businessDescriptor/inventory/',
)({
  staticData: {
    titleKey: 'inventory.title',
  },
  validateSearch: (search): InventorySearch => {
    return InventorySearchSchema.parse(search)
  },
  loader: async ({ context, params, location }) => {
    const { queryClient } = context as any

    const searchParams = InventorySearchSchema.parse(location.search)
    await inventoryListLoader({
      queryClient,
      businessDescriptor: params.businessDescriptor,
      search: searchParams,
    })
  },
  component: () => (
    <Suspense fallback={<InventoryListSkeleton />}>
      <InventoryListPage />
    </Suspense>
  ),
})
