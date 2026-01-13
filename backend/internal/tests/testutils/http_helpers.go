package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// HTTPClient provides reusable HTTP request helpers for E2E tests
type HTTPClient struct {
	BaseURL string
	Client  *http.Client
}

// NewHTTPClient creates a new HTTP client helper
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		// E2E tests often assert redirect responses (e.g., invoice download).
		// Disable automatic redirect following to keep tests deterministic and
		// avoid making external HTTP requests (e.g., to Stripe hosted invoice pages).
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (h *HTTPClient) fullURL(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return h.BaseURL + path
}

// Get makes a GET request to the specified path
func (h *HTTPClient) Get(path string) (*http.Response, error) {
	return h.Client.Get(h.fullURL(path))
}

// Post makes a POST request with JSON payload
func (h *HTTPClient) Post(path string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", h.fullURL(path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return h.Client.Do(req)
}

// Put makes a PUT request with JSON payload
func (h *HTTPClient) Put(path string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", h.fullURL(path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return h.Client.Do(req)
}

// Patch makes a PATCH request with JSON payload
func (h *HTTPClient) Patch(path string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", h.fullURL(path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return h.Client.Do(req)
}

// Delete makes a DELETE request
func (h *HTTPClient) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", h.fullURL(path), nil)
	if err != nil {
		return nil, err
	}
	return h.Client.Do(req)
}

// Request makes a custom HTTP request with optional JSON body
func (h *HTTPClient) Request(method, path string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(method, h.fullURL(path), body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return h.Client.Do(req)
}

// AuthenticatedRequest makes an authenticated request with Authorization: Bearer.
func (h *HTTPClient) AuthenticatedRequest(method, path string, payload interface{}, token string) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(method, h.fullURL(path), body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return h.Client.Do(req)
}

// PostRaw makes a POST request with raw body bytes.
// Useful for webhook tests where payload and signature must match exactly.
func (h *HTTPClient) PostRaw(path string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", h.fullURL(path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return h.Client.Do(req)
}

// RequestRaw makes a request with raw body bytes.
// Useful for endpoints that accept non-JSON payloads (e.g., direct file uploads).
func (h *HTTPClient) RequestRaw(method, path string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, h.fullURL(path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return h.Client.Do(req)
}

// AuthenticatedRequestRaw makes an authenticated request with raw body bytes.
func (h *HTTPClient) AuthenticatedRequestRaw(method, path string, body []byte, headers map[string]string, token string) (*http.Response, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}
	return h.RequestRaw(method, path, body, headers)
}

// DecodeJSON decodes JSON response body into target
func DecodeJSON(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// ReadBody reads the entire response body as string
func ReadBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ProblemResponse represents an RFC7807 Problem JSON response
type ProblemResponse struct {
	Type       string                 `json:"type"`
	Title      string                 `json:"title"`
	Status     int                    `json:"status"`
	Detail     string                 `json:"detail"`
	Instance   string                 `json:"instance"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// DecodeErrorResponse decodes an error response and returns the Problem JSON.
// Note: This reads the response body. Make sure to call this before closing resp.Body.
func DecodeErrorResponse(resp *http.Response) (*ProblemResponse, error) {
	// Read all content first
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var problem ProblemResponse
	if err := json.Unmarshal(body, &problem); err != nil {
		return nil, err
	}
	return &problem, nil
}

// GetErrorCode extracts the error code from a Problem JSON response.
// Note: This reads the response body, so call it before closing resp.Body or
// reading the body elsewhere.
func GetErrorCode(resp *http.Response) (string, error) {
	// Read all content first so we don't consume the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Restore the body so others can still read it if needed
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	var problem ProblemResponse
	if err := json.Unmarshal(body, &problem); err != nil {
		return "", err
	}

	if problem.Extensions != nil {
		if code, ok := problem.Extensions["code"].(string); ok {
			return code, nil
		}
	}
	return "", nil
}
