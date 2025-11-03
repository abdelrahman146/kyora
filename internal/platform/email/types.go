package email

import (
	"context"
	"time"
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
