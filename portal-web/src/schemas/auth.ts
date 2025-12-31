import { z } from 'zod'

/**
 * Login Form Validation Schema
 */
export const LoginSchema = z.object({
  email: z
    .string()
    .min(1, 'validation.required')
    .email('validation.invalid_email'),
  password: z.string().min(1, 'validation.required'),
})

export type LoginFormData = z.infer<typeof LoginSchema>

/**
 * Login Form Validators
 *
 * Multi-mode validation support:
 * - onBlur: Full validation when field loses focus
 */
export const createLoginValidators = () => ({
  email: {
    onBlur: LoginSchema.shape.email,
  },
  password: {
    onBlur: LoginSchema.shape.password,
  },
})

/**
 * Forgot Password Form Validation Schema
 */
export const ForgotPasswordSchema = z.object({
  email: z
    .string()
    .min(1, 'validation.required')
    .email('validation.invalid_email'),
})

export type ForgotPasswordFormData = z.infer<typeof ForgotPasswordSchema>

/**
 * Forgot Password Form Validators
 */
export const createForgotPasswordValidators = () => ({
  email: {
    onBlur: ForgotPasswordSchema.shape.email,
  },
})

/**
 * Reset Password Form Validation Schema
 */
export const ResetPasswordSchema = z
  .object({
    password: z.string().min(8, 'validation.password_min_length'),
    confirmPassword: z.string().min(1, 'validation.required'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'validation.passwords_must_match',
    path: ['confirmPassword'],
  })

export type ResetPasswordFormData = z.infer<typeof ResetPasswordSchema>

/**
 * Reset Password Form Validators
 *
 * With linked field validation:
 * - confirmPassword listens to password changes for real-time matching
 */
export const createResetPasswordValidators = () => ({
  password: {
    onBlur: ResetPasswordSchema.shape.password,
    // Optional: Add onChange for real-time password strength feedback
    onChange: ({ value }: { value: string }) => {
      if (value.length === 0) return undefined
      if (value.length < 8) return 'validation.password_min_length'
      return undefined
    },
  },
  confirmPassword: {
    onBlur: ResetPasswordSchema.shape.confirmPassword,
    // Listen to password field changes for real-time matching validation
    onChangeListenTo: ['password'],
    onChange: ({ value, fieldApi }: { value: string; fieldApi: any }) => {
      const password = fieldApi.form.getFieldValue('password')
      if (value.length === 0) return undefined
      if (value !== password) return 'validation.passwords_must_match'
      return undefined
    },
  },
})
