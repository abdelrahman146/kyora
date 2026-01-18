import { z } from 'zod'
import { UserSchema } from './auth'

/**
 * Zod schemas for onboarding API based on backend implementation
 * All schemas validate API request/response data for the onboarding flow
 */

// Onboarding Session Stage & Payment Status

export const SessionStageSchema = z.enum([
  'plan_selected',
  'identity_pending',
  'identity_verified',
  'business_staged',
  'payment_pending',
  'ready_to_commit',
  'committed',
])

export type SessionStage = z.infer<typeof SessionStageSchema>

export const PaymentStatusSchema = z.enum(['skipped', 'pending', 'succeeded'])

export type PaymentStatus = z.infer<typeof PaymentStatusSchema>

export const IdentityMethodSchema = z.enum(['email', 'google'])

export type IdentityMethod = z.infer<typeof IdentityMethodSchema>

// Plan Schema

export const PlanFeaturesSchema = z.object({
  customerManagement: z.boolean(),
  inventoryManagement: z.boolean(),
  orderManagement: z.boolean(),
  expenseManagement: z.boolean(),
  accounting: z.boolean(),
  basicAnalytics: z.boolean(),
  financialReports: z.boolean(),
  dataImport: z.boolean(),
  dataExport: z.boolean(),
  advancedAnalytics: z.boolean(),
  advancedFinancialReports: z.boolean(),
  orderPaymentLinks: z.boolean(),
  invoiceGeneration: z.boolean(),
  exportAnalyticsData: z.boolean(),
  aiBusinessAssistant: z.boolean(),
})

export type PlanFeatures = z.infer<typeof PlanFeaturesSchema>

export const PlanLimitsSchema = z.object({
  maxOrdersPerMonth: z.number(),
  maxTeamMembers: z.number(),
  maxBusinesses: z.number(),
})

export type PlanLimits = z.infer<typeof PlanLimitsSchema>

export const PlanSchema = z.object({
  id: z.string(),
  descriptor: z.string(),
  name: z.string(),
  description: z.string().optional().nullable(),
  price: z.string(),
  currency: z.string(),
  billingCycle: z.string(),
  stripePlanId: z.string().optional().nullable(),
  features: PlanFeaturesSchema,
  limits: PlanLimitsSchema,
  createdAt: z.string().optional(),
  updatedAt: z.string().optional(),
})

export type Plan = z.infer<typeof PlanSchema>

// Start Session - POST /v1/onboarding/start

export const StartSessionRequestSchema = z.object({
  email: z
    .string()
    .email({ message: 'validation.invalid_email' })
    .min(1, 'validation.required'),
  planDescriptor: z.string().min(1, 'validation.required'),
})

export type StartSessionRequest = z.infer<typeof StartSessionRequestSchema>

export const StartSessionResponseSchema = z.object({
  sessionToken: z.string(),
  stage: SessionStageSchema,
  isPaid: z.boolean(),
})

export type StartSessionResponse = z.infer<typeof StartSessionResponseSchema>

// Get Session - GET /v1/onboarding/session

export const GetSessionResponseSchema = z.object({
  sessionToken: z.string(),
  email: z.string(),
  planId: z.string(),
  planDescriptor: z.string(),
  stage: SessionStageSchema,
  isPaidPlan: z.boolean(),
  method: IdentityMethodSchema.nullable().optional(),
  emailVerified: z.boolean().optional(),
  firstName: z.string().nullable().optional(),
  lastName: z.string().nullable().optional(),
  businessName: z.string().nullable().optional(),
  businessDescriptor: z.string().nullable().optional(),
  businessCountry: z.string().nullable().optional(),
  businessCurrency: z.string().nullable().optional(),
  paymentStatus: PaymentStatusSchema,
  checkoutSessionId: z.string().nullable().optional(),
  otpExpiry: z.string().nullable().optional(),
  expiresAt: z.string().nullable().optional(),
})

