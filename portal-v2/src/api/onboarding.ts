import {
  useMutation,
  useQuery,
  type UseMutationOptions,
} from '@tanstack/react-query'
import { get, post, del } from './client'
import {
  StartSessionRequestSchema,
  StartSessionResponseSchema,
  GetSessionResponseSchema,
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
  type GetSessionResponse,
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
} from './types/onboarding'
import { z } from 'zod'

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
  /**
   * Initializes or resumes an onboarding session for an email and plan.
   * @returns Session token and stage information
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async startSession(data: StartSessionRequest): Promise<StartSessionResponse> {
    const validatedRequest = StartSessionRequestSchema.parse(data)

    const response = await post<unknown>('v1/onboarding/start', {
      json: validatedRequest,
    })

    return StartSessionResponseSchema.parse(response)
  },

  /**
   * Retrieves the current onboarding session state by token.
   * Use this to restore/resume onboarding flow when user returns.
   * @returns Complete session state including stage, user data, and business details
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async getSession(sessionToken: string): Promise<GetSessionResponse> {
    const response = await get<unknown>(
      `v1/onboarding/session?sessionToken=${sessionToken}`,
    )

    return GetSessionResponseSchema.parse(response)
  },

  /**
   * Generates a 6-digit OTP and sends it to the user's email
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async sendEmailOTP(data: SendOTPRequest): Promise<SendOTPResponse> {
    const validatedRequest = SendOTPRequestSchema.parse(data)

    const response = await post<unknown>('v1/onboarding/email/otp', {
      json: validatedRequest,
    })

    return SendOTPResponseSchema.parse(response)
  },

  /**
   * Validates OTP code and stores user profile and password
   * @returns Updated session stage
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async verifyEmail(data: VerifyEmailRequest): Promise<VerifyEmailResponse> {
    const validatedRequest = VerifyEmailRequestSchema.parse(data)

    const response = await post<unknown>('v1/onboarding/email/verify', {
      json: validatedRequest,
    })

    return VerifyEmailResponseSchema.parse(response)
  },

  /**
   * Sets OAuth identity from Google and stages user profile
   * @returns Updated session stage
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async oauthGoogle(data: OAuthGoogleRequest): Promise<OAuthGoogleResponse> {
    const validatedRequest = OAuthGoogleRequestSchema.parse(data)

    const response = await post<unknown>('v1/onboarding/oauth/google', {
      json: validatedRequest,
    })

    return OAuthGoogleResponseSchema.parse(response)
  },

  /**
   * Stages business details for the onboarding session
   * @returns Updated session stage
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async setBusiness(data: SetBusinessRequest): Promise<SetBusinessResponse> {
    const validatedRequest = SetBusinessRequestSchema.parse(data)

    const response = await post<unknown>('v1/onboarding/business', {
      json: validatedRequest,
    })

    return SetBusinessResponseSchema.parse(response)
  },

  /**
   * Creates Stripe checkout session for paid plans
   * @returns Checkout URL (null for free plans)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async startPayment(data: PaymentStartRequest): Promise<PaymentStartResponse> {
    const validatedRequest = PaymentStartRequestSchema.parse(data)

    const response = await post<unknown>('v1/onboarding/payment/start', {
      json: validatedRequest,
    })

    return PaymentStartResponseSchema.parse(response)
  },

  /**
   * Finalizes onboarding and commits all staged data to permanent tables
   * @returns User data and JWT tokens
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async complete(
    data: CompleteOnboardingRequest,
  ): Promise<CompleteOnboardingResponse> {
    const validatedRequest = CompleteOnboardingRequestSchema.parse(data)

    const response = await post<unknown>('v1/onboarding/complete', {
      json: validatedRequest,
    })

    return CompleteOnboardingResponseSchema.parse(response)
  },

  /**
   * Fetches all available billing plans
   * @returns Array of plans
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async listPlans(): Promise<Plan[]> {
    const response = await get<unknown>('v1/billing/plans')

    return z.array(PlanSchema).parse(response)
  },

  /**
   * Fetches a specific plan by descriptor
   * @param descriptor - Plan descriptor (e.g., "free", "starter", "professional")
   * @returns Plan details
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async getPlan(descriptor: string): Promise<Plan> {
    const response = await get<unknown>(`v1/billing/plans/${descriptor}`)

    return PlanSchema.parse(response)
  },

  /**
   * Deletes an onboarding session (cancel/restart flow).
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async deleteSession(sessionToken: string): Promise<void> {
    await del('v1/onboarding/session', {
      json: { sessionToken },
    })
  },
}

/**
 * TanStack Query Hooks for Onboarding
 */

