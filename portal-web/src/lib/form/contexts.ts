/**
 * Form Composition Contexts
 *
 * Creates React contexts for TanStack Form composition layer.
 * Provides `useFieldContext` and `useFormContext` hooks for accessing
 * field and form state within custom field components.
 *
 * @see https://tanstack.com/form/latest/docs/framework/react/guides/form-composition
 */

import { createFormHookContexts } from '@tanstack/react-form'

/**
 * Form composition contexts
 *
 * Usage:
 * - `useFieldContext()`: Access current field state within field components
 * - `useFormContext()`: Access form-level state and methods
 * - `fieldContext`: Pass to `createFormHook` for field component binding
 * - `formContext`: Pass to `createFormHook` for form component binding
 */
export const { fieldContext, formContext, useFieldContext, useFormContext } =
  createFormHookContexts()
