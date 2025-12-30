/**
 * SubmitButton Component - Form Submission Button
 *
 * Smart button that subscribes to form state and automatically disables
 * when form is invalid or submitting. Uses form.Subscribe to prevent
 * component re-renders (better performance than useStore).
 *
 * Usage:
 * ```tsx
 * <form.SubmitButton>Login</form.SubmitButton>
 * <form.SubmitButton loadingText="Logging in...">Login</form.SubmitButton>
 * ```
 */

import { Loader2 } from 'lucide-react'
import { Button } from '@/components/atoms/Button'
import { useFormContext } from '../contexts'
import type { ButtonHTMLAttributes } from 'react'

interface SubmitButtonProps extends Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'type'> {
  children: React.ReactNode
  /** Text to show while form is submitting */
  loadingText?: string
  /** Variant for button styling */
  variant?: 'primary' | 'secondary' | 'accent' | 'ghost' | 'link'
  /** Size variant */
  size?: 'xs' | 'sm' | 'md' | 'lg'
}

export function SubmitButton({
  children,
  loadingText,
  variant = 'primary',
  size = 'md',
  disabled,
  ...props
}: SubmitButtonProps) {
  const form = useFormContext()

  return (
    <form.Subscribe
      selector={(state) => ({
        canSubmit: state.canSubmit,
        isSubmitting: state.isSubmitting,
      })}
    >
      {({ canSubmit, isSubmitting }) => (
        <Button
          type="submit"
          variant={variant}
          size={size}
          disabled={!canSubmit || isSubmitting || disabled}
          {...props}
        >
          {isSubmitting && <Loader2 className="animate-spin" size={16} />}
          {isSubmitting && loadingText ? loadingText : children}
        </Button>
      )}
    </form.Subscribe>
  )
}
