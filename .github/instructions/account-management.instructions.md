---
description: "Kyora account management SSOT (backend + portal-web): auth sessions, tokens, users/workspaces, invitations, RBAC + plan gates"
---

# Kyora Account Management — Single Source of Truth (SSOT)

This file documents **account/auth/workspace** behavior implemented today across:

- Backend (source of truth): `backend/internal/domain/account/**` + wiring in `backend/internal/server/routes.go`
- Portal Web (current consumer): `portal-web/src/api/auth.ts`, `portal-web/src/api/client.ts`, `portal-web/src/lib/auth.ts`, `portal-web/src/api/user.ts`, `portal-web/src/routes/auth/**`

If you change authentication/session behavior or route contracts, update backend + portal-web together.

## Non-negotiables

- **Workspace is the tenant for auth:** JWTs carry `workspaceId` and requests must not be able to access other workspaces.
- **Never trust workspace IDs from URL params** for these routes. Workspace is derived from the authenticated user.
- **No user enumeration:** forgot-password and request-email-verification return success even if the email does not exist.
- **Session security model:** access token is short-lived and stateless (JWT); refresh token is long-lived and stateful (DB session, revocable).
- **Token rotation:** refresh always revokes the old session token hash and issues a new one.
- **Auth invalidation:** on sensitive changes (password reset, logout-all), bump `User.AuthVersion` to invalidate all access tokens.

## Backend: route surface (authoritative)

### Public auth routes (no auth required)

All routes are under `/v1/auth`.

- `POST /login`
  - Body: `{ email, password }`
  - Returns: `LoginResponse { user, token, refreshToken }`
  - Rate limited (best-effort): 20 attempts / 10 minutes per `(email, ip)`.

- `POST /refresh`
  - Body: `{ refreshToken }`
  - Returns: `RefreshResponse { token, refreshToken }`
  - Refresh tokens are rotated: old refresh session is revoked first.

- `POST /logout`
  - Body: `{ refreshToken }`
  - Returns: `204`
  - Revokes exactly the provided refresh token (session).

- `POST /logout-all`
  - Body: `{ refreshToken }`
  - Returns: `204`
  - Revokes all refresh sessions for the user and increments `authVersion` (invalidates all access JWTs).

- `POST /logout-others`
  - Body: `{ refreshToken }`
  - Returns: `204`
  - Revokes all refresh sessions except the provided refresh token.

- Google OAuth
  - `GET /google/url` → returns `{ url, state }`
  - `POST /google/login` body: `{ code }` → returns `LoginResponse`
  - Current behavior: Google login succeeds only if a user already exists with that email.

- Password reset
  - `POST /forgot-password` body: `{ email }` → `204`
    - If email does not exist: still returns `204`.
    - If rate limited: returns `429`.
  - `POST /reset-password` body: `{ token, newPassword }` → `204`
    - On success: updates password, increments `authVersion`, revokes all sessions.

- Email verification
  - `POST /verify-email/request` body: `{ email }` → `204`
    - If email does not exist: still returns `204`.
    - If rate limited: returns `429`.
  - `POST /verify-email` body: `{ token }` → `204`

### Public invitation acceptance routes (no auth required)

All routes are under `/v1/invitations`.

- `POST /accept?token=...`
  - Body: `{ firstName, lastName, password }`
  - Behavior: creates a user in the invitation workspace, marks invitation accepted, consumes token, then issues tokens.
  - Returns: `LoginResponse`.

- `GET /accept/google?token=...&code=...`
  - Behavior: exchanges Google code, requires Google email to match invitation email, creates user (no password), marks invitation accepted, consumes token, issues tokens.
  - Returns: `LoginResponse`.

### Protected routes (auth required)

Middleware chain used for protected routes:

- `auth.EnforceAuthentication` (JWT)
- `account.EnforceValidActor(accountService)`
  - Loads the user from DB.
  - Rejects if JWT `authVersion` does not match DB `User.AuthVersion`.
