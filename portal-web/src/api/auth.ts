import apiClient, { setTokens, clearTokens, getRefreshToken } from "./client";
import {
  LoginRequestSchema,
  LoginResponseSchema,
  RefreshRequestSchema,
  RefreshResponseSchema,
  LogoutRequestSchema,
  LogoutAllRequestSchema,
  ForgotPasswordRequestSchema,
  ResetPasswordRequestSchema,
  GoogleLoginRequestSchema,
  RequestEmailVerificationSchema,
  VerifyEmailRequestSchema,
  type LoginRequest,
  type LoginResponse,
  type RefreshRequest,
  type RefreshResponse,
  type LogoutRequest,
  type LogoutAllRequest,
  type ForgotPasswordRequest,
  type ResetPasswordRequest,
  type GoogleLoginRequest,
  type RequestEmailVerification,
  type VerifyEmailRequest,
} from "./types/auth";

/**
 * Authentication API Service
 *
 * All methods validate request/response data using Zod schemas
 * and provide type-safe interfaces to backend auth endpoints.
 */

export const authApi = {
  // ==========================================================================
  // Login with Email & Password - POST /v1/auth/login
  // ==========================================================================

  /**
   * Authenticates a user with email and password
   * @returns LoginResponse with access token, refresh token, and user data
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    // Validate request data
    const validatedRequest = LoginRequestSchema.parse(credentials);

    // Make API call
    const response = await apiClient
      .post("v1/auth/login", { json: validatedRequest })
      .json();

    // Validate response data
    const validatedResponse = LoginResponseSchema.parse(response);

    // Store tokens in memory
    setTokens(validatedResponse.token, validatedResponse.refreshToken);

    return validatedResponse;
  },

  // ==========================================================================
  // Refresh Access Token - POST /v1/auth/refresh
  // ==========================================================================

  /**
   * Refreshes the access token using a refresh token
   * @returns RefreshResponse with new access token and refresh token
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async refreshToken(request: RefreshRequest): Promise<RefreshResponse> {
    // Validate request data
    const validatedRequest = RefreshRequestSchema.parse(request);

    // Make API call
    const response = await apiClient
      .post("v1/auth/refresh", { json: validatedRequest })
      .json();

    // Validate response data
    const validatedResponse = RefreshResponseSchema.parse(response);

    // Update tokens in memory
    setTokens(validatedResponse.token, validatedResponse.refreshToken);

    return validatedResponse;
  },

  // ==========================================================================
  // Logout - POST /v1/auth/logout
  // ==========================================================================

  /**
   * Logs out the current user by revoking the refresh token
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async logout(request: LogoutRequest): Promise<void> {
    // Validate request data
    const validatedRequest = LogoutRequestSchema.parse(request);

    // Make API call
    await apiClient.post("v1/auth/logout", { json: validatedRequest });

    // Clear tokens from memory
    clearTokens();
  },

  // ==========================================================================
  // Logout All Devices - POST /v1/auth/logout-all
  // ==========================================================================

  /**
   * Logs out the user from all devices by revoking all refresh tokens
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async logoutAll(request: LogoutAllRequest): Promise<void> {
    // Validate request data
    const validatedRequest = LogoutAllRequestSchema.parse(request);

    // Make API call
    await apiClient.post("v1/auth/logout-all", { json: validatedRequest });

    // Clear tokens from memory
    clearTokens();
  },

  // ==========================================================================
  // Forgot Password - POST /v1/auth/forgot-password
  // ==========================================================================

  /**
   * Sends a password reset email to the user
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async forgotPassword(request: ForgotPasswordRequest): Promise<void> {
    // Validate request data
    const validatedRequest = ForgotPasswordRequestSchema.parse(request);

    // Make API call
    await apiClient.post("v1/auth/forgot-password", {
      json: validatedRequest,
    });
  },

  // ==========================================================================
  // Reset Password - POST /v1/auth/reset-password
  // ==========================================================================

  /**
   * Resets user password using a valid reset token
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async resetPassword(request: ResetPasswordRequest): Promise<void> {
    // Validate request data
    const validatedRequest = ResetPasswordRequestSchema.parse(request);

    // Make API call
    await apiClient.post("v1/auth/reset-password", {
      json: validatedRequest,
    });
  },

  // ==========================================================================
  // Google OAuth Login - POST /v1/auth/google/login
  // ==========================================================================

  /**
   * Authenticates a user using Google OAuth code
   * @returns LoginResponse with access token, refresh token, and user data
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async loginWithGoogle(request: GoogleLoginRequest): Promise<LoginResponse> {
    // Validate request data
    const validatedRequest = GoogleLoginRequestSchema.parse(request);

    // Make API call
    const response = await apiClient
      .post("v1/auth/google/login", { json: validatedRequest })
      .json();

    // Validate response data
    const validatedResponse = LoginResponseSchema.parse(response);

    // Store tokens in memory
    setTokens(validatedResponse.token, validatedResponse.refreshToken);

    return validatedResponse;
  },

  // ==========================================================================
  // Get Google OAuth URL - GET /v1/auth/google/url
  // ==========================================================================

  /**
   * Gets the Google OAuth authorization URL for user authentication
   * @returns Object containing the authorization URL
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async getGoogleAuthUrl(): Promise<{ url: string }> {
    // Make API call
    const response = await apiClient.get("v1/auth/google/url").json();

    // Validate response structure
    if (
      typeof response !== "object" ||
      response === null ||
      !("url" in response)
    ) {
      throw new Error("Invalid response from Google OAuth URL endpoint");
    }

    return response as { url: string };
  },

  // ==========================================================================
  // Request Email Verification - POST /v1/auth/request-email-verification
  // ==========================================================================

  /**
   * Sends an email verification link to the user
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async requestEmailVerification(
    request: RequestEmailVerification
  ): Promise<void> {
    // Validate request data
    const validatedRequest = RequestEmailVerificationSchema.parse(request);

    // Make API call
    await apiClient.post("v1/auth/request-email-verification", {
      json: validatedRequest,
    });
  },

  // ==========================================================================
  // Verify Email - POST /v1/auth/verify-email
  // ==========================================================================

  /**
   * Verifies user email using a verification token
   * @returns void (204 No Content)
   * @throws HTTPError with parsed ProblemDetails on failure
   */
  async verifyEmail(request: VerifyEmailRequest): Promise<void> {
    // Validate request data
    const validatedRequest = VerifyEmailRequestSchema.parse(request);

    // Make API call
    await apiClient.post("v1/auth/verify-email", {
      json: validatedRequest,
    });
  },

  // ==========================================================================
  // Helper: Logout with current refresh token
  // ==========================================================================

  /**
   * Convenience method to logout using the current stored refresh token
   * @returns void (204 No Content)
   * @throws Error if no refresh token is available
   */
  async logoutCurrent(): Promise<void> {
    const refreshToken = getRefreshToken();
    if (!refreshToken) {
      throw new Error("No refresh token available");
    }
    return this.logout({ refreshToken });
  },

  /**
   * Convenience method to logout all devices using the current stored refresh token
   * @returns void (204 No Content)
   * @throws Error if no refresh token is available
   */
  async logoutAllCurrent(): Promise<void> {
    const refreshToken = getRefreshToken();
    if (!refreshToken) {
      throw new Error("No refresh token available");
    }
    return this.logoutAll({ refreshToken });
  },
};
