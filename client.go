package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError represents an error response from the Forgejo API.
type APIError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("forge: %s %d: %s", e.URL, e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 response.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// IsForbidden returns true if the error is a 403 response.
func IsForbidden(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden
}

// IsConflict returns true if the error is a 409 response.
func IsConflict(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusConflict
}

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// Client is a low-level HTTP client for the Forgejo API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new Forgejo API client.
func NewClient(url, token string, opts ...Option) *Client {
	c := &Client{
		baseURL:    strings.TrimRight(url, "/"),
		token:      token,
		httpClient: http.DefaultClient,
		userAgent:  "go-forge/0.1",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPost, path, body, out)
}

// Patch performs a PATCH request.
func (c *Client) Patch(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPatch, path, body, out)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPut, path, body, out)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// DeleteWithBody performs a DELETE request with a JSON body.
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any) error {
	return c.do(ctx, http.MethodDelete, path, body, nil)
}

// PostRaw performs a POST request with a JSON body and returns the raw
// response body as bytes instead of JSON-decoding. Useful for endpoints
// such as /markdown that return raw HTML text.
func (c *Client) PostRaw(ctx context.Context, path string, body any) ([]byte, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("forge: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("forge: create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forge: request POST %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("forge: read response body: %w", err)
	}

	return data, nil
}

// GetRaw performs a GET request and returns the raw response body as bytes
// instead of JSON-decoding. Useful for endpoints that return raw file content.
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("forge: create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forge: request GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("forge: read response body: %w", err)
	}

	return data, nil
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("forge: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("forge: create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("forge: request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.parseError(resp, path)
	}

	if out != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("forge: decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) parseError(resp *http.Response, path string) error {
	var errBody struct {
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&errBody)
	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    errBody.Message,
		URL:        path,
	}
}