- For workspace routes: `account.EnforceWorkspaceMembership(accountService)`
  - Loads workspace by the actor’s `WorkspaceID` (never from URL).

#### User profile

Under `/v1/users`:

- `GET /me` → returns the authenticated `User`.
- `PATCH /me` body: `{ firstName?, lastName? }` → returns updated `User`.

#### Workspace

Under `/v1/workspaces`:

- `GET /me` → returns the authenticated user’s `Workspace` (preloads users).

Workspace users (permission: `role.ActionView` on `role.ResourceAccount`):

- `GET /users` → returns `[]User` ordered by `created_at ASC`.
- `GET /users/:userId` → returns `User`.
  - Must be scoped to the actor’s workspace; probing other workspace user IDs returns 404.

Workspace user management (permission: `role.ActionManage` on `role.ResourceAccount`):

- `PATCH /users/:userId/role` body: `{ role: 'user'|'admin' }`
  - Cannot update your own role.
  - Cannot update the workspace owner’s role.

- `DELETE /users/:userId`
  - Cannot delete yourself.
  - Cannot delete the workspace owner.
  - Soft-deletes the user.

Workspace invitations management (permission: `role.ActionManage` on `role.ResourceAccount`):

- `POST /invitations`
  - Plan gates applied:
    - `billing.EnforceActiveSubscription`
    - `billing.EnforcePlanWorkspaceLimits(billing.PlanSchema.MaxTeamMembers, accountService.CountWorkspaceUsersForPlanLimit)`
  - Body: `{ email, role }` where role is `user|admin`.
  - Returns: `UserInvitation`.

- `GET /invitations?status=pending|accepted|expired|revoked` → returns `[]UserInvitation`.

- `DELETE /invitations/:invitationId` → `204`
  - Only pending invitations can be revoked.

## Backend: token/session storage semantics

- Access token: JWT created by `auth.NewJwtToken(userID, workspaceID, authVersion)`.
- Refresh token: random opaque token; only **hash** is stored.
  - Stored entity: `Session { userId, workspaceId, tokenHash, expiresAt, createdIP, userAgent }`.
  - TTL defaults to 30 days if not configured.
- Password reset, email verification, workspace invitation tokens:
  - Stored in cache with prefixes:
    - `pwreset:`, `emailverify:`, `invitation:`
  - Payloads include `expAt`, and tokens are consumed (deleted) after use.
- Abuse protection:
  - Login: cache-backed throttle per `(email, ip)`.
  - Forgot-password + verify-email/request: cache-backed throttle per `email`.

## Portal Web: expected client behavior

- Token storage:
  - Access token is in memory.
  - Refresh token is stored in a cookie (`kyora_refresh_token`).
- Request auth:
  - `Authorization: Bearer <accessToken>` is added by `portal-web/src/api/client.ts`.
- Refresh flow:
  - On any non-auth request returning `401`, portal attempts `POST /v1/auth/refresh` with `{ refreshToken }`.
  - If refresh succeeds, it retries the original request with the new access token.
  - If refresh fails, it clears tokens and navigates to `/auth/login`.
- Session restoration:
  - `restoreSession()` refreshes tokens, then calls `GET /v1/users/me`.

### Frontend features implemented today

- Login page uses `authApi.login`.
- Forgot password (`/auth/forgot-password`) uses `authApi.forgotPassword`.
- Reset password (`/auth/reset-password?token=...`) uses `authApi.resetPassword`.
- Current user profile fetch/update uses `userApi.getCurrentUser` + `userApi.updateCurrentUser`.

### Frontend features not implemented yet (but backend supports)

When building team/workspace management UI, follow backend contract above:

- Workspace info: `GET /v1/workspaces/me`.
- Workspace users list/details: `GET /v1/workspaces/users`, `GET /v1/workspaces/users/:userId`.
- Workspace invitation management: `POST/GET/DELETE /v1/workspaces/invitations`.
- Invitation acceptance pages for `POST /v1/invitations/accept?token=...` and `GET /v1/invitations/accept/google?...`.
- Logout other devices: `POST /v1/auth/logout-others`.
