import { z } from 'zod'

/**
 * Onboarding Form Validation Schemas
 */

/**
 * Email Form Schema
 */
export const EmailFormSchema = z.object({
  email: z
    .string()
    .min(1, 'validation.required')
    .email('validation.invalid_email'),
})

export type EmailFormData = z.infer<typeof EmailFormSchema>

/**
 * OTP Verification Schema
 */
export const OTPVerificationSchema = z.object({
  code: z
    .string()
    .length(6, 'validation.otp_length')
    .regex(/^\d+$/, 'validation.otp_digits_only'),
  fullName: z
    .string()
    .min(1, 'validation.required')
    .max(100, 'validation.max_length'),
  password: z.string().min(8, 'validation.password_min_length'),
})

export type OTPVerificationFormData = z.infer<typeof OTPVerificationSchema>

/**
 * Business Setup Schema
 */
export const BusinessSetupSchema = z.object({
  name: z
    .string()
    .min(1, 'validation.required')
    .max(100, 'validation.max_length'),
  descriptor: z
    .string()
    .min(3, 'validation.min_length')
    .max(20, 'validation.max_length')
    .regex(/^[a-z0-9-]+$/, 'validation.business_descriptor_format'),
  country: z.string().length(2, 'validation.country_code'),
  currency: z.string().length(3, 'validation.currency_code'),
})

export type BusinessSetupFormData = z.infer<typeof BusinessSetupSchema>
