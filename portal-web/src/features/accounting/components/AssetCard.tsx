/**
 * AssetCard Component
 *
 * Mobile-first card for displaying fixed assets.
 * Shows asset name, type, value, purchase date, and optional note.
 *
 * Features:
 * - Type-based icon (equipment, software, vehicle, furniture, other)
 * - Value formatting with currency
 * - Date with localized formatting
 * - Quick actions dropdown for edit/delete
 * - RTL-compatible with logical properties
 * - Mobile-optimized touch targets
 */

import { useTranslation } from 'react-i18next'
import { Car, Laptop, Monitor, Package, Sofa } from 'lucide-react'

import { AssetQuickActions } from './AssetQuickActions'
import type { Asset, AssetType } from '@/api/accounting'
import { formatCurrency } from '@/lib/formatCurrency'
import { formatDateShort } from '@/lib/formatDate'

export interface AssetCardProps {
  asset: Asset
  currency: string
  businessDescriptor: string
  onActionComplete?: () => void
}

/**
 * Get the appropriate icon for an asset type
 */
function getAssetTypeIcon(type: AssetType) {
  switch (type) {
    case 'software':
      return Monitor
    case 'equipment':
      return Laptop
    case 'vehicle':
      return Car
    case 'furniture':
      return Sofa
    case 'other':
    default:
      return Package
  }
}

/**
 * Get the appropriate color class for an asset type
 */
function getAssetTypeColorClass(type: AssetType) {
  switch (type) {
    case 'software':
      return 'bg-info/10 text-info'
    case 'equipment':
      return 'bg-primary/10 text-primary'
    case 'vehicle':
      return 'bg-warning/10 text-warning'
    case 'furniture':
      return 'bg-secondary/10 text-secondary'
    case 'other':
    default:
      return 'bg-base-content/10 text-base-content'
  }
}

export function AssetCard({
  asset,
  currency,
  businessDescriptor,
  onActionComplete,
}: AssetCardProps) {
  const { t } = useTranslation('accounting')

  const value = parseFloat(asset.value)
  const Icon = getAssetTypeIcon(asset.type)
  const colorClass = getAssetTypeColorClass(asset.type)

  return (
    <div className="card bg-base-100 border border-base-300">
      <div className="card-body p-4 flex flex-row items-start gap-3">
        {/* Icon */}
        <div className={`rounded-full p-2 ${colorClass} shrink-0`}>
          <Icon className="h-5 w-5" />
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0 space-y-1">
          {/* Name & Value */}
          <div className="flex items-center justify-between gap-2">
            <h3 className="font-medium text-base-content truncate">
              {asset.name}
            </h3>
            <span className="font-semibold text-base-content text-sm shrink-0">
              {formatCurrency(value, currency)}
            </span>
          </div>

          {/* Type & Date */}
          <div className="flex items-center gap-2 flex-wrap">
            <span className="badge badge-sm badge-ghost">
              {t(`asset_type.${asset.type}`)}
            </span>
            <span className="text-xs text-base-content/60">
              {t('list.purchased_on', {
                date: formatDateShort(asset.purchasedAt),
              })}
            </span>
          </div>

          {/* Note (if present) */}
          {asset.note && (
            <p className="text-sm text-base-content/70 line-clamp-2">
              {asset.note}
            </p>
          )}
        </div>

        {/* Quick Actions Menu */}
        <AssetQuickActions
          asset={asset}
          businessDescriptor={businessDescriptor}
          onActionComplete={onActionComplete}
        />
      </div>
    </div>
  )
}
