/**
 * Kyora Form System - Production-Grade Form Management
 *
 * Unified form composition layer built on TanStack Form providing:
 * - Type-safe form definitions with Zod validation
 * - Automatic error translation (i18n)
 * - Progressive validation disclosure (revalidateLogic)
 * - Focus management on validation errors
 * - Pre-bound field components (75% less boilerplate)
 * - Server error injection
 * - Accessible form controls
 *
 * @example
 * ```tsx
 * import { useKyoraForm } from '@/lib/form'
 *
 * const form = useKyoraForm({
 *   defaultValues: { email: '', password: '' },
 *   onSubmit: async ({ value }) => {
 *     await api.login(value)
 *   },
 * })
 *
 * return (
 *   <form.FormRoot>
 *     <form.FormError />
 *     <form.Field name="email">
 *       {(field) => (
 *         <field.TextField label="Email" type="email" />
 *       )}
 *     </form.Field>
 *     <form.SubmitButton>Login</form.SubmitButton>
 *   </form.FormRoot>
 * )
 * ```
 */

export { useKyoraForm, withKyoraForm } from './useKyoraForm'
export { useFieldContext, useFormContext } from './contexts'
export {
  useFocusOnError,
  useAutoFocus,
  createFocusManagement,
} from './useFocusManagement'
export { useServerErrors, translateServerError } from './useServerErrors'
export * from './types'
export * from './components'
