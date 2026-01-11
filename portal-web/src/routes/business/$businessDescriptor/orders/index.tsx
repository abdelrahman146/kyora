import { createFileRoute } from '@tanstack/react-router'

import type { Order } from '@/api/order'
import type { SocialPlatform } from '@/api/customer'
import { orderQueries } from '@/api/order'
import { OrdersListPage } from '@/features/orders/components/OrdersListPage'
import { OrdersSearchSchema } from '@/features/orders/schema/ordersSearch'

export const Route = createFileRoute('/business/$businessDescriptor/orders/')({
  validateSearch: (search) => OrdersSearchSchema.parse(search),
  loaderDeps: ({ search }) => search,
  loader: async ({ context, deps: search, params }) => {
    const { queryClient } = context as any

    const orderByArray = search.sortBy
      ? [`${search.sortOrder === 'desc' ? '-' : ''}${search.sortBy}`]
      : ['-orderedAt']

    await queryClient.ensureQueryData(
      orderQueries.list(params.businessDescriptor, {
        search: search.search,
        page: search.page,
        pageSize: search.pageSize,
        orderBy: orderByArray,
        status: search.status as Array<Order['status']>,
        paymentStatus: search.paymentStatus as Array<Order['paymentStatus']>,
        socialPlatforms: search.socialPlatforms as Array<SocialPlatform>,
        customerId: search.customerId,
        from: search.from,
        to: search.to,
      }),
    )
  },
  component: OrdersListPage,
})
