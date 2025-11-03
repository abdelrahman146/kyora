package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ResendClient sends emails via Resend HTTP API
// Docs reference: https://resend.com/docs/api-reference/emails/send-email
// Only the basic fields are mapped; unsupported fields can be added as needed.

type ResendClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type resendSendRequest struct {
	From        string             `json:"from"`
	To          []string           `json:"to"`
	Cc          []string           `json:"cc,omitempty"`
	Bcc         []string           `json:"bcc,omitempty"`
	Subject     string             `json:"subject"`
	HTML        string             `json:"html,omitempty"`
	Text        string             `json:"text,omitempty"`
	ReplyTo     []string           `json:"reply_to,omitempty"`
	Headers     map[string]string  `json:"headers,omitempty"`
	Attachments []resendAttachment `json:"attachments,omitempty"`
}

type resendAttachment struct {
	Filename  string `json:"filename"`
	Content   string `json:"content,omitempty"`
	Path      string `json:"path,omitempty"`
	ContentID string `json:"content_id,omitempty"`
}

type resendSendResponse struct {
	ID string `json:"id"`
}

func (c *ResendClient) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if msg == nil {
		return nil, fmt.Errorf("email message is nil")
	}

	// build payload
	reqPayload := resendSendRequest{
		From:    msg.From,
		To:      msg.To,
		Cc:      msg.Cc,
		Bcc:     msg.Bcc,
		Subject: msg.Subject,
		HTML:    msg.HTML,
		Text:    msg.Text,
		ReplyTo: msg.ReplyTo,
		Headers: msg.Headers,
	}
	if len(msg.Attachments) > 0 {
		reqPayload.Attachments = make([]resendAttachment, 0, len(msg.Attachments))
		for _, a := range msg.Attachments {
			reqPayload.Attachments = append(reqPayload.Attachments, resendAttachment{
				Filename:  a.Filename,
				Content:   a.Content,
				Path:      a.Path,
				ContentID: a.ContentID,
			})
		}
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	url := c.baseURL + "/emails"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("resend send failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var rr resendSendResponse
	if err := json.Unmarshal(respBody, &rr); err != nil {
		// tolerate schema change: return raw body
		return &SendResult{Provider: "resend", CreatedAt: time.Now(), Raw: string(respBody)}, nil
	}

	return &SendResult{
		Provider:  "resend",
		ID:        rr.ID,
		CreatedAt: time.Now(),
		Raw:       json.RawMessage(respBody),
	}, nil
}

// SendTemplate renders an embedded template then sends via Resend
func (c *ResendClient) SendTemplate(ctx context.Context, id TemplateID, to []string, from string, subject string, data map[string]any) (*SendResult, error) {
	html, err := RenderTemplate(id, data)
	if err != nil {
		return nil, err
	}
	if subject == "" {
		subject = SubjectFor(id)
	}
	return c.Send(ctx, &Message{
		From:    from,
		To:      to,
		Subject: subject,
		HTML:    html,
	})
}
