import { Languages } from 'lucide-react'
import { useLanguage } from '@/hooks/useLanguage'
import { cn } from '@/lib/utils'

export interface LanguageSwitcherProps {
  variant?: 'full' | 'icon'
  className?: string
}

/**
 * LanguageSwitcher Component
 *
 * Toggles between Arabic and English with cookie persistence.
 */
export function LanguageSwitcher({
  variant = 'full',
  className,
}: LanguageSwitcherProps) {
  const { language, toggleLanguage } = useLanguage()

  return (
    <button
      onClick={toggleLanguage}
      className={cn(
        'btn btn-ghost',
        variant === 'icon' ? 'btn-circle btn-sm' : 'btn-sm',
        className,
      )}
      aria-label={`Switch to ${language === 'ar' ? 'English' : 'Arabic'}`}
    >
      {variant === 'full' ? (
        <>
          <Languages size={18} />
          <span>{language === 'ar' ? 'English' : 'العربية'}</span>
        </>
      ) : (
        <Languages size={18} />
      )}
    </button>
  )
}
