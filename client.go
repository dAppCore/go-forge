package forge

import (
	// Note: AX-6 — request cancellation and deadlines are structural to HTTP calls; no core context primitive exists.
	"context"
	// Note: goccy/go-json — faster drop-in encoding/json replacement; no core equivalent for JSON-decoder performance profile.
	json "github.com/goccy/go-json"
	// Note: AX-6 — multipart upload bodies require the standard multipart writer; core-io provides storage, not multipart encoding.
	"mime/multipart"
	// Note: AX-6 — this low-level Forgejo client owns the HTTP boundary; no core.Client equivalent exists.
	"net/http"
	// Note: AX-6 — stream types, multipart content piping, and bounded HTTP error reads require io; coreio Medium only replaces multipart buffering below.
	goio "io"

	core "dappco.re/go/core"
	coreio "dappco.re/go/io"
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

// Error returns the formatted Forge API error string.
//
// Usage:
//
//	err := (&forge.APIError{StatusCode: 404, Message: "not found", URL: "/api/v1/repos/x/y"}).Error()
func (e *APIError) Error() string {
	if e == nil {
		return "forge.APIError{<nil>}"
	}
	return core.Concat("forge: ", e.URL, " ", core.Itoa(e.StatusCode), ": ", e.Message)
}

// String returns a safe summary of the API error.
//
// Usage:
//
//	s := err.String()
func (e *APIError) String() string { return e.Error() }

// GoString returns a safe Go-syntax summary of the API error.
//
// Usage:
//
//	s := fmt.Sprintf("%#v", err)
func (e *APIError) GoString() string { return e.Error() }

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

// String returns a safe summary of the rate limit state.
//
// Usage:
//
//	rl := client.RateLimit()
//	_ = rl.String()
func (r RateLimit) String() string {
	return core.Concat(
		"forge.RateLimit{limit=",
		core.Itoa(r.Limit),
		", remaining=",
		core.Itoa(r.Remaining),
		", reset=",
		core.FormatInt(r.Reset, 10),
		"}",
	)
}

// GoString returns a safe Go-syntax summary of the rate limit state.
//
// Usage:
//
//	_ = fmt.Sprintf("%#v", client.RateLimit())
func (r RateLimit) GoString() string { return r.String() }

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

// BaseURL returns the configured Forgejo base URL.
//
// Usage:
//
//	baseURL := client.BaseURL()
func (c *Client) BaseURL() string {
	if c == nil {
		return ""
	}
	return c.baseURL
}

// RateLimit returns the last known rate limit information.
//
// Usage:
//
//	rl := client.RateLimit()
func (c *Client) RateLimit() RateLimit {
	if c == nil {
		return RateLimit{}
	}
	return c.rateLimit
}

// UserAgent returns the configured User-Agent header value.
//
// Usage:
//
//	ua := client.UserAgent()
func (c *Client) UserAgent() string {
	if c == nil {
		return ""
	}
	return c.userAgent
}

// HTTPClient returns the configured underlying HTTP client.
//
// Usage:
//
//	hc := client.HTTPClient()
func (c *Client) HTTPClient() *http.Client {
	if c == nil {
		return nil
	}
	return c.httpClient
}

// String returns a safe summary of the client configuration.
//
// Usage:
//
//	s := client.String()
func (c *Client) String() string {
	if c == nil {
		return "forge.Client{<nil>}"
	}
	tokenState := "unset"
	if c.HasToken() {
		tokenState = "set"
	}
	return core.Concat("forge.Client{baseURL=", core.Sprintf("%q", c.baseURL), ", token=", tokenState, ", userAgent=", core.Sprintf("%q", c.userAgent), "}")
}

// GoString returns a safe Go-syntax summary of the client configuration.
//
// Usage:
//
//	s := fmt.Sprintf("%#v", client)
func (c *Client) GoString() string { return c.String() }

// HasToken reports whether the client was configured with an API token.
//
// Usage:
//
//	if c.HasToken() {
//	    _ = "authenticated"
//	}
func (c *Client) HasToken() bool {
	if c == nil {
		return false
	}
	return c.token != ""
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
//
// Usage:
//
//	var out map[string]string
//	err := client.Get(ctx, "/api/v1/user", &out)
func (c *Client) Get(ctx context.Context, path string, out any) error {
	_, err := c.doJSON(ctx, http.MethodGet, path, nil, out)
	return err
}

// Post performs a POST request.
//
// Usage:
//
//	var out map[string]any
//	err := client.Post(ctx, "/api/v1/orgs/core/repos", body, &out)
func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	_, err := c.doJSON(ctx, http.MethodPost, path, body, out)
	return err
}

// Patch performs a PATCH request.
//
// Usage:
//
//	var out map[string]any
//	err := client.Patch(ctx, "/api/v1/repos/core/go-forge", body, &out)
func (c *Client) Patch(ctx context.Context, path string, body, out any) error {
	_, err := c.doJSON(ctx, http.MethodPatch, path, body, out)
	return err
}

// Put performs a PUT request.
//
// Usage:
//
//	var out map[string]any
//	err := client.Put(ctx, "/api/v1/repos/core/go-forge", body, &out)
func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	_, err := c.doJSON(ctx, http.MethodPut, path, body, out)
	return err
}

