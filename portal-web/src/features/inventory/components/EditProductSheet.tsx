import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import { ArrowLeft, Layers, Package, Trash2 } from 'lucide-react'

import { UpdateProductForm } from './forms/UpdateProductForm'
import { UpdateProductWithVariantsForm } from './forms/UpdateProductWithVariantsForm'
import { BottomSheet } from '@/components/molecules/BottomSheet'
import { ConfirmDialog } from '@/components/molecules/ConfirmDialog'
import { useProductQuery } from '@/api/inventory'

export interface EditProductSheetProps {
  isOpen: boolean
  onClose: () => void
  productId: string | null
  businessDescriptor: string
  onSuccess?: () => void
  onDelete?: () => void
}

export function EditProductSheet({
  isOpen,
  onClose,
  productId,
  businessDescriptor,
  onSuccess,
  onDelete,
}: EditProductSheetProps) {
  const { t } = useTranslation()
  const [showVariantsForm, setShowVariantsForm] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)

  const { data: product, isLoading } = useProductQuery(
    businessDescriptor,
    productId || '',
  )

  const handleClose = () => {
    setShowVariantsForm(false)
    onClose()
  }

  const handleSuccess = () => {
    setShowVariantsForm(false)
    onSuccess?.()
    onClose()
  }

  const handleAddVariants = () => {
    setShowVariantsForm(true)
  }

  const handleDelete = () => {
    setShowDeleteConfirm(true)
  }

  const handleConfirmDelete = () => {
    setShowDeleteConfirm(false)
    onDelete?.()
    onClose()
  }

  if (!productId) return null

  const hasSingleVariant = product?.variants && product.variants.length === 1
  const hasMultipleVariants = product?.variants && product.variants.length > 1

  return (
    <>
      <BottomSheet
        isOpen={isOpen}
        onClose={handleClose}
        title={t('edit_product', { ns: 'inventory' })}
        size="lg"
        side="end"
      >
        {showVariantsForm && (
          <div className="border-b border-base-300 bg-base-50">
            <button
              type="button"
              onClick={() => setShowVariantsForm(false)}
              className="flex items-center gap-2 px-4 py-3 w-full text-start hover:bg-base-200/50 transition-colors"
            >
              <ArrowLeft
                size={18}
                className="rtl:rotate-180 text-base-content/70"
              />
              <span className="text-sm font-medium text-base-content/70">
                {t('common.back')}
              </span>
            </button>
          </div>
        )}

        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        ) : !product ? (
          <div className="flex flex-col items-center justify-center py-12 text-center">
            <Package size={64} className="text-base-content/40 mb-4" />
            <p className="text-base-content/70">
              {t('product_not_found', { ns: 'inventory' })}
            </p>
          </div>
        ) : (
          <>
            {!showVariantsForm && hasSingleVariant && (
              <div className="space-y-4">
                <UpdateProductForm
                  product={product}
                  businessDescriptor={businessDescriptor}
                  onSuccess={handleSuccess}
                  onCancel={handleClose}
                />
                <div className="px-4 pb-4">
                  <button
                    type="button"
                    onClick={handleAddVariants}
                    className="btn btn-outline btn-sm w-full"
                  >
                    <Layers className="w-4 h-4" />
                    {t('add_more_variants', { ns: 'inventory' })}
                  </button>
                  <button
                    type="button"
                    onClick={handleDelete}
                    className="btn btn-ghost btn-sm w-full text-error mt-2"
                  >
                    <Trash2 className="w-4 h-4" />
                    {t('delete_product', { ns: 'inventory' })}
                  </button>
                </div>
              </div>
            )}

            {(showVariantsForm || hasMultipleVariants) && (
              <div className="space-y-4">
                <UpdateProductWithVariantsForm
                  product={product}
                  businessDescriptor={businessDescriptor}
                  onSuccess={handleSuccess}
                  onCancel={handleClose}
                />
                <div className="px-4 pb-4">
                  <button
                    type="button"
                    onClick={handleDelete}
                    className="btn btn-ghost btn-sm w-full text-error"
                  >
                    <Trash2 className="w-4 h-4" />
                    {t('delete_product', { ns: 'inventory' })}
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </BottomSheet>

      <ConfirmDialog
        isOpen={showDeleteConfirm}
        onClose={() => setShowDeleteConfirm(false)}
        onConfirm={handleConfirmDelete}
        title={t('delete_confirm_title', { ns: 'inventory' })}
        message={t('delete_confirm_message', {
          ns: 'inventory',
          name: product?.name || '',
        })}
        confirmText={t('common.delete')}
        cancelText={t('common.cancel')}
        variant="error"
      />
    </>
  )
}
