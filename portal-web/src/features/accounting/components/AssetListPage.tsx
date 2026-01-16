/**
 * AssetListPage Component
 *
 * List view for fixed assets.
 * Features pagination, empty states with helpful CTAs.
 *
 * Layout:
 * - Header with Add Asset action
 * - List of asset cards
 * - Pagination
 * - Empty state with CTA
 *
 * Mobile-first with card-based layout.
 */

import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Package, Plus } from 'lucide-react'

import { AssetCard } from './AssetCard'
import { CreateAssetSheet } from './sheets/CreateAssetSheet'
import type { Asset } from '@/api/accounting'
import { useAssetsQuery } from '@/api/accounting'
import { Button } from '@/components/atoms/Button'
import { Pagination } from '@/components/molecules/Pagination'

interface AssetListPageProps {
  businessDescriptor: string
  currency: string
}

export function AssetListPage({
  businessDescriptor,
  currency,
}: AssetListPageProps) {
  const { t } = useTranslation('accounting')
  const [page, setPage] = useState(1)
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  const pageSize = 20

  // Fetch assets
  const {
    data: assetsData,
    isLoading,
    refetch: refetchAssets,
  } = useAssetsQuery(businessDescriptor, {
    page,
    pageSize,
    orderBy: ['-purchasedAt'],
  })

  const assets = assetsData?.items ?? []
  const isEmpty = assets.length === 0 && !isLoading

  const handleActionComplete = async () => {
    await refetchAssets()
  }

  return (
    <div className="space-y-4">
      {/* Header with Action Button */}
      <div className="flex items-center justify-between gap-4">
        <h1 className="text-2xl font-bold text-base-content">
          {t('header.assets')}
        </h1>

        {/* Add Asset Button (Desktop) */}
        <div className="hidden sm:block">
          <Button
            variant="primary"
            size="md"
            onClick={() => setIsCreateOpen(true)}
            className="gap-2"
          >
            <Plus className="h-5 w-5" />
            {t('actions.add_asset')}
          </Button>
        </div>
      </div>

      {/* List Content */}
      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div
              key={i}
              className="skeleton h-24 w-full rounded-lg bg-base-200"
            />
          ))}
        </div>
      ) : isEmpty ? (
        <div className="card bg-base-100 p-8 text-center">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-base-200">
            <Package className="h-8 w-8 text-base-content/60" />
          </div>
          <h3 className="mb-2 text-lg font-semibold text-base-content">
            {t('empty.assets_title')}
          </h3>
          <p className="mb-6 text-sm text-base-content/70">
            {t('empty.assets_description')}
          </p>
          <Button
            variant="primary"
            size="md"
            onClick={() => setIsCreateOpen(true)}
            className="gap-2"
          >
            <Plus className="h-5 w-5" />
            {t('actions.add_first_asset')}
          </Button>
        </div>
      ) : (
        <>
          {/* Assets List */}
          <div className="space-y-3">
            {assets.map((asset: Asset) => (
              <AssetCard
                key={asset.id}
                asset={asset}
                currency={currency}
                businessDescriptor={businessDescriptor}
                onActionComplete={handleActionComplete}
              />
            ))}
          </div>

          {/* Pagination */}
          {assetsData && (
            <Pagination
              currentPage={page}
              totalPages={Math.ceil(
                // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
                (assetsData.totalCount ?? 0) / pageSize,
              )}
              pageSize={pageSize}
              // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
              totalItems={assetsData.totalCount ?? 0}
              itemsName={t('items.assets')}
              onPageChange={setPage}
            />
          )}
        </>
      )}

      {/* Mobile FAB */}
      <div className="sm:hidden fixed bottom-6 end-6 z-10">
        <Button
          variant="primary"
          size="lg"
          className="btn-circle shadow-lg"
          onClick={() => setIsCreateOpen(true)}
          aria-label={t('actions.add_asset')}
        >
          <Plus className="h-6 w-6" />
        </Button>
      </div>

      {/* Create Asset Sheet */}
      <CreateAssetSheet
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        businessDescriptor={businessDescriptor}
        onCreated={handleActionComplete}
      />
    </div>
  )
}