// Delete performs a DELETE request.
//
// Usage:
//
//	err := client.Delete(ctx, "/api/v1/repos/core/go-forge")
func (c *Client) Delete(ctx context.Context, path string) error {
	_, err := c.doJSON(ctx, http.MethodDelete, path, nil, nil)
	return err
}

// DeleteWithBody performs a DELETE request with a JSON body.
//
// Usage:
//
//	err := client.DeleteWithBody(ctx, "/api/v1/repos/core/go-forge/labels", body)
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any) error {
	_, err := c.doJSON(ctx, http.MethodDelete, path, body, nil)
	return err
}

// PostRaw performs a POST request with a JSON body and returns the raw
// response body as bytes instead of JSON-decoding. Useful for endpoints
// such as /markdown that return raw HTML text.
//
// Usage:
//
//	body, err := client.PostRaw(ctx, "/api/v1/markdown", payload)
func (c *Client) PostRaw(ctx context.Context, path string, body any) ([]byte, error) {
	return c.postRawJSON(ctx, path, body)
}

func (c *Client) postRawJSON(ctx context.Context, path string, body any) ([]byte, error) {
	url := c.baseURL + path

	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, core.E("Client.PostRaw", "forge: marshal body", err)
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, core.NewReader(string(data)))
	}
	if err != nil {
		return nil, core.E("Client.PostRaw", "forge: create request", err)
	}

	if auth := c.authorizationHeader(); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, core.E("Client.PostRaw", "forge: request POST "+path, err)
	}
	defer resp.Body.Close()

	c.updateRateLimit(resp)

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := readBody(resp.Body)
	if err != nil {
		return nil, core.E("Client.PostRaw", "forge: read response body", err)
	}

	return data, nil
}

func (c *Client) postRawText(ctx context.Context, path, body string) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, core.NewReader(body))
	if err != nil {
		return nil, core.E("Client.PostText", "forge: create request", err)
	}

	if auth := c.authorizationHeader(); auth != "" {
		req.Header.Set("Authorization", auth)
	}
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

	c.updateRateLimit(resp)

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := readBody(resp.Body)
	if err != nil {
		return nil, core.E("Client.PostText", "forge: read response body", err)
	}

	return data, nil
}

