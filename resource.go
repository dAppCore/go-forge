package forge

import "context"

// Resource provides generic CRUD operations for a Forgejo API resource.
// T is the resource type, C is the create options type, U is the update options type.
type Resource[T any, C any, U any] struct {
	client *Client
	path   string
}

// NewResource creates a new Resource for the given path pattern.
// The path may contain {placeholders} that are resolved via Params.
func NewResource[T any, C any, U any](c *Client, path string) *Resource[T, C, U] {
	return &Resource[T, C, U]{client: c, path: path}
}

// List returns a single page of resources.
func (r *Resource[T, C, U]) List(ctx context.Context, params Params, opts ListOptions) (*PagedResult[T], error) {
	return ListPage[T](ctx, r.client, ResolvePath(r.path, params), nil, opts)
}

// ListAll returns all resources across all pages.
func (r *Resource[T, C, U]) ListAll(ctx context.Context, params Params) ([]T, error) {
	return ListAll[T](ctx, r.client, ResolvePath(r.path, params), nil)
}

// Get returns a single resource by appending id to the path.
func (r *Resource[T, C, U]) Get(ctx context.Context, params Params) (*T, error) {
	var out T
	if err := r.client.Get(ctx, ResolvePath(r.path, params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a new resource.
func (r *Resource[T, C, U]) Create(ctx context.Context, params Params, body *C) (*T, error) {
	var out T
	if err := r.client.Post(ctx, ResolvePath(r.path, params), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update modifies an existing resource.
func (r *Resource[T, C, U]) Update(ctx context.Context, params Params, body *U) (*T, error) {
	var out T
	if err := r.client.Patch(ctx, ResolvePath(r.path, params), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a resource.
func (r *Resource[T, C, U]) Delete(ctx context.Context, params Params) error {
	return r.client.Delete(ctx, ResolvePath(r.path, params))
}