export type GetSessionResponse = z.infer<typeof GetSessionResponseSchema>

// Send OTP - POST /v1/onboarding/email/otp

export const SendOTPRequestSchema = z.object({
  sessionToken: z.string().min(1, 'validation.required'),
})

export type SendOTPRequest = z.infer<typeof SendOTPRequestSchema>

export const SendOTPResponseSchema = z.object({
  retryAfterSeconds: z.number().int().nonnegative().optional(),
})

export type SendOTPResponse = z.infer<typeof SendOTPResponseSchema>

// Verify Email - POST /v1/onboarding/email/verify

export const VerifyOTPRequestSchema = z.object({
  sessionToken: z.string().min(1, 'validation.required'),
  code: z.string().length(6, 'validation.otp_length'),
  firstName: z.string().min(1, 'validation.required'),
  lastName: z.string().min(1, 'validation.required'),
  password: z.string().min(8, 'validation.password_min_length'),
})

export type VerifyOTPRequest = z.infer<typeof VerifyOTPRequestSchema>

export const VerifyOTPResponseSchema = z.object({
  stage: SessionStageSchema,
})

export type VerifyOTPResponse = z.infer<typeof VerifyOTPResponseSchema>

// OAuth Google - POST /v1/onboarding/oauth/google

export const OAuthGoogleRequestSchema = z.object({
  sessionToken: z.string().min(1, 'validation.required'),
  code: z.string().min(1, 'validation.required'),
})

export type OAuthGoogleRequest = z.infer<typeof OAuthGoogleRequestSchema>

export const OAuthGoogleResponseSchema = z.object({
  stage: SessionStageSchema,
})

export type OAuthGoogleResponse = z.infer<typeof OAuthGoogleResponseSchema>

// Set Business - POST /v1/onboarding/business

export const SetBusinessRequestSchema = z.object({
  sessionToken: z.string().min(1, 'validation.required'),
  businessName: z.string().min(1, 'validation.required'),
  businessDescriptor: z
    .string()
    .min(3, 'validation.business_descriptor_min_length')
    .max(20, 'validation.business_descriptor_max_length')
    .regex(/^[a-z0-9-]+$/, 'validation.business_descriptor_format'),
  country: z.string().min(1, 'validation.country_required'),
  currency: z.string().min(1, 'validation.required'),
})

export type SetBusinessRequest = z.infer<typeof SetBusinessRequestSchema>

export const SetBusinessResponseSchema = z.object({
  stage: SessionStageSchema,
})

export type SetBusinessResponse = z.infer<typeof SetBusinessResponseSchema>

// Payment Start - POST /v1/onboarding/payment/start

export const PaymentStartRequestSchema = z.object({
  sessionToken: z.string().min(1, 'validation.required'),
  successUrl: z.string().url('validation.invalid_url'),
  cancelUrl: z.string().url('validation.invalid_url'),
})

export type PaymentStartRequest = z.infer<typeof PaymentStartRequestSchema>

export const PaymentStartResponseSchema = z.object({
  checkoutUrl: z.string().url(),
})

export type PaymentStartResponse = z.infer<typeof PaymentStartResponseSchema>

// Complete Onboarding - POST /v1/onboarding/complete

export const CompleteOnboardingRequestSchema = z.object({
  sessionToken: z.string().min(1, 'validation.required'),
})

export type CompleteOnboardingRequest = z.infer<
  typeof CompleteOnboardingRequestSchema
>

export const CompleteOnboardingResponseSchema = z.object({
  token: z.string(),
  refreshToken: z.string(),
  user: UserSchema,
})

export type CompleteOnboardingResponse = z.infer<
  typeof CompleteOnboardingResponseSchema
>

// Get Plans - GET /v1/onboarding/plans

export const GetPlansResponseSchema = z.array(PlanSchema)

export type GetPlansResponse = z.infer<typeof GetPlansResponseSchema>
