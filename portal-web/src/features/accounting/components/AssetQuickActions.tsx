/**
 * AssetQuickActions Component
 *
 * Dropdown menu for asset actions.
 * Provides edit and delete operations.
 */

import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'react-hot-toast'
import { MoreVertical, Pencil, Trash2 } from 'lucide-react'

import { EditAssetSheet } from './sheets/EditAssetSheet'
import type { Asset } from '@/api/accounting'
import { useDeleteAssetMutation } from '@/api/accounting'
import { Button } from '@/components/atoms/Button'
import { ConfirmDialog } from '@/components/molecules/ConfirmDialog'
import { translateErrorAsync } from '@/lib/translateError'

interface AssetQuickActionsProps {
  asset: Asset
  businessDescriptor: string
  onActionComplete?: () => void
}

export function AssetQuickActions({
  asset,
  businessDescriptor,
  onActionComplete,
}: AssetQuickActionsProps) {
  const { t } = useTranslation('accounting')
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = useState(false)

  const deleteAssetMutation = useDeleteAssetMutation(
    businessDescriptor,
    asset.id,
  )

  const handleDelete = async () => {
    try {
      await deleteAssetMutation.mutateAsync()
      toast.success(t('toast.asset_deleted'))
      await onActionComplete?.()
      setIsDeleteOpen(false)
    } catch (error) {
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    }
  }

  return (
    <>
      {/* Dropdown Menu */}
      <div className="dropdown dropdown-end">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="btn-square"
          aria-label="More actions"
          onClick={() => setIsMenuOpen(!isMenuOpen)}
        >
          <MoreVertical className="h-4 w-4" />
        </Button>

        {isMenuOpen && (
          <ul className="dropdown-content z-10 menu p-2 shadow bg-base-100 border border-base-300 rounded-box w-52 mt-2">
            <li>
              <button
                type="button"
                onClick={() => {
                  setIsMenuOpen(false)
                  setIsEditOpen(true)
                }}
                className="flex items-center gap-2"
              >
                <Pencil className="h-4 w-4" />
                {t('actions.edit_asset')}
              </button>
            </li>
            <li>
              <button
                type="button"
                onClick={() => {
                  setIsMenuOpen(false)
                  setIsDeleteOpen(true)
                }}
                className="flex items-center gap-2 text-error"
              >
                <Trash2 className="h-4 w-4" />
                {t('actions.delete_asset')}
              </button>
            </li>
          </ul>
        )}
      </div>

      {/* Edit Sheet */}
      <EditAssetSheet
        isOpen={isEditOpen}
        onClose={() => setIsEditOpen(false)}
        asset={asset}
        businessDescriptor={businessDescriptor}
        onUpdated={onActionComplete}
      />

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={isDeleteOpen}
        onClose={() => setIsDeleteOpen(false)}
        onConfirm={handleDelete}
        title={t('delete.asset_title')}
        message={t('delete.asset_message')}
        confirmText={t('common:actions.delete')}
        cancelText={t('common:actions.cancel')}
        variant="error"
      />
    </>
  )
}
