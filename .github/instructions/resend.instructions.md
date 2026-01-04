---
description: Resend Email API — Transactional Emails for Kyora
applyTo: "backend/**"
---

# Resend Email Integration

**Purpose**: Send transactional emails (verification, password reset, invitations)  
**Scope**: Backend only (email sending)  
**Reference**: [backend-core.instructions.md](./backend-core.instructions.md) for platform patterns

## Kyora-Specific Implementation

**Platform**: `backend/internal/platform/email`  
**Current Usage**:

- Email verification
- Password reset
- Workspace invitations
- Welcome emails
- Subscription confirmations

**Key Files** (reference implementations):

- `client.go` — Email client interface
- `resend_client.go` — Resend API implementation
- `mock_client.go` — Mock for testing
- `templates/` — HTML email templates (embedded via `go:embed`)
- `README.md` — Template list and usage guide

## Configuration

**Config Keys**:

- `email.provider` — `resend` (production) or `mock` (dev/test)
- `email.resend.api_key` — Resend API key
- `email.resend.base_url` — Resend API base URL (default: `https://api.resend.com`)
- `email.from_email` — Sender email address
- `email.from_name` — Sender name
- `email.support_email` — Support contact email
- `email.help_url` — Help/documentation URL

**Initialization** (`internal/server/server.go`):

```go
emailClient := email.NewClient(
    cfg.GetString(config.EmailProvider),
    cfg.GetString(config.ResendAPIKey),
    cfg.GetString(config.ResendBaseURL),
)
```

## Sending Emails

**Template-Based Sending** (recommended):

```go
err := emailClient.SendTemplate(ctx,
    email.TemplateEmailVerification,
    user.Email,
    email.From{Email: cfg.FromEmail, Name: cfg.FromName},
    "Verify your email",
    map[string]interface{}{
        "Name": user.FirstName,
        "VerificationLink": verificationURL,
    },
)
```

**Custom HTML Sending**:

```go
err := emailClient.Send(ctx,
    user.Email,
    email.From{Email: cfg.FromEmail, Name: cfg.FromName},
    "Subject",
    htmlBody,
)
```

## Available Templates

**Template Constants** (defined in `email/templates.go`):

- `TemplateForgotPassword` — Password reset link
- `TemplateEmailVerification` — Email verification link
- `TemplateWelcome` — Welcome message after signup
- `TemplateWorkspaceInvitation` — Workspace invitation link
- `TemplateSubscriptionConfirmed` — Subscription activation
- `TemplateSubscriptionCancelled` — Subscription cancellation
- `TemplatePaymentFailed` — Payment failure notification

**Template Data** (common fields):

- `Name` — User first name
- `WorkspaceName` — Workspace name
- `Link` / `VerificationLink` / `ResetLink` / `InvitationLink` — Action URL
- `SupportEmail` — Support contact
- `HelpURL` — Documentation link

## Template Structure

**Files**: `backend/internal/platform/email/templates/*.html`  
**Embedded**: Templates embedded via `go:embed` directive  
**Variables**: Use `{{.VariableName}}` for data interpolation  
**Styling**: Inline CSS (email clients strip `<style>` tags)

**Example Template**:

```html
<!DOCTYPE html>
<html>
  <body style="font-family: Arial, sans-serif;">
    <h1>Hello {{.Name}}</h1>
    <p>Click below to verify your email:</p>
    <a href="{{.VerificationLink}}" style="color: #007bff;">Verify Email</a>
  </body>
</html>
```

## Event Bus Integration

**Email triggers** (via event bus):

```go
// Emit event to trigger email
bus.Emit(bus.VerifyEmailTopic, map[string]interface{}{
    "email": user.Email,
    "token": verificationToken,
})

// Listen for email events (in server setup)
bus.Listen(bus.VerifyEmailTopic, func(payload any) {
    data := payload.(map[string]interface{})
    // Send verification email
    emailClient.SendTemplate(ctx, email.TemplateEmailVerification, ...)
})
```

**Built-in Topics** (`backend/internal/platform/bus/events.go`):

- `VerifyEmailTopic` — Email verification request
- `ResetPasswordTopic` — Password reset request
- `WorkspaceInvitationTopic` — Workspace invitation sent

## Testing

**Mock Provider** (dev/test):

- Set `email.provider=mock` in config
- Emails logged to console instead of sent
- No API calls made
- Useful for E2E tests

**Example Mock Usage**:

