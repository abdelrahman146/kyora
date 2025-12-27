import { get, post } from "./client";
import {
  StartSessionRequestSchema,
  StartSessionResponseSchema,
  SendOTPRequestSchema,
  SendOTPResponseSchema,
  VerifyEmailRequestSchema,
  VerifyEmailResponseSchema,
  OAuthGoogleRequestSchema,
  OAuthGoogleResponseSchema,
  SetBusinessRequestSchema,
  SetBusinessResponseSchema,
  PaymentStartRequestSchema,
  PaymentStartResponseSchema,
  CompleteOnboardingRequestSchema,
  CompleteOnboardingResponseSchema,
  PlanSchema,
  type StartSessionRequest,
  type StartSessionResponse,
  type SendOTPRequest,
  type SendOTPResponse,
  type VerifyEmailRequest,
  type VerifyEmailResponse,
  type OAuthGoogleRequest,
  type OAuthGoogleResponse,
  type SetBusinessRequest,
  type SetBusinessResponse,
  type PaymentStartRequest,
  type PaymentStartResponse,
  type CompleteOnboardingRequest,
  type CompleteOnboardingResponse,
  type Plan,
} from "./types/onboarding";
import { z } from "zod";

/**
 * Onboarding API Service
 *
 * All methods validate request/response data using Zod schemas
 * and provide type-safe interfaces to backend onboarding endpoints.
 *
 * Onboarding Flow:
 * 1. Start Session - Select plan and provide email
 * 2. Verify Identity - Email OTP or Google OAuth
 * 3. Set Business - Business name, descriptor, country, currency
 * 4. Payment - Stripe checkout for paid plans (skipped for free)
 * 5. Complete - Finalize onboarding and receive JWT tokens
 */

export const onboardingApi = {
  // ==========================================================================
  // Start Session - POST /v1/onboarding/start
  // ==========================================================================

  /**
   * Initializes or resumes an onboarding session for an email and plan.
   * @returns Session token and stage information
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async startSession(data: StartSessionRequest): Promise<StartSessionResponse> {
    const validatedRequest = StartSessionRequestSchema.parse(data);

    const response = await post<unknown>("v1/onboarding/start", {
      json: validatedRequest,
    });

    return StartSessionResponseSchema.parse(response);
  },

  // ==========================================================================
  // Send Email OTP - POST /v1/onboarding/email/otp
  // ==========================================================================

  /**
   * Generates a 6-digit OTP and sends it to the user's email
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async sendEmailOTP(data: SendOTPRequest): Promise<SendOTPResponse> {
    const validatedRequest = SendOTPRequestSchema.parse(data);

    const response = await post<unknown>("v1/onboarding/email/otp", {
      json: validatedRequest,
    });

    return SendOTPResponseSchema.parse(response);
  },

  // ==========================================================================
  // Verify Email - POST /v1/onboarding/email/verify
  // ==========================================================================

  /**
   * Validates OTP code and stores user profile and password
   * @returns Updated session stage
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async verifyEmail(data: VerifyEmailRequest): Promise<VerifyEmailResponse> {
    const validatedRequest = VerifyEmailRequestSchema.parse(data);

    const response = await post<unknown>("v1/onboarding/email/verify", {
      json: validatedRequest,
    });

    return VerifyEmailResponseSchema.parse(response);
  },

  // ==========================================================================
  // OAuth Google - POST /v1/onboarding/oauth/google
  // ==========================================================================

  /**
   * Sets OAuth identity from Google and stages user profile
   * @returns Updated session stage
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async oauthGoogle(data: OAuthGoogleRequest): Promise<OAuthGoogleResponse> {
    const validatedRequest = OAuthGoogleRequestSchema.parse(data);

    const response = await post<unknown>("v1/onboarding/oauth/google", {
      json: validatedRequest,
    });

    return OAuthGoogleResponseSchema.parse(response);
  },

  // ==========================================================================
  // Set Business - POST /v1/onboarding/business
  // ==========================================================================

  /**
   * Stages business details for the onboarding session
   * @returns Updated session stage
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async setBusiness(data: SetBusinessRequest): Promise<SetBusinessResponse> {
    const validatedRequest = SetBusinessRequestSchema.parse(data);

    const response = await post<unknown>("v1/onboarding/business", {
      json: validatedRequest,
    });

    return SetBusinessResponseSchema.parse(response);
  },

  // ==========================================================================
  // Payment Start - POST /v1/onboarding/payment/start
  // ==========================================================================

  /**
   * Creates Stripe checkout session for paid plans
   * @returns Checkout URL (null for free plans)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async startPayment(data: PaymentStartRequest): Promise<PaymentStartResponse> {
    const validatedRequest = PaymentStartRequestSchema.parse(data);

    const response = await post<unknown>("v1/onboarding/payment/start", {
      json: validatedRequest,
    });

    return PaymentStartResponseSchema.parse(response);
  },

  // ==========================================================================
  // Complete Onboarding - POST /v1/onboarding/complete
  // ==========================================================================

  /**
   * Finalizes onboarding and commits all staged data to permanent tables
   * @returns User data and JWT tokens
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async complete(
    data: CompleteOnboardingRequest
  ): Promise<CompleteOnboardingResponse> {
    const validatedRequest = CompleteOnboardingRequestSchema.parse(data);

    const response = await post<unknown>("v1/onboarding/complete", {
      json: validatedRequest,
    });

    return CompleteOnboardingResponseSchema.parse(response);
  },

  // ==========================================================================
  // List Plans - GET /v1/billing/plans
  // ==========================================================================

  /**
   * Fetches all available billing plans
   * @returns Array of plans
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async listPlans(): Promise<Plan[]> {
    const response = await get<unknown>("v1/billing/plans");

    return z.array(PlanSchema).parse(response);
  },

  // ==========================================================================
  // Get Plan - GET /v1/billing/plans/:descriptor
  // ==========================================================================

  /**
   * Fetches a specific plan by descriptor
   * @param descriptor - Plan descriptor (e.g., "free", "starter", "professional")
   * @returns Plan details
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async getPlan(descriptor: string): Promise<Plan> {
    const response = await get<unknown>(`v1/billing/plans/${descriptor}`);

    return PlanSchema.parse(response);
  },
};
