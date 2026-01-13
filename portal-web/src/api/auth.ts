import { useMutation } from '@tanstack/react-query'

import {
  clearTokens,
  get,
  getRefreshToken,
  post,
  postVoid,
  setTokens,
} from './client'
import {
  ForgotPasswordRequestSchema,
  GoogleLoginRequestSchema,
  LoginRequestSchema,
  LoginResponseSchema,
  LogoutAllRequestSchema,
  LogoutRequestSchema,
  RefreshRequestSchema,
  RefreshResponseSchema,
  RequestEmailVerificationSchema,
  ResetPasswordRequestSchema,
  VerifyEmailRequestSchema,
} from './types/auth'
import type {
  ForgotPasswordRequest,
  GoogleLoginRequest,
  LoginRequest,
  LoginResponse,
  LogoutAllRequest,
  LogoutRequest,
  RefreshRequest,
  RefreshResponse,
  RequestEmailVerification,
  ResetPasswordRequest,
  VerifyEmailRequest,
} from './types/auth'

import type { UseMutationOptions } from '@tanstack/react-query'

/**
 * Authentication API Service
 *
 * All methods validate request/response data using Zod schemas
 * and provide type-safe interfaces to backend auth endpoints.
 */

export const authApi = {
  /**
   * Authenticates a user with email and password
   * @returns LoginResponse with access token, refresh token, and user data
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    // Validate request data
    const validatedRequest = LoginRequestSchema.parse(credentials)

    // Make API call
    const response = await post<unknown>('v1/auth/login', {
      json: validatedRequest,
    })

    // Validate response data
    const validatedResponse = LoginResponseSchema.parse(response)

    // Store tokens in memory
    setTokens(validatedResponse.token, validatedResponse.refreshToken)

    return validatedResponse
  },

  /**
   * Refreshes the access token using a refresh token
   * @returns RefreshResponse with new access token and refresh token
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async refreshToken(request: RefreshRequest): Promise<RefreshResponse> {
    // Validate request data
    const validatedRequest = RefreshRequestSchema.parse(request)

    // Make API call
    const response = await post<unknown>('v1/auth/refresh', {
      json: validatedRequest,
    })

    // Validate response data
    const validatedResponse = RefreshResponseSchema.parse(response)

    // Update tokens in memory
    setTokens(validatedResponse.token, validatedResponse.refreshToken)

    return validatedResponse
  },

  /**
   * Logs out the current user by revoking the refresh token
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async logout(request: LogoutRequest): Promise<void> {
    // Validate request data
    const validatedRequest = LogoutRequestSchema.parse(request)

    // Make API call
    await postVoid('v1/auth/logout', { json: validatedRequest })

    // Clear tokens from memory
    clearTokens()
  },

  /**
   * Logs out the user from all devices by revoking all refresh tokens
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async logoutAll(request: LogoutAllRequest): Promise<void> {
    // Validate request data
    const validatedRequest = LogoutAllRequestSchema.parse(request)

    // Make API call
    await postVoid('v1/auth/logout-all', { json: validatedRequest })

    // Clear tokens from memory
    clearTokens()
  },

  /**
   * Sends a password reset email to the user
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async forgotPassword(request: ForgotPasswordRequest): Promise<void> {
    // Validate request data
    const validatedRequest = ForgotPasswordRequestSchema.parse(request)

    // Make API call
    await postVoid('v1/auth/forgot-password', { json: validatedRequest })
  },

  /**
   * Resets user password using a valid reset token
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async resetPassword(request: ResetPasswordRequest): Promise<void> {
    // Validate request data
    const validatedRequest = ResetPasswordRequestSchema.parse(request)

    // Make API call
    await postVoid('v1/auth/reset-password', { json: validatedRequest })
  },

  /**
   * Authenticates a user using Google OAuth code
   * @returns LoginResponse with access token, refresh token, and user data
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async loginWithGoogle(request: GoogleLoginRequest): Promise<LoginResponse> {
    // Validate request data
    const validatedRequest = GoogleLoginRequestSchema.parse(request)

    // Make API call
    const response = await post<unknown>('v1/auth/google/login', {
      json: validatedRequest,
    })

    // Validate response data
    const validatedResponse = LoginResponseSchema.parse(response)

    // Store tokens in memory
    setTokens(validatedResponse.token, validatedResponse.refreshToken)

    return validatedResponse
  },

  /**
   * Gets the Google OAuth authorization URL for user authentication
   * @returns Object containing the authorization URL
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async getGoogleAuthUrl(): Promise<{ url: string }> {
    // Make API call
    const response = await get<unknown>('v1/auth/google/url')

    // Validate response structure
    if (
      typeof response !== 'object' ||
      response === null ||
      !('url' in response)
    ) {
      throw new Error('Invalid response from Google OAuth URL endpoint')
    }

    return response as { url: string }
  },

  /**
   * Sends an email verification link to the user
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async requestEmailVerification(
    request: RequestEmailVerification,
  ): Promise<void> {
    // Validate request data
    const validatedRequest = RequestEmailVerificationSchema.parse(request)

    // Make API call
    await postVoid('v1/auth/verify-email/request', {
      json: validatedRequest,
    })
  },

  /**
   * Verifies user email using a verification token
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async verifyEmail(request: VerifyEmailRequest): Promise<void> {
    // Validate request data
    const validatedRequest = VerifyEmailRequestSchema.parse(request)

    // Make API call
    await postVoid('v1/auth/verify-email', { json: validatedRequest })
  },

  /**
   * Convenience method to logout using the current stored refresh token
   * @returns void (204 No Content)
   * @throws Error if no refresh token is available
   */
  async logoutCurrent(): Promise<void> {
    const refreshToken = getRefreshToken()
    if (!refreshToken) {
      throw new Error('No refresh token available')
    }
    return this.logout({ refreshToken })
  },

  /**
   * Convenience method to logout all devices using the current stored refresh token
   * @returns void (204 No Content)
   * @throws Error if no refresh token is available
   */
  async logoutAllCurrent(): Promise<void> {
    const refreshToken = getRefreshToken()
    if (!refreshToken) {
      throw new Error('No refresh token available')
    }
    return this.logoutAll({ refreshToken })
  },
}

/**
 * Mutation Hooks
 *
 * UI code must consume auth endpoints via TanStack Query mutations.
 */

export function useLoginMutation(
  options?: UseMutationOptions<LoginResponse, Error, LoginRequest>,
) {
  return useMutation({
    mutationFn: (credentials: LoginRequest) => authApi.login(credentials),
    ...options,
  })
}

export function useForgotPasswordMutation(
  options?: UseMutationOptions<void, Error, ForgotPasswordRequest>,
) {
  return useMutation({
    mutationFn: (request: ForgotPasswordRequest) =>
      authApi.forgotPassword(request),
    ...options,
  })
}

export function useResetPasswordMutation(
  options?: UseMutationOptions<void, Error, ResetPasswordRequest>,
) {
  return useMutation({
    mutationFn: (request: ResetPasswordRequest) =>
      authApi.resetPassword(request),
    ...options,
  })
}

export function useLoginWithGoogleMutation(
  options?: UseMutationOptions<LoginResponse, Error, GoogleLoginRequest>,
) {
  return useMutation({
    mutationFn: (request: GoogleLoginRequest) =>
      authApi.loginWithGoogle(request),
    ...options,
  })
}

export function useGoogleAuthUrlMutation(
  options?: UseMutationOptions<{ url: string }, Error, void>,
) {
  return useMutation({
    mutationFn: () => authApi.getGoogleAuthUrl(),
    ...options,
  })
}
