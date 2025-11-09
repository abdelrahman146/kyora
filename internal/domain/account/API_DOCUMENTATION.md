# Account Domain HTTP API Documentation

This document provides a comprehensive overview of all HTTP endpoints available in the account domain.

## Table of Contents

- [Authentication Endpoints](#authentication-endpoints)
- [User Profile Endpoints](#user-profile-endpoints)
- [Workspace Endpoints](#workspace-endpoints)
- [Invitation Endpoints](#invitation-endpoints)
- [User Management Endpoints](#user-management-endpoints)

---

## Authentication Endpoints

All authentication endpoints are **public** (no authentication required).

### 1. Login with Email and Password

**Endpoint:** `POST /v1/auth/login`

**Description:** Authenticates a user with email and password credentials.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": "usr_...",
    "workspaceId": "wrk_...",
    "role": "admin",
    "firstName": "John",
    "lastName": "Doe",
    "email": "user@example.com",
    "isEmailVerified": true
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Features:**
- Sends login notification email with client IP and user agent
- Returns JWT token for subsequent authenticated requests

---

### 2. Get Google OAuth URL

**Endpoint:** `GET /v1/auth/google/url`

**Description:** Returns the Google OAuth authorization URL for user authentication.

**Response:** `200 OK`
```json
{
  "url": "https://accounts.google.com/o/oauth2/auth?...",
  "state": "random_state_string"
}
```

---

### 3. Login with Google OAuth

**Endpoint:** `POST /v1/auth/google/login`

**Description:** Authenticates a user using Google OAuth authorization code.

**Request Body:**
```json
{
  "code": "google_auth_code"
}
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": "usr_...",
    "workspaceId": "wrk_...",
    "role": "admin",
    "firstName": "John",
    "lastName": "Doe",
    "email": "user@example.com",
    "isEmailVerified": true
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

### 4. Forgot Password

**Endpoint:** `POST /v1/auth/forgot-password`

**Description:** Initiates password reset process by sending a reset email.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:** `204 No Content`

**Note:** Always returns success to prevent email enumeration.

---

### 5. Reset Password

**Endpoint:** `POST /v1/auth/reset-password`

**Description:** Resets user password using a valid reset token.

**Request Body:**
```json
{
  "token": "reset_token_from_email",
  "newPassword": "newSecurePassword123"
}
```

**Response:** `204 No Content`

**Features:**
- Sends password reset confirmation email
- Token is consumed after successful reset

---

### 6. Request Email Verification

**Endpoint:** `POST /v1/auth/verify-email/request`

**Description:** Sends an email verification link to the user.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:** `204 No Content`

**Note:** Always returns success to prevent email enumeration.

---

### 7. Verify Email

**Endpoint:** `POST /v1/auth/verify-email`

**Description:** Verifies user's email address using a verification token.

**Request Body:**
```json
{
  "token": "verification_token_from_email"
}
```

**Response:** `204 No Content`

---

## User Profile Endpoints

All user profile endpoints require authentication.

### 8. Get Current User

**Endpoint:** `GET /v1/users/me`

**Description:** Returns the profile of the currently authenticated user.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Response:** `200 OK`
```json
{
  "id": "usr_...",
  "workspaceId": "wrk_...",
  "role": "admin",
  "firstName": "John",
  "lastName": "Doe",
  "email": "user@example.com",
  "isEmailVerified": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

---

### 9. Update Current User

**Endpoint:** `PATCH /v1/users/me`

**Description:** Updates the profile of the currently authenticated user.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Request Body:**
```json
{
  "firstName": "Jane",
  "lastName": "Smith"
}
```

**Response:** `200 OK`
```json
{
  "id": "usr_...",
  "workspaceId": "wrk_...",
  "role": "admin",
  "firstName": "Jane",
  "lastName": "Smith",
  "email": "user@example.com",
  "isEmailVerified": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-02T00:00:00Z"
}
```

---

## Workspace Endpoints

All workspace endpoints require authentication.

### 10. Get Current Workspace

**Endpoint:** `GET /v1/workspaces/me`

**Description:** Returns the workspace of the currently authenticated user.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Response:** `200 OK`
```json
{
  "id": "wrk_...",
  "ownerId": "usr_...",
  "stripeCustomerId": "cus_...",
  "stripePaymentMethodId": "pm_...",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

---

### 11. Get Workspace Users

**Endpoint:** `GET /v1/workspaces/users`

**Description:** Returns all users in the workspace.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Permissions:** Requires `view` permission for `account` resource.

**Response:** `200 OK`
```json
[
  {
    "id": "usr_...",
    "workspaceId": "wrk_...",
    "role": "admin",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "isEmailVerified": true
  },
  {
    "id": "usr_...",
    "workspaceId": "wrk_...",
    "role": "user",
    "firstName": "Jane",
    "lastName": "Smith",
    "email": "jane@example.com",
    "isEmailVerified": true
  }
]
```

---

### 12. Get Workspace User

**Endpoint:** `GET /v1/workspaces/users/{userId}`

**Description:** Returns a specific user in the workspace by ID.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Permissions:** Requires `view` permission for `account` resource.

**URL Parameters:**
- `userId` (required): The ID of the user to retrieve

**Response:** `200 OK`
```json
{
  "id": "usr_...",
  "workspaceId": "wrk_...",
  "role": "user",
  "firstName": "Jane",
  "lastName": "Smith",
  "email": "jane@example.com",
  "isEmailVerified": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

---

## Invitation Endpoints

### 13. Accept Invitation

**Endpoint:** `POST /v1/invitations/accept`

**Description:** Accepts a workspace invitation and creates a new user account.

**Query Parameters:**
- `token` (required): The invitation token from the email

**Request Body:**
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": "usr_...",
    "workspaceId": "wrk_...",
    "role": "user",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "isEmailVerified": true
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Features:**
- Auto-verifies email for invited users
- Returns JWT token for auto-login

---

### 14. Accept Invitation with Google

**Endpoint:** `GET /v1/invitations/accept/google`

**Description:** Accepts a workspace invitation using Google OAuth.

**Query Parameters:**
- `token` (required): The invitation token from the email
- `code` (required): The Google OAuth authorization code

**Response:** `200 OK`
```json
{
  "user": {
    "id": "usr_...",
    "workspaceId": "wrk_...",
    "role": "user",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "isEmailVerified": true
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

### 15. Invite User to Workspace

**Endpoint:** `POST /v1/workspaces/invitations`

**Description:** Sends an invitation email to a user to join the workspace.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Permissions:** 
- Requires `manage` permission for `account` resource
- Requires active subscription
- Subject to plan limits for team members

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "role": "user"
}
```

**Response:** `201 Created`
```json
{
  "id": "inv_...",
  "workspaceId": "wrk_...",
  "email": "newuser@example.com",
  "role": "user",
  "inviterId": "usr_...",
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

---

### 16. Get Workspace Invitations

**Endpoint:** `GET /v1/workspaces/invitations`

**Description:** Returns all invitations for the workspace.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Permissions:** Requires `manage` permission for `account` resource.

**Query Parameters:**
- `status` (optional): Filter by invitation status (`pending`, `accepted`, `expired`, `revoked`)

**Response:** `200 OK`
```json
[
  {
    "id": "inv_...",
    "workspaceId": "wrk_...",
    "email": "newuser@example.com",
    "role": "user",
    "inviterId": "usr_...",
    "status": "pending",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

---

### 17. Revoke Invitation

**Endpoint:** `DELETE /v1/workspaces/invitations/{invitationId}`

**Description:** Revokes a pending workspace invitation.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Permissions:** Requires `manage` permission for `account` resource.

**URL Parameters:**
- `invitationId` (required): The ID of the invitation to revoke

**Response:** `204 No Content`

---

## User Management Endpoints

All user management endpoints require authentication and admin permissions.

### 18. Update User Role

**Endpoint:** `PATCH /v1/workspaces/users/{userId}/role`

**Description:** Updates the role of a user within the workspace.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Permissions:** Requires `manage` permission for `account` resource.

**URL Parameters:**
- `userId` (required): The ID of the user to update

**Request Body:**
```json
{
  "role": "admin"
}
```

**Response:** `200 OK`
```json
{
  "id": "usr_...",
  "workspaceId": "wrk_...",
  "role": "admin",
  "firstName": "Jane",
  "lastName": "Smith",
  "email": "jane@example.com",
  "isEmailVerified": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-02T00:00:00Z"
}
```

**Restrictions:**
- Cannot update your own role
- Cannot update workspace owner's role

---

### 19. Remove User from Workspace

**Endpoint:** `DELETE /v1/workspaces/users/{userId}`

**Description:** Removes a user from the workspace (soft delete).

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Permissions:** Requires `manage` permission for `account` resource.

**URL Parameters:**
- `userId` (required): The ID of the user to remove

**Response:** `204 No Content`

**Restrictions:**
- Cannot remove yourself
- Cannot remove workspace owner

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "type": "https://tools.ietf.org/html/rfc7231#section-6.5.1",
  "title": "Bad Request",
  "status": 400,
  "detail": "validation failed",
  "instance": "/v1/auth/login"
}
```

### 401 Unauthorized
```json
{
  "type": "https://tools.ietf.org/html/rfc7235#section-3.1",
  "title": "Unauthorized",
  "status": 401,
  "detail": "invalid email or password",
  "instance": "/v1/auth/login"
}
```

### 403 Forbidden
```json
{
  "type": "https://tools.ietf.org/html/rfc7231#section-6.5.3",
  "title": "Forbidden",
  "status": 403,
  "detail": "you cannot update your own role",
  "instance": "/v1/workspaces/users/usr_123/role"
}
```

### 404 Not Found
```json
{
  "type": "https://tools.ietf.org/html/rfc7231#section-6.5.4",
  "title": "Not Found",
  "status": 404,
  "detail": "user is not a member of this workspace",
  "instance": "/v1/workspaces/users/usr_123"
}
```

### 409 Conflict
```json
{
  "type": "https://tools.ietf.org/html/rfc7231#section-6.5.8",
  "title": "Conflict",
  "status": 409,
  "detail": "user with this email already exists",
  "instance": "/v1/invitations/accept"
}
```

### 500 Internal Server Error
```json
{
  "type": "https://tools.ietf.org/html/rfc7231#section-6.6.1",
  "title": "Internal Server Error",
  "status": 500,
  "detail": "an unexpected error occurred",
  "instance": "/v1/auth/login"
}
```

---

## Security Notes

1. **JWT Authentication**: All protected endpoints require a valid JWT token in the Authorization header
2. **Workspace Isolation**: All data is strictly scoped by workspace ID to ensure multi-tenancy
3. **Role-Based Permissions**: Different operations require different permission levels
4. **Email Enumeration Protection**: Auth endpoints return success even for non-existent emails
5. **Token Expiration**: All tokens (password reset, email verification, invitations) have TTL
6. **Login Notifications**: Security emails are sent for successful logins

---

## Implementation Notes

- All endpoints follow REST conventions
- All responses use consistent JSON format
- All errors follow RFC 7807 Problem Details specification
- All timestamps are in ISO 8601 format (UTC)
- workspaceId is never exposed in URLs - it's always derived from the authenticated user
- All handlers include proper OpenAPI documentation comments