func (c *Client) postMultipartJSON(ctx context.Context, path string, query map[string]string, fields map[string]string, fieldName, fileName string, content goio.Reader, out any) error {
	target := c.baseURL + path
	parsedTarget := core.URLParse(target)
	if !parsedTarget.OK {
		var err error
		if parseErr, ok := parsedTarget.Value.(error); ok {
			err = parseErr
		}
		return core.E("Client.PostMultipart", "forge: parse url", err)
	}
	targetStringer, ok := parsedTarget.Value.(interface{ String() string })
	if !ok {
		return core.E("Client.PostMultipart", "forge: parse url", nil)
	}
	target = targetStringer.String()
	if len(query) > 0 {
		values := newQueryBuilder()
		for key, value := range query {
			values.Set(key, value)
		}
		if encoded := values.Encode(); encoded != "" {
			separator := "?"
			if core.Contains(target, "?") {
				separator = "&"
			}
			target = core.Concat(target, separator, encoded)
		}
	}

	body := coreio.NewMemoryMedium()
	bodyWriter, err := body.Create("multipart")
	if err != nil {
		return core.E("Client.PostMultipart", "forge: create multipart body", err)
	}
	bodyWriterOpen := true
	defer func() {
		if bodyWriterOpen {
			_ = bodyWriter.Close()
		}
	}()

	writer := multipart.NewWriter(bodyWriter)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return core.E("Client.PostMultipart", "forge: create multipart form field", err)
		}
	}
	if fieldName != "" {
		part, err := writer.CreateFormFile(fieldName, fileName)
		if err != nil {
			return core.E("Client.PostMultipart", "forge: create multipart form file", err)
		}
		if content != nil {
			if _, err := goio.Copy(part, content); err != nil {
				return core.E("Client.PostMultipart", "forge: write multipart form file", err)
			}
		}
	}
	if err := writer.Close(); err != nil {
		return core.E("Client.PostMultipart", "forge: close multipart writer", err)
	}
	if err := bodyWriter.Close(); err != nil {
		bodyWriterOpen = false
		return core.E("Client.PostMultipart", "forge: close multipart body", err)
	}
	bodyWriterOpen = false

	bodyReader, err := body.ReadStream("multipart")
	if err != nil {
		return core.E("Client.PostMultipart", "forge: read multipart body", err)
	}
	defer bodyReader.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bodyReader)
	if err != nil {
		return core.E("Client.PostMultipart", "forge: create request", err)
	}

	if auth := c.authorizationHeader(); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return core.E("Client.PostMultipart", "forge: request POST "+path, err)
	}
	defer resp.Body.Close()

	c.updateRateLimit(resp)

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
//
// Usage:
//
//	body, err := client.GetRaw(ctx, "/api/v1/signing-key.gpg")
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, core.E("Client.GetRaw", "forge: create request", err)
	}

	if auth := c.authorizationHeader(); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, core.E("Client.GetRaw", "forge: request GET "+path, err)
	}
	defer resp.Body.Close()

	c.updateRateLimit(resp)

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp, path)
	}

	data, err := readBody(resp.Body)
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

	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, core.E("Client.doJSON", "forge: marshal body", err)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, core.NewReader(string(data)))
	}
	if err != nil {
		return nil, core.E("Client.doJSON", "forge: create request", err)
	}

	if auth := c.authorizationHeader(); auth != "" {
		req.Header.Set("Authorization", auth)
	}
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
	data, _ := readBody(goio.LimitReader(resp.Body, 1024))
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

func readBody(reader any) ([]byte, error) {
	result := core.ReadAll(reader)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, err
		}
		return nil, core.E("Client.readBody", "forge: read body", nil)
	}

	body, ok := result.Value.(string)
	if !ok {
		return nil, core.E("Client.readBody", "forge: read body", nil)
	}
	return []byte(body), nil
}

func (c *Client) updateRateLimit(resp *http.Response) {
	if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
		c.rateLimit.Limit = parseRateLimitInt(limit)
	}
	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		c.rateLimit.Remaining = parseRateLimitInt(remaining)
	}
	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		c.rateLimit.Reset = parseRateLimitInt64(reset)
	}
}

func parseRateLimitInt(value string) int {
	result := core.Atoi(value)
	if !result.OK {
		return 0
	}
	if parsed, ok := result.Value.(int); ok {
		return parsed
	}
	return 0
}

func parseRateLimitInt64(value string) int64 {
	result := core.ParseInt(value, 10, 64)
	if !result.OK {
		return 0
	}
	if parsed, ok := result.Value.(int64); ok {
		return parsed
	}
	return 0
}

func (c *Client) authorizationHeader() string {
	if c == nil || c.token == "" {
		return ""
	}
	return "Bearer " + c.token
}
