package email

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/logger"
)

// MockClient prints email contents to console instead of sending
// It uses slog JSON logger, but formats the email in a clear multi-line block.
type MockClient struct{}

func (m *MockClient) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if msg == nil {
		return nil, fmt.Errorf("email message is nil")
	}
	lg := logger.FromContext(ctx)

	// build a pretty block
	var b strings.Builder
	b.WriteString("\n==================== MOCK EMAIL ====================\n")
	b.WriteString(fmt.Sprintf("Time:     %s\n", time.Now().Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("From:     %s\n", msg.From))
	b.WriteString(fmt.Sprintf("To:       %s\n", strings.Join(msg.To, ", ")))
	if len(msg.Cc) > 0 {
		b.WriteString(fmt.Sprintf("Cc:       %s\n", strings.Join(msg.Cc, ", ")))
	}
	if len(msg.Bcc) > 0 {
		b.WriteString(fmt.Sprintf("Bcc:      %s\n", strings.Join(msg.Bcc, ", ")))
	}
	if len(msg.ReplyTo) > 0 {
		b.WriteString(fmt.Sprintf("Reply-To: %s\n", strings.Join(msg.ReplyTo, ", ")))
	}
	b.WriteString(fmt.Sprintf("Subject:  %s\n", msg.Subject))
	if msg.Text != "" {
		b.WriteString("\n-- TEXT ------------------------------\n" + msg.Text + "\n")
	}
	if msg.HTML != "" {
		b.WriteString("\n-- HTML ------------------------------\n" + msg.HTML + "\n")
	}
	if len(msg.Attachments) > 0 {
		b.WriteString("\n-- ATTACHMENTS -----------------------\n")
		for i, a := range msg.Attachments {
			kind := "content(base64)"
			if a.Path != "" {
				kind = "path(url)"
			}
			cid := ""
			if a.ContentID != "" {
				cid = " (cid=" + a.ContentID + ")"
			}
			b.WriteString(fmt.Sprintf("%d) %s %s%s\n", i+1, a.Filename, kind, cid))
		}
	}
	if len(msg.Headers) > 0 {
		b.WriteString("\n-- HEADERS ---------------------------\n")
		for k, v := range msg.Headers {
			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}
	b.WriteString("====================================================\n")

	lg.Info("mock email send", "provider", "mock", "preview", b.String())

	return &SendResult{
		Provider:  "mock",
		ID:        fmt.Sprintf("mock-%d", time.Now().UnixNano()),
		CreatedAt: time.Now(),
		Raw:       map[string]any{"preview": b.String()},
	}, nil
}

func (m *MockClient) SendTemplate(ctx context.Context, id TemplateID, to []string, from string, subject string, data map[string]any) (*SendResult, error) {
	html, err := RenderTemplate(id, data)
	if err != nil {
		return nil, err
	}
	if subject == "" {
		subject = SubjectFor(id)
	}
	return m.Send(ctx, &Message{
		From:    from,
		To:      to,
		Subject: subject,
		HTML:    html,
	})
}
