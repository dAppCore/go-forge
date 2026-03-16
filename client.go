package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	coreerr "forge.lthn.ai/core/go-log"
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

// RateLimit represents the rate limit information from the Forgejo API.
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int64
}

// Client is a low-level HTTP client for the Forgejo API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	userAgent  string
	rateLimit  RateLimit
}

// RateLimit returns the last known rate limit information.
func (c *Client) RateLimit() RateLimit {
	return c.rateLimit
}

// NewClient creates a new Forgejo API client.
func NewClient(url, token string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(url, "/"),
		token:   token,
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		userAgent: "go-forge/0.1",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, out any) error {
	_, err := c.doJSON(ctx, http.MethodGet, path, nil, out)
	return err
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	_, err := c.doJSON(ctx, http.MethodPost, path, body, out)
	return err
}

// Patch performs a PATCH request.
func (c *Client) Patch(ctx context.Context, path string, body, out any) error {
	_, err := c.doJSON(ctx, http.MethodPatch, path, body, out)
	return err
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	_, err := c.doJSON(ctx, http.MethodPut, path, body, out)
	return err
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) error {
	_, err := c.doJSON(ctx, http.MethodDelete, path, nil, nil)
	return err
}

// DeleteWithBody performs a DELETE request with a JSON body.
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any) error {
	_, err := c.doJSON(ctx, http.MethodDelete, path, body, nil)
	return err
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
			return nil, coreerr.E("Client.PostRaw", "forge: marshal body", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, coreerr.E("Client.PostRaw", "forge: create request", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, coreerr.E("Client.PostRaw", "forge: request POST "+path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, coreerr.E("Client.PostRaw", "forge: read response body", err)
	}

	return data, nil
}

// GetRaw performs a GET request and returns the raw response body as bytes
// instead of JSON-decoding. Useful for endpoints that return raw file content.
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, coreerr.E("Client.GetRaw", "forge: create request", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, coreerr.E("Client.GetRaw", "forge: request GET "+path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, coreerr.E("Client.GetRaw", "forge: read response body", err)
	}

	return data, nil
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	_, err := c.doJSON(ctx, method, path, body, out)
	return err
}

func (c *Client) doJSON(ctx context.Context, method, path string, body, out any) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, coreerr.E("Client.doJSON", "forge: marshal body", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, coreerr.E("Client.doJSON", "forge: create request", err)
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
		return nil, coreerr.E("Client.doJSON", "forge: request "+method+" "+path, err)
	}
	defer resp.Body.Close()

	c.updateRateLimit(resp)

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	if out != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return nil, coreerr.E("Client.doJSON", "forge: decode response", err)
		}
	}

	return resp, nil
}

func (c *Client) parseError(resp *http.Response, path string) error {
	var errBody struct {
		Message string `json:"message"`
	}

	// Read a bit of the body to see if we can get a message
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	_ = json.Unmarshal(data, &errBody)

	msg := errBody.Message
	if msg == "" && len(data) > 0 {
		msg = string(data)
	}
	if msg == "" {
		msg = http.StatusText(resp.StatusCode)
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    msg,
		URL:        path,
	}
}

func (c *Client) updateRateLimit(resp *http.Response) {
	if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
		c.rateLimit.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		c.rateLimit.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		c.rateLimit.Reset, _ = strconv.ParseInt(reset, 10, 64)
	}
}
