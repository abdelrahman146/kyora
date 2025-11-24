package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
		Client:  &http.Client{},
	}
}

// Get makes a GET request to the specified path
func (h *HTTPClient) Get(path string) (*http.Response, error) {
	return h.Client.Get(h.BaseURL + path)
}

// Post makes a POST request with JSON payload
func (h *HTTPClient) Post(path string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", h.BaseURL+path, bytes.NewBuffer(body))
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
	req, err := http.NewRequest("PUT", h.BaseURL+path, bytes.NewBuffer(body))
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
	req, err := http.NewRequest("PATCH", h.BaseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return h.Client.Do(req)
}

// Delete makes a DELETE request
func (h *HTTPClient) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", h.BaseURL+path, nil)
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
	req, err := http.NewRequest(method, h.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return h.Client.Do(req)
}

// AuthenticatedRequest makes an authenticated request with JWT cookie
func (h *HTTPClient) AuthenticatedRequest(method, path string, payload interface{}, token string) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(method, h.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.AddCookie(&http.Cookie{Name: "jwt", Value: token})
	}
	return h.Client.Do(req)
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

// ExtractJWTCookie extracts JWT token from response cookies
func ExtractJWTCookie(resp *http.Response) string {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "jwt" {
			return cookie.Value
		}
	}
	return ""
}
