/**
 * Form System Type Definitions
 *
 * Core types for the Kyora form composition layer built on TanStack Form.
 * Provides type-safe field props, validation error shapes, and form values.
 */

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
  /** Hint for which keyboard to show on mobile */
  inputMode?:
    | 'none'
    | 'text'
    | 'tel'
    | 'url'
    | 'email'
    | 'numeric'
    | 'decimal'
    | 'search'
  /** Mobile keyboard action button label */
  enterKeyHint?:
    | 'enter'
    | 'done'
    | 'go'
    | 'next'
    | 'previous'
    | 'search'
    | 'send'
  /** Control mobile/browser auto-capitalization */
  autoCapitalize?: 'none' | 'sentences' | 'words' | 'characters'
  /** Control browser autocorrect */
  autoCorrect?: 'on' | 'off'
  /** Control spellcheck */
  spellCheck?: boolean
  /** Force direction for mixed RTL/LTR fields (e.g., phone/order id) */
  dir?: 'ltr' | 'rtl' | 'auto'
  /** Optional pattern for numeric/codes (use with inputMode) */
  pattern?: string
  /** Maximum character length */
  maxLength?: number
}

/**
 * Price input specific props
 */
export interface PriceFieldProps extends BaseFieldProps {
  /** Currency code to display inside the field */
  currencyCode?: string
  /** Autocomplete attribute for browser autofill */
  autoComplete?: string
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
  options: Array<{
    value: T
    label: string
    description?: string
    icon?: React.ReactNode
    disabled?: boolean
    renderCustom?: () => React.ReactNode
  }>
  /** Whether to allow search/filtering */
  searchable?: boolean
  /** Whether to allow multiple selection */
  multiSelect?: boolean
  /** Whether to show clear button */
  clearable?: boolean
  /** Maximum height of dropdown menu */
  maxHeight?: number
  /** Size variant */
  size?: 'sm' | 'md' | 'lg'
  /** Style variant */
  variant?: 'default' | 'filled' | 'ghost'
  /** Help text displayed below select */
  helperText?: string
  /** Whether to take full width */
  fullWidth?: boolean
}

/**
 * Checkbox specific props
 */
export interface CheckboxFieldProps extends BaseFieldProps {
  /** Description text displayed next to checkbox */
  description?: string
}

/**
 * Checkbox group specific props
 */
export interface CheckboxGroupFieldProps<
  T extends string = string,
> extends BaseFieldProps {
  /** Description text displayed below label */
  description?: string
  /** Options for the checkbox group */
  options: Array<{
    value: T
    label: string
    /** Optional icon/emoji to display next to the label */
    icon?: string
    description?: string
    disabled?: boolean
  }>
}

/**
 * Customer select field specific props
 */
export interface CustomerSelectFieldProps extends BaseFieldProps {
  /** Business descriptor for fetching customers */
  businessDescriptor: string
  /** Placeholder text for the select */
  placeholder?: string
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
 * Radio button group specific props
 */
export interface RadioFieldProps extends BaseFieldProps {
  /** Options for the radio group */
  options: Array<{
    value: string
    label: string
    description?: string
    disabled?: boolean
  }>
  /** Layout orientation */
  orientation?: 'vertical' | 'horizontal'
  /** Size variant */
  size?: 'sm' | 'md' | 'lg'
  /** Color variant */
  variant?: 'default' | 'primary' | 'secondary'
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

/**
 * Array item operations provided to render function in FieldArray
 */
export interface ArrayItemOperations {
  /** Remove this item from the array */
  remove: () => void
  /** Move this item up in the list */
  moveUp: () => void
  /** Move this item down in the list */
  moveDown: () => void
}

/**
 * FieldArray specific props for dynamic array field management
 */
export interface FieldArrayProps<T = any> extends BaseFieldProps {
  /** Minimum number of items required */
  minItems?: number
  /** Maximum number of items allowed */
  maxItems?: number
  /** Factory function to create new item with default values */
  defaultValue?: () => T
  /** Label for the add button */
  addButtonLabel?: string
  /** Message to show when array is empty */
  emptyMessage?: string
  /** Whether items can be reordered via drag-and-drop */
  reorderable?: boolean
  /** Helper text displayed below array field */
  helperText?: string
  /** Render function for each array item */
  render: (
    item: T,
    index: number,
    operations: ArrayItemOperations,
  ) => React.ReactNode
}

/**
 * FileUploadField props for generic file uploads
 */
export interface FileUploadFieldProps extends BaseFieldProps {
  /** Accepted file types (MIME types or extensions) */
  accept?: string
  /** Maximum number of files */
  maxFiles?: number
  /** Maximum file size (e.g., "10MB") */
  maxSize?: string
  /** Allow multiple file selection */
  multiple?: boolean
  /** Enable drag-and-drop reordering */
  reorderable?: boolean
  /** Visual required indicator */
  required?: boolean
  /** Callback when upload completes */
  onUploadComplete?: (references: Array<any>) => void
}

/**
 * ImageUploadField props for image-only uploads
 */
export interface ImageUploadFieldProps extends Omit<
  FileUploadFieldProps,
  'accept' | 'multiple'
> {
  /** Single image mode (default: multiple) */
  single?: boolean
}
