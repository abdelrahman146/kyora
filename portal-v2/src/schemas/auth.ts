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