/**
 * Query to fetch onboarding session by token
 */
export function useOnboardingSessionQuery(sessionToken: string | null) {
  return useQuery({
    queryKey: ['onboarding', 'session', sessionToken],
    queryFn: () => onboardingApi.getSession(sessionToken!),
    enabled: !!sessionToken,
    staleTime: 0, // Always fresh - onboarding is time-sensitive
    gcTime: 0, // Don't cache - session may be invalidated server-side
  })
}

/**
 * Query to fetch all available plans
 */
export function usePlansQuery() {
  return useQuery({
    queryKey: ['onboarding', 'plans'],
    queryFn: () => onboardingApi.listPlans(),
    staleTime: 5 * 60 * 1000, // 5 minutes - plans are semi-static
  })
}

/**
 * Query to fetch a specific plan by descriptor
 */
export function usePlanQuery(descriptor: string | null) {
  return useQuery({
    queryKey: ['onboarding', 'plan', descriptor],
    queryFn: () => onboardingApi.getPlan(descriptor!),
    enabled: !!descriptor,
    staleTime: 5 * 60 * 1000, // 5 minutes - plans are semi-static
  })
}

/**
 * Mutation to start onboarding session
 */
export function useStartSessionMutation(
  options?: UseMutationOptions<
    StartSessionResponse,
    Error,
    StartSessionRequest
  >,
) {
  return useMutation({
    mutationFn: (data: StartSessionRequest) => onboardingApi.startSession(data),
    ...options,
  })
}

/**
 * Mutation to send email OTP
 */
export function useSendOTPMutation(
  options?: UseMutationOptions<SendOTPResponse, Error, SendOTPRequest>,
) {
  return useMutation({
    mutationFn: (data: SendOTPRequest) => onboardingApi.sendEmailOTP(data),
    ...options,
  })
}

/**
 * Mutation to verify email with OTP
 */
export function useVerifyEmailMutation(
  options?: UseMutationOptions<VerifyEmailResponse, Error, VerifyEmailRequest>,
) {
  return useMutation({
    mutationFn: (data: VerifyEmailRequest) => onboardingApi.verifyEmail(data),
    ...options,
  })
}

/**
 * Mutation to handle Google OAuth
 */
export function useOAuthGoogleMutation(
  options?: UseMutationOptions<OAuthGoogleResponse, Error, OAuthGoogleRequest>,
) {
  return useMutation({
    mutationFn: (data: OAuthGoogleRequest) => onboardingApi.oauthGoogle(data),
    ...options,
  })
}

/**
 * Mutation to set business details
 */
export function useSetBusinessMutation(
  options?: UseMutationOptions<SetBusinessResponse, Error, SetBusinessRequest>,
) {
  return useMutation({
    mutationFn: (data: SetBusinessRequest) => onboardingApi.setBusiness(data),
    ...options,
  })
}

/**
 * Mutation to start payment process
 */
export function useStartPaymentMutation(
  options?: UseMutationOptions<
    PaymentStartResponse,
    Error,
    PaymentStartRequest
  >,
) {
  return useMutation({
    mutationFn: (data: PaymentStartRequest) => onboardingApi.startPayment(data),
    ...options,
  })
}

/**
 * Mutation to complete onboarding
 */
export function useCompleteOnboardingMutation(
  options?: UseMutationOptions<
    CompleteOnboardingResponse,
    Error,
    CompleteOnboardingRequest
  >,
) {
  return useMutation({
    mutationFn: (data: CompleteOnboardingRequest) =>
      onboardingApi.complete(data),
    ...options,
  })
}

/**
 * Mutation to delete onboarding session
 */
export function useDeleteSessionMutation(
  options?: UseMutationOptions<void, Error, string>,
) {
  return useMutation({
    mutationFn: (sessionToken: string) =>
      onboardingApi.deleteSession(sessionToken),
    ...options,
  })
}
