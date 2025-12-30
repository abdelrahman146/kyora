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

import { createFormHook, revalidateLogic } from '@tanstack/react-form'
import { zodValidator } from '@tanstack/zod-form-adapter'
import { fieldContext, formContext } from './contexts'
import { TextField, PasswordField, ErrorInfo } from './components'
import { FormRoot } from './components/FormRoot'
import { SubmitButton } from './components/SubmitButton'
import { FormError } from './components/FormError'

/**
 * Create Kyora Form hook with pre-bound components and default configuration
 */
export const { useKyoraForm, withKyoraForm } = createFormHook({
  fieldContext,
  formContext,

  // Pre-bind field components for use in forms
  fieldComponents: {
    TextField,
    PasswordField,
    ErrorInfo,
  },

  // Pre-bind form-level components
  formComponents: {
    FormRoot,
    SubmitButton,
    FormError,
  },

  // Default form configuration
  defaultValues: {},

  // Enable Zod schema validation
  validatorAdapter: zodValidator(),

  // Progressive validation disclosure for optimal UX:
  // - Before first submit: only validate on submit (mode: 'submit')
  // - After first submit: validate on blur (modeAfterSubmission: 'blur')
  // This prevents annoying errors while typing, but provides feedback after interaction
  validationLogic: revalidateLogic({
    mode: 'submit',
    modeAfterSubmission: 'blur',
  }),
})
