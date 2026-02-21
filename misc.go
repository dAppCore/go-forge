package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// MiscService handles miscellaneous Forgejo API endpoints such as
// markdown rendering, licence templates, gitignore templates, and
// server metadata.
// No Resource embedding — heterogeneous read-only endpoints.
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

// ListLicenses returns all available licence templates.
func (s *MiscService) ListLicenses(ctx context.Context) ([]types.LicensesTemplateListEntry, error) {
	var out []types.LicensesTemplateListEntry
	if err := s.client.Get(ctx, "/api/v1/licenses", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLicense returns a single licence template by name.
func (s *MiscService) GetLicense(ctx context.Context, name string) (*types.LicenseTemplateInfo, error) {
	path := fmt.Sprintf("/api/v1/licenses/%s", name)
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

// GetGitignoreTemplate returns a single gitignore template by name.
func (s *MiscService) GetGitignoreTemplate(ctx context.Context, name string) (*types.GitignoreTemplateInfo, error) {
	path := fmt.Sprintf("/api/v1/gitignore/templates/%s", name)
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
