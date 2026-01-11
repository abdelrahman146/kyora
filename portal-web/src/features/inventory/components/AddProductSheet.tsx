import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import { ArrowLeft, Layers, Package } from 'lucide-react'

import { SingleVariantProductForm } from './forms/SingleVariantProductForm'
import { MultiVariantProductForm } from './forms/MultiVariantProductForm'
import { BottomSheet } from '@/components/molecules/BottomSheet'

export interface AddProductSheetProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: () => void
}

type VariantMode = 'selection' | 'single' | 'multiple'

export function AddProductSheet({
  isOpen,
  onClose,
  onSuccess,
}: AddProductSheetProps) {
  const { t: tInventory } = useTranslation('inventory')
  const { t: tCommon } = useTranslation('common')
  const [mode, setMode] = useState<VariantMode>('selection')

  const handleClose = () => {
    setMode('selection')
    onClose()
  }

  const handleSuccess = () => {
    setMode('selection')
    onSuccess?.()
    onClose()
  }

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={handleClose}
      title={
        mode === 'selection'
          ? tInventory('add_product')
          : mode === 'single'
            ? tInventory('add_single_variant_product')
            : tInventory('add_multi_variant_product')
      }
      size="lg"
      side="end"
    >
      {mode !== 'selection' && (
        <div className="border-b border-base-300 bg-base-50">
          <button
            type="button"
            onClick={() => setMode('selection')}
            className="flex items-center gap-2 px-4 py-3 w-full text-start hover:bg-base-200/50 transition-colors cursor-pointer"
          >
            <ArrowLeft
              size={18}
              className="rtl:rotate-180 text-base-content/70"
            />
            <span className="text-sm font-medium text-base-content/70">
              {tCommon('back')}
            </span>
          </button>
        </div>
      )}

      {mode === 'selection' && (
        <div className="flex flex-col gap-4 p-4">
          <p className="text-base-content/70 mb-2">
            {tInventory('select_variant_mode_description')}
          </p>

          <button
            type="button"
            onClick={() => setMode('single')}
            className="flex items-start gap-4 p-6 rounded-xl border-2 border-base-300 hover:border-primary hover:bg-primary/5 transition-all text-start cursor-pointer"
          >
            <div className="p-3 rounded-lg bg-primary/10">
              <Package className="w-6 h-6 text-primary" />
            </div>
            <div className="flex-1">
              <h3 className="font-semibold text-lg mb-1">
                {tInventory('single_variant_product')}
              </h3>
              <p className="text-sm text-base-content/70">
                {tInventory('single_variant_description')}
              </p>
            </div>
          </button>

          <button
            type="button"
            onClick={() => setMode('multiple')}
            className="flex items-start gap-4 p-6 rounded-xl border-2 border-base-300 hover:border-primary hover:bg-primary/5 transition-all text-start cursor-pointer"
          >
            <div className="p-3 rounded-lg bg-secondary/10">
              <Layers className="w-6 h-6 text-secondary" />
            </div>
            <div className="flex-1">
              <h3 className="font-semibold text-lg mb-1">
                {tInventory('multi_variant_product')}
              </h3>
              <p className="text-sm text-base-content/70">
                {tInventory('multi_variant_description')}
              </p>
            </div>
          </button>
        </div>
      )}

      {mode === 'single' && (
        <SingleVariantProductForm
          onSuccess={handleSuccess}
          onCancel={handleClose}
        />
      )}

      {mode === 'multiple' && (
        <MultiVariantProductForm
          onSuccess={handleSuccess}
          onCancel={handleClose}
        />
      )}
    </BottomSheet>
  )
}
