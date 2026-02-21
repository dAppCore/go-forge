package forge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListOptions controls pagination.
type ListOptions struct {
	Page  int // 1-based page number
	Limit int // items per page (default 50)
}

// DefaultList returns sensible default pagination.
var DefaultList = ListOptions{Page: 1, Limit: 50}

// PagedResult holds a single page of results with metadata.
type PagedResult[T any] struct {
	Items      []T
	TotalCount int
	Page       int
	HasMore    bool
}

// ListPage fetches a single page of results.
// Extra query params can be passed via the query map.
func ListPage[T any](ctx context.Context, c *Client, path string, query map[string]string, opts ListOptions) (*PagedResult[T], error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 {
		opts.Limit = 50
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("forge: parse url: %w", err)
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("limit", strconv.Itoa(opts.Limit))
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("forge: create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/json")
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

	var items []T
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("forge: decode response: %w", err)
	}

	totalCount, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))

	return &PagedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       opts.Page,
		HasMore:    len(items) >= opts.Limit && opts.Page*opts.Limit < totalCount,
	}, nil
}

// ListAll fetches all pages of results.
func ListAll[T any](ctx context.Context, c *Client, path string, query map[string]string) ([]T, error) {
	var all []T
	page := 1

	for {
		result, err := ListPage[T](ctx, c, path, query, ListOptions{Page: page, Limit: 50})
		if err != nil {
			return nil, err
		}
		all = append(all, result.Items...)
		if !result.HasMore {
			break
		}
		page++
	}

	return all, nil
}
