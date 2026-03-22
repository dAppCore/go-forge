package forge

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"

	coreerr "dappco.re/go/core/log"
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

	u, err := url.Parse(path)
	if err != nil {
		return nil, coreerr.E("ListPage", "forge: parse path", err)
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("limit", strconv.Itoa(opts.Limit))
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	var items []T
	resp, err := c.doJSON(ctx, http.MethodGet, u.String(), nil, &items)
	if err != nil {
		return nil, err
	}

	totalCount, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))

	return &PagedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       opts.Page,
		// If totalCount is provided, use it to determine if there are more items.
		// Otherwise, assume there are more if we got a full page.
		HasMore: (totalCount > 0 && (opts.Page-1)*opts.Limit+len(items) < totalCount) ||
			(totalCount == 0 && len(items) >= opts.Limit),
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

// ListIter returns an iterator over all resources across all pages.
func ListIter[T any](ctx context.Context, c *Client, path string, query map[string]string) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		page := 1
		for {
			result, err := ListPage[T](ctx, c, path, query, ListOptions{Page: page, Limit: 50})
			if err != nil {
				yield(*new(T), err)
				return
			}
			for _, item := range result.Items {
				if !yield(item, nil) {
					return
				}
			}
			if !result.HasMore {
				break
			}
			page++
		}
	}
}