```go
// In tests, use mock provider
cfg := &config.Config{
    EmailProvider: "mock",
}
emailClient := email.NewClient("mock", "", "")

// Emails will be logged, not sent
emailClient.SendTemplate(ctx, email.TemplateWelcome, user.Email, ...)
```

## Resend API Details

**Endpoint**: `POST https://api.resend.com/emails`  
**Authentication**: Bearer token (`Authorization: Bearer YOUR_API_KEY`)  
**Rate Limits**: Check Resend docs for current limits  
**Response**:

```json
{
  "id": "re_abc123",
  "from": "noreply@kyora.app",
  "to": "user@example.com",
  "created_at": "2026-01-04T..."
}
```

**Error Handling**:

- HTTP 400: Invalid request (check email format, template data)
- HTTP 401: Invalid API key
- HTTP 429: Rate limit exceeded
- HTTP 500: Resend server error

## Common Patterns

**Send Verification Email**:

```go
verificationToken := generateToken()
verificationURL := fmt.Sprintf("%s/verify-email?token=%s", baseURL, verificationToken)

err := emailClient.SendTemplate(ctx,
    email.TemplateEmailVerification,
    user.Email,
    email.From{Email: cfg.FromEmail, Name: cfg.FromName},
    "Verify your email - Kyora",
    map[string]interface{}{
        "Name": user.FirstName,
        "VerificationLink": verificationURL,
        "SupportEmail": cfg.SupportEmail,
    },
)
```

**Send Password Reset**:

```go
resetToken := generateToken()
resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, resetToken)

err := emailClient.SendTemplate(ctx,
    email.TemplateForgotPassword,
    user.Email,
    email.From{Email: cfg.FromEmail, Name: cfg.FromName},
    "Reset your password - Kyora",
    map[string]interface{}{
        "Name": user.FirstName,
        "ResetLink": resetURL,
        "SupportEmail": cfg.SupportEmail,
    },
)
```

**Send Workspace Invitation**:

```go
invitationURL := fmt.Sprintf("%s/accept-invitation?token=%s", baseURL, invitationToken)

err := emailClient.SendTemplate(ctx,
    email.TemplateWorkspaceInvitation,
    invitedEmail,
    email.From{Email: cfg.FromEmail, Name: cfg.FromName},
    fmt.Sprintf("%s invited you to join their workspace", inviter.Name),
    map[string]interface{}{
        "InviterName": inviter.FirstName,
        "WorkspaceName": workspace.Name,
        "InvitationLink": invitationURL,
        "SupportEmail": cfg.SupportEmail,
    },
)
```

## Security Best Practices

**API Key Management**:

- Store in environment variables or secrets manager
- Never commit API keys to git
- Rotate keys periodically
- Use different keys for dev/staging/production

**Email Content**:

- Sanitize user input before including in emails
- Use URL encoding for query parameters
- Validate email addresses before sending
- Include unsubscribe links for marketing emails

**Token Handling**:

- Use cryptographically secure token generation
- Set expiration times (e.g., 24h for verification, 1h for password reset)
- Invalidate tokens after use
- Store token hashes in database, not plain text

## Troubleshooting

**Email not received**:

- Check spam/junk folder
- Verify sender email is configured in Resend
- Check Resend Dashboard for delivery status
- Verify recipient email format is valid

**Template rendering errors**:

- Verify template file exists in `templates/` folder
- Check variable names match template placeholders
- Ensure all required fields are provided in data map

**Rate limit errors**:

- Implement exponential backoff retry logic
- Cache frequently sent emails
- Upgrade Resend plan if needed
- Use batch sending for multiple recipients

## Documentation References

**Resend Docs** (external):

- [Send Email API](https://resend.com/docs/api-reference/emails/send-email)
- [HTML Best Practices](https://resend.com/docs/knowledge-base/html-email-best-practices)
- [Rate Limits](https://resend.com/docs/knowledge-base/rate-limits)

**Kyora Internal**:

- [backend-core.instructions.md](./backend-core.instructions.md) — Platform patterns, error handling
- [backend-testing.instructions.md](./backend-testing.instructions.md) — Testing email integration
- `backend/internal/platform/email/README.md` — Complete template list

## Anti-Patterns

❌ Never send emails synchronously in HTTP handlers — use event bus  
❌ Never expose API keys in frontend code  
❌ Never send sensitive data in email subject lines  
❌ Never use production email provider in tests — use mock  
❌ Never skip email validation before sending  
❌ Never send emails without proper error handling
