/**
 * Form System Type Definitions
 *
 * Core types for the Kyora form composition layer built on TanStack Form.
 * Provides type-safe field props, validation error shapes, and form values.
 */

import type { FieldApi, FormApi } from '@tanstack/react-form'

/**
 * Supported validation error types
 *
 * - string: Translation keys (most common) - e.g., "validation.required"
 * - object: Structured errors with message/code - e.g., { message: "...", code: 1001 }
 * - undefined: No error
 */
export type ValidationError =
  | string
  | { message: string; code?: number }
  | undefined

/**
 * Base form field props shared across all field components
 */
export interface BaseFieldProps {
  /** Field label displayed above input */
  label?: string
  /** Placeholder text for empty state */
  placeholder?: string
  /** Help text displayed below input */
  hint?: string
  /** Whether the field is required (visual indicator only) */
  required?: boolean
  /** Whether the field is disabled */
  disabled?: boolean
  /** Additional CSS classes */
  className?: string
  /** Auto-focus the field on mount */
  autoFocus?: boolean
}

/**
 * Text input specific props
 */
export interface TextFieldProps extends BaseFieldProps {
  /** Input type */
  type?: 'text' | 'email' | 'url' | 'tel' | 'search'
  /** Icon component to display at start of input */
  startIcon?: React.ReactNode
  /** Icon component to display at end of input */
  endIcon?: React.ReactNode
  /** Autocomplete attribute for browser autofill */
  autoComplete?: string
  /** Maximum character length */
  maxLength?: number
}

/**
 * Password input specific props
 */
export interface PasswordFieldProps extends BaseFieldProps {
  /** Autocomplete attribute for browser autofill */
  autoComplete?: string
  /** Show password strength indicator */
  showStrength?: boolean
}

/**
 * Select dropdown specific props
 */
export interface SelectFieldProps<T = string> extends BaseFieldProps {
  /** Options for the select dropdown */
  options: Array<{ value: T; label: string; disabled?: boolean }>
  /** Whether to allow search/filtering */
  searchable?: boolean
  /** Whether to allow multiple selection */
  multiple?: boolean
  /** Placeholder for search input */
  searchPlaceholder?: string
  /** Custom rendering for selected option */
  renderValue?: (value: T) => React.ReactNode
  /** Custom rendering for dropdown option */
  renderOption?: (option: { value: T; label: string }) => React.ReactNode
}

/**
 * Checkbox specific props
 */
export interface CheckboxFieldProps extends BaseFieldProps {
  /** Description text displayed next to checkbox */
  description?: string
}

/**
 * Textarea specific props
 */
export interface TextareaFieldProps extends BaseFieldProps {
  /** Number of visible rows */
  rows?: number
  /** Maximum character length */
  maxLength?: number
  /** Whether to show character counter */
  showCounter?: boolean
  /** Whether to auto-resize based on content */
  autoResize?: boolean
}

/**
 * Toggle/Switch specific props
 */
export interface ToggleFieldProps extends BaseFieldProps {
  /** Description text displayed next to toggle */
  description?: string
  /** Size variant */
  size?: 'sm' | 'md' | 'lg'
}

/**
 * Generic form values type
 *
 * Can be extended for type-safe form definitions:
 * @example
 * interface LoginFormValues extends FormValues {
 *   email: string
 *   password: string
 * }
 */
export type FormValues = Record<string, unknown>

/**
 * Field context value providing access to field state and methods
 */
export type FieldContextValue<TValue = unknown> = FieldApi<
  any,
  any,
  any,
  any,
  TValue
>

/**
 * Form context value providing access to form state and methods
 */
export type FormContextValue = FormApi<any, any>

/**
 * Validation mode timing options
 *
 * Controls when validation runs for a field:
 * - onBlur: Validate when field loses focus (default, best UX)
 * - onChange: Validate on every keystroke (real-time feedback)
 * - onSubmit: Validate only when form is submitted
 * - onDynamic: Validate based on form state changes (cross-field)
 */
export type ValidationMode = 'onBlur' | 'onChange' | 'onSubmit' | 'onDynamic'

/**
 * Server error shape from backend API
 *
 * Maps field names to error translation keys returned from backend
 * validation failures (400/422 responses)
 */
export type ServerErrors = Record<string, string>
