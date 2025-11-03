# Email package

Minimal email abstraction with Resend integration and console mock.

## Features

- Provider factory (mock or Resend) via Viper config.
- Simple `Send` with provider-agnostic `Message`.
- Embedded HTML templates with `go:embed` and `SendTemplate`.
- Graceful handling of missing template data (renders empty or defaults via `default` func).

## Config keys

- `email.provider`: `resend` (default) or `mock`
- `email.mock.enabled`: `true` to force mock regardless of provider
- `email.resend.api_key`: Resend API key (required if provider is `resend`)
- `email.resend.base_url`: optional, defaults to `https://api.resend.com`

Flat env fallbacks are also read: `email_provider`, `email_mock_enabled`, `email_resend_api_key`, `email_resend_base_url`.

## Usage

```go
import (
    "context"
    "github.com/abdelrahman146/kyora/internal/platform/email"
)

func sendForgotPassword(ctx context.Context, to string, resetURL string, from string) error {
    c, err := email.New()
    if err != nil { return err }
    data := map[string]any{
        "userName":    "there",            // optional
        "productName": "Kyora",           // optional
        "resetURL":    resetURL,            // required
        "title":       "Reset your password", // optional override
    }
    _, err = c.SendTemplate(ctx, email.TemplateForgotPassword, []string{to}, from, "", data)
    return err
}
```

If `subject` is empty, a sensible default is used per template id.

### Available templates

- `TemplateForgotPassword` → `templates/forgot_password.html`

Template helpers:

- `default DEF VALUE` → returns `VALUE` if it is a non-empty string/bytes; otherwise `DEF`.

```html
<p>Hi {{default "there" .userName}},</p>
```
