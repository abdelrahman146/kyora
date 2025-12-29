import { z } from "zod";
import { UserSchema } from "./auth";

/**
 * Zod schemas for onboarding API based on backend implementation
 * All schemas validate API request/response data for the onboarding flow
 */

// Onboarding Session Stage & Payment Status

export const SessionStageSchema = z.enum([
  "plan_selected",
  "identity_pending",
  "identity_verified",
  "business_staged",
  "payment_pending",
  "payment_confirmed",
  "ready_to_commit",
  "committed",
]);

export type SessionStage = z.infer<typeof SessionStageSchema>;

export const PaymentStatusSchema = z.enum(["skipped", "pending", "succeeded"]);

export type PaymentStatus = z.infer<typeof PaymentStatusSchema>;

export const IdentityMethodSchema = z.enum(["email", "google"]);

export type IdentityMethod = z.infer<typeof IdentityMethodSchema>;

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
});

export type PlanFeatures = z.infer<typeof PlanFeaturesSchema>;

export const PlanLimitsSchema = z.object({
  maxOrdersPerMonth: z.number(),
  maxTeamMembers: z.number(),
  maxBusinesses: z.number(),
});

export type PlanLimits = z.infer<typeof PlanLimitsSchema>;

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
});

export type Plan = z.infer<typeof PlanSchema>;

// Start Session - POST /v1/onboarding/start

export const StartSessionRequestSchema = z.object({
  email: z.email({ message: "Invalid email address" }),
  planDescriptor: z.string().min(1, "Plan is required"),
});

export type StartSessionRequest = z.infer<typeof StartSessionRequestSchema>;

export const StartSessionResponseSchema = z.object({
  sessionToken: z.string(),
  stage: SessionStageSchema,
  isPaid: z.boolean(),
});

export type StartSessionResponse = z.infer<typeof StartSessionResponseSchema>;

// Get Session - GET /v1/onboarding/session

export const GetSessionResponseSchema = z.object({
  sessionToken: z.string(),
  email: z.string(),
  stage: SessionStageSchema,
  emailVerified: z.boolean(),
  method: IdentityMethodSchema,
  firstName: z.string().optional(),
  lastName: z.string().optional(),
  planId: z.string(),
  planDescriptor: z.string(),
  isPaidPlan: z.boolean(),
  businessName: z.string().optional(),
  businessDescriptor: z.string().optional(),
  businessCountry: z.string().optional(),
  businessCurrency: z.string().optional(),
  paymentStatus: PaymentStatusSchema,
  checkoutSessionId: z.string().optional(),
  otpExpiry: z.string().optional().nullable(),
  expiresAt: z.string(),
});

export type GetSessionResponse = z.infer<typeof GetSessionResponseSchema>;

// Send Email OTP - POST /v1/onboarding/email/otp

export const SendOTPRequestSchema = z.object({
  sessionToken: z.string().min(1, "Session token is required"),
});

export type SendOTPRequest = z.infer<typeof SendOTPRequestSchema>;

export const SendOTPResponseSchema = z.object({
  retryAfterSeconds: z.number().int().nonnegative(),
});

export type SendOTPResponse = z.infer<typeof SendOTPResponseSchema>;

// Verify Email - POST /v1/onboarding/email/verify

export const VerifyEmailRequestSchema = z.object({
  sessionToken: z.string().min(1, "Session token is required"),
  code: z.string().length(6, "Code must be 6 digits"),
  firstName: z.string().min(1, "First name is required").max(100),
  lastName: z.string().min(1, "Last name is required").max(100),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export type VerifyEmailRequest = z.infer<typeof VerifyEmailRequestSchema>;

export const VerifyEmailResponseSchema = z.object({
  stage: SessionStageSchema,
});

export type VerifyEmailResponse = z.infer<typeof VerifyEmailResponseSchema>;

// OAuth Google - POST /v1/onboarding/oauth/google

export const OAuthGoogleRequestSchema = z.object({
  sessionToken: z.string().min(1, "Session token is required"),
  code: z.string().min(1, "OAuth code is required"),
});

export type OAuthGoogleRequest = z.infer<typeof OAuthGoogleRequestSchema>;

export const OAuthGoogleResponseSchema = z.object({
  stage: SessionStageSchema,
});

export type OAuthGoogleResponse = z.infer<typeof OAuthGoogleResponseSchema>;

// Set Business - POST /v1/onboarding/business

export const SetBusinessRequestSchema = z.object({
  sessionToken: z.string().min(1, "Session token is required"),
  name: z.string().min(1, "Business name is required"),
  descriptor: z
    .string()
    .min(3, "Business descriptor must be at least 3 characters")
    .max(50, "Business descriptor must be at most 50 characters")
    .regex(
      /^[a-z0-9-]+$/,
      "Business descriptor must contain only lowercase letters, numbers, and hyphens"
    ),
  country: z.string().length(2, "Country code must be 2 characters"),
  currency: z.string().length(3, "Currency code must be 3 characters"),
});

export type SetBusinessRequest = z.infer<typeof SetBusinessRequestSchema>;

export const SetBusinessResponseSchema = z.object({
  stage: SessionStageSchema,
});

export type SetBusinessResponse = z.infer<typeof SetBusinessResponseSchema>;

// Payment Start - POST /v1/onboarding/payment/start

export const PaymentStartRequestSchema = z.object({
  sessionToken: z.string().min(1, "Session token is required"),
  successUrl: z.url({ message: "Invalid success URL" }),
  cancelUrl: z.url({ message: "Invalid cancel URL" }),
});

export type PaymentStartRequest = z.infer<typeof PaymentStartRequestSchema>;

export const PaymentStartResponseSchema = z.object({
  checkoutUrl: z.url().optional().nullable(),
});

export type PaymentStartResponse = z.infer<typeof PaymentStartResponseSchema>;

// Complete Onboarding - POST /v1/onboarding/complete

export const CompleteOnboardingRequestSchema = z.object({
  sessionToken: z.string().min(1, "Session token is required"),
});

export type CompleteOnboardingRequest = z.infer<
  typeof CompleteOnboardingRequestSchema
>;

export const CompleteOnboardingResponseSchema = z.object({
  user: UserSchema,
  token: z.string(),
  refreshToken: z.string(),
});

export type CompleteOnboardingResponse = z.infer<
  typeof CompleteOnboardingResponseSchema
>;
