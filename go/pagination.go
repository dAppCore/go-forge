package forge

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"

	core "dappco.re/go"
)

const defaultPageSize = 50

// defaultPageLimit is retained for compatibility with existing call sites.
const defaultPageLimit = defaultPageSize

// ListOptions controls pagination.
//
// Usage:
//
//	opts := forge.ListOptions{Page: 1, Limit: 50}
//	_ = opts
type ListOptions struct {
	Page     int // 1-based page number
	PageSize int // items per page (default 50)
	// Limit is a compatibility alias for PageSize.
	Limit int
}

// String returns a safe summary of the pagination options.
//
// Usage:
//
//	_ = forge.DefaultList.String()
func (o ListOptions) String() string {
	pageSize := o.PageSize
	if pageSize == 0 {
		pageSize = o.Limit
	}
	return core.Concat(
		"forge.ListOptions{page=",
		strconv.Itoa(o.Page),
		", page_size=",
		strconv.Itoa(pageSize),
		"}",
	)
}

// GoString returns a safe Go-syntax summary of the pagination options.
//
// Usage:
//
//	_ = fmt.Sprintf("%#v", forge.DefaultList)
func (o ListOptions) GoString() string { return o.String() }

// DefaultList provides sensible default pagination.
//
// Usage:
//
//	page, err := forge.ListPage[types.Repository](ctx, client, path, nil, forge.DefaultList)
//	_ = page
var DefaultList = ListOptions{Page: 1, PageSize: defaultPageSize}

// PagedResult holds a single page of results with metadata.
//
// Usage:
//
//	page, err := forge.ListPage[types.Repository](ctx, client, path, nil, forge.DefaultList)
//	_ = page
type PagedResult[T any] struct {
	Items      []T
	TotalCount int
	Page       int
	HasMore    bool
}

// String returns a safe summary of a page of results.
//
// Usage:
//
//	page, _ := forge.ListPage[types.Repository](...)
//	_ = page.String()
func (r PagedResult[T]) String() string {
	items := 0
	if r.Items != nil {
		items = len(r.Items)
	}
	return core.Concat(
		"forge.PagedResult{items=",
		strconv.Itoa(items),
		", totalCount=",
		strconv.Itoa(r.TotalCount),
		", page=",
		strconv.Itoa(r.Page),
		", hasMore=",
		strconv.FormatBool(r.HasMore),
		"}",
	)
}

// GoString returns a safe Go-syntax summary of a page of results.
//
// Usage:
//
//	_ = fmt.Sprintf("%#v", page)
func (r PagedResult[T]) GoString() string { return r.String() }

// ListPage fetches a single page of results.
// Extra query params can be passed via the query map.
//
// Usage:
//
//	page, err := forge.ListPage[types.Repository](ctx, client, "/api/v1/user/repos", nil, forge.DefaultList)
//	_ = page
func ListPage[T any](ctx context.Context, c *Client, path string, query map[string]string, opts ListOptions) (*PagedResult[T], error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	pageSize := opts.PageSize
	if pageSize < 1 {
		pageSize = opts.Limit
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, core.E("ListPage", "forge: parse path", err)
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("limit", strconv.Itoa(pageSize))
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
		HasMore: (totalCount > 0 && (opts.Page-1)*pageSize+len(items) < totalCount) ||
			(totalCount == 0 && len(items) >= pageSize),
	}, nil
}

// ListAll fetches all pages of results.
//
// Usage:
//
//	items, err := forge.ListAll[types.Repository](ctx, client, "/api/v1/user/repos", nil)
//	_ = items
func ListAll[T any](ctx context.Context, c *Client, path string, query map[string]string) ([]T, error) {
	var all []T
	page := 1

	for {
		result, err := ListPage[T](ctx, c, path, query, ListOptions{Page: page, PageSize: defaultPageSize})
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
//
// Usage:
//
//	for item, err := range forge.ListIter[types.Repository](ctx, client, "/api/v1/user/repos", nil) {
//	    _, _ = item, err
//	}
func ListIter[T any](ctx context.Context, c *Client, path string, query map[string]string) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		page := 1
		for {
			result, err := ListPage[T](ctx, c, path, query, ListOptions{Page: page, PageSize: defaultPageSize})
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
