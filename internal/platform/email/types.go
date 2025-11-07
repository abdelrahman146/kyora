package email

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/spf13/viper"
)

// Message represents an email to be sent via the email provider
// This intentionally mirrors common Resend fields while staying provider-agnostic.
type Message struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Cc          []string          `json:"cc,omitempty"`
	Bcc         []string          `json:"bcc,omitempty"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	ReplyTo     []string          `json:"reply_to,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
}

// Attachment supports either base64 content or a remote path (URL)
// For inline images, set ContentID and reference using cid:CONTENT_ID in HTML
// Only one of Content or Path should be set.
type Attachment struct {
	Filename  string `json:"filename"`
	Content   string `json:"content,omitempty"`    // base64-encoded
	Path      string `json:"path,omitempty"`       // remote URL
	ContentID string `json:"content_id,omitempty"` // for inline
}

// SendResult is a provider-agnostic response summary for a send operation.
type SendResult struct {
	Provider  string    `json:"provider"`
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	Raw       any       `json:"raw,omitempty"` // optional raw provider response
}

// Client describes an email client capability
// Implementations must be safe for concurrent use
// and should honor context cancellation.
type Client interface {
	Send(ctx context.Context, msg *Message) (*SendResult, error)
	// SendTemplate renders an embedded HTML template with data, then sends.
	// If subject is empty, a sensible default for the template id is used.
	SendTemplate(ctx context.Context, id TemplateID, to []string, from string, subject string, data map[string]any) (*SendResult, error)
}

// TemplateID identifies an embedded template.
type TemplateID string

// EmailInfo holds common, app-level email metadata and URLs used when building templates
// in domain notification integrations. Values are sourced from the config package
// and can be overridden via options when constructing.
type EmailInfo struct {
	FromEmail    string
	FromName     string
	SupportEmail string
	HelpURL      string
	BaseURL      string
	ProductName  string
}

// FormattedFrom returns a display-friendly from header, e.g. "Acme <noreply@acme.com>"
func (e EmailInfo) FormattedFrom() string {
	if e.FromName != "" {
		return e.FromName + " <" + e.FromEmail + ">"
	}
	return e.FromEmail
}

// EmailOption allows overriding specific fields when creating EmailInfo
type EmailOption func(*EmailInfo)

// WithFrom overrides from email and optionally name.
func WithFrom(emailAddr, name string) EmailOption {
	return func(info *EmailInfo) {
		if emailAddr != "" {
			info.FromEmail = emailAddr
		}
		if name != "" {
			info.FromName = name
		}
	}
}

// NewEmail constructs EmailInfo from configuration defaults with optional overrides.
// It reads typed keys from the config package to avoid ad-hoc strings.
func NewEmail(opts ...EmailOption) EmailInfo {
	// local import to avoid circulars in type file; kept minimal
	fromEmail := viper.GetString(config.EmailFromEmail)
	if fromEmail == "" {
		fromEmail = "no-reply@kyora.com"
	}
	fromName := viper.GetString(config.EmailFromName)
	if fromName == "" {
		fromName = viper.GetString(config.AppName)
		if fromName == "" {
			fromName = "Kyora"
		}
	}
	supportEmail := viper.GetString(config.EmailSupportEmail)
	if supportEmail == "" {
		supportEmail = "support@kyora.com"
	}
	helpURL := viper.GetString(config.EmailHelpURL)
	if helpURL == "" {
		helpURL = "https://help.kyora.com"
	}
	baseURL := viper.GetString(config.HTTPBaseURL)
	if baseURL == "" {
		baseURL = "https://app.kyora.com"
	}
	productName := viper.GetString(config.AppName)
	if productName == "" {
		productName = "Kyora"
	}

	info := EmailInfo{
		FromEmail:    fromEmail,
		FromName:     fromName,
		SupportEmail: supportEmail,
		HelpURL:      helpURL,
		BaseURL:      baseURL,
		ProductName:  productName,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&info)
		}
	}
	return info
}
