package forge

import (
	"context"
	"iter"
	"strconv"

	core "dappco.re/go/core"
)

// Resource provides generic CRUD operations for a Forgejo API resource.
// T is the resource type, C is the create options type, U is the update options type.
//
// Usage:
//
//	r := forge.NewResource[types.Issue, types.CreateIssueOption, types.EditIssueOption](client, "/api/v1/repos/{owner}/{repo}/issues/{index}")
//	_ = r
type Resource[T any, C any, U any] struct {
	client     *Client
	path       string // item path: /api/v1/repos/{owner}/{repo}/issues/{index}
	collection string // collection path: /api/v1/repos/{owner}/{repo}/issues
}

// String returns a safe summary of the resource configuration.
//
// Usage:
//
//	s := res.String()
func (r *Resource[T, C, U]) String() string {
	return core.Concat(
		"forge.Resource{path=",
		strconv.Quote(r.path),
		", collection=",
		strconv.Quote(r.collection),
		"}",
	)
}

// GoString returns a safe Go-syntax summary of the resource configuration.
//
// Usage:
//
//	s := fmt.Sprintf("%#v", res)
func (r *Resource[T, C, U]) GoString() string { return r.String() }

// NewResource creates a new Resource for the given path pattern.
// The path should be the item path (e.g., /repos/{owner}/{repo}/issues/{index}).
// The collection path is derived by stripping the last /{placeholder} segment.
//
// Usage:
//
//	r := forge.NewResource[types.Issue, types.CreateIssueOption, types.EditIssueOption](client, "/api/v1/repos/{owner}/{repo}/issues/{index}")
//	_ = r
func NewResource[T any, C any, U any](c *Client, path string) *Resource[T, C, U] {
	collection := path
	// Strip last segment if it's a pure placeholder like /{index}
	// Don't strip if mixed like /repos or /{org}/repos
	if i := lastIndexByte(path, '/'); i >= 0 {
		lastSeg := path[i+1:]
		if core.HasPrefix(lastSeg, "{") && core.HasSuffix(lastSeg, "}") {
			collection = path[:i]
		}
	}
	return &Resource[T, C, U]{client: c, path: path, collection: collection}
}

// List returns a single page of resources.
//
// Usage:
//
//	page, err := res.List(ctx, forge.Params{"owner": "core"}, forge.DefaultList)
func (r *Resource[T, C, U]) List(ctx context.Context, params Params, opts ListOptions) (*PagedResult[T], error) {
	return ListPage[T](ctx, r.client, ResolvePath(r.collection, params), nil, opts)
}

// ListAll returns all resources across all pages.
//
// Usage:
//
//	items, err := res.ListAll(ctx, forge.Params{"owner": "core"})
func (r *Resource[T, C, U]) ListAll(ctx context.Context, params Params) ([]T, error) {
	return ListAll[T](ctx, r.client, ResolvePath(r.collection, params), nil)
}

// Iter returns an iterator over all resources across all pages.
//
// Usage:
//
//	for item, err := range res.Iter(ctx, forge.Params{"owner": "core"}) {
//	    _, _ = item, err
//	}
func (r *Resource[T, C, U]) Iter(ctx context.Context, params Params) iter.Seq2[T, error] {
	return ListIter[T](ctx, r.client, ResolvePath(r.collection, params), nil)
}

// Get returns a single resource by appending id to the path.
//
// Usage:
//
//	item, err := res.Get(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
func (r *Resource[T, C, U]) Get(ctx context.Context, params Params) (*T, error) {
	var out T
	if err := r.client.Get(ctx, ResolvePath(r.path, params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a new resource.
//
// Usage:
//
//	item, err := res.Create(ctx, forge.Params{"owner": "core"}, body)
func (r *Resource[T, C, U]) Create(ctx context.Context, params Params, body *C) (*T, error) {
	var out T
	if err := r.client.Post(ctx, ResolvePath(r.collection, params), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update modifies an existing resource.
//
// Usage:
//
//	item, err := res.Update(ctx, forge.Params{"owner": "core", "repo": "go-forge"}, body)
func (r *Resource[T, C, U]) Update(ctx context.Context, params Params, body *U) (*T, error) {
	var out T
	if err := r.client.Patch(ctx, ResolvePath(r.path, params), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a resource.
//
// Usage:
//
//	err := res.Delete(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
func (r *Resource[T, C, U]) Delete(ctx context.Context, params Params) error {
	return r.client.Delete(ctx, ResolvePath(r.path, params))
}
