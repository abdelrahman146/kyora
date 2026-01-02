/**
 * AddProductSheet Component
 *
 * Placeholder component for adding products.
 * Displays "Coming Soon" message until product creation is implemented.
 */

import { useTranslation } from 'react-i18next'
import { Construction } from 'lucide-react'

import { BottomSheet } from '../molecules/BottomSheet'
import { Button } from '../atoms/Button'

export interface AddProductSheetProps {
  isOpen: boolean
  onClose: () => void
}

export function AddProductSheet({ isOpen, onClose }: AddProductSheetProps) {
  const { t } = useTranslation()

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={onClose}
      title={t('inventory.add_product')}
      size="lg"
      side="end"
      footer={
        <div className="flex gap-2">
          <Button variant="primary" fullWidth onClick={onClose}>
            {t('inventory.coming_soon_button')}
          </Button>
        </div>
      }
    >
      <div className="flex flex-col items-center justify-center min-h-64 gap-6 text-center p-6">
        <Construction size={64} className="text-base-content/40" />
        <div>
          <h3 className="text-2xl font-bold mb-2">
            {t('inventory.coming_soon_title')}
          </h3>
          <p className="text-base-content/70 max-w-md">
            {t('inventory.coming_soon_message')}
          </p>
        </div>
      </div>
    </BottomSheet>
  )
}
