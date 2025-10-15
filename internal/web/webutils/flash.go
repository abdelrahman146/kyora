package webutils

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FlashMessage represents a toast message to show on the client.
// Type should be one of: success | error | info
type FlashMessage struct {
	Type      string `json:"type"`
	Text      string `json:"text"`
	TimeoutMs int    `json:"timeoutMs,omitempty"`
}

// TriggerFlash sends HTMX trigger header to push toast messages without a full reload.
// It uses HX-Trigger to dispatch a "flash" custom event on the client with the payload.
func TriggerFlash(c *gin.Context, msgs ...FlashMessage) {
	if len(msgs) == 0 {
		return
	}
	payload := map[string]any{
		"flash": map[string]any{
			"messages": msgs,
		},
	}
	b, _ := json.Marshal(payload)
	c.Writer.Header().Add("HX-Trigger", string(b))
}

// TriggerFlashAfterSwap behaves like TriggerFlash but fires after swap lifecycle.
func TriggerFlashAfterSwap(c *gin.Context, msgs ...FlashMessage) {
	if len(msgs) == 0 {
		return
	}
	payload := map[string]any{
		"flash": map[string]any{
			"messages": msgs,
		},
	}
	b, _ := json.Marshal(payload)
	c.Writer.Header().Add("HX-Trigger-After-Swap", string(b))
}

// RedirectWithFlash queues messages across redirect by sending HX-Redirect and a cookie-like session queue
// using a small HTML body with a script to sessionStorage. This avoids a full page reload for HTMX requests.
// For non-HTMX requests, it falls back to HTTP redirect and sets a short-lived body script.
func RedirectWithFlash(c *gin.Context, location string, msgs ...FlashMessage) {
	if len(msgs) > 0 {
		// Inject a tiny script in the body to queue messages into sessionStorage before redirect
		// Works for both HTMX (since we instruct redirect) and normal navigation as content is returned.
		queueJS := `
          <script>
            try {
              const k = 'kyora.flash.queue';
              const newMsgs = ` + string(mustJSON(msgs)) + `;
              const existing = JSON.parse(sessionStorage.getItem(k) || '[]');
              sessionStorage.setItem(k, JSON.stringify(existing.concat(newMsgs)));
            } catch(e) {}
          </script>`
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write([]byte(queueJS))
	}
	// If HTMX request, instruct client via HX-Redirect header; else do a normal HTTP 302
	if c.GetHeader("HX-Request") != "" {
		Redirect(c, location)
	} else {
		c.Redirect(http.StatusFound, location)
	}
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
