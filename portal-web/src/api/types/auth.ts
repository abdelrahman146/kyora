import { z } from "zod";

/**
 * Zod schemas for authentication API based on backend swagger.json
 * All schemas validate API request/response data
 */

// Role & User Schemas

export const RoleSchema = z.enum(["user", "admin"]);
export type Role = z.infer<typeof RoleSchema>;

const NullableStringSchema = z
  .object({
    String: z.string().optional(),
    Valid: z.boolean().optional(),
  })
  .loose();

const WorkspaceSchema = z
  .object({
    // backend returns ids like "wrk..." not UUIDs
    id: z.string().min(1),
    ownerId: z.string().min(1).optional(),
    // name is not present in the login response today
    name: z.string().optional(),

    // timestamps are returned in PascalCase on some endpoints/structs
    createdAt: z.string().optional(),
    updatedAt: z.string().optional(),
    CreatedAt: z.string().optional(),
    UpdatedAt: z.string().optional(),

    stripeCustomerId: NullableStringSchema.optional(),
    stripePaymentMethodId: NullableStringSchema.optional(),
  })
  .loose();

export const UserSchema = z
  .object({
    // backend returns ids like "usr..." not UUIDs
    id: z.string().min(1),
    email: z.email(),
    firstName: z.string(),
    lastName: z.string(),
    isEmailVerified: z.boolean(),
    role: RoleSchema,
    workspaceId: z.string().min(1),

    // Allow both camelCase and PascalCase timestamp keys.
    createdAt: z.string().optional(),
    updatedAt: z.string().optional(),
    CreatedAt: z.string().optional(),
    UpdatedAt: z.string().optional(),

    // gorm.DeletedAt sometimes comes as null.
    deletedAt: z.unknown().optional().nullable(),
    DeletedAt: z.unknown().optional().nullable(),

    workspace: WorkspaceSchema.optional(),

    // gorm also returns these, ignore them
    ID: z.number().optional(),
  })
  .loose();

export type User = z.infer<typeof UserSchema>;

// Login Schemas - POST /v1/auth/login

export const LoginRequestSchema = z.object({
  email: z.email({ message: "Invalid email address" }),
  password: z.string().min(1, "Password is required"),
});

export type LoginRequest = z.infer<typeof LoginRequestSchema>;

export const LoginResponseSchema = z.object({
  token: z.string(),
  refreshToken: z.string(),
  user: UserSchema,
});

export type LoginResponse = z.infer<typeof LoginResponseSchema>;

// Refresh Token Schemas - POST /v1/auth/refresh

export const RefreshRequestSchema = z.object({
  refreshToken: z.string().min(1, "Refresh token is required"),
});

export type RefreshRequest = z.infer<typeof RefreshRequestSchema>;

export const RefreshResponseSchema = z.object({
  token: z.string(),
  refreshToken: z.string(),
});

export type RefreshResponse = z.infer<typeof RefreshResponseSchema>;

// Logout Schemas - POST /v1/auth/logout

export const LogoutRequestSchema = z.object({
  refreshToken: z.string().min(1, "Refresh token is required"),
});

export type LogoutRequest = z.infer<typeof LogoutRequestSchema>;

// Response is 204 No Content (void)

// Logout All Devices Schemas - POST /v1/auth/logout-all

export const LogoutAllRequestSchema = z.object({
  refreshToken: z.string().min(1, "Refresh token is required"),
});

export type LogoutAllRequest = z.infer<typeof LogoutAllRequestSchema>;

// Response is 204 No Content (void)

// ProblemDetails Schema - RFC 7807 (Error Response)

export const ProblemDetailsSchema = z.object({
  type: z.string().optional(),
  title: z.string().optional(),
  status: z.number().int().optional(),
  detail: z.string().optional(),
  instance: z.string().optional(),
  extensions: z.record(z.string(), z.unknown()).optional(),
});

export type ProblemDetails = z.infer<typeof ProblemDetailsSchema>;

// Forgot Password Schemas - POST /v1/auth/forgot-password

export const ForgotPasswordRequestSchema = z.object({
  email: z.email({ message: "Invalid email address" }),
});

export type ForgotPasswordRequest = z.infer<typeof ForgotPasswordRequestSchema>;

// Response is 204 No Content (void)

// Reset Password Schemas - POST /v1/auth/reset-password

export const ResetPasswordRequestSchema = z.object({
  token: z.string().min(1, "Reset token is required"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export type ResetPasswordRequest = z.infer<typeof ResetPasswordRequestSchema>;

// Response is 204 No Content (void)

// Google OAuth Schemas - POST /v1/auth/google/login

export const GoogleLoginRequestSchema = z.object({
  code: z.string().min(1, "Google OAuth code is required"),
});

export type GoogleLoginRequest = z.infer<typeof GoogleLoginRequestSchema>;

// Response uses same LoginResponseSchema

// Email Verification Schemas

export const RequestEmailVerificationSchema = z.object({
  email: z.email({ message: "Invalid email address" }),
});

export type RequestEmailVerification = z.infer<
  typeof RequestEmailVerificationSchema
>;

export const VerifyEmailRequestSchema = z.object({
  token: z.string().min(1, "Verification token is required"),
});

export type VerifyEmailRequest = z.infer<typeof VerifyEmailRequestSchema>;

// Response is 204 No Content (void)
