package forge

import (
	"bytes"
	"context"
	json "github.com/goccy/go-json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	core "dappco.re/go/core"
)

// APIError represents an error response from the Forgejo API.
//
// Usage:
//
//	if apiErr, ok := err.(*forge.APIError); ok {
//	    _ = apiErr.StatusCode
//	}
type APIError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e *APIError) Error() string {
	return core.Concat("forge: ", e.URL, " ", strconv.Itoa(e.StatusCode), ": ", e.Message)
}

// IsNotFound returns true if the error is a 404 response.
//
// Usage:
//
//	if forge.IsNotFound(err) {
//	    return nil
//	}
func IsNotFound(err error) bool {
	var apiErr *APIError
	return core.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// IsForbidden returns true if the error is a 403 response.
//
// Usage:
//
//	if forge.IsForbidden(err) {
//	    return nil
//	}
func IsForbidden(err error) bool {
	var apiErr *APIError
	return core.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden
}

// IsConflict returns true if the error is a 409 response.
//
// Usage:
//
//	if forge.IsConflict(err) {
//	    return nil
//	}
func IsConflict(err error) bool {
	var apiErr *APIError
	return core.As(err, &apiErr) && apiErr.StatusCode == http.StatusConflict
}

// Option configures the Client.
//
// Usage:
//
//	opts := []forge.Option{forge.WithUserAgent("go-forge/1.0")}
type Option func(*Client)

// WithHTTPClient sets a custom http.Client.
//
// Usage:
//
//	c := forge.NewClient(url, token, forge.WithHTTPClient(http.DefaultClient))
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets the User-Agent header.
//
// Usage:
//
//	c := forge.NewClient(url, token, forge.WithUserAgent("go-forge/1.0"))
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// RateLimit represents the rate limit information from the Forgejo API.
//
// Usage:
//
//	rl := client.RateLimit()
//	_ = rl.Remaining
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int64
}

// Client is a low-level HTTP client for the Forgejo API.
//
// Usage:
//
//	c := forge.NewClient("https://forge.lthn.ai", "token")
//	_ = c
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
//
// Usage:
//
//	c := forge.NewClient("https://forge.lthn.ai", "token")
//	_ = c
func NewClient(url, token string, opts ...Option) *Client {
	c := &Client{
		baseURL: trimTrailingSlashes(url),
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
	return c.postRawJSON(ctx, path, body)
}

func (c *Client) postRawJSON(ctx context.Context, path string, body any) ([]byte, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, core.E("Client.PostRaw", "forge: marshal body", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, core.E("Client.PostRaw", "forge: create request", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, core.E("Client.PostRaw", "forge: request POST "+path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, core.E("Client.PostRaw", "forge: read response body", err)
	}

	return data, nil
}

func (c *Client) postRawText(ctx context.Context, path, body string) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, core.E("Client.PostText", "forge: create request", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Content-Type", "text/plain")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, core.E("Client.PostText", "forge: request POST "+path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, core.E("Client.PostText", "forge: read response body", err)
	}

	return data, nil
}

func (c *Client) postMultipartJSON(ctx context.Context, path string, query map[string]string, fieldName, fileName string, content io.Reader, out any) error {
	target, err := url.Parse(c.baseURL + path)
	if err != nil {
		return core.E("Client.PostMultipart", "forge: parse url", err)
	}
	if len(query) > 0 {
		values := target.Query()
		for key, value := range query {
			values.Set(key, value)
		}
		target.RawQuery = values.Encode()
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return core.E("Client.PostMultipart", "forge: create multipart form file", err)
	}
	if content != nil {
		if _, err := io.Copy(part, content); err != nil {
			return core.E("Client.PostMultipart", "forge: write multipart form file", err)
		}
	}
	if err := writer.Close(); err != nil {
		return core.E("Client.PostMultipart", "forge: close multipart writer", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target.String(), &body)
	if err != nil {
		return core.E("Client.PostMultipart", "forge: create request", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return core.E("Client.PostMultipart", "forge: request POST "+path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.parseError(resp, path)
	}

	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return core.E("Client.PostMultipart", "forge: decode response body", err)
	}
	return nil
}

// GetRaw performs a GET request and returns the raw response body as bytes
// instead of JSON-decoding. Useful for endpoints that return raw file content.
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, core.E("Client.GetRaw", "forge: create request", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, core.E("Client.GetRaw", "forge: request GET "+path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, core.E("Client.GetRaw", "forge: read response body", err)
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
			return nil, core.E("Client.doJSON", "forge: marshal body", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, core.E("Client.doJSON", "forge: create request", err)
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
		return nil, core.E("Client.doJSON", "forge: request "+method+" "+path, err)
	}
	defer resp.Body.Close()

	c.updateRateLimit(resp)

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	if out != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return nil, core.E("Client.doJSON", "forge: decode response", err)
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
