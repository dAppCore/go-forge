package forge

// Forge is the top-level client for the Forgejo API.
type Forge struct {
	client *Client
}

// NewForge creates a new Forge client.
func NewForge(url, token string, opts ...Option) *Forge {
	c := NewClient(url, token, opts...)
	return &Forge{client: c}
}
