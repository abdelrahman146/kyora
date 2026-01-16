import {
  createFileRoute,
  useParams,
  useRouteContext,
} from '@tanstack/react-router'

import { AssetListPage } from '@/features/accounting/components/AssetListPage'

/**
 * Assets Management Route
 *
 * Manage fixed assets: equipment, vehicles, software, furniture, etc.
 * Tracks business-owned assets and their values.
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/accounting/assets',
)({
  staticData: {
    titleKey: 'accounting:header.assets',
  },
  component: AssetsRoute,
})

function AssetsRoute() {
  const { businessDescriptor } = useParams({
    from: '/business/$businessDescriptor/accounting/assets',
  })
  const { business } = useRouteContext({
    from: '/business/$businessDescriptor',
  })

  return (
    <AssetListPage
      businessDescriptor={businessDescriptor}
      currency={business.currency}
    />
  )
}
