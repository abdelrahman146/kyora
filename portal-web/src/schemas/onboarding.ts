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
 * Email Form Validators
 */
export const createEmailFormValidators = () => ({
  email: {
    onBlur: EmailFormSchema.shape.email,
  },
})

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
 * OTP Verification Form Validators
 */
export const createOTPVerificationValidators = () => ({
  code: {
    onBlur: OTPVerificationSchema.shape.code,
  },
  fullName: {
    onBlur: OTPVerificationSchema.shape.fullName,
  },
  password: {
    onBlur: OTPVerificationSchema.shape.password,
    // Real-time password strength feedback
    onChange: ({ value }: { value: string }) => {
      if (value.length === 0) return undefined
      if (value.length < 8) return 'validation.password_min_length'
      return undefined
    },
  },
})

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

/**
 * Business Setup Form Validators
 *
 * With dynamic descriptor generation:
 * - descriptor field uses listeners to auto-generate from business name
 * - onChange validation for slug format
 */
export const createBusinessSetupValidators = () => ({
  name: {
    onBlur: BusinessSetupSchema.shape.name,
  },
  descriptor: {
    onBlur: BusinessSetupSchema.shape.descriptor,
    // Real-time slug validation for immediate feedback
    onChange: ({ value }: { value: string }) => {
      if (value.length === 0) return undefined
      if (value.length < 3) return 'validation.min_length'
      if (value.length > 20) return 'validation.max_length'
      if (!/^[a-z0-9-]+$/.test(value))
        return 'validation.business_descriptor_format'
      return undefined
    },
  },
  country: {
    onBlur: BusinessSetupSchema.shape.country,
  },
  currency: {
    onBlur: BusinessSetupSchema.shape.currency,
  },
})
