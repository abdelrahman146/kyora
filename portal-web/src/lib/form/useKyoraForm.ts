/**
 * Kyora Form Hook - Production-Grade Form Composition
 *
 * Custom form hook built on TanStack Form with pre-bound field components,
 * automatic validation, focus management, and i18n error translation.
 *
 * Features:
 * - Pre-bound field components (TextField, PasswordField, etc.)
 * - Automatic error translation from validation keys
 * - Progressive validation disclosure (revalidateLogic)
 * - Focus management on submission errors
 * - Server error injection support
 * - Type-safe field definitions
 * - Minimal boilerplate (~7 lines per field vs ~27 lines)
 *
 * @example
 * ```tsx
 * const form = useKyoraForm({
 *   defaultValues: {
 *     email: '',
 *     password: '',
 *   },
 *   onSubmit: async ({ value }) => {
 *     await loginMutation.mutateAsync(value)
 *   },
 * })
 *
 * return (
 *   <form.FormRoot>
 *     <form.Field name="email">
 *       {(field) => (
 *         <field.TextField
 *           label="Email"
 *           type="email"
 *           startIcon={<Mail />}
 *         />
 *       )}
 *     </form.Field>
 *
 *     <form.Field name="password">
 *       {(field) => (
 *         <field.PasswordField label="Password" />
 *       )}
 *     </form.Field>
 *
 *     <form.SubmitButton>Login</form.SubmitButton>
 *   </form.FormRoot>
 * )
 * ```
 */

import { createFormHook } from '@tanstack/react-form'
import { fieldContext, formContext } from './contexts'
import {
  AddressSelectField,
  CategorySelectField,
  CheckboxField,
  CheckboxGroupField,
  CustomerSelectField,
  DateField,
  DateRangeField,
  DateTimeField,
  ErrorInfo,
  FieldArray,
  FileUploadField,
  ImageUploadField,
  PasswordField,
  PriceField,
  ProductVariantSelectField,
  QuantityField,
  RadioField,
  SelectField,
  TextField,
  TextareaField,
  TimeField,
  ToggleField,
} from './components'
import { FormRoot } from './components/FormRoot'
import { SubmitButton } from './components/SubmitButton'
import { FormError } from './components/FormError'

/**
 * Create Kyora Form hook with pre-bound components and default configuration
 */
const { useAppForm, withForm } = createFormHook({
  fieldContext,
  formContext,

  // Pre-bind field components for use in forms
  fieldComponents: {
    AddressSelectField,
    CategorySelectField,
    TextField,
    PasswordField,
    TextareaField,
    SelectField,
    CheckboxField,
    CheckboxGroupField,
    CustomerSelectField,
    ProductVariantSelectField,
    RadioField,
    ToggleField,
    DateField,
    TimeField,
    DateTimeField,
    DateRangeField,
    FieldArray,
    FileUploadField,
    ImageUploadField,
    PriceField,
    QuantityField,
    ErrorInfo,
  },

  // Pre-bind form-level components
  formComponents: {
    FormRoot,
    SubmitButton,
    FormError,
  },
})

// Export with Kyora naming for consistency
export const useKyoraForm = useAppForm
export const withKyoraForm = withForm
