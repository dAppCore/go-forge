package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// MiscService handles miscellaneous Forgejo API endpoints such as
// markdown rendering, licence templates, gitignore templates, and
// server metadata.
// No Resource embedding — heterogeneous read-only endpoints.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Misc.GetVersion(ctx)
type MiscService struct {
	client *Client
}

func newMiscService(c *Client) *MiscService {
	return &MiscService{client: c}
}

// RenderMarkdown renders markdown text to HTML. The response is raw HTML
// text, not JSON.
func (s *MiscService) RenderMarkdown(ctx context.Context, text, mode string) (string, error) {
	body := types.MarkdownOption{Text: text, Mode: mode}
	data, err := s.client.PostRaw(ctx, "/api/v1/markdown", body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// RenderMarkup renders markup text to HTML. The response is raw HTML text,
// not JSON.
func (s *MiscService) RenderMarkup(ctx context.Context, text, mode string) (string, error) {
	body := types.MarkupOption{Text: text, Mode: mode}
	data, err := s.client.PostRaw(ctx, "/api/v1/markup", body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// RenderMarkdownRaw renders raw markdown text to HTML. The request body is
// sent as text/plain and the response is raw HTML text, not JSON.
func (s *MiscService) RenderMarkdownRaw(ctx context.Context, text string) (string, error) {
	data, err := s.client.postRawText(ctx, "/api/v1/markdown/raw", text)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ListLicenses returns all available licence templates.
func (s *MiscService) ListLicenses(ctx context.Context) ([]types.LicensesTemplateListEntry, error) {
	var out []types.LicensesTemplateListEntry
	if err := s.client.Get(ctx, "/api/v1/licenses", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IterLicenses returns an iterator over all available licence templates.
func (s *MiscService) IterLicenses(ctx context.Context) iter.Seq2[types.LicensesTemplateListEntry, error] {
	return func(yield func(types.LicensesTemplateListEntry, error) bool) {
		items, err := s.ListLicenses(ctx)
		if err != nil {
			yield(*new(types.LicensesTemplateListEntry), err)
			return
		}
		for _, item := range items {
			if !yield(item, nil) {
				return
			}
		}
	}
}

// GetLicense returns a single licence template by name.
func (s *MiscService) GetLicense(ctx context.Context, name string) (*types.LicenseTemplateInfo, error) {
	path := ResolvePath("/api/v1/licenses/{name}", pathParams("name", name))
	var out types.LicenseTemplateInfo
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListGitignoreTemplates returns all available gitignore template names.
func (s *MiscService) ListGitignoreTemplates(ctx context.Context) ([]string, error) {
	var out []string
	if err := s.client.Get(ctx, "/api/v1/gitignore/templates", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IterGitignoreTemplates returns an iterator over all available gitignore template names.
func (s *MiscService) IterGitignoreTemplates(ctx context.Context) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		items, err := s.ListGitignoreTemplates(ctx)
		if err != nil {
			yield("", err)
			return
		}
		for _, item := range items {
			if !yield(item, nil) {
				return
			}
		}
	}
}

// GetGitignoreTemplate returns a single gitignore template by name.
func (s *MiscService) GetGitignoreTemplate(ctx context.Context, name string) (*types.GitignoreTemplateInfo, error) {
	path := ResolvePath("/api/v1/gitignore/templates/{name}", pathParams("name", name))
	var out types.GitignoreTemplateInfo
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNodeInfo returns the NodeInfo metadata for the Forgejo instance.
func (s *MiscService) GetNodeInfo(ctx context.Context) (*types.NodeInfo, error) {
	var out types.NodeInfo
	if err := s.client.Get(ctx, "/api/v1/nodeinfo", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetVersion returns the server version.
func (s *MiscService) GetVersion(ctx context.Context) (*types.ServerVersion, error) {
	var out types.ServerVersion
	if err := s.client.Get(ctx, "/api/v1/version", &out); err != nil {
		return nil, err
	}
	return &out, nil
}
